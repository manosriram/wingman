package ast

import (
	"errors"
	"log"
	"os"

	"github.com/manosriram/wingman/internal/algorithm"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/language"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

func NewAST(nodePath string, PkgPaths map[string][]string, parser utils.TreeSitterParserType) *AST {
	data, err := os.ReadFile(nodePath)
	if err != nil {
		log.Fatalf("Error initializing AST")
	}
	return &AST{
		NodeData:     data,
		NodePath:     nodePath,
		NodeLanguage: utils.GetLanguage(nodePath),
		Parser:       parser,
		PkgPaths:     PkgPaths,
		Algorithm:    algorithm.NewPageRankAlgorithm(),
	}
}

func (a *AST) GetNodeImports() ([]types.NodeImport, error) {
	return language.GetStrategy(language.StrategyArgs{
		NodeData:         a.NodeData,
		NodePath:         a.NodePath,
		Parser:           a.Parser,
		PkgPaths:         a.PkgPaths,
		StrategyLanguage: utils.GetLanguage(a.NodePath),
	}).GetNodeImportList()
}

func (a *AST) CalculateASTNodesScore(g *graph.Graph) error {
	if a.Algorithm.GetAlgorithmType() == types.PAGERANK_CONTEXT_ALGORITHM {
		a.Algorithm.CalculateScore(g)
	} else {
		return errors.New("Algorithm not implemented")
	}
	return nil
}
