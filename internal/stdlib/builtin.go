package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to builtin_objectFlatMerge: %d, expected 1", len(args))
	}

	val, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !val.IsArray() {
		return evaluator.Value{}, fmt.Errorf("(builtin objectFlatMerge) unexpected type of arg 1: %s, expected array", val.Type().String())
	}

	// TODO: Think. Either just add layers or put all fields in a single layer. Test this later
	inputArr := val.Array(ctx)
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
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to builtin flatMapArray: %d, expected 2", len(args))
	}

	mapperFunc, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !mapperFunc.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type of arg 0: %s, expected function", mapperFunc.Type().String())
	}

	val, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !val.IsArray() {
		return evaluator.Value{}, fmt.Errorf("(builtin flatMapArray) unexpected type of arg 1: %s, expected array", mapperFunc.Type().String())
	}

	inputArr := val.Array(ctx)

	// Create the array once and mutate it to reduce objects on the heap
	mapperFuncInput := []evaluator.NamedValue{{}}
	mFunc := mapperFunc.Function(ctx)

	res := make([]evaluator.Value, 0, len(inputArr))
	for _, v := range inputArr {
		mapperFuncInput[0] = evaluator.NamedValue{Value: v}
		out, err := mFunc(mapperFuncInput, ctx)
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
