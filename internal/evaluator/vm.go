package evaluator

import (
	"errors"
	"fmt"
	"strings"
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
	prog *Program

	Stack Stack
	// ip    int

	Registry *Registry
	Interner *Interner

	callFrames []CallFrame
	fp         int

	ctx   Context
	state CallFrame
}

type CallFrame struct {
	ip      uint32
	scopeId uint32
	self    Value

	bp int

	superOffset int
	thunkId     uint32
}

func NewVM(interner *Interner, registry *Registry) *VM {
	return &VM{
		// prog:  program,
		Stack:      NewStack(1024),
		Interner:   interner,
		Registry:   registry,
		callFrames: make([]CallFrame, 1024),
		ctx: Context{
			Interner: interner,
			Registry: registry,
		},
	}
}

func (vm *VM) Run(prog *Program) (Value, error) {
	vm.prog = prog

	vm.fp = 1
	vm.callFrames[vm.fp] = CallFrame{
		ip: 0,
		bp: 0,
	}

	err := vm.runLoop(0)
	if err != nil {
		return Value{}, err
	}

	if vm.Stack.Length() == 1 {
		return vm.Stack.Pop(), nil
	}

	return Value{}, fmt.Errorf("execution finished with invalid stack state")
}

func (vm *VM) runLoop(targetFp int) error {

	for {
		inst := vm.prog.Instructions[vm.state.ip]

		switch inst.op {
		default:
			return fmt.Errorf("unhandled op code: %d", inst.op)

		// Execution

		case OpReturn:
			res := vm.Stack.Pop()

			vm.Stack.p = vm.state.bp
			vm.Stack.Push(res)

			if vm.fp == targetFp {
				return nil
			}

			vm.fp--
			vm.state = vm.callFrames[vm.fp]

			continue

		// case OpJump:
		// case OpJumpIfFalse:

		// Values

		case OpPushNull:
			vm.Stack.Push(MakeNull())

		case OpPushString:
			s := vm.prog.Strings[inst.operand]
			id := vm.Registry.Strings.Alloc(s)
			vm.Stack.Push(Value{t: ValueTypeString, refId: id})

		case OpPushNumber:
			n := vm.prog.Numbers[inst.operand]
			vm.Stack.Push(MakeNumber(n))

		case OpAdd:
			right := vm.Stack.Pop()
			left := vm.Stack.Pop()

			res, err := bopPlus(right, left, vm.ctx)
			if err != nil {
				return err
			}

			vm.Stack.Push(res)

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

			childScopeId := Scope{
				ParentId: vm.state.scopeId,
				Bindings: vm.Registry.NamedValueBufs.Alloc(0, int(inst.operand)),
			}

			vm.state.scopeId = vm.Registry.Scopes.Alloc(childScopeId)

		case OpPopScope:
			current := vm.Registry.Scopes.GetPtr(vm.state.scopeId)
			vm.state.scopeId = current.ParentId

		case OpLocalSet:
			v := vm.Stack.Pop()
			nv := NamedValue{inst.operand, v}

			activeScope := vm.Registry.Scopes.GetPtr(vm.state.scopeId)
			activeScope.Bindings = append(activeScope.Bindings, nv)

		case OpLocalGet:
			slotIndex := inst.operand
			depthDiff := inst.operand2

			targetScopeId := vm.state.scopeId
			for range depthDiff {
				s := vm.Registry.Scopes.GetPtr(targetScopeId)
				targetScopeId = s.ParentId
			}

			s := vm.Registry.Scopes.GetPtr(targetScopeId)
			val := s.Bindings[slotIndex].Value

			vm.Stack.Push(val)

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
}

func (vm *VM) Force(value Value) (Value, error) {
	if !value.IsThunk() {
		return value, nil
	}

	thunk := value.Thunk(vm.ctx)
	if !thunk.Value.IsNone() {
		return thunk.Value, nil // already evaluated (memoized)
	}

	// savedState := vm.state

	targetFp := vm.fp

	if vm.fp >= len(vm.callFrames) {
		return Value{}, MakeRuntimeError(errors.New("maximum call stack size exceeded"))
	}

	vm.callFrames[vm.fp] = vm.state
	vm.fp++

	vm.state = CallFrame{
		ip: thunk.NodeId, // We stored the IP here!

		bp: vm.Stack.p,

		scopeId:     thunk.ScopeId,
		self:        thunk.CapturedSelf,
		superOffset: thunk.CapturedSuperOffset,
	}

	err := vm.runLoop(targetFp)
	if err != nil {
		return Value{}, err
	}

	result := vm.Stack.Pop()

	thunk = value.Thunk(vm.ctx)
	thunk.Value = result

	return result, nil

}

func (vm *VM) ManifestValue(value Value) (any, error) {
	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type '%s'", value.Type().String())
	case ValueTypeNull:
		return nil, nil
	case ValueTypeString:
		s := vm.Registry.Strings.Get(value.RefId())
		return strings.Clone(s), nil
	case ValueTypeNumber:
		n := value.Number()
		if float64(int64(n)) == n {
			return int64(n), nil
		}
		return n, nil
	case ValueTypeBool:
		return value.Bool(), nil
	// case ValueTypeObject:
	// 	subCtx := ctx
	// 	subCtx.Self = value
	// 	return manifestObject(value.Object(ctx), subCtx)
	case ValueTypeArray:
		arr := value.Array(vm.ctx)
		res := make([]any, 0, len(arr))
		for i := range arr {
			ev, err := vm.ManifestValue(arr[i])
			if err != nil {
				return nil, err
			}
			res = append(res, ev)
		}
		return res, nil
	// case ValueTypeFunction:
	// 	res, err := value.Function(ctx).Exec(nil, ctx)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return ManifestValue(res)
	case ValueTypeThunk:
		v, err := vm.Force(value)
		if err != nil {
			return nil, err
		}
		return vm.ManifestValue(v)
	}
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
