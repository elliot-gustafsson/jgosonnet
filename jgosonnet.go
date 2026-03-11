package jgosonnet

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
)

type Evaluator struct {
	interner *evaluator.Interner
	jpaths   []string
	traceOut io.Writer
	// nativeFuncs map[string]*NativeFunction // TODO: make it possible to pass custom funcs
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		interner: evaluator.NewInterner(),
		traceOut: os.Stdout,
	}
}

func (t *Evaluator) JPaths(paths []string) {
	t.jpaths = paths
}

// Get output as a go struct, map[string]any || []any ...
func (t *Evaluator) Evaluate(file string) (any, error) {
	value, ctx, err := t.evaluate(file)
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
	value, ctx, err := t.evaluate(file)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	err = evaluator.ManifestJson(&b, value, ctx, "   ", "\n", ": ")
	if err != nil {
		return "", err
	}

	b.WriteByte('\n')

	return b.String(), nil
}

func (t *Evaluator) EvaluateYaml(file string) (string, error) {
	value, ctx, err := t.evaluate(file)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	err = evaluator.ManifestYaml(&b, value, ctx, true, false, false, true)
	if err != nil {
		return "", err
	}

	b.WriteByte('\n')

	return b.String(), nil
}

func (t *Evaluator) EvaluateYamlMulti(file string) (map[string]string, error) {
	value, ctx, err := t.evaluate(file)
	if err != nil {
		return nil, err
	}

	if !value.IsObject() {
		return nil, fmt.Errorf("root object must be of type object for EvaluateYamlMulti, got: %s", value.Type().String())
	}

	// value.Object(ctx)

	root, err := evaluator.ManifestObjectRoot(value.Object(ctx), ctx)
	if err != nil {
		return nil, err
	}

	evalCtx := ctx
	evalCtx.Self = value

	res := make(map[string]string, len(root))
	for key, v := range root {

		var b strings.Builder
		b.Grow(64 * 1024)

		err = evaluator.ManifestYaml(&b, v, evalCtx, true, false, false, true)
		if err != nil {
			return nil, err
		}

		b.WriteByte('\n')

		res[key] = b.String()
	}

	return res, nil
}

func (t *Evaluator) evaluate(file string) (evaluator.Value, evaluator.Context, error) {

	node, err := evaluator.ResolveImport(file)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, err
	}

	ctx := evaluator.Context{
		Interner: t.interner,
		Arena:    evaluator.NewArena(),
	}

	std, err := stdlib.InitStdLib(ctx)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, err
	}

	scopeId := ctx.Arena.NewScope(0, 2)
	ctx.Arena.AddScopeBind(scopeId, ctx.Interner.Intern("$std"), std)
	ctx.Arena.AddScopeBind(scopeId, ctx.Interner.Intern("std"), std)

	ctx.Importer = evaluator.NewImporter(scopeId, t.jpaths)

	value, err := evaluator.EvaluateNodeStrict(node, scopeId, ctx)
	if err != nil {
		return evaluator.Value{}, evaluator.Context{}, err
	}

	if value.IsFunction() {
		res, err := value.Function(ctx)(nil, ctx)
		if err != nil {
			return evaluator.Value{}, evaluator.Context{}, err
		}
		return res, ctx, nil
	}

	// println("Arrays", len(ctx.Arena.Arrays))
	// println("Objects", len(ctx.Arena.Objects))
	// println("Funcs", len(ctx.Arena.Funcs))
	// println("Thunks", len(ctx.Arena.Thunks))
	// println("Scopes", len(ctx.Arena.Scopes))

	return value, ctx, nil
}
