package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func liftObjectToValueErr(f func(evaluator.Value, evaluator.Context) (evaluator.Value, error), name string) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !a.IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected object", name, a.Type().String())
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
		if len(args) != 2 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 2", name, len(args))
		}
		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !a.IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected object", name, a.Type().String())
		}
		b, err := args[1].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !b.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 1): %s, expected string", name, b.Type().String())
		}
		res, err := f(a, b.String(ctx), ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		return res, nil
	}
}

func std_get(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.get(o, f, default=null, inc_hidden=true)
	if len(args) < 2 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.get %d, expected 2-4", len(args))
	}

	obj, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !obj.IsObject() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 0): %s, expected object", obj.Type().String())
	}

	field, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !field.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 1): %s, expected string", field.Type().String())
	}

	defaultVal := evaluator.MakeNull()
	if len(args) > 2 {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		defaultVal = v
	}

	inclHidden := true
	if len(args) > 3 {
		v, err := args[3].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.get (arg 3): %s, expected boolean", v.Type().String())
		}
		inclHidden = v.Bool()
	}

	keyId := ctx.Interner.Intern(field.String(ctx))

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
