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
	return manifestJson(value, ctx, b, "", config)
}

func manifestJson(value Value, ctx Context, b *strings.Builder, cindent string, config *JsonManifestConfig) error {

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
			b.WriteString(strconv.FormatFloat(data, 'f', -1, 64))
			return nil
		}
		b.WriteString(unparseNumber(data))
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
				b.WriteString(cindent)
				b.WriteByte(']')
				return nil
			}

			b.WriteString("[]")
			return nil
		}

		b.WriteByte('[')
		nextIndent := cindent + config.IndentStep

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
				b.WriteString(nextIndent)
			}

			err = manifestJson(v, ctx, b, nextIndent, config)
			if err != nil {
				return err
			}

		}

		b.WriteString(config.Newline)
		b.WriteString(cindent)
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
				b.WriteString(cindent)
				b.WriteByte('}')
				return nil
			}

			b.WriteString("{}")
			return nil
		}

		b.WriteByte('{')
		nextIndent := cindent + config.IndentStep
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

			b.WriteString(nextIndent)

			writeJsonString(b, subCtx.Interner.Get(p.KeyId))
			b.WriteString(config.KeyValSep)

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return err
			}

			err = manifestJson(fieldValue, subCtx, b, nextIndent, config)
			if err != nil {
				return err
			}

			hasWritten = true
		}

		b.WriteString(config.Newline)
		b.WriteString(cindent)
		b.WriteByte('}')

		return nil
	}

}

// Borrowed from go-jsonnet
func unparseNumber(v float64) string {
	if v == math.Floor(v) {
		return fmt.Sprintf("%.0f", v)
	}

	// See "What Every Computer Scientist Should Know About Floating-Point Arithmetic"
	// Theorem 15
	// http://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
	return fmt.Sprintf("%.17g", v)
}

const (
	hexChars = "0123456789abcdef"
)

func writeJsonString(b *strings.Builder, s string) {
	needsEscape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {
			needsEscape = true
			break
		}
	}

	if !needsEscape {
		// Fast path: Just wrap in double quotes
		b.WriteByte('"')
		b.WriteString(s)
		b.WriteByte('"')
		return
	}

	// SLOW PATH: Full Builder Escaping
	// var b strings.Builder
	b.Grow(len(s) + 8)
	b.WriteByte('"')

	start := 0
	for i := 0; i < len(s); i++ {
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

}
