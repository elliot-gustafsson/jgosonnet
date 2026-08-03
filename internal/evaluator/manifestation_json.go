package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

func ManifestJson(b *strings.Builder, value Value, ctx Context, config *JsonManifestConfig) error {
	config.hasNewline = config.Newline != ""
	return manifestJson(value, ctx, b, 0, config)
}

func manifestJson(value Value, ctx Context, b *strings.Builder, indentLevel int, config *JsonManifestConfig) error {

	value, err := value.Eval(ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		var p [64]byte
		if config.StrictFloat {
			b.Write(strconv.AppendFloat(p[:0], data, 'f', -1, 64))
			return nil
		}
		b.Write(unparseNumber(p[:0], data))
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
		writeJsonString(b, data)
		return nil
	case ValueTypeArray:
		data := value.Array()
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
		obj := value.Object()
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

func writeIndent(b *strings.Builder, indentLevel int, step string) {
	if step == "" {
		return
	}

	// 64 spaces
	const maxIndentSpaces = "                                                                "

	if step == " " || step == "  " || step == "   " || step == "    " {
		totalSpaces := indentLevel * len(step)

		for totalSpaces > 0 {
			if totalSpaces <= len(maxIndentSpaces) {
				b.WriteString(maxIndentSpaces[:totalSpaces])
				break
			}

			b.WriteString(maxIndentSpaces)
			totalSpaces -= len(maxIndentSpaces)
		}
		return

	}

	for range indentLevel {
		b.WriteString(step)
	}
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

func writeJsonString(b *strings.Builder, s string) {
	if len(s) == 0 {
		b.WriteString(`""`)
		return
	}

	_ = s[len(s)-1]

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {

			// Fast path failed, proceed to slow path
			b.WriteByte('"')
			b.WriteString(s[:i])

			start := i
			for ; i < len(s); i++ {
				c := s[i]
				if c >= 0x20 && c != '"' && c != '\\' {
					continue
				}

				if start < i {
					b.WriteString(s[start:i])
				}

				switch c {
				case '"', '\\':
					b.WriteByte('\\')
					b.WriteByte(c)
				case '\n':
					b.WriteString(`\n`)
				case '\r':
					b.WriteString(`\r`)
				case '\t':
					b.WriteString(`\t`)
				case '\b':
					b.WriteString(`\b`)
				case '\f':
					b.WriteString(`\f`)
				default:
					b.WriteString(`\u00`)
					b.WriteByte(hexChars[c>>4])
					b.WriteByte(hexChars[c&0xF])
				}
				start = i + 1
			}

			if start < len(s) {
				b.WriteString(s[start:])
			}
			b.WriteByte('"')
			return
		}
	}

	// Fast path: No escapes needed. Just wrap in double quotes
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	b.WriteString(s)
	b.WriteByte('"')
}
