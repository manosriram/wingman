package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manosriram/wingman/internal/algorithm"
	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestGoFile(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "ast-test-*")
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

func setupTestGoFileWithImports(t *testing.T) (string, string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "ast-test-*")
	require.NoError(t, err)

	goModContent := `module testmodule

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	// Create a helper package
	helperDir := filepath.Join(tmpDir, "helper")
	err = os.MkdirAll(helperDir, 0755)
	require.NoError(t, err)

	helperGoContent := `package helper

func DoSomething() string {
	return "hello"
}
`
	helperGoPath := filepath.Join(helperDir, "helper.go")
	err = os.WriteFile(helperGoPath, []byte(helperGoContent), 0644)
	require.NoError(t, err)

	// Create main file that imports helper
	mainGoContent := `package main

import (
	"fmt"
	"testmodule/helper"
)

func main() {
	fmt.Println(helper.DoSomething())
}
`
	mainGoPath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(mainGoPath, []byte(mainGoContent), 0644)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return mainGoPath, tmpDir, cleanup
}

func Test_GetNodeImportsShouldReturnNonEmpty(t *testing.T) {
	testFilePath, cleanup := setupTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)
	imports, err := a.GetNodeImports()

	assert.NoError(t, err)
	assert.NotEmpty(t, imports, "imports should not be empty for a file with imports")

	d := graph.NewGraph()
	d.BuildGraphFromImports(imports)

	// The test file imports "fmt" and "os", so it should have outgoing edges
	outNodes := d.GetOutNodesOfNode(testFilePath)
	assert.NotEmpty(t, outNodes, "should have outgoing edges for imports")

	// The test file is not imported by anything, so no incoming edges
	inNodes := d.GetInNodesOfNode(testFilePath)
	assert.Empty(t, inNodes, "should have no incoming edges")
}

func Test_GetNodeImportsEmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "ast-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	goModContent := `module testmodule

go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	// File with no imports
	noImportsContent := `package main

func main() {
}
`
	noImportsPath := filepath.Join(tmpDir, "noimports.go")
	err = os.WriteFile(noImportsPath, []byte(noImportsContent), 0644)
	require.NoError(t, err)

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(noImportsPath, nil, treeSitterLanguageParser)
	imports, err := a.GetNodeImports()

	assert.NoError(t, err)
	// File with no imports should return empty slice
	assert.Empty(t, imports)
}

func Test_CalculateASTNodesScore(t *testing.T) {
	testFilePath, cleanup := setupTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)

	// Initialize the Algorithm field - this is required before calling CalculateASTNodesScore
	a.Algorithm = algorithm.NewPageRankAlgorithm()

	imports, err := a.GetNodeImports()
	assert.NoError(t, err)
	assert.NotEmpty(t, imports)

	d := graph.NewGraph()
	d.BuildGraphFromImports(imports)

	err = a.CalculateASTNodesScore(d)
	assert.NoError(t, err)

	// Verify that scores were calculated
	hasScores := false
	for k := range a.Algorithm.NodeScores {
		if a.Algorithm.NodeScores[k] > 0.0 {
			hasScores = true
			break
		}
	}
	assert.True(t, hasScores, "should have calculated positive scores for nodes")
}

func Test_CalculateASTNodesScoreEmptyGraph(t *testing.T) {
	testFilePath, cleanup := setupTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)

	// Initialize the Algorithm field
	a.Algorithm = algorithm.NewPageRankAlgorithm()

	// Empty graph
	d := graph.NewGraph()

	err := a.CalculateASTNodesScore(d)
	assert.NoError(t, err)

	// With empty graph, no scores should be calculated
	assert.Empty(t, a.Algorithm.NodeScores)
}

func Test_GraphBuildFromImports(t *testing.T) {
	testFilePath, cleanup := setupTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)
	imports, err := a.GetNodeImports()

	assert.NoError(t, err)

	d := graph.NewGraph()
	d.BuildGraphFromImports(imports)

	// Verify graph structure
	assert.NotEmpty(t, d.G, "graph should not be empty after building from imports")

	// The source file should be in the graph
	_, exists := d.G[testFilePath]
	assert.True(t, exists, "source file should be a node in the graph")
}

func Test_NewASTCreation(t *testing.T) {
	testFilePath, cleanup := setupTestGoFile(t)
	defer cleanup()

	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	a := ast.NewAST(testFilePath, nil, treeSitterLanguageParser)

	assert.NotNil(t, a)
	assert.Equal(t, testFilePath, a.NodePath)
	assert.NotEmpty(t, a.NodeData)
}
