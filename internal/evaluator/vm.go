package evaluator

import (
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

func (s *Stack) Peek() Value {
	return s.s[s.p-1]
}

func (s *Stack) Replace(v Value) {
	s.s[s.p-1] = v
}

func (s *Stack) Pop2() (v1 Value, v2 Value) {
	v1 = s.s[s.p-1]
	v2 = s.s[s.p-2]
	s.p = s.p - 2
	return
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

	bp      int
	scopeBp int

	superOffset int32
	thunkId     uint32

	posCount   int
	namedCount int
}

type ArgMapper []Value

func (m *ArgMapper) Prepare(size int) {
	if cap(*m) < size {
		*m = make([]Value, size)
	}
	*m = (*m)[:size]
	clear(*m)
}

func NewVM(interner *Interner, registry *Registry) *VM {
	return &VM{
		// prog:  program,
		Stack:      NewStack(1024),
		Interner:   interner,
		Registry:   registry,
		callFrames: make([]CallFrame, 1024),
		ctx: Context{
			State: &ContextState{
				Interner: interner,
				Registry: registry,
			},
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
		return ValueNone, err
	}

	if vm.Stack.Length() == 1 {
		return vm.Stack.Pop(), nil
	}

	return ValueNone, fmt.Errorf("execution finished with invalid stack state")
}

func (vm *VM) runLoop(targetFp int) error {

	var mapper ArgMapper = make([]Value, 0, 32)

	for {
		inst := vm.prog.Instructions[vm.state.ip]

		switch inst.op {
		default:
			return fmt.Errorf("unhandled op code: %d", inst.op)

		// Execution

		case OpReturn:
			res := vm.Stack.Pop()

			vm.Stack.p = vm.state.bp

			isThunkEval := vm.state.thunkId != 0

			if isThunkEval && !res.IsThunk() {
				t := vm.ctx.State.Registry.Thunks.GetPtr(vm.state.thunkId)
				t.Value = res
			}

			vm.fp--
			vm.state = vm.callFrames[vm.fp]

			if vm.fp == targetFp {
				// target hit, push res and exit
				vm.Stack.Push(res)
				return nil
			}

			// replace top stack value with res and continue
			// happends after forcing thunks
			vm.Stack.Replace(res)

			if isThunkEval && res.IsThunk() {
				// keep looping until we get a solid value
				continue
			}

		case OpForce:

			val := vm.Stack.Peek()

			// Fast path: not a thunk, do nothing!
			if !val.IsThunk() {
				vm.state.ip++
				continue
			}

			thunk := val.Thunk(vm.ctx)

			if !thunk.Value.IsNone() {
				vm.Stack.Replace(thunk.Value)
				vm.state.ip++
				continue
			}

			// 1. Save current state
			vm.callFrames[vm.fp] = vm.state
			vm.fp++

			vm.state = CallFrame{
				ip:          thunk.NodeId,
				bp:          vm.Stack.p,
				scopeBp:     thunk.CapturedBp,
				self:        thunk.CapturedSelf,
				superOffset: thunk.CapturedSuperOffset,
				scopeId:     thunk.ScopeId,
				thunkId:     val.RefId(),
			}

			continue

		case OpJump:
			vm.state.ip += inst.operand
			// continue

		// case OpJumpIfFalse:

		// Values

		case OpPushNull:
			vm.Stack.Push(MakeNull())

		case OpPushString:
			s := vm.prog.Strings[inst.operand]
			id := vm.Registry.Strings.Alloc(s)
			vm.Stack.Push(MakeStringValue(id))

		case OpPushNumber:
			n := vm.prog.Numbers[inst.operand]
			vm.Stack.Push(MakeNumber(n))

		case OpPushFalse:
			vm.Stack.Push(MakeFalse())

		case OpPushTrue:
			vm.Stack.Push(MakeTrue())

		// ---------------------------------------------------------
		// Var resolving
		// ---------------------------------------------------------

		case OpAllocLocals:
			count := int(inst.operand)
			vm.Stack.p += count

		case OpPopLocals:
			// The body result is currently at the top of the stack.
			// We need to save it, remove the locals, and put it back.
			res := vm.Stack.Pop()
			vm.Stack.p -= int(inst.operand)
			vm.Stack.Push(res)

		case OpLocalSet:
			val := vm.Stack.Pop()
			// slotIndex := vm.state.bp + int(inst.operand)
			slotIndex := vm.state.scopeBp + int(inst.operand)
			vm.Stack.s[slotIndex] = val

		case OpLocalGet:
			// slotIndex := vm.state.bp + int(inst.operand)
			slotIndex := vm.state.scopeBp + int(inst.operand)
			val := vm.Stack.s[slotIndex]
			vm.Stack.Push(val)

		// ---------------------------------------------------------
		// Make objects
		// ---------------------------------------------------------

		case OpMakeThunk:

			t := Thunk{
				NodeId:              vm.state.ip + 1,
				ScopeId:             vm.state.scopeId,
				CapturedSelf:        vm.state.self,
				CapturedSuperOffset: vm.state.superOffset,
				CapturedBp:          vm.state.scopeBp,
			}
			tv := MakeThunk(t, vm.ctx)

			vm.Stack.Push(tv)

			vm.state.ip += inst.operand

		case OpMakeArray:
			length := int(inst.operand)

			elements, refId := vm.Registry.Arrays.Make(length)

			for i := length - 1; i >= 0; i-- {
				elements[i] = vm.Stack.Pop()
			}

			vm.Stack.Push(MakeArrayValue(refId))

		case OpMakeFunction:
			f := Function{
				ip:        inst.operand,
				argsCount: int(inst.operand2),
				scopeBp:   vm.state.scopeBp,
				metaIp:    vm.state.ip + 1, // Metadata starts right here!
			}
			vm.Stack.Push(MakeFunction(f, vm.ctx))

			// Skip over the inline metadata (2 instructions per parameter)
			vm.state.ip += uint32(f.argsCount * 2)

		// ---------------------------------------------------------
		// Functions
		// ---------------------------------------------------------

		case OpCall:
			numPos := int(inst.operand)
			numNamed := int(inst.operand2)

			// MAGIC: Named args only consume 1 stack slot now (the value Thunk)!
			totalArgSlots := numPos + numNamed

			funcStackIdx := vm.Stack.p - 1 - totalArgSlots
			val := vm.Stack.s[funcStackIdx]

			if !val.IsFunction() {
				return TypeErrorSpecific(ValueTypeFunction, val.Type())
			}

			fn := val.Function(vm.ctx)

			if numPos > fn.argsCount {
				return MakeRuntimeError(fmt.Errorf("function expected at most %v positional argument(s), but got %v", fn.argsCount, numPos))
			}

			mapper.Prepare(fn.argsCount)

			// 1. Map Positional
			argsBase := funcStackIdx + 1
			for i := 0; i < numPos; i++ {
				mapper[i] = vm.Stack.s[argsBase+i]
			}

			// 2. Map Named (Reading NameIDs directly from bytecode!)
			namedBase := argsBase + numPos
			for i := 0; i < numNamed; i++ {
				nameId := vm.prog.Instructions[vm.state.ip+1+uint32(i)].operand
				argVal := vm.Stack.s[namedBase+i]

				slotIdx := -1
				for p := 0; p < fn.argsCount; p++ {
					metaInst := vm.prog.Instructions[fn.metaIp+uint32(p*2)]
					if metaInst.operand == nameId {
						slotIdx = p
						break
					}
				}

				if slotIdx == -1 {
					return MakeRuntimeError(fmt.Errorf("function has no parameter for named arg"))
				}
				if !mapper[slotIdx].IsNone() {
					return MakeRuntimeError(fmt.Errorf("argument bound multiple times"))
				}

				mapper[slotIdx] = argVal
			}

			// 3. Map Defaults
			for i := 0; i < fn.argsCount; i++ {
				if mapper[i].IsNone() {
					metaInst1 := vm.prog.Instructions[fn.metaIp+uint32(i*2)]
					if metaInst1.operand2 != 1 {
						return MakeRuntimeError(fmt.Errorf("missing required parameter %d", i))
					}

					metaInst2 := vm.prog.Instructions[fn.metaIp+uint32(i*2)+1]
					t := Thunk{
						NodeId:              metaInst2.operand,
						ScopeId:             vm.state.scopeId,
						CapturedSelf:        vm.state.self,
						CapturedSuperOffset: vm.state.superOffset,
						CapturedBp:          funcStackIdx,
					}
					mapper[i] = MakeThunk(t, vm.ctx)
				}
			}

			// 4. Commit to VM Stack
			vm.Stack.p = funcStackIdx + 1
			for i := 0; i < fn.argsCount; i++ {
				vm.Stack.Push(mapper[i])
			}

			// 5. Advance the OLD frame's IP over the inline caller metadata
			vm.state.ip += uint32(numNamed)

			// 6. Jump into the function body
			vm.callFrames[vm.fp] = vm.state
			vm.fp++

			vm.state = CallFrame{
				ip:          fn.ip,
				bp:          funcStackIdx + 1,
				scopeBp:     funcStackIdx + 1,
				self:        vm.state.self,
				superOffset: vm.state.superOffset,
				scopeId:     vm.state.scopeId,
			}

			continue

			// fmt.Println(finalArgs)
			// fmt.Println(paramCount)
			// fmt.Println(posCount)
			// fmt.Println(namedCount)
			// fmt.Println(passedArgs)

		// ---------------------------------------------------------
		// Binary ops
		// ---------------------------------------------------------

		case OpPlus:
			right, left := vm.Stack.Pop2()

			var res Value
			if left.IsNumber() && right.IsNumber() {
				res = MakeNumber(left.Number() + right.Number())
			} else {
				var err error
				res, err = bopPlus(left, right, vm.ctx)
				if err != nil {
					return err
				}
			}
			vm.Stack.Push(res)

		case OpMinus:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			vm.Stack.Push(MakeNumber(left.Number() - right.Number()))

		case OpDiv:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			vm.Stack.Push(MakeNumber(left.Number() / right.Number()))

		case OpMult:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			vm.Stack.Push(MakeNumber(left.Number() * right.Number()))

		case OpBitwiseAnd:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			val, err := builtinBitwiseAnd(left.Number(), right.Number())
			if err != nil {
				return err
			}
			vm.Stack.Push(MakeNumber(val))

		case OpBitwiseOr:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			val, err := builtinBitwiseOr(left.Number(), right.Number())
			if err != nil {
				return err
			}
			vm.Stack.Push(MakeNumber(val))

		case OpBitwiseXor:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			val, err := builtinBitwiseXor(left.Number(), right.Number())
			if err != nil {
				return err
			}
			vm.Stack.Push(MakeNumber(val))

		case OpShiftL:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			val, err := builtinShiftL(left.Number(), right.Number())
			if err != nil {
				return err
			}
			vm.Stack.Push(MakeNumber(val))

		case OpShiftR:
			right, left := vm.Stack.Pop2()

			if !left.IsNumber() || !right.IsNumber() {
				return typeErrorNotNumber(left, right)
			}
			val, err := builtinShiftR(left.Number(), right.Number())
			if err != nil {
				return err
			}
			vm.Stack.Push(MakeNumber(val))
		}

		vm.state.ip++
	}
}

func (vm *VM) RunForce(value Value) (Value, error) {
	if !value.IsThunk() {
		return value, nil
	}

	thunk := value.Thunk(vm.ctx)
	if !thunk.Value.IsNone() {
		return thunk.Value, nil
	}

	targetFp := vm.fp

	vm.callFrames[vm.fp] = vm.state
	vm.fp++
	vm.Stack.Push(value)

	vm.state = CallFrame{
		ip:          thunk.NodeId,
		bp:          vm.Stack.p,
		scopeBp:     thunk.CapturedBp,
		self:        thunk.CapturedSelf,
		superOffset: thunk.CapturedSuperOffset,
		scopeId:     thunk.ScopeId,
		thunkId:     value.RefId(),
	}

	err := vm.runLoop(targetFp)
	if err != nil {
		return ValueNone, err
	}

	return vm.Stack.Pop(), nil
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
	case ValueTypeNativeFunction:
		res, err := value.NativeFunction(vm.ctx).Exec(nil, vm.ctx)
		if err != nil {
			return nil, err
		}
		return vm.ManifestValue(res)
	case ValueTypeThunk:
		v, err := vm.RunForce(value)
		if err != nil {
			return nil, err
		}
		return vm.ManifestValue(v)
	}
}

func typeErrorNotNumber(left, right Value) error {
	if !left.IsNumber() {
		return TypeErrorSpecific(ValueTypeNumber, left.Type())
	}
	if !right.IsNumber() {
		return TypeErrorSpecific(ValueTypeNumber, right.Type())
	}
	return nil
}

// func (vm *VM) handleMakeObject(fieldCount int, flags uint16) {

// 	var localsCount int
// 	var assertsCount int

// 	// Because the compiler emitted them in order, we just step the IP!
// 	if (flags & FlagObjectHasLocals) != 0 {
// 		vm.state.ip++
// 		localsCount = int(vm.prog.Instructions[vm.state.ip].operand)
// 	}

// 	if (flags & FlagObjectHasAsserts) != 0 {
// 		vm.state.ip++
// 		assertsCount = int(vm.prog.Instructions[vm.state.ip].operand)
// 	}

// 	layer, _ := vm.Registry.Layers.New()

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
