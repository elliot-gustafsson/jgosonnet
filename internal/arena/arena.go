package arena

const blockSize = 4096

type Arena[T any] struct {
	blocks  []*[blockSize]T
	current int
	offset  int
}

func NewArena[T any]() *Arena[T] {
	return &Arena[T]{
		blocks:  []*[blockSize]T{{}},
		current: 0,
		offset:  1, // Burn index 0 so it can be used as a "nil" check
	}
}

func (a *Arena[T]) Alloc(val T) (id uint32) {
	if a.offset >= blockSize {
		a.grow()
	}

	a.blocks[a.current][a.offset] = val

	id = uint32(a.current*blockSize + a.offset)
	a.offset++
	return
}

func (a *Arena[T]) New() (ptr *T, id uint32) {
	if a.offset >= blockSize {
		a.grow()
	}

	ptr = &a.blocks[a.current][a.offset]

	id = uint32(a.current*blockSize + a.offset)
	a.offset++
	return
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
	a.offset = 1 // Re-burn index 0 so it can be used as a "nil" check
}

//go:noinline
func (a *Arena[T]) grow() {
	a.current++
	a.offset = 0
	if a.current >= len(a.blocks) {
		a.blocks = append(a.blocks, new([blockSize]T))
	}
}
