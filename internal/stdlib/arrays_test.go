package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceString(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		start, end, step int
		expected         string
	}{
		{
			name:     "simple",
			input:    "123456789",
			start:    1,
			end:      3,
			step:     1,
			expected: "23",
		},
		{
			name:     "start greater than length",
			input:    "123456789",
			start:    100,
			end:      2,
			step:     1,
			expected: "",
		},
		{
			name:     "end greater than length",
			input:    "123456789",
			start:    1,
			end:      18,
			step:     1,
			expected: "23456789",
		},
		{
			name:     "negative end",
			input:    "123456789",
			start:    1,
			end:      -1,
			step:     1,
			expected: "2345678",
		},
		{
			name:     "very negative end",
			input:    "123456789",
			start:    1,
			end:      -100,
			step:     1,
			expected: "",
		},
		{
			name:     "simple with steps",
			input:    "123456789",
			start:    1,
			end:      8,
			step:     2,
			expected: "2468",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			res, err := sliceArr([]rune(tt.input), tt.start, tt.end, tt.step)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, string(res))
		})
	}
}
