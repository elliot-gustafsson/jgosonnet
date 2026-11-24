package evaluator

import (
	"cmp"
	"path/filepath"
	"slices"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"github.com/google/go-jsonnet/ast"
)

const (
	// TemplateIndexThreshold is the field count above which a compile-time
	// LayerTemplate builds a DescriptorTable hash index once for O(1) lookups.
	TemplateIndexThreshold = 8

	// MaxInternLiteralLen caps the length of AST string literals that are
	// eagerly interned at compile time; longer string literals use a zero-copy
	// static string pointer without hashing or polluting the interner.
	MaxInternLiteralLen = 64
)

type IRKind uint8

const (
	IRKindConst IRKind = iota
	IRKindSelf
	IRKindVar
	IRKindVarUnresolved
	IRKindIndexConst
	IRKindIndex
	IRKindBinary
	IRKindUnary
	IRKindConditional
	IRKindApply
	IRKindLocal
	IRKindObject
	IRKindArray
	IRKindFunction
	IRKindSuperIndex
	IRKindInSuper
	IRKindImport
	IRKindImportStr
	IRKindImportBin
	IRKindError
	IRKindInvalidNumber
	IRKindGoCallback
)

type IRNode struct {
	Orig     ast.Node
	ConstVal Value
	Child0   *IRNode
	Child1   *IRNode
	Extra    unsafe.Pointer
	A        uint32
	B        uint32
	Kind     IRKind
	SubOp    uint8
}

type IRApply struct {
	Positional []*IRNode
	NamedKeys  []uint32
	NamedArgs  []*IRNode
	TailStrict bool
}

type IRLocal struct {
	Binds []*IRNode
	Body  *IRNode
}

type IRArray struct {
	Elements []*IRNode
}

type IRFunction struct {
	Body        *IRNode
	ParamKeyIds []uint32
	DefaultArgs []*IRNode
	RequiredCnt int
}

type IRObjectField struct {
	Name  *IRNode // nil when KeyId != 0 (static key)
	Body  *IRNode
	KeyId uint32
	Meta  uint8
}

type IRObject struct {
	StaticTemplate *LayerTemplate // non-nil when all fields have static string names
	Fields         []IRObjectField
	LocalNodes     []*IRNode
	Asserts        []*IRNode
}

type IRImport struct {
	Dir  string
	File string
	Abs  string
}

// LayerTemplate holds the immutable, compile-time-computed structure of a
// static object literal so every runtime instantiation is O(1).
type LayerTemplate struct {
	Keys         []uint32
	Nodes        []*IRNode
	Meta         []uint8
	Index        *utils.DescriptorTable
	LocalNodes   []*IRNode
	Asserts      []*IRNode
	VisibleCount int
	HasHidden    bool

	// Lazily computed, lock-free cached visible field orderings for fast
	// single-layer manifestation without sorting or allocating FieldPlans.
	sortedVisible  atomic.Pointer[[]uint32]
	naturalVisible atomic.Pointer[[]uint32]
}

func (t *LayerTemplate) GetSortedVisible(inter *interner.Interner) []uint32 {
	if p := t.sortedVisible.Load(); p != nil {
		return *p
	}
	order := make([]uint32, 0, t.VisibleCount)
	for f := range t.Keys {
		if ast.ObjectFieldHide(t.Meta[f]&MaskVisibility) != ast.ObjectFieldHidden {
			order = append(order, uint32(f))
		}
	}
	if len(order) > 1 {
		keys := t.Keys
		slices.SortFunc(order, func(a, b uint32) int {
			return cmp.Compare(inter.Get(keys[a]), inter.Get(keys[b]))
		})
	}
	t.sortedVisible.Store(&order)
	return order
}

func (t *LayerTemplate) GetNaturalVisible(inter *interner.Interner) []uint32 {
	if p := t.naturalVisible.Load(); p != nil {
		return *p
	}
	order := make([]uint32, 0, t.VisibleCount)
	for f := range t.Keys {
		if ast.ObjectFieldHide(t.Meta[f]&MaskVisibility) != ast.ObjectFieldHidden {
			order = append(order, uint32(f))
		}
	}
	if len(order) > 1 {
		keys := t.Keys
		slices.SortFunc(order, func(a, b uint32) int {
			return naturalStringSort(inter.Get(keys[a]), inter.Get(keys[b]))
		})
	}
	t.naturalVisible.Store(&order)
	return order
}

