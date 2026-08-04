package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	// TODO: benchmark if stack or arena arrays are better
	objs := arena.Alloc[*evaluator.Object](ctx.State.Allocator, len(inputArr))
	arena.MemclrSlice(objs)

	totalLayers := 0
	for i, v := range inputArr {
		obj, err := v.EvalObject(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		objs[i] = obj
		totalLayers += len(obj.GetLayers(ctx))
	}

	layers := arena.Alloc[*evaluator.Layer](ctx.State.Allocator, totalLayers)
	arena.MemclrSlice(layers)

	index := 0
	for _, o := range objs {
		objLayers := o.GetLayers(ctx)
		copy(layers[index:], objLayers)
		index += len(objLayers)
	}

	obj := arena.Create[evaluator.Object](ctx.State.Allocator)
	arena.Memclr(obj)
	obj.Layers = layers

	return evaluator.MakeObjectValue(obj), nil
}

func builtin_flatMapArray(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	mapperFunc, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	inputArr, err := args[1].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	mapperFuncInput := arena.Alloc[evaluator.NamedValue](ctx.State.Allocator, 1)

	// TODO: benchmark if stack or arena arrays are better
	subArrayValues := make([]evaluator.Value, len(inputArr))
	totalLen := 0
	for i, v := range inputArr {
		mapperFuncInput[0] = evaluator.NamedValue{Value: v}
		out, err := mapperFunc.FunctionExec(mapperFuncInput, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if !out.IsArray() {
			return evaluator.ValueNone, fmt.Errorf("unexpected response type of builtin_flatMapArray map func call: %s, expected array", out.Type().String())
		}
		subArrayValues[i] = out
		totalLen += len(out.Array())
	}

	res, val := evaluator.MakeArraySized(totalLen, ctx)
	index := 0
	for _, out := range subArrayValues {
		arr := out.Array()
		copy(res[index:], arr)
		index += len(arr)
	}

	return val, nil
}
