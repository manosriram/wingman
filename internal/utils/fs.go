package utils

import (
	"io/fs"
	"strings"

	"github.com/manosriram/wingman/internal/types"
)

func WalkDir(targetDir string, f func(path string, d fs.DirEntry, err error)) error {
	return nil
}

func GetLanguage(path string) types.Language {
	if strings.HasSuffix(path, ".go") {
		return types.GOLANG
	} else if strings.Contains(path, ".py") {
		return types.PYTHON
	} else if strings.Contains(path, ".js") {
		return types.JAVASCRIPT
	}
	return types.UNKNOWN
}
