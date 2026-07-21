package ast

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type parser struct {
	data string
	pos  uint32

	col  uint32
	line uint32

	Nodes []Node

	SideTable []uint32

	Locations []NodeContext

	Interner *interner.Interner
}

func Parse(filename string, interner *interner.Interner) (*AST, error) {

	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := unsafe.String(unsafe.SliceData(rawData), len(rawData))

	p := parser{
		data: data,

		Interner:  interner,
		Nodes:     make([]Node, 1, 8192),
		Locations: make([]NodeContext, 1, 8192),
		SideTable: make([]uint32, 0, 4096),
	}

	err = p.parse()
	if err != nil {
		return nil, err
	}

	if len(p.Nodes) == 0 {
		return nil, fmt.Errorf("empyy file")
	}

	return &AST{
		RootId:    uint32(len(p.Nodes) - 1),
		Nodes:     p.Nodes,
		SideTable: p.SideTable,
		Locations: p.Locations,
	}, nil
}

func (t *parser) emit(n Node) uint32 {
	id := uint32(len(t.Nodes))
	t.Nodes = append(t.Nodes, n)
	t.Locations = append(t.Locations, NodeContext{}) // TODO: add real ctx
	return id
}

func (t *parser) parse() error {

	for {
		t.consumeWhitespace()

		r := t.peekByte()

		if r == 0 {
			return nil
		}

		switch r {
		default:
			return fmt.Errorf("unhandled token: %c", r)

		// Comments
		case '/':
			t.pos++ // consume first slash

			next := t.peekByte()

			if next == '/' {
				t.pos++ // consume second slash
				newLineIdx := strings.IndexByte(t.data[t.pos:], '\n')
				if newLineIdx == -1 {
					return nil // all comment
				}
				t.pos += uint32(newLineIdx) + 1 // jump to and consume newline
				continue
			}

			if next == '*' {
				t.pos++ // consume star
				for {
					// Search from current position
					idx := strings.IndexByte(t.data[t.pos:], '*')
					if idx == -1 {
						// TODO: add location data
						return fmt.Errorf("unterminated block comment")
					}
					t.pos += uint32(idx) // jump to the '*'

					if t.peekByteOffset(1) == '/' {
						t.pos += 2 // consume "*/"
						break
					}

					t.pos++
				}
				continue
			}

			return fmt.Errorf("unexpected token '/'")
		case '#':
			t.pos++ // consume hash
			newLineIdx := strings.IndexByte(t.data[t.pos:], '\n')
			if newLineIdx == -1 {
				return nil // all comment
			}
			t.pos += uint32(newLineIdx) + 1 // jump to and consume newline

		// Numbers
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			// t.pos++

			err := t.parseNumber()
			if err != nil {
				return err
			}

		// Strings
		case '"':
			err := t.parseString('"')
			if err != nil {
				return err
			}
		case '\'':
			err := t.parseString('\'')
			if err != nil {
				return err
			}

		// Tokens
		case '|':

			if t.sliceOffset(2) == "||" {
				t.pos += 3 // comsume all three pipes
				// TODO: this is broken... doesnt strip space prefix properly
				err := t.parseBlockString()
				if err != nil {
					return err
				}
			}

			// TODO: Handle other | stuff

		case '@':
			next := t.peekByteOffset(1)
			if next == '"' || next == '\'' {
				t.pos += 2 // Consume `@` and the quote
				err := t.parseVerbatimString(next)
				if err != nil {
					return err
				}
				continue
			}

			// TODO: Handle other uses of `@` if any, or return error

		// Array
		case '[':

		}

	}

}

func (t *parser) peekByte() byte {
	if int(t.pos) >= len(t.data) {
		return 0
	}
	return t.data[t.pos]
}

func (t *parser) peekByteOffset(offset uint32) byte {
	if int(t.pos+offset) >= len(t.data) {
		return 0
	}
	return t.data[t.pos+offset]
}

func (t *parser) sliceOffset(offset uint32) string {
	if int(t.pos+offset) >= len(t.data) {
		return ""
	}
	return t.data[t.pos : t.pos+offset]
}

func (t *parser) peek() rune {

	if int(t.pos) >= len(t.data) {
		return -1
	}

	r, _ := utf8.DecodeRuneInString(t.data[t.pos:])
	return r
}

func (t *parser) consumeWhitespace() {
	for int(t.pos) < len(t.data) {
		c := t.data[t.pos]
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			t.pos++

			if c == '\n' {
				t.line++
				t.col = 0
			} else {
				t.col++
			}
		} else {
			break
		}
	}
}

func (t *parser) parseNumber() error {
	start := t.pos
	dataLen := uint32(len(t.data))

	// Ensure the first character is a digit or leading minus/dot if you support that
	if t.data[start] == '0' && start+1 < dataLen && t.data[start+1] >= '0' && t.data[start+1] <= '9' {
		// reject if leading zero
		return fmt.Errorf("leading zero is not allowed in numbers")
	}

	for t.pos < dataLen {
		c := t.data[t.pos]

		// Fast character class check for valid number components
		if (c >= '0' && c <= '9') || c == '.' || c == 'e' || c == 'E' || c == '_' {
			t.pos++
		} else if c == '+' || c == '-' {
			// Signs are ONLY valid immediately following an 'e' or 'E'
			if t.pos > start && (t.data[t.pos-1] == 'e' || t.data[t.pos-1] == 'E') {
				t.pos++
			} else {
				// It's a binary operator like `1 + 2`, so stop lexing the number
				break
			}
		} else {
			// We hit a space, bracket, or other syntax character
			break
		}
	}

	// Extract the zero-allocation substring
	numStr := t.data[start:t.pos]

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return fmt.Errorf("invalid number literal '%s'", numStr)
	}

	bits := math.Float64bits(num)

	t.emit(Node{
		Type: NodeTypeNumber,
		A:    uint32(bits >> 32),
		B:    uint32(bits),
	})

	return nil
}

