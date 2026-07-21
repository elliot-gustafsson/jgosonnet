package ast

import (
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet/ast"
)

type NodeType uint8

const (
	NodeTypeNone NodeType = iota

	// literals
	NodeTypeString
	NodeTypeNull
	NodeTypeTrue
	NodeTypeFalse
	NodeTypeNumber

	NodeTypeObject
	NodeTypeArray
	NodeTypeLocal
	NodeTypeApply
	NodeTypeIndex
	NodeTypeVar
	NodeTypeFunction
	NodeTypeConditional
	NodeTypeImport
	NodeTypeImportStr
	NodeTypeSelf
	NodeTypeSuperIndex
	NodeTypeInSuper
	NodeTypeError

	// binary ops
	NodeTypeBinaryAnd
	NodeTypeBinaryOr
	NodeTypeBinaryEqual
	NodeTypeBinaryPlus
	NodeTypeBinaryMinus
	NodeTypeBinaryDiv
	NodeTypeBinaryMult
	NodeTypeBinaryBitwiseAnd
	NodeTypeBinaryBitwiseOr
	NodeTypeBinaryBitwiseXor
	NodeTypeBinaryShiftL
	NodeTypeBinaryShiftR
	NodeTypeBinaryUnequal
	NodeTypeBinaryGreater
	NodeTypeBinaryGreaterEq
	NodeTypeBinaryLess
	NodeTypeBinaryLessEq

	// unary ops
	NodeTypeUnaryNot
	NodeTypeUnaryMinus
	NodeTypeUnaryBitwiseNot
)

type FieldVisibility uint8

const (
	FieldVisibilityHidden FieldVisibility = iota
	FieldVisibilityInherit
	FieldVisibilityVisible
)

const (
	FlagObjectHasLocals  uint8 = 1 << 0
	FlagObjectHasAsserts uint8 = 1 << 1
)

type Node struct {
	Type  NodeType
	Flags uint8
	_     uint16
	A     uint32
	B     uint32
	C     uint32
}

type AST struct {
	Nodes     []Node
	SideTable []uint32
	RootId    uint32

	Locations []NodeContext
}

type AstBuilder struct {
	Nodes []Node

	SideTable []uint32

	Locations []NodeContext

	Interner *interner.Interner
}

func NewAstBuilder(interner *interner.Interner) *AstBuilder {
	return &AstBuilder{
		Interner:  interner,
		Nodes:     make([]Node, 1, 8192),
		Locations: make([]NodeContext, 1, 8192),
		SideTable: make([]uint32, 0, 4096),
	}
}

func (b *AstBuilder) Parse(n ast.Node) (*AST, error) {
	rootId, err := b.visit(n)
	if err != nil {
		return nil, err
	}

	tree := &AST{
		RootId:    rootId,
		Nodes:     b.Nodes,
		SideTable: b.SideTable,
		Locations: b.Locations,
	}

	return tree, nil
}

func (b *AstBuilder) emit(n Node, loc NodeContext) uint32 {
	id := uint32(len(b.Nodes))
	b.Nodes = append(b.Nodes, n)
	b.Locations = append(b.Locations, loc)
	return id
}