type compileScope struct {
	names  []ast.Identifier
	parent *compileScope
}

func (s *compileScope) resolve(name ast.Identifier) (upDepth uint32, slotIdx uint32, found bool) {
	depth := uint32(0)
	for curr := s; curr != nil; curr = curr.parent {
		for i := len(curr.names) - 1; i >= 0; i-- {
			if curr.names[i] == name {
				return depth, uint32(i), true
			}
		}
		depth++
	}
	return 0, 0, false
}

func (s *compileScope) isRootStdVar(name ast.Identifier) bool {
	if name != "$std" && name != "std" {
		return false
	}
	for curr := s; curr != nil; curr = curr.parent {
		for i := len(curr.names) - 1; i >= 0; i-- {
			if curr.names[i] == name {
				return curr == rootCompileScope
			}
		}
	}
	return false
}

var rootCompileScope = &compileScope{
	names: []ast.Identifier{"$std", "std"},
}

type StdLookupFunc func(name string) (Value, bool)

type astCompiler struct {
	inter     *interner.Interner
	stdLookup StdLookupFunc
	slab      []IRNode
}

func (c *astCompiler) allocNode() *IRNode {
	if len(c.slab) == 0 {
		c.slab = make([]IRNode, 64)
	}
	n := &c.slab[0]
	c.slab = c.slab[1:]
	return n
}

func CompileAST(root ast.Node, inter *interner.Interner) *IRNode {
	return CompileASTWithStd(root, inter, nil)
}

func CompileASTWithStd(root ast.Node, inter *interner.Interner, stdLookup StdLookupFunc) *IRNode {
	c := astCompiler{inter: inter, stdLookup: stdLookup}
	return c.compile(root, rootCompileScope)
}

