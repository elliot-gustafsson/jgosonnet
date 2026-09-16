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

func execFunction(funcVal Value, args []NamedValue, ctx Context, tailstrict bool) (Value, error) {
	t := uint64(funcVal) >> typeShift

	if t == uint64(ValueTypeNativeFunction) {
		return execNativeFunction(funcVal, args, ctx, tailstrict)
	}

	if t == uint64(ValueTypeFunction) {
		return execUserFunction(funcVal, args, ctx, tailstrict)
	}

	return ValueNone, TypeErrorSpecific(ValueTypeFunction, funcVal.Type())
}

//go:noinline
func execUserFunction(funcVal Value, args []NamedValue, callCtx Context, tailstrict bool) (Value, error) {
	f := funcVal.Function()

	paramCount := len(f.Node.Parameters)
	node := f.Node
	scopePtr := f.ScopePtr
	paramKeyIds := f.ParamKeyIds

	ctx := callCtx
	if !tailstrict {
		if callCtx.State.MaxStack > 0 && callCtx.Depth >= callCtx.State.MaxStack {
			return ValueNone, MakeRuntimeError(fmt.Errorf("max stack frames exceeded."))
		}
		ctx.Depth++
	}
	ctx.Self = f.Self
	ctx.SuperOffset = f.SuperOffset

	if paramCount == 0 && len(args) == 0 {
		return EvaluateNode(node.Body, scopePtr, ctx)
	}

	if paramCount == len(args) && (paramCount == 0 || args[paramCount-1].Key == 0) {
		if paramCount == 0 {
			return EvaluateNode(node.Body, scopePtr, ctx)
		}
		s, childScopeId := ctx.NewScope(scopePtr, paramCount)
		for i := range args {
			a := args[i]
			if tailstrict {
				var err error
				a.Value, err = a.Eval(ctx)
				if err != nil {
					return ValueNone, err
				}
			}
			a.Key = paramKeyIds[i]
			s.Bindings[i] = a
		}
		return EvaluateNode(node.Body, childScopeId, ctx)
	}

	return execUserFunctionSlow(f, args, ctx, tailstrict)
}

//go:noinline
func execUserFunctionSlow(f *Function, args []NamedValue, ctx Context, tailstrict bool) (Value, error) {
	paramCount := len(f.Node.Parameters)
	node := f.Node
	scopePtr := f.ScopePtr
	paramKeyIds := f.ParamKeyIds

	s, childScopeId := ctx.NewScope(scopePtr, paramCount)

	var onNamedArgs bool
	var posIdx int
	var bound int

OuterLoop:
	for _, inputArg := range args {

		if tailstrict {
			v, err := inputArg.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
			inputArg.Value = v
		}

		if inputArg.Key == 0 {
			// Positional Argument
			if onNamedArgs {
				return ValueNone, fmt.Errorf("Positional argument after a named argument is not allowed")
			}
			if posIdx+1 > paramCount {
				numPos := 0
				for _, a := range args {
					if a.Key == 0 {
						numPos++
					}
				}
				return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", paramCount, numPos))
			}

			bindVal := inputArg
			bindVal.Key = paramKeyIds[posIdx]
			s.Bindings[posIdx] = bindVal

			posIdx++
			bound++
			continue
		}
		onNamedArgs = true

		for j := posIdx; j < paramCount; j++ {
			if paramKeyIds[j] == inputArg.Key && s.Bindings[j].IsNone() {
				s.Bindings[j] = inputArg
				bound++
				continue OuterLoop
			}
		}

		// Should only end up here in case of error
		passedName := ctx.State.Interner.Get(inputArg.Key)
		for _, bind := range s.Bindings {
			if bind.Key == inputArg.Key {
				return ValueNone, MakeRuntimeError(fmt.Errorf("Argument %s already provided", passedName))
			}
		}
		return ValueNone, MakeRuntimeError(fmt.Errorf("function has no parameter %s", passedName))

	}

	missingCount := paramCount - bound
	if missingCount == 0 {
		return EvaluateNode(node.Body, childScopeId, ctx)
	}

	for i := posIdx; i < paramCount; i++ {
		binding := s.Bindings[i]
		if !binding.IsNone() {
			continue
		}

		n := node.Parameters[i].DefaultArg
		if n != nil {
			da, err := evaluateNodeLazy(n, childScopeId, ctx)
			if err != nil {
				return ValueNone, err
			}

			if tailstrict {
				da, err = da.Eval(ctx)
				if err != nil {
					return ValueNone, err
				}
			}

			s.Bindings[i] = NamedValue{paramKeyIds[i], da}
			missingCount--
			if missingCount == 0 {
				break // all bound
			}
			continue
		}

		missingName := ctx.State.Interner.Get(paramKeyIds[i])
		return ValueNone, MakeRuntimeError(fmt.Errorf("Missing argument: %s", missingName))
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
func execNativeFunction(funcVal Value, args []NamedValue, ctx Context, tailstrict bool) (Value, error) {
	f := funcVal.NativeFunction()
	paramCount := len(f.Params)
	argsCount := len(args)
	fn := f.Func

	if argsCount == paramCount {
		if argsCount == 0 || args[argsCount-1].Key == 0 {
			// If all params are passed and all is positional
			if tailstrict {
				for i := range args {
					v, err := args[i].Eval(ctx)
					if err != nil {
						return ValueNone, err
					}
					args[i].Value = v
				}
			}
			return fn(args, ctx)
		}
	}

	params := f.Params
	optStart := int(f.OptStart)

	orderedArgs := arena.Alloc[NamedValue](ctx.State.Allocator, paramCount)
	clear(orderedArgs)

	var onNamedArgs bool
	posIdx := 0

	for _, na := range args {

		if tailstrict {
			v, err := na.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
			na.Value = v
		}

		if na.Key == 0 {
			// Positional Argument
			if onNamedArgs {
				return ValueNone, fmt.Errorf("Positional argument after a named argument is not allowed")
			}
			if posIdx+1 > paramCount {
				numPos := 0
				for _, a := range args {
					if a.Key == 0 {
						numPos++
					}
				}
				return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", paramCount, numPos))
			}

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
				if !orderedArgs[j].IsNone() {
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

	for i := range optStart {
		if orderedArgs[i].Value.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Missing argument: %s", params[i]))
		}
	}

	args = orderedArgs

	return fn(args, ctx)
}
