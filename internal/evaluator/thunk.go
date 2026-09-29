package evaluator

import (
	"errors"
	"math"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

type ThunkType uint8

const (
	ThunkTypeNone ThunkType = iota
	// Literals should never be thunked
	// ThunkTypeNull
	// ThunkTypeString
	// ThunkTypeNumber
	// ThunkTypeBool

	// ThunkTypeSelf

	ThunkTypeObject
	ThunkTypeArray
	ThunkTypeLocal
	ThunkTypeApply
	ThunkTypeIndex
	ThunkTypeVar
	ThunkTypeFunction
	ThunkTypeConditional
	ThunkTypeBinary
	ThunkTypeUnary
	ThunkTypeImport
	ThunkTypeImportStr
	ThunkTypeImportBin
	ThunkTypeSuperIndex
	ThunkTypeInSuper
	ThunkTypeError

	ThunkTypeGoCallback
)

const (
	ThunkEvaled  = math.MaxUint32
	ThunkEvaling = math.MaxUint32 - 1
)

type ThunkNodePtr uint64

type Thunk struct {
	// NodePtr   ThunkNodePtr
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
	ctx.SuperOffset = uint32(t.EvalState)
	t.EvalState = ThunkEvaling

	value, err = evaluateNode(unsafe.Pointer(t.Bp), uint32(t.Bp), t.ScopePtr, ctx)
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
