package benchmarks

import (
	"path/filepath"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet"
)

func BenchmarkLargeStringJoinYaml(b *testing.B) {
	file := filepath.Join("resources", "large_string_join.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkLargeStringTemplateYaml(b *testing.B) {
	file := filepath.Join("resources", "large_string_template.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkRealisticBenchmark1Yaml(b *testing.B) {
	file := filepath.Join("resources", "realistic_benchmark1.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkRealisticBenchmark2Yaml(b *testing.B) {
	file := filepath.Join("resources", "realistic_benchmark2.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkTailCallYaml(b *testing.B) {
	file := filepath.Join("resources", "tail_call.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkInheritenceRecursionYaml(b *testing.B) {
	file := filepath.Join("resources", "inheritence_recursion.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}

func BenchmarkComparisonsPrimitivesYaml(b *testing.B) {
	file := filepath.Join("resources", "comparisons_primitives.jsonnet")

	b.ResetTimer()

	for b.Loop() {
		ev := jgosonnet.NewEvaluator()
		_, err := ev.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}
}
