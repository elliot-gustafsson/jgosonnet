package evaluator

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	old_ast "github.com/google/go-jsonnet/ast"

	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
)

func EvaluateNode(tree *ast.AST, nodeId, scopeId uint32, ctx Context) (Value, error) {
	val, err := evaluateNode(tree, nodeId, scopeId, ctx)
	if err != nil {
		// return ValueNone, WrapError(err, n)
		return ValueNone, err // TODO: Fix WrapError
	}
	if val.IsThunk() {
		val, err = val.Eval(ctx)
		if err != nil {
			// return ValueNone, WrapError(err, n)
			return ValueNone, err // TODO: Fix WrapError
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
	keyId := ctx.State.Interner.Intern("thisFile")

	layer := &Layer{
		Keys:   []uint32{keyId},
		Values: []Value{MakeString(filename, ctx)},
		Meta:   []uint8{0},
	}

	fileObjId := NewObject([]*Layer{layer}, ctx)
	fileObjVal := MakeObjectValue(fileObjId)

	mergedObjId := MergeObjects(baseStd.RefId(), fileObjVal.RefId(), ctx)
	fileStd := MakeObjectValue(mergedObjId)

	s, scopeId := ctx.NewScope(0, 2)
	// s := ctx.State.Registry.Scopes.GetPtr(scopeId)

	s.Bindings[0] = NamedValue{Key: ctx.State.Interner.Intern("$std"), Value: fileStd}
	s.Bindings[1] = NamedValue{Key: ctx.State.Interner.Intern("std"), Value: fileStd}

	return scopeId
}

func evaluateNodeLazy(tree *ast.AST, nodeId, scopeId uint32, ctx Context) (Value, error) {
	node := tree.Nodes[nodeId]

	switch node.Type {
	case ast.NodeTypeString:
		return MakeStringValue(node.A | StringConstFlag), nil
	case ast.NodeTypeNull:
		return MakeNull(), nil
	case ast.NodeTypeFalse:
		return MakeFalse(), nil
	case ast.NodeTypeTrue:
		return MakeTrue(), nil
	case ast.NodeTypeNumber:
		bits := (uint64(node.A) << 32) | uint64(node.B)
		num := math.Float64frombits(bits)
		return MakeNumber(num), nil
	case ast.NodeTypeSelf:
		// if ctx.Self.IsNone() {
		// 	return ValueNone, errors.New("self not set")
		// }
		return ctx.Self, nil
	case ast.NodeTypeLocal:
		return handleLocal(tree, node, scopeId, ctx)
	default:
		return MakeThunk(NewThunk(nodeId, scopeId, ctx), ctx), nil
	}
}

func evaluateNode(tree *ast.AST, nodeId, scopeId uint32, ctx Context) (Value, error) {

	node := tree.Nodes[nodeId]

	switch node.Type {
	default:
		return ValueNone, fmt.Errorf("unhandled node type: %T", node)
	case ast.NodeTypeString:
		return MakeStringValue(node.A | StringConstFlag), nil
	case ast.NodeTypeNull:
		return MakeNull(), nil
	case ast.NodeTypeFalse:
		return MakeFalse(), nil
	case ast.NodeTypeTrue:
		return MakeTrue(), nil
	case ast.NodeTypeNumber:
		bits := (uint64(node.A) << 32) | uint64(node.B)
		num := math.Float64frombits(bits)
		return MakeNumber(num), nil
	case ast.NodeTypeObject:
		return handleDesugaredObject(tree, node, scopeId, ctx)
	case ast.NodeTypeArray:
		arr, val := MakeArraySized(int(node.B), ctx)

		elements := tree.SideTable[node.A : node.A+node.B]
		for i, nodeId := range elements {
			ev, err := evaluateNodeLazy(tree, nodeId, scopeId, ctx)
			if err != nil {
				return ValueNone, err
			}
			arr[i] = ev
		}
		return val, nil
	case ast.NodeTypeLocal:
		return handleLocal(tree, node, scopeId, ctx)

	case ast.NodeTypeApply:

		val, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !val.IsFunction() {
			return ValueNone, TypeErrorSpecific(ValueTypeFunction, val.Type())
		}

		callArgs := tree.SideTable[node.B : node.B+node.C*2]

		args := ctx.State.Registry.NamedValueBufs.Alloc(int(node.C), int(node.C))
		for i := range args {
			// TODO: think abt maybe not wasting space for pos args
			keyId, argNodeId := callArgs[i*2], callArgs[i*2+1]

			v, err := evaluateNodeLazy(tree, argNodeId, scopeId, ctx)
			if err != nil {
				return ValueNone, err
			}

			args[i] = NamedValue{keyId, v}
		}

		res, err := val.Function(ctx).Exec(args, ctx)
		if err != nil {
			return ValueNone, err
		}
		return res, nil
	case ast.NodeTypeIndex:
		return handleIndex(tree, node, scopeId, ctx)
	case ast.NodeTypeVar:
		val, found := ctx.GetScopeBind(scopeId, node.A)
		if !found {
			name := ctx.State.Interner.Get(node.A)
			return ValueNone, MakeRuntimeError(fmt.Errorf("variable not found in scope, name: %s", name))
		}
		return val, nil
	case ast.NodeTypeFunction:
		return handleFunction(tree, node, scopeId, ctx)
	case ast.NodeTypeConditional:
		cond, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !cond.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, cond.Type())
		}

		if cond.Bool() {
			bt, err := EvaluateNode(tree, node.B, scopeId, ctx)
			if err != nil {
				return ValueNone, err
			}
			return bt, nil
		}

		bf, err := EvaluateNode(tree, node.C, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		return bf, nil

	case ast.NodeTypeImport:
		return handleImport(tree, node, scopeId, ctx)
	case ast.NodeTypeImportStr:

		// TODO: make string cache
		// importer := ctx.Environment.Importer

		// TODO: take full path here?
		filePath := ctx.State.Interner.Get(node.A)

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

		res := MakeString(string(fileData), ctx)

		// importer.Set(fp, res)

		return res, nil
	case ast.NodeTypeSelf:
		return ctx.Self, nil
	case ast.NodeTypeSuperIndex:
		return handleSuperIndex(tree, node, scopeId, ctx)
	case ast.NodeTypeError:
		msg, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !msg.IsString() {
			return ValueNone, TypeErrorSpecific(ValueTypeString, msg.Type())
		}
		return ValueNone, MakeRuntimeError(errors.New(msg.String(ctx)))

	case ast.NodeTypeCallbackFunction:
		// TODO: fix this --------------------------------------------------------------------------------------------
		// TODO: fix this --------------------------------------------------------------------------------------------
		// TODO: fix this --------------------------------------------------------------------------------------------
		// TODO: fix this --------------------------------------------------------------------------------------------
		res, err := node.Func.Exec(node.Args, ctx)
		if err != nil {
			return ValueNone, err
		}
		return res, nil

	// -------------------------------------------------------------------------------
	// Unary ops
	// -------------------------------------------------------------------------------
	case ast.NodeTypeUnaryNot:
		unary, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, unary.Type())
		}
		return MakeBool(!unary.Bool()), nil

	case ast.NodeTypeUnaryMinus:
		unary, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		res := -unary.Number()
		return MakeNumber(res), nil

	case ast.NodeTypeUnaryBitwiseNot:
		unary, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !unary.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, unary.Type())
		}
		val32 := int64(unary.Number())
		notVal32 := ^val32
		return MakeNumber(float64(notVal32)), nil

	// -------------------------------------------------------------------------------
	// Binary ops
	// -------------------------------------------------------------------------------
	case ast.NodeTypeBinaryAnd:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !left.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, left.Type())
		}

		if !left.Bool() {
			return MakeBool(false), nil
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, right.Type())
		}
		res := left.Bool() && right.Bool()
		return MakeBool(res), nil

	case ast.NodeTypeBinaryOr:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}
		if !left.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, left.Type())
		}

		if !left.Bool() {
			return MakeBool(false), nil
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, TypeErrorSpecific(ValueTypeBool, right.Type())
		}
		// res := left.Bool() || right.Bool()
		return MakeBool(right.Bool()), nil
	case ast.NodeTypeBinaryEqual:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(eq), nil

	case ast.NodeTypeBinaryPlus:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		return bopPlus(left, right, ctx)

	case ast.NodeTypeBinaryMinus:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		return MakeNumber(left.Number() - right.Number()), nil

	case ast.NodeTypeBinaryDiv:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		return MakeNumber(left.Number() / right.Number()), nil

	case ast.NodeTypeBinaryMult:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		return MakeNumber(left.Number() * right.Number()), nil

	case ast.NodeTypeBinaryBitwiseAnd:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		res, err := builtinBitwiseAnd(left.Number(), right.Number())
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}
		return MakeNumber(res), nil

	case ast.NodeTypeBinaryBitwiseOr:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		res, err := builtinBitwiseOr(left.Number(), right.Number())
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}
		return MakeNumber(res), nil

	case ast.NodeTypeBinaryBitwiseXor:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		res, err := builtinBitwiseXor(left.Number(), right.Number())
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}
		return MakeNumber(res), nil

	case ast.NodeTypeBinaryShiftL:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		res, err := builtinShiftL(left.Number(), right.Number())
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}
		return MakeNumber(res), nil

	case ast.NodeTypeBinaryShiftR:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !left.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, left.Type())
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsNumber() {
			return ValueNone, TypeErrorSpecific(ValueTypeNumber, right.Type())
		}

		res, err := builtinShiftR(left.Number(), right.Number())
		if err != nil {
			return ValueNone, MakeRuntimeError(err)
		}
		return MakeNumber(res), nil

	case ast.NodeTypeBinaryUnequal:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		eq, err := left.Equal(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(!eq), nil

	case ast.NodeTypeBinaryGreater:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x > 0), nil

	case ast.NodeTypeBinaryGreaterEq:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x >= 0), nil

	case ast.NodeTypeBinaryLess:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x < 0), nil

	case ast.NodeTypeBinaryLessEq:
		left, err := EvaluateNode(tree, node.A, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		right, err := EvaluateNode(tree, node.B, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		x, err := left.Compare(right, ctx)
		if err != nil {
			return ValueNone, err
		}
		return MakeBool(x <= 0), nil

	}
}

