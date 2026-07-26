package stdlib

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

type param struct {
	Name     string
	Optional bool
}

func req(name string) param {
	return param{Name: name}
}

func opt(name string) param {
	return param{Name: name, Optional: true}
}

var constants = map[string]evaluator.Value{
	"pi": evaluator.MakeNumber(math.Pi),
}

var functions = map[string]func(evaluator.Context) evaluator.Function{
	// --- General ---
	"$flatMapArray":    f(builtin_flatMapArray, req("func"), req("arr")),
	"$objectFlatMerge": f(builtin_objectFlatMerge, req("arr")),
	"extVar":           f(std_extVar, req("x")),
	"trace":            f(std_trace, req("str"), req("rest")),
	"assertEqual":      f(std_assertEqual, req("a"), req("b")),
	"toString":         f(std_toString, req("a")),
	"length":           f(std_length, req("x")),
	"mod":              f(std_mod, req("a"), req("b")),

	// --- Types ---
	"type":       f(std_type, req("x")),
	"isString":   f(std_isString, req("v")),
	"isNumber":   f(std_isNumber, req("v")),
	"isBoolean":  f(std_isBoolean, req("v")),
	"isObject":   f(std_isObject, req("v")),
	"isArray":    f(std_isArray, req("v")),
	"isFunction": f(std_isFunction, req("v")),
	"isNull":     f(std_isNull, req("x")),
	"prune":      f(std_prune, req("a")),

	// --- Parse ---
	"parseInt":   f(std_parseInt, req("str")),
	"parseOctal": f(std_parseOctal, req("str")),
	"parseHex":   f(std_parseHex, req("str")),
	"parseJson":  f(std_parseJson, req("str")),
	"parseYaml":  f(std_parseYaml, req("str")),
	"encodeUTF8": f(std_encodeUTF8, req("str")),
	"decodeUTF8": f(std_decodeUTF8, req("arr")),

	// --- Math ---
	"floor":     f(std_floor, req("x")),
	"ceil":      f(std_ceil, req("x")),
	"round":     f(std_round, req("x")),
	"pow":       f(std_pow, req("x"), req("n")),
	"sqrt":      f(std_sqrt, req("x")),
	"hypot":     f(std_hypot, req("a"), req("b")),
	"modulo":    f(std_modulo, req("a"), req("b")),
	"mantissa":  f(std_mantissa, req("x")),
	"exponent":  f(std_exponent, req("x")),
	"sin":       f(std_sin, req("x")),
	"cos":       f(std_cos, req("x")),
	"tan":       f(std_tan, req("x")),
	"asin":      f(std_asin, req("x")),
	"acos":      f(std_acos, req("x")),
	"atan":      f(std_atan, req("x")),
	"atan2":     f(std_atan2, req("y"), req("x")),
	"deg2rad":   f(std_deg2rad, req("x")),
	"rad2deg":   f(std_rad2deg, req("x")),
	"log":       f(std_log, req("x")),
	"log2":      f(std_log2, req("x")),
	"log10":     f(std_log10, req("x")),
	"exp":       f(std_exp, req("x")),
	"isEven":    f(std_isEven, req("x")),
	"isOdd":     f(std_isOdd, req("x")),
	"isInteger": f(std_isInteger, req("x")),
	"isDecimal": f(std_isDecimal, req("x")),
	"max":       f(std_max, req("a"), req("b")),
	"min":       f(std_min, req("a"), req("b")),
	"abs":       f(std_abs, req("n")),
	"sign":      f(std_sign, req("n")),
	"clamp":     f(std_clamp, req("x"), req("minVal"), req("maxVal")),

	// --- Strings ---
	"format":              f(std_format, req("str"), req("vals")),
	"stringChars":         f(std_stringChars, req("str")),
	"startsWith":          f(std_startsWith, req("a"), req("b")),
	"endsWith":            f(std_endsWith, req("a"), req("b")),
	"substr":              f(std_substr, req("str"), req("from"), req("len")),
	"findSubstr":          f(std_findSubstr, req("pat"), req("str")),
	"strReplace":          f(std_strReplace, req("str"), req("from"), req("to")),
	"split":               f(std_split, req("str"), req("c")),
	"splitLimit":          f(std_splitLimit, req("str"), req("c"), req("maxsplits")),
	"splitLimitR":         f(std_splitLimitR, req("str"), req("c"), req("maxsplits")),
	"stripChars":          f(std_stripChars, req("str"), req("chars")),
	"rstripChars":         f(std_rstripChars, req("str"), req("chars")),
	"lstripChars":         f(std_lstripChars, req("str"), req("chars")),
	"isEmpty":             f(std_isEmpty, req("str")),
	"trim":                f(std_trim, req("str")),
	"md5":                 f(std_md5, req("s")),
	"sha1":                f(std_sha1, req("s")),
	"sha256":              f(std_sha256, req("s")),
	"sha512":              f(std_sha512, req("s")),
	"sha3":                f(std_sha3, req("s")),
	"char":                f(std_char, req("n")),
	"codepoint":           f(std_codepoint, req("str")),
	"base64":              f(std_base64, req("input")),
	"base64Decode":        f(std_base64Decode, req("str")),
	"base64DecodeBytes":   f(std_base64DecodeBytes, req("str")),
	"asciiLower":          f(std_asciiLower, req("str")),
	"asciiUpper":          f(std_asciiUpper, req("str")),
	"escapeStringBash":    f(std_escapeStringBash, req("str_")),
	"escapeStringDollars": f(std_escapeStringDollars, req("str_")),
	"escapeStringJson":    f(std_escapeStringJson, req("str_")),
	"escapeStringPython":  f(std_escapeStringJson, req("str")), // Intentionally same as function as escapeStringJson
	"escapeStringXML":     f(std_escapeStringXML, req("str_")),
	"equalsIgnoreCase":    f(std_equalsIgnoreCase, req("str1"), req("str2")),

	// --- Arrays ---
	"join":             f(std_join, req("sep"), req("arr")),
	"deepJoin":         f(std_deepJoin, req("arr")),
	"range":            f(std_range, req("from"), req("to")),
	"makeArray":        f(std_makeArray, req("sz"), req("func")),
	"filter":           f(std_filter, req("func"), req("arr")),
	"uniq":             f(std_uniq, req("arr"), opt("keyF")),
	"sort":             f(std_sort, req("arr"), opt("keyF")),
	"map":              f(std_map, req("func"), req("arr")),
	"mapWithIndex":     f(std_mapWithIndex, req("func"), req("arr")),
	"flatMap":          f(std_flatMap, req("func"), req("arr")),
	"filterMap":        f(std_filterMap, req("filter_func"), req("map_func"), req("arr")),
	"member":           f(std_member, req("arr"), req("x")),
	"setMember":        f(std_setMember, req("x"), req("arr"), opt("keyF")),
	"slice":            f(std_slice, req("indexable"), opt("index"), opt("end"), opt("step")),
	"count":            f(std_count, req("arr"), req("x")),
	"lines":            f(std_lines, req("arr")),
	"reverse":          f(std_reverse, req("arrs")),
	"foldl":            f(std_foldl, req("func"), req("arr"), req("init")),
	"foldr":            f(std_foldr, req("func"), req("arr"), req("init")),
	"sum":              f(std_sum, req("arr")),
	"flattenArrays":    f(std_flattenArrays, req("arr")),
	"flattenDeepArray": f(std_flattenDeepArray, req("arr")),
	"repeat":           f(std_repeat, req("what"), req("count")),
	"setUnion":         f(std_setUnion, req("a"), req("b"), opt("keyF")),
	"setInter":         f(std_setInter, req("a"), req("b"), opt("keyF")),
	"setDiff":          f(std_setDiff, req("a"), req("b"), opt("keyF")),
	"find":             f(std_find, req("value"), req("arr")),
	"any":              f(std_any, req("arr")),
	"all":              f(std_all, req("arr")),
	"avg":              f(std_avg, req("arr")),
	"minArray":         f(std_minArray, req("arr"), opt("keyF"), opt("onEmpty")),
	"maxArray":         f(std_maxArray, req("arr"), opt("keyF"), opt("onEmpty")),
	"contains":         f(std_contains, req("arr"), req("elem")),
	"remove":           f(std_remove, req("arr"), req("elem")),
	"removeAt":         f(std_removeAt, req("arr"), req("idx")),

	// -- Booleans ---
	"xor":  f(std_xor, req("x"), req("y")),
	"xnor": f(std_xnor, req("x"), req("y")),

	// -- Sets ---
	"set": f(std_set, req("arr"), opt("keyF")),

	// --- Objects ---
	"get":                 f(std_get, req("o"), req("f"), opt("default"), opt("inc_hidden")),
	"objectFields":        f(std_objectFields, req("o")),
	"objectFieldsAll":     f(std_objectFieldsAll, req("o")),
	"objectFieldsEx":      f(std_objectFieldsEx, req("obj"), req("hidden")),
	"objectHas":           f(std_objectHas, req("o"), req("f")),
	"objectHasAll":        f(std_objectHasAll, req("o"), req("f")),
	"objectHasEx":         f(std_objectHasEx, req("obj"), req("fname"), req("hidden")),
	"objectValues":        f(std_objectValues, req("o")),
	"objectValuesAll":     f(std_objectValuesAll, req("o")),
	"objectKeysValues":    f(std_objectKeysValues, req("o")),
	"objectKeysValuesAll": f(std_objectKeysValuesAll, req("o")),
	"mapWithKey":          f(std_mapWithkey, req("func"), req("obj")),
	"objectRemoveKey":     f(std_objectRemoveKey, req("obj"), req("key")),
	"mergePatch":          f(std_mergePatch, req("target"), req("patch")),

	// --- Manifestation ---
	"manifestYamlDoc":      f(std_manifestYamlDoc, req("value"), opt("indent_array_in_object"), opt("quote_keys")),
	"manifestYamlStream":   f(std_manifestYamlStream, req("value"), opt("indent_array_in_object"), opt("c_document_end"), opt("quote_keys")),
	"manifestJson":         f(std_manifestJson, req("value")),
	"manifestJsonEx":       f(std_manifestJsonEx, req("value"), req("indent"), opt("newline"), opt("key_val_sep")),
	"manifestJsonMinified": f(std_manifestJsonMinified, req("value")),
	"manifestIni":          f(std_manifestIni, req("ini")),
	"manifestPython":       f(std_manifestPython, req("v")),
	"manifestPythonVars":   f(std_manifestPythonVars, req("conf")),
	"manifestXmlJsonml":    f(std_manifestXmlJsonml, req("value")),
	"manifestTomlEx":       f(std_manifestTomlEx, req("value"), req("indent")),
}

