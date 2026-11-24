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
	MaxLayerLinearKeys = 32 // Threshold for eager DynamicLayer.Index map fallback
	LazyLayerIndexMin  = 8  // Threshold for lazy DynamicLayer.Index on first findField

	MaskVisibility = 0x03 // Binary 00000011
	FlagPlusSuper  = 0x04 // Binary 00000100
	FlagTombstone  = 0x08 // Binary 00001000
)

// DynamicLayer holds per-instance field tables for non-static layers
// (computed keys, objectFlatMerge, std.parseJson/Yaml, std.prune, std.mapWithKey,
// std.objectRemoveKey, std.mergePatch, stdlib, file thisFile layer).
type DynamicLayer struct {
	Keys       []uint32
	Nodes      []*IRNode
	Values     []Value
	Meta       []uint8
	Index      *utils.DescriptorTable
	LocalNodes []*IRNode
	Asserts    []*IRNode
	HasHidden  bool
}

// Layer is a compact 24-byte object layer. Static object literals set Template
// (shared immutable compile-time data) and leave Dyn nil.
type Layer struct {
	Template       *LayerTemplate
	Dyn            *DynamicLayer
	ParentScopePtr uintptr
}

func NewDynamicLayer(allocator *arena.Allocator) (*Layer, *DynamicLayer) {
	dl := arena.CreateZeroed[DynamicLayer](allocator)
	l := arena.CreateZeroed[Layer](allocator)
	l.Dyn = dl
	return l, dl
}

func (l *Layer) Keys() []uint32 {
	if l.Template != nil {
		return l.Template.Keys
	}
	return l.Dyn.Keys
}

func (l *Layer) Meta() []uint8 {
	if l.Template != nil {
		return l.Template.Meta
	}
	return l.Dyn.Meta
}

func (l *Layer) Nodes() []*IRNode {
	if l.Template != nil {
		return l.Template.Nodes
	}
	return l.Dyn.Nodes
}

func (l *Layer) LocalNodes() []*IRNode {
	if l.Template != nil {
		return l.Template.LocalNodes
	}
	return l.Dyn.LocalNodes
}

func (l *Layer) Asserts() []*IRNode {
	if l.Template != nil {
		return l.Template.Asserts
	}
	return l.Dyn.Asserts
}

func (l *Layer) HasHidden() bool {
	if l.Template != nil {
		return l.Template.HasHidden
	}
	return l.Dyn.HasHidden
}

func (l *Layer) findField(key uint32, allocator *arena.Allocator) int {
	if t := l.Template; t != nil {
		if t.Index != nil {
			return t.Index.GetIndex(key)
		}
		for i, k := range t.Keys {
			if k == key {
				return i
			}
		}
		return -1
	}
	d := l.Dyn
	if d.Index != nil {
		return d.Index.GetIndex(key)
	}
	if len(d.Keys) > LazyLayerIndexMin && allocator != nil {
		d.Index = utils.NewDescriptorTable(allocator, d.Keys)
		return d.Index.GetIndex(key)
	}
	for i, k := range d.Keys {
		if k == key {
			return i
		}
	}
	return -1
}

func (l *Layer) UnpackFieldScopes() []uintptr {
	if l.ParentScopePtr&1 == 0 {
		return nil
	}
	ptr := (*uintptr)(unsafe.Pointer(l.ParentScopePtr &^ 1))
	return unsafe.Slice(ptr, len(l.Keys()))
}

type singleLayerObject struct {
	obj   Object
	layer [1]uintptr
}

type staticSingleLayerObject struct {
	obj      Object
	layerPtr [1]uintptr
	layer    Layer
}

// NewStaticSingleLayerObject allocates a 64-byte Object + embedded 24-byte Layer
// in a single 96-byte arena allocation for static object literals.
func NewStaticSingleLayerObject(allocator *arena.Allocator, tmpl *LayerTemplate, scopePtr uintptr) *Object {
	s := arena.CreateZeroed[staticSingleLayerObject](allocator)
	s.layer.Template = tmpl
	s.layer.ParentScopePtr = scopePtr
	s.layerPtr[0] = uintptr(unsafe.Pointer(&s.layer))
	s.obj.Layers = s.layerPtr[:]
	s.obj.LayerCount = 1
	if len(tmpl.Asserts) == 0 {
		s.obj.AssertionState = AssertStatusChecked
	}
	if tmpl.HasHidden {
		s.obj.Flags = objFlagHasHidden
	}
	return &s.obj
}

