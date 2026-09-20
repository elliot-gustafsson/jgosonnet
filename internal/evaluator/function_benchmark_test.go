package evaluator

import (
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet/ast"
)

func setupBenchContext() Context {
	alloc := arena.NewAllocator()
	inter := interner.NewInterner()
	return Context{
		State: &ContextState{
			// MaxStack:  10000,
			Interner:  inter,
			Allocator: alloc,
		},
	}
}

func makeBenchFunc(ctx Context, paramNames []string, defaultArg ast.Node) Value {
	params := make([]ast.Parameter, len(paramNames))
	paramKeyIds := arena.Alloc[uint32](ctx.State.Allocator, len(paramNames))
	for i, name := range paramNames {
		params[i] = ast.Parameter{Name: ast.Identifier(name)}
		if i == len(paramNames)-1 && defaultArg != nil {
			params[i].DefaultArg = defaultArg
		}
		paramKeyIds[i] = ctx.State.Interner.Intern(name)
	}

	body := &ast.LiteralNull{}
	node := &ast.Function{
		Parameters: params,
		Body:       body,
	}

	f := arena.Create[Function](ctx.State.Allocator)
	arena.Memclr(f)
	*f = Function{
		Node:        node,
		ScopePtr:    0,
		ParamKeyIds: paramKeyIds,
	}
	return MakeFunctionValue(f, ctx)
}

func BenchmarkExecUserFunction_1Positional(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn := makeBenchFunc(ctx, []string{"x"}, nil)
		args := []NamedValue{{Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkExecUserFunction_2Positional(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn := makeBenchFunc(ctx, []string{"x", "y"}, nil)
		args := []NamedValue{{Value: MakeNumber(1)}, {Value: MakeNumber(2)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkExecUserFunction_DefaultArg(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn := makeBenchFunc(ctx, []string{"x", "y"}, &ast.LiteralNull{})
		args := []NamedValue{{Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkExecUserFunction_NamedArg(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ctx := setupBenchContext()
		fn := makeBenchFunc(ctx, []string{"x", "y"}, nil)
		yKey := ctx.State.Interner.Intern("y")
		xKey := ctx.State.Interner.Intern("x")
		args := []NamedValue{{Key: yKey, Value: MakeNumber(2)}, {Key: xKey, Value: MakeNumber(1)}}
		for j := 0; j < 1000; j++ {
			_, err := execUserFunction(fn, args, ctx, false)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
