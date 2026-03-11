package stdlib

import (
	"fmt"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

const (
	yamlSeparator = "---"
)

func std_manifestYamlDoc(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestYamlDoc(value, indent_array_in_object=false, quote_keys=true)
	if len(args) < 1 || len(args) > 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlDoc: %d, expected 1-3", len(args))
	}

	indent_array_in_object := false
	if len(args) > 1 {
		if !args[1].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlDoc (arg 1): %s, expected boolean", args[1].Type().String())
		}
		indent_array_in_object = args[1].Bool()
	}

	quote_keys := true
	if len(args) > 2 {
		if !args[2].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlDoc (arg 2): %s, expected boolean", args[2].Type().String())
		}
		quote_keys = args[2].Bool()
	}

	var b strings.Builder
	b.Grow(1024)

	err := evaluator.ManifestYaml(&b, args[0], ctx, indent_array_in_object, quote_keys, true, false)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestYamlStream(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestYamlStream(value, indent_array_in_object=false, c_document_end=false, quote_keys=true)
	if len(args) < 1 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlStream: %d, expected 1-4", len(args))
	}

	inputArr := args[0]
	if !inputArr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 0): %s, expected array", args[0].Type().String())
	}

	indent_array_in_object := false
	if len(args) > 1 {
		if !args[1].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 1): %s, expected boolean", args[1].Type().String())
		}
		indent_array_in_object = args[1].Bool()
	}

	c_document_end := false
	if len(args) > 2 {
		if !args[2].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 2): %s, expected boolean", args[2].Type().String())
		}
		c_document_end = args[2].Bool()
	}

	quote_keys := true
	if len(args) > 3 {
		if !args[3].IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 3): %s, expected boolean", args[2].Type().String())
		}
		quote_keys = args[3].Bool()
	}

	var b strings.Builder
	b.Grow(1024)

	for _, v := range inputArr.Array(ctx) {
		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		b.WriteString(yamlSeparator)
		b.WriteByte('\n')

		err = evaluator.ManifestYaml(&b, v, ctx, indent_array_in_object, quote_keys, true, false)
		if err != nil {
			return evaluator.Value{}, err
		}
		b.WriteByte('\n')
	}

	if c_document_end {
		b.WriteString("...")
		b.WriteByte('\n')
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJson(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 1", len(args))
	}

	var b strings.Builder
	b.Grow(1024)

	err := evaluator.ManifestJson(&b, args[0], ctx, "    ", "\n", ": ")
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonMinified(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 1", len(args))
	}

	var b strings.Builder
	b.Grow(1024)

	err := evaluator.ManifestJson(&b, args[0], ctx, "", "", ":")
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonEx(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) < 2 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 2-4", len(args))
	}

	indent := args[1]
	if !indent.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 1): %s, expected string", indent.Type().String())
	}

	newline := "\n"
	if len(args) > 2 {
		if !args[2].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 2): %s, expected string", args[2].Type().String())
		}
		newline = args[2].String(ctx)
	}

	key_val_sep := ": "
	if len(args) > 3 {
		if !args[3].IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 3): %s, expected string", args[3].Type().String())
		}
		key_val_sep = args[3].String(ctx)
	}

	var b strings.Builder
	b.Grow(1024)

	err := evaluator.ManifestJson(&b, args[0], ctx, indent.String(ctx), newline, key_val_sep)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

// func std_manifestIni(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
// 	if len(args) != 1 {
// 		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestIni: %d, expected 1", len(args))
// 	}

// 	cfg := ini.Empty()

// }
