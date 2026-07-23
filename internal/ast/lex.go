package ast

import (
	"fmt"
	"strings"
	"unsafe"
)

type TokenKind uint8

const (
	TokenEof TokenKind = iota

	TokenNumber
	TokenString
	TokenIdent

	TokenBraceL
	TokenBraceR
	TokenBracketL
	TokenBracketR
	TokenComma
	TokenDollar
	TokenDot
	TokenParenL
	TokenParenR
	TokenSemicolon
	TokenColon

	TokenOperator

	// Keywords
	TokenAssert
	TokenElse
	TokenError
	TokenFalse
	TokenFor
	TokenFunction
	TokenIf
	TokenImport
	TokenImportStr
	TokenImportBin
	TokenIn
	TokenLocal
	TokenNull
	TokenSelf
	TokenSuper
	TokenTailStrict
	TokenThen
	TokenTrue
)

var keywords = map[string]TokenKind{
	"assert":     TokenAssert,
	"else":       TokenElse,
	"error":      TokenError,
	"false":      TokenFalse,
	"for":        TokenFor,
	"function":   TokenFunction,
	"if":         TokenIf,
	"import":     TokenImport,
	"importstr":  TokenImportStr,
	"importbin":  TokenImportBin,
	"in":         TokenIn,
	"local":      TokenLocal,
	"null":       TokenNull,
	"self":       TokenSelf,
	"super":      TokenSuper,
	"tailstrict": TokenTailStrict,
	"then":       TokenThen,
	"true":       TokenTrue,
}

type Token struct {
	Kind TokenKind
	Data string
	Pos  uint32
}

type lexer struct {
	data string
	pos  uint32
}

func (l *lexer) Next() (Token, error) {
	err := l.consumeWhitespaceAndComments()
	if err != nil {
		return Token{}, err
	}

	if l.pos >= uint32(len(l.data)) {
		return Token{Kind: TokenEof, Pos: l.pos}, nil
	}

	start := l.pos
	c := l.data[l.pos]
	l.pos++

	switch c {
	case '{':
		return Token{Kind: TokenBraceL, Data: "{", Pos: start}, nil
	case '}':
		return Token{Kind: TokenBraceR, Data: "}", Pos: start}, nil
	case '[':
		return Token{Kind: TokenBracketL, Data: "[", Pos: start}, nil
	case ']':
		return Token{Kind: TokenBracketR, Data: "]", Pos: start}, nil
	case ',':
		return Token{Kind: TokenComma, Data: ",", Pos: start}, nil
	case '.':
		return Token{Kind: TokenDot, Data: ".", Pos: start}, nil
	case '(':
		return Token{Kind: TokenParenL, Data: "(", Pos: start}, nil
	case ')':
		return Token{Kind: TokenParenR, Data: ")", Pos: start}, nil
	case ';':
		return Token{Kind: TokenSemicolon, Data: ";", Pos: start}, nil
	case '$':
		return Token{Kind: TokenDollar, Data: "$", Pos: start}, nil
	case '"', '\'':
		return l.lexString(c, start)
	case '@':
		if l.pos < uint32(len(l.data)) && (l.data[l.pos] == '"' || l.data[l.pos] == '\'') {
			quote := l.data[l.pos]
			l.pos++
			return l.lexVerbatimString(quote, start)
		}
		return Token{}, fmt.Errorf("unexpected '@'")
	case '|':
		if l.pos+1 < uint32(len(l.data)) && l.data[l.pos] == '|' && l.data[l.pos+1] == '|' {
			l.pos += 2
			return l.lexBlockString(start)
		}
		// Otherwise, it's an operator starting with |
		l.pos--
		return l.lexOperator(start)
	}

	if isDigit(c) {
		l.pos--
		return l.lexNumber(start)
	}
	if isIdentFirst(c) {
		l.pos--
		return l.lexIdent(start)
	}
	if isSymbol(c) {
		l.pos--
		return l.lexOperator(start)
	}

	return Token{}, fmt.Errorf("unexpected character '%c'", c)
}

func (l *lexer) consumeWhitespaceAndComments() error {
	for l.pos < uint32(len(l.data)) {
		c := l.data[l.pos]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			l.pos++
			continue
		}

		if c == '/' {
			if l.pos+1 < uint32(len(l.data)) && l.data[l.pos+1] == '/' {
				l.pos += 2
				idx := strings.IndexByte(l.data[l.pos:], '\n')
				if idx == -1 {
					l.pos = uint32(len(l.data))
				} else {
					l.pos += uint32(idx) + 1
				}
				continue
			}
			if l.pos+1 < uint32(len(l.data)) && l.data[l.pos+1] == '*' {
				l.pos += 2
				idx := strings.Index(l.data[l.pos:], "*/")
				if idx == -1 {
					return fmt.Errorf("unterminated block comment")
				}
				l.pos += uint32(idx) + 2
				continue
			}
		}

		if c == '#' {
			l.pos++
			idx := strings.IndexByte(l.data[l.pos:], '\n')
			if idx == -1 {
				l.pos = uint32(len(l.data))
			} else {
				l.pos += uint32(idx) + 1
			}
			continue
		}

		break
	}
	return nil
}

