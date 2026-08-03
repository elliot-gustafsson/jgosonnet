package evaluator

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strings"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"github.com/google/go-jsonnet/ast"
)

const (
	MaxPlanLinearKeys  = 32 // Threshold for compileObjectPlan map fallback
	MaxLayerLinearKeys = 32 // Threshold for Layer.Index map fallback

	MaskVisibility = 0x03 // Binary 00000011
	FlagPlusSuper  = 0x04 // Binary 00000100
	FlagTombstone  = 0x08 // Binary 00001000
)

type Layer struct {
	Keys   []uint32
	Nodes  ast.Nodes
	Values []Value
	Meta   []uint8

	Index *utils.DescriptorTable

	LocalKeys  []uint32
	LocalNodes ast.Nodes

	AssertsPtr uintptr

	ParentScopePtr uintptr
}

func (l *Layer) findField(key uint32) (layerId int) {

	if l.Index != nil {
		if i, ok := l.Index.Get(key); ok {
			return int(i)
		}
	} else {
		keys := l.Keys
		for i := range keys {
			if keys[i] == key {
				return i
			}
		}
	}

	return -1
}

func (l *Layer) packAsserts(s []ast.Node) {
	if len(s) > math.MaxUint16 {
		panic("asserts length exceeds uint16 max (65535)")
	}

	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(s)))

	const addrMask = (1 << 48) - 1
	l.AssertsPtr = (uintptr(len(s)) << 48) | (ptr & addrMask)
}

func (l *Layer) unpackAsserts() []ast.Node {
	if l.AssertsPtr == 0 {
		return nil
	}

	length := int(l.AssertsPtr >> 48)

	// mask of the top 16 bits
	const lenMask = 0xFFFF << 48

	ptr := unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&l.AssertsPtr))) &^ lenMask)

	// Note: The row above is equal to "ptr := unsafe.Pointer(l.AssertsPtr &^ lenMask)",
	// 	they result the same assembly code. Its done this way to not make "go vet"
	// 	flag it as "possible misuse of unsafe.Pointer".

	return unsafe.Slice((*ast.Node)(ptr), length)
}

func NewSingleLayerObject(allocator *arena.Allocator, layer *Layer) *Object {
	o := arena.Create[Object](allocator)
	o.Layers = arena.Alloc[*Layer](allocator, 1)
	o.Layers[0] = layer
	return o
}

const (
	AssertStatusUnchecked uint8 = 0
	AssertStatusChecking  uint8 = 1
	AssertStatusChecked   uint8 = 2
)

const FieldCacheInlineCount = 4

type FieldCache struct {
	inlineKeys    [FieldCacheInlineCount]uint32
	inlineVals    [FieldCacheInlineCount]Value
	fieldCache    *utils.PropertyMap[Value]
	inlineCount   uint8
	inlineVisible uint8
}

func (c *FieldCache) Get(key uint32) (v Value, visible bool, ok bool) {
	for i := range c.inlineCount {
		if c.inlineKeys[i] == key {
			// extract visibility from the bitmask
			visible := (c.inlineVisible & (1 << i)) != 0
			return c.inlineVals[i], visible, true
		}
	}
	if c.fieldCache != nil {
		if entry, meta, ok := c.fieldCache.GetEx(key); ok {
			return entry, (meta == 1), true
		}
	}
	return ValueNone, false, false
}

func (c *FieldCache) Set(key uint32, val Value, visible bool, ctx Context) {
	for i := range c.inlineCount {
		if c.inlineKeys[i] == key {
			c.inlineVals[i] = val
			if visible {
				c.inlineVisible |= (1 << i) // Set bit
			} else {
				c.inlineVisible &= ^(1 << i) // Clear bit
			}
			return
		}
	}
	if c.inlineCount < FieldCacheInlineCount {
		c.inlineKeys[c.inlineCount] = key
		c.inlineVals[c.inlineCount] = val
		if visible {
			c.inlineVisible |= (1 << c.inlineCount)
		}
		c.inlineCount++
		return
	}

	allocator := ctx.State.Registry.Allocator

	if c.fieldCache == nil {
		c.fieldCache = utils.NewPropertyMap[Value](allocator, 8)
	}
	var meta uint8
	if visible {
		meta = 1
	}
	c.fieldCache.PutEx(allocator, key, val, meta)
}

