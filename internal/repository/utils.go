package repository

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

func (r *Repository) populateRepositoryNodeImports(path string, d fs.DirEntry, err error) error {
	if d.IsDir() {
		if d.Name() == ".git" || d.Name() == ".aider" {
			return filepath.SkipDir
		}
	} else {
		// TODO: support different languages, remove this static check and use dynamic languages support
		switch utils.GetLanguage(path) {
		case types.GO:
			imports, err := ast.NewAST(path, types.GO, r.TreeSitterLanguageParser).GetNodeImports()
			if err != nil {
				return err
			}

			r.NodeImports[path] = imports
		default:
			fmt.Println("To be implemented")

		}
	}

	return nil
}
