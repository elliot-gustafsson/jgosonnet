package evaluator

type Context struct {
	// Cwd      string
	Interner *Interner
	Importer *Importer
	Arena    *Arena

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
	Objects []Object
	Arrays  [][]Value
	Thunks  []Thunk
	Funcs   []Func

	Scopes []Scope
}

type Binding struct {
	Key uint32
	Val Value
}

type Scope struct {
	Bindings []Binding

	ParentId uint32
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

func (a *Arena) Reset() {

	clear(a.Thunks)
	clear(a.Objects)
	clear(a.Arrays)
	clear(a.Funcs)
	clear(a.Scopes)

	a.Thunks = a.Thunks[:0]
	a.Objects = a.Objects[:0]
	a.Arrays = a.Arrays[:0]
	a.Funcs = a.Funcs[:0]
	a.Scopes = a.Scopes[:0]
}
