package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		args        any // Can be []any, map[string]any, or single value
		expected    string
		expectedErr string
	}{
		// --- Basic Types ---
		{
			name:     "Simple String",
			format:   "Hello %s",
			args:     []any{"World"},
			expected: "Hello World",
		},
		{
			name:     "Integer",
			format:   "Value: %d",
			args:     []any{42},
			expected: "Value: 42",
		},
		{
			name:     "Integer Alias %i",
			format:   "%i",
			args:     []any{42},
			expected: "42",
		},
		{
			name:     "Unsigned Alias %u",
			format:   "%u",
			args:     []any{42},
			expected: "42",
		},
		{
			name:     "Float",
			format:   "Pi: %f",
			args:     []any{3.14159},
			expected: "Pi: 3.141590", // Default Go/C precision is 6
		},
		{
			name:     "Octal %o",
			format:   "%o",
			args:     []any{64}, // 64 in decimal is 100 in octal
			expected: "100",
		},
		{
			name:     "Hex Lower %x",
			format:   "%x",
			args:     []any{255},
			expected: "ff",
		},
		{
			name:     "Hex Upper %X",
			format:   "%X",
			args:     []any{255},
			expected: "FF",
		},
		{
			name:     "Char from Int",
			format:   "Char: %c",
			args:     []any{65},
			expected: "Char: A",
		},
		{
			name:     "Char from String",
			format:   "Char: %c",
			args:     []any{"A"},
			expected: "Char: A",
		},
		{
			name:     "Literal Percent",
			format:   "100%% sure",
			args:     []any{},
			expected: "100% sure",
		},

		// --- Flags & Padding ---
		{
			name:     "Left Align",
			format:   "|%-5s|",
			args:     []any{"foo"},
			expected: "|foo  |",
		},
		{
			name:     "Right Align (Default)",
			format:   "|%5s|",
			args:     []any{"foo"},
			expected: "|  foo|",
		},
		{
			name:     "Zero Pad Int",
			format:   "%05d",
			args:     []any{42},
			expected: "00042",
		},
		{
			name:     "Space Flag (Positive)",
			format:   "% d",
			args:     []any{42},
			expected: " 42",
		},
		{
			name:     "Space Flag (Negative)",
			format:   "% d",
			args:     []any{-42},
			expected: "-42",
		},
		{
			name:     "Plus Flag",
			format:   "%+d",
			args:     []any{42},
			expected: "+42",
		},
		{
			name:     "Alt Form Hex",
			format:   "%#x",
			args:     []any{255},
			expected: "0xff",
		},

		// --- Width & Precision ---
		{
			name:     "Precision Float",
			format:   "%.2f",
			args:     []any{3.14159},
			expected: "3.14",
		},
		{
			name:     "Width and Precision",
			format:   "%10.2f",
			args:     []any{3.14159},
			expected: "      3.14",
		},
		{
			name:     "Scientific Lower %e",
			format:   "%e",
			args:     []any{1000.0},
			expected: "1.000000e+03",
		},
		{
			name:     "Scientific Upper %E",
			format:   "%E",
			args:     []any{1000.0},
			expected: "1.000000E+03",
		},
		{
			name:     "General %g (Compact)",
			format:   "%g",
			args:     []any{123.456},
			expected: "123.456",
		},
		{
			name:     "General %g (Large Exponent)",
			format:   "%.2g",
			args:     []any{1234567.0},
			expected: "1.2e+06", // Go fmt behavior for %g switches to sci notation here
		},

		// --- Dynamic Width/Precision (*) ---
		{
			name:     "Dynamic Width",
			format:   "%*d",
			args:     []any{5, 10},
			expected: "   10",
		},
		{
			name:     "Dynamic Precision",
			format:   "%.*f",
			args:     []any{2, 1.23456},
			expected: "1.23",
		},
		{
			name:     "Dynamic Width & Precision",
			format:   "%*.*f",
			args:     []any{10, 2, 1.23456},
			expected: "      1.23",
		},

		// --- Named Arguments (Python Style) ---
		{
			name:   "Named String",
			format: "Hello %(name)s",
			args: map[string]any{
				"name": "Alice",
			},
			expected: "Hello Alice",
		},
		{
			name:   "Named Int & Float",
			format: "%(x)d + %(y).2f",
			args: map[string]any{
				"x": 10,
				"y": 20.555,
			},
			expected: "10 + 20.56",
		},
		{
			name:   "Named with Flags",
			format: "Score: %(score)05d",
			args: map[string]any{
				"score": 7,
			},
			expected: "Score: 00007",
		},
		{
			name:     "Mixed Named/Positional (Named in Format, List passed)",
			format:   "aaaa %(asdf)s %(asdf)s bbbb",
			args:     []any{"cccc", "dddd"},
			expected: "aaaa cccc dddd bbbb",
		},

		// --- Type Coercion ---
		{
			name:     "Float to Int (%d)",
			format:   "%d",
			args:     []any{3.9}, // Should truncate or cast
			expected: "3",
		},
		{
			name:     "Int to Float (%f)",
			format:   "%.1f",
			args:     []any{42},
			expected: "42.0",
		},
		{
			name:     "Bool to String",
			format:   "%s",
			args:     []any{true},
			expected: "true",
		},
		{
			name:     "Null to String",
			format:   "%s",
			args:     []any{nil},
			expected: "null",
		},

		// --- Single Value Wrapper ---
		{
			name:     "Single Value Input (Not Slice)",
			format:   "Val: %d",
			args:     123, // Pass raw int, should be wrapped to [123]
			expected: "Val: 123",
		},

		// --- Python/Jsonnet Specifics ---
		{
			name:     "Repr %r (Fallback to String)",
			format:   "Repr: %r",
			args:     []any{"test"},
			expected: "Repr: test",
		},
		{
			name:     "Float to Char (%c)",
			format:   "%c",
			args:     []any{65.0}, // Should cast to int 65 -> 'A'
			expected: "A",
		},
		{
			name:     "Complex Struct to String",
			format:   "Struct: %s",
			args:     []any{[]int{1, 2}}, // Should use fmt.Sprint
			expected: "Struct: [1 2]",
		},

		// --- Error Cases ---
		{
			name:        "Not Enough Arguments",
			format:      "%s %s",
			args:        []any{"One"},
			expectedErr: "not enough arguments for format string",
		},
		{
			name:        "Too Many Arguments",
			format:      "%s",
			args:        []any{"One", "Two"},
			expectedErr: "not all arguments converted during string formatting",
		},
		{
			name:        "Missing Key",
			format:      "%(missing)s",
			args:        map[string]any{"found": 1},
			expectedErr: "key 'missing' not found",
		},
		{
			name:        "Mixed Named/Positional (Positional in Format, Map passed)",
			format:      "%s",
			args:        map[string]any{"key": "val"},
			expectedErr: "format requires a mapping (%(key)s) when a dictionary is passed",
		},
		{
			name:        "Invalid Dynamic Width Type",
			format:      "%*d",
			args:        []any{"NOT INT", 5},
			expectedErr: "width requires integer",
		},
		{
			name:        "Incomplete Format String",
			format:      "Hello %",
			args:        []any{},
			expectedErr: "incomplete format string",
		},
		{
			name:        "Unknown Verb",
			format:      "%z",
			args:        []any{1},
			expectedErr: "unsupported format character 'z'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Format(tt.format, tt.args)
			if err != nil {
				assert.EqualError(t, err, tt.expectedErr)
				assert.Equal(t, "", got)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
