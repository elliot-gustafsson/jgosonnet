package ast

import (
	"math"
	"unsafe"
)

type Tag uint32

const (
	TagNone Tag = iota // tag 1B

	// Literals
	TagNull   // tag 1B
	TagTrue   // tag 1B
	TagFalse  // tag 1B
	TagNumber // tag 1B, val 8B
	TagString // tag 1B, id 4B

	// Scoping & Identifiers
	TagVar        // tag 1B, id 4B
	TagSelf       // tag 1B
	TagDollar     // tag 1B
	TagSuperIndex // tag 1B, index 4B
	TagInSuper    // tag 1B, index 4B

	TagLocal1 // tag 1B, key 4B, val 4B, body 4B
	TagLocal  // tag 1B, bind_count 4B, binds N*(key 4B, val 4B), body 4B

	// Indexing
	TagIndex // tag 1B, target 4B, index 4B

	// Unary Operators
	TagUnaryNot        // tag 1B, child 4B
	TagUnaryMinus      // tag 1B, child 4B
	TagUnaryPlus       // tag 1B, child 4B
	TagUnaryBitwiseNot // tag 1B, child 4B

	// Binary Math Operators
	TagBinaryAdd // tag 1B, left 4B, right 4B
	TagBinarySub // tag 1B, left 4B, right 4B
	TagBinaryMul // tag 1B, left 4B, right 4B
	TagBinaryDiv // tag 1B, left 4B, right 4B
	TagBinaryMod // tag 1B, left 4B, right 4B

	// Binary Bitwise Operators
	TagBinaryBitAnd // tag 1B, left 4B, right 4B
	TagBinaryBitOr  // tag 1B, left 4B, right 4B
	TagBinaryBitXor // tag 1B, left 4B, right 4B
	TagBinaryShiftL // tag 1B, left 4B, right 4B
	TagBinaryShiftR // tag 1B, left 4B, right 4B

	// Binary Comparison Operators
	TagBinaryEq  // tag 1B, left 4B, right 4B
	TagBinaryNeq // tag 1B, left 4B, right 4B
	TagBinaryLt  // tag 1B, left 4B, right 4B
	TagBinaryLte // tag 1B, left 4B, right 4B
	TagBinaryGt  // tag 1B, left 4B, right 4B
	TagBinaryGte // tag 1B, left 4B, right 4B
	TagBinaryIn  // tag 1B, left 4B, right 4B

	// Short-Circuiting Logical Operators
	TagBinaryAnd // tag 1B, left 4B, right 4B
	TagBinaryOr  // tag 1B, left 4B, right 4B

	// Control Flow & Errors
	TagConditional // tag 1B, cond 4B, true 4B, false 4B
	TagError       // tag 1B, msg 4B
	TagAssert      // tag 1B, cond 4B, msg 4B, rest 4B

	// Functions & Calls
	TagFunction // tag 1B, param_count 2B, body 4B, params N*(name 4B, default 4B)
	TagApply0   // tag 1B, func 4B
	TagApply1   // tag 1B, func 4B, arg0 4B
	TagApply2   // tag 1B, func 4B, arg0 4B, arg1 4B
	TagApply3   // tag 1B, func 4B, arg0 4B, arg1 4B, arg2 4B
	TagApply4   // tag 1B, func 4B, arg0 4B, arg1 4B, arg2 4B, arg3 4B
	TagApply    // tag 1B, pos_count 2B, named_count 2B, func 4B, pos_args N*4B, named_args M*(key 4B, val 4B)

	// Arrays
	TagArray0 // tag 1B
	// TagArray1 // tag 1B, elem0 4B
	// TagArray2 // tag 1B, elem0 4B, elem1 4B
	TagArray // tag 1B, count 2B, elems N*4B

	TagArrayComp // tag 1B, count 2B, elems N*4B

	// Objects
	TagObject0      // tag 1B                                                                                    -> {} (empty)
	TagObjectSimple // tag 1B, field_count 2B, fields N*(key 4B, val 4B, meta 1B)                                -> purely static fields (<= 65k)
	TagObject       // tag 1B, field_count 4B, local_count 4B, assert_count 4B, fields..., locals..., asserts... -> general object

	// Imports
	TagImport    // tag 1B, file_id 4B
	TagImportStr // tag 1B, file_id 4B
	TagImportBin // tag 1B, file_id 4B
)

type AST []uint32

// func (a *AST) alloc(n int) (uint32, unsafe.Pointer) {
// 	offset := len(*a)
// 	newLen := offset + n

// 	if newLen <= cap(*a) {
// 		*a = (*a)[:newLen]
// 		return uint32(offset), unsafe.Pointer(&(*a)[offset])
// 	}

// 	*a = append(*a, make([]uint32, n)...)
// 	return uint32(offset), unsafe.Pointer(&(*a)[offset])
// }

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

// func (a *AST) emitNumber(num float64) uint32 {
// 	offset, p := a.alloc(9)
// 	*(*uint8)(p) = uint8(TagNumber)
// 	*(*float64)(unsafe.Add(p, 1)) = num
// 	return offset
// }

func (a *AST) emitNumber(num float64) uint32 {
	idx := uint32(len(*a))
	b := math.Float64bits(num)
	*a = append(*a, uint32(TagNumber), uint32(b), uint32(b>>32))
	return idx

	// idx := uint32(len(*a))
	// v := uint64(MakeNumber(num)) // NaN-box once at parse time
	// *a = append(*a, uint32(TagNumber), uint32(v), uint32(v>>32))

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

func (a *AST) emitFunction(args []uint32, body uint32) uint32 {
	idx := uint32(len(*a))
	count := uint32(len(args))
	*a = append(*a, uint32(TagFunction), count)
	*a = append(*a, args...)
	*a = append(*a, body)
	return idx
}

func (a *AST) emitFunction1(arg, d_arg uint32, body uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagFunction), 1, arg, d_arg, body)
	return idx
}

func (a *AST) emitFunction2(arg1, d_arg1, arg2, d_arg2 uint32, body uint32) uint32 {
	idx := uint32(len(*a))
	*a = append(*a, uint32(TagFunction), 2, arg1, d_arg1, arg2, d_arg2, body)
	return idx
}

func (a AST) ReadTag(idx uint32) Tag         { return ReadTag(a.Ptr(idx), 0) }
func (a AST) ReadUint32(idx uint32) uint32   { return ReadU32(a.Ptr(idx), 0) }
func (a AST) ReadFloat64(idx uint32) float64 { return ReadF64(a.Ptr(idx), 0) }

func ReadTag(p unsafe.Pointer, wordOffset uint32) Tag {
	return Tag(*(*uint32)(unsafe.Add(p, wordOffset*4)))
}

func ReadU32(p unsafe.Pointer, wordOffset uint32) uint32 {
	return *(*uint32)(unsafe.Add(p, wordOffset*4))
}

func ReadF64(p unsafe.Pointer, wordOffset uint32) float64 {
	return *(*float64)(unsafe.Add(p, wordOffset*4))
}

// func ReadValue(p unsafe.Pointer, wordOffset uintptr) evaluator.Value {
// 	return *(*evaluator.Value)(unsafe.Add(p, wordOffset*4))
// }
