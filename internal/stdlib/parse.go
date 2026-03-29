package stdlib

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/google/go-jsonnet/ast"
	"gopkg.in/yaml.v3"
)

func liftStringToValueErr(f func(string) (evaluator.Value, error), name string) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !a.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected string", name, a.Type().String())
		}
		return f(a.String(ctx))
	}
}

var std_parseInt = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("error: %s is not a base 10 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
}, "std.parseInt")

var std_parseOctal = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 8, 64)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("error: %s is not a base 8 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
}, "std.parseOctal")

var std_parseHex = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("error: %s is not a base 16 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
}, "std.parseHex")

func std_parseJson(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.parseJson: %d, expected 1", len(args))
	}
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !a.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.parseJson (arg 0): %s, expected string", a.Type().String())
	}
	var data any
	err = json.Unmarshal([]byte(a.String(ctx)), &data)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("failed to parse json, err: %w", err)
	}
	return rawDataToValue(data, ctx)
}

func std_parseYaml(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.parseYaml: %d, expected 1", len(args))
	}
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !a.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.parseYaml (arg 0): %s, expected string", a.Type().String())
	}
	var data any
	err = yaml.Unmarshal([]byte(a.String(ctx)), &data)
	if err != nil {
		return evaluator.Value{}, fmt.Errorf("failed to parse yaml, err: %w", err)
	}
	return rawDataToValue(data, ctx)
}

func rawDataToValue(x any, ctx evaluator.Context) (evaluator.Value, error) {
	switch data := x.(type) {
	case nil:
		return evaluator.MakeNull(), nil
	case float64:
		return evaluator.MakeNumber(data), nil
	case int:
		return evaluator.MakeNumber(float64(data)), nil
	case string:
		return evaluator.MakeString(data, ctx), nil
	case bool:
		return evaluator.MakeBool(data), nil
	case []any:
		res := make([]evaluator.Value, 0, len(data))
		for _, rv := range data {
			v, err := rawDataToValue(rv, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			res = append(res, v)
		}
		return evaluator.MakeArray(res, ctx), nil
	case map[string]any:

		fieldCount := len(data)

		layer := &evaluator.Layer{
			Keys:  make([]uint32, 0, fieldCount),
			Nodes: make(ast.Nodes, 0, fieldCount),
			Meta:  make([]uint8, 0, fieldCount),

			Index: make(map[uint32]int, fieldCount),
		}

		obj := evaluator.NewObject([]*evaluator.Layer{layer})

		index := 0
		for keyName, value := range data {
			keyId := ctx.Interner.Intern(keyName)

			v, err := rawDataToValue(value, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}

			layer.Keys = append(layer.Keys, keyId)
			layer.Values = append(layer.Values, v)
			layer.Meta = append(layer.Meta, evaluator.CreateFieldMeta(ast.ObjectFieldVisible, false))
			layer.Index[keyId] = index

			index++
		}

		return evaluator.MakeObject(obj, ctx), nil
	default:
		return evaluator.Value{}, fmt.Errorf("unahandled type %T when converting data to value", x)
	}

}
