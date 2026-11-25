package interner

import (
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	shardBits = 4
	numShards = 1 << shardBits // 16
	shardMask = numShards - 1

	// Each shard packs length-prefixed strings directly into pointer-free byte slabs.
	// Layout of localAddr (28 bits):
	//   bit 27 == 0: regular slab string -> slabIdx = localAddr >> slabShift, off = localAddr & slabMask
	//   bit 27 == 1: jumbo string (> maxSlabStringLen) -> jumboIdx = localAddr &^ jumboFlag
	slabShift        = 12 // 4 KB slabs (noscan for Go GC)
	slabSize         = 1 << slabShift
	slabMask         = slabSize - 1
	maxSlabStringLen = slabSize - 4 // leaves >=2B so &slab[off+2] is always strictly in-bounds
	jumboFlag        = uint32(1) << 27

	initialSlots = 128
)

type tableData struct {
	slotsPtr unsafe.Pointer // unsafe.SliceData(slots)
	slab0    unsafe.Pointer // &[slabSize]byte
	mask     uint32
	slots    []uint64
}

// shard uses two 64-byte cache lines (128 bytes total):
//   - Line 0: lock-free read-hot fields, only written on resize or new slab allocation.
//   - Line 1: writer-only fields, isolated so locks/inserts never dirty reader cache lines.
type shard struct {
	tData      atomic.Pointer[tableData]
	slab0      atomic.Pointer[[slabSize]byte]
	extraSlabs atomic.Pointer[[]*[slabSize]byte]
	jumbos     atomic.Pointer[[]string]
	_pad0      [32]byte

	mu      sync.Mutex
	curSlab *[slabSize]byte
	count   atomic.Uint32 // written under mu; atomic load for Len()
	slabIdx uint32        // current active slab index (0 = slab0, >=1 = extraSlabs[slabIdx-1])
	slabOff uint32        // next 2-byte aligned offset in current active slab
	_pad1   [36]byte
}

type Interner struct {
	shards [numShards]shard
}

func NewInterner() *Interner {
	return &Interner{}
}

func (in *Interner) Intern(s string) uint32 {
	sp := unsafe.StringData(s)
	slen := len(s)

	var h64 uint64
	seed := rapidSeedInit ^ uint64(slen)
	if slen == 0 {
		h64 = rapidMix(seed, rapidSecret2)
	} else if slen <= 16 {
		var a, b uint64
		if slen >= 8 {
			a = read64(sp, 0)
			b = read64(sp, slen-8)
		} else if slen >= 4 {
			a = uint64(read32(sp, 0))
			b = uint64(read32(sp, slen-4))
		} else {
			a = (uint64(read8(sp, 0)) << 16) | (uint64(read8(sp, slen>>1)) << 8) | uint64(read8(sp, slen-1))
		}
		h64 = rapidMix(seed^a, rapidSecret2^b)
	} else if slen <= 32 {
		seed = rapidMix(read64(sp, 0)^rapidSecret2, read64(sp, 8)^seed^rapidSecret1)
		h64 = rapidMix(seed^read64(sp, slen-16), rapidSecret2^read64(sp, slen-8))
	} else {
		h64 = hashStringOver32(sp, slen, seed)
	}

	shardIdx := uint32(h64 >> (64 - shardBits))
	sh := &in.shards[shardIdx]
	h := uint32(h64)

	if td := sh.tData.Load(); td != nil {
		mask := td.mask
		slotsPtr := td.slotsPtr
		slab0 := td.slab0
		idx := h & mask

		for {
			e := atomic.LoadUint64((*uint64)(unsafe.Add(slotsPtr, uintptr(idx)*8)))
			if e == 0 {
				break
			}
			if uint32(e>>32) == h {
				localAddr := uint32(e)
				if localAddr < slabSize && slen <= 32 {
					p := unsafe.Add(slab0, uintptr(localAddr))
					if int(*(*uint16)(p)) == slen {
						b := (*byte)(unsafe.Add(p, 2))
						if slen >= 17 {
							if ((read64(b, 0) ^ read64(sp, 0)) |
								(read64(b, 8) ^ read64(sp, 8)) |
								(read64(b, slen-16) ^ read64(sp, slen-16)) |
								(read64(b, slen-8) ^ read64(sp, slen-8))) == 0 {
								return (localAddr << shardBits) | shardIdx
							}
						} else if slen >= 8 {
							if ((read64(b, 0) ^ read64(sp, 0)) |
								(read64(b, slen-8) ^ read64(sp, slen-8))) == 0 {
								return (localAddr << shardBits) | shardIdx
							}
						} else if slen >= 4 {
							if ((read32(b, 0) ^ read32(sp, 0)) |
								(read32(b, slen-4) ^ read32(sp, slen-4))) == 0 {
								return (localAddr << shardBits) | shardIdx
							}
						} else if slen == 0 ||
							(read8(b, 0) == read8(sp, 0) &&
								read8(b, slen>>1) == read8(sp, slen>>1) &&
								read8(b, slen-1) == read8(sp, slen-1)) {
							return (localAddr << shardBits) | shardIdx
						}
					}
				} else if sh.matchSlow(td.slab0, localAddr, s) {
					return (localAddr << shardBits) | shardIdx
				}
			}
			idx = (idx + 1) & mask
		}
	}

	return sh.internSlow(s, h, shardIdx)
}

