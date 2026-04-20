package tests

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet"
	"github.com/stretchr/testify/assert"
)

func TestGoJsonnetTests(t *testing.T) {
	testsLoc := filepath.Join("resources", "go-jsonnet")

	err := os.Chdir(testsLoc)
	assert.NoError(t, err)

	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := filepath.Join("testdata", d.Name())

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(name) != ".jsonnet" {
			return nil
		}

		if strings.HasPrefix(name, "error.") {
			fmt.Println("skipping error test for now", name)
			return nil
		}

		expectedOutputFile := strings.TrimSuffix(name, ".jsonnet") + ".golden"
		_, err = os.Stat(expectedOutputFile)
		if err != nil {
			return nil
		}

		// fmt.Println(fi.Name())

		t.Run(name, func(t *testing.T) {

			eofInfo, err := os.Stat(expectedOutputFile)
			if err != nil {
				assert.FailNow(t, err.Error())
			}

			if eofInfo.IsDir() {
				fmt.Println("expectedOutputFile is a dir, skipping for now", name)
				return
			}

			expectedOut, err := os.ReadFile(expectedOutputFile)
			assert.NoError(t, err)

			if strings.HasPrefix(string(expectedOut), "RUNTIME ERROR:") {
				fmt.Println("skipping error test for now", name)
				return
			}

			ev := jgosonnet.NewEvaluator()

			out, err := ev.EvaluateJson(name)
			if err != nil {
				assert.Contains(t, string(expectedOut), err.Error())
				return
			}
			// assert.NoError(t, err)

			assert.Equal(t, string(expectedOut), out)
		})

		return nil

	})
}

func TestJsonnetCppTests(t *testing.T) {
	testsLoc := filepath.Join("resources", "jsonnet-cpp", "test_suite")

	err := os.Chdir(testsLoc)
	assert.NoError(t, err)

	filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		name := d.Name()

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(name) != ".jsonnet" {
			return nil
		}

		if strings.HasPrefix(name, "error.") {
			fmt.Println("skipping error test for now", name)
			return nil
		}

		expectedOutputFile := name + ".golden"
		_, err = os.Stat(expectedOutputFile)
		if err != nil {
			return nil
		}

		// fmt.Println(fi.Name())

		t.Run(name, func(t *testing.T) {
			ev := jgosonnet.NewEvaluator()

			ev.ExtVar("var1", "test")
			ev.ExtCode("var2", `{"x": 1, "y": 2}`)

			out, err := ev.EvaluateJson(name)
			assert.NoError(t, err)

			expectedOut, err := os.ReadFile(expectedOutputFile)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedOut), out)
		})

		return nil

	})
}

func TestSpecific(t *testing.T) {
	var testsLoc, name, expectedOutputFile string

	testsLoc = filepath.Join("resources", "jsonnet-cpp", "test_suite")
	name = "stdlib.jsonnet"
	expectedOutputFile = "stdlib.jsonnet.golden"
	// name = "trace.jsonnet"
	// expectedOutputFile = "trace.jsonnet.golden"

	// testsLoc = filepath.Join("resources", "go-jsonnet", "testdata")
	// // name = "stdlib_smoke_test.jsonnet"
	// // expectedOutputFile = "stdlib_smoke_test.golden"
	// name = "builtinObjectRemoveKey_super.jsonnet"
	// expectedOutputFile = "builtinObjectRemoveKey_super.golden"
	// name = "builtinObjectRemoveKey.jsonnet"
	// expectedOutputFile = "builtinObjectRemoveKey.golden"

	err := os.Chdir(testsLoc)
	assert.NoError(t, err)

	ev := jgosonnet.NewEvaluator()

	ev.ExtVar("var1", "test")
	ev.ExtCode("var2", `{"x": 1, "y": 2}`)

	out, err := ev.EvaluateJson(name)
	assert.NoError(t, err)

	expectedOut, err := os.ReadFile(expectedOutputFile)
	assert.NoError(t, err)

	assert.Equal(t, string(expectedOut), out)

}