type Object struct {
	Layers []*Layer

	// Used ONLY for lazy merging (A + B)
	LeftPtr  uintptr
	RightPtr uintptr

	Cache FieldCache

	Scopes []uintptr

	AssertionState uint8
}

func (t *Object) GetField(key uint32, ctx Context) (val Value, vis bool, err error) {
	val, vis, ok := t.Cache.Get(key)
	if ok {
		return
	}
	val, vis, err = t.getField(key, ctx, 0)

	if err == nil {
		t.Cache.Set(key, val, vis, ctx)
	}
	return
}

func (t *Object) GetFieldWithOffset(key uint32, ctx Context, offset int) (Value, bool, error) {
	return t.getField(key, ctx, offset)
}

//go:noinline
func (t *Object) getField(key uint32, ctx Context, offset int) (res Value, visible bool, err error) {

	err = runAssertions(t, ctx)
	if err != nil {
		return ValueNone, false, err
	}

	currentVisibility := ast.ObjectFieldInherit

	layers := t.GetLayers(ctx)

	var layerOffset int

	for layerOffset = len(layers) - (1 + offset); layerOffset >= 0; layerOffset-- {
		layer := layers[layerOffset]

		fieldIndex := layer.findField(key)
		if fieldIndex == -1 {
			continue
		}

		visibility, plusSuper, tombstone := EvalFieldMeta(layer.Meta[fieldIndex])

		if tombstone {
			layerOffset -= int(layer.Values[fieldIndex].RefId())
			continue
		}

		val, err := getValue(t, layerOffset, fieldIndex, ctx)
		if err != nil {
			return ValueNone, false, err
		}

		if visibility != ast.ObjectFieldInherit {
			currentVisibility = visibility
		}

		if res.IsNone() && !plusSuper {
			// Fast exit if its the first time we encounter the key and it shouldnt merge with super
			res = val
			break
		}

		if res.IsNone() {
			res = val
		} else {
			v, err := bopPlus(val, res, ctx)
			if err != nil {
				return ValueNone, false, err
			}
			res = v
		}

		if !plusSuper {
			// If field does not have plus, just break since later layers dont matter
			break
		}
	}

	if res.IsNone() {
		return ValueNone, false, nil
	}

	if currentVisibility == ast.ObjectFieldInherit {
		// If not explicitly visible, search for same key in lower layers to determine final visibility

		for j := layerOffset - 1; j >= 0; j-- {
			layer := layers[j]

			fieldIndex := layer.findField(key)
			if fieldIndex == -1 {
				continue
			}

			visibility, _, tombstone := EvalFieldMeta(layer.Meta[fieldIndex])

			if tombstone {
				j -= int(layer.Values[fieldIndex].RefId())
				continue
			}

			if visibility != ast.ObjectFieldInherit {
				currentVisibility = visibility
				break
			}
		}
	}

	visible = currentVisibility != ast.ObjectFieldHidden

	return res, visible, nil
}

func (t *Object) getScope(layerIndex int, layer *Layer, ctx Context) (uintptr, error) {

	// no locals no need to create another scope
	if len(layer.LocalKeys) == 0 {
		return layer.ParentScopePtr, nil
	}

	return t.createLayerScope(layerIndex, layer, ctx)
}

//go:noinline
func (t *Object) createLayerScope(layerIndex int, layer *Layer, ctx Context) (uintptr, error) {
	if t.Scopes == nil {
		t.Scopes = arena.Alloc[uintptr](ctx.State.Registry.Allocator, len(t.GetLayers(ctx)))
	}

	scopePtr := t.Scopes[layerIndex]
	if scopePtr != 0 {
		return scopePtr, nil
	}

	s, scopeId := ctx.NewScope(layer.ParentScopePtr, len(layer.LocalKeys))

	for i := range layer.LocalKeys {
		node := layer.LocalNodes[i]

		val, err := evaluateNodeLazy(node, scopeId, ctx)
		if err != nil {
			return 0, err
		}

		s.Bindings[i] = NamedValue{layer.LocalKeys[i], val}
	}

	t.Scopes[layerIndex] = scopeId

	return scopeId, nil
}

