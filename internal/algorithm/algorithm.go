package algorithm

import (
	"github.com/manosriram/wingman/internal/dag"
	"github.com/manosriram/wingman/internal/types"
)

type ContextAlgorithm interface {
	CalculateScore(*dag.DAG)
	GetAlgorithmType() types.ContextAlgorithmType
}
