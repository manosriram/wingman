package language

import (
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type LangStrategy interface {
	/*
		Return the list of imports of a given node (file)

		This interface is implemented for each language since each has its own way of imports
		Refer language/<language>.go for specific implementation of this strategy
	*/
	GetNodeImportList() ([]types.NodeImport, error)

	/*
		Internal method which parses the code repository and lists the imports
	*/
	resolveImportNodes(ResolveImportNodesArgs) []types.NodeImport
}

func GetStrategy(args StrategyArgs) LangStrategy {
	switch args.StrategyLanguage {
	case types.GOLANG:
		return NewGolangStrategy(args)
	}
	return NewDefaultStrategy(args)
}

type StrategyArgs struct {
	NodeData         []byte
	NodePath         string
	Parser           utils.TreeSitterParserType
	StrategyLanguage types.Language
}

type ResolveImportNodesArgs struct {
	RootNode          *tree_sitter.Node
	GolangModFileData string
}
