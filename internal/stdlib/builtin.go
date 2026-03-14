package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func builtin_objectFlatMerge(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to builtin_objectFlatMerge: %d, expected 1", len(args))
	}

	val := args[0]

	if !val.IsArray() {
		return evaluator.Value{}, fmt.Errorf("(builtin objectFlatMerge) unexpected type of arg 1: %s, expected array", val.Type().String())
	}

	// TODO: Think. Either just add layers or put all fields in a single layer. Test this later
	inputArr := val.Array(ctx)
	layers := make([]*evaluator.Layer, 0, len(inputArr))
	for _, v := range inputArr {

		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		if !v.IsObject() {
			return evaluator.Value{}, fmt.Errorf("unexpected type of builtin_objectFlatMerge arg: %s, expected object", v.Type().String())
		}

		layers = append(layers, v.Object(ctx).Layers...)
	}

	obj := evaluator.NewObject(layers)

	return evaluator.MakeObject(obj, ctx), nil
}

func builtin_flatMapArray(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to builtin flatMapArray: %d, expected 2", len(args))
	}

	mapperFunc := args[0]

	if !mapperFunc.IsFunction() {
		return evaluator.Value{}, fmt.Errorf("unexpected type of arg 0: %s, expected function", mapperFunc.Type().String())
	}

	val := args[1]

	err := evaluator.EvaluateValueStrict(&val, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !val.IsArray() {
		return evaluator.Value{}, fmt.Errorf("(builtin flatMapArray) unexpected type of arg 1: %s, expected array", mapperFunc.Type().String())
	}

	inputArr := val.Array(ctx)

	// Create the array once and mutate it to reduce object on the heap
	mapperFuncInput := []evaluator.Value{{}}

	res := make([]evaluator.Value, 0, len(inputArr))
	for _, v := range inputArr {
		mapperFuncInput[0] = v
		out, err := mapperFunc.Function(ctx)(mapperFuncInput, ctx)
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

// func builtin_mod(args []Value, scope Scope, ctx Context) (Value, error) {

// 	if len(args) != 2 {
// 		return Value{}, fmt.Errorf("unexpected amount of arguments passed to builtin_mod: %d, expected 2", len(args))
// 	}

// 	str := args[0]
// 	if !str.IsString() {
// 		return Value{}, fmt.Errorf("unexpected type of builtin_mod arg 0: %s, expected string", str.Type.String())
// 	}

// 	input := args[1]

// 	var res string

// 	switch input.Type {
// 	case ValueTypeString:
// 		res = fmt.Sprintf(str.String(), input.String())
// 	case ValueTypeArray:
// 		sprintfArgs := make([]any, len(input.Array()))
// 		for i, v := range input.Array() {
// 			raw, err := getValueRaw(v, ctx)
// 			if err != nil {
// 				return Value{}, err
// 			}
// 			sprintfArgs[i] = raw
// 		}
// 		res = fmt.Sprintf(str.String(), sprintfArgs...)
// 	// case ValueTypeObject:

// 	// 	objRaw, err := getObjectRaw(input.Object(), ctx)
// 	// 	if err != nil {
// 	// 		return Value{}, err
// 	// 	}

// 	// 	fmt.Sprintf()

// 	// 	res = "asdf"
// 	default:
// 		return Value{}, fmt.Errorf("unexpected argument type to string format '%s', expected string/array/object", input.Type.String())
// 	}

// 	return Value{Type: ValueTypeString, StrVal: res}, nil
// }
