package evaluator

import (
	"cmp"
	"fmt"
	"math"
	"strings"

	"github.com/google/go-jsonnet/ast"
)

type ValueType uint8

const (
	ValueTypeNone ValueType = iota
	ValueTypeNull
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeObject
	ValueTypeArray
	ValueTypeFunction
	ValueTypeThunk
)

func (t ValueType) IsLiteral() bool {
	return t == ValueTypeString || t == ValueTypeNumber || t == ValueTypeBool || t == ValueTypeNull
}

func (t ValueType) String() string {
	switch t {
	case ValueTypeNone:
		return "none"
	case ValueTypeNull:
		return "null"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBool:
		return "boolean"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeFunction:
		return "function"
	case ValueTypeThunk:
		return "thunk"
	default:
		return fmt.Sprintf("unknown (%d)", t)
	}
}

type ThunkEvalFunc func(Context) (Value, error)

type Thunk struct {
	NodeId              uint32
	ScopeId             uint32
	CapturedSelf        Value
	CapturedSuperOffset int32

	Value Value
}

func NewThunk(node ast.Node, scopeId uint32, ctx Context) Thunk {
	return Thunk{
		NodeId:              ctx.State.Registry.Nodes.Alloc(node),
		ScopeId:             scopeId,
		CapturedSelf:        ctx.Self,
		CapturedSuperOffset: ctx.SuperOffset,
	}
}

type Func = func(args []NamedValue, ctx Context) (Value, error)

type Function struct {
	argsCount int

	fn Func
}

func NewFunction(argsCount int, fn Func) Function {
	return Function{
		argsCount: argsCount,
		fn:        fn,
	}
}

func (t Function) Exec(args []NamedValue, ctx Context) (Value, error) {
	if t.fn == nil {
		return ValueNone, fmt.Errorf("function not instantiated")
	}
	return t.fn(args, ctx)
}

func (t Function) Noop() bool {
	return t.fn == nil
}

func (t Function) Length() int {
	return t.argsCount
}

// Value represents a NaN-boxed 64-bit value.
// Bit layout for non-float types:
//
// 63                   51 50         40 39       32 31                              0
// |----------------------|-------------|-----------|--------------------------------|
// |  Float NaN Boundary  |    Flags    | ValueType |         RefId / Payload        |
// |       (13 bits)      |  (11 bits)  |  (8 bits) |            (32 bits)           |
// |----------------------|-------------|-----------|--------------------------------|
//
// - Bits 63-51: Reserved for IEEE-754 Quiet NaN float64 detection.
// - Bits 50-40: Available for custom boolean flags (e.g., FlagStringConst).
// - Bits 39-32: ValueType enum.
// - Bits 31-0 : uint32 Arena ID (or 1/0 for booleans).
type Value uint64

const (
	nanTag         uint64 = 0xFFF8000000000000
	floatThreshold uint64 = 1 << 19

	ValueNone Value = 0

	ValueFlagStringConst uint64 = 1 << 40
)

func box(t ValueType, refId uint32) Value {
	// pack the type into bits 32-47, and the refId into bits 0-31
	return Value((uint64(t) << 32) | uint64(refId))
}

type NamedValue struct {
	Key uint32

	Value
}

type CachedValue struct {
	Value

	Visible bool
}

func MakeNull() Value {
	return box(ValueTypeNull, 0)
}

func MakeString(v string, ctx Context) Value {
	id := ctx.State.Registry.Strings.Alloc(v)
	return box(ValueTypeString, id)
}

func MakeStringValue(id uint32) Value {
	return box(ValueTypeString, id)
}

func MakeStringConst(id uint32) Value {
	return Value(uint64(box(ValueTypeString, id)) | ValueFlagStringConst)
}

func MakeStringConcat(v1, v2 string, ctx Context) Value {
	id := ctx.State.Registry.Strings.AllocConcat(v1, v2)
	return box(ValueTypeString, id)
}

func MakeNumber(v float64) Value {
	bits := math.Float64bits(v)
	if math.IsNaN(v) {
		bits = 0x7FF8000000000000
	}
	return Value(bits ^ nanTag)
}

