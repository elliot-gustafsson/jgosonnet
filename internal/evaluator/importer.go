package evaluator

import (
	"sync"
	"unsafe"

	fast "github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
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
	prog *fast.Program
	err  error
}

type inFlightParse struct {
	key  string
	call *parseCall
}

type AstImporter struct {
	interner  *interner.Interner
	mu        sync.RWMutex
	astCache  map[string]*fast.Program
	byBase    map[uintptr]*fast.Program
	inFlights []inFlightParse
}

func NewImporter(jPaths []string, baseStd Value, astImporter *AstImporter) *Importer {
	return &Importer{
		JPaths:      jPaths,
		BaseStd:     baseStd,
		cache:       make(map[string]Value, 32),
		astImporter: astImporter,
	}
}

func NewAstImporter(in *interner.Interner) *AstImporter {
	return &AstImporter{
		interner:  in,
		astCache:  make(map[string]*fast.Program, 64),
		byBase:    make(map[uintptr]*fast.Program, 64),
		inFlights: make([]inFlightParse, 0, 16),
	}
}

func (i *Importer) Set(path string, v Value) {
	i.cache[path] = v
}

func (i *Importer) Get(path string) Value {
	return i.cache[path]
}

func (i *Importer) ResolveSnippet(name, data string) (*fast.Program, error) {
	return i.astImporter.ResolveSnippet(name, data)
}

func (i *Importer) ResolveImport(filePath string) (*fast.Program, error) {
	return i.astImporter.ResolveImport(filePath)
}

func (i *Importer) LookupProgram(bp unsafe.Pointer) *fast.Program {
	if i == nil || i.astImporter == nil {
		return nil
	}
	return i.astImporter.LookupProgram(bp)
}

func (t *AstImporter) LookupProgram(bp unsafe.Pointer) *fast.Program {
	t.mu.RLock()
	prog := t.byBase[uintptr(bp)]
	t.mu.RUnlock()
	return prog
}

func (t *AstImporter) ResolveImport(filePath string) (*fast.Program, error) {
	t.mu.RLock()
	prog, exist := t.astCache[filePath]
	t.mu.RUnlock()
	if exist {
		return prog, nil
	}
	return t.resolveImportSlow(filePath)
}

//go:noinline
func (t *AstImporter) resolveImportSlow(filePath string) (prog *fast.Program, err error) {
	prog, slot, call, isLeader, err := t.claim(filePath)
	if !isLeader {
		return
	}
	defer func() {
		t.release(slot, filePath, call, prog, err)
	}()

	prog, err = fast.ParseProgram(filePath, t.interner)
	return
}

func (t *AstImporter) ResolveSnippet(name, data string) (*fast.Program, error) {
	t.mu.RLock()
	prog, exist := t.astCache[name]
	t.mu.RUnlock()
	if exist {
		return prog, nil
	}

	return t.resolveSnippetSlow(name, data)
}

//go:noinline
func (t *AstImporter) resolveSnippetSlow(name, data string) (prog *fast.Program, err error) {
	prog, slot, call, isLeader, err := t.claim(name)
	if !isLeader {
		return
	}
	defer func() {
		t.release(slot, name, call, prog, err)
	}()

	prog, err = fast.ParseProgramSnippet(name, data, t.interner)
	return
}

func (t *AstImporter) claim(key string) (cached *fast.Program, slot int, call *parseCall, isLeader bool, err error) {
	t.mu.Lock()
	if existing, exist := t.astCache[key]; exist {
		t.mu.Unlock()
		return existing, 0, nil, false, nil
	}

	freeSlot := -1
	for i := range t.inFlights {
		if t.inFlights[i].key == key {
			if t.inFlights[i].call.done == nil {
				t.inFlights[i].call.done = make(chan struct{})
			}
			waitingCall := t.inFlights[i].call
			t.mu.Unlock()
			<-waitingCall.done
			return waitingCall.prog, 0, nil, false, waitingCall.err
		}
		if freeSlot == -1 && t.inFlights[i].key == "" {
			freeSlot = i
		}
	}

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

func (t *AstImporter) release(slot int, key string, call *parseCall, prog *fast.Program, err error) {
	call.prog = prog
	call.err = err

	t.mu.Lock()
	if err == nil && prog != nil {
		t.astCache[key] = prog
		t.byBase[uintptr(prog.Base())] = prog
	}
	t.inFlights[slot] = inFlightParse{}
	t.mu.Unlock()

	if call.done != nil {
		close(call.done)
		return
	}

	*call = parseCall{}
	parseCallPool.Put(call)
}
