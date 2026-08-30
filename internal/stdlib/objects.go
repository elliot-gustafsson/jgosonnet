package stdlib

import (
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
)

func liftObjectToValueErr(f func(evaluator.Value, evaluator.Context) (evaluator.Value, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if !a.IsObject() {
			return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, a.Type())
		}
		res, err := f(a, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		return res, nil
	}
}

func liftObjectStringToValueErr(f func(evaluator.Value, string, evaluator.Context) (evaluator.Value, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if !a.IsObject() {
			return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, a.Type())
		}
		b, err := args[1].EvalString(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		res, err := f(a, b, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		return res, nil
	}
}

func std_get(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.get(o, f, default=null, inc_hidden=true)

	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !objVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}

	field, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	defaultVal := evaluator.MakeNull()
	if !args[2].IsNone() {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		defaultVal = v
	}

	inclHidden := true
	if !args[3].IsNone() {
		v, err := args[3].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		inclHidden = v
	}

	keyId := ctx.State.Interner.Intern(field)

	childCtx := ctx
	childCtx.Self = objVal

	val, visible, err := objVal.Object().GetField(keyId, childCtx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if val.IsNone() || !visible && !inclHidden {
		return defaultVal, nil
	}
	return val, nil
}

var std_objectFields = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res := evaluator.GetObjectFields(v.Object(), ctx, false)
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectFieldsAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res := evaluator.GetObjectFields(v.Object(), ctx, true)
	return evaluator.MakeArray(res, ctx), nil
})

func std_objectFieldsEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	inclHidden, err := args[1].EvalBool(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res := evaluator.GetObjectFields(obj, ctx, inclHidden)
	return evaluator.MakeArray(res, ctx), nil
}

var std_objectValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(), ctx, false)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(), ctx, true)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectKeysValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(), ctx, false)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectKeysValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(), ctx, true)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectHas = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.State.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object().GetField(keyId, subCtx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
})

var std_objectHasAll = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.State.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object().GetField(keyId, subCtx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
})

func std_objectHasEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !objVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}

	fname, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	hidden, err := args[2].EvalBool(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	keyId := ctx.State.Interner.Intern(fname)

	subCtx := ctx
	subCtx.Self = objVal

	value, visible, err := objVal.Object().GetField(keyId, subCtx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if value.IsNone() || (!visible && !hidden) {
		return evaluator.MakeBool(false), nil
	}

	return evaluator.MakeBool(true), nil
}

func std_mapWithkey(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	obj, err := args[1].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	keys, vals, err := evaluator.GetObjectKeysValuesArray(obj, ctx, false)
	if err != nil {
		return evaluator.ValueNone, err
	}

	fieldCount := len(keys)
	allocator := ctx.State.Allocator

	layer := arena.Create[evaluator.Layer](allocator)
	arena.Memclr(layer)

	layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
	layer.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
	layer.Meta = arena.Alloc[uint8](allocator, fieldCount)

	resObj := evaluator.NewSingleLayerObject(allocator, layer)
	resObjVal := evaluator.MakeObjectValue(resObj)

	mapCtx := ctx
	mapCtx.Self = resObjVal

	allArgs := arena.Alloc[evaluator.NamedValue](allocator, len(keys)*2)
	for i, k := range keys {
		v := vals[i]
		idx := i * 2

		keyString := ctx.State.Interner.Get(k)

		allArgs[idx] = evaluator.NamedValue{Value: evaluator.MakeString(keyString, ctx)}
		allArgs[idx+1] = evaluator.NamedValue{Value: v}

		n := arena.Create[evaluator.GoCallbackNode](allocator)
		arena.Memclr(n)
		*n = evaluator.GoCallbackNode{
			FuncVal: mapFunc,
			Args:    allArgs[idx : idx+2],
		}

		layer.Keys[i] = k
		layer.Values[i] = evaluator.NewThunk(evaluator.ThunkTypeGoCallback, unsafe.Pointer(n), 0, mapCtx)
		layer.Meta[i] = evaluator.DefaultFieldMeta

	}

	return resObjVal, nil
}

func std_objectRemoveKey(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !objVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}
	obj := objVal.Object()

	key, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	allocator := ctx.State.Allocator

	keyId := ctx.State.Interner.Intern(key)

	existingLayers := obj.GetLayers(ctx)

	tombstoneLayer := arena.Create[evaluator.Layer](allocator)
	arena.Memclr(tombstoneLayer)

	tombstoneLayer.Keys = arena.Alloc[uint32](allocator, 1)
	tombstoneLayer.Keys[0] = keyId

	tombstoneLayer.Values = arena.Alloc[evaluator.Value](allocator, 1)
	tombstoneLayer.Values[0] = evaluator.MakeTombstoneValue(len(existingLayers))

	tombstoneLayer.Meta = arena.Alloc[uint8](allocator, 1)
	tombstoneLayer.Meta[0] = evaluator.FlagTombstone | evaluator.DefaultFieldMeta

	newLen := len(existingLayers) + 1
	newLayers := arena.Alloc[*evaluator.Layer](allocator, newLen)
	arena.MemclrSlice(newLayers)
	copy(newLayers, existingLayers)
	newLayers[newLen-1] = tombstoneLayer

	resObj := arena.Create[evaluator.Object](allocator)
	arena.Memclr(resObj)

	resObj.Layers = newLayers
	return evaluator.MakeObjectValue(resObj), nil
}

