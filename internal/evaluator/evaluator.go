package evaluator

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	fast "github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
)

func EvaluateNode(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := evaluateNode(bp, offset, scopePtr, ctx)
	if err != nil {
		return ValueNone, WrapError(err, bp, offset, ctx)
	}
	if val.IsThunk() {
		val, err = val.Eval(ctx)
		if err != nil {
			return ValueNone, WrapError(err, bp, offset, ctx)
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

	s, scopePtr := ctx.NewScope(0, 2)

	s.Bindings[0] = NamedValue{Key: ctx.State.Interner.Intern("$std"), Value: fileStd}
	s.Bindings[1] = NamedValue{Key: ctx.State.Interner.Intern("std"), Value: fileStd}

	return scopePtr
}

func evaluateNodeLazy(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	tag := fast.ReadTag(bp, offset)

	switch tag {
	case fast.TagString:
		id := fast.ReadU32(bp, offset+1)
		return MakeStringConst(id), nil
	case fast.TagNull:
		return ValueNull, nil
	case fast.TagFalse:
		return ValueFalse, nil
	case fast.TagTrue:
		return ValueTrue, nil
	case fast.TagNumber:
		num := fast.ReadF64(bp, offset+1)
		return MakeNumber(num), nil
	case fast.TagSelf:
		return ctx.Self, nil
	default:
		return NewThunk(bp, offset, scopePtr, ctx), nil
	}
}

func evaluateNode(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	tag := fast.ReadTag(bp, offset)

	switch tag {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %d", tag)
	case fast.TagString:
		id := fast.ReadU32(bp, offset+1)
		return MakeStringConst(id), nil
	case fast.TagNull:
		return ValueNull, nil
	case fast.TagFalse:
		return ValueFalse, nil
	case fast.TagTrue:
		return ValueTrue, nil
	case fast.TagNumber:
		num := fast.ReadF64(bp, offset+1)
		return MakeNumber(num), nil
	case fast.TagSelf:
		return ctx.Self, nil

	case fast.TagVar:
		return handleVar(fast.ReadU32(bp, offset+1), scopePtr, ctx)
	case fast.TagLocal1:
		return handleLocal1(bp, offset+1, scopePtr, ctx)
	case fast.TagLocal:
		return handleLocal(bp, offset+1, scopePtr, ctx)

	case fast.TagObject0:
		return handleObject0(ctx), nil
	case fast.TagObjectSimple:
		return handleObjectSimple(bp, offset+1, scopePtr, ctx), nil
	case fast.TagObject:
		return handleObject(bp, offset+1, scopePtr, ctx)
	case fast.TagObjectComp:
		return handleObjectComp(bp, offset+1, scopePtr, ctx)

	case fast.TagArray0:
		_, val := MakeArraySized(0, ctx)
		return val, nil
	case fast.TagArray:
		return handleArray(bp, offset+1, scopePtr, ctx)
	case fast.TagArrayComp:
		return handleArrayComp(bp, offset+1, scopePtr, ctx)

	case fast.TagFunction:
		return handleFunction(bp, offset, scopePtr, ctx), nil
	case fast.TagApply0:
		return handleApply0(bp, offset, scopePtr, ctx)
	case fast.TagApply1:
		return handleApply1(bp, offset, scopePtr, ctx)
	case fast.TagApply2:
		return handleApply2(bp, offset, scopePtr, ctx)
	case fast.TagApply3:
		return handleApply3(bp, offset, scopePtr, ctx)
	case fast.TagApply4:
		return handleApply4(bp, offset, scopePtr, ctx)
	case fast.TagApply:
		return handleApply(bp, offset, scopePtr, ctx)

	case fast.TagIndex:
		return handleIndex(bp, offset+1, scopePtr, ctx)
	case fast.TagIndexConst:
		return handleIndexConst(bp, offset+1, scopePtr, ctx)
	case fast.TagSuperIndex:
		return handleSuperIndex(bp, offset+1, scopePtr, ctx)
	case fast.TagSuperIndexConst:
		return handleSuperIndexConst(fast.ReadU32(bp, offset+1), ctx)
	case fast.TagInSuper:
		return handleInSuper(bp, offset+1, scopePtr, ctx)

	case fast.TagConditional:
		return handleConditional(bp, offset+1, scopePtr, ctx)
	case fast.TagAssert:
		return handleAssert(bp, offset+1, scopePtr, ctx)
	case fast.TagError:
		return handleError(bp, offset+1, scopePtr, ctx)

	case fast.TagImport:
		return handleImport(bp, offset+1, ctx)
	case fast.TagImportStr:
		return handleImportStr(bp, offset+1, ctx)
	case fast.TagImportBin:
		return handleImportBin(bp, offset+1, ctx)

	case fast.TagBinaryAdd:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			val := left.Number() + right.Number()
			if val > math.MaxFloat64 || val < -math.MaxFloat64 {
				return ValueNone, OverflowErr
			}
			return MakeNumber(val), nil
		}
		return bopPlus(left, right, ctx)
	case fast.TagBinarySub:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(left - right)
	case fast.TagBinaryMul:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(left * right)
	case fast.TagBinaryDiv:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if right == 0 {
			return ValueNone, DivisionByZeroErr
		}
		return makeCheckedNumber(left / right)
	case fast.TagBinaryMod:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			rn := right.Number()
			if rn == 0 {
				return ValueNone, DivisionByZeroErr
			}
			res := math.Mod(left.Number(), rn)
			if math.IsNaN(res) {
				return ValueNone, NanErr
			}
			return MakeNumber(res), nil
		}
		if left.IsString() {
			args := arena.Alloc[NamedValue](ctx.State.Allocator, 2)
			args[0] = NamedValue{Value: left}
			args[1] = NamedValue{Value: right}
			return StdFormatFunc(args, ctx)
		}
		return ValueNone, MakeRuntimeError(fmt.Errorf("Operator %% cannot be used on types %s and %s.", left.Type().String(), right.Type().String()))

	case fast.TagBinaryBitAnd:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		val, err := builtinBitwiseAnd(left, right)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(val)
	case fast.TagBinaryBitOr:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		val, err := builtinBitwiseOr(left, right)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(val)
	case fast.TagBinaryBitXor:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		val, err := builtinBitwiseXor(left, right)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(val)
	case fast.TagBinaryShiftL:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		val, err := builtinShiftL(left, right)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(val)
	case fast.TagBinaryShiftR:
		left, right, err := extractBinaryArgsNums(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		val, err := builtinShiftR(left, right)
		if err != nil {
			return ValueNone, err
		}
		return makeCheckedNumber(val)

	case fast.TagBinaryEq:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() == right.Number()), nil
		}
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(eq), nil
	case fast.TagBinaryNeq:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() != right.Number()), nil
		}
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(!eq), nil
	case fast.TagBinaryLt:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() < right.Number()), nil
		}
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x < 0), nil
	case fast.TagBinaryLte:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() <= right.Number()), nil
		}
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x <= 0), nil
	case fast.TagBinaryGt:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() > right.Number()), nil
		}
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x > 0), nil
	case fast.TagBinaryGte:
		left, right, err := extractBinaryArgs(bp, offset+1, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if left.IsNumber() && right.IsNumber() {
			return MakeBool(left.Number() >= right.Number()), nil
		}
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x >= 0), nil
	case fast.TagBinaryIn:
		// Desugared semantics match std.objectHasAll(right, left): evaluate right (object) first, then left (string).
		right, err := EvaluateNode(bp, fast.ReadU32(bp, offset+2), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !right.IsObject() {
			return ValueNone, TypeErrorSpecific(ValueTypeObject, right.Type())
		}
		left, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !left.IsString() {
			return ValueNone, TypeErrorSpecific(ValueTypeString, left.Type())
		}
		keyId := left.AsStringConst()
		if keyId == 0 {
			keyId = ctx.State.Interner.Intern(left.String(ctx))
		}
		return MakeBool(right.Object().HasField(keyId, true, ctx)), nil

	case fast.TagBinaryAnd:
		left, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !left.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, left.Type())
		}
		if !left.Bool() {
			return ValueFalse, nil
		}
		right, err := EvaluateNode(bp, fast.ReadU32(bp, offset+2), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !right.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, right.Type())
		}
		return right, nil
	case fast.TagBinaryOr:
		left, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !left.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, left.Type())
		}
		if left.Bool() {
			return ValueTrue, nil
		}
		right, err := EvaluateNode(bp, fast.ReadU32(bp, offset+2), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !right.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, right.Type())
		}
		return right, nil

	case fast.TagUnaryNot:
		unary, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, unary.Type())
		}
		return MakeBool(!unary.Bool()), nil
	case fast.TagUnaryMinus:
		unary, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		return MakeNumber(-unary.Number()), nil
	case fast.TagUnaryPlus:
		unary, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		return unary, nil
	case fast.TagUnaryBitwiseNot:
		unary, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		val64 := int64(unary.Number())
		return MakeNumber(float64(^val64)), nil
	}
}

