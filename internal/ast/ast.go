package ast

import "github.com/manosriram/wingman/internal/types"

type AST struct {
	Code     string
	Language types.Language
}

func NewAST(code string, language types.Language) *AST {
	return &AST{
		Code:     code,
		Language: language,
	}
}
