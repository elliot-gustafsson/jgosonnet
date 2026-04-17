package stdlib

import (
	"math"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func liftNumeric(f func(float64) float64) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalNumber(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(a)
		return evaluator.MakeNumber(res), nil
	}
}

func liftNumeric2(f func(float64, float64) float64) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalNumber(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		b, err := args[1].EvalNumber(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(a, b)
		return evaluator.MakeNumber(res), nil
	}
}

func liftNumericToBoolean(f func(float64) bool) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].EvalNumber(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(a)
		return evaluator.MakeBool(res), nil
	}
}

var std_floor = liftNumeric(math.Floor)
var std_pow = liftNumeric2(math.Pow)
var std_modulo = liftNumeric2(math.Mod)
var std_sqrt = liftNumeric(math.Sqrt)
var std_hypot = liftNumeric2(math.Hypot)
var std_ceil = liftNumeric(math.Ceil)
var std_sin = liftNumeric(math.Sin)
var std_cos = liftNumeric(math.Cos)
var std_tan = liftNumeric(math.Tan)
var std_asin = liftNumeric(math.Asin)
var std_acos = liftNumeric(math.Acos)
var std_atan = liftNumeric(math.Atan)
var std_atan2 = liftNumeric2(math.Atan2)
var std_log = liftNumeric(math.Log)
var std_log2 = liftNumeric(math.Log2)
var std_log10 = liftNumeric(math.Log10)
var std_exp = liftNumeric(func(f float64) float64 {
	res := math.Exp(f)
	if res == 0 && f > 0 {
		return math.Inf(1)
	}
	return res
})
var std_mantissa = liftNumeric(func(f float64) float64 {
	mantissa, _ := math.Frexp(f)
	return mantissa
})
var std_exponent = liftNumeric(func(f float64) float64 {
	_, exponent := math.Frexp(f)
	return float64(exponent)
})
var std_round = liftNumeric(math.Round)
var std_isEven = liftNumericToBoolean(func(f float64) bool {
	i, _ := math.Modf(f) // Get the integral part of the float
	return math.Mod(i, 2) == 0
})
var std_isOdd = liftNumericToBoolean(func(f float64) bool {
	i, _ := math.Modf(f) // Get the integral part of the float
	return math.Mod(i, 2) != 0
})
var std_isInteger = liftNumericToBoolean(func(f float64) bool {
	_, frac := math.Modf(f) // Get the fraction part of the float
	return frac == 0
})
var std_isDecimal = liftNumericToBoolean(func(f float64) bool {
	_, frac := math.Modf(f) // Get the fraction part of the float
	return frac != 0
})
var std_max = liftNumeric2(math.Max)
var std_min = liftNumeric2(math.Max)
var std_abs = liftNumeric(math.Abs)
var std_sign = liftNumeric(func(f float64) float64 {
	if f == 0 {
		return 0
	}
	if f > 0 {
		return 1
	}
	return -1
})
var std_deg2rad = liftNumeric(func(f float64) float64 {
	return f * (math.Pi / 180.0)
})
var std_rad2deg = liftNumeric(func(f float64) float64 {
	return f * (180.0 / math.Pi)
})

func std_clamp(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	x, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !x.IsNumber() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, x.Type())
	}
	minVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !minVal.IsNumber() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, minVal.Type())
	}
	maxVal, err := args[2].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !maxVal.IsNumber() {
		return evaluator.Value{}, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, maxVal.Type())
	}

	if x.Number() <= minVal.Number() {
		return minVal, nil
	}

	if x.Number() >= maxVal.Number() {
		return maxVal, nil
	}
	return x, nil
}
