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
	"github.com/google/go-jsonnet/ast"
)

//go:noinline
func wrapNodeErr(err error, n *IRNode) (Value, error) {
	return ValueNone, WrapError(err, n.Orig)
}

func EvaluateNode(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	switch n.Kind {
	case IRKindConst:
		return n.ConstVal, nil
	case IRKindSelf:
		return ctx.Self, nil
	case IRKindVar:
		slotPtr := getScopeSlotPtr(scopePtr, n.A, n.B)
		val := *slotPtr
		if val.IsThunk() {
			var err error
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
			// Indirection shortcutting: cache resolved value directly in scope slot
			*slotPtr = val
		}
		return val, nil
	case IRKindIndexConst:
		val, err := handleIndexConst(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	case IRKindBinary:
		val, err := handleBinary(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		return val, nil
	case IRKindConditional:
		val, err := handleConditional(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	case IRKindApply:
		val, err := handleApply(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	case IRKindLocal:
		val, err := handleLocal(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	case IRKindObject:
		val, err := handleDesugaredObject(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		return val, nil
	case IRKindArray:
		return handleArray(n, scopePtr, ctx)
	case IRKindFunction:
		return handleFunction(n, scopePtr, ctx)
	case IRKindIndex:
		val, err := handleIndex(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	case IRKindUnary:
		val, err := handleUnary(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		return val, nil
	default:
		val, err := evaluateIR(n, scopePtr, ctx)
		if err != nil {
			return wrapNodeErr(err, n)
		}
		if val.IsThunk() {
			val, err = val.Eval(ctx)
			if err != nil {
				return wrapNodeErr(err, n)
			}
		}
		return val, nil
	}
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

	keyId := ctx.State.Interner.InternStatic("thisFile")

	layer, dl := NewDynamicLayer(allocator)

	dl.Keys = arena.Alloc[uint32](allocator, 1)
	dl.Keys[0] = keyId

	dl.Values = arena.Alloc[Value](allocator, 1)
	dl.Values[0] = MakeString(filename, ctx)

	dl.Meta = arena.Alloc[uint8](allocator, 1)
	dl.Meta[0] = 0
	dl.HasHidden = true

	fileObj := NewSingleLayerObject(allocator, layer)

	mergedObj := MergeObjects(baseStd.Payload(), uintptr(unsafe.Pointer(fileObj)), ctx)
	fileStd := MakeObjectValue(mergedObj)

	// Root scope layout matches rootCompileScope: slot 0 = "$std", slot 1 = "std".
	slots, scopePtr := ctx.NewScope(0, 2)
	slots[0] = fileStd
	slots[1] = fileStd

	return scopePtr
}

// evaluateIRLazy defers evaluation when scopePtr itself is still being populated
// (e.g. mutually-recursive local/object-local/default-arg initializers).
// Outer scopes (n.A > 0) are already fully initialized and can be read directly
// without allocating a wrapper Thunk (preserving laziness if the slot holds a Thunk).
func evaluateIRLazy(n *IRNode, scopePtr uintptr, ctx Context) Value {
	switch n.Kind {
	case IRKindConst:
		return n.ConstVal
	case IRKindSelf:
		return ctx.Self
	case IRKindVar:
		if n.A > 0 {
			return getScopeSlot(scopePtr, n.A, n.B)
		}
		return NewThunk(n.Kind, unsafe.Pointer(n), scopePtr, ctx)
	default:
		return NewThunk(n.Kind, unsafe.Pointer(n), scopePtr, ctx)
	}
}

// evaluateIRLazyReady defers evaluation when all bindings in scopePtr are already
// initialized (function call arguments, array elements). Variable references copy
// the existing Value/Thunk directly instead of allocating a redundant wrapper Thunk.
func evaluateIRLazyReady(n *IRNode, scopePtr uintptr, ctx Context) Value {
	switch n.Kind {
	case IRKindConst:
		return n.ConstVal
	case IRKindSelf:
		return ctx.Self
	case IRKindVar:
		return getScopeSlot(scopePtr, n.A, n.B)
	default:
		return NewThunk(n.Kind, unsafe.Pointer(n), scopePtr, ctx)
	}
}

func evaluateIR(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	switch n.Kind {
	case IRKindConst:
		return n.ConstVal, nil
	case IRKindSelf:
		return ctx.Self, nil
	case IRKindVar:
		return getScopeSlot(scopePtr, n.A, n.B), nil
	case IRKindIndexConst:
		return handleIndexConst(n, scopePtr, ctx)
	case IRKindIndex:
		return handleIndex(n, scopePtr, ctx)
	case IRKindBinary:
		return handleBinary(n, scopePtr, ctx)
	case IRKindUnary:
		return handleUnary(n, scopePtr, ctx)
	case IRKindConditional:
		return handleConditional(n, scopePtr, ctx)
	case IRKindApply:
		return handleApply(n, scopePtr, ctx)
	case IRKindLocal:
		return handleLocal(n, scopePtr, ctx)
	case IRKindObject:
		return handleDesugaredObject(n, scopePtr, ctx)
	case IRKindArray:
		return handleArray(n, scopePtr, ctx)
	case IRKindFunction:
		return handleFunction(n, scopePtr, ctx)
	case IRKindSuperIndex:
		return handleSuperIndex(n, scopePtr, ctx)
	case IRKindInSuper:
		return handleInSuper(n, scopePtr, ctx)
	case IRKindImport:
		return handleImport(n, ctx)
	case IRKindImportStr:
		return handleImportStr(n, ctx)
	case IRKindImportBin:
		return handleImportBin(n, ctx)
	case IRKindError:
		return handleError(n, scopePtr, ctx)
	case IRKindInvalidNumber:
		numNode := n.Orig.(*ast.LiteralNumber)
		num, err := strconv.ParseFloat(numNode.OriginalString, 64)
		if err != nil {
			return ValueNone, fmt.Errorf("(%T) failed to parse float val (%s), err: %w", numNode, numNode.OriginalString, err)
		}
		return MakeNumber(num), nil
	case IRKindVarUnresolved:
		varNode := n.Orig.(*ast.Var)
		return ValueNone, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", string(varNode.Id)))
	default:
		return ValueNone, fmt.Errorf("unhandled IR node kind: %d", n.Kind)
	}
}

type GoCallbackNode struct {
	FuncVal Value
	Args    []NamedValue
}

func handleArray(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	elems := (*IRArray)(n.Extra).Elements
	arr, val := MakeArraySized(len(elems), ctx)
	for i, e := range elems {
		arr[i] = evaluateIRLazyReady(e, scopePtr, ctx)
	}
	return val, nil
}

func handleDesugaredObject(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	irObj := (*IRObject)(n.Extra)
	allocator := ctx.State.Allocator

	if tmpl := irObj.StaticTemplate; tmpl != nil {
		obj := NewStaticSingleLayerObject(allocator, tmpl, scopePtr)
		return MakeObjectValue(obj), nil
	}

	fieldCount := len(irObj.Fields)
	layer, dl := NewDynamicLayer(allocator)
	layer.ParentScopePtr = scopePtr

	if fieldCount > 0 {
		dl.Keys = arena.Alloc[uint32](allocator, fieldCount)
		dl.Nodes = arena.AllocZeroed[*IRNode](allocator, fieldCount)
		dl.Meta = arena.Alloc[uint8](allocator, fieldCount)
	}

	dl.LocalNodes = irObj.LocalNodes
	dl.Asserts = irObj.Asserts

	index := 0
	for i := range irObj.Fields {
		f := &irObj.Fields[i]
		keyId := f.KeyId
		if f.Name != nil {
			name, err := EvaluateNode(f.Name, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
			if name.IsNull() {
				continue
			}
			if !name.IsString() {
				return ValueNone, fmt.Errorf("unexpected field name type %s, expected string", name.Type().String())
			}
			keyId = name.AsStringConst()
			if keyId == 0 {
				keyId = ctx.State.Interner.Intern(name.String(ctx))
			}
		}

		dl.Keys[index] = keyId
		dl.Nodes[index] = f.Body
		dl.Meta[index] = f.Meta
		if ast.ObjectFieldHide(f.Meta&MaskVisibility) == ast.ObjectFieldHidden {
			dl.HasHidden = true
		}
		index++
	}

	if index < fieldCount {
		dl.Keys = dl.Keys[:index]
		dl.Nodes = dl.Nodes[:index]
		dl.Meta = dl.Meta[:index]
	}

	obj := NewSingleLayerObject(allocator, layer)
	return MakeObjectValue(obj), nil
}

func handleLocal(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	irLocal := (*IRLocal)(n.Extra)
	slots, childScopePtr := ctx.NewScope(scopePtr, len(irLocal.Binds))

	for i, b := range irLocal.Binds {
		slots[i] = evaluateIRLazy(b, childScopePtr, ctx)
	}

	return EvaluateNode(irLocal.Body, childScopePtr, ctx)
}

func handleApply(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	var val Value
	var err error
	if n.Child0.Kind == IRKindConst {
		val = n.Child0.ConstVal
	} else {
		val, err = EvaluateNode(n.Child0, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
	}

	irApply := (*IRApply)(n.Extra)
	posNodes := irApply.Positional
	namedNodes := irApply.NamedArgs
	posCount := len(posNodes)
	nameCount := len(namedNodes)

	t := uint64(val) >> typeShift

	// Fast path: direct positional binding into user function scope (no intermediate args slice!)
	if t == uint64(ValueTypeFunction) && nameCount == 0 && !irApply.TailStrict {
		f := val.Function()
		ir := f.IR
		paramCount := len(ir.ParamKeyIds)
		if posCount == paramCount {
			if ctx.Depth >= ctx.State.MaxStack {
				return errMaxStackExceeded()
			}
			callCtx := ctx
			callCtx.Depth++
			callCtx.Self = f.Self
			callCtx.SuperOffset = f.SuperOffset

			if paramCount == 0 {
				return EvaluateNode(ir.Body, f.ScopePtr, callCtx)
			}

			slots, childScopeId := callCtx.NewScope(f.ScopePtr, paramCount)
			for i, a := range posNodes {
				slots[i] = evaluateIRLazyReady(a, scopePtr, ctx)
			}
			return EvaluateNode(ir.Body, childScopeId, callCtx)
		}
	}

	if t != uint64(ValueTypeFunction) && t != uint64(ValueTypeNativeFunction) {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}

	var args []NamedValue
	totalArgs := posCount + nameCount
	if totalArgs > 0 {
		args = arena.Alloc[NamedValue](ctx.State.Allocator, totalArgs)
		for i, a := range posNodes {
			args[i] = NamedValue{Value: evaluateIRLazyReady(a, scopePtr, ctx)}
		}
		namedKeys := irApply.NamedKeys
		for i, a := range namedNodes {
			args[i+posCount] = NamedValue{
				Key:   namedKeys[i],
				Value: evaluateIRLazyReady(a, scopePtr, ctx),
			}
		}
	}

	var res Value
	if t == uint64(ValueTypeNativeFunction) {
		res, err = execNativeFunction(val, args, ctx, irApply.TailStrict)
	} else {
		res, err = execUserFunction(val, args, ctx, irApply.TailStrict)
	}

	if err != nil {
		if traceSig, ok := err.(*TraceSignal); ok {
			err := traceSig.Print(ctx.State.Environment.TraceOut, n.Orig)
			if err != nil {
				return ValueNone, err
			}
			return traceSig.Rest, nil
		}
		return ValueNone, err
	}
	return res, nil
}

func handleBinary(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	var left Value
	var err error
	if n.Child0.Kind == IRKindConst {
		left = n.Child0.ConstVal
	} else {
		left, err = EvaluateNode(n.Child0, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
	}

	op := ast.BinaryOp(n.SubOp)

	switch op {
	case ast.BopAnd:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", left.Type().String())
		}
		if !left.Bool() {
			return MakeBool(false), nil
		}
		var right Value
		if n.Child1.Kind == IRKindConst {
			right = n.Child1.ConstVal
		} else {
			right, err = EvaluateNode(n.Child1, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
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
		var right Value
		if n.Child1.Kind == IRKindConst {
			right = n.Child1.ConstVal
		} else {
			right, err = EvaluateNode(n.Child1, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
		}
		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}
		return right, nil

	default:
		var right Value
		if n.Child1.Kind == IRKindConst {
			right = n.Child1.ConstVal
		} else {
			right, err = EvaluateNode(n.Child1, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
		}
		return handleBinaryOp(op, left, right, ctx)
	}
}

func handleUnary(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	unary, err := EvaluateNode(n.Child0, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	op := ast.UnaryOp(n.SubOp)
	switch op {
	default:
		return ValueNone, fmt.Errorf("unhandled unary type: %s", op.String())
	case ast.UopNot:
		if !unary.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, unary.Type())
		}
		return MakeBool(!unary.Bool()), nil
	case ast.UopMinus:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		return MakeNumber(-unary.Number()), nil
	case ast.UopBitwiseNot:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		val32 := int64(unary.Number())
		return MakeNumber(float64(^val32)), nil
	case ast.UopPlus:
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		return unary, nil
	}
}

func handleFunction(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	f := arena.CreateZeroed[Function](ctx.State.Allocator)
	*f = Function{
		IR:          (*IRFunction)(n.Extra),
		ScopePtr:    scopePtr,
		Self:        ctx.Self,
		SuperOffset: ctx.SuperOffset,
	}
	return MakeFunctionValue(f, ctx), nil
}

func handleConditional(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	cond, err := EvaluateNode(n.Child0, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !cond.IsBool() {
		return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
	}

	var branch *IRNode
	if cond.Bool() {
		branch = n.Child1
	} else {
		branch = (*IRNode)(n.Extra)
	}
	if branch.Kind == IRKindConst {
		return branch.ConstVal, nil
	}
	return EvaluateNode(branch, scopePtr, ctx)
}

func handleIndexConst(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	var target Value
	var err error
	if n.Child0.Kind == IRKindSelf {
		target = ctx.Self
	} else {
		target, err = EvaluateNode(n.Child0, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
	}

	if target.IsObject() {
		keyId := n.A
		obj := target.Object()

		subCtx := ctx
		subCtx.Self = target

		val, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return ValueNone, err
		}
		if val.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", ctx.State.Interner.Get(keyId)))
		}
		return val.Eval(subCtx)
	}

	switch target.Type() {
	case ValueTypeString:
		return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing string, expected number, got %s", ValueTypeString.String()))
	case ValueTypeArray:
		return ValueNone, TypeErrorSpecific(ValueTypeNumber, ValueTypeString)
	default:
		return ValueNone, fmt.Errorf("value not indexable: %s", target.Type().String())
	}
}

func handleIndex(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	var index Value
	var err error
	if n.Child1.Kind == IRKindConst {
		index = n.Child1.ConstVal
	} else {
		index, err = EvaluateNode(n.Child1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
	}

	target, err := EvaluateNode(n.Child0, scopePtr, ctx)
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

		subCtx := ctx
		subCtx.Self = target

		val, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return ValueNone, err
		}
		if val.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", index.String(ctx)))
		}
		return val.Eval(subCtx)
	case ValueTypeArray:
		if !index.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, index.Type())
		}
		i := int(index.Number())

		arr := target.Array()
		if i < 0 || len(arr) <= i {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", i, len(arr)))
		}
		return arr[i], nil
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

		if (w & swarNonAsciiMask) == 0 {
			if currRuneIndex+8 <= targetRuneIndex {
				currRuneIndex += 8
				byteOffset += 8
				continue
			}
			offset := targetRuneIndex - currRuneIndex
			return s[byteOffset+offset : byteOffset+offset+1], nil
		}

		cont := (w & swarNonAsciiMask) & ^((w << 1) & swarNonAsciiMask)
		numRunes := 8 - bits.OnesCount64(cont)

		if currRuneIndex+numRunes <= targetRuneIndex {
			if byteOffset+8 < n && (s[byteOffset+8]&0xC0) == 0x80 {
				break
			}
			currRuneIndex += numRunes
			byteOffset += 8
			continue
		}
		break
	}

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

func handleSuperIndex(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	keyId := n.A
	if keyId == 0 {
		index, err := EvaluateNode(n.Child0, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if ctx.Self.IsNone() {
			return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
		}
		keyId = index.AsStringConst()
		if keyId == 0 {
			keyId = ctx.State.Interner.Intern(index.String(ctx))
		}
	} else if ctx.Self.IsNone() {
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	obj := ctx.Self.Object()
	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	if val.IsNone() {
		return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", ctx.State.Interner.Get(keyId)))
	}

	return val.Eval(ctx)
}

func handleInSuper(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	keyId := n.A
	if keyId == 0 {
		index, err := EvaluateNode(n.Child0, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if ctx.Self.IsNone() {
			return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
		}
		keyId = index.AsStringConst()
		if keyId == 0 {
			keyId = ctx.State.Interner.Intern(index.String(ctx))
		}
	} else if ctx.Self.IsNone() {
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	obj := ctx.Self.Object()
	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	return MakeBool(!val.IsNone()), nil
}

func handleError(n *IRNode, scopePtr uintptr, ctx Context) (Value, error) {
	msg, err := EvaluateNode(n.Child0, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !msg.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
	}
	return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))
}

func handleImport(n *IRNode, ctx Context) (Value, error) {
	importer := ctx.State.Environment.Importer
	irImp := (*IRImport)(n.Extra)
	importedPath := irImp.File
	if importedPath == "" {
		return ValueNone, MakeRuntimeError(errors.New("couldn't open import \"\": the empty string is not a valid filename"))
	}

	absPath := irImp.Abs

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

	v = evaluateIRLazy(importedNode, importScope, importCtx)

	importer.Set(finalPath, v)

	return v, nil
}

func handleImportStr(n *IRNode, ctx Context) (Value, error) {
	irImp := (*IRImport)(n.Extra)
	f, size, err := openImportFile(ctx.State.Environment.Importer, irImp)
	if err != nil {
		return ValueNone, err
	}
	if size == 0 {
		return MakeString("", ctx), nil
	}
	defer f.Close()

	ptr, buf := AllocStringBuilder(ctx, size)
	buf = buf[:size]

	numRead, err := io.ReadFull(f, buf)
	if err != nil && err != io.EOF {
		return ValueNone, MakeRuntimeError(err)
	}

	return FinalizeStringBuilder(ptr, numRead), nil
}

func handleImportBin(n *IRNode, ctx Context) (Value, error) {
	irImp := (*IRImport)(n.Extra)
	f, size, err := openImportFile(ctx.State.Environment.Importer, irImp)
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
		numRead, err := f.Read(chunk[:])
		for i := range numRead {
			arr[offset+i] = MakeNumber(float64(chunk[i]))
		}
		offset += numRead
		if err != nil {
			if err == io.EOF {
				break
			}
			return ValueNone, MakeRuntimeError(err)
		}
	}

	return arrVal, nil
}

func openImportFile(importer *Importer, irImp *IRImport) (*os.File, int, error) {
	importedPath := irImp.File
	if importedPath == "" {
		return nil, 0, MakeRuntimeError(errors.New("couldn't open import \"\": the empty string is not a valid filename"))
	}

	absPath := irImp.Abs
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
