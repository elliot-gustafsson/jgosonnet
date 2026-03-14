package evaluator

import (
	"fmt"
	"os"
	"sync"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

type Importer struct {
	JPaths      []string
	ImportScope uint32

	astImporter *AstImporter
	cache       map[string]Value
}

type AstImporter struct {
	cacheMu  sync.RWMutex
	astCache map[string]ast.Node
}

func NewImporter(scopeId uint32, jPaths []string, astImporter *AstImporter) *Importer {
	return &Importer{
		ImportScope: scopeId,
		JPaths:      jPaths,
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

func (i *Importer) ResolveImport(filePath string) (ast.Node, error) {
	return i.astImporter.ResolveImport(filePath)
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

	importedNode, err = jsonnet.SnippetToAST(filePath, string(fileData))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve import %s, err: %w", filePath, err)
	}

	t.cacheMu.Lock()
	t.astCache[filePath] = importedNode
	t.cacheMu.Unlock()

	return importedNode, nil
}
