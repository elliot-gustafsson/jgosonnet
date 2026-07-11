package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type YamlManifestConfig struct {
	IndentArrayInObjects bool
	QuoteKeys            bool
	QuoteValues          bool
	NaturalSort          bool
	FormatIntegers       bool
	UseBlockScalars      bool
	Modern               bool
}

func ManifestYaml(b *strings.Builder, value Value, ctx Context, config YamlManifestConfig) error {
	return manifestYaml(value, ctx, b, 0, config)
}

func manifestYaml(value Value, ctx Context, buf *strings.Builder, indentLevel int, config YamlManifestConfig) error {
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
			buf.WriteString(strconv.FormatFloat(data, 'f', 0, 64))
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

		n := len(data)
		if n == 0 {
			buf.WriteString(`""`)
			return nil
		}

		var multiline bool
		lastByte := data[n-1]

		if config.UseBlockScalars {
			// Check if the string contains newlines to determine if it can be a block scalar.
			// Since YAML block scalars do not safely preserve trailing whitespace on lines,
			// we must safely fall back to a quoted single-line string if any line ends with a space or tab.
			for i := 0; i < n; i++ {
				if data[i] != '\n' {
					continue
				}

				multiline = true
				if i > 0 && (data[i-1] == ' ' || data[i-1] == '\t') {
					multiline = false
					break
				}
			}
			if multiline && (lastByte == ' ' || lastByte == '\t') {
				multiline = false
			}

		} else {
			multiline = lastByte == '\n'
		}

		if !multiline {
			if config.QuoteValues {
				writeYamlString(buf, data, true, false, config.Modern)
				return nil
			}

			writeYamlString(buf, data, false, true, config.Modern)
			return nil
		}

		buf.WriteByte('|')
		if config.UseBlockScalars {
			firstByte := data[0]
			if firstByte == ' ' || /* data[0] == '\t' || */ firstByte == '\n' {
				buf.WriteString(yamlIndentNumber)
			}

			if lastByte != '\n' {
				buf.WriteByte('-')
			} else if n >= 2 && data[n-2] == '\n' || n == 1 {
				buf.WriteByte('+')
			}

			if firstByte == '\n' {
				data = data[1:] // prefix trim
				n--             // update length
				if n == 0 {
					return nil
				}
			}
		}

		start := 0
		end := n
		if end > 0 && data[end-1] == '\n' {
			end--
		}

		// read data up to newlines, print newline then print line if not empty
		for start <= end {
			idx := strings.IndexByte(data[start:end], '\n')

			var line string
			if idx == -1 {
				line = data[start:end]
				start = end + 1 // Break next loop
			} else {
				line = data[start : start+idx]
				start += idx + 1
			}

			buf.WriteByte('\n')
			if line != "" || !config.UseBlockScalars {
				writeYamlIndent(buf, indentLevel+1)
				buf.WriteString(line)
			}
		}

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
				writeYamlIndent(buf, indentLevel)
			}
			buf.WriteByte('-')

			if v.IsArray() && len(v.Array(ctx)) > 0 {
				buf.WriteByte('\n')
				writeYamlIndent(buf, indentLevel+1)
			} else {
				buf.WriteByte(' ')
			}

			nextIndentLevel := indentLevel
			switch v.Type() {
			case ValueTypeArray, ValueTypeObject:
				nextIndentLevel++
			}

			err = manifestYaml(v, ctx, buf, nextIndentLevel, config)
			if err != nil {
				return err
			}
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
				writeYamlIndent(buf, indentLevel)
			}

			keyStr := ctx.State.Interner.Get(p.KeyId)
			writeYamlString(buf, keyStr, config.QuoteKeys, config.Modern, config.Modern)

			buf.WriteByte(':')

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return err
			}

			nextIndentLevel := indentLevel

			if fieldValue.IsArray() && len(fieldValue.Array(subCtx)) > 0 {
				buf.WriteByte('\n')
				writeYamlIndent(buf, indentLevel)
				if config.IndentArrayInObjects {
					writeYamlIndent(buf, 1)
					nextIndentLevel++
				}

			} else if fieldValue.IsObject() {
				hasFields := false

				plans := compileObjectPlan(fieldValue.Object(subCtx), subCtx)
				for i := range plans {
					if !plans[i].IsHidden() {
						hasFields = true
						break
					}
				}

				if hasFields {
					buf.WriteByte('\n')
					writeYamlIndent(buf, indentLevel+1)
					nextIndentLevel++
				} else {
					buf.WriteByte(' ')
				}
			} else {
				buf.WriteByte(' ')
			}

			err = manifestYaml(fieldValue, subCtx, buf, nextIndentLevel, config)
			if err != nil {
				return err
			}
			hasWritten = true
		}

		return nil
	}
}

const (
	yamlIndentSpaces = 2
)

var (
	yamlIndentNumber = strconv.Itoa(yamlIndentSpaces)
)

func writeYamlIndent(b *strings.Builder, indentLevel int) {
	// 64 spaces
	const maxIndentString = "                                                                "

	totalSpaces := indentLevel * yamlIndentSpaces

	for totalSpaces > 0 {
		// If the remaining spaces fit in our pre-allocated string,
		// slice it, write it, and we are done!
		if totalSpaces <= len(maxIndentString) {
			b.WriteString(maxIndentString[:totalSpaces])
			break
		}

		// Otherwise, write the max chunk of 64 spaces and subtract it
		b.WriteString(maxIndentString)
		totalSpaces -= len(maxIndentString)
	}
}

