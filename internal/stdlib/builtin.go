package stdlib

import (
	"errors"
	"fmt"
	"math"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"github.com/google/go-jsonnet/ast"
)

func builtin_objectFlatMerge(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	n := len(inputArr)
	allocator := ctx.State.Allocator

	layer := arena.Create[evaluator.Layer](allocator)
	arena.Memclr(layer)
	layer.Keys = arena.Alloc[uint32](allocator, n)
	layer.Nodes = arena.Alloc[ast.Node](allocator, n)
	arena.MemclrSlice(layer.Nodes)
	layer.Meta = arena.Alloc[uint8](allocator, n)

	scopes := arena.Alloc[uintptr](allocator, n)
	layer.ParentScopePtr = uintptr(unsafe.Pointer(unsafe.SliceData(scopes))) | 1

	var dt *utils.DescriptorTable
	if n > evaluator.MaxLayerLinearKeys {
		dt = utils.NewEmptyDescriptorTable(allocator, n)
		layer.Index = dt
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

		innerLayer := l[0]
		if len(innerLayer.Keys) == 0 {
			continue
		}
		if len(innerLayer.Keys) != 1 {
			return evaluator.ValueNone, evaluator.MakeRuntimeError(errors.New("Object comprehension can only have one field"))
		}

		key := innerLayer.Keys[0]
		if dt != nil {
			if dt.Append(key) == math.MaxUint32 {
				return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
			}
		} else {
			for j := 0; j < index; j++ {
				if layer.Keys[j] == key {
					return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
				}
			}
		}

		layer.Keys[index] = key
		layer.Nodes[index] = innerLayer.Nodes[0]
		layer.Meta[index] = innerLayer.Meta[0]

		scopes[index] = innerLayer.ParentScopePtr

		index++
	}

	if index < n {
		layer.Keys = layer.Keys[:index]
		layer.Nodes = layer.Nodes[:index]
		layer.Meta = layer.Meta[:index]
	}

	obj := evaluator.NewSingleLayerObject(allocator, layer)
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
