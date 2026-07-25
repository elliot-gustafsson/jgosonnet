package jgosonnet

import (
	"fmt"
	"io"
	"iter"
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

	astImporter *evaluator.AstImporter
	extVars     map[string]string
	extCodes    map[string]string
	nativeFuncs map[string]evaluator.Function
}

type NativeFunction struct {
	Args map[string]any
	Fn   func(args []any) (any, error)
}

type FileOutput struct {
	Filename string
	Content  string
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
	b.Grow(16 * 1024)

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
	root, ctx, cleanup, err := t.evaluateMulti(file)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	c := &evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	res := make(map[string]string, len(root))
	for _, v := range root {

		var b strings.Builder
		b.Grow(16 * 1024)

		err = evaluator.ManifestJson(&b, v.Value, ctx, c)
		if err != nil {
			return nil, wrapManifestationErr(err)
		}

		b.WriteByte('\n')

		// Clone string due to string arena being reset in defer
		kClone := strings.Clone(ctx.State.Interner.Get(v.Key))

		res[kClone] = b.String()
	}

	return res, nil
}

// The caller MUST range over the returned iterator to execute the manifestation and release
// underlying evaluation resources back to the pool.
func (t *Evaluator) EvaluateJsonMultiIter(file string) (iter.Seq2[FileOutput, error], error) {
	root, ctx, cleanup, err := t.evaluateMulti(file)
	if err != nil {
		cleanup()
		return nil, err
	}

	c := &evaluator.JsonManifestConfig{
		IndentStep: "   ",
		Newline:    "\n",
		KeyValSep:  ": ",
		SpaceComma: true,
	}

	iterator := func(yield func(FileOutput, error) bool) {
		defer cleanup() // deferred until the caller finishes iterating

		for _, v := range root {
			var b strings.Builder
			b.Grow(16 * 1024)

			err := evaluator.ManifestJson(&b, v.Value, ctx, c)
			if err != nil {
				yield(FileOutput{}, wrapManifestationErr(err))
				return
			}
			b.WriteByte('\n')

			// Clone string due to string arena being reset in defer
			kClone := strings.Clone(ctx.State.Interner.Get(v.Key))

			out := FileOutput{
				Filename: kClone,
				Content:  b.String(),
			}

			if !yield(out, nil) {
				return // Caller broke early
			}
		}
	}

	return iterator, nil
}

func (t *Evaluator) EvaluateYaml(file string) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", wrapEvaluationErr(err)
	}

	var b strings.Builder
	b.Grow(16 * 1024)

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

func (t *Evaluator) EvaluateYamlMulti(file string) (map[string]string, error) {

	root, ctx, cleanup, err := t.evaluateMulti(file)
	defer cleanup()
	if err != nil {
		return nil, err
	}

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
		Modern:               true,
	}

	res := make(map[string]string, len(root))
	for _, v := range root {

		var b strings.Builder
		b.Grow(16 * 1024)

		err = evaluator.ManifestYaml(&b, v.Value, ctx, c)
		if err != nil {
			return nil, wrapManifestationErr(err)
		}

		b.WriteByte('\n')

		// Clone string due to string arena being reset in defer
		kClone := strings.Clone(ctx.State.Interner.Get(v.Key))

		res[kClone] = b.String()
	}

	return res, nil
}

// The caller MUST range over the returned iterator to execute the manifestation and release
// underlying evaluation resources back to the pool.
func (t *Evaluator) EvaluateYamlMultiIter(file string) (iter.Seq2[FileOutput, error], error) {

	root, ctx, cleanup, err := t.evaluateMulti(file)
	if err != nil {
		cleanup()
		return nil, err
	}

	c := evaluator.YamlManifestConfig{
		IndentArrayInObjects: true,
		NaturalSort:          true,
		FormatIntegers:       true,
		UseBlockScalars:      true,
		Modern:               true,
	}

	iterator := func(yield func(FileOutput, error) bool) {
		defer cleanup() // deferred until the caller finishes iterating

		for _, v := range root {
			var b strings.Builder
			b.Grow(16 * 1024)

			err := evaluator.ManifestYaml(&b, v.Value, ctx, c)
			if err != nil {
				yield(FileOutput{}, wrapManifestationErr(err))
				return
			}
			b.WriteByte('\n')

			// Clone string due to string arena being reset in defer
			kClone := strings.Clone(ctx.State.Interner.Get(v.Key))

			out := FileOutput{
				Filename: kClone,
				Content:  b.String(),
			}

			if !yield(out, nil) {
				return
			}
		}
	}

	return iterator, nil
}

type EvaluationEngine struct {
	Registry *evaluator.Registry
	Interner *interner.Interner
}

var enginePool = sync.Pool{
	New: func() any {
		return &EvaluationEngine{
			Registry: evaluator.NewRegistry(),
			Interner: interner.NewInterner(),
		}
	},
}

func (t *Evaluator) evaluate(file string) (evaluator.Value, evaluator.Context, func(), error) {

	node, err := t.astImporter.ResolveImport(file)
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
		engine.Interner.Reset()
		enginePool.Put(engine)

		// pprof.StopCPUProfile()
		// f.Close()
	}

	ctx := evaluator.Context{
		State: &evaluator.ContextState{
			Interner: engine.Interner,
			Registry: engine.Registry,
		},
	}

	std, err := stdlib.InitStdLib(ctx)
	if err != nil {
		return evaluator.ValueNone, evaluator.Context{}, cleanup, err
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

	ctx.State.Environment = env

	scopeId := evaluator.CreateFileScope(file, std, ctx)

	value, err := evaluator.EvaluateNode(node, scopeId, ctx)
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

func (t *Evaluator) evaluateMulti(file string) ([]evaluator.NamedValue, evaluator.Context, func(), error) {
	value, ctx, cleanup, err := t.evaluate(file)
	// defer cleanup()
	if err != nil {
		return nil, evaluator.Context{}, cleanup, wrapEvaluationErr(err)
	}

	if !value.IsObject() {
		return nil, evaluator.Context{}, cleanup, wrapEvaluationErr(evaluator.TypeErrorSpecific(evaluator.ValueTypeObject, value.Type()))
	}

	evalCtx := ctx
	evalCtx.Self = value

	root, err := evaluator.ManifestObjectRoot(value.Object(evalCtx), evalCtx)
	if err != nil {
		return nil, evaluator.Context{}, cleanup, wrapManifestationErr(err)
	}

	return root, ctx, cleanup, nil
}

func wrapEvaluationErr(err error) error {
	return fmt.Errorf("%w\tDuring evaluation", err)
}

func wrapManifestationErr(err error) error {
	return fmt.Errorf("%w\tDuring manifestation", err)
}
