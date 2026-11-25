package ast

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

const anonymousBind = "anonymous"

type parser struct {
	lex       lexer
	peekBuf   [2]Token
	peekCount int

	data     AST
	locs     []LocEntry
	contexts []string

	stack []uint32

	Interner *interner.Interner

	objLevel   int
	usedDollar bool
	curCtxIdx  uint32 // 1-based index into contexts; 0 means ""
	curBind    string

	lastStartByte uint32
	lastEndByte   uint32
	lastBareStr   bool
	lastBlockStr  bool
}

func ParseProgramSnippet(name, data string, in *interner.Interner) (*Program, error) {
	p := parser{
		lex:      lexer{filename: name, data: data},
		data:     make(AST, 0, max(64, len(data)/2)),
		locs:     make([]LocEntry, 0, max(32, len(data)/8)),
		contexts: []string{"$"},
		stack:    make([]uint32, 0, 64),
		Interner: in,
		curCtxIdx: 1, // "$"
		curBind:   anonymousBind,
	}

	rootIdx, err := p.parseExpr(0)
	if err != nil {
		return nil, err
	}

	tok, err := p.peek()
	if err != nil {
		return nil, err
	}
	if tok.Kind != TokenEof {
		return nil, p.staticErrTok(tok, fmt.Sprintf("Did not expect: %s", tok.Display()))
	}

	if len(p.data) == 0 {
		return nil, StaticError{Loc: byteRangeToLoc(name, data, 0, 0), Msg: "Unexpected end of file"}
	}

	prog := &Program{
		Filename: name,
		Source:   data,
		Nodes:    p.data,
		RootIdx:  rootIdx,
		Locs:     p.locs,
		Contexts: p.contexts,
	}

	if err := analyzeProgram(prog, in); err != nil {
		return nil, err
	}

	return prog, nil
}

func ParseProgram(filename string, in *interner.Interner) (*Program, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	data := unsafe.String(unsafe.SliceData(rawData), len(rawData))
	return ParseProgramSnippet(filename, data, in)
}

func ParseSnippet(name, data string, in *interner.Interner) (AST, uint32, error) {
	prog, err := ParseProgramSnippet(name, data, in)
	if err != nil {
		return nil, 0, err
	}
	return prog.Nodes, prog.RootIdx, nil
}

func Parse(filename string, in *interner.Interner) (AST, uint32, error) {
	prog, err := ParseProgram(filename, in)
	if err != nil {
		return nil, 0, err
	}
	return prog.Nodes, prog.RootIdx, nil
}

func (p *parser) addContext(name string) uint32 {
	p.contexts = append(p.contexts, name)
	return uint32(len(p.contexts))
}

func (p *parser) currentContextStr() string {
	if p.curCtxIdx == 0 || int(p.curCtxIdx) > len(p.contexts) {
		return ""
	}
	return p.contexts[p.curCtxIdx-1]
}

func (p *parser) recordLoc(off, startByte, endByte, ctxIdx uint32) {
	p.locs = append(p.locs, LocEntry{
		NodeOffset: off,
		StartByte:  startByte,
		EndByte:    endByte,
		ContextIdx: ctxIdx,
	})
	p.lastStartByte = startByte
	p.lastEndByte = endByte
}

func (p *parser) staticErr(start, end uint32, msg string) error {
	return StaticError{
		Loc: byteRangeToLoc(p.lex.filename, p.lex.data, start, end),
		Msg: msg,
	}
}

func (p *parser) staticErrTok(tok Token, msg string) error {
	return p.staticErr(tok.Pos, tok.End, msg)
}

func (p *parser) nextToken() (Token, error) {
	if p.peekCount > 0 {
		t := p.peekBuf[0]
		p.peekBuf[0] = p.peekBuf[1]
		p.peekCount--
		return t, nil
	}
	return p.lex.Next()
}

func (p *parser) peek() (Token, error) {
	if p.peekCount >= 1 {
		return p.peekBuf[0], nil
	}
	t, err := p.lex.Next()
	if err != nil {
		return Token{}, err
	}
	p.peekBuf[0] = t
	p.peekCount = 1
	return t, nil
}

func (p *parser) peek2() (Token, error) {
	if p.peekCount < 1 {
		if _, err := p.peek(); err != nil {
			return Token{}, err
		}
	}
	if p.peekCount >= 2 {
		return p.peekBuf[1], nil
	}
	t, err := p.lex.Next()
	if err != nil {
		return Token{}, err
	}
	p.peekBuf[1] = t
	p.peekCount = 2
	return t, nil
}

