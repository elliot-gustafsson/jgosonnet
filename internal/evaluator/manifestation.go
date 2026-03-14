package evaluator

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type JsonManifestConfig struct {
	IndentStep string
	Newline    string
	KeyValSep  string
	SpaceComma bool
	hasNewline bool
}

// Pre-defined configurations matching Jsonnet's standard library.
var (
	JsonConfigPretty   = JsonManifestConfig{IndentStep: "    ", Newline: "\n", KeyValSep: ": ", SpaceComma: false}
	JsonConfigMinified = JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ":", SpaceComma: false}
	JsonConfigToString = JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ": ", SpaceComma: true}
)

type YamlManifestConfig struct {
	IndentArrayInObjects bool
	QuoteKeys            bool
	QuoteValues          bool
	SingleQuoteEscape    bool
}

func ManifestYaml(b *strings.Builder, value Value, ctx Context, config YamlManifestConfig) error {
	return manifestYaml(value, ctx, b, "", config)
}

// value Value, ctx Context, b *strings.Builder, cindent string, indent, newline, key_val_sep string
func ManifestJson(b *strings.Builder, value Value, ctx Context, config JsonManifestConfig) error {
	config.hasNewline = config.Newline != ""
	return manifestJson(value, ctx, b, "", config)
}

func manifestYaml(value Value, ctx Context, buf *strings.Builder, cindent string, config YamlManifestConfig) error {
	err := EvaluateValueStrict(&value, ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		// buf.WriteString(strconv.FormatFloat(data, 'f', -1, 64))
		buf.WriteString(unparseNumber(data))
		return nil
	case ValueTypeNull:
		buf.WriteString("null")
		return nil
	case ValueTypeBool:
		if value.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
		return nil
	case ValueTypeString:
		data := value.String(ctx)
		if data == "" {
			buf.WriteString(`""`)
			return nil
		}

		if strings.Contains(data, "\n") {

			buf.WriteByte('|')
			if !strings.HasSuffix(data, "\n") {
				buf.WriteByte('-')
			} else if strings.HasSuffix(data, "\n\n") || data == "\n" {
				buf.WriteByte('+')
			}

			for line := range strings.SplitSeq(strings.TrimSuffix(data, "\n"), "\n") {
				buf.WriteByte('\n')
				if line != "" {
					buf.WriteString(cindent)
					buf.WriteString(yamlIndent)
					buf.WriteString(line)
				}
			}
			return nil
		}

		if config.QuoteValues {
			writeYamlString(buf, data, true, false)
			return nil
		}

		writeYamlString(buf, data, false, true)
		return nil
	case ValueTypeArray:
		data := value.Array(ctx)
		if len(data) == 0 {
			buf.WriteString("[]")
			return nil
		}
		for i, v := range data {
			err := EvaluateValueStrict(&v, ctx)
			if err != nil {
				return err
			}

			if i != 0 {
				buf.WriteByte('\n')
				buf.WriteString(cindent)
			}
			buf.WriteByte('-')

			if v.IsArray() && len(v.Array(ctx)) > 0 {
				buf.WriteByte('\n')
				buf.WriteString(cindent)
				buf.WriteString(yamlIndent)
			} else {
				buf.WriteByte(' ')
			}

			prevIndent := cindent
			switch v.Type() {
			case ValueTypeArray, ValueTypeObject:
				cindent = cindent + yamlIndent
			}

			err = manifestYaml(v, ctx, buf, cindent, config)
			if err != nil {
				return err
			}
			cindent = prevIndent
		}
		return nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlan(obj, ctx)
		if len(plans) == 0 {
			buf.WriteString("{}")
			return nil
		}

		subCtx := ctx
		subCtx.Self = value

		hasWritten := false
		for _, p := range plans {
			if p.IsHidden() {
				continue
			}
			if hasWritten {
				buf.WriteByte('\n')
				buf.WriteString(cindent)
			}

			keyStr := ctx.Interner.Get(p.KeyId)
			if config.QuoteKeys || !yamlBareSafe(keyStr) {
				// buf.WriteByte('"')
				writeYamlString(buf, keyStr, true, false)
				// buf.WriteByte('"')
			} else {
				buf.WriteString(keyStr)
			}
			buf.WriteByte(':')
			prevIndent := cindent

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return err
			}

			if fieldValue.IsArray() && len(fieldValue.Array(subCtx)) > 0 {
				buf.WriteByte('\n')
				buf.WriteString(cindent)
				if config.IndentArrayInObjects {
					buf.WriteString(yamlIndent)
					cindent = cindent + yamlIndent
				}

			} else if fieldValue.IsObject() {
				// TODO: Write object isEmpty && isEmptyAll
				if len(GetObjectFields(fieldValue.Object(subCtx), subCtx, false)) > 0 {
					buf.WriteByte('\n')
					buf.WriteString(cindent)
					buf.WriteString(yamlIndent)
					cindent = cindent + yamlIndent
				} else {
					buf.WriteByte(' ')
				}
			} else {
				buf.WriteByte(' ')
			}

			err = manifestYaml(fieldValue, subCtx, buf, cindent, config)
			if err != nil {
				return err
			}
			hasWritten = true
			cindent = prevIndent
		}

		return nil
	}
}

