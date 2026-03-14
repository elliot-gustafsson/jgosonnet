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
	// interner *evaluator.Interner
	jpaths   []string
	traceOut io.Writer

	astImporter *evaluator.AstImporter
	// nativeFuncs map[string]*NativeFunction // TODO: make it possible to pass custom funcs
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		// interner:    evaluator.NewInterner(),
		traceOut:    os.Stdout,
		astImporter: evaluator.NewAstImporter(),
	}
}

func (t *Evaluator) JPaths(paths []string) {
	t.jpaths = paths
}

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
		SingleQuoteEscape:    true,
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

	root, err := evaluator.ManifestObjectRoot(value.Object(ctx), ctx)
	if err != nil {
		return nil, err
	}

	evalCtx := ctx
	evalCtx.Self = value

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		SingleQuoteEscape:    true,
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

	scopeId := ctx.Arena.NewScope(0, 2)
	ctx.Arena.AddScopeBind(scopeId, ctx.Interner.Intern("$std"), std)
	ctx.Arena.AddScopeBind(scopeId, ctx.Interner.Intern("std"), std)

	ctx.Importer = evaluator.NewImporter(scopeId, t.jpaths, t.astImporter)

	value, err := evaluator.EvaluateNodeStrict(node, scopeId, ctx)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, cleanup, err
	}

	if value.IsFunction() {
		res, err := value.Function(ctx)(nil, ctx)
		if err != nil {
			return evaluator.Value{}, evaluator.Context{}, cleanup, err
		}
		return res, ctx, cleanup, nil
	}

	return value, ctx, cleanup, nil
}
