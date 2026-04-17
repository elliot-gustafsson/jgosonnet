package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	// TODO: Think. Either just add layers or put all fields in a single layer. Test this later
	layers := make([]*evaluator.Layer, 0, len(inputArr))
	for _, v := range inputArr {

		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !v.IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type of builtin_objectFlatMerge arg: %s, expected object", v.Type().String())
		}

		layers = append(layers, v.Object(ctx).GetLayers()...)
	}

	obj := evaluator.NewObject(layers)

	return evaluator.MakeObject(obj, ctx), nil
}

func builtin_flatMapArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapperFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	inputArr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	mapperFuncInput := []evaluator.NamedValue{{}}

	res := make([]evaluator.Value, 0, len(inputArr))
	for _, v := range inputArr {
		mapperFuncInput[0] = evaluator.NamedValue{Value: v}
		out, err := mapperFunc.Exec(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !out.IsArray() {
			return evaluator.Value{}, fmt.Errorf("unexpected response type of builtin_flatMapArray map func call: %s, expected array", out.Type().String())
		}
		res = append(res, out.Array(ctx)...)
	}

	return evaluator.MakeArray(res, ctx), nil
}
