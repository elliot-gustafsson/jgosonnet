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

	indent_array_in_object := false
	if !args[1].IsNone() {
		b, err := args[1].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		indent_array_in_object = b
	}

	quote_keys := true
	if !args[2].IsNone() {
		b, err := args[2].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		quote_keys = b
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
		return evaluator.ValueNone, err
	}
	err = evaluator.ManifestYaml(&b, v, ctx, c)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestYamlStream(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	inputArr, err := args[0].EvalArray(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	indent_array_in_object := false
	if !args[1].IsNone() {
		b, err := args[1].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		indent_array_in_object = b
	}

	c_document_end := true
	if !args[2].IsNone() {
		b, err := args[2].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		c_document_end = b
	}

	quote_keys := true
	if !args[3].IsNone() {
		b, err := args[3].EvalBool(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		quote_keys = b
	}

	var b strings.Builder
	b.Grow(1024)

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: indent_array_in_object,
		QuoteKeys:            quote_keys,
		QuoteValues:          true,
	}

	for _, v := range inputArr {
		v, err := v.Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		b.WriteString(yamlSeparator)
		b.WriteByte('\n')

		err = evaluator.ManifestYaml(&b, v, ctx, c)
		if err != nil {
			return evaluator.ValueNone, err
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

	var b strings.Builder
	b.Grow(1024)

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	err = evaluator.ManifestJson(&b, a, ctx, evaluator.JsonConfigPretty)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonMinified(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)

	var b strings.Builder
	b.Grow(1024)

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	err = evaluator.ManifestJson(&b, a, ctx, evaluator.JsonConfigMinified)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestJsonEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	// std.manifestJsonEx(value, indent, newline, key_val_sep)

	indent, err := args[1].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	newline := "\n"
	if !args[2].IsNone() {
		v, err := args[2].EvalString(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		newline = v
	}

	key_val_sep := ": "
	if !args[3].IsNone() {
		v, err := args[3].EvalString(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
		key_val_sep = v
	}

	var b strings.Builder
	b.Grow(1024)

	c := &evaluator.JsonManifestConfig{
		IndentStep:  indent,
		Newline:     newline,
		KeyValSep:   key_val_sep,
		StrictFloat: true,
	}
	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	err = evaluator.ManifestJson(&b, v, ctx, c)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestIni(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	iniObj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	var b strings.Builder
	b.Grow(512)

	mainKeyId := ctx.State.Interner.Intern("main")
	mainVal, visible, err := iniObj.GetField(mainKeyId, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	// std.objectHas matches even hidden fields, so we only check if it exists (!IsNone)
	if visible && !mainVal.IsNone() {
		mainVal, err := mainVal.EvalObject(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		err = printIniSection(&b, mainVal, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
	}

	sectionsKeyId := ctx.State.Interner.Intern("sections")
	sectionsVal, _, err := iniObj.GetField(sectionsKeyId, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if sectionsVal.IsNone() {
		return evaluator.ValueNone, fmt.Errorf("expected field 'sections' does not exist on object passed to std.manifestIni")
	}

	if !sectionsVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, sectionsVal.Type())
	}

	sectionsObj := sectionsVal.Object(ctx)
	plans := evaluator.CompileObjectPlan(sectionsObj, ctx)

	// TODO: look over error messages
	for _, plan := range plans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(sectionsObj, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		val, err = val.Eval(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		keyStr := ctx.State.Interner.Get(plan.KeyId)
		if !val.IsObject() {
			return evaluator.ValueNone, fmt.Errorf("expected object for ini section field '%s', got %s", keyStr, val.Type().String())
		}

		b.WriteByte('[')
		b.WriteString(keyStr)
		b.WriteString("]\n")
		err = printIniSection(&b, val.Object(ctx), ctx)
		if err != nil {
			return evaluator.ValueNone, err
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

		keyStr := ctx.State.Interner.Get(plan.KeyId)

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

	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	var b strings.Builder
	err = evaluator.ManifestJson(&b, v, ctx, evaluator.JsonConfigPython)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestPythonVars(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	plans := evaluator.CompileObjectPlan(obj, ctx)

	var b strings.Builder
	for _, plan := range plans {
		keyId := plan.KeyId

		if plan.IsHidden() {
			continue
		}

		name := ctx.State.Interner.Get(keyId)
		b.WriteString(name)
		b.WriteString(" = ")

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		err = evaluator.ManifestJson(&b, value, ctx, evaluator.JsonConfigPython)
		if err != nil {
			return evaluator.ValueNone, err
		}

		b.WriteByte('\n')
	}
	return evaluator.MakeString(b.String(), ctx), nil
}

func std_manifestXmlJsonml(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !a.IsArray() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, a.Type())
	}
	var b strings.Builder
	b.Grow(1024)
	err = auxManifestXmlJsonml(a, ctx, &b)
	if err != nil {
		return evaluator.ValueNone, err
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

	tag, err := arr[0].EvalString(ctx)
	if err != nil {
		return err
	}

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
			keyStr := ctx.State.Interner.Get(fp.KeyId)
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

	// Evaluate the value
	val, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	// TOML root must be an object
	if !val.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, val.Type())
	}
	// Evaluate the indent string
	sindent, err := args[1].Value.EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	var b strings.Builder

	evalCtx := ctx
	evalCtx.Self = val

	err = evaluator.ManifestToml(&b, val, evalCtx, sindent)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeString(b.String(), ctx), nil
}