func handleArray(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	count := fast.ReadU32(bp, offset)
	arr, val := MakeArraySized(int(count), ctx)
	for i := range count {
		elOffset := fast.ReadU32(bp, offset+1+i)
		ev, err := evaluateNodeLazy(bp, elOffset, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		arr[i] = ev
	}
	return val, nil
}

func handleArrayComp(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	bodyOff := fast.ReadU32(bp, offset)
	clauseCount := fast.ReadU32(bp, offset+1)

	// Fast path: single `for` clause with no `if` filters
	if clauseCount == 1 && fast.ReadU32(bp, offset+4) == 0 {
		varId := fast.ReadU32(bp, offset+2)
		iterOff := fast.ReadU32(bp, offset+3)
		iterVal, err := EvaluateNode(bp, iterOff, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !iterVal.IsArray() {
			return ValueNone, TypeErrorSpecific(ValueTypeArray, iterVal.Type())
		}
		iterArr := iterVal.Array()
		res, arrVal := MakeArraySized(len(iterArr), ctx)
		for i, elem := range iterArr {
			s, iterScope := ctx.NewScope(scopePtr, 1)
			s.Bindings[0] = NamedValue{Key: varId, Value: elem}
			ev, err := evaluateNodeLazy(bp, bodyOff, iterScope, ctx)
			if err != nil {
				return ValueNone, err
			}
			res[i] = ev
		}
		return arrVal, nil
	}

	var out []Value
	if err := evalArrayCompClause(bp, offset+2, clauseCount, bodyOff, scopePtr, ctx, &out); err != nil {
		return ValueNone, err
	}
	return MakeArray(out, ctx), nil
}

func evalArrayCompClause(bp unsafe.Pointer, clausePos, remaining, bodyOff uint32, scopePtr uintptr, ctx Context, out *[]Value) error {
	varId := fast.ReadU32(bp, clausePos)
	iterOff := fast.ReadU32(bp, clausePos+1)
	condCount := fast.ReadU32(bp, clausePos+2)
	nextPos := clausePos + 3 + condCount

	iterVal, err := EvaluateNode(bp, iterOff, scopePtr, ctx)
	if err != nil {
		return err
	}
	if !iterVal.IsArray() {
		return TypeErrorSpecific(ValueTypeArray, iterVal.Type())
	}
	iterArr := iterVal.Array()

	if remaining == 1 && condCount == 0 {
		needed := len(*out) + len(iterArr)
		if needed > cap(*out) {
			newCap := max(needed, cap(*out)*2)
			grown := make([]Value, len(*out), newCap)
			copy(grown, *out)
			*out = grown
		}
		for _, elem := range iterArr {
			s, iterScope := ctx.NewScope(scopePtr, 1)
			s.Bindings[0] = NamedValue{Key: varId, Value: elem}
			ev, err := evaluateNodeLazy(bp, bodyOff, iterScope, ctx)
			if err != nil {
				return err
			}
			*out = append(*out, ev)
		}
		return nil
	}

	for _, elem := range iterArr {
		s, iterScope := ctx.NewScope(scopePtr, 1)
		s.Bindings[0] = NamedValue{Key: varId, Value: elem}

		passed := true
		for j := range condCount {
			condVal, err := EvaluateNode(bp, fast.ReadU32(bp, clausePos+3+j), iterScope, ctx)
			if err != nil {
				return err
			}
			if !condVal.IsBool() {
				return TypeErrorSpecific(ValueTypeBool, condVal.Type())
			}
			if !condVal.Bool() {
				passed = false
				break
			}
		}
		if !passed {
			continue
		}

		if remaining == 1 {
			ev, err := evaluateNodeLazy(bp, bodyOff, iterScope, ctx)
			if err != nil {
				return err
			}
			*out = append(*out, ev)
		} else {
			if err := evalArrayCompClause(bp, nextPos, remaining-1, bodyOff, iterScope, ctx, out); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleObject0(ctx Context) Value {
	allocator := ctx.State.Allocator
	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)
	obj := NewSingleLayerObject(allocator, layer)
	return MakeObjectValue(obj)
}

func handleObjectSimple(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) Value {
	allocator := ctx.State.Allocator
	fieldCount := int(fast.ReadU32(bp, offset))

	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)
	layer.Bp = bp
	layer.ParentScopePtr = scopePtr

	layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
	layer.NodeOffsets = arena.Alloc[uint32](allocator, fieldCount)
	layer.Meta = arena.Alloc[uint8](allocator, fieldCount)

	useMap := fieldCount > MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	base := offset + 1
	for i := range fieldCount {
		p := base + uint32(i)*3
		keyId := fast.ReadU32(bp, p)
		layer.Keys[i] = keyId
		layer.NodeOffsets[i] = fast.ReadU32(bp, p+1)
		layer.Meta[i] = uint8(fast.ReadU32(bp, p+2))
		if useMap {
			layer.Index.Append(keyId)
		}
	}

	obj := NewSingleLayerObject(allocator, layer)
	return MakeObjectValue(obj)
}

func handleObject(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	allocator := ctx.State.Allocator

	fieldCount := int(fast.ReadU32(bp, offset))
	localsCount := int(fast.ReadU32(bp, offset+1))
	assertsCount := fast.ReadU32(bp, offset+2)

	fieldsBase := offset + 3
	localsBase := fieldsBase + uint32(fieldCount)*3
	assertsBase := localsBase + uint32(localsCount)*2

	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)
	layer.Bp = bp
	layer.ParentScopePtr = scopePtr

	if fieldCount > 0 {
		layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
		layer.NodeOffsets = arena.Alloc[uint32](allocator, fieldCount)
		layer.Meta = arena.Alloc[uint8](allocator, fieldCount)
	}

	if localsCount > 0 {
		layer.LocalKeys = arena.Alloc[uint32](allocator, localsCount)
		layer.LocalNodeOffsets = arena.Alloc[uint32](allocator, localsCount)
		for i := range localsCount {
			p := localsBase + uint32(i)*2
			layer.LocalKeys[i] = fast.ReadU32(bp, p)
			layer.LocalNodeOffsets[i] = fast.ReadU32(bp, p+1)
		}
	}

	if assertsCount > 0 {
		layer.AssertOffset = assertsBase
		layer.AssertCount = assertsCount
	}

	useMap := fieldCount > MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	index := 0
	for i := range fieldCount {
		p := fieldsBase + uint32(i)*3
		keyOrOff := fast.ReadU32(bp, p)
		bodyOff := fast.ReadU32(bp, p+1)
		rawMeta := uint8(fast.ReadU32(bp, p+2))

		var keyId uint32
		if (rawMeta & fast.FieldFlagStaticKey) != 0 {
			keyId = keyOrOff
		} else {
			name, err := EvaluateNode(bp, keyOrOff, scopePtr, ctx)
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

		if useMap {
			if layer.Index.Append(keyId) == math.MaxUint32 {
				return ValueNone, MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(keyId)))
			}
		} else {
			for j := 0; j < index; j++ {
				if layer.Keys[j] == keyId {
					return ValueNone, MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(keyId)))
				}
			}
		}

		layer.Keys[index] = keyId
		layer.NodeOffsets[index] = bodyOff
		layer.Meta[index] = rawMeta &^ fast.FieldFlagStaticKey
		index++
	}

	if index < fieldCount {
		layer.Keys = layer.Keys[:index]
		layer.NodeOffsets = layer.NodeOffsets[:index]
		layer.Meta = layer.Meta[:index]
	}

	obj := NewSingleLayerObject(allocator, layer)
	return MakeObjectValue(obj), nil
}