func (l *lexer) lexNumber(start uint32) (Token, error) {
	if l.data[l.pos] == '0' && l.pos+1 < uint32(len(l.data)) && l.data[l.pos+1] >= '0' && l.data[l.pos+1] <= '9' {
		return Token{}, fmt.Errorf("leading zero is not allowed in numbers")
	}

	for l.pos < uint32(len(l.data)) {
		c := l.data[l.pos]
		if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' {
			l.pos++
		} else if c == '+' || c == '-' {
			if l.pos > start && (l.data[l.pos-1] == 'e' || l.data[l.pos-1] == 'E') {
				l.pos++
			} else {
				break
			}
		} else {
			break
		}
	}
	return Token{Kind: TokenNumber, Data: l.data[start:l.pos], Pos: start}, nil
}

func (l *lexer) lexIdent(start uint32) (Token, error) {
	for l.pos < uint32(len(l.data)) && isIdent(l.data[l.pos]) {
		l.pos++
	}
	data := l.data[start:l.pos]
	if kind, ok := keywords[data]; ok {
		return Token{Kind: kind, Data: data, Pos: start}, nil
	}
	return Token{Kind: TokenIdent, Data: data, Pos: start}, nil
}

func (l *lexer) lexOperator(start uint32) (Token, error) {
	for l.pos < uint32(len(l.data)) && isSymbol(l.data[l.pos]) {
		l.pos++
	}
	data := l.data[start:l.pos]
	if data == ":" {
		return Token{Kind: TokenColon, Data: data, Pos: start}, nil
	}
	return Token{Kind: TokenOperator, Data: data, Pos: start}, nil
}

func (l *lexer) lexString(quoteType byte, start uint32) (Token, error) {
	for {
		idx := strings.IndexByte(l.data[l.pos:], quoteType)
		if idx == -1 {
			return Token{}, fmt.Errorf("unterminated string literal")
		}

		absoluteIdx := l.pos + uint32(idx)

		backslashes := 0
		for j := absoluteIdx - 1; j >= start && l.data[j] == '\\'; j-- {
			backslashes++
		}

		l.pos = absoluteIdx + 1

		if backslashes%2 == 0 {
			break
		}
	}
	return Token{Kind: TokenString, Data: l.data[start+1 : l.pos-1], Pos: start}, nil
}

func (l *lexer) lexVerbatimString(quoteType byte, start uint32) (Token, error) {
	strStart := l.pos
	for {
		idx := strings.IndexByte(l.data[l.pos:], quoteType)
		if idx == -1 {
			return Token{}, fmt.Errorf("unterminated verbatim string literal")
		}

		absoluteIdx := l.pos + uint32(idx)

		if int(absoluteIdx+1) < len(l.data) && l.data[absoluteIdx+1] == quoteType {
			l.pos = absoluteIdx + 2
			continue
		}

		l.pos = absoluteIdx + 1
		break
	}

	strVal := l.data[strStart : l.pos-1]

	if quoteType == '"' && strings.Contains(strVal, `""`) {
		strVal = strings.ReplaceAll(strVal, `""`, `"`)
	} else if quoteType == '\'' && strings.Contains(strVal, `''`) {
		strVal = strings.ReplaceAll(strVal, `''`, `'`)
	}

	return Token{Kind: TokenString, Data: strVal, Pos: start}, nil
}

func (l *lexer) lexBlockString(start uint32) (Token, error) {
	strStart := l.pos
	for {
		idx := strings.IndexByte(l.data[l.pos:], '|')
		if idx == -1 {
			return Token{}, fmt.Errorf("unterminated block string")
		}

		l.pos += uint32(idx)

		if int(l.pos+2) >= len(l.data) {
			return Token{}, fmt.Errorf("unterminated block string")
		}

		if l.data[l.pos:l.pos+3] == "|||" {
			break
		}

		if l.data[l.pos+1] == '|' {
			l.pos += 2
		} else {
			l.pos++
		}
	}

	rawBlock := l.data[strStart:l.pos]
	l.pos += 3 // consume |||

	lastNL := strings.LastIndexByte(rawBlock, '\n')
	var indentStr string
	if lastNL != -1 {
		indentStr = rawBlock[lastNL+1:]
		rawBlock = rawBlock[:lastNL+1]
	}

	if len(rawBlock) > 0 && rawBlock[0] == '\n' {
		rawBlock = rawBlock[1:]
	}

	if len(indentStr) == 0 {
		return Token{Kind: TokenString, Data: rawBlock, Pos: start}, nil
	}

	indentLen := len(indentStr)
	out := make([]byte, 0, len(rawBlock))
	pos := 0

	for pos < len(rawBlock) {
		nextNL := strings.IndexByte(rawBlock[pos:], '\n')
		var line string

		if nextNL == -1 {
			line = rawBlock[pos:]
			pos = len(rawBlock)
		} else {
			line = rawBlock[pos : pos+nextNL]
			pos += nextNL + 1
		}

		if strings.HasPrefix(line, indentStr) {
			out = append(out, line[indentLen:]...)
		} else {
			out = append(out, line...)
		}

		if nextNL != -1 {
			out = append(out, '\n')
		}
	}

	finalStr := unsafe.String(unsafe.SliceData(out), len(out))
	return Token{Kind: TokenString, Data: finalStr, Pos: start}, nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isIdentFirst(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isIdent(c byte) bool {
	return isIdentFirst(c) || isDigit(c)
}

func isSymbol(c byte) bool {
	switch c {
	case '!', ':', '~', '+', '-', '&', '|', '^', '=', '<', '>', '*', '/', '%':
		return true
	}
	return false
}
