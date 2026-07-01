package stdlib

import (
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func liftObjectToValueErr(f func(evaluator.Value, evaluator.Context) (evaluator.Value, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !a.IsObject() {
			return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, a.Type())
		}
		res, err := f(a, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		return res, nil
	}
}

func liftObjectStringToValueErr(f func(evaluator.Value, string, evaluator.Context) (evaluator.Value, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !a.IsObject() {
			return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, a.Type())
		}
		b, err := args[1].EvalString(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res, err := f(a, b, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		return res, nil
	}
}

func std_get(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.get(o, f, default=null, inc_hidden=true)

	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !objVal.IsObject() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}

	field, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	defaultVal := evaluator.MakeNull()
	if !args[2].IsNone() {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		defaultVal = v
	}

	inclHidden := true
	if !args[3].IsNone() {
		v, err := args[3].EvalBool(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		inclHidden = v
	}

	keyId := ctx.Interner.Intern(field)

	childCtx := ctx
	childCtx.Self = objVal

	val, visible, err := objVal.Object(ctx).GetField(keyId, childCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if val.IsNone() || !visible && !inclHidden {
		return defaultVal, nil
	}
	return val, nil
}

var std_objectFields = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res := evaluator.GetObjectFields(v.Object(ctx), ctx, false)
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectFieldsAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res := evaluator.GetObjectFields(v.Object(ctx), ctx, true)
	return evaluator.MakeArray(res, ctx), nil
})

func std_objectFieldsEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	inclHidden, err := args[1].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	res := evaluator.GetObjectFields(obj, ctx, inclHidden)
	return evaluator.MakeArray(res, ctx), nil
}

var std_objectValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(ctx), ctx, false)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(ctx), ctx, true)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectKeysValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(ctx), ctx, false)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectKeysValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(ctx), ctx, true)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
})

var std_objectHas = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object(ctx).GetField(keyId, subCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
})

var std_objectHasAll = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object(ctx).GetField(keyId, subCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
})

func std_objectHasEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !objVal.IsObject() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}

	fname, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	hidden, err := args[2].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	keyId := ctx.Interner.Intern(fname)

	subCtx := ctx
	subCtx.Self = objVal

	value, visible, err := objVal.Object(ctx).GetField(keyId, subCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if value.IsNone() || (!visible && !hidden) {
		return evaluator.MakeBool(false), nil
	}

	return evaluator.MakeBool(true), nil
}

func std_mapWithkey(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	obj, err := args[1].EvalObject(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	keys, vals, err := evaluator.GetObjectKeysValuesArray(obj, ctx, false)
	if err != nil {
		return evaluator.Value{}, err
	}

	fieldCount := len(keys)

	layer := &evaluator.Layer{
		Keys:   make([]uint32, 0, fieldCount),
		Values: make([]evaluator.Value, 0, fieldCount),
		Meta:   make([]uint8, 0, fieldCount),
	}

	resId := evaluator.NewObject([]*evaluator.Layer{layer}, ctx)

	resObjVal := evaluator.MakeObjectValue(resId)

	mapCtx := ctx
	mapCtx.Self = resObjVal

	allArgs := make([]evaluator.NamedValue, len(keys)*2)
	for i, k := range keys {
		v := vals[i]
		idx := i * 2

		keyString := ctx.Interner.Get(k)

		allArgs[idx] = evaluator.NamedValue{Value: evaluator.MakeString(keyString, ctx)}
		allArgs[idx+1] = evaluator.NamedValue{Value: v}

		n := &evaluator.GoCallbackNode{
			Func: mapFunc,
			Args: allArgs[idx : idx+2],
		}

		thunk := evaluator.NewThunk(n, 0, mapCtx)

		layer.Keys = append(layer.Keys, k)
		layer.Values = append(layer.Values, evaluator.MakeThunk(thunk, mapCtx))
		layer.Meta = append(layer.Meta, evaluator.DefaultFieldMeta)

	}

	return resObjVal, nil
}

func std_objectRemoveKey(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !objVal.IsObject() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, objVal.Type())
	}
	obj := objVal.Object(ctx)

	key, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	keyId := ctx.Interner.Intern(key)

	existingLayers := obj.GetLayers(ctx)

	tombstoneLayer, _ := ctx.Registry.Layers.New()
	tombstoneLayer.Keys = []uint32{keyId}
	tombstoneLayer.Values = []evaluator.Value{evaluator.MakeTombstoneValue(len(existingLayers))}
	tombstoneLayer.Meta = []uint8{evaluator.FlagTombstone | evaluator.DefaultFieldMeta}

	newLen := len(existingLayers) + 1
	newLayers := ctx.Registry.LayerBufs.Alloc(newLen, newLen)
	copy(newLayers, existingLayers)
	newLayers[newLen-1] = tombstoneLayer

	resObjId := evaluator.NewObject(newLayers, ctx)
	return evaluator.MakeObjectValue(resObjId), nil
}

func std_mergePatch(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	patchVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !patchVal.IsObject() {
		return patchVal, nil
	}
	// patchObj := patchVal.Object(ctx)

	targetVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var objVal evaluator.Value

	// TODO: think abt this
	// std.mergePatch({ a: 1 }, { b: 2, c: self.a }) should not work according to go-jsonnet
	// but with this kinda merge it does... so some solution where the objects are individually evaled is needed
	if targetVal.IsObject() {
		mergedObjId := evaluator.MergeObjects(targetVal.RefId(), patchVal.RefId(), ctx)
		objVal = evaluator.MakeObjectValue(mergedObjId)
	} else {
		// If target val is not an object, just scrap it and only use the patch object
		objVal = patchVal
	}
	obj := objVal.Object(ctx)

	plans := evaluator.CompileObjectPlan(obj, ctx)

	fieldCount := len(plans)

	layer := &evaluator.Layer{
		Keys:   make([]uint32, 0, fieldCount),
		Values: make([]evaluator.Value, 0, fieldCount),
		Meta:   make([]uint8, 0, fieldCount),
	}

	subCtx := ctx
	subCtx.Self = objVal

	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, subCtx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if val.IsNull() {
			continue
		}

		layer.Keys = append(layer.Keys, plan.KeyId)
		layer.Values = append(layer.Values, val)
		layer.Meta = append(layer.Meta, evaluator.DefaultFieldMeta)

	}

	resObjId := evaluator.NewObject([]*evaluator.Layer{layer}, ctx)

	return evaluator.MakeObjectValue(resObjId), nil

}
