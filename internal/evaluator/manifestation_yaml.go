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

func ManifestYaml(b []byte, value Value, ctx Context, config YamlManifestConfig) ([]byte, error) {
	return manifestYaml(value, ctx, b, 0, config)
}

func manifestYaml(value Value, ctx Context, b []byte, indentLevel int, config YamlManifestConfig) ([]byte, error) {
	value, err := value.Eval(ctx)
	if err != nil {
		return nil, err
	}

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()

		if config.FormatIntegers && data == math.Floor(data) {
			return strconv.AppendFloat(b, data, 'f', 0, 64), nil
		}
		return strconv.AppendFloat(b, data, 'f', -1, 64), nil
	case ValueTypeNull:
		return append(b, "null"...), nil
	case ValueTypeBool:
		if value.Bool() {
			return append(b, "true"...), nil
		}
		return append(b, "false"...), nil
	case ValueTypeString:
		data := value.String(ctx)

		n := len(data)
		if n == 0 {
			return append(b, `""`...), nil
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
				b = writeYamlString(b, data, true, false, config.Modern)
				return b, nil
			}

			b = writeYamlString(b, data, false, true, config.Modern)
			return b, nil
		}

		b = append(b, '|')
		if config.UseBlockScalars {
			firstByte := data[0]
			if firstByte == ' ' || /* data[0] == '\t' || */ firstByte == '\n' {
				b = append(b, yamlIndentNumber...)
			}

			if lastByte != '\n' {
				b = append(b, '-')
			} else if n >= 2 && data[n-2] == '\n' || n == 1 {
				b = append(b, '+')
			}

			if firstByte == '\n' {
				data = data[1:] // prefix trim
				n--             // update length
				if n == 0 {
					return b, nil
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

			b = append(b, '\n')
			if line != "" || !config.UseBlockScalars {
				b = writeYamlIndent(b, indentLevel+1)
				b = append(b, line...)
			}
		}

		return b, nil
	case ValueTypeArray:
		data := value.Array(ctx)
		if len(data) == 0 {
			return append(b, "[]"...), nil
		}
		for i, v := range data {
			v, err := v.Eval(ctx)
			if err != nil {
				return nil, err
			}

			if i != 0 {
				b = append(b, '\n')
				b = writeYamlIndent(b, indentLevel)
			}
			b = append(b, '-')

			if v.IsArray() && len(v.Array(ctx)) > 0 {
				b = append(b, '\n')
				b = writeYamlIndent(b, indentLevel+1)
			} else {
				b = append(b, ' ')
			}

			nextIndentLevel := indentLevel
			switch v.Type() {
			case ValueTypeArray, ValueTypeObject:
				nextIndentLevel++
			}

			b, err = manifestYaml(v, ctx, b, nextIndentLevel, config)
			if err != nil {
				return nil, err
			}
		}
		return b, nil
	case ValueTypeObject:
		obj := value.Object(ctx)
		plans := CompileObjectPlanEx(obj, ctx, config.NaturalSort)
		if len(plans) == 0 {
			return append(b, "{}"...), nil
		}

		subCtx := ctx
		subCtx.Self = value

		hasWritten := false
		for _, p := range plans {
			if p.IsHidden() {
				continue
			}
			if hasWritten {
				b = append(b, '\n')
				b = writeYamlIndent(b, indentLevel)
			}

			keyStr := ctx.State.Interner.Get(p.KeyId)
			b = writeYamlString(b, keyStr, config.QuoteKeys, config.Modern, config.Modern)

			b = append(b, ':')

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return nil, err
			}

			nextIndentLevel := indentLevel

			if fieldValue.IsArray() && len(fieldValue.Array(subCtx)) > 0 {
				b = append(b, '\n')
				b = writeYamlIndent(b, indentLevel)
				if config.IndentArrayInObjects {
					b = writeYamlIndent(b, 1)
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
					b = append(b, '\n')
					b = writeYamlIndent(b, indentLevel+1)
					nextIndentLevel++
				} else {
					b = append(b, ' ')
				}
			} else {
				b = append(b, ' ')
			}

			b, err = manifestYaml(fieldValue, subCtx, b, nextIndentLevel, config)
			if err != nil {
				return nil, err
			}
			hasWritten = true
		}

		return b, nil
	}
}

const (
	yamlIndentSpaces = 2
)

var (
	yamlIndentNumber = strconv.Itoa(yamlIndentSpaces)
)

func writeYamlIndent(b []byte, indentLevel int) []byte {
	// 64 spaces
	const maxIndentString = "                                                                "

	totalSpaces := indentLevel * yamlIndentSpaces

	for totalSpaces > 0 {
		// If the remaining spaces fit in our pre-allocated string,
		// slice it, write it, and we are done!
		if totalSpaces <= len(maxIndentString) {
			b = append(b, maxIndentString[:totalSpaces]...)
			break
		}

		// Otherwise, write the max chunk of 64 spaces and subtract it
		b = append(b, maxIndentString...)
		totalSpaces -= len(maxIndentString)
	}

	return b
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

func writeYamlString(b []byte, s string, forceQuotes, preferSingleQuotes, modern bool) []byte {
	if len(s) == 0 {
		if preferSingleQuotes {
			return append(b, "''"...)
		}
		return append(b, `""`...)
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
				b = append(b, '"')
				b = append(b, s...)
				b = append(b, '"')
				return b
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
		return append(b, s...)
	}

	if !useSingle {
		return writeJsonString(b, s)
	}

	if !hasSingleQuote {
		b = append(b, '\'')
		b = append(b, s...)
		b = append(b, '\'')
		return b
	}

	// escape single quotes
	b = append(b, '\'')
	remaining := s
	for {
		idx := strings.IndexByte(remaining, '\'')
		if idx == -1 {
			// No more quotes found, write the rest of the string
			b = append(b, remaining...)
			break
		}
		b = append(b, remaining[:idx]...)
		b = append(b, "''"...)
		remaining = remaining[idx+1:]
	}
	b = append(b, '\'')
	return b
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
