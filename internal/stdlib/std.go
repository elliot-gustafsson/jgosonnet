package stdlib

import (
	"fmt"
	"math"
	"os"
	"unicode/utf8"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

var constants = map[string]evaluator.Value{
	"pi": evaluator.MakeNumber(math.Pi),
}

var functions = map[string]func(evaluator.Context) evaluator.Function{
	// --- General ---
	"$flatMapArray":    f(builtin_flatMapArray, "func", "arr"),
	"$objectFlatMerge": f(builtin_objectFlatMerge, "arr"),
	"trace":            f(std_trace, "str", "rest"),
	"assertEqual":      f(std_assertEqual, "a", "b"),
	"toString":         f(std_toString, "a"),
	"length":           f(std_length, "x"),
	"mod":              f(std_mod, "a", "b"),

	// --- Types ---
	"type":       f(std_type, "x"),
	"isString":   f(std_isString, "v"),
	"isNumber":   f(std_isNumber, "v"),
	"isBoolean":  f(std_isBoolean, "v"),
	"isObject":   f(std_isObject, "v"),
	"isArray":    f(std_isArray, "v"),
	"isFunction": f(std_isFunction, "v"),
	"isNull":     f(std_isNull, "x"),
	"prune":      f(std_prune, "a"),

	// --- Parse ---
	"parseInt":   f(std_parseInt, "str"),
	"parseOctal": f(std_parseOctal, "str"),
	"parseHex":   f(std_parseHex, "str"),
	"parseJson":  f(std_parseJson, "str"),
	"parseYaml":  f(std_parseYaml, "str"),

	// --- Math ---
	"floor":     f(std_floor, "x"),
	"ceil":      f(std_ceil, "x"),
	"round":     f(std_round, "x"),
	"pow":       f(std_pow, "x", "n"),
	"sqrt":      f(std_sqrt, "x"),
	"hypot":     f(std_hypot, "a", "b"),
	"modulo":    f(std_modulo, "a", "b"),
	"mantissa":  f(std_mantissa, "x"),
	"exponent":  f(std_exponent, "x"),
	"sin":       f(std_sin, "x"),
	"cos":       f(std_cos, "x"),
	"tan":       f(std_tan, "x"),
	"asin":      f(std_asin, "x"),
	"acos":      f(std_acos, "x"),
	"atan":      f(std_atan, "x"),
	"atan2":     f(std_atan2, "y", "x"),
	"deg2rad":   f(std_deg2rad, "x"),
	"rad2deg":   f(std_rad2deg, "x"),
	"log":       f(std_log, "x"),
	"log2":      f(std_log2, "x"),
	"log10":     f(std_log10, "x"),
	"exp":       f(std_exp, "x"),
	"isEven":    f(std_isEven, "x"),
	"isOdd":     f(std_isOdd, "x"),
	"isInteger": f(std_isInteger, "x"),
	"isDecimal": f(std_isDecimal, "x"),
	"max":       f(std_max, "a", "b"),
	"min":       f(std_min, "a", "b"),
	"abs":       f(std_abs, "n"),
	"sign":      f(std_sign, "n"),
	"clamp":     f(std_clamp, "x", "minVal", "maxVal"),

	// --- Strings ---
	"format":      f(std_format, "str", "vals"),
	"stringChars": f(std_stringChars, "str"),
	"startsWith":  f(std_startsWith, "a", "b"),
	"endsWith":    f(std_endsWith, "a", "b"),
	"substr":      f(std_substr, "str", "from", "len"),
	"findSubstr":  f(std_findSubstr, "pat", "str"),
	"strReplace":  f(std_strReplace, "str", "from", "to"),
	"split":       f(std_split, "str", "c"),
	"splitLimit":  f(std_splitLimit, "str", "c", "maxsplits"),
	"stripChars":  f(std_stripChars, "str", "chars"),
	"rstripChars": f(std_rstripChars, "str", "chars"),
	"lstripChars": f(std_lstripChars, "str", "chars"),
	"isEmpty":     f(std_isEmpty, "str"),
	"trim":        f(std_trim, "str"),
	"md5":         f(std_md5, "s"),
	"sha1":        f(std_sha1, "s"),
	"sha256":      f(std_sha256, "s"),
	"sha512":      f(std_sha512, "s"),
	"sha3":        f(std_sha3, "s"),
	"char":        f(std_char, "n"),
	"codepoint":   f(std_codepoint, "str"),
	"base64":      f(std_base64, "input"),
	"asciiLower":  f(std_asciiLower, "str"),
	"asciiUpper":  f(std_asciiUpper, "str"),

	// --- Arrays ---
	"join":             f(std_join, "sep", "arr"),
	"range":            f(std_range, "from", "to"),
	"makeArray":        f(std_makeArray, "sz", "func"),
	"filter":           f(std_filter, "func", "arr"),
	"uniq":             f(std_uniq, "arr", "keyF"),
	"sort":             f(std_sort, "arr", "keyF"),
	"map":              f(std_map, "func", "arr"),
	"mapWithIndex":     f(std_mapWithIndex, "func", "arr"),
	"flatMap":          f(std_flatMap, "func", "arr"),
	"filterMap":        f(std_filterMap, "filter_func", "map_func", "arr"),
	"member":           f(std_member, "arr", "x"),
	"setMember":        f(std_setMember, "x", "arr", "keyF"),
	"slice":            f(std_slice, "indexable", "index", "end", "step"),
	"count":            f(std_count, "arr", "x"),
	"lines":            f(std_lines, "arr"),
	"reverse":          f(std_reverse, "arrs"),
	"foldl":            f(std_foldl, "func", "arr", "init"),
	"foldr":            f(std_foldr, "func", "arr", "init"),
	"sum":              f(std_sum, "arr"),
	"flattenArrays":    f(std_flattenArrays, "arr"),
	"flattenDeepArray": f(std_flattenDeepArray, "arr"),
	"repeat":           f(std_repeat, "what", "count"),

	// -- Sets ---
	"set": f(std_set, "arr", "keyF"),

	// --- Objects ---
	"get":                 f(std_get, "o", "f", "default", "inc_hidden"),
	"objectFields":        f(std_objectFields, "o"),
	"objectFieldsAll":     f(std_objectFieldsAll, "o"),
	"objectHas":           f(std_objectHas, "o", "f"),
	"objectHasAll":        f(std_objectHasAll, "o", "f"),
	"objectValues":        f(std_objectValues, "o"),
	"objectValuesAll":     f(std_objectValuesAll, "o"),
	"objectKeysValues":    f(std_objectKeysValues, "o"),
	"objectKeysValuesAll": f(std_objectKeysValuesAll, "o"),
	"mapWithKey":          f(std_mapWithkey, "func", "obj"),

	// --- Manifestation ---
	"manifestYamlDoc":      f(std_manifestYamlDoc, "value", "indent_array_in_object", "quote_keys"),
	"manifestYamlStream":   f(std_manifestYamlStream, "value", "indent_array_in_object", "c_document_end", "quote_keys"),
	"manifestJson":         f(std_manifestJson, "value"),
	"manifestJsonEx":       f(std_manifestJsonEx, "value", "indent", "newline", "key_val_sep"),
	"manifestJsonMinified": f(std_manifestJsonMinified, "value"),
	"manifestIni":          f(std_manifestIni, "ini"),
}

func f(f evaluator.Func, argn ...string) func(evaluator.Context) evaluator.Function {

	return func(ctx evaluator.Context) evaluator.Function {

		argIds := make([]uint32, 0, len(argn))
		for _, v := range argn {
			id := ctx.Interner.Intern(v)
			argIds = append(argIds, id)
		}

		var fn evaluator.Func = func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

			// TODO: do arg reordering in evaluator.Function.Exec

			var onNamedArgs bool

			orderedArgs := make([]evaluator.NamedValue, len(argIds))
			for i, na := range args {
				if !onNamedArgs && na.Key != 0 {
					onNamedArgs = true
				}

				if na.Key == 0 {
					if onNamedArgs {
						return evaluator.Value{}, fmt.Errorf("positional argument after a named argument is not allowed")
					}
					orderedArgs[i] = na
					continue
				}

				for ii, aid := range argIds {
					if na.Key == aid {
						orderedArgs[ii] = na
					}
				}

			}

			return f(orderedArgs, ctx)
		}

		return evaluator.NewFunction(argIds, fn)
	}
}

