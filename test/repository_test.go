package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/manosriram/wingman/internal/repository"
)

func writeFileRepo(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func TestRepository_Run_BuildsGraphAndASTs(t *testing.T) {
	tmp := t.TempDir()

	// Minimal module so internal/language GolangStrategy can resolve internal imports.
	modName := "example.com/myrepo"
	goModFilePath := filepath.Join(tmp, "go.mod")
	writeFileRepo(t, goModFilePath, "module "+modName+"\n\ngo 1.22\n")

	// A small repo with internal imports.
	// Note: repository.populateRepositoryNodeImports uses the first line "package X" as the key.
	aFilePath := filepath.Join(tmp, "a.go")
	writeFileRepo(t, aFilePath, `package a

import (
	"`+modName+`/internal/foo"
)

func A() {}
`)

	bFilePath := filepath.Join(tmp, "b.go")
	writeFileRepo(t, bFilePath, `package b

import (
	"`+modName+`/internal/foo"
	"`+modName+`/internal/bar"
)

func B() {}
`)

	// Add a file under .git that should be skipped.
	writeFileRepo(t, filepath.Join(tmp, ".git", "ignored.go"), `package ignored

func Ignored() {}
`)

	r := repository.NewRepository(tmp)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error: %v", err)
	}

	// Basic sanity: ASTs created for packages a and b.
	if _, ok := r.RepositoryNodesAST["a"]; !ok {
		t.Fatalf("expected RepositoryNodesAST to contain key %q", "a")
	}
	if _, ok := r.RepositoryNodesAST["b"]; !ok {
		t.Fatalf("expected RepositoryNodesAST to contain key %q", "b")
	}

	// Ensure .git was skipped.
	if _, ok := r.RepositoryNodesAST["ignored"]; ok {
		t.Fatalf("did not expect RepositoryNodesAST to contain package from .git directory")
	}

	// NodeImports should be populated for a and b.
	if imps, ok := r.NodeImports["a"]; !ok {
		t.Fatalf("expected NodeImports to contain key %q", "a")
	} else if len(imps) == 0 {
		t.Fatalf("expected NodeImports[%q] to be non-empty", "a")
	}

	if imps, ok := r.NodeImports["b"]; !ok {
		t.Fatalf("expected NodeImports to contain key %q", "b")
	} else if len(imps) == 0 {
		t.Fatalf("expected NodeImports[%q] to be non-empty", "b")
	}

	// Graph should have nodes/edges after BuildGraphFromImports.
	if r.Graph == nil {
		t.Fatalf("expected Graph to be initialized")
	}
	if len(r.Graph.G) == 0 {
		t.Fatalf("expected Graph.G to be non-empty after Run()")
	}

	// After Run(), PageRank should have produced scores for graph nodes.
	// Repository.Run prints scores using NodeScores[k] where k is a package key ("a", "b").
	// We assert those keys exist in the score map.
	if r.RepositoryNodesAST["a"].Algorithm == nil {
		t.Fatalf("expected Algorithm to be initialized for AST %q", "a")
	}
	if _, ok := r.RepositoryNodesAST["a"].Algorithm.NodeScores[aFilePath]; !ok {
		t.Fatalf("expected NodeScores to contain key %q", "a")
	}
	if _, ok := r.RepositoryNodesAST["b"].Algorithm.NodeScores[bFilePath]; !ok {
		t.Fatalf("expected NodeScores to contain key %q", "b")
	}
}

func TestRepository_Run_EmptyDir_NoError(t *testing.T) {
	tmp := t.TempDir()

	// Even with no files, Run should not error.
	r := repository.NewRepository(tmp)
	if err := r.Run(); err != nil {
		t.Fatalf("Run() error on empty dir: %v", err)
	}

	if len(r.NodeImports) != 0 {
		t.Fatalf("expected NodeImports to be empty, got %d", len(r.NodeImports))
	}
	if len(r.RepositoryNodesAST) != 0 {
		t.Fatalf("expected RepositoryNodesAST to be empty, got %d", len(r.RepositoryNodesAST))
	}
	// Graph is initialized but should remain empty.
	if r.Graph == nil {
		t.Fatalf("expected Graph to be initialized")
	}
	if len(r.Graph.G) != 0 {
		t.Fatalf("expected Graph.G to be empty, got %d", len(r.Graph.G))
	}
}
