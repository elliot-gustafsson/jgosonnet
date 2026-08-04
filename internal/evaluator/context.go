package evaluator

import (
	"io"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

type ContextState struct {
	Interner    *interner.Interner
	Registry    *Registry
	Environment *Environment
}

type Context struct {
	State *ContextState

	Self Value // self

	SuperOffset uint32
}

type Registry struct {
	// Objects   *arena.Arena[Object]
	// Arrays    *arena.SliceArena[Value]
	// Strings   *arena.StringArena
	// Thunks    *arena.Arena[Thunk]
	// Functions *arena.Arena[Function]

	// Scopes *arena.Arena[Scope]

	// GoCallbackNodes *arena.Arena[GoCallbackNode]

	Allocator *arena.Allocator
}

type Scope struct {
	Bindings  []NamedValue
	ParentPtr uintptr
}

func NewRegistry() *Registry {
	return &Registry{
		// Functions: arena.NewArena[Function](),

		// GoCallbackNodes: arena.NewArena[GoCallbackNode](),

		Allocator: arena.NewAllocator(),
	}
}

func (t *Registry) Reset() {

	// t.Functions.Reset()

	// t.GoCallbackNodes.Reset()

	t.Allocator.Reset()
}

func (c Context) NewScope(parentPtr uintptr, length int) (*Scope, uintptr) {
	s := arena.Create[Scope](c.State.Registry.Allocator)
	arena.Memclr(s)

	s.ParentPtr = parentPtr
	s.Bindings = arena.Alloc[NamedValue](c.State.Registry.Allocator, length)
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
