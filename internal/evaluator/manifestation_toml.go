package evaluator

import (
	"fmt"
	"strings"
)

func ManifestToml(b *strings.Builder, value Value, ctx Context, sindent string) error {

	if !value.IsObject() {
		return fmt.Errorf("root value must be object for toml manifestation, got: %s", value.Type().String())
	}

	obj := value.Object(ctx)

	err := renderTomlTable(b, obj, ctx, sindent, []string{}, "", false)
	if err != nil {
		return err
	}
	return nil
}

func renderTomlTable(b *strings.Builder, obj *Object, ctx Context, sindent string, path []string, cindent string, initNewline bool) error {

	fieldPlans := CompileObjectPlan(obj, ctx)

	hasWritten := false

	complexValues := make([]NamedValue, 0, len(fieldPlans)/2)

	for _, plan := range fieldPlans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return err
		}

		isSection, err := tomlIsSection(val, ctx)
		if err != nil {
			return err
		}

		if isSection {
			complexValues = append(complexValues, NamedValue{plan.KeyId, val})
			continue
		}

		if initNewline {
			b.WriteByte('\n')
			initNewline = false
		}

		if hasWritten {
			b.WriteByte('\n')
		}

		fieldName := ctx.State.Interner.Get(plan.KeyId)

		b.WriteString(cindent)
		writeTomlKey(b, fieldName)
		b.WriteString(" = ")

		err = writeTomlValue(b, val, ctx, cindent, sindent, false)
		if err != nil {
			return err
		}

		hasWritten = true
	}

	if /* hasWritten && */ len(complexValues) > 0 {
		b.WriteString("\n\n")
	}

	for i, val := range complexValues {
		if i == 0 && !hasWritten && initNewline {
			b.WriteByte('\n')
		} else if i > 0 {
			b.WriteString("\n\n")
		}

		fieldName := ctx.State.Interner.Get(val.Key)

		childPath := tomlAddToPath(path, fieldName)

		switch val.Type() {
		case ValueTypeObject:
			err := writeTomlTable(b, val.Object(ctx), ctx, sindent, childPath, cindent)
			if err != nil {
				return err
			}
		case ValueTypeArray:
			err := writeTomlTableArray(b, val.Array(ctx), ctx, sindent, childPath, cindent)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid type for section: %s", val.Type().String())
		}
		continue
	}

	return nil
}

func writeTomlValue(b *strings.Builder, value Value, ctx Context, cindent, sindent string, inline bool) error {
	value, err := value.Eval(ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unsupported value type %s for toml manifestation", value.Type().String())
	case ValueTypeBool:
		if value.Bool() {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		return nil
	case ValueTypeNumber:
		var p [64]byte
		b.Write(unparseNumber(p[:0], value.Number()))
		return nil
	case ValueTypeString:
		writeJsonString(b, value.String(ctx))
		return nil
	case ValueTypeArray:
		arr := value.Array(ctx)

		if len(arr) == 0 {
			b.WriteString("[]")
			return nil
		}

		b.WriteByte('[')

		newIndent := cindent + sindent
		separator := "\n"
		if inline {
			newIndent = ""
			separator = " "
		}

		b.WriteString(separator)

		for i, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return err
			}

			if i > 0 {
				b.WriteByte(',')
				b.WriteString(separator)
			}

			b.WriteString(newIndent)
			err = writeTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return err
			}

		}
		b.WriteString(separator)
		if !inline {
			b.WriteString(cindent)
		}
		b.WriteByte(']')

		return nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlan(obj, ctx)

		if len(plans) == 0 {
			b.WriteString("{}")
			return nil
		}

		b.WriteString("{ ")

		hasWritten := false

		for _, plan := range plans {
			if plan.IsHidden() {
				continue
			}

			v, err := plan.GetValue(obj, ctx)
			if err != nil {
				return err
			}

			if hasWritten {
				b.WriteString(", ")
			}

			fieldName := ctx.State.Interner.Get(plan.KeyId)
			writeTomlKey(b, fieldName)
			b.WriteString(" = ")

			err = writeTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return err
			}

			hasWritten = true
		}

		b.WriteString(" }")
		return nil
	}

}

func writeTomlTable(b *strings.Builder, obj *Object, ctx Context, sindent string, path []string, cindent string) error {

	b.WriteString(cindent)
	b.WriteByte('[')

	for i, el := range path {
		if i > 0 {
			b.WriteByte('.')
		}
		writeTomlKey(b, el)

	}

	b.WriteByte(']')

	return renderTomlTable(b, obj, ctx, sindent, path, cindent+sindent, true)
}

func writeTomlTableArray(b *strings.Builder, arr []Value, ctx Context, sindent string, path []string, cindent string) error {

	hasWritten := false
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return err
		}

		if !v.IsObject() {
			return fmt.Errorf("invalid type for section: %s", v.Type().String())
		}
		if hasWritten {
			b.WriteString("\n\n")
		}

		b.WriteString(cindent)
		b.WriteString("[[")
		for i, el := range path {
			if i > 0 {
				b.WriteByte('.')
			}
			writeTomlKey(b, el)
		}
		b.WriteString("]]")

		err = renderTomlTable(b, v.Object(ctx), ctx, sindent, path, cindent+sindent, true)
		if err != nil {
			return err
		}

		hasWritten = true

	}

	return nil
}

func writeTomlKey(b *strings.Builder, s string) {
	bareAllowed := true

	// for empty string, return ''
	if len(s) == 0 {
		b.WriteString("''")
		return
	}

	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			continue
		}

		bareAllowed = false
		break
	}

	if bareAllowed {
		b.WriteString(s)
		return
	}

	writeJsonString(b, s)
}

func tomlIsSection(val Value, ctx Context) (bool, error) {
	if val.IsObject() {
		return true, nil
	} else if val.IsArray() {
		arr := val.Array(ctx)
		if len(arr) == 0 {
			return false, nil
		}
		for _, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return false, err
			}
			if !v.IsObject() {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}

func tomlAddToPath(path []string, tail string) []string {
	result := make([]string, len(path)+1)

	copy(result, path)

	result[len(path)] = tail

	return result
}
