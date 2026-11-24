package interner

import (
	"math/bits"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	initialTableSize = 512 // Must be a power of 2
	slab0Size        = 512 // First slab inline on Interner
)

type tableSnap struct {
	entries []atomic.Uint64
	mask    uint32
}

type Interner struct {
	snap atomic.Pointer[tableSnap]

	slab0 [slab0Size]string
	slabs [24]atomic.Pointer[[]string]

	mu      sync.Mutex
	curSnap *tableSnap
	nextId  uint32
}

func NewInterner() *Interner {
	return newInternerWithCap(initialTableSize)
}

func newInternerWithCap(capSize uint32) *Interner {
	ts := &tableSnap{
		entries: make([]atomic.Uint64, capSize),
		mask:    capSize - 1,
	}
	i := &Interner{
		curSnap: ts,
		nextId:  1, // burn 0 index
	}
	i.snap.Store(ts)
	return i
}

// Clone creates an independent copy of a base interner (e.g. pre-populated with stdlib symbols).
// Caller must ensure the source interner is not concurrently mutated.
func (src *Interner) Clone() *Interner {
	srcSnap := src.curSnap
	newEntries := make([]atomic.Uint64, len(srcSnap.entries))
	for j := range srcSnap.entries {
		if v := srcSnap.entries[j].Load(); v != 0 {
			newEntries[j].Store(v)
		}
	}
	ts := &tableSnap{
		entries: newEntries,
		mask:    srcSnap.mask,
	}
	dst := &Interner{
		slab0:   src.slab0,
		curSnap: ts,
		nextId:  src.nextId,
	}
	for k := range src.slabs {
		if p := src.slabs[k].Load(); p != nil {
			cp := make([]string, len(*p))
			copy(cp, *p)
			dst.slabs[k].Store(&cp)
		}
	}
	dst.snap.Store(ts)
	return dst
}

// Intern interns a string (cloning on miss so arena-backed strings remain valid after Reset).
func (i *Interner) Intern(s string) uint32 {
	h := uint32(HashString(s))
	snap := i.snap.Load()
	idx := h & snap.mask

	for {
		v := snap.entries[idx].Load()
		id := uint32(v)
		if id == 0 {
			return i.internSlow(s, h, true)
		}
		if uint32(v>>32) == h && i.Get(id) == s {
			return id
		}
		idx = (idx + 1) & snap.mask
	}
}

// InternStatic interns an immutable heap/literal string without cloning on miss.
func (i *Interner) InternStatic(s string) uint32 {
	h := uint32(HashString(s))
	snap := i.snap.Load()
	idx := h & snap.mask

	for {
		v := snap.entries[idx].Load()
		id := uint32(v)
		if id == 0 {
			return i.internSlow(s, h, false)
		}
		if uint32(v>>32) == h && i.Get(id) == s {
			return id
		}
		idx = (idx + 1) & snap.mask
	}
}

//go:noinline
func (i *Interner) internSlow(s string, h uint32, clone bool) uint32 {
	i.mu.Lock()
	defer i.mu.Unlock()

	snap := i.curSnap
	idx := h & snap.mask
	for {
		v := snap.entries[idx].Load()
		id := uint32(v)
		if id == 0 {
			break
		}
		if uint32(v>>32) == h && i.Get(id) == s {
			return id
		}
		idx = (idx + 1) & snap.mask
	}

	// Resize before inserting if load factor would exceed 75%.
	if i.nextId > (snap.mask*3)>>2 {
		snap = i.resizeLocked()
		idx = h & snap.mask
		for snap.entries[idx].Load() != 0 {
			idx = (idx + 1) & snap.mask
		}
	}

	stored := s
	if clone {
		stored = strings.Clone(s)
	}

	newId := i.nextId
	i.setStringLocked(newId, stored)
	i.nextId++

	// Release-store publishes both the string in slab and the hash entry.
	snap.entries[idx].Store((uint64(h) << 32) | uint64(newId))
	return newId
}

func (i *Interner) setStringLocked(id uint32, s string) {
	if id < slab0Size {
		i.slab0[id] = s
		return
	}
	b := bits.Len32(id) - 1
	slabIdx := b - 9
	p := i.slabs[slabIdx].Load()
	if p == nil {
		sl := make([]string, 1<<b)
		p = &sl
		i.slabs[slabIdx].Store(p)
	}
	(*p)[id^(uint32(1)<<b)] = s
}

func (i *Interner) Get(id uint32) string {
	if id < slab0Size {
		return i.slab0[id]
	}
	b := bits.Len32(id) - 1
	p := i.slabs[b-9].Load()
	if p == nil {
		return ""
	}
	return (*p)[id^(uint32(1)<<b)]
}

func (i *Interner) Reset() {
	i.mu.Lock()
	defer i.mu.Unlock()
	ts := &tableSnap{
		entries: make([]atomic.Uint64, initialTableSize),
		mask:    initialTableSize - 1,
	}
	i.slab0 = [slab0Size]string{}
	for k := range i.slabs {
		i.slabs[k].Store(nil)
	}
	i.curSnap = ts
	i.nextId = 1
	i.snap.Store(ts)
}

func (i *Interner) resizeLocked() *tableSnap {
	oldSnap := i.curSnap
	newMask := (oldSnap.mask << 1) | 1
	newEntries := make([]atomic.Uint64, newMask+1)

	for j := range oldSnap.entries {
		v := oldSnap.entries[j].Load()
		if uint32(v) == 0 {
			continue
		}
		h := uint32(v >> 32)
		idx := h & newMask
		for newEntries[idx].Load() != 0 {
			idx = (idx + 1) & newMask
		}
		newEntries[idx].Store(v)
	}

	newSnap := &tableSnap{
		entries: newEntries,
		mask:    newMask,
	}
	i.curSnap = newSnap
	i.snap.Store(newSnap)
	return newSnap
}
