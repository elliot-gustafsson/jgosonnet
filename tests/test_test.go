package tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestTest(t *testing.T) {
	interner := evaluator.NewInterner()
	registry := evaluator.NewRegistry()

	ctx := evaluator.Context{
		Interner: interner,
		Registry: registry,
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

	fileName := "resources/test.jsonnet"
	fileData, err := os.ReadFile(fileName)
	assert.NoError(t, err)

	// scopeId := evaluator.CreateFileScope(filename, std, ctx)

	node, err := astImporter.ResolveSnippet(fileName, string(fileData))
	assert.NoError(t, err)

	assert.NotNil(t, node)

	globalCtx := &evaluator.GlobalContext{
		Interner:     ctx.Interner,
		ProgramCache: &evaluator.ProgramCache{},
	}

	compiler := evaluator.NewCompiler(globalCtx)
	program, err := compiler.Compile(node)
	assert.NoError(t, err)

	vm := evaluator.NewVM(interner, registry)

	result, err := vm.Run(program)
	assert.NoError(t, err)

	res, err := vm.ManifestValue(result)
	assert.NoError(t, err)

	b, err := json.MarshalIndent(res, "", "  ")
	assert.NoError(t, err)

	fmt.Println()
	fmt.Println(string(b))
	fmt.Println()
}
