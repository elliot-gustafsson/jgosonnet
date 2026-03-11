package benchmarks

import (
	"path/filepath"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet"
)

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
