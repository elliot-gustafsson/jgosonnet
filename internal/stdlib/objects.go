package stdlib

import (
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/google/go-jsonnet/ast"
)

func liftObjectToValueErr(f func(evaluator.Value, evaluator.Context) (evaluator.Value, error), name string) evaluator.Func {
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

func liftObjectStringToValueErr(f func(evaluator.Value, string, evaluator.Context) (evaluator.Value, error), name string) evaluator.Func {
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

	obj, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !obj.IsObject() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, obj.Type())
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
	childCtx.Self = obj

	val, visible, err := obj.Object(ctx).GetField(keyId, childCtx)
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
}, "std.objectFields")

var std_objectFieldsAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res := evaluator.GetObjectFields(v.Object(ctx), ctx, true)
	return evaluator.MakeArray(res, ctx), nil
}, "std.objectFieldsAll")

var std_objectValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(ctx), ctx, false)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
}, "std.objectValues")

var std_objectValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectValues(v.Object(ctx), ctx, true)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
}, "std.objectValuesAll")

var std_objectKeysValues = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(ctx), ctx, false)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
}, "std.objectKeysValues")

var std_objectKeysValuesAll = liftObjectToValueErr(func(v evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	res, err := evaluator.GetObjectKeysValues(v.Object(ctx), ctx, true)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeArray(res, ctx), nil
}, "std.objectKeysValuesAll")

var std_objectHas = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object(ctx).GetField(keyId, subCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
}, "std.objectHas")

var std_objectHasAll = liftObjectStringToValueErr(func(v evaluator.Value, s string, ctx evaluator.Context) (evaluator.Value, error) {
	keyId := ctx.Interner.Intern(s)
	subCtx := ctx
	subCtx.Self = v
	value, _, err := v.Object(ctx).GetField(keyId, subCtx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeBool(!value.IsNone()), nil
}, "std.objectHasAll")

func std_mapWithkey(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	obj, err := args[1].EvalObject(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapFuncArgs := []evaluator.NamedValue{{}, {}}

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

	res := evaluator.NewObject([]*evaluator.Layer{layer})

	m := evaluator.CreateFieldMeta(ast.ObjectFieldInherit, false)
	for i, k := range keys {
		v := vals[i]

		keyString := ctx.Interner.Get(k)

		mapFuncArgs[0] = evaluator.NamedValue{Value: evaluator.MakeString(keyString, ctx)}
		mapFuncArgs[1] = evaluator.NamedValue{Value: v}

		x, err := mapFunc.Exec(mapFuncArgs, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		layer.Keys = append(layer.Keys, k)
		layer.Values = append(layer.Values, x)
		layer.Meta = append(layer.Meta, m)

	}

	return evaluator.MakeObject(res, ctx), nil
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

	// ----------------

	// objLayers := obj.GetLayers()

	// resLayers := make([]*evaluator.Layer, 0, len(objLayers))

	// for _, l := range objLayers {

	// 	fieldCount := len(l.Keys) - 1

	// 	newLayer := &evaluator.Layer{
	// 		Keys:  make([]uint32, 0, fieldCount),
	// 		Nodes: make(ast.Nodes, 0, fieldCount),
	// 		Meta:  make([]uint8, 0, fieldCount),
	// 	}

	// 	useMap := len(l.Keys) > evaluator.MaxLinearKeys
	// 	if useMap {
	// 		newLayer.Index = make(map[uint32]int, fieldCount)
	// 	}

	// 	index := 0
	// 	for i, k := range l.Keys {
	// 		if keyId == k {
	// 			continue
	// 		}

	// 		var v evaluator.Value
	// 		if l.Values != nil {
	// 			v = l.Values[fieldId]
	// 		} else {
	// 			x, err := evaluator.EvaluateNode()
	// 		}

	// 		newLayer.Keys = append(newLayer.Keys, k)
	// 		newLayer.Meta = append(newLayer.Meta, l.Meta[i])
	// 		newLayer.Nodes = append(newLayer.Nodes, v)

	// 		if useMap {
	// 			newLayer.Index[k] = index
	// 		}
	// 		index++
	// 	}
	// 	resLayers = append(resLayers, newLayer)
	// }

	// resObj := evaluator.NewObject(resLayers)
	// return evaluator.MakeObject(resObj, ctx), nil

	// ----------------

	// TODO: look over this, fairly sure is doesnt work for all cases...

	val, _, err := obj.GetField(keyId, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if val.IsNone() {
		// If key doesnt exist, just copy the object
		objLayers := obj.GetLayers()
		newLayers := make([]*evaluator.Layer, len(objLayers))
		copy(newLayers, objLayers)
		res := evaluator.NewObject(newLayers)
		return evaluator.MakeObject(res, ctx), nil
	}

	plans := evaluator.CompileObjectPlan(obj, ctx)

	fieldCount := len(plans) - 1
	layer := &evaluator.Layer{
		Keys:   make([]uint32, 0, fieldCount),
		Values: make([]evaluator.Value, 0, fieldCount),
		Meta:   make([]uint8, 0, fieldCount),
	}

	useMap := len(plans) > evaluator.MaxLinearKeys
	if useMap {
		layer.Index = make(map[uint32]int, fieldCount)
	}

	evalCtx := ctx
	evalCtx.Self = objVal

	for i, plan := range plans {
		if plan.KeyId == keyId {
			continue
		}

		v, err := plan.GetValue(obj, evalCtx)
		if err != nil {
			return evaluator.Value{}, err
		}

		layer.Keys = append(layer.Keys, plan.KeyId)
		layer.Values = append(layer.Values, v)
		layer.Meta = append(layer.Meta, evaluator.CreateFieldMeta(plan.Visibility, false))

		if useMap {
			layer.Index[plan.KeyId] = i
		}

	}

	resObj := evaluator.NewObject([]*evaluator.Layer{layer})
	return evaluator.MakeObject(resObj, ctx), nil
}
