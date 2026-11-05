package stdlib

import (
	"math"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func makeCheckedNumber(res float64) (evaluator.Value, error) {
	if math.IsNaN(res) {
		return evaluator.ValueNone, evaluator.NanErr
	}
	if res > math.MaxFloat64 || res < -math.MaxFloat64 {
		return evaluator.ValueNone, evaluator.OverflowErr
	}
	return evaluator.MakeNumber(res), nil
}

func std_floor(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Floor(a))
}

func std_pow(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Pow(a, b))
}

func std_sqrt(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Sqrt(a))
}

func std_hypot(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Hypot(a, b))
}

func std_ceil(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Ceil(a))
}

func std_sin(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Sin(a))
}

func std_cos(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Cos(a))
}

func std_tan(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Tan(a))
}

func std_asin(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Asin(a))
}

func std_acos(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Acos(a))
}

func std_atan(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Atan(a))
}

func std_atan2(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Atan2(a, b))
}

func std_log(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Log(a))
}

func std_log2(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Log2(a))
}

func std_log10(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Log10(a))
}

func std_exp(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	res := math.Exp(f)
	if res == 0 && f > 0 {
		return evaluator.ValueNone, evaluator.OverflowErr
	}
	return makeCheckedNumber(res)
}

func std_mantissa(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	mantissa, _ := math.Frexp(f)
	return makeCheckedNumber(mantissa)
}

func std_exponent(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	_, exponent := math.Frexp(f)
	return makeCheckedNumber(float64(exponent))
}

func std_round(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Round(a))
}

func std_isEven(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	i, _ := math.Modf(f) // Get the integral part of the float
	return evaluator.MakeBool(math.Mod(i, 2) == 0), nil
}

func std_isOdd(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	i, _ := math.Modf(f) // Get the integral part of the float
	return evaluator.MakeBool(math.Mod(i, 2) != 0), nil
}

func std_isInteger(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	_, frac := math.Modf(f) // Get the fraction part of the float
	return evaluator.MakeBool(frac == 0), nil
}

func std_isDecimal(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	_, frac := math.Modf(f) // Get the fraction part of the float
	return evaluator.MakeBool(frac != 0), nil
}

func std_max(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Max(a, b))
}

func std_min(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Min(a, b))
}

func std_abs(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(math.Abs(a))
}

func std_sign(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	var res float64
	if f == 0 {
		res = 0
	} else if f > 0 {
		res = 1
	} else {
		res = -1
	}
	return makeCheckedNumber(res)
}

func std_deg2rad(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(f * (math.Pi / 180.0))
}

func std_rad2deg(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	f, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return makeCheckedNumber(f * (180.0 / math.Pi))
}

func std_modulo(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].EvalNumber(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if b == 0 {
		return evaluator.ValueNone, evaluator.DivisionByZeroErr
	}
	res := math.Mod(a, b)
	if math.IsNaN(res) {
		return evaluator.ValueNone, evaluator.NanErr
	}
	return evaluator.MakeNumber(res), nil
}

func std_clamp(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	x, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !x.IsNumber() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, x.Type())
	}
	minVal, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !minVal.IsNumber() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, minVal.Type())
	}
	maxVal, err := args[2].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !maxVal.IsNumber() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeNumber, maxVal.Type())
	}

	if x.Number() <= minVal.Number() {
		return minVal, nil
	}

	if x.Number() >= maxVal.Number() {
		return maxVal, nil
	}
	return x, nil
}
