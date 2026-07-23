package tests

import (
	"path/filepath"
	"strings"
	"testing"
	
	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
)

func TestParseSuite(t *testing.T) {
	files, err := filepath.Glob("resources/go-jsonnet/testdata/*.jsonnet")
	if err != nil {
		t.Fatal(err)
	}
	
	passed := 0
	failed := 0
	
	for _, file := range files {
		interner := interner.NewInterner()
		_, err := ast.Parse(file, interner)
		if err != nil {
			failed++
			base := filepath.Base(file)
			if !strings.Contains(base, "comp") {
				t.Logf("Failed parsing %s: %v", base, err)
			}
		} else {
			passed++
		}
	}
	
	t.Logf("Passed: %d, Failed: %d", passed, failed)
}
