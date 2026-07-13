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
		return a.headers.Alloc(s)
	}

	if a.offset+length > a.blockSize {
		a.grow()
	}

	targetBytes := a.elementBlocks[a.activeIdx][a.offset : a.offset+length]
	a.offset += length

	copy(targetBytes, s)

	res := unsafe.String(unsafe.SliceData(targetBytes), length)
	return a.headers.Alloc(res)
}

func (a *StringArena) AllocConcat(s1, s2 string) uint32 {
	length := len(s1) + len(s2)

	if length == 0 {
		return a.headers.Alloc("")
	}

	if length > a.blockSize {
		return a.allocJumbo(s1, s2, length)
	}

	if a.offset+length > a.blockSize {
		a.grow()
	}

	targetBytes := a.elementBlocks[a.activeIdx][a.offset : a.offset+length]
	a.offset += length

	n := copy(targetBytes, s1)
	copy(targetBytes[n:], s2)

	res := unsafe.String(unsafe.SliceData(targetBytes), length)
	return a.headers.Alloc(res)
}

func (a *StringArena) Get(id uint32) string {
	return a.headers.GetValue(id)
}

func (a *StringArena) Reset() {
	a.activeIdx = 0
	a.offset = 0

	a.headers.Reset()
}

//go:noinline
func (a *StringArena) grow() {
	a.activeIdx++
	a.offset = 0
	if a.activeIdx < len(a.elementBlocks) {
		return
	}
	a.elementBlocks = append(a.elementBlocks, make([]byte, a.blockSize))
}

//go:noinline
func (a *StringArena) allocJumbo(s1, s2 string, length int) uint32 {
	jumboBlock := make([]byte, length)
	n := copy(jumboBlock, s1)
	copy(jumboBlock[n:], s2)

	res := unsafe.String(unsafe.SliceData(jumboBlock), length)
	return a.headers.Alloc(res)
}
