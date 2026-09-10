package jgosonnet

import (
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
	"sync"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/evaluator"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/elliot-gustafsson/jgosonnet/internal/stdlib"
)

var jsonManifestConfig = &evaluator.JsonManifestConfig{
	IndentStep: "   ",
	Newline:    "\n",
	KeyValSep:  ": ",
	SpaceComma: true,
}

var yamlManifestConfig = evaluator.YamlManifestConfig{
	IndentArrayInObjects: true,
	NaturalSort:          true,
	FormatIntegers:       true,
	UseBlockScalars:      true,
	Modern:               true,
}

type Evaluator struct {
	jpaths   []string
	traceOut io.Writer

	astImporter *evaluator.AstImporter
	extVars     map[string]string
	extCodes    map[string]string
	tlaVars     map[string]string
	tlaCodes    map[string]string
	nativeFuncs map[string]evaluator.Function

	noNewline bool
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
		tlaVars:     make(map[string]string),
		tlaCodes:    make(map[string]string),
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

func (t *Evaluator) TLAVar(key, val string) {
	t.tlaVars[key] = val
}

func (t *Evaluator) TLACode(key, val string) {
	t.tlaCodes[key] = val
}

// func (t *Evaluator) NativeFunction(key string, f NativeFunction) {
// 		t.nativeFuncs[key] = f
// }

func (t *Evaluator) NoNewline(v bool) {
	t.noNewline = v
}

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

// JSON
func (t *Evaluator) EvaluateJson(file string) (string, error) {
	return t.manifestSingle(file, formatJSON)
}
func (t *Evaluator) EvaluateJsonMulti(file string) (map[string]string, error) {
	return t.manifestMulti(file, formatJSON)
}
func (t *Evaluator) EvaluateJsonMultiIter(file string) (iter.Seq2[FileOutput, error], error) {
	return t.manifestMultiIter(file, formatJSON)
}

// YAML
func (t *Evaluator) EvaluateYaml(file string) (string, error) {
	return t.manifestSingle(file, formatYAML)
}
func (t *Evaluator) EvaluateYamlMulti(file string) (map[string]string, error) {
	return t.manifestMulti(file, formatYAML)
}
func (t *Evaluator) EvaluateYamlMultiIter(file string) (iter.Seq2[FileOutput, error], error) {
	return t.manifestMultiIter(file, formatYAML)
}

// String
func (t *Evaluator) EvaluateString(file string) (string, error) {
	return t.manifestSingle(file, formatString)
}
func (t *Evaluator) EvaluateStringMulti(file string) (map[string]string, error) {
	return t.manifestMulti(file, formatString)
}
func (t *Evaluator) EvaluateStringMultiIter(file string) (iter.Seq2[FileOutput, error], error) {
	return t.manifestMultiIter(file, formatString)
}

type manifestFormat uint8

const (
	formatJSON manifestFormat = iota
	formatYAML
	formatString
)

func formatValue(b *strings.Builder, val evaluator.Value, ctx evaluator.Context, fmtType manifestFormat) error {
	switch fmtType {
	case formatJSON:
		return evaluator.ManifestJson(b, val, ctx, jsonManifestConfig)
	case formatYAML:
		return evaluator.ManifestYaml(b, val, ctx, yamlManifestConfig)
	case formatString:
		if !val.IsString() {
			return evaluator.TypeErrorSpecific(evaluator.ValueTypeString, val.Type())
		}
		b.WriteString(val.String(ctx))
		return nil
	}
	return nil
}

func (t *Evaluator) manifestSingle(file string, fmtType manifestFormat) (string, error) {
	value, ctx, cleanup, err := t.evaluate(file)
	defer cleanup()
	if err != nil {
		return "", wrapEvaluationErr(err)
	}

	var b strings.Builder
	b.Grow(16 * 1024)

	if err := formatValue(&b, value, ctx, fmtType); err != nil {
		return "", wrapManifestationErr(err)
	}

	if !t.noNewline {
		b.WriteByte('\n')
	}

	return b.String(), nil
}

func (t *Evaluator) manifestMultiIter(file string, fmtType manifestFormat) (iter.Seq2[FileOutput, error], error) {
	root, ctx, cleanup, err := t.evaluateMulti(file)
	if err != nil {
		cleanup()
		return nil, err
	}

	iterator := func(yield func(FileOutput, error) bool) {
		defer cleanup()

		for _, v := range root {
			var b strings.Builder
			b.Grow(16 * 1024)

			if err := formatValue(&b, v.Value, ctx, fmtType); err != nil {
				yield(FileOutput{}, wrapManifestationErr(err))
				return
			}

			if !t.noNewline {
				b.WriteByte('\n')
			}

			kClone := strings.Clone(ctx.State.Interner.Get(v.Key))
			if !yield(FileOutput{Filename: kClone, Content: b.String()}, nil) {
				return
			}
		}
	}

	return iterator, nil
}

func (t *Evaluator) manifestMulti(file string, fmtType manifestFormat) (map[string]string, error) {
	it, err := t.manifestMultiIter(file, fmtType)
	if err != nil {
		return nil, err
	}

	res := make(map[string]string)
	for out, err := range it {
		if err != nil {
			return nil, err
		}
		res[out.Filename] = out.Content
	}
	return res, nil
}

type EvaluationEngine struct {
	Allocator *arena.Allocator
	Interner  *interner.Interner
}

var enginePool = sync.Pool{
	New: func() any {
		return &EvaluationEngine{
			Allocator: arena.NewAllocator(),
			Interner:  interner.NewInterner(),
		}
	},
}

func (t *Evaluator) evaluate(file string) (evaluator.Value, evaluator.Context, func(), error) {

	node, err := t.astImporter.ResolveImport(file, file)
	if err != nil {
		return evaluator.ValueNone, evaluator.Context{}, func() {}, evaluator.MakeRuntimeError(fmt.Errorf("%w\n", err))
	}

	// f, err := os.Create("cpu.prof")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(f)

	engine := enginePool.Get().(*EvaluationEngine)
	cleanup := func() {
		engine.Allocator.Reset()
		engine.Interner.Reset()
		enginePool.Put(engine)

		// pprof.StopCPUProfile()
		// f.Close()
	}

	ctx := evaluator.Context{
		State: &evaluator.ContextState{
			Interner:  engine.Interner,
			Allocator: engine.Allocator,
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
		var tlaArgs []evaluator.NamedValue
		tlaCount := len(t.tlaVars) + len(t.tlaCodes)

		if tlaCount > 0 {
			tlaArgs = make([]evaluator.NamedValue, 0, tlaCount)
			for k, v := range t.tlaVars {
				tlaArgs = append(tlaArgs, evaluator.NamedValue{
					Key:   ctx.State.Interner.Intern(k),
					Value: evaluator.MakeString(v, ctx),
				})
			}
			for k, code := range t.tlaCodes {
				codeNode, err := t.astImporter.ResolveSnippet("<tla:"+k+">", code)
				if err != nil {
					return evaluator.ValueNone, evaluator.Context{}, cleanup, err
				}
				val, err := evaluator.EvaluateNode(codeNode, scopeId, ctx)
				if err != nil {
					return evaluator.ValueNone, evaluator.Context{}, cleanup, err
				}
				tlaArgs = append(tlaArgs, evaluator.NamedValue{
					Key:   ctx.State.Interner.Intern(k),
					Value: val,
				})
			}
		}

		res, err := value.FunctionExec(tlaArgs, ctx)
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

	root, err := evaluator.ManifestObjectRoot(value.Object(), evalCtx)
	if err != nil {
		return nil, evaluator.Context{}, cleanup, wrapManifestationErr(err)
	}

	return root, ctx, cleanup, nil
}

func manifestJSON(b *strings.Builder, val evaluator.Value, ctx evaluator.Context) error {
	return evaluator.ManifestJson(b, val, ctx, jsonManifestConfig)
}

func manifestYAML(b *strings.Builder, val evaluator.Value, ctx evaluator.Context) error {
	return evaluator.ManifestYaml(b, val, ctx, yamlManifestConfig)
}

func manifestString(b *strings.Builder, val evaluator.Value, ctx evaluator.Context) error {
	if !val.IsString() {
		return evaluator.TypeErrorSpecific(evaluator.ValueTypeString, val.Type())
	}
	b.WriteString(val.String(ctx))
	return nil
}

func wrapEvaluationErr(err error) error {
	return fmt.Errorf("%w\tDuring evaluation", err)
}

func wrapManifestationErr(err error) error {
	return fmt.Errorf("%w\tDuring manifestation", err)
}
