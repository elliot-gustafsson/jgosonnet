package evaluator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

func ManifestYaml(b *[]byte, value Value, ctx Context, config YamlManifestConfig) error {
	var err error
	*b, err = manifestYaml(value, ctx, *b, 0, config)
	return err
}

func manifestYaml(value Value, ctx Context, buf []byte, indentLevel int, config YamlManifestConfig) ([]byte, error) {
	if ctx.Depth >= ctx.State.MaxStack {
		return buf, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
	}
	ctx.Depth++

	value, err := value.Eval(ctx)
	if err != nil {
		return buf, err
	}

	switch value.Type() {
	default:
		return buf, fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		if config.FormatIntegers && data == math.Floor(data) {
			return strconv.AppendFloat(buf, data, 'f', 0, 64), nil
		}
		return strconv.AppendFloat(buf, data, 'f', -1, 64), nil
	case ValueTypeNull:
		return append(buf, "null"...), nil
	case ValueTypeBool:
		if value.Bool() {
			return append(buf, "true"...), nil
		}
		return append(buf, "false"...), nil
	case ValueTypeString:
		data := value.String(ctx)

		n := len(data)
		if n == 0 {
			return append(buf, `""`...), nil
		}

		var multiline bool
		lastByte := data[n-1]

		if config.UseBlockScalars {
			// Check if the string contains newlines to determine if it can be a block scalar.
			// Since YAML block scalars do not safely preserve trailing whitespace on lines,
			// we must safely fall back to a quoted single-line string if any line ends with a space or tab.
			if idx := strings.IndexByte(data, '\n'); idx != -1 {
				multiline = true
				for {
					if idx > 0 && (data[idx-1] == ' ' || data[idx-1] == '\t') {
						multiline = false
						break
					}
					next := strings.IndexByte(data[idx+1:], '\n')
					if next == -1 {
						break
					}
					idx += 1 + next
				}
				if multiline && (lastByte == ' ' || lastByte == '\t') {
					multiline = false
				}
			}
		} else {
			multiline = lastByte == '\n'
		}

		if !multiline {
			if config.QuoteValues {
				return appendYamlString(buf, data, true, false, config.Modern), nil
			}
			return appendYamlString(buf, data, false, true, config.Modern), nil
		}

		buf = append(buf, '|')
		if config.UseBlockScalars {
			firstByte := data[0]
			if firstByte == ' ' || firstByte == '\n' {
				buf = append(buf, yamlIndentNumber...)
			}

			if lastByte != '\n' {
				buf = append(buf, '-')
			} else if n >= 2 && data[n-2] == '\n' || n == 1 {
				buf = append(buf, '+')
			}

			if firstByte == '\n' {
				data = data[1:] // prefix trim
				n--             // update length
				if n == 0 {
					return buf, nil
				}
			}
		}

		start := 0
		end := n
		if end > 0 && data[end-1] == '\n' {
			end--
		}

		// Precompute "\n" + indent for block lines
		var nlIndentBuf [65]byte
		var nlIndent []byte
		indentBytes := (indentLevel + 1) * yamlIndentSpaces
		if indentBytes+1 <= len(nlIndentBuf) {
			nlIndent = nlIndentBuf[:indentBytes+1]
			nlIndent[0] = '\n'
			copy(nlIndent[1:], yamlMaxIndentString[:indentBytes])
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

			if line != "" || !config.UseBlockScalars {
				if nlIndent != nil {
					buf = append(buf, nlIndent...)
				} else {
					buf = append(buf, '\n')
					buf = appendYamlIndent(buf, indentLevel+1)
				}
				buf = append(buf, line...)
			} else {
				buf = append(buf, '\n')
			}
		}

		return buf, nil
	case ValueTypeArray:
		data := value.Array()
		if len(data) == 0 {
			return append(buf, "[]"...), nil
		}
		var nlIndentBuf [65]byte
		var nlIndent []byte
		indentBytes := indentLevel * yamlIndentSpaces
		if indentBytes+1 <= len(nlIndentBuf) {
			nlIndent = nlIndentBuf[:indentBytes+1]
			nlIndent[0] = '\n'
			copy(nlIndent[1:], yamlMaxIndentString[:indentBytes])
		}
		for i, v := range data {
			v, err := v.Eval(ctx)
			if err != nil {
				return buf, err
			}

			if i != 0 {
				if nlIndent != nil {
					buf = append(buf, nlIndent...)
				} else {
					buf = append(buf, '\n')
					buf = appendYamlIndent(buf, indentLevel)
				}
			}
			buf = append(buf, '-')

			if v.IsArray() && len(v.Array()) > 0 {
				buf = append(buf, '\n')
				buf = appendYamlIndent(buf, indentLevel+1)
			} else {
				buf = append(buf, ' ')
			}

			nextIndentLevel := indentLevel
			switch v.Type() {
			case ValueTypeArray, ValueTypeObject:
				nextIndentLevel++
			}

			buf, err = manifestYaml(v, ctx, buf, nextIndentLevel, config)
			if err != nil {
				return buf, err
			}
		}
		return buf, nil
	case ValueTypeObject:
		obj := value.Object()
		evalCtx := ctx
		evalCtx.Self = value
		err := runAssertions(obj, evalCtx)
		if err != nil {
			return buf, err
		}
		layers := obj.GetLayers(ctx)
		if len(layers) == 1 && layerAt(layers, 0).Template != nil {
			return manifestYamlStaticObject(value, obj, layerAt(layers, 0).Template, ctx, buf, indentLevel, config)
		}
		plans := CompileObjectPlanEx(obj, ctx, config.NaturalSort)
		return manifestYamlObject(value, obj, plans, ctx, buf, indentLevel, config)
	}
}

