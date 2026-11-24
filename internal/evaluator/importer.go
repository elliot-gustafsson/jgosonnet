package evaluator

import (
	"os"
	"sync"
	"unsafe"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet"
)

var parseCallPool = sync.Pool{
	New: func() any { return new(parseCall) },
}

type Importer struct {
	JPaths  []string
	BaseStd Value

	astImporter *AstImporter
	cache       map[string]Value
}

type parseCall struct {
	done chan struct{}
	node *IRNode
	err  error
}

type inFlightParse struct {
	key  string
	call *parseCall
}

type AstImporter struct {
	Interner  *interner.Interner
	stdLookup StdLookupFunc
	mu        sync.RWMutex
	astCache  map[string]*IRNode
	inFlights []inFlightParse
}

func NewImporter(jPaths []string, baseStd Value, astImporter *AstImporter) *Importer {
	return &Importer{
		JPaths:      jPaths,
		BaseStd:     baseStd,
		cache:       make(map[string]Value, 8),
		astImporter: astImporter,
	}
}

func NewAstImporter(inter *interner.Interner, stdLookup StdLookupFunc) *AstImporter {
	if inter == nil {
		inter = interner.NewInterner()
	}
	return &AstImporter{
		Interner:  inter,
		stdLookup: stdLookup,
		astCache:  make(map[string]*IRNode, 8),
		inFlights: make([]inFlightParse, 0, 4),
	}
}

func (i *Importer) Set(path string, v Value) {
	i.cache[path] = v
}

func (i *Importer) Get(path string) Value {
	return i.cache[path]
}

func (i *Importer) ResolveSnippet(name, data string) (*IRNode, error) {
	return i.astImporter.ResolveSnippet(name, data)
}

func (i *Importer) ResolveImport(filePath string) (*IRNode, error) {
	return i.astImporter.ResolveImport(filePath)
}

func (t *AstImporter) ResolveImport(filePath string) (*IRNode, error) {
	t.mu.RLock()
	node, exist := t.astCache[filePath]
	t.mu.RUnlock()
	if exist {
		return node, nil
	}
	return t.resolveImportSlow(filePath)
}

//go:noinline
func (t *AstImporter) resolveImportSlow(filePath string) (node *IRNode, err error) {

	node, slot, call, isLeader, err := t.claim(filePath)
	if !isLeader {
		return
	}
	defer func() {
		t.release(slot, filePath, call, node, err)
	}()

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	dataStr := unsafe.String(unsafe.SliceData(fileData), len(fileData))

	rawNode, err := jsonnet.SnippetToAST(filePath, dataStr)
	if err != nil {
		return nil, err
	}

	node = CompileASTWithStd(rawNode, t.Interner, t.stdLookup)
	return node, nil
}

func (t *AstImporter) ResolveSnippet(name, data string) (*IRNode, error) {
	t.mu.RLock()
	node, exist := t.astCache[name]
	t.mu.RUnlock()
	if exist {
		return node, nil
	}

	return t.resolveSnippetSlow(name, data)
}

//go:noinline
func (t *AstImporter) resolveSnippetSlow(name, data string) (node *IRNode, err error) {
	node, slot, call, isLeader, err := t.claim(name)
	if !isLeader {
		return
	}
	defer func() {
		t.release(slot, name, call, node, err)
	}()

	rawNode, err := jsonnet.SnippetToAST(name, data)
	if err != nil {
		return nil, err
	}

	node = CompileASTWithStd(rawNode, t.Interner, t.stdLookup)
	return node, nil
}

func (t *AstImporter) claim(key string) (cached *IRNode, slot int, call *parseCall, isLeader bool, err error) {
	// slow path
	t.mu.Lock()
	// double check cache under write lock
	if existing, exist := t.astCache[key]; exist {
		t.mu.Unlock()
		return existing, 0, nil, false, nil
	}

	freeSlot := -1
	for i := range t.inFlights {
		if t.inFlights[i].key == key {
			// create channel on first waiter
			if t.inFlights[i].call.done == nil {
				t.inFlights[i].call.done = make(chan struct{})
			}
			// an in flight parse is running, wait for it
			waitingCall := t.inFlights[i].call
			t.mu.Unlock()
			<-waitingCall.done
			return waitingCall.node, 0, nil, false, waitingCall.err
		}
		if freeSlot == -1 && t.inFlights[i].key == "" {
			freeSlot = i
		}
	}

	// register an in flight parse, use first free slot or append to the end
	call = parseCallPool.Get().(*parseCall)
	entry := inFlightParse{key: key, call: call}
	slot = freeSlot
	if freeSlot >= 0 {
		t.inFlights[freeSlot] = entry
	} else {
		slot = len(t.inFlights)
		t.inFlights = append(t.inFlights, entry)
	}
	t.mu.Unlock()

	return nil, slot, call, true, nil
}

func (t *AstImporter) release(slot int, key string, call *parseCall, node *IRNode, err error) {
	call.node = node
	call.err = err

	// register in cache and clear in flight parse slot
	t.mu.Lock()
	if err == nil {
		t.astCache[key] = node
	}
	t.inFlights[slot] = inFlightParse{}
	t.mu.Unlock()

	// release waiters if any
	if call.done != nil {
		close(call.done)
		return
	}

	// recycle call if noone waited
	*call = parseCall{}
	parseCallPool.Put(call)
}
