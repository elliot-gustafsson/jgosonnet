package ast

import (
	"fmt"
	"os"
	"strconv"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type parser struct {
	lex       lexer
	peekToken Token
	hasPeek   bool

	data AST

	stack []uint32

	Interner *interner.Interner
}

func ParseSnippet(name, data string, interner *interner.Interner) (AST, uint32, error) {
	p := parser{
		lex:   lexer{data: data},
		data:  make(AST, 0, 1024),
		stack: make([]uint32, 0, 32),

		Interner: interner,
	}

	rootIdx, err := p.parseExpr(0)
	if err != nil {
		return nil, 0, err
	}

	tok, err := p.peek()
	if err != nil {
		return nil, 0, err
	}
	if tok.Kind != TokenEof {
		return nil, 0, fmt.Errorf("unexpected token after root expression: %v", tok.Data)
	}

	if len(p.data) == 0 {
		return nil, 0, fmt.Errorf("empty snippet")
	}

	return p.data, rootIdx, nil
}

func Parse(filename string, interner *interner.Interner) (AST, uint32, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, 0, err
	}

	data := unsafe.String(unsafe.SliceData(rawData), len(rawData))

	p := parser{
		lex:   lexer{data: data},
		data:  make(AST, 0, 1024),
		stack: make([]uint32, 0, 32),

		Interner: interner,
	}

	rootIdx, err := p.parseExpr(0)
	if err != nil {
		return nil, 0, err
	}

	tok, err := p.peek()
	if err != nil {
		return nil, 0, err
	}
	if tok.Kind != TokenEof {
		return nil, 0, fmt.Errorf("unexpected token after root expression: %v", tok.Data)
	}

	if len(p.data) == 0 {
		return nil, 0, fmt.Errorf("empty file")
	}

	return p.data, rootIdx, nil
}

func (p *parser) nextToken() (Token, error) {
	if p.hasPeek {
		p.hasPeek = false
		return p.peekToken, nil
	}
	return p.lex.Next()
}

func (p *parser) peek() (Token, error) {
	if p.hasPeek {
		return p.peekToken, nil
	}
	t, err := p.lex.Next()
	if err != nil {
		return Token{}, err
	}
	p.peekToken = t
	p.hasPeek = true
	return t, nil
}

