package language

import (
	"errors"
	"slices"
	"strings"

	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
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
	PkgPaths map[string][]string
	Parser   utils.TreeSitterParserType
	Visit    map[*tree_sitter.Node]bool
}

// func NewGolangStrategy(data []byte, path string, parser utils.TreeSitterParserType) *GolangStrategy {
func NewGolangStrategy(args StrategyArgs) *GolangStrategy {
	return &GolangStrategy{
		NodeData: args.NodeData,
		NodePath: args.NodePath,
		Parser:   args.Parser,
		Visit:    make(map[*tree_sitter.Node]bool),
		PkgPaths: args.PkgPaths,
	}
}

func (g *GolangStrategy) resolveImportNodes(args ResolveImportNodesArgs) []types.NodeImport {

	rootNode := args.RootNode

	g.Visit[rootNode] = true
	// modFileData := args.GolangModFileData

	imports := []types.NodeImport{}

	if rootNode.Kind() == "call_expression" {
		funcNode := rootNode.ChildByFieldName("function")
		if funcNode != nil {
			callExpr := string(g.NodeData[funcNode.StartByte():funcNode.EndByte()])
			if strings.Contains(callExpr, ".") {
				pkg := strings.Split(callExpr, ".")[0]

				// We now know which file/pkgs this imports from
				// Build the graph upon this knowledge, such that the imports being populated are meaningful
				// Return valid imports from this function on calculation, and the graph should build itself correctly
				if _, ok := g.PkgPaths[pkg]; ok {
					for _, z := range g.PkgPaths[pkg] {
						n := types.NodeImport{
							ImportPackage: z,
							FilePath:      g.NodePath,
						}

						if !slices.Contains(imports, n) {
							imports = append(imports, n)
						}
					}
				}
			}
		}
	}

	for i := uint(0); i < rootNode.ChildCount(); i++ {
		args.RootNode = rootNode.Child(i)
		imports = append(imports, g.resolveImportNodes(args)...)
		// for _, i := range imports {
		// if !slices.Contains(im, i) {
		// imports = append(imports, i)
		// }
		// }
	}
	g.Visit[rootNode] = false
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
