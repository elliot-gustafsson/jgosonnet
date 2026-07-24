package evaluator

import (
	"cmp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func naturalStringSort(a, b string) int {
	i, j := 0, 0
	lenA, lenB := len(a), len(b)

	for i < lenA && j < lenB {
		cA, cB := a[i], b[j]

		// Fallback to slow UTF-8 path for non-ASCII characters
		if cA >= utf8.RuneSelf || cB >= utf8.RuneSelf {
			return naturalStringSortUTF8(a, b, i, j)
		}

		if isDigit(cA) && isDigit(cB) {
			startI, startJ := i, j

			for i < lenA && isDigit(a[i]) {
				i++
			}
			for j < lenB && isDigit(b[j]) {
				j++
			}

			num1 := a[startI:i]
			num2 := b[startJ:j]

			trim1, trim2 := 0, 0
			for trim1 < len(num1) && num1[trim1] == '0' {
				trim1++
			}
			for trim2 < len(num2) && num2[trim2] == '0' {
				trim2++
			}

			lenTrim1 := len(num1) - trim1
			lenTrim2 := len(num2) - trim2

			if lenTrim1 != lenTrim2 {
				return cmp.Compare(lenTrim1, lenTrim2)
			}

			for k := 0; k < lenTrim1; k++ {
				if num1[trim1+k] != num2[trim2+k] {
					return cmp.Compare(num1[trim1+k], num2[trim2+k])
				}
			}

			if len(num1) != len(num2) {
				return cmp.Compare(len(num1), len(num2))
			}
			continue
		}

		alphaA := isAlphanumeric(cA)
		alphaB := isAlphanumeric(cB)

		if alphaA != alphaB {
			if alphaA {
				return 1
			}
			return -1
		}

		if cA != cB {
			return cmp.Compare(cA, cB)
		}

		i++
		j++
	}

	return cmp.Compare(lenA, lenB)
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isAlphanumeric(b byte) bool {
	return isDigit(b) || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}

//go:noinline
func naturalStringSortUTF8(a, b string, i, j int) int {

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

		alphaA := unicode.IsLetter(rA) || unicode.IsDigit(rA)
		alphaB := unicode.IsLetter(rB) || unicode.IsDigit(rB)

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
