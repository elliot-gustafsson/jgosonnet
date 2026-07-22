package evaluator

import (
	"fmt"
)

type Func = func(args []NamedValue, ctx Context) (Value, error)

var BuiltinFunctions []BuiltinFunction

type BuiltinFunction struct {
	Func     Func
	Params   []string
	OptStart int
}

func (t BuiltinFunction) Noop() bool {
	return t.Func == nil
}

func (t BuiltinFunction) Length() int {
	return t.OptStart // Length should be the number of required args
}

func (t BuiltinFunction) Exec(args []NamedValue, ctx Context) (Value, error) {

	paramCount := len(t.Params)

	if len(args) > paramCount {
		return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", len(t.Params), len(args)))
	}

	var stackArgs [4]NamedValue // Max args for any stdlib function is 4
	var orderedArgs []NamedValue

	if paramCount <= len(stackArgs) {
		orderedArgs = stackArgs[:paramCount]
	} else {
		// Just in case a user defines a custom native extension with 5+ args
		orderedArgs = make([]NamedValue, paramCount)
	}

	var onNamedArgs bool
	posIdx := 0

	for _, na := range args {
		if na.Key == 0 { // Positional Argument
			if onNamedArgs {
				return ValueNone, fmt.Errorf("positional argument after a named argument is not allowed")
			}
			orderedArgs[posIdx] = na
			posIdx++
			continue
		}

		// Named Argument
		onNamedArgs = true
		found := false

		// Look up the string once from the caller's interner
		passedName := ctx.State.Interner.Get(na.Key)

		for j, paramName := range t.Params {
			if passedName == paramName { // Raw Go string comparison
				if !orderedArgs[j].Value.IsNone() {
					argName := ctx.State.Interner.Get(na.Key)
					return ValueNone, MakeRuntimeError(fmt.Errorf("argument '%s' provided multiple times", argName))
				}
				orderedArgs[j] = na
				found = true
				break
			}
		}
		if !found {
			return ValueNone, MakeRuntimeError(fmt.Errorf("function has no parameter %s", passedName))
		}
	}

	// Check required arguments
	for i := 0; i < t.OptStart; i++ {
		if orderedArgs[i].Value.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("missing argument: %s", t.Params[i]))
		}
	}

	return t.Func(orderedArgs, ctx)
}
