package repository

import (
	"os"
	"sort"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/llm"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Repository struct {
	TargetDir                string
	Graph                    *graph.Graph
	TreeSitterLanguageParser utils.TreeSitterParserType
	PkgPaths                 map[string][]string
	NodeImports              map[string][]types.NodeImport // Pkg vs Imports
	RepositoryNodesAST       map[string]*ast.AST
	Signatures               map[string][]string
	AddedFiles               map[string]string
}

type KeyValue struct {
	Key   string
	Value float64
}

func NewRepository(targetDir string) *Repository {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()

	return &Repository{
		TargetDir:                targetDir,
		TreeSitterLanguageParser: treeSitterLanguageParser,
		Graph:                    graph.NewGraph(),
		NodeImports:              make(map[string][]types.NodeImport),
		RepositoryNodesAST:       make(map[string]*ast.AST),
		PkgPaths:                 make(map[string][]string),
		Signatures:               make(map[string][]string),
		AddedFiles:               make(map[string]string),
	}
}

func getFunctionInfo(node *tree_sitter.Node, source []byte) (name string, params string) {
	if nameNode := node.ChildByFieldName("name"); nameNode != nil {
		name = string(source[nameNode.StartByte():nameNode.EndByte()])
	}
	if paramsNode := node.ChildByFieldName("parameters"); paramsNode != nil {
		params = string(source[paramsNode.StartByte():paramsNode.EndByte()])
	}
	return name, params
}

func getNodeSignatures(node *tree_sitter.Node, source []byte) []string {
	var signatures []string

	if node.Kind() == "function_declaration" || node.Kind() == "method_declaration" {
		// Extract the signature text (from start to the opening brace)
		// Or extract individual parts: name, parameters, result

		fnName, fnParams := getFunctionInfo(node, source)

		signatures = append(signatures, fnName+fnParams)
	}

	for i := uint(0); i < node.ChildCount(); i++ {
		signatures = append(signatures, getNodeSignatures(node.Child(i), source)...)
	}
	return signatures
}

func (r *Repository) GetNodeSignatures(path string) []string {
	p := r.TreeSitterLanguageParser.Parsers[types.GOLANG]

	d, _ := os.ReadFile(path)

	tree := p.Parse(d, nil)
	root := tree.RootNode()

	return getNodeSignatures(root, d)
}

func (r *Repository) Run() error {
	if err := r.walkDirAndPopulateRepositoryPkgPaths(); err != nil {
		return err
	}
	if err := r.walkDirAndPopulateNodeImports(); err != nil {
		return err
	}

	for _, v := range r.NodeImports {
		r.Graph.BuildGraphFromImports(v)
	}

	for k := range r.NodeImports {
		err := r.RepositoryNodesAST[k].CalculateASTNodesScore(r.Graph)
		if err != nil {
			return err
		}
	}

	// Sort scores descending
	nodeVsScores := make(map[string]float64)
	for k := range r.NodeImports {
		nodeVsScores[k] = r.RepositoryNodesAST[k].Algorithm.GetScoreForNode(k)
		// fmt.Printf("path = %s, score = %f\n", k, r.RepositoryNodesAST[k].Algorithm.GetScoreForNode(k))
	}

	var sorted []KeyValue
	for k, v := range nodeVsScores {
		sorted = append(sorted, KeyValue{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	for _, v := range sorted {
		r.Signatures[v.Key] = r.GetNodeSignatures(v.Key)
	}

	/* TODO
	We now have scores and the file paths.

	1. Create Repo Map by:
		a. Reading the path, and appending the signatures to the final LLM prompt

	2. Loop (1) until
		a. All paths are done
		b. Token limit is exhausted for the selected specific LLM (default 500K tokens)
	*/

	return nil
}

func (r *Repository) AddFile(path string) error {
	d, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	r.AddedFiles[path] = string(d)
	return nil
}

func (r *Repository) AddFiles(paths []string) error {
	for _, path := range paths {
		d, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		r.AddedFiles[path] = string(d)
	}
	return nil
}

func (r *Repository) CreateMasterPrompt(input string) string {
	return llm.CreateMasterPrompt(r.Signatures, r.AddedFiles, input)
}
