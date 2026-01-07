package ast

import (
	"github.com/manosriram/wingman/internal/algorithm"
	"github.com/manosriram/wingman/internal/language"
	"github.com/manosriram/wingman/internal/types"
	"github.com/manosriram/wingman/internal/utils"
)

type Signature struct {
}

type AST struct {
	Parser           utils.TreeSitterParserType
	NodePath         string
	NodeData         []byte
	NodeLanguage     types.Language
	PkgPaths         map[string][]string
	Algorithm        *algorithm.PageRankAlgorithm
	LanguageStrategy *language.LangStrategy
}
