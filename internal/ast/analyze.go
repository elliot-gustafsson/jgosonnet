package ast

import (
	"fmt"
	"slices"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type analyzer struct {
	prog     *Program
	interner *interner.Interner
	vars     []uint32
}

// analyzeProgram statically checks variable bindings and object-only constructs (self/super).
func analyzeProgram(prog *Program, in *interner.Interner) error {
	a := analyzer{
		prog:     prog,
		interner: in,
		vars:     make([]uint32, 0, 32),
	}
	a.vars = append(a.vars, in.Intern("std"), in.Intern("$std"))
	return a.visit(prog.RootIdx, false)
}

func (a *analyzer) hasVar(id uint32) bool {
	return slices.Contains(a.vars, id)
}

func (a *analyzer) visit(off uint32, inObject bool) error {
	if off == NoNode {
		return nil
	}
	nodes := a.prog.Nodes
	tag := Tag(nodes[off])

	switch tag {
	case TagNone, TagNull, TagTrue, TagFalse, TagNumber, TagString, TagDollar,
		TagArray0, TagObject0, TagImport, TagImportStr, TagImportBin:
		return nil

	case TagSelf:
		if !inObject {
			loc, _ := a.prog.NodeInfo(off)
			return StaticError{Loc: loc, Msg: "Can't use self outside of an object."}
		}
		return nil

	case TagSuperIndexConst:
		if !inObject {
			loc, _ := a.prog.NodeInfo(off)
			return StaticError{Loc: loc, Msg: "Can't use super outside of an object."}
		}
		return nil

	case TagSuperIndex, TagInSuper:
		if !inObject {
			loc, _ := a.prog.NodeInfo(off)
			return StaticError{Loc: loc, Msg: "Can't use super outside of an object."}
		}
		return a.visit(nodes[off+1], inObject)

	case TagVar:
		id := nodes[off+1]
		if !a.hasVar(id) {
			loc, _ := a.prog.NodeInfo(off)
			return StaticError{
				Loc: loc,
				Msg: fmt.Sprintf("Unknown variable: %s", a.interner.Get(id)),
			}
		}
		return nil

	case TagError, TagUnaryNot, TagUnaryMinus, TagUnaryPlus, TagUnaryBitwiseNot, TagIndexConst:
		return a.visit(nodes[off+1], inObject)

	case TagIndex,
		TagBinaryAdd, TagBinarySub, TagBinaryMul, TagBinaryDiv, TagBinaryMod,
		TagBinaryBitAnd, TagBinaryBitOr, TagBinaryBitXor, TagBinaryShiftL, TagBinaryShiftR,
		TagBinaryEq, TagBinaryNeq, TagBinaryLt, TagBinaryLte, TagBinaryGt, TagBinaryGte, TagBinaryIn,
		TagBinaryAnd, TagBinaryOr:
		if err := a.visit(nodes[off+1], inObject); err != nil {
			return err
		}
		return a.visit(nodes[off+2], inObject)

	case TagConditional, TagAssert:
		if err := a.visit(nodes[off+1], inObject); err != nil {
			return err
		}
		if err := a.visit(nodes[off+2], inObject); err != nil {
			return err
		}
		return a.visit(nodes[off+3], inObject)

	case TagLocal1:
		mark := len(a.vars)
		a.vars = append(a.vars, nodes[off+1])
		if err := a.visit(nodes[off+2], inObject); err != nil {
			return err
		}
		if err := a.visit(nodes[off+3], inObject); err != nil {
			return err
		}
		a.vars = a.vars[:mark]
		return nil

	case TagLocal:
		mark := len(a.vars)
		bindCount := nodes[off+1]
		for i := range bindCount {
			a.vars = append(a.vars, nodes[off+2+i*2])
		}
		for i := range bindCount {
			if err := a.visit(nodes[off+3+i*2], inObject); err != nil {
				return err
			}
		}
		bodyOff := nodes[off+2+bindCount*2]
		if err := a.visit(bodyOff, inObject); err != nil {
			return err
		}
		a.vars = a.vars[:mark]
		return nil

	case TagFunction:
		mark := len(a.vars)
		paramCount := nodes[off+1]
		bodyOff := nodes[off+2]
		for i := range paramCount {
			a.vars = append(a.vars, nodes[off+3+i*2])
		}
		for i := range paramCount {
			if err := a.visit(nodes[off+4+i*2], inObject); err != nil {
				return err
			}
		}
		if err := a.visit(bodyOff, inObject); err != nil {
			return err
		}
		a.vars = a.vars[:mark]
		return nil

	case TagApply0:
		return a.visit(nodes[off+1], inObject)
	case TagApply1:
		if err := a.visit(nodes[off+1], inObject); err != nil {
			return err
		}
		return a.visit(nodes[off+2], inObject)
	case TagApply2:
		for i := uint32(1); i <= 3; i++ {
			if err := a.visit(nodes[off+i], inObject); err != nil {
				return err
			}
		}
		return nil
	case TagApply3:
		for i := uint32(1); i <= 4; i++ {
			if err := a.visit(nodes[off+i], inObject); err != nil {
				return err
			}
		}
		return nil
	case TagApply4:
		for i := uint32(1); i <= 5; i++ {
			if err := a.visit(nodes[off+i], inObject); err != nil {
				return err
			}
		}
		return nil
	case TagApply:
		posCount := nodes[off+1]
		namedCount := nodes[off+2]
		if err := a.visit(nodes[off+4], inObject); err != nil {
			return err
		}
		base := off + 5
		for i := range posCount {
			if err := a.visit(nodes[base+i], inObject); err != nil {
				return err
			}
		}
		namedBase := base + posCount
		for i := range namedCount {
			if err := a.visit(nodes[namedBase+i*2+1], inObject); err != nil {
				return err
			}
		}
		return nil

	case TagArray:
		count := nodes[off+1]
		for i := range count {
			if err := a.visit(nodes[off+2+i], inObject); err != nil {
				return err
			}
		}
		return nil

	case TagArrayComp:
		bodyOff := nodes[off+1]
		clauseCount := nodes[off+2]
		mark := len(a.vars)
		pos := off + 3
		for range clauseCount {
			varId := nodes[pos]
			iterOff := nodes[pos+1]
			condCount := nodes[pos+2]
			if err := a.visit(iterOff, inObject); err != nil {
				return err
			}
			a.vars = append(a.vars, varId)
			for j := range condCount {
				if err := a.visit(nodes[pos+3+j], inObject); err != nil {
					return err
				}
			}
			pos += 3 + condCount
		}
		if err := a.visit(bodyOff, inObject); err != nil {
			return err
		}
		a.vars = a.vars[:mark]
		return nil

	case TagObjectSimple:
		fieldCount := nodes[off+1]
		for i := range fieldCount {
			if err := a.visit(nodes[off+2+i*3+1], true); err != nil {
				return err
			}
		}
		return nil

	case TagObject:
		fieldCount := nodes[off+1]
		localCount := nodes[off+2]
		assertCount := nodes[off+3]
		fieldsBase := off + 4
		localsBase := fieldsBase + fieldCount*3
		assertsBase := localsBase + localCount*2

		// Object locals, field bodies, and asserts are inside the object
		mark := len(a.vars)
		for i := range localCount {
			a.vars = append(a.vars, nodes[localsBase+i*2])
		}
		for i := range localCount {
			if err := a.visit(nodes[localsBase+i*2+1], true); err != nil {
				return err
			}
		}
		for i := range fieldCount {
			if err := a.visit(nodes[fieldsBase+i*3+1], true); err != nil {
				return err
			}
		}
		for i := range assertCount {
			if err := a.visit(nodes[assertsBase+i], true); err != nil {
				return err
			}
		}
		a.vars = a.vars[:mark]

		// Computed field names are evaluated OUTSIDE the object
		for i := range fieldCount {
			meta := uint8(nodes[fieldsBase+i*3+2])
			if (meta & FieldFlagStaticKey) == 0 {
				if err := a.visit(nodes[fieldsBase+i*3], inObject); err != nil {
					return err
				}
			}
		}
		return nil

	case TagObjectComp:
		keyOff := nodes[off+1]
		bodyOff := nodes[off+2]
		clauseCount := nodes[off+4]
		mark := len(a.vars)
		pos := off + 5
		for range clauseCount {
			varId := nodes[pos]
			iterOff := nodes[pos+1]
			condCount := nodes[pos+2]
			if err := a.visit(iterOff, inObject); err != nil {
				return err
			}
			a.vars = append(a.vars, varId)
			for j := range condCount {
				if err := a.visit(nodes[pos+3+j], inObject); err != nil {
					return err
				}
			}
			pos += 3 + condCount
		}
		if err := a.visit(keyOff, inObject); err != nil {
			return err
		}
		if err := a.visit(bodyOff, true); err != nil {
			return err
		}
		a.vars = a.vars[:mark]
		return nil

	default:
		return fmt.Errorf("analyzer: unhandled tag %d", tag)
	}
}
