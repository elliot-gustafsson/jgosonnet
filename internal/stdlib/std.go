package stdlib

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

var constants = []struct {
	name string
	val  evaluator.Value
}{
	{"pi", evaluator.MakeNumber(math.Pi)},
}

var functions = []struct {
	name string
	fn   evaluator.BuiltinFunction
}{
	{"$flatMapArray", evaluator.BuiltinFunction{Func: builtin_flatMapArray, Params: []string{"func", "arr"}, OptStart: 2}},
	{"$objectFlatMerge", evaluator.BuiltinFunction{Func: builtin_objectFlatMerge, Params: []string{"arr"}, OptStart: 1}},
	{"extVar", evaluator.BuiltinFunction{Func: std_extVar, Params: []string{"x"}, OptStart: 1}},
	{"trace", evaluator.BuiltinFunction{Func: std_trace, Params: []string{"str", "rest"}, OptStart: 2}},
	{"assertEqual", evaluator.BuiltinFunction{Func: std_assertEqual, Params: []string{"a", "b"}, OptStart: 2}},
	{"toString", evaluator.BuiltinFunction{Func: std_toString, Params: []string{"a"}, OptStart: 1}},
	{"length", evaluator.BuiltinFunction{Func: std_length, Params: []string{"x"}, OptStart: 1}},
	{"mod", evaluator.BuiltinFunction{Func: std_mod, Params: []string{"a", "b"}, OptStart: 2}},

	// --- Types ---
	{"type", evaluator.BuiltinFunction{Func: std_type, Params: []string{"x"}, OptStart: 1}},
	{"isString", evaluator.BuiltinFunction{Func: std_isString, Params: []string{"v"}, OptStart: 1}},
	{"isNumber", evaluator.BuiltinFunction{Func: std_isNumber, Params: []string{"v"}, OptStart: 1}},
	{"isBoolean", evaluator.BuiltinFunction{Func: std_isBoolean, Params: []string{"v"}, OptStart: 1}},
	{"isObject", evaluator.BuiltinFunction{Func: std_isObject, Params: []string{"v"}, OptStart: 1}},
	{"isArray", evaluator.BuiltinFunction{Func: std_isArray, Params: []string{"v"}, OptStart: 1}},
	{"isFunction", evaluator.BuiltinFunction{Func: std_isFunction, Params: []string{"v"}, OptStart: 1}},
	{"isNull", evaluator.BuiltinFunction{Func: std_isNull, Params: []string{"x"}, OptStart: 1}},
	{"prune", evaluator.BuiltinFunction{Func: std_prune, Params: []string{"a"}, OptStart: 1}},

	// --- Parse ---
	{"parseInt", evaluator.BuiltinFunction{Func: std_parseInt, Params: []string{"str"}, OptStart: 1}},
	{"parseOctal", evaluator.BuiltinFunction{Func: std_parseOctal, Params: []string{"str"}, OptStart: 1}},
	{"parseHex", evaluator.BuiltinFunction{Func: std_parseHex, Params: []string{"str"}, OptStart: 1}},
	{"parseJson", evaluator.BuiltinFunction{Func: std_parseJson, Params: []string{"str"}, OptStart: 1}},
	{"parseYaml", evaluator.BuiltinFunction{Func: std_parseYaml, Params: []string{"str"}, OptStart: 1}},
	{"encodeUTF8", evaluator.BuiltinFunction{Func: std_encodeUTF8, Params: []string{"str"}, OptStart: 1}},
	{"decodeUTF8", evaluator.BuiltinFunction{Func: std_decodeUTF8, Params: []string{"arr"}, OptStart: 1}},

	// --- Math ---
	{"floor", evaluator.BuiltinFunction{Func: std_floor, Params: []string{"x"}, OptStart: 1}},
	{"ceil", evaluator.BuiltinFunction{Func: std_ceil, Params: []string{"x"}, OptStart: 1}},
	{"round", evaluator.BuiltinFunction{Func: std_round, Params: []string{"x"}, OptStart: 1}},
	{"pow", evaluator.BuiltinFunction{Func: std_pow, Params: []string{"x", "n"}, OptStart: 2}},
	{"sqrt", evaluator.BuiltinFunction{Func: std_sqrt, Params: []string{"x"}, OptStart: 1}},
	{"hypot", evaluator.BuiltinFunction{Func: std_hypot, Params: []string{"a", "b"}, OptStart: 2}},
	{"modulo", evaluator.BuiltinFunction{Func: std_modulo, Params: []string{"a", "b"}, OptStart: 2}},
	{"mantissa", evaluator.BuiltinFunction{Func: std_mantissa, Params: []string{"x"}, OptStart: 1}},
	{"exponent", evaluator.BuiltinFunction{Func: std_exponent, Params: []string{"x"}, OptStart: 1}},
	{"sin", evaluator.BuiltinFunction{Func: std_sin, Params: []string{"x"}, OptStart: 1}},
	{"cos", evaluator.BuiltinFunction{Func: std_cos, Params: []string{"x"}, OptStart: 1}},
	{"tan", evaluator.BuiltinFunction{Func: std_tan, Params: []string{"x"}, OptStart: 1}},
	{"asin", evaluator.BuiltinFunction{Func: std_asin, Params: []string{"x"}, OptStart: 1}},
	{"acos", evaluator.BuiltinFunction{Func: std_acos, Params: []string{"x"}, OptStart: 1}},
	{"atan", evaluator.BuiltinFunction{Func: std_atan, Params: []string{"x"}, OptStart: 1}},
	{"atan2", evaluator.BuiltinFunction{Func: std_atan2, Params: []string{"y", "x"}, OptStart: 2}},
	{"deg2rad", evaluator.BuiltinFunction{Func: std_deg2rad, Params: []string{"x"}, OptStart: 1}},
	{"rad2deg", evaluator.BuiltinFunction{Func: std_rad2deg, Params: []string{"x"}, OptStart: 1}},
	{"log", evaluator.BuiltinFunction{Func: std_log, Params: []string{"x"}, OptStart: 1}},
	{"log2", evaluator.BuiltinFunction{Func: std_log2, Params: []string{"x"}, OptStart: 1}},
	{"log10", evaluator.BuiltinFunction{Func: std_log10, Params: []string{"x"}, OptStart: 1}},
	{"exp", evaluator.BuiltinFunction{Func: std_exp, Params: []string{"x"}, OptStart: 1}},
	{"isEven", evaluator.BuiltinFunction{Func: std_isEven, Params: []string{"x"}, OptStart: 1}},
	{"isOdd", evaluator.BuiltinFunction{Func: std_isOdd, Params: []string{"x"}, OptStart: 1}},
	{"isInteger", evaluator.BuiltinFunction{Func: std_isInteger, Params: []string{"x"}, OptStart: 1}},
	{"isDecimal", evaluator.BuiltinFunction{Func: std_isDecimal, Params: []string{"x"}, OptStart: 1}},
	{"max", evaluator.BuiltinFunction{Func: std_max, Params: []string{"a", "b"}, OptStart: 2}},
	{"min", evaluator.BuiltinFunction{Func: std_min, Params: []string{"a", "b"}, OptStart: 2}},
	{"abs", evaluator.BuiltinFunction{Func: std_abs, Params: []string{"n"}, OptStart: 1}},
	{"sign", evaluator.BuiltinFunction{Func: std_sign, Params: []string{"n"}, OptStart: 1}},
	{"clamp", evaluator.BuiltinFunction{Func: std_clamp, Params: []string{"x", "minVal", "maxVal"}, OptStart: 3}},

	// --- Strings ---
	{"format", evaluator.BuiltinFunction{Func: std_format, Params: []string{"str", "vals"}, OptStart: 2}},
	{"stringChars", evaluator.BuiltinFunction{Func: std_stringChars, Params: []string{"str"}, OptStart: 1}},
	{"startsWith", evaluator.BuiltinFunction{Func: std_startsWith, Params: []string{"a", "b"}, OptStart: 2}},
	{"endsWith", evaluator.BuiltinFunction{Func: std_endsWith, Params: []string{"a", "b"}, OptStart: 2}},
	{"substr", evaluator.BuiltinFunction{Func: std_substr, Params: []string{"str", "from", "len"}, OptStart: 3}},
	{"findSubstr", evaluator.BuiltinFunction{Func: std_findSubstr, Params: []string{"pat", "str"}, OptStart: 2}},
	{"strReplace", evaluator.BuiltinFunction{Func: std_strReplace, Params: []string{"str", "from", "to"}, OptStart: 3}},
	{"split", evaluator.BuiltinFunction{Func: std_split, Params: []string{"str", "c"}, OptStart: 2}},
	{"splitLimit", evaluator.BuiltinFunction{Func: std_splitLimit, Params: []string{"str", "c", "maxsplits"}, OptStart: 3}},
	{"splitLimitR", evaluator.BuiltinFunction{Func: std_splitLimitR, Params: []string{"str", "c", "maxsplits"}, OptStart: 3}},
	{"stripChars", evaluator.BuiltinFunction{Func: std_stripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"rstripChars", evaluator.BuiltinFunction{Func: std_rstripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"lstripChars", evaluator.BuiltinFunction{Func: std_lstripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"isEmpty", evaluator.BuiltinFunction{Func: std_isEmpty, Params: []string{"str"}, OptStart: 1}},
	{"trim", evaluator.BuiltinFunction{Func: std_trim, Params: []string{"str"}, OptStart: 1}},
	{"md5", evaluator.BuiltinFunction{Func: std_md5, Params: []string{"s"}, OptStart: 1}},
	{"sha1", evaluator.BuiltinFunction{Func: std_sha1, Params: []string{"s"}, OptStart: 1}},
	{"sha256", evaluator.BuiltinFunction{Func: std_sha256, Params: []string{"s"}, OptStart: 1}},
	{"sha512", evaluator.BuiltinFunction{Func: std_sha512, Params: []string{"s"}, OptStart: 1}},
	{"sha3", evaluator.BuiltinFunction{Func: std_sha3, Params: []string{"s"}, OptStart: 1}},
	{"char", evaluator.BuiltinFunction{Func: std_char, Params: []string{"n"}, OptStart: 1}},
	{"codepoint", evaluator.BuiltinFunction{Func: std_codepoint, Params: []string{"str"}, OptStart: 1}},
	{"base64", evaluator.BuiltinFunction{Func: std_base64, Params: []string{"input"}, OptStart: 1}},
	{"base64Decode", evaluator.BuiltinFunction{Func: std_base64Decode, Params: []string{"str"}, OptStart: 1}},
	{"base64DecodeBytes", evaluator.BuiltinFunction{Func: std_base64DecodeBytes, Params: []string{"str"}, OptStart: 1}},
	{"asciiLower", evaluator.BuiltinFunction{Func: std_asciiLower, Params: []string{"str"}, OptStart: 1}},
	{"asciiUpper", evaluator.BuiltinFunction{Func: std_asciiUpper, Params: []string{"str"}, OptStart: 1}},
	{"escapeStringBash", evaluator.BuiltinFunction{Func: std_escapeStringBash, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringDollars", evaluator.BuiltinFunction{Func: std_escapeStringDollars, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringJson", evaluator.BuiltinFunction{Func: std_escapeStringJson, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringPython", evaluator.BuiltinFunction{Func: std_escapeStringJson, Params: []string{"str"}, OptStart: 1}}, // Intentionally same as function as escapeStringJson
	{"escapeStringXML", evaluator.BuiltinFunction{Func: std_escapeStringXML, Params: []string{"str_"}, OptStart: 1}},
	{"equalsIgnoreCase", evaluator.BuiltinFunction{Func: std_equalsIgnoreCase, Params: []string{"str1", "str2"}, OptStart: 2}},

	// --- Arrays ---
	{"join", evaluator.BuiltinFunction{Func: std_join, Params: []string{"sep", "arr"}, OptStart: 2}},
	{"deepJoin", evaluator.BuiltinFunction{Func: std_deepJoin, Params: []string{"arr"}, OptStart: 1}},
	{"range", evaluator.BuiltinFunction{Func: std_range, Params: []string{"from", "to"}, OptStart: 2}},
	{"makeArray", evaluator.BuiltinFunction{Func: std_makeArray, Params: []string{"sz", "func"}, OptStart: 2}},
	{"filter", evaluator.BuiltinFunction{Func: std_filter, Params: []string{"func", "arr"}, OptStart: 2}},
	{"uniq", evaluator.BuiltinFunction{Func: std_uniq, Params: []string{"arr", "keyF"}, OptStart: 1}},
	{"sort", evaluator.BuiltinFunction{Func: std_sort, Params: []string{"arr", "keyF"}, OptStart: 1}},
	{"map", evaluator.BuiltinFunction{Func: std_map, Params: []string{"func", "arr"}, OptStart: 2}},
	{"mapWithIndex", evaluator.BuiltinFunction{Func: std_mapWithIndex, Params: []string{"func", "arr"}, OptStart: 2}},
	{"flatMap", evaluator.BuiltinFunction{Func: std_flatMap, Params: []string{"func", "arr"}, OptStart: 2}},
	{"filterMap", evaluator.BuiltinFunction{Func: std_filterMap, Params: []string{"filter_func", "map_func", "arr"}, OptStart: 3}},
	{"member", evaluator.BuiltinFunction{Func: std_member, Params: []string{"arr", "x"}, OptStart: 2}},
	{"setMember", evaluator.BuiltinFunction{Func: std_setMember, Params: []string{"x", "arr", "keyF"}, OptStart: 2}},
	{"slice", evaluator.BuiltinFunction{Func: std_slice, Params: []string{"indexable", "index", "end", "step"}, OptStart: 1}},
	{"count", evaluator.BuiltinFunction{Func: std_count, Params: []string{"arr", "x"}, OptStart: 2}},
	{"lines", evaluator.BuiltinFunction{Func: std_lines, Params: []string{"arr"}, OptStart: 1}},
	{"reverse", evaluator.BuiltinFunction{Func: std_reverse, Params: []string{"arrs"}, OptStart: 1}},
	{"foldl", evaluator.BuiltinFunction{Func: std_foldl, Params: []string{"func", "arr", "init"}, OptStart: 3}},
	{"foldr", evaluator.BuiltinFunction{Func: std_foldr, Params: []string{"func", "arr", "init"}, OptStart: 3}},
	{"sum", evaluator.BuiltinFunction{Func: std_sum, Params: []string{"arr"}, OptStart: 1}},
	{"flattenArrays", evaluator.BuiltinFunction{Func: std_flattenArrays, Params: []string{"arr"}, OptStart: 1}},
	{"flattenDeepArray", evaluator.BuiltinFunction{Func: std_flattenDeepArray, Params: []string{"arr"}, OptStart: 1}},
	{"repeat", evaluator.BuiltinFunction{Func: std_repeat, Params: []string{"what", "count"}, OptStart: 2}},
	{"setUnion", evaluator.BuiltinFunction{Func: std_setUnion, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"setInter", evaluator.BuiltinFunction{Func: std_setInter, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"setDiff", evaluator.BuiltinFunction{Func: std_setDiff, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"find", evaluator.BuiltinFunction{Func: std_find, Params: []string{"value", "arr"}, OptStart: 2}},
	{"any", evaluator.BuiltinFunction{Func: std_any, Params: []string{"arr"}, OptStart: 1}},
	{"all", evaluator.BuiltinFunction{Func: std_all, Params: []string{"arr"}, OptStart: 1}},
	{"avg", evaluator.BuiltinFunction{Func: std_avg, Params: []string{"arr"}, OptStart: 1}},
	{"minArray", evaluator.BuiltinFunction{Func: std_minArray, Params: []string{"arr", "keyF", "onEmpty"}, OptStart: 1}},
	{"maxArray", evaluator.BuiltinFunction{Func: std_maxArray, Params: []string{"arr", "keyF", "onEmpty"}, OptStart: 1}},
	{"contains", evaluator.BuiltinFunction{Func: std_contains, Params: []string{"arr", "elem"}, OptStart: 2}},
	{"remove", evaluator.BuiltinFunction{Func: std_remove, Params: []string{"arr", "elem"}, OptStart: 2}},
	{"removeAt", evaluator.BuiltinFunction{Func: std_removeAt, Params: []string{"arr", "idx"}, OptStart: 2}},

	// -- Booleans ---
	{"xor", evaluator.BuiltinFunction{Func: std_xor, Params: []string{"x", "y"}, OptStart: 2}},
	{"xnor", evaluator.BuiltinFunction{Func: std_xnor, Params: []string{"x", "y"}, OptStart: 2}},

	// -- Sets ---
	{"set", evaluator.BuiltinFunction{Func: std_set, Params: []string{"arr", "keyF"}, OptStart: 1}},

	// --- Objects ---
	{"get", evaluator.BuiltinFunction{Func: std_get, Params: []string{"o", "f", "default", "inc_hidden"}, OptStart: 2}},
	{"objectFields", evaluator.BuiltinFunction{Func: std_objectFields, Params: []string{"o"}, OptStart: 1}},
	{"objectFieldsAll", evaluator.BuiltinFunction{Func: std_objectFieldsAll, Params: []string{"o"}, OptStart: 1}},
	{"objectFieldsEx", evaluator.BuiltinFunction{Func: std_objectFieldsEx, Params: []string{"obj", "hidden"}, OptStart: 2}},
	{"objectHas", evaluator.BuiltinFunction{Func: std_objectHas, Params: []string{"o", "f"}, OptStart: 2}},
	{"objectHasAll", evaluator.BuiltinFunction{Func: std_objectHasAll, Params: []string{"o", "f"}, OptStart: 2}},
	{"objectHasEx", evaluator.BuiltinFunction{Func: std_objectHasEx, Params: []string{"obj", "fname", "hidden"}, OptStart: 3}},
	{"objectValues", evaluator.BuiltinFunction{Func: std_objectValues, Params: []string{"o"}, OptStart: 1}},
	{"objectValuesAll", evaluator.BuiltinFunction{Func: std_objectValuesAll, Params: []string{"o"}, OptStart: 1}},
	{"objectKeysValues", evaluator.BuiltinFunction{Func: std_objectKeysValues, Params: []string{"o"}, OptStart: 1}},
	{"objectKeysValuesAll", evaluator.BuiltinFunction{Func: std_objectKeysValuesAll, Params: []string{"o"}, OptStart: 1}},
	{"mapWithKey", evaluator.BuiltinFunction{Func: std_mapWithkey, Params: []string{"func", "obj"}, OptStart: 2}},
	{"objectRemoveKey", evaluator.BuiltinFunction{Func: std_objectRemoveKey, Params: []string{"obj", "key"}, OptStart: 2}},
	{"mergePatch", evaluator.BuiltinFunction{Func: std_mergePatch, Params: []string{"target", "patch"}, OptStart: 2}},

	// --- Manifestation ---
	{"manifestYamlDoc", evaluator.BuiltinFunction{Func: std_manifestYamlDoc, Params: []string{"value", "indent_array_in_object", "quote_keys"}, OptStart: 1}},
	{"manifestYamlStream", evaluator.BuiltinFunction{Func: std_manifestYamlStream, Params: []string{"value", "indent_array_in_object", "c_document_end", "quote_keys"}, OptStart: 1}},
	{"manifestJson", evaluator.BuiltinFunction{Func: std_manifestJson, Params: []string{"value"}, OptStart: 1}},
	{"manifestJsonEx", evaluator.BuiltinFunction{Func: std_manifestJsonEx, Params: []string{"value", "indent", "newline", "key_val_sep"}, OptStart: 2}},
	{"manifestJsonMinified", evaluator.BuiltinFunction{Func: std_manifestJsonMinified, Params: []string{"value"}, OptStart: 1}},
	{"manifestIni", evaluator.BuiltinFunction{Func: std_manifestIni, Params: []string{"ini"}, OptStart: 1}},
	{"manifestPython", evaluator.BuiltinFunction{Func: std_manifestPython, Params: []string{"v"}, OptStart: 1}},
	{"manifestPythonVars", evaluator.BuiltinFunction{Func: std_manifestPythonVars, Params: []string{"conf"}, OptStart: 1}},
	{"manifestXmlJsonml", evaluator.BuiltinFunction{Func: std_manifestXmlJsonml, Params: []string{"value"}, OptStart: 1}},
	{"manifestTomlEx", evaluator.BuiltinFunction{Func: std_manifestTomlEx, Params: []string{"value", "indent"}, OptStart: 2}},
}

var (
	stdBuiltins = initBuiltins()
)

func initBuiltins() []evaluator.BuiltinFunction {
	res := make([]evaluator.BuiltinFunction, len(functions))
	for i := range evaluator.BuiltinFunctions {
		res[i] = functions[i].fn
	}
	return res
}

func InitStdLib(ctx evaluator.Context) (evaluator.Value, error) {

	ctx.State.Registry.BuiltinFunctions = stdBuiltins

	fieldCount := len(functions) + len(constants)

	layer := &evaluator.Layer{
		Keys:   make([]uint32, 0, fieldCount),
		Values: make([]evaluator.Value, 0, fieldCount),
		Meta:   make([]uint8, 0, fieldCount),

		// Index: make(map[uint32]int, fieldCount),
	}

	objId := evaluator.NewObject([]*evaluator.Layer{layer}, ctx)

	// index := 0
	for id, f := range functions {
		keyId := ctx.State.Interner.Intern(f.name)

		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, evaluator.MakeNativeFunction(uint32(id)))
		layer.Meta = append(layer.Meta, 0)
		// layer.Index[keyId] = index

		// index++
	}

	for _, c := range constants {
		keyId := ctx.State.Interner.Intern(c.name)

		layer.Keys = append(layer.Keys, keyId)
		layer.Values = append(layer.Values, c.val)
		layer.Meta = append(layer.Meta, 0)
		// layer.Index[keyId] = index

		// index++
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
