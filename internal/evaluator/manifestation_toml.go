package evaluator

import (
	"fmt"
)

func ManifestToml(b *[]byte, value Value, ctx Context, sindent string) error {

	if !value.IsObject() {
		return fmt.Errorf("root value must be object for toml manifestation, got: %s", value.Type().String())
	}

	obj := value.Object()

	var err error
	*b, err = renderTomlTable(*b, obj, ctx, sindent, []string{}, "", false)
	return err
}

func renderTomlTable(b []byte, obj *Object, ctx Context, sindent string, path []string, cindent string, initNewline bool) ([]byte, error) {
	if ctx.Depth >= ctx.State.MaxStack {
		return b, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
	}
	ctx.Depth++

	fieldPlans := CompileObjectPlan(obj, ctx)

	hasWritten := false

	complexValues := make([]NamedValue, 0, len(fieldPlans)/2)

	for i := range fieldPlans {
		plan := &fieldPlans[i]
		if plan.IsHidden() {
			continue
		}

		val, err := plan.GetValue(obj, ctx)
		if err != nil {
			return b, err
		}

		isSection, err := tomlIsSection(val, ctx)
		if err != nil {
			return b, err
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
		b = appendTomlKey(b, fieldName)
		b = append(b, " = "...)

		b, err = appendTomlValue(b, val, ctx, cindent, sindent, false)
		if err != nil {
			return b, err
		}

		hasWritten = true
	}

	if len(complexValues) > 0 {
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

		var err error
		switch val.Type() {
		case ValueTypeObject:
			b, err = appendTomlTable(b, val.Object(), ctx, sindent, childPath, cindent)
			if err != nil {
				return b, err
			}
		case ValueTypeArray:
			b, err = appendTomlTableArray(b, val.Array(), ctx, sindent, childPath, cindent)
			if err != nil {
				return b, err
			}
		default:
			return b, fmt.Errorf("invalid type for section: %s", val.Type().String())
		}
	}

	return b, nil
}

func appendTomlValue(b []byte, value Value, ctx Context, cindent, sindent string, inline bool) ([]byte, error) {
	if ctx.Depth >= ctx.State.MaxStack {
		return b, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
	}
	ctx.Depth++

	value, err := value.Eval(ctx)
	if err != nil {
		return b, err
	}

	switch value.Type() {
	default:
		return b, fmt.Errorf("unsupported value type %s for toml manifestation", value.Type().String())
	case ValueTypeBool:
		if value.Bool() {
			return append(b, "true"...), nil
		}
		return append(b, "false"...), nil
	case ValueTypeNumber:
		return unparseNumber(b, value.Number()), nil
	case ValueTypeString:
		return appendJsonString(b, value.String(ctx)), nil
	case ValueTypeArray:
		arr := value.Array()

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
				return b, err
			}

			if i > 0 {
				b = append(b, ',')
				b = append(b, separator...)
			}

			b = append(b, newIndent...)
			b, err = appendTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return b, err
			}
		}
		b = append(b, separator...)
		if !inline {
			b = append(b, cindent...)
		}
		b = append(b, ']')

		return b, nil
	case ValueTypeObject:
		obj := value.Object()
		evalCtx := ctx
		evalCtx.Self = value
		err := runAssertions(obj, evalCtx)
		if err != nil {
			return b, err
		}
		plans := CompileObjectPlan(obj, ctx)

		if len(plans) == 0 {
			return append(b, "{}"...), nil
		}

		b = append(b, "{ "...)

		hasWritten := false

		for i := range plans {
			plan := &plans[i]
			if plan.IsHidden() {
				continue
			}

			v, err := plan.GetValue(obj, ctx)
			if err != nil {
				return b, err
			}

			if hasWritten {
				b = append(b, ", "...)
			}

			fieldName := ctx.State.Interner.Get(plan.KeyId)
			b = appendTomlKey(b, fieldName)
			b = append(b, " = "...)

			b, err = appendTomlValue(b, v, ctx, "", sindent, true)
			if err != nil {
				return b, err
			}

			hasWritten = true
		}

		b = append(b, " }"...)
		return b, nil
	}
}

func appendTomlTable(b []byte, obj *Object, ctx Context, sindent string, path []string, cindent string) ([]byte, error) {
	b = append(b, cindent...)
	b = append(b, '[')

	for i, el := range path {
		if i > 0 {
			b = append(b, '.')
		}
		b = appendTomlKey(b, el)
	}

	b = append(b, ']')

	return renderTomlTable(b, obj, ctx, sindent, path, cindent+sindent, true)
}

func appendTomlTableArray(b []byte, arr []Value, ctx Context, sindent string, path []string, cindent string) ([]byte, error) {
	hasWritten := false
	for _, v := range arr {
		v, err := v.Eval(ctx)
		if err != nil {
			return b, err
		}

		if !v.IsObject() {
			return b, fmt.Errorf("invalid type for section: %s", v.Type().String())
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
			b = appendTomlKey(b, el)
		}
		b = append(b, "]]"...)

		b, err = renderTomlTable(b, v.Object(), ctx, sindent, path, cindent+sindent, true)
		if err != nil {
			return b, err
		}

		hasWritten = true
	}

	return b, nil
}

func appendTomlKey(b []byte, s string) []byte {
	bareAllowed := true

	if len(s) == 0 {
		return append(b, "''"...)
	}

	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			continue
		}

		bareAllowed = false
		break
	}

	if bareAllowed {
		return append(b, s...)
	}

	return appendJsonString(b, s)
}

func tomlIsSection(val Value, ctx Context) (bool, error) {
	if val.IsObject() {
		return true, nil
	} else if val.IsArray() {
		arr := val.Array()
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
