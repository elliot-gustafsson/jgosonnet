package evaluator

import (
	"errors"
	"math"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

const (
	ThunkEvaled  = math.MaxUint32
	ThunkEvaling = math.MaxUint32 - 1

	ThunkOffsetGoCallback = math.MaxUint32
)

type GoCallbackNode struct {
	FuncVal Value
	Args    []NamedValue
}

type Thunk struct {
	Bp        uintptr
	ScopePtr  uintptr
	Value     Value
	Offset    uint32
	EvalState uint32
}

func NewThunk(nodePtr unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) Value {
	t := arena.Create[Thunk](ctx.State.Allocator)
	t.Bp = uintptr(nodePtr)
	t.Offset = offset
	t.ScopePtr = scopePtr
	t.Value = ctx.Self
	t.EvalState = ctx.SuperOffset
	return MakeThunkValue(t, ctx)
}

func NewCallbackThunk(cb *GoCallbackNode, ctx Context) Value {
	t := arena.Create[Thunk](ctx.State.Allocator)
	t.Bp = uintptr(unsafe.Pointer(cb))
	t.Offset = ThunkOffsetGoCallback
	t.ScopePtr = 0
	t.Value = ctx.Self
	t.EvalState = ctx.SuperOffset
	return MakeThunkValue(t, ctx)
}

//go:noinline
func (t *Thunk) Eval(baseCtx Context) (value Value, err error) {
	if t.EvalState == ThunkEvaled {
		return t.Value, nil
	}
	if t.EvalState == ThunkEvaling {
		return ValueNone, MakeRuntimeError(errors.New("infinite loop detected"))
	}

	ctx := baseCtx
	ctx.Self = t.Value
	ctx.SuperOffset = t.EvalState
	t.EvalState = ThunkEvaling

	if t.Offset == ThunkOffsetGoCallback {
		cb := (*GoCallbackNode)(resolveUintptr(t.Bp))
		value, err = cb.FuncVal.FunctionExec(cb.Args, ctx)
	} else {
		value, err = evaluateNode(resolveUintptr(t.Bp), t.Offset, t.ScopePtr, ctx)
	}
	if err != nil {
		return ValueNone, err
	}

	value, err = value.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	t.Value = value
	t.EvalState = ThunkEvaled
	return
}
