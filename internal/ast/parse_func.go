package ast

import "fmt"

func (p *parser) parseFunction() (uint32, error) {
	// already consumed 'function'
	
	tok, err := p.nextToken()
	if err != nil { return 0, err }
	if tok.Kind != TokenParenL {
		return 0, fmt.Errorf("expected '(' after 'function'")
	}

	type param struct {
		nameId  uint32
		defaultId uint32
	}
	var params []param

	for {
		tok, err := p.peek()
		if err != nil { return 0, err }
		
		if tok.Kind == TokenParenR {
			p.nextToken()
			break
		}
		
		identTok, err := p.nextToken()
		if err != nil { return 0, err }
		if identTok.Kind != TokenIdent {
			return 0, fmt.Errorf("expected identifier in function parameter")
		}
		
		nameId := p.Interner.Intern(identTok.Data)
		var defaultId uint32 // 0 implies no default (which is usually safe in these side tables unless 0 happens to be a valid AST index... wait, root/first node is index 0. Wait, AST node 0 is the root node usually? Actually in `jgosonnet`, 0 means empty? Let's check).
		// Wait, if 0 means no default, what if AST index 0 is valid?
		// Usually Node 0 is a dummy node, let's look at `Nodes: make([]Node, 1, 8192)`
		// Yes, Nodes starts with len 1! So id 0 is naturally the "null pointer". Perfect!

		tok, err = p.peek()
		if err != nil { return 0, err }
		if tok.Kind == TokenOperator && tok.Data == "=" {
			p.nextToken() // consume '='
			defaultId, err = p.parseExpr(0)
			if err != nil { return 0, err }
		}
		
		params = append(params, param{nameId, defaultId})
		
		tok, err = p.peek()
		if err != nil { return 0, err }
		if tok.Kind == TokenComma {
			p.nextToken()
			continue
		} else if tok.Kind == TokenParenR {
			p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("expected ',' or ')' in function parameters, got %v", tok.Data)
		}
	}
	
	bodyId, err := p.parseExpr(0)
	if err != nil { return 0, err }
	
	startIdx := uint32(len(p.SideTable))
	for _, param := range params {
		p.SideTable = append(p.SideTable, param.nameId, param.defaultId)
	}
	
	return p.emit(Node{
		Type: NodeTypeFunction,
		A:    bodyId,
		B:    startIdx,
		C:    uint32(len(params)),
	}), nil
}