func manifestYamlStaticObject(value Value, obj *Object, tmpl *LayerTemplate, ctx Context, buf []byte, indentLevel int, config YamlManifestConfig) ([]byte, error) {
	if len(tmpl.Keys) == 0 {
		return append(buf, "{}"...), nil
	}
	if tmpl.VisibleCount == 0 {
		return buf, nil
	}

	subCtx := ctx
	subCtx.Self = value
	inter := ctx.State.Interner

	var order []uint32
	if config.NaturalSort {
		order = tmpl.GetNaturalVisible(inter)
	} else {
		order = tmpl.GetSortedVisible(inter)
	}

	var nlIndentBuf [65]byte
	var nlIndent []byte
	indentBytes := indentLevel * yamlIndentSpaces
	if indentBytes+1 <= len(nlIndentBuf) {
		nlIndent = nlIndentBuf[:indentBytes+1]
		nlIndent[0] = '\n'
		copy(nlIndent[1:], yamlMaxIndentString[:indentBytes])
	}

	keys := tmpl.Keys
	for i, fIdx := range order {
		if i > 0 {
			if nlIndent != nil {
				buf = append(buf, nlIndent...)
			} else {
				buf = append(buf, '\n')
				buf = appendYamlIndent(buf, indentLevel)
			}
		}

		keyId := keys[fIdx]
		buf = appendYamlString(buf, inter.Get(keyId), config.QuoteKeys, config.Modern, config.Modern)
		buf = append(buf, ':')

		fieldValue, err := getSingleLayerFieldValue(obj, tmpl, keyId, int(fIdx), subCtx)
		if err != nil {
			return buf, err
		}

		buf, err = manifestYamlFieldValue(fieldValue, subCtx, buf, indentLevel, nlIndent, config)
		if err != nil {
			return buf, err
		}
	}

	return buf, nil
}

func manifestYamlFieldValue(fieldValue Value, subCtx Context, buf []byte, indentLevel int, nlIndent []byte, config YamlManifestConfig) ([]byte, error) {
	nextIndentLevel := indentLevel

	if fieldValue.IsArray() && len(fieldValue.Array()) > 0 {
		if nlIndent != nil {
			buf = append(buf, nlIndent...)
		} else {
			buf = append(buf, '\n')
			buf = appendYamlIndent(buf, indentLevel)
		}
		if config.IndentArrayInObjects {
			buf = append(buf, "  "...)
			nextIndentLevel++
		}
		return manifestYaml(fieldValue, subCtx, buf, nextIndentLevel, config)
	}

	if fieldValue.IsObject() {
		if subCtx.Depth >= subCtx.State.MaxStack {
			return buf, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
		}
		childCtx := subCtx
		childCtx.Depth++
		childObj := fieldValue.Object()
		childEvalCtx := childCtx
		childEvalCtx.Self = fieldValue
		if err := runAssertions(childObj, childEvalCtx); err != nil {
			return buf, err
		}

		childLayers := childObj.GetLayers(childCtx)
		if len(childLayers) == 1 && layerAt(childLayers, 0).Template != nil {
			childTmpl := layerAt(childLayers, 0).Template
			if childTmpl.VisibleCount > 0 {
				buf = append(buf, '\n')
				buf = appendYamlIndent(buf, indentLevel+1)
				nextIndentLevel++
			} else {
				buf = append(buf, ' ')
			}
			return manifestYamlStaticObject(fieldValue, childObj, childTmpl, childCtx, buf, nextIndentLevel, config)
		}

		childPlans := CompileObjectPlanEx(childObj, childCtx, config.NaturalSort)
		hasFields := false
		for j := range childPlans {
			if !childPlans[j].IsHidden() {
				hasFields = true
				break
			}
		}

		if hasFields {
			buf = append(buf, '\n')
			buf = appendYamlIndent(buf, indentLevel+1)
			nextIndentLevel++
		} else {
			buf = append(buf, ' ')
		}

		return manifestYamlObject(fieldValue, childObj, childPlans, childCtx, buf, nextIndentLevel, config)
	}

	buf = append(buf, ' ')
	return manifestYaml(fieldValue, subCtx, buf, nextIndentLevel, config)
}

