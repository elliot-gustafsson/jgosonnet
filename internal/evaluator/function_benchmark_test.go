package evaluator

import (
	"strings"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	fast "github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

func setupBenchContext() Context {
	alloc := arena.NewAllocator()
	inter := interner.NewInterner()
	return Context{
		State: &ContextState{
			MaxStack:  10000,
			Interner:  inter,
			Allocator: alloc,
		},
	}
}

func makeBenchFunc(ctx Context, paramNames []string, hasDefaultOnLast bool) (Value, fast.AST) {
	parts := make([]string, len(paramNames))
	copy(parts, paramNames)
	if hasDefaultOnLast && len(parts) > 0 {
		parts[len(parts)-1] = parts[len(parts)-1] + "=null"
	}
	snippet := "function(" + strings.Join(parts, ", ") + ") null"
	a, rootIdx, err := fast.ParseSnippet("bench.jsonnet", snippet, ctx.State.Interner)
	if err != nil {
		panic(err)
	}

	f := arena.Create[Function](ctx.State.Allocator)
	arena.Memclr(f)
	*f = Function{
		Bp:       a.Base(),
		Offset:   rootIdx,
		ScopePtr: 0,
	}
	return MakeFunctionValue(f, ctx), a
}

func BenchmarkExecUserFunction_1Positional(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn, a := makeBenchFunc(ctx, []string{"x"}, false)
		args := []NamedValue{{Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
		_ = a
	}
}

func BenchmarkExecUserFunction_2Positional(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn, a := makeBenchFunc(ctx, []string{"x", "y"}, false)
		args := []NamedValue{{Value: MakeNumber(1)}, {Value: MakeNumber(2)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
		_ = a
	}
}

func BenchmarkExecUserFunction_DefaultArg(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn, a := makeBenchFunc(ctx, []string{"x", "y"}, true)
		args := []NamedValue{{Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
		_ = a
	}
}

func BenchmarkExecUserFunction_NamedArg(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn, a := makeBenchFunc(ctx, []string{"x", "y"}, false)
		yKey := ctx.State.Interner.Intern("y")
		xKey := ctx.State.Interner.Intern("x")
		args := []NamedValue{{Key: yKey, Value: MakeNumber(2)}, {Key: xKey, Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
		_ = a
	}
}
