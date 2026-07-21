package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {

	interner := interner.NewInterner()

	tree, err := ast.Parse("resources/parse.jsonnet", interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	str := evalTree(t, tree, interner)
	fmt.Println()
	fmt.Print(str)
	fmt.Println()

	expectedStr, err := GetExpected("resources/parse.jsonnet")
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, expectedStr, str)
}

func evalTree(t *testing.T, tree *ast.AST, interner *interner.Interner) string {
	registry := evaluator.NewRegistry()

	ctx := evaluator.Context{
		State: &evaluator.ContextState{
			Registry: registry,
			Interner: interner,
		},
	}

	std, err := stdlib.InitStdLib(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	baseScopeId := evaluator.CreateFileScope("resources/parse.jsonnet", std, ctx)

	env := &evaluator.Environment{
		BaseScopeId: baseScopeId,
	}
	ctx.State.Environment = env

	ctx.AstId = uint32(len(ctx.State.Registry.ASTs))
	ctx.State.Registry.ASTs = append(ctx.State.Registry.ASTs, tree)

	value, err := evaluator.EvaluateNode(tree, tree.RootId, baseScopeId, ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	c := &evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	var b strings.Builder

	err = evaluator.ManifestJson(&b, value, ctx, c)
	if err != nil {
		t.Fatal(err.Error())
	}

	b.WriteByte('\n')

	return b.String()
}
