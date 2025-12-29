package algorithm

import (
	"github.com/manosriram/wingman/internal/types"
)

type ContextAlgorithm interface {
	CalculateScore()
	GetAlgorithmType() types.ContextAlgorithmType
}
