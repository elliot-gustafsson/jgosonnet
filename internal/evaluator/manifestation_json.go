package evaluator

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"sync"
	"unsafe"
)

const (
	defaultManifestBufCap = 4 * 1024
	maxPooledManifestCap  = 64 * 1024 * 1024
)

var manifestBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, defaultManifestBufCap)
		return &b
	},
}

func GetManifestBuf() *[]byte {
	bp := manifestBufPool.Get().(*[]byte)
	*bp = (*bp)[:0]
	return bp
}

func PutManifestBuf(bp *[]byte) {
	if cap(*bp) <= maxPooledManifestCap {
		*bp = (*bp)[:0]
		manifestBufPool.Put(bp)
	}
}

type JsonManifestConfig struct {
	IndentStep  string
	Newline     string
	KeyValSep   string
	SpaceComma  bool
	Python      bool
	StrictFloat bool
}

// Pre-defined configurations matching Jsonnet's standard library.
var (
	JsonConfigPretty   = &JsonManifestConfig{IndentStep: "    ", Newline: "\n", KeyValSep: ": ", SpaceComma: false, StrictFloat: true}
	JsonConfigMinified = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ":", SpaceComma: false, StrictFloat: true}
	JsonConfigToString = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ": ", SpaceComma: true}
	JsonConfigPython   = &JsonManifestConfig{IndentStep: "", Newline: "", KeyValSep: ": ", SpaceComma: true, Python: true}
)

func ManifestJson(b *[]byte, value Value, ctx Context, config *JsonManifestConfig) error {
	var err error
	*b, err = manifestJson(value, ctx, *b, 0, config)
	return err
}

//go:noinline
func growDouble(b []byte, need int) []byte {
	grow := cap(b)
	if grow < need {
		grow = need
	}
	if grow < 256 {
		grow = 256
	}
	return slices.Grow(b, grow)
}

