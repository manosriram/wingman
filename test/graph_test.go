package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGraphTestGoFile(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "graph-test-*")
	require.NoError(t, err)

	goModContent := `module testmodule

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	mainGoContent := `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello")
	os.Exit(0)
}
`
	mainGoPath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(mainGoPath, []byte(mainGoContent), 0644)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return mainGoPath, cleanup
}

func Test_BuildGraph(t *testing.T) {
	testFilePath, cleanup := setupGraphTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)
	assert.NotNil(t, a)

	imports, err := a.GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()

	assert.Empty(t, d.G[testFilePath])
	d.BuildGraphFromImports(imports)
	assert.NotEmpty(t, d.G[testFilePath])
}

func Test_NewGraph(t *testing.T) {
	g := graph.NewGraph()

	assert.NotNil(t, g)
	assert.NotNil(t, g.G)
	assert.Empty(t, g.G)
}

func Test_BuildGraphFromImports(t *testing.T) {
	g := graph.NewGraph()

	imports := []types.NodeImport{
		{FilePath: "/path/to/main.go", ImportPackage: "fmt"},
		{FilePath: "/path/to/main.go", ImportPackage: "os"},
		{FilePath: "/path/to/other.go", ImportPackage: "fmt"},
	}

	g.BuildGraphFromImports(imports)

	// main.go should have 2 outgoing edges (fmt, os)
	assert.Len(t, g.G["/path/to/main.go"], 2)

	// other.go should have 1 outgoing edge (fmt)
	assert.Len(t, g.G["/path/to/other.go"], 1)

	// fmt should exist as a node (destination)
	_, exists := g.G["fmt"]
	assert.True(t, exists)
}

func Test_GetOutNodesOfNode(t *testing.T) {
	g := graph.NewGraph()

	imports := []types.NodeImport{
		{FilePath: "/path/to/main.go", ImportPackage: "fmt"},
		{FilePath: "/path/to/main.go", ImportPackage: "os"},
	}

	g.BuildGraphFromImports(imports)

	outNodes := g.GetOutNodesOfNode("/path/to/main.go")
	assert.Len(t, outNodes, 2)

	// Verify the out nodes contain fmt and os
	outValues := make([]string, len(outNodes))
	for i, node := range outNodes {
		outValues[i] = node.NodeValue
	}
	assert.Contains(t, outValues, "fmt")
	assert.Contains(t, outValues, "os")
}

func Test_GetOutNodesOfNode_Empty(t *testing.T) {
	g := graph.NewGraph()

	outNodes := g.GetOutNodesOfNode("/nonexistent/path.go")
	assert.Empty(t, outNodes)
}

func Test_GetInNodesOfNode(t *testing.T) {
	g := graph.NewGraph()

	imports := []types.NodeImport{
		{FilePath: "/path/to/main.go", ImportPackage: "fmt"},
		{FilePath: "/path/to/other.go", ImportPackage: "fmt"},
		{FilePath: "/path/to/another.go", ImportPackage: "fmt"},
	}

	g.BuildGraphFromImports(imports)

	// fmt should have 3 incoming edges
	inNodes := g.GetInNodesOfNode("fmt")
	assert.Len(t, inNodes, 3)
	assert.Contains(t, inNodes, "/path/to/main.go")
	assert.Contains(t, inNodes, "/path/to/other.go")
	assert.Contains(t, inNodes, "/path/to/another.go")
}

func Test_GetInNodesOfNode_Empty(t *testing.T) {
	g := graph.NewGraph()

	imports := []types.NodeImport{
		{FilePath: "/path/to/main.go", ImportPackage: "fmt"},
	}

	g.BuildGraphFromImports(imports)

	// main.go has no incoming edges (nothing imports it)
	inNodes := g.GetInNodesOfNode("/path/to/main.go")
	assert.Empty(t, inNodes)
}

func Test_NewGraphNode(t *testing.T) {
	node := graph.NewGraphNode("test-value")

	assert.Equal(t, "test-value", node.NodeValue)
}

func Test_GraphWithEmptyImports(t *testing.T) {
	g := graph.NewGraph()

	imports := []types.NodeImport{}

	g.BuildGraphFromImports(imports)

	assert.Empty(t, g.G)
}

func Test_GraphCircularDependency(t *testing.T) {
	g := graph.NewGraph()

	// Simulate circular dependency: A imports B, B imports A
	imports := []types.NodeImport{
		{FilePath: "A", ImportPackage: "B"},
		{FilePath: "B", ImportPackage: "A"},
	}

	g.BuildGraphFromImports(imports)

	// A should have B as out node
	outA := g.GetOutNodesOfNode("A")
	assert.Len(t, outA, 1)
	assert.Equal(t, "B", outA[0].NodeValue)

	// B should have A as out node
	outB := g.GetOutNodesOfNode("B")
	assert.Len(t, outB, 1)
	assert.Equal(t, "A", outB[0].NodeValue)

	// A should have B as in node
	inA := g.GetInNodesOfNode("A")
	assert.Len(t, inA, 1)
	assert.Contains(t, inA, "B")

	// B should have A as in node
	inB := g.GetInNodesOfNode("B")
	assert.Len(t, inB, 1)
	assert.Contains(t, inB, "A")
}
