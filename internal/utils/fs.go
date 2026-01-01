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
	} else if strings.HasSuffix(path, ".py") {
		return types.PYTHON
	} else if strings.HasSuffix(path, ".js") {
		return types.JAVASCRIPT
	}
	return types.UNKNOWN
}
