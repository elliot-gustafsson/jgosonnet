package evaluator

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/google/go-jsonnet/ast"
)

type OpCode uint8

const (
	OpReturn OpCode = iota
	OpJump

	OpPushNull
	OpPushString
	OpPushNumber
	OpPushFalse
	OpPushTrue

	OpAdd

	OpCall

	OpLocalGet
	OpLocalSet

	OpPushScope
	OpPopScope

	OpJumpIfFalse

	OpMakeThunk
	OpMakeObject
	OpMakeArray
	OpMakeFunction

	OpMetaObjectLocals
	OpMetaObjectAsserts

	OpBindDefaultArg
	OpFuncParamData
)

type GlobalInterner struct {
	lock    sync.RWMutex
	mapping map[string]uint32
	strings []string
}

func NewGlobalInterner() *GlobalInterner {
	return &GlobalInterner{
		mapping: make(map[string]uint32, 8192),
		strings: make([]string, 0, 8192),
	}
}

func (i *GlobalInterner) Intern(s string) uint32 {
	i.lock.RLock()
	defer i.lock.RUnlock()

	if id, ok := i.mapping[s]; ok {
		return id
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	id := uint32(len(i.strings))
	i.strings = append(i.strings, s)
	i.mapping[s] = id

	return id
}

func (i *GlobalInterner) Get(id uint32) string {
	i.lock.RLock()
	defer i.lock.RUnlock()

	if id >= uint32(len(i.strings)) {
		return ""
	}
	return i.strings[id]
}

func (i *GlobalInterner) Reset() {
	i.lock.Lock()
	defer i.lock.Unlock()

	clear(i.mapping)
	i.strings = i.strings[:0]
}

type ProgramCache struct {
	lock  sync.RWMutex
	cache map[string]*Program
}

type GlobalContext struct {
	Interner     *Interner
	ProgramCache *ProgramCache
}

type Program struct {
	Instructions []Instruction
	Strings      []string
	Numbers      []float64

	// Templates *arena.Arena[Layer]

	// Interner    *Interner
	// Registry    *Registry
	// Environment *Environment
}

type CompilerVar struct {
	Name       string
	ScopeDepth uint32
	SlotIndex  uint32
}

type Compiler struct {
	Interner     *Interner
	ProgramCache *ProgramCache

	Instructions []Instruction
	Strings      []string
	Numbers      []float64

	Locals     []CompilerVar
	ScopeDepth uint32
}

type Instruction struct {
	op OpCode

	operand2 uint16
	operand  uint32
}

func NewCompiler(ctx *GlobalContext) *Compiler {
	return &Compiler{
		Interner:     ctx.Interner,
		ProgramCache: ctx.ProgramCache,
	}
}

func (c *Compiler) Compile(node ast.Node) (*Program, error) {
	err := c.visit(node)
	if err != nil {
		return nil, err
	}

	c.emit(OpReturn, 0)

	return &Program{
		Instructions: c.Instructions,
		Strings:      c.Strings,
		Numbers:      c.Numbers,
	}, nil
}

func (c *Compiler) visit(n ast.Node) error {
	switch node := n.(type) {
	default:
		return fmt.Errorf("unsupported AST node type: %T", n)
	case *ast.LiteralNull:
		c.emit(OpPushNull, 0)

	case *ast.LiteralString:
		id := uint32(len(c.Strings))
		c.Strings = append(c.Strings, node.Value)
		c.emit(OpPushString, id)

	case *ast.LiteralBoolean:
		if node.Value {
			c.emit(OpPushTrue, 0)
		} else {
			c.emit(OpPushFalse, 0)
		}

	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		id := uint32(len(c.Numbers))
		c.Numbers = append(c.Numbers, num)
		c.emit(OpPushNumber, id)

	case *ast.DesugaredObject:
		c.makeObject(node)

	case *ast.Array:
		for i := range node.Elements {
			err := c.visitLazy(node.Elements[i].Expr)
			if err != nil {
				return err
			}
		}
		c.emit(OpMakeArray, uint32(len(node.Elements)))

	case *ast.Local:

		// TODO: Dont use scopes, we dont want to scope travese during exec...

		c.ScopeDepth++
		c.emit(OpPushScope, uint32(len(node.Binds)))

		startLocalsLen := len(c.Locals)

		for i := range node.Binds {
			c.Locals = append(c.Locals, CompilerVar{
				Name:       string(node.Binds[i].Variable),
				ScopeDepth: c.ScopeDepth,
				SlotIndex:  uint32(i),
			})
		}

		for i := range node.Binds {
			c.visitLazy(node.Binds[i].Body)
			c.emit(OpLocalSet, uint32(i))
		}

		c.visit(node.Body)

		c.emit(OpPopScope, 0)

		c.Locals = c.Locals[:startLocalsLen]
		c.ScopeDepth--

	case *ast.Var:
		targetName := string(node.Id)

		for i := len(c.Locals) - 1; i >= 0; i-- {
			if c.Locals[i].Name == targetName {

				// 2. Calculate the distance for the VM to jump
				depthDiff := uint16(c.ScopeDepth - c.Locals[i].ScopeDepth)
				slotIndex := c.Locals[i].SlotIndex

				// 3. Emit the instruction: OpLocalGet <slotIndex>, <depthDiff>
				c.emit2(OpLocalGet, slotIndex, depthDiff)
				return nil
			}
		}
		return fmt.Errorf("undefined variable: %s", targetName)

	case *ast.Function:
		paramCount := len(node.Parameters)

		jumpIdx := len(c.Instructions)
		c.emit(OpJump, 0) // placeholder

		// the functions bytecode starts here
		funcIP := uint32(len(c.Instructions))

		c.ScopeDepth++
		startLocalsLen := len(c.Locals)

		for i, p := range node.Parameters {
			c.Locals = append(c.Locals, CompilerVar{
				Name:       string(p.Name),
				ScopeDepth: c.ScopeDepth,
				SlotIndex:  uint32(i),
			})
		}

		for i, p := range node.Parameters {
			if p.DefaultArg != nil {
				c.makeThunk(p.DefaultArg)
				c.emit(OpBindDefaultArg, uint32(i))
			}
		}

		// The body is evaluated in tail position!
		c.visit(node.Body) // true (tail eval)

		c.emit(OpPopScope, 0)
		c.emit(OpReturn, 0)

		c.Locals = c.Locals[:startLocalsLen]
		c.ScopeDepth--

		c.Instructions[jumpIdx].operand = uint32(len(c.Instructions) - (jumpIdx + 1))

		c.emit2(OpMakeFunction, funcIP, uint16(paramCount))

		// 8. Inline Data Stream (For Named Arguments)
		// We emit the interned string IDs for the parameter names so the VM
		// knows how to match named arguments when the function is called!
		for _, p := range node.Parameters {
			nameId := c.Interner.Intern(string(p.Name))
			c.emit(OpFuncParamData, nameId)
		}

	case *ast.Apply:
		c.visit(node.Target)

		for i := range node.Arguments.Positional {
			err := c.visitLazy(node.Arguments.Positional[i].Expr)
			if err != nil {
				return err
			}
		}
		for i := range node.Arguments.Named {
			err := c.visitLazy(node.Arguments.Named[i].Arg)
			if err != nil {
				return err
			}
		}

		c.emit2(OpCall, uint32(len(node.Arguments.Positional)), uint16(len(node.Arguments.Named)))

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

func (c *Compiler) visitLazy(n ast.Node) error {
	switch node := n.(type) {
	default:
		c.makeThunk(node)

	case *ast.LiteralNull:
		c.emit(OpPushNull, 0)

	case *ast.LiteralString:
		id := uint32(len(c.Strings))
		c.Strings = append(c.Strings, node.Value)
		c.emit(OpPushString, id)

	case *ast.LiteralBoolean:
		if node.Value {
			c.emit(OpPushTrue, 0)
		} else {
			c.emit(OpPushFalse, 0)
		}

	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		id := uint32(len(c.Numbers))
		c.Numbers = append(c.Numbers, num)
		c.emit(OpPushNumber, id)

	case *ast.Self:
		// TODO: Fix this
		// c.pushConst(c.prog.)
	}

	return nil
}

func (c *Compiler) emit(op OpCode, operand uint32) {
	c.emit2(op, operand, 0)
}

func (c *Compiler) emit2(op OpCode, operand uint32, operand2 uint16) {
	i := Instruction{
		op:       op,
		operand:  operand,
		operand2: operand2,
	}
	c.Instructions = append(c.Instructions, i)
}

func (c *Compiler) makeThunk(node ast.Node) {
	thunkIp := len(c.Instructions)

	c.emit(OpMakeThunk, 0)

	c.visit(node)

	c.emit(OpReturn, 0)

	bodyLength := len(c.Instructions) - (thunkIp + 1)

	c.Instructions[thunkIp].operand = uint32(bodyLength)
}

const (
	FlagObjectHasLocals  uint16 = 1 << 0 // 1
	FlagObjectHasAsserts uint16 = 1 << 1 // 2
)

func (c *Compiler) makeObject(node *ast.DesugaredObject) {
	numFields := uint32(len(node.Fields))
	numLocals := uint32(len(node.Locals))
	numAsserts := uint32(len(node.Asserts))

	var flags uint16
	if numLocals > 0 {
		flags |= FlagObjectHasLocals
	}
	if numAsserts > 0 {
		flags |= FlagObjectHasAsserts
	}

	c.emit2(OpMakeObject, numFields, flags)

	if numLocals > 0 {
		c.emit(OpMetaObjectLocals, numLocals)
	}
	if numAsserts > 0 {
		c.emit(OpMetaObjectAsserts, numAsserts)
	}

}
