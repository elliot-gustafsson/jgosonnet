package stdlib

import (
	"fmt"

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

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: indent_array_in_object,
		QuoteKeys:            quote_keys,
		QuoteValues:          true,
	}

	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = evaluator.ManifestYaml(bp, v, ctx, c)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
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

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

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
		*bp = append(*bp, yamlSeparator...)
		*bp = append(*bp, '\n')

		err = evaluator.ManifestYaml(bp, v, ctx, c)
		if err != nil {
			return evaluator.ValueNone, err
		}
		*bp = append(*bp, '\n')
	}

	if c_document_end {
		*bp = append(*bp, "...\n"...)
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestJson(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = evaluator.ManifestJson(bp, a, ctx, evaluator.JsonConfigPretty)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestJsonMinified(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = evaluator.ManifestJson(bp, a, ctx, evaluator.JsonConfigMinified)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestJsonEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
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

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = evaluator.ManifestJson(bp, v, ctx, c)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestIni(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	iniObj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	mainKeyId := ctx.State.Interner.InternStatic("main")
	mainVal, visible, err := iniObj.GetFieldVis(mainKeyId, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if visible && !mainVal.IsNone() {
		mainVal, err := mainVal.EvalObject(ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		err = printIniSection(bp, mainVal, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
	}

	sectionsKeyId := ctx.State.Interner.InternStatic("sections")
	sectionsVal, err := iniObj.GetField(sectionsKeyId, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if sectionsVal.IsNone() {
		return evaluator.ValueNone, fmt.Errorf("expected field 'sections' does not exist on object passed to std.manifestIni")
	}

	if !sectionsVal.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, sectionsVal.Type())
	}

	sectionsObj := sectionsVal.Object()
	plans := evaluator.CompileObjectPlan(sectionsObj, ctx)

	for i := range plans {
		plan := &plans[i]
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

		*bp = append(*bp, '[')
		*bp = append(*bp, keyStr...)
		*bp = append(*bp, "]\n"...)
		err = printIniSection(bp, val.Object(), ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}
	}
	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func printIniSection(b *[]byte, obj *evaluator.Object, ctx evaluator.Context) error {

	plans := evaluator.CompileObjectPlan(obj, ctx)

	for i := range plans {
		plan := &plans[i]
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
			for _, v := range val.Array() {
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

func printIniVal(b *[]byte, k string, v evaluator.Value, ctx evaluator.Context) error {
	strVal, err := v.ToString(ctx)
	if err != nil {
		return err
	}
	*b = append(*b, k...)
	*b = append(*b, " = "...)
	*b = append(*b, strVal...)
	*b = append(*b, '\n')
	return nil
}

func std_manifestPython(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	v, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = evaluator.ManifestJson(bp, v, ctx, evaluator.JsonConfigPython)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestPythonVars(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	obj, err := args[0].EvalObject(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	plans := evaluator.CompileObjectPlan(obj, ctx)

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	for i := range plans {
		plan := &plans[i]
		if plan.IsHidden() {
			continue
		}

		name := ctx.State.Interner.Get(plan.KeyId)
		*bp = append(*bp, name...)
		*bp = append(*bp, " = "...)

		value, err := plan.GetValue(obj, ctx)
		if err != nil {
			return evaluator.ValueNone, err
		}

		err = evaluator.ManifestJson(bp, value, ctx, evaluator.JsonConfigPython)
		if err != nil {
			return evaluator.ValueNone, err
		}

		*bp = append(*bp, '\n')
	}
	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func std_manifestXmlJsonml(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	a, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !a.IsArray() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeArray, a.Type())
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	err = auxManifestXmlJsonml(a, ctx, bp)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}

func auxManifestXmlJsonml(v evaluator.Value, ctx evaluator.Context, b *[]byte) error {
	v, err := v.Eval(ctx)
	if err != nil {
		return err
	}

	if v.IsString() {
		*b = append(*b, v.String(ctx)...)
		return nil
	}

	if !v.IsArray() {
		return fmt.Errorf("unexpected type in array passed to std.manifestXmlJsonml: %s, expected string or array value", v.Type().String())
	}

	arr := v.Array()
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
			attrs = secondVal.Object()
			childrenStart = 2
		}
	}
	*b = append(*b, '<')
	*b = append(*b, tag...)
	if hasAttrs {
		fps := evaluator.CompileObjectPlan(attrs, ctx)
		for i := range fps {
			fp := &fps[i]
			if fp.IsHidden() {
				continue
			}
			keyStr := ctx.State.Interner.Get(fp.KeyId)
			fieldValue, err := fp.GetValue(attrs, ctx)
			if err != nil {
				return err
			}
			strVal, err := fieldValue.ToString(ctx)
			if err != nil {
				return err
			}
			*b = append(*b, ' ')
			*b = append(*b, keyStr...)
			*b = append(*b, `="`...)
			*b = append(*b, strVal...)
			*b = append(*b, '"')
		}
	}
	*b = append(*b, '>')
	for i := childrenStart; i < len(arr); i++ {
		err = auxManifestXmlJsonml(arr[i], ctx, b)
		if err != nil {
			return err
		}
	}
	*b = append(*b, "</"...)
	*b = append(*b, tag...)
	*b = append(*b, '>')
	return nil
}

func std_manifestTomlEx(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {
	val, err := args[0].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	if !val.IsObject() {
		return evaluator.ValueNone, evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, val.Type())
	}
	sindent, err := args[1].Value.EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	bp := evaluator.GetManifestBuf()
	defer evaluator.PutManifestBuf(bp)

	evalCtx := ctx
	evalCtx.Self = val

	err = evaluator.ManifestToml(bp, val, evalCtx, sindent)
	if err != nil {
		return evaluator.ValueNone, err
	}

	return evaluator.MakeStringFromBytes(*bp, ctx), nil
}