func MakeBool(v bool) Value {
	if v {
		return box(ValueTypeBool, 1)
	}
	return box(ValueTypeBool, 0)
}

func MakeObject(v Object, ctx Context) Value {
	refId := ctx.State.Registry.Objects.Alloc(v)
	return box(ValueTypeObject, refId)
}

func MakeObjectValue(id uint32) Value {
	return box(ValueTypeObject, id)
}

func MakeArray(v []Value, ctx Context) Value {
	refId := ctx.State.Registry.Arrays.Alloc(v)
	return box(ValueTypeArray, refId)
}

func MakeArraySized(l int, ctx Context) ([]Value, Value) {
	arr, refId := ctx.State.Registry.Arrays.Make(l)
	return arr, box(ValueTypeArray, refId)
}

func MakeFunction(v Function, ctx Context) Value {
	refId := ctx.State.Registry.Functions.Alloc(v)
	return box(ValueTypeFunction, refId)
}

func MakeThunk(v Thunk, ctx Context) Value {
	refId := ctx.State.Registry.Thunks.Alloc(v)
	return box(ValueTypeThunk, refId)
}

func MakeTombstoneValue(scope int) Value {
	return box(ValueTypeNone, uint32(scope))
}

func (v Value) Type() ValueType {
	t := uint64(v) >> 32
	if t >= floatThreshold {
		return ValueTypeNumber
	}
	return ValueType(t)
}

func (v Value) RefId() uint32 {
	return uint32(v)
}

func (v Value) String(ctx Context) string {
	if (uint64(v) & ValueFlagStringConst) != 0 {
		return ctx.State.Interner.Get(v.RefId())
	}
	return ctx.State.Registry.Strings.Get(v.RefId())
}

func (v Value) Number() float64 {
	return math.Float64frombits(uint64(v) ^ nanTag)
}

func (v Value) Bool() bool {
	// Because MakeBool stored 1 or 0 in the RefId slot
	return uint32(v) == 1
}

func (v Value) Array(ctx Context) []Value {
	return ctx.State.Registry.Arrays.Get(v.RefId())
}

func (v Value) Object(ctx Context) *Object {
	return ctx.State.Registry.Objects.GetPtr(v.RefId())
}

func (v Value) Function(ctx Context) Function {
	return ctx.State.Registry.Functions.GetValue(v.RefId())
}

func (v Value) Thunk(ctx Context) *Thunk {
	return ctx.State.Registry.Thunks.GetPtr(v.RefId())
}

func (v Value) Eval(ctx Context) (Value, error) {
	if !v.IsThunk() {
		return v, nil
	}
	return v.evalThunk(ctx)
}

//go:noinline
func (v Value) evalThunk(ctx Context) (Value, error) {
	thunk := v.Thunk(ctx)

	if !thunk.Value.IsNone() {
		return thunk.Value, nil
	}

	evalCtx := ctx
	evalCtx.Self = thunk.CapturedSelf
	evalCtx.SuperOffset = thunk.CapturedSuperOffset

	node := ctx.State.Registry.Nodes.GetValue(thunk.NodeId)

	evaledVal, err := EvaluateNode(node, thunk.ScopeId, evalCtx)
	if err != nil {
		return ValueNone, err
	}

	thunk = v.Thunk(ctx)
	thunk.Value = evaledVal
	return evaledVal, nil
}

type RuntimeError struct {
	err error
}

func (e RuntimeError) Error() string {
	return "RUNTIME ERROR: " + e.err.Error()
}

func MakeRuntimeError(err error) RuntimeError {
	return RuntimeError{err}
}

func TypeErrorSpecific(good, bad ValueType) error {
	return MakeRuntimeError(fmt.Errorf("Unexpected type %s, expected %s", bad, good))
}

func TypeErrorGeneral(bad ValueType) error {
	return MakeRuntimeError(fmt.Errorf("Unexpected type %s", bad))
}

func (v Value) EvalString(ctx Context) (string, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return "", err
	}
	if !x.IsString() {
		return "", TypeErrorSpecific(ValueTypeString, x.Type())
	}
	return x.String(ctx), nil
}