func NewSingleLayerObject(allocator *arena.Allocator, layer *Layer) *Object {
	s := arena.CreateZeroed[singleLayerObject](allocator)
	s.layer[0] = uintptr(unsafe.Pointer(layer))
	s.obj.Layers = s.layer[:]
	s.obj.LayerCount = 1
	if len(layer.Asserts()) == 0 {
		s.obj.AssertionState = AssertStatusChecked
	}
	if layer.HasHidden() {
		s.obj.Flags = objFlagHasHidden
	}
	return &s.obj
}

func NewMultiLayerObject(allocator *arena.Allocator, existingLayers []uintptr, topLayer *Layer) *Object {
	newLen := len(existingLayers) + 1
	newLayers := arena.Alloc[uintptr](allocator, newLen)
	copy(newLayers, existingLayers)
	newLayers[newLen-1] = uintptr(unsafe.Pointer(topLayer))

	resObj := arena.CreateZeroed[Object](allocator)
	resObj.Layers = newLayers
	resObj.LayerCount = uint32(newLen)

	hasHidden := topLayer.HasHidden()
	if !hasHidden {
		for _, lp := range existingLayers {
			if (*Layer)(resolveUintptr(lp)).HasHidden() {
				hasHidden = true
				break
			}
		}
	}
	if hasHidden {
		resObj.Flags = objFlagHasHidden
	}
	return resObj
}

func LayerAt(layers []uintptr, idx int) *Layer {
	return (*Layer)(resolveUintptr(layers[idx]))
}

func layerAt(layers []uintptr, idx int) *Layer {
	return (*Layer)(resolveUintptr(layers[idx]))
}

const (
	AssertStatusUnchecked uint8 = 0
	AssertStatusChecking  uint8 = 1
	AssertStatusChecked   uint8 = 2

	objFlagHasHidden uint8 = 1 << 0
)

const (
	FieldCacheInlineCount = 4
	FieldCacheUnevaled    = Value(1)

	cacheVisUnknown uint8 = 0
	cacheVisHidden  uint8 = 1
	cacheVisVisible uint8 = 2
)

type FieldCache struct {
	inlineKeys  [FieldCacheInlineCount]uint32
	inlineVals  [FieldCacheInlineCount]Value
	fieldCache  *utils.PropertyMap[Value]
	inlineCount uint8
	inlineVis   uint8 // 2 bits per inline slot: cacheVisUnknown/Hidden/Visible
}

func (c *FieldCache) getInline(key uint32) (Value, uint8, bool) {
	for i := range c.inlineCount {
		if c.inlineKeys[i] == key {
			return c.inlineVals[i], (c.inlineVis >> (i * 2)) & 3, true
		}
	}
	return ValueNone, cacheVisUnknown, false
}

//go:noinline
func (c *FieldCache) getRawSlow(key uint32) (Value, uint8, bool) {
	if entry, meta, ok := c.fieldCache.GetEx(key); ok {
		return entry, meta, true
	}
	return ValueNone, cacheVisUnknown, false
}

func (c *FieldCache) getRaw(key uint32) (Value, uint8, bool) {
	if v, vis, ok := c.getInline(key); ok {
		return v, vis, true
	}
	if c.fieldCache != nil {
		return c.getRawSlow(key)
	}
	return ValueNone, cacheVisUnknown, false
}

func (c *FieldCache) setRaw(key uint32, val Value, vis uint8, ctx Context) {
	for i := range c.inlineCount {
		if c.inlineKeys[i] == key {
			if val != FieldCacheUnevaled || c.inlineVals[i] == ValueNone {
				c.inlineVals[i] = val
			}
			if vis != cacheVisUnknown {
				shift := i * 2
				c.inlineVis = (c.inlineVis &^ (3 << shift)) | (vis << shift)
			}
			return
		}
	}
	if c.inlineCount < FieldCacheInlineCount {
		idx := c.inlineCount
		c.inlineKeys[idx] = key
		c.inlineVals[idx] = val
		shift := idx * 2
		c.inlineVis = (c.inlineVis &^ (3 << shift)) | (vis << shift)
		c.inlineCount++
		return
	}

	allocator := ctx.State.Allocator
	if c.fieldCache == nil {
		c.fieldCache = utils.NewPropertyMap[Value](allocator, 8)
	}
	if val == FieldCacheUnevaled {
		if existing, _, ok := c.fieldCache.GetEx(key); ok && existing != ValueNone {
			val = existing
		}
	}
	c.fieldCache.PutEx(allocator, key, val, vis)
}

