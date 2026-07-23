package evaluator

import (
	"fmt"
	"os"
	"sync"

	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	
)

type Importer struct {
	JPaths  []string
	BaseStd Value

	astImporter *AstImporter
	cache       map[string]Value
}

type AstImporter struct {
	cache    sync.Map // map[string]*cacheEntry
	interner *interner.Interner
}

type cacheEntry struct {
	once sync.Once
	tree *ast.AST
	err  error
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

func NewAstImporter(interner *interner.Interner) *AstImporter {
	return &AstImporter{
		interner: interner,
	}
}

func (i *Importer) Set(path string, v Value) {
	i.cache[path] = v
}

func (i *Importer) Get(path string) Value {
	return i.cache[path]
}

func (i *Importer) ResolveSnippet(name, data string) (*ast.AST, error) {
	return i.astImporter.ResolveSnippet(name, data)
}

func (i *Importer) ResolveImport(filePath string) (*ast.AST, error) {
	return i.astImporter.ResolveImport(filePath)
}

func (t *AstImporter) ResolveSnippet(name, data string) (*ast.AST, error) {

	actual, _ := t.cache.LoadOrStore(name, &cacheEntry{})
	entry := actual.(*cacheEntry)

	entry.once.Do(func() {
		entry.tree, entry.err = ast.ParseSnippet(name, data, t.interner)
	})

	return entry.tree, entry.err
}

func (t *AstImporter) ResolveImport(filePath string) (*ast.AST, error) {

	actual, _ := t.cache.LoadOrStore(filePath, &cacheEntry{})
	entry := actual.(*cacheEntry)

	entry.once.Do(func() {
		_, err := os.ReadFile(filePath)
		if err != nil {
			if !os.IsNotExist(err) {
				entry.err = fmt.Errorf("failed importing file: %s, err: %w", filePath, err)
			} else {
				entry.err = err
			}
			return
		}

		entry.tree, entry.err = ast.Parse(filePath, t.interner)
	})

	return entry.tree, entry.err
}
