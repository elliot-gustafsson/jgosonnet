package evaluator

import (
	"fmt"

	"github.com/google/go-jsonnet/ast"
)

func handleBinaryOp(op ast.BinaryOp, left, right Value, ctx Context) (res Value, err error) {

	switch op {
	default:
		return Value{}, fmt.Errorf("unhandled binary operation '%s'", op.String())
	case ast.BopPlus:
		res, err = bopPlus(left, right, ctx)
	case ast.BopMinus:
		res, err = bopMinus(left, right)
	case ast.BopDiv:
		res, err = bopDiv(left, right)
	case ast.BopManifestEqual:
		res, err = BopManifestEqual(left, right, ctx)
	case ast.BopManifestUnequal:
		res, err = bopManifestUnequal(left, right, ctx)
	case ast.BopMult:
		res, err = bopMultiply(left, right)
	case ast.BopGreater:
		res, err = BopGreater(left, right, ctx)
	case ast.BopLess:
		res, err = BopGreater(right, left, ctx)
	case ast.BopGreaterEq:
		res, err = bopGreaterEq(left, right, ctx)
	case ast.BopLessEq:
		res, err = bopGreaterEq(right, left, ctx)
	}
	return res, err
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

func BopManifestEqual(left, right Value, ctx Context) (Value, error) {

	if left.Type() != right.Type() {
		return MakeBool(false), nil
	}

	var res bool

	switch left.Type() {
	case ValueTypeNull:
		res = true
	case ValueTypeString:
		res = left.refId == right.refId
	case ValueTypeNumber:
		res = left.Number() == right.Number()
	case ValueTypeBool:
		res = left.Bool() == right.Bool()
	// case ValueTypeArray:
	// 	val := make([]Value, len(left.Array())+len(right.Array()))
	// 	copy(val, left.Array())
	// 	copy(val[len(left.Array()):], right.Array())
	// 	return &Value{Type: ValueTypeArray, ArrVal: val}, nil
	// case ValueTypeObject:
	// 	// Virtually combine objects
	// 	obj := MergeObjects(left.Object(), right.Object())
	// 	return &Value{Type: ValueTypeObject, ObjVal: obj}, nil
	default:
		return Value{}, fmt.Errorf("bop equal: unexpected type %s", left.Type().String())
	}

	return MakeBool(res), nil
}

func bopManifestUnequal(left, right Value, ctx Context) (Value, error) {

	if left.Type() != right.Type() {
		return MakeBool(true), nil
	}

	var res bool

	switch left.Type() {
	case ValueTypeNull:
		res = false
	case ValueTypeString:
		res = left.refId != right.refId
	case ValueTypeNumber:
		res = left.Number() != right.Number()
	case ValueTypeObject:
		planAs := CompileObjectPlan(left.Object(ctx), ctx)
		planBs := CompileObjectPlan(right.Object(ctx), ctx)

		if len(planAs) == 0 && len(planBs) == 0 {
			return MakeBool(false), nil
		}

		if len(planAs) != len(planBs) {
			return MakeBool(true), nil
		}

		return Value{}, fmt.Errorf("unsupported stuff bop unequal for object")
	default:
		return Value{}, fmt.Errorf("bop unequal: unexpected type %s", left.Type().String())
	}

	return MakeBool(res), nil
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

func BopGreater(left, right Value, ctx Context) (Value, error) {

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op greater (%s,%s)", left.Type().String(), right.Type().String())
	}

	var val bool
	switch left.Type() {
	case ValueTypeNumber:
		val = left.Number() > right.Number()
	case ValueTypeString:
		val = left.String(ctx) > right.String(ctx)
	case ValueTypeArray:
		val = len(left.Array(ctx)) > len(right.Array(ctx))
	case ValueTypeObject:
		val = left.Object(ctx).GetLength() > right.Object(ctx).GetLength()
	default:
		return Value{}, fmt.Errorf("bop greater: unhandled type %s", left.Type().String())
	}

	return MakeBool(val), nil
}

func bopGreaterEq(left, right Value, ctx Context) (Value, error) {

	if left.Type() != right.Type() {
		return Value{}, fmt.Errorf("non matching types passed to binary op greater eq (%s,%s)", left.Type().String(), right.Type().String())
	}

	var val bool
	switch left.Type() {
	case ValueTypeNumber:
		val = left.Number() >= right.Number()
	case ValueTypeString:
		val = left.String(ctx) >= right.String(ctx)
	case ValueTypeArray:
		val = len(left.Array(ctx)) >= len(right.Array(ctx))
	default:
		return Value{}, fmt.Errorf("bop greaterEq: unexpected type %s", left.Type().String())
	}

	return MakeBool(val), nil
}
