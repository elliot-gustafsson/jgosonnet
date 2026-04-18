package jgosonnet

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
)

type Evaluator struct {
	jpaths   []string
	traceOut io.Writer

	astImporter *evaluator.AstImporter
	extVars     map[string]string
	extCodes    map[string]string
	nativeFuncs map[string]evaluator.Function
}

type NativeFunction struct {
	Args map[string]any
	Fn   func(args []any) (any, error)
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		traceOut:    os.Stderr,
		astImporter: evaluator.NewAstImporter(),
		extVars:     make(map[string]string),
		extCodes:    make(map[string]string),
		nativeFuncs: make(map[string]evaluator.Function),
	}
}

func (t *Evaluator) JPaths(paths []string) {
	t.jpaths = paths
}

func (t *Evaluator) TraceOut(w io.Writer) {
	t.traceOut = w
}

func (t *Evaluator) ExtVar(key, val string) {
	t.extVars[key] = val
}

func (t *Evaluator) ExtCode(key, val string) {
	t.extCodes[key] = val
}

// func (t *Evaluator) NativeFunction(key string, f NativeFunction) {
// 		t.nativeFuncs[key] = f
// }

// Get output as a go struct, map[string]any || []any ...
func (t *Evaluator) Evaluate(file string) (any, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	raw, err := evaluator.ManifestValue(value, ctx)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (t *Evaluator) EvaluateJson(file string) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	c := evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	err = evaluator.ManifestJson(&b, value, ctx, c)
	if err != nil {
		return "", err
	}

	b.WriteByte('\n')

	return b.String(), nil
}

func (t *Evaluator) EvaluateYaml(file string) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
	}

	err = evaluator.ManifestYaml(&b, value, ctx, c)
	if err != nil {
		return "", err
	}

	b.WriteByte('\n')

	return b.String(), nil
}

func (t *Evaluator) EvaluateYamlMulti(file string) (map[string]string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	if !value.IsObject() {
		return nil, fmt.Errorf("root object must be of type object, got: %s", value.Type().String())
	}

	evalCtx := ctx
	evalCtx.Self = value

	root, err := evaluator.ManifestObjectRoot(value.Object(evalCtx), evalCtx)
	if err != nil {
		return nil, err
	}

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
	}

	res := make(map[string]string, len(root))
	for key, v := range root {

		var b strings.Builder
		b.Grow(64 * 1024)

		err = evaluator.ManifestYaml(&b, v, evalCtx, c)
		if err != nil {
			return nil, err
		}

		b.WriteByte('\n')

		res[key] = b.String()
	}

	return res, nil
}

var arenaPool = sync.Pool{
	New: func() any {
		return evaluator.NewArena()
	},
}

func (t *Evaluator) evaluate(file string) (evaluator.Value, evaluator.Context, func(), error) {

	node, err := t.astImporter.ResolveImport(file)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, func() {}, err
	}

	arena := arenaPool.Get().(*evaluator.Arena)
	cleanup := func() {
		arena.Reset()
		arenaPool.Put(arena)
	}

	ctx := evaluator.Context{
		Interner: evaluator.NewInterner(),
		Arena:    arena,
	}

	std, err := stdlib.InitStdLib(ctx)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, cleanup, err
	}

	env := &evaluator.Environment{
		TraceOut:        t.traceOut,
		Importer:        evaluator.NewImporter(t.jpaths, std, t.astImporter),
		ExtVars:         t.extVars,
		ExtCodes:        t.extCodes,
		NativeFunctions: make(map[string]evaluator.Function, len(t.nativeFuncs)),
	}

	// TODO: handle native funcs
	// for k, v := range t.nativeFuncs {
	// 	vars
	// 	env.NativeFunctions[k] = evaluator.Function{}
	// }

	ctx.Environment = env

	scopeId := evaluator.CreateFileScope(file, std, ctx)

	value, err := evaluator.EvaluateNode(node, scopeId, ctx)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, cleanup, err
	}

	if value.IsFunction() {
		res, err := value.Function(ctx).Exec(nil, ctx)
		if err != nil {
			return evaluator.Value{}, evaluator.Context{}, cleanup, err
		}
		return res, ctx, cleanup, nil
	}

	return value, ctx, cleanup, nil
}
