package repository

import (
	"fmt"
	"path/filepath"

	"github.com/manosriram/wingman/internal/dag"
	"github.com/manosriram/wingman/internal/utils"
)

type Repository struct {
	TargetDir                string
	Graph                    *dag.DAG
	NodeImports              map[string][]string // File vs Imports
	TreeSitterLanguageParser utils.TreeSitterParserType
	// AST         *ast.AST -- map[string]*ast.AST?
}

func NewRepository(targetDir string) *Repository {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()

	return &Repository{
		TargetDir:                targetDir,
		NodeImports:              make(map[string][]string),
		TreeSitterLanguageParser: treeSitterLanguageParser,
	}
}

func (r *Repository) Run() error {
	err := filepath.WalkDir(r.TargetDir, r.populateRepositoryNodeImports)

	fmt.Println(r.NodeImports)
	return err
}