func manifestJson(value Value, ctx Context, b []byte, indentLevel int, config *JsonManifestConfig) ([]byte, error) {
	if ctx.Depth >= ctx.State.MaxStack {
		return b, MakeRuntimeError(fmt.Errorf("max manifest depth exceeded, possible infinite recursion"))
	}
	ctx.Depth++

	hasNewline := config.Newline != ""

	value, err := value.Eval(ctx)
	if err != nil {
		return b, err
	}

	switch value.Type() {
	default:
		return b, fmt.Errorf("unhandled value type: %s", value.Type().String())
	case ValueTypeNumber:
		data := value.Number()
		if config.StrictFloat {
			return strconv.AppendFloat(b, data, 'f', -1, 64), nil
		}
		return unparseNumber(b, data), nil
	case ValueTypeNull:
		if config.Python {
			return append(b, "None"...), nil
		}
		return append(b, "null"...), nil
	case ValueTypeBool:
		if config.Python {
			if value.Bool() {
				return append(b, "True"...), nil
			}
			return append(b, "False"...), nil
		}
		if value.Bool() {
			return append(b, "true"...), nil
		}
		return append(b, "false"...), nil
	case ValueTypeString:
		data := value.String(ctx)
		return appendJsonString(b, data), nil
	case ValueTypeArray:
		data := value.Array()
		if len(data) == 0 {
			if config.SpaceComma && !config.Python {
				return append(b, "[ ]"...), nil
			}

			if hasNewline {
				b = append(b, '[')
				b = append(b, config.Newline...)
				b = append(b, config.Newline...)
				b = appendIndent(b, indentLevel, config.IndentStep)
				b = append(b, ']')
				return b, nil
			}

			return append(b, "[]"...), nil
		}

		b = append(b, '[')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}

		b = append(b, config.Newline...)

		// Precompute separator bytes for subsequent elements: "," + newline/space + indent
		var sepBuf [128]byte
		var sep []byte
		if nextIndentLevel*len(config.IndentStep)+len(config.Newline)+2 <= len(sepBuf) {
			sep = sepBuf[:0]
			sep = append(sep, ',')
			if hasNewline {
				sep = append(sep, config.Newline...)
			} else if config.SpaceComma {
				sep = append(sep, ' ')
			}
			sep = appendIndent(sep, nextIndentLevel, config.IndentStep)
		}

		for i, v := range data {
			v, err := v.Eval(ctx)
			if err != nil {
				return b, err
			}

			if cap(b)-len(b) < 64 {
				b = growDouble(b, 64)
			}

			if i > 0 {
				if sep != nil {
					b = append(b, sep...)
				} else {
					b = append(b, ',')
					if hasNewline {
						b = append(b, config.Newline...)
					} else if config.SpaceComma {
						b = append(b, ' ')
					}
					b = appendIndent(b, nextIndentLevel, config.IndentStep)
				}
			} else if hasNewline {
				b = appendIndent(b, nextIndentLevel, config.IndentStep)
			}

			switch v.Type() {
			case ValueTypeString:
				b = appendJsonString(b, v.String(ctx))
			case ValueTypeBool:
				if config.Python {
					if v.Bool() {
						b = append(b, "True"...)
					} else {
						b = append(b, "False"...)
					}
				} else if v.Bool() {
					b = append(b, "true"...)
				} else {
					b = append(b, "false"...)
				}
			case ValueTypeNumber:
				if config.StrictFloat {
					b = strconv.AppendFloat(b, v.Number(), 'f', -1, 64)
				} else {
					b = unparseNumber(b, v.Number())
				}
			case ValueTypeNull:
				if config.Python {
					b = append(b, "None"...)
				} else {
					b = append(b, "null"...)
				}
			default:
				b, err = manifestJson(v, ctx, b, nextIndentLevel, config)
				if err != nil {
					return b, err
				}
			}
		}

		b = append(b, config.Newline...)
		b = appendIndent(b, indentLevel, config.IndentStep)
		b = append(b, ']')
		return b, nil
	case ValueTypeObject:
		obj := value.Object()
		evalCtx := ctx
		evalCtx.Self = value
		err := runAssertions(obj, evalCtx)
		if err != nil {
			return b, err
		}

		// Fast path for single-layer static objects: uses cached pre-sorted
		// field order on LayerTemplate with zero FieldPlan allocations or sorting.
		layers := obj.GetLayers(ctx)
		if len(layers) == 1 && layerAt(layers, 0).Template != nil {
			tmpl := layerAt(layers, 0).Template
			if tmpl.VisibleCount == 0 {
				if config.SpaceComma && !config.Python {
					return append(b, "{ }"...), nil
				}
				if hasNewline {
					b = append(b, '{')
					b = append(b, config.Newline...)
					b = append(b, config.Newline...)
					b = appendIndent(b, indentLevel, config.IndentStep)
					b = append(b, '}')
					return b, nil
				}
				return append(b, "{}"...), nil
			}

			b = append(b, '{')
			nextIndentLevel := indentLevel
			if config.IndentStep != "" {
				nextIndentLevel++
			}
			b = append(b, config.Newline...)

			var sepBuf [128]byte
			var sep []byte
			if nextIndentLevel*len(config.IndentStep)+len(config.Newline)+2 <= len(sepBuf) {
				sep = sepBuf[:0]
				sep = append(sep, ',')
				if hasNewline {
					sep = append(sep, config.Newline...)
				} else if config.SpaceComma {
					sep = append(sep, ' ')
				}
				sep = appendIndent(sep, nextIndentLevel, config.IndentStep)
			}

			subCtx := ctx
			subCtx.Self = value
			inter := subCtx.State.Interner
			order := tmpl.GetSortedVisible(inter)
			keys := tmpl.Keys

			for i, fIdx := range order {
				if i > 0 {
					if sep != nil {
						b = append(b, sep...)
					} else {
						b = append(b, ',')
						if hasNewline {
							b = append(b, config.Newline...)
						} else if config.SpaceComma {
							b = append(b, ' ')
						}
						b = appendIndent(b, nextIndentLevel, config.IndentStep)
					}
				} else {
					b = appendIndent(b, nextIndentLevel, config.IndentStep)
				}

				keyId := keys[fIdx]
				b = appendJsonString(b, inter.Get(keyId))
				b = append(b, config.KeyValSep...)

				fieldValue, err := getSingleLayerFieldValue(obj, tmpl, keyId, int(fIdx), subCtx)
				if err != nil {
					return b, err
				}

				if fieldValue.IsThunk() {
					fieldValue, err = fieldValue.Eval(subCtx)
					if err != nil {
						return b, err
					}
				}
				switch fieldValue.Type() {
				case ValueTypeString:
					b = appendJsonString(b, fieldValue.String(subCtx))
				case ValueTypeBool:
					if config.Python {
						if fieldValue.Bool() {
							b = append(b, "True"...)
						} else {
							b = append(b, "False"...)
						}
					} else if fieldValue.Bool() {
						b = append(b, "true"...)
					} else {
						b = append(b, "false"...)
					}
				case ValueTypeNumber:
					if config.StrictFloat {
						b = strconv.AppendFloat(b, fieldValue.Number(), 'f', -1, 64)
					} else {
						b = unparseNumber(b, fieldValue.Number())
					}
				case ValueTypeNull:
					if config.Python {
						b = append(b, "None"...)
					} else {
						b = append(b, "null"...)
					}
				default:
					b, err = manifestJson(fieldValue, subCtx, b, nextIndentLevel, config)
					if err != nil {
						return b, err
					}
				}
			}

			b = append(b, config.Newline...)
			b = appendIndent(b, indentLevel, config.IndentStep)
			b = append(b, '}')
			return b, nil
		}

		plans := CompileObjectPlan(obj, ctx)

		empty := true
		for i := range plans {
			if !plans[i].IsHidden() {
				empty = false
				break
			}
		}
		if empty {
			if config.SpaceComma && !config.Python {
				return append(b, "{ }"...), nil
			}

			if hasNewline {
				b = append(b, '{')
				b = append(b, config.Newline...)
				b = append(b, config.Newline...)
				b = appendIndent(b, indentLevel, config.IndentStep)
				b = append(b, '}')
				return b, nil
			}

			return append(b, "{}"...), nil
		}

		b = append(b, '{')
		nextIndentLevel := indentLevel
		if config.IndentStep != "" {
			nextIndentLevel++
		}
		b = append(b, config.Newline...)

		// Precompute separator bytes for subsequent fields: "," + newline/space + indent
		var sepBuf [128]byte
		var sep []byte
		if nextIndentLevel*len(config.IndentStep)+len(config.Newline)+2 <= len(sepBuf) {
			sep = sepBuf[:0]
			sep = append(sep, ',')
			if hasNewline {
				sep = append(sep, config.Newline...)
			} else if config.SpaceComma {
				sep = append(sep, ' ')
			}
			sep = appendIndent(sep, nextIndentLevel, config.IndentStep)
		}

		subCtx := ctx
		subCtx.Self = value

		hasWritten := false

		for i := range plans {
			p := &plans[i]
			if p.IsHidden() {
				continue
			}

			if hasWritten {
				if sep != nil {
					b = append(b, sep...)
				} else {
					b = append(b, ',')
					if hasNewline {
						b = append(b, config.Newline...)
					} else if config.SpaceComma {
						b = append(b, ' ')
					}
					b = appendIndent(b, nextIndentLevel, config.IndentStep)
				}
			} else {
				b = appendIndent(b, nextIndentLevel, config.IndentStep)
			}

			b = appendJsonString(b, subCtx.State.Interner.Get(p.KeyId))
			b = append(b, config.KeyValSep...)

			fieldValue, err := p.GetValue(obj, subCtx)
			if err != nil {
				return b, err
			}

			b, err = manifestJson(fieldValue, subCtx, b, nextIndentLevel, config)
			if err != nil {
				return b, err
			}

			hasWritten = true
		}

		b = append(b, config.Newline...)
		b = appendIndent(b, indentLevel, config.IndentStep)
		b = append(b, '}')

		return b, nil
	}
}

