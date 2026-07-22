package stdlib

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
)

func std_format(args []evaluator.NamedValue, ctx evaluator.Context) (evaluator.Value, error) {

	format, err := args[0].EvalString(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	arg, err := args[1].Eval(ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}

	str, err := formatString(format, arg, ctx)
	if err != nil {
		return evaluator.ValueNone, err
	}
	return evaluator.MakeString(str, ctx), nil
}

const (
	FormatFlagAlternate   = 0x01 // Binary 00000001 (#)
	FormatFlagLeftJustify = 0x02 // Binary 00000010 (-)
	FormatFlagForceSign   = 0x04 // Binary 00000100 (+)
	FormatFlagSpaceSign   = 0x08 // Binary 00001000 ( )
	FormatFlagZeroPad     = 0x10 // Binary 00010000 (0)
)

// Sprintf formats a string using Python-style format specifiers.
// It supports:
// - Positional args: Sprintf("Value: %d", 10)
// - Named args:      Sprintf("Value: %(val)d", map[string]any{"val": 10})
// - Flags:           %#0- +
// - Width/Prec:      %10.5f, %*.2f (dynamic width), %.*f (dynamic prec)
func formatString(str string, data evaluator.Value, ctx evaluator.Context) (string, error) {
	n := len(str)

	sp, b := evaluator.GetBufferSlice(n + 128)

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
		argBuf := [1]evaluator.Value{data}
		args = argBuf[:]
	}

	argIdx := 0
	for i := 0; i < n; {

		if str[i] != '%' {
			next := strings.IndexByte(str[i:], '%')
			if next == -1 {
				b = append(b, str[i:]...)
				break
			}
			b = append(b, str[i:i+next]...)
			i += next
		}

		if i < n && str[i] == '%' {
			i++
		}

		if i >= n {
			return "", fmt.Errorf("incomplete format string")
		}

		// Handle "%%" (Literal Percent)
		if str[i] == '%' {
			b = append(b, '%')
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

		var flags uint8

	ParseFlags:
		for i < n {
			switch str[i] {
			case '#':
				flags |= FormatFlagAlternate
			case '0':
				flags |= FormatFlagZeroPad
			case '-':
				flags |= FormatFlagLeftJustify
			case ' ':
				flags |= FormatFlagSpaceSign
			case '+':
				flags |= FormatFlagForceSign
			default:
				break ParseFlags
			}
			i++
		}

		// 3. Parse Width
		widthVal := -1

		if i < n && str[i] == '*' {
			// Dynamic Width
			if useNamed {
				return "", fmt.Errorf("* width not supported with mapping")
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
			hasWidth := false
			w := 0
			for i < n && str[i] >= '0' && str[i] <= '9' {
				w = w*10 + int(str[i]-'0')
				hasWidth = true
				i++
			}
			if hasWidth {
				widthVal = w
			}
		}

		// 4. Parse Precision
		precVal := -1
		if i < n && str[i] == '.' {
			i++
			if i < n && str[i] == '*' {
				// Dynamic Precision
				if useNamed {
					return "", fmt.Errorf("* precision not supported with mapping")
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
				// Inline ASCII to Integer
				hasPrec := false
				p := 0
				for i < n && str[i] >= '0' && str[i] <= '9' {
					p = p*10 + int(str[i]-'0')
					hasPrec = true
					i++
				}
				// If a dot is present but no digits (e.g. "%.f"), precision is explicitly 0
				if hasPrec {
					precVal = p
				} else {
					precVal = 0
				}
			}
		}

		// 5. length modifier, ignored in Jsonnet.
		for i < n && (str[i] == 'h' || str[i] == 'l' || str[i] == 'L') {
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
			keyId := ctx.State.Interner.Intern(key)
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

		currentArg, err := currentArg.Eval(ctx)
		if err != nil {
			return "", err
		}

		width := 0
		if widthVal != -1 {
			width = widthVal
		}

		// Negative width implies left-justify
		if width < 0 {
			width = -width
			flags |= FormatFlagLeftJustify
		}

		prec := precVal

		switch verb {
		case 's':

			strVal, err := currentArg.ToString(ctx)
			if err != nil {
				return "", err
			}

			if width == 0 && flags == 0 {
				b = append(b, strVal...)
				continue
			}

			// note: jsonnet doesnt support precision on strings
			b = writeFormatString(b, strVal, width, flags)

		case 'd', 'i', 'u': // Integer types

			if !currentArg.IsNumber() {
				return "", fmt.Errorf("format %%%c requires number", verb)
			}
			num := int64(currentArg.Number())

			if width == 0 && flags == 0 {
				b = strconv.AppendInt(b, num, 10)
				continue
			}

			// note: jsonnet doesnt support precision on integer types
			b = writeFormatInteger(b, num, width, flags)

		case 'o': // Octal

			if !currentArg.IsNumber() {
				return "", fmt.Errorf("format %%%c requires number", verb)
			}
			num := int64(currentArg.Number())

			if width == 0 && prec == -1 && flags == 0 {
				b = strconv.AppendInt(b, num, 8)
				continue
			}

			b = writeFormatOctal(b, num, width, prec, flags)

		case 'x', 'X': // Hex

			if !currentArg.IsNumber() {
				return "", fmt.Errorf("format %%%c requires number", verb)
			}
			num := int64(currentArg.Number())
			uppercase := verb == 'X'

			if width == 0 && prec == -1 && flags == 0 {
				startIdx := len(b)
				b = strconv.AppendInt(b, num, 16)
				if uppercase {
					toUppercase(b[startIdx:])
				}
				continue
			}

			b = writeFormatHex(b, num, width, flags, uppercase)

		case 'f', 'F', 'e', 'E', 'g', 'G': // Float types

			if !currentArg.IsNumber() {
				return "", fmt.Errorf("format %%%c requires number", verb)
			}

			var uppercase bool
			fmt := byte(verb)

			switch verb {
			case 'F':
				fmt = 'f'
				uppercase = true // Need this for NaN/Inf
			case 'E':
				uppercase = true
			case 'G':
				uppercase = true
			}

			if prec < 0 {
				prec = 6
			}

			num := currentArg.Number()
			if width == 0 && flags == 0 {
				startIdx := len(b)
				b = strconv.AppendFloat(b, num, fmt, prec, 64)
				if uppercase {
					toUppercase(b[startIdx:])
				}
				continue
			}

			b = writeFormatFloat(b, num, fmt, width, prec, flags, uppercase)

		case 'c':
			// Character
			var char rune
			switch {
			case currentArg.IsNumber():
				n := currentArg.Number()
				if n > codepointMax {
					return "", fmt.Errorf("invalid unicode codepoint, got %v", n)
				}
				char = rune(n)
			case currentArg.IsString():
				s := currentArg.String(ctx)
				if utf8.RuneCountInString(s) == 1 {
					char, _ = utf8.DecodeRuneInString(s)
					break
				}
				fallthrough
			default:
				return "", fmt.Errorf("format %%c requires integer or char")
			}

			if width == 0 && flags == 0 {
				b = utf8.AppendRune(b, char)
				continue
			}

			b = writeFormatChar(b, char, width, flags)

		default:
			return "", evaluator.MakeRuntimeError(fmt.Errorf("Unrecognised conversion type: %s", string(verb)))
		}
	}

	if !useNamed && argIdx < len(args) {
		return "", fmt.Errorf("not all arguments converted during string formatting")
	}

	res := string(b)

	evaluator.PutBufferSlice(sp, b)

	return res, nil
}

func writePad(b []byte, zeroPad bool, count int) []byte {
	const (
		spaces = "                                                                "
		zeros  = "0000000000000000000000000000000000000000000000000000000000000000"
	)

	var chunk string
	if zeroPad {
		chunk = zeros
	} else {
		chunk = spaces
	}

	for count > 0 {
		if count <= len(chunk) {
			b = append(b, chunk[:count]...)
			break
		}

		b = append(b, chunk...)
		count -= len(chunk)
	}
	return b
}

func writePadded(b []byte, content []byte, prefix string, padLen int, flags uint8) []byte {
	leftJustify := flags&FormatFlagLeftJustify != 0
	zeroPad := !leftJustify && flags&FormatFlagZeroPad != 0

	// right justified, space padded
	if !leftJustify && !zeroPad && padLen > 0 {
		b = writePad(b, false, padLen)
	}

	if prefix != "" {
		b = append(b, prefix...)
	}

	// right justified, zero padded
	if zeroPad && padLen > 0 {
		b = writePad(b, true, padLen)
	}

	b = append(b, content...)

	// left justified, space padded
	if leftJustify && padLen > 0 {
		b = writePad(b, false, padLen)
	}

	return b
}

//go:noinline
func writeFormatString(b []byte, s string, width int, flags uint8) []byte {
	padLen := width - utf8.RuneCountInString(s)

	content := unsafe.Slice(unsafe.StringData(s), len(s))

	flags &^= FormatFlagZeroPad
	return writePadded(b, content, "", padLen, flags)
}

//go:noinline
func writeFormatInteger(b []byte, num int64, width int, flags uint8) []byte {

	var prefix string
	u := uint64(num)

	if num < 0 {
		prefix = "-"
		u = -u
	} else if flags&FormatFlagForceSign != 0 {
		prefix = "+"
	} else if flags&FormatFlagSpaceSign != 0 {
		prefix = " "
	}

	if width == 0 {
		if prefix != "" {
			b = append(b, prefix...)
		}
		return strconv.AppendUint(b, u, 10)
	}

	var dst [64]byte
	res := strconv.AppendUint(dst[:0], u, 10)

	padLen := width - len(res) - len(prefix)
	return writePadded(b, res, prefix, padLen, flags)
}

//go:noinline
func writeFormatOctal(b []byte, num int64, width, prec int, flags uint8) []byte {

	u := uint64(num)
	if num < 0 {
		u = -u
	}

	if width == 0 && prec == -1 && flags&FormatFlagAlternate == 0 {
		return strconv.AppendUint(b, u, 8)
	}

	var dst [64]byte
	res := strconv.AppendUint(dst[:0], u, 8)

	// alternate flag for octal ensures it starts with '0'
	alt := flags&FormatFlagAlternate != 0 && res[0] != '0'
	forceSign := flags&FormatFlagForceSign != 0
	spaceSign := flags&FormatFlagSpaceSign != 0

	var prefix string
	switch {
	case num < 0 && alt:
		prefix = "-0"
	case num < 0:
		prefix = "-"

	case forceSign && alt:
		prefix = "+0"
	case forceSign:
		prefix = "+"

	case spaceSign && alt:
		prefix = " 0"
	case spaceSign:
		prefix = " "

	case alt:
		prefix = "0"
	}

	padLen := width - len(res) - len(prefix)
	return writePadded(b, res, prefix, padLen, flags)
}

//go:noinline
func writeFormatHex(b []byte, num int64, width int, flags uint8, uppercase bool) []byte {

	u := uint64(num)
	isNeg := num < 0
	if isNeg {
		u = -u
	}

	var dst [64]byte
	res := strconv.AppendUint(dst[:0], u, 16)
	if uppercase {
		toUppercase(res)
	}

	alt := flags&FormatFlagAlternate != 0
	forceSign := flags&FormatFlagForceSign != 0
	spaceSign := flags&FormatFlagSpaceSign != 0

	var prefix string
	switch {
	case isNeg && alt && uppercase:
		prefix = "-0X"
	case isNeg && alt:
		prefix = "-0x"
	case isNeg:
		prefix = "-"

	case forceSign && alt && uppercase:
		prefix = "+0X"
	case forceSign && alt:
		prefix = "+0x"
	case forceSign:
		prefix = "+"

	case spaceSign && alt && uppercase:
		prefix = " 0X"
	case spaceSign && alt:
		prefix = " 0x"
	case spaceSign:
		prefix = " "

	case alt && uppercase:
		prefix = "0X"
	case alt:
		prefix = "0x"
	}

	padLen := width - len(res) - len(prefix)
	return writePadded(b, res, prefix, padLen, flags)
}

//go:noinline
func writeFormatFloat(b []byte, num float64, fmt byte, width, prec int, flags uint8, uppercase bool) []byte {
	var dst [128]byte

	res := strconv.AppendFloat(dst[:0], num, fmt, prec, 64)

	if len(res) > 0 && (res[0] == '-' || res[0] == '+') {
		res = res[1:]
	}

	if uppercase {
		toUppercase(res)
	}

	alt := flags&FormatFlagAlternate != 0
	if alt && prec == 0 && !math.IsNaN(num) && !math.IsInf(num, 0) {

		if fmt == 'f' {
			res = append(res, '.')

		} else {
			var c byte = 'e'
			if uppercase {
				c = 'E'
			}
			idx := bytes.IndexByte(res, c)
			if idx != -1 {
				res = append(res, 0)
				copy(res[idx+1:], res[idx:])
				res[idx] = '.'
			}
		}

	}

	isNeg := math.Signbit(num)
	forceSign := flags&FormatFlagForceSign != 0
	spaceSign := flags&FormatFlagSpaceSign != 0

	var prefix string
	switch {
	case isNeg:
		prefix = "-"
	case forceSign:
		prefix = "+"
	case spaceSign:
		prefix = " "
	}

	padLen := width - len(res) - len(prefix)
	return writePadded(b, res, prefix, padLen, flags)
}

//go:noinline
func writeFormatChar(b []byte, c rune, width int, flags uint8) []byte {

	if width == 0 {
		return utf8.AppendRune(b, c)
	}

	var dst [utf8.UTFMax]byte
	n := utf8.EncodeRune(dst[:], c)

	padLen := width - 1
	flags &^= FormatFlagZeroPad
	return writePadded(b, dst[:n], "", padLen, flags)
}

func toUppercase(x []byte) {
	for i := range x {
		// If it's a lowercase letter (a-z), shift it to uppercase (A-Z)
		if x[i] >= 'a' && x[i] <= 'z' {
			x[i] -= 32
		}
	}
}
