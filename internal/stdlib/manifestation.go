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
	if len(args) != 3 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlDoc: %d, expected 3", len(args))
	}

	indent_array_in_object := false
	if !args[1].IsNone() {
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
	if !args[2].IsNone() {
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
	if len(args) != 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestYamlStream: %d, expected 4", len(args))
	}

	inputArr, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !inputArr.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestYamlStream (arg 0): %s, expected array", inputArr.Type().String())
	}

	indent_array_in_object := false
	if !args[1].IsNone() {
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
	if !args[1].IsNone() {
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
	if !args[3].IsNone() {
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
	if len(args) != 4 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestJsonEx: %d, expected 4", len(args))
	}

	indent, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !indent.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestJsonEx (arg 1): %s, expected string", indent.Type().String())
	}

	newline := "\n"
	if !args[2].IsNone() {
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
	if !args[3].IsNone() {
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

func std_manifestIni(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestIni: %d, expected 1", len(args))
	}
	iniObjVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !iniObjVal.IsObject() {
		return evaluator.Value{}, fmt.Errorf("expected object passed to std.manifestIni, got %s", iniObjVal.Type().String())
	}
	var b strings.Builder
	b.Grow(512)

	iniObj := iniObjVal.Object(ctx)
	// Handle 'main' section (top-level properties)
	mainKeyId := ctx.Interner.Intern("main")
	mainVal, _, err := iniObj.GetField(mainKeyId, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	// std.objectHas matches even hidden fields, so we only check if it exists (!IsNone)
	if !mainVal.IsNone() {
		err = evaluator.EvaluateValueStrict(&mainVal, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		err = printIniSection(&b, mainVal, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
	}
	// Handle 'sections' section (grouped properties)
	sectionsKeyId := ctx.Interner.Intern("sections")
	sectionsVal, _, err := iniObj.GetField(sectionsKeyId, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !sectionsVal.IsNone() {
		err = evaluator.EvaluateValueStrict(&sectionsVal, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !sectionsVal.IsObject() {
			return evaluator.Value{}, fmt.Errorf("expected object for 'sections' field")
		}
		sectionsObj := sectionsVal.Object(ctx)
		plans := evaluator.CompileObjectPlan(sectionsObj, ctx)

		for _, plan := range plans {
			// jsonnet's std.objectFields skips hidden fields, so we do too!
			if plan.IsHidden() {
				continue
			}
			secVal, err := plan.GetValue(sectionsObj, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			err = evaluator.EvaluateValueStrict(&secVal, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
			keyStr := ctx.Interner.Get(plan.KeyId)
			b.WriteString("[")
			b.WriteString(keyStr)
			b.WriteString("]\n")
			err = printIniSection(&b, secVal, ctx)
			if err != nil {
				return evaluator.Value{}, err
			}
		}
	}
	return evaluator.MakeString(b.String(), ctx), nil
}

// printIniSection is our custom helper to loop through and print properties
func printIniSection(b *strings.Builder, objVal evaluator.Value, ctx evaluator.Context) error {
	if !objVal.IsObject() {
		return fmt.Errorf("expected object for INI section")
	}

	obj := objVal.Object(ctx)
	plans := evaluator.CompileObjectPlan(obj, ctx)

	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return err
		}

		err = evaluator.EvaluateValueStrict(&val, ctx)
		if err != nil {
			return err
		}

		// val.ToString gives us the raw representation just like `%s` does in Jsonnet
		strVal, err := val.ToString(ctx)
		if err != nil {
			return err
		}

		keyStr := ctx.Interner.Get(plan.KeyId)
		b.WriteString(keyStr)
		b.WriteString(" = ")
		b.WriteString(strVal)
		b.WriteByte('\n')
	}
	return nil
}
