package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elliot-gustafsson/jgosonnet"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFloats(t *testing.T) {
	var myInt64 int64 = 9223372036854774784
	var myFloat64 float64 = math.Float64frombits(uint64(myInt64))

	// fmt.Println(int64(math.Float64bits(myFloat64)))

	assert.Equal(t, int64(9223372036854774784), int64(math.Float64bits(myFloat64)))

	d, err := yaml.Marshal(map[string]string{
		"!asdf asdf": "asdf",
	})
	assert.NoError(t, err)
	fmt.Println(string(d))
}

func TestEvaluator(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	file := filepath.Join("resources", "test.jsonnet")
	// file := filepath.Join(infraDir, "mimir-alerts-dashboards.jsonnet")

	// data, err := json.Marshal(9223372036854774784)
	// assert.NoError(t, err)
	// fmt.Println(string(data))

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	jgosonnetStart := time.Now()
	stuff, err := interpreter.EvaluateJson(file)
	jgosonnetDur := time.Since(jgosonnetStart)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	println()
	println("jgosonnet:", jgosonnetDur.String())

	// println("Actual:")
	// println(string(raw))
	// println("")

	goJsonnetStart := time.Now()
	og, err := GetExpected(file, filepath.Join(infraDir, "vendor"))
	goJsonnetDur := time.Since(goJsonnetStart)
	assert.NoError(t, err)

	// println("Expected:")
	// println(string(og))
	// println("")

	// newFile, err := os.Create("/tmp/jgosonnet")
	// assert.NoError(t, err)

	// newFile.WriteString(stuff)

	// oldFile, err := os.Create("/tmp/go-jsonnet")
	// assert.NoError(t, err)

	// oldFile.WriteString(og)

	println("go-jsonnet:", goJsonnetDur.String())
	println()
	println(jgosonnetDur.String(), "/", goJsonnetDur.String(), "~", fmt.Sprintf("%.2f", GetChange(jgosonnetDur, goJsonnetDur)), "times faster")
	println()

	assert.Equal(t, og, stuff)
	// if og != stuff {
	// 	assert.FailNow(t, "output not equal")
	// 	return
	// }

	// println("--- out ---")
	println(stuff)
	// println("--- end ---")
	// println()
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

/*

var b strings.Builder
	b.Grow(len(s) + 8)
	// b.WriteByte('"')
	b.WriteByte(escapeChar)
	for i := 0; i < len(s); i++ {
		c := s[i]

		if c == escapeChar {

		}
	}

*/

func TestEvaluatorReal(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	file := filepath.Join(infraDir, "sto1-prod001.jsonnet")

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	jgosonnetStart := time.Now()
	stuff, err := interpreter.EvaluateYamlMulti(file)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	jgosonnetDur := time.Since(jgosonnetStart)

	println()
	println("jgosonnet:", jgosonnetDur.String())

	dir := filepath.Join(infraDir, "manifests", "sto1-prod001")
	for k, v := range stuff {

		// err := os.WriteFile(filepath.Join(dir, k), []byte(v), 0600)
		f, err := os.Create(filepath.Join(dir, k+".yaml"))
		assert.NoError(t, err)

		// enc := yaml.NewEncoder(f)
		// enc.SetIndent(2)
		// err = enc.Encode(&v)
		// assert.NoError(t, err)
		// err = f.Close()

		_, err = f.WriteString(v)
		assert.NoError(t, err)

		assert.NoError(t, err)
	}

}