func (p *parser) parseExpr(prec int) (uint32, error) {
	tok, err := p.nextToken()
	if err != nil {
		return 0, err
	}

	var lhs uint32

	switch tok.Kind {
	case TokenNumber:
		num, err := strconv.ParseFloat(tok.Data, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number '%s'", tok.Data)
		}
		lhs = p.data.emitNumber(num)
	case TokenString:
		stringId := p.Interner.Intern(tok.Data)
		return p.data.emitTagU32(TagString, stringId), nil
	case TokenTrue:
		return p.data.emitTag(TagTrue), nil
	case TokenFalse:
		return p.data.emitTag(TagFalse), nil
	case TokenNull:
		return p.data.emitTag(TagNull), nil
	case TokenSelf:
		return p.data.emitTag(TagSelf), nil
	case TokenIdent:
		stringId := p.Interner.Intern(tok.Data)
		return p.data.emitTagU32(TagVar, stringId), nil
	case TokenOperator:
		return p.parseOperator(tok.Data)
	case TokenLocal:
		return p.parseLocal()
	// case TokenIf:
	// 	return p.parseIf()
	case TokenFunction:
		return p.parseFunction()
	case TokenError:
		expr, err := p.parseExpr(0)
		if err != nil {
			return 0, err
		}
		return p.data.emitTagU32(TagError, expr), nil
	case TokenImport:
		strTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if strTok.Kind != TokenString {
			return 0, fmt.Errorf("expected string after import")
		}
		strId := p.Interner.Intern(strTok.Data)
		return p.data.emitTagU32(TagImport, strId), nil
	case TokenImportStr:
		strTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if strTok.Kind != TokenString {
			return 0, fmt.Errorf("expected string after importstr")
		}
		strId := p.Interner.Intern(strTok.Data)
		return p.data.emitTagU32(TagImportStr, strId), nil
	case TokenImportBin:
		strTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if strTok.Kind != TokenString {
			return 0, fmt.Errorf("expected string after importbin")
		}
		strId := p.Interner.Intern(strTok.Data)
		return p.data.emitTagU32(TagImportBin, strId), nil
	case TokenBracketL:
		lhs, err = p.parseArray()
		if err != nil {
			return 0, err
		}
	// case TokenBraceL:
	// 	lhs, err = p.parseObject()
	// 	if err != nil {
	// 		return 0, err
	// 	}
	case TokenParenL:
		expr, err := p.parseExpr(0)
		if err != nil {
			return 0, err
		}
		tok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if tok.Kind != TokenParenR {
			return 0, fmt.Errorf("expected ')' after expression")
		}
		lhs = expr
	default:
		return 0, fmt.Errorf("unexpected token in expression: %v (%v)", tok.Data, tok.Kind)
	}

	for {
		peekTok, err := p.peek()
		if err != nil {
			return 0, err
		}

		if peekTok.Kind == TokenEof {
			break
		}

		var opPrec int
		var isBinary bool
		var isApply bool
		var isIndex bool

		if peekTok.Kind == TokenOperator || peekTok.Kind == TokenIn {
			isBinary = true
			opPrec = getOpPrec(peekTok.Data)
		} else if peekTok.Kind == TokenParenL || peekTok.Kind == TokenBraceL {
			isApply = true
			opPrec = 13
		} else if peekTok.Kind == TokenBracketL || peekTok.Kind == TokenDot {
			isIndex = true
			opPrec = 13
		}

		if (!isBinary && !isApply && !isIndex) || opPrec <= prec {
			break
		}

		_, _ = p.nextToken()

		// if isBinary {
		// 	rhsPrec := opPrec
		// 	rhs, err := p.parseExpr(rhsPrec)
		// 	if err != nil {
		// 		return 0, err
		// 	}

		// 	nodeType := getBinaryNodeType(peekTok.Data)
		// 	lhs = p.emit(Node{Type: nodeType, A: lhs, B: rhs})
		// } else if isApply {
		// 	if peekTok.Kind == TokenParenL {
		// 		lhs, err = p.parseApply(lhs)
		// 		if err != nil {
		// 			return 0, err
		// 		}
		// 	} else {
		// 		// foo { bar: 2 } is syntactic sugar for foo + { bar: 2 }
		// 		rhs, err := p.parseObject()
		// 		if err != nil {
		// 			return 0, err
		// 		}
		// 		lhs = p.emit(Node{Type: NodeTypeBinaryPlus, A: lhs, B: rhs})
		// 	}
		// } else if isIndex {
		// 	if peekTok.Kind == TokenDot {
		// 		identTok, err := p.nextToken()
		// 		if err != nil {
		// 			return 0, err
		// 		}
		// 		if identTok.Kind != TokenIdent {
		// 			return 0, fmt.Errorf("expected identifier after dot")
		// 		}

		// 		stringId := p.Interner.Intern(identTok.Data)
		// 		stringNode := p.emit(Node{Type: NodeTypeString, A: stringId})
		// 		lhs = p.emit(Node{Type: NodeTypeIndex, A: stringNode, B: lhs})
		// 	} else {
		// 		rhs, err := p.parseExpr(0)
		// 		if err != nil {
		// 			return 0, err
		// 		}

		// 		closeTok, err := p.nextToken()
		// 		if err != nil {
		// 			return 0, err
		// 		}
		// 		if closeTok.Kind != TokenBracketR {
		// 			return 0, fmt.Errorf("expected ']' after index")
		// 		}

		// 		lhs = p.emit(Node{Type: NodeTypeIndex, A: rhs, B: lhs})
		// 	}
		// }
	}

	return lhs, nil
}

func getOpPrec(op string) int {
	switch op {
	case "||":
		return 2
	case "&&":
		return 3
	case "|":
		return 4
	case "^":
		return 5
	case "&":
		return 6
	case "==", "!=":
		return 7
	case "<", ">", "<=", ">=", "in":
		return 8
	case "<<", ">>":
		return 9
	case "+", "-":
		return 10
	case "*", "/", "%":
		return 11
	}
	return 0
}

func getBinaryTag(op string) Tag {
	switch op {
	case "*":
		return TagBinaryMul
	case "/":
		return TagBinaryDiv
	case "%":
		return TagBinaryMod
	case "+":
		return TagBinaryAdd
	case "-":
		return TagBinarySub
	case "<<":
		return TagBinaryShiftL
	case ">>":
		return TagBinaryShiftR
	case ">":
		return TagBinaryGt
	case ">=":
		return TagBinaryGte
	case "<":
		return TagBinaryLt
	case "<=":
		return TagBinaryLte
	case "==":
		return TagBinaryEq
	case "!=":
		return TagBinaryNeq
	case "in":
		return TagBinaryIn
	case "&":
		return TagBinaryBitAnd
	case "^":
		return TagBinaryBitXor
	case "|":
		return TagBinaryBitOr
	case "&&":
		return TagBinaryAnd
	case "||":
		return TagBinaryOr
	}
	return TagNone
}

