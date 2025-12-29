package algorithm

import (
	"github.com/manosriram/wingman/internal/types"
)

type PageRankAlgorithm struct {
	NodeScores map[string]int64
}

func NewPageRankAlgorithm() *PageRankAlgorithm {
	return &PageRankAlgorithm{}
}

func (p *PageRankAlgorithm) CalculateScore() {
}

func (p *PageRankAlgorithm) GetAlgorithmType() types.ContextAlgorithmType {
	return types.PAGERANK_CONTEXT_ALGORITHM
}
