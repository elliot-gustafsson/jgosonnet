package arena

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type DummyStruct struct {
	A, B, C, D uint64
}

var dummySize = unsafe.Sizeof(DummyStruct{})
var dummyAlign = unsafe.Alignof(DummyStruct{})

// --- 4. Benchmarks ---

func BenchmarkGenericArena_Alloc(b *testing.B) {
	runtime.GOMAXPROCS(1)

	a := NewArena[DummyStruct]()
	b.ResetTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x, _ := a.New()
		_ = x.A
	}
}

func BenchmarkByteArena_Alloc(b *testing.B) {
	runtime.GOMAXPROCS(1)

	a := NewAllocator()
	b.ResetTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x := Create[DummyStruct](a)
		*x = DummyStruct{}
		_ = x.A
	}
}

// const accessCount = 10000

// func BenchmarkGenericArena_Access(b *testing.B) {
// 	a := NewArena[DummyStruct]()
// 	ids := make([]uint32, accessCount)
// 	for i := range accessCount {
// 		_, id := a.New()
// 		ids[i] = id
// 	}

// 	for i := 0; b.Loop(); i++ {
// 		id := ids[i%accessCount]

// 		// The math required to dereference an ID inside GenericArena
// 		blockIdx := id / blockSize
// 		offsetIdx := id % blockSize
// 		ptr := &a.blocks[blockIdx][offsetIdx]

// 		_ = ptr.A
// 	}
// }

// func BenchmarkByteArena_Access(b *testing.B) {
// 	a := NewByteArena()
// 	ptrs := make([]*DummyStruct, accessCount)
// 	for i := range accessCount {
// 		ptrs[i] = (*DummyStruct)(allocRaw(a, dummySize, dummyAlign))
// 	}

// 	for i := 0; b.Loop(); i++ {
// 		// A direct pointer dereference (What SpiderMonkey NaN boxing allows)
// 		ptr := ptrs[i%accessCount]
// 		_ = ptr.A
// 	}
// }

type SomeStruct struct {
	x int
	y float64
}

func TestByteArena(t *testing.T) {
	a := NewAllocator()

	z := Create[SomeStruct](a)

	slice := Alloc[uint32](a, 65536)

	// slice[2] = 100

	z.x = 10
	z.y = 10.23

	assert.Equal(t, &SomeStruct{
		10, 10.23,
	}, z)

	// assert.Equal(t, []uint32{0, 0, 100, 0, 0}, slice)

	assert.NotNil(t, slice)
}
