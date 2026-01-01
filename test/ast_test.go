package test

import (
	"testing"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TARGET_DIR     = "/Users/manosriram/go/src/go2java"
	TARGET_GO_NODE = "../internal/repository/repository.go"
)

func Test_GetNodeImportsShouldReturnNonEmpty(t *testing.T) {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	imports, err := ast.NewAST(TARGET_GO_NODE, treeSitterLanguageParser).GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()
	d.BuildGraphFromImports(imports)

	assert.NotEmpty(t, d.GetOutNodesOfNode(TARGET_GO_NODE))
	assert.Empty(t, d.GetInNodesOfNode(TARGET_GO_NODE))
}

func Test_CalculateASTNodesScore(t *testing.T) {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(TARGET_GO_NODE, treeSitterLanguageParser)
	imports, err := a.GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()
	d.BuildGraphFromImports(imports)

	for range imports {
		err = a.CalculateASTNodesScore(d)
		assert.Nil(t, err)
	}
	for k := range a.Algorithm.NodeScores {
		assert.Greater(t, a.Algorithm.NodeScores[k], 0.0)
	}
}
