package ast

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type parser struct {
	lex       lexer
	peekToken *Token

	Nodes     []Node
	SideTable []uint32
	Locations []NodeContext

	Interner *interner.Interner
}


func ParseSnippet(name, data string, interner *interner.Interner) (*AST, error) {
	p := parser{
		lex: lexer{data: data},
		Interner:  interner,
		Nodes:     make([]Node, 1, max(64, len(data)/4)),
		Locations: make([]NodeContext, 1, max(64, len(data)/4)),
		SideTable: make([]uint32, 0, max(32, len(data)/8)),
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
		return nil, fmt.Errorf("unexpected token after root expression: %v", tok.Data)
	}

	if len(p.Nodes) == 0 {
		return nil, fmt.Errorf("empty snippet")
	}

	return &AST{
		RootId:    rootIdx,
		Nodes:     p.Nodes,
		SideTable: p.SideTable,
		Locations: p.Locations,
	}, nil
}

func Parse(filename string, interner *interner.Interner) (*AST, error) {
	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := unsafe.String(unsafe.SliceData(rawData), len(rawData))

	p := parser{
		lex: lexer{data: data},

		Interner:  interner,
		Nodes:     make([]Node, 1, max(64, len(data)/4)),
		Locations: make([]NodeContext, 1, max(64, len(data)/4)),
		SideTable: make([]uint32, 0, max(32, len(data)/8)),
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
		return nil, fmt.Errorf("unexpected token after root expression: %v", tok.Data)
	}

	if len(p.Nodes) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	return &AST{
		RootId:    rootIdx,
		Nodes:     p.Nodes,
		SideTable: p.SideTable,
		Locations: p.Locations,
	}, nil
}

func (p *parser) nextToken() (Token, error) {
	if p.peekToken != nil {
		t := *p.peekToken
		p.peekToken = nil
		return t, nil
	}
	return p.lex.Next()
}

func (p *parser) peek() (Token, error) {
	if p.peekToken != nil {
		return *p.peekToken, nil
	}
	t, err := p.lex.Next()
	if err != nil {
		return Token{}, err
	}
	p.peekToken = &t
	return t, nil
}

func (p *parser) emit(n Node) uint32 {
	id := uint32(len(p.Nodes))
	p.Nodes = append(p.Nodes, n)
	p.Locations = append(p.Locations, NodeContext{}) 
	return id
}

