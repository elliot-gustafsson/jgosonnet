package evaluator

import (
	"errors"
	"fmt"
	"io"
	"math/bits"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"github.com/google/go-jsonnet/ast"
)

func EvaluateNode(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := evaluateNode(n, scopePtr, ctx)
	if err != nil {
		return ValueNone, WrapError(err, n)
	}
	if val.IsThunk() {
		val, err = val.Eval(ctx)
		if err != nil {
			return ValueNone, WrapError(err, n)
		}
	}
	return val, nil
}

func ManifestValue(value Value, ctx Context) (any, error) {

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type '%s'", value.Type().String())
	case ValueTypeNull:
		return nil, nil
	case ValueTypeString:
		s := value.String(ctx)
		return strings.Clone(s), nil
	case ValueTypeNumber:
		n := value.Number()
		if float64(int64(n)) == n {
			return int64(n), nil
		}
		return n, nil
	case ValueTypeBool:
		return value.Bool(), nil
	case ValueTypeObject:
		subCtx := ctx
		subCtx.Self = value
		return manifestObject(value.Object(), subCtx)
	case ValueTypeArray:
		res := make([]any, 0, len(value.Array()))
		for _, v := range value.Array() {
			ev, err := ManifestValue(v, ctx)
			if err != nil {
				return nil, err
			}
			res = append(res, ev)
		}
		return res, nil
	case ValueTypeFunction, ValueTypeNativeFunction:
		res, err := value.FunctionExec(nil, ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(res, ctx)
	case ValueTypeThunk:
		v, err := value.Eval(ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(v, ctx)
	}
}

func CreateFileScope(filename string, baseStd Value, ctx Context) uintptr {
	allocator := ctx.State.Allocator

	keyId := ctx.State.Interner.Intern("thisFile")

	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)

	layer.Keys = arena.Alloc[uint32](allocator, 1)
	layer.Keys[0] = keyId

	layer.Values = arena.Alloc[Value](allocator, 1)
	layer.Values[0] = MakeString(filename, ctx)

	layer.Meta = arena.Alloc[uint8](allocator, 1)
	layer.Meta[0] = 0

	fileObj := NewSingleLayerObject(allocator, layer)

	mergedObjId := MergeObjects(baseStd.Payload(), uintptr(unsafe.Pointer(fileObj)), ctx)
	fileStd := MakeObjectValue(mergedObjId)

	s, scopePtr := ctx.NewScope(2, 2)

	s.Bindings[0] = NamedValue{Key: ctx.State.Interner.Intern("$std"), Value: fileStd}
	s.Bindings[1] = NamedValue{Key: ctx.State.Interner.Intern("std"), Value: fileStd}

	return scopePtr
}

