package repository

import (
	"fmt"
	"path/filepath"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/dag"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Repository struct {
	TargetDir                string
	Graph                    *dag.DAG
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
		Graph:                    dag.NewDAG(),
		RepositoryNodesAST:       make(map[string]*ast.AST),
	}
}

func (r *Repository) Run() error {
	err := filepath.WalkDir(r.TargetDir, r.populateRepositoryNodeImports)

	for _, v := range r.NodeImports {
		r.Graph.BuildGraphFromImports(v)
	}

	fmt.Println(r.Graph)
	fmt.Println(r.Graph.GetInNodesOfNode("os"))
	fmt.Println(r.Graph.GetOutNodesOfNode("os"))

	return err
}
