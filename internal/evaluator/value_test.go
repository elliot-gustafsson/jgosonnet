package evaluator

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// Value represents a NaN-boxed 64-bit value.
// Bit layout for non-float types:
//
// 63                   51 50    47 46                               0
// |----------------------|--------|---------------------------------|
// |  Float NaN Boundary  |  Type  |             Payload             |
// |       (13 bits)      |(4 bits)|            (47 bits)            |
// |----------------------|--------|---------------------------------|
//
// - Bits 63-51: Reserved for IEEE-754 Quiet NaN float64 detection.
// - Bits 50-47: ValueType enum.
// - Bits 46-0 : Payload (pointer info, interned string id, 1/0 for bools)

func (v Value) testIsString() bool {
	return ValueType(uint64(v)>>typeShift) == ValueTypeString
}

func (v Value) testIsNumber() bool {
	const floatThreshold uint64 = 1 << 51
	return uint64(v) >= floatThreshold
}

func (v Value) unboxNumber() float64 {
	const nanTag uint64 = 0xFFF8000000000000
	return math.Float64frombits(uint64(v) ^ nanTag)
}

func testMakeNumber(f float64) Value {
	const nanTag uint64 = 0xFFF8000000000000

	bits := math.Float64bits(f)

	if math.IsNaN(f) {
		bits = 0x7FF8000000000000
	}
	return Value(bits ^ nanTag)
}

func TestBoxing(t *testing.T) {
	someString := "hello"

	value := boxPtr(ValueTypeString, unsafe.Pointer(&someString))

	isString := value.testIsString()
	assert.Equal(t, true, isString)

	unboxedStringPtr := value.unboxPtr()
	unboxedString := *(*string)(unboxedStringPtr)

	assert.Equal(t, "hello", unboxedString)

	assert.NotNil(t, someString)
}

func TestBoxingNumber(t *testing.T) {
	someNum := 13.37

	value := testMakeNumber(someNum)

	isNum := value.testIsNumber()
	assert.Equal(t, true, isNum)

	unboxedNum := value.unboxNumber()

	assert.Equal(t, 13.37, unboxedNum)
}
