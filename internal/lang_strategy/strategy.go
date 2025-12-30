package langstrategy

import "github.com/manosriram/wingman/internal/types"

type LangStrategy interface {
	GetNodeImports() ([]types.NodeImport, error)
}