const DefaultFieldMeta = uint8(ast.ObjectFieldInherit) & MaskVisibility

func CreateFieldMeta(visibility ast.ObjectFieldHide, plusSuper bool) uint8 {
	m := uint8(visibility) & MaskVisibility
	if plusSuper {
		m |= FlagPlusSuper
	}
	return m
}

func EvalFieldMeta(m uint8) (visibility ast.ObjectFieldHide, plusSuper, tombstone bool) {
	visibility = ast.ObjectFieldHide(m & MaskVisibility)
	plusSuper = (m & FlagPlusSuper) != 0
	tombstone = (m & FlagTombstone) != 0
	return
}

func (t *Object) Length(ctx Context) int {
	fps := CompileObjectPlan(t, ctx)

	length := 0
	for _, v := range fps {
		if v.IsHidden() {
			continue
		}
		length++
	}
	return length
}

func (t *Object) appendLayers(dest []*Layer, ctx Context) []*Layer {
	if t.Layers != nil {
		targetLen := len(dest) + len(t.Layers)
		if targetLen > cap(dest) {
			newCap := max(targetLen, cap(dest)*2)
			dest = arena.Realloc(ctx.State.Registry.Allocator, dest, newCap)[:len(dest)]
		}

		return append(dest, t.Layers...)
	}

	if t.LeftPtr == 0 || t.RightPtr == 0 {
		return dest
	}

	// copy lefts layer ids
	l := (*Object)(resolveUintptr(t.LeftPtr))
	dest = l.appendLayers(dest, ctx)

	// copy rights layer ids
	r := (*Object)(resolveUintptr(t.RightPtr))
	return r.appendLayers(dest, ctx)
}

func (t *Object) GetLayers(ctx Context) []*Layer {
	if t.Layers != nil {
		return t.Layers
	}

	layers := arena.Alloc[*Layer](ctx.State.Registry.Allocator, 8)[:0]
	layers = t.appendLayers(layers, ctx)

	t.Layers = layers
	t.LeftPtr = 0
	t.RightPtr = 0

	return t.Layers
}

func MergeObjects(leftPtr, rightPtr uintptr, ctx Context) *Object {
	o := arena.Create[Object](ctx.State.Registry.Allocator)
	o.LeftPtr = leftPtr
	o.RightPtr = rightPtr
	return o
}

type FieldPlan struct {
	Layers           []LayerRef
	KeyId            uint32
	ShadowUntilLayer uint16
	Visibility       uint8
	IsClosed         bool
}

func (fp FieldPlan) IsHidden() bool {
	return fp.Visibility == uint8(ast.ObjectFieldHidden)
}

func (fp FieldPlan) IsInherit() bool {
	return fp.Visibility == uint8(ast.ObjectFieldInherit)
}

func (fp FieldPlan) IsVisible() bool {
	return fp.Visibility == uint8(ast.ObjectFieldVisible)
}

type LayerRef struct {
	LayerIdx int32
	FieldIdx int32
}

func CompileObjectPlan(obj *Object, ctx Context) []FieldPlan {
	return CompileObjectPlanEx(obj, ctx, false)
}

func CompileObjectPlanEx(obj *Object, ctx Context, naturalSort bool) []FieldPlan {
	plans := compileObjectPlan(obj, ctx)

	interner := ctx.State.Interner

	if naturalSort {
		slices.SortFunc(plans, func(a, b FieldPlan) int {
			aName := interner.Get(a.KeyId)
			bName := interner.Get(b.KeyId)
			return naturalStringSort(aName, bName)
		})
		return plans
	}

	slices.SortFunc(plans, func(a, b FieldPlan) int {
		aName := interner.Get(a.KeyId)
		bName := interner.Get(b.KeyId)
		return cmp.Compare(aName, bName)

	})

	return plans
}

