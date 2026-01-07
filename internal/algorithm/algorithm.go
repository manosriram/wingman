package algorithm

import (
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
)

type ContextAlgorithm interface {
	CalculateScore(*graph.Graph)
	GetScoreForNode(string) float64
	GetAlgorithmType() types.ContextAlgorithmType
}