func (c *astCompiler) compile(n ast.Node, scope *compileScope) *IRNode {
	ir := c.allocNode()
	ir.Orig = n

	switch node := n.(type) {
	case *ast.LiteralString:
		ir.Kind = IRKindConst
		if len(node.Value) > MaxInternLiteralLen {
			strPtr := &node.Value
			ir.Extra = unsafe.Pointer(strPtr)
			ir.ConstVal = MakeStaticString(strPtr)
			return ir
		}
		ir.ConstVal = MakeStringConst(c.inter.InternStatic(node.Value))
		return ir

	case *ast.LiteralNull:
		ir.Kind = IRKindConst
		ir.ConstVal = MakeNull()
		return ir

	case *ast.LiteralBoolean:
		ir.Kind = IRKindConst
		ir.ConstVal = MakeBool(node.Value)
		return ir

	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			ir.Kind = IRKindInvalidNumber
			return ir
		}
		ir.Kind = IRKindConst
		ir.ConstVal = MakeNumber(num)
		return ir

	case *ast.Self:
		ir.Kind = IRKindSelf
		return ir

	case *ast.Var:
		upDepth, slotIdx, found := scope.resolve(node.Id)
		if !found {
			ir.Kind = IRKindVarUnresolved
			return ir
		}
		ir.Kind = IRKindVar
		ir.A = upDepth
		ir.B = slotIdx
		return ir

	case *ast.Index:
		if ls, ok := node.Index.(*ast.LiteralString); ok {
			if c.stdLookup != nil && ls.Value != "thisFile" {
				if v, ok := node.Target.(*ast.Var); ok && scope.isRootStdVar(v.Id) {
					if stdVal, found := c.stdLookup(ls.Value); found {
						ir.Kind = IRKindConst
						ir.ConstVal = stdVal
						return ir
					}
				}
			}
			ir.Child0 = c.compile(node.Target, scope)
			ir.Kind = IRKindIndexConst
			ir.A = c.inter.InternStatic(ls.Value)
			return ir
		}
		ir.Child0 = c.compile(node.Target, scope)
		ir.Kind = IRKindIndex
		ir.Child1 = c.compile(node.Index, scope)
		return ir

	case *ast.Binary:
		ir.Kind = IRKindBinary
		ir.Child0 = c.compile(node.Left, scope)
		ir.Child1 = c.compile(node.Right, scope)
		ir.SubOp = uint8(node.Op)
		return ir

	case *ast.Unary:
		ir.Kind = IRKindUnary
		ir.Child0 = c.compile(node.Expr, scope)
		ir.SubOp = uint8(node.Op)
		return ir

	case *ast.Conditional:
		ir.Kind = IRKindConditional
		ir.Child0 = c.compile(node.Cond, scope)
		ir.Child1 = c.compile(node.BranchTrue, scope)
		ir.Extra = unsafe.Pointer(c.compile(node.BranchFalse, scope))
		return ir

	case *ast.Local:
		ir.Kind = IRKindLocal
		bindNames := make([]ast.Identifier, len(node.Binds))
		for i, b := range node.Binds {
			bindNames[i] = b.Variable
		}
		childScope := &compileScope{names: bindNames, parent: scope}

		binds := make([]*IRNode, len(node.Binds))
		for i, b := range node.Binds {
			binds[i] = c.compile(b.Body, childScope)
		}
		body := c.compile(node.Body, childScope)
		ir.Extra = unsafe.Pointer(&IRLocal{
			Binds: binds,
			Body:  body,
		})
		return ir

	case *ast.Apply:
		ir.Kind = IRKindApply
		ir.Child0 = c.compile(node.Target, scope)

		posCount := len(node.Arguments.Positional)
		nameCount := len(node.Arguments.Named)

		var pos []*IRNode
		if posCount > 0 {
			pos = make([]*IRNode, posCount)
			for i, a := range node.Arguments.Positional {
				pos[i] = c.compile(a.Expr, scope)
			}
		}

		var namedKeys []uint32
		var namedArgs []*IRNode
		if nameCount > 0 {
			namedKeys = make([]uint32, nameCount)
			namedArgs = make([]*IRNode, nameCount)
			for i, a := range node.Arguments.Named {
				namedKeys[i] = c.inter.InternStatic(string(a.Name))
				namedArgs[i] = c.compile(a.Arg, scope)
			}
		}

		ir.Extra = unsafe.Pointer(&IRApply{
			Positional: pos,
			NamedKeys:  namedKeys,
			NamedArgs:  namedArgs,
			TailStrict: node.TailStrict,
		})
		return ir

	case *ast.Array:
		ir.Kind = IRKindArray
		elems := make([]*IRNode, len(node.Elements))
		for i := range node.Elements {
			elems[i] = c.compile(node.Elements[i].Expr, scope)
		}
		ir.Extra = unsafe.Pointer(&IRArray{Elements: elems})
		return ir

	case *ast.Function:
		ir.Kind = IRKindFunction
		paramCount := len(node.Parameters)
		var paramKeyIds []uint32
		var defaultArgs []*IRNode
		requiredCnt := 0
		bodyScope := scope

		if paramCount > 0 {
			paramKeyIds = make([]uint32, paramCount)
			paramNames := make([]ast.Identifier, paramCount)
			hasDefaults := false
			for i, p := range node.Parameters {
				paramNames[i] = p.Name
				paramKeyIds[i] = c.inter.InternStatic(string(p.Name))
				if p.DefaultArg == nil {
					requiredCnt++
				} else {
					hasDefaults = true
				}
			}
			bodyScope = &compileScope{names: paramNames, parent: scope}
			if hasDefaults {
				defaultArgs = make([]*IRNode, paramCount)
				for i, p := range node.Parameters {
					if p.DefaultArg != nil {
						defaultArgs[i] = c.compile(p.DefaultArg, bodyScope)
					}
				}
			}
		}

		body := c.compile(node.Body, bodyScope)
		ir.Extra = unsafe.Pointer(&IRFunction{
			Body:        body,
			ParamKeyIds: paramKeyIds,
			DefaultArgs: defaultArgs,
			RequiredCnt: requiredCnt,
		})
		return ir

	case *ast.DesugaredObject:
		ir.Kind = IRKindObject
		innerScope := scope
		var localNodes []*IRNode
		if len(node.Locals) > 0 {
			localNames := make([]ast.Identifier, len(node.Locals))
			for i, l := range node.Locals {
				localNames[i] = l.Variable
			}
			innerScope = &compileScope{names: localNames, parent: scope}
			localNodes = make([]*IRNode, len(node.Locals))
			for i, l := range node.Locals {
				localNodes[i] = c.compile(l.Body, innerScope)
			}
		}

		var asserts []*IRNode
		if len(node.Asserts) > 0 {
			asserts = make([]*IRNode, len(node.Asserts))
			for i, a := range node.Asserts {
				asserts[i] = c.compile(a, innerScope)
			}
		}

		allStatic := true
		for i := range node.Fields {
			if _, ok := node.Fields[i].Name.(*ast.LiteralString); !ok {
				allStatic = false
				break
			}
		}

		fieldCount := len(node.Fields)
		if allStatic {
			var keys []uint32
			var nodes []*IRNode
			var meta []uint8
			visibleCount := 0

			if fieldCount > 0 {
				keys = make([]uint32, fieldCount)
				nodes = make([]*IRNode, fieldCount)
				meta = make([]uint8, fieldCount)
				for i, f := range node.Fields {
					ls := f.Name.(*ast.LiteralString)
					keys[i] = c.inter.InternStatic(ls.Value)
					nodes[i] = c.compile(f.Body, innerScope)
					meta[i] = CreateFieldMeta(f.Hide, f.PlusSuper)
					if f.Hide != ast.ObjectFieldHidden {
						visibleCount++
					}
				}
			}

			var idx *utils.DescriptorTable
			if fieldCount > TemplateIndexThreshold {
				idx = utils.NewHeapDescriptorTable(keys)
			}

			tmpl := &LayerTemplate{
				Keys:         keys,
				Nodes:        nodes,
				Meta:         meta,
				Index:        idx,
				LocalNodes:   localNodes,
				Asserts:      asserts,
				VisibleCount: visibleCount,
				HasHidden:    visibleCount < fieldCount,
			}
			ir.Extra = unsafe.Pointer(&IRObject{StaticTemplate: tmpl})
			return ir
		}

		fields := make([]IRObjectField, fieldCount)
		for i, f := range node.Fields {
			m := CreateFieldMeta(f.Hide, f.PlusSuper)
			body := c.compile(f.Body, innerScope)
			if ls, ok := f.Name.(*ast.LiteralString); ok {
				fields[i] = IRObjectField{
					KeyId: c.inter.InternStatic(ls.Value),
					Body:  body,
					Meta:  m,
				}
			} else {
				// Dynamic field names are evaluated in the OUTER scope!
				fields[i] = IRObjectField{
					Name: c.compile(f.Name, scope),
					Body: body,
					Meta: m,
				}
			}
		}

		ir.Extra = unsafe.Pointer(&IRObject{
			Fields:     fields,
			LocalNodes: localNodes,
			Asserts:    asserts,
		})
		return ir

	case *ast.SuperIndex:
		ir.Kind = IRKindSuperIndex
		if ls, ok := node.Index.(*ast.LiteralString); ok {
			ir.A = c.inter.InternStatic(ls.Value)
			return ir
		}
		ir.Child0 = c.compile(node.Index, scope)
		return ir

	case *ast.InSuper:
		ir.Kind = IRKindInSuper
		if ls, ok := node.Index.(*ast.LiteralString); ok {
			ir.A = c.inter.InternStatic(ls.Value)
			return ir
		}
		ir.Child0 = c.compile(node.Index, scope)
		return ir

	case *ast.Import:
		ir.Kind = IRKindImport
		ir.Extra = unsafe.Pointer(makeIRImport(node.Loc().FileName, node.File.Value))
		return ir

	case *ast.ImportStr:
		ir.Kind = IRKindImportStr
		ir.Extra = unsafe.Pointer(makeIRImport(node.Loc().FileName, node.File.Value))
		return ir

	case *ast.ImportBin:
		ir.Kind = IRKindImportBin
		ir.Extra = unsafe.Pointer(makeIRImport(node.Loc().FileName, node.File.Value))
		return ir

	case *ast.Error:
		ir.Kind = IRKindError
		ir.Child0 = c.compile(node.Expr, scope)
		return ir

	default:
		ir.Kind = IRKindError
		return ir
	}
}

func makeIRImport(fromFile, importedPath string) *IRImport {
	dir := filepath.Dir(fromFile)
	var absPath string
	if importedPath != "" {
		if filepath.IsAbs(importedPath) {
			absPath = importedPath
		} else {
			absPath = filepath.Join(dir, importedPath)
		}
	}
	return &IRImport{
		Dir:  dir,
		File: importedPath,
		Abs:  absPath,
	}
}
