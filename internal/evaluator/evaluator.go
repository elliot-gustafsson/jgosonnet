package evaluator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/go-jsonnet/ast"
)

func EvaluateNodeStrict(n ast.Node, scopeId uint32, ctx Context) (Value, error) {
	val, err := evaluateNode(n, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, n.Loc())
	}
	err = EvaluateValueStrict(&val, ctx)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

func EvaluateValueStrict(value *Value, ctx Context) error {
	err := evaluateValue(value, ctx)
	if err != nil {
		return err
	}
	if value.IsThunk() {
		err := EvaluateValueStrict(value, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ManifestValue(value Value, ctx Context) (any, error) {

	switch value.Type() {
	default:
		return nil, fmt.Errorf("unhandled value type '%s'", value.Type().String())
	case ValueTypeNull:
		return nil, nil
	case ValueTypeString:
		return value.String(ctx), nil
	case ValueTypeNumber:
		return value.Number(), nil
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
		res, err := value.Function(ctx)(nil, ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(res, ctx)
	case ValueTypeThunk:
		err := evaluateValue(&value, ctx)
		if err != nil {
			return nil, err
		}
		return ManifestValue(value, ctx)
	}
}

func createErrorWithContext(err error, loc *ast.LocationRange) error {
	return fmt.Errorf("%w\n\nlocation: %s", err, loc.String())
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
	// res := MakeThunk(Thunk{Node: node, ScopeId: scopeId, SkipMemoize: true}, ctx)
	// return res, nil
	// case *ast.Local:
	// 	return handleLocal(node, scopeId, ctx)
	default:
		return MakeThunk(Thunk{
			Node:                node,
			ScopeId:             scopeId,
			CapturedSelf:        ctx.Self,
			CapturedSuperOffset: ctx.SuperOffset,
		}, ctx), nil
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

		val, err := EvaluateNodeStrict(node.Target, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}

		if !val.IsFunction() {
			// return Value{}, fmt.Errorf("(%T) unexpected value type '%s'", node, val.Type.String())
			return Value{}, createErrorWithContext(fmt.Errorf("(%T) unexpected value type '%s'", node, val.Type().String()), &node.LocRange)
		}

		args := make([]Value, 0, len(node.Arguments.Positional)+len(node.Arguments.Named))
		for _, a := range node.Arguments.Positional {
			v, err := EvaluateNodeStrict(a.Expr, scopeId, ctx)
			if err != nil {
				return Value{}, createErrorWithContext(err, &node.LocRange)
			}
			args = append(args, v)
		}
		for _, a := range node.Arguments.Named {
			// a.Name
			v, err := EvaluateNodeStrict(a.Arg, scopeId, ctx)
			if err != nil {
				return Value{}, createErrorWithContext(err, &node.LocRange)
			}
			args = append(args, v)
		}

		res, err := val.Function(ctx)(args, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}
		return res, nil
	case *ast.Index:
		return handleIndex(node, scopeId, ctx)
	case *ast.Var:
		name := string(node.Id)

		// if name == "$" {
		// 	if ctx.Root.IsNone() {
		// 		return Value{}, createErrorWithContext(fmt.Errorf("root value is not defined"), &node.LocRange)
		// 	}
		// 	return ctx.Root, nil
		// }

		keyId := ctx.Interner.Intern(name)

		// val, found := scopeId.FindBinding(keyId)
		// val, found := FindScopeBinding(scopeId, keyId, ctx)
		val, found := ctx.Arena.GetScopeBind(scopeId, keyId)
		if !found {
			val, _ := ctx.Arena.GetScopeBind(scopeId, keyId)
			if val.IsNone() {

			}
			return Value{}, createErrorWithContext(fmt.Errorf("variable not found in scope, name: %s", name), &node.LocRange)
		}

		// if val.IsObject() && name != "std" {
		// 	obj := val.Object(ctx)
		// 	clone := obj.Clone()
		// 	return MakeObject(clone, ctx), nil
		// }

		return val, nil
	case *ast.Function:
		return handleFunction(node, scopeId, ctx)
	case *ast.Conditional:
		cond, err := EvaluateNodeStrict(node.Cond, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}
		if !cond.IsBool() {
			return Value{}, fmt.Errorf("(%T) unexpected conditional type '%s'", node, cond.Type().String())
		}

		if cond.Bool() {
			bt, err := evaluateNode(node.BranchTrue, scopeId, ctx)
			if err != nil {
				return Value{}, createErrorWithContext(err, &node.LocRange)
			}
			return bt, nil
		}

		bf, err := evaluateNode(node.BranchFalse, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}
		return bf, nil
	case *ast.Binary:
		return handleBinary(node, scopeId, ctx)
	case *ast.Unary:
		unary, err := EvaluateNodeStrict(node.Expr, scopeId, ctx)
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
		}
	case *ast.Import:
		return handleImport(node, scopeId, ctx)
	case *ast.Self:
		// if ctx.Self.IsNone() {
		// 	return Value{}, errors.New("self not set")
		// }
		return ctx.Self, nil
	case *ast.SuperIndex:
		return handleSuperIndex(node, scopeId, ctx)
	case *ast.Error:
		msg, err := EvaluateNodeStrict(node.Expr, scopeId, ctx)
		if err != nil {
			return Value{}, err
		}
		if !msg.IsString() {
			return Value{}, createErrorWithContext(fmt.Errorf("unexpected value type '%s', expected string", msg.Type().String()), &node.LocRange)
		}
		return Value{}, errors.New(msg.String(ctx))
	}
}

func evaluateValue(value *Value, ctx Context) error {
	if !value.IsThunk() {
		return nil
	}
	thunk := value.Thunk(ctx)
	if thunk == nil {
		return nil
	}
	if !thunk.Value.IsNone() {
		*value = thunk.Value
		return nil
	}

	// if thunk.Busy {
	// 	return createErrorWithContext(fmt.Errorf("error during thunk evaluation, infinite recursion detected"), thunk.Node.Loc())
	// }

	// thunk.Busy = true

	evalCtx := ctx
	evalCtx.Self = thunk.CapturedSelf
	evalCtx.SuperOffset = thunk.CapturedSuperOffset

	evaledVal, err := evaluateNode(thunk.Node, thunk.ScopeId, evalCtx)
	if err != nil {
		return createErrorWithContext(err, thunk.Node.Loc())
	}

	// thunk = value.Thunk(ctx)
	// thunk.Busy = false

	if evaledVal.IsThunk() {
		err = evaluateValue(&evaledVal, ctx)
		if err != nil {
			return err
		}
	}

	*value = evaledVal
	thunk.Value = evaledVal
	return nil
}

func handleDesugaredObject(node *ast.DesugaredObject, scopeId uint32, ctx Context) (Value, error) {

	fieldCount := len(node.Fields)
	localsCount := len(node.Locals)

	layer := &Layer{
		ParentScopeId: scopeId,

		Keys:  make([]uint32, 0, fieldCount),
		Nodes: make(ast.Nodes, 0, fieldCount),
		Meta:  make([]uint8, 0, fieldCount),

		LocalKeys:  make([]uint32, 0, localsCount),
		LocalNodes: make(ast.Nodes, 0, localsCount),

		Asserts: make(ast.Nodes, 0, len(node.Asserts)),
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
		name, err := EvaluateNodeStrict(v.Name, scopeId, ctx)
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

		// if n == "mapRuleGroups" {
		// 	log.Println("hej")
		// }

		keyId := ctx.Interner.Intern(n)

		layer.Keys = append(layer.Keys, keyId)
		layer.Nodes = append(layer.Nodes, v.Body)
		layer.Meta = append(layer.Meta, CreateFieldMeta(v.Hide, v.PlusSuper))

		if useMap {
			layer.Index[keyId] = index
			index++
		}

	}

	for _, v := range node.Asserts {
		// ass, err := evaluateNode(v, childScopeId, ctx)
		// if err != nil {
		// 	return Value{}, err
		// }
		// if !ass.IsBool() {
		// 	return Value{}, fmt.Errorf("(%T) unexpected assert return, err: %w", v, err)
		// }

		layer.Asserts = append(layer.Asserts, v)
	}

	obj := NewObject([]*Layer{layer})

	return MakeObject(obj, ctx), nil
}

func handleLocal(node *ast.Local, scopeId uint32, ctx Context) (Value, error) {

	childScopeId := ctx.Arena.NewScope(scopeId, len(node.Binds))

	for _, v := range node.Binds {
		vname := string(v.Variable)
		keyId := ctx.Interner.Intern(vname)
		t, err := evaluateNodeLazy(v.Body, childScopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		ctx.Arena.AddScopeBind(childScopeId, keyId, t)
	}

	val, err := evaluateNodeLazy(node.Body, childScopeId, ctx)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

func handleBinary(node *ast.Binary, scopeId uint32, ctx Context) (Value, error) {
	left, err := EvaluateNodeStrict(node.Left, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
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

		right, err := EvaluateNodeStrict(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
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

		right, err := EvaluateNodeStrict(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}

		if !right.IsBool() {
			return Value{}, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}
		res := left.Bool() || right.Bool()
		return MakeBool(res), nil
	default:
		right, err := EvaluateNodeStrict(node.Right, scopeId, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}

		res, err := handleBinaryOp(node.Op, left, right, ctx)
		if err != nil {
			return Value{}, createErrorWithContext(fmt.Errorf("(%T) failed to handle binary op, err: %w", node, err), &node.LocRange)
		}

		// if res.IsObject() {
		// 	// Self needs to be reset to point to the virtually merged object
		// 	ctx.Self = &res
		// }

		return res, nil
	}

}

func handleFunction(node *ast.Function, scopeId uint32, ctx Context) (Value, error) {

	paramCount := len(node.Parameters)

	paramKeyIds := make([]uint32, paramCount)
	for i, p := range node.Parameters {
		paramKeyIds[i] = ctx.Interner.Intern(string(p.Name))
	}

	f := func(args []Value, _ Context) (Value, error) {
		if len(args) > paramCount {
			return Value{}, fmt.Errorf("unexpected amount of args passed to function")
		}

		// Note: Using callingCtx here to ensure we use the active Arena/State
		childScopeId := ctx.Arena.NewScope(scopeId, paramCount)

		for i := range paramCount {
			keyId := paramKeyIds[i]

			if i < len(args) {
				// Arg was passed, bind it directly
				ctx.Arena.AddScopeBind(childScopeId, keyId, args[i])
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

			ctx.Arena.AddScopeBind(childScopeId, keyId, da)
		}

		val, err := evaluateNode(node.Body, childScopeId, ctx)
		if err != nil {
			return Value{}, err
		}

		return val, nil
	}

	return MakeFunction(f, ctx), nil
}

func handleIndex(node *ast.Index, scopeId uint32, ctx Context) (Value, error) {
	index, err := EvaluateNodeStrict(node.Index, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}

	target, err := EvaluateNodeStrict(node.Target, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}

	switch target.Type() {
	default:
		return Value{}, fmt.Errorf("value not indexable: %s", target.Type().String())
	case ValueTypeString:
		if !index.IsNumber() {
			return Value{}, createErrorWithContext(fmt.Errorf("unexpected index type for indexing string, expected number, got %s", index.Type().String()), &node.LocRange)
		}
		i := int(index.Number())
		if len(target.String(ctx)) <= i {
			return Value{}, createErrorWithContext(fmt.Errorf("index (%d) out of bounds, string length %d", i, len(target.Array(ctx))), &node.LocRange)
		}
		s := target.String(ctx)
		return MakeString(string(s[i]), ctx), nil
	case ValueTypeObject:
		if !index.IsString() {
			return Value{}, createErrorWithContext(fmt.Errorf("unexpected index type for indexing object, expected string, got %s", index.Type().String()), &node.LocRange)
		}

		name := index.String(ctx)

		keyId := ctx.Interner.Intern(name)

		obj := target.Object(ctx)

		// Reset self to point to correct obj
		subCtx := ctx
		subCtx.Self = target

		// if name == "mapRuleGroups" {
		// 	log.Printf("asdf")
		// }

		val, _, err := obj.GetField(keyId, subCtx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}
		if val.IsNone() {
			// fields := getObjectFields(obj, subCtx, true)
			// str := "\n"
			// for _, v := range fields {
			// 	str += "\n" + v.String(subCtx)
			// }
			// return Value{}, createErrorWithContext(fmt.Errorf("index not found on object, index: %s, %s", name, str), &node.LocRange)
			return Value{}, createErrorWithContext(fmt.Errorf("index not found on object, index: %s", name), &node.LocRange)
		}
		err = EvaluateValueStrict(&val, subCtx)
		if err != nil {
			return Value{}, createErrorWithContext(err, &node.LocRange)
		}
		return val, nil
	case ValueTypeArray:
		if !index.IsNumber() {
			return Value{}, createErrorWithContext(fmt.Errorf("unexpected index type for indexing array, expected number, got %s", index.Type().String()), &node.LocRange)
		}
		i := int(index.Number())
		if len(target.Array(ctx)) <= i {
			return Value{}, createErrorWithContext(fmt.Errorf("index (%d) out of bounds, array length %d", i, len(target.Array(ctx))), &node.LocRange)
		}
		return target.Array(ctx)[i], nil
	}

}

func handleSuperIndex(node *ast.SuperIndex, scopeId uint32, ctx Context) (Value, error) {
	index, err := evaluateNode(node.Index, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}

	if ctx.Self.IsNone() {
		return Value{}, errors.New("ctx.Self not set")
	}

	name := index.String(ctx)

	keyId := ctx.Interner.Intern(name)

	obj := ctx.Self.Object(ctx)

	targetOffset := ctx.SuperOffset + 1

	ctx.SuperOffset++
	val, _, err := obj.GetFieldWithOffset(keyId, ctx, targetOffset)
	if err != nil {
		return Value{}, err
	}

	if val.IsNone() {
		return Value{}, createErrorWithContext(fmt.Errorf("super index not found, index: %s", name), &node.LocRange)
	}

	err = EvaluateValueStrict(&val, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}

	return val, nil
}

func handleImport(node *ast.Import, scopeId uint32, ctx Context) (Value, error) {
	fileVal, err := EvaluateNodeStrict(node.File, scopeId, ctx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}
	if !fileVal.IsString() {
		return Value{}, fmt.Errorf("(%T) unexpected file data type '%s'", node, fileVal.Type().String())
	}

	file := fileVal.String(ctx)

	currentFileDir := filepath.Dir(node.NodeBase.LocRange.FileName)

	var importedNode ast.Node
	var finalPath string

	dirs := []string{""}
	if !filepath.IsAbs(file) {
		dirs = []string{currentFileDir}
		dirs = append(dirs, ctx.Importer.JPaths...)
	}

	var rangeErr error
	for _, dir := range dirs {
		fp := filepath.Join(dir, file)

		v := ctx.Importer.Get(fp)
		if !v.IsNone() {
			return v, nil
		}

		// TODO: check and mark fp loading to catch import loops

		in, innerErr := ctx.Importer.ResolveImport(fp)
		if os.IsNotExist(innerErr) {
			rangeErr = errors.Join(rangeErr, innerErr)
			continue
		}

		if innerErr != nil {
			return Value{}, createErrorWithContext(innerErr, &node.LocRange)
		}

		importedNode = in
		finalPath = fp
		break

	}

	if importedNode == nil {
		return Value{}, errors.Join(errors.New("error resolving import"), rangeErr)
	}

	importScope := ctx.Importer.ImportScope

	importCtx := ctx
	importCtx.Self = Value{}

	v, err := evaluateNodeLazy(importedNode, importScope, importCtx)
	if err != nil {
		return Value{}, createErrorWithContext(err, &node.LocRange)
	}

	ctx.Importer.Set(finalPath, v)

	return v, nil

}
