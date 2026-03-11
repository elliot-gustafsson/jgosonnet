package evaluator

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-jsonnet/ast"
)

type Interpreter struct {
	interner *Interner
	jpaths   []string
	traceOut io.Writer
	// nativeFuncs map[string]*NativeFunction // TODO: make it possible to pass custom funcs
}

type Importer struct {
	JPaths      []string
	ImportScope uint32

	lock  sync.Mutex
	cache map[string]Value
}

func NewImporter(scopeId uint32, jPaths []string) *Importer {
	return &Importer{
		ImportScope: scopeId,
		JPaths:      jPaths,
		cache:       make(map[string]Value, 32),
	}
}

func (i *Importer) Set(path string, v Value) {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.cache[path] = v
}

func (i *Importer) Get(path string) Value {
	i.lock.Lock()
	defer i.lock.Unlock()
	return i.cache[path]
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

type ValueType uint8

const (
	ValueTypeNone ValueType = iota
	ValueTypeNull
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeObject
	ValueTypeArray
	ValueTypeFunction
	ValueTypeThunk
)

func (t ValueType) IsLiteral() bool {
	return t == ValueTypeString || t == ValueTypeNumber || t == ValueTypeBool || t == ValueTypeNull
}

func (t ValueType) String() string {
	switch t {
	case ValueTypeNone:
		return "none"
	case ValueTypeNull:
		return "null"
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBool:
		return "boolean"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeFunction:
		return "function"
	case ValueTypeThunk:
		return "thunk"
	default:
		return fmt.Sprintf("unknown (%d)", t)
	}
}

type Binding struct {
	Key uint32
	Val Value
}

type Scope struct {
	Bindings []Binding

	ParentId uint32
}

type Arena struct {
	Objects []Object
	Arrays  [][]Value
	Thunks  []Thunk
	Funcs   []Func

	Scopes []Scope
}

func NewArena() *Arena {
	return &Arena{
		Thunks:  make([]Thunk, 0, 32*1024),
		Objects: make([]Object, 0, 8*1024),
		Arrays:  make([][]Value, 0, 16*1024),
		Funcs:   make([]Func, 0, 2*1024),
		Scopes:  make([]Scope, 0, 32*1024),
	}
}

type Context struct {
	// Cwd      string
	Interner *Interner
	Importer *Importer
	Arena    *Arena

	// Root Value // $
	Self Value // self

	SuperOffset int
}

func (a *Arena) NewScope(parentId uint32, cap int) uint32 {
	id := uint32(len(a.Scopes))

	a.Scopes = append(a.Scopes, Scope{
		ParentId: parentId,
		Bindings: make([]Binding, 0, cap),
	})

	return id
}

func (a *Arena) GetScope(id uint32) *Scope {
	return &a.Scopes[id]
}

func (a *Arena) AddScopeBind(scopeId, keyId uint32, val Value) {
	s := &a.Scopes[scopeId]

	s.Bindings = append(s.Bindings, Binding{
		Key: keyId,
		Val: val,
	})
}

func (a *Arena) GetScopeBind(scopeId, key uint32) (Value, bool) {
	currId := scopeId

	for {
		scope := &a.Scopes[currId]

		for i := len(scope.Bindings) - 1; i >= 0; i-- {
			if scope.Bindings[i].Key == key {
				return scope.Bindings[i].Val, true
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

// func FindScopeBinding(scopeId, key uint32, ctx Context) (Value, bool) {
// 	currId := scopeId

// 	for {
// 		scope := &ctx.Arena.Scopes[currId]

// 		for i := len(scope.Bindings) - 1; i >= 0; i-- {
// 			if scope.Bindings[i].Key == key {
// 				return scope.Bindings[i].Val, true
// 			}
// 		}

// 		if currId == 0 {
// 			break
// 		}

// 		if scope.ParentId == currId {
// 			break
// 		}

// 		currId = scope.ParentId
// 	}
// 	return Value{}, false
// }

type Thunk struct {
	Node                ast.Node
	ScopeId             uint32
	CapturedSelf        Value
	CapturedSuperOffset int

	Value Value
}

func CreateThunk(node ast.Node, scopeId uint32, ctx Context) Value {
	return MakeThunk(Thunk{
		Node:                node,
		ScopeId:             scopeId,
		CapturedSelf:        ctx.Self,
		CapturedSuperOffset: ctx.SuperOffset,
	}, ctx)
}

type Func = func(args []Value, ctx Context) (Value, error)

type Value struct {
	t ValueType

	// Reference to id in arena for string, object, array, function, thunk
	refId uint32

	// Also holds 1.0 and 0.0 for bool
	num float64
}

func MakeNull() Value {
	return Value{t: ValueTypeNull}
}

func MakeString(v string, ctx Context) Value {
	refId := ctx.Interner.Intern(v)
	return Value{t: ValueTypeString, refId: refId}
}

func MakeNumber(v float64) Value {
	return Value{t: ValueTypeNumber, num: v}
}

func MakeBool(v bool) Value {
	if v {
		return Value{t: ValueTypeBool, num: 1}
	}
	return Value{t: ValueTypeBool, num: 0}

}

func MakeObject(v Object, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Objects))
	ctx.Arena.Objects = append(ctx.Arena.Objects, v)
	return Value{t: ValueTypeObject, refId: refId}
}

func MakeArray(v []Value, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Arrays))
	ctx.Arena.Arrays = append(ctx.Arena.Arrays, v)
	return Value{t: ValueTypeArray, refId: refId}
}

func MakeFunction(v Func, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Funcs))
	ctx.Arena.Funcs = append(ctx.Arena.Funcs, v)
	return Value{t: ValueTypeFunction, refId: refId}
}

