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

func (r *Repository) walkDirAndPopulateRepositoryPkgPaths() error {
	return filepath.WalkDir(r.TargetDir, r.populateRepositoryPkgPaths)
}

func (r *Repository) walkDirAndPopulateNodeImports() error {
	return filepath.WalkDir(r.TargetDir, r.populateRepositoryNodeImports)
}

// func (g *GolangStrategy) resolveImportNodes(args ResolveImportNodesArgs) []types.NodeImport {

// rootNode := args.RootNode
// // modFileData := args.GolangModFileData

// if rootNode.Kind() == "call_expression" {
// funcNode := rootNode.ChildByFieldName("function")
// if funcNode != nil {
// fmt.Println("call_expression found:", string(g.NodeData[funcNode.StartByte():funcNode.EndByte()]))
// }
// }

// imports := []types.NodeImport{}
// for i := uint(0); i < rootNode.ChildCount(); i++ {
// child := rootNode.Child(i)
// args.RootNode = child
// g.resolveImportNodes(args)
// }
// return imports
// }

// func (r *Repository) getSignaturesFromPath(path string) []string {

// return []string{}
// }

func (r *Repository) populateRepositoryPkgPaths(path string, d fs.DirEntry, err error) error {
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

		}
		// } else {
		// pkg = path
		// }

		if pkg != "" {
			if utils.GetLanguage(path) == types.GOLANG && len(strings.Split(pkg, " ")) > 1 {
				pkg = strings.Split(pkg, " ")[1]
			}

			r.PkgPaths[pkg] = append(r.PkgPaths[pkg], path)

		}
	}
	return nil
}

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

		}
		// } else {
		// pkg = path
		// }

		if pkg != "" {
			if utils.GetLanguage(path) == types.GOLANG && len(strings.Split(pkg, " ")) > 1 {
				pkg = strings.Split(pkg, " ")[1]
			}

			r.RepositoryNodesAST[path] = ast.NewAST(path, r.PkgPaths, r.TreeSitterLanguageParser)
			imports, err := r.RepositoryNodesAST[path].GetNodeImports()
			if err != nil {
				return err
			}
			r.NodeImports[path] = imports
		}
	}
	return nil
}
