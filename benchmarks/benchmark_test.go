package benchmarks

import (
	"path/filepath"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet"
)

func BenchmarkLargeStringJoin(b *testing.B) {
	file := filepath.Join("resources", "large_string_join.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkLargeStringTemplate(b *testing.B) {
	file := filepath.Join("resources", "large_string_template.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkRealisticBenchmark1(b *testing.B) {
	file := filepath.Join("resources", "realistic_benchmark1.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkRealisticBenchmark2(b *testing.B) {
	file := filepath.Join("resources", "realistic_benchmark2.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkTailCall(b *testing.B) {
	file := filepath.Join("resources", "tail_call.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkInheritenceRecursion(b *testing.B) {
	file := filepath.Join("resources", "inheritence_recursion.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkComparisonsPrimitives(b *testing.B) {
	file := filepath.Join("resources", "comparisons_primitives.jsonnet")

	ev := jgosonnet.NewEvaluator()

	b.ResetTimer()

	for b.Loop() {
		_, err := ev.EvaluateJson(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
