package arena

import (
	"unsafe"
)

const jumboFlag uint32 = 1 << 31

// New allocates memory and returns a 32-bit packed ID instead of a pointer.
// High 16 bits = Chunk Index. Low 16 bits = Byte Offset.
func Alloc[T any](a *ByteArena) uint32 {
	var zero T
	size := int(unsafe.Sizeof(zero))
	align := int(unsafe.Alignof(zero))

	if size == 0 {
		return 0 // ID 0 represents nil/invalid
	}

	alignedOffset := (a.offset + align - 1) &^ (align - 1)

	// Fast path bounds check
	if alignedOffset+size > ChunkSize || size > MaxSmallAlloc {
		_ = allocRawSlow(a, size, align)

		if a.curr >= 32768 {
			panic("ByteArena out of memory: Exceeded maximum chunks (2.14 GB limit)")
		}

		if size > MaxSmallAlloc {
			jumboIdx := uint32(len(a.jumbos) - 1)

			if jumboIdx >= jumboFlag {
				panic("ByteArena out of memory: Exceeded maximum jumbo allocations")
			}

			return jumboFlag | jumboIdx
		}

		// allocRawSlow created a new chunk. Grab the updated state.
		chunkIdx := a.curr
		offsetBeforeSize := a.offset - size
		return uint32(chunkIdx<<16) | uint32(offsetBeforeSize)
	}

	chunkIdx := a.curr
	a.offset = alignedOffset + size

	// Ensure we never return ID 0 for a valid object.
	// If it's the very first object in the very first chunk, burn 8 bytes.
	if chunkIdx == 0 && alignedOffset == 0 {
		alignedOffset = 8
		a.offset = alignedOffset + size
	}

	return uint32(chunkIdx<<16) | uint32(alignedOffset)
}

func AllocSlice[T any](a *ByteArena, length int) uint32 {
	if length <= 0 {
		return 0
	}

	var zero T
	elemSize := int(unsafe.Sizeof(zero))
	align := int(unsafe.Alignof(zero))

	if elemSize == 0 {
		return 0
	}

	totalSize := elemSize * length
	alignedOffset := (a.offset + align - 1) &^ (align - 1)

	if alignedOffset+totalSize > ChunkSize || totalSize > MaxSmallAlloc {
		_ = allocRawSlow(a, totalSize, align)

		if a.curr >= 32768 {
			panic("ByteArena out of memory: Exceeded maximum chunks (2.14 GB limit)")
		}

		if totalSize > MaxSmallAlloc {
			jumboIdx := uint32(len(a.jumbos) - 1)

			if jumboIdx >= jumboFlag {
				panic("ByteArena out of memory: Exceeded maximum jumbo allocations")
			}

			return jumboFlag | jumboIdx
		}

		chunkIdx := a.curr
		offsetBeforeSize := a.offset - totalSize
		return uint32(chunkIdx<<16) | uint32(offsetBeforeSize)
	}

	chunkIdx := a.curr
	a.offset = alignedOffset + totalSize

	if chunkIdx == 0 && alignedOffset == 0 {
		alignedOffset = 8
		a.offset = alignedOffset + totalSize
	}

	return uint32(chunkIdx<<16) | uint32(alignedOffset)
}

func Get[T any](a *ByteArena, id uint32) *T {
	if id == 0 {
		return nil
	}

	chunkIdx := id >> 16
	offset := id & 0xFFFF

	ptr := unsafe.Add(unsafe.Pointer(a.chunks[chunkIdx]), offset)
	return (*T)(ptr)
}

func GetSlice[T any](a *ByteArena, id uint32, length int) []T {
	if id == 0 || length <= 0 {
		return nil
	}

	chunkIdx := id >> 16
	offset := id & 0xFFFF

	ptr := unsafe.Add(unsafe.Pointer(a.chunks[chunkIdx]), offset)
	return unsafe.Slice((*T)(ptr), length)
}
