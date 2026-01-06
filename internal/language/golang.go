package language

import (
	"errors"
	"fmt"
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
	// modFileData := args.GolangModFileData

	if rootNode.Kind() == "call_expression" {
		funcNode := rootNode.ChildByFieldName("function")
		if funcNode != nil {
			fmt.Println("call_expression found:", string(g.NodeData[funcNode.StartByte():funcNode.EndByte()]))
		}
	}

	imports := []types.NodeImport{}
	for i := uint(0); i < rootNode.ChildCount(); i++ {
		child := rootNode.Child(i)
		args.RootNode = child
		g.resolveImportNodes(args)
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