func (p *parser) parseExpr(prec int) (uint32, error) {
	tok, err := p.nextToken()
	if err != nil {
		return 0, err
	}

	p.lastBareStr = false
	p.lastBlockStr = false

	var lhs uint32
	lhsStart := tok.Pos

	switch tok.Kind {
	case TokenEof:
		return 0, p.staticErrTok(tok, "Unexpected end of file")

	case TokenNumber:
		num, err := strconv.ParseFloat(tok.Data, 64)
		if err != nil {
			if math.IsInf(num, 0) || !errors.Is(err, strconv.ErrRange) {
				return 0, p.staticErrTok(tok, fmt.Sprintf("failed to parse float val (%s), err: %v", tok.Data, err))
			}
			// Subnormal underflow to 0 is allowed in Jsonnet
		}
		lhs = p.data.emitNumber(num)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenString:
		stringId := p.Interner.Intern(tok.Data)
		lhs = p.data.emitTagU32(TagString, stringId)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)
		p.lastBareStr = true

	case TokenStringBlock:
		stringId := p.Interner.Intern(tok.Data)
		lhs = p.data.emitTagU32(TagString, stringId)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)
		p.lastBareStr = true
		p.lastBlockStr = true

	case TokenTrue:
		lhs = p.data.emitTag(TagTrue)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenFalse:
		lhs = p.data.emitTag(TagFalse)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenNull:
		lhs = p.data.emitTag(TagNull)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenSelf:
		lhs = p.data.emitTag(TagSelf)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenDollar:
		if p.objLevel == 0 {
			return 0, p.staticErrTok(tok, "No top-level object found.")
		}
		p.usedDollar = true
		dollarId := p.Interner.Intern("$")
		lhs = p.data.emitTagU32(TagVar, dollarId)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenSuper:
		next, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		switch next.Kind {
		case TokenDot:
			fieldTok, err := p.nextToken()
			if err != nil {
				return 0, err
			}
			if fieldTok.Kind != TokenIdent {
				return 0, p.staticErrTok(fieldTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", fieldTok.Display()))
			}
			keyId := p.Interner.Intern(fieldTok.Data)
			lhs = p.data.emitTagU32(TagSuperIndexConst, keyId)
			p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)
		case TokenBracketL:
			prevBind := p.curBind
			p.curBind = anonymousBind
			idxOff, err := p.parseExpr(0)
			p.curBind = prevBind
			if err != nil {
				return 0, err
			}
			closeTok, err := p.nextToken()
			if err != nil {
				return 0, err
			}
			if closeTok.Kind != TokenBracketR {
				return 0, p.staticErrTok(closeTok, fmt.Sprintf("Expected token %s but got %s", TokenBracketR.String(), closeTok.Display()))
			}
			lhs = p.data.emitTagU32(TagSuperIndex, idxOff)
			p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)
		default:
			return 0, p.staticErrTok(tok, "Expected . or [ after super")
		}

	case TokenIdent:
		stringId := p.Interner.Intern(tok.Data)
		lhs = p.data.emitTagU32(TagVar, stringId)
		p.recordLoc(lhs, tok.Pos, tok.End, p.curCtxIdx)

	case TokenOperator:
		lhs, err = p.parseUnary(tok)
		if err != nil {
			return 0, err
		}

	case TokenLocal:
		return p.parseLocal(tok)

	case TokenIf:
		return p.parseIf(tok)

	case TokenAssert:
		return p.parseAssert(tok)

	case TokenFunction:
		return p.parseFunction(tok)

	case TokenError:
		prevBind := p.curBind
		p.curBind = anonymousBind
		expr, err := p.parseExpr(0)
		p.curBind = prevBind
		if err != nil {
			return 0, err
		}
		endByte := p.lastEndByte
		off := p.data.emitTagU32(TagError, expr)
		p.recordLoc(off, tok.Pos, endByte, p.curCtxIdx)
		return off, nil

	case TokenImport:
		return p.parseImport(tok, TagImport)
	case TokenImportStr:
		return p.parseImport(tok, TagImportStr)
	case TokenImportBin:
		return p.parseImport(tok, TagImportBin)

	case TokenBracketL:
		prevBind := p.curBind
		p.curBind = anonymousBind
		lhs, err = p.parseArray(tok)
		p.curBind = prevBind
		if err != nil {
			return 0, err
		}

	case TokenBraceL:
		lhs, err = p.parseObject(tok)
		if err != nil {
			return 0, err
		}

	case TokenParenL:
		prevBind := p.curBind
		p.curBind = anonymousBind
		expr, err := p.parseExpr(0)
		p.curBind = prevBind
		if err != nil {
			return 0, err
		}
		closeTok, err := p.nextToken()
		if err != nil {
			return 0, err
		}
		if closeTok.Kind != TokenParenR {
			return 0, p.staticErrTok(closeTok, fmt.Sprintf("Expected token %s but got %s", TokenParenR.String(), closeTok.Display()))
		}
		p.lastStartByte = tok.Pos
		p.lastEndByte = closeTok.End
		p.lastBareStr = false
		p.lastBlockStr = false
		lhs = expr

	default:
		return 0, p.staticErrTok(tok, fmt.Sprintf("Unexpected: %s while parsing terminal", tok.Display()))
	}

	// Pratt infix & postfix loop
	for {
		peekTok, err := p.peek()
		if err != nil {
			return 0, err
		}

		if peekTok.Kind == TokenEof {
			break
		}

		var opPrec int
		var kind uint8 // 1=binary, 2=apply/brace, 3=index/dot

		switch peekTok.Kind {
		case TokenOperator:
			if peekTok.Data == ":" || peekTok.Data == "::" {
				// Delimiter for assert message or slice syntax
				return lhs, nil
			}
			opPrec = getOpPrec(peekTok.Data)
			if opPrec == 0 {
				return 0, p.staticErrTok(peekTok, fmt.Sprintf("Not a binary operator: %s", peekTok.Data))
			}
			kind = 1
		case TokenIn:
			opPrec = 8
			kind = 1
		case TokenParenL, TokenBraceL:
			opPrec = 13
			kind = 2
		case TokenBracketL, TokenDot:
			opPrec = 13
			kind = 3
		default:
			return lhs, nil
		}

		if opPrec <= prec {
			break
		}

		_, _ = p.nextToken()
		p.lastBareStr = false
		p.lastBlockStr = false

		prevBind := p.curBind
		p.curBind = anonymousBind

		switch kind {
		case 1: // Binary or `in super`
			if peekTok.Kind == TokenIn {
				next, err := p.peek()
				if err != nil {
					p.curBind = prevBind
					return 0, err
				}
				if next.Kind == TokenSuper {
					superTok, _ := p.nextToken()
					lhs = p.data.emitTagU32(TagInSuper, lhs)
					p.recordLoc(lhs, lhsStart, superTok.End, p.curCtxIdx)
					p.lastBareStr = false
					p.lastBlockStr = false
					p.curBind = prevBind
					continue
				}
			}

			rhs, err := p.parseExpr(opPrec)
			if err != nil {
				p.curBind = prevBind
				return 0, err
			}
			endByte := p.lastEndByte
			tag := getBinaryTag(peekTok.Data)
			lhs = p.data.emitTag2U32(tag, lhs, rhs)
			ctxIdx := p.curCtxIdx
			if tag == TagBinaryMod {
				ctxIdx = 0
			}
			p.recordLoc(lhs, lhsStart, endByte, ctxIdx)

		case 2: // Apply `(...)` or ApplyBrace `{...}`
			if peekTok.Kind == TokenParenL {
				lhs, err = p.parseApply(lhs, lhsStart)
				if err != nil {
					p.curBind = prevBind
					return 0, err
				}
			} else {
				// `lhs { ... }` desugars to `lhs + { ... }`
				rhs, err := p.parseObject(peekTok)
				if err != nil {
					p.curBind = prevBind
					return 0, err
				}
				endByte := p.lastEndByte
				lhs = p.data.emitTag2U32(TagBinaryAdd, lhs, rhs)
				p.recordLoc(lhs, lhsStart, endByte, p.curCtxIdx)
			}

		case 3: // `.ident` or `[index]` or `[begin:end:step]`
			if peekTok.Kind == TokenDot {
				identTok, err := p.nextToken()
				if err != nil {
					p.curBind = prevBind
					return 0, err
				}
				if identTok.Kind != TokenIdent {
					p.curBind = prevBind
					return 0, p.staticErrTok(identTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", identTok.Display()))
				}
				keyId := p.Interner.Intern(identTok.Data)
				lhs = p.data.emitTag2U32(TagIndexConst, lhs, keyId)
				p.recordLoc(lhs, lhsStart, identTok.End, p.curCtxIdx)
			} else {
				lhs, err = p.parseIndexOrSlice(lhs, lhsStart)
				if err != nil {
					p.curBind = prevBind
					return 0, err
				}
			}
		}

		p.lastBareStr = false
		p.lastBlockStr = false
		p.curBind = prevBind
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

func (p *parser) parseUnary(tok Token) (uint32, error) {
	var tag Tag
	switch tok.Data {
	case "+":
		tag = TagUnaryPlus
	case "-":
		tag = TagUnaryMinus
	case "!":
		tag = TagUnaryNot
	case "~":
		tag = TagUnaryBitwiseNot
	default:
		return 0, p.staticErrTok(tok, fmt.Sprintf("Not a unary operator: %s", tok.Data))
	}

	prevBind := p.curBind
	p.curBind = anonymousBind
	expr, err := p.parseExpr(12)
	p.curBind = prevBind
	if err != nil {
		return 0, err
	}

	endByte := p.lastEndByte
	off := p.data.emitTagU32(tag, expr)
	p.recordLoc(off, tok.Pos, endByte, p.curCtxIdx)
	return off, nil
}

func (p *parser) parseImport(tok Token, tag Tag) (uint32, error) {
	dataStart := len(p.data)
	locsStart := len(p.locs)

	prevBind := p.curBind
	p.curBind = anonymousBind
	bodyOff, err := p.parseExpr(0)
	p.curBind = prevBind
	if err != nil {
		return 0, err
	}

	bodyStart := p.lastStartByte
	bodyEnd := p.lastEndByte

	if !p.lastBareStr {
		return 0, p.staticErr(bodyStart, bodyEnd, "Computed imports are not allowed")
	}
	if p.lastBlockStr {
		return 0, p.staticErr(bodyStart, bodyEnd, "Block string literals not allowed in imports")
	}

	// Extract the interned string ID from the bare TagString node and replace it with TagImport*.
	strId := p.data[bodyOff+1]
	p.data = p.data[:dataStart]
	p.locs = p.locs[:locsStart]

	off := p.data.emitTagU32(tag, strId)
	p.recordLoc(off, tok.Pos, bodyEnd, p.curCtxIdx)
	return off, nil
}

func (p *parser) parseIf(tok Token) (uint32, error) {
	prevBind := p.curBind
	p.curBind = anonymousBind

	condOff, err := p.parseExpr(0)
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}

	thenTok, err := p.nextToken()
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	if thenTok.Kind != TokenThen {
		p.curBind = prevBind
		return 0, p.staticErrTok(thenTok, fmt.Sprintf("Expected token then but got %s", thenTok.Display()))
	}

	trueOff, err := p.parseExpr(0)
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	endByte := p.lastEndByte

	var falseOff uint32
	next, err := p.peek()
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	if next.Kind == TokenElse {
		_, _ = p.nextToken()
		falseOff, err = p.parseExpr(0)
		if err != nil {
			p.curBind = prevBind
			return 0, err
		}
		endByte = p.lastEndByte
	} else {
		falseOff = p.data.emitTag(TagNull)
	}

	p.curBind = prevBind
	off := p.data.emitTag3U32(TagConditional, condOff, trueOff, falseOff)
	p.recordLoc(off, tok.Pos, endByte, p.curCtxIdx)
	return off, nil
}

func (p *parser) parseAssert(tok Token) (uint32, error) {
	prevBind := p.curBind
	p.curBind = anonymousBind

	condOff, err := p.parseExpr(0)
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}

	msgOff := NoNode
	next, err := p.peek()
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	if next.Kind == TokenOperator && next.Data == ":" {
		_, _ = p.nextToken()
		msgOff, err = p.parseExpr(0)
		if err != nil {
			p.curBind = prevBind
			return 0, err
		}
	}

	semiTok, err := p.nextToken()
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	if semiTok.Kind != TokenSemicolon {
		p.curBind = prevBind
		return 0, p.staticErrTok(semiTok, fmt.Sprintf("Expected token %s but got %s", TokenSemicolon.String(), semiTok.Display()))
	}

	restOff, err := p.parseExpr(0)
	if err != nil {
		p.curBind = prevBind
		return 0, err
	}
	endByte := p.lastEndByte

	p.curBind = prevBind
	off := p.data.emitTag3U32(TagAssert, condOff, msgOff, restOff)
	// Context 0 ("") on TagAssert itself matches go-jsonnet's desugared error frame
	p.recordLoc(off, tok.Pos, endByte, 0)
	return off, nil
}