func manifestJson(value Value, ctx Context, b *strings.Builder, cindent string, config JsonManifestConfig) error {

	err := EvaluateValueStrict(&value, ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		// buf.WriteString(strconv.FormatFloat(data, 'f', -1, 64))
		b.WriteString(unparseNumber(data))
		return nil
	case ValueTypeNull:
		b.WriteString("null")
		return nil
	case ValueTypeBool:
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
		// escaped := writeJsonString(data)
		// buf.WriteString(escaped)
		writeJsonString(b, data)
		return nil
	case ValueTypeArray:
		data := value.Array(ctx)
		if len(data) == 0 {
			if config.SpaceComma {
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
			err := EvaluateValueStrict(&v, ctx)
			if err != nil {
				return err
			}

			if i > 0 {
				b.WriteByte(',')
				b.WriteString(config.Newline)
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
			if config.SpaceComma {
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

const (
	yamlIndent = "  "
	hexChars   = "0123456789abcdef"
)

var yamlReserved = []string{
	// Boolean types taken from https://yaml.org/type/bool.html
	"y", "Y", "n", "N",
	"yes", "Yes", "YES", "no", "No", "NO",
	"true", "True", "TRUE", "false", "False", "FALSE",
	"on", "On", "ON", "off", "Off", "OFF",

	// Null types taken from https://yaml.org/type/null.html
	"null", "Null", "NULL", "~",

	// Numerical words taken from https://yaml.org/type/float.html
	".nan", ".NaN", ".NAN",
	".inf", ".Inf", ".INF",
	"+.inf", "+.Inf", "+.INF",
	"-.inf", "-.Inf", "-.INF",

	// Invalid keys that contain no invalid characters / Document markers
	"-", "---", "...", "''",
}

func yamlBareSafe(s string) bool {
	if len(s) == 0 {
		return false
	}

	if slices.Contains(yamlReserved, s) {
		return false
	}

	hasAlpha := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		isAlpha := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
		isDigit := c >= '0' && c <= '9'

		if isAlpha {
			hasAlpha = true
		}

		if !isAlpha && !isDigit && c != '_' && c != '-' && c != '/' && c != '.' && c != ':' {
			return false
		}
	}

	if hasAlpha {
		if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
			return false
		}
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			return false
		}
		return true
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return false
	}

	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			return false
		}
	}

	return true
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

func writeYamlString(b *strings.Builder, s string, forceQuotes, preferSingleQuotes bool) {
	needsQuotes := forceQuotes
	useSingle := preferSingleQuotes

	// --- PHASE 1: Determine if we can leave it bare ---
	if !needsQuotes {
		if len(s) == 0 {
			needsQuotes = true
		} else if slices.Contains(yamlReserved, s) {
			needsQuotes = true
			useSingle = false // Reserved words (true/null) need double quotes
		} else if _, err := strconv.ParseFloat(s, 64); err == nil {
			needsQuotes = true
			useSingle = false // Numbers as strings need double quotes
		} else if strings.TrimSpace(s) != s {
			needsQuotes = true
		} else {
			// Check for control characters
			for i := 0; i < len(s); i++ {
				if s[i] < 0x20 && s[i] != '\t' {
					needsQuotes = true
					useSingle = false // Control chars strictly require double quotes
					break
				}
			}
			// Check for structural indicators at the start
			if !needsQuotes {
				switch s[0] {
				case '[', ']', '{', '}', ',', '#', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
					needsQuotes = true
				case '-', '?', ':':
					if len(s) == 1 || s[1] == ' ' || s[1] == '\t' || s[1] == '\n' {
						needsQuotes = true
					}
				}
			}
			// Check for inline indicators and trailing colons (your fixes!)
			if !needsQuotes && (strings.Contains(s, ": ") || strings.Contains(s, ":\n") || strings.Contains(s, " #") || strings.HasSuffix(s, ":")) {
				needsQuotes = true
			}
		}
	}

	// If forced to quote, we STILL must check for control chars to override single quotes
	if needsQuotes && useSingle {
		for i := 0; i < len(s); i++ {
			if s[i] < 0x20 && s[i] != '\t' {
				useSingle = false
				break
			}
		}
	}

	// If it passed all checks, emit bare!
	if !needsQuotes {
		b.WriteString(s)
		return
	}

	// --- PHASE 2: Apply Quotes ---
	if useSingle {
		// Fast path for YAML single quotes
		b.WriteByte('\'')
		remaining := s
		for {
			idx := strings.IndexByte(remaining, '\'')
			if idx == -1 {
				// No more quotes found, write the rest of the string
				b.WriteString(remaining)
				break
			}
			b.WriteString(remaining[:idx])
			b.WriteString("''")
			remaining = remaining[idx+1:]
		}
		b.WriteByte('\'')

		return
	}

	// Fallback to JSON logic for double quotes!
	writeJsonString(b, s)
}