func f(f evaluator.Func, params ...param) func(evaluator.Context) evaluator.Function {

	return func(ctx evaluator.Context) evaluator.Function {

		argIds := make([]uint32, 0, len(params))
		optStart := -1

		for i, p := range params {
			id := ctx.State.Interner.Intern(p.Name)
			argIds = append(argIds, id)

			if p.Optional && optStart == -1 {
				optStart = i
			}
		}
		if optStart == -1 {
			optStart = len(params)
		}

		var fn evaluator.Func
		fn = func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

			if len(args) > len(argIds) {
				return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", len(argIds), len(args)))
			}

			var onNamedArgs bool

			orderedArgs := ctx.State.Registry.NamedValueBufs.Alloc(len(argIds), len(argIds))
			posIdx := 0

			for _, na := range args {
				if na.Key == 0 {
					// Positional Argument
					if onNamedArgs {
						return evaluator.ValueNone, fmt.Errorf("Positional argument after a named argument is not allowed")
					}
					na.Key = argIds[posIdx]
					orderedArgs[posIdx] = na
					posIdx++
					continue
				}

				// Named Argument
				onNamedArgs = true
				found := false
				for j, aid := range argIds {
					if aid == na.Key {
						if !orderedArgs[j].Value.IsNone() {
							argName := ctx.State.Interner.Get(na.Key)
							return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Argument %s already provided", argName))
						}
						orderedArgs[j] = na
						found = true
						break
					}
				}
				if !found {
					argName := ctx.State.Interner.Get(na.Key)
					return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("function has no parameter %s", argName))
				}
			}

			for i := 0; i < optStart; i++ {
				if orderedArgs[i].Value.IsNone() {
					return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Missing argument: %s", params[i].Name))
				}
			}

			return f(orderedArgs, ctx)
		}

		return evaluator.NewFunction(len(argIds), fn)
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

	objId := evaluator.NewObject([]*evaluator.Layer{layer}, ctx)

	index := 0
	for name, f := range functions {
		keyId := ctx.State.Interner.Intern(name)

		fVal := f(ctx)

		v := evaluator.MakeFunction(fVal, ctx)
		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, v)
		layer.Meta = append(layer.Meta, 0)
		layer.Index[keyId] = index

		index++
	}

	for name, v := range constants {
		keyId := ctx.State.Interner.Intern(name)

		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, v)
		layer.Meta = append(layer.Meta, 0)
		layer.Index[keyId] = index

		index++
	}

	val := evaluator.MakeObjectValue(objId)

	return val, nil
}