func InitStdLib(ctx evaluator.Context) (evaluator.Value, error) {

	fieldCount := len(functions) + len(constants)

	layer := &evaluator.Layer{
		Keys:   make([]uint32, 0, fieldCount),
		Values: make([]evaluator.Value, 0, fieldCount),
		Meta:   make([]uint8, 0, fieldCount),

		Index: make(map[uint32]int, fieldCount),
	}

	obj := evaluator.NewObject([]*evaluator.Layer{layer})

	index := 0
	for name, f := range functions {
		keyId := ctx.Interner.Intern(name)

		fVal := f(ctx)

		v := evaluator.MakeFunction(fVal, ctx)
		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, v)
		layer.Meta = append(layer.Meta, 0)
		layer.Index[keyId] = index

		index++
	}

	for name, v := range constants {
		keyId := ctx.Interner.Intern(name)

		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, v)
		layer.Meta = append(layer.Meta, 0)
		layer.Index[keyId] = index

		index++
	}

	val := evaluator.MakeObject(obj, ctx)

	return val, nil
}

func liftValueToBool(f func(evaluator.NamedValue) bool, name string) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		res := f(evaluator.NamedValue{Value: a})
		return evaluator.MakeBool(res), nil
	}
}

func std_trace(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.trace: %d, expected 2", len(args))
	}

	str, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !str.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.trace (arg 0): %s, expected string", str.Type().String())
	}

	_, err = fmt.Fprint(os.Stdout, "TRACE: "+str.String(ctx))
	if err != nil {
		return evaluator.Value{}, err
	}
	return args[1].Value, nil
}

