package repository

import (
	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Repository struct {
	TargetDir                string
	Graph                    *graph.Graph
	TreeSitterLanguageParser utils.TreeSitterParserType
	NodeImports              map[string][]types.NodeImport // Pkg vs Imports
	RepositoryNodesAST       map[string]*ast.AST
}

func NewRepository(targetDir string) *Repository {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()

	return &Repository{
		TargetDir:                targetDir,
		TreeSitterLanguageParser: treeSitterLanguageParser,
		Graph:                    graph.NewGraph(),
		NodeImports:              make(map[string][]types.NodeImport),
		RepositoryNodesAST:       make(map[string]*ast.AST),
	}
}

func (r *Repository) Run() error {
	err := r.walkDirAndPopulateNodeImports()

	for _, v := range r.NodeImports {
		r.Graph.BuildGraphFromImports(v)
	}
	for k := range r.NodeImports {
		err = r.RepositoryNodesAST[k].CalculateASTNodesScore(r.Graph)
		if err != nil {
			return err
		}
	}
	// for k := range r.NodeImports {
	// fmt.Printf("path = %s, score = %f\n", k, r.RepositoryNodesAST[k].Algorithm.NodeScores[k])
	// }

	return err
}