func yamlReserved(s string) bool {
	switch s {
	case
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
		"-", "---", "...", "''":
		return true
	}
	return false
}

func writeYamlString(b *strings.Builder, s string, forceQuotes, preferSingleQuotes, modern bool) {
	if len(s) == 0 {
		if preferSingleQuotes {
			b.WriteString("''")
		} else {
			b.WriteString(`""`)
		}
		return
	}

	needsQuotes := forceQuotes
	useSingle := preferSingleQuotes

	// prefix and suffix checks
	if !needsQuotes {

		if unicode.IsSpace(rune(s[0])) || unicode.IsSpace(rune(s[len(s)-1])) {
			needsQuotes = true
		} else {
			// structural indicators at the start
			switch s[0] {
			case '[', ']', '{', '}', ',', '#', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
				needsQuotes = true
			case '-', '?', ':':
				if len(s) == 1 || s[1] == ' ' || s[1] == '\t' || s[1] == '\n' {
					needsQuotes = true
				}
			}
		}

		if !needsQuotes {

			if yamlReserved(s) || isYamlNumber(s) || isYamlTimestamp(s) {
				b.WriteByte('"')
				b.WriteString(s)
				b.WriteByte('"')
				return
			}
		}
	}

	hasSingleQuote := false
	for i := 0; i < len(s); i++ {
		if needsQuotes && !useSingle {
			break
		}

		c := s[i]

		// force quotes on control chars
		if c < 0x20 || c == 0x7F {
			useSingle = false
			needsQuotes = true
			break
		}

		if c == '\'' {
			hasSingleQuote = true
		}

		// look for internal structural markers (": ", ":\n", trailing ":", " #")
		if !needsQuotes {

			if modern {
				if c == ':' {
					if i == len(s)-1 || s[i+1] == ' ' || s[i+1] == '\n' {
						needsQuotes = true
					}
				} else if c == ' ' && i+1 < len(s) && s[i+1] == '#' {
					needsQuotes = true
				}
			} else {

				isAlpha := (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
				isDigit := c >= '0' && c <= '9'
				if !isAlpha && !isDigit && c != '_' && c != '-' && c != '/' && c != '.' {
					needsQuotes = true
					useSingle = false
					break
				}
			}

		}
	}

	if !needsQuotes {
		b.WriteString(s)
		return
	}

	if !useSingle {
		writeJsonString(b, s)
		return
	}

	if !hasSingleQuote {
		b.WriteByte('\'')
		b.WriteString(s)
		b.WriteByte('\'')
		return
	}

	// escape single quotes
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
}

func isYamlNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	first := s[0]
	if first != '+' && first != '-' && first != '.' && (first < '0' || first > '9') {
		return false
	}

	i := 0
	if s[i] == '+' || s[i] == '-' {
		i++
	}

	// binary, hex, octal checks
	if i+2 < len(s) && s[i] == '0' {
		base := s[i+1]
		if base == 'b' || base == 'x' || base == 'o' {
			return checkBaseFormat(s, i+2, base)
		}
	}

	hasDot, hasE, hasColon, hasDigit := false, false, false, false
	digitsSinceColon, colonValue := 0, 0

	for j := i; j < len(s); j++ {
		c := s[j]
		if c >= '0' && c <= '9' {
			hasDigit = true

			// validate base-60 number
			if hasColon && !hasDot {
				digitsSinceColon++
				if digitsSinceColon > 2 {
					return false
				}
				if digitsSinceColon == 1 {
					colonValue = int(c - '0')
				} else {
					colonValue = colonValue*10 + int(c-'0')
				}
				if colonValue >= 60 {
					return false
				}
			}
		} else if c == '_' {
			// underscores are treated as spacers
			continue
		} else if c == '.' {
			if hasDot || hasE {
				// only a single dot can be present in a valid number. scientific notation numbers cannot have a dot.
				return false
			}
			hasDot = true
		} else if c == 'e' || c == 'E' {
			// scientific notation numbers cannot have multiple e's or contain colon
			// e's must come after a digit
			if hasE || !hasDigit || hasColon {
				return false
			}
			hasE = true
			if j+1 < len(s) && (s[j+1] == '+' || s[j+1] == '-') {
				j++
			}
			hasDigit = false // force the exponent to have trailing digits
		} else if c == ':' {
			// base 60 numbers must have digits before colon, cannot contains anything but colons and numbers
			if !hasDigit || hasDot || hasE {
				return false
			}
			hasColon = true
			hasDigit = false // force digits after colon
			digitsSinceColon = 0
			colonValue = 0
		} else {
			return false
		}
	}

	return hasDigit
}

func checkBaseFormat(s string, start int, base byte) bool {
	seenDigit := false
	for j := start; j < len(s); j++ {
		c := s[j]
		if c == '_' {
			continue
		}

		isValid := false
		switch base {
		case 'b':
			isValid = (c == '0' || c == '1')
		case 'x':
			isValid = isHex(c)
		case 'o':
			isValid = (c >= '0' && c <= '7')
		}

		if !isValid {
			return false
		}
		seenDigit = true
	}
	return seenDigit
}

func isHex(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isYamlTimestamp(s string) bool {
	if len(s) < 10 {
		return false
	}

	if s[4] != '-' || s[7] != '-' {
		return false
	}

	for i := range 10 {
		if i == 4 || i == 7 {
			continue
		}
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	if len(s) == 10 {
		return true
	}

	c := s[10]
	if c == 'T' || c == 't' || c == ' ' {
		return true
	}

	return false
}
