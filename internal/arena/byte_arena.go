package arena

import (
	"unsafe"
)

const (
	// ChunkSize defines the standard block size (64 KB).
	ChunkSize = 65536
	// MaxSmallAlloc defines the threshold above which allocations use the jumbo pool.
	MaxSmallAlloc = ChunkSize / 2
)

type ByteArena struct {
	chunks []*[ChunkSize]byte
	curr   int
	offset int
	jumbos [][]byte
}

func NewByteArena() (a *ByteArena) {
	a = &ByteArena{
		chunks: make([]*[ChunkSize]byte, 1, 64),
	}
	a.chunks[0] = new([ChunkSize]byte)
	return
}

func New[T any](a *ByteArena) (ptr *T) {
	var zero T
	size := int(unsafe.Sizeof(zero))
	align := int(unsafe.Alignof(zero))

	if size == 0 {
		var empty struct{}
		return (*T)(unsafe.Pointer(&empty))
	}

	ptr = (*T)(allocRaw(a, size, align))
	if ptr != nil {
		return
	}

	ptr = (*T)(allocRawSlow(a, size, align))
	return
}

func NewSlice[T any](a *ByteArena, length int) (s []T) {
	if length <= 0 {
		return nil
	}

	var zero T
	elemSize := int(unsafe.Sizeof(zero))
	align := int(unsafe.Alignof(zero))

	if elemSize == 0 {
		return make([]T, length)
	}

	totalSize := elemSize * length
	ptr := allocRaw(a, totalSize, align)
	if ptr != nil {
		s = unsafe.Slice((*T)(ptr), length)
		return
	}

	ptr = allocRawSlow(a, totalSize, align)
	s = unsafe.Slice((*T)(ptr), length)
	return
}

func (a *ByteArena) Reset() {
	for i := 0; i <= a.curr; i++ {
		limit := ChunkSize
		if i == a.curr {
			limit = int(a.offset)
		}
		clear(a.chunks[i][:limit])
	}

	a.curr = 0
	a.offset = 0

	// Release jumbo buffers to the Go garbage collector
	if len(a.jumbos) > 0 {
		a.jumbos = nil
	}
}

func allocRaw(a *ByteArena, size, align int) (ptr unsafe.Pointer) {
	alignedOffset := (a.offset + align - 1) &^ (align - 1)

	if alignedOffset+size > ChunkSize || size > MaxSmallAlloc {
		return nil
	}

	ptr = unsafe.Add(unsafe.Pointer(a.chunks[a.curr]), alignedOffset)
	a.offset = alignedOffset + size
	return
}

//go:noinline
func allocRawSlow(a *ByteArena, size, align int) unsafe.Pointer {
	if size > MaxSmallAlloc {
		// jumbo alloc
		buf := make([]byte, size)
		a.jumbos = append(a.jumbos, buf)
		return unsafe.Pointer(unsafe.SliceData(buf))
	}

	alignedOffset := (a.offset + align - 1) &^ (align - 1)
	if alignedOffset+size > ChunkSize {
		// not enough space in current block, add a new one
		a.curr++
		if a.curr >= len(a.chunks) {
			a.chunks = append(a.chunks, new([ChunkSize]byte))
		}
		a.offset = 0
		alignedOffset = 0
	}

	ptr := unsafe.Add(unsafe.Pointer(a.chunks[a.curr]), alignedOffset)
	a.offset = alignedOffset + size
	return ptr
}