func (in *Interner) Get(id uint32) string {
	sh := &in.shards[id&shardMask]
	localAddr := id >> shardBits
	if b := sh.slab0.Load(); b != nil && localAddr <= maxSlabStringLen {
		p := unsafe.Add(unsafe.Pointer(b), uintptr(localAddr))
		return unsafe.String((*byte)(unsafe.Add(p, 2)), int(*(*uint16)(p)))
	}
	return sh.getSlow(localAddr)
}

func (in *Interner) Len() int {
	total := 0
	for i := range numShards {
		total += int(in.shards[i].count.Load())
	}
	return total
}

//go:noinline
func (sh *shard) matchSlow(slab0 unsafe.Pointer, localAddr uint32, s string) bool {
	if localAddr < slabSize {
		p := unsafe.Add(slab0, uintptr(localAddr))
		return unsafe.String((*byte)(unsafe.Add(p, 2)), int(*(*uint16)(p))) == s
	}
	return sh.getSlow(localAddr) == s
}

//go:noinline
func (sh *shard) internSlow(s string, h, shardIdx uint32) uint32 {
	sh.mu.Lock()

	td := sh.tData.Load()
	if td == nil {
		slab0 := new([slabSize]byte)
		sh.curSlab = slab0
		sh.slabOff = 2 // reserve offset 0 so localAddr is always > 0
		sh.slab0.Store(slab0)

		localAddr := sh.appendString(s)

		slots := make([]uint64, initialSlots)
		mask := uint32(initialSlots - 1)
		slots[h&mask] = (uint64(h) << 32) | uint64(localAddr)

		sh.count.Store(1)
		sh.tData.Store(&tableData{
			slotsPtr: unsafe.Pointer(unsafe.SliceData(slots)),
			slab0:    unsafe.Pointer(slab0),
			mask:     mask,
			slots:    slots,
		})
		sh.mu.Unlock()
		return (localAddr << shardBits) | shardIdx
	}

	mask := td.mask
	slots := td.slots
	idx := h & mask

	for {
		e := slots[idx]
		if e == 0 {
			localAddr := sh.appendString(s)
			atomic.StoreUint64(&slots[idx], (uint64(h)<<32)|uint64(localAddr))
			newCount := sh.count.Load() + 1
			sh.count.Store(newCount)

			if newCount > mask-mask/4 {
				sh.resize(td)
			}
			sh.mu.Unlock()
			return (localAddr << shardBits) | shardIdx
		}

		if uint32(e>>32) == h {
			localAddr := uint32(e)
			if sh.matchSlow(td.slab0, localAddr, s) {
				sh.mu.Unlock()
				return (localAddr << shardBits) | shardIdx
			}
		}
		idx = (idx + 1) & mask
	}
}