func (t *Object) ensureCache(ctx Context) *FieldCache {
	if t.Cache == nil {
		t.Cache = arena.CreateZeroed[FieldCache](ctx.State.Allocator)
	}
	return t.Cache
}

type Object struct {
	Layers []uintptr

	// Used ONLY for lazy merging (A + B)
	LeftPtr  uintptr
	RightPtr uintptr

	Cache     *FieldCache
	ScopesPtr uintptr

	LayerCount     uint32
	AssertionState uint8
	Flags          uint8
}

func (t *Object) cacheSet(key uint32, val Value, vis uint8, ctx Context) {
	t.ensureCache(ctx).setRaw(key, val, vis, ctx)
}

// GetField retrieves a field value without scanning lower layers for visibility.
func (t *Object) GetField(key uint32, ctx Context) (val Value, err error) {
	if c := t.Cache; c != nil {
		if v, _, found := c.getInline(key); found {
			if v != FieldCacheUnevaled {
				return v, nil
			}
		} else if c.fieldCache != nil {
			if v, _, found := c.getRawSlow(key); found && v != FieldCacheUnevaled {
				return v, nil
			}
		}
	}
	var visState uint8
	val, visState, err = t.getField(key, ctx, 0, false)
	if err == nil {
		t.cacheSet(key, val, visState, ctx)
	}
	return val, err
}

// GetFieldVis retrieves both the field value and its final visibility.
func (t *Object) GetFieldVis(key uint32, ctx Context) (val Value, visible bool, err error) {
	if t.Cache != nil {
		if v, visState, found := t.Cache.getRaw(key); found && v != FieldCacheUnevaled {
			if v.IsNone() {
				return ValueNone, false, nil
			}
			if visState != cacheVisUnknown {
				return v, visState == cacheVisVisible, nil
			}
			vis := t.resolveFieldVisibility(key, ctx)
			t.Cache.setRaw(key, v, vis, ctx)
			return v, vis == cacheVisVisible, nil
		}
	}
	var visState uint8
	val, visState, err = t.getField(key, ctx, 0, true)
	if err == nil {
		t.cacheSet(key, val, visState, ctx)
	}
	return val, visState == cacheVisVisible, err
}

func (t *Object) GetFieldWithOffset(key uint32, ctx Context, offset int) (Value, bool, error) {
	val, visState, err := t.getField(key, ctx, offset, false)
	return val, visState == cacheVisVisible, err
}

func (t *Object) resolveFieldVisibility(key uint32, ctx Context) uint8 {
	if (t.Flags & objFlagHasHidden) == 0 {
		return cacheVisVisible
	}
	layers := t.GetLayers(ctx)
	alloc := ctx.State.Allocator
	for layerOffset := len(layers) - 1; layerOffset >= 0; layerOffset-- {
		layer := layerAt(layers, layerOffset)
		fieldIndex := layer.findField(key, alloc)
		if fieldIndex == -1 {
			continue
		}
		vis, _, tombstone := EvalFieldMeta(layer.Meta()[fieldIndex])
		if tombstone {
			layerOffset -= int(uint32(layer.Dyn.Values[fieldIndex]))
			continue
		}
		if vis == ast.ObjectFieldHidden {
			return cacheVisHidden
		}
		if vis == ast.ObjectFieldVisible {
			return cacheVisVisible
		}
	}
	return cacheVisVisible
}

