package utils

import "github.com/elliot-gustafsson/jgosonnet/internal/arena"

type Property[T any] struct {
	Key   uint32
	Meta  uint8
	Value T
}

// PropertyMap is a dynamic, arena-backed open-addressed hash table.
type PropertyMap[T any] struct {
	entries []Property[T]
	capMask uint32
	count   uint32
}

// NewPropertyMap initializes a PropertyMap with arena-allocated storage.
func NewPropertyMap[T any](a *arena.Allocator, initialCap uint32) *PropertyMap[T] {
	cap := nextPowerOf2(initialCap)
	if cap == 0 {
		cap = 8
	}

	pm := arena.Create[PropertyMap[T]](a)
	arena.Memclr(pm)
	pm.entries = arena.Alloc[Property[T]](a, int(cap))
	arena.MemclrSlice(pm.entries)
	pm.capMask = cap - 1
	pm.count = 0
	return pm
}

func (pm *PropertyMap[T]) Get(sym uint32) (t T, v bool) {
	t, _, v = pm.GetEx(sym)
	return
}

func (pm *PropertyMap[T]) GetEx(sym uint32) (T, uint8, bool) {
	if len(pm.entries) == 0 {
		var empty T
		return empty, 0, false
	}
	idx := hashUint32(sym) & pm.capMask
	for {
		entry := &pm.entries[idx]
		if entry.Key == sym {
			return entry.Value, entry.Meta, true
		}
		if entry.Key == 0 {
			var empty T
			return empty, 0, false
		}
		idx = (idx + 1) & pm.capMask
	}
}

func (pm *PropertyMap[T]) Put(a *arena.Allocator, key uint32, val T) {
	pm.PutEx(a, key, val, 0)
}

// Put inserts or updates a property, growing the arena allocation if load factor >= 75%.
func (pm *PropertyMap[T]) PutEx(a *arena.Allocator, key uint32, val T, meta uint8) {
	if len(pm.entries) == 0 || (pm.count*4 >= (pm.capMask+1)*3) {
		pm.grow(a)
	}

	idx := hashUint32(key) & pm.capMask
	for {
		entry := &pm.entries[idx]
		if entry.Key == 0 {
			entry.Key = key
			entry.Meta = meta
			entry.Value = val
			pm.count++
			return
		}
		if entry.Key == key {
			entry.Value = val
			entry.Meta = meta
			return
		}
		idx = (idx + 1) & pm.capMask
	}
}

// grow doubles capacity by allocating a new slice from the arena and re-inserting active entries.
//
//go:noinline
func (pm *PropertyMap[T]) grow(a *arena.Allocator) {
	oldEntries := pm.entries
	newCap := uint32(len(oldEntries)) * 2
	if newCap == 0 {
		newCap = 8
	}
	newMask := newCap - 1

	// Fast allocation of the larger backing table
	newEntries := arena.Alloc[Property[T]](a, int(newCap))
	clear(newEntries)

	for i := range oldEntries {
		old := &oldEntries[i]
		if old.Key != 0 {
			idx := hashUint32(old.Key) & newMask
			for newEntries[idx].Key != 0 {
				idx = (idx + 1) & newMask
			}
			newEntries[idx] = *old
		}
	}

	pm.entries = newEntries
	pm.capMask = newMask
}
