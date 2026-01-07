package algorithm

import (
	"github.com/manosriram/wingman/internal/graph"
	"github.com/manosriram/wingman/internal/types"
)

/*
	for N iterations {
		PR(u) = (1-d)/N + d * (v = for each file that imports file u) PR(v)/L(v)

		d = damping factor 0.85
		L(v) = The number of imports inside file v
	}
*/

type PageRankAlgorithm struct {
	NodeScores map[string]float64
}

func NewPageRankAlgorithm() *PageRankAlgorithm {
	return &PageRankAlgorithm{
		NodeScores: make(map[string]float64),
	}
}

func (p *PageRankAlgorithm) GetScoreForNode(node string) float64 {
	return p.NodeScores[node]
}

func (p *PageRankAlgorithm) CalculateScore(graph *graph.Graph) {
	d := 0.85
	iters := 10
	N := len(graph.G)
	if N == 0 {
		return
	}
	n := float64(N)

	prev := make(map[string]float64, N)
	next := make(map[string]float64, N)

	for node := range graph.G {
		prev[node] = 1.0 / n
	}

	for range iters {
		base := (1.0 - d) / n
		for node := range graph.G {
			next[node] = base
		}

		var dangling float64
		for v := range graph.G {
			out := len(graph.G[v])
			if out == 0 {
				dangling += prev[v]
			}
		}

		for v, outs := range graph.G {
			outDegree := len(outs)
			if outDegree == 0 {
				continue
			}
			share := d * prev[v] / float64(outDegree)
			for _, u := range outs {
				next[u.NodeValue] += share
			}
		}

		if dangling != 0 {
			add := d * dangling / n
			for node := range graph.G {
				next[node] += add
			}
		}

		prev, next = next, prev
	}

	p.NodeScores = prev
}

func (p *PageRankAlgorithm) GetAlgorithmType() types.ContextAlgorithmType {
	return types.PAGERANK_CONTEXT_ALGORITHM
}