//go:noinline
func (t *Object) getField(key uint32, ctx Context, offset int, needVis bool) (res Value, visState uint8, err error) {

	err = runAssertions(t, ctx)
	if err != nil {
		return ValueNone, cacheVisUnknown, err
	}

	currentVisibility := ast.ObjectFieldInherit

	layers := t.GetLayers(ctx)
	alloc := ctx.State.Allocator

	var layerOffset int

	for layerOffset = len(layers) - (1 + offset); layerOffset >= 0; layerOffset-- {
		layer := layerAt(layers, layerOffset)

		fieldIndex := layer.findField(key, alloc)
		if fieldIndex == -1 {
			continue
		}

		visibility, plusSuper, tombstone := EvalFieldMeta(layer.Meta()[fieldIndex])

		if tombstone {
			layerOffset -= int(uint32(layer.Dyn.Values[fieldIndex]))
			continue
		}

		val, err := getValue(t, layerOffset, fieldIndex, ctx)
		if err != nil {
			return ValueNone, cacheVisUnknown, err
		}

		if visibility != ast.ObjectFieldInherit {
			currentVisibility = visibility
		}

		if res.IsNone() && !plusSuper {
			res = val
			break
		}

		if res.IsNone() {
			res = val
		} else {
			v, err := bopPlus(val, res, ctx)
			if err != nil {
				return ValueNone, cacheVisUnknown, err
			}
			res = v
		}

		if !plusSuper {
			break
		}
	}

	if res.IsNone() {
		return ValueNone, cacheVisHidden, nil
	}

	if currentVisibility == ast.ObjectFieldInherit {
		if (t.Flags & objFlagHasHidden) == 0 {
			return res, cacheVisVisible, nil
		}
		if !needVis {
			return res, cacheVisUnknown, nil
		}
		for j := layerOffset - 1; j >= 0; j-- {
			layer := layerAt(layers, j)

			fieldIndex := layer.findField(key, alloc)
			if fieldIndex == -1 {
				continue
			}

			visibility, _, tombstone := EvalFieldMeta(layer.Meta()[fieldIndex])

			if tombstone {
				j -= int(uint32(layer.Dyn.Values[fieldIndex]))
				continue
			}

			if visibility != ast.ObjectFieldInherit {
				currentVisibility = visibility
				break
			}
		}
	}

	if currentVisibility == ast.ObjectFieldHidden {
		return res, cacheVisHidden, nil
	}
	return res, cacheVisVisible, nil
}

func (t *Object) HasField(key uint32, includeHidden bool, ctx Context) bool {
	if c := t.Cache; c != nil {
		var val Value
		var visState uint8
		var found bool
		if val, visState, found = c.getInline(key); !found && c.fieldCache != nil {
			val, visState, found = c.getRawSlow(key)
		}
		if found {
			if val == ValueNone {
				return false
			}
			if includeHidden {
				return true
			}
			if visState != cacheVisUnknown {
				return visState == cacheVisVisible
			}
			vis := t.resolveFieldVisibility(key, ctx)
			c.setRaw(key, val, vis, ctx)
			return vis == cacheVisVisible
		}
	}

	var exists bool
	var visState uint8 = cacheVisVisible

	layers := t.GetLayers(ctx)
	alloc := ctx.State.Allocator
	for layerOffset := len(layers) - 1; layerOffset >= 0; layerOffset-- {
		layer := layerAt(layers, layerOffset)
		fieldIndex := layer.findField(key, alloc)
		if fieldIndex == -1 {
			continue
		}

		visibility, _, tombstone := EvalFieldMeta(layer.Meta()[fieldIndex])
		if tombstone {
			layerOffset -= int(uint32(layer.Dyn.Values[fieldIndex]))
			continue
		}

		if visibility == ast.ObjectFieldInherit {
			if (t.Flags & objFlagHasHidden) != 0 {
				if includeHidden {
					visState = cacheVisUnknown
				} else {
					for j := layerOffset - 1; j >= 0; j-- {
						l := layerAt(layers, j)
						fi := l.findField(key, alloc)
						if fi == -1 {
							continue
						}
						v, _, ts := EvalFieldMeta(l.Meta()[fi])
						if ts {
							j -= int(uint32(l.Dyn.Values[fi]))
							continue
						}
						if v != ast.ObjectFieldInherit {
							visibility = v
							break
						}
					}
					if visibility == ast.ObjectFieldHidden {
						visState = cacheVisHidden
					}
				}
			}
		} else if visibility == ast.ObjectFieldHidden {
			visState = cacheVisHidden
		}

		exists = true
		break
	}

	if exists {
		t.cacheSet(key, FieldCacheUnevaled, visState, ctx)
		return includeHidden || visState == cacheVisVisible
	}

	t.cacheSet(key, ValueNone, cacheVisHidden, ctx)
	return false
}

func (t *Object) getScope(layerIndex int, layer *Layer, ctx Context) (uintptr, error) {
	if len(layer.LocalNodes()) > 0 {
		return t.createLayerScope(layerIndex, ctx)
	}
	return layer.ParentScopePtr, nil
}

