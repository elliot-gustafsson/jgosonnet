package stdlib

import (
	"fmt"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/google/go-jsonnet/ast"
)

var functions = map[string]evaluator.Func{
	// --- General ---
	"$flatMapArray":    builtin_flatMapArray,
	"$objectFlatMerge": builtin_objectFlatMerge,
	"toString":         std_toString,
	"length":           std_length,
	"mod":              std_mod,

	// --- Types ---
	"type":       std_type,
	"isString":   std_isString,
	"isNumber":   std_isNumber,
	"isBoolean":  std_isBoolean,
	"isObject":   std_isObject,
	"isArray":    std_isArray,
	"isFunction": std_isFunction,
	"prune":      std_prune,

	// --- Math ---
	"floor":     std_floor,
	"ceil":      std_ceil,
	"round":     std_round,
	"pow":       std_pow,
	"sqrt":      std_sqrt,
	"hypot":     std_hypot,
	"modulo":    std_modulo,
	"mantissa":  std_mantissa,
	"exponent":  std_exponent,
	"sin":       std_sin,
	"cos":       std_cos,
	"tan":       std_tan,
	"asin":      std_asin,
	"acos":      std_acos,
	"atan":      std_atan,
	"atan2":     std_atan2,
	"log":       std_log,
	"exp":       std_exp,
	"isEven":    std_isEven,
	"isOdd":     std_isOdd,
	"isInteger": std_isInteger,
	"isDecimal": std_isDecimal,

	// --- Strings ---
	"format":      std_format,
	"stringChars": std_stringChars,
	"startsWith":  std_startsWith,
	"endsWith":    std_endsWith,
	"substr":      std_substr,
	"findSubstr":  std_findSubstr,
	"split":       std_split,
	"splitLimit":  std_splitLimit,
	"rstripChars": std_rstripChars,
	"lstripChars": std_lstripChars,
	"isEmpty":     std_isEmpty,
	"trim":        std_trim,
	"md5":         std_md5,
	"sha1":        std_sha1,
	"sha256":      std_sha256,
	"sha512":      std_sha512,
	"sha3":        std_sha3,
	"char":        std_char,
	"codepoint":   std_codepoint,
	"parseInt":    std_parseInt,
	"base64":      std_base64,

	// --- Arrays ---
	"join":      std_join,
	"range":     std_range,
	"makeArray": std_makeArray,
	"filter":    std_filter,
	"uniq":      std_uniq,
	"sort":      std_sort,
	"map":       std_map,
	"filterMap": std_filterMap,
	"member":    std_member,
	"setMember": std_setMember,
	"slice":     std_slice,
	"count":     std_count,
	"lines":     std_lines,
	"reverse":   std_reverse,

	// -- Sets ---
	"set": std_set,

	// --- Objects ---
	"get":                 std_get,
	"objectFields":        std_objectFields,
	"objectFieldsAll":     std_objectFieldsAll,
	"objectHas":           std_objectHas,
	"objectHasAll":        std_objectHasAll,
	"objectValues":        std_objectValues,
	"objectValuesAll":     std_objectValuesAll,
	"objectKeysValues":    std_objectKeysValues,
	"objectKeysValuesAll": std_objectKeysValuesAll,

	// --- Manifestation ---
	"manifestYamlDoc":      std_manifestYamlDoc,
	"manifestYamlStream":   std_manifestYamlStream,
	"manifestJson":         std_manifestJson,
	"manifestJsonEx":       std_manifestJsonEx,
	"manifestJsonMinified": std_manifestJsonMinified,
}

func InitStdLib(ctx evaluator.Context) (evaluator.Value, error) {

	fieldCount := len(functions)

	layer := &evaluator.Layer{
		Keys:  make([]uint32, 0, fieldCount),
		Nodes: make(ast.Nodes, 0, fieldCount),
		Meta:  make([]uint8, 0, fieldCount),

		Index: make(map[uint32]int, fieldCount),
	}

	obj := evaluator.NewObject([]*evaluator.Layer{layer})

	obj.Values = make([]evaluator.Value, fieldCount)

	index := 0
	for name, f := range functions {
		keyId := ctx.Interner.Intern(name)

		v := evaluator.MakeFunction(f, ctx)
		layer.Keys = append(layer.Keys, keyId)
		layer.Meta = append(layer.Meta, 0)
		layer.Index[keyId] = index

		obj.Values[index] = v

		index++
	}

	val := evaluator.MakeObject(obj, ctx)

	return val, nil
}

func liftValueToBool(f func(evaluator.Value) bool, name string) evaluator.Func {
	return func(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
		if len(args) != 1 {
			return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to %s: %d, expected 1", name, len(args))
		}
		res := f(args[0])
		return evaluator.MakeBool(res), nil
	}
}

func std_type(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.type: %d, expected 1", len(args))
	}

	return evaluator.MakeString(args[0].Type().String(), ctx), nil
}

var std_isString = liftValueToBool(func(v evaluator.Value) bool { return v.IsString() }, "std.isString")
var std_isNumber = liftValueToBool(func(v evaluator.Value) bool { return v.IsNumber() }, "std.isNumber")
var std_isBoolean = liftValueToBool(func(v evaluator.Value) bool { return v.IsBool() }, "std.isBoolean")
var std_isObject = liftValueToBool(func(v evaluator.Value) bool { return v.IsObject() }, "std.isObject")
var std_isArray = liftValueToBool(func(v evaluator.Value) bool { return v.IsArray() }, "std.isArray")
var std_isFunction = liftValueToBool(func(v evaluator.Value) bool { return v.IsFunction() }, "std.isFunction")

func std_toString(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.toString: %d, expected 1", len(args))
	}

	s, err := args[0].ToString(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(s, ctx), nil
}

func std_length(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.length: %d, expected 1", len(args))
	}

	arg := args[0]

	var res float64
	switch arg.Type() {
	case evaluator.ValueTypeString:
		res = float64(len(arg.String(ctx)))
	case evaluator.ValueTypeArray:
		res = float64(len(arg.Array(ctx)))
	// case ValueTypeObject:
	// 	res = float64(arg.Object().GetLength())
	default:
		return evaluator.Value{}, fmt.Errorf("std.length: unexpected type %s", arg.Type().String())
	}

	return evaluator.MakeNumber(res), nil
}

func std_mod(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.mod: %d, expected 2", len(args))
	}
	if args[0].IsNumber() && args[1].IsNumber() {
		return std_modulo(args, ctx)
	}

	if args[0].IsString() {
		return std_format(args, ctx)
	}
	return evaluator.Value{}, fmt.Errorf("'Operator %% cannot be used on types %s and %s", args[0].Type().String(), args[1].Type().String())
}

func std_prune(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.prune: %d, expected 1", len(args))
	}
	arg := args[0]

	err := evaluator.EvaluateValueStrict(&arg, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	res, err := arg.Prune(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return res, nil
}
