package program

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manosriram/wingman/internal/dag"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Program struct {
	TargetDir   string
	Parser      utils.TreeSitterParserType
	Graph       dag.DAG
	NodeImports map[string][]string // File vs Imports
}

func NewProgram(targetDir string) *Program {
	parser := utils.NewTreeSitterParserType()

	return &Program{
		TargetDir:   targetDir,
		NodeImports: make(map[string][]string),
		Parser:      parser,
	}
}

func (p *Program) getImports(path string, d fs.DirEntry, err error) error {
	// Check if it's a directory or a file
	if d.IsDir() {

		// Optional: Skip specific directories (e.g., .git or node_modules)
		if d.Name() == ".git" || d.Name() == ".aider" {
			// fmt.Printf("Directory: %s\n", path)
			return filepath.SkipDir
		}
	} else {
		// TODO: support different languages, remove this static check and use dynamic languages support
		switch utils.GetLanguage(path) {
		case types.GO:
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			tree := p.Parser.GetLanguageParser(types.GO).Parse(data, nil)
			defer tree.Close()

			rootNode := tree.RootNode()

			imports := []string{}

			// Iterate through child nodes of the root
			// TODO: Move this from here
			for i := 0; i < int(rootNode.ChildCount()); i++ {
				child := rootNode.Child(uint(i))

				// Check if the node is an import_declaration
				if child.Kind() == "import_declaration" {
					// Handle single import: import "fmt"
					importSpec := child.ChildByFieldName("spec")
					if importSpec != nil && importSpec.Kind() == "import_spec" {
						path := importSpec.ChildByFieldName("path")
						if path != nil {
							importPath := data[path.StartByte():path.EndByte()]
							imports = append(imports, string(importPath))
						}
					}

					// Handle grouped imports: import ( ... )
					for j := 0; j < int(child.ChildCount()); j++ {
						specList := child.Child(uint(j))
						if specList.Kind() == "import_spec_list" {
							for k := 0; k < int(specList.ChildCount()); k++ {
								spec := specList.Child(uint(k))
								if spec.Kind() == "import_spec" {
									path := spec.ChildByFieldName("path")
									if path != nil {
										importPath := data[path.StartByte():path.EndByte()]
										imports = append(imports, string(importPath))
									}
								}
							}
						}
					}
				}
			}

			p.NodeImports[path] = imports

		default:
			fmt.Println("To be implemented")

		}
	}

	return nil
}

func (p *Program) Run() error {
	err := filepath.WalkDir(p.TargetDir, p.getImports)
	return err
}
