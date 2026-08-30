package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/elliot-gustafsson/jgosonnet"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	file := filepath.Join("resources", "test.jsonnet")

	interpreter := jgosonnet.NewEvaluator()

	jgosonnetStart := time.Now()
	stuff, err := interpreter.EvaluateJson(file)
	jgosonnetDur := time.Since(jgosonnetStart)
	if err != nil {
		t.Fatal(err.Error())
	}

	println()
	println("jgosonnet:", jgosonnetDur.String())

	goJsonnetStart := time.Now()
	og, err := GetExpected(file)
	goJsonnetDur := time.Since(goJsonnetStart)
	if err != nil {
		t.Fatal(err.Error())
	}

	println("go-jsonnet:", goJsonnetDur.String())
	println()
	println(jgosonnetDur.String(), "/", goJsonnetDur.String(), "~", fmt.Sprintf("%.2f", GetChange(jgosonnetDur, goJsonnetDur)), "times faster")
	println()

	assert.Equal(t, og, stuff)

	println("expected")
	println(og)
	println("actual")
	println(stuff)

}

func GetExpected(file string, jpaths ...string) (string, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{
		JPaths: jpaths,
	})

	node, _, err := vm.ImportAST("", file)
	if err != nil {
		return "", err
	}

	og, err := vm.Evaluate(node)
	if err != nil {
		return "", err
	}

	// return PrettifyJson(og), nil

	return og, nil
}

func DePrettifyJson(t *testing.T, val string) string {
	dst := &bytes.Buffer{}

	if err := json.Compact(dst, []byte(val)); err != nil {
		assert.FailNowf(t, "error compacting json, err: %s", err.Error())
		return ""
	}

	return dst.String()
}

func PrettifyJson(val string) string {
	var data any

	// 3. Unmarshal the ugly JSON into the interface
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	// 4. Marshal it back with indentation
	// "" is the prefix (usually left empty)
	// "  " is the indent (2 spaces is standard, or use "\t" for tab)
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling JSON: %v", err)
	}

	return string(prettyJSON)
}

func GetChange(old, new time.Duration) float64 {
	if old == 0 {
		return 0.0 // Avoid division by zero
	}

	// Convert both to float64 (nanoseconds) for precise division
	diff := float64(new - old)
	baseline := float64(old)

	return (diff / baseline)
}

func BenchmarkEvaluatorLoop(b *testing.B) {
	b.ReportAllocs()

	runtime.GOMAXPROCS(1)

	// originalGC := debug.SetGCPercent(-1)
	// defer debug.SetGCPercent(originalGC)

	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	if err != nil {
		b.Fatal(err.Error())
	}

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	file := filepath.Join(infraDir, "sto3-prod001.jsonnet")
	// file := "../benchmarks/resources/realistic_benchmark2.jsonnet"

	_, err = interpreter.Evaluate(file)
	if err != nil {
		b.Fatal(err.Error())
	}

	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal(err.Error())
	}
	defer cpuFile.Close()
	err = pprof.StartCPUProfile(cpuFile)
	if err != nil {
		b.Fatal(err.Error())
	}
	defer pprof.StopCPUProfile()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		b.Fatalf("could not create memory profile: %v", err)
	}
	defer memFile.Close()

	// jgosonnetStart := time.Now()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := interpreter.EvaluateYaml(file)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	// jgosonnetDur := time.Since(jgosonnetStart)

	runtime.GC()
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		b.Fatalf("could not write memory profile: %v", err)
	}

	// println()
	// println("jgosonnet:", jgosonnetDur.String())

}
