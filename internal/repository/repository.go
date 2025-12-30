package repository

import (
	"path/filepath"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Repository struct {
	TargetDir                string
	Graph                    *graph.Graph
	NodeImports              map[string][]types.NodeImport // File vs Imports
	TreeSitterLanguageParser utils.TreeSitterParserType
	RepositoryNodesAST       map[string]*ast.AST
}

func NewRepository(targetDir string) *Repository {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()

	return &Repository{
		TargetDir:                targetDir,
		NodeImports:              make(map[string][]types.NodeImport),
		TreeSitterLanguageParser: treeSitterLanguageParser,
		Graph:                    graph.NewGraph(),
		RepositoryNodesAST:       make(map[string]*ast.AST),
	}
}

func (r *Repository) Run() error {
	err := filepath.WalkDir(r.TargetDir, r.populateRepositoryNodeImports)

	for _, v := range r.NodeImports {
		r.Graph.BuildGraphFromImports(v)
	}
	for k := range r.NodeImports {
		r.RepositoryNodesAST[k].CalculateASTNodesScore(r.Graph)
	}
	// for k := range r.NodeImports {
	// fmt.Printf("path = %s, score = %f\n", k, r.RepositoryNodesAST[k].Algorithm.NodeScores[k])
	// }

	return err
}