func compileObjectPlan(obj *Object, ctx Context) []FieldPlan {
	allocator := ctx.State.Registry.Allocator

	layers := obj.GetLayers(ctx)

	maxKeys := 0
	for i := range layers {
		maxKeys += len(layers[i].Keys)
	}

	plans := arena.Alloc[FieldPlan](allocator, maxKeys)[:0]

	var planIdxMap *utils.DescriptorTable
	useMap := maxKeys > MaxPlanLinearKeys
	if useMap {
		planIdxMap = utils.NewEmptyDescriptorTable(allocator, maxKeys)
	}

	var validate bool

	for l := len(layers) - 1; l >= 0; l-- {
		layer := layers[l]

		for f, keyID := range layer.Keys {

			pIdx := -1
			if useMap {
				if idx, exists := planIdxMap.Get(keyID); exists {
					pIdx = int(idx)
				}
			} else {
				for i := 0; i < len(plans); i++ {
					if plans[i].KeyId == keyID {
						pIdx = i
						break
					}
				}
			}

			if pIdx == -1 {
				plans = append(plans, FieldPlan{
					KeyId:            keyID,
					Visibility:       uint8(ast.ObjectFieldInherit),
					Layers:           arena.Alloc[LayerRef](allocator, 4)[:0],
					ShadowUntilLayer: math.MaxUint16,
				})
				pIdx = len(plans) - 1

				if useMap {
					planIdxMap.Append(keyID)
				}

			}

			vis, plus, tombstone := EvalFieldMeta(layer.Meta[f])

			plan := &plans[pIdx]

			if plan.IsInherit() {
				switch vis {
				case ast.ObjectFieldHidden:
					plan.Visibility = uint8(ast.ObjectFieldHidden)
				case ast.ObjectFieldVisible:
					plan.Visibility = uint8(ast.ObjectFieldVisible)
				}
			}

			if plan.IsClosed {
				continue
			}

			if plan.ShadowUntilLayer != math.MaxUint16 && l >= int(plan.ShadowUntilLayer) {
				continue
			}

			if tombstone {
				validate = true
				plan.ShadowUntilLayer = uint16(l - int(layer.Values[f].RefId()))
				continue
			}

			if len(plan.Layers) == cap(plan.Layers) {
				n := len(plan.Layers)
				plan.Layers = arena.Realloc(allocator, plan.Layers, n*2)[:n]
			}
			plan.Layers = append(plan.Layers, LayerRef{int32(l), int32(f)})

			if !plus {
				plan.IsClosed = true
			}
		}
	}

	if !validate {
		return plans
	}

	validPlans := plans[:0]
	for _, p := range plans {
		if len(p.Layers) > 0 {
			validPlans = append(validPlans, p)
		}
	}

	return validPlans
}

func (t *FieldPlan) GetValue(obj *Object, ctx Context) (Value, error) {
	value, _, ok := obj.Cache.Get(t.KeyId)
	if ok {
		return value, nil
	}

	err := runAssertions(obj, ctx)
	if err != nil {
		return ValueNone, err
	}

	layersCount := len(t.Layers)

	if layersCount == 0 {
		return ValueNone, fmt.Errorf("no layers passed to plan.getValue")
	}

	// lastIdx := len(t.Layers) - 1

	layerRef := t.Layers[0]

	value, err = getValue(obj, int(layerRef.LayerIdx), int(layerRef.FieldIdx), ctx)
	if err != nil {
		return ValueNone, err
	}

	for i := 1; i < layersCount; i++ {

		overlayRef := t.Layers[i]

		innerVal, err := getValue(obj, int(overlayRef.LayerIdx), int(overlayRef.FieldIdx), ctx)
		if err != nil {
			return ValueNone, err
		}

		innerVal, err = innerVal.Eval(ctx)
		if err != nil {
			return ValueNone, err
		}

		res, err := bopPlus(innerVal, value, ctx)
		if err != nil {
			return ValueNone, err
		}

		res, err = res.Eval(ctx)
		if err != nil {
			return ValueNone, err
		}
		value = res
	}

	obj.Cache.Set(t.KeyId, value, !t.IsHidden(), ctx)

	return value, nil
}

