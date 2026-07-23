package tests

import (
	"fmt"
	"testing"
)

func TestQuick(t *testing.T) {
	str, err := GetExpected("resources/parse.jsonnet")
	fmt.Printf("Expected: %v, err: %v\n", str, err)
}
