package evaluator

import (
	"io"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
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

	AstId uint32

	SuperOffset uint32
}

type Registry struct {
	ASTs []*ast.AST

	Objects   *arena.Arena[Object]
	Arrays    *arena.SliceArena[Value]
	Strings   *arena.StringArena
	Thunks    *arena.Arena[Thunk]
	Functions *arena.Arena[Function]

	Scopes         *arena.Arena[Scope]
	Layers         *arena.Arena[Layer]
	CallbackThunks *arena.Arena[CallbackThunk]

	Uint8Bufs      *arena.BufferArena[uint8]
	Uint32Bufs     *arena.BufferArena[uint32]
	LayerBufs      *arena.BufferArena[*Layer]
	NamedValueBufs *arena.BufferArena[NamedValue]
	// NodesBufs      *arena.BufferArena[ast.Node]

	LayerRefBufs  *arena.BufferArena[LayerRef]
	FieldPlanBufs *arena.BufferArena[FieldPlan]
}

type Scope struct {
	Bindings []NamedValue

	ParentId uint32
}

const sliceArenaChunkSize = 4096
const stringArenaBlockSize = 4096
const bufferArenaBlockSize = 4096

func NewRegistry() *Registry {
	return &Registry{
		ASTs: make([]*ast.AST, 32),

		Objects:   arena.NewArena[Object](),
		Arrays:    arena.NewSliceArena[Value](sliceArenaChunkSize),
		Strings:   arena.NewStringArena(stringArenaBlockSize),
		Thunks:    arena.NewArena[Thunk](),
		Functions: arena.NewArena[Function](),

		Scopes:         arena.NewArena[Scope](),
		Layers:         arena.NewArena[Layer](),
		CallbackThunks: arena.NewArena[CallbackThunk](),

		Uint8Bufs:      arena.NewBufferArena[uint8](bufferArenaBlockSize),
		Uint32Bufs:     arena.NewBufferArena[uint32](bufferArenaBlockSize),
		LayerBufs:      arena.NewBufferArena[*Layer](bufferArenaBlockSize),
		NamedValueBufs: arena.NewBufferArena[NamedValue](bufferArenaBlockSize),

		LayerRefBufs:  arena.NewBufferArena[LayerRef](bufferArenaBlockSize),
		FieldPlanBufs: arena.NewBufferArena[FieldPlan](bufferArenaBlockSize),
	}
}

func (t *Registry) Reset() {
	t.ASTs = t.ASTs[:0]

	t.Objects.Reset()
	t.Arrays.Reset()
	t.Strings.Reset()
	t.Thunks.Reset()
	t.Functions.Reset()

	t.Scopes.Reset()
	t.Layers.Reset()
	t.CallbackThunks.Reset()

	t.Uint8Bufs.Reset()
	t.Uint32Bufs.Reset()
	t.LayerBufs.Reset()
	t.NamedValueBufs.Reset()

	t.LayerRefBufs.Reset()
	t.FieldPlanBufs.Reset()
}

func (c Context) NewScope(parentId uint32, length int) (*Scope, uint32) {

	s, id := c.State.Registry.Scopes.New()

	s.ParentId = parentId
	s.Bindings = c.State.Registry.NamedValueBufs.Alloc(length, length)

	return s, id
}

func (c Context) GetScopeBind(scopeId, key uint32) (Value, bool) {
	currId := scopeId

	for {
		scope := c.State.Registry.Scopes.GetPtr(currId)

		for i := len(scope.Bindings) - 1; i >= 0; i-- {
			if scope.Bindings[i].Key == key {
				return scope.Bindings[i].Value, true
			}
		}

		if currId == 0 {
			break
		}

		if scope.ParentId == currId {
			break
		}

		currId = scope.ParentId
	}
	return ValueNone, false
}

type Environment struct {
	BaseScopeId     uint32
	TraceOut        io.Writer
	Importer        *Importer
	ExtVars         map[string]string
	ExtCodes        map[string]string
	NativeFunctions map[string]Function
}