func (p *parser) parseOperator(data string) (uint32, error) {
	var tag Tag
	switch data {
	case "+":
		tag = TagUnaryPlus
	case "-":
		tag = TagUnaryMinus
	case "!":
		tag = TagUnaryNot
	case "~":
		tag = TagUnaryBitwiseNot
	default:
		return 0, fmt.Errorf("unexpected unary operator '%s'", data)
	}

	expr, err := p.parseExpr(12)
	if err != nil {
		return 0, err
	}

	if tag == TagUnaryPlus {
		// plus is noop
		return expr, nil
	}

	return p.data.emitTagU32(tag, expr), nil
}

func (p *parser) parseLocal() (uint32, error) {
	// already consumed 'local'

	bindStart := len(p.stack)

	for {
		identTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if identTok.Kind != TokenIdent {
			return 0, fmt.Errorf("expected identifier in local binding")
		}

		nameId := p.Interner.Intern(identTok.Data)

		eqTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if eqTok.Kind != TokenOperator || eqTok.Data != "=" {
			return 0, fmt.Errorf("expected '=' in local binding")
		}

		exprId, err := p.parseExpr(0)
		if err != nil {
			return 0, err
		}
		p.stack = append(p.stack, nameId, exprId)

		tok, err := p.peek()
		if err != nil {
			return 0, err
		}

		if tok.Kind == TokenComma {
			p.nextToken()
			continue
		} else if tok.Kind == TokenSemicolon {
			p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("expected ',' or ';' after local binding, got %v", tok.Data)
		}
	}

	bodyId, err := p.parseExpr(0)
	if err != nil {
		return 0, err
	}

	bindCount := (len(p.stack) - bindStart) / 2

	var offset uint32
	if bindCount == 1 {
		offset = p.data.emitLocal(p.stack[bindStart], p.stack[bindStart+1], bodyId)
	} else {
		offset = p.data.emitLocals(p.stack[bindStart:], bodyId)
	}

	p.stack = p.stack[:bindStart]
	return offset, nil
}

func (p *parser) parseFunction() (uint32, error) {
	// already consumed 'function'

	t, err := p.nextToken()
	if err != nil {
		return 0, err
	}

	if t.Kind != TokenParenL {
		return 0, fmt.Errorf("expected '(' after 'function', got: %s", t.Data)
	}

	stackStart := len(p.stack)

	for {
		t, err := p.peek()
		if err != nil {
			return 0, err
		}

		if t.Kind == TokenParenR {
			_, _ = p.nextToken()
			break
		}

		t, err = p.nextToken()
		if err != nil {
			return 0, err
		}

		if t.Kind != TokenIdent {
			return 0, fmt.Errorf("expected indentifier for function parameter")
		}
		name_id := p.Interner.Intern(t.Data)
		// TODO: track locals and error on parse instead of runtime?

		t, err = p.peek()
		if err != nil {
			return 0, err
		}

		has_default := t.Kind == TokenOperator && t.Data == "="

		var defaultNodeOffset uint32
		if has_default {
			_, _ = p.nextToken() // consume =
			defaultNodeOffset, err = p.parseExpr(0)
			if err != nil {
				return 0, err
			}
		}

		p.stack = append(p.stack, name_id, defaultNodeOffset)

		t, err = p.peek()
		if err != nil {
			return 0, err
		}

		if t.Kind == TokenComma {
			_, _ = p.nextToken()
			continue
		} else if t.Kind == TokenParenR {
			_, _ = p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("unexpected token: %s", t.Data)
		}
	}

	bodyId, err := p.parseExpr(0)
	if err != nil {
		return 0, err
	}

	// argCount := (len(p.stack) - stackStart) / 2

	offset := p.data.emitFunction(p.stack[stackStart:], bodyId)

	p.stack = p.stack[:stackStart]

	return offset, nil
}

func (p *parser) parseArray() (uint32, error) {
	// [ consumed

	stackStart := len(p.stack)

	for {
		t, err := p.peek()
		if err != nil {
			return 0, err
		}

		if t.Kind == TokenBracketR {
			_, _ = p.nextToken()
			break
		}

		id, err := p.parseExpr(0)
		if err != nil {
			return 0, err
		}

		p.stack = append(p.stack, id)

		t, err = p.peek()
		if err != nil {
			return 0, err
		}

		if t.Kind == TokenComma {
			_, _ = p.nextToken()
			continue
		} else if t.Kind == TokenBracketR {
			_, _ = p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("unexpected token: %s", t.Data)
		}
	}

	elementCount := (len(p.stack) - stackStart)

	var offset uint32
	if elementCount == 0 {
		offset = p.data.emitTag(TagArray0)
	} else {
		offset = p.data.emitArray(p.stack[stackStart:])
	}

	p.stack = p.stack[:stackStart]

	return offset, nil
}
