package arena

type SliceArena[T any] struct {
	elementBlocks [][]T
	blockSize     int
	activeIdx     int
	offset        int

	headers *Arena[[]T]
}

func NewSliceArena[T any](elementBlockSize int) *SliceArena[T] {
	return &SliceArena[T]{
		elementBlocks: [][]T{make([]T, elementBlockSize)},
		blockSize:     elementBlockSize,
		activeIdx:     0,
		offset:        0,

		headers: NewArena[[]T](),
	}
}

func (a *SliceArena[T]) Alloc(items []T) (id uint32) {
	length := len(items)
	slice, id := a.Make(length)
	copy(slice, items)
	return
}

func (a *SliceArena[T]) Make(length int) ([]T, uint32) {
	if length == 0 {
		return nil, a.headers.Alloc(nil)
	}

	if length > a.blockSize {
		// If length is larger than blocksize just create it on the heap
		jumboBlock := make([]T, length)
		return jumboBlock, a.headers.Alloc(jumboBlock)
	}

	if a.offset+length > a.blockSize {
		a.grow()
	}

	currBlock := a.elementBlocks[a.activeIdx]
	targetSlice := currBlock[a.offset : a.offset+length : a.offset+length]
	a.offset += length

	return targetSlice, a.headers.Alloc(targetSlice)
}

func (a *SliceArena[T]) Get(id uint32) []T {
	return a.headers.GetValue(id)
}

func (a *SliceArena[T]) Reset() {
	// clear used memory so it can be gc:ed
	for i := 0; i <= a.activeIdx && i < len(a.elementBlocks); i++ {
		limit := a.blockSize
		if i == a.activeIdx {
			limit = a.offset
		}
		clear(a.elementBlocks[i][:limit])
	}

	a.activeIdx = 0
	a.offset = 0

	a.headers.Reset()
}

//go:noinline
func (a *SliceArena[T]) grow() {
	a.activeIdx++
	a.offset = 0
	if a.activeIdx < len(a.elementBlocks) {
		return
	}
	a.elementBlocks = append(a.elementBlocks, make([]T, a.blockSize))
}
