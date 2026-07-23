package ast

import "fmt"

func (p *parser) parseArray() (uint32, error) {
	// already consumed '['

	var elements []uint32

	for {
		tok, err := p.peek()
		if err != nil { return 0, err }
		if tok.Kind == TokenBracketR {
			p.nextToken()
			break
		}

		exprId, err := p.parseExpr(0)
		if err != nil { return 0, err }
		elements = append(elements, exprId)

		tok, err = p.peek()
		if err != nil { return 0, err }
		
		if tok.Kind == TokenComma {
			p.nextToken()
			continue
		} else if tok.Kind == TokenBracketR {
			p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("expected ',' or ']' in array, got %v", tok.Data)
		}
	}

	startIdx := uint32(len(p.SideTable))
	p.SideTable = append(p.SideTable, elements...)
	
	return p.emit(Node{
		Type: NodeTypeArray,
		A:    startIdx,
		B:    uint32(len(elements)),
	}), nil
}
