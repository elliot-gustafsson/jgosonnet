package evaluator

import (
	"fmt"
	"math"
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

func ManifestJson(b *Builder, value Value, ctx Context, config *JsonManifestConfig) error {
	config.hasNewline = config.Newline != ""
	return manifestJson(value, ctx, b, 0, config)
}

func manifestJson(value Value, ctx Context, b *Builder, indentLevel int, config *JsonManifestConfig) error {

	value, err := value.Eval(ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		if config.StrictFloat {
			b.AppendFloat(data, 'f', -1, 64)
			return nil
		}
		unparseNumber(b, data)
		return nil
	case ValueTypeNull:
		if config.Python {
			b.WriteString("None")
			return nil
		}
		b.WriteString("null")
		return nil
	case ValueTypeBool:
		if config.Python {
			if value.Bool() {
				b.WriteString("True")
			} else {
				b.WriteString("False")
			}
			return nil
		}
		if value.Bool() {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		return nil
	case ValueTypeString:
		data := value.String(ctx)
		if data == "" {
			b.WriteString(`""`)
			return nil
		}
		writeJsonString(b, data)
		return nil
	case ValueTypeArray:
		data := value.Array(ctx)
		if len(data) == 0 {
			if config.SpaceComma && !config.Python {
				b.WriteString("[ ]")
				return nil
			}

			if config.hasNewline {
				b.WriteByte('[')
				b.WriteString(config.Newline)
				b.WriteString(config.Newline)
				writeIndent(b, indentLevel, config.IndentStep)
				b.WriteByte(']')
				return nil
			}

			b.WriteString("[]")
			return nil
		}

		b.WriteByte('[')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}

		b.WriteString(config.Newline)

		for i, v := range data {
			v, err := v.Eval(ctx)
			if err != nil {
				return err
			}

			if i > 0 {
				b.WriteByte(',')
				if config.hasNewline {
					b.WriteString(config.Newline)
				} else if config.SpaceComma {
					b.WriteByte(' ')
				}
			}

			if i != 0 || config.hasNewline {
				writeIndent(b, nextIndentLevel, config.IndentStep)
			}

			err = manifestJson(v, ctx, b, nextIndentLevel, config)
			if err != nil {
				return err
			}

		}

		b.WriteString(config.Newline)
		writeIndent(b, indentLevel, config.IndentStep)
		b.WriteByte(']')
		return nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlan(obj, ctx)
		if len(plans) == 0 {
			if config.SpaceComma && !config.Python {
				b.WriteString("{ }")
				return nil
			}

			if config.hasNewline {
				b.WriteByte('{')
				b.WriteString(config.Newline)
				b.WriteString(config.Newline)
				writeIndent(b, indentLevel, config.IndentStep)
				b.WriteByte('}')
				return nil
			}

			b.WriteString("{}")
			return nil
		}

		b.WriteByte('{')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}
		b.WriteString(config.Newline)

		subCtx := ctx
		subCtx.Self = value

		hasWritten := false

		for _, p := range plans {
			if p.IsHidden() {
				continue
			}

			if hasWritten {
				b.WriteByte(',')
				if config.hasNewline {
					b.WriteString(config.Newline)
				} else if config.SpaceComma {
					b.WriteByte(' ')
				}
			}

			writeIndent(b, nextIndentLevel, config.IndentStep)

			writeJsonString(b, subCtx.State.Interner.Get(p.KeyId))
			b.WriteString(config.KeyValSep)

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return err
			}

			err = manifestJson(fieldValue, subCtx, b, nextIndentLevel, config)
			if err != nil {
				return err
			}

			hasWritten = true
		}

		b.WriteString(config.Newline)
		writeIndent(b, indentLevel, config.IndentStep)
		b.WriteByte('}')

		return nil
	}

}

func writeIndent(b *Builder, indentLevel int, step string) {
	buf := b.Bytes()
	for range indentLevel {
		buf = append(buf, step...)
	}
	b.Set(buf)
}

// Borrowed from go-jsonnet, optimized to avoid fmt reflection
func unparseNumber(b *Builder, v float64) {
	if v == math.Floor(v) {
		// return fmt.Sprintf("%.0f", v)
		b.AppendFloat(v, 'f', 0, 64)
		return
	}

	// See "What Every Computer Scientist Should Know About Floating-Point Arithmetic"
	// Theorem 15
	// http://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
	// return fmt.Sprintf("%.17g", v)
	b.AppendFloat(v, 'g', 17, 64)
}

const (
	hexChars = "0123456789abcdef"
)

func writeJsonString(b *Builder, s string) {
	buf := b.Bytes()

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {

			// Fast path failed, proceed to slow path
			buf = append(buf, '"')
			buf = append(buf, s[:i]...)

			start := i
			for ; i < len(s); i++ {
				c := s[i]
				if c >= 0x20 && c != '"' && c != '\\' {
					continue
				}

				if start < i {
					buf = append(buf, s[start:i]...)
				}

				switch c {
				case '"', '\\':
					buf = append(buf, '\\')
					buf = append(buf, c)
				case '\n':
					buf = append(buf, `\n`...)
				case '\r':
					buf = append(buf, `\r`...)
				case '\t':
					buf = append(buf, `\t`...)
				case '\b':
					buf = append(buf, `\b`...)
				case '\f':
					buf = append(buf, `\f`...)
				default:
					buf = append(buf, `\u00`...)
					buf = append(buf, hexChars[c>>4])
					buf = append(buf, hexChars[c&0xF])
				}
				start = i + 1
			}

			if start < len(s) {
				buf = append(buf, s[start:]...)
			}
			buf = append(buf, '"')
			b.Set(buf)
			return
		}
	}

	// Fast path: No escapes needed. Just wrap in double quotes
	buf = append(buf, '"')
	buf = append(buf, s...)
	buf = append(buf, '"')

	b.Set(buf)
}
