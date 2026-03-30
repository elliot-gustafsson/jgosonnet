package evaluator

import (
	"cmp"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
