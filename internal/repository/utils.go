package repository

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

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

		var pkg string
		if utils.GetLanguage(path) == types.GOLANG {

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				pkg = scanner.Text()
				break
			}
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}

			if utils.GetLanguage(path) == types.GOLANG && len(strings.Split(pkg, " ")) > 1 {
				pkg = strings.Split(pkg, " ")[1]
			}
		} else {
			pkg = path
		}

		r.RepositoryNodesAST[pkg] = ast.NewAST(path, r.TreeSitterLanguageParser)
		imports, err := r.RepositoryNodesAST[pkg].GetNodeImports()
		if err != nil {
			return err
		}
		r.NodeImports[pkg] = imports
	}
	return nil
}