// type GoCallbackNode struct {
// 	Func Function
// 	Args []NamedValue
// }

// func (n *GoCallbackNode) Context() ast.Context             { return nil }
// func (n *GoCallbackNode) Loc() *old_ast.LocationRange      { return nil }
// func (n *GoCallbackNode) FreeVariables() ast.Identifiers   { return nil }
// func (n *GoCallbackNode) SetFreeVariables(ast.Identifiers) {}
// func (n *GoCallbackNode) SetContext(ast.Context)           {}
// func (n *GoCallbackNode) OpenFodder() *old_ast.Fodder      { return nil }

type CallbackFunction struct{}

func handleDesugaredObject(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {

	fieldCount := len(node.Fields)
	localsCount := len(node.Locals)

	layer, _ := ctx.State.Registry.Layers.New()

	layer.ParentScopeId = scopeId

	if fieldCount > 0 {
		layer.Keys = ctx.State.Registry.Uint32Bufs.Alloc(0, fieldCount)
		layer.Nodes = ctx.State.Registry.NodesBufs.Alloc(0, fieldCount)
		layer.Meta = ctx.State.Registry.Uint8Bufs.Alloc(0, fieldCount)
	}

	if localsCount > 0 {
		layer.LocalKeys = ctx.State.Registry.Uint32Bufs.Alloc(0, localsCount)
		layer.LocalNodes = ctx.State.Registry.NodesBufs.Alloc(0, localsCount)
	}

	if len(node.Asserts) > 0 {
		layer.Asserts = ctx.State.Registry.NodesBufs.Alloc(len(node.Asserts), len(node.Asserts))
		copy(layer.Asserts, node.Asserts)
	}

	for _, v := range node.Locals {

		name := string(v.Variable)
		keyId := ctx.State.Interner.Intern(name)

		layer.LocalKeys = append(layer.LocalKeys, keyId)
		layer.LocalNodes = append(layer.LocalNodes, v.Body)

	}

	useMap := fieldCount > MaxLayerLinearKeys

	if useMap {
		layer.Index = make(map[uint32]int, fieldCount)
	}

	index := 0
	for _, v := range node.Fields {
		name, err := EvaluateNode(v.Name, scopeId, ctx)
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

		n := name.String(ctx)

		keyId := ctx.State.Interner.Intern(n)

		layer.Keys = append(layer.Keys, keyId)
		layer.Nodes = append(layer.Nodes, v.Body)
		layer.Meta = append(layer.Meta, CreateFieldMeta(v.Hide, v.PlusSuper))

		if useMap {
			layer.Index[keyId] = index
			index++
		}

	}

	layers := ctx.State.Registry.LayerBufs.Alloc(1, 1)
	layers[0] = layer
	objId := NewObject(layers, ctx)

	return MakeObjectValue(objId), nil
}

func handleLocal(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {

	s, childScopeId := ctx.NewScope(scopeId, int(node.C))

	binds := tree.SideTable[node.B : node.B+node.C*2]

	for i := range s.Bindings {
		keyId, bindNodeId := binds[i*2], binds[i*2+1]

		t, err := evaluateNodeLazy(tree, bindNodeId, childScopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		s.Bindings[i] = NamedValue{keyId, t}
	}

	val, err := EvaluateNode(tree, node.A, childScopeId, ctx)
	if err != nil {
		return ValueNone, err
	}
	return val, nil
}

func handleBinary(node *old_ast.Binary, scopeId uint32, ctx Context) (Value, error) {
	left, err := EvaluateNode(node.Left, scopeId, ctx)
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

		right, err := EvaluateNode(node.Right, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for && op, expected boolean", right.Type().String())
		}
		res := left.Bool() && right.Bool()
		return MakeBool(res), nil

	case ast.BopOr:
		if !left.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", left.Type().String())
		}

		if left.Bool() {
			return MakeBool(true), nil
		}

		right, err := EvaluateNode(node.Right, scopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		if !right.IsBool() {
			return ValueNone, fmt.Errorf("unexpected type %s for || op, expected boolean", right.Type().String())
		}
		res := left.Bool() || right.Bool()
		return MakeBool(res), nil
	default:
		right, err := EvaluateNode(node.Right, scopeId, ctx)
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

func handleFunction(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {

	paramCount := int(node.C)
	params := tree.SideTable[node.B : node.B+node.C*2]

	paramKeyIds := ctx.State.Registry.Uint32Bufs.Alloc(paramCount, paramCount)
	for i := range paramKeyIds {
		paramKeyIds[i] = params[i*2]
	}

	var fn Func = func(args []NamedValue, _ Context) (Value, error) {
		if len(args) > paramCount {
			return ValueNone, fmt.Errorf("unexpected amount of args passed to function")
		}

		if paramCount == 0 {
			return EvaluateNode(tree, node.A, scopeId, ctx)
		}

		s, childScopeId := ctx.NewScope(scopeId, paramCount)

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
			defArgNodeId := params[i*2+1]
			if defArgNodeId == 0 {
				return ValueNone, fmt.Errorf("arg (%d) with no default arg had no value passed", i)
			}

			da, err := evaluateNodeLazy(tree, defArgNodeId, childScopeId, ctx)
			if err != nil {
				return ValueNone, err
			}

			s.Bindings[i] = NamedValue{keyId, da}
		}

		return EvaluateNode(tree, node.A, childScopeId, ctx)
	}

	f := Function{
		argsCount: len(paramKeyIds),
		fn:        fn,
	}

	return MakeFunction(f, ctx), nil
}

func handleIndex(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {
	index, err := EvaluateNode(tree, node.A, scopeId, ctx)
	if err != nil {
		return ValueNone, err
	}

	target, err := EvaluateNode(tree, node.B, scopeId, ctx)
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

		keyId := index.AsStringConst()
		if keyId == 0 {
			// if its a dynamic string we need to intern it
			name := index.String(ctx)
			keyId = ctx.State.Interner.Intern(name)
		}

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

func handleSuperIndex(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {
	index, err := EvaluateNode(tree, node.A, scopeId, ctx)
	if err != nil {
		return ValueNone, err
	}

	if ctx.Self.IsNone() {
		return ValueNone, errors.New("ctx.Self not set")
	}

	keyId := index.AsStringConst()
	if keyId == 0 {
		// if its a dynamic string we need to intern it
		name := index.String(ctx)
		keyId = ctx.State.Interner.Intern(name)
	}

	obj := ctx.Self.Object(ctx)

	targetOffset := ctx.SuperOffset + 1

	val, _, err := obj.GetFieldWithOffset(keyId, ctx, int(targetOffset))
	if err != nil {
		return ValueNone, err
	}

	if val.IsNone() {
		name := index.String(ctx)
		return ValueNone, MakeRuntimeError(fmt.Errorf("Field does not exist: %s", name))
	}

	val, err = val.Eval(ctx)
	if err != nil {
		return ValueNone, err
	}

	return val, nil
}

func handleImport(tree *ast.AST, node ast.Node, scopeId uint32, ctx Context) (Value, error) {
	// fileVal, err := EvaluateNode(tree, node.File, scopeId, ctx)
	// if err != nil {
	// 	return ValueNone, err
	// }
	// if !fileVal.IsString() {
	// 	return ValueNone, fmt.Errorf("(%T) unexpected file data type '%s'", node, fileVal.Type().String())
	// }

	file := ctx.State.Interner.Get(node.A)
	fileLocation := ctx.State.Interner.Get(node.A)

	currentFileDir := filepath.Dir(fileLocation)

	var importedAst *ast.AST
	var finalPath string

	importer := ctx.State.Environment.Importer

	dirs := []string{""}
	if !filepath.IsAbs(file) {
		dirs = []string{currentFileDir}
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

		inTree, innerErr := importer.ResolveImport(fp)
		if os.IsNotExist(innerErr) {
			rangeErr = errors.Join(rangeErr, innerErr)
			continue
		}

		if innerErr != nil {
			return ValueNone, innerErr
		}

		importedAst = inTree
		finalPath = fp
		break

	}

	if importedAst == nil {
		return ValueNone, errors.Join(errors.New("error resolving import"), rangeErr)
	}

	// importScope := CreateFileScope(file, importer.BaseStd, ctx)
	importScope := ctx.State.Environment.BaseScopeId

	importCtx := ctx
	importCtx.Self = ValueNone

	v, err := evaluateNodeLazy(importedAst, importedAst.RootId, importScope, importCtx)
	if err != nil {
		return ValueNone, err
	}

	importer.Set(finalPath, v)

	return v, nil

}