func (v Value) EvalNumber(ctx Context) (float64, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return 0, err
	}
	if !x.IsNumber() {
		return 0, TypeErrorSpecific(ValueTypeNumber, x.Type())
	}
	return x.Number(), nil
}

func (v Value) EvalInteger(ctx Context) (int, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return 0, err
	}
	if !x.IsNumber() {
		return 0, TypeErrorSpecific(ValueTypeNumber, x.Type())
	}
	intNum := int(x.Number())
	if float64(intNum) != x.Number() {
		return 0, MakeRuntimeError(fmt.Errorf("Expected an integer, but got %v", x.Number()))
	}
	return intNum, nil
}

func (v Value) EvalBool(ctx Context) (bool, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return false, err
	}
	if !x.IsBool() {
		return false, TypeErrorSpecific(ValueTypeBool, x.Type())
	}
	return x.Bool(), nil
}

func (v Value) EvalArray(ctx Context) ([]Value, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if !x.IsArray() {
		return nil, TypeErrorSpecific(ValueTypeArray, x.Type())
	}
	return x.Array(ctx), nil
}

func (v Value) EvalObject(ctx Context) (*Object, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return nil, err
	}
	if !x.IsObject() {
		return nil, TypeErrorSpecific(ValueTypeObject, x.Type())
	}
	return x.Object(ctx), nil
}

func (v Value) EvalFunction(ctx Context) (Function, error) {
	x, err := v.Eval(ctx)
	if err != nil {
		return Function{}, err
	}
	if !x.IsFunction() {
		return Function{}, TypeErrorSpecific(ValueTypeFunction, x.Type())
	}
	return x.Function(ctx), nil
}

func (v Value) ToString(ctx Context) (string, error) {

	switch v.Type() {
	default:
		return "", fmt.Errorf("unhandled type %s, string conversion not available", v.Type().String())
	case ValueTypeNull:
		return "null", nil
	case ValueTypeString:
		return v.String(ctx), nil
	case ValueTypeNumber:
		var p [64]byte
		res := unparseNumber(p[:0], v.Number())
		return string(res), nil
	case ValueTypeBool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case ValueTypeArray:
		var b strings.Builder
		err := ManifestJson(&b, v, ctx, JsonConfigToString)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	case ValueTypeObject:
		subCtx := ctx
		subCtx.Self = v
		var b strings.Builder
		err := ManifestJson(&b, v, subCtx, JsonConfigToString)
		if err != nil {
			return "", err
		}
		return b.String(), nil
	case ValueTypeThunk:
		v, err := v.Eval(ctx)
		if err != nil {
			return "", err
		}
		return v.ToString(ctx)
	}

}

func (v Value) IsNone() bool {
	return v == 0
}

func (v Value) IsLiteral() bool {
	return v.Type().IsLiteral()
}

func (v Value) IsNull() bool {
	return ValueType(uint64(v)>>32) == ValueTypeNull
}

func (v Value) IsString() bool {
	return ValueType(uint64(v)>>32) == ValueTypeString
}

func (v Value) IsNumber() bool {
	return (uint64(v) >> 32) >= floatThreshold
}

func (v Value) IsBool() bool {
	return ValueType(uint64(v)>>32) == ValueTypeBool
}

func (v Value) IsThunk() bool {
	return ValueType(uint64(v)>>32) == ValueTypeThunk
}

func (v Value) IsObject() bool {
	return ValueType(uint64(v)>>32) == ValueTypeObject
}

func (v Value) IsFunction() bool {
	return ValueType(uint64(v)>>32) == ValueTypeFunction
}

func (v Value) IsArray() bool {
	return ValueType(uint64(v)>>32) == ValueTypeArray
}

func (v Value) IsStringConst() bool {
	return ValueType(uint64(v)>>32) == ValueTypeString && (uint64(v)&ValueFlagStringConst) != 0
}

func (v Value) IsEmpty(ctx Context) bool {
	switch v.Type() {
	default:
		return false
	case ValueTypeNull:
		return true
	case ValueTypeObject:
		return v.Object(ctx).Length(ctx) == 0
	case ValueTypeArray:
		return len(v.Array(ctx)) == 0
	}
}

