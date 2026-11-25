package interner

import (
	"crypto/rand"
	"fmt"
	"strings"
	"sync"
	"testing"
	"unsafe"
)

var (
	benchKeysUnique   []string
	benchKeysRepeated []string
	benchKeysShort    []string
)

func init() {
	benchKeysUnique = make([]string, 10000)
	for i := range benchKeysUnique {
		b := make([]byte, 10)
		rand.Read(b)
		benchKeysUnique[i] = fmt.Sprintf("key_%x", b)
	}

	benchKeysRepeated = make([]string, 100000)
	for i := range benchKeysRepeated {
		benchKeysRepeated[i] = benchKeysUnique[i%100]
	}

	benchKeysShort = make([]string, 100000)
	shortPool := []string{
		"self", "std", "name", "metadata", "labels", "spec", "value", "key",
		"x", "y", "z", "items", "data", "config", "host", "port",
		"namespace", "annotations", "containers", "image", "env", "args",
	}
	for i := range benchKeysShort {
		benchKeysShort[i] = shortPool[i%len(shortPool)]
	}
}

func TestRapidSeedInit(t *testing.T) {
	want := rapidSeed ^ rapidMix(rapidSeed^rapidSecret0, rapidSecret1)
	if rapidSeedInit != want {
		t.Fatalf("rapidSeedInit = 0x%x, want 0x%x", rapidSeedInit, want)
	}
}

func TestShardSize(t *testing.T) {
	if got := unsafe.Sizeof(shard{}); got != 128 {
		t.Fatalf("sizeof(shard) = %d, want 128", got)
	}
}

func TestInternerBasic(t *testing.T) {
	in := NewInterner()

	if in.Get(0) != "" {
		t.Fatalf("Get(0) = %q, want empty", in.Get(0))
	}
	if in.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", in.Len())
	}

	// Empty string must get a non-zero ID
	emptyID := in.Intern("")
	if emptyID == 0 {
		t.Fatal("Intern(\"\") returned 0")
	}
	if in.Get(emptyID) != "" {
		t.Fatalf("Get(emptyID) = %q, want empty", in.Get(emptyID))
	}
	if in.Intern("") != emptyID {
		t.Fatal("second Intern(\"\") returned different ID")
	}

	// Large string (> maxArenaStringSize)
	large := strings.Repeat("abcdef0123456789", 200)
	largeID := in.Intern(large)
	if largeID == 0 || in.Get(largeID) != large || in.Intern(large) != largeID {
		t.Fatal("large string interning failed")
	}

	// Detachment from mutable source buffer
	buf := []byte("hello_detach")
	s := unsafe.String(unsafe.SliceData(buf), len(buf))
	id := in.Intern(s)
	copy(buf, "XXXXXXXXXXXX")
	if got := in.Get(id); got != "hello_detach" {
		t.Fatalf("interned string mutated with source buffer: got %q", got)
	}

	// Boundary lengths across all hash/comparison/slab branches (0..40 and around maxSlabStringLen)
	boundaryLens := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 23, 24, 31, 32, 33, 64, 96, 97, 200, maxSlabStringLen - 1, maxSlabStringLen, maxSlabStringLen + 1, maxSlabStringLen + 2}
	bIDs := make([]uint32, len(boundaryLens))
	for i, l := range boundaryLens {
		s := strings.Repeat("x", l)
		if l > 0 {
			s = "a" + s[1:]
		}
		if l > 1 {
			s = s[:l-1] + "z"
		}
		id := in.Intern(s)
		if id == 0 {
			t.Fatalf("Intern(len=%d) returned 0", l)
		}
		if got := in.Get(id); got != s {
			t.Fatalf("Get(len=%d) mismatch: len=%d", l, len(got))
		}
		if gotID := in.Intern(s); gotID != id {
			t.Fatalf("re-Intern(len=%d) = %d, want %d", l, gotID, id)
		}
		bIDs[i] = id
	}
	// Verify all distinct lengths got distinct IDs (except len 0 already interned as emptyID)
	if bIDs[0] != emptyID {
		t.Fatalf("len 0 got %d, want %d", bIDs[0], emptyID)
	}
	for i := 1; i < len(bIDs); i++ {
		for j := i + 1; j < len(bIDs); j++ {
			if bIDs[i] == bIDs[j] {
				t.Fatalf("duplicate ID for lengths %d and %d", boundaryLens[i], boundaryLens[j])
			}
		}
	}

	// Overflow past shardChunkSize (insert 20,000 strings)
	ids := make([]uint32, 20000)
	for i := range ids {
		str := fmt.Sprintf("overflow_key_%06d", i)
		ids[i] = in.Intern(str)
		if ids[i] == 0 {
			t.Fatalf("Intern returned 0 for %q", str)
		}
	}
	for i, id := range ids {
		want := fmt.Sprintf("overflow_key_%06d", i)
		if got := in.Get(id); got != want {
			t.Fatalf("Get(%d) = %q, want %q", id, got, want)
		}
		if gotID := in.Intern(want); gotID != id {
			t.Fatalf("re-Intern(%q) = %d, want %d", want, gotID, id)
		}
	}
	if in.Len() != 20003+len(boundaryLens)-1 {
		t.Fatalf("Len() = %d, want %d", in.Len(), 20003+len(boundaryLens)-1)
	}
}

