package ast

import (
	"fmt"
	"math"
	"sort"
	"unsafe"
)

type Tag uint32

// NoNode is the sentinel value indicating an absent optional AST node offset
// (since offset 0 is a valid node offset).
const NoNode uint32 = math.MaxUint32

type FieldHide uint8

const (
	FieldHidden  FieldHide = 0 // ::
	FieldInherit FieldHide = 1 // :
	FieldVisible FieldHide = 2 // :::
)

const (
	FieldMaskVisibility uint8 = 0x03
	FieldFlagPlusSuper  uint8 = 0x04
	FieldFlagTombstone  uint8 = 0x08
	FieldFlagStaticKey  uint8 = 0x10 // AST-only: field key word is an interned string ID directly
)

const (
	ApplyFlagTailStrict uint32 = 1 << 0
)

const (
	TagNone Tag = iota

	// Literals
	TagNull   // [tag]
	TagTrue   // [tag]
	TagFalse  // [tag]
	TagNumber // [tag, lo32, hi32]
	TagString // [tag, strId]

	// Scoping & Identifiers
	TagVar             // [tag, varId]
	TagSelf            // [tag]
	TagDollar          // [tag] (desugared to TagVar("$"))
	TagSuperIndex      // [tag, indexOff]
	TagSuperIndexConst // [tag, keyId]
	TagInSuper         // [tag, indexOff]

	TagLocal1 // [tag, keyId, valOff, bodyOff]
	TagLocal  // [tag, bindCount, binds N*(keyId, valOff), bodyOff]

	// Indexing
	TagIndex      // [tag, targetOff, indexOff]
	TagIndexConst // [tag, targetOff, keyId]

	// Unary Operators
	TagUnaryNot        // [tag, childOff]
	TagUnaryMinus      // [tag, childOff]
	TagUnaryPlus       // [tag, childOff]
	TagUnaryBitwiseNot // [tag, childOff]

	// Binary Math Operators
	TagBinaryAdd // [tag, leftOff, rightOff]
	TagBinarySub // [tag, leftOff, rightOff]
	TagBinaryMul // [tag, leftOff, rightOff]
	TagBinaryDiv // [tag, leftOff, rightOff]
	TagBinaryMod // [tag, leftOff, rightOff]

	// Binary Bitwise Operators
	TagBinaryBitAnd // [tag, leftOff, rightOff]
	TagBinaryBitOr  // [tag, leftOff, rightOff]
	TagBinaryBitXor // [tag, leftOff, rightOff]
	TagBinaryShiftL // [tag, leftOff, rightOff]
	TagBinaryShiftR // [tag, leftOff, rightOff]

	// Binary Comparison Operators
	TagBinaryEq  // [tag, leftOff, rightOff]
	TagBinaryNeq // [tag, leftOff, rightOff]
	TagBinaryLt  // [tag, leftOff, rightOff]
	TagBinaryLte // [tag, leftOff, rightOff]
	TagBinaryGt  // [tag, leftOff, rightOff]
	TagBinaryGte // [tag, leftOff, rightOff]
	TagBinaryIn  // [tag, leftOff, rightOff]

	// Short-Circuiting Logical Operators
	TagBinaryAnd // [tag, leftOff, rightOff]
	TagBinaryOr  // [tag, leftOff, rightOff]

	// Control Flow & Errors
	TagConditional // [tag, condOff, trueOff, falseOff]
	TagError       // [tag, msgOff]
	TagAssert      // [tag, condOff, msgOff (or NoNode), restOff]

	// Functions & Calls
	TagFunction // [tag, paramCount, bodyOff, params N*(nameId, defaultOff or NoNode)]
	TagApply0   // [tag, funcOff]
	TagApply1   // [tag, funcOff, arg0]
	TagApply2   // [tag, funcOff, arg0, arg1]
	TagApply3   // [tag, funcOff, arg0, arg1, arg2]
	TagApply4   // [tag, funcOff, arg0, arg1, arg2, arg3]
	TagApply    // [tag, posCount, namedCount, flags, funcOff, posArgs..., namedArgs M*(keyId, valOff)]

	// Arrays
	TagArray0    // [tag]
	TagArray     // [tag, count, elems...]
	TagArrayComp // [tag, bodyOff, clauseCount, clauses...(varId, iterOff, condCount, conds...)]

	// Objects
	TagObject0      // [tag]
	TagObjectSimple // [tag, fieldCount, fields N*(keyId, bodyOff, meta)]
	TagObject       // [tag, fieldCount, localCount, assertCount, fields N*(keyOrOff, bodyOff, meta), locals L*(keyId, bodyOff), asserts A*(assertOff)]
	TagObjectComp   // [tag, keyOff, bodyOff, meta, clauseCount, clauses...(varId, iterOff, condCount, conds...)]

	// Imports
	TagImport    // [tag, fileId]
	TagImportStr // [tag, fileId]
	TagImportBin // [tag, fileId]
)

