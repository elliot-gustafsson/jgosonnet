package evaluator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/utils"
	"github.com/google/go-jsonnet/ast"
)

func EvaluateNode(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := evaluateNode(n, scopePtr, ctx)
	if err != nil {
		return ValueNone, WrapError(err, n)
	}
	if val.IsThunk() {
		val, err = val.Eval(ctx)
		if err != nil {
			return ValueNone, WrapError(err, n)
		}
	}
	return val, nil
}

func ManifestValue(value Value, ctx Context) (any, error) {

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type '%s'", value.Type().String())
	case ValueTypeNull:
		return nil, nil
	case ValueTypeString:
		s := value.String(ctx)
		return strings.Clone(s), nil
	case ValueTypeNumber:
		n := value.Number()
		if float64(int64(n)) == n {
			return int64(n), nil
		}
		return n, nil
	case ValueTypeBool:
		return value.Bool(), nil
	case ValueTypeObject:
		subCtx := ctx
		subCtx.Self = value
		return manifestObject(value.Object(ctx), subCtx)
	case ValueTypeArray:
		res := make([]any, 0, len(value.Array(ctx)))
		for _, v := range value.Array(ctx) {
			ev, err := ManifestValue(v, ctx)
			if err != nil {
				return nil, err
			}
			res = append(res, ev)
		}
		return res, nil
	case ValueTypeFunction:
		res, err := value.Function(ctx).Exec(nil, ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(res, ctx)
	case ValueTypeThunk:
		v, err := value.Eval(ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(v, ctx)
	}
}

func CreateFileScope(filename string, baseStd Value, ctx Context) uintptr {
	allocator := ctx.State.Registry.Allocator

	keyId := ctx.State.Interner.Intern("thisFile")

	layer := arena.Create[Layer](allocator)
	layer.Keys = arena.Alloc[uint32](allocator, 1)
	layer.Keys[0] = keyId

	layer.Values = arena.Alloc[Value](allocator, 1)
	layer.Values[0] = MakeString(filename, ctx)

	layer.Meta = arena.Alloc[uint8](allocator, 1)

	fileObj := NewSingleLayerObject(allocator, layer)

	mergedObjId := MergeObjects(baseStd.Payload(), uintptr(unsafe.Pointer(fileObj)), ctx)
	fileStd := MakeObjectValue(mergedObjId)

	s, scopePtr := ctx.NewScope(2, 2)

	s.Bindings[0] = NamedValue{Key: ctx.State.Interner.Intern("$std"), Value: fileStd}
	s.Bindings[1] = NamedValue{Key: ctx.State.Interner.Intern("std"), Value: fileStd}

	return scopePtr
}

func evaluateNodeLazy(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	switch node := n.(type) {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %T (lazy eval)", node)
	case *ast.LiteralString:
		// id := ctx.State.Interner.Intern(node.Value)
		// return MakeStringConst(id), nil
		return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return ValueNone, fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.Self:
		return ctx.Self, nil

	case *ast.DesugaredObject:
		return NewThunk(ThunkTypeObject, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Array:
		return NewThunk(ThunkTypeArray, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Local:
		return NewThunk(ThunkTypeLocal, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Apply:
		return NewThunk(ThunkTypeApply, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Index:
		return NewThunk(ThunkTypeIndex, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Var:
		return NewThunk(ThunkTypeVar, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Function:
		return NewThunk(ThunkTypeFunction, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Conditional:
		return NewThunk(ThunkTypeConditional, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Binary:
		return NewThunk(ThunkTypeBinary, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Unary:
		return NewThunk(ThunkTypeUnary, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Import:
		return NewThunk(ThunkTypeImport, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.ImportStr:
		return NewThunk(ThunkTypeImportStr, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.SuperIndex:
		return NewThunk(ThunkTypeSuperIndex, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.InSuper:
		return NewThunk(ThunkTypeInSuper, unsafe.Pointer(node), scopePtr, ctx), nil
	case *ast.Error:
		return NewThunk(ThunkTypeError, unsafe.Pointer(node), scopePtr, ctx), nil

	case *GoCallbackNode:
		return NewThunk(ThunkTypeGoCallback, unsafe.Pointer(node), scopePtr, ctx), nil
	}
}

func evaluateNode(n ast.Node, scopePtr uintptr, ctx Context) (Value, error) {
	switch node := n.(type) {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %T", node)
	case *ast.LiteralString:
		// id := ctx.State.Interner.Intern(node.Value)
		// return MakeStringConst(id), nil
		return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return ValueNone, fmt.Errorf("(%T) failed to parse float val (%s), err: %w", node, node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.DesugaredObject:
		return handleDesugaredObject(node, scopePtr, ctx)
	case *ast.Array:
		return handleArray(node, scopePtr, ctx)
	case *ast.Local:
		return handleLocal(node, scopePtr, ctx)
	case *ast.Apply:
		return handleApply(node, scopePtr, ctx)
	case *ast.Index:
		return handleIndex(node, scopePtr, ctx)
	case *ast.Var:
		return handleVar(node, scopePtr, ctx)
	case *ast.Function:
		return handleFunction(node, scopePtr, ctx)
	case *ast.Conditional:
		return handleConditional(node, scopePtr, ctx)
	case *ast.Binary:
		return handleBinary(node, scopePtr, ctx)
	case *ast.Unary:
		return handleUnary(node, scopePtr, ctx)
	case *ast.Import:
		return handleImport(node, ctx)
	case *ast.ImportStr:
		return handleImportStr(node, ctx)
	case *ast.Self:
		return ctx.Self, nil
	case *ast.SuperIndex:
		return handleSuperIndex(node, scopePtr, ctx)
	case *ast.InSuper:
		return handleInSuper(node, scopePtr, ctx)
	case *ast.Error:
		return handleError(node, scopePtr, ctx)
	case *GoCallbackNode:
		return node.Func.Exec(node.Args, ctx)
	}
}

type GoCallbackNode struct {
	Func Function
	Args []NamedValue
}

func (n *GoCallbackNode) Context() ast.Context             { return nil }
func (n *GoCallbackNode) Loc() *ast.LocationRange          { return nil }
func (n *GoCallbackNode) FreeVariables() ast.Identifiers   { return nil }
func (n *GoCallbackNode) SetFreeVariables(ast.Identifiers) {}
func (n *GoCallbackNode) SetContext(ast.Context)           {}
func (n *GoCallbackNode) OpenFodder() *ast.Fodder          { return nil }

func handleArray(node *ast.Array, scopePtr uintptr, ctx Context) (Value, error) {
	arr, val := MakeArraySized(len(node.Elements), ctx)
	for i := range node.Elements {
		ev, err := evaluateNodeLazy(node.Elements[i].Expr, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		arr[i] = ev
	}
	return val, nil
}

func handleDesugaredObject(node *ast.DesugaredObject, scopePtr uintptr, ctx Context) (Value, error) {

	allocator := ctx.State.Registry.Allocator

	fieldCount := len(node.Fields)
	localsCount := len(node.Locals)

	layer := arena.Create[Layer](allocator)

	layer.ParentScopePtr = scopePtr

	if fieldCount > 0 {
		layer.Keys = arena.Alloc[uint32](allocator, fieldCount)
		layer.Nodes = arena.Alloc[ast.Node](allocator, fieldCount)
		layer.Meta = arena.Alloc[uint8](allocator, fieldCount)
	}

	if localsCount > 0 {
		layer.LocalKeys = arena.Alloc[uint32](allocator, localsCount)
		layer.LocalNodes = arena.Alloc[ast.Node](allocator, localsCount)
	}

	if len(node.Asserts) > 0 {
		layer.packAsserts(node.Asserts)
	}

	for i, v := range node.Locals {
		layer.LocalKeys[i] = ctx.State.Interner.Intern(string(v.Variable))
		layer.LocalNodes[i] = v.Body
	}

	useMap := fieldCount > MaxLayerLinearKeys
	if useMap {
		layer.Index = utils.NewEmptyDescriptorTable(allocator, fieldCount)
	}

	index := 0
	for _, v := range node.Fields {
		var nameStr string

		// Fast path literal string, most keys are just static strings
		if ls, ok := v.Name.(*ast.LiteralString); ok {

			nameStr = ls.Value
		} else {
			var err error
			name, err := EvaluateNode(v.Name, scopePtr, ctx)
			if err != nil {
				return ValueNone, err
			}

			if name.IsNull() {
				// Omitted field
				continue
			}

			if !name.IsString() {
				return ValueNone, fmt.Errorf("unexpected field name type %s, expected string", name.Type().String())
			}

			nameStr = name.String(ctx)
		}

		keyId := ctx.State.Interner.Intern(nameStr)

		layer.Keys[index] = keyId
		layer.Nodes[index] = v.Body
		layer.Meta[index] = CreateFieldMeta(v.Hide, v.PlusSuper)

		if useMap {
			layer.Index.Append(keyId)
		}
		index++
	}

	if index < fieldCount {
		layer.Keys = layer.Keys[:index]
		layer.Nodes = layer.Nodes[:index]
		layer.Meta = layer.Meta[:index]
	}

	obj := NewSingleLayerObject(allocator, layer)

	return MakeObjectValue(obj), nil
}

func handleLocal(node *ast.Local, scopePtr uintptr, ctx Context) (Value, error) {

	s, childScopePtr := ctx.NewScope(scopePtr, len(node.Binds))

	for i, v := range node.Binds {
		vname := string(v.Variable)
		keyId := ctx.State.Interner.Intern(vname)
		t, err := evaluateNodeLazy(v.Body, childScopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		s.Bindings[i] = NamedValue{keyId, t}
	}

	val, err := EvaluateNode(node.Body, childScopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	return val, nil
}

func handleApply(node *ast.Apply, scopePtr uintptr, ctx Context) (Value, error) {
	val, err := EvaluateNode(node.Target, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if !val.IsFunction() {
		return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
	}

	posCount := len(node.Arguments.Positional)
	nameCount := len(node.Arguments.Named)

	args := arena.Alloc[NamedValue](ctx.State.Registry.Allocator, posCount+nameCount)
	for i, a := range node.Arguments.Positional {
		// v, err := EvaluateNodeStrict(a.Expr, scopeId, ctx)
		v, err := evaluateNodeLazy(a.Expr, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		args[i] = NamedValue{Value: v}
	}
	for i, a := range node.Arguments.Named {
		// a.Name
		// v, err := EvaluateNodeStrict(a.Arg, scopeId, ctx)
		v, err := evaluateNodeLazy(a.Arg, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		nameKeyId := ctx.State.Interner.Intern(string(a.Name))
		args[i+posCount] = NamedValue{nameKeyId, v}
	}

	res, err := val.Function(ctx).Exec(args, ctx)
	if err != nil {
		return ValueNone, err
	}
	return res, nil
}

func handleBinary(node *ast.Binary, scopePtr uintptr, ctx Context) (Value, error) {
	left, err := EvaluateNode(node.Left, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	// Check if fast exit is possible
	switch node.Op {
	case ast.BopAnd:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", left.Type().String())
		}

		if !left.Bool() {
			return MakeBool(false), nil
		}

		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", right.Type().String())
		}

		return right, nil

	case ast.BopOr:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", left.Type().String())
		}

		if left.Bool() {
			return MakeBool(true), nil
		}

		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}

		return right, nil

	default:
		right, err := EvaluateNode(node.Right, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}

		res, err := handleBinaryOp(node.Op, left, right, ctx)
		if err != nil {
			return ValueNone, err
		}

		return res, nil
	}

}

func handleUnary(node *ast.Unary, scopePtr uintptr, ctx Context) (Value, error) {
	unary, err := EvaluateNode(node.Expr, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	switch node.Op {
	default:
		return ValueNone, fmt.Errorf("unhandler unary type: %s", node.Op.String())
	case ast.UopNot:
		if !unary.IsBool() {
			return ValueNone, fmt.Errorf("unexpected unary type %s for op %s, expected boolean", unary.Type().String(), node.Op.String())
		}
		return MakeBool(!unary.Bool()), nil
	case ast.UopMinus:
		if !unary.IsNumber() {
			return ValueNone, fmt.Errorf("unexpected unary type %s for op %s, expected number", unary.Type().String(), node.Op.String())
		}
		res := -unary.Number()
		return MakeNumber(res), nil
	case ast.UopBitwiseNot:
		if !unary.IsNumber() {
			return ValueNone, fmt.Errorf("unexpected unary type %s for op %s, expected number", unary.Type().String(), node.Op.String())
		}
		val32 := int64(unary.Number())
		notVal32 := ^val32
		return MakeNumber(float64(notVal32)), nil
	}
}

func handleVar(node *ast.Var, scopePtr uintptr, ctx Context) (Value, error) {
	name := string(node.Id)

	keyId := ctx.State.Interner.Intern(name)

	val, found := ctx.GetScopeBind(scopePtr, keyId)
	if !found {
		return ValueNone, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", name))
	}

	return val, nil
}

func handleFunction(node *ast.Function, scopePtr uintptr, ctx Context) (Value, error) {

	paramCount := len(node.Parameters)

	paramKeyIds := arena.Alloc[uint32](ctx.State.Registry.Allocator, paramCount)
	for i, p := range node.Parameters {
		paramKeyIds[i] = ctx.State.Interner.Intern(string(p.Name))
	}

	var fn Func = func(args []NamedValue, _ Context) (Value, error) {
		if len(args) > paramCount {
			return ValueNone, fmt.Errorf("unexpected amount of args passed to function")
		}

		if paramCount == 0 {
			return EvaluateNode(node.Body, scopePtr, ctx)
		}

		s, childScopeId := ctx.NewScope(scopePtr, paramCount)

		// TODO: Throw err on argument x already provided

		posIndex := 0
		for i := range paramCount {
			keyId := paramKeyIds[i]

			var bindVal NamedValue

			if posIndex < len(args) && args[posIndex].Key == 0 {
				bindVal = args[posIndex]
				bindVal.Key = keyId
				posIndex++
			} else {
				for _, a := range args {
					if a.Key == keyId {
						bindVal = a
						break
					}
				}
			}

			if !bindVal.IsNone() {
				s.Bindings[i] = bindVal
				continue
			}

			// No arg was passed, fallback to default arg
			defArgNode := node.Parameters[i].DefaultArg
			if defArgNode == nil {
				return ValueNone, fmt.Errorf("arg (%d) with no default arg had no value passed", i)
			}

			da, err := evaluateNodeLazy(defArgNode, childScopeId, ctx)
			if err != nil {
				return ValueNone, err
			}

			s.Bindings[i] = NamedValue{keyId, da}
		}

		return EvaluateNode(node.Body, childScopeId, ctx)
	}

	f := Function{
		argsCount: len(paramKeyIds),
		fn:        fn,
	}

	return MakeFunction(f, ctx), nil
}

func handleConditional(node *ast.Conditional, scopePtr uintptr, ctx Context) (Value, error) {
	cond, err := EvaluateNode(node.Cond, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !cond.IsBool() {
		return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
	}

	if cond.Bool() {
		bt, err := EvaluateNode(node.BranchTrue, scopePtr, ctx)
		if err != nil {
			return ValueNone, err
		}
		return bt, nil
	}

	bf, err := EvaluateNode(node.BranchFalse, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	return bf, nil
}

func handleIndex(node *ast.Index, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	target, err := EvaluateNode(node.Target, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	switch target.Type() {
	default:
		return ValueNone, fmt.Errorf("value not indexable: %s", target.Type().String())
	case ValueTypeString:
		if !index.IsNumber() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing string, expected number, got %s", index.Type().String()))
		}
		i := int(index.Number())
		if len(target.String(ctx)) <= i {
			return ValueNone, MakeRuntimeError(fmt.Errorf("index (%d) out of bounds, string length %d", i, len(target.Array(ctx))))
		}
		s := target.String(ctx)
		return MakeString(string(s[i]), ctx), nil
	case ValueTypeObject:
		if !index.IsString() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing object, expected string, got %s", index.Type().String()))
		}

		// var keyId uint32
		// if index.IsStringConst() {
		// 	keyId = index.RefId()
		// } else {
		// 	name := index.String(ctx)
		// 	keyId = ctx.State.Interner.Intern(name)
		// }

		// TODO: think abt this again
		name := index.String(ctx)
		keyId := ctx.State.Interner.Intern(name)

		obj := target.Object(ctx)

		// Reset self to point to correct obj
		subCtx := ctx
		subCtx.Self = target

		val, _, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return ValueNone, err
		}
		if val.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", index.String(ctx)))
		}
		val, err = val.Eval(subCtx)
		if err != nil {
			return ValueNone, err
		}
		return val, nil
	case ValueTypeArray:
		if !index.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, index.Type())
		}
		i := int(index.Number())

		arr := target.Array(ctx)
		if i < 0 || len(arr) <= i {
			return ValueNone, MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", i, len(arr)))
		}
		return target.Array(ctx)[i], nil
	}

}

func handleSuperIndex(node *ast.SuperIndex, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		return ValueNone, errors.New("ctx.Self not set")
	}

	// var keyId uint32
	// if index.IsStringConst() {
	// 	keyId = index.RefId()
	// } else {
	// 	name := index.String(ctx)
	// 	keyId = ctx.State.Interner.Intern(name)
	// }

	// TODO: think abt this again
	name := index.String(ctx)
	keyId := ctx.State.Interner.Intern(name)

	obj := ctx.Self.Object(ctx)

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	if val.IsNone() {
		return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", index.String(ctx)))
	}

	val, err = val.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	return val, nil
}

func handleInSuper(node *ast.InSuper, scopePtr uintptr, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		return ValueNone, errors.New("ctx.Self not set")
	}

	// var keyId uint32
	// if index.IsStringConst() {
	// 	keyId = index.RefId()
	// } else {
	// 	name := index.String(ctx)
	// 	keyId = ctx.State.Interner.Intern(name)
	// }

	// TODO: think abt this again
	name := index.String(ctx)
	keyId := ctx.State.Interner.Intern(name)

	obj := ctx.Self.Object(ctx)

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	return MakeBool(!val.IsNone()), nil
}

func handleError(node *ast.Error, scopePtr uintptr, ctx Context) (Value, error) {
	msg, err := EvaluateNode(node.Expr, scopePtr, ctx)
	if err != nil {
		return ValueNone, err
	}
	if !msg.IsString() {
		return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
	}
	return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))
}

func handleImport(node *ast.Import, ctx Context) (Value, error) {

	// TODO: optimize import loop below

	file := node.File.Value

	var importedNode ast.Node
	var finalPath string

	importer := ctx.State.Environment.Importer

	dirs := []string{""}
	if !filepath.IsAbs(file) {
		dirs = []string{filepath.Dir(node.NodeBase.LocRange.FileName)}
		dirs = append(dirs, importer.JPaths...)
	}

	var rangeErr error
	for _, dir := range dirs {
		fp, err := filepath.Abs(filepath.Join(dir, file))
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}

		v := importer.Get(fp)
		if !v.IsNone() {
			return v, nil
		}

		// TODO: check and mark fp loading to catch import loops

		in, innerErr := importer.ResolveImport(fp)
		if os.IsNotExist(innerErr) {
			rangeErr = errors.Join(rangeErr, innerErr)
			continue
		}

		if innerErr != nil {
			return ValueNone, innerErr
		}

		importedNode = in
		finalPath = fp
		break

	}

	if importedNode == nil {
		return ValueNone, errors.Join(errors.New("error resolving import"), rangeErr)
	}

	importScope := CreateFileScope(file, importer.BaseStd, ctx)

	importCtx := ctx
	importCtx.Self = ValueNone

	v, err := evaluateNodeLazy(importedNode, importScope, importCtx)
	if err != nil {
		return ValueNone, err
	}

	importer.Set(finalPath, v)

	return v, nil

}

func handleImportStr(node *ast.ImportStr, ctx Context) (Value, error) {
	// TODO: make string cache
	// importer := ctx.Environment.Importer

	// TODO: take full path here?
	filePath := string(node.File.Value)

	// currentFileDir := filepath.Dir(node.NodeBase.LocRange.FileName)

	// fp := filepath.Join(currentFileDir, filePath)
	fp := filePath
	// fmt.Println(currentFileDir)

	// val := importer.Get(fp)
	// if !val.IsNone() {
	// 	return val, nil
	// }

	fileData, err := os.ReadFile(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return ValueNone, err
		}
		return ValueNone, fmt.Errorf("failed importing file: %s, err: %w", fp, err)
	}

	res := unsafe.String(unsafe.SliceData(fileData), len(fileData))

	// importer.Set(fp, res)

	return MakeString(res, ctx), nil
}
