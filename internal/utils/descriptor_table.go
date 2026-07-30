package utils

import (
	"math"

	"github.com/elliot-gustafsson/jgosonnet/internal/arena"
)

type Descriptor struct {
	Key   uint32
	Value uint32
}

type DescriptorTable struct {
	entries []Descriptor
	capMask uint32
	count   uint32
}

func NewDescriptorTable(a *arena.Allocator, symbols []uint32) *DescriptorTable {
	n := uint32(len(symbols))
	// 1.4x capacity guarantees a max load factor of ~70% for minimal probing
	cap := nextPowerOf2(n + (n >> 1))
	capMask := cap - 1

	entries := arena.Alloc[Descriptor](a, int(cap))

	for i, sym := range symbols {
		idx := hashUint32(sym) & capMask
		for entries[idx].Key != 0 {
			idx = (idx + 1) & capMask
		}
		entries[idx].Key = sym
		entries[idx].Value = uint32(i)
	}

	dt := arena.Create[DescriptorTable](a)
	dt.entries = entries
	dt.capMask = capMask
	dt.count = n
	return dt
}

func NewEmptyDescriptorTable(a *arena.Allocator, l int) *DescriptorTable {
	n := uint32(l)
	// 1.4x capacity guarantees a max load factor of ~70% for minimal probing
	cap := nextPowerOf2(n + (n >> 1))
	capMask := cap - 1

	dt := arena.Create[DescriptorTable](a)
	dt.entries = arena.Alloc[Descriptor](a, int(cap))
	dt.capMask = capMask
	dt.count = 0
	return dt
}

func (dt *DescriptorTable) Length() int {
	return int(dt.count)
}

func (dt *DescriptorTable) Append(key uint32) uint32 {
	idx := hashUint32(key) & dt.capMask
	for {
		entry := &dt.entries[idx]
		if entry.Key == 0 {
			entry.Key = key
			assignedIdx := dt.count
			entry.Value = assignedIdx
			dt.count++
			return assignedIdx
		}
		if entry.Key == key {
			// already exist, return max value as way of signaling err
			return math.MaxUint32
		}
		idx = (idx + 1) & dt.capMask
	}
}

func (dt *DescriptorTable) Get(key uint32) (uint32, bool) {
	idx := hashUint32(key) & dt.capMask
	for {
		entry := &dt.entries[idx]
		if entry.Key == key {
			return entry.Value, true
		}
		if entry.Key == 0 {
			return 0, false
		}
		idx = (idx + 1) & dt.capMask
	}
}

// hashUint32 is a multiplicative Fibonacci hash
func hashUint32(x uint32) uint32 {
	return x * 2654435769
}

// nextPowerOf2 rounds up to the nearest power of 2.
func nextPowerOf2(n uint32) uint32 {
	if n < 2 {
		return 2
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}