func TestInternerConcurrent(t *testing.T) {
	in := NewInterner()
	const goroutines = 8
	const perG = 3000

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := range goroutines {
		go func(g int) {
			defer wg.Done()
			for i := range perG {
				// Mix of shared keys (high contention) and unique keys (triggers resize + extra chunks)
				shared := fmt.Sprintf("shared_%04d", i%500)
				id1 := in.Intern(shared)
				if in.Get(id1) != shared {
					t.Errorf("Get(shared) mismatch: got %q want %q", in.Get(id1), shared)
					return
				}

				uniq := fmt.Sprintf("uniq_%d_%05d", g, i)
				id2 := in.Intern(uniq)
				if in.Get(id2) != uniq {
					t.Errorf("Get(uniq) mismatch: got %q want %q", in.Get(id2), uniq)
					return
				}
				if in.Intern(uniq) != id2 {
					t.Errorf("re-Intern(uniq) mismatch for %q", uniq)
					return
				}
			}
		}(g)
	}
	wg.Wait()

	wantLen := 500 + goroutines*perG
	if in.Len() != wantLen {
		t.Fatalf("Len() = %d, want %d", in.Len(), wantLen)
	}
}

func BenchmarkInterner_Unique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		interner := NewInterner()
		for _, k := range benchKeysUnique {
			interner.Intern(k)
		}
	}
}

func BenchmarkInterner_Repeated(b *testing.B) {
	interner := NewInterner()
	for _, k := range benchKeysRepeated {
		interner.Intern(k)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range benchKeysRepeated {
			interner.Intern(k)
		}
	}
}

func BenchmarkInterner_RepeatedShort(b *testing.B) {
	interner := NewInterner()
	for _, k := range benchKeysShort {
		interner.Intern(k)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range benchKeysShort {
			interner.Intern(k)
		}
	}
}

func BenchmarkInterner_SmallNew(b *testing.B) {
	keys := benchKeysUnique[:200]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		in := NewInterner()
		for _, k := range keys {
			in.Intern(k)
		}
	}
}

func BenchmarkInterner_Get(b *testing.B) {
	in := NewInterner()
	ids := make([]uint32, 1000)
	for i := range ids {
		ids[i] = in.Intern(benchKeysUnique[i])
	}

	b.ResetTimer()
	var sink string
	for i := 0; i < b.N; i++ {
		for _, id := range ids {
			sink = in.Get(id)
		}
	}
	_ = sink
}

func BenchmarkInterner_ParallelRepeated(b *testing.B) {
	in := NewInterner()
	for _, k := range benchKeysUnique[:1000] {
		in.Intern(k)
	}
	keys := benchKeysUnique[:1000]

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			in.Intern(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkInterner_ParallelUnique(b *testing.B) {
	for i := 0; i < b.N; i++ {
		in := NewInterner()
		var wg sync.WaitGroup
		wg.Add(4)
		for g := range 4 {
			go func(g int) {
				defer wg.Done()
				start := g * 2500
				for _, k := range benchKeysUnique[start : start+2500] {
					in.Intern(k)
				}
			}(g)
		}
		wg.Wait()
	}
}