func (p *parser) emitSideTable(vals ...uint32) uint32 {
	start := uint32(len(p.SideTable))
	p.SideTable = append(p.SideTable, vals...)
	return start
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
		bits := math.Float64bits(num)
		lhs = p.emit(Node{
			Type: NodeTypeNumber,
			A:    uint32(bits >> 32),
			B:    uint32(bits),
		})
	case TokenString:
		stringId := p.Interner.Intern(tok.Data)
		lhs = p.emit(Node{Type: NodeTypeString, A: stringId})
	case TokenTrue:
		lhs = p.emit(Node{Type: NodeTypeTrue})
	case TokenFalse:
		lhs = p.emit(Node{Type: NodeTypeFalse})
	case TokenNull:
		lhs = p.emit(Node{Type: NodeTypeNull})
	case TokenSelf:
		lhs = p.emit(Node{Type: NodeTypeSelf})
	case TokenIdent:
		stringId := p.Interner.Intern(tok.Data)
		lhs = p.emit(Node{Type: NodeTypeVar, A: stringId})
	case TokenLocal:
		return p.parseLocal()
	case TokenIf:
		return p.parseIf()
	case TokenFunction:
		return p.parseFunction()
	case TokenError:
		expr, err := p.parseExpr(0)
		if err != nil { return 0, err }
		lhs = p.emit(Node{Type: NodeTypeError, A: expr})
	case TokenImport:
		strTok, err := p.nextToken()
		if err != nil { return 0, err }
		if strTok.Kind != TokenString { return 0, fmt.Errorf("expected string after import") }
		strId := p.Interner.Intern(strTok.Data)
		lhs = p.emit(Node{Type: NodeTypeImport, A: strId})
	case TokenImportStr:
		strTok, err := p.nextToken()
		if err != nil { return 0, err }
		if strTok.Kind != TokenString { return 0, fmt.Errorf("expected string after importstr") }
		strId := p.Interner.Intern(strTok.Data)
		lhs = p.emit(Node{Type: NodeTypeImportStr, A: strId})
	case TokenBracketL:
		lhs, err = p.parseArray()
		if err != nil { return 0, err }
	case TokenBraceL:
		lhs, err = p.parseObject()
		if err != nil { return 0, err }
	case TokenParenL:
		expr, err := p.parseExpr(0)
		if err != nil { return 0, err }
		tok, err := p.nextToken()
		if err != nil { return 0, err }
		if tok.Kind != TokenParenR { return 0, fmt.Errorf("expected ')' after expression") }
		lhs = expr

	case TokenOperator:
		if tok.Data == "!" || tok.Data == "-" || tok.Data == "+" || tok.Data == "~" {
			expr, err := p.parseExpr(12)
			if err != nil { return 0, err }
			if tok.Data == "+" {
				lhs = expr
			} else {
				nodeType := NodeTypeUnaryMinus
				if tok.Data == "!" {
					nodeType = NodeTypeUnaryNot
				} else if tok.Data == "~" {
					nodeType = NodeTypeUnaryBitwiseNot
				}
				lhs = p.emit(Node{Type: nodeType, A: expr})
			}
		} else {
			return 0, fmt.Errorf("unexpected unary operator '%s'", tok.Data)
		}
	default:
		return 0, fmt.Errorf("unexpected token in expression: %v (%v)", tok.Data, tok.Kind)
	}

	for {
		peekTok, err := p.peek()
		if err != nil { return 0, err }

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

		if isBinary {
			rhsPrec := opPrec
			rhs, err := p.parseExpr(rhsPrec)
			if err != nil { return 0, err }
			
			nodeType := getBinaryNodeType(peekTok.Data)
			lhs = p.emit(Node{Type: nodeType, A: lhs, B: rhs})
		} else if isApply {
			if peekTok.Kind == TokenParenL {
				lhs, err = p.parseApply(lhs)
				if err != nil { return 0, err }
			} else {
								// foo { bar: 2 } is syntactic sugar for foo + { bar: 2 }
				rhs, err := p.parseObject()
				if err != nil { return 0, err }
				lhs = p.emit(Node{Type: NodeTypeBinaryPlus, A: lhs, B: rhs})
			}
		} else if isIndex {
			if peekTok.Kind == TokenDot {
				identTok, err := p.nextToken()
				if err != nil { return 0, err }
				if identTok.Kind != TokenIdent { return 0, fmt.Errorf("expected identifier after dot") }
				
				stringId := p.Interner.Intern(identTok.Data)
				stringNode := p.emit(Node{Type: NodeTypeString, A: stringId})
				lhs = p.emit(Node{Type: NodeTypeIndex, A: stringNode, B: lhs})
			} else {
				rhs, err := p.parseExpr(0)
				if err != nil { return 0, err }
				
				closeTok, err := p.nextToken()
				if err != nil { return 0, err }
				if closeTok.Kind != TokenBracketR { return 0, fmt.Errorf("expected ']' after index") }
				
				lhs = p.emit(Node{Type: NodeTypeIndex, A: rhs, B: lhs})
			}
		}
	}

	return lhs, nil
}

func getOpPrec(op string) int {
	switch op {
	case "||": return 2
	case "&&": return 3
	case "|": return 4
	case "^": return 5
	case "&": return 6
	case "==", "!=": return 7
	case "<", ">", "<=", ">=", "in": return 8
	case "<<", ">>": return 9
	case "+", "-": return 10
	case "*", "/", "%": return 11
	}
	return 0
}

func getBinaryNodeType(op string) NodeType {
	switch op {
	case "*": return NodeTypeBinaryMult
	case "/": return NodeTypeBinaryDiv
	case "+": return NodeTypeBinaryPlus
	case "-": return NodeTypeBinaryMinus
	case "<<": return NodeTypeBinaryShiftL
	case ">>": return NodeTypeBinaryShiftR
	case ">": return NodeTypeBinaryGreater
	case ">=": return NodeTypeBinaryGreaterEq
	case "<": return NodeTypeBinaryLess
	case "<=": return NodeTypeBinaryLessEq
	case "==": return NodeTypeBinaryEqual
	case "!=": return NodeTypeBinaryUnequal
	case "&": return NodeTypeBinaryBitwiseAnd
	case "^": return NodeTypeBinaryBitwiseXor
	case "|": return NodeTypeBinaryBitwiseOr
	case "&&": return NodeTypeBinaryAnd
	case "||": return NodeTypeBinaryOr
	}
	return NodeTypeError
}
