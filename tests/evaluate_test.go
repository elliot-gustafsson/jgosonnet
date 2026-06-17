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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elliot-gustafsson/jgosonnet"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEvaluator(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	// file := filepath.Join("resources", "test.jsonnet")
	file := "/home/elliot.gustafsson@fnox.it/Projects/jgosonnet/tests/resources/test.jsonnet"

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	jgosonnetStart := time.Now()
	stuff, err := interpreter.EvaluateJson(file)
	jgosonnetDur := time.Since(jgosonnetStart)
	assert.NoError(t, err)
	// if err != nil {
	// 	return
	// }

	println()
	println("jgosonnet:", jgosonnetDur.String())

	goJsonnetStart := time.Now()
	og, err := GetExpected(file, filepath.Join(infraDir, "vendor"))
	goJsonnetDur := time.Since(goJsonnetStart)
	assert.NoError(t, err)

	println("go-jsonnet:", goJsonnetDur.String())
	println()
	println(jgosonnetDur.String(), "/", goJsonnetDur.String(), "~", fmt.Sprintf("%.2f", GetChange(jgosonnetDur, goJsonnetDur)), "times faster")
	println()

	assert.Equal(t, og, stuff)

	println(stuff)

}

func TestJgosonnetYaml(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	file := filepath.Join("resources", "test.jsonnet")

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	// stuff, err := interpreter.EvaluateJson(file)
	stuff, err := interpreter.EvaluateYaml(file)
	assert.NoError(t, err)

	fmt.Println("jgosonnet:")
	fmt.Println()
	fmt.Print(stuff)

	fmt.Println()
	fmt.Println()

	fmt.Println("yaml v3:")
	fmt.Println()
	data := map[string]string{
		// "a":  "\n",
		// "a2": "\n\n",
		// "a3": "\n\nasdf",
		// "a4": "\n\n ",
		// "b":  "\ta",
		// "b2": "\ta\n",
		// "c":  " ",
		// "d":  " asdf",
		// "e":  " asdf\n",

		"0X123": "asdf",
		"0O123": "asdf",
		"0B123": "asdf",
		// "0b1":    "asdf",
		// "0o1":    "asdf",
		// "123_123.2": "asdf",
	}

	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)

	err = enc.Encode(data)
	assert.NoError(t, err)
	fmt.Print(b.String())

	assert.Equal(t, b.String(), string(stuff))

}

