package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/manosriram/wingman/internal/types"
)

func WalkDir(targetDir string, f func(path string, d fs.DirEntry, err error)) error {
	return nil
}

func GetLanguage(path string) types.Language {
	switch filepath.Ext(path) {
	case ".go":
		return types.GOLANG
	case ".py":
		return types.PYTHON
	case ".js":
		return types.JAVASCRIPT
	default:
		return types.UNKNOWN
	}
}

func FindGoModPath(startFilePath string) (string, error) {
	dir := startFilePath
	fi, err := os.Stat(startFilePath)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		dir = filepath.Dir(startFilePath)
	}

	for {
		candidate := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir { // root
			break
		}
		dir = parent
	}

	return "", errors.New("go.mod not found in any parent directory")
}

func ReadGoModFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