type objCompEntry struct {
	keyId uint32
	scope uintptr
}

func handleObjectComp(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	keyOff := fast.ReadU32(bp, offset)
	bodyOff := fast.ReadU32(bp, offset+1)
	meta := uint8(fast.ReadU32(bp, offset+2))
	clauseCount := fast.ReadU32(bp, offset+3)

	var entries []objCompEntry
	if err := evalObjectCompClause(bp, offset+4, clauseCount, keyOff, scopePtr, ctx, &entries); err != nil {
		return ValueNone, err
	}

	allocator := ctx.State.Allocator
	n := len(entries)

	layer := arena.Create[Layer](allocator)
	arena.Memclr(layer)
	layer.Bp = bp

	if n > 0 {
		layer.Keys = arena.Alloc[uint32](allocator, n)
		layer.NodeOffsets = arena.Alloc[uint32](allocator, n)
		layer.Meta = arena.Alloc[uint8](allocator, n)
		scopes := arena.Alloc[uintptr](allocator, n)
		layer.ParentScopePtr = uintptr(unsafe.Pointer(unsafe.SliceData(scopes))) | 1

		var dt *utils.DescriptorTable
		if n > MaxLayerLinearKeys {
			dt = utils.NewEmptyDescriptorTable(allocator, n)
			layer.Index = dt
		}

		for i, ent := range entries {
			key := ent.keyId
			if dt != nil {
				if dt.Append(key) == math.MaxUint32 {
					return ValueNone, MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
				}
			} else {
				for j := 0; j < i; j++ {
					if layer.Keys[j] == key {
						return ValueNone, MakeRuntimeError(fmt.Errorf("Duplicate field name: %q", ctx.State.Interner.Get(key)))
					}
				}
			}
			layer.Keys[i] = key
			layer.NodeOffsets[i] = bodyOff
			layer.Meta[i] = meta
			scopes[i] = ent.scope
		}
	}

	obj := NewSingleLayerObject(allocator, layer)
	return MakeObjectValue(obj), nil
}

