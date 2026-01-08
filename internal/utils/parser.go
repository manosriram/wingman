package utils

import (
	"github.com/manosriram/wingman/internal/types"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

type TreeSitterParserType struct {
	Language types.Language
	Parsers  map[types.Language]*tree_sitter.Parser
}

func NewTreeSitterParserType() TreeSitterParserType {
	parsers := make(map[types.Language]*tree_sitter.Parser)

	goParser := tree_sitter.NewParser()
	pythonParser := tree_sitter.NewParser()
	javascriptParser := tree_sitter.NewParser()

	golangLanguage := tree_sitter.NewLanguage(tree_sitter_go.Language())
	javascriptLanguage := tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	pythonLanguage := tree_sitter.NewLanguage(tree_sitter_python.Language())

	javascriptParser.SetLanguage(javascriptLanguage)
	goParser.SetLanguage(golangLanguage)
	pythonParser.SetLanguage(pythonLanguage)

	parsers[types.GOLANG] = goParser
	parsers[types.JAVASCRIPT] = javascriptParser
	parsers[types.PYTHON] = pythonParser

	return TreeSitterParserType{
		Parsers: parsers,
	}
}

func (p TreeSitterParserType) GetLanguageParser(language types.Language) *tree_sitter.Parser {
	switch language {
	case types.GOLANG:
		return p.Parsers[types.GOLANG]
	case types.JAVASCRIPT:
		return p.Parsers[types.JAVASCRIPT]
	case types.PYTHON:
		return p.Parsers[types.PYTHON]
	}
	return nil
}

func (p TreeSitterParserType) Close() {}
