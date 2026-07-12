package evaluator

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/google/go-jsonnet/ast"
)

type OpCode uint8

const (
	OpNone OpCode = iota

	OpReturn
	OpForce
	OpJump

	OpPushNull
	OpPushString
	OpPushNumber
	OpPushFalse
	OpPushTrue

	// binary ops
	OpPlus
	OpMinus
	OpDiv
	OpMult
	OpBitwiseAnd
	OpBitwiseOr
	OpBitwiseXor
	OpShiftL
	OpShiftR
	OpEqual
	OpUnequal
	OpGreater
	OpGreaterEq
	OpLess
	OpLessEq

	OpIndex
	OpCall
	OpParamMeta

	OpJumpIfFalse

	OpAllocLocals
	OpPopLocals

	OpLocalSet
	OpLocalGet

	OpMakeClosure

	OpAllocateParams
	OpParamSet

	OpMakeThunk
	OpMakeObject
	OpMakeArray
	OpMakeFunction

	OpBindArgs
	OpPushNameId

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
}

type paramMeta struct {
	nameId     uint32
	hasDefault uint16
	defaultIp  uint32
}

type CompilerFrame struct {
	LocalBase int // The index in c.Locals where this frame's variables begin
}

type Compiler struct {
	Interner     *Interner
	ProgramCache *ProgramCache

	Instructions []Instruction
	Strings      []string
	Numbers      []float64

	Locals []uint32
	Frames []CompilerFrame

	paramMetaBuf []paramMeta
}

type Instruction struct {
	op OpCode

	operand2 uint16
	operand  uint32
}

