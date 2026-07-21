package interner

import (
	"strings"
	"sync"
	"sync/atomic"
)

const (
	defaultInternerSize = 16384 // Must be a power of 2
	defaultInternerMask = defaultInternerSize - 1
	chunkSize           = 1024
	chunkShift          = 10
	chunkMask           = chunkSize - 1
)

type entry struct {
	hash uint32
	idx  uint32
}

type Interner struct {
	mu     sync.Mutex
	table  []entry
	chunks atomic.Pointer[[]*[chunkSize]string]
	mask   uint32
	count  uint32
}

func NewInterner() *Interner {
	i := &Interner{
		mu:    sync.Mutex{},
		table: make([]entry, defaultInternerSize),
		mask:  defaultInternerMask,
	}

	firstChunk := new([chunkSize]string)
	firstChunk[0] = "" // burn 0 index

	initialChunks := make([]*[chunkSize]string, 1)
	initialChunks[0] = firstChunk

	i.chunks.Store(&initialChunks)
	i.count = 1

	return i
}

func (i *Interner) Intern(s string) uint32 {
	i.mu.Lock()
	defer i.mu.Unlock()

	h := fnv1a(s)
	idx := h & i.mask

	for {
		e := i.table[idx]

		if e.idx == 0 {
			newId := i.count
			chunkIdx := newId >> chunkShift
			itemIdx := newId & chunkMask

			oldChunks := *i.chunks.Load()
			if chunkIdx >= uint32(len(oldChunks)) {
				// Allocate a new slice of pointers. Since this only happens once
				// every 1024 strings, copying this tiny array is virtually free.
				newChunks := make([]*[chunkSize]string, len(oldChunks)+1)
				copy(newChunks, oldChunks)

				newChunk := new([chunkSize]string)
				newChunk[itemIdx] = strings.Clone(s)
				newChunks[len(oldChunks)] = newChunk

				i.chunks.Store(&newChunks) // Atomically swap in the new slice header
			} else {
				// It is safe to mutate the array element because readers only read old, fully
				// initialized strings, and we bounds check on count before reading this index.
				oldChunks[chunkIdx][itemIdx] = strings.Clone(s)
			}

			i.table[idx] = entry{h, newId}
			atomic.StoreUint32(&i.count, newId+1)

			// resize if table becomes more than 50% full
			if i.count > i.mask/2 {
				i.resize()
			}
			return newId
		}

		if e.hash == h {
			chunks := *i.chunks.Load()
			if chunks[e.idx>>chunkShift][e.idx&chunkMask] == s {
				return e.idx
			}
		}

		idx = (idx + 1) & i.mask
	}
}

func (i *Interner) Get(id uint32) string {
	// Bounds check ensures readers never read ahead of what writers have committed
	if id == 0 || id >= atomic.LoadUint32(&i.count) {
		return ""
	}

	chunks := *i.chunks.Load()
	return chunks[id>>chunkShift][id&chunkMask]
}

func (i *Interner) Reset() {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.table = make([]entry, defaultInternerSize)
	i.mask = defaultInternerMask

	firstChunk := new([chunkSize]string)
	firstChunk[0] = ""

	initialChunks := make([]*[chunkSize]string, 1)
	initialChunks[0] = firstChunk

	i.chunks.Store(&initialChunks)
	atomic.StoreUint32(&i.count, 1)
}

func (i *Interner) resize() {
	newMask := (i.mask << 1) | 1
	newTable := make([]entry, newMask+1)

	for j := range i.table {
		e := i.table[j]
		if e.idx == 0 {
			continue
		}

		idx := e.hash & newMask
		for {
			if newTable[idx].idx == 0 {
				newTable[idx] = e
				break
			}
			idx = (idx + 1) & newMask
		}
	}

	i.table = newTable
	i.mask = newMask
}

// fnv1a hash function
func fnv1a(s string) uint32 {
	var hash uint32 = 2166136261
	for j := 0; j < len(s); j++ {
		hash ^= uint32(s[j])
		hash *= 16777619
	}
	return hash
}
