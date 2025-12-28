package test

import (
	"testing"

	"github.com/manosriram/wingman/internal/ast"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
	"github.com/stretchr/testify/assert"
)

const (
	TARGET_DIR     = "/Users/manosriram/go/src/go2java"
	TARGET_GO_NODE = "/Users/manosriram/go/src/go2java/ast.go"
)

func Test_GetNodeImports(t *testing.T) {
	treeSitterLanguageParser := utils.NewTreeSitterParserType()
	imports, err := ast.NewAST(TARGET_GO_NODE, types.GO, treeSitterLanguageParser).GetNodeImports()
	assert.Nil(t, err)
	assert.NotEmpty(t, imports)
}
