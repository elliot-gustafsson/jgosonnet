package interner

import (
	"strings"
	"sync"
	"sync/atomic"
)

const (
	ChunkShift     = 12
	ChunkMask      = (1 << ChunkShift) - 1
	MaxFixedChunks = 256
	NumShards      = 32
	ShardBits      = 5
	InitialSlots   = 128
)

type tableData struct {
	slots []uint64
	mask  uint32
}

// shard is padded to exactly 64 bytes (8B mu + 8B tData + 4B count + 40B pad)
// to prevent false sharing and cache-line invalidation between CPU cores.
type shard struct {
	mu       sync.Mutex
	tData    atomic.Pointer[tableData]
	count    uint32
	_padding [44]byte
}

type Interner struct {
	shards   [NumShards]shard
	chunks   [MaxFixedChunks]atomic.Pointer[[1 << ChunkShift]string]
	chunksMu sync.Mutex
	overflow atomic.Pointer[[]*[1 << ChunkShift]string] // Used if > 256*4096 strings
	nextId   atomic.Uint32
}

func NewInterner() *Interner {
	in := &Interner{}
	for i := range NumShards {
		td := &tableData{
			slots: make([]uint64, InitialSlots),
			mask:  InitialSlots - 1,
		}
		in.shards[i].tData.Store(td)
	}

	firstChunk := new([1 << ChunkShift]string)
	in.chunks[0].Store(firstChunk)
	in.nextId.Store(1)
	return in
}

func (in *Interner) Intern(s string) uint32 {
	h64 := HashString(s)
	shardIdx := (h64 >> (64 - ShardBits))
	sh := &in.shards[shardIdx]
	h := uint32(h64)

	// try lock free
	td := sh.tData.Load()
	mask := td.mask
	idx := h & mask

	for {
		e := atomic.LoadUint64(&td.slots[idx])
		if e == 0 {
			break // fall back to slow path
		}
		if uint32(e>>32) == h {
			eIdx := uint32(e)
			cIdx := eIdx >> ChunkShift
			cOff := eIdx & ChunkMask
			chunk := in.loadChunk(cIdx)
			if chunk[cOff] == s {
				return eIdx
			}
		}
		idx = (idx + 1) & mask
	}

	// mutex protected slow path
	sh.mu.Lock()

	// reload table pointer in case of a resize while waiting on the lock
	td = sh.tData.Load()
	mask = td.mask
	idx = h & mask

	for {
		e := td.slots[idx]
		eIdx := uint32(e)

		if eIdx == 0 {
			newId := in.nextId.Add(1) - 1
			cIdx := newId >> ChunkShift
			cOff := newId & ChunkMask

			chunk := in.loadChunk(cIdx)
			if chunk == nil {
				in.chunksMu.Lock()
				chunk = in.ensureChunk(cIdx)
				in.chunksMu.Unlock()
			}

			// clone only on insertion to detach from large source buffers
			chunk[cOff] = strings.Clone(s)

			atomic.StoreUint64(&td.slots[idx], (uint64(h)<<32)|uint64(newId))
			sh.count++

			if sh.count > mask-mask/4 {
				sh.resize(td)
			}
			sh.mu.Unlock()
			return newId
		}

		if uint32(e>>32) == h {
			cIdx := eIdx >> ChunkShift
			cOff := eIdx & ChunkMask
			chunk := in.loadChunk(cIdx)
			if chunk[cOff] == s {
				sh.mu.Unlock()
				return eIdx
			}
		}
		idx = (idx + 1) & mask
	}
}

func (in *Interner) Get(id uint32) string {
	c := in.loadChunk(id >> ChunkShift)
	if c == nil {
		return ""
	}
	return c[id&ChunkMask]
}

func (in *Interner) Len() int {
	return int(in.nextId.Load()) - 1
}

func (in *Interner) loadChunk(cIdx uint32) *[1 << ChunkShift]string {
	if cIdx < MaxFixedChunks {
		return in.chunks[cIdx].Load()
	}
	ov := in.overflow.Load()
	if ov != nil && int(cIdx-MaxFixedChunks) < len(*ov) {
		return (*ov)[cIdx-MaxFixedChunks]
	}
	return nil
}

func (in *Interner) ensureChunk(cIdx uint32) *[1 << ChunkShift]string {
	if cIdx < MaxFixedChunks {
		chunk := in.chunks[cIdx].Load()
		if chunk == nil {
			chunk = new([1 << ChunkShift]string)
			in.chunks[cIdx].Store(chunk)
		}
		return chunk
	}

	ovIdx := int(cIdx - MaxFixedChunks)
	ovPtr := in.overflow.Load()
	var oldOv []*[1 << ChunkShift]string
	if ovPtr != nil {
		oldOv = *ovPtr
	}

	if ovIdx < len(oldOv) && oldOv[ovIdx] != nil {
		return oldOv[ovIdx]
	}

	newLen := len(oldOv)
	if ovIdx >= newLen {
		newLen = max(ovIdx+1, len(oldOv)*2)
	}

	newOv := make([]*[1 << ChunkShift]string, newLen)
	copy(newOv, oldOv)

	chunk := newOv[ovIdx]
	if chunk == nil {
		chunk = new([1 << ChunkShift]string)
		newOv[ovIdx] = chunk
	}

	in.overflow.Store(&newOv)
	return chunk
}

func (sh *shard) resize(oldTd *tableData) {
	newMask := (oldTd.mask << 1) | 1
	newSlots := make([]uint64, newMask+1)

	for _, e := range oldTd.slots {
		if uint32(e) == 0 {
			continue
		}
		h := uint32(e >> 32)
		idx := h & newMask
		for {
			if newSlots[idx] == 0 {
				newSlots[idx] = e
				break
			}
			idx = (idx + 1) & newMask
		}
	}
	sh.tData.Store(&tableData{
		slots: newSlots,
		mask:  newMask,
	})
}
