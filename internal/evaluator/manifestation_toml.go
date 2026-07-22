package evaluator

import (
	"fmt"
)

func ManifestToml(b []byte, value Value, ctx Context, sindent string) ([]byte, error) {

	if !value.IsObject() {
		return nil, fmt.Errorf("root value must be object for toml manifestation, got: %s", value.Type().String())
	}

	obj := value.Object(ctx)

	return renderTomlTable(b, obj, ctx, sindent, []string{}, "", false)
}

func renderTomlTable(b []byte, obj *Object, ctx Context, sindent string, path []string, cindent string, initNewline bool) ([]byte, error) {

	fieldPlans := CompileObjectPlan(obj, ctx)

	hasWritten := false

	complexValues := make([]NamedValue, 0, len(fieldPlans)/2)

	for _, plan := range fieldPlans {
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return nil, err
		}

		isSection, err := tomlIsSection(val, ctx)
		if err != nil {
			return nil, err
		}

		if isSection {
			complexValues = append(complexValues, NamedValue{plan.KeyId, val})
			continue
		}

		if initNewline {
			b = append(b, '\n')
			initNewline = false
		}

		if hasWritten {
			b = append(b, '\n')
		}

		fieldName := ctx.State.Interner.Get(plan.KeyId)

		b = append(b, cindent...)
		b = writeTomlKey(b, fieldName)
		b = append(b, " = "...)

		b, err = writeTomlValue(b, val, ctx, cindent, sindent, false)
		if err != nil {
			return nil, err
		}

		hasWritten = true
	}

	if /* hasWritten && */ len(complexValues) > 0 {
		b = append(b, "\n\n"...)
	}

	for i, val := range complexValues {
		if i == 0 && !hasWritten && initNewline {
			b = append(b, '\n')
		} else if i > 0 {
			b = append(b, "\n\n"...)
		}

		fieldName := ctx.State.Interner.Get(val.Key)

		childPath := tomlAddToPath(path, fieldName)

		switch val.Type() {
		case ValueTypeObject:
			var err error
			b, err = writeTomlTable(b, val.Object(ctx), ctx, sindent, childPath, cindent)
			if err != nil {
				return nil, err
			}
		case ValueTypeArray:
			var err error
			b, err = writeTomlTableArray(b, val.Array(ctx), ctx, sindent, childPath, cindent)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("invalid type for section: %s", val.Type().String())
		}
		continue
	}

	return b, nil
}

func writeTomlValue(b []byte, value Value, ctx Context, cindent, sindent string, inline bool) ([]byte, error) {
	value, err := value.Eval(ctx)
	if err != nil {
		return nil, err
	}

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unsupported value type %s for toml manifestation", value.Type().String())
	case ValueTypeBool:
		if value.Bool() {
			return append(b, "true"...), nil
		}
		return append(b, "false"...), nil
	case ValueTypeNumber:
		return unparseNumber(b, value.Number()), nil
	case ValueTypeString:
		return writeJsonString(b, value.String(ctx)), nil
	case ValueTypeArray:
		arr := value.Array(ctx)

		if len(arr) == 0 {
			return append(b, "[]"...), nil
		}

		b = append(b, '[')

		newIndent := cindent + sindent
		separator := "\n"
		if inline {
			newIndent = ""
			separator = " "
		}

		b = append(b, separator...)

		for i, v := range arr {
			v, err := v.Eval(ctx)
			if err != nil {
				return nil, err
			}

			if i > 0 {
				b = append(b, ',')
				b = append(b, separator...)
			}

			b = append(b, newIndent...)
			b, err = writeTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return nil, err
			}

		}
		b = append(b, separator...)
		if !inline {
			b = append(b, cindent...)
		}
		b = append(b, ']')

		return b, nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlan(obj, ctx)

		if len(plans) == 0 {
			return append(b, "{}"...), nil
		}

		b = append(b, "{ "...)

		hasWritten := false

		for _, plan := range plans {
			if plan.IsHidden() {
				continue
			}

			v, err := plan.GetValue(obj, ctx)
			if err != nil {
				return nil, err
			}

			if hasWritten {
				b = append(b, ", "...)
			}

			fieldName := ctx.State.Interner.Get(plan.KeyId)
			b = writeTomlKey(b, fieldName)
			b = append(b, " = "...)

			b, err = writeTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return nil, err
			}

			hasWritten = true
		}

		b = append(b, " }"...)
		return b, nil
	}

}

func writeTomlTable(b []byte, obj *Object, ctx Context, sindent string, path []string, cindent string) ([]byte, error) {

	b = append(b, cindent...)
	b = append(b, '[')

	for i, el := range path {
		if i > 0 {
			b = append(b, '.')
		}
		b = writeTomlKey(b, el)
	}

	b = append(b, ']')

	return renderTomlTable(b, obj, ctx, sindent, path, cindent+sindent, true)
}

func writeTomlTableArray(b []byte, arr []Value, ctx Context, sindent string, path []string, cindent string) ([]byte, error) {

	hasWritten := false
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return nil, err
		}

		if !v.IsObject() {
			return nil, fmt.Errorf("invalid type for section: %s", v.Type().String())
		}
		if hasWritten {
			b = append(b, "\n\n"...)
		}

		b = append(b, cindent...)
		b = append(b, "[["...)
		for i, el := range path {
			if i > 0 {
				b = append(b, '.')
			}
			b = writeTomlKey(b, el)
		}
		b = append(b, "]]"...)

		b, err = renderTomlTable(b, v.Object(ctx), ctx, sindent, path, cindent+sindent, true)
		if err != nil {
			return nil, err
		}

		hasWritten = true
	}

	return b, nil
}

func writeTomlKey(b []byte, s string) []byte {
	bareAllowed := true

	// for empty string, return ''
	if len(s) == 0 {
		return append(b, "''"...)
	}

	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			continue
		}

		bareAllowed = false
		break
	}

	if bareAllowed {
		return append(b, s...)
	}

	return writeJsonString(b, s)
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
