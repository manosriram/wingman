package test

import (
	"testing"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
)

func Test_BuildGraph(t *testing.T) {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(TARGET_GO_NODE, nil, treeSitterLanguageParser)
	assert.NotNil(t, a)

	imports, err := a.GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()

	assert.Empty(t, d.G[TARGET_GO_NODE])
	d.BuildGraphFromImports(imports)
	assert.NotEmpty(t, d.G[TARGET_GO_NODE])
}
