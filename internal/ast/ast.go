package ast

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/manosriram/wingman/internal/algorithm"
	"github.com/manosriram/wingman/internal/dag"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Signature struct {
}

type AST struct {
	Parser       utils.TreeSitterParserType
	NodePath     string
	NodeData     []byte
	NodeLanguage types.Language
	Signatures   map[string]Signature
	// Algorithm    algorithm.ContextAlgorithm
	Algorithm *algorithm.PageRankAlgorithm
}

func NewAST(nodePath string, parser utils.TreeSitterParserType) *AST {
	// if types.DEFAULT_CONTEXT_ALGORITHM == types.PAGERANK_CONTEXT_ALGORITHM {

	// } else {

	// }

	data, err := os.ReadFile(nodePath)
	if err != nil {
		log.Fatalf("Error initializing AST")
	}
	return &AST{
		NodeData:     data,
		NodePath:     nodePath,
		NodeLanguage: utils.GetLanguage(nodePath),
		Parser:       parser,
		Signatures:   make(map[string]Signature),
		Algorithm:    algorithm.NewPageRankAlgorithm(),
	}
}

func (a *AST) BuildTree() (*tree_sitter.Tree, error) {
	return a.Parser.GetLanguageParser(a.NodeLanguage).Parse(a.NodeData, nil), nil
}

func (a *AST) GetNodeImports() ([]types.NodeImport, error) {
	switch utils.GetLanguage(a.NodePath) {
	case types.GOLANG:
		tree, err := a.BuildTree()
		if err != nil {
			return []types.NodeImport{}, err
		}
		defer tree.Close()

		modFile, err := os.ReadFile("go.mod")
		mo := strings.Split(strings.Split(string(modFile), "\n")[0], " ")[1]

		rootNode := tree.RootNode()
		imports := []types.NodeImport{}

		for i := uint(0); i < rootNode.ChildCount(); i++ {
			child := rootNode.Child(i)

			// Check if the node is an import_declaration
			if child.Kind() == "import_declaration" {
				// Handle single import: import "fmt"
				importSpec := child.ChildByFieldName("spec")
				if importSpec != nil && importSpec.Kind() == "import_spec" {
					path := importSpec.ChildByFieldName("path")
					if path != nil {
						importPath := strings.Trim(string(a.NodeData[path.StartByte():path.EndByte()]), "\"")
						// Only if the import is internal
						if strings.HasPrefix(importPath, mo) {
							imports = append(imports, types.NodeImport{
								ImportPath: importPath,
								FilePath:   a.NodePath,
							})
						}

					}
				}

				// Handle grouped imports: import ( ... )
				for j := uint(0); j < child.ChildCount(); j++ {
					specList := child.Child(j)
					if specList.Kind() == "import_spec_list" {
						for k := uint(0); k < specList.ChildCount(); k++ {
							spec := specList.Child(k)
							if spec.Kind() == "import_spec" {
								path := spec.ChildByFieldName("path")
								if path != nil {
									importPath := strings.Trim(string(a.NodeData[path.StartByte():path.EndByte()]), "\"")

									// Only if the import is internal
									if strings.HasPrefix(importPath, mo) {
										imports = append(imports, types.NodeImport{
											ImportPath: importPath,
											FilePath:   a.NodePath,
										})
									}
								}
							}
						}
					}
				}
			}
		}

		return imports, nil
	default:
		return []types.NodeImport{}, nil // Unsupported language is not an error, hence "nil"
	}
}

func (a *AST) CalculateASTNodesScore(dag *dag.DAG) error {
	// switch?
	if a.Algorithm.GetAlgorithmType() == types.PAGERANK_CONTEXT_ALGORITHM {
		a.Algorithm.CalculateScore(dag)
	} else {
		return errors.New("Algorithm not implemented")
	}
	return nil
}
