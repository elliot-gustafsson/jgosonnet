package evaluator

const blockSize = 1024

type Arena[T any] struct {
	blocks  []*[blockSize]T
	current int
	offset  int
}

func NewArena[T any]() *Arena[T] {
	return &Arena[T]{
		blocks: []*[blockSize]T{{}},
	}
}

func (a *Arena[T]) Alloc(val T) uint32 {
	if a.offset >= blockSize {
		a.blocks = append(a.blocks, new([blockSize]T))
		a.current++
		a.offset = 0
	}

	a.blocks[a.current][a.offset] = val

	id := uint32(a.current*blockSize + a.offset)
	a.offset++
	return id
}

func (a *Arena[T]) GetPtr(id uint32) *T {
	blockIdx := id / blockSize
	offsetIdx := id % blockSize
	return &a.blocks[blockIdx][offsetIdx]
}

func (a *Arena[T]) GetValue(id uint32) T {
	blockIdx := id / blockSize
	offsetIdx := id % blockSize
	return a.blocks[blockIdx][offsetIdx]
}

func (a *Arena[T]) Reset() {

	for i := 0; i <= a.current; i++ {
		limit := blockSize
		if i == a.current {
			limit = a.offset
		}

		clear(a.blocks[i][:limit])
	}

	a.current = 0
	a.offset = 0
}
