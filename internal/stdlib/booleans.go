package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func std_xor(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.xor %d, expected 3", len(args))
	}

	x, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	y, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !x.IsBool() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.xor (arg 0): %s, expected boolean", x.Type().String())
	}

	if !y.IsBool() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.xor (arg 1): %s, expected boolean", x.Type().String())
	}

	return evaluator.MakeBool(x.Bool() != y.Bool()), nil
}

func std_xnor(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected number of args passed to std.xnor %d, expected 3", len(args))
	}

	x, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	y, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !x.IsBool() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.xnor (arg 0): %s, expected boolean", x.Type().String())
	}

	if !y.IsBool() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.xnor (arg 1): %s, expected boolean", x.Type().String())
	}

	return evaluator.MakeBool(x.Bool() == y.Bool()), nil
}
