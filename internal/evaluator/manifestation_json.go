package evaluator

import (
	"fmt"
	"math"
	"strconv"
)

type JsonManifestConfig struct {
	IndentStep  string
	Newline     string
	KeyValSep   string
	SpaceComma  bool
	Python      bool
	StrictFloat bool

	hasNewline bool
}

// Pre-defined configurations matching Jsonnet's standard library.
var (
	JsonConfigPretty   = &JsonManifestConfig{IndentStep: "    ", Newline: "\n", KeyValSep: ": ", SpaceComma: false, StrictFloat: true}
	JsonConfigMinified = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ":", SpaceComma: false, StrictFloat: true}
	JsonConfigToString = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ": ", SpaceComma: true}
	JsonConfigPython   = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ": ", SpaceComma: true, Python: true}
)

func ManifestJson(b []byte, value Value, ctx Context, config *JsonManifestConfig) ([]byte, error) {
	config.hasNewline = config.Newline != ""
	return manifestJson(value, ctx, b, 0, config)
}

func manifestJson(value Value, ctx Context, b []byte, indentLevel int, config *JsonManifestConfig) ([]byte, error) {

	value, err := value.Eval(ctx)
	if err != nil {
		return nil, err
	}

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		if config.StrictFloat {
			return strconv.AppendFloat(b, data, 'f', -1, 64), nil
		}
		return unparseNumber(b, data), nil
	case ValueTypeNull:
		if config.Python {
			return append(b, "None"...), nil
		}
		return append(b, "null"...), nil
	case ValueTypeBool:
		if config.Python {
			if value.Bool() {
				return append(b, "True"...), nil
			}
			return append(b, "False"...), nil
		}
		if value.Bool() {
			return append(b, "true"...), nil
		}
		return append(b, "false"...), nil
	case ValueTypeString:
		data := value.String(ctx)
		if data == "" {
			return append(b, `""`...), nil
		}
		return writeJsonString(b, data), nil
	case ValueTypeArray:
		data := value.Array(ctx)
		if len(data) == 0 {
			if config.SpaceComma && !config.Python {
				return append(b, "[ ]"...), nil
			}

			if config.hasNewline {
				b = append(b, '[')
				b = append(b, config.Newline...)
				b = append(b, config.Newline...)
				b = writeJsonIndent(b, indentLevel, config.IndentStep)
				b = append(b, ']')
				return b, nil
			}

			return append(b, "[]"...), nil
		}

		b = append(b, '[')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}

		b = append(b, config.Newline...)

		for i, v := range data {
			v, err := v.Eval(ctx)
			if err != nil {
				return nil, err
			}

			if i > 0 {
				b = append(b, ',')
				if config.hasNewline {
					b = append(b, config.Newline...)
				} else if config.SpaceComma {
					b = append(b, ' ')
				}
			}

			if i != 0 || config.hasNewline {
				b = writeJsonIndent(b, nextIndentLevel, config.IndentStep)
			}

			b, err = manifestJson(v, ctx, b, nextIndentLevel, config)
			if err != nil {
				return nil, err
			}
		}

		b = append(b, config.Newline...)
		b = writeJsonIndent(b, indentLevel, config.IndentStep)
		b = append(b, ']')
		return b, nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlan(obj, ctx)
		if len(plans) == 0 {
			if config.SpaceComma && !config.Python {
				return append(b, "{ }"...), nil
			}

			if config.hasNewline {
				b = append(b, '{')
				b = append(b, config.Newline...)
				b = append(b, config.Newline...)
				b = writeJsonIndent(b, indentLevel, config.IndentStep)
				b = append(b, '}')
				return b, nil
			}

			return append(b, "{}"...), nil
		}

		b = append(b, '{')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}
		b = append(b, config.Newline...)

		subCtx := ctx
		subCtx.Self = value

		hasWritten := false

		for _, p := range plans {
			if p.IsHidden() {
				continue
			}

			if hasWritten {
				b = append(b, ',')
				if config.hasNewline {
					b = append(b, config.Newline...)
				} else if config.SpaceComma {
					b = append(b, ' ')
				}
			}

			b = writeJsonIndent(b, nextIndentLevel, config.IndentStep)

			b = writeJsonString(b, subCtx.State.Interner.Get(p.KeyId))
			b = append(b, config.KeyValSep...)

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return nil, err
			}

			b, err = manifestJson(fieldValue, subCtx, b, nextIndentLevel, config)
			if err != nil {
				return nil, err
			}

			hasWritten = true
		}

		b = append(b, config.Newline...)
		b = writeJsonIndent(b, indentLevel, config.IndentStep)
		b = append(b, '}')
		return b, nil
	}

}

func writeJsonIndent(b []byte, indentLevel int, step string) []byte {
	for range indentLevel {
		b = append(b, step...)
	}
	return b
}

// Borrowed from go-jsonnet, optimized to avoid fmt reflection
func unparseNumber(dst []byte, v float64) []byte {
	if v == math.Floor(v) {
		// return fmt.Sprintf("%.0f", v)
		return strconv.AppendFloat(dst, v, 'f', 0, 64)
	}

	// See "What Every Computer Scientist Should Know About Floating-Point Arithmetic"
	// Theorem 15
	// http://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
	// return fmt.Sprintf("%.17g", v)
	return strconv.AppendFloat(dst, v, 'g', 17, 64)
}

const (
	hexChars = "0123456789abcdef"
)

func writeJsonString(b []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {

			// Fast path failed, proceed to slow path
			b = append(b, '"')
			b = append(b, s[:i]...)

			start := i
			for ; i < len(s); i++ {
				c := s[i]
				if c >= 0x20 && c != '"' && c != '\\' {
					continue
				}

				if start < i {
					b = append(b, s[start:i]...)
				}

				switch c {
				case '"', '\\':
					b = append(b, '\\')
					b = append(b, c)
				case '\n':
					b = append(b, `\n`...)
				case '\r':
					b = append(b, `\r`...)
				case '\t':
					b = append(b, `\t`...)
				case '\b':
					b = append(b, `\b`...)
				case '\f':
					b = append(b, `\f`...)
				default:
					b = append(b, `\u00`...)
					b = append(b, hexChars[c>>4])
					b = append(b, hexChars[c&0xF])
				}
				start = i + 1
			}

			if start < len(s) {
				b = append(b, s[start:]...)
			}
			b = append(b, '"')
			return b
		}
	}

	// Fast path: No escapes needed. Just wrap in double quotes
	b = append(b, '"')
	b = append(b, s...)
	b = append(b, '"')
	return b
}