func appendIndent(b []byte, indentLevel int, step string) []byte {
	if step == "" || indentLevel == 0 {
		return b
	}

	// 64 spaces
	const maxIndentSpaces = "                                                                "

	if step == " " || step == "  " || step == "   " || step == "    " {
		totalSpaces := indentLevel * len(step)

		for totalSpaces > 0 {
			if totalSpaces <= len(maxIndentSpaces) {
				b = append(b, maxIndentSpaces[:totalSpaces]...)
				break
			}

			b = append(b, maxIndentSpaces...)
			totalSpaces -= len(maxIndentSpaces)
		}
		return b
	}

	for range indentLevel {
		b = append(b, step...)
	}
	return b
}

const maxSafeInt = float64((1 << 53) - 1)

// Borrowed from go-jsonnet, optimized for speed and to avoid fmt reflection
func unparseNumber(dst []byte, v float64) []byte {
	if v == math.Floor(v) {
		if v == 0 {
			if math.Signbit(v) {
				return append(dst, "-0"...)
			}
			return append(dst, '0')
		}
		if math.Abs(v) <= maxSafeInt {
			return strconv.AppendInt(dst, int64(v), 10)
		}
		return strconv.AppendFloat(dst, v, 'f', 0, 64)
	}

	// See "What Every Computer Scientist Should Know About Floating-Point Arithmetic"
	// Theorem 15
	// http://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html
	return strconv.AppendFloat(dst, v, 'g', 17, 64)
}

