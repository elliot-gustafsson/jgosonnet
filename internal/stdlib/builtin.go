package stdlib

import (
	"errors"
	"fmt"
	"math"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	n := len(inputArr)
	allocator := ctx.State.Allocator

	layers := arena.Alloc[*evaluator.Layer](allocator, n)
	arena.MemclrSlice(layers)

	var dt *utils.DescriptorTable
	if n > evaluator.MaxLayerLinearKeys {
		dt = utils.NewEmptyDescriptorTable(allocator, n)
	}

	index := 0
	for _, v := range inputArr {
		obj, err := v.EvalObject(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		l := obj.GetLayers(ctx)
		if len(l) == 0 {
			continue
		}
		if len(l) != 1 {
			return evaluator.ValueNone, evaluator.MakeRuntimeError(errors.New("Object comprehension can only have one layer"))
		}

		if len(l[0].Keys) == 0 {
			continue
		}
		if len(l[0].Keys) != 1 {
			return evaluator.ValueNone, evaluator.MakeRuntimeError(errors.New("Object comprehension can only have one field"))
		}

		key := l[0].Keys[0]

		if dt != nil {
			if dt.Append(key) == math.MaxUint32 {
				return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
			}
		} else {
			for j := 0; j < index; j++ {
				if layers[j].Keys[0] == key {
					return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
				}
			}
		}

		layers[index] = l[0]

		index++
	}

	if index < n {
		layers = layers[:index]
	}

	obj := arena.Create[evaluator.Object](allocator)
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