func (sh *shard) appendString(s string) uint32 {
	n := len(s)
	off := int(sh.slabOff)
	if n <= maxSlabStringLen && off+2+n <= slabSize-2 {
		slab := sh.curSlab
		*(*uint16)(unsafe.Pointer(&slab[off])) = uint16(n)
		if n > 0 {
			copy(slab[off+2:], s)
		}
		sh.slabOff = uint32((off + 2 + n + 1) &^ 1)
		return (sh.slabIdx << slabShift) | uint32(off)
	}
	return sh.appendStringSlow(s)
}

//go:noinline
func (sh *shard) appendStringSlow(s string) uint32 {
	n := len(s)
	if n > maxSlabStringLen {
		cloned := strings.Clone(s)
		oldPtr := sh.jumbos.Load()
		var oldJumbos []string
		if oldPtr != nil {
			oldJumbos = *oldPtr
		}
		jIdx := uint32(len(oldJumbos))
		newJumbos := make([]string, jIdx+1)
		copy(newJumbos, oldJumbos)
		newJumbos[jIdx] = cloned
		sh.jumbos.Store(&newJumbos)
		return jumboFlag | jIdx
	}

	newSlab := new([slabSize]byte)
	*(*uint16)(unsafe.Pointer(&newSlab[0])) = uint16(n)
	if n > 0 {
		copy(newSlab[2:], s)
	}

	sh.slabIdx++
	sh.curSlab = newSlab
	sh.slabOff = uint32((2 + n + 1) &^ 1)

	oldPtr := sh.extraSlabs.Load()
	var oldSlabs []*[slabSize]byte
	if oldPtr != nil {
		oldSlabs = *oldPtr
	}
	newSlabs := make([]*[slabSize]byte, sh.slabIdx)
	copy(newSlabs, oldSlabs)
	newSlabs[sh.slabIdx-1] = newSlab
	sh.extraSlabs.Store(&newSlabs)

	return sh.slabIdx << slabShift
}

//go:noinline
func (sh *shard) getSlow(localAddr uint32) string {
	if localAddr&jumboFlag != 0 {
		jPtr := sh.jumbos.Load()
		if jPtr == nil {
			return ""
		}
		jIdx := int(localAddr &^ jumboFlag)
		if jIdx >= len(*jPtr) {
			return ""
		}
		return (*jPtr)[jIdx]
	}

	if localAddr < slabSize {
		return ""
	}
	exPtr := sh.extraSlabs.Load()
	if exPtr == nil {
		return ""
	}
	sIdx := int(localAddr>>slabShift) - 1
	if sIdx >= len(*exPtr) {
		return ""
	}
	off := localAddr & slabMask
	b := (*exPtr)[sIdx]
	p := unsafe.Add(unsafe.Pointer(b), uintptr(off))
	n := int(*(*uint16)(p))
	if n == 0 {
		return ""
	}
	return unsafe.String((*byte)(unsafe.Add(p, 2)), n)
}

//go:noinline
func (sh *shard) resize(oldTd *tableData) {
	newMask := (oldTd.mask << 1) | 1
	newSlots := make([]uint64, newMask+1)

	for _, e := range oldTd.slots {
		if e == 0 {
			continue
		}
		h := uint32(e >> 32)
		idx := h & newMask
		for newSlots[idx] != 0 {
			idx = (idx + 1) & newMask
		}
		newSlots[idx] = e
	}
	sh.tData.Store(&tableData{
		slotsPtr: unsafe.Pointer(unsafe.SliceData(newSlots)),
		slab0:    oldTd.slab0,
		mask:     newMask,
		slots:    newSlots,
	})
}