func manifestObject(obj *Object, ctx Context) (map[string]any, error) {

	err := runAssertions(obj, ctx)
	if err != nil {
		return nil, err
	}

	plans := CompileObjectPlan(obj, ctx)

	res := make(map[string]any, len(plans))
	for _, plan := range plans {
		// if len(values) == 0 {
		// 	continue
		// }
		keyId := plan.KeyId

		if plan.IsHidden() {
			continue
		}

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		name := ctx.State.Interner.Get(keyId)

		rawVal, err := ManifestValue(value, ctx)
		if err != nil {
			return nil, err
		}

		// Clone due to arenas being reset
		nameClone := strings.Clone(name)
		res[nameClone] = rawVal

	}

	return res, nil
}

func ManifestObjectRoot(obj *Object, ctx Context) ([]NamedValue, error) {

	err := runAssertions(obj, ctx)
	if err != nil {
		return nil, err
	}

	plans := CompileObjectPlan(obj, ctx)

	res := arena.Alloc[NamedValue](ctx.State.Registry.Allocator, len(plans))
	var index int
	for _, plan := range plans {

		if plan.IsHidden() {
			continue
		}

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		res[index] = NamedValue{plan.KeyId, value}
		index++
	}

	if index < len(plans) {
		res = res[:index]
	}

	return res, nil
}

func getValue(obj *Object, layerId, fieldId int, ctx Context) (Value, error) {

	layers := obj.GetLayers(ctx)

	var val Value

	l := layers[layerId]

	if l.Values != nil {
		val = l.Values[fieldId]
	} else {

		n := l.Nodes[fieldId]

		evalCtx := ctx
		evalCtx.SuperOffset = uint32(len(layers) - 1 - layerId)

		scopeId, err := obj.getScope(layerId, l, evalCtx)
		if err != nil {
			return ValueNone, err
		}

		val, err = EvaluateNode(n, scopeId, evalCtx)
		if err != nil {
			return ValueNone, err
		}
	}

	return val, nil
}

func runAssertions(obj *Object, ctx Context) error {
	if obj.AssertionState == AssertStatusChecked {
		return nil
	}
	if obj.AssertionState == AssertStatusChecking {
		return nil
	}

	obj.AssertionState = AssertStatusChecking

	layers := obj.GetLayers(ctx)

	for i := len(layers) - 1; i >= 0; i-- {
		layer := layers[i]

		asserts := layer.unpackAsserts()
		if asserts == nil {
			continue
		}

		evalCtx := ctx
		evalCtx.SuperOffset = uint32(len(layers) - 1 - i)

		scopeId, err := obj.getScope(i, layer, evalCtx)
		if err != nil {
			return err
		}

		for _, n := range asserts {
			val, err := EvaluateNode(n, scopeId, ctx)
			if err != nil {
				return err
			}
			if !val.IsBool() {
				return TypeErrorSpecific(ValueTypeBool, val.Type())
			}
		}
	}

	obj.AssertionState = AssertStatusChecked
	return nil
}

func GetObjectFields(obj *Object, ctx Context, inclhidden bool) []Value {

	plans := CompileObjectPlan(obj, ctx)

	res := make([]Value, 0, len(plans))
	for _, fp := range plans {
		if inclhidden || !fp.IsHidden() {
			res = append(res, MakeString(ctx.State.Interner.Get(fp.KeyId), ctx))
		}
	}

	return res
}

func GetObjectValues(obj *Object, ctx Context, inclHidden bool) ([]Value, error) {
	plans := CompileObjectPlan(obj, ctx)

	res := make([]Value, 0, len(plans))
	for _, plan := range plans {

		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		res = append(res, val)

	}

	return res, nil
}

func GetObjectKeysValues(obj *Object, ctx Context, inclHidden bool) ([]Value, error) {
	allocator := ctx.State.Registry.Allocator

	plans := CompileObjectPlan(obj, ctx)

	keyIdx := ctx.State.Interner.Intern("key")
	valueIdx := ctx.State.Interner.Intern("value")

	res := arena.Alloc[Value](allocator, len(plans))

	var index int
	for _, plan := range plans {

		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		layer := arena.Create[Layer](allocator)

		layer.Keys = arena.Alloc[uint32](allocator, 2)
		layer.Keys[0] = keyIdx
		layer.Keys[1] = valueIdx

		layer.Values = arena.Alloc[Value](allocator, 2)
		layer.Values[0] = MakeString(ctx.State.Interner.Get(plan.KeyId), ctx)
		layer.Values[1] = val

		layer.Meta = arena.Alloc[uint8](allocator, 2)
		layer.Meta[0] = DefaultFieldMeta
		layer.Meta[1] = DefaultFieldMeta

		obj := NewSingleLayerObject(allocator, layer)
		kv := MakeObjectValue(obj)

		res[index] = kv

		index++

	}

	if index < len(plans) {
		res = res[:index]
	}

	return res, nil
}

