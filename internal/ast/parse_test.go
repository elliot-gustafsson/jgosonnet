package ast

import (
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

func TestParseNumber(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
1337.42069
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Len(t, ast, 3)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagNumber, tag)

	num := ReadF64(p, 1)
	assert.Equal(t, 1337.42069, num)
}

func TestParseString(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
"banan"
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Len(t, ast, 2)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagString, tag)

	idx := ReadU32(p, 1)
	str := interner.Get(idx)
	assert.Equal(t, "banan", str)
}

func TestParseArray0(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
[]
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Len(t, ast, 1)

	p := ast.Ptr(rootIdx)
	tag := ReadTag(p, 0)
	assert.Equal(t, TagArray0, tag)
}

func TestParseArray(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
[
	true,
	false,
	123.123,
	"asdf",
]
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedLen := 1 + 1 + (1 + 2) + (1 + 1) + (1 + 1 + 4)
	assert.Len(t, ast, expectedLen)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagArray, tag)

	elemCount := ReadU32(p, 1)
	assert.Equal(t, uint32(4), elemCount)
}

func TestParseArrayComprehensions(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
[
  x
  for x in [1, 2, 3]
]
`

	node, err := jsonnet.SnippetToAST("test.jsonnet", snippet)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.NotNil(t, node)

	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	expectedLen := 1 + 1 + (1 + 2) + (1 + 1) + (1 + 1 + 4)
	assert.Len(t, ast, expectedLen)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagArray, tag)

	elemCount := ReadU32(p, 1)
	assert.Equal(t, uint32(4), elemCount)
}

func TestParseFunction(t *testing.T) {
	interner := interner.NewInterner()

	const snippet = `
function(x, y=10, z) 1
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, interner)
	if err != nil {
		t.Fatal(err.Error())
	}

	// expectedLen := 1 + 1 + (1 + 2) + (1 + 1) + (1 + 1 + 4)
	// assert.Len(t, ast, expectedLen)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagFunction, tag)

	// elemCount := ReadU32(p, 1)
	// assert.Equal(t, uint32(4), elemCount)

	// err = disassemble(ast.Base(), rootIdx, interner)
	// if err != nil {
	// 	t.Fatal(err.Error())
	// }

	// idx := ReadU32(p, 1)
	// str := interner.Get(idx)
	// assert.Equal(t, "banan", str)
}