func (t *parser) parseString(quoteType byte) error {
	// quoteType is either '"' or '\''
	start := t.pos

	for {
		idx := strings.IndexByte(t.data[t.pos:], quoteType)
		if idx == -1 {
			return fmt.Errorf("unterminated string literal")
		}

		absoluteIdx := t.pos + uint32(idx)

		// count consecutive backslashes immediately preceding the quote
		backslashes := 0
		for j := absoluteIdx - 1; j >= start && t.data[j] == '\\'; j-- {
			backslashes++
		}

		// consume the quote
		t.pos = absoluteIdx + 1

		// if even number of backslashes the quote is not escaped
		if backslashes%2 == 0 {
			break
		}

	}
	strVal := t.data[start : t.pos-1]

	// TODO: do we need to unescape newlines here?

	stringId := t.Interner.Intern(strVal)

	t.emit(Node{
		Type: NodeTypeString,
		A:    stringId,
	})

	return nil
}

func (t *parser) parseBlockString() error {

	start := t.pos

	for {
		// search for the first pipe
		idx := strings.IndexByte(t.data[t.pos:], '|')
		if idx == -1 {
			return fmt.Errorf("unterminated block string")
		}

		// jump to the found pipe
		t.pos += uint32(idx)

		// do we have enough characters left for "|||"?
		if int(t.pos+2) >= len(t.data) {
			return fmt.Errorf("unterminated block string")
		}

		if t.sliceOffset(3) == "|||" {
			break
		}

		// False alarm.
		// If it was a double pipe "||", we can skip past both!
		if t.data[t.pos+1] == '|' {
			t.pos += 2
		} else {
			t.pos++
		}
	}

	rawBlock := t.data[start:t.pos]
	t.pos += 3 // consume the three pipes

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
		stringId := t.Interner.Intern(rawBlock)
		t.emit(Node{Type: NodeTypeString, A: stringId})
		return nil
	}

	indentLen := len(indentStr)
	out := make([]byte, 0, len(rawBlock))
	pos := 0

	for pos < len(rawBlock) {
		// Find the end of the current line
		nextNL := strings.IndexByte(rawBlock[pos:], '\n')
		var line string

		if nextNL == -1 {
			line = rawBlock[pos:]
			pos = len(rawBlock)
		} else {
			line = rawBlock[pos : pos+nextNL]
			pos += nextNL + 1
		}

		// Strip the indent if the line starts with it
		// (Jsonnet spec: only strip if the line has the exact indent prefix.
		// If it's a blank line or differently indented, leave it alone or handle errors based on spec strictness).
		if strings.HasPrefix(line, indentStr) {
			out = append(out, line[indentLen:]...)
		} else {
			out = append(out, line...)
		}

		// Add the newline back
		if nextNL != -1 {
			out = append(out, '\n')
		}
	}

	finalStr := unsafe.String(unsafe.SliceData(out), len(out))

	stringId := t.Interner.Intern(finalStr)

	t.emit(Node{
		Type: NodeTypeString,
		A:    stringId,
	})

	return nil
}

func (t *parser) parseVerbatimString(quoteType byte) error {
	// quoteType is either '"' or '\''
	start := t.pos

	for {
		// Fast SIMD search for the quote
		idx := strings.IndexByte(t.data[t.pos:], quoteType)
		if idx == -1 {
			return fmt.Errorf("unterminated verbatim string literal")
		}

		absoluteIdx := t.pos + uint32(idx)

		// Peek at the NEXT character after the quote
		if int(absoluteIdx+1) < len(t.data) && t.data[absoluteIdx+1] == quoteType {
			// It's a double quote (e.g. "" or '')!
			// This is how verbatim strings escape quotes.
			// Skip over BOTH quotes and keep searching.
			t.pos = absoluteIdx + 2
			continue
		}

		// It's a single quote! We found the end of the string.
		t.pos = absoluteIdx + 1
		break
	}

	// Extract the string contents
	strVal := t.data[start : t.pos-1]

	// IMPORTANT: You must process `strVal` to replace double quotes with single quotes.
	// E.g., `""` becomes `"`, and `''` becomes `'`.
	// Since this is relatively rare, we only do the allocation/replace if necessary.
	if quoteType == '"' && strings.Contains(strVal, `""`) {
		strVal = strings.ReplaceAll(strVal, `""`, `"`)
	} else if quoteType == '\'' && strings.Contains(strVal, `''`) {
		strVal = strings.ReplaceAll(strVal, `''`, `'`)
	}

	stringId := t.Interner.Intern(strVal)

	t.emit(Node{
		Type: NodeTypeString,
		A:    stringId,
	})

	return nil
}
