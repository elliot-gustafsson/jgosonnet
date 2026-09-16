package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"log/slog"
	"path/filepath"
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

	// jgosonnetStart := time.Now()
	stuff, err := interpreter.EvaluateJson(file)
	// jgosonnetDur := time.Since(jgosonnetStart)
	if err != nil {
		t.Fatal(err.Error())
	}

	// println()
	// println("jgosonnet:", jgosonnetDur.String())

	// goJsonnetStart := time.Now()
	og, err := GetExpected(file)
	// goJsonnetDur := time.Since(goJsonnetStart)
	if err != nil {
		t.Fatal(err.Error())
	}

	// println("go-jsonnet:", goJsonnetDur.String())
	// println()
	// println(jgosonnetDur.String(), "/", goJsonnetDur.String(), "~", fmt.Sprintf("%.2f", GetChange(jgosonnetDur, goJsonnetDur)), "times faster")
	// println()

	assert.Equal(t, og, stuff)

	if og != stuff {
		println("expected")
		println(og)
		println("actual")
		println(stuff)
	}

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
