package language

import (
	"errors"
	"strings"

	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

/*
GolangStrategy implements LangStrategy.

 1. Internal imports are considered by getting the imports of each node and
    checking the import, if it starts with the module name from go.mod file,
    it is an internal import and added to the graph

 2. The list of imports is used to run an algorithm (pagerank) and find out the most important files.
    The goal of this method is to send the signatures of the most important files for the
    repo context tree to the LLM.
*/
type GolangStrategy struct {
	NodeData []byte
	NodePath string
	Parser   utils.TreeSitterParserType
}

// func NewGolangStrategy(data []byte, path string, parser utils.TreeSitterParserType) *GolangStrategy {
func NewGolangStrategy(args StrategyArgs) *GolangStrategy {
	return &GolangStrategy{
		NodeData: args.NodeData,
		NodePath: args.NodePath,
		Parser:   args.Parser,
	}
}

func (g *GolangStrategy) resolveImportNodes(args ResolveImportNodesArgs) []types.NodeImport {
	rootNode := args.RootNode
	modFileData := args.GolangModFileData

	imports := []types.NodeImport{}
	for i := uint(0); i < rootNode.ChildCount(); i++ {
		child := rootNode.Child(i)

		// Check if the node is an import_declaration
		if child.Kind() == "import_declaration" {
			importSpec := child.ChildByFieldName("spec")
			if importSpec != nil && importSpec.Kind() == "import_spec" {
				path := importSpec.ChildByFieldName("path")
				if path != nil {
					importPath := strings.Trim(string(g.NodeData[path.StartByte():path.EndByte()]), "\"")
					// If the import is internal, get the import pkg and add it to "imports"
					if strings.HasPrefix(importPath, modFileData) {
						importPath, found := strings.CutPrefix(importPath, modFileData)
						if found {
							importPathSplit := strings.Split(strings.TrimLeft(importPath, "/"), "/")
							if len(importPathSplit) > 1 {
								importPath = importPathSplit[1]
							} /*  else { */
						}
						imports = append(imports, types.NodeImport{
							ImportPackage: importPath,
							FilePath:      g.NodePath,
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
								importPath := strings.Trim(string(g.NodeData[path.StartByte():path.EndByte()]), "\"")

								// Only if the import is internal
								if strings.HasPrefix(importPath, modFileData) {

									importPath, found := strings.CutPrefix(importPath, modFileData)
									if found {
										importPathSplit := strings.Split(strings.TrimLeft(importPath, "/"), "/")
										if len(importPathSplit) > 1 {
											importPath = importPathSplit[1]
										}
									}
									imports = append(imports, types.NodeImport{
										ImportPackage: importPath,
										FilePath:      g.NodePath,
									})
								}
							}
						}
					}
				}
			}
		}
	}
	return imports
}

func (g *GolangStrategy) GetNodeImportList() ([]types.NodeImport, error) {

	tree := g.Parser.GetLanguageParser(types.GOLANG).Parse(g.NodeData, nil)
	if tree == nil {
		return []types.NodeImport{}, errors.New("Error initializing Parser")
	}
	defer tree.Close()

	modFilePath, err := utils.FindGoModPath(g.NodePath)
	if err != nil {
		return []types.NodeImport{}, errors.New("Error reading go.mod file")
	}
	modFile, err := utils.ReadGoModFile(modFilePath)
	if err != nil {
		return []types.NodeImport{}, errors.New("Error reading go.mod file")
	}

	modFileSplit := strings.Split(string(modFile), "\n")
	if len(modFileSplit) < 1 {
		return []types.NodeImport{}, errors.New("Error reading go.mod file")
	}

	moduleNameSplit := strings.Split(modFileSplit[0], " ")
	if len(moduleNameSplit) < 1 {
		return []types.NodeImport{}, errors.New("Error reading go.mod file")
	}
	modFileData := moduleNameSplit[1]

	rootNode := tree.RootNode()

	return g.resolveImportNodes(ResolveImportNodesArgs{
		RootNode:          rootNode,
		GolangModFileData: modFileData,
	}), nil
}
