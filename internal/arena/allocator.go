package arena

import (
	"unsafe"
)

const (
	KB = 1024
	// ChunkSize defines the standard block size.
	ChunkSize = 256 * KB
	// MaxSmallAlloc defines the threshold above which allocations use the jumbo pool.
	MaxSmallAlloc = 32 * KB
)

type Allocator struct {
	base   unsafe.Pointer
	chunks []*[ChunkSize]byte
	curr   int
	offset uintptr
	jumbos [][]byte
}

func NewAllocator() (a *Allocator) {
	a = &Allocator{
		chunks: make([]*[ChunkSize]byte, 1, 64),
	}
	a.chunks[0] = new([ChunkSize]byte)
	a.base = unsafe.Pointer(a.chunks[0])
	return
}

func Create[T any](a *Allocator) (ptr *T) {
	var zero T
	size := unsafe.Sizeof(zero)
	align := unsafe.Alignof(zero)

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

func Alloc[T any](a *Allocator, length int) (s []T) {
	if length <= 0 {
		return nil
	}

	var zero T
	elemSize := unsafe.Sizeof(zero)
	align := unsafe.Alignof(zero)

	if elemSize == 0 {
		var empty struct{}
		return unsafe.Slice((*T)(unsafe.Pointer(&empty)), length)
	}

	totalSize := elemSize * uintptr(length)
	ptr := allocRaw(a, totalSize, align)
	if ptr != nil {
		s = unsafe.Slice((*T)(ptr), length)
		return
	}

	ptr = allocRawSlow(a, totalSize, align)
	s = unsafe.Slice((*T)(ptr), length)
	return
}

// Realloc takes an existing slice and a new desired length. It allocates a
// new slice in the arena, copies over the existing values, and returns it.
// If the requested length is less than or equal to the current length, it
// simply returns the original slice.
func Realloc[T any](a *Allocator, slice []T, length int) (s []T) {
	// If the slice already has enough capacity, we can just reslice it.
	if length <= cap(slice) {
		return slice[:length]
	}

	var zero T
	elemSize := unsafe.Sizeof(zero)

	if elemSize == 0 {
		var empty struct{}
		return unsafe.Slice((*T)(unsafe.Pointer(&empty)), length)
	}

	if cap(slice) > 0 {
		// calculate the memory address at the end of the slice's capacity
		slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
		capBytes := uintptr(cap(slice)) * elemSize
		sliceEnd := uintptr(slicePtr) + capBytes

		// calculate the memory address of the arena's current offset
		arenaEnd := uintptr(a.base) + a.offset

		// if they match, this slice was the very last allocation in the current chunk.
		if uintptr(slicePtr) >= uintptr(a.base) && sliceEnd == arenaEnd {
			totalSize := uintptr(length) * elemSize
			additionalBytes := totalSize - capBytes
			newOffset := a.offset + additionalBytes

			// ensure it still fits in the chunk and respects the max small alloc threshold
			if totalSize <= MaxSmallAlloc && newOffset <= ChunkSize {
				// zero newly claimed memory
				extendedPtr := unsafe.Add(a.base, a.offset)
				clear(unsafe.Slice((*byte)(unsafe.Pointer(extendedPtr)), additionalBytes))

				a.offset = newOffset
				s = unsafe.Slice((*T)(slicePtr), length)
				return
			}
		}
	}

	// allocate new space and copy the existing elements over
	s = Alloc[T](a, length)
	copy(s, slice)
	return
}

func (a *Allocator) Reset() {
	a.curr = 0
	a.base = unsafe.Pointer(a.chunks[0])
	a.offset = 0

	// Release jumbo buffers to the Go garbage collector
	if len(a.jumbos) > 0 {
		a.jumbos = nil
	}
}

func allocRaw(a *Allocator, size, align uintptr) (ptr unsafe.Pointer) {
	alignedOffset := (a.offset + align - 1) &^ (align - 1)
	end := alignedOffset + size

	if end > ChunkSize || size > MaxSmallAlloc {
		return nil
	}

	ptr = unsafe.Add(a.base, alignedOffset)
	a.offset = end

	clear(unsafe.Slice((*byte)(ptr), size))
	return
}

//go:noinline
func allocRawSlow(a *Allocator, size, align uintptr) (ptr unsafe.Pointer) {
	if size > MaxSmallAlloc {
		// jumbo alloc
		buf := make([]byte, size+align)
		a.jumbos = append(a.jumbos, buf)

		ptr := unsafe.Pointer(unsafe.SliceData(buf))
		// calculate alignment offset
		mask := uintptr(align) - 1
		aligned := (uintptr(ptr) + mask) &^ mask
		return unsafe.Add(ptr, aligned-uintptr(ptr))

	}

	a.curr++
	if a.curr >= len(a.chunks) {
		a.chunks = append(a.chunks, new([ChunkSize]byte))
	}

	ptr = unsafe.Pointer(a.chunks[a.curr])

	a.offset = size
	a.base = ptr

	clear(unsafe.Slice((*byte)(ptr), size))
	return
}
