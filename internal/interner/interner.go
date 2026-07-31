package interner

import (
	"hash/maphash"
)

const (
	defaultInternerSize = 16384 // Must be a power of 2
	defaultInternerMask = defaultInternerSize - 1
)

type entry struct {
	hash uint32
	idx  uint32
}

type Interner struct {
	table   []entry
	strings []string
	mask    uint32
	count   uint32
	seed    maphash.Seed
}

func NewInterner() *Interner {
	interner := &Interner{
		table:   make([]entry, defaultInternerSize),
		strings: make([]string, 1, defaultInternerSize/2), // burn 0 index
		mask:    defaultInternerMask,
		seed:    maphash.MakeSeed(),
	}
	interner.strings[0] = ""
	return interner
}

func (i *Interner) Intern(s string) uint32 {

	h := uint32(maphash.String(i.seed, s))
	idx := h & i.mask

	for {
		e := i.table[idx]

		if e.idx == 0 {
			newId := uint32(len(i.strings))
			i.strings = append(i.strings, s)
			i.table[idx] = entry{h, newId}

			i.count++

			// resize if table becomes more than 50% full
			if i.count > i.mask/2 {
				i.resize()
			}
			return newId
		}

		if e.hash == h && i.strings[e.idx] == s {
			return e.idx
		}

		idx = (idx + 1) & i.mask
	}
}

func (i *Interner) Get(id uint32) string {
	if id == 0 || id >= uint32(len(i.strings)) {
		return ""
	}
	return i.strings[id]
}

func (i *Interner) Reset() {
	clear(i.table)
	i.strings = i.strings[:1]
	i.count = 0
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