func manifestYamlObject(value Value, obj *Object, plans []FieldPlan, ctx Context, buf []byte, indentLevel int, config YamlManifestConfig) ([]byte, error) {
	if len(plans) == 0 {
		return append(buf, "{}"...), nil
	}

	subCtx := ctx
	subCtx.Self = value

	var nlIndentBuf [65]byte
	var nlIndent []byte
	indentBytes := indentLevel * yamlIndentSpaces
	if indentBytes+1 <= len(nlIndentBuf) {
		nlIndent = nlIndentBuf[:indentBytes+1]
		nlIndent[0] = '\n'
		copy(nlIndent[1:], yamlMaxIndentString[:indentBytes])
	}

	hasWritten := false
	for i := range plans {
		p := &plans[i]
		if p.IsHidden() {
			continue
		}
		if hasWritten {
			if nlIndent != nil {
				buf = append(buf, nlIndent...)
			} else {
				buf = append(buf, '\n')
				buf = appendYamlIndent(buf, indentLevel)
			}
		}

		keyStr := ctx.State.Interner.Get(p.KeyId)
		buf = appendYamlString(buf, keyStr, config.QuoteKeys, config.Modern, config.Modern)

		buf = append(buf, ':')

		fieldValue, err := p.GetValue(obj, subCtx)
		if err != nil {
			return buf, err
		}

		buf, err = manifestYamlFieldValue(fieldValue, subCtx, buf, indentLevel, nlIndent, config)
		if err != nil {
			return buf, err
		}
		hasWritten = true
	}

	return buf, nil
}

const (
	yamlIndentSpaces    = 2
	yamlMaxIndentString = "                                                                "
)

var (
	yamlIndentNumber = strconv.Itoa(yamlIndentSpaces)
)

func appendYamlIndent(b []byte, indentLevel int) []byte {
	totalSpaces := indentLevel * yamlIndentSpaces

	for totalSpaces > 0 {
		if totalSpaces <= len(yamlMaxIndentString) {
			b = append(b, yamlMaxIndentString[:totalSpaces]...)
			break
		}

		b = append(b, yamlMaxIndentString...)
		totalSpaces -= len(yamlMaxIndentString)
	}
	return b
}

func isSpaceByte(b byte) bool {
	return b == ' ' || (b >= '\t' && b <= '\r') || b == 0x85 || b == 0xA0
}

func yamlReserved(s string) bool {
	if len(s) > 5 {
		return false
	}
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

func appendYamlString(b []byte, s string, forceQuotes, preferSingleQuotes, modern bool) []byte {
	n := len(s)
	if n == 0 {
		if preferSingleQuotes {
			return append(b, "''"...)
		}
		return append(b, `""`...)
	}

	needsQuotes := forceQuotes
	useSingle := preferSingleQuotes

	// prefix and suffix checks
	if !needsQuotes {
		c0 := s[0]
		if isSpaceByte(c0) || isSpaceByte(s[n-1]) {
			needsQuotes = true
		} else {
			// structural indicators at the start
			switch c0 {
			case '[', ']', '{', '}', ',', '#', '&', '*', '!', '|', '>', '\'', '"', '%', '@', '`':
				needsQuotes = true
			case '-', '?', ':':
				if n == 1 || s[1] == ' ' || s[1] == '\t' || s[1] == '\n' {
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
	if !(needsQuotes && !useSingle) {
		if modern {
			for i := 0; i < n; i++ {
				c := s[i]

				// force quotes on control chars
				if c < 0x20 || c == 0x7F {
					useSingle = false
					needsQuotes = true
					break
				}

				if c == '\'' {
					hasSingleQuote = true
				} else if !needsQuotes {
					if c == ':' {
						if i == n-1 || s[i+1] == ' ' || s[i+1] == '\n' {
							needsQuotes = true
							if !useSingle {
								break
							}
						}
					} else if c == ' ' && i+1 < n && s[i+1] == '#' {
						needsQuotes = true
						if !useSingle {
							break
						}
					}
				}
			}
		} else {
			for i := 0; i < n; i++ {
				c := s[i]

				if c < 0x20 || c == 0x7F {
					useSingle = false
					needsQuotes = true
					break
				}

				if c == '\'' {
					hasSingleQuote = true
				}

				if !needsQuotes {
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
	}

	if !needsQuotes {
		return append(b, s...)
	}

	if !useSingle {
		return appendJsonString(b, s)
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
			b = append(b, remaining...)
			break
		}
		b = append(b, remaining[:idx]...)
		b = append(b, '\'', '\'')
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
	if first == '+' || first == '-' {
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
				return false
			}
			hasDot = true
		} else if c == 'e' || c == 'E' {
			if hasE || !hasDigit || hasColon {
				return false
			}
			hasE = true
			if j+1 < len(s) && (s[j+1] == '+' || s[j+1] == '-') {
				j++
			}
			hasDigit = false
		} else if c == ':' {
			if !hasDigit || hasDot || hasE {
				return false
			}
			hasColon = true
			hasDigit = false
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
	if len(s) < 10 || s[0] < '0' || s[0] > '9' {
		return false
	}

	if s[4] != '-' || s[7] != '-' {
		return false
	}

	for i := 1; i < 10; i++ {
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
