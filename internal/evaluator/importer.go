package evaluator

import (
	"fmt"
	"os"
	"sync"

	"github.com/elliot-gustafsson/jgosonnet/internal/ast"
	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet"
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
		node, err := jsonnet.SnippetToAST(name, data)
		if err != nil {
			entry.err = fmt.Errorf("failed to resolve snippet %s, err: %w", name, err)
			return
		}

		builder := ast.NewAstBuilder(t.interner)

		entry.tree, entry.err = builder.Parse(node)
	})

	return entry.tree, entry.err
}

func (t *AstImporter) ResolveImport(filePath string) (*ast.AST, error) {

	actual, _ := t.cache.LoadOrStore(filePath, &cacheEntry{})
	entry := actual.(*cacheEntry)

	entry.once.Do(func() {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			if !os.IsNotExist(err) {
				entry.err = fmt.Errorf("failed importing file: %s, err: %w", filePath, err)
			} else {
				entry.err = err
			}
			return
		}

		node, err := jsonnet.SnippetToAST(filePath, string(fileData))
		if err != nil {
			entry.err = fmt.Errorf("failed to resolve import %s, err: %w", filePath, err)
			return
		}

		builder := ast.NewAstBuilder(t.interner)

		entry.tree, entry.err = builder.Parse(node)
	})

	return entry.tree, entry.err
}
