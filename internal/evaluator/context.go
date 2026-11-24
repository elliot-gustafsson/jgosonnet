package evaluator

import (
	"io"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type ContextState struct {
	MaxStack    uint32
	Interner    *interner.Interner
	Allocator   *arena.Allocator
	Environment *Environment
}

type Context struct {
	State *ContextState

	Self Value // self

	SuperOffset uint32

	Depth uint32
}

// NewScope allocates a pointer-free lexical scope frame in the arena:
//   [0..8):      parentPtr (uintptr)
//   [8..8+8*N):  slots[0..N-1] (Value)
// Callers are responsible for initializing all slots (or calling clear(slots[pos:])
// if only a prefix is immediately written).
func (c Context) NewScope(parentPtr uintptr, length int) (slots []Value, scopePtr uintptr) {
	totalSize := uintptr(8) + uintptr(length)*valueSize
	ptr := arena.AlignedAlloc(c.State.Allocator, totalSize, 8)
	*(*uintptr)(ptr) = parentPtr
	if length > 0 {
		slots = unsafe.Slice((*Value)(unsafe.Add(ptr, 8)), length)
	}
	return slots, uintptr(ptr)
}

func getScopeSlotPtr(scopePtr uintptr, upDepth, slotIdx uint32) *Value {
	p := scopePtr
	switch upDepth {
	case 0:
	case 1:
		p = *(*uintptr)(resolveUintptr(p))
	case 2:
		p = *(*uintptr)(resolveUintptr(p))
		p = *(*uintptr)(resolveUintptr(p))
	default:
		for range upDepth {
			p = *(*uintptr)(resolveUintptr(p))
		}
	}
	return (*Value)(unsafe.Add(resolveUintptr(p), 8+uintptr(slotIdx)*valueSize))
}

func getScopeSlot(scopePtr uintptr, upDepth, slotIdx uint32) Value {
	return *getScopeSlotPtr(scopePtr, upDepth, slotIdx)
}

type Environment struct {
	TraceOut        io.Writer
	Importer        *Importer
	ExtVars         map[string]string
	ExtCodes        map[string]string
	NativeFunctions map[string]Function
}
