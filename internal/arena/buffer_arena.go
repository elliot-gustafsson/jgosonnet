package arena

// BufferArena provides contiguous, reusable slice memory.
// It is designed purely to replace native Go make() calls for temporary buffers.
type BufferArena[T any] struct {
	elementBlocks [][]T
	blockSize     int
	activeIdx     int
	offset        int
}

func NewBufferArena[T any](elementBlockSize int) *BufferArena[T] {
	return &BufferArena[T]{
		elementBlocks: [][]T{make([]T, elementBlockSize)},
		blockSize:     elementBlockSize,
		activeIdx:     0,
		offset:        0,
	}
}

// Alloc mimics make([]T, length, capacity).
func (a *BufferArena[T]) Alloc(length, capacity int) []T {
	if capacity < length {
		capacity = length
	}
	if capacity == 0 {
		return nil
	}
	if capacity > a.blockSize {
		// If length is larger than blocksize just create it on the heap
		return make([]T, length, capacity)
	}

	currBlock := a.elementBlocks[a.activeIdx]
	if a.offset+capacity > len(currBlock) {
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

	// slice off the requested memory
	targetSlice := currBlock[a.offset : a.offset+length : a.offset+capacity]

	// move the arena offset forward by the capacity claimed
	a.offset += capacity

	return targetSlice
}

func (a *BufferArena[T]) Reset() {
	// Clear used memory so it can be gc:ed
	for i := 0; i <= a.activeIdx && i < len(a.elementBlocks); i++ {
		limit := a.blockSize
		if i == a.activeIdx {
			limit = a.offset
		}
		clear(a.elementBlocks[i][:limit])
	}
	a.activeIdx = 0
	a.offset = 0
}
