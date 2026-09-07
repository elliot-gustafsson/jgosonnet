package alloc

import (
	"runtime"
	"testing"
)

type DummyStruct struct {
	A, B, C, D uint64
}

func BenchmarkByteArena_Alloc(b *testing.B) {
	runtime.GOMAXPROCS(1)

	a := NewAllocator()
	b.ResetTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x := a.Create[DummyStruct]()
		Memclr(x)
		_ = x.A
	}
}
