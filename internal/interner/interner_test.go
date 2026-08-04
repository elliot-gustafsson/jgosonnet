package interner

import (
	"crypto/rand"
	"fmt"
	"runtime"
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
	runtime.GOMAXPROCS(1)

	for i := 0; i < b.N; i++ {
		interner := NewInterner()
		for _, k := range benchKeysUnique {
			interner.Intern(k)
		}
	}
}

func BenchmarkInterner_Repeated(b *testing.B) {
	runtime.GOMAXPROCS(1)

	interner := NewInterner()
	for _, k := range benchKeysRepeated {
		interner.Intern(k)
	}

	for i := 0; i < b.N; i++ {
		for _, k := range benchKeysRepeated {
			interner.Intern(k)
		}
	}
}