func evalObjectCompClause(bp unsafe.Pointer, clausePos, remaining, keyOff uint32, scopePtr uintptr, ctx Context, out *[]objCompEntry) error {
	varId := fast.ReadU32(bp, clausePos)
	iterOff := fast.ReadU32(bp, clausePos+1)
	condCount := fast.ReadU32(bp, clausePos+2)
	nextPos := clausePos + 3 + condCount

	iterVal, err := EvaluateNode(bp, iterOff, scopePtr, ctx)
	if err != nil {
		return err
	}
	if !iterVal.IsArray() {
		return TypeErrorSpecific(ValueTypeArray, iterVal.Type())
	}
	iterArr := iterVal.Array()

	for _, elem := range iterArr {
		s, iterScope := ctx.NewScope(scopePtr, 1)
		s.Bindings[0] = NamedValue{Key: varId, Value: elem}

		passed := true
		for j := range condCount {
			condVal, err := EvaluateNode(bp, fast.ReadU32(bp, clausePos+3+j), iterScope, ctx)
			if err != nil {
				return err
			}
			if !condVal.IsBool() {
				return TypeErrorSpecific(ValueTypeBool, condVal.Type())
			}
			if !condVal.Bool() {
				passed = false
				break
			}
		}
		if !passed {
			continue
		}

		if remaining == 1 {
			keyVal, err := EvaluateNode(bp, keyOff, iterScope, ctx)
			if err != nil {
				return err
			}
			if keyVal.IsNull() {
				continue
			}
			if !keyVal.IsString() {
				return fmt.Errorf("unexpected field name type %s, expected string", keyVal.Type().String())
			}
			keyId := keyVal.AsStringConst()
			if keyId == 0 {
				keyId = ctx.State.Interner.Intern(keyVal.String(ctx))
			}
			*out = append(*out, objCompEntry{keyId: keyId, scope: iterScope})
		} else {
			if err := evalObjectCompClause(bp, nextPos, remaining-1, keyOff, iterScope, ctx, out); err != nil {
				return err
			}
		}
	}
	return nil
}

