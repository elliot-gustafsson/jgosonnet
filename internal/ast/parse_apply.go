package ast

import "fmt"

func (p *parser) parseApply(funcId uint32) (uint32, error) {
	// already consumed '('

	type arg struct {
		nameId uint32
		exprId uint32
	}
	var args []arg
	
	gotNamed := false

	for {
		tok, err := p.peek()
		if err != nil { return 0, err }
		if tok.Kind == TokenParenR {
			p.nextToken()
			break
		}
		
		// Could be named arg or positional arg.
		// If it's `id = expr`, it's named.
		// Since `id = expr` isn't a valid expression in Jsonnet itself (no assignment expressions),
		// we can peek ahead.
		
		var nameId uint32 = 0
		var exprId uint32
		
		// To distinguish, we need to peek 2 tokens. Or we just look at the first token.
		// Since we only have 1 token peek, we could do a manual lookahead, or just try parsing an expression
		// but wait! If we parse an expression `x`, it might just be the variable `x`.
		// But if the next token is `=`, we've gone too far.
		// So if `tok.Kind == TokenIdent`, let's peek the next one manually if possible.
		// Let's implement double peek or just buffer. 
		// For simplicity, since the lexer is stateless other than `pos`, we can just peek and peek next, but we don't have that method.
		
		if tok.Kind == TokenIdent {
			p.nextToken() // temporarily consume it
			nextTok, _ := p.peek()
			if nextTok.Kind == TokenOperator && nextTok.Data == "=" {
				p.nextToken() // consume '='
				nameId = p.Interner.Intern(tok.Data)
				exprId, err = p.parseExpr(0)
				if err != nil { return 0, err }
				gotNamed = true
			} else {
				// It wasn't named. We need to parse an expression, but we already consumed the first identifier!
				// A Pratt parser normally wants us to pass the token in, or we can backtrack the lexer.
				// Lexer backtrack is super easy: `p.lex.pos = tok.Pos` and clear `peekToken`.
				p.lex.pos = tok.Pos
				p.peekToken = nil
				
				if gotNamed { return 0, fmt.Errorf("positional argument after named argument") }
				
				exprId, err = p.parseExpr(0)
				if err != nil { return 0, err }
			}
		} else {
			if gotNamed { return 0, fmt.Errorf("positional argument after named argument") }
			exprId, err = p.parseExpr(0)
			if err != nil { return 0, err }
		}

		args = append(args, arg{nameId, exprId})
		
		tok, err = p.peek()
		if err != nil { return 0, err }
		
		if tok.Kind == TokenComma {
			p.nextToken()
			continue
		} else if tok.Kind == TokenParenR {
			p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("expected ',' or ')' in function arguments")
		}
	}
	
	startIdx := uint32(len(p.SideTable))
	for _, a := range args {
		p.SideTable = append(p.SideTable, a.nameId, a.exprId)
	}
	
	return p.emit(Node{
		Type: NodeTypeApply,
		A:    funcId,
		B:    startIdx,
		C:    uint32(len(args)),
	}), nil
}
