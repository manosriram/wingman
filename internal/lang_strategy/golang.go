package langstrategy

import (
	"strings"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type GolangStrategy struct {
	a *ast.AST
}

func (g *GolangStrategy) GetNodeImports() ([]types.NodeImport, error) {
	switch utils.GetLanguage(g.a.NodePath) {
	case types.GOLANG:
		tree, err := g.a.BuildTree()
		if err != nil {
			return []types.NodeImport{}, err
		}
		defer tree.Close()

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
						importPath := strings.Trim(string(g.a.NodeData[path.StartByte():path.EndByte()]), "\"")
						imports = append(imports, types.NodeImport{
							ImportPath: importPath,
							FilePath:   g.a.NodePath,
						})
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
									importPath := strings.Trim(string(g.a.NodeData[path.StartByte():path.EndByte()]), "\"")
									imports = append(imports, types.NodeImport{
										ImportPath: importPath,
										FilePath:   g.a.NodePath,
									})
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
