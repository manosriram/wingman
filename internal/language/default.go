package language

import (
	"github.com/manosriram/wingman/internal/types"
)

type DefaultStrategy struct{}

func NewDefaultStrategy(args StrategyArgs) *DefaultStrategy {
	return &DefaultStrategy{}
}

func (d *DefaultStrategy) resolveImportNodes(args ResolveImportNodesArgs) []types.NodeImport {
	return []types.NodeImport{}
}

func (d *DefaultStrategy) GetNodeImportList() ([]types.NodeImport, error) {
	return d.resolveImportNodes(ResolveImportNodesArgs{}), nil
}
