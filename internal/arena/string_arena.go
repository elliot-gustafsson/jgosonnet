package arena

import "unsafe"

// StringArena manages raw contiguous memory for string bytes (avoiding GC overhead),
type StringArena struct {
	elementBlocks [][]byte
	blockSize     int
	activeIdx     int
	offset        int

	headers *Arena[string]
}

func NewStringArena(blockSize int) *StringArena {
	return &StringArena{
		elementBlocks: [][]byte{make([]byte, blockSize)},
		blockSize:     blockSize,
		activeIdx:     0,
		offset:        0,
		headers:       NewArena[string](),
	}
}

func (a *StringArena) Alloc(s string) uint32 {
	length := len(s)

	if length == 0 {
		return a.headers.Alloc("")
	}

	if length > a.blockSize {
		return a.headers.Alloc(s) // Avoid copying jumbo strings
	}

	targetBytes := a.reserve(length)
	copy(targetBytes, s)

	arenaStr := unsafe.String(&targetBytes[0], length)
	return a.headers.Alloc(arenaStr)
}

func (a *StringArena) AllocConcat(s1, s2 string) uint32 {
	length := len(s1) + len(s2)

	if length == 0 {
		return a.headers.Alloc("")
	}
	if length > a.blockSize {
		jumboBlock := make([]byte, length)
		copy(jumboBlock, s1)
		copy(jumboBlock[len(s1):], s2)

		arenaStr := unsafe.String(&jumboBlock[0], length)
		return a.headers.Alloc(arenaStr)
	}

	targetBytes := a.reserve(length)
	copy(targetBytes, s1)
	copy(targetBytes[len(s1):], s2)

	arenaStr := unsafe.String(&targetBytes[0], length)
	return a.headers.Alloc(arenaStr)
}

func (a *StringArena) Get(id uint32) string {
	return a.headers.GetValue(id)
}

func (a *StringArena) Reset() {
	a.activeIdx = 0
	a.offset = 0

	a.headers.Reset()
}

func (a *StringArena) reserve(length int) []byte {
	currBlock := a.elementBlocks[a.activeIdx]
	if a.offset+length > len(currBlock) {
		a.activeIdx++
		a.offset = 0
		if a.activeIdx < len(a.elementBlocks) {
			currBlock = a.elementBlocks[a.activeIdx]
		} else {
			currBlock = make([]byte, a.blockSize)
			a.elementBlocks = append(a.elementBlocks, currBlock)
		}
	}

	targetBytes := currBlock[a.offset : a.offset+length]
	a.offset += length
	return targetBytes
}
