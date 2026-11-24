package evaluator

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

type Func = func(args []NamedValue, ctx Context) (Value, error)

type Function struct {
	IR          *IRFunction
	ScopePtr    uintptr
	Self        Value
	SuperOffset uint32
}

type NativeFunction struct {
	Params   []string
	Func     Func
	OptStart uint32
}

// Length of a function is the number of required args
func (t *Function) Length() int {
	return t.IR.RequiredCnt
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
	ir := f.IR
	paramKeyIds := ir.ParamKeyIds
	paramCount := len(paramKeyIds)
	scopePtr := f.ScopePtr

	ctx := callCtx
	if !tailstrict {
		if callCtx.Depth >= callCtx.State.MaxStack {
			return errMaxStackExceeded()
		}
		ctx.Depth++
	}
	ctx.Self = f.Self
	ctx.SuperOffset = f.SuperOffset

	argsLen := len(args)

	if paramCount == 0 {
		if argsLen == 0 {
			return EvaluateNode(ir.Body, scopePtr, ctx)
		}
		if args[0].Key == 0 {
			return errTooManyPositional(0, countPositional(args))
		}
		return errNoSuchParam(args[0].Key, ctx)
	}

	slots, childScopeId := ctx.NewScope(scopePtr, paramCount)

	// bind positional args
	posIdx := 0
	for posIdx < argsLen && args[posIdx].Key == 0 {
		if posIdx >= paramCount {
			numPos := argsLen
			for i := posIdx; i < argsLen; i++ {
				if args[i].Key != 0 {
					numPos = i
					break
				}
			}
			return errTooManyPositional(paramCount, numPos)
		}
		v := args[posIdx].Value
		if tailstrict {
			var err error
			v, err = v.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
		}
		slots[posIdx] = v
		posIdx++
	}

	if posIdx == argsLen && posIdx == paramCount {
		return EvaluateNode(ir.Body, childScopeId, ctx)
	}

	// Zero unbound trailing slots before binding named/default args
	clear(slots[posIdx:])

	// bind named args
	if posIdx < argsLen {
	NamedArgsOuterLoop:
		for i := posIdx; i < argsLen; i++ {
			na := args[i]
			if na.Key == 0 {
				return errPosArgAfterNamed()
			}

			val := na.Value
			if tailstrict {
				var err error
				val, err = val.Eval(ctx)
				if err != nil {
					return ValueNone, err
				}
			}

			for j := range paramCount {
				if paramKeyIds[j] == na.Key {
					if !slots[j].IsNone() {
						return errArgAlreadyProvided(na.Key, ctx)
					}
					slots[j] = val
					continue NamedArgsOuterLoop
				}
			}

			return errNoSuchParam(na.Key, ctx)
		}
	}

	// fill default args
	defaultArgs := ir.DefaultArgs
	for i := posIdx; i < paramCount; i++ {
		if !slots[i].IsNone() {
			continue
		}

		if defaultArgs == nil || defaultArgs[i] == nil {
			return errMissingArg(paramKeyIds[i], ctx)
		}

		da := evaluateIRLazy(defaultArgs[i], childScopeId, ctx)
		if tailstrict {
			var err error
			da, err = da.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
		}

		slots[i] = da
	}

	return EvaluateNode(ir.Body, childScopeId, ctx)
}

func (t *NativeFunction) Length() int {
	return int(t.OptStart)
}

//go:noinline
func execNativeFunction(funcVal Value, args []NamedValue, ctx Context, tailstrict bool) (Value, error) {
	f := funcVal.NativeFunction()
	paramCount := len(f.Params)
	params := f.Params
	argsLen := len(args)
	optStart := int(f.OptStart)
	fn := f.Func

	if paramCount == 0 {
		if argsLen == 0 {
			return fn(args, ctx)
		}
		if args[0].Key == 0 {
			return errTooManyPositional(0, countPositional(args))
		}
		return errNoSuchParam(args[0].Key, ctx)
	}

	// bind positional args
	posIdx := 0
	for posIdx < argsLen && args[posIdx].Key == 0 {
		if posIdx >= paramCount {
			return errTooManyPositional(paramCount, countPositional(args))
		}
		if tailstrict {
			v, err := args[posIdx].Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
			args[posIdx].Value = v
		}
		posIdx++
	}

	if posIdx == argsLen && posIdx == paramCount {
		return fn(args, ctx)
	}

	// alloc temp slice only when named arguments are used or defaults needed
	orderedArgs := arena.Alloc[NamedValue](ctx.State.Allocator, paramCount)
	clear(orderedArgs)

	// copy verified positional arguments directly
	copy(orderedArgs[:posIdx], args[:posIdx])

	requiredMissing := optStart - posIdx

	// bind named arguments
	if posIdx < argsLen {
	NamedArgsOuterLoop:
		for i := posIdx; i < argsLen; i++ {
			na := args[i]
			if na.Key == 0 {
				return errPosArgAfterNamed()
			}

			if tailstrict {
				v, err := na.Eval(ctx)
				if err != nil {
					return ValueNone, err
				}
				na.Value = v
			}

			passedName := ctx.State.Interner.Get(na.Key)
			for j, paramName := range params {
				if passedName == paramName {
					if !orderedArgs[j].IsNone() {
						return errArgAlreadyProvided(na.Key, ctx)
					}
					orderedArgs[j] = na
					if j < optStart {
						requiredMissing--
					}
					continue NamedArgsOuterLoop
				}
			}

			return errNoSuchParam(na.Key, ctx)
		}
	}

	if requiredMissing > 0 {
		for i := posIdx; i < optStart; i++ {
			if orderedArgs[i].Value.IsNone() {
				return errMissingArgName(params[i])
			}
		}
	}

	return fn(orderedArgs, ctx)
}

//go:noinline
func errMaxStackExceeded() (Value, error) {
	return ValueNone, MakeRuntimeError(fmt.Errorf("max stack frames exceeded."))
}

//go:noinline
func errPosArgAfterNamed() (Value, error) {
	return ValueNone, fmt.Errorf("Positional argument after a named argument is not allowed")
}

//go:noinline
func errTooManyPositional(expected, got int) (Value, error) {
	return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", expected, got))
}

//go:noinline
func errArgAlreadyProvided(key uint32, ctx Context) (Value, error) {
	name := ctx.State.Interner.Get(key)
	return ValueNone, MakeRuntimeError(fmt.Errorf("Argument %s already provided", name))
}

//go:noinline
func errNoSuchParam(key uint32, ctx Context) (Value, error) {
	name := ctx.State.Interner.Get(key)
	return ValueNone, MakeRuntimeError(fmt.Errorf("function has no parameter %s", name))
}

//go:noinline
func errMissingArgName(name string) (Value, error) {
	return ValueNone, MakeRuntimeError(fmt.Errorf("Missing argument: %s", name))
}

//go:noinline
func errMissingArg(key uint32, ctx Context) (Value, error) {
	name := ctx.State.Interner.Get(key)
	return errMissingArgName(name)
}

func countPositional(args []NamedValue) int {
	for i, a := range args {
		if a.Key != 0 {
			return i
		}
	}
	return len(args)
}
