package evaluator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/go-jsonnet/ast"
)

func EvaluateNode(n ast.Node, scopeId uint32, ctx Context) (Value, error) {
	val, err := evaluateNode(n, scopeId, ctx)
	if err != nil {
		return Value{}, WrapError(err, n)
	}
	if val.IsThunk() {
		val, err = val.Eval(ctx)
		if err != nil {
			return Value{}, WrapError(err, n)
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

func CreateFileScope(filename string, baseStd Value, ctx Context) uint32 {
	keyId := ctx.Interner.Intern("thisFile")

	layer := &Layer{
		Keys:   []uint32{keyId},
		Values: []Value{MakeString(filename, ctx)},
		Meta:   []uint8{0},
	}

	fileObj := NewObject([]*Layer{layer})
	fileObjVal := MakeObject(fileObj, ctx)

	mergedObj := MergeObjects(baseStd.RefId(), fileObjVal.RefId())
	fileStd := MakeObject(mergedObj, ctx)

	scopeId := ctx.NewScope(2, 2)
	s := ctx.Registry.Scopes.GetPtr(scopeId)

	s.Bindings[0] = NamedValue{Key: ctx.Interner.Intern("$std"), Value: fileStd}
	s.Bindings[1] = NamedValue{Key: ctx.Interner.Intern("std"), Value: fileStd}

	return scopeId
}

func evaluateNodeLazy(n ast.Node, scopeId uint32, ctx Context) (Value, error) {
	switch node := n.(type) {
	case *ast.LiteralString:
		return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return Value{}, fmt.Errorf("failed to parse float val (%s), err: %w", node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.Self:
		// if ctx.Self.IsNone() {
		// 	return Value{}, errors.New("self not set")
		// }
		return ctx.Self, nil
	case *ast.Local:
		return handleLocal(node, scopeId, ctx)
	default:
		return MakeThunk(NewThunk(node, scopeId, ctx), ctx), nil
	}
}

func evaluateNode(n ast.Node, scopeId uint32, ctx Context) (Value, error) {
	switch node := n.(type) {
	default:
		return Value{}, fmt.Errorf("unhandled node type: %T", node)
	case *ast.LiteralString:
		return MakeString(node.Value, ctx), nil
	case *ast.LiteralNull:
		return MakeNull(), nil
	case *ast.LiteralBoolean:
		return MakeBool(node.Value), nil
	case *ast.LiteralNumber:
		num, err := strconv.ParseFloat(node.OriginalString, 64)
		if err != nil {
			return Value{}, fmt.Errorf("(%T) failed to parse float val (%s), err: %w", node, node.OriginalString, err)
		}
		return MakeNumber(num), nil
	case *ast.DesugaredObject:
		return handleDesugaredObject(node, scopeId, ctx)
	case *ast.Array:
		res := make([]Value, 0, len(node.Elements))
		for _, v := range node.Elements {
			ev, err := evaluateNodeLazy(v.Expr, scopeId, ctx)
			if err != nil {
				return Value{}, err
			}
			res = append(res, ev)
		}
		return MakeArray(res, ctx), nil
	case *ast.Local:
		return handleLocal(node, scopeId, ctx)
	case *ast.Apply:

		val, err := EvaluateNode(node.Target, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		if !val.IsFunction() {
			return Value{}, TypeErrorSpecific(ValueTypeFunction, val.Type())
		}

		args := ctx.Registry.NamedValueBufs.Alloc(0, len(node.Arguments.Positional)+len(node.Arguments.Named))
		for _, a := range node.Arguments.Positional {
			// v, err := EvaluateNodeStrict(a.Expr, scopeId, ctx)
			v, err := evaluateNodeLazy(a.Expr, scopeId, ctx)
			if err != nil {
				return Value{}, err
			}
			args = append(args, NamedValue{Value: v})
		}
		for _, a := range node.Arguments.Named {
			// a.Name
			// v, err := EvaluateNodeStrict(a.Arg, scopeId, ctx)
			v, err := evaluateNodeLazy(a.Arg, scopeId, ctx)
			if err != nil {
				return Value{}, err
			}
			nameKeyId := ctx.Interner.Intern(string(a.Name))
			args = append(args, NamedValue{nameKeyId, v})
		}

		res, err := val.Function(ctx).Exec(args, ctx)
		if err != nil {
			return Value{}, err
		}
		return res, nil
	case *ast.Index:
		return handleIndex(node, scopeId, ctx)
	case *ast.Var:
		name := string(node.Id)

		keyId := ctx.Interner.Intern(name)

		val, found := ctx.GetScopeBind(scopeId, keyId)
		if !found {
			return Value{}, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", name))
		}

		return val, nil
	case *ast.Function:
		return handleFunction(node, scopeId, ctx)
	case *ast.Conditional:
		cond, err := EvaluateNode(node.Cond, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}
		if !cond.IsBool() {
			return Value{}, TypeErrorSpecific(ValueTypeBool, cond.Type())
		}

		if cond.Bool() {
			bt, err := EvaluateNode(node.BranchTrue, scopeId, ctx)
			if err != nil {
				return Value{}, err
			}
			return bt, nil
		}

		bf, err := EvaluateNode(node.BranchFalse, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}
		return bf, nil
	case *ast.Binary:
		return handleBinary(node, scopeId, ctx)
	case *ast.Unary:
		unary, err := EvaluateNode(node.Expr, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		switch node.Op {
		default:
			return Value{}, fmt.Errorf("unhandler unary type: %s", node.Op.String())
		case ast.UopNot:
			if !unary.IsBool() {
				return Value{}, fmt.Errorf("unexpected unary type %s for op %s, expected boolean", unary.Type().String(), node.Op.String())
			}
			return MakeBool(!unary.Bool()), nil
		case ast.UopMinus:
			if !unary.IsNumber() {
				return Value{}, fmt.Errorf("unexpected unary type %s for op %s, expected number", unary.Type().String(), node.Op.String())
			}
			res := -unary.Number()
			return MakeNumber(res), nil
		case ast.UopBitwiseNot:
			if !unary.IsNumber() {
				return Value{}, fmt.Errorf("unexpected unary type %s for op %s, expected number", unary.Type().String(), node.Op.String())
			}
			val32 := int64(unary.Number())
			notVal32 := ^val32
			return MakeNumber(float64(notVal32)), nil
		}
	case *ast.Import:
		return handleImport(node, scopeId, ctx)
	case *ast.ImportStr:

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
				return Value{}, err
			}
			return Value{}, fmt.Errorf("failed importing file: %s, err: %w", fp, err)
		}

		res := MakeString(string(fileData), ctx)

		// importer.Set(fp, res)

		return res, nil
	case *ast.Self:
		// if ctx.Self.IsNone() {
		// 	return Value{}, errors.New("self not set")
		// }
		return ctx.Self, nil
	case *ast.SuperIndex:
		return handleSuperIndex(node, scopeId, ctx)
	case *ast.Error:
		msg, err := EvaluateNode(node.Expr, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}
		if !msg.IsString() {
			return Value{}, TypeErrorSpecific(ValueTypeString, msg.Type())
		}
		return Value{}, MakeRuntimeError(errors.New(msg.String(ctx)))
	case *GoCallbackNode:
		res, err := node.Func.Exec(node.Args, ctx)
		if err != nil {
			return Value{}, err
		}
		return res, nil
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

func handleDesugaredObject(node *ast.DesugaredObject, scopeId uint32, ctx Context) (Value, error) {

	fieldCount := len(node.Fields)
	localsCount := len(node.Locals)

	layer := &Layer{ParentScopeId: scopeId}

	if fieldCount > 0 {
		layer.Keys = ctx.Registry.Uint32Bufs.Alloc(0, fieldCount)
		layer.Nodes = ctx.Registry.NodesBufs.Alloc(0, fieldCount)
		layer.Meta = ctx.Registry.Uint8Bufs.Alloc(0, fieldCount)
	}

	if localsCount > 0 {
		layer.LocalKeys = ctx.Registry.Uint32Bufs.Alloc(0, localsCount)
		layer.LocalNodes = ctx.Registry.NodesBufs.Alloc(0, localsCount)
	}

	if len(node.Asserts) > 0 {
		layer.Asserts = ctx.Registry.NodesBufs.Alloc(len(node.Asserts), len(node.Asserts))
		copy(layer.Asserts, node.Asserts)
	}

	for _, v := range node.Locals {

		name := string(v.Variable)
		keyId := ctx.Interner.Intern(name)

		layer.LocalKeys = append(layer.LocalKeys, keyId)
		layer.LocalNodes = append(layer.LocalNodes, v.Body)

	}

	useMap := fieldCount > MaxLinearKeys

	if useMap {
		layer.Index = make(map[uint32]int, fieldCount)
	}

	index := 0
	for _, v := range node.Fields {
		name, err := EvaluateNode(v.Name, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		if name.IsNull() {
			// Omitted field
			continue
		}

		if !name.IsString() {
			return Value{}, fmt.Errorf("unexpected field name type %s, expected string", name.Type().String())
		}

		n := name.String(ctx)

		keyId := ctx.Interner.Intern(n)

		layer.Keys = append(layer.Keys, keyId)
		layer.Nodes = append(layer.Nodes, v.Body)
		layer.Meta = append(layer.Meta, CreateFieldMeta(v.Hide, v.PlusSuper))

		if useMap {
			layer.Index[keyId] = index
			index++
		}

	}

	layers := ctx.Registry.LayerBufs.Alloc(1, 1)
	layers[0] = layer
	obj := NewObject(layers)

	return MakeObject(obj, ctx), nil
}

func handleLocal(node *ast.Local, scopeId uint32, ctx Context) (Value, error) {

	childScopeId := ctx.NewScope(scopeId, len(node.Binds))

	s := ctx.Registry.Scopes.GetPtr(childScopeId)

	for i, v := range node.Binds {
		vname := string(v.Variable)
		keyId := ctx.Interner.Intern(vname)
		t, err := evaluateNodeLazy(v.Body, childScopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		s.Bindings[i] = NamedValue{keyId, t}
	}

	val, err := EvaluateNode(node.Body, childScopeId, ctx)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

func handleBinary(node *ast.Binary, scopeId uint32, ctx Context) (Value, error) {
	left, err := EvaluateNode(node.Left, scopeId, ctx)
	if err != nil {
		return Value{}, err
	}

	// Check if fast exit is possible
	switch node.Op {
	case ast.BopAnd:
		if !left.IsBool() {
			return Value{}, fmt.Errorf("unexpected type %s for && op, expected boolean", left.Type().String())
		}

		if !left.Bool() {
			return MakeBool(false), nil
		}

		right, err := EvaluateNode(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		if !right.IsBool() {
			return Value{}, fmt.Errorf("unexpected type %s for && op, expected boolean", right.Type().String())
		}
		res := left.Bool() && right.Bool()
		return MakeBool(res), nil

	case ast.BopOr:
		if !left.IsBool() {
			return Value{}, fmt.Errorf("unexpected type %s for || op, expected boolean", left.Type().String())
		}

		if left.Bool() {
			return MakeBool(true), nil
		}

		right, err := EvaluateNode(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		if !right.IsBool() {
			return Value{}, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}
		res := left.Bool() || right.Bool()
		return MakeBool(res), nil
	default:
		right, err := EvaluateNode(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		res, err := handleBinaryOp(node.Op, left, right, ctx)
		if err != nil {
			return Value{}, err
		}

		return res, nil
	}

}

func handleFunction(node *ast.Function, scopeId uint32, ctx Context) (Value, error) {

	paramCount := len(node.Parameters)

	paramKeyIds := make([]uint32, paramCount)
	for i, p := range node.Parameters {
		paramKeyIds[i] = ctx.Interner.Intern(string(p.Name))
	}

	var fn Func = func(args []NamedValue, _ Context) (Value, error) {
		if len(args) > paramCount {
			return Value{}, fmt.Errorf("unexpected amount of args passed to function")
		}

		childScopeId := ctx.NewScope(scopeId, paramCount)
		s := ctx.Registry.Scopes.GetPtr(childScopeId)

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
				return Value{}, fmt.Errorf("arg (%d) with no default arg had no value passed", i)
			}

			da, err := evaluateNodeLazy(defArgNode, childScopeId, ctx)
			if err != nil {
				return Value{}, err
			}

			s.Bindings[i] = NamedValue{keyId, da}
		}

		val, err := EvaluateNode(node.Body, childScopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		return val, nil
	}

	f := Function{
		Args: paramKeyIds,
		fn:   fn,
	}

	return MakeFunction(f, ctx), nil
}

func handleIndex(node *ast.Index, scopeId uint32, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopeId, ctx)
	if err != nil {
		return Value{}, err
	}

	target, err := EvaluateNode(node.Target, scopeId, ctx)
	if err != nil {
		return Value{}, err
	}

	switch target.Type() {
	default:
		return Value{}, fmt.Errorf("value not indexable: %s", target.Type().String())
	case ValueTypeString:
		if !index.IsNumber() {
			return Value{}, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing string, expected number, got %s", index.Type().String()))
		}
		i := int(index.Number())
		if len(target.String(ctx)) <= i {
			return Value{}, MakeRuntimeError(fmt.Errorf("index (%d) out of bounds, string length %d", i, len(target.Array(ctx))))
		}
		s := target.String(ctx)
		return MakeString(string(s[i]), ctx), nil
	case ValueTypeObject:
		if !index.IsString() {
			return Value{}, MakeRuntimeError(fmt.Errorf("unexpected index type for indexing object, expected string, got %s", index.Type().String()))
		}

		// TODO: can we optimize this? Since the index is a string we can take the id directly. If DataString is implemented that wont work...
		name := index.String(ctx)

		keyId := ctx.Interner.Intern(name)

		obj := target.Object(ctx)

		// Reset self to point to correct obj
		subCtx := ctx
		subCtx.Self = target

		val, _, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return Value{}, err
		}
		if val.IsNone() {
			return Value{}, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", name))
		}
		val, err = val.Eval(subCtx)
		if err != nil {
			return Value{}, err
		}
		return val, nil
	case ValueTypeArray:
		if !index.IsNumber() {
			return Value{}, TypeErrorSpecific(ValueTypeNumber, index.Type())
		}
		i := int(index.Number())

		arr := target.Array(ctx)
		if i < 0 || len(arr) <= i {
			return Value{}, MakeRuntimeError(fmt.Errorf("Index %d out of bounds, not within [0, %d)", i, len(arr)))
		}
		return target.Array(ctx)[i], nil
	}

}

func handleSuperIndex(node *ast.SuperIndex, scopeId uint32, ctx Context) (Value, error) {
	index, err := EvaluateNode(node.Index, scopeId, ctx)
	if err != nil {
		return Value{}, err
	}

	if ctx.Self.IsNone() {
		return Value{}, errors.New("ctx.Self not set")
	}

	name := index.String(ctx)

	keyId := ctx.Interner.Intern(name)

	obj := ctx.Self.Object(ctx)

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, targetOffset)
	if err != nil {
		return Value{}, err
	}

	if val.IsNone() {
		return Value{}, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", name))
	}

	val, err = val.Eval(ctx)
	if err != nil {
		return Value{}, err
	}

	return val, nil
}

func handleImport(node *ast.Import, scopeId uint32, ctx Context) (Value, error) {
	fileVal, err := EvaluateNode(node.File, scopeId, ctx)
	if err != nil {
		return Value{}, err
	}
	if !fileVal.IsString() {
		return Value{}, fmt.Errorf("(%T) unexpected file data type '%s'", node, fileVal.Type().String())
	}

	file := fileVal.String(ctx)

	currentFileDir := filepath.Dir(node.NodeBase.LocRange.FileName)

	var importedNode ast.Node
	var finalPath string

	importer := ctx.Environment.Importer

	dirs := []string{""}
	if !filepath.IsAbs(file) {
		dirs = []string{currentFileDir}
		dirs = append(dirs, importer.JPaths...)
	}

	var rangeErr error
	for _, dir := range dirs {
		fp, err := filepath.Abs(filepath.Join(dir, file))
		if err != nil {
			return Value{}, MakeRuntimeError(err)
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
			return Value{}, innerErr
		}

		importedNode = in
		finalPath = fp
		break

	}

	if importedNode == nil {
		return Value{}, errors.Join(errors.New("error resolving import"), rangeErr)
	}

	importScope := CreateFileScope(file, importer.BaseStd, ctx)

	importCtx := ctx
	importCtx.Self = Value{}

	v, err := evaluateNodeLazy(importedNode, importScope, importCtx)
	if err != nil {
		return Value{}, err
	}

	importer.Set(finalPath, v)

	return v, nil

}
