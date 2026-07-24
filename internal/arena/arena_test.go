package arena

import (
	"testing"
)

// go test -gcflags="-d=ssa/check_bce" -bench=. ./internal/arena

func BenchmarkArena_GetPtr(b *testing.B) {
	a := NewArena[int]()

	const numElements = 1 << 16
	ids := make([]uint32, numElements)
	for i := range numElements {
		ids[i] = a.Alloc(i)
	}

	var sum int
	for i := 0; b.Loop(); i++ {
		id := ids[i&(numElements-1)]
		sum += *a.GetPtr(id)
	}

	if sum == -1 {
		b.Fatal("impossible")
	}
}

func BenchmarkArena_GetValue(b *testing.B) {
	a := NewArena[int]()

	const numElements = 1 << 16
	ids := make([]uint32, numElements)
	for i := range numElements {
		ids[i] = a.Alloc(i)
	}

	var sum int
	for i := 0; b.Loop(); i++ {
		id := ids[i&(numElements-1)]
		sum += a.GetValue(id)
	}

	if sum == -1 {
		b.Fatal("impossible")
	}
}
