package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func liftObjectToValueErr(f func(evaluator.Value, evaluator.Context) (evaluator.Value, error), name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		if !args[0].IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected object", name, args[0].Type().String())
		}
		res, err := f(args[0], ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		return res, nil
	}
}

func liftObjectStringToValueErr(f func(evaluator.Value, string, evaluator.Context) (evaluator.Value, error), name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 2 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 2", name, len(args))
		}
		if !args[0].IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected object", name, args[0].Type().String())
		}
		if !args[1].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 1): %s, expected string", name, args[0].Type().String())
		}
		res, err := f(args[0], args[1].String(ctx), ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		return res, nil
	}
}

func std_get(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.get(o, f, default=null, inc_hidden=true)
	if len(args) < 2 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.get %d, expected 2-4", len(args))
	}

	obj := args[0]
	if !obj.IsObject() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 0): %s, expected object", obj.Type().String())
	}

	field := args[1]
	if !field.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 1): %s, expected string", field.Type().String())
	}

	defaultVal := evaluator.MakeNull()
	if len(args) > 2 {
		defaultVal = args[2]
	}

	inclHidden := true
	if len(args) > 3 {
		if !args[3].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 3): %s, expected boolean", args[3].Type().String())
		}
		inclHidden = args[3].Bool()
	}

	keyId := ctx.Interner.Intern(field.String(ctx))

	childCtx := ctx
	childCtx.Self = args[0]

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
