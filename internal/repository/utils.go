package repository

import (
	"io/fs"
	"path/filepath"

	"github.com/manosriram/wingman/internal/ast"
)

func (r *Repository) populateRepositoryNodeImports(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		if d.Name() == ".git" || d.Name() == ".aider" {
			return filepath.SkipDir
		}
	} else {
		r.RepositoryNodesAST[path] = ast.NewAST(path, r.TreeSitterLanguageParser)
		imports, err := r.RepositoryNodesAST[path].GetNodeImports()
		if err != nil {
			return err
		}
		r.NodeImports[path] = imports

		r.RepositoryNodesAST[path].CalculateASTNodesScore(r.Graph)
	}
	return nil
}
