package evaluator

import (
	"errors"
	"fmt"
	"math"
	"unsafe"

	fast "github.com/elliot-gustafsson/jgosonnet/internal/ast"
)

// StdFormatFunc is wired by stdlib during package init to support `"..." % vals`
// without creating an import cycle between evaluator and stdlib.
var StdFormatFunc Func

func extractBinaryArgs(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (left, right Value, err error) {
	left, err = EvaluateNode(bp, fast.ReadU32(bp, offset), scopePtr, ctx)
	if err != nil {
		return
	}
	right, err = EvaluateNode(bp, fast.ReadU32(bp, offset+1), scopePtr, ctx)
	if err != nil {
		return
	}
	return
}

func extractBinaryArgsNums(bp unsafe.Pointer, offset uint32, scopePtr uintptr, ctx Context) (left, right float64, err error) {
	leftV, rightV, err := extractBinaryArgs(bp, offset, scopePtr, ctx)
	if err != nil {
		return
	}
	if !leftV.IsNumber() {
		err = TypeErrorSpecific(ValueTypeNumber, leftV.Type())
		return
	}
	if !rightV.IsNumber() {
		err = TypeErrorSpecific(ValueTypeNumber, rightV.Type())
		return
	}
	left = leftV.Number()
	right = rightV.Number()
	return
}

func makeCheckedNumber(val float64) (Value, error) {
	if math.IsInf(val, 0) {
		return 0, OverflowErr
	}
	if math.IsNaN(val) {
		return 0, NanErr
	}
	return MakeNumber(val), nil
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
		if math.IsInf(val, 0) {
			return ValueNone, OverflowErr
		}
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

const (
	maxSafeIntValue float64 = (1 << 53) - 1
	minSafeIntValue float64 = -maxSafeIntValue
)

var DivisionByZeroErr = MakeRuntimeError(errors.New("Division by zero."))
var OverflowErr = MakeRuntimeError(errors.New("Overflow"))
var NanErr = MakeRuntimeError(errors.New("Not a number"))

func checkBitwiseArgs(left, right float64, positiveRightArg bool) (int64, int64, error) {
	if left < minSafeIntValue || left > maxSafeIntValue {
		err := fmt.Errorf("Bitwise operator argument %v outside of range [%v, %v]", left, int64(minSafeIntValue), int64(maxSafeIntValue))
		return 0, 0, MakeRuntimeError(err)
	}
	if right < minSafeIntValue || right > maxSafeIntValue {
		err := fmt.Errorf("Bitwise operator argument %v outside of range [%v, %v]", right, int64(minSafeIntValue), int64(maxSafeIntValue))
		return 0, 0, MakeRuntimeError(err)
	}
	if positiveRightArg && right < 0 {
		return 0, 0, MakeRuntimeError(errors.New("Shift by negative exponent."))
	}
	return int64(left), int64(right), nil
}

func finishBitwise(res float64) (float64, error) {
	if math.IsNaN(res) {
		return 0, NanErr
	}
	if res > math.MaxFloat64 || res < -math.MaxFloat64 {
		return 0, OverflowErr
	}
	return res, nil
}

func builtinShiftL(left, right float64) (float64, error) {
	x, y, err := checkBitwiseArgs(left, right, true)
	if err != nil {
		return 0, err
	}
	return finishBitwise(float64(x << uint(y%64)))
}

func builtinShiftR(left, right float64) (float64, error) {
	x, y, err := checkBitwiseArgs(left, right, true)
	if err != nil {
		return 0, err
	}
	return finishBitwise(float64(x >> uint(y%64)))
}

func builtinBitwiseAnd(left, right float64) (float64, error) {
	x, y, err := checkBitwiseArgs(left, right, false)
	if err != nil {
		return 0, err
	}
	return finishBitwise(float64(x & y))
}

func builtinBitwiseOr(left, right float64) (float64, error) {
	x, y, err := checkBitwiseArgs(left, right, false)
	if err != nil {
		return 0, err
	}
	return finishBitwise(float64(x | y))
}

func builtinBitwiseXor(left, right float64) (float64, error) {
	x, y, err := checkBitwiseArgs(left, right, false)
	if err != nil {
		return 0, err
	}
	return finishBitwise(float64(x ^ y))
}
