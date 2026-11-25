package ast

import (
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNumber(t *testing.T) {
	in := interner.NewInterner()

	const snippet = `
1337.42069
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
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
	in := interner.NewInterner()

	const snippet = `
"banan"
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Len(t, ast, 2)

	p := ast.Ptr(rootIdx)

	tag := ReadTag(p, 0)
	assert.Equal(t, TagString, tag)

	idx := ReadU32(p, 1)
	str := in.Get(idx)
	assert.Equal(t, "banan", str)
}

func TestParseArray0(t *testing.T) {
	in := interner.NewInterner()

	const snippet = `
[]
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Len(t, ast, 1)

	p := ast.Ptr(rootIdx)
	tag := ReadTag(p, 0)
	assert.Equal(t, TagArray0, tag)
}

func TestParseArray(t *testing.T) {
	in := interner.NewInterner()

	const snippet = `
[
	true,
	false,
	123.123,
	"asdf",
]
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
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
	in := interner.NewInterner()

	const snippet = `
[
  x
  for x in [1, 2, 3]
]
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
	require.NoError(t, err)

	p := ast.Ptr(rootIdx)
	assert.Equal(t, TagArrayComp, ReadTag(p, 0))
	assert.Equal(t, uint32(1), ReadU32(p, 2)) // clauseCount
	assert.Equal(t, "x", in.Get(ReadU32(p, 3)))
	assert.Equal(t, uint32(0), ReadU32(p, 5)) // condCount
}

func TestParseFunction(t *testing.T) {
	in := interner.NewInterner()

	const snippet = `
function(x, y=10, z) 1
`
	ast, rootIdx, err := ParseSnippet("test.jsonnet", snippet, in)
	require.NoError(t, err)

	p := ast.Ptr(rootIdx)
	assert.Equal(t, TagFunction, ReadTag(p, 0))
	assert.Equal(t, uint32(3), ReadU32(p, 1)) // paramCount

	bodyOff := ReadU32(p, 2)
	assert.Equal(t, TagNumber, ast.ReadTag(bodyOff))
	assert.Equal(t, 1.0, ReadF64(ast.Ptr(bodyOff), 1))

	// x: no default
	assert.Equal(t, "x", in.Get(ReadU32(p, 3)))
	assert.Equal(t, NoNode, ReadU32(p, 4))

	// y: default 10 (at offset 0!)
	assert.Equal(t, "y", in.Get(ReadU32(p, 5)))
	yDef := ReadU32(p, 6)
	assert.Equal(t, uint32(0), yDef)
	assert.Equal(t, TagNumber, ast.ReadTag(yDef))
	assert.Equal(t, 10.0, ReadF64(ast.Ptr(yDef), 1))

	// z: no default
	assert.Equal(t, "z", in.Get(ReadU32(p, 7)))
	assert.Equal(t, NoNode, ReadU32(p, 8))
}

func TestParseObjectsAndComps(t *testing.T) {
	in := interner.NewInterner()

	// Empty object -> TagObject0
	a0, r0, err := ParseSnippet("t.jsonnet", `{}`, in)
	require.NoError(t, err)
	assert.Equal(t, TagObject0, a0.ReadTag(r0))

	// Static simple object -> TagObjectSimple
	a1, r1, err := ParseSnippet("t.jsonnet", `{ a: 1, b:: 2, c::: 3 }`, in)
	require.NoError(t, err)
	assert.Equal(t, TagObjectSimple, a1.ReadTag(r1))
	assert.Equal(t, uint32(3), ReadU32(a1.Ptr(r1), 1))

	// Object with local -> TagObject
	a2, r2, err := ParseSnippet("t.jsonnet", `{ local x = 1, a: x, assert self.a == 1 }`, in)
	require.NoError(t, err)
	assert.Equal(t, TagObject, a2.ReadTag(r2))
	assert.Equal(t, uint32(1), ReadU32(a2.Ptr(r2), 1)) // fields
	assert.Equal(t, uint32(1), ReadU32(a2.Ptr(r2), 2)) // locals
	assert.Equal(t, uint32(1), ReadU32(a2.Ptr(r2), 3)) // asserts

	// Object comprehension -> TagObjectComp
	a3, r3, err := ParseSnippet("t.jsonnet", `{ [x]: x for x in ["a", "b"] if x != "b" }`, in)
	require.NoError(t, err)
	assert.Equal(t, TagObjectComp, a3.ReadTag(r3))
}

func TestParseStaticChecks(t *testing.T) {
	in := interner.NewInterner()

	_, _, err := ParseSnippet("t.jsonnet", `unknownVar`, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown variable: unknownVar")

	_, _, err = ParseSnippet("t.jsonnet", `self`, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Can't use self outside of an object.")

	_, _, err = ParseSnippet("t.jsonnet", `$`, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "No top-level object found.")

	_, _, err = ParseSnippet("t.jsonnet", `super.x`, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Can't use super outside of an object.")

	_, _, err = ParseSnippet("t.jsonnet", `{ local x = 1, [x]: 2 }`, in)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Unknown variable: x")
}
