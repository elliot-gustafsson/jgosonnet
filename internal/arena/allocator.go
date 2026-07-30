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
	chunks []*[ChunkSize]byte
	curr   int
	offset int
	jumbos [][]byte
}

func NewAllocator() (a *Allocator) {
	a = &Allocator{
		chunks: make([]*[ChunkSize]byte, 1, 64),
	}
	a.chunks[0] = new([ChunkSize]byte)
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
func Realloc[T any](a *Allocator, slice []T, length int) []T {
	// If the slice already has enough capacity, we can just reslice it.
	if length <= cap(slice) {
		return slice[:length]
	}

	var zero T
	elemSize := int(unsafe.Sizeof(zero))

	if elemSize == 0 {
		var empty struct{}
		return unsafe.Slice((*T)(unsafe.Pointer(&empty)), length)
	}

	if cap(slice) > 0 {
		// calculate the memory address at the end of the slice's capacity
		slicePtr := unsafe.Pointer(unsafe.SliceData(slice))
		sliceEnd := uintptr(slicePtr) + uintptr(cap(slice)*elemSize)

		// calculate the memory address of the arena's current offset
		arenaEnd := uintptr(unsafe.Pointer(a.chunks[a.curr])) + uintptr(a.offset)

		// if they match, this slice was the very last allocation in the current chunk.
		if sliceEnd == arenaEnd {
			additionalBytes := (length - cap(slice)) * elemSize
			totalSize := length * elemSize

			// ensure it still fits in the chunk and respects the max small alloc threshold
			if totalSize <= MaxSmallAlloc && (a.offset)+additionalBytes <= ChunkSize {
				// zero newly claimed memory
				extendedPtr := unsafe.Add(unsafe.Pointer(a.chunks[a.curr]), a.offset)
				clear(unsafe.Slice((*byte)(unsafe.Pointer(extendedPtr)), additionalBytes))

				a.offset += additionalBytes
				return unsafe.Slice((*T)(slicePtr), length)
			}
		}
	}

	// Fallback: allocate new space and copy the existing elements over
	newSlice := Alloc[T](a, length)
	copy(newSlice, slice)
	return newSlice
}

func (a *Allocator) Reset() {
	a.curr = 0
	a.offset = 0

	// Release jumbo buffers to the Go garbage collector
	if len(a.jumbos) > 0 {
		a.jumbos = nil
	}
}

func allocRaw(a *Allocator, size, align uintptr) (ptr unsafe.Pointer) {
	alignedOffset := (uintptr(a.offset) + align - 1) &^ (align - 1)

	if alignedOffset+size > ChunkSize || size > MaxSmallAlloc {
		return nil
	}

	ptr = unsafe.Add(unsafe.Pointer(a.chunks[a.curr]), alignedOffset)
	a.offset = int(alignedOffset + size)

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
		offset := aligned - uintptr(ptr)
		return unsafe.Add(ptr, offset)

	}

	alignedOffset := (uintptr(a.offset) + align - 1) &^ (align - 1)
	if alignedOffset+size > ChunkSize {
		// not enough space in current block, add a new one
		a.curr++
		if a.curr >= len(a.chunks) {
			a.chunks = append(a.chunks, new([ChunkSize]byte))
		}
		a.offset = 0
		alignedOffset = 0
	}

	ptr = unsafe.Add(unsafe.Pointer(a.chunks[a.curr]), alignedOffset)
	a.offset = int(alignedOffset + size)

	clear(unsafe.Slice((*byte)(ptr), size))
	return
}
