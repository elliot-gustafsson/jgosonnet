package evaluator

import (
	"io"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Context struct {
	Interner    *Interner
	Registry    *Registry
	Environment *Environment

	Self Value // self

	SuperOffset int
}

type Registry struct {
	Objects   *arena.Arena[Object]
	Arrays    *arena.SliceArena[Value]
	Strings   *arena.StringArena
	Thunks    *arena.Arena[Thunk]
	Functions *arena.Arena[Function]

	Scopes *arena.Arena[Scope]
	Layers *arena.Arena[Layer]

	Nodes           *arena.Arena[ast.Node]
	GoCallbackNodes *arena.Arena[GoCallbackNode]

	Uint8Bufs      *arena.BufferArena[uint8]
	Uint32Bufs     *arena.BufferArena[uint32]
	LayerBufs      *arena.BufferArena[*Layer]
	NamedValueBufs *arena.BufferArena[NamedValue]
	NodesBufs      *arena.BufferArena[ast.Node]

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

		Nodes:           arena.NewArena[ast.Node](),
		GoCallbackNodes: arena.NewArena[GoCallbackNode](),

		Uint8Bufs:      arena.NewBufferArena[uint8](bufferArenaBlockSize),
		Uint32Bufs:     arena.NewBufferArena[uint32](bufferArenaBlockSize),
		LayerBufs:      arena.NewBufferArena[*Layer](bufferArenaBlockSize),
		NamedValueBufs: arena.NewBufferArena[NamedValue](bufferArenaBlockSize),
		NodesBufs:      arena.NewBufferArena[ast.Node](bufferArenaBlockSize),

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

	t.Nodes.Reset()
	t.GoCallbackNodes.Reset()

	t.Uint8Bufs.Reset()
	t.Uint32Bufs.Reset()
	t.LayerBufs.Reset()
	t.NamedValueBufs.Reset()
	t.NodesBufs.Reset()

	t.LayerRefBufs.Reset()
	t.FieldPlanBufs.Reset()
}

func (c Context) NewScope(parentId uint32, length int) (*Scope, uint32) {

	s, id := c.Registry.Scopes.New()

	s.ParentId = parentId
	s.Bindings = c.Registry.NamedValueBufs.Alloc(length, length)

	return s, id
}

func (c Context) GetScopeBind(scopeId, key uint32) (Value, bool) {
	currId := scopeId

	for {
		scope := c.Registry.Scopes.GetPtr(currId)

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
	return Value{}, false
}

type ExtVar interface {
	Eval(scopeId uint32, ctx Context) (Value, error)
}

type ExtString struct {
	Key, Val string
	v        Value
}

func (t *ExtString) Eval(scopeId uint32, ctx Context) (Value, error) {
	if !t.v.IsNone() {
		return t.v, nil
	}
	return MakeString(t.Val, ctx), nil
}

type ExtCode struct {
	Key, Val string
	n        ast.Node
	v        Value
}

func (t *ExtCode) Eval(scopeId uint32, ctx Context) (Value, error) {
	if !t.v.IsNone() {
		return t.v, nil
	}

	if t.n == nil {
		n, err := jsonnet.SnippetToAST(t.Key, t.Val)
		if err != nil {
			return Value{}, err
		}
		t.n = n
	}
	return EvaluateNode(t.n, scopeId, ctx)
}

type Environment struct {
	TraceOut        io.Writer
	Importer        *Importer
	ExtVars         map[string]string
	ExtCodes        map[string]string
	NativeFunctions map[string]Function
}
