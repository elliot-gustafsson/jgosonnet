package stdlib

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

var bufPool = sync.Pool{
	New: func() any {
		// Start with a small capacity to avoid wasting RAM on tiny strings
		return bytes.NewBuffer(make([]byte, 0, 64))
	},
}

func std_format(args []evaluator.Value, ctx evaluator.Context) (evaluator.Value, error) {
	if len(args) != 2 {
		return evaluator.Value{}, fmt.Errorf("unexpected amount of arguments passed to std.format: %d, expected 2", len(args))
	}

	format := args[0]
	arg := args[1]

	if !format.IsString() {
		return evaluator.Value{}, fmt.Errorf("unexpected type %s, expected string (std.format arg 0)", format.Type().String())
	}

	// v, err := manifestValue(arg, ctx)
	// if err != nil {
	// 	return evaluator.Value{}, err
	// }

	str, err := formatString(format.String(ctx), arg, ctx)
	if err != nil {
		return evaluator.Value{}, err
	}
	return evaluator.MakeString(str, ctx), nil
}

// Sprintf formats a string using Python-style format specifiers.
// It supports:
// - Positional args: Sprintf("Value: %d", 10)
// - Named args:      Sprintf("Value: %(val)d", map[string]any{"val": 10})
// - Flags:           %#0- +
// - Width/Prec:      %10.5f, %*.2f (dynamic width), %.*f (dynamic prec)
func formatString(str string, data evaluator.Value, ctx evaluator.Context) (string, error) {

	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	// pre greow to hopefully avoid reallocs
	if cap := buf.Cap(); cap < len(str) {
		buf.Grow(len(str) + len(str)/5 - cap)
	}

	// 1. Normalize input data into List or Map
	var args []evaluator.Value
	// var dict map[string]any
	var dict *evaluator.Object
	useNamed := false

	// switch v := data.(type) {
	// case []any:
	// 	args = v
	// case map[string]any:
	// 	dict = v
	// 	useNamed = true
	// default:
	// 	args = []any{v}
	// }

	switch data.Type() {
	default:
		return "", fmt.Errorf("unsupported data type passed to format: %s, expected string, array, object", data.Type().String())
	case evaluator.ValueTypeArray:
		args = data.Array(ctx)
	case evaluator.ValueTypeObject:
		dict = data.Object(ctx)
		useNamed = true
	case evaluator.ValueTypeString, evaluator.ValueTypeBool, evaluator.ValueTypeNumber, evaluator.ValueTypeNull:
		args = []evaluator.Value{data}
	}

	i := 0
	n := len(str)
	argIdx := 0

	for i < n {
		// Find the next '%' from the current position
		remain := str[i:]
		idx := strings.IndexByte(remain, '%')

		if idx == -1 {
			// No more verbs found, write and exit
			buf.WriteString(remain)
			break
		}

		// Write the chunk of text before the '%'
		buf.WriteString(remain[:idx])
		i += idx

		// We are now standing on '%'. Advance past it.
		i++
		if i >= n {
			return "", fmt.Errorf("incomplete format string")
		}

		// Handle "%%" (Literal Percent)
		if str[i] == '%' {
			buf.WriteByte('%')
			i++
			continue
		}

		// 1. Parse Mapping Key: %(key)
		var key string
		hasKey := false
		if str[i] == '(' {
			end := strings.IndexByte(str[i:], ')')
			if end == -1 {
				return "", fmt.Errorf("incomplete format key")
			}
			key = str[i+1 : i+end]
			hasKey = true
			i += end + 1 // Move past ')'
		}

		// Validate Mode (Named vs Positional)
		if useNamed && !hasKey {
			return "", fmt.Errorf("format requires a mapping (%%(key)s) when a dictionary is passed")
		}
		// if !useNamed && hasKey {
		// 	return "", fmt.Errorf("format requires a tuple/list (no named keys) when a list is passed")
		// }

		// 2. Parse Flags
		var flags strings.Builder
		for i < n && strings.ContainsRune("#0- +", rune(str[i])) {
			flags.WriteString(string(str[i]))
			i++
		}

		// 3. Parse Width
		widthVal := -1
		widthStr := ""

		if i < n && str[i] == '*' {
			// Dynamic Width
			if useNamed {
				return "", fmt.Errorf("width '*' cannot be used with dictionary arguments")
			}
			if argIdx >= len(args) {
				return "", fmt.Errorf("not enough arguments for format string")
			}
			if v := args[argIdx]; v.IsNumber() {
				widthVal = int(v.Number())
			} else {
				return "", fmt.Errorf("width requires integer, got %s", v.Type().String())
			}
			argIdx++
			i++
		} else {
			// Static Width
			start := i
			for i < n && str[i] >= '0' && str[i] <= '9' {
				i++
			}
			if i > start {
				widthStr = str[start:i]
			}
		}

		// 4. Parse Precision
		precVal := -1
		precStr := ""

		if i < n && str[i] == '.' {
			i++ // Skip '.'
			if i < n && str[i] == '*' {
				// Dynamic Precision
				if useNamed {
					return "", fmt.Errorf("precision '*' cannot be used with dictionary arguments")
				}
				if argIdx >= len(args) {
					return "", fmt.Errorf("not enough arguments for format string")
				}
				if v := args[argIdx]; v.IsNumber() {
					precVal = int(v.Number())
				} else {
					return "", fmt.Errorf("precision requires integer, got %s", v.Type().String())
				}
				argIdx++
				i++
			} else {
				// Static Precision
				start := i
				for i < n && str[i] >= '0' && str[i] <= '9' {
					i++
				}
				precStr = str[start:i]
			}
		}

		// 5. Length Modifier (Ignored in Go/Jsonnet, e.g. 'l', 'h')
		for i < n && strings.ContainsRune("hlL", rune(str[i])) {
			i++
		}

		// 6. Parse Verb
		if i >= n {
			return "", fmt.Errorf("incomplete format string")
		}
		verb, size := utf8.DecodeRuneInString(str[i:])
		i += size

		// --- RETRIEVE ARGUMENT ---
		var currentArg evaluator.Value

		if useNamed {
			// val, ok := dict[key]
			// if !ok {
			// 	return "", fmt.Errorf("key '%s' not found", key)
			// }
			keyId := ctx.Interner.Intern(key)
			subCtx := ctx
			subCtx.Self = data
			val, _, err := dict.GetField(keyId, subCtx)
			if err != nil {
				return "", err
			}
			if val.IsNone() {
				return "", fmt.Errorf("key '%s' not found", key)
			}
			currentArg = val
		} else {
			if argIdx >= len(args) {
				return "", fmt.Errorf("not enough arguments for format string")
			}
			currentArg = args[argIdx]
			argIdx++
		}

		err := evaluator.EvaluateValueStrict(&currentArg, ctx)
		if err != nil {
			return "", err
		}

		// Rebuild format string
		fmtBuilder := strings.Builder{}
		fmtBuilder.WriteByte('%')
		fmtBuilder.WriteString(flags.String())

		if widthVal != -1 {
			fmtBuilder.WriteString(strconv.Itoa(widthVal))
		} else {
			fmtBuilder.WriteString(widthStr)
		}

		if precVal != -1 {
			fmtBuilder.WriteByte('.')
			fmtBuilder.WriteString(strconv.Itoa(precVal))
		} else if precStr != "" {
			fmtBuilder.WriteByte('.')
			fmtBuilder.WriteString(precStr)
		}

		switch verb {
		case 's', 'r':
			// %s: String
			fmtBuilder.WriteByte('s')
			// var strVal string
			// if currentArg == nil {
			// 	strVal = "null" // Jsonnet style null
			// } else if s, ok := currentArg.(string); ok {
			// 	strVal = s
			// } else {
			// 	strVal = fmt.Sprint(currentArg) // Fallback for bool/numbers
			// }

			strVal, err := currentArg.ToString(ctx)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(buf, fmtBuilder.String(), strVal)

		case 'd', 'i', 'u', 'o', 'x', 'X':
			// Integer types

			// Go does not support %i or %u, map to d
			if verb == 'i' || verb == 'u' {
				fmtBuilder.WriteRune('d')
			} else {
				fmtBuilder.WriteRune(verb)
			}

			// Jsonnet numbers are float64. Convert to Int64.
			if currentArg.IsNumber() {
				fmt.Fprintf(buf, fmtBuilder.String(), int64(currentArg.Number()))
			} else {
				return "", fmt.Errorf("format %%%c requires integer", verb)
			}

		case 'f', 'F', 'e', 'E', 'g', 'G':
			// Float types
			fmtBuilder.WriteRune(verb)
			if currentArg.IsNumber() {
				n := currentArg.Number()
				if verb == 'f' || verb == 'F' {
					// Add small epsilon to "fix" IEEE 754
					n += 1e-9
				}
				fmt.Fprintf(buf, fmtBuilder.String(), n) // +1e-9
			} else {
				return "", fmt.Errorf("format %%%c requires number", verb)
			}

		case 'c':
			// Character
			fmtBuilder.WriteByte('c')
			if currentArg.IsNumber() {
				n := currentArg.Number()
				if n > codepointMax {
					return "", fmt.Errorf("invalid unicode codepoint, got %v", n)
				}
				fmt.Fprintf(buf, fmtBuilder.String(), rune(n))
			} else if currentArg.IsString() && len(currentArg.String(ctx)) == 1 {
				r, _ := utf8.DecodeRuneInString(currentArg.String(ctx))
				fmt.Fprintf(buf, fmtBuilder.String(), r)
			} else {
				return "", fmt.Errorf("format %%c requires integer or char")
			}

		default:
			return "", fmt.Errorf("unsupported format character '%c'", verb)
		}
	}

	if !useNamed && argIdx < len(args) {
		return "", fmt.Errorf("not all arguments converted during string formatting")
	}

	return buf.String(), nil
}