func MakeThunk(v Thunk, ctx Context) Value {
	refId := uint32(len(ctx.Arena.Thunks))
	ctx.Arena.Thunks = append(ctx.Arena.Thunks, v)
	return Value{t: ValueTypeThunk, refId: refId}
}

func (v Value) Type() ValueType {
	return v.t
}

func (v Value) String(ctx Context) string {
	return ctx.Interner.Get(v.refId)
}

func (v Value) Number() float64 {
	return v.num
}

func (v Value) Bool() bool {
	return v.num == 1
}

func (v Value) Array(ctx Context) []Value {
	return ctx.Arena.Arrays[v.refId]
}

func (v Value) Object(ctx Context) *Object {
	return &ctx.Arena.Objects[v.refId]
}

func (v Value) Function(ctx Context) Func {
	return ctx.Arena.Funcs[v.refId]
}

func (v Value) Thunk(ctx Context) *Thunk {
	return &ctx.Arena.Thunks[v.refId]
}

func (v Value) Eval(ctx Context) (Value, error) {
	if !v.IsThunk() {
		return v, nil
	}
	thunk := v.Thunk(ctx)
	if !thunk.Value.IsNone() {
		return thunk.Value, nil
	}

	evalCtx := ctx
	evalCtx.Self = thunk.CapturedSelf
	evalCtx.SuperOffset = thunk.CapturedSuperOffset

	evaledVal, err := EvaluateNodeStrict(thunk.Node, thunk.ScopeId, evalCtx)
	if err != nil {
		return Value{}, err
	}
	thunk.Value = evaledVal
	return evaledVal, nil
}

func (v Value) ToString(ctx Context) (string, error) {

	switch v.t {
	default:
		return "", fmt.Errorf("unhandled type %s, string conversion not available", v.t.String())
	case ValueTypeNull:
		return "null", nil
	case ValueTypeString:
		return v.String(ctx), nil
	case ValueTypeNumber:
		res := strconv.FormatFloat(v.Number(), 'f', -1, 64)
		return res, nil
	case ValueTypeBool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case ValueTypeObject, ValueTypeArray:
		var b strings.Builder
		err := ManifestJson(&b, v, ctx, "", "", ": ")
		if err != nil {
			return "", err
		}
		return b.String(), nil
	case ValueTypeThunk:
		err := EvaluateValueStrict(&v, ctx)
		if err != nil {
			return "", err
		}
		return v.ToString(ctx)
	}

}

func (v Value) IsNone() bool {
	return v.t == ValueTypeNone
}

func (v Value) IsLiteral() bool {
	return v.t.IsLiteral()
}

func (v Value) IsNull() bool {
	return v.t == ValueTypeNull
}

func (v Value) IsString() bool {
	return v.t == ValueTypeString
}

func (v Value) IsNumber() bool {
	return v.t == ValueTypeNumber
}

func (v Value) IsBool() bool {
	return v.t == ValueTypeBool
}

func (v Value) IsThunk() bool {
	return v.t == ValueTypeThunk
}

func (v Value) IsObject() bool {
	return v.t == ValueTypeObject
}

func (v Value) IsFunction() bool {
	return v.t == ValueTypeFunction
}

func (v Value) IsArray() bool {
	return v.t == ValueTypeArray
}

func (v Value) IsEmpty(ctx Context) bool {
	switch v.Type() {
	default:
		return false
	case ValueTypeNull:
		return true
	case ValueTypeObject:
		return v.Object(ctx).GetLength() != 0
	case ValueTypeArray:
		return len(v.Array(ctx)) != 0
	}
}

func (v Value) Prune(ctx Context) (Value, error) {
	switch v.Type() {
	default:
		return Value{}, fmt.Errorf("unhandled type (%s) in Value.Prune()", v.Type())
	case ValueTypeNull:
		return MakeNull(), nil
	case ValueTypeString:
		return Value{t: ValueTypeString, refId: v.refId}, nil
	case ValueTypeNumber:
		return Value{t: ValueTypeNumber, num: v.num}, nil
	case ValueTypeBool:
		return Value{t: ValueTypeBool, num: v.num}, nil
	case ValueTypeObject:
		return Value{}, fmt.Errorf("Value.Prune() object not handled yet")
	case ValueTypeArray:
		arr := v.Array(ctx)
		res := make([]Value, 0, len(arr))
		for _, v := range arr {
			err := EvaluateValueStrict(&v, ctx)
			if err != nil {
				return Value{}, err
			}
			out, err := v.Prune(ctx)
			if err != nil {
				return Value{}, err
			}
			if out.IsEmpty(ctx) {
				continue
			}
			res = append(res, out)
		}
		return MakeArray(res, ctx), nil
	}

}