func GetObjectKeysValuesArray(obj *Object, ctx Context, inclHidden bool) ([]uint32, []Value, error) {
	plans := CompileObjectPlan(obj, ctx)

	keys := make([]uint32, 0, len(plans))
	vals := make([]Value, 0, len(plans))
	for _, plan := range plans {

		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, nil, err
		}

		keys = append(keys, plan.KeyId)
		vals = append(vals, val)
	}

	return keys, vals, nil
}

func (t *Object) Prune(ctx Context) (Value, error) {
	err := runAssertions(t, ctx)
	if err != nil {
		return ValueNone, err
	}

	plans := CompileObjectPlan(t, ctx)

	n := len(plans)
	allocator := ctx.State.Registry.Allocator

	layer := arena.Create[Layer](allocator)
	layer.Keys = arena.Alloc[uint32](allocator, n)
	layer.Values = arena.Alloc[Value](allocator, n)
	layer.Meta = arena.Alloc[uint8](allocator, n)

	useMap := n > MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, n)
	}

	index := 0
	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(t, ctx)
		if err != nil {
			return ValueNone, err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return ValueNone, err

		}

		prunedVal, err := val.Prune(ctx)
		if err != nil {
			return ValueNone, err
		}

		if prunedVal.IsEmpty(ctx) {
			continue
		}

		layer.Keys[index] = plan.KeyId
		layer.Values[index] = prunedVal
		layer.Meta[index] = CreateFieldMeta(ast.ObjectFieldHide(plan.Visibility), false)

		if useMap {
			layer.Index.Append(plan.KeyId)
		}
		index++
	}

	if index < n {
		layer.Keys = layer.Keys[:index]
		layer.Values = layer.Values[:index]
		layer.Meta = layer.Meta[:index]
	}

	obj := NewSingleLayerObject(allocator, layer)

	return MakeObjectValue(obj), nil
}

func (a *Object) Equal(b *Object, ctx Context) (bool, error) {
	if a == b {
		return true, nil
	}

	planAs := CompileObjectPlan(a, ctx)
	planBs := CompileObjectPlan(b, ctx)
	if len(planAs) == 0 && len(planBs) == 0 {
		return true, nil
	}

	i, j := 0, 0
	for i < len(planAs) && j < len(planBs) {
		// Skip hidden fields for both objects
		if planAs[i].IsHidden() {
			i++
			continue
		}
		if planBs[j].IsHidden() {
			j++
			continue
		}

		if planAs[i].KeyId != planBs[j].KeyId {
			return false, nil
		}

		subCtx := ctx
		subCtx.Self = MakeObjectValue(a)

		valA, err := planAs[i].GetValue(a, subCtx)
		if err != nil {
			return false, err
		}
		valA, err = valA.Eval(ctx)
		if err != nil {
			return false, err
		}

		subCtx.Self = MakeObjectValue(b)
		valB, err := planBs[j].GetValue(b, subCtx)
		if err != nil {
			return false, err
		}
		valB, err = valB.Eval(ctx)
		if err != nil {
			return false, err
		}

		eq, err := valA.Equal(valB, ctx)
		if err != nil {
			return false, err
		}
		if !eq {
			return false, nil
		}
		i++
		j++
	}

	// Make sure object A has no remaining visible fields
	for i < len(planAs) {
		if !planAs[i].IsHidden() {
			return false, nil
		}
		i++
	}
	// Make sure object B has no remaining visible fields
	for j < len(planBs) {
		if !planBs[j].IsHidden() {
			return false, nil
		}
		j++
	}
	return true, nil
}
