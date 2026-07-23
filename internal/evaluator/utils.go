package evaluator

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

var builderPool = sync.Pool{
	New: func() any {
		return new(Builder)
	},
}

func GetBuilder() *Builder {
	return builderPool.Get().(*Builder)
}

func PutBuilder(b *Builder) {
	b.Reset()
	builderPool.Put(b)
}

type Builder struct {
	buf []byte
}

func (b *Builder) Grow(n int) {
	b.buf = slices.Grow(b.buf, n)
}

func (b *Builder) Length() int {
	return len(b.buf)
}

func (b *Builder) Write(p []byte) {
	b.buf = append(b.buf, p...)
}

func (b *Builder) Set(p []byte) {
	b.buf = p
}

func (b *Builder) WriteString(s string) {
	b.buf = append(b.buf, s...)
}

func (b *Builder) WriteByte(c byte) {
	b.buf = append(b.buf, c)
}

func (b *Builder) WriteRune(c rune) {
	b.buf = utf8.AppendRune(b.buf, c)
}

func (b *Builder) AppendInt(i int64, base int) {
	b.buf = strconv.AppendInt(b.buf, i, base)
}

func (b *Builder) AppendUint(i uint64, base int) {
	b.buf = strconv.AppendUint(b.buf, i, base)
}

func (b *Builder) AppendFloat(f float64, fmt byte, prec, bitSize int) {
	b.buf = strconv.AppendFloat(b.buf, f, fmt, prec, bitSize)
}

func (b *Builder) Bytes() []byte {
	return b.buf
}

func (b *Builder) String() string {
	return string(b.buf)
}

func (b *Builder) UnsafeString() string {
	if len(b.buf) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))
}

func (b *Builder) Reset() {
	b.buf = b.buf[:0]
}

func naturalStringSort(a, b string) int {
	i, j := 0, 0

	for i < len(a) && j < len(b) {
		rA, widthA := utf8.DecodeRuneInString(a[i:])
		rB, widthB := utf8.DecodeRuneInString(b[j:])

		if unicode.IsDigit(rA) && unicode.IsDigit(rB) {
			startI, startJ := i, j

			for i < len(a) {
				r, w := utf8.DecodeRuneInString(a[i:])
				if !unicode.IsDigit(r) {
					break
				}
				i += w
			}

			for j < len(b) {
				r, w := utf8.DecodeRuneInString(b[j:])
				if !unicode.IsDigit(r) {
					break
				}
				j += w
			}

			aNumStr := a[startI:i]
			bNumStr := b[startJ:j]

			aTrim := strings.TrimLeft(aNumStr, "0")
			bTrim := strings.TrimLeft(bNumStr, "0")

			if len(aTrim) != len(bTrim) {
				return cmp.Compare(len(aTrim), len(bTrim))
			}
			if res := cmp.Compare(aTrim, bTrim); res != 0 {
				return res
			}

			if res := cmp.Compare(len(aNumStr), len(bNumStr)); res != 0 {
				return res
			}
			continue
		}

		alphaA := isAlphanumeric(rA)
		alphaB := isAlphanumeric(rB)

		if alphaA != alphaB {
			if alphaA {
				return 1 // A is alphanumeric, B is a sign. B comes first.
			}
			return -1 // A is a sign, B is alphanumeric. A comes first.
		}

		if rA != rB {
			return cmp.Compare(rA, rB)
		}

		i += widthA
		j += widthB
	}

	return cmp.Compare(len(a), len(b))
}

func isAlphanumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
