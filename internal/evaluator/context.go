package evaluator

import (
	"io"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Context struct {
	Interner    *Interner
	Arena       *Arena
	Environment *Environment

	Self Value // self

	SuperOffset int
}

type Interner struct {
	mapping map[string]uint32
	strings []string
}

func NewInterner() *Interner {
	return &Interner{
		mapping: make(map[string]uint32, 8192),
		strings: make([]string, 0, 8192),
	}
}

func (i *Interner) Intern(s string) uint32 {
	if id, ok := i.mapping[s]; ok {
		return id
	}

	id := uint32(len(i.strings))
	i.strings = append(i.strings, s)
	i.mapping[s] = id

	return id
}

func (i *Interner) Get(id uint32) string {
	if id >= uint32(len(i.strings)) {
		return ""
	}
	return i.strings[id]
}

type Arena struct {
	Objects   []Object
	Arrays    [][]Value
	Thunks    []Thunk
	Functions []Function

	Scopes []Scope

	bindings []NamedValue
}

type Scope struct {
	Bindings []NamedValue

	ParentId uint32
}

func NewArena() *Arena {
	return &Arena{
		Thunks:    make([]Thunk, 0, 32*1024),
		Objects:   make([]Object, 0, 8*1024),
		Arrays:    make([][]Value, 0, 16*1024),
		Functions: make([]Function, 0, 2*1024),
		Scopes:    make([]Scope, 0, 32*1024),

		bindings: make([]NamedValue, 0, 128*1024),
	}
}

func (a *Arena) NewScope(parentId uint32, cap int) uint32 {
	id := uint32(len(a.Scopes))

	a.Scopes = append(a.Scopes, Scope{
		ParentId: parentId,
		Bindings: a.makeBindings(cap),
	})

	return id
}

func (a *Arena) makeBindings(n int) []NamedValue {
	if n == 0 {
		return nil
	}

	start := len(a.bindings)
	a.bindings = append(a.bindings, make([]NamedValue, n)...)
	total := start + n

	return a.bindings[start:start:total]
}

func (a *Arena) GetScope(id uint32) *Scope {
	return &a.Scopes[id]
}

func (a *Arena) AddScopeBind(scopeId uint32, val NamedValue) {
	s := &a.Scopes[scopeId]

	s.Bindings = append(s.Bindings, val)
}

func (a *Arena) GetScopeBind(scopeId, key uint32) (Value, bool) {
	currId := scopeId

	for {
		scope := &a.Scopes[currId]

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

func (a *Arena) Reset() {

	clear(a.Thunks)
	clear(a.Objects)
	clear(a.Arrays)
	clear(a.Functions)
	clear(a.Scopes)
	clear(a.bindings)

	a.Thunks = a.Thunks[:0]
	a.Objects = a.Objects[:0]
	a.Arrays = a.Arrays[:0]
	a.Functions = a.Functions[:0]
	a.Scopes = a.Scopes[:0]
	a.bindings = a.bindings[:0]
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

	Location *ast.LocationRange
}
