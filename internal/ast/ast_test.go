package ast

import (
	"fmt"
	"testing"

	"github.com/elliot-gustafsson/jgosonnet/internal/interner"
	"github.com/google/go-jsonnet"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

	interner := interner.NewInterner()

	fileData := `

1 + "_hej"

`

	node, err := jsonnet.SnippetToAST("test.jsonnet", fileData)
	if err != nil {
		t.Fatal(err.Error())
	}

	builder := AstBuilder{
		Interner: interner,
	}

	tree, err := builder.Parse(node)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.NotNil(t, tree)
	fmt.Println(tree.Nodes)

}
