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

func (a *SliceArena[T]) Alloc(items []T) uint32 {
	length := len(items)
	var targetSlice []T

	if length == 0 {
		return a.headers.Alloc(targetSlice)
	}

	if length > a.blockSize {
		// If length is larger than blocksize just create it on the heap
		jumboBlock := make([]T, length)
		copy(jumboBlock, items)
		return a.headers.Alloc(jumboBlock)
	}

	currBlock := a.elementBlocks[a.activeIdx]
	if a.offset+length > len(currBlock) {
		a.activeIdx++
		a.offset = 0
		// do we already have a block allocated?
		if a.activeIdx < len(a.elementBlocks) {
			currBlock = a.elementBlocks[a.activeIdx]
		} else {
			currBlock = make([]T, a.blockSize)
			a.elementBlocks = append(a.elementBlocks, currBlock)
		}
	}

	targetSlice = currBlock[a.offset : a.offset+length : a.offset+length]
	copy(targetSlice, items)
	a.offset += length

	return a.headers.Alloc(targetSlice)
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