func TestYamlV3(t *testing.T) {
	data := map[string]string{
		"0O123": "banan",
	}

	out, err := yaml.Marshal(data)
	assert.NoError(t, err)

	fmt.Print(string(out))

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

func TestEvaluatorParallel(t *testing.T) {
	// curr := debug.SetGCPercent(-1)
	// defer func() {
	// 	debug.SetGCPercent(curr)
	// }()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	// x := []string{
	// 	"sto1-prod001",
	// 	"sto2-prod001",
	// 	"sto3-prod001",
	// }
	x := []string{
		"it-it001",
		"it-rancher",
		"sto1-acce001",
		"sto1-build001",
		"sto1-dev-analytics001",
		"sto1-dev001",
		"sto1-harvester001",
		"sto1-infra-edge001",
		"sto1-infra-public001",
		"sto1-infra001",
		"sto1-lb001",
		"sto1-prod-analytics001",
		"sto1-prod-gpu001",
		"sto1-prod001",
		"sto1-rancher",
		"sto2-acce001",
		"sto2-build001",
		"sto2-dev-gpu001",
		"sto2-dev001",
		"sto2-harvester001",
		"sto2-infra-edge001",
		"sto2-infra-public001",
		"sto2-infra001",
		"sto2-lb001",
		"sto2-prod-analytics001",
		"sto2-prod-gpu001",
		"sto2-prod001",
		"sto2-rancher",
		"sto3-acce001",
		"sto3-build-gpu001",
		"sto3-build001",
		"sto3-dev001",
		"sto3-harvester001",
		"sto3-infra-edge001",
		"sto3-infra-public001",
		"sto3-infra001",
		"sto3-lb001",
		"sto3-prod-analytics001",
		"sto3-prod-gpu001",
		"sto3-prod001",
		"sto3-rancher",
	}
	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	wg := sync.WaitGroup{}

	jgosonnetStart := time.Now()

	jobCh := make(chan string)

	for range 4 {
		wg.Go(func() {

			for c := range jobCh {

				fmt.Println("running", c)

				stuff, err := interpreter.EvaluateYamlMulti(filepath.Join(infraDir, c+".jsonnet"))
				assert.NoError(t, err)
				if err != nil {
					return
				}

				dir := filepath.Join(infraDir, "manifests", c)
				for k, v := range stuff {

					f, err := os.Create(filepath.Join(dir, k+".yaml"))
					assert.NoError(t, err)

					_, err = f.WriteString(v)
					assert.NoError(t, err)
				}

				fmt.Println("done", c)

			}
		})
	}

	for _, c := range x {
		jobCh <- c
	}
	close(jobCh)

	wg.Wait()

	jgosonnetDur := time.Since(jgosonnetStart)

	println()
	println("jgosonnet:", jgosonnetDur.String())

}

func TestEvaluatorSerial(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	// x := []string{
	// 	"sto1-prod001",
	// 	"sto2-prod001",
	// 	"sto3-prod001",
	// }
	x := []string{
		// "it-it001",
		// "it-rancher",
		// "mimir-alerts-dashboards",
		// "sto1-acce001",
		// "sto1-build001",
		// "sto1-dev-analytics001",
		// "sto1-dev001",
		// "sto1-harvester001",
		// "sto1-infra-edge001",
		// "sto1-infra-public001",
		// "sto1-infra001",
		// "sto1-lb001",
		// "sto1-prod-analytics001",
		// "sto1-prod-gpu001",
		"sto1-prod001",
		// "sto1-rancher",
		// "sto2-acce001",
		// "sto2-build001",
		// "sto2-dev-gpu001",
		// "sto2-dev001",
		// "sto2-harvester001",
		// "sto2-infra-edge001",
		// "sto2-infra-public001",
		// "sto2-infra001",
		// "sto2-lb001",
		// "sto2-prod-analytics001",
		// "sto2-prod-gpu001",
		"sto2-prod001",
		// "sto2-rancher",
		// "sto3-acce001",
		// "sto3-build-gpu001",
		// "sto3-build001",
		// "sto3-dev001",
		// "sto3-harvester001",
		// "sto3-infra-edge001",
		// "sto3-infra-public001",
		// "sto3-infra001",
		// "sto3-lb001",
		// "sto3-prod-analytics001",
		// "sto3-prod-gpu001",
		"sto3-prod001",
		// "sto3-rancher",
		// "vms",
	}
	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	wg := sync.WaitGroup{}

	jgosonnetStart := time.Now()

	for _, c := range x {
		fmt.Println("running", c)

		stuff, err := interpreter.EvaluateYamlMulti(filepath.Join(infraDir, c+".jsonnet"))
		assert.NoError(t, err)
		if err != nil {
			return
		}

		dir := filepath.Join(infraDir, "manifests", c)
		for k, v := range stuff {

			f, err := os.Create(filepath.Join(dir, k+".yaml"))
			assert.NoError(t, err)

			_, err = f.WriteString(v)
			assert.NoError(t, err)
		}

		fmt.Println("done", c)
	}

	wg.Wait()

	jgosonnetDur := time.Since(jgosonnetStart)

	println()
	println("jgosonnet:", jgosonnetDur.String())

}

func TestEvaluatorLoop(t *testing.T) {
	runtime.GOMAXPROCS(1)

	// originalGC := debug.SetGCPercent(-1)
	// defer debug.SetGCPercent(originalGC)

	slog.SetLogLoggerLevel(slog.LevelDebug)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	infraDir := filepath.Join(filepath.Dir(filepath.Dir(cwd)), "infra", "jsonnet", "proact")
	assert.NotEmpty(t, infraDir)

	interpreter := jgosonnet.NewEvaluator()
	interpreter.JPaths([]string{filepath.Join(infraDir, "vendor")})

	file := "sto3-prod001.jsonnet"

	_, err = interpreter.Evaluate(filepath.Join(infraDir, file))
	assert.NoError(t, err)
	if err != nil {
		return
	}

	f, err := os.Create("cpu.prof")
	assert.NoError(t, err)
	defer f.Close()
	err = pprof.StartCPUProfile(f)
	assert.NoError(t, err)
	defer pprof.StopCPUProfile()

	jgosonnetStart := time.Now()

	for range 30 {
		_, err := interpreter.Evaluate(filepath.Join(infraDir, file))
		assert.NoError(t, err)
		if err != nil {
			return
		}
	}

	jgosonnetDur := time.Since(jgosonnetStart)

	println()
	println("jgosonnet:", jgosonnetDur.String())

}