type AST []uint32

func (a AST) Base() unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(a))
}

func (a AST) Ptr(idx uint32) unsafe.Pointer {
	return unsafe.Add(a.Base(), uintptr(idx)*4)
}

func (a *AST) offset() uint32 {
	return uint32(len(*a))
}

func (a *AST) emitTag(tag Tag) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(tag))
	return idx
}

func (a *AST) emitTagU32(tag Tag, v uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(tag), v)
	return idx
}

func (a *AST) emitTag2U32(tag Tag, v0, v1 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(tag), v0, v1)
	return idx
}

func (a *AST) emitTag3U32(tag Tag, v0, v1, v2 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(tag), v0, v1, v2)
	return idx
}

func (a *AST) emitNumber(num float64) uint32 {
	idx := uint32(len(*a))
	b := math.Float64bits(num)
	*a = append(*a, uint32(TagNumber), uint32(b), uint32(b>>32))
	return idx
}

func (a *AST) emitLocal(key, val, body uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagLocal1), key, val, body)
	return idx
}

func (a *AST) emitLocals(binds []uint32, body uint32) uint32 {
	idx := uint32(len(*a))
	count := uint32(len(binds) / 2)
	*a = append(*a, uint32(TagLocal), count)
	*a = append(*a, binds...)
	*a = append(*a, body)
	return idx
}

func (a *AST) emitArray(elements []uint32) uint32 {
	idx := uint32(len(*a))
	count := uint32(len(elements))
	*a = append(*a, uint32(TagArray), count)
	*a = append(*a, elements...)
	return idx
}

func (a *AST) emitFunction(params []uint32, body uint32) uint32 {
	idx := uint32(len(*a))
	count := uint32(len(params) / 2)
	*a = append(*a, uint32(TagFunction), count, body)
	*a = append(*a, params...)
	return idx
}

func (a *AST) emitApply0(fn uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagApply0), fn)
	return idx
}

func (a *AST) emitApply1(fn, a0 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagApply1), fn, a0)
	return idx
}

func (a *AST) emitApply2(fn, a0, a1 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagApply2), fn, a0, a1)
	return idx
}

func (a *AST) emitApply3(fn, a0, a1, a2 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagApply3), fn, a0, a1, a2)
	return idx
}

func (a *AST) emitApply4(fn, a0, a1, a2, a3 uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagApply4), fn, a0, a1, a2, a3)
	return idx
}

func (a *AST) emitApply(fn uint32, posArgs, namedArgs []uint32, tailStrict bool) uint32 {
	if !tailStrict && len(namedArgs) == 0 {
		switch len(posArgs) {
		case 0:
			return a.emitApply0(fn)
		case 1:
			return a.emitApply1(fn, posArgs[0])
		case 2:
			return a.emitApply2(fn, posArgs[0], posArgs[1])
		case 3:
			return a.emitApply3(fn, posArgs[0], posArgs[1], posArgs[2])
		case 4:
			return a.emitApply4(fn, posArgs[0], posArgs[1], posArgs[2], posArgs[3])
		}
	}
	idx := uint32(len(*a))
	var flags uint32
	if tailStrict {
		flags |= ApplyFlagTailStrict
	}
	*a = append(*a, uint32(TagApply), uint32(len(posArgs)), uint32(len(namedArgs)/2), flags, fn)
	*a = append(*a, posArgs...)
	*a = append(*a, namedArgs...)
	return idx
}

func (a *AST) emitObjectSimple(fields []uint32) uint32 {
	idx := uint32(len(*a))
	count := uint32(len(fields) / 3)
	*a = append(*a, uint32(TagObjectSimple), count)
	*a = append(*a, fields...)
	return idx
}

func (a *AST) emitObject(fields, locals, asserts []uint32) uint32 {
	idx := uint32(len(*a))
	fieldCount := uint32(len(fields) / 3)
	localCount := uint32(len(locals) / 2)
	assertCount := uint32(len(asserts))
	*a = append(*a, uint32(TagObject), fieldCount, localCount, assertCount)
	*a = append(*a, fields...)
	*a = append(*a, locals...)
	*a = append(*a, asserts...)
	return idx
}

func (a *AST) emitArrayComp(body uint32, clauseCount uint32, clauses []uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagArrayComp), body, clauseCount)
	*a = append(*a, clauses...)
	return idx
}

