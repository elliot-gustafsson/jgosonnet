package stdlib

import (
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func std_xor(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	x, err := args[0].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	y, err := args[1].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeBool(x != y), nil
}

func std_xnor(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	x, err := args[0].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	y, err := args[1].EvalBool(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeBool(x == y), nil
}