func handleLocal1(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	keyId := fast.ReadU32(bp, offset)
	valOff := fast.ReadU32(bp, offset+1)
	bodyOff := fast.ReadU32(bp, offset+2)

	s, childScopePtr := ctx.NewScope(scopePtr, 1)
	t, err := evaluateNodeLazy(bp, valOff, childScopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	s.Bindings[0] = NamedValue{keyId, t}

	return EvaluateNode(bp, bodyOff, childScopePtr, ctx)
}

func handleLocal(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	bindCount := fast.ReadU32(bp, offset)
	s, childScopePtr := ctx.NewScope(scopePtr, int(bindCount))

	base := offset + 1
	for i := range bindCount {
		keyId := fast.ReadU32(bp, base+i*2)
		valOff := fast.ReadU32(bp, base+i*2+1)
		t, err := evaluateNodeLazy(bp, valOff, childScopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		s.Bindings[i] = NamedValue{keyId, t}
	}

	bodyOff := fast.ReadU32(bp, base+bindCount*2)
	return EvaluateNode(bp, bodyOff, childScopePtr, ctx)
}

func execApplyCall(fnVal Value, args []NamedValue, tailStrict bool, bp unsafe.Pointer, callOffset uint32, ctx Context) (Value, error) {
	res, err := fnVal.FunctionExecEx(args, ctx, tailStrict)
	if err != nil {
		if traceSig, ok := err.(*TraceSignal); ok {
			if perr := traceSig.Print(ctx.State.Environment.TraceOut, bp, callOffset, ctx); perr != nil {
				return ValueNone, perr
			}
			return traceSig.Rest, nil
		}
		return ValueNone, err
	}
	return res, nil
}

func handleApply0(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(bp, fast.ReadU32(bp, callOffset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}
	return execApplyCall(val, nil, false, bp, callOffset, ctx)
}

func handleApply1(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(bp, fast.ReadU32(bp, callOffset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}
	args := arena.Alloc[NamedValue](ctx.State.Allocator, 1)
	v0, err := evaluateNodeLazy(bp, fast.ReadU32(bp, callOffset+2), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	args[0] = NamedValue{Value: v0}
	return execApplyCall(val, args, false, bp, callOffset, ctx)
}

func handleApply2(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(bp, fast.ReadU32(bp, callOffset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}
	args := arena.Alloc[NamedValue](ctx.State.Allocator, 2)
	v0, err := evaluateNodeLazy(bp, fast.ReadU32(bp, callOffset+2), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	v1, err := evaluateNodeLazy(bp, fast.ReadU32(bp, callOffset+3), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	args[0] = NamedValue{Value: v0}
	args[1] = NamedValue{Value: v1}
	return execApplyCall(val, args, false, bp, callOffset, ctx)
}

func handleApply3(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(bp, fast.ReadU32(bp, callOffset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}
	args := arena.Alloc[NamedValue](ctx.State.Allocator, 3)
	for i := range uint32(3) {
		v, err := evaluateNodeLazy(bp, fast.ReadU32(bp, callOffset+2+i), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		args[i] = NamedValue{Value: v}
	}
	return execApplyCall(val, args, false, bp, callOffset, ctx)
}

func handleApply4(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(bp, fast.ReadU32(bp, callOffset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}
	args := arena.Alloc[NamedValue](ctx.State.Allocator, 4)
	for i := range uint32(4) {
		v, err := evaluateNodeLazy(bp, fast.ReadU32(bp, callOffset+2+i), scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		args[i] = NamedValue{Value: v}
	}
	return execApplyCall(val, args, false, bp, callOffset, ctx)
}

func handleApply(bp unsafe.Pointer, callOffset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	posCount := fast.ReadU32(bp, callOffset+1)
	namedCount := fast.ReadU32(bp, callOffset+2)
	flags := fast.ReadU32(bp, callOffset+3)
	funcOff := fast.ReadU32(bp, callOffset+4)

	val, err := EvaluateNode(bp, funcOff, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}

	totalArgs := int(posCount + namedCount)
	var args []NamedValue
	if totalArgs > 0 {
		args = arena.Alloc[NamedValue](ctx.State.Allocator, totalArgs)
		posBase := callOffset + 5
		for i := range posCount {
			v, err := evaluateNodeLazy(bp, fast.ReadU32(bp, posBase+i), scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
			args[i] = NamedValue{Value: v}
		}
		namedBase := posBase + posCount
		for i := range namedCount {
			keyId := fast.ReadU32(bp, namedBase+i*2)
			v, err := evaluateNodeLazy(bp, fast.ReadU32(bp, namedBase+i*2+1), scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}
			args[posCount+i] = NamedValue{Key: keyId, Value: v}
		}
	}

	tailStrict := (flags & fast.ApplyFlagTailStrict) != 0
	return execApplyCall(val, args, tailStrict, bp, callOffset, ctx)
}

func handleVar(keyId uint32, scopePtr uintptr, ctx Context) (Value, error) {
	val, found := ctx.GetScopeBind(scopePtr, keyId)
	if !found {
		return ValueNone, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", ctx.State.Interner.Get(keyId)))
	}
	return val, nil
}

func handleFunction(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) Value {
	f := arena.Create[Function](ctx.State.Allocator)
	f.Bp = bp
	f.ScopePtr = scopePtr
	f.Self = ctx.Self
	f.Offset = offset
	f.SuperOffset = ctx.SuperOffset
	return MakeFunctionValue(f, ctx)
}

func handleConditional(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	cond, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !cond.IsBool() {
		return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
	}

	if cond.Bool() {
		return EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
	}
	return EvaluateNode(bp, fast.ReadU32(bp, offset+2), scopePtr, ctx)
}

func handleAssert(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	cond, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !cond.IsBool() {
		return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
	}
	if cond.Bool() {
		return EvaluateNode(bp, fast.ReadU32(bp, offset+2), scopePtr, ctx)
	}

	msgOff := fast.ReadU32(bp, offset+1)
	if msgOff == fast.NoNode {
		return ValueNone, MakeRuntimeError(errors.New("Assertion failed"))
	}
	msg, err := EvaluateNode(bp, msgOff, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !msg.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
	}
	return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))
}

func handleIndexConst(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	target, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	keyId := fast.ReadU32(bp, offset+1)

	switch target.Type() {
	case ValueTypeObject:
		obj := target.Object()
		subCtx := ctx
		subCtx.Self = target

		val, _, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return ValueNone, err
		}
		if val.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", ctx.State.Interner.Get(keyId)))
		}
		val, err = val.Eval(subCtx)
		if err != nil {
			return ValueNone, err
		}
		return val, nil
	case ValueTypeString:
		return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing string, expected number, got string"))
	case ValueTypeArray:
		return ValueNone, TypeErrorSpecific(ValueTypeNumber, ValueTypeString)
	default:
		return ValueNone, fmt.Errorf("value not indexable: %s", target.Type().String())
	}
}

func handleIndex(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	target, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
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

func handleSuperIndexConst(keyId uint32, ctx Context) (Value, error) {
	if ctx.Self.IsNone() {
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
	val, err = val.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}
	return val, nil
}

func handleSuperIndex(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	if !index.IsString() {
		return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing object, expected string, got %s", index.Type().String()))
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

func handleInSuper(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		return ValueNone, MakeRuntimeError(errors.New("Attempt to use super when there is no super class."))
	}

	if !index.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, index.Type())
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

func handleError(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (Value, error) {
	msg, err := EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !msg.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
	}
	return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))
}

func handleImport(bp unsafe.Pointer, offset uint32, ctx Context) (Value, error) {
	importer := ctx.State.Environment.Importer

	var fromFile string
	if prog := lookupProgram(bp, ctx); prog != nil {
		fromFile = prog.Filename
	}
	dir := filepath.Dir(fromFile)
	importedPath := ctx.State.Interner.Get(fast.ReadU32(bp, offset))
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

	importedProg, err := importer.ResolveImport(absPath)
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

				importedProg, err = importer.ResolveImport(jpf)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return ValueNone, MakeRuntimeError(err)
				}

				if importedProg != nil {
					finalPath = jpf
					break
				}
			}
		}
	}

	if importedProg == nil {
		return ValueNone, MakeRuntimeError(fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath))
	}

	importScope := CreateFileScope(finalPath, importer.BaseStd, ctx)

	importCtx := ctx
	importCtx.Self = ValueNone
	importCtx.SuperOffset = 0

	v, err = evaluateNodeLazy(importedProg.Base(), importedProg.RootIdx, importScope, importCtx)
	if err != nil {
		return ValueNone, err
	}

	importer.Set(finalPath, v)
	return v, nil
}

func handleImportStr(bp unsafe.Pointer, offset uint32, ctx Context) (Value, error) {
	var fromFile string
	if prog := lookupProgram(bp, ctx); prog != nil {
		fromFile = prog.Filename
	}
	importedPath := ctx.State.Interner.Get(fast.ReadU32(bp, offset))

	f, size, err := openImportFile(ctx.State.Environment.Importer, fromFile, importedPath)
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

func handleImportBin(bp unsafe.Pointer, offset uint32, ctx Context) (Value, error) {
	var fromFile string
	if prog := lookupProgram(bp, ctx); prog != nil {
		fromFile = prog.Filename
	}
	importedPath := ctx.State.Interner.Get(fast.ReadU32(bp, offset))

	f, size, err := openImportFile(ctx.State.Environment.Importer, fromFile, importedPath)
	if err != nil {
		return ValueNone, err
	}
	if size == 0 {
		return MakeArray(nil, ctx), nil
	}
	defer f.Close()

	arr, arrVal := MakeArraySized(size, ctx)

	var chunk [32 * 1024]byte
	offsetBytes := 0
	for offsetBytes < size {
		n, err := f.Read(chunk[:])
		for i := range n {
			arr[offsetBytes+i] = MakeNumber(float64(chunk[i]))
		}
		offsetBytes += n
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
