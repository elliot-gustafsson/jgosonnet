package evaluator

import (
	"io"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/alloc"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type ContextState struct {
	Interner    *interner.Interner
	Allocator   *alloc.Allocator
	Environment *Environment
}

type Context struct {
	State *ContextState

	Self Value // self

	SuperOffset uint32
}

type Scope struct {
	Bindings  []NamedValue
	ParentPtr uintptr
}

func (c Context) NewScope(parentPtr uintptr, length int) (*Scope, uintptr) {
	s := c.State.Allocator.Create[Scope]()
	alloc.Memclr(s)

	s.ParentPtr = parentPtr
	s.Bindings = c.State.Allocator.Alloc[NamedValue](length)
	clear(s.Bindings)

	return s, uintptr(unsafe.Pointer(s))
}

func (c Context) GetScopeBind(scopePtr uintptr, key uint32) (val Value, found bool) {
	currPtr := scopePtr

	for {
		scope := (*Scope)(resolveUintptr(currPtr))
		bindings := scope.Bindings

		for i := len(bindings) - 1; i >= 0; i-- {
			if bindings[i].Key == key {
				val, found = bindings[i].Value, true
				return
			}
		}

		if currPtr == 0 {
			break
		}

		if scope.ParentPtr == currPtr {
			break
		}

		currPtr = scope.ParentPtr
	}
	return
}

type Environment struct {
	TraceOut        io.Writer
	Importer        *Importer
	ExtVars         map[string]string
	ExtCodes        map[string]string
	NativeFunctions map[string]Function
}
