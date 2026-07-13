package evaluator

import "fmt"

type Func = func(args []NamedValue, ctx Context) (Value, error)

type Function struct {
	AstId   uint32
	NodeId  uint32
	ScopeId uint32
}

func (t Function) Length(ctx Context) int {
	tree := ctx.State.Registry.ASTs[t.AstId]
	return int(tree.Nodes[t.NodeId].C)
}

func (t Function) Noop() bool {
	return t.NodeId == 0 || t.AstId == 0
}

func (t Function) Exec(args []NamedValue, ctx Context) (Value, error) {
	tree := ctx.State.Registry.ASTs[t.AstId]
	node := tree.Nodes[t.NodeId]

	paramCount := int(node.C)
	params := tree.SideTable[node.B : node.B+uint32(paramCount*2)]

	if len(args) > paramCount {
		return ValueNone, fmt.Errorf("unexpected amount of args passed to function")
	}

	if paramCount == 0 {
		return EvaluateNode(tree, node.A, t.ScopeId, ctx)
	}

	s, childScopeId := ctx.NewScope(t.ScopeId, paramCount)

	posIndex := 0
	for i := 0; i < paramCount; i++ {
		keyId := params[i*2]

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

		defArgNodeId := params[i*2+1]
		if defArgNodeId == 0 {
			return ValueNone, fmt.Errorf("arg (%d) with no default arg had no value passed", i)
		}

		// Evaluate default arg in the new child scope
		da, err := evaluateNodeLazy(tree, defArgNodeId, childScopeId, ctx)
		if err != nil {
			return ValueNone, err
		}

		s.Bindings[i] = NamedValue{Key: keyId, Value: da}
	}

	// Make sure we pass the correct AstId in the context when jumping into the function body!
	execCtx := ctx
	execCtx.AstId = t.AstId

	return EvaluateNode(tree, node.A, childScopeId, execCtx)
}

var GlobalNativeFunctions []NativeFunction

type NativeFunction struct {
	Params   []string
	OptStart int
	Func     Func
}

func (t NativeFunction) Exec(args []NamedValue, ctx Context) (Value, error) {

	// TODO: Argument x already provided

	paramCount := len(t.Params)

	if len(args) > paramCount {
		return ValueNone, MakeRuntimeError(fmt.Errorf("function expected %d positional argument(s), but got %d", len(t.Params), len(args)))
	}

	var stackArgs [4]NamedValue // Max args for any stdlib function is 4
	var orderedArgs []NamedValue

	if paramCount <= len(stackArgs) {
		orderedArgs = stackArgs[:paramCount]
	} else {
		// Just in case a user defines a custom native extension with 5+ args
		orderedArgs = make([]NamedValue, paramCount)
	}

	var onNamedArgs bool
	posIdx := 0

	for _, na := range args {
		if na.Key == 0 { // Positional Argument
			if onNamedArgs {
				return ValueNone, fmt.Errorf("positional argument after a named argument is not allowed")
			}
			orderedArgs[posIdx] = na
			posIdx++
			continue
		}

		// Named Argument
		onNamedArgs = true
		found := false

		// Look up the string once from the caller's interner
		passedName := ctx.State.Interner.Get(na.Key)

		for j, paramName := range t.Params {
			if passedName == paramName { // Raw Go string comparison
				orderedArgs[j] = na
				found = true
				break
			}
		}
		if !found {
			return ValueNone, MakeRuntimeError(fmt.Errorf("function has no parameter %s", passedName))
		}
	}

	// Check required arguments
	for i := 0; i < t.OptStart; i++ {
		if orderedArgs[i].Value.IsNone() {
			return ValueNone, MakeRuntimeError(fmt.Errorf("missing argument: %s", t.Params[i]))
		}
	}

	return t.Func(orderedArgs, ctx)
}
