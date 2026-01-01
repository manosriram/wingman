package test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/manosriram/wingman/internal/language"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

const (
	EXAMPLE_GO_FILE_CONTENT = `package main
		import (
			"fmt"
			"%s/internal/foo"
			"%s/pkg/bar"
		)

		func main() {
			fmt.Println("hi")
		}
`
)

func writeFile(t *testing.T, path, contents string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func sortImports(imps []types.NodeImport) {
	sort.Slice(imps, func(i, j int) bool {
		if imps[i].FilePath == imps[j].FilePath {
			return imps[i].ImportPackage < imps[j].ImportPackage
		}
		return imps[i].FilePath < imps[j].FilePath
	})
}

func TestDefaultStrategy_ReturnsEmpty(t *testing.T) {
	s := language.NewDefaultStrategy(language.StrategyArgs{})
	imps, err := s.GetNodeImportList()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(imps) != 0 {
		t.Fatalf("expected 0 imports, got %d: %#v", len(imps), imps)
	}
}

func TestGolangStrategy_GetNodeImportList_FindsInternalImports(t *testing.T) {
	tmp := t.TempDir()

	// Create a minimal Go module.
	modName := "example.com/myrepo"
	writeFile(t, filepath.Join(tmp, "go.mod"), "module "+modName+"\n\ngo 1.22\n")

	// Create a Go file with both internal and external imports.
	// The strategy should only return imports that start with the module name.
	mainPath := filepath.Join(tmp, "cmd", "app", "main.go")
	writeFile(t, mainPath, fmt.Sprintf(EXAMPLE_GO_FILE_CONTENT, modName, modName))

	data, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}

	parser := utils.NewTreeSitterParserType()

	s := language.NewGolangStrategy(language.StrategyArgs{
		NodeData:         data,
		NodePath:         mainPath,
		Parser:           parser,
		StrategyLanguage: types.GOLANG,
	})

	imps, err := s.GetNodeImportList()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// NOTE: The current implementation extracts the second path segment after the module name:
	// e.g. example.com/myrepo/internal/foo -> "foo"
	//      example.com/myrepo/pkg/bar      -> "bar"
	// This test asserts the current behavior.
	want := []types.NodeImport{
		{ImportPackage: "bar", FilePath: mainPath},
		{ImportPackage: "foo", FilePath: mainPath},
	}

	sortImports(imps)
	sortImports(want)

	if len(imps) != len(want) {
		t.Fatalf("expected %d imports, got %d: %#v", len(want), len(imps), imps)
	}

	for i := range want {
		if imps[i] != want[i] {
			t.Fatalf("import mismatch at %d: want %#v, got %#v", i, want[i], imps[i])
		}
	}
}

func TestGetStrategy_ReturnsGolangStrategy(t *testing.T) {
	parser := utils.NewTreeSitterParserType()
	s := language.GetStrategy(language.StrategyArgs{
		NodeData:         []byte("package main\n"),
		NodePath:         "main.go",
		Parser:           parser,
		StrategyLanguage: types.GOLANG,
	})

	if _, ok := s.(*language.GolangStrategy); !ok {
		t.Fatalf("expected *language.GolangStrategy, got %T", s)
	}
}

func TestGetStrategy_ReturnsDefaultStrategyForUnknown(t *testing.T) {
	parser := utils.NewTreeSitterParserType()
	s := language.GetStrategy(language.StrategyArgs{
		NodeData:         []byte(""),
		NodePath:         "file.unknown",
		Parser:           parser,
		StrategyLanguage: types.UNKNOWN,
	})

	if _, ok := s.(*language.DefaultStrategy); !ok {
		t.Fatalf("expected *language.DefaultStrategy, got %T", s)
	}
}