func evaluateNodeLazy(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	switch node := n.(type) {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %T (lazy eval)", node)
	case *ast.LiteralString:
		id := ctx.State.Interner.Intern(node.Value)
		return MakeStringConst(id), nil
		// return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return ValueNone, fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.Self:
		// if !ctx.Self.IsObject() {
		// 	return ValueNone, MakeRuntimeError(fmt.Errorf("self not initialized"))
		// }
		return ctx.Self, nil

	case *ast.DesugaredObject:
		return NewThunk(ThunkTypeObject, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Array:
		return NewThunk(ThunkTypeArray, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Local:
		return NewThunk(ThunkTypeLocal, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Apply:
		return NewThunk(ThunkTypeApply, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Index:
		return NewThunk(ThunkTypeIndex, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Var:
		return NewThunk(ThunkTypeVar, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Function:
		return NewThunk(ThunkTypeFunction, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Conditional:
		return NewThunk(ThunkTypeConditional, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Binary:
		return NewThunk(ThunkTypeBinary, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Unary:
		return NewThunk(ThunkTypeUnary, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Import:
		return NewThunk(ThunkTypeImport, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.ImportStr:
		return NewThunk(ThunkTypeImportStr, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.ImportBin:
		return NewThunk(ThunkTypeImportBin, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.SuperIndex:
		return NewThunk(ThunkTypeSuperIndex, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.InSuper:
		return NewThunk(ThunkTypeInSuper, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Error:
		return NewThunk(ThunkTypeError, unsafe.Pointer(node), scopePtr, ctx), nil

	case *GoCallbackNode:
		return NewThunk(ThunkTypeGoCallback, unsafe.Pointer(node), scopePtr, ctx), nil
	}
}

func evaluateNode(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	switch node := n.(type) {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %T", node)
	case *ast.LiteralString:
		id := ctx.State.Interner.Intern(node.Value)
		return MakeStringConst(id), nil
		// return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return ValueNone, fmt.Errorf("(%T) failed to parse float val (%s), err: %w", node, node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.DesugaredObject:
		return handleDesugaredObject(node, scopePtr, ctx)
	case *ast.Array:
		return handleArray(node, scopePtr, ctx)
	case *ast.Local:
		return handleLocal(node, scopePtr, ctx)
	case *ast.Apply:
		return handleApply(node, scopePtr, ctx)
	case *ast.Index:
		return handleIndex(node, scopePtr, ctx)
	case *ast.Var:
		return handleVar(node, scopePtr, ctx)
	case *ast.Function:
		return handleFunction(node, scopePtr, ctx)
	case *ast.Conditional:
		return handleConditional(node, scopePtr, ctx)
	case *ast.Binary:
		return handleBinary(node, scopePtr, ctx)
	case *ast.Unary:
		return handleUnary(node, scopePtr, ctx)
	case *ast.Import:
		return handleImport(node, ctx)
	case *ast.ImportStr:
		return handleImportStr(node, ctx)
	case *ast.ImportBin:
		return handleImportBin(node, ctx)
	case *ast.Self:
		// if !ctx.Self.IsObject() {
		// 	return ValueNone, MakeRuntimeError(fmt.Errorf("self not initialized"))
		// }
		return ctx.Self, nil
	case *ast.SuperIndex:
		return handleSuperIndex(node, scopePtr, ctx)
	case *ast.InSuper:
		return handleInSuper(node, scopePtr, ctx)
	case *ast.Error:
		return handleError(node, scopePtr, ctx)
	case *GoCallbackNode:
		return node.FuncVal.FunctionExec(node.Args, ctx)
	}
}

type GoCallbackNode struct {
	FuncVal Value
	Args    []NamedValue
}

func (n *GoCallbackNode) Context() ast.Context             { return nil }
func (n *GoCallbackNode) Loc() *ast.LocationRange          { return nil }
func (n *GoCallbackNode) FreeVariables() ast.Identifiers   { return nil }
func (n *GoCallbackNode) SetFreeVariables(ast.Identifiers) {}
func (n *GoCallbackNode) SetContext(ast.Context)           {}
func (n *GoCallbackNode) OpenFodder() *ast.Fodder          { return nil }

func handleArray(node *ast.Array, scopePtr uintptr, ctx Context) (Value, error) {
	arr, val := MakeArraySized(len(node.Elements), ctx)
	for i := range node.Elements {
		ev, err := evaluateNodeLazy(node.Elements[i].Expr, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		arr[i] = ev
	}
	return val, nil
}

func handleDesugaredObject(node *ast.DesugaredObject, scopePtr uintptr, ctx Context) (Value, error) {

	allocator := ctx.State.Allocator

	fieldCount := len(node.Fields)
	localsCount := len(node.Locals)

	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)

	layer.ParentScopePtr = scopePtr

	if fieldCount > 0 {
		layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
		layer.Nodes = arena.Alloc[ast.Node](allocator, fieldCount)
		arena.MemclrSlice(layer.Nodes)
		layer.Meta = arena.Alloc[uint8](allocator, fieldCount)
	}

	if localsCount > 0 {
		layer.LocalKeys = arena.Alloc[uint32](allocator, localsCount)
		layer.LocalNodes = arena.Alloc[ast.Node](allocator, localsCount)
		arena.MemclrSlice(layer.LocalNodes)
	}

	if len(node.Asserts) > 0 {
		layer.packAsserts(node.Asserts)
	}

	for i, v := range node.Locals {
		layer.LocalKeys[i] = ctx.State.Interner.Intern(string(v.Variable))
		layer.LocalNodes[i] = v.Body
	}

	useMap := fieldCount > MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	index := 0
	for _, v := range node.Fields {
		var nameStr string

		// Fast path literal string, most keys are just static strings
		if ls, ok := v.Name.(*ast.LiteralString); ok {

			nameStr = ls.Value
		} else {
			var err error
			name, err := EvaluateNode(v.Name, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}

			if name.IsNull() {
				// Omitted field
				continue
			}

			if !name.IsString() {
				return ValueNone, fmt.Errorf("unexpected field name type %s, expected string", name.Type().String())
			}

			nameStr = name.String(ctx)
		}

		keyId := ctx.State.Interner.Intern(nameStr)

		layer.Keys[index] = keyId
		layer.Nodes[index] = v.Body
		layer.Meta[index] = CreateFieldMeta(v.Hide, v.PlusSuper)

		if useMap {
			layer.Index.Append(keyId)
		}
		index++
	}

	if index < fieldCount {
		layer.Keys = layer.Keys[:index]
		layer.Nodes = layer.Nodes[:index]
		layer.Meta = layer.Meta[:index]
	}

	obj := NewSingleLayerObject(allocator, layer)

	return MakeObjectValue(obj), nil
}

func handleLocal(node *ast.Local, scopePtr uintptr, ctx Context) (Value, error) {

	s, childScopePtr := ctx.NewScope(scopePtr, len(node.Binds))

	for i, v := range node.Binds {
		vname := string(v.Variable)
		keyId := ctx.State.Interner.Intern(vname)
		t, err := evaluateNodeLazy(v.Body, childScopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		s.Bindings[i] = NamedValue{keyId, t}
	}

	val, err := EvaluateNode(node.Body, childScopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	return val, nil
}

func handleApply(node *ast.Apply, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(node.Target, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}

	posCount := len(node.Arguments.Positional)
	nameCount := len(node.Arguments.Named)

	args := arena.Alloc[NamedValue](ctx.State.Allocator, posCount+nameCount)
	clear(args)
	for i, a := range node.Arguments.Positional {
		// v, err := EvaluateNodeStrict(a.Expr, scopeId, ctx)
		v, err := evaluateNodeLazy(a.Expr, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		args[i] = NamedValue{Value: v}
	}
	for i, a := range node.Arguments.Named {
		// a.Name
		// v, err := EvaluateNodeStrict(a.Arg, scopeId, ctx)
		v, err := evaluateNodeLazy(a.Arg, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		nameKeyId := ctx.State.Interner.Intern(string(a.Name))
		args[i+posCount] = NamedValue{nameKeyId, v}
	}

	res, err := val.FunctionExecEx(args, ctx, node.TailStrict)
	if err != nil {
		// Check for the trace signal
		if traceSig, ok := err.(*TraceSignal); ok {
			err := traceSig.Print(ctx.State.Environment.TraceOut, node)
			if err != nil {
				return ValueNone, err
			}
			// Resume normal execution with the return value
			return traceSig.Rest, nil
		}
		return ValueNone, err
	}
	return res, nil
}

func handleBinary(node *ast.Binary, scopePtr uintptr, ctx Context) (Value, error) {
	left, err := EvaluateNode(node.Left, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	// Check if fast exit is possible
	switch node.Op {
	case ast.BopAnd:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", left.Type().String())
		}

		if !left.Bool() {
			return MakeBool(false), nil
		}

		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", right.Type().String())
		}

		return right, nil

	case ast.BopOr:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", left.Type().String())
		}

		if left.Bool() {
			return MakeBool(true), nil
		}

		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}

		return right, nil

	default:
		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		res, err := handleBinaryOp(node.Op, left, right, ctx)
		if err != nil {
			return ValueNone, err
		}

		return res, nil
	}

}

func handleUnary(node *ast.Unary, scopePtr uintptr, ctx Context) (Value, error) {
	unary, err := EvaluateNode(node.Expr, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	switch node.Op {
	default:
		return ValueNone, fmt.Errorf("unhandled unary type: %s", node.Op.String())
	case ast.UopNot:
		if !unary.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, unary.Type())
		}
		return MakeBool(!unary.Bool()), nil
	case ast.UopMinus:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		res := -unary.Number()
		return MakeNumber(res), nil
	case ast.UopBitwiseNot:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		val32 := int64(unary.Number())
		notVal32 := ^val32
		return MakeNumber(float64(notVal32)), nil
	case ast.UopPlus:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		return unary, nil
	}
}

func handleVar(node *ast.Var, scopePtr uintptr, ctx Context) (Value, error) {
	name := string(node.Id)

	keyId := ctx.State.Interner.Intern(name)

	val, found := ctx.GetScopeBind(scopePtr, keyId)
	if !found {
		return ValueNone, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", name))
	}

	return val, nil
}

func handleFunction(node *ast.Function, scopePtr uintptr, ctx Context) (Value, error) {

	paramKeyIds := arena.Alloc[uint32](ctx.State.Allocator, len(node.Parameters))
	for i, p := range node.Parameters {
		paramKeyIds[i] = ctx.State.Interner.Intern(string(p.Name))
	}

	f := arena.Create[Function](ctx.State.Allocator)
	arena.Memclr(f)
	*f = Function{
		Node:        node,
		ScopePtr:    scopePtr,
		ParamKeyIds: paramKeyIds,
		Self:        ctx.Self,
		SuperOffset: ctx.SuperOffset,
	}

	return MakeFunctionValue(f, ctx), nil
}

func handleConditional(node *ast.Conditional, scopePtr uintptr, ctx Context) (Value, error) {
	cond, err := EvaluateNode(node.Cond, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !cond.IsBool() {
		return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
	}

	if cond.Bool() {
		bt, err := EvaluateNode(node.BranchTrue, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		return bt, nil
	}

	bf, err := EvaluateNode(node.BranchFalse, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	return bf, nil
}

func handleIndex(node *ast.Index, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	target, err := EvaluateNode(node.Target, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	switch target.Type() {
	default:
		return ValueNone, fmt.Errorf("value not indexable: %s", target.Type().String())
	case ValueTypeString:
		if !index.IsNumber() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing string, expected number, got %s", index.Type().String()))
		}
		i := int(index.Number())
		s := target.String(ctx)

		runeBytes, err := indexStringRune(s, i)
		if err != nil {
			return ValueNone, err
		}
		return MakeString(runeBytes, ctx), nil
	case ValueTypeObject:
		if !index.IsString() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing object, expected string, got %s", index.Type().String()))
		}

		keyId := index.AsStringConst()
		if keyId == 0 {
			name := index.String(ctx)
			keyId = ctx.State.Interner.Intern(name)
		}

		obj := target.Object()

		// Reset self to point to correct obj
		subCtx := ctx
		subCtx.Self = target

		val, _, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return ValueNone, err
		}
		if val.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", index.String(ctx)))
		}
		val, err = val.Eval(subCtx)
		if err != nil {
			return ValueNone, err
		}
		return val, nil
	case ValueTypeArray:
		if !index.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, index.Type())
		}
		i := int(index.Number())

		arr := target.Array()
		if i < 0 || len(arr) <= i {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", i, len(arr)))
		}
		return target.Array()[i], nil
	}

}

const (
	broadcastByte           = 0x0101010101010101
	swarNonAsciiMask uint64 = utf8.RuneSelf * broadcastByte
)

func indexStringRune(s string, targetRuneIndex int) (string, error) {
	if targetRuneIndex < 0 || len(s) == 0 {
		return "", MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", targetRuneIndex, utf8.RuneCountInString(s)))
	}
	n := len(s)
	byteOffset := 0
	currRuneIndex := 0
	basePtr := unsafe.Pointer(unsafe.StringData(s))

	for byteOffset+8 <= n {
		w := *(*uint64)(unsafe.Add(basePtr, byteOffset))

		// check if w includes any runes
		if (w & swarNonAsciiMask) == 0 {
			// all ascii bytes
			if currRuneIndex+8 <= targetRuneIndex {
				// 8 or more bytes left, continue to search
				currRuneIndex += 8
				byteOffset += 8
				continue
			}
			offset := targetRuneIndex - currRuneIndex
			return s[byteOffset+offset : byteOffset+offset+1], nil
		}

		// isolate continuation bytes, tests bit 7 == 1 and bit 6 == 0
		cont := (w & swarNonAsciiMask) & ^((w << 1) & swarNonAsciiMask)
		numRunes := 8 - bits.OnesCount64(cont)

		if currRuneIndex+numRunes <= targetRuneIndex {
			// check if the byte right across the boundary a continuation byte
			if byteOffset+8 < n && (s[byteOffset+8]&0xC0) == 0x80 {
				break
			}
			currRuneIndex += numRunes
			byteOffset += 8
			continue
		}
		break
	}

	// remainder loop for rest of string, always len < 8
	for relOffset, r := range s[byteOffset:] {
		if currRuneIndex == targetRuneIndex {
			start := byteOffset + relOffset
			if r == utf8.RuneError {
				if _, size := utf8.DecodeRuneInString(s[start:]); size == 1 {
					return "", MakeRuntimeError(errors.New("invalid UTF-8 sequence in string"))
				}
			}
			return s[start : start+utf8.RuneLen(r)], nil
		}
		currRuneIndex++
	}

	return "", MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", targetRuneIndex, currRuneIndex))
}

func handleSuperIndex(node *ast.SuperIndex, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		// return ValueNone, errors.New("ctx.Self not set")
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	keyId := index.AsStringConst()
	if keyId == 0 {
		name := index.String(ctx)
		keyId = ctx.State.Interner.Intern(name)
	}

	obj := ctx.Self.Object()

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	if val.IsNone() {
		return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", index.String(ctx)))
	}

	val, err = val.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	return val, nil
}

func handleInSuper(node *ast.InSuper, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		// return ValueNone, errors.New("ctx.Self not set")
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	keyId := index.AsStringConst()
	if keyId == 0 {
		name := index.String(ctx)
		keyId = ctx.State.Interner.Intern(name)
	}

	obj := ctx.Self.Object()

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	return MakeBool(!val.IsNone()), nil
}

func handleError(node *ast.Error, scopePtr uintptr, ctx Context) (Value, error) {
	msg, err := EvaluateNode(node.Expr, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !msg.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
	}
	return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))
}

func handleImport(node *ast.Import, ctx Context) (Value, error) {

	importer := ctx.State.Environment.Importer

	dir := filepath.Dir(node.Loc().FileName)
	importedPath := node.File.Value
	if importedPath == "" {
		return ValueNone, MakeRuntimeError(errors.New("couldn't open import \"\": the empty string is not a valid filename"))
	}

	var absPath string
	if filepath.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = filepath.Join(dir, importedPath)
	}

	v := importer.Get(absPath)
	if !v.IsNone() {
		return v, nil
	}

	finalPath := absPath

	importedNode, err := importer.ResolveImport(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return ValueNone, MakeRuntimeError(err)
		}

		if !filepath.IsAbs(importedPath) {
			for i := len(importer.JPaths) - 1; i >= 0; i-- {
				jpf := filepath.Join(importer.JPaths[i], importedPath)

				v := importer.Get(jpf)
				if !v.IsNone() {
					return v, nil
				}

				importedNode, err = importer.ResolveImport(jpf)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return ValueNone, MakeRuntimeError(err)
				}

				if importedNode != nil {
					finalPath = jpf
					break
				}
			}
		}
	}

	if importedNode == nil {
		return ValueNone, MakeRuntimeError(fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath))
	}

	importScope := CreateFileScope(finalPath, importer.BaseStd, ctx)

	importCtx := ctx
	importCtx.Self = ValueNone

	v, err = evaluateNodeLazy(importedNode, importScope, importCtx)
	if err != nil {
		return ValueNone, err
	}

	importer.Set(finalPath, v)

	return v, nil

}

func handleImportStr(node *ast.ImportStr, ctx Context) (Value, error) {

	f, size, err := openImportFile(ctx.State.Environment.Importer, node.Loc().FileName, node.File.Value)
	if err != nil {
		return ValueNone, err
	}
	if size == 0 {
		return MakeString("", ctx), nil
	}
	defer f.Close()

	ptr, buf := AllocStringBuilder(ctx, size)
	buf = buf[:size]

	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.EOF {
		return ValueNone, MakeRuntimeError(err)
	}

	return FinalizeStringBuilder(ptr, n), nil
}

func handleImportBin(node *ast.ImportBin, ctx Context) (Value, error) {

	f, size, err := openImportFile(ctx.State.Environment.Importer, node.Loc().FileName, node.File.Value)
	if err != nil {
		return ValueNone, err
	}
	if size == 0 {
		return MakeArray(nil, ctx), nil
	}
	defer f.Close()

	arr, arrVal := MakeArraySized(size, ctx)

	var chunk [32 * 1024]byte
	offset := 0
	for offset < size {
		n, err := f.Read(chunk[:])
		for i := range n {
			arr[offset+i] = MakeNumber(float64(chunk[i]))
		}
		offset += n
		if err != nil {
			if err == io.EOF {
				break
			}
			return ValueNone, MakeRuntimeError(err)
		}
	}

	return arrVal, nil
}

func openImportFile(importer *Importer, importedFrom, importedPath string) (*os.File, int, error) {
	if importedPath == "" {
		return nil, 0, MakeRuntimeError(errors.New("couldn't open import \"\": the empty string is not a valid filename"))
	}
	dir := filepath.Dir(importedFrom)

	var absPath string
	if filepath.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = filepath.Join(dir, importedPath)
	}

	finalPath := absPath

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, 0, MakeRuntimeError(err)
		}

		if !filepath.IsAbs(importedPath) {
			for i := len(importer.JPaths) - 1; i >= 0; i-- {
				jpf := filepath.Join(importer.JPaths[i], importedPath)

				fileInfo, err = os.Stat(jpf)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return nil, 0, MakeRuntimeError(err)
				}

				if fileInfo != nil {
					finalPath = jpf
					break
				}
			}
		}
	}

	if fileInfo == nil {
		return nil, 0, MakeRuntimeError(fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath))
	}

	if fileInfo.IsDir() {
		return nil, 0, MakeRuntimeError(fmt.Errorf("read %s: is a directory", finalPath))
	}

	size := int(fileInfo.Size())
	if size == 0 {
		return nil, 0, nil
	}

	f, err := os.Open(finalPath)
	if err != nil {
		return nil, 0, MakeRuntimeError(err)
	}

	return f, size, nil
}