func std_mergePatch(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	targetVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	patchVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return doMergePatch(targetVal, patchVal, ctx)
}

func doMergePatch(target, patch evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	var err error
	patch, err = patch.Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	// RFC 7396: If patch is not an object, it replaces target completely.
	if !patch.IsObject() {
		return patch, nil
	}

	patchObj := patch.Object()
	patchPlans := evaluator.CompileObjectPlan(patchObj, ctx)

	allocator := ctx.State.Allocator

	var targetObj *evaluator.Object
	var targetLayers []*evaluator.Layer
	if !target.IsNone() {
		target, err = target.Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if target.IsObject() {
			targetObj = target.Object()
			targetLayers = targetObj.GetLayers(ctx)
		}
	}

	fieldCount := len(patchPlans)

	layer := arena.Create[evaluator.Layer](allocator)
	arena.Memclr(layer)

	layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
	layer.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
	layer.Meta = arena.Alloc[uint8](allocator, fieldCount)

	useMap := fieldCount > evaluator.MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	var index int
	for _, plan := range patchPlans {
		if plan.IsHidden() {
			continue
		}

		// Evaluate patch field in patch's context
		val, err := plan.GetValue(patchObj, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		keyId := plan.KeyId

		// RFC 7396: Null patch value removes the field from target
		if val.IsNull() {
			if targetObj != nil {
				layer.Keys[index] = keyId
				layer.Values[index] = evaluator.MakeTombstoneValue(len(targetLayers))
				layer.Meta[index] = evaluator.FlagTombstone | evaluator.DefaultFieldMeta
				if useMap {
					layer.Index.Append(keyId)
				}
				index++
			}
			continue
		}

		// RFC 7396: Object patch value recursively merges with target[keyId] if target[keyId] is an object
		if val.IsObject() {
			var mergedVal evaluator.Value
			if targetObj != nil {
				targetFieldVal, vis, err := targetObj.GetField(keyId, ctx)
				if err != nil {
					return evaluator.ValueNone, err
				}
				if vis {
					targetFieldVal, err = targetFieldVal.Eval(ctx)
					if err != nil {
						return evaluator.ValueNone, err
					}
				}
				if vis && targetFieldVal.IsObject() {
					m, err := doMergePatch(targetFieldVal, val, ctx)
					if err != nil {
						return evaluator.ValueNone, err
					}
					mergedVal = m
				} else {
					m, err := doMergePatch(evaluator.ValueNone, val, ctx)
					if err != nil {
						return evaluator.ValueNone, err
					}
					mergedVal = m
				}
			} else {
				m, err := doMergePatch(evaluator.ValueNone, val, ctx)
				if err != nil {
					return evaluator.ValueNone, err
				}
				mergedVal = m
			}

			layer.Keys[index] = keyId
			layer.Values[index] = mergedVal
			layer.Meta[index] = evaluator.DefaultFieldMeta
			if useMap {
				layer.Index.Append(keyId)
			}
			index++
			continue
		}

		// Primitive, array, or boolean values replace target field directly
		layer.Keys[index] = keyId
		layer.Values[index] = val
		layer.Meta[index] = evaluator.DefaultFieldMeta
		if useMap {
			layer.Index.Append(keyId)
		}
		index++
	}

	if index < fieldCount {
		layer.Keys = layer.Keys[:index]
		layer.Values = layer.Values[:index]
		layer.Meta = layer.Meta[:index]
	}

	// If target was not an object, return a single-layer object containing non-null patch fields
	if targetObj == nil {
		res := evaluator.NewSingleLayerObject(allocator, layer)
		return evaluator.MakeObjectValue(res), nil
	}

	// Stack patch layer on top of existing target layers
	newLen := len(targetLayers) + 1
	newLayers := arena.Alloc[*evaluator.Layer](allocator, newLen)
	arena.MemclrSlice(newLayers)
	copy(newLayers, targetLayers)
	newLayers[newLen-1] = layer

	resObj := arena.Create[evaluator.Object](allocator)
	arena.Memclr(resObj)
	resObj.Layers = newLayers

	return evaluator.MakeObjectValue(resObj), nil
}