func NewCompiler(ctx *GlobalContext) *Compiler {
	c := &Compiler{
		Interner:     ctx.Interner,
		ProgramCache: ctx.ProgramCache,
		// Frames:       make([]CompilerFrame, 0, 128),
		// Locals:       make([]string, 0, 1024),
		// Upvalues:     make([]CompilerUpvalue, 0, 1024),

		Frames:       make([]CompilerFrame, 0, 128),
		paramMetaBuf: make([]paramMeta, 0, 16),
	}

	c.Frames = append(c.Frames, CompilerFrame{
		LocalBase: 0,
	})

	return c
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
		bindCount := uint32(len(node.Binds))
		localsIdx := len(c.Locals)

		for i := range node.Binds {
			nid := c.Interner.Intern(string(node.Binds[i].Variable))
			c.Locals = append(c.Locals, nid)
		}

		c.emit(OpAllocLocals, bindCount)

		for i := range node.Binds {
			err := c.visitLazy(node.Binds[i].Body)
			if err != nil {
				return err
			}
			slotIndex := uint32(localsIdx+i) - uint32(c.currentFrame().LocalBase)
			c.emit(OpLocalSet, slotIndex)
		}

		err := c.visit(node.Body)
		if err != nil {
			return err
		}

		c.Locals = c.Locals[:localsIdx]

		c.emit(OpPopLocals, bindCount)

	case *ast.Var:
		tid := c.Interner.Intern(string(node.Id))
		frame := c.currentFrame()

		for i := len(c.Locals) - 1; i >= frame.LocalBase; i-- {
			if c.Locals[i] != tid {
				continue
			}

			slotIndex := uint32(i - frame.LocalBase)
			c.emit(OpLocalGet, slotIndex)
			c.emit(OpForce, 0)
			return nil
		}
		return fmt.Errorf("variable not found in current function: %s", node.Id)

	case *ast.Binary:
		if err := c.visit(node.Left); err != nil {
			return err
		}
		c.emit(OpForce, 0)

		if err := c.visit(node.Right); err != nil {
			return err
		}
		c.emit(OpForce, 0)

		// TODO: handle literals directly at compile time

		switch node.Op {
		default:
			return fmt.Errorf("unsupported binary operator: %s", node.Op.String())
		case ast.BopPlus:
			c.emit(OpPlus, 0)
		case ast.BopMinus:
			c.emit(OpMinus, 0)
		case ast.BopDiv:
			c.emit(OpDiv, 0)
		case ast.BopMult:
			c.emit(OpMult, 0)
		case ast.BopBitwiseAnd:
			c.emit(OpBitwiseAnd, 0)
		case ast.BopBitwiseOr:
			c.emit(OpBitwiseOr, 0)
		case ast.BopBitwiseXor:
			c.emit(OpBitwiseXor, 0)
		case ast.BopShiftL:
			c.emit(OpShiftL, 0)
		case ast.BopShiftR:
			c.emit(OpShiftR, 0)
		case ast.BopManifestEqual:
			c.emit(OpEqual, 0)
		case ast.BopManifestUnequal:
			c.emit(OpUnequal, 0)
		case ast.BopGreater:
			c.emit(OpGreater, 0)
		case ast.BopGreaterEq:
			c.emit(OpGreaterEq, 0)
		case ast.BopLess:
			c.emit(OpLess, 0)
		case ast.BopLessEq:
			c.emit(OpLessEq, 0)
		}

	case *ast.Function:
		jumpIdx := len(c.Instructions)
		c.emit(OpJump, 0) // Placeholder to jump over body and defaults

		// funcIp := uint32(len(c.Instructions))

		localsLen := len(c.Locals)
		frameIdx := len(c.Frames)
		c.Frames = append(c.Frames, CompilerFrame{LocalBase: localsLen})

		// Clear and reuse compiler buffer
		c.paramMetaBuf = c.paramMetaBuf[:0]

		for i := range node.Parameters {
			nid := c.Interner.Intern(string(node.Parameters[i].Name))
			c.Locals = append(c.Locals, nid)

			meta := paramMeta{nameId: nid}

			if node.Parameters[i].DefaultArg != nil {
				meta.hasDefault = 1
				meta.defaultIp = uint32(len(c.Instructions))

				// Thunks act as isolated executable blocks ending in Return
				c.visit(node.Parameters[i].DefaultArg)
				c.emit(OpReturn, 0)
			}
			c.paramMetaBuf = append(c.paramMetaBuf, meta)
		}

		bodyIp := uint32(len(c.Instructions))
		err := c.visit(node.Body)
		if err != nil {
			return err
		}
		c.emit(OpReturn, 0)

		// Scope cleanup
		c.Locals = c.Locals[:localsLen]
		c.Frames = c.Frames[:frameIdx]

		bodyLength := uint32(len(c.Instructions) - (jumpIdx + 1))
		c.patchOperand(jumpIdx, bodyLength)

		c.emit2(OpMakeFunction, bodyIp, uint16(len(node.Parameters)))

		// Emit Inline Metadata (2 instructions per parameter)
		for _, meta := range c.paramMetaBuf {
			c.emit2(OpParamMeta, meta.nameId, meta.hasDefault)
			c.emit(OpParamMeta, meta.defaultIp)
		}

	case *ast.Index:
		if err := c.visit(node.Target); err != nil {
			return err
		}
		c.emit(OpForce, 0)

		if err := c.visit(node.Index); err != nil {
			return err
		}
		c.emit(OpForce, 0)

		c.emit(OpIndex, 0)
		c.emit(OpForce, 0)

	case *ast.Apply:
		if err := c.visit(node.Target); err != nil {
			return err
		}

		// Positional arguments (Thunks pushed to stack)
		for i := range node.Arguments.Positional {
			if err := c.visitLazy(node.Arguments.Positional[i].Expr); err != nil {
				return err
			}
		}

		// Named arguments (Thunks pushed to stack)
		for i := range node.Arguments.Named {
			if err := c.visitLazy(node.Arguments.Named[i].Arg); err != nil {
				return err
			}
		}

		numPos := uint32(len(node.Arguments.Positional))
		numNamed := uint16(len(node.Arguments.Named))

		c.emit2(OpCall, numPos, numNamed)

		// Emit Inline Metadata for caller
		for i := range node.Arguments.Named {
			name := node.Arguments.Named[i].Name
			nameId := c.Interner.Intern(string(name))
			c.emit(OpParamMeta, nameId)
		}

	}

	return nil
}

