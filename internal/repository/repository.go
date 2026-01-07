package repository

import (
	"fmt"
	"sort"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Repository struct {
	TargetDir                string
	Graph                    *graph.Graph
	TreeSitterLanguageParser utils.TreeSitterParserType
	PkgPaths                 map[string][]string
	NodeImports              map[string][]types.NodeImport // Pkg vs Imports
	RepositoryNodesAST       map[string]*ast.AST
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
	}
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
		fmt.Println(v.Key, v.Value)
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
