package test

import (
	"fmt"
	"testing"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TARGET_DIR     = "/Users/manosriram/go/src/go2java"
	TARGET_GO_NODE = "/Users/manosriram/go/src/go2java/ast.go"
)

func Test_GetNodeImports(t *testing.T) {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	imports, err := ast.NewAST(TARGET_GO_NODE, treeSitterLanguageParser).GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()
	d.BuildGraphFromImports([]types.NodeImport{
		{
			FilePath:      "a.go",
			ImportPackage: "os",
		},
		{
			FilePath:      "a.go",
			ImportPackage: "fmt",
		},
		{
			FilePath:      "a.go",
			ImportPackage: "testing",
		},
	})

	fmt.Println(d.GetInNodesOfNode("os"))

}
