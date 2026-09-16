package evaluator

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
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
