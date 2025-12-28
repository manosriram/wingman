package ast

import (
	"errors"
	"log"
	"os"

	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type AST struct {
	Parser       utils.TreeSitterParserType
	NodePath     string
	NodeData     []byte
	NodeLanguage types.Language
}

func NewAST(nodePath string, nodeLanguage types.Language, parser utils.TreeSitterParserType) *AST {
	data, err := os.ReadFile(nodePath)
	if err != nil {
		log.Fatalf("Error initializing AST")
	}
	return &AST{
		NodeData:     data,
		NodePath:     nodePath,
		Parser:       parser,
		NodeLanguage: nodeLanguage,
	}
}

func (a *AST) BuildTree() (*tree_sitter.Tree, error) {
	return a.Parser.GetLanguageParser(a.NodeLanguage).Parse(a.NodeData, nil), nil
}

func (a *AST) GetNodeImports() ([]string, error) {
	switch utils.GetLanguage(a.NodePath) {
	case types.GO:
		tree, err := a.BuildTree()
		if err != nil {
			return []string{}, err
		}
		defer tree.Close()

		rootNode := tree.RootNode()
		imports := []string{}

		// Iterate through child nodes of the root
		// TODO: Move this from here
		for i := 0; i < int(rootNode.ChildCount()); i++ {
			child := rootNode.Child(uint(i))

			// Check if the node is an import_declaration
			if child.Kind() == "import_declaration" {
				// Handle single import: import "fmt"
				importSpec := child.ChildByFieldName("spec")
				if importSpec != nil && importSpec.Kind() == "import_spec" {
					path := importSpec.ChildByFieldName("path")
					if path != nil {
						importPath := a.NodeData[path.StartByte():path.EndByte()]
						imports = append(imports, string(importPath))
					}
				}

				// Handle grouped imports: import ( ... )
				for j := 0; j < int(child.ChildCount()); j++ {
					specList := child.Child(uint(j))
					if specList.Kind() == "import_spec_list" {
						for k := 0; k < int(specList.ChildCount()); k++ {
							spec := specList.Child(uint(k))
							if spec.Kind() == "import_spec" {
								path := spec.ChildByFieldName("path")
								if path != nil {
									importPath := a.NodeData[path.StartByte():path.EndByte()]
									imports = append(imports, string(importPath))
								}
							}
						}
					}
				}
			}
		}

		return imports, nil
	default:
		return []string{}, errors.New("Language support to be implemented")
	}
}
