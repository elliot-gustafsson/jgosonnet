package evaluator

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

type YamlManifestConfig struct {
	IndentArrayInObjects bool
	QuoteKeys            bool
	QuoteValues          bool
	SingleQuoteEscape    bool
	NaturalSort          bool
	FormatIntegers       bool
}

func ManifestYaml(b *strings.Builder, value Value, ctx Context, config YamlManifestConfig) error {
	return manifestYaml(value, ctx, b, "", config)
}

func manifestYaml(value Value, ctx Context, buf *strings.Builder, cindent string, config YamlManifestConfig) error {
	value, err := value.Eval(ctx)
	if err != nil {
		return err
	}

	switch value.Type() {
	default:
		return fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		if config.FormatIntegers && data == math.Floor(data) {
			fmt.Fprintf(buf, "%.0f", data)
			return nil
		}
		buf.WriteString(strconv.FormatFloat(data, 'f', -1, 64))
		// buf.WriteString(unparseNumber(data))
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
			v, err := v.Eval(ctx)
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
		plans := CompileObjectPlanEx(obj, ctx, config.NaturalSort)
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

const (
	yamlIndent = "  "
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

		if !isAlpha && !isDigit && c != '_' && c != '-' && c != '/' && c != '.' /*&& c != ':'*/ {
			return false
		}
	}

	if hasAlpha {
		// if s[0] == '0' && len(s) > 1 && (s[1] == 'x' || s[1] == 'X') {
		// 	return false
		// }

		offset := 0
		if s[0] == '-' || s[0] == '+' {
			offset = 1
		}
		if len(s) > offset+1 && s[offset] == '0' && (s[offset+1] == 'x' || s[offset+1] == 'b') {
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

	// for i := 0; i < len(s); i++ {
	// 	if s[i] == '-' {
	// 		return false
	// 	}
	// }

	// Catch dates (e.g. 2001-12-14) but allow phone numbers (1-234-567-8901)
	hyphens := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			hyphens++
		}
	}
	if hyphens == 2 {
		return false
	}

	return true
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
		} else if n, err := strconv.ParseFloat(s, 64); err == nil {
			if !math.IsInf(n, 0) && !math.IsNaN(n) {
				needsQuotes = true
				useSingle = false // Numbers as strings need double quotes
			}
		} else if _, err := strconv.ParseInt(s, 0, 64); err == nil {
			// Handles hex and octal numbers
			needsQuotes = true
			useSingle = false
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
