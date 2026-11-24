package evaluator

import (
	"errors"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

type ThunkNodePtr uint64

type Thunk struct {
	NodePtr   ThunkNodePtr
	ScopePtr  uintptr
	Value     Value
	EvalState int
}

func NewThunk(nodeType IRKind, nodePtr unsafe.Pointer, scopePtr uintptr, ctx Context) Value {
	t := arena.Create[Thunk](ctx.State.Allocator)
	t.NodePtr = boxThunkNodePtr(nodeType, nodePtr)
	t.ScopePtr = scopePtr
	t.Value = ctx.Self
	t.EvalState = int(ctx.SuperOffset)
	return MakeThunkValue(t, ctx)
}

func boxThunkNodePtr(nodeType IRKind, nodePtr unsafe.Pointer) ThunkNodePtr {
	return ThunkNodePtr((uint64(nodeType) << 56) | uint64(uintptr(nodePtr)))
}

func (v ThunkNodePtr) unbox() (t IRKind, p unsafe.Pointer) {
	t = IRKind(v >> 56)
	p = unsafe.Pointer(uintptr(v &^ tagMask))
	return
}

//go:noinline
func (t *Thunk) Eval(baseCtx Context) (value Value, err error) {
	if t.EvalState == -1 {
		return t.Value, nil
	}
	if t.EvalState == -2 {
		return ValueNone, MakeRuntimeError(errors.New("infinite loop detected"))
	}

	nodeType, nodePtr := t.NodePtr.unbox()
	scopeId := t.ScopePtr

	ctx := baseCtx
	ctx.Self = t.Value
	ctx.SuperOffset = uint32(t.EvalState)
	t.EvalState = -2

	if nodeType == IRKindGoCallback {
		node := (*GoCallbackNode)(nodePtr)
		value, err = node.FuncVal.FunctionExec(node.Args, ctx)
	} else {
		value, err = evaluateIR((*IRNode)(nodePtr), scopeId, ctx)
	}

	if err != nil {
		return ValueNone, err
	}

	value, err = value.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	t.Value = value
	t.EvalState = -1
	return
}