func (v Value) Prune(ctx Context) (Value, error) {
	switch v.Type() {
	default:
		return ValueNone, fmt.Errorf("unhandled type (%s) in Value.Prune()", v.Type())
	case ValueTypeNull, ValueTypeString, ValueTypeNumber, ValueTypeBool:
		return v, nil
	case ValueTypeObject:
		return v.Object(ctx).Prune(ctx)
	case ValueTypeArray:
		arr := v.Array(ctx)
		res := make([]Value, 0, len(arr))
		for _, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return ValueNone, err
			}
			out, err := v.Prune(ctx)
			if err != nil {
				return ValueNone, err
			}
			if out.IsEmpty(ctx) {
				continue
			}
			res = append(res, out)
		}
		return MakeArray(res, ctx), nil
	}

}

func (a Value) Equal(b Value, ctx Context) (bool, error) {
	if a.Type() != b.Type() {
		return false, nil
	}

	switch a.Type() {
	case ValueTypeNull:
		return true, nil
	case ValueTypeString:
		if a.IsStringConst() && b.IsStringConst() {
			return a.RefId() == b.RefId(), nil
		}

		if !a.IsStringConst() && !b.IsStringConst() && a.RefId() == b.RefId() {
			return true, nil
		}

		return a.String(ctx) == b.String(ctx), nil
	case ValueTypeNumber:
		return a.Number() == b.Number(), nil
	case ValueTypeBool:
		return a.Bool() == b.Bool(), nil
	case ValueTypeObject:
		return a.Object(ctx).Equal(b.Object(ctx), ctx)
	case ValueTypeArray:
		aArr := a.Array(ctx)
		bArr := b.Array(ctx)

		if len(aArr) != len(bArr) {
			return false, nil
		}

		for i, av := range aArr {
			av, err := av.Eval(ctx)
			if err != nil {
				return false, err
			}
			bv, err := bArr[i].Eval(ctx)
			if err != nil {
				return false, err
			}

			eq, err := av.Equal(bv, ctx)
			if err != nil {
				return false, err
			}

			if !eq {
				return false, nil
			}

		}

		return true, nil
	case ValueTypeThunk:
		a, err := a.Eval(ctx)
		if err != nil {
			return false, err
		}
		b, err := b.Eval(ctx)
		if err != nil {
			return false, err
		}
		return a.Equal(b, ctx)
	default:
		return false, fmt.Errorf("comparing types %s is not supported", a.Type().String())
	}

}

func (a Value) Compare(b Value, ctx Context) (int, error) {
	if a.Type() != b.Type() {
		return 0, fmt.Errorf("comparing %s with %s is not supported", a.Type().String(), b.Type().String())
	}

	switch a.Type() {
	case ValueTypeNull:
		return 0, nil
	case ValueTypeString:
		if a.IsStringConst() && b.IsStringConst() {
			if a.RefId() == b.RefId() {
				return 0, nil
			}
		} else if !a.IsStringConst() && !b.IsStringConst() {
			if a.RefId() == b.RefId() {
				return 0, nil
			}
		}
		return cmp.Compare(a.String(ctx), b.String(ctx)), nil
	case ValueTypeNumber:
		return cmp.Compare(a.Number(), b.Number()), nil
	// case ValueTypeBool:
	// 	return 0, fmt.Errorf("comparing boolean values is not supported")
	// case ValueTypeObject:
	// 	return a.Object(ctx).Compare(b.Object(ctx), ctx)
	case ValueTypeArray:
		aArr := a.Array(ctx)
		bArr := b.Array(ctx)

		i, j := 0, 0
		for i < len(aArr) && j < len(bArr) {

			x, err := aArr[i].Compare(bArr[j], ctx)
			if err != nil {
				return 0, err
			}

			if x != 0 {
				return x, nil
			}
			i++
			j++
		}
		return cmp.Compare(len(aArr), len(bArr)), nil
	case ValueTypeThunk:
		a, err := a.Eval(ctx)
		if err != nil {
			return 0, err
		}
		b, err := b.Eval(ctx)
		if err != nil {
			return 0, err
		}
		return a.Compare(b, ctx)
	default:
		return 0, fmt.Errorf("comparing types %s is not supported", a.Type().String())
	}

}