func std_type(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.type: %d, expected 1", len(args))
	}

	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(v.Type().String(), ctx), nil
}

func std_assertEqual(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.assertEqual: %d, expected 2", len(args))
	}

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	b, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	eq, err := a.Equal(b, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if eq {
		return evaluator.MakeBool(true), nil
	}

	aStr, err := a.ToString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	bStr, err := b.ToString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeBool(false), fmt.Errorf("assertion failed %s != %s", aStr, bStr)
}

var std_isString = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsString() }, "std.isString")
var std_isNumber = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsNumber() }, "std.isNumber")
var std_isBoolean = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsBool() }, "std.isBoolean")
var std_isObject = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsObject() }, "std.isObject")
var std_isArray = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsArray() }, "std.isArray")
var std_isFunction = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsFunction() }, "std.isFunction")
var std_isNull = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsNull() }, "std.isNull")

func std_toString(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.toString: %d, expected 1", len(args))
	}

	s, err := args[0].ToString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(s, ctx), nil
}

func std_length(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.length: %d, expected 1", len(args))
	}

	arg, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	var res float64
	switch arg.Type() {
	case evaluator.ValueTypeString:
		// res = float64(len(arg.String(ctx)))
		res = float64(utf8.RuneCountInString(arg.String(ctx)))
	case evaluator.ValueTypeArray:
		res = float64(len(arg.Array(ctx)))
	case evaluator.ValueTypeObject:
		res = float64(arg.Object(ctx).Length(ctx))
	case evaluator.ValueTypeFunction:
		res = float64(len(arg.Function(ctx).Args))
	default:
		return evaluator.Value{}, fmt.Errorf("std.length: unexpected type %s", arg.Type().String())
	}

	return evaluator.MakeNumber(res), nil
}

func std_mod(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.mod: %d, expected 2", len(args))
	}

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	b, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if a.IsNumber() && b.IsNumber() {
		return std_modulo(args, ctx)
	}

	if a.IsString() {
		x := []evaluator.NamedValue{{Value: a}, {Value: b}}
		return std_format(x, ctx)
	}
	return evaluator.Value{}, fmt.Errorf("'Operator %% cannot be used on types %s and %s", a.Type().String(), b.Type().String())
}

func std_prune(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.prune: %d, expected 1", len(args))
	}
	arg, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res, err := arg.Prune(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return res, nil
}
