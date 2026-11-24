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

	layer, dl := evaluator.NewDynamicLayer(allocator)
	dl.Keys = arena.Alloc[uint32](allocator, n)
	dl.Nodes = arena.AllocZeroed[*evaluator.IRNode](allocator, n)
	dl.Meta = arena.Alloc[uint8](allocator, n)

	scopes := arena.Alloc[uintptr](allocator, n)
	layer.ParentScopePtr = uintptr(unsafe.Pointer(unsafe.SliceData(scopes))) | 1

	var dt *utils.DescriptorTable
	if n > evaluator.MaxLayerLinearKeys {
		dt = utils.NewEmptyDescriptorTable(allocator, n)
		dl.Index = dt
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

		innerLayer := evaluator.LayerAt(l, 0)
		innerKeys := innerLayer.Keys()
		if len(innerKeys) == 0 {
			continue
		}
		if len(innerKeys) != 1 {
			return evaluator.ValueNone, evaluator.MakeRuntimeError(errors.New("Object comprehension can only have one field"))
		}

		key := innerKeys[0]
		if dt != nil {
			if dt.Append(key) == math.MaxUint32 {
				return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
			}
		} else {
			for j := 0; j < index; j++ {
				if dl.Keys[j] == key {
					return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
				}
			}
		}

		meta := innerLayer.Meta()[0]
		dl.Keys[index] = key
		dl.Nodes[index] = innerLayer.Nodes()[0]
		dl.Meta[index] = meta
		if ast.ObjectFieldHide(meta&evaluator.MaskVisibility) == ast.ObjectFieldHidden {
			dl.HasHidden = true
		}

		scopes[index] = innerLayer.ParentScopePtr

		index++
	}

	if index < n {
		dl.Keys = dl.Keys[:index]
		dl.Nodes = dl.Nodes[:index]
		dl.Meta = dl.Meta[:index]
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

	n := len(inputArr)
	var stackSub [32]evaluator.Value
	var subArrayValues []evaluator.Value
	if n <= len(stackSub) {
		subArrayValues = stackSub[:n]
	} else {
		subArrayValues = arena.Alloc[evaluator.Value](ctx.State.Allocator, n)
	}
	totalLen := 0

	// Fast path: 1-parameter user function (all desugared array & object comprehensions)
	if mapperFunc.Type() == evaluator.ValueTypeFunction {
		f := mapperFunc.Function()
		ir := f.IR
		if len(ir.ParamKeyIds) == 1 {
			if ctx.Depth >= ctx.State.MaxStack {
				return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("max stack frames exceeded."))
			}
			callCtx := ctx
			callCtx.Depth++
			callCtx.Self = f.Self
			callCtx.SuperOffset = f.SuperOffset
			body := ir.Body
			parentScope := f.ScopePtr

			for i, v := range inputArr {
				slots, childScopeId := callCtx.NewScope(parentScope, 1)
				slots[0] = v
				out, err := evaluator.EvaluateNode(body, childScopeId, callCtx)
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
	}

	mapperFuncInput := arena.Alloc[evaluator.NamedValue](ctx.State.Allocator, 1)
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