func (p *parser) parseLocal(localTok Token) (uint32, error) {
	bindStart := len(p.stack)

	for {
		identTok, err := p.nextToken()
		if err != nil {
			p.stack = p.stack[:bindStart]
			return 0, err
		}
		if identTok.Kind != TokenIdent {
			p.stack = p.stack[:bindStart]
			return 0, p.staticErrTok(identTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", identTok.Display()))
		}

		nameId := p.Interner.Intern(identTok.Data)

		// Check duplicate local var in same local binding list
		for i := bindStart; i < len(p.stack); i += 2 {
			if p.stack[i] == nameId {
				p.stack = p.stack[:bindStart]
				return 0, p.staticErrTok(identTok, fmt.Sprintf("Duplicate local var: %s", identTok.Data))
			}
		}

		// Check for local function sugar: `local f(x, y=1) = expr`
		next, err := p.peek()
		if err != nil {
			p.stack = p.stack[:bindStart]
			return 0, err
		}

		var exprId uint32
		prevCtx := p.curCtxIdx
		prevBind := p.curBind

		if next.Kind == TokenParenL {
			_, _ = p.nextToken() // consume '('
			p.curBind = identTok.Data
			exprId, err = p.parseMethodBody(identTok.Pos, true)
		} else {
			eqTok, err2 := p.nextToken()
			if err2 != nil {
				p.stack = p.stack[:bindStart]
				return 0, err2
			}
			if eqTok.Kind != TokenOperator || eqTok.Data != "=" {
				p.stack = p.stack[:bindStart]
				return 0, p.staticErrTok(eqTok, fmt.Sprintf("Expected operator = but got %s", eqTok.Display()))
			}
			p.curBind = identTok.Data
			p.curCtxIdx = p.addContext("thunk <" + identTok.Data + "> from <" + p.currentContextStr() + ">")
			exprId, err = p.parseExpr(0)
		}

		p.curCtxIdx = prevCtx
		p.curBind = prevBind
		if err != nil {
			p.stack = p.stack[:bindStart]
			return 0, err
		}

		p.stack = append(p.stack, nameId, exprId)

		delim, err := p.nextToken()
		if err != nil {
			p.stack = p.stack[:bindStart]
			return 0, err
		}

		if delim.Kind == TokenComma {
			continue
		} else if delim.Kind == TokenSemicolon {
			break
		} else {
			p.stack = p.stack[:bindStart]
			return 0, p.staticErrTok(delim, fmt.Sprintf("Expected , or ; but got %s", delim.Display()))
		}
	}

	bodyId, err := p.parseExpr(0)
	if err != nil {
		p.stack = p.stack[:bindStart]
		return 0, err
	}
	endByte := p.lastEndByte

	bindCount := (len(p.stack) - bindStart) / 2
	var offset uint32
	if bindCount == 1 {
		offset = p.data.emitLocal(p.stack[bindStart], p.stack[bindStart+1], bodyId)
	} else {
		offset = p.data.emitLocals(p.stack[bindStart:], bodyId)
	}

	p.stack = p.stack[:bindStart]
	p.recordLoc(offset, localTok.Pos, endByte, p.curCtxIdx)
	return offset, nil
}

// parseParameters parses function/method parameters after '(' has been consumed.
// Appends pairs (nameId, defaultOffOrNoNode) to p.stack and returns the closing ')' token.
func (p *parser) parseParameters(elementKind string) (Token, error) {
	gotComma := false
	first := true
	for {
		t, err := p.peek()
		if err != nil {
			return Token{}, err
		}
		if t.Kind == TokenParenR {
			return p.nextToken()
		}
		if !first && !gotComma {
			return Token{}, p.staticErrTok(t, fmt.Sprintf("Expected a comma before next %s, got %s", elementKind, t.Display()))
		}

		identTok, err := p.nextToken()
		if err != nil {
			return Token{}, err
		}
		if identTok.Kind != TokenIdent {
			return Token{}, p.staticErrTok(identTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", identTok.Display()))
		}
		nameId := p.Interner.Intern(identTok.Data)

		defaultOff := NoNode
		t, err = p.peek()
		if err != nil {
			return Token{}, err
		}
		if t.Kind == TokenOperator && t.Data == "=" {
			_, _ = p.nextToken() // consume '='
			defaultOff, err = p.parseExpr(0)
			if err != nil {
				return Token{}, err
			}
		}

		p.stack = append(p.stack, nameId, defaultOff)

		t, err = p.peek()
		if err != nil {
			return Token{}, err
		}
		if t.Kind == TokenComma {
			_, _ = p.nextToken()
			gotComma = true
		} else {
			gotComma = false
		}
		first = false
	}
}

func (p *parser) parseFunction(fnTok Token) (uint32, error) {
	t, err := p.nextToken()
	if err != nil {
		return 0, err
	}
	if t.Kind != TokenParenL {
		return 0, p.staticErrTok(t, fmt.Sprintf("Expected ( but got %s", t.Display()))
	}

	fnCtxIdx := p.addContext("function <" + p.curBind + ">")
	prevCtx := p.curCtxIdx
	prevBind := p.curBind
	p.curCtxIdx = fnCtxIdx
	p.curBind = anonymousBind

	stackStart := len(p.stack)
	_, err = p.parseParameters("function parameter")
	if err != nil {
		p.stack = p.stack[:stackStart]
		p.curCtxIdx = prevCtx
		p.curBind = prevBind
		return 0, err
	}

	bodyId, err := p.parseExpr(0)
	endByte := p.lastEndByte
	p.curCtxIdx = prevCtx
	p.curBind = prevBind
	if err != nil {
		p.stack = p.stack[:stackStart]
		return 0, err
	}

	offset := p.data.emitFunction(p.stack[stackStart:], bodyId)
	p.stack = p.stack[:stackStart]
	p.recordLoc(offset, fnTok.Pos, endByte, prevCtx)
	return offset, nil
}

// parseMethodBody parses `(params) = expr` (for local functions) after '(' has been consumed.
func (p *parser) parseMethodBody(startPos uint32, isLocalFn bool) (uint32, error) {
	fnCtxIdx := p.addContext("function <" + p.curBind + ">")
	prevCtx := p.curCtxIdx
	prevBind := p.curBind
	p.curCtxIdx = fnCtxIdx
	p.curBind = anonymousBind

	stackStart := len(p.stack)
	_, err := p.parseParameters("function parameter")
	if err != nil {
		p.stack = p.stack[:stackStart]
		p.curCtxIdx = prevCtx
		p.curBind = prevBind
		return 0, err
	}

	if isLocalFn {
		eqTok, err := p.nextToken()
		if err != nil {
			p.stack = p.stack[:stackStart]
			p.curCtxIdx = prevCtx
			p.curBind = prevBind
			return 0, err
		}
		if eqTok.Kind != TokenOperator || eqTok.Data != "=" {
			p.stack = p.stack[:stackStart]
			p.curCtxIdx = prevCtx
			p.curBind = prevBind
			return 0, p.staticErrTok(eqTok, fmt.Sprintf("Expected operator = but got %s", eqTok.Display()))
		}
	}

	bodyId, err := p.parseExpr(0)
	endByte := p.lastEndByte
	p.curCtxIdx = prevCtx
	p.curBind = prevBind
	if err != nil {
		p.stack = p.stack[:stackStart]
		return 0, err
	}

	offset := p.data.emitFunction(p.stack[stackStart:], bodyId)
	p.stack = p.stack[:stackStart]
	p.recordLoc(offset, startPos, endByte, prevCtx)
	return offset, nil
}

func (p *parser) parseApply(fnOff, lhsStart uint32) (uint32, error) {
	// '(' already consumed
	posStart := len(p.stack)
	var namedArgs []uint32

	gotComma := false
	namedAdded := false
	first := true
	var endTok Token

	for {
		next, err := p.peek()
		if err != nil {
			p.stack = p.stack[:posStart]
			return 0, err
		}
		if next.Kind == TokenParenR {
			endTok, _ = p.nextToken()
			break
		}
		if !first && !gotComma {
			p.stack = p.stack[:posStart]
			return 0, p.staticErrTok(next, fmt.Sprintf("Expected a comma before next function argument, got %s", next.Display()))
		}

		// Lookahead for named argument: `ident = expr`
		isNamed := false
		var nameId uint32
		if next.Kind == TokenIdent {
			t2, err := p.peek2()
			if err != nil {
				p.stack = p.stack[:posStart]
				return 0, err
			}
			if t2.Kind == TokenOperator && t2.Data == "=" {
				isNamed = true
				identTok, _ := p.nextToken()
				_, _ = p.nextToken() // consume '='
				nameId = p.Interner.Intern(identTok.Data)
			}
		}

		if !isNamed && namedAdded {
			p.stack = p.stack[:posStart]
			return 0, p.staticErrTok(next, "Positional argument after a named argument is not allowed")
		}

		argOff, err := p.parseExpr(0)
		if err != nil {
			p.stack = p.stack[:posStart]
			return 0, err
		}

		if isNamed {
			namedAdded = true
			namedArgs = append(namedArgs, nameId, argOff)
		} else {
			p.stack = append(p.stack, argOff)
		}

		t, err := p.peek()
		if err != nil {
			p.stack = p.stack[:posStart]
			return 0, err
		}
		if t.Kind == TokenComma {
			_, _ = p.nextToken()
			gotComma = true
		} else {
			gotComma = false
		}
		first = false
	}

	tailStrict := false
	t, err := p.peek()
	if err != nil {
		p.stack = p.stack[:posStart]
		return 0, err
	}
	if t.Kind == TokenTailStrict {
		_, _ = p.nextToken()
		tailStrict = true
	}

	off := p.data.emitApply(fnOff, p.stack[posStart:], namedArgs, tailStrict)
	p.stack = p.stack[:posStart]
	p.recordLoc(off, lhsStart, endTok.End, p.curCtxIdx)
	return off, nil
}

func (p *parser) parseIndexOrSlice(targetOff, lhsStart uint32) (uint32, error) {
	// '[' already consumed
	indexes := [3]uint32{NoNode, NoNode, NoNode}
	colonsConsumed := 0
	readyForNextIndex := true
	var endTok Token

	for colonsConsumed < 3 {
		pt, err := p.peek()
		if err != nil {
			return 0, err
		}
		if pt.Kind == TokenBracketR {
			endTok, _ = p.nextToken()
			break
		} else if pt.Kind == TokenOperator && pt.Data == ":" {
			endTok, _ = p.nextToken()
			colonsConsumed++
			readyForNextIndex = true
		} else if pt.Kind == TokenOperator && pt.Data == "::" {
			endTok, _ = p.nextToken()
			colonsConsumed += 2
			readyForNextIndex = true
		} else if readyForNextIndex {
			idxOff, err := p.parseExpr(0)
			if err != nil {
				return 0, err
			}
			indexes[colonsConsumed] = idxOff
			readyForNextIndex = false
		} else {
			return 0, p.staticErrTok(pt, fmt.Sprintf("Expected token %s but got %s", TokenBracketR.String(), pt.Display()))
		}
	}

	if colonsConsumed > 2 {
		return 0, p.staticErrTok(endTok, "Invalid slice: too many colons")
	}
	if colonsConsumed == 0 && readyForNextIndex {
		return 0, p.staticErrTok(endTok, "ast.Index requires an expression")
	}

	if colonsConsumed == 0 {
		off := p.data.emitTag2U32(TagIndex, targetOff, indexes[0])
		p.recordLoc(off, lhsStart, endTok.End, p.curCtxIdx)
		return off, nil
	}

	// Slice desugars to `$std.slice(target, begin, end, step)`
	for i := range 3 {
		if indexes[i] == NoNode {
			indexes[i] = p.data.emitTag(TagNull)
		}
	}
	stdVar := p.data.emitTagU32(TagVar, p.Interner.Intern("$std"))
	sliceFn := p.data.emitTag2U32(TagIndexConst, stdVar, p.Interner.Intern("slice"))
	off := p.data.emitApply4(sliceFn, targetOff, indexes[0], indexes[1], indexes[2])
	p.recordLoc(off, lhsStart, endTok.End, p.curCtxIdx)
	return off, nil
}

func (p *parser) parseArray(openTok Token) (uint32, error) {
	// '[' already consumed
	t, err := p.peek()
	if err != nil {
		return 0, err
	}
	if t.Kind == TokenBracketR {
		closeTok, _ := p.nextToken()
		off := p.data.emitTag(TagArray0)
		p.recordLoc(off, openTok.Pos, closeTok.End, p.curCtxIdx)
		return off, nil
	}

	firstOff, err := p.parseExpr(0)
	if err != nil {
		return 0, err
	}

	gotComma := false
	t, err = p.peek()
	if err != nil {
		return 0, err
	}
	if t.Kind == TokenComma {
		_, _ = p.nextToken()
		gotComma = true
		t, err = p.peek()
		if err != nil {
			return 0, err
		}
	}

	if t.Kind == TokenFor {
		_, _ = p.nextToken() // consume 'for'
		clauseCount, clauses, closeTok, err := p.parseComprehensionSpecs(TokenBracketR)
		if err != nil {
			return 0, err
		}
		off := p.data.emitArrayComp(firstOff, clauseCount, clauses)
		p.recordLoc(off, openTok.Pos, closeTok.End, 0)
		return off, nil
	}

	stackStart := len(p.stack)
	p.stack = append(p.stack, firstOff)

	var closeTok Token
	for {
		t, err := p.peek()
		if err != nil {
			p.stack = p.stack[:stackStart]
			return 0, err
		}
		if t.Kind == TokenBracketR {
			closeTok, _ = p.nextToken()
			break
		}
		if !gotComma {
			p.stack = p.stack[:stackStart]
			return 0, p.staticErrTok(t, "Expected a comma before next array element")
		}

		elemOff, err := p.parseExpr(0)
		if err != nil {
			p.stack = p.stack[:stackStart]
			return 0, err
		}
		p.stack = append(p.stack, elemOff)

		t, err = p.peek()
		if err != nil {
			p.stack = p.stack[:stackStart]
			return 0, err
		}
		if t.Kind == TokenComma {
			_, _ = p.nextToken()
			gotComma = true
		} else {
			gotComma = false
		}
	}

	off := p.data.emitArray(p.stack[stackStart:])
	p.stack = p.stack[:stackStart]
	p.recordLoc(off, openTok.Pos, closeTok.End, p.curCtxIdx)
	return off, nil
}

// parseComprehensionSpecs parses after the first 'for' has been consumed:
//   ident in expr [if cond]* [for ident in expr [if cond]*]* <endKind>
// Returns (clauseCount, flatClauses, closeTok, err).
func (p *parser) parseComprehensionSpecs(endKind TokenKind) (uint32, []uint32, Token, error) {
	var clauses []uint32
	var clauseCount uint32

	for {
		varTok, err := p.nextToken()
		if err != nil {
			return 0, nil, Token{}, err
		}
		if varTok.Kind != TokenIdent {
			return 0, nil, Token{}, p.staticErrTok(varTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", varTok.Display()))
		}
		varId := p.Interner.Intern(varTok.Data)

		inTok, err := p.nextToken()
		if err != nil {
			return 0, nil, Token{}, err
		}
		if inTok.Kind != TokenIn {
			return 0, nil, Token{}, p.staticErrTok(inTok, fmt.Sprintf("Expected token in but got %s", inTok.Display()))
		}

		iterOff, err := p.parseExpr(0)
		if err != nil {
			return 0, nil, Token{}, err
		}

		clauseStart := len(clauses)
		clauses = append(clauses, varId, iterOff, 0) // placeholder for condCount
		var condCount uint32

		next, err := p.nextToken()
		if err != nil {
			return 0, nil, Token{}, err
		}
		for next.Kind == TokenIf {
			condOff, err := p.parseExpr(0)
			if err != nil {
				return 0, nil, Token{}, err
			}
			clauses = append(clauses, condOff)
			condCount++

			next, err = p.nextToken()
			if err != nil {
				return 0, nil, Token{}, err
			}
		}

		clauses[clauseStart+2] = condCount
		clauseCount++

		if next.Kind == endKind {
			return clauseCount, clauses, next, nil
		}
		if next.Kind != TokenFor {
			return 0, nil, Token{}, p.staticErrTok(next, fmt.Sprintf("Expected for, if or %s after for clause, got: %s", endKind.String(), next.Display()))
		}
	}
}

func (p *parser) parseObjectAssignmentOp() (plusSuper bool, hide FieldHide, err error) {
	op, err := p.nextToken()
	if err != nil {
		return false, 0, err
	}
	if op.Kind != TokenOperator {
		return false, 0, p.staticErrTok(op, fmt.Sprintf("Expected token OPERATOR but got %s", op.Display()))
	}

	s := op.Data
	if len(s) > 0 && s[0] == '+' {
		plusSuper = true
		s = s[1:]
	}

	switch s {
	case ":":
		return plusSuper, FieldInherit, nil
	case "::":
		return plusSuper, FieldHidden, nil
	case ":::":
		return plusSuper, FieldVisible, nil
	default:
		return false, 0, p.staticErrTok(op, fmt.Sprintf("Expected one of :, ::, :::, +:, +::, +:::, got: %s", op.Data))
	}
}

func (p *parser) parseObject(openTok Token) (uint32, error) {
	// '{' already consumed
	outerObjLevel := p.objLevel
	outerUsedDollar := p.usedDollar
	outerCtxIdx := p.curCtxIdx
	outerBind := p.curBind

	objCtxIdx := p.addContext("object <" + outerBind + ">")

	if outerObjLevel == 0 {
		p.usedDollar = false
	}

	restoreOuter := func() {
		p.objLevel = outerObjLevel
		if outerObjLevel == 0 {
			p.usedDollar = outerUsedDollar
		}
		p.curCtxIdx = outerCtxIdx
		p.curBind = outerBind
	}

	var fields []uint32  // triples: (keyOrOff, bodyOff, meta)
	var locals []uint32  // pairs:   (keyId, bodyOff)
	var asserts []uint32 // offsets of assert nodes
	allFieldsStatic := true

	gotComma := false
	first := true

	for {
		next, err := p.nextToken()
		if err != nil {
			restoreOuter()
			return 0, err
		}

		if next.Kind == TokenBraceR {
			innerUsedDollar := p.usedDollar
			restoreOuter()

			if outerObjLevel == 0 && innerUsedDollar {
				dollarId := p.Interner.Intern("$")
				selfOff := p.data.emitTag(TagSelf)
				locals = append(locals, dollarId, selfOff)
			}

			var off uint32
			if len(fields) == 0 && len(locals) == 0 && len(asserts) == 0 {
				off = p.data.emitTag(TagObject0)
			} else if allFieldsStatic && len(locals) == 0 && len(asserts) == 0 {
				// Strip FieldFlagStaticKey from meta for TagObjectSimple
				for i := 2; i < len(fields); i += 3 {
					fields[i] &^= uint32(FieldFlagStaticKey)
				}
				off = p.data.emitObjectSimple(fields)
			} else {
				off = p.data.emitObject(fields, locals, asserts)
			}
			p.recordLoc(off, openTok.Pos, next.End, outerCtxIdx)
			return off, nil
		}

		if next.Kind == TokenFor {
			// Object comprehension!
			innerUsedDollar := p.usedDollar

			if len(asserts) > 0 {
				restoreOuter()
				return 0, p.staticErrTok(next, "Object comprehension cannot have asserts")
			}
			if len(fields) != 3 {
				restoreOuter()
				return 0, p.staticErrTok(next, "Object comprehension can only have one field")
			}
			fieldKeyOff := fields[0]
			fieldBodyOff := fields[1]
			fieldMeta := fields[2]
			if FieldHide(uint8(fieldMeta)&FieldMaskVisibility) != FieldInherit {
				restoreOuter()
				return 0, p.staticErrTok(next, "Object comprehensions cannot have hidden fields")
			}
			if (uint8(fieldMeta) & FieldFlagStaticKey) != 0 {
				restoreOuter()
				return 0, p.staticErrTok(next, "Object comprehensions can only have [e] fields")
			}

			// Parse comprehension clauses at OUTER objLevel & context
			p.objLevel = outerObjLevel
			p.curCtxIdx = outerCtxIdx
			p.curBind = anonymousBind
			clauseCount, clauses, closeTok, err := p.parseComprehensionSpecs(TokenBraceR)
			if err != nil {
				restoreOuter()
				return 0, err
			}

			if outerObjLevel == 0 && innerUsedDollar {
				dollarId := p.Interner.Intern("$")
				selfOff := p.data.emitTag(TagSelf)
				locals = append(locals, dollarId, selfOff)
			}

			// Wrap object locals around field body (each iteration binds its own locals)
			if len(locals) == 2 {
				fieldBodyOff = p.data.emitLocal(locals[0], locals[1], fieldBodyOff)
			} else if len(locals) > 2 {
				fieldBodyOff = p.data.emitLocals(locals, fieldBodyOff)
			}

			restoreOuter()
			off := p.data.emitObjectComp(fieldKeyOff, fieldBodyOff, fieldMeta&uint32(FieldMaskVisibility|FieldFlagPlusSuper), clauseCount, clauses)
			p.recordLoc(off, openTok.Pos, closeTok.End, 0)
			return off, nil
		}

		if !gotComma && !first {
			restoreOuter()
			return 0, p.staticErrTok(next, "Expected a comma before next field")
		}

		switch next.Kind {
		case TokenBracketL, TokenIdent, TokenString, TokenStringBlock:
			var keyOrOff uint32
			var isStatic bool
			var fieldNameStr string

			if next.Kind == TokenBracketL {
				// Computed field name is evaluated OUTSIDE this object!
				p.objLevel = outerObjLevel
				p.curCtxIdx = outerCtxIdx
				p.curBind = anonymousBind
				keyOrOff, err = p.parseExpr(0)
				if err != nil {
					restoreOuter()
					return 0, err
				}
				closeBracket, err := p.nextToken()
				if err != nil {
					restoreOuter()
					return 0, err
				}
				if closeBracket.Kind != TokenBracketR {
					restoreOuter()
					return 0, p.staticErrTok(closeBracket, fmt.Sprintf("Expected token %s but got %s", TokenBracketR.String(), closeBracket.Display()))
				}
				allFieldsStatic = false
			} else {
				isStatic = true
				fieldNameStr = next.Data
				keyOrOff = p.Interner.Intern(fieldNameStr)
			}

			// Method sugar: `foo(a, b=1): body`
			p.objLevel = outerObjLevel + 1
			p.curCtxIdx = objCtxIdx
			p.curBind = anonymousBind

			isMethod := false
			paramStart := len(p.stack)
			var methodCtxIdx uint32
			pt, err := p.peek()
			if err != nil {
				restoreOuter()
				return 0, err
			}
			if pt.Kind == TokenParenL {
				_, _ = p.nextToken() // consume '('
				isMethod = true
				methodCtxIdx = p.addContext("function <anonymous>")
				p.curCtxIdx = methodCtxIdx
				_, err = p.parseParameters("method parameter")
				if err != nil {
					p.stack = p.stack[:paramStart]
					restoreOuter()
					return 0, err
				}
			}

			plusSuper, hide, err := p.parseObjectAssignmentOp()
			if err != nil {
				p.stack = p.stack[:paramStart]
				restoreOuter()
				return 0, err
			}

			if plusSuper && isMethod {
				p.stack = p.stack[:paramStart]
				restoreOuter()
				return 0, p.staticErrTok(next, fmt.Sprintf("Cannot use +: syntax sugar in a method: %s", next.Data))
			}

			if isStatic {
				for i := 0; i < len(fields); i += 3 {
					if (uint8(fields[i+2])&FieldFlagStaticKey) != 0 && fields[i] == keyOrOff {
						p.stack = p.stack[:paramStart]
						restoreOuter()
						return 0, p.staticErrTok(next, fmt.Sprintf("Duplicate field: %s", fieldNameStr))
					}
				}
			}

			bodyOff, err := p.parseExpr(0)
			if err != nil {
				p.stack = p.stack[:paramStart]
				restoreOuter()
				return 0, err
			}
			bodyEnd := p.lastEndByte

			if isMethod {
				fnOff := p.data.emitFunction(p.stack[paramStart:], bodyOff)
				p.stack = p.stack[:paramStart]
				p.recordLoc(fnOff, next.Pos, bodyEnd, objCtxIdx)
				bodyOff = fnOff
			}

			meta := uint8(hide) & FieldMaskVisibility
			if plusSuper {
				meta |= FieldFlagPlusSuper
			}
			if isStatic {
				meta |= FieldFlagStaticKey
			}
			fields = append(fields, keyOrOff, bodyOff, uint32(meta))

		case TokenLocal:
			p.objLevel = outerObjLevel + 1
			p.curCtxIdx = objCtxIdx
			p.curBind = anonymousBind

			varTok, err := p.nextToken()
			if err != nil {
				restoreOuter()
				return 0, err
			}
			if varTok.Kind != TokenIdent {
				restoreOuter()
				return 0, p.staticErrTok(varTok, fmt.Sprintf("Expected token IDENTIFIER but got %s", varTok.Display()))
			}
			nameId := p.Interner.Intern(varTok.Data)

			for i := 0; i < len(locals); i += 2 {
				if locals[i] == nameId {
					restoreOuter()
					return 0, p.staticErrTok(varTok, fmt.Sprintf("Duplicate local var: %s", varTok.Data))
				}
			}

			var bodyOff uint32
			pt, err := p.peek()
			if err != nil {
				restoreOuter()
				return 0, err
			}
			if pt.Kind == TokenParenL {
				_, _ = p.nextToken() // consume '('
				bodyOff, err = p.parseMethodBody(varTok.Pos, true)
			} else {
				eqTok, err2 := p.nextToken()
				if err2 != nil {
					restoreOuter()
					return 0, err2
				}
				if eqTok.Kind != TokenOperator || eqTok.Data != "=" {
					restoreOuter()
					return 0, p.staticErrTok(eqTok, fmt.Sprintf("Expected operator = but got %s", eqTok.Display()))
				}
				bodyOff, err = p.parseExpr(0)
			}
			if err != nil {
				restoreOuter()
				return 0, err
			}

			locals = append(locals, nameId, bodyOff)

		case TokenAssert:
			p.objLevel = outerObjLevel + 1
			p.curCtxIdx = objCtxIdx
			p.curBind = anonymousBind

			condOff, err := p.parseExpr(0)
			if err != nil {
				restoreOuter()
				return 0, err
			}
			assertEnd := p.lastEndByte

			msgOff := NoNode
			pt, err := p.peek()
			if err != nil {
				restoreOuter()
				return 0, err
			}
			if pt.Kind == TokenOperator && pt.Data == ":" {
				_, _ = p.nextToken()
				msgOff, err = p.parseExpr(0)
				if err != nil {
					restoreOuter()
					return 0, err
				}
				assertEnd = p.lastEndByte
			} else {
				msgId := p.Interner.Intern("Object assertion failed.")
				msgOff = p.data.emitTagU32(TagString, msgId)
			}

			trueOff := p.data.emitTag(TagTrue)
			assertOff := p.data.emitTag3U32(TagAssert, condOff, msgOff, trueOff)
			p.recordLoc(assertOff, next.Pos, assertEnd, 0)
			asserts = append(asserts, assertOff)

		default:
			restoreOuter()
			return 0, p.staticErrTok(next, fmt.Sprintf("Unexpected: %s while parsing field definition", next.Display()))
		}

		pt, err := p.peek()
		if err != nil {
			restoreOuter()
			return 0, err
		}
		if pt.Kind == TokenComma {
			_, _ = p.nextToken()
			gotComma = true
		} else {
			gotComma = false
		}
		first = false
	}
}