func liftValueToBool(f func(evaluator.NamedValue) bool) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
		a, err := args[0].Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		res := f(evaluator.NamedValue{Value: a})
		return evaluator.MakeBool(res), nil
	}
}

func std_extVar(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	name, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	s, ok := ctx.State.Environment.ExtVars[name]
	if ok {
		return evaluator.MakeString(s, ctx), nil
	}

	s, ok = ctx.State.Environment.ExtCodes[name]
	if ok {
		name := "<extvar:" + name + ">"

		importer := ctx.State.Environment.Importer

		val := ctx.State.Environment.Importer.Get(name)
		if !val.IsNone() {
			return val, nil
		}

		n, err := importer.ResolveSnippet(name, s)
		if err != nil {
			return evaluator.ValueNone, err
		}

		importCtx := ctx
		importCtx.Self = evaluator.ValueNone
		importCtx.SuperOffset = 0

		scopeId := evaluator.CreateFileScope(name, importer.BaseStd, importCtx)

		val, err = evaluator.EvaluateNode(n, scopeId, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		importer.Set(name, val)

		return val, nil
	}

	return evaluator.ValueNone, fmt.Errorf("undefined external variable: %s", name)
}

func std_trace(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	str, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	_, err = fmt.Fprint(ctx.State.Environment.TraceOut, "TRACE: "+str+"\n")
	if err != nil {
		return evaluator.ValueNone, err
	}
	return args[1].Value, nil
}

func std_type(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(v.Type().String(), ctx), nil
}

func std_assertEqual(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.ValueNone, fmt.Errorf("unexpected amount of arguments passed to std.assertEqual: %d, expected 2", len(args))
	}

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	eq, err := a.Equal(b, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	if eq {
		return evaluator.MakeBool(true), nil
	}

	aStr, err := a.ToString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	bStr, err := b.ToString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeBool(false), fmt.Errorf("assertion failed %s != %s", aStr, bStr)
}

var std_isString = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsString() })
var std_isNumber = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsNumber() })
var std_isBoolean = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsBool() })
var std_isObject = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsObject() })
var std_isArray = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsArray() })
var std_isFunction = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsFunction() })
var std_isNull = liftValueToBool(func(v evaluator.NamedValue) bool { return v.IsNull() })

func std_toString(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	s, err := args[0].ToString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeString(s, ctx), nil
}

func std_length(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arg, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
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
		res = float64(arg.Function(ctx).Length())
	default:
		return evaluator.ValueNone, evaluator.TypeErrorGeneral(arg.Type())
	}

	return evaluator.MakeNumber(res), nil
}

func std_mod(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	b, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if a.IsNumber() && b.IsNumber() {
		return std_modulo(args, ctx)
	}

	if a.IsString() {
		x := []evaluator.NamedValue{{Value: a}, {Value: b}}
		return std_format(x, ctx)
	}
	return evaluator.ValueNone, evaluator.MakeRuntimeError(fmt.Errorf("Operator %% cannot be used on types %s and %s.", a.Type().String(), b.Type().String()))
}

func std_prune(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	arg, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	res, err := arg.Prune(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return res, nil
}
