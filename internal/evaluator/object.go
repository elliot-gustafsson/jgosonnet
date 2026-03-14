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
	Keys  []uint32
	Nodes ast.Nodes
	Meta  []uint8

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
		// Values: make([][]Value, len(layers)),
		// Scopes: make([]uint32, len(layers)),
	}
}

const (
	AssertStatusUnchecked uint8 = 0
	AssertStatusChecking  uint8 = 1
	AssertStatusChecked   uint8 = 2
)

type Object struct {
	Layers []*Layer

	Values       []Value
	layerOffsets []int

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

func (t *Object) getField(key uint32, ctx Context, offset int) (Value, bool, error) {

	err := runAssertions(t, ctx)
	if err != nil {
		return Value{}, false, err
	}

	var res Value

	currentVisibility := ast.ObjectFieldInherit

	for i := len(t.Layers) - (1 + offset); i >= 0; i-- {
		layer := t.Layers[i]

		fieldIndex := layer.findField(key)
		if fieldIndex == -1 {
			continue
		}

		evalCtx := ctx
		evalCtx.SuperOffset = len(t.Layers) - 1 - i

		val, err := getValue(t, i, fieldIndex, evalCtx)
		if err != nil {
			return Value{}, false, err
		}

		visibility, plusSuper := EvalFieldMeta(layer.Meta[fieldIndex])

		currentVisibility = getObjVisibility(currentVisibility, visibility)

		// Fast exit if its the first time we encounter the key and it shouldnt merge with super
		if res.IsNone() && !plusSuper {
			return val, false, nil
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
			// If field does not have plus, just return value since later layers dont matter
			// return CreateField(res, currentVisibility, plusSuper), true, nil
			return res, currentVisibility != ast.ObjectFieldHidden, nil
		}
	}

	if res.IsNone() {
		return Value{}, false, nil
	}

	return res, currentVisibility != ast.ObjectFieldHidden, nil
}

func (t *Object) getScope(layerIndex int, layer *Layer, ctx Context) (uint32, error) {

	if t.Scopes == nil {
		t.Scopes = make([]uint32, len(t.Layers))
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

	scopeId := ctx.Arena.NewScope(layer.ParentScopeId, len(layer.LocalKeys))

	// rootKeyId := ctx.Interner.Intern("$")

	for i, keyId := range layer.LocalKeys {
		node := layer.LocalNodes[i]

		// if keyId == rootKeyId {
		// 	// Only set root if not already set
		// 	_, found := ctx.Arena.GetScopeBind(layer.ParentScopeId, keyId)
		// 	if found {
		// 		continue
		// 	}
		// }

		val, err := evaluateNodeLazy(node, scopeId, ctx)
		if err != nil {
			return 0, err
		}

		ctx.Arena.AddScopeBind(scopeId, keyId, val)
	}

	return scopeId, nil
}

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

func (t *Object) AddLayers(layers []*Layer) {
	t.Layers = append(t.Layers, layers...)
}

func (t *Object) GetLength() int {
	res := make(map[uint32]any)
	for i := len(t.Layers) - 1; i >= 0; i-- {
		layer := t.Layers[i]

		for _, k := range layer.Keys {
			res[k] = nil
		}
	}
	return len(res)
}

func (t *Object) Clone() Object {
	layers := make([]*Layer, len(t.Layers))
	copy(layers, t.Layers)
	obj := NewObject(layers)
	return obj
}

func MergeObjects(left, right *Object) Object {
	layers := make([]*Layer, len(left.Layers)+len(right.Layers))
	copy(layers, left.Layers)
	copy(layers[len(left.Layers):], right.Layers)
	obj := NewObject(layers)
	return obj
}

func getObjVisibility(curr, inc ast.ObjectFieldHide) ast.ObjectFieldHide {
	if inc == ast.ObjectFieldVisible {
		return ast.ObjectFieldVisible
	}

	if curr == ast.ObjectFieldHidden {
		return ast.ObjectFieldHidden
	}

	return inc
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

	length := len(obj.Layers) * 5

	plans := make([]*FieldPlan, 0, length)

	keyToIndex := make(map[uint32]int, length)

	for l := len(obj.Layers) - 1; l >= 0; l-- {
		layer := obj.Layers[l]

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

	// Sort by field name
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

func (t *FieldPlan) GetValue(obj *Object, ctx Context) (Value, error) {
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

		err := EvaluateValueStrict(&value, ctx)
		if err != nil {
			return Value{}, err
		}

		overlayRef := t.Layers[i]

		innerVal, err := getValue(obj, overlayRef.LayerIdx, overlayRef.FieldIdx, ctx)
		if err != nil {
			return Value{}, err
		}

		err = EvaluateValueStrict(&innerVal, ctx)
		if err != nil {
			return Value{}, err
		}

		res, err := bopPlus(innerVal, value, ctx)
		if err != nil {
			return Value{}, err
		}
		value = res
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

	if obj.layerOffsets == nil {
		obj.layerOffsets = make([]int, len(obj.Layers))
	}

	if obj.Values == nil {
		totalFields := 0

		for i, layer := range obj.Layers {
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

	l := obj.Layers[layerId]

	n := l.Nodes[fieldId]

	scopeId, err := obj.getScope(layerId, l, ctx)
	if err != nil {
		return Value{}, err
	}

	val, err = EvaluateNodeStrict(n, scopeId, ctx)
	if err != nil {
		return Value{}, err
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

	for i := len(obj.Layers) - 1; i >= 0; i-- {
		layer := obj.Layers[i]

		scopeId, err := obj.getScope(i, layer, ctx)
		if err != nil {
			return err
		}

		for _, n := range layer.Asserts {
			val, err := EvaluateNodeStrict(n, scopeId, ctx)
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
		m := CreateFieldMeta(ast.ObjectFieldInherit, false)
		layer.Meta = []uint8{m, m}

		obj.Values = []Value{
			MakeString(ctx.Interner.Get(plan.KeyId), ctx),
			val,
		}

		kv := MakeObject(obj, ctx)

		res = append(res, kv)

	}

	return res, nil
}
