package benchmarks

import (
	"os/exec"
	"path/filepath"
	"testing"
)

var binaries = []struct {
	name string
	path string
}{
	{"go-jsonnet", "jsonnet"},
	{"jrsonnet", "jrsonnet"},
	{"jgosonnet", "../cmd/jgosonnet/jgosonnet"},
}

var testFiles = []string{
	"large_string_join.jsonnet",
	"large_string_template.jsonnet",
	"realistic_benchmark1.jsonnet",
	"realistic_benchmark2.jsonnet",
	"tail_call.jsonnet",
	"inheritence_recursion.jsonnet",
	"comparisons_primitives.jsonnet",
}

func BenchmarkBinaries(b *testing.B) {
	for _, bin := range binaries {
		cmdPath, err := exec.LookPath(bin.path)
		if err != nil {
			b.Logf("Skipping %s: binary %q not found", bin.name, bin.path)
			continue
		}

		for _, file := range testFiles {
			b.Run(bin.name+"_"+file, func(b *testing.B) {
				filePath := filepath.Join("resources", file)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					cmd := exec.Command(cmdPath, filePath)
					err := cmd.Run()
					if err != nil {
						b.Logf("Execution failed for %s on %s: %v", bin.name, file, err)
					}
				}
			})
		}
	}
}
