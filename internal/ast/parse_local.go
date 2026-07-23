package ast

import "fmt"

func (p *parser) parseLocal() (uint32, error) {
	// already consumed 'local'

	var binds []uint32

	for {
		identTok, err := p.nextToken()
		if err != nil { return 0, err }
		if identTok.Kind != TokenIdent {
			return 0, fmt.Errorf("expected identifier in local binding")
		}

		nameId := p.Interner.Intern(identTok.Data)
		binds = append(binds, nameId)

		eqTok, err := p.nextToken()
		if err != nil { return 0, err }
		if eqTok.Kind != TokenOperator || eqTok.Data != "=" {
			return 0, fmt.Errorf("expected '=' in local binding")
		}

		exprId, err := p.parseExpr(0)
		if err != nil { return 0, err }
		binds = append(binds, exprId)

		tok, err := p.peek()
		if err != nil { return 0, err }
		
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
	if err != nil { return 0, err }

	startIdx := p.emitSideTable(binds...)
	numBinds := uint32(len(binds) / 2)

	return p.emit(Node{
		Type: NodeTypeLocal,
		A:    bodyId,
		B:    startIdx,
		C:    numBinds,
	}), nil
}