func (c *Compiler) visitLazy(n ast.Node) error {
	switch node := n.(type) {
	default:
		err := c.makeThunk(node)
		if err != nil {
			return err
		}

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

func (c *Compiler) patchOperand(id int, operand uint32) {
	c.Instructions[id].operand = operand
}

func (c *Compiler) patchOperand2(id int, operand2 uint16) {
	c.Instructions[id].operand2 = operand2
}

func (c *Compiler) currentFrame() *CompilerFrame {
	return &c.Frames[len(c.Frames)-1]
}

func (c *Compiler) makeThunk(node ast.Node) error {
	thunkIp := len(c.Instructions)

	c.emit(OpMakeThunk, 0)

	err := c.visit(node)
	if err != nil {
		return err
	}

	c.emit(OpReturn, 0)

	bodyLength := len(c.Instructions) - (thunkIp + 1)

	c.Instructions[thunkIp].operand = uint32(bodyLength)
	return nil
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

// case *ast.Function:
// 	paramCount := len(node.Parameters)

// 	jumpIdx := len(c.Instructions)
// 	c.emit(OpJump, 0) // placeholder

// 	// the functions bytecode starts here
// 	funcIP := uint32(len(c.Instructions))

// 	c.ScopeDepth++
// 	startLocalsLen := len(c.Locals)

// 	for i, p := range node.Parameters {
// 		c.Locals = append(c.Locals, CompilerVar{
// 			Name:       string(p.Name),
// 			ScopeDepth: c.ScopeDepth,
// 			SlotIndex:  uint32(i),
// 		})
// 	}

// 	for i, p := range node.Parameters {
// 		if p.DefaultArg != nil {
// 			c.makeThunk(p.DefaultArg)
// 			c.emit(OpBindDefaultArg, uint32(i))
// 		}
// 	}

// 	// The body is evaluated in tail position!
// 	c.visit(node.Body) // true (tail eval)

// 	c.emit(OpPopScope, 0)
// 	c.emit(OpReturn, 0)

// 	c.Locals = c.Locals[:startLocalsLen]
// 	c.ScopeDepth--

// 	c.patchOperand(jumpIdx, uint32(len(c.Instructions)-(jumpIdx+1)))

// 	c.emit2(OpMakeFunction, funcIP, uint16(paramCount))

// 	// 8. Inline Data Stream (For Named Arguments)
// 	// We emit the interned string IDs for the parameter names so the VM
// 	// knows how to match named arguments when the function is called!
// 	for _, p := range node.Parameters {
// 		nameId := c.Interner.Intern(string(p.Name))
// 		c.emit(OpFuncParamData, nameId)
// 	}

// case *ast.Apply:
// 	c.visit(node.Target)

// 	for i := range node.Arguments.Positional {
// 		err := c.visitLazy(node.Arguments.Positional[i].Expr)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	for i := range node.Arguments.Named {
// 		err := c.visitLazy(node.Arguments.Named[i].Arg)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	c.emit2(OpCall, uint32(len(node.Arguments.Positional)), uint16(len(node.Arguments.Named)))

// ---------------------------------------
// locals

// frame := &c.Frames[c.FrameIdx]
// prevLocalCount := frame.LocalCount

// // Register the locals
// for i := range node.Binds {
// 	c.Locals = append(c.Locals, string(node.Binds[i].Variable))
// 	frame.LocalCount++
// }

// for i := range node.Binds {
// 	slotIndex := prevLocalCount + uint32(i)

// 	c.pushFrame()
// 	closureFrame := c.getFrame()

// 	initIdx := len(c.Instructions)
// 	c.emit(OpInitThunk, slotIndex)

// 	jumpIdx := initIdx + 1
// 	c.emit(OpJump, 0) // placeholder

// 	bodyStart := jumpIdx + 1

// 	c.visit(node.Binds[i].Body)
// 	c.emit(OpReturn, 0)

// 	c.patchOperand(jumpIdx, uint32(len(c.Instructions)-bodyStart))

// 	c.patchOperand2(initIdx, uint16(closureFrame.UpvalCount))

// 	c.popFrame()
// }

// c.visit(node.Body)

// // 4. Pop locals
// c.emit(OpPopLocals, bindCount)

// frame = &c.Frames[c.FrameIdx]
// frame.LocalCount = prevLocalCount
// c.Locals = c.Locals[:frame.LocalStart+prevLocalCount]

// -----------------------------
// var

// for i := len(c.Locals) - 1; i >= 0; i-- {
// 	if c.Locals[i].Name == targetName {

// 		// 2. Calculate the distance for the VM to jump
// 		depthDiff := uint16(c.ScopeDepth - c.Locals[i].ScopeDepth)
// 		slotIndex := c.Locals[i].SlotIndex

// 		// 3. Emit the instruction: OpLocalGet <slotIndex>, <depthDiff>
// 		c.emit2(OpLocalGet, slotIndex, depthDiff)
// 		return nil
// 	}
// }
// return fmt.Errorf("undefined variable: %s", targetName)
