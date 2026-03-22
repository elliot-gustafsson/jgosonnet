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

func TestCompare(t *testing.T) {

	file := filepath.Join("resources", "test.jsonnet")

	ev := jgosonnet.NewEvaluator()

	out, err := ev.EvaluateJson(file)
	assert.NoError(t, err)

	exp, err := GetExpected(file)
	assert.NoError(t, err)

	assert.Equal(t, exp, out)

	// if exp != out {
	// 	println("--------------------")
	// 	println("expected:")
	// 	println("--------------------")
	// 	println(out)
	// 	println("--------------------")

	// 	println("actual:")
	// 	println("--------------------")
	// 	println(exp)
	// 	println("--------------------")
	// }
}

func TestStuff(t *testing.T) {

	// file := filepath.Join("resources", "jsonnet-cpp", "test_suite", "array_comparison.jsonnet")
	file := filepath.Join("resources", "jsonnet-cpp", "test_suite", "stdlib.jsonnet")

	ev := jgosonnet.NewEvaluator()

	out, err := ev.EvaluateJson(file)
	assert.NoError(t, err)

	expectedOut, err := os.ReadFile(file + ".golden")
	assert.NoError(t, err)

	assert.Equal(t, string(expectedOut), out)
}

func TestJsonnet(t *testing.T) {
	testsLoc := filepath.Join("resources", "jsonnet-cpp", "test_suite")

	filepath.WalkDir(testsLoc, func(path string, d fs.DirEntry, err error) error {
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

		// fmt.Println(name)
		inputFile := filepath.Join(testsLoc, name)

		expectedOutputFile := filepath.Join(testsLoc, name) + ".golden"
		_, err = os.Stat(expectedOutputFile)
		if err != nil {
			return nil
		}

		// fmt.Println(fi.Name())

		t.Run(name, func(t *testing.T) {
			ev := jgosonnet.NewEvaluator()

			out, err := ev.EvaluateJson(inputFile)
			assert.NoError(t, err)

			expectedOut, err := os.ReadFile(expectedOutputFile)
			assert.NoError(t, err)

			assert.Equal(t, string(expectedOut), out)
		})

		return nil

	})
}
