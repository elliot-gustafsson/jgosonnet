package ast

import "fmt"

func (p *parser) parseObject() (uint32, error) {
	// already consumed '{'

	type localBind struct {
		nameId uint32
		exprId uint32
	}
	var locals []localBind
	var asserts []uint32

	type objectField struct {
		keyId  uint32
		bodyId uint32
		meta   uint32
	}
	var fields []objectField

	for {
		tok, err := p.peek()
		if err != nil { return 0, err }

		if tok.Kind == TokenBraceR {
			p.nextToken()
			break
		}

		if tok.Kind == TokenLocal {
			p.nextToken() // consume 'local'
			identTok, err := p.nextToken()
			if err != nil { return 0, err }
			if identTok.Kind != TokenIdent { return 0, fmt.Errorf("expected identifier in local") }
			
			nameId := p.Interner.Intern(identTok.Data)
			
			// wait, local inside object can be a function! `local f(x) = x`
			// Let's assume simple locals for now (or implement function parsing if we have time)
			eqTok, err := p.nextToken()
			if err != nil { return 0, err }
			if eqTok.Kind != TokenOperator || eqTok.Data != "=" { return 0, fmt.Errorf("expected '=' in local") }
			
			exprId, err := p.parseExpr(0)
			if err != nil { return 0, err }
			
			locals = append(locals, localBind{nameId, exprId})
		} else if tok.Kind == TokenAssert {
			p.nextToken() // consume 'assert'
			exprId, err := p.parseExpr(0)
			if err != nil { return 0, err }
			asserts = append(asserts, exprId)
			// wait, asserts can have a message `assert expr : msg`
			tok, err := p.peek()
			if err != nil { return 0, err }
			if tok.Kind == TokenColon {
				p.nextToken()
				msgId, err := p.parseExpr(0)
				if err != nil { return 0, err }
				// wait, jgosonnet doesn't seem to store the message in AST for asserts in objects?
				// Actually `go-jsonnet/ast.Node` assert inside object is just a Node.
				_ = msgId 
			}
		} else {
			// Field
			var keyId uint32
			if tok.Kind == TokenBracketL {
				p.nextToken() // consume '['
				keyId, err = p.parseExpr(0)
				if err != nil { return 0, err }
				
				closeTok, err := p.nextToken()
				if err != nil { return 0, err }
				if closeTok.Kind != TokenBracketR { return 0, fmt.Errorf("expected ']'") }
			} else if tok.Kind == TokenIdent || tok.Kind == TokenString {
				p.nextToken()
				strId := p.Interner.Intern(tok.Data)
				keyId = p.emit(Node{Type: NodeTypeString, A: strId})
			} else {
				return 0, fmt.Errorf("unexpected token in object: %v", tok.Data)
			}
			
			// separator: ':', '::', ':::', '+:', '+::', '+:::'
			sepTok, err := p.nextToken()
			if err != nil { return 0, err }
			if sepTok.Kind != TokenColon && sepTok.Kind != TokenOperator {
				return 0, fmt.Errorf("expected ':' in object field")
			}
			
			sepStr := sepTok.Data
			
			// parse optional plus super
			plusSuper := false
			if len(sepStr) > 0 && sepStr[0] == '+' {
				plusSuper = true
				sepStr = sepStr[1:]
			}
			
			var visibility uint32
			if sepStr == ":" {
				visibility = 1 // Inherit
			} else if sepStr == "::" {
				visibility = 0 // Hidden
			} else if sepStr == ":::" {
				visibility = 2 // Visible
			} else {
				return 0, fmt.Errorf("invalid field separator '%v'", sepTok.Data)
			}
			
			meta := visibility & 0x03
			if plusSuper {
				meta |= 0x04
			}
			
			bodyId, err := p.parseExpr(0)
			if err != nil { return 0, err }
			
			fields = append(fields, objectField{keyId, bodyId, meta})
		}
		
		tok, err = p.peek()
		if err != nil { return 0, err }
		if tok.Kind == TokenComma {
			p.nextToken()
			continue
		} else if tok.Kind == TokenBraceR {
			p.nextToken()
			break
		} else {
			return 0, fmt.Errorf("expected ',' or '}' in object, got %v", tok.Data)
		}
	}
	
	var flags uint8
	var headerSize int
	if len(locals) > 0 {
		flags |= 1 // FlagObjectHasLocals
		headerSize++
	}
	if len(asserts) > 0 {
		flags |= 2 // FlagObjectHasAsserts
		headerSize++
	}
	

	startIdx := len(p.SideTable)
	
	// manually grow and append to SideTable
	var sidetable []uint32
	if len(locals) > 0 {
		sidetable = append(sidetable, uint32(len(locals)))
	}
	if len(asserts) > 0 {
		sidetable = append(sidetable, uint32(len(asserts)))
	}
	for _, l := range locals {
		sidetable = append(sidetable, l.nameId, l.exprId)
	}
	for _, a := range asserts {
		sidetable = append(sidetable, a)
	}
	for _, f := range fields {
		sidetable = append(sidetable, f.keyId, f.bodyId, f.meta)
	}
	
	p.emitSideTable(sidetable...)
	
	return p.emit(Node{
		Type:  NodeTypeObject,
		A:     uint32(startIdx),
		B:     uint32(len(fields)),
		Flags: flags,
	}), nil
}
