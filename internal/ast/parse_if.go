package ast

import "fmt"

func (p *parser) parseIf() (uint32, error) {
	// already consumed 'if'
	condId, err := p.parseExpr(0)
	if err != nil { return 0, err }

	thenTok, err := p.nextToken()
	if err != nil { return 0, err }
	if thenTok.Kind != TokenThen {
		return 0, fmt.Errorf("expected 'then' after if condition")
	}

	trueId, err := p.parseExpr(0)
	if err != nil { return 0, err }

	var falseId uint32
	tok, err := p.peek()
	if err != nil { return 0, err }
	if tok.Kind == TokenElse {
		p.nextToken() // consume 'else'
		falseId, err = p.parseExpr(0)
		if err != nil { return 0, err }
	} else {
		falseId = p.emit(Node{Type: NodeTypeNull}) // implicitly null if no else
	}

	return p.emit(Node{
		Type: NodeTypeConditional,
		A:    condId,
		B:    trueId,
		C:    falseId,
	}), nil
}