//go:noinline
func (t *Object) createLayerScope(layerIndex int, ctx Context) (uintptr, error) {
	layers := t.GetLayers(ctx)
	layer := layerAt(layers, layerIndex)
	var scopes []uintptr
	if t.ScopesPtr == 0 {
		scopes = arena.AllocZeroed[uintptr](ctx.State.Allocator, len(layers))
		t.ScopesPtr = uintptr(unsafe.Pointer(unsafe.SliceData(scopes)))
	} else {
		scopes = unsafe.Slice((*uintptr)(resolveUintptr(t.ScopesPtr)), len(layers))
	}

	scopePtr := scopes[layerIndex]
	if scopePtr != 0 {
		return scopePtr, nil
	}

	localNodes := layer.LocalNodes()
	slots, scopeId := ctx.NewScope(layer.ParentScopePtr, len(localNodes))

	for i, node := range localNodes {
		slots[i] = evaluateIRLazy(node, scopeId, ctx)
	}

	scopes[layerIndex] = scopeId

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

func (t *Object) HasVisibleFields(ctx Context) bool {
	layers := t.GetLayers(ctx)
	if len(layers) == 1 {
		layer := layerAt(layers, 0)
		if tmpl := layer.Template; tmpl != nil {
			return tmpl.VisibleCount > 0
		}
		for _, m := range layer.Dyn.Meta {
			if (m&FlagTombstone) == 0 && ast.ObjectFieldHide(m&MaskVisibility) != ast.ObjectFieldHidden {
				return true
			}
		}
		return false
	}

	fps := compileObjectPlan(t, ctx)
	for i := range fps {
		if !fps[i].IsHidden() {
			return true
		}
	}
	return false
}

func (t *Object) Length(ctx Context) int {
	layers := t.GetLayers(ctx)
	if len(layers) == 1 {
		layer := layerAt(layers, 0)
		if tmpl := layer.Template; tmpl != nil {
			return tmpl.VisibleCount
		}
		length := 0
		for _, m := range layer.Dyn.Meta {
			if (m&FlagTombstone) == 0 && ast.ObjectFieldHide(m&MaskVisibility) != ast.ObjectFieldHidden {
				length++
			}
		}
		return length
	}

	fps := compileObjectPlan(t, ctx)

	length := 0
	for i := range fps {
		if fps[i].IsHidden() {
			continue
		}
		length++
	}
	return length
}

func (t *Object) fillLayers(dest []uintptr, idx int) int {
	if t.Layers != nil {
		return idx + copy(dest[idx:], t.Layers)
	}

	if t.LeftPtr == 0 || t.RightPtr == 0 {
		return idx
	}

	l := (*Object)(resolveUintptr(t.LeftPtr))
	idx = l.fillLayers(dest, idx)

	r := (*Object)(resolveUintptr(t.RightPtr))
	return r.fillLayers(dest, idx)
}

func (t *Object) GetLayers(ctx Context) []uintptr {
	if t.Layers != nil {
		return t.Layers
	}

	total := int(t.LayerCount)
	if total == 0 {
		return nil
	}

	layers := arena.Alloc[uintptr](ctx.State.Allocator, total)
	n := t.fillLayers(layers, 0)
	layers = layers[:n]

	t.Layers = layers
	t.LeftPtr = 0
	t.RightPtr = 0

	return t.Layers
}

func MergeObjects(leftPtr, rightPtr uintptr, ctx Context) *Object {
	o := arena.CreateZeroed[Object](ctx.State.Allocator)
	o.LeftPtr = leftPtr
	o.RightPtr = rightPtr
	l := (*Object)(resolveUintptr(leftPtr))
	r := (*Object)(resolveUintptr(rightPtr))
	o.LayerCount = l.LayerCount + r.LayerCount
	o.AssertionState = l.AssertionState & r.AssertionState & AssertStatusChecked
	o.Flags = l.Flags | r.Flags
	return o
}

const (
	fpFlagHasFirst uint8 = 1 << 0
	fpFlagIsClosed uint8 = 1 << 1
)

type FieldPlan struct {
	First            LayerRef
	ExtraPtr         uintptr
	ExtraLen         uint32
	ExtraCap         uint32
	KeyId            uint32
	ShadowUntilLayer uint16
	Visibility       uint8
	Flags            uint8
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

func (fp *FieldPlan) appendLayerRef(allocator *arena.Allocator, ref LayerRef) {
	if (fp.Flags & fpFlagHasFirst) == 0 {
		fp.First = ref
		fp.Flags |= fpFlagHasFirst
		return
	}
	if fp.ExtraLen == fp.ExtraCap {
		newCap := fp.ExtraCap * 2
		if newCap == 0 {
			newCap = 4
		}
		var oldSlice []LayerRef
		if fp.ExtraLen > 0 {
			oldSlice = unsafe.Slice((*LayerRef)(resolveUintptr(fp.ExtraPtr)), fp.ExtraLen)
		}
		newSlice := arena.Alloc[LayerRef](allocator, int(newCap))
		copy(newSlice, oldSlice)
		fp.ExtraPtr = uintptr(unsafe.Pointer(unsafe.SliceData(newSlice)))
		fp.ExtraCap = newCap
	}
	extra := unsafe.Slice((*LayerRef)(resolveUintptr(fp.ExtraPtr)), fp.ExtraCap)
	extra[fp.ExtraLen] = ref
	fp.ExtraLen++
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
	if len(plans) <= 1 {
		return plans
	}

	inter := ctx.State.Interner

	if naturalSort {
		slices.SortFunc(plans, func(a, b FieldPlan) int {
			return naturalStringSort(inter.Get(a.KeyId), inter.Get(b.KeyId))
		})
		return plans
	}

	slices.SortFunc(plans, func(a, b FieldPlan) int {
		return cmp.Compare(inter.Get(a.KeyId), inter.Get(b.KeyId))
	})

	return plans
}

func compileObjectPlan(obj *Object, ctx Context) []FieldPlan {
	allocator := ctx.State.Allocator

	layers := obj.GetLayers(ctx)

	if len(layers) == 1 {
		layer := layerAt(layers, 0)
		keys := layer.Keys()
		n := len(keys)
		if n == 0 {
			return nil
		}
		meta := layer.Meta()
		plans := arena.Alloc[FieldPlan](allocator, n)
		out := 0
		for f, keyID := range keys {
			vis, _, tombstone := EvalFieldMeta(meta[f])
			if tombstone {
				continue
			}
			plans[out] = FieldPlan{
				First:            LayerRef{0, int32(f)},
				KeyId:            keyID,
				ShadowUntilLayer: math.MaxUint16,
				Visibility:       uint8(vis),
				Flags:            fpFlagHasFirst | fpFlagIsClosed,
			}
			out++
		}
		return plans[:out]
	}

	maxKeys := 0
	for i := range layers {
		maxKeys += len(layerAt(layers, i).Keys())
	}
	if maxKeys == 0 {
		return nil
	}

	plans := arena.Alloc[FieldPlan](allocator, maxKeys)
	plans = plans[:0]

	var planIdxMap *utils.DescriptorTable
	useMap := maxKeys > MaxPlanLinearKeys
	if useMap {
		planIdxMap = utils.NewEmptyDescriptorTable(allocator, maxKeys)
	}

	var validate bool

	for l := len(layers) - 1; l >= 0; l-- {
		layer := layerAt(layers, l)
		keys := layer.Keys()
		meta := layer.Meta()

		for f, keyID := range keys {

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
					ShadowUntilLayer: math.MaxUint16,
				})
				pIdx = len(plans) - 1

				if useMap {
					planIdxMap.Append(keyID)
				}
			}

			vis, plus, tombstone := EvalFieldMeta(meta[f])

			plan := &plans[pIdx]

			if plan.IsInherit() {
				switch vis {
				case ast.ObjectFieldHidden:
					plan.Visibility = uint8(ast.ObjectFieldHidden)
				case ast.ObjectFieldVisible:
					plan.Visibility = uint8(ast.ObjectFieldVisible)
				}
			}

			if (plan.Flags & fpFlagIsClosed) != 0 {
				continue
			}

			if plan.ShadowUntilLayer != math.MaxUint16 && l >= int(plan.ShadowUntilLayer) {
				continue
			}

			if tombstone {
				validate = true
				plan.ShadowUntilLayer = uint16(l) - uint16(layer.Dyn.Values[f])
				continue
			}

			plan.appendLayerRef(allocator, LayerRef{int32(l), int32(f)})

			if !plus {
				plan.Flags |= fpFlagIsClosed
			}
		}
	}

	if !validate {
		return plans
	}

	validPlans := plans[:0]
	for _, p := range plans {
		if (p.Flags & fpFlagHasFirst) != 0 {
			validPlans = append(validPlans, p)
		}
	}

	return validPlans
}

