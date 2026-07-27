package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	objs := make([]*evaluator.Object, len(inputArr))
	totalLayers := 0
	for i, v := range inputArr {
		obj, err := v.EvalObject(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		objs[i] = obj
		totalLayers += len(obj.GetLayers(ctx))
	}

	layers := make([]*evaluator.Layer, 0, totalLayers)
	for _, o := range objs {
		layers = append(layers, o.GetLayers(ctx)...)
	}

	objId := evaluator.NewObject(layers, ctx)

	return evaluator.MakeObjectValue(objId), nil
}

func builtin_flatMapArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapperFunc, err := args[0].EvalFunction(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	inputArr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	mapperFuncInput := ctx.State.Registry.NamedValueBufs.Alloc(1, 1)

	subArrays := make([][]evaluator.Value, len(inputArr))
	totalLen := 0
	for i, v := range inputArr {
		mapperFuncInput[0] = evaluator.NamedValue{Value: v}
		out, err := mapperFunc.Exec(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if !out.IsArray() {
			return evaluator.ValueNone, fmt.Errorf("unexpected response type of builtin_flatMapArray map func call: %s, expected array", out.Type().String())
		}
		arr := out.Array(ctx)
		subArrays[i] = arr
		totalLen += len(arr)
	}

	res := make([]evaluator.Value, 0, totalLen)
	for _, arr := range subArrays {
		res = append(res, arr...)
	}

	return evaluator.MakeArray(res, ctx), nil
}
