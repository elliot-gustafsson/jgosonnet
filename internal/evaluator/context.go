package evaluator

import (
	"io"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet/ast"
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
	Objects   *arena.Arena[Object]
	Arrays    *arena.SliceArena[Value]
	Strings   *arena.StringArena
	Thunks    *arena.Arena[Thunk]
	Functions *arena.Arena[Function]

	Scopes *arena.Arena[Scope]
	Layers *arena.Arena[Layer]

	GoCallbackNodes *arena.Arena[GoCallbackNode]

	NodeSlices *arena.SliceArena[ast.Node]

	Uint8Bufs      *arena.BufferArena[uint8]
	Uint32Bufs     *arena.BufferArena[uint32]
	LayerBufs      *arena.BufferArena[*Layer]
	NamedValueBufs *arena.BufferArena[NamedValue]
	NodeBufs       *arena.BufferArena[ast.Node]

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
		Objects:   arena.NewArena[Object](),
		Arrays:    arena.NewSliceArena[Value](sliceArenaChunkSize),
		Strings:   arena.NewStringArena(stringArenaBlockSize),
		Thunks:    arena.NewArena[Thunk](),
		Functions: arena.NewArena[Function](),

		Scopes: arena.NewArena[Scope](),
		Layers: arena.NewArena[Layer](),

		GoCallbackNodes: arena.NewArena[GoCallbackNode](),

		NodeSlices: arena.NewSliceArena[ast.Node](128),

		Uint8Bufs:      arena.NewBufferArena[uint8](bufferArenaBlockSize),
		Uint32Bufs:     arena.NewBufferArena[uint32](bufferArenaBlockSize),
		LayerBufs:      arena.NewBufferArena[*Layer](bufferArenaBlockSize),
		NamedValueBufs: arena.NewBufferArena[NamedValue](bufferArenaBlockSize),
		NodeBufs:       arena.NewBufferArena[ast.Node](bufferArenaBlockSize),

		LayerRefBufs:  arena.NewBufferArena[LayerRef](bufferArenaBlockSize),
		FieldPlanBufs: arena.NewBufferArena[FieldPlan](bufferArenaBlockSize),
	}
}

func (t *Registry) Reset() {
	t.Objects.Reset()
	t.Arrays.Reset()
	t.Strings.Reset()
	t.Thunks.Reset()
	t.Functions.Reset()

	t.Scopes.Reset()
	t.Layers.Reset()

	t.GoCallbackNodes.Reset()

	t.NodeSlices.Reset()

	t.Uint8Bufs.Reset()
	t.Uint32Bufs.Reset()
	t.LayerBufs.Reset()
	t.NamedValueBufs.Reset()
	t.NodeBufs.Reset()

	t.LayerRefBufs.Reset()
	t.FieldPlanBufs.Reset()
}

func (c Context) NewScope(parentId uint32, length int) (*Scope, uint32) {

	s, id := c.State.Registry.Scopes.New()

	s.ParentId = parentId
	s.Bindings = c.State.Registry.NamedValueBufs.Alloc(length, length)

	return s, id
}

func (c Context) GetScopeBind(scopeId, key uint32) (val Value, found bool) {
	currId := scopeId

	scopes := c.State.Registry.Scopes

	for {
		scope := scopes.GetPtr(currId)
		bindings := scope.Bindings

		for i := len(bindings) - 1; i >= 0; i-- {
			if bindings[i].Key == key {
				val, found = bindings[i].Value, true
				return
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
	return
}

type Environment struct {
	TraceOut        io.Writer
	Importer        *Importer
	ExtVars         map[string]string
	ExtCodes        map[string]string
	NativeFunctions map[string]Function
}
