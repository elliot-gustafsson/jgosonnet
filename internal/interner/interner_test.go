package interner

import (
	"crypto/rand"
	"fmt"
	"testing"
)

var (
	benchKeysUnique   []string
	benchKeysRepeated []string
)

func init() {
	benchKeysUnique = make([]string, 10000)
	for i := range benchKeysUnique {
		b := make([]byte, 10)
		rand.Read(b)
		benchKeysUnique[i] = fmt.Sprintf("key_%x", b)
	}

	benchKeysRepeated = make([]string, 100000)
	for i := range benchKeysRepeated {
		benchKeysRepeated[i] = benchKeysUnique[i%100]
	}
}

func BenchmarkInterner_Unique(b *testing.B) {
	for b.Loop() {
		interner := NewInterner()
		for _, k := range benchKeysUnique {
			interner.Intern(k)
		}
	}
}

func BenchmarkInterner_Repeated(b *testing.B) {
	interner := NewInterner()
	for _, k := range benchKeysRepeated {
		interner.Intern(k)
	}

	for b.Loop() {
		for _, k := range benchKeysRepeated {
			interner.Intern(k)
		}
	}
}
