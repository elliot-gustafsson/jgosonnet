package stdlib

import (
	"fmt"
	"math"
	"unicode/utf8"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
)

var constants = []struct {
	name string
	val  evaluator.Value
}{
	{"pi", evaluator.MakeNumber(math.Pi)},
}

var functions = []struct {
	name string
	fn   evaluator.NativeFunction
}{
	{"$flatMapArray", evaluator.NativeFunction{Func: builtin_flatMapArray, Params: []string{"func", "arr"}, OptStart: 2}},
	{"$objectFlatMerge", evaluator.NativeFunction{Func: builtin_objectFlatMerge, Params: []string{"arr"}, OptStart: 1}},
	{"extVar", evaluator.NativeFunction{Func: std_extVar, Params: []string{"x"}, OptStart: 1}},
	{"trace", evaluator.NativeFunction{Func: std_trace, Params: []string{"str", "rest"}, OptStart: 2}},
	{"assertEqual", evaluator.NativeFunction{Func: std_assertEqual, Params: []string{"a", "b"}, OptStart: 2}},
	{"toString", evaluator.NativeFunction{Func: std_toString, Params: []string{"a"}, OptStart: 1}},
	{"length", evaluator.NativeFunction{Func: std_length, Params: []string{"x"}, OptStart: 1}},
	{"mod", evaluator.NativeFunction{Func: std_mod, Params: []string{"a", "b"}, OptStart: 2}},

	// --- Types ---
	{"type", evaluator.NativeFunction{Func: std_type, Params: []string{"x"}, OptStart: 1}},
	{"isString", evaluator.NativeFunction{Func: std_isString, Params: []string{"v"}, OptStart: 1}},
	{"isNumber", evaluator.NativeFunction{Func: std_isNumber, Params: []string{"v"}, OptStart: 1}},
	{"isBoolean", evaluator.NativeFunction{Func: std_isBoolean, Params: []string{"v"}, OptStart: 1}},
	{"isObject", evaluator.NativeFunction{Func: std_isObject, Params: []string{"v"}, OptStart: 1}},
	{"isArray", evaluator.NativeFunction{Func: std_isArray, Params: []string{"v"}, OptStart: 1}},
	{"isFunction", evaluator.NativeFunction{Func: std_isFunction, Params: []string{"v"}, OptStart: 1}},
	{"isNull", evaluator.NativeFunction{Func: std_isNull, Params: []string{"x"}, OptStart: 1}},
	{"prune", evaluator.NativeFunction{Func: std_prune, Params: []string{"a"}, OptStart: 1}},

	// --- Parse ---
	{"parseInt", evaluator.NativeFunction{Func: std_parseInt, Params: []string{"str"}, OptStart: 1}},
	{"parseOctal", evaluator.NativeFunction{Func: std_parseOctal, Params: []string{"str"}, OptStart: 1}},
	{"parseHex", evaluator.NativeFunction{Func: std_parseHex, Params: []string{"str"}, OptStart: 1}},
	{"parseJson", evaluator.NativeFunction{Func: std_parseJson, Params: []string{"str"}, OptStart: 1}},
	{"parseYaml", evaluator.NativeFunction{Func: std_parseYaml, Params: []string{"str"}, OptStart: 1}},
	{"encodeUTF8", evaluator.NativeFunction{Func: std_encodeUTF8, Params: []string{"str"}, OptStart: 1}},
	{"decodeUTF8", evaluator.NativeFunction{Func: std_decodeUTF8, Params: []string{"arr"}, OptStart: 1}},

	// --- Math ---
	{"floor", evaluator.NativeFunction{Func: std_floor, Params: []string{"x"}, OptStart: 1}},
	{"ceil", evaluator.NativeFunction{Func: std_ceil, Params: []string{"x"}, OptStart: 1}},
	{"round", evaluator.NativeFunction{Func: std_round, Params: []string{"x"}, OptStart: 1}},
	{"pow", evaluator.NativeFunction{Func: std_pow, Params: []string{"x", "n"}, OptStart: 2}},
	{"sqrt", evaluator.NativeFunction{Func: std_sqrt, Params: []string{"x"}, OptStart: 1}},
	{"hypot", evaluator.NativeFunction{Func: std_hypot, Params: []string{"a", "b"}, OptStart: 2}},
	{"modulo", evaluator.NativeFunction{Func: std_modulo, Params: []string{"a", "b"}, OptStart: 2}},
	{"mantissa", evaluator.NativeFunction{Func: std_mantissa, Params: []string{"x"}, OptStart: 1}},
	{"exponent", evaluator.NativeFunction{Func: std_exponent, Params: []string{"x"}, OptStart: 1}},
	{"sin", evaluator.NativeFunction{Func: std_sin, Params: []string{"x"}, OptStart: 1}},
	{"cos", evaluator.NativeFunction{Func: std_cos, Params: []string{"x"}, OptStart: 1}},
	{"tan", evaluator.NativeFunction{Func: std_tan, Params: []string{"x"}, OptStart: 1}},
	{"asin", evaluator.NativeFunction{Func: std_asin, Params: []string{"x"}, OptStart: 1}},
	{"acos", evaluator.NativeFunction{Func: std_acos, Params: []string{"x"}, OptStart: 1}},
	{"atan", evaluator.NativeFunction{Func: std_atan, Params: []string{"x"}, OptStart: 1}},
	{"atan2", evaluator.NativeFunction{Func: std_atan2, Params: []string{"y", "x"}, OptStart: 2}},
	{"deg2rad", evaluator.NativeFunction{Func: std_deg2rad, Params: []string{"x"}, OptStart: 1}},
	{"rad2deg", evaluator.NativeFunction{Func: std_rad2deg, Params: []string{"x"}, OptStart: 1}},
	{"log", evaluator.NativeFunction{Func: std_log, Params: []string{"x"}, OptStart: 1}},
	{"log2", evaluator.NativeFunction{Func: std_log2, Params: []string{"x"}, OptStart: 1}},
	{"log10", evaluator.NativeFunction{Func: std_log10, Params: []string{"x"}, OptStart: 1}},
	{"exp", evaluator.NativeFunction{Func: std_exp, Params: []string{"x"}, OptStart: 1}},
	{"isEven", evaluator.NativeFunction{Func: std_isEven, Params: []string{"x"}, OptStart: 1}},
	{"isOdd", evaluator.NativeFunction{Func: std_isOdd, Params: []string{"x"}, OptStart: 1}},
	{"isInteger", evaluator.NativeFunction{Func: std_isInteger, Params: []string{"x"}, OptStart: 1}},
	{"isDecimal", evaluator.NativeFunction{Func: std_isDecimal, Params: []string{"x"}, OptStart: 1}},
	{"max", evaluator.NativeFunction{Func: std_max, Params: []string{"a", "b"}, OptStart: 2}},
	{"min", evaluator.NativeFunction{Func: std_min, Params: []string{"a", "b"}, OptStart: 2}},
	{"abs", evaluator.NativeFunction{Func: std_abs, Params: []string{"n"}, OptStart: 1}},
	{"sign", evaluator.NativeFunction{Func: std_sign, Params: []string{"n"}, OptStart: 1}},
	{"clamp", evaluator.NativeFunction{Func: std_clamp, Params: []string{"x", "minVal", "maxVal"}, OptStart: 3}},

	// --- Strings ---
	{"format", evaluator.NativeFunction{Func: std_format, Params: []string{"str", "vals"}, OptStart: 2}},
	{"stringChars", evaluator.NativeFunction{Func: std_stringChars, Params: []string{"str"}, OptStart: 1}},
	{"startsWith", evaluator.NativeFunction{Func: std_startsWith, Params: []string{"a", "b"}, OptStart: 2}},
	{"endsWith", evaluator.NativeFunction{Func: std_endsWith, Params: []string{"a", "b"}, OptStart: 2}},
	{"substr", evaluator.NativeFunction{Func: std_substr, Params: []string{"str", "from", "len"}, OptStart: 3}},
	{"findSubstr", evaluator.NativeFunction{Func: std_findSubstr, Params: []string{"pat", "str"}, OptStart: 2}},
	{"strReplace", evaluator.NativeFunction{Func: std_strReplace, Params: []string{"str", "from", "to"}, OptStart: 3}},
	{"split", evaluator.NativeFunction{Func: std_split, Params: []string{"str", "c"}, OptStart: 2}},
	{"splitLimit", evaluator.NativeFunction{Func: std_splitLimit, Params: []string{"str", "c", "maxsplits"}, OptStart: 3}},
	{"splitLimitR", evaluator.NativeFunction{Func: std_splitLimitR, Params: []string{"str", "c", "maxsplits"}, OptStart: 3}},
	{"stripChars", evaluator.NativeFunction{Func: std_stripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"rstripChars", evaluator.NativeFunction{Func: std_rstripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"lstripChars", evaluator.NativeFunction{Func: std_lstripChars, Params: []string{"str", "chars"}, OptStart: 2}},
	{"isEmpty", evaluator.NativeFunction{Func: std_isEmpty, Params: []string{"str"}, OptStart: 1}},
	{"trim", evaluator.NativeFunction{Func: std_trim, Params: []string{"str"}, OptStart: 1}},
	{"md5", evaluator.NativeFunction{Func: std_md5, Params: []string{"s"}, OptStart: 1}},
	{"sha1", evaluator.NativeFunction{Func: std_sha1, Params: []string{"s"}, OptStart: 1}},
	{"sha256", evaluator.NativeFunction{Func: std_sha256, Params: []string{"s"}, OptStart: 1}},
	{"sha512", evaluator.NativeFunction{Func: std_sha512, Params: []string{"s"}, OptStart: 1}},
	{"sha3", evaluator.NativeFunction{Func: std_sha3, Params: []string{"s"}, OptStart: 1}},
	{"char", evaluator.NativeFunction{Func: std_char, Params: []string{"n"}, OptStart: 1}},
	{"codepoint", evaluator.NativeFunction{Func: std_codepoint, Params: []string{"str"}, OptStart: 1}},
	{"base64", evaluator.NativeFunction{Func: std_base64, Params: []string{"input"}, OptStart: 1}},
	{"base64Decode", evaluator.NativeFunction{Func: std_base64Decode, Params: []string{"str"}, OptStart: 1}},
	{"base64DecodeBytes", evaluator.NativeFunction{Func: std_base64DecodeBytes, Params: []string{"str"}, OptStart: 1}},
	{"asciiLower", evaluator.NativeFunction{Func: std_asciiLower, Params: []string{"str"}, OptStart: 1}},
	{"asciiUpper", evaluator.NativeFunction{Func: std_asciiUpper, Params: []string{"str"}, OptStart: 1}},
	{"escapeStringBash", evaluator.NativeFunction{Func: std_escapeStringBash, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringDollars", evaluator.NativeFunction{Func: std_escapeStringDollars, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringJson", evaluator.NativeFunction{Func: std_escapeStringJson, Params: []string{"str_"}, OptStart: 1}},
	{"escapeStringPython", evaluator.NativeFunction{Func: std_escapeStringJson, Params: []string{"str"}, OptStart: 1}}, // Intentionally same as function as escapeStringJson
	{"escapeStringXML", evaluator.NativeFunction{Func: std_escapeStringXML, Params: []string{"str_"}, OptStart: 1}},
	{"equalsIgnoreCase", evaluator.NativeFunction{Func: std_equalsIgnoreCase, Params: []string{"str1", "str2"}, OptStart: 2}},

	// --- Arrays ---
	{"join", evaluator.NativeFunction{Func: std_join, Params: []string{"sep", "arr"}, OptStart: 2}},
	{"deepJoin", evaluator.NativeFunction{Func: std_deepJoin, Params: []string{"arr"}, OptStart: 1}},
	{"range", evaluator.NativeFunction{Func: std_range, Params: []string{"from", "to"}, OptStart: 2}},
	{"makeArray", evaluator.NativeFunction{Func: std_makeArray, Params: []string{"sz", "func"}, OptStart: 2}},
	{"filter", evaluator.NativeFunction{Func: std_filter, Params: []string{"func", "arr"}, OptStart: 2}},
	{"uniq", evaluator.NativeFunction{Func: std_uniq, Params: []string{"arr", "keyF"}, OptStart: 1}},
	{"sort", evaluator.NativeFunction{Func: std_sort, Params: []string{"arr", "keyF"}, OptStart: 1}},
	{"map", evaluator.NativeFunction{Func: std_map, Params: []string{"func", "arr"}, OptStart: 2}},
	{"mapWithIndex", evaluator.NativeFunction{Func: std_mapWithIndex, Params: []string{"func", "arr"}, OptStart: 2}},
	{"flatMap", evaluator.NativeFunction{Func: std_flatMap, Params: []string{"func", "arr"}, OptStart: 2}},
	{"filterMap", evaluator.NativeFunction{Func: std_filterMap, Params: []string{"filter_func", "map_func", "arr"}, OptStart: 3}},
	{"member", evaluator.NativeFunction{Func: std_member, Params: []string{"arr", "x"}, OptStart: 2}},
	{"setMember", evaluator.NativeFunction{Func: std_setMember, Params: []string{"x", "arr", "keyF"}, OptStart: 2}},
	{"slice", evaluator.NativeFunction{Func: std_slice, Params: []string{"indexable", "index", "end", "step"}, OptStart: 1}},
	{"count", evaluator.NativeFunction{Func: std_count, Params: []string{"arr", "x"}, OptStart: 2}},
	{"lines", evaluator.NativeFunction{Func: std_lines, Params: []string{"arr"}, OptStart: 1}},
	{"reverse", evaluator.NativeFunction{Func: std_reverse, Params: []string{"arrs"}, OptStart: 1}},
	{"foldl", evaluator.NativeFunction{Func: std_foldl, Params: []string{"func", "arr", "init"}, OptStart: 3}},
	{"foldr", evaluator.NativeFunction{Func: std_foldr, Params: []string{"func", "arr", "init"}, OptStart: 3}},
	{"sum", evaluator.NativeFunction{Func: std_sum, Params: []string{"arr"}, OptStart: 1}},
	{"flattenArrays", evaluator.NativeFunction{Func: std_flattenArrays, Params: []string{"arr"}, OptStart: 1}},
	{"flattenDeepArray", evaluator.NativeFunction{Func: std_flattenDeepArray, Params: []string{"arr"}, OptStart: 1}},
	{"repeat", evaluator.NativeFunction{Func: std_repeat, Params: []string{"what", "count"}, OptStart: 2}},
	{"setUnion", evaluator.NativeFunction{Func: std_setUnion, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"setInter", evaluator.NativeFunction{Func: std_setInter, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"setDiff", evaluator.NativeFunction{Func: std_setDiff, Params: []string{"a", "b", "keyF"}, OptStart: 2}},
	{"find", evaluator.NativeFunction{Func: std_find, Params: []string{"value", "arr"}, OptStart: 2}},
	{"any", evaluator.NativeFunction{Func: std_any, Params: []string{"arr"}, OptStart: 1}},
	{"all", evaluator.NativeFunction{Func: std_all, Params: []string{"arr"}, OptStart: 1}},
	{"avg", evaluator.NativeFunction{Func: std_avg, Params: []string{"arr"}, OptStart: 1}},
	{"minArray", evaluator.NativeFunction{Func: std_minArray, Params: []string{"arr", "keyF", "onEmpty"}, OptStart: 1}},
	{"maxArray", evaluator.NativeFunction{Func: std_maxArray, Params: []string{"arr", "keyF", "onEmpty"}, OptStart: 1}},
	{"contains", evaluator.NativeFunction{Func: std_contains, Params: []string{"arr", "elem"}, OptStart: 2}},
	{"remove", evaluator.NativeFunction{Func: std_remove, Params: []string{"arr", "elem"}, OptStart: 2}},
	{"removeAt", evaluator.NativeFunction{Func: std_removeAt, Params: []string{"arr", "idx"}, OptStart: 2}},

	// -- Booleans ---
	{"xor", evaluator.NativeFunction{Func: std_xor, Params: []string{"x", "y"}, OptStart: 2}},
	{"xnor", evaluator.NativeFunction{Func: std_xnor, Params: []string{"x", "y"}, OptStart: 2}},

	// -- Sets ---
	{"set", evaluator.NativeFunction{Func: std_set, Params: []string{"arr", "keyF"}, OptStart: 1}},

	// --- Objects ---
	{"get", evaluator.NativeFunction{Func: std_get, Params: []string{"o", "f", "default", "inc_hidden"}, OptStart: 2}},
	{"objectFields", evaluator.NativeFunction{Func: std_objectFields, Params: []string{"o"}, OptStart: 1}},
	{"objectFieldsAll", evaluator.NativeFunction{Func: std_objectFieldsAll, Params: []string{"o"}, OptStart: 1}},
	{"objectFieldsEx", evaluator.NativeFunction{Func: std_objectFieldsEx, Params: []string{"obj", "hidden"}, OptStart: 2}},
	{"objectHas", evaluator.NativeFunction{Func: std_objectHas, Params: []string{"o", "f"}, OptStart: 2}},
	{"objectHasAll", evaluator.NativeFunction{Func: std_objectHasAll, Params: []string{"o", "f"}, OptStart: 2}},
	{"objectHasEx", evaluator.NativeFunction{Func: std_objectHasEx, Params: []string{"obj", "fname", "hidden"}, OptStart: 3}},
	{"objectValues", evaluator.NativeFunction{Func: std_objectValues, Params: []string{"o"}, OptStart: 1}},
	{"objectValuesAll", evaluator.NativeFunction{Func: std_objectValuesAll, Params: []string{"o"}, OptStart: 1}},
	{"objectKeysValues", evaluator.NativeFunction{Func: std_objectKeysValues, Params: []string{"o"}, OptStart: 1}},
	{"objectKeysValuesAll", evaluator.NativeFunction{Func: std_objectKeysValuesAll, Params: []string{"o"}, OptStart: 1}},
	{"mapWithKey", evaluator.NativeFunction{Func: std_mapWithkey, Params: []string{"func", "obj"}, OptStart: 2}},
	{"objectRemoveKey", evaluator.NativeFunction{Func: std_objectRemoveKey, Params: []string{"obj", "key"}, OptStart: 2}},
	{"mergePatch", evaluator.NativeFunction{Func: std_mergePatch, Params: []string{"target", "patch"}, OptStart: 2}},

	// --- Manifestation ---
	{"manifestYamlDoc", evaluator.NativeFunction{Func: std_manifestYamlDoc, Params: []string{"value", "indent_array_in_object", "quote_keys"}, OptStart: 1}},
	{"manifestYamlStream", evaluator.NativeFunction{Func: std_manifestYamlStream, Params: []string{"value", "indent_array_in_object", "c_document_end", "quote_keys"}, OptStart: 1}},
	{"manifestJson", evaluator.NativeFunction{Func: std_manifestJson, Params: []string{"value"}, OptStart: 1}},
	{"manifestJsonEx", evaluator.NativeFunction{Func: std_manifestJsonEx, Params: []string{"value", "indent", "newline", "key_val_sep"}, OptStart: 2}},
	{"manifestJsonMinified", evaluator.NativeFunction{Func: std_manifestJsonMinified, Params: []string{"value"}, OptStart: 1}},
	{"manifestIni", evaluator.NativeFunction{Func: std_manifestIni, Params: []string{"ini"}, OptStart: 1}},
	{"manifestPython", evaluator.NativeFunction{Func: std_manifestPython, Params: []string{"v"}, OptStart: 1}},
	{"manifestPythonVars", evaluator.NativeFunction{Func: std_manifestPythonVars, Params: []string{"conf"}, OptStart: 1}},
	{"manifestXmlJsonml", evaluator.NativeFunction{Func: std_manifestXmlJsonml, Params: []string{"value"}, OptStart: 1}},
	{"manifestTomlEx", evaluator.NativeFunction{Func: std_manifestTomlEx, Params: []string{"value", "indent"}, OptStart: 2}},
}

func InitStdLib(ctx evaluator.Context) (evaluator.Value, error) {

	fieldCount := len(functions) + len(constants)
	allocator := ctx.State.Allocator

	layer := arena.Create[evaluator.Layer](allocator)
	arena.Memclr(layer)

	layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
	layer.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
	layer.Meta = arena.Alloc[uint8](allocator, fieldCount)
	layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)

	nativeFunctions := arena.Alloc[evaluator.NativeFunction](allocator, len(functions))
	arena.MemclrSlice(nativeFunctions)

	index := 0
	for i := range functions {
		nativeFunctions[i] = functions[i].fn

		keyId := ctx.State.Interner.Intern(functions[i].name)

		layer.Keys[index] = keyId
		layer.Values[index] = evaluator.MakeNativeFunctionValue(&nativeFunctions[i], ctx)

		layer.Index.Append(keyId)

		index++
	}

	for _, t := range constants {
		keyId := ctx.State.Interner.Intern(t.name)

		layer.Keys[index] = keyId
		layer.Values[index] = t.val

		layer.Index.Append(keyId)

		index++
	}

	obj := evaluator.NewSingleLayerObject(allocator, layer)
	val := evaluator.MakeObjectValue(obj)

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
		res = float64(len(arg.Array()))
	case evaluator.ValueTypeObject:
		res = float64(arg.Object().Length(ctx))
	case evaluator.ValueTypeFunction:
		res = float64(arg.Function().Length())
	case evaluator.ValueTypeNativeFunction:
		res = float64(arg.NativeFunction().Length())
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
