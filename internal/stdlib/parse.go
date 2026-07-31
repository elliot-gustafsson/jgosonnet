package stdlib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"gopkg.in/yaml.v3"
)

func liftStringToValueErr(f func(string) (evaluator.Value, error)) evaluator.Func {
	return func(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

		a, err := args[0].EvalString(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		return f(a)
	}
}

var std_parseInt = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return evaluator.ValueNone, fmt.Errorf("error: %s is not a base 10 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
})

var std_parseOctal = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 8, 64)
	if err != nil {
		return evaluator.ValueNone, fmt.Errorf("error: %s is not a base 8 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
})

var std_parseHex = liftStringToValueErr(func(s string) (evaluator.Value, error) {
	num, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return evaluator.ValueNone, fmt.Errorf("error: %s is not a base 16 integer", s)
	}
	return evaluator.MakeNumber(float64(num)), nil
})

func std_parseJson(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	jsonString, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	jsonBytes := unsafe.Slice(unsafe.StringData(jsonString), len(jsonString))

	var data any
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return evaluator.ValueNone, fmt.Errorf("failed to parse json, err: %w", err)
	}
	return rawDataToValue(data, ctx)
}

func std_parseYaml(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	yamlString, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	yamlBytes := unsafe.Slice(unsafe.StringData(yamlString), len(yamlString))
	dec := yaml.NewDecoder(bytes.NewReader(yamlBytes))

	var documents []any
	for {

		var data any
		err := dec.Decode(&data)
		if err == io.EOF {
			break
		} else if err != nil {
			return evaluator.ValueNone, fmt.Errorf("failed to parse yaml, err: %w", err)
		}

		documents = append(documents, data)
	}
	if len(documents) == 0 {
		return evaluator.MakeNull(), nil
	}

	if len(documents) > 1 {
		return rawDataToValue(documents, ctx)
	}
	return rawDataToValue(documents[0], ctx)
}

func rawDataToValue(x any, ctx evaluator.Context) (evaluator.Value, error) {
	switch data := x.(type) {
	case nil:
		return evaluator.MakeNull(), nil
	case float64:
		return evaluator.MakeNumber(data), nil
	case int:
		return evaluator.MakeNumber(float64(data)), nil
	case string:
		return evaluator.MakeString(data, ctx), nil
	case bool:
		return evaluator.MakeBool(data), nil
	case []any:
		res, arrVal := evaluator.MakeArraySized(len(data), ctx)
		for i, rv := range data {
			v, err := rawDataToValue(rv, ctx)
			if err != nil {
				return evaluator.ValueNone, err
			}
			res[i] = v
		}
		return arrVal, nil
	case map[string]any:

		fieldCount := len(data)
		allocator := ctx.State.Registry.Allocator

		layer := arena.Create[evaluator.Layer](allocator)
		layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
		layer.Values = arena.Alloc[evaluator.Value](allocator, fieldCount)
		layer.Meta = arena.Alloc[uint8](allocator, fieldCount)

		useMap := fieldCount > evaluator.MaxLayerLinearKeys
		if useMap {
			layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
		}

		index := 0
		for keyName, value := range data {
			keyId := ctx.State.Interner.Intern(keyName)

			v, err := rawDataToValue(value, ctx)
			if err != nil {
				return evaluator.ValueNone, err
			}

			layer.Keys[index] = keyId
			layer.Values[index] = v
			layer.Meta[index] = evaluator.DefaultFieldMeta

			if useMap {
				layer.Index.Append(keyId)
			}

			index++
		}

		if index < fieldCount {
			layer.Keys = layer.Keys[:index]
			layer.Values = layer.Values[:index]
			layer.Meta = layer.Meta[:index]
		}

		layers := arena.Alloc[*evaluator.Layer](allocator, 1)
		layers[0] = layer
		objId := evaluator.NewObject(layers, ctx)

		return evaluator.MakeObjectValue(objId), nil
	default:
		return evaluator.ValueNone, fmt.Errorf("unahandled type %T when converting data to value", x)
	}

}

func std_encodeUTF8(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	str, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	strBytes := unsafe.Slice(unsafe.StringData(str), len(str))
	res, arrVal := evaluator.MakeArraySized(len(strBytes), ctx)
	for i, b := range strBytes {
		res[i] = evaluator.MakeNumber(float64(b))
	}

	return arrVal, nil
}

func std_decodeUTF8(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	arr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	var b strings.Builder
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		if !v.IsNumber() {
			return evaluator.ValueNone, fmt.Errorf("unexpected type in std.encodeUTF8 (arg 0) array: %s, expected number", v.Type().String())
		}
		num := v.Number()
		if num < 0 || num > 255 {
			return evaluator.ValueNone, fmt.Errorf("Bytes must be integers in range [0, 255], got %.0f", v.Number())
		}
		b.WriteByte(byte(num))
	}

	return evaluator.MakeString(b.String(), ctx), nil
}