func (b *AstBuilder) visit(n ast.Node) (uint32, error) {

	var loc NodeContext
	if lr := n.Loc(); lr != nil {
		loc = NodeContext{
			Context: n.Context(),
			File:    lr.FileName,
			Begin:   Location{Line: uint32(lr.Begin.Line), Column: uint32(lr.Begin.Column)},
			End:     Location{Line: uint32(lr.End.Line), Column: uint32(lr.End.Column)},
		}

		// if lr.File != nil {
		// 	loc.File = string(lr.File.DiagnosticFileName)
		// }

	}
	// b.Locations = append(b.Locations, loc)

	switch node := n.(type) {
	default:
		return 0, fmt.Errorf("unhandled node type: %T", node)
	case *ast.LiteralString:
		id := b.Interner.Intern(node.Value)
		return b.emit(Node{Type: NodeTypeString, A: id}, loc), nil

	case *ast.LiteralNull:
		return b.emit(Node{Type: NodeTypeNull}, loc), nil

	case *ast.LiteralBoolean:
		if node.Value {
			return b.emit(Node{Type: NodeTypeTrue}, loc), nil
		}
		return b.emit(Node{Type: NodeTypeFalse}, loc), nil

	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		bits := math.Float64bits(num)

		res := Node{
			Type: NodeTypeNumber,
			A:    uint32(bits >> 32),
			B:    uint32(bits),
		}
		return b.emit(res, loc), nil

	case *ast.DesugaredObject:

		numFields := len(node.Fields)
		numLocals := len(node.Locals)
		numAsserts := len(node.Asserts)

		if numFields > math.MaxUint16 {
			return 0, fmt.Errorf("object exceeds maximum of %d fields", math.MaxUint16)
		}

		var flags uint8
		var headerSize int

		if numLocals > 0 {
			flags |= FlagObjectHasLocals
			headerSize++
		}
		if numAsserts > 0 {
			flags |= FlagObjectHasAsserts
			headerSize++
		}

		totalSize := headerSize + (3 * numFields) + (2 * numLocals) + numAsserts

		startIdx := len(b.SideTable)
		b.SideTable = slices.Grow(b.SideTable, totalSize)[:startIdx+totalSize]

		offset := startIdx

		if numLocals > 0 {
			b.SideTable[offset] = uint32(numLocals)
			offset++
		}
		if numAsserts > 0 {
			b.SideTable[offset] = uint32(numAsserts)
			offset++
		}

		for i := range node.Locals {
			nameId := b.Interner.Intern(string(node.Locals[i].Variable))
			bodyId, err := b.visit(node.Locals[i].Body)
			if err != nil {
				return 0, err
			}

			b.SideTable[offset+0] = nameId
			b.SideTable[offset+1] = bodyId
			offset += 2
		}

		for i := range node.Asserts {
			assertId, err := b.visit(node.Asserts[i])
			if err != nil {
				return 0, err
			}
			b.SideTable[offset] = assertId
			offset++
		}

		for i := range node.Fields {
			keyId, err := b.visit(node.Fields[i].Name)
			if err != nil {
				return 0, err
			}
			bodyId, err := b.visit(node.Fields[i].Body)
			if err != nil {
				return 0, err
			}

			meta := createFieldMeta(node.Fields[i].Hide, node.Fields[i].PlusSuper)

			b.SideTable[offset+0] = keyId
			b.SideTable[offset+1] = bodyId
			b.SideTable[offset+2] = meta
			offset += 3
		}

		return b.emit(Node{Type: NodeTypeObject, A: uint32(startIdx), B: uint32(numFields), Flags: flags}, loc), nil

	case *ast.Array:

		arrSize := len(node.Elements)

		startIdx := len(b.SideTable)

		b.SideTable = slices.Grow(b.SideTable, arrSize)[:startIdx+arrSize]

		for i := range node.Elements {
			offset := startIdx + i

			id, err := b.visit(node.Elements[i].Expr)
			if err != nil {
				return 0, err
			}
			b.SideTable[offset] = id
		}

		return b.emit(Node{Type: NodeTypeArray, A: uint32(startIdx), B: uint32(arrSize)}, loc), nil

	case *ast.Local:

		binds := len(node.Binds)
		slots := binds * 2

		startIdx := len(b.SideTable)

		b.SideTable = slices.Grow(b.SideTable, slots)[:startIdx+slots]

		offset := startIdx
		for i := range node.Binds {

			keyId := b.Interner.Intern(string(node.Binds[i].Variable))
			b.SideTable[offset] = keyId

			bodyId, err := b.visit(node.Binds[i].Body)
			if err != nil {
				return 0, err
			}
			b.SideTable[offset+1] = bodyId

			offset += 2
		}

		bodyId, err := b.visit(node.Body)
		if err != nil {
			return 0, err
		}

		res := Node{
			Type: NodeTypeLocal,
			A:    bodyId,
			B:    uint32(startIdx),
			C:    uint32(binds),
		}

		return b.emit(res, loc), nil

	case *ast.Apply:

		posArgs := len(node.Arguments.Positional)
		namedArgs := len(node.Arguments.Named)

		paramCount := posArgs + namedArgs
		paramSlots := paramCount * 2

		startIdx := len(b.SideTable)

		b.SideTable = slices.Grow(b.SideTable, paramSlots)[:startIdx+paramSlots]

		offset := startIdx
		for i := range node.Arguments.Positional {
			// TODO: think abt maybe not wasting space here, most args are probably positional
			b.SideTable[offset] = 0

			id, err := b.visit(node.Arguments.Positional[i].Expr)
			if err != nil {
				return 0, err
			}
			b.SideTable[offset+1] = id

			offset += 2
		}

		for i := range node.Arguments.Named {

			paramKeyId := b.Interner.Intern(string(node.Arguments.Named[i].Name))
			b.SideTable[offset] = paramKeyId

			id, err := b.visit(node.Arguments.Named[i].Arg)
			if err != nil {
				return 0, err
			}
			b.SideTable[offset+1] = id

			offset += 2
		}

		funcId, err := b.visit(node.Target)
		if err != nil {
			return 0, err
		}

		return b.emit(Node{Type: NodeTypeApply, A: funcId, B: uint32(startIdx), C: uint32(paramCount)}, loc), nil

	case *ast.Index:

		index, err := b.visit(node.Index)
		if err != nil {
			return 0, err
		}

		target, err := b.visit(node.Target)
		if err != nil {
			return 0, err
		}

		return b.emit(Node{Type: NodeTypeIndex, A: index, B: target}, loc), nil

	case *ast.Var:
		keyId := b.Interner.Intern(string(node.Id))
		return b.emit(Node{Type: NodeTypeVar, A: keyId}, loc), nil

	case *ast.Function:

		paramCount := len(node.Parameters)
		paramSlots := paramCount * 2

		startIdx := len(b.SideTable)

		b.SideTable = slices.Grow(b.SideTable, paramSlots)[:startIdx+paramSlots]

		offset := startIdx
		for i := range paramCount {

			paramKeyId := b.Interner.Intern(string(node.Parameters[i].Name))
			b.SideTable[offset] = paramKeyId

			var defaultArgId uint32
			if node.Parameters[i].DefaultArg != nil {
				id, err := b.visit(node.Parameters[i].DefaultArg)
				if err != nil {
					return 0, err
				}
				defaultArgId = id
			}
			b.SideTable[offset+1] = defaultArgId

			offset += 2
		}

		body, err := b.visit(node.Body)
		if err != nil {
			return 0, err
		}

		return b.emit(Node{Type: NodeTypeFunction, A: body, B: uint32(startIdx), C: uint32(paramCount)}, loc), nil

	case *ast.Conditional:

		cond, err := b.visit(node.Cond)
		if err != nil {
			return 0, err
		}

		bTrue, err := b.visit(node.BranchTrue)
		if err != nil {
			return 0, err
		}

		bFalse, err := b.visit(node.BranchFalse)
		if err != nil {
			return 0, err
		}

		return b.emit(Node{Type: NodeTypeConditional, A: cond, B: bTrue, C: bFalse}, loc), nil

	case *ast.Binary:

		left, err := b.visit(node.Left)
		if err != nil {
			return 0, err
		}

		right, err := b.visit(node.Right)
		if err != nil {
			return 0, err
		}

		var t NodeType

		switch node.Op {
		default:
			return 0, fmt.Errorf("unhandled binary operation: %s", node.Op.String())
		case ast.BopMult:
			t = NodeTypeBinaryMult
		case ast.BopDiv:
			t = NodeTypeBinaryDiv
		// case ast.BopPercent:
		// 	t = NodeTypeBinaryPercent
		case ast.BopPlus:
			t = NodeTypeBinaryPlus
		case ast.BopMinus:
			t = NodeTypeBinaryMinus
		case ast.BopShiftL:
			t = NodeTypeBinaryShiftL
		case ast.BopShiftR:
			t = NodeTypeBinaryShiftR
		case ast.BopGreater:
			t = NodeTypeBinaryGreater
		case ast.BopGreaterEq:
			t = NodeTypeBinaryGreaterEq
		case ast.BopLess:
			t = NodeTypeBinaryLess
		case ast.BopLessEq:
			t = NodeTypeBinaryLessEq
		case ast.BopManifestEqual:
			t = NodeTypeBinaryEqual
		case ast.BopManifestUnequal:
			t = NodeTypeBinaryUnequal
		case ast.BopBitwiseAnd:
			t = NodeTypeBinaryBitwiseAnd
		case ast.BopBitwiseXor:
			t = NodeTypeBinaryBitwiseXor
		case ast.BopBitwiseOr:
			t = NodeTypeBinaryBitwiseOr
		case ast.BopAnd:
			t = NodeTypeBinaryAnd
		case ast.BopOr:
			t = NodeTypeBinaryOr
		}

		return b.emit(Node{Type: t, A: left, B: right}, loc), nil

	case *ast.Unary:

		switch node.Op {
		default:
			return 0, fmt.Errorf("unhandled unary type: %s", node.Op.String())
		case ast.UopNot:
			if n, ok := node.Expr.(*ast.LiteralBoolean); ok {
				if n.Value {
					// note: false here is correct since we flip the bool
					return b.emit(Node{Type: NodeTypeFalse}, loc), nil
				}
				return b.emit(Node{Type: NodeTypeTrue}, loc), nil
			}

			id, err := b.visit(node.Expr)
			if err != nil {
				return 0, err
			}
			return b.emit(Node{Type: NodeTypeUnaryNot, A: id}, loc), nil

		case ast.UopMinus:

			if n, ok := node.Expr.(*ast.LiteralNumber); ok {
				num, err := strconv.ParseFloat(n.OriginalString, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to parse float val (%s), err: %w", n.OriginalString, err)
				}
				bits := math.Float64bits(-num)

				res := Node{
					Type: NodeTypeNumber,
					A:    uint32(bits >> 32),
					B:    uint32(bits),
				}
				return b.emit(res, loc), nil
			}

			id, err := b.visit(node.Expr)
			if err != nil {
				return 0, err
			}
			return b.emit(Node{Type: NodeTypeUnaryMinus, A: id}, loc), nil

		case ast.UopBitwiseNot:

			if n, ok := node.Expr.(*ast.LiteralNumber); ok {
				num, err := strconv.ParseFloat(n.OriginalString, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to parse float val (%s), err: %w", n.OriginalString, err)
				}
				bits := math.Float64bits(float64(^int64(num)))

				res := Node{
					Type: NodeTypeNumber,
					A:    uint32(bits >> 32),
					B:    uint32(bits),
				}
				return b.emit(res, loc), nil
			}

			id, err := b.visit(node.Expr)
			if err != nil {
				return 0, err
			}
			return b.emit(Node{Type: NodeTypeUnaryBitwiseNot, A: id}, loc), nil
		}

	case *ast.Import:
		id := b.Interner.Intern(node.File.Value)
		fileLoc := b.Interner.Intern(node.NodeBase.LocRange.FileName)
		return b.emit(Node{Type: NodeTypeImport, A: id, B: fileLoc}, loc), nil

	case *ast.ImportStr:
		id := b.Interner.Intern(node.File.Value)
		return b.emit(Node{Type: NodeTypeImportStr, A: id}, loc), nil

	case *ast.Self:
		return b.emit(Node{Type: NodeTypeSelf}, loc), nil

	case *ast.SuperIndex:
		// TODO: can only be LiteralString?
		id, err := b.visit(node.Index)
		if err != nil {
			return 0, err
		}
		return b.emit(Node{Type: NodeTypeSuperIndex, A: id}, loc), nil

	case *ast.InSuper:
		// TODO: can only be LiteralString?
		id, err := b.visit(node.Index)
		if err != nil {
			return 0, err
		}
		return b.emit(Node{Type: NodeTypeInSuper, A: id}, loc), nil

	case *ast.Error:
		id, err := b.visit(node.Expr)
		if err != nil {
			return 0, err
		}
		return b.emit(Node{Type: NodeTypeError, A: id}, loc), nil

	}

}

const (
	MaskVisibility = 0x03 // Binary 00000011
	FlagPlusSuper  = 0x04 // Binary 00000100
)

func createFieldMeta(visibility ast.ObjectFieldHide, plusSuper bool) uint32 {
	m := uint32(visibility) & MaskVisibility
	if plusSuper {
		m |= FlagPlusSuper
	}
	return m
}
