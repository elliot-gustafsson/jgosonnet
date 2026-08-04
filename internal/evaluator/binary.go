package evaluator

import (
	"errors"
	"fmt"
	"math"

	"github.com/google/go-jsonnet/ast"
)

func handleBinaryOp(op ast.BinaryOp, left, right Value, ctx Context) (Value, error) {

	switch op {
	case ast.BopPlus:
		return bopPlus(left, right, ctx)
	case ast.BopManifestEqual:
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(eq), nil
	case ast.BopManifestUnequal:
		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(!eq), nil
	case ast.BopGreater:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x > 0), nil
	case ast.BopGreaterEq:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x >= 0), nil
	case ast.BopLess:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x < 0), nil
	case ast.BopLessEq:
		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x <= 0), nil
	}

	if left.IsNumber() && right.IsNumber() {
		res, err := handleNumberOp(left.Number(), right.Number(), op)
		if err != nil {
			return ValueNone, err
		}
		return MakeNumber(res), nil
	}

	return ValueNone, fmt.Errorf("unhandled binary operation %s %s %s", left.Type().String(), op.String(), right.Type().String())
}

func bopPlus(left, right Value, ctx Context) (Value, error) {

	// Allow 123 + '123', should return '123123'
	if left.IsString() {
		rs, err := right.ToString(ctx)
		if err != nil {
			return ValueNone, err
		}
		ls := left.String(ctx)
		return MakeStringConcat(ls, rs, ctx), nil
	}

	if right.IsString() {
		ls, err := left.ToString(ctx)
		if err != nil {
			return ValueNone, err
		}
		rs := right.String(ctx)
		return MakeStringConcat(ls, rs, ctx), nil
	}

	if left.Type() != right.Type() {
		return ValueNone, TypeErrorSpecific(left.Type(), right.Type())
	}

	switch left.Type() {
	case ValueTypeNull:
		return MakeNull(), nil

	case ValueTypeString:
		ls := left.String(ctx)
		rs := right.String(ctx)
		return MakeStringConcat(ls, rs, ctx), nil

	case ValueTypeNumber:
		val := left.Number() + right.Number()
		return MakeNumber(val), nil

	case ValueTypeArray:
		leftArr := left.Array()
		if len(leftArr) == 0 {
			return right, nil
		}
		rightArr := right.Array()
		if len(rightArr) == 0 {
			return left, nil
		}

		val, arrVal := MakeArraySized(len(leftArr)+len(rightArr), ctx)
		copy(val, leftArr)
		copy(val[len(leftArr):], rightArr)
		return arrVal, nil

	case ValueTypeObject:
		// Virtually combine objects
		id := MergeObjects(left.Payload(), right.Payload(), ctx)
		return MakeObjectValue(id), nil
	default:
		return ValueNone, fmt.Errorf("bop plus: unexpected type %s", left.Type().String())
	}
}

func handleNumberOp(left, right float64, op ast.BinaryOp) (val float64, err error) {

	switch op {
	default:
		return 0, fmt.Errorf("unhandled number operation: %s", op.String())
	case ast.BopMinus:
		val = left - right
	case ast.BopDiv:
		val = left / right
	case ast.BopMult:
		val = left * right
	case ast.BopBitwiseAnd:
		val, err = builtinBitwiseAnd(left, right)
	case ast.BopBitwiseOr:
		val, err = builtinBitwiseOr(left, right)
	case ast.BopBitwiseXor:
		val, err = builtinBitwiseXor(left, right)
	case ast.BopShiftL:
		val, err = builtinShiftL(left, right)
	case ast.BopShiftR:
		val, err = builtinShiftR(left, right)
	}
	return
}

const (
	maxSafeIntValue float64 = (1 << 53) - 1
	minSafeIntValue float64 = -maxSafeIntValue
)

func liftBitwise(f func(int64, int64) int64, positiveRightArg bool) func(float64, float64) (float64, error) {
	return func(left, right float64) (float64, error) {

		if left < minSafeIntValue || left > maxSafeIntValue {
			err := fmt.Errorf("Bitwise operator argument %v outside of range [%v, %v]", left, int64(minSafeIntValue), int64(maxSafeIntValue))
			return 0, MakeRuntimeError(err)
		}
		if right < minSafeIntValue || right > maxSafeIntValue {
			err := fmt.Errorf("Bitwise operator argument %v outside of range [%v, %v]", right, int64(minSafeIntValue), int64(maxSafeIntValue))
			return 0, MakeRuntimeError(err)
		}
		if positiveRightArg && right < 0 {
			return 0, MakeRuntimeError(errors.New("Shift by negative exponent."))
		}
		res := float64(f(int64(left), int64(right)))
		if math.IsNaN(res) {
			return 0, MakeRuntimeError(errors.New("Not a number"))
		}
		if math.IsInf(res, 0) {
			return 0, MakeRuntimeError(errors.New("Overflow"))
		}
		return res, nil
	}
}

var builtinShiftL = liftBitwise(func(x, y int64) int64 { return x << uint(y%64) }, true)
var builtinShiftR = liftBitwise(func(x, y int64) int64 { return x >> uint(y%64) }, true)
var builtinBitwiseAnd = liftBitwise(func(x, y int64) int64 { return x & y }, false)
var builtinBitwiseOr = liftBitwise(func(x, y int64) int64 { return x | y }, false)
var builtinBitwiseXor = liftBitwise(func(x, y int64) int64 { return x ^ y }, false)
