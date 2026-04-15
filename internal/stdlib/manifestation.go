package stdlib

import (
	"fmt"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/google/go-jsonnet/ast"
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

	c_document_end := true
	if !args[2].IsNone() {
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
		v, err := v.Eval(ctx)
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

	mainKeyId := ctx.Interner.Intern("main")
	mainVal, visible, err := iniObj.GetField(mainKeyId, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	// std.objectHas matches even hidden fields, so we only check if it exists (!IsNone)
	if visible && !mainVal.IsNone() {
		mainVal, err := mainVal.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
		if !mainVal.IsObject() {
			return evaluator.Value{}, fmt.Errorf("expected object for ini main section, got %s", mainVal.Type().String())
		}

		err = printIniSection(&b, mainVal.Object(ctx), ctx)
		if err != nil {
			return evaluator.Value{}, err
		}
	}

	sectionsKeyId := ctx.Interner.Intern("sections")
	sectionsVal, _, err := iniObj.GetField(sectionsKeyId, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if sectionsVal.IsNone() {
		return evaluator.Value{}, fmt.Errorf("expected field 'sections' does not exist on object passed to std.manifestIni")
	}

	if !sectionsVal.IsObject() {
		return evaluator.Value{}, fmt.Errorf("expected object for ini 'sections' field, got %s", sectionsVal.Type().String())
	}

	sectionsObj := sectionsVal.Object(ctx)
	plans := evaluator.CompileObjectPlan(sectionsObj, ctx)

	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(sectionsObj, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		keyStr := ctx.Interner.Get(plan.KeyId)
		if !val.IsObject() {
			return evaluator.Value{}, fmt.Errorf("expected object for ini section field '%s', got %s", keyStr, val.Type().String())
		}

		b.WriteByte('[')
		b.WriteString(keyStr)
		b.WriteString("]\n")
		err = printIniSection(&b, val.Object(ctx), ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

	}
	return evaluator.MakeString(b.String(), ctx), nil
}

// printIniSection is our custom helper to loop through and print properties
func printIniSection(b *strings.Builder, obj *evaluator.Object, ctx evaluator.Context) error {

	plans := evaluator.CompileObjectPlan(obj, ctx)

	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return err
		}

		keyStr := ctx.Interner.Get(plan.KeyId)

		if val.IsArray() {
			for _, v := range val.Array(ctx) {
				err = printIniVal(b, keyStr, v, ctx)
				if err != nil {
					return err
				}
			}
			continue
		}

		err = printIniVal(b, keyStr, val, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func printIniVal(b *strings.Builder, k string, v evaluator.Value, ctx evaluator.Context) error {

	strVal, err := v.ToString(ctx)
	if err != nil {
		return err
	}
	b.WriteString(k)
	b.WriteString(" = ")
	b.WriteString(strVal)
	b.WriteByte('\n')

	return nil
}

func std_manifestPython(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestPython: %d, expected 1", len(args))
	}
	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	var b strings.Builder
	err = evaluator.ManifestJson(&b, v, ctx, evaluator.JsonConfigPython)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestPythonVars(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestPythonVars: %d, expected 1", len(args))
	}
	objVal, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}

	if !objVal.IsObject() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestPythonVars (arg 0): %s, expected object", objVal.Type().String())
	}
	obj := objVal.Object(ctx)

	plans := evaluator.CompileObjectPlan(obj, ctx)

	var b strings.Builder
	for _, plan := range plans {
		keyId := plan.KeyId

		if plan.Visibility == ast.ObjectFieldHidden {
			continue
		}

		name := ctx.Interner.Get(keyId)
		b.WriteString(name)
		b.WriteString(" = ")

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return evaluator.Value{}, err
		}

		err = evaluator.ManifestJson(&b, value, ctx, evaluator.JsonConfigPython)
		if err != nil {
			return evaluator.Value{}, err
		}

		b.WriteByte('\n')
	}
	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestXmlJsonml(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 1 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestXmlJsonml: %d, expected 1", len(args))
	}
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !a.IsArray() {
		return evaluator.Value{}, fmt.Errorf("unexpected type passed to std.manifestXmlJsonml (arg 0): %s, expected array", a.Type().String())
	}
	var b strings.Builder
	b.Grow(1024)
	err = auxManifestXmlJsonml(a, ctx, &b)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeString(b.String(), ctx), nil
}

func auxManifestXmlJsonml(v evaluator.Value, ctx evaluator.Context, b *strings.Builder) error {
	v, err := v.Eval(ctx)
	if err != nil {
		return err
	}

	if v.IsString() {
		b.WriteString(v.String(ctx))
		return nil
	}

	if !v.IsArray() {
		return fmt.Errorf("unexpected type in array passed to std.manifestXmlJsonml: %s, expected string or array value", v.Type().String())
	}

	arr := v.Array(ctx)
	if len(arr) == 0 {
		return fmt.Errorf("empty array passed to std.manifestXmlJsonml")
	}

	tagVal, err := arr[0].Eval(ctx)
	if err != nil {
		return err
	}

	if !tagVal.IsString() {
		return fmt.Errorf("unexpected type of tag value passed to std.manifestXmlJsonml: %s, expected string", tagVal.Type().String())
	}

	tag := tagVal.String(ctx)
	hasAttrs := false
	var attrs *evaluator.Object
	childrenStart := 1
	if len(arr) > 1 {
		secondVal, err := arr[1].Eval(ctx)
		if err != nil {
			return err
		}
		if secondVal.IsObject() {
			hasAttrs = true
			attrs = secondVal.Object(ctx)
			childrenStart = 2
		}
	}
	b.WriteByte('<')
	b.WriteString(tag)
	if hasAttrs {
		// CompileObjectPlan inherently sorts object keys alphabetically!
		fps := evaluator.CompileObjectPlan(attrs, ctx)
		for _, fp := range fps {
			if fp.IsHidden() {
				continue
			}
			keyStr := ctx.Interner.Get(fp.KeyId)
			fieldValue, err := fp.GetValue(attrs, ctx)
			if err != nil {
				return err
			}
			// ToString identically mimics jsonnet's `%s` stringification coercion
			strVal, err := fieldValue.ToString(ctx)
			if err != nil {
				return err
			}
			b.WriteByte(' ')
			b.WriteString(keyStr)
			b.WriteString(`="`)
			b.WriteString(strVal)
			b.WriteByte('"')
		}
	}
	b.WriteByte('>')
	for i := childrenStart; i < len(arr); i++ {
		err = auxManifestXmlJsonml(arr[i], ctx, b)
		if err != nil {
			return err
		}
	}
	b.WriteString("</")
	b.WriteString(tag)
	b.WriteByte('>')
	return nil
}

func std_manifestTomlEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.manifestTomlEx: %d, expected 2", len(args))
	}
	// Evaluate the value
	val, err := args[0].Value.Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	// TOML root must be an object
	if !val.IsObject() {
		return evaluator.Value{}, fmt.Errorf("TOML body must be an object. Got %s", val.Type().String())
	}
	// Evaluate the indent string
	vindent, err := args[1].Value.Eval(ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	if !vindent.IsString() {
		return evaluator.Value{}, fmt.Errorf("std.manifestTomlEx expects indent to be a string")
	}
	sindent := vindent.String(ctx)

	var b strings.Builder

	err = evaluator.ManifestToml(&b, val, ctx, sindent)
	if err != nil {
		return evaluator.Value{}, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}