func (t *FieldPlan) GetValue(obj *Object, ctx Context) (Value, error) {
	if obj.Cache != nil {
		if value, _, found := obj.Cache.getRaw(t.KeyId); found && value != FieldCacheUnevaled {
			return value, nil
		}
	}

	err := runAssertions(obj, ctx)
	if err != nil {
		return ValueNone, err
	}

	if (t.Flags & fpFlagHasFirst) == 0 {
		return ValueNone, fmt.Errorf("no layers passed to plan.getValue")
	}

	value, err := getValue(obj, int(t.First.LayerIdx), int(t.First.FieldIdx), ctx)
	if err != nil {
		return ValueNone, err
	}

	if t.ExtraLen > 0 {
		extra := unsafe.Slice((*LayerRef)(resolveUintptr(t.ExtraPtr)), t.ExtraLen)
		for _, overlayRef := range extra {
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
	}

	vis := cacheVisVisible
	if t.IsHidden() {
		vis = cacheVisHidden
	}
	obj.cacheSet(t.KeyId, value, vis, ctx)

	return value, nil
}

// getSingleLayerFieldValue evaluates a field of a single-layer static object.
// Constant fields bypass FieldCache allocation and getValue completely.
func getSingleLayerFieldValue(obj *Object, tmpl *LayerTemplate, keyId uint32, fieldIdx int, ctx Context) (Value, error) {
	node := tmpl.Nodes[fieldIdx]
	if node.Kind == IRKindConst {
		return node.ConstVal, nil
	}
	if obj.Cache != nil {
		if val, _, found := obj.Cache.getRaw(keyId); found && val != FieldCacheUnevaled {
			return val, nil
		}
	}
	val, err := getValue(obj, 0, fieldIdx, ctx)
	if err != nil {
		return ValueNone, err
	}
	obj.cacheSet(keyId, val, cacheVisVisible, ctx)
	return val, nil
}

func manifestObject(obj *Object, ctx Context) (map[string]any, error) {
	if ctx.Depth >= ctx.State.MaxStack {
		return nil, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
	}
	ctx.Depth++

	err := runAssertions(obj, ctx)
	if err != nil {
		return nil, err
	}

	plans := CompileObjectPlan(obj, ctx)

	res := make(map[string]any, len(plans))
	for i := range plans {
		plan := &plans[i]
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

	res := arena.Alloc[NamedValue](ctx.State.Allocator, len(plans))
	var index int
	for i := range plans {
		plan := &plans[i]

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

	l := layerAt(layers, layerId)

	var n *IRNode
	if t := l.Template; t != nil {
		n = t.Nodes[fieldId]
		if n.Kind == IRKindConst {
			return n.ConstVal, nil
		}
	} else {
		d := l.Dyn
		if d.Values != nil {
			return d.Values[fieldId], nil
		}
		n = d.Nodes[fieldId]
		if n.Kind == IRKindConst {
			return n.ConstVal, nil
		}
	}

	evalCtx := ctx
	evalCtx.SuperOffset = uint32(len(layers) - 1 - layerId)

	if evalCtx.Depth >= evalCtx.State.MaxStack {
		return ValueNone, MakeRuntimeError(fmt.Errorf("max stack frames exceeded."))
	}
	evalCtx.Depth++

	var scopeId uintptr
	var err error

	if l.Template == nil {
		if s := l.UnpackFieldScopes(); s != nil {
			scopeId = s[fieldId]
		} else {
			scopeId, err = obj.getScope(layerId, l, evalCtx)
			if err != nil {
				return ValueNone, err
			}
		}
	} else {
		scopeId, err = obj.getScope(layerId, l, evalCtx)
		if err != nil {
			return ValueNone, err
		}
	}

	return EvaluateNode(n, scopeId, evalCtx)
}

func runAssertions(obj *Object, ctx Context) error {
	if obj.AssertionState != AssertStatusUnchecked {
		return nil
	}
	return runAssertionsSlow(obj, ctx)
}

//go:noinline
func runAssertionsSlow(obj *Object, ctx Context) error {
	obj.AssertionState = AssertStatusChecking

	layers := obj.GetLayers(ctx)

	for i := len(layers) - 1; i >= 0; i-- {
		layer := layerAt(layers, i)

		asserts := layer.Asserts()
		if len(asserts) == 0 {
			continue
		}

		evalCtx := ctx
		evalCtx.SuperOffset = uint32(len(layers) - 1 - i)

		scopeId, err := obj.getScope(i, layer, evalCtx)
		if err != nil {
			return err
		}

		for _, n := range asserts {
			val, err := EvaluateNode(n, scopeId, evalCtx)
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
	layers := obj.GetLayers(ctx)
	if !inclhidden && len(layers) == 1 && layerAt(layers, 0).Template != nil {
		tmpl := layerAt(layers, 0).Template
		if tmpl.VisibleCount == 0 {
			return nil
		}
		order := tmpl.GetSortedVisible(ctx.State.Interner)
		res := arena.Alloc[Value](ctx.State.Allocator, len(order))
		for i, fIdx := range order {
			res[i] = MakeStringConst(tmpl.Keys[fIdx])
		}
		return res
	}

	plans := CompileObjectPlan(obj, ctx)
	if len(plans) == 0 {
		return nil
	}

	res := arena.Alloc[Value](ctx.State.Allocator, len(plans))
	index := 0
	for i := range plans {
		fp := &plans[i]
		if inclhidden || !fp.IsHidden() {
			res[index] = MakeStringConst(fp.KeyId)
			index++
		}
	}

	return res[:index]
}

func GetObjectValues(obj *Object, ctx Context, inclHidden bool) ([]Value, error) {
	plans := CompileObjectPlan(obj, ctx)
	if len(plans) == 0 {
		return nil, nil
	}

	res := arena.Alloc[Value](ctx.State.Allocator, len(plans))
	index := 0
	for i := range plans {
		plan := &plans[i]
		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		res[index] = val
		index++
	}

	return res[:index], nil
}

func GetObjectKeysValues(obj *Object, ctx Context, inclHidden bool) ([]Value, error) {
	allocator := ctx.State.Allocator

	plans := CompileObjectPlan(obj, ctx)
	if len(plans) == 0 {
		return nil, nil
	}

	keyIdx := ctx.State.Interner.InternStatic("key")
	valueIdx := ctx.State.Interner.InternStatic("value")

	sharedKeys := arena.Alloc[uint32](allocator, 2)
	sharedKeys[0] = keyIdx
	sharedKeys[1] = valueIdx

	sharedMeta := arena.Alloc[uint8](allocator, 2)
	sharedMeta[0] = DefaultFieldMeta
	sharedMeta[1] = DefaultFieldMeta

	res := arena.Alloc[Value](allocator, len(plans))

	var index int
	for i := range plans {
		plan := &plans[i]
		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		layer, dl := NewDynamicLayer(allocator)
		dl.Keys = sharedKeys
		dl.Values = arena.Alloc[Value](allocator, 2)
		dl.Values[0] = MakeStringConst(plan.KeyId)
		dl.Values[1] = val
		dl.Meta = sharedMeta

		kvObj := NewSingleLayerObject(allocator, layer)
		res[index] = MakeObjectValue(kvObj)
		index++
	}

	return res[:index], nil
}

func GetObjectKeysValuesArray(obj *Object, ctx Context, inclHidden bool) ([]uint32, []Value, error) {
	plans := CompileObjectPlan(obj, ctx)
	if len(plans) == 0 {
		return nil, nil, nil
	}

	allocator := ctx.State.Allocator
	keys := arena.Alloc[uint32](allocator, len(plans))
	vals := arena.Alloc[Value](allocator, len(plans))
	index := 0
	for i := range plans {
		plan := &plans[i]
		if !inclHidden && plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, nil, err
		}

		keys[index] = plan.KeyId
		vals[index] = val
		index++
	}

	return keys[:index], vals[:index], nil
}

func (t *Object) Prune(ctx Context) (Value, error) {
	err := runAssertions(t, ctx)
	if err != nil {
		return ValueNone, err
	}

	plans := compileObjectPlan(t, ctx)

	n := len(plans)
	allocator := ctx.State.Allocator

	layer, dl := NewDynamicLayer(allocator)
	dl.Keys = arena.Alloc[uint32](allocator, n)
	dl.Values = arena.Alloc[Value](allocator, n)
	dl.Meta = arena.Alloc[uint8](allocator, n)

	index := 0
	for i := range plans {
		plan := &plans[i]
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

		dl.Keys[index] = plan.KeyId
		dl.Values[index] = prunedVal
		dl.Meta[index] = CreateFieldMeta(ast.ObjectFieldHide(plan.Visibility), false)
		index++
	}

	if index < n {
		dl.Keys = dl.Keys[:index]
		dl.Values = dl.Values[:index]
		dl.Meta = dl.Meta[:index]
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

	for i < len(planAs) {
		if !planAs[i].IsHidden() {
			return false, nil
		}
		i++
	}
	for j < len(planBs) {
		if !planBs[j].IsHidden() {
			return false, nil
		}
		j++
	}
	return true, nil
}
