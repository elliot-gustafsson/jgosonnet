package evaluator

import (
	"fmt"
	"strconv"

	"github.com/google/go-jsonnet/ast"
)

type OpCode int

const (
	OpPushConst OpCode = iota
	OpAdd
	OpLocalGet
	OpLocalSet

	OpObjectOp
	OpJumpIfFalse
)

type Program struct {
	Instructions []Instruction
	Constants    []Value // All parsed numbers, strings, etc.
}

type Compiler struct {
	prog Program
}

type Instruction struct {
	op      OpCode
	operand uint32
}

func (c *Compiler) Compile(node ast.Node) (Program, error) {
	err := c.visit(node)
	if err != nil {
		return Program{}, err
	}
	return c.prog, nil
}

func (c *Compiler) visit(n ast.Node) error {
	switch node := n.(type) {
	default:
		return fmt.Errorf("unsupported AST node type: %T", n)
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		val := MakeNumber(num)

		// 2. Put it in the constants pool
		constIdx := uint32(len(c.prog.Constants))
		c.prog.Constants = append(c.prog.Constants, val)

		// 3. Emit instruction to push the constant at that index!
		c.emit(OpPushConst, constIdx)
	case *ast.Binary:
		if err := c.visit(node.Left); err != nil {
			return err
		}
		if err := c.visit(node.Right); err != nil {
			return err
		}

		switch node.Op {
		default:
			return fmt.Errorf("unsupported binary operator: %s", node.Op.String())
		case ast.BopPlus:
			c.emit(OpAdd, 0)
		}
	}

	return nil
}

func (c *Compiler) emit(op OpCode, operand uint32) {
	c.prog.Instructions = append(c.prog.Instructions, Instruction{op, operand})
}

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
	prog  Program
	stack Stack
	ip    int
}

func NewVM(program Program) *VM {
	return &VM{
		prog:  program,
		stack: NewStack(1024),
	}
}

func (vm *VM) Run() (Value, error) {
	for vm.ip < len(vm.prog.Instructions) {
		inst := vm.prog.Instructions[vm.ip]

		switch inst.op {
		case OpPushConst:
			vm.stack.Push(vm.prog.Constants[inst.operand])
		case OpAdd:
			right := vm.stack.Pop()
			left := vm.stack.Pop()

			vm.stack.Push(MakeNumber(left.Number() + right.Number()))
		}

		vm.ip++
	}

	if vm.stack.Length() == 1 {
		return vm.stack.Pop(), nil
	}
	return Value{}, fmt.Errorf("execution finished with invalid stack state")
}
