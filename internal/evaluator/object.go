package evaluator

import (
	"fmt"
	"slices"

	"github.com/google/go-jsonnet/ast"
)

const (
	MaxLinearKeys = 16

	MaskVisibility = 0x03 // Binary 00000011
	FlagPlusSuper  = 0x04 // Binary 00000100
)

type Field struct {
	Key        uint32
	Node       ast.Node
	Visibility ast.ObjectFieldHide
	PlusSuper  bool
}

type Layer struct {
	Keys   []uint32
	Nodes  ast.Nodes
	Values []Value
	Meta   []uint8

	Index map[uint32]int

	LocalKeys  []uint32
	LocalNodes ast.Nodes

	Asserts ast.Nodes

	ParentScopeId uint32
}

func (l *Layer) findField(key uint32) (layerId int) {
	idx := -1

	if l.Index != nil {
		if i, ok := l.Index[key]; ok {
			idx = i
		}
	} else {
		for i, k := range l.Keys {
			if k == key {
				idx = i
				break
			}
		}
	}

	return idx
}

func NewObject(layers []*Layer) Object {
	return Object{
		Layers: layers,
	}
}

const (
	AssertStatusUnchecked uint8 = 0
	AssertStatusChecking  uint8 = 1
	AssertStatusChecked   uint8 = 2
)

type Object struct {
	Layers []*Layer

	// Used ONLY for lazy merging (A + B)
	Left  *Object
	Right *Object

	Values       []Value
	layerOffsets []int

	fieldCache map[uint32]CachedValue

	Scopes []uint32

	AssertionState uint8
}

func (t *Field) Hidden() bool {
	return t.Visibility == ast.ObjectFieldHidden
}

func (t *Field) Visible() bool {
	return t.Visibility == ast.ObjectFieldVisible
}

func (t *Field) Inherit() bool {
	return t.Visibility == ast.ObjectFieldInherit
}

func (t *Object) GetField(key uint32, ctx Context) (Value, bool, error) {
	return t.getField(key, ctx, 0)
}

func (t *Object) GetSuperField(key uint32, ctx Context) (Value, bool, error) {
	return t.getField(key, ctx, 1)
}

func (t *Object) GetFieldWithOffset(key uint32, ctx Context, offset int) (Value, bool, error) {
	return t.getField(key, ctx, offset)
}

func (t *Object) getField(key uint32, ctx Context, offset int) (res Value, visible bool, err error) {

	err = runAssertions(t, ctx)
	if err != nil {
		return Value{}, false, err
	}

	if offset == 0 && t.fieldCache != nil {
		if cached, ok := t.fieldCache[key]; ok {
			return cached.Value, cached.Visible, nil
		}
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

		val, err := getValue(t, layerOffset, fieldIndex, ctx)
		if err != nil {
			return Value{}, false, err
		}

		visibility, plusSuper := EvalFieldMeta(layer.Meta[fieldIndex])

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
				return Value{}, false, err
			}
			res = v
		}

		if !plusSuper {
			// If field does not have plus, just beak since later layers dont matter
			break
		}
	}

	if res.IsNone() {
		return Value{}, false, nil
	}

	if currentVisibility == ast.ObjectFieldInherit {
		// If not explicitly visible, search for same key in lower layers to determine final visibility

		for j := layerOffset - 1; j >= 0; j-- {
			layer := layers[j]

			fieldIndex := layer.findField(key)
			if fieldIndex == -1 {
				continue
			}

			visibility, _ := EvalFieldMeta(layer.Meta[fieldIndex])
			if visibility != ast.ObjectFieldInherit {
				currentVisibility = visibility
				break
			}
		}
	}

	visible = currentVisibility != ast.ObjectFieldHidden

	if offset == 0 {
		if t.fieldCache == nil {
			t.fieldCache = make(map[uint32]CachedValue)
		}
		t.fieldCache[key] = CachedValue{Value: res, Visible: visible}
	}

	return res, visible, nil
}