const (
	hexChars = "0123456789abcdef"

	swarLo01 uint64 = 0x0101010101010101
	swarHi80 uint64 = 0x8080808080808080
	swarSp20 uint64 = 0x2020202020202020
	swarQt22 uint64 = 0x2222222222222222
	swarBs5c uint64 = 0x5c5c5c5c5c5c5c5c
)

// scanJsonSafe advances i past any run of bytes that need no JSON escaping (c >= 0x20 && c != '"' && c != '\\').
func scanJsonSafe(s string, i int) int {
	n := len(s)
	if n-i >= 8 {
		basePtr := unsafe.Pointer(unsafe.StringData(s))
		for i+8 <= n {
			w := *(*uint64)(unsafe.Add(basePtr, i))
			hasCtrl := (w - swarSp20) &^ w & swarHi80
			xq := w ^ swarQt22
			hasQt := (xq - swarLo01) &^ xq & swarHi80
			xs := w ^ swarBs5c
			hasBs := (xs - swarLo01) &^ xs & swarHi80
			if (hasCtrl | hasQt | hasBs) != 0 {
				break
			}
			i += 8
		}
	}
	for i < n {
		c := s[i]
		if c < 0x20 || c == '"' || c == '\\' {
			return i
		}
		i++
	}
	return n
}

func appendJsonString(b []byte, s string) []byte {
	n := len(s)
	if n == 0 {
		return append(b, `""`...)
	}

	i := scanJsonSafe(s, 0)
	if i == n {
		// Fast path: no escapes needed.
		if cap(b)-len(b) < n+2 {
			b = growDouble(b, n+2)
		}
		b = append(b, '"')
		b = append(b, s...)
		b = append(b, '"')
		return b
	}

	return appendJsonStringSlow(b, s, i)
}

//go:noinline
func appendJsonStringSlow(b []byte, s string, i int) []byte {
	n := len(s)
	if cap(b)-len(b) < n+16 {
		b = growDouble(b, n+16)
	}
	b = append(b, '"')
	b = append(b, s[:i]...)

	for i < n {
		c := s[i]
		switch c {
		case '"', '\\':
			b = append(b, '\\', c)
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		case '\b':
			b = append(b, '\\', 'b')
		case '\f':
			b = append(b, '\\', 'f')
		default:
			b = append(b, '\\', 'u', '0', '0', hexChars[c>>4], hexChars[c&0xF])
		}
		start := i + 1
		i = scanJsonSafe(s, start)
		if start < i {
			b = append(b, s[start:i]...)
		}
	}

	b = append(b, '"')
	return b
}
