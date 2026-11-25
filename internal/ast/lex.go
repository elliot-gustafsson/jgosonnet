package ast

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
	"unsafe"
)

type TokenKind uint8

const (
	TokenEof TokenKind = iota

	TokenNumber
	TokenString
	TokenStringBlock
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

func (k TokenKind) String() string {
	switch k {
	case TokenBraceL:
		return `"{"`
	case TokenBraceR:
		return `"}"`
	case TokenBracketL:
		return `"["`
	case TokenBracketR:
		return `"]"`
	case TokenComma:
		return `","`
	case TokenDollar:
		return `"$"`
	case TokenDot:
		return `"."`
	case TokenParenL:
		return `"("`
	case TokenParenR:
		return `")"`
	case TokenSemicolon:
		return `";"`
	case TokenIdent:
		return "IDENTIFIER"
	case TokenNumber:
		return "NUMBER"
	case TokenOperator:
		return "OPERATOR"
	case TokenString:
		return "STRING"
	case TokenStringBlock:
		return "STRING_BLOCK"
	case TokenAssert:
		return "assert"
	case TokenElse:
		return "else"
	case TokenError:
		return "error"
	case TokenFalse:
		return "false"
	case TokenFor:
		return "for"
	case TokenFunction:
		return "function"
	case TokenIf:
		return "if"
	case TokenImport:
		return "import"
	case TokenImportStr:
		return "importstr"
	case TokenImportBin:
		return "importbin"
	case TokenIn:
		return "in"
	case TokenLocal:
		return "local"
	case TokenNull:
		return "null"
	case TokenSelf:
		return "self"
	case TokenSuper:
		return "super"
	case TokenTailStrict:
		return "tailstrict"
	case TokenThen:
		return "then"
	case TokenTrue:
		return "true"
	case TokenEof:
		return "end of file"
	default:
		return "unknown"
	}
}

func (t Token) Display() string {
	switch t.Kind {
	case TokenEof:
		return "end of file"
	case TokenNumber, TokenIdent, TokenOperator, TokenString, TokenStringBlock:
		return fmt.Sprintf("(%s, %q)", t.Kind.String(), t.Data)
	default:
		return fmt.Sprintf("%q", t.Data)
	}
}

func keywords(s string) (TokenKind, bool) {
	switch s {
	case "assert":
		return TokenAssert, true
	case "else":
		return TokenElse, true
	case "error":
		return TokenError, true
	case "false":
		return TokenFalse, true
	case "for":
		return TokenFor, true
	case "function":
		return TokenFunction, true
	case "if":
		return TokenIf, true
	case "import":
		return TokenImport, true
	case "importstr":
		return TokenImportStr, true
	case "importbin":
		return TokenImportBin, true
	case "in":
		return TokenIn, true
	case "local":
		return TokenLocal, true
	case "null":
		return TokenNull, true
	case "self":
		return TokenSelf, true
	case "super":
		return TokenSuper, true
	case "tailstrict":
		return TokenTailStrict, true
	case "then":
		return TokenThen, true
	case "true":
		return TokenTrue, true
	default:
		return 0, false
	}
}

type Token struct {
	Kind TokenKind
	Data string
	Pos  uint32
	End  uint32
}

type lexer struct {
	filename string
	data     string
	pos      uint32
}

func (l *lexer) staticErr(start, end uint32, msg string) error {
	return StaticError{
		Loc: byteRangeToLoc(l.filename, l.data, start, end),
		Msg: msg,
	}
}

func (l *lexer) staticErrPoint(pos uint32, msg string) error {
	return l.staticErr(pos, pos, msg)
}

func (l *lexer) Next() (Token, error) {
	err := l.consumeWhitespaceAndComments()
	if err != nil {
		return Token{}, err
	}

	if l.pos >= uint32(len(l.data)) {
		return Token{Kind: TokenEof, Pos: l.pos, End: l.pos}, nil
	}

	start := l.pos
	c := l.data[l.pos]
	l.pos++

	switch c {
	case '{':
		return Token{Kind: TokenBraceL, Data: "{", Pos: start, End: l.pos}, nil
	case '}':
		return Token{Kind: TokenBraceR, Data: "}", Pos: start, End: l.pos}, nil
	case '[':
		return Token{Kind: TokenBracketL, Data: "[", Pos: start, End: l.pos}, nil
	case ']':
		return Token{Kind: TokenBracketR, Data: "]", Pos: start, End: l.pos}, nil
	case ',':
		return Token{Kind: TokenComma, Data: ",", Pos: start, End: l.pos}, nil
	case '.':
		return Token{Kind: TokenDot, Data: ".", Pos: start, End: l.pos}, nil
	case '(':
		return Token{Kind: TokenParenL, Data: "(", Pos: start, End: l.pos}, nil
	case ')':
		return Token{Kind: TokenParenR, Data: ")", Pos: start, End: l.pos}, nil
	case ';':
		return Token{Kind: TokenSemicolon, Data: ";", Pos: start, End: l.pos}, nil
	case '"', '\'':
		return l.lexString(c, start)
	case '@':
		if l.pos < uint32(len(l.data)) && (l.data[l.pos] == '"' || l.data[l.pos] == '\'') {
			quote := l.data[l.pos]
			l.pos++
			return l.lexVerbatimString(quote, start)
		}
		var r rune = -1
		if l.pos < uint32(len(l.data)) {
			r, _ = utf8.DecodeRuneInString(l.data[l.pos:])
		}
		return Token{}, l.staticErrPoint(start, fmt.Sprintf("Couldn't lex verbatim string, junk after '@': %v", r))
	case '|':
		if l.pos+1 < uint32(len(l.data)) && l.data[l.pos] == '|' && l.data[l.pos+1] == '|' {
			l.pos += 2
			return l.lexBlockString(start)
		}
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

	r, _ := utf8.DecodeRuneInString(l.data[start:])
	return Token{}, l.staticErrPoint(start, fmt.Sprintf("Could not lex the character %s", strconv.QuoteRuneToASCII(r)))
}

func (l *lexer) consumeWhitespaceAndComments() error {
	n := uint32(len(l.data))
	for l.pos < n {
		c := l.data[l.pos]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			l.pos++
			continue
		}

		if c == '/' {
			if l.pos+1 < n && l.data[l.pos+1] == '/' {
				l.pos += 2
				idx := strings.IndexByte(l.data[l.pos:], '\n')
				if idx == -1 {
					l.pos = n
				} else {
					l.pos += uint32(idx) + 1
				}
				continue
			}
			if l.pos+1 < n && l.data[l.pos+1] == '*' {
				commentStart := l.pos
				l.pos += 2
				idx := strings.Index(l.data[l.pos:], "*/")
				if idx == -1 {
					return l.staticErrPoint(commentStart, "Multi-line comment has no terminating */")
				}
				l.pos += uint32(idx) + 2
				continue
			}
		}

		if c == '#' {
			l.pos++
			idx := strings.IndexByte(l.data[l.pos:], '\n')
			if idx == -1 {
				l.pos = n
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
	type numState uint8
	const (
		numBegin numState = iota
		numAfterZero
		numAfterOneToNine
		numAfterIntUnderscore
		numAfterDot
		numAfterDigit
		numAfterFracUnderscore
		numAfterE
		numAfterExpSign
		numAfterExpDigit
		numAfterExpUnderscore
	)

	n := uint32(len(l.data))
	state := numBegin
	hasUnderscore := false

	peekRune := func() rune {
		if l.pos >= n {
			return -1
		}
		r, _ := utf8.DecodeRuneInString(l.data[l.pos:])
		return r
	}

outer:
	for {
		var c byte
		if l.pos < n {
			c = l.data[l.pos]
		}

		switch state {
		case numBegin:
			if c == '0' {
				state = numAfterZero
			} else {
				state = numAfterOneToNine
			}
		case numAfterZero:
			switch c {
			case '.':
				state = numAfterDot
			case 'e', 'E':
				state = numAfterE
			case '_':
				return Token{}, l.staticErrPoint(l.pos, "Couldn't lex number, _ not allowed after leading 0")
			default:
				break outer
			}
		case numAfterOneToNine:
			switch {
			case c == '.':
				state = numAfterDot
			case c == 'e' || c == 'E':
				state = numAfterE
			case c >= '0' && c <= '9':
				state = numAfterOneToNine
			case c == '_':
				state = numAfterIntUnderscore
				hasUnderscore = true
			default:
				break outer
			}
		case numAfterIntUnderscore:
			if c >= '0' && c <= '9' {
				state = numAfterOneToNine
			} else {
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after '_': %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		case numAfterDot:
			if c >= '0' && c <= '9' {
				state = numAfterDigit
			} else {
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after decimal point: %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		case numAfterDigit:
			switch {
			case c == 'e' || c == 'E':
				state = numAfterE
			case c >= '0' && c <= '9':
				state = numAfterDigit
			case c == '_':
				state = numAfterFracUnderscore
				hasUnderscore = true
			default:
				break outer
			}
		case numAfterFracUnderscore:
			if c >= '0' && c <= '9' {
				state = numAfterDigit
			} else {
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after '_': %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		case numAfterE:
			switch {
			case c == '+' || c == '-':
				state = numAfterExpSign
			case c >= '0' && c <= '9':
				state = numAfterExpDigit
			default:
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after 'E': %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		case numAfterExpSign:
			if c >= '0' && c <= '9' {
				state = numAfterExpDigit
			} else {
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after exponent sign: %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		case numAfterExpDigit:
			switch {
			case c >= '0' && c <= '9':
				state = numAfterExpDigit
			case c == '_':
				state = numAfterExpUnderscore
				hasUnderscore = true
			default:
				break outer
			}
		case numAfterExpUnderscore:
			if c >= '0' && c <= '9' {
				state = numAfterExpDigit
			} else {
				return Token{}, l.staticErrPoint(l.pos, fmt.Sprintf("Couldn't lex number, junk after '_': %v", strconv.QuoteRuneToASCII(peekRune())))
			}
		}

		l.pos++
	}

	raw := l.data[start:l.pos]
	if hasUnderscore {
		buf := make([]byte, 0, len(raw))
		for i := 0; i < len(raw); i++ {
			if raw[i] != '_' {
				buf = append(buf, raw[i])
			}
		}
		raw = unsafe.String(unsafe.SliceData(buf), len(buf))
	}

	return Token{Kind: TokenNumber, Data: raw, Pos: start, End: l.pos}, nil
}

func (l *lexer) lexIdent(start uint32) (Token, error) {
	n := uint32(len(l.data))
	for l.pos < n && isIdent(l.data[l.pos]) {
		l.pos++
	}
	data := l.data[start:l.pos]
	if kind, ok := keywords(data); ok {
		return Token{Kind: kind, Data: data, Pos: start, End: l.pos}, nil
	}
	return Token{Kind: TokenIdent, Data: data, Pos: start, End: l.pos}, nil
}

func (l *lexer) lexOperator(start uint32) (Token, error) {
	n := uint32(len(l.data))
	l.pos++ // consume first symbol
	for l.pos < n && isSymbol(l.data[l.pos]) {
		c := l.data[l.pos]
		// Stop before // or /* comments
		if c == '/' && l.pos+1 < n && (l.data[l.pos+1] == '/' || l.data[l.pos+1] == '*') {
			break
		}
		// Stop before ||| text blocks
		if c == '|' && l.pos+2 < n && l.data[l.pos+1] == '|' && l.data[l.pos+2] == '|' {
			break
		}
		l.pos++
	}

	// Operators may not end with +, -, ~, !, $ unless length is 1.
	for l.pos > start+1 {
		switch l.data[l.pos-1] {
		case '+', '-', '~', '!', '$':
			l.pos--
			continue
		}
		break
	}

	data := l.data[start:l.pos]
	if data == "$" {
		return Token{Kind: TokenDollar, Data: data, Pos: start, End: l.pos}, nil
	}
	return Token{Kind: TokenOperator, Data: data, Pos: start, End: l.pos}, nil
}

func (l *lexer) lexString(quoteType byte, start uint32) (Token, error) {
	hasEscape := false
	for {
		idx := strings.IndexByte(l.data[l.pos:], quoteType)
		if idx == -1 {
			return Token{}, l.staticErrPoint(start, "Unterminated String")
		}

		absoluteIdx := l.pos + uint32(idx)

		backslashes := 0
		for j := int(absoluteIdx) - 1; j > int(start) && l.data[j] == '\\'; j-- {
			backslashes++
		}
		if backslashes > 0 {
			hasEscape = true
		}

		l.pos = absoluteIdx + 1

		if backslashes%2 == 0 {
			break
		}
	}

	raw := l.data[start+1 : l.pos-1]
	if !hasEscape && strings.IndexByte(raw, '\\') != -1 {
		hasEscape = true
	}
	if hasEscape {
		unescaped, err := l.unescapeString(raw, start, l.pos)
		if err != nil {
			return Token{}, err
		}
		return Token{Kind: TokenString, Data: unescaped, Pos: start, End: l.pos}, nil
	}
	return Token{Kind: TokenString, Data: raw, Pos: start, End: l.pos}, nil
}

func (l *lexer) unescapeString(s string, start, end uint32) (string, error) {
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); {
		c := s[i]
		if c != '\\' {
			buf = append(buf, c)
			i++
			continue
		}
		i++
		if i >= len(s) {
			return "", l.staticErr(start, end, "Truncated escape sequence in string literal.")
		}
		r2, w := utf8.DecodeRuneInString(s[i:])
		i += w
		switch r2 {
		case '"':
			buf = append(buf, '"')
		case '\'':
			buf = append(buf, '\'')
		case '\\':
			buf = append(buf, '\\')
		case '/':
			buf = append(buf, '/')
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			if i+4 > len(s) {
				return "", l.staticErr(start, end, "Truncated unicode escape sequence in string literal.")
			}
			code, ok := parseHex4(s[i : i+4])
			if !ok {
				return "", l.staticErr(start, end, fmt.Sprintf("Unicode escape sequence was malformed: %s", s[i:i+4]))
			}
			i += 4
			if utf16.IsSurrogate(code) {
				highSurrogate := code
				if i+6 > len(s) {
					return "", l.staticErr(start, end, "Truncated unicode surrogate pair escape sequence in string literal.")
				}
				if s[i] != '\\' || s[i+1] != 'u' {
					return "", l.staticErr(start, end, "Unicode surrogate pair escape sequence missing low surrogate in string literal.")
				}
				i += 2
				lowSurrogate, ok := parseHex4(s[i : i+4])
				if !ok {
					return "", l.staticErr(start, end, fmt.Sprintf("Unicode low surrogate escape sequence was malformed: %s", s[i:i+4]))
				}
				i += 4
				code = utf16.DecodeRune(highSurrogate, lowSurrogate)
			}
			buf = utf8.AppendRune(buf, code)
		default:
			return "", l.staticErr(start, end, fmt.Sprintf("Unknown escape sequence in string literal: \\%c", r2))
		}
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf)), nil
}

func parseHex4(s string) (rune, bool) {
	var v rune
	for i := 0; i < 4; i++ {
		c := s[i]
		v <<= 4
		switch {
		case c >= '0' && c <= '9':
			v |= rune(c - '0')
		case c >= 'a' && c <= 'f':
			v |= rune(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			v |= rune(c - 'A' + 10)
		default:
			return 0, false
		}
	}
	return v, true
}

func (l *lexer) lexVerbatimString(quoteType byte, start uint32) (Token, error) {
	strStart := l.pos
	hasDoubled := false
	for {
		idx := strings.IndexByte(l.data[l.pos:], quoteType)
		if idx == -1 {
			return Token{}, l.staticErrPoint(start, "Unterminated String")
		}

		absoluteIdx := l.pos + uint32(idx)

		if int(absoluteIdx+1) < len(l.data) && l.data[absoluteIdx+1] == quoteType {
			hasDoubled = true
			l.pos = absoluteIdx + 2
			continue
		}

		l.pos = absoluteIdx + 1
		break
	}

	strVal := l.data[strStart : l.pos-1]
	if hasDoubled {
		if quoteType == '"' {
			strVal = strings.ReplaceAll(strVal, `""`, `"`)
		} else {
			strVal = strings.ReplaceAll(strVal, `''`, `'`)
		}
	}

	return Token{Kind: TokenString, Data: strVal, Pos: start, End: l.pos}, nil
}

func checkBlockWhitespace(indent, s string) int {
	i := 0
	for ; i < len(indent); i++ {
		if indent[i] != ' ' && indent[i] != '\t' {
			return i
		}
		if i >= len(s) || indent[i] != s[i] {
			return 0
		}
	}
	return i
}

func (l *lexer) lexBlockString(start uint32) (Token, error) {
	n := uint32(len(l.data))

	chompTrailingNL := false
	if l.pos < n && l.data[l.pos] == '-' {
		chompTrailingNL = true
		l.pos++
	}

	// Skip horizontal whitespace up to required '\n'
	for l.pos < n && (l.data[l.pos] == ' ' || l.data[l.pos] == '\t' || l.data[l.pos] == '\r') {
		l.pos++
	}
	if l.pos >= n || l.data[l.pos] != '\n' {
		return Token{}, l.staticErrPoint(start, "Text block requires new line after |||.")
	}
	l.pos++ // consume '\n'

	var out []byte

	// Preserve leading blank lines before calculating indent
	for l.pos < n && l.data[l.pos] == '\n' {
		out = append(out, '\n')
		l.pos++
	}

	numWS := checkBlockWhitespace(l.data[l.pos:], l.data[l.pos:])
	if numWS == 0 {
		return Token{}, l.staticErrPoint(start, "Text block's first line must start with whitespace")
	}
	blockIndent := l.data[l.pos : l.pos+uint32(numWS)]

	for {
		l.pos += uint32(numWS)
		nlIdx := strings.IndexByte(l.data[l.pos:], '\n')
		if nlIdx == -1 {
			return Token{}, l.staticErrPoint(start, "Unexpected EOF")
		}
		out = append(out, l.data[l.pos:l.pos+uint32(nlIdx)+1]...)
		l.pos += uint32(nlIdx) + 1

		// Preserve blank lines
		for l.pos < n && l.data[l.pos] == '\n' {
			out = append(out, '\n')
			l.pos++
		}

		numWS = checkBlockWhitespace(blockIndent, l.data[l.pos:])
		if numWS == 0 {
			for l.pos < n && (l.data[l.pos] == ' ' || l.data[l.pos] == '\t') {
				l.pos++
			}
			if !strings.HasPrefix(l.data[l.pos:], "|||") {
				return Token{}, l.staticErrPoint(start, "Text block not terminated with |||")
			}
			l.pos += 3 // consume '|||'

			if chompTrailingNL && len(out) > 0 {
				out = out[:len(out)-1]
			}
			finalStr := unsafe.String(unsafe.SliceData(out), len(out))
			return Token{Kind: TokenStringBlock, Data: finalStr, Pos: start, End: l.pos}, nil
		}
	}
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
	case '!', '$', ':', '~', '+', '-', '&', '|', '^', '=', '<', '>', '*', '/', '%':
		return true
	}
	return false
}
