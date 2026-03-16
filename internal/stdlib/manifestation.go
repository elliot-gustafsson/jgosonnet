package stdlib

import (
	"fmt"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

const (
	yamlSeparator = "---"
)

func std_manifestYamlDoc(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestYamlDoc(value, indent_array_in_object=false, quote_keys=true)
	if len(args) < 1 || len(args) > 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlDoc: %d, expected 1-3", len(args))
	}

	indent_array_in_object := false
	if len(args) > 1 {
		v, err := args[1].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlDoc (arg 1): %s, expected boolean", v.Type().String())
		}
		indent_array_in_object = v.Bool()
	}

	quote_keys := true
	if len(args) > 2 {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlDoc (arg 2): %s, expected boolean", v.Type().String())
		}
		quote_keys = v.Bool()
	}

	var b strings.Builder
	b.Grow(1024)

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: indent_array_in_object,
		QuoteKeys:            quote_keys,
		QuoteValues:          true,
	}

	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	err = evaluator.ManifestYaml(&b, v, ctx, c)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestYamlStream(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestYamlStream(value, indent_array_in_object=false, c_document_end=false, quote_keys=true)
	if len(args) < 1 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlStream: %d, expected 1-4", len(args))
	}

	inputArr, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !inputArr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 0): %s, expected array", inputArr.Type().String())
	}

	indent_array_in_object := false
	if len(args) > 1 {
		v, err := args[1].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 1): %s, expected boolean", v.Type().String())
		}
		indent_array_in_object = v.Bool()
	}

	c_document_end := false
	if len(args) > 2 {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 2): %s, expected boolean", v.Type().String())
		}
		c_document_end = v.Bool()
	}

	quote_keys := true
	if len(args) > 3 {
		v, err := args[3].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsBool() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 3): %s, expected boolean", v.Type().String())
		}
		quote_keys = v.Bool()
	}

	var b strings.Builder
	b.Grow(1024)

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: indent_array_in_object,
		QuoteKeys:            quote_keys,
		QuoteValues:          true,
	}

	for _, v := range inputArr.Array(ctx) {
		err := evaluator.EvaluateValueStrict(&v, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		b.WriteString(yamlSeparator)
		b.WriteByte('\n')

		err = evaluator.ManifestYaml(&b, v, ctx, c)
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

func std_manifestJson(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 1", len(args))
	}

	var b strings.Builder
	b.Grow(1024)

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	err = evaluator.ManifestJson(&b, a, ctx, evaluator.JsonConfigPretty)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonMinified(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 1", len(args))
	}

	var b strings.Builder
	b.Grow(1024)

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	err = evaluator.ManifestJson(&b, a, ctx, evaluator.JsonConfigMinified)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)
	if len(args) < 2 || len(args) > 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 2-4", len(args))
	}

	indent, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !indent.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 1): %s, expected string", indent.Type().String())
	}

	newline := "\n"
	if len(args) > 2 {
		v, err := args[2].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 2): %s, expected string", v.Type().String())
		}
		newline = v.String(ctx)
	}

	key_val_sep := ": "
	if len(args) > 3 {
		v, err := args[3].Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !v.IsString() {
			return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 3): %s, expected string", v.Type().String())
		}
		key_val_sep = v.String(ctx)
	}

	var b strings.Builder
	b.Grow(1024)

	c := evaluator.JsonManifestConfig{
		IndentStep: indent.String(ctx),
		Newline:    newline,
		KeyValSep:  key_val_sep,
	}
	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	err = evaluator.ManifestJson(&b, v, ctx, c)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

// func std_manifestIni(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
// 	if len(args) != 1 {
// 		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestIni: %d, expected 1", len(args))
// 	}

// 	cfg := ini.Empty()

// }
