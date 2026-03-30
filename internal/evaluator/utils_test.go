package evaluator

import (
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNaturalStringSort(t *testing.T) {

	data := []string{
		"asdf100",
		"asdf20_v2",
		"Asdf20", // Uppercase 'A' comes before lowercase 'a'
		"asdf20",
		"asdf_20",
		"asdf-20",
		"asdf2",
		"asdf20-v10",
		"asdf1",
		"asdf01",
	}

	slices.SortFunc(data, naturalStringSort)

	assert.Equal(t, []string{
		"Asdf20",
		"asdf-20",
		"asdf_20",
		"asdf1",
		"asdf01",
		"asdf2",
		"asdf20",
		"asdf20-v10",
		"asdf20_v2",
		"asdf100",
	}, data)

}

func TestNaturalStringSort2(t *testing.T) {
	data := map[string]string{
		"asdf100":    "1",
		"asdf20_v2":  "1",
		"Asdf20":     "1", // Uppercase 'A' comes before lowercase 'a'
		"asdf20":     "1",
		"asdf_20":    "1",
		"asdf-20":    "1",
		"asdf2":      "1",
		"asdf20-v10": "1",
		"asdf1":      "1",
		"asdf01":     "1",
	}

	out, err := yaml.Marshal(data)
	assert.NoError(t, err)

	fmt.Println(string(out))
}
