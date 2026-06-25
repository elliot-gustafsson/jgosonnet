package evaluator

import (
	"fmt"
)

type Stack struct {
	s []Value
	p int
}

func NewStack(size int) Stack {
	return Stack{
		s: make([]Value, size),
	}
}

func (s *Stack) Push(v Value) {
	s.s[s.p] = v
	s.p++
}

func (s *Stack) Pop() Value {
	s.p--
	return s.s[s.p]
}

func (s *Stack) Length() int {
	return s.p
}

type VM struct {
	// prog  Program
	Stack Stack
	// ip    int

	Registry *Registry

	callFrames []CallFrame
	fp         int

	ctx   Context
	state CallFrame
}

type CallFrame struct {
	ip      uint32
	scopeId uint32
	self    Value

	superOffset int
	thunkId     uint32
}

func NewVM() *VM {
	return &VM{
		// prog:  program,
		Stack:    NewStack(1024),
		Registry: NewRegistry(),
		// ctx: Context{
		// 	Interner: program.Interner,
		// 	Registry: program.Registry,
		// },
	}
}

func (vm *VM) Run(prog *Program) (Value, error) {
	for vm.state.ip < uint32(len(prog.Instructions)) {
		inst := prog.Instructions[vm.state.ip]

		switch inst.op {
		default:
			return Value{}, fmt.Errorf("unhandled op code: %d", inst.op)

		case OpPushNull:
			vm.Stack.Push(MakeNull())

		case OpPushString:
			s := prog.Strings[inst.operand]
			id := vm.Registry.Strings.Alloc(s)
			vm.Stack.Push(Value{t: ValueTypeString, refId: id})

		case OpAdd:
			right := vm.Stack.Pop()
			left := vm.Stack.Pop()

			vm.Stack.Push(MakeNumber(left.Number() + right.Number()))

		case OpMakeThunk:

			t := Thunk{
				NodeId:              vm.state.ip + 1,
				ScopeId:             vm.state.scopeId,
				CapturedSelf:        vm.state.self,
				CapturedSuperOffset: vm.state.superOffset,
			}
			tv := MakeThunk(t, vm.ctx)

			vm.Stack.Push(tv)

			vm.state.ip += inst.operand

		case OpPushScope:
			childScopeId := vm.ctx.NewScope(vm.state.scopeId, int(inst.operand))
			// vm.state.scopeId = childScopeId
			// s := vm.ctx.Registry.Scopes.GetPtr(childScopeId)

			vm.state.scopeId = childScopeId

		case OpLocalSet:
			v := vm.Stack.Pop()
			nv := NamedValue{inst.operand, v}

			activeScope := vm.ctx.Registry.Scopes.GetPtr(vm.state.scopeId)
			activeScope.Bindings = append(activeScope.Bindings, nv)

		case OpMakeObject:
			// vm.handleMakeObject(prog, int(inst.operand), inst.operand2)

		case OpMakeArray:
			length := int(inst.operand)

			elements, refId := vm.Registry.Arrays.Make(length)

			for i := length - 1; i >= 0; i-- {
				elements[i] = vm.Stack.Pop()
			}

			vm.Stack.Push(Value{t: ValueTypeArray, refId: refId})

		}

		vm.state.ip++
	}

	if vm.Stack.Length() == 1 {
		return vm.Stack.Pop(), nil
	}
	return Value{}, fmt.Errorf("execution finished with invalid stack state")
}

// func (vm *VM) handleMakeObject(prog *Program, fieldCount int, flags uint16) {

// 	var localsCount int
// 	var assertsCount int

// 	// Because the compiler emitted them in order, we just step the IP!
// 	if (flags & FlagObjectHasLocals) != 0 {
// 		vm.state.ip++
// 		localsCount = int(prog.Instructions[vm.state.ip].operand)
// 	}

// 	if (flags & FlagObjectHasAsserts) != 0 {
// 		vm.state.ip++
// 		assertsCount = int(prog.Instructions[vm.state.ip].operand)
// 	}

// 	layerId := vm.Registry.Layers.Alloc(Layer{})
// 	layer := vm.Registry.Layers.GetPtr(layerId)

// 	layer.ParentScopeId = scopeId

// 	if fieldCount > 0 {
// 		layer.Keys = vm.Registry.Uint32Bufs.Alloc(0, fieldCount)
// 		layer.Nodes = vm.Registry.NodesBufs.Alloc(0, fieldCount)
// 		layer.Meta = vm.Registry.Uint8Bufs.Alloc(0, fieldCount)
// 	}

// 	if localsCount > 0 {
// 		layer.LocalKeys = vm.Registry.Uint32Bufs.Alloc(0, localsCount)
// 		layer.LocalNodes = vm.Registry.NodesBufs.Alloc(0, localsCount)
// 	}

// 	if assertsCount > 0 {
// 		layer.Asserts = vm.Registry.NodesBufs.Alloc(len(node.Asserts), len(node.Asserts))
// 		copy(layer.Asserts, node.Asserts)
// 	}

// 	for _, v := range node.Locals {

// 		name := string(v.Variable)
// 		keyId := ctx.Interner.Intern(name)

// 		layer.LocalKeys = append(layer.LocalKeys, keyId)
// 		layer.LocalNodes = append(layer.LocalNodes, v.Body)

// 	}

// 	useMap := fieldCount > MaxLinearKeys

// 	if useMap {
// 		layer.Index = make(map[uint32]int, fieldCount)
// 	}

// 	index := 0
// 	for _, v := range node.Fields {
// 		name, err := EvaluateNode(v.Name, scopeId, ctx)
// 		if err != nil {
// 			return Value{}, err
// 		}

// 		if name.IsNull() {
// 			// Omitted field
// 			continue
// 		}

// 		if !name.IsString() {
// 			return Value{}, fmt.Errorf("unexpected field name type %s, expected string", name.Type().String())
// 		}

// 		n := name.String(ctx)

// 		keyId := ctx.Interner.Intern(n)

// 		layer.Keys = append(layer.Keys, keyId)
// 		layer.Nodes = append(layer.Nodes, v.Body)
// 		layer.Meta = append(layer.Meta, CreateFieldMeta(v.Hide, v.PlusSuper))

// 		if useMap {
// 			layer.Index[keyId] = index
// 			index++
// 		}

// 	}

// 	layers := ctx.Registry.LayerBufs.Alloc(1, 1)
// 	layers[0] = layer
// 	obj := NewObject(layers)
// 	refId := prog.Registry.Objects.Alloc(obj)
// 	val := Value{t: ValueTypeObject, refId: refId}

// 	vm.Stack.Push(val)
// }
