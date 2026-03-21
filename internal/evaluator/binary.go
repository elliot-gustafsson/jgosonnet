package evaluator

import (
	"fmt"

	"github.com/google/go-jsonnet/ast"
)

func handleBinaryOp(op ast.BinaryOp, left, right Value, ctx Context) (Value, error) {

	switch op {
	default:
		return Value{}, fmt.Errorf("unhandled binary operation '%s'", op.String())
	case ast.BopPlus:
		return bopPlus(left, right, ctx)
	case ast.BopMinus:
		return bopMinus(left, right)
	case ast.BopDiv:
		return bopDiv(left, right)
	case ast.BopMult:
		return bopMultiply(left, right)
	case ast.BopManifestEqual:
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(eq), nil
	case ast.BopManifestUnequal:
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(!eq), nil
	case ast.BopGreater:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(x > 0), nil
	case ast.BopGreaterEq:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(x >= 0), nil
	case ast.BopLess:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(x < 0), nil
	case ast.BopLessEq:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return Value{}, err
		}
		return MakeBool(x <= 0), nil
	}
}

func bopPlus(left, right Value, ctx Context) (Value, error) {

	// Allow 123 + '123', should return '123123'
	if left.IsString() {
		rs, err := right.ToString(ctx)
		if err != nil {
			return Value{}, err
		}
		res := left.String(ctx) + rs
		return MakeString(res, ctx), nil
	}

	if right.IsString() {
		ls, err := left.ToString(ctx)
		if err != nil {
			return Value{}, err
		}
		res := ls + right.String(ctx)
		return MakeString(res, ctx), nil
	}

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op plus (%s,%s)", left.Type().String(), right.Type().String())
	}

	switch left.Type() {
	case ValueTypeNull:
		return MakeNull(), nil

	case ValueTypeString:
		val := left.String(ctx) + right.String(ctx)
		return MakeString(val, ctx), nil

	case ValueTypeNumber:
		val := left.Number() + right.Number()
		return MakeNumber(val), nil

	case ValueTypeArray:
		leftArr := left.Array(ctx)
		rightArr := right.Array(ctx)
		val := make([]Value, len(leftArr)+len(rightArr))
		copy(val, leftArr)
		copy(val[len(leftArr):], rightArr)
		return MakeArray(val, ctx), nil

	case ValueTypeObject:
		// Virtually combine objects
		obj := MergeObjects(left.Object(ctx), right.Object(ctx))
		return MakeObject(obj, ctx), nil
	default:
		return Value{}, fmt.Errorf("bop plus: unexpected type %s", left.Type().String())
	}
}

func bopMinus(left, right Value) (Value, error) {

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op minus (%s,%s)", left.Type().String(), right.Type().String())
	}

	switch left.Type() {

	case ValueTypeNumber:
		val := left.Number() - right.Number()
		return MakeNumber(val), nil

	default:
		return Value{}, fmt.Errorf("bop minus: unexpected type %s", left.Type().String())
	}
}

func bopDiv(left, right Value) (Value, error) {

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op div (%s,%s)", left.Type().String(), right.Type().String())
	}

	switch left.Type() {

	case ValueTypeNumber:
		val := left.Number() / right.Number()
		return MakeNumber(val), nil

	default:
		return Value{}, fmt.Errorf("bop div: unexpected type %s", left.Type().String())
	}
}

func bopMultiply(left, right Value) (Value, error) {

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op multiply (%s,%s)", left.Type().String(), right.Type().String())
	}

	switch left.Type() {
	case ValueTypeNumber:
		val := left.Number() * right.Number()
		return MakeNumber(val), nil
	default:
		return Value{}, fmt.Errorf("bop multiply: unhandled type %s", left.Type().String())
	}
}