func (a *AST) emitObjectComp(keyOff, bodyOff, meta, clauseCount uint32, clauses []uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagObjectComp), keyOff, bodyOff, meta, clauseCount)
	*a = append(*a, clauses...)
	return idx
}

func (a AST) ReadTag(idx uint32) Tag         { return ReadTag(a.Ptr(idx), 0) }
func (a AST) ReadUint32(idx uint32) uint32   { return ReadU32(a.Ptr(idx), 0) }
func (a AST) ReadFloat64(idx uint32) float64 { return ReadF64(a.Ptr(idx), 0) }

func ReadTag(p unsafe.Pointer, wordOffset uint32) Tag {
	return Tag(*(*uint32)(unsafe.Add(p, uintptr(wordOffset)*4)))
}

func ReadU32(p unsafe.Pointer, wordOffset uint32) uint32 {
	return *(*uint32)(unsafe.Add(p, uintptr(wordOffset)*4))
}

func ReadF64(p unsafe.Pointer, wordOffset uint32) float64 {
	return *(*float64)(unsafe.Add(p, uintptr(wordOffset)*4))
}

// --- Cold source metadata (only touched on errors, std.trace, and imports) ---

type Location struct {
	Line   int
	Column int
}

func (l Location) String() string {
	return fmt.Sprintf("%d:%d", l.Line, l.Column)
}

type LocationRange struct {
	FileName string
	Begin    Location
	End      Location
}

func (lr LocationRange) IsSet() bool {
	return lr.Begin.Line != 0
}

func (lr LocationRange) String() string {
	var filePrefix string
	if len(lr.FileName) > 0 {
		filePrefix = lr.FileName + ":"
	}
	if !lr.IsSet() {
		return lr.FileName
	}
	if lr.Begin.Line == lr.End.Line {
		if lr.Begin.Column == lr.End.Column {
			return fmt.Sprintf("%s%d:%d", filePrefix, lr.Begin.Line, lr.Begin.Column)
		}
		return fmt.Sprintf("%s%d:%d-%d", filePrefix, lr.Begin.Line, lr.Begin.Column, lr.End.Column)
	}
	return fmt.Sprintf("%s(%d:%d)-(%d:%d)", filePrefix, lr.Begin.Line, lr.Begin.Column, lr.End.Line, lr.End.Column)
}

type StaticError struct {
	Loc LocationRange
	Msg string
}

func (e StaticError) Error() string {
	if e.Loc.IsSet() {
		return fmt.Sprintf("%s %s", e.Loc.String(), e.Msg)
	}
	return e.Msg
}

type LocEntry struct {
	NodeOffset uint32
	StartByte  uint32
	EndByte    uint32
	ContextIdx uint32 // 0 = no context (""), else 1-based index into Program.Contexts
}

type Program struct {
	Filename string
	Source   string
	Nodes    AST
	RootIdx  uint32
	Locs     []LocEntry // sorted by NodeOffset
	Contexts []string   // 1-based via ContextIdx
}

func (prog *Program) Base() unsafe.Pointer {
	return prog.Nodes.Base()
}

// NodeInfo returns the cold LocationRange and Context name for a node offset.
func (prog *Program) NodeInfo(offset uint32) (LocationRange, string) {
	if prog == nil || len(prog.Locs) == 0 {
		return LocationRange{}, ""
	}
	i := sort.Search(len(prog.Locs), func(i int) bool {
		return prog.Locs[i].NodeOffset >= offset
	})
	if i >= len(prog.Locs) || prog.Locs[i].NodeOffset != offset {
		return LocationRange{}, ""
	}
	entry := prog.Locs[i]
	var ctx string
	if entry.ContextIdx > 0 && int(entry.ContextIdx) <= len(prog.Contexts) {
		ctx = prog.Contexts[entry.ContextIdx-1]
	}
	return prog.ByteRangeToLoc(entry.StartByte, entry.EndByte), ctx
}

func (prog *Program) ByteRangeToLoc(startByte, endByte uint32) LocationRange {
	return byteRangeToLoc(prog.Filename, prog.Source, startByte, endByte)
}

func byteRangeToLoc(filename, src string, startByte, endByte uint32) LocationRange {
	begin := byteOffsetToLoc(src, startByte)
	end := byteOffsetToLoc(src, endByte)
	return LocationRange{
		FileName: filename,
		Begin:    begin,
		End:      end,
	}
}

func byteOffsetToLoc(src string, byteOff uint32) Location {
	off := int(byteOff)
	if off > len(src) {
		off = len(src)
	}
	line := 1
	lineStart := 0
	for i := 0; i < off; i++ {
		if src[i] == '\n' {
			line++
			lineStart = i + 1
		}
	}
	return Location{
		Line:   line,
		Column: off - lineStart + 1,
	}
}
