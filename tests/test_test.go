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
		State: &evaluator.ContextState{
			Interner: interner,
			Registry: registry,
		},
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

	ctx.State.Environment = env

	// -------------------------------------------------

	fileName := "resources/test.jsonnet"
	fileData, err := os.ReadFile(fileName)
	assert.NoError(t, err)

	// scopeId := evaluator.CreateFileScope(filename, std, ctx)

	node, err := astImporter.ResolveSnippet(fileName, string(fileData))
	assert.NoError(t, err)

	assert.NotNil(t, node)

	globalCtx := &evaluator.GlobalContext{
		Interner:     ctx.State.Interner,
		ProgramCache: &evaluator.ProgramCache{},
	}

	compiler := evaluator.NewCompiler(globalCtx)
	program, err := compiler.Compile(node)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println("Instructions:")
	for _, v := range program.Instructions {
		fmt.Println(v)
	}
	fmt.Println()

	vm := evaluator.NewVM(interner, registry)

	result, err := vm.Run(program)
	if err != nil {
		t.Fatal(err.Error())
	}

	res, err := vm.ManifestValue(result)
	if err != nil {
		t.Fatal(err.Error())
	}

	b, err := json.MarshalIndent(res, "", "  ")
	assert.NoError(t, err)

	fmt.Println()
	fmt.Println(string(b))
	fmt.Println()
}
