package evaluator

import (
	"fmt"
	"unsafe"

	"github.com/google/go-jsonnet/ast"
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
	ThunkTypeSuperIndex
	ThunkTypeInSuper
	ThunkTypeError

	ThunkTypeGoCallback
)

type ThunkNodePtr uint64

type Thunk struct {
	NodePtr             ThunkNodePtr
	CapturedSelf        Value
	CapturedSuperOffset uint32
	ScopeId             uint32

	Value Value
}

func NewThunk(nodeType ThunkType, nodePtr unsafe.Pointer, scopeId uint32, ctx Context) Value {
	t := Thunk{
		NodePtr:             boxThunkNodePtr(nodeType, nodePtr),
		ScopeId:             scopeId,
		CapturedSelf:        ctx.Self,
		CapturedSuperOffset: ctx.SuperOffset,
	}
	return MakeThunk(t, ctx)
}

func boxThunkNodePtr(nodeType ThunkType, nodePtr unsafe.Pointer) ThunkNodePtr {
	return ThunkNodePtr((uint64(nodeType) << 56) | uint64(uintptr(nodePtr)))
}

func (v ThunkNodePtr) unbox() (t ThunkType, p unsafe.Pointer) {
	const typeMask = 0xFF00000000000000

	t = ThunkType(v >> 56)
	p = unsafe.Pointer(uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&v))) &^ typeMask)

	// Note: The row above is equal to "p = unsafe.Pointer(uintptr(v &^ typeMask))",
	// 	they result the same assembly code. Its done this way to not make "go vet"
	// 	flag it as "possible misuse of unsafe.Pointer".
	return
}

//go:noinline
func (t *Thunk) Eval(baseCtx Context) (value Value, err error) {

	nodeType, nodePtr := t.NodePtr.unbox()
	scopeId := t.ScopeId

	ctx := baseCtx
	ctx.Self = t.CapturedSelf
	ctx.SuperOffset = t.CapturedSuperOffset

	switch nodeType {
	default:
		return ValueNone, fmt.Errorf("unhandled thunk node type (%d)", nodeType)

	case ThunkTypeObject:
		node := (*ast.DesugaredObject)(nodePtr)
		value, err = handleDesugaredObject(node, scopeId, ctx)
	case ThunkTypeArray:
		node := (*ast.Array)(nodePtr)
		value, err = handleArray(node, scopeId, ctx)
	case ThunkTypeLocal:
		node := (*ast.Local)(nodePtr)
		value, err = handleLocal(node, scopeId, ctx)
	case ThunkTypeApply:
		node := (*ast.Apply)(nodePtr)
		value, err = handleApply(node, scopeId, ctx)
	case ThunkTypeIndex:
		node := (*ast.Index)(nodePtr)
		value, err = handleIndex(node, scopeId, ctx)
	case ThunkTypeVar:
		node := (*ast.Var)(nodePtr)
		value, err = handleVar(node, scopeId, ctx)
	case ThunkTypeFunction:
		node := (*ast.Function)(nodePtr)
		value, err = handleFunction(node, scopeId, ctx)
	case ThunkTypeConditional:
		node := (*ast.Conditional)(nodePtr)
		value, err = handleConditional(node, scopeId, ctx)
	case ThunkTypeBinary:
		node := (*ast.Binary)(nodePtr)
		value, err = handleBinary(node, scopeId, ctx)
	case ThunkTypeUnary:
		node := (*ast.Unary)(nodePtr)
		value, err = handleUnary(node, scopeId, ctx)
	case ThunkTypeImport:
		node := (*ast.Import)(nodePtr)
		value, err = handleImport(node, ctx)
	case ThunkTypeImportStr:
		node := (*ast.ImportStr)(nodePtr)
		value, err = handleImportStr(node, ctx)
	case ThunkTypeSuperIndex:
		node := (*ast.SuperIndex)(nodePtr)
		value, err = handleSuperIndex(node, scopeId, ctx)
	case ThunkTypeInSuper:
		node := (*ast.InSuper)(nodePtr)
		value, err = handleInSuper(node, scopeId, ctx)
	case ThunkTypeError:
		node := (*ast.Error)(nodePtr)
		value, err = handleError(node, scopeId, ctx)

	case ThunkTypeGoCallback:
		node := (*GoCallbackNode)(nodePtr)
		value, err = node.Func.Exec(node.Args, ctx)
	}

	if err != nil {
		return ValueNone, err
	}

	value, err = value.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	t.Value = value
	return
}
