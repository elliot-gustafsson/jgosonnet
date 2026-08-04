package evaluator

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/google/go-jsonnet/ast"
)

type Func = func(args []NamedValue, ctx Context) (Value, error)

type Function struct {
	Node     *ast.Function
	ScopePtr uintptr

	ParamKeyIds []uint32

	Self        Value
	SuperOffset uint32
}

type NativeFunction struct {
	Params   []string
	Func     Func
	OptStart uint32
}

// func (t *Function) Noop() bool {
// 	return t.Node == nil
// }

// Length of a function is the number of required args
func (t *Function) Length() int {
	params := t.Node.Parameters

	var n int
	for _, p := range params {
		if p.DefaultArg == nil {
			n++
		}
	}
	return n
}

func execFunction(funcVal Value, args []NamedValue, ctx Context) (Value, error) {
	t := uint64(funcVal) >> typeShift

	if t == uint64(ValueTypeNativeFunction) {
		return execNativeFunction(funcVal, args, ctx)
	}

	if t == uint64(ValueTypeFunction) {
		return execUserFunction(funcVal, args, ctx)
	}

	return ValueNone, TypeErrorSpecific(ValueTypeFunction, funcVal.Type())
}

//go:noinline
func execUserFunction(funcVal Value, args []NamedValue, callCtx Context) (Value, error) {
	f := funcVal.Function()

	paramCount := len(f.Node.Parameters)
	node := f.Node
	scopePtr := f.ScopePtr
	paramKeyIds := f.ParamKeyIds

	ctx := callCtx
	ctx.Self = f.Self
	ctx.SuperOffset = f.SuperOffset

	if len(args) > paramCount {
		return ValueNone, fmt.Errorf("unexpected amount of args passed to function")
	}

	if paramCount == 0 {
		return EvaluateNode(node.Body, scopePtr, ctx)
	}

	s, childScopeId := ctx.NewScope(scopePtr, paramCount)

	// TODO: Throw err on argument x already provided

	posIndex := 0
	for i := range paramCount {
		keyId := paramKeyIds[i]

		var bindVal NamedValue

		if posIndex < len(args) && args[posIndex].Key == 0 {
			bindVal = args[posIndex]
			bindVal.Key = keyId
			posIndex++
		} else {
			for _, a := range args {
				if a.Key == keyId {
					bindVal = a
					break
				}
			}
		}

		if !bindVal.IsNone() {
			s.Bindings[i] = bindVal
			continue
		}

		// No arg was passed, fallback to default arg
		defArgNode := node.Parameters[i].DefaultArg
		if defArgNode == nil {
			return ValueNone, fmt.Errorf("arg (%d) with no default arg had no value passed", i)
		}

		da, err := evaluateNodeLazy(defArgNode, childScopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		s.Bindings[i] = NamedValue{keyId, da}
	}

	return EvaluateNode(node.Body, childScopeId, ctx)
}

// func (t NativeFunction) Noop() bool {
// 	return t.fn == nil
// }

func (t *NativeFunction) Length() int {
	return int(t.OptStart)
}

//go:noinline
func execNativeFunction(funcVal Value, args []NamedValue, ctx Context) (Value, error) {
	f := funcVal.NativeFunction()
	paramCount := len(f.Params)
	fn := f.Func
	params := f.Params
	optStart := int(f.OptStart)

	if len(args) > paramCount {
		return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", paramCount, len(args)))
	}

	var onNamedArgs bool

	// Use stack arr and overrite args to avoid allocating a slice for all stdlib funcs
	var stackArgs [4]NamedValue // Max args for any stdlib function is 4
	var orderedArgs []NamedValue

	if paramCount <= len(stackArgs) {
		orderedArgs = stackArgs[:paramCount]
	} else {
		// just in case a user defines a custom native extension with 5+ args
		orderedArgs = arena.Alloc[NamedValue](ctx.State.Allocator, paramCount)
		clear(orderedArgs)
	}

	posIdx := 0

	for _, na := range args {
		if na.Key == 0 {
			// Positional Argument
			if onNamedArgs {
				return ValueNone, fmt.Errorf("Positional argument after a named argument is not allowed")
			}
			// na.Key = argIds[posIdx]
			orderedArgs[posIdx] = na
			posIdx++
			continue
		}

		// Named Argument
		passedName := ctx.State.Interner.Get(na.Key)

		onNamedArgs = true
		found := false
		for j, paramName := range params {
			if passedName == paramName {
				if !orderedArgs[j].Value.IsNone() {
					argName := ctx.State.Interner.Get(na.Key)
					return ValueNone, MakeRuntimeError(fmt.Errorf("Argument %s already provided", argName))
				}
				orderedArgs[j] = na
				found = true
				break
			}
		}
		if !found {
			return ValueNone, MakeRuntimeError(fmt.Errorf("function has no parameter %s", passedName))
		}
	}

	for i := 0; i < optStart; i++ {
		if orderedArgs[i].Value.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Missing argument: %s", params[i]))
		}
	}

	args = orderedArgs

	return fn(args, ctx)
}
