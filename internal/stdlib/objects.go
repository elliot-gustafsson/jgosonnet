package stdlib

import (
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
)

func internStringVal(arg evaluator.NamedValue, ctx evaluator.Context) (uint32, error) {
	v, err := arg.Eval(ctx)
	if err != nil {
		return 0, err
	}
	if !v.IsString() {
		return 0, evaluator.TypeErrorSpecific(evaluator.ValueTypeString, v.Type())
	}
	if keyId := v.AsStringConst(); keyId != 0 {
		return keyId, nil
	}
	return ctx.State.Interner.Intern(v.String(ctx)), nil
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

	keyId, err := internStringVal(args[1], ctx)
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

	childCtx := ctx
	childCtx.Self = objVal

	if inclHidden {
		val, err := objVal.Object().GetField(keyId, childCtx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if val.IsNone() {
			return defaultVal, nil
		}
		return val, nil
	}

	val, visible, err := objVal.Object().GetFieldVis(keyId, childCtx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if val.IsNone() || !visible {
		return defaultVal, nil
	}
	return val, nil
}

func std_objectFields(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res := evaluator.GetObjectFields(obj, ctx, false)
	return evaluator.MakeArray(res, ctx), nil
}

func std_objectFieldsAll(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res := evaluator.GetObjectFields(obj, ctx, true)
	return evaluator.MakeArray(res, ctx), nil
}

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

func std_objectValues(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res, err := evaluator.GetObjectValues(obj, ctx, false)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_objectValuesAll(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res, err := evaluator.GetObjectValues(obj, ctx, true)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_objectKeysValues(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res, err := evaluator.GetObjectKeysValues(obj, ctx, false)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_objectKeysValuesAll(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res, err := evaluator.GetObjectKeysValues(obj, ctx, true)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeArray(res, ctx), nil
}

func std_objectHas(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	keyId, err := internStringVal(args[1], ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	hasField := obj.HasField(keyId, false, ctx)
	return evaluator.MakeBool(hasField), nil
}

func std_objectHasAll(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	keyId, err := internStringVal(args[1], ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	hasField := obj.HasField(keyId, true, ctx)
	return evaluator.MakeBool(hasField), nil
}

func std_objectHasEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !objVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}

	keyId, err := internStringVal(args[1], ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	hidden, err := args[2].EvalBool(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	hasField := objVal.Object().HasField(keyId, hidden, ctx)
	return evaluator.MakeBool(hasField), nil
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

	layer, dl := evaluator.NewDynamicLayer(allocator)

	dl.Keys = arena.Alloc[uint32](allocator, fieldCount)
	dl.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
	dl.Meta = arena.Alloc[uint8](allocator, fieldCount)

	resObj := evaluator.NewSingleLayerObject(allocator, layer)
	resObjVal := evaluator.MakeObjectValue(resObj)

	mapCtx := ctx
	mapCtx.Self = resObjVal

	allArgs := arena.Alloc[evaluator.NamedValue](allocator, len(keys)*2)
	for i, k := range keys {
		v := vals[i]
		idx := i * 2

		allArgs[idx] = evaluator.NamedValue{Value: evaluator.MakeStringConst(k)}
		allArgs[idx+1] = evaluator.NamedValue{Value: v}

		n := arena.CreateZeroed[evaluator.GoCallbackNode](allocator)
		*n = evaluator.GoCallbackNode{
			FuncVal: mapFunc,
			Args:    allArgs[idx : idx+2],
		}

		dl.Keys[i] = k
		dl.Values[i] = evaluator.NewThunk(evaluator.IRKindGoCallback, unsafe.Pointer(n), 0, mapCtx)
		dl.Meta[i] = evaluator.DefaultFieldMeta

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

	keyId, err := internStringVal(args[1], ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	allocator := ctx.State.Allocator

	existingLayers := obj.GetLayers(ctx)

	tombstoneLayer, dl := evaluator.NewDynamicLayer(allocator)

	dl.Keys = arena.Alloc[uint32](allocator, 1)
	dl.Keys[0] = keyId

	dl.Values = arena.Alloc[evaluator.Value](allocator, 1)
	dl.Values[0] = evaluator.MakeTombstoneValue(len(existingLayers))

	dl.Meta = arena.Alloc[uint8](allocator, 1)
	dl.Meta[0] = evaluator.FlagTombstone | evaluator.DefaultFieldMeta

	resObj := evaluator.NewMultiLayerObject(allocator, existingLayers, tombstoneLayer)
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
	var targetLayers []uintptr
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

	layer, dl := evaluator.NewDynamicLayer(allocator)

	dl.Keys = arena.Alloc[uint32](allocator, fieldCount)
	dl.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
	dl.Meta = arena.Alloc[uint8](allocator, fieldCount)

	useMap := fieldCount > evaluator.MaxLayerLinearKeys
	if useMap {
		dl.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	var index int
	for i := range patchPlans {
		plan := &patchPlans[i]
		if plan.IsHidden() {
			continue
		}

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
				dl.Keys[index] = keyId
				dl.Values[index] = evaluator.MakeTombstoneValue(len(targetLayers))
				dl.Meta[index] = evaluator.FlagTombstone | evaluator.DefaultFieldMeta
				if useMap {
					dl.Index.Append(keyId)
				}
				index++
			}
			continue
		}

		// RFC 7396: Object patch value recursively merges with target[keyId] if target[keyId] is an object
		if val.IsObject() {
			var mergedVal evaluator.Value
			if targetObj != nil {
				targetFieldVal, vis, err := targetObj.GetFieldVis(keyId, ctx)
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

			dl.Keys[index] = keyId
			dl.Values[index] = mergedVal
			dl.Meta[index] = evaluator.DefaultFieldMeta
			if useMap {
				dl.Index.Append(keyId)
			}
			index++
			continue
		}

		// Primitive, array, or boolean values replace target field directly
		dl.Keys[index] = keyId
		dl.Values[index] = val
		dl.Meta[index] = evaluator.DefaultFieldMeta
		if useMap {
			dl.Index.Append(keyId)
		}
		index++
	}

	if index < fieldCount {
		dl.Keys = dl.Keys[:index]
		dl.Values = dl.Values[:index]
		dl.Meta = dl.Meta[:index]
	}

	// If target was not an object, return a single-layer object containing non-null patch fields
	if targetObj == nil {
		res := evaluator.NewSingleLayerObject(allocator, layer)
		return evaluator.MakeObjectValue(res), nil
	}

	// Stack patch layer on top of existing target layers
	resObj := evaluator.NewMultiLayerObject(allocator, targetLayers, layer)
	return evaluator.MakeObjectValue(resObj), nil
}
