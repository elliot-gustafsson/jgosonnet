package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdSlice(t *testing.T) {

	str := "123456789"

	res, err := StdSlice([]rune(str), 1, 18, 2)
	assert.NoError(t, err)

	assert.Equal(t, "2468", string(res))

}
