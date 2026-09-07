package evaluator

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Importer struct {
	JPaths  []string
	BaseStd Value

	astImporter *AstImporter
	cache       map[string]Value
}

const (
	defaultImporterTableSize = 2048 // Must be a power of 2
	chunkSize                = 1024
	chunkShift               = 10
	chunkMask                = chunkSize - 1
)

// Lean entry struct (no 16-byte error interface!)
type entry struct {
	key  string
	node ast.Node
	once sync.Once
}

type AstImporter struct {
	mu       sync.Mutex
	table    atomic.Pointer[[]uint64] // []uint64: top 32 bits = hash, bottom 32 bits = idx (0 GC pointers!)
	chunks   atomic.Pointer[[]*[chunkSize]entry]
	count    uint32           // Plain uint32 (mutated only under mu.Lock)
	errCache map[string]error // Separate error cache for rare import failures
}

func NewImporter(jPaths []string, baseStd Value, astImporter *AstImporter) *Importer {
	return &Importer{
		JPaths:  jPaths,
		BaseStd: baseStd,
		// TODO: maybe use slices?
		cache:       make(map[string]Value, 32),
		astImporter: astImporter,
	}
}

func NewAstImporter() *AstImporter {
	i := &AstImporter{
		errCache: make(map[string]error),
	}

	firstChunk := new([chunkSize]entry)
	initialChunks := make([]*[chunkSize]entry, 1)
	initialChunks[0] = firstChunk

	i.chunks.Store(&initialChunks)
	i.count = 1 // Burn index 0 as sentinel

	initialTable := make([]uint64, defaultImporterTableSize)
	i.table.Store(&initialTable)

	return i
}

func (i *Importer) Set(path string, v Value) {
	i.cache[path] = v
}

func (i *Importer) Get(path string) Value {
	return i.cache[path]
}

func (i *Importer) ResolveSnippet(name, data string) (ast.Node, error) {
	return i.astImporter.ResolveSnippet(name, data)
}

func (i *Importer) ResolveImport(filePath string) (ast.Node, error) {
	return i.astImporter.ResolveImport(filePath)
}

func (t *AstImporter) ResolveSnippet(name, data string) (ast.Node, error) {
	e := t.getOrCreateEntry(name)

	e.once.Do(func() {
		node, err := jsonnet.SnippetToAST(name, data)
		if err != nil {
			t.setErr(name, fmt.Errorf("failed to resolve snippet %s, err: %w", name, err))
			return
		}
		e.node = node
	})

	if e.node == nil {
		return nil, t.getErr(name)
	}
	return e.node, nil
}

func (t *AstImporter) ResolveImport(filePath string) (ast.Node, error) {
	e := t.getOrCreateEntry(filePath)

	e.once.Do(func() {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			t.setErr(filePath, err)
			return
		}

		dataStr := unsafe.String(unsafe.SliceData(fileData), len(fileData))
		node, err := jsonnet.SnippetToAST(filePath, dataStr)
		if err != nil {
			t.setErr(filePath, fmt.Errorf("failed to resolve import %s, err: %w", filePath, err))
			return
		}
		e.node = node
	})

	if e.node == nil {
		return nil, t.getErr(filePath)
	}
	return e.node, nil
}

func (t *AstImporter) getOrCreateEntry(key string) *entry {
	h := uint32(interner.HashString(key)) // 2ns Rapidhash

	// 1. Lock-Free Read Fast-Path for warm cache hits
	tablePtr := t.table.Load()
	table := *tablePtr
	mask := uint32(len(table) - 1)
	idx := h & mask
	chunks := *t.chunks.Load()

	for {
		entry64 := atomic.LoadUint64(&table[idx])
		if entry64 == 0 {
			break // Not found in fast path
		}

		eHash := uint32(entry64 >> 32)
		eIdx := uint32(entry64)

		if eHash == h {
			e := &chunks[eIdx>>chunkShift][eIdx&chunkMask]
			if e.key == key {
				return e // 100% Lock-free read hit!
			}
		}
		idx = (idx + 1) & mask
	}

	// 2. Cold-Path: Acquire mutex for new entry creation
	t.mu.Lock()
	defer t.mu.Unlock()

	tablePtr = t.table.Load()
	table = *tablePtr
	mask = uint32(len(table) - 1)
	idx = h & mask
	chunks = *t.chunks.Load()

	for {
		entry64 := atomic.LoadUint64(&table[idx])

		if entry64 == 0 {
			newId := t.count
			chunkIdx := newId >> chunkShift
			itemIdx := newId & chunkMask

			oldChunks := *t.chunks.Load()
			if chunkIdx >= uint32(len(oldChunks)) {
				newChunks := make([]*[chunkSize]entry, len(oldChunks)+1)
				copy(newChunks, oldChunks)

				newChunk := new([chunkSize]entry)
				newChunks[len(oldChunks)] = newChunk

				t.chunks.Store(&newChunks)
				chunks = newChunks
			}

			e := &chunks[chunkIdx][itemIdx]
			e.key = key

			// Pack hash (top 32 bits) and idx (bottom 32 bits) into uint64
			packed := (uint64(h) << 32) | uint64(newId)
			atomic.StoreUint64(&table[idx], packed)

			t.count++

			// Resize table if >50% full
			if t.count > uint32(len(table))/2 {
				t.resize(table)
			}
			return e
		}

		eHash := uint32(entry64 >> 32)
		eIdx := uint32(entry64)

		if eHash == h {
			e := &chunks[eIdx>>chunkShift][eIdx&chunkMask]
			if e.key == key {
				return e
			}
		}

		idx = (idx + 1) & mask
	}
}

func (t *AstImporter) resize(oldTable []uint64) {
	newSize := len(oldTable) * 2
	newMask := uint32(newSize - 1)
	newTable := make([]uint64, newSize)

	for _, entry64 := range oldTable {
		if entry64 == 0 {
			continue
		}

		h := uint32(entry64 >> 32)
		idx := h & newMask
		for {
			if newTable[idx] == 0 {
				newTable[idx] = entry64
				break
			}
			idx = (idx + 1) & newMask
		}
	}

	t.table.Store(&newTable)
}

func (t *AstImporter) setErr(key string, err error) {
	t.mu.Lock()
	t.errCache[key] = err
	t.mu.Unlock()
}

func (t *AstImporter) getErr(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.errCache[key]
}
