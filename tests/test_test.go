package tests

import (
	"os"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	ctx := evaluator.Context{
		Interner: evaluator.NewInterner(),
		Registry: evaluator.NewRegistry(),
	}

	std, err := stdlib.InitStdLib(ctx)
	assert.NoError(t, err)

	astImporter := evaluator.NewAstImporter()

	env := &evaluator.Environment{
		TraceOut: os.Stdout,
		Importer: evaluator.NewImporter([]string{}, std, astImporter),
		ExtVars:  make(map[string]string),
		ExtCodes: make(map[string]string),
	}

	ctx.Environment = env

	// -------------------------------------------------

	filename := "main.jsonnet"
	snippet := "2 + 3"

	// scopeId := evaluator.CreateFileScope(filename, std, ctx)

	node, err := astImporter.ResolveSnippet(filename, snippet)
	assert.NoError(t, err)

	assert.NotNil(t, node)

	compiler := evaluator.Compiler{}

	program, err := compiler.Compile(node)
	assert.NoError(t, err)

	vm := evaluator.NewVM(program)

	result, err := vm.Run()
	assert.NoError(t, err)

	res, err := evaluator.ManifestValue(result, ctx)
	assert.NoError(t, err)

	assert.Equal(t, int64(5), res)
}
