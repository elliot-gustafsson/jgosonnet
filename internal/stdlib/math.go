package stdlib

import (
	"fmt"
	"math"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func liftNumeric(f func(float64) float64, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		if !args[0].IsNumber() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s: %s, expected number", name, args[0].Type().String())
		}
		res := f(args[0].Number())
		return evaluator.MakeNumber(res), nil
	}
}

func liftNumeric2(f func(float64, float64) float64, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 2 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 2", name, len(args))
		}
		if !args[0].IsNumber() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected number", name, args[0].Type().String())
		}
		if !args[1].IsNumber() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 1): %s, expected number", name, args[1].Type().String())
		}
		res := f(args[0].Number(), args[1].Number())
		return evaluator.MakeNumber(res), nil
	}
}

func liftNumericToBoolean(f func(float64) bool, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		if !args[0].IsNumber() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to %s (arg 0): %s, expected number", name, args[0].Type().String())
		}
		res := f(args[0].Number())
		return evaluator.MakeBool(res), nil
	}
}

var std_floor = liftNumeric(math.Floor, "std.floor")
var std_pow = liftNumeric2(math.Pow, "std.pow")
var std_modulo = liftNumeric2(math.Mod, "std.modulo")
var std_sqrt = liftNumeric(math.Sqrt, "std.sqrt")
var std_hypot = liftNumeric2(math.Hypot, "std.hypot")
var std_ceil = liftNumeric(math.Ceil, "std.ceil")
var std_sin = liftNumeric(math.Sin, "std.sin")
var std_cos = liftNumeric(math.Cos, "std.cos")
var std_tan = liftNumeric(math.Tan, "std.tan")
var std_asin = liftNumeric(math.Asin, "std.asin")
var std_acos = liftNumeric(math.Acos, "std.acos")
var std_atan = liftNumeric(math.Atan, "std.atan")
var std_atan2 = liftNumeric2(math.Atan2, "std.atan2")
var std_log = liftNumeric(math.Log, "std.log")
var std_exp = liftNumeric(func(f float64) float64 {
	res := math.Exp(f)
	if res == 0 && f > 0 {
		return math.Inf(1)
	}
	return res
}, "std.exp")
var std_mantissa = liftNumeric(func(f float64) float64 {
	mantissa, _ := math.Frexp(f)
	return mantissa
}, "std.mantissa")
var std_exponent = liftNumeric(func(f float64) float64 {
	_, exponent := math.Frexp(f)
	return float64(exponent)
}, "std.exponent")
var std_round = liftNumeric(math.Round, "std.round")
var std_isEven = liftNumericToBoolean(func(f float64) bool {
	i, _ := math.Modf(f) // Get the integral part of the float
	return math.Mod(i, 2) == 0
}, "std.isEven")
var std_isOdd = liftNumericToBoolean(func(f float64) bool {
	i, _ := math.Modf(f) // Get the integral part of the float
	return math.Mod(i, 2) != 0
}, "std.isOdd")
var std_isInteger = liftNumericToBoolean(func(f float64) bool {
	_, frac := math.Modf(f) // Get the fraction part of the float
	return frac == 0
}, "std.isInteger")
var std_isDecimal = liftNumericToBoolean(func(f float64) bool {
	_, frac := math.Modf(f) // Get the fraction part of the float
	return frac != 0
}, "std.isDecimal")
