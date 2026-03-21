package evaluator

import (
	"cmp"
	"fmt"
	"strconv"
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

type Thunk struct {
	Node                ast.Node
	ScopeId             uint32
	CapturedSelf        Value
	CapturedSuperOffset int

	Value Value
}

type Func = func(args []NamedValue, ctx Context) (Value, error)

type Value struct {
	t ValueType

	// Reference to id in arena for string, object, array, function, thunk
	refId uint32

	// Also holds 1.0 and 0.0 for bool
	num float64
}

type NamedValue struct {
	Key uint32

	Value
}

func MakeNull() Value {
	return Value{t: ValueTypeNull}
}

func MakeString(v string, ctx Context) Value {
	refId := ctx.Interner.Intern(v)
	return Value{t: ValueTypeString, refId: refId}
}

func MakeNumber(v float64) Value {
	return Value{t: ValueTypeNumber, num: v}
}

func MakeBool(v bool) Value {
	if v {
		return Value{t: ValueTypeBool, num: 1}
	}
	return Value{t: ValueTypeBool, num: 0}

}

func MakeObject(v Object, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Objects))
	ctx.Arena.Objects = append(ctx.Arena.Objects, v)
	return Value{t: ValueTypeObject, refId: refId}
}

func MakeArray(v []Value, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Arrays))
	ctx.Arena.Arrays = append(ctx.Arena.Arrays, v)
	return Value{t: ValueTypeArray, refId: refId}
}

func MakeFunction(v Func, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Funcs))
	ctx.Arena.Funcs = append(ctx.Arena.Funcs, v)
	return Value{t: ValueTypeFunction, refId: refId}
}

func MakeThunk(v Thunk, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Thunks))
	ctx.Arena.Thunks = append(ctx.Arena.Thunks, v)
	return Value{t: ValueTypeThunk, refId: refId}
}

func (v Value) Type() ValueType {
	return v.t
}

func (v Value) String(ctx Context) string {
	return ctx.Interner.Get(v.refId)
}

func (v Value) Number() float64 {
	return v.num
}

func (v Value) Bool() bool {
	return v.num == 1
}

func (v Value) Array(ctx Context) []Value {
	return ctx.Arena.Arrays[v.refId]
}

func (v Value) Object(ctx Context) *Object {
	return &ctx.Arena.Objects[v.refId]
}

func (v Value) Function(ctx Context) Func {
	return ctx.Arena.Funcs[v.refId]
}

func (v Value) Thunk(ctx Context) *Thunk {
	return &ctx.Arena.Thunks[v.refId]
}

func (v Value) Eval(ctx Context) (Value, error) {
	if !v.IsThunk() {
		return v, nil
	}
	thunk := v.Thunk(ctx)
	if !thunk.Value.IsNone() {
		return thunk.Value, nil
	}

	evalCtx := ctx
	evalCtx.Self = thunk.CapturedSelf
	evalCtx.SuperOffset = thunk.CapturedSuperOffset

	evaledVal, err := evaluateNode(thunk.Node, thunk.ScopeId, evalCtx)
	if err != nil {
		return Value{}, err
	}

	if evaledVal.IsThunk() {
		evaledVal, err = evaledVal.Eval(ctx)
		if err != nil {
			return Value{}, err
		}
	}

	thunk = v.Thunk(ctx)
	thunk.Value = evaledVal
	return evaledVal, nil
}

func (v Value) ToString(ctx Context) (string, error) {

	switch v.t {
	default:
		return "", fmt.Errorf("unhandled type %s, string conversion not available", v.t.String())
	case ValueTypeNull:
		return "null", nil
	case ValueTypeString:
		return v.String(ctx), nil
	case ValueTypeNumber:
		res := strconv.FormatFloat(v.Number(), 'f', -1, 64)
		return res, nil
	case ValueTypeBool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case ValueTypeObject, ValueTypeArray:
		var b strings.Builder
		err := ManifestJson(&b, v, ctx, JsonConfigToString)
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
	return v.t == ValueTypeNone
}

func (v Value) IsLiteral() bool {
	return v.t.IsLiteral()
}

func (v Value) IsNull() bool {
	return v.t == ValueTypeNull
}

func (v Value) IsString() bool {
	return v.t == ValueTypeString
}

func (v Value) IsNumber() bool {
	return v.t == ValueTypeNumber
}

func (v Value) IsBool() bool {
	return v.t == ValueTypeBool
}

func (v Value) IsThunk() bool {
	return v.t == ValueTypeThunk
}

func (v Value) IsObject() bool {
	return v.t == ValueTypeObject
}

func (v Value) IsFunction() bool {
	return v.t == ValueTypeFunction
}

func (v Value) IsArray() bool {
	return v.t == ValueTypeArray
}

func (v Value) IsEmpty(ctx Context) bool {
	switch v.Type() {
	default:
		return false
	case ValueTypeNull:
		return true
	case ValueTypeObject:
		return len(v.Object(ctx).GetLayers()) == 0
	case ValueTypeArray:
		return len(v.Array(ctx)) == 0
	}
}

func (v Value) Prune(ctx Context) (Value, error) {
	switch v.Type() {
	default:
		return Value{}, fmt.Errorf("unhandled type (%s) in Value.Prune()", v.Type())
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
				return Value{}, err
			}
			out, err := v.Prune(ctx)
			if err != nil {
				return Value{}, err
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
		return a.refId == b.refId, nil
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
		if a.refId == b.refId {
			return 0, nil
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
