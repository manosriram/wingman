package algorithm

import (
	"github.com/manosriram/wingman/internal/dag"
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

func (p *PageRankAlgorithm) CalculateScore(dag *dag.DAG) {
	d := 0.85
	iters := 50
	N := len(dag.Graph)
	if N == 0 {
		return
	}
	n := float64(N)

	// prev and next rank vectors
	prev := make(map[string]float64, N)
	next := make(map[string]float64, N)

	// init uniform
	for node := range dag.Graph {
		prev[node] = 1.0 / n
	}

	for range iters {
		// base teleportation
		base := (1.0 - d) / n
		for node := range dag.Graph {
			next[node] = base
		}

		// dangling mass (nodes with outDegree 0)
		var dangling float64
		for v := range dag.Graph {
			out := len(dag.Graph[v])
			if out == 0 {
				dangling += prev[v]
			}
		}

		// distribute rank along edges
		for v, outs := range dag.Graph {
			outDegree := len(outs)
			if outDegree == 0 {
				continue
			}
			share := d * prev[v] / float64(outDegree)
			for _, u := range outs {
				// u.NodeValue is the destination key
				next[u.NodeValue] += share
			}
		}

		// redistribute dangling mass uniformly
		if dangling != 0 {
			add := d * dangling / n
			for node := range dag.Graph {
				next[node] += add
			}
		}

		// swap prev/next
		prev, next = next, prev
	}

	p.NodeScores = prev
}

func (p *PageRankAlgorithm) GetAlgorithmType() types.ContextAlgorithmType {
	return types.PAGERANK_CONTEXT_ALGORITHM
}
