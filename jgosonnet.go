package jgosonnet

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
)

type Evaluator struct {
	jpaths   []string
	traceOut io.Writer

	interner    *interner.Interner
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
	interner := interner.NewInterner()
	return &Evaluator{
		traceOut:    os.Stderr,
		interner:    interner,
		astImporter: evaluator.NewAstImporter(interner),
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
		return nil, wrapEvaluationErr(err)
	}

	raw, err := evaluator.ManifestValue(value, ctx)
	if err != nil {
		return nil, wrapManifestationErr(err)
	}

	return raw, nil
}

func (t *Evaluator) EvaluateJson(file string) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", wrapEvaluationErr(err)
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	c := &evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	err = evaluator.ManifestJson(&b, value, ctx, c)
	if err != nil {
		return "", wrapManifestationErr(err)
	}

	b.WriteByte('\n')

	return b.String(), nil
}

func (t *Evaluator) EvaluateJsonMulti(file string) (map[string]string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return nil, wrapEvaluationErr(err)
	}

	if !value.IsObject() {
		return nil, wrapEvaluationErr(evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, value.Type()))
	}

	evalCtx := ctx
	evalCtx.Self = value

	root, err := evaluator.ManifestObjectRoot(value.Object(evalCtx), evalCtx)
	if err != nil {
		return nil, wrapManifestationErr(err)
	}

	c := &evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	res := make(map[string]string, len(root))
	for key, v := range root {

		var b strings.Builder
		b.Grow(64 * 1024)

		err = evaluator.ManifestJson(&b, v, evalCtx, c)
		if err != nil {
			return nil, wrapManifestationErr(err)
		}

		b.WriteByte('\n')

		// Clone string due to string arena being reset in defer
		kClone := strings.Clone(key)

		res[kClone] = b.String()
	}

	return res, nil
}

// NOTE: Maybe fully compliant, use with caution
func (t *Evaluator) EvaluateYaml(file string) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", wrapEvaluationErr(err)
	}

	var b strings.Builder
	b.Grow(1024 * 1024)

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
		Modern:               true,
	}

	err = evaluator.ManifestYaml(&b, value, ctx, c)
	if err != nil {
		return "", wrapManifestationErr(err)
	}

	b.WriteByte('\n')

	return b.String(), nil
}

// NOTE: Maybe fully compliant, use with caution
func (t *Evaluator) EvaluateYamlMulti(file string) (map[string]string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return nil, wrapEvaluationErr(err)
	}

	if !value.IsObject() {
		return nil, wrapEvaluationErr(evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, value.Type()))
	}

	evalCtx := ctx
	evalCtx.Self = value

	root, err := evaluator.ManifestObjectRoot(value.Object(evalCtx), evalCtx)
	if err != nil {
		return nil, wrapManifestationErr(err)
	}

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
		Modern:               true,
	}

	res := make(map[string]string, len(root))
	for key, v := range root {

		var b strings.Builder
		b.Grow(64 * 1024)

		err = evaluator.ManifestYaml(&b, v, evalCtx, c)
		if err != nil {
			return nil, wrapManifestationErr(err)
		}

		b.WriteByte('\n')

		// Clone string due to string arena being reset in defer
		kClone := strings.Clone(key)

		res[kClone] = b.String()
	}

	return res, nil
}

type EvaluationEngine struct {
	Registry *evaluator.Registry
}

var enginePool = sync.Pool{
	New: func() any {
		return &EvaluationEngine{
			Registry: evaluator.NewRegistry(),
		}
	},
}

func (t *Evaluator) evaluate(file string) (evaluator.Value, evaluator.Context, func(), error) {

	tree, err := t.astImporter.ResolveImport(file)
	if err != nil {
		return evaluator.ValueNone, evaluator.Context{}, func() {}, err
	}

	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)

	engine := enginePool.Get().(*EvaluationEngine)
	cleanup := func() {
		engine.Registry.Reset()
		enginePool.Put(engine)

		// pprof.StopCPUProfile()
		// f.Close()
	}

	ctx := evaluator.Context{
		State: &evaluator.ContextState{
			Registry: engine.Registry,
		},
	}

	std, err := stdlib.InitStdLib(ctx)
	if err != nil {
		return evaluator.ValueNone, evaluator.Context{}, cleanup, err
	}

	baseScopeId := evaluator.CreateFileScope(file, std, ctx)

	env := &evaluator.Environment{
		BaseScopeId:     baseScopeId,
		TraceOut:        t.traceOut,
		Importer:        evaluator.NewImporter(t.jpaths, std, t.astImporter),
		ExtVars:         t.extVars,
		ExtCodes:        t.extCodes,
		NativeFunctions: make(map[string]evaluator.Function, len(t.nativeFuncs)),
	}
	ctx.State.Environment = env

	// TODO: handle native funcs
	// for k, v := range t.nativeFuncs {
	// 	vars
	// 	env.NativeFunctions[k] = evaluator.Function{}
	// }

	value, err := evaluator.EvaluateNode(tree, tree.RootId, baseScopeId, ctx)
	if err != nil {
		return evaluator.ValueNone, evaluator.Context{}, cleanup, err
	}

	if value.IsFunction() {
		res, err := value.Function(ctx).Exec(nil, ctx)
		if err != nil {
			return evaluator.ValueNone, evaluator.Context{}, cleanup, err
		}
		return res, ctx, cleanup, nil
	}

	return value, ctx, cleanup, nil
}

func wrapEvaluationErr(err error) error {
	return fmt.Errorf("%w\tDuring evaluation", err)
}

func wrapManifestationErr(err error) error {
	return fmt.Errorf("%w\tDuring manifestation", err)
}
