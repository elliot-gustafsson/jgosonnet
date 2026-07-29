package evaluator

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Importer struct {
	JPaths  []string
	BaseStd Value

	astImporter *AstImporter
	cache       map[string]Value
}

type AstImporter struct {
	cacheMu  sync.RWMutex
	astCache map[string]ast.Node
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
	return &AstImporter{
		// TODO: maybe use slices?
		astCache: make(map[string]ast.Node, 32),
	}
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

	t.cacheMu.RLock()
	importedNode, exist := t.astCache[name]
	t.cacheMu.RUnlock()

	if exist {
		return importedNode, nil
	}

	importedNode, err := jsonnet.SnippetToAST(name, data)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve snippet %s, err: %w", name, err)
	}

	t.cacheMu.Lock()
	defer t.cacheMu.Unlock()

	// Double check import cache again
	existing, exist := t.astCache[name]
	if exist {
		return existing, nil
	}

	t.astCache[name] = importedNode

	return importedNode, nil
}

func (t *AstImporter) ResolveImport(filePath string) (ast.Node, error) {

	t.cacheMu.RLock()
	importedNode, exist := t.astCache[filePath]
	t.cacheMu.RUnlock()

	if exist {
		return importedNode, nil
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed importing file: %s, err: %w", filePath, err)
	}

	dataStr := unsafe.String(unsafe.SliceData(fileData), len(fileData))

	importedNode, err = jsonnet.SnippetToAST(filePath, dataStr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve import %s, err: %w", filePath, err)
	}

	t.cacheMu.Lock()
	defer t.cacheMu.Unlock()

	// Double check import cache again
	existing, exist := t.astCache[filePath]
	if exist {
		return existing, nil
	}

	t.astCache[filePath] = importedNode

	return importedNode, nil
}