func (t *Object) getScope(layerIndex int, layer *Layer, ctx Context) (uint32, error) {

	if t.Scopes == nil {
		t.Scopes = make([]uint32, len(t.GetLayers(ctx)))
	}

	scopeId := t.Scopes[layerIndex]
	if scopeId == 0 {
		sid, err := createScope(layer, ctx)
		if err != nil {
			return 0, err
		}
		t.Scopes[layerIndex] = sid
		scopeId = sid
	}

	return scopeId, nil
}

func createScope(layer *Layer, ctx Context) (uint32, error) {

	scopeId := ctx.NewScope(layer.ParentScopeId, len(layer.LocalKeys))

	for i, keyId := range layer.LocalKeys {
		node := layer.LocalNodes[i]

		val, err := evaluateNodeLazy(node, scopeId, ctx)
		if err != nil {
			return 0, err
		}

		ctx.AddScopeBind(scopeId, NamedValue{keyId, val})
	}

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

func EvalFieldMeta(m uint8) (visibility ast.ObjectFieldHide, plusSuper bool) {
	visibility = ast.ObjectFieldHide(m & MaskVisibility)
	plusSuper = (m & FlagPlusSuper) != 0
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

func (t *Object) totalLayerCount() int {
	if t.Layers != nil {
		return len(t.Layers)
	}
	return t.Left.totalLayerCount() + t.Right.totalLayerCount()
}

func (t *Object) populateLayers(dest []*Layer, offset int) int {
	if t.Layers != nil {
		copied := copy(dest[offset:], t.Layers)
		return offset + copied
	}
	offset = t.Left.populateLayers(dest, offset)
	return t.Right.populateLayers(dest, offset)
}

func (t *Object) GetLayers(ctx Context) []*Layer {
	if t.Layers != nil {
		return t.Layers
	}

	total := t.totalLayerCount()

	// layers := make([]*Layer, total)
	layers := ctx.Registry.LayerBufs.Alloc(total, total)

	t.populateLayers(layers, 0)

	t.Layers = layers
	t.Left = nil
	t.Right = nil

	return t.Layers
}

func MergeObjects(left, right *Object) Object {
	return Object{
		Left:  left,
		Right: right,
	}
}

type FieldPlan struct {
	KeyId      uint32
	Visibility ast.ObjectFieldHide
	IsClosed   bool
	Layers     []LayerRef
}

func (fp FieldPlan) IsHidden() bool {
	return fp.Visibility == ast.ObjectFieldHidden
}

type LayerRef struct {
	LayerIdx int
	FieldIdx int
}

func CompileObjectPlan(obj *Object, ctx Context) []*FieldPlan {
	return CompileObjectPlanEx(obj, ctx, false)
}

func CompileObjectPlanEx(obj *Object, ctx Context, naturalSort bool) []*FieldPlan {
	plans := compileObjectPlan(obj, ctx)

	if naturalSort {
		slices.SortFunc(plans, func(a, b *FieldPlan) int {
			aName := ctx.Interner.Get(a.KeyId)
			bName := ctx.Interner.Get(b.KeyId)
			return naturalStringSort(aName, bName)
		})
		return plans
	}

	slices.SortFunc(plans, func(a, b *FieldPlan) int {
		aName := ctx.Interner.Get(a.KeyId)
		bName := ctx.Interner.Get(b.KeyId)
		if aName > bName {
			return 1
		}
		if aName < bName {
			return -1
		}
		return 0
	})

	return plans
}

func compileObjectPlan(obj *Object, ctx Context) []*FieldPlan {
	layers := obj.GetLayers(ctx)

	length := len(layers) * 5

	plans := make([]*FieldPlan, 0, length)

	keyToIndex := make(map[uint32]int, length)

	for l := len(layers) - 1; l >= 0; l-- {
		layer := layers[l]

		for f, keyID := range layer.Keys {
			pIdx, exists := keyToIndex[keyID]
			var plan *FieldPlan

			if !exists {
				plan = &FieldPlan{
					KeyId:      keyID,
					Visibility: ast.ObjectFieldInherit,
					Layers:     make([]LayerRef, 0, 4),
				}
				plans = append(plans, plan)
				keyToIndex[keyID] = len(plans) - 1
			} else {
				plan = plans[pIdx]
			}

			vis, plus := EvalFieldMeta(layer.Meta[f])

			if plan.Visibility == ast.ObjectFieldInherit {
				switch vis {
				case ast.ObjectFieldHidden:
					plan.Visibility = ast.ObjectFieldHidden
				case ast.ObjectFieldVisible:
					plan.Visibility = ast.ObjectFieldVisible
				}
			}

			if plan.IsClosed {
				continue
			}

			plan.Layers = append(plan.Layers, LayerRef{l, f})

			if !plus {
				plan.IsClosed = true
			}
		}
	}

	return plans
}

func (t *FieldPlan) GetValue(obj *Object, ctx Context) (Value, error) {
	if obj.fieldCache != nil {
		if cached, ok := obj.fieldCache[t.KeyId]; ok {
			return cached.Value, nil
		}
	}

	layersCount := len(t.Layers)

	if layersCount == 0 {
		return Value{}, fmt.Errorf("no layers passed to plan.getValue")
	}

	// lastIdx := len(t.Layers) - 1

	layerRef := t.Layers[0]

	value, err := getValue(obj, layerRef.LayerIdx, layerRef.FieldIdx, ctx)
	if err != nil {
		return Value{}, err
	}

	for i := 1; i < layersCount; i++ {

		overlayRef := t.Layers[i]

		innerVal, err := getValue(obj, overlayRef.LayerIdx, overlayRef.FieldIdx, ctx)
		if err != nil {
			return Value{}, err
		}

		innerVal, err = innerVal.Eval(ctx)
		if err != nil {
			return Value{}, err
		}

		res, err := bopPlus(innerVal, value, ctx)
		if err != nil {
			return Value{}, err
		}

		res, err = res.Eval(ctx)
		if err != nil {
			return Value{}, err
		}
		value = res
	}

	if obj.fieldCache == nil {
		obj.fieldCache = make(map[uint32]CachedValue)
	}
	obj.fieldCache[t.KeyId] = CachedValue{
		Value:   value,
		Visible: t.Visibility != ast.ObjectFieldHidden,
	}

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

		if plan.Visibility == ast.ObjectFieldHidden {
			continue
		}

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		name := ctx.Interner.Get(keyId)

		rawVal, err := ManifestValue(value, ctx)
		if err != nil {
			return nil, err
		}

		res[name] = rawVal

	}

	return res, nil
}

func ManifestObjectRoot(obj *Object, ctx Context) (map[string]Value, error) {

	err := runAssertions(obj, ctx)
	if err != nil {
		return nil, err
	}

	plans := CompileObjectPlan(obj, ctx)

	res := make(map[string]Value, len(plans))
	for _, plan := range plans {
		// if len(values) == 0 {
		// 	continue
		// }
		keyId := plan.KeyId

		if plan.Visibility == ast.ObjectFieldHidden {
			continue
		}

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		name := ctx.Interner.Get(keyId)
		res[name] = value
	}

	return res, nil
}

func getValue(obj *Object, layerId, fieldId int, ctx Context) (Value, error) {

	layers := obj.GetLayers(ctx)

	if obj.layerOffsets == nil {
		obj.layerOffsets = make([]int, len(layers))
	}

	if obj.Values == nil {
		totalFields := 0

		for i, layer := range layers {
			obj.layerOffsets[i] = totalFields
			totalFields += len(layer.Keys)
		}
		obj.Values = make([]Value, totalFields)
	}

	flatIndex := obj.layerOffsets[layerId] + fieldId

	val := obj.Values[flatIndex]
	if !val.IsNone() {
		return val, nil
	}

	l := layers[layerId]

	if l.Values != nil {
		val = l.Values[fieldId]
	} else {

		n := l.Nodes[fieldId]

		evalCtx := ctx
		evalCtx.SuperOffset = len(layers) - 1 - layerId

		scopeId, err := obj.getScope(layerId, l, evalCtx)
		if err != nil {
			return Value{}, err
		}

		val, err = EvaluateNode(n, scopeId, evalCtx)
		if err != nil {
			return Value{}, err
		}
	}

	obj.Values[flatIndex] = val
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

		evalCtx := ctx
		evalCtx.SuperOffset = len(layers) - 1 - i

		scopeId, err := obj.getScope(i, layer, evalCtx)
		if err != nil {
			return err
		}

		for _, n := range layer.Asserts {
			val, err := EvaluateNode(n, scopeId, ctx)
			if err != nil {
				return err
			}
			if !val.IsBool() {
				return fmt.Errorf("unexpected assert return type: %s, expected bool", val.Type().String())
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
			res = append(res, MakeString(ctx.Interner.Get(fp.KeyId), ctx))
		}
	}

	return res
}

func GetObjectValues(obj *Object, ctx Context, inclHidden bool) ([]Value, error) {
	plans := CompileObjectPlan(obj, ctx)

	res := make([]Value, 0, len(plans))
	for _, plan := range plans {

		if !inclHidden && plan.Visibility == ast.ObjectFieldHidden {
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
	plans := CompileObjectPlan(obj, ctx)

	res := make([]Value, 0, len(plans))
	for _, plan := range plans {

		if !inclHidden && plan.Visibility == ast.ObjectFieldHidden {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		layer := &Layer{}
		obj := NewObject([]*Layer{layer})

		layer.Keys = []uint32{
			ctx.Interner.Intern("key"),
			ctx.Interner.Intern("value"),
		}
		layer.Values = []Value{
			MakeString(ctx.Interner.Get(plan.KeyId), ctx),
			val,
		}

		layer.Meta = []uint8{DefaultFieldMeta, DefaultFieldMeta}

		kv := MakeObject(obj, ctx)

		res = append(res, kv)

	}

	return res, nil
}

func GetObjectKeysValuesArray(obj *Object, ctx Context, inclHidden bool) ([]uint32, []Value, error) {
	plans := CompileObjectPlan(obj, ctx)

	keys := make([]uint32, 0, len(plans))
	vals := make([]Value, 0, len(plans))
	for _, plan := range plans {

		if !inclHidden && plan.Visibility == ast.ObjectFieldHidden {
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
		return Value{}, err
	}

	plans := CompileObjectPlan(t, ctx)

	layer := &Layer{
		Keys: make([]uint32, 0, len(plans)),
		Meta: make([]uint8, 0, len(plans)),
	}

	resObj := NewObject([]*Layer{layer})

	useMap := len(plans) > MaxLinearKeys
	if useMap {
		layer.Index = make(map[uint32]int, len(layer.Keys))
	}

	index := 0
	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(t, ctx)
		if err != nil {
			return Value{}, err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return Value{}, err

		}

		prunedVal, err := val.Prune(ctx)
		if err != nil {
			return Value{}, err
		}

		if prunedVal.IsEmpty(ctx) {
			continue
		}

		layer.Keys = append(layer.Keys, plan.KeyId)
		layer.Values = append(layer.Values, prunedVal)
		layer.Meta = append(layer.Meta, CreateFieldMeta(plan.Visibility, false))

		if useMap {
			layer.Index[plan.KeyId] = index
		}
		index++
	}

	return MakeObject(resObj, ctx), nil
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

		// TODO: Think over this, feels hacky...
		subCtx := ctx
		subCtx.Self = MakeObject(*a, ctx)

		valA, err := planAs[i].GetValue(a, subCtx)
		if err != nil {
			return false, err
		}
		valA, err = valA.Eval(ctx)
		if err != nil {
			return false, err
		}

		// TODO: Think over this, feels hacky...
		subCtx.Self = MakeObject(*b, ctx)
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
