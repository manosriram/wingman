package test

import (
	"math"
	"testing"

	"github.com/manosriram/wingman/internal/algorithm"
	"github.com/manosriram/wingman/internal/graph"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func sumScores(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

func TestPageRankAlgorithm_EmptyGraph_NoScores(t *testing.T) {
	g := graph.NewGraph()

	pr := algorithm.NewPageRankAlgorithm()
	pr.CalculateScore(g)

	if len(pr.NodeScores) != 0 {
		t.Fatalf("expected no scores for empty graph, got %d: %#v", len(pr.NodeScores), pr.NodeScores)
	}
}

func TestPageRankAlgorithm_SimpleChain_ProducesDistribution(t *testing.T) {
	g := graph.NewGraph()

	// Build a simple chain: A -> B -> C
	g.G["A"] = []graph.GraphNode{graph.NewGraphNode("B")}
	g.G["B"] = []graph.GraphNode{graph.NewGraphNode("C")}
	g.G["C"] = []graph.GraphNode{}

	pr := algorithm.NewPageRankAlgorithm()
	pr.CalculateScore(g)

	// Should have a score for every node in the graph map.
	if _, ok := pr.NodeScores["A"]; !ok {
		t.Fatalf("expected score for node A")
	}
	if _, ok := pr.NodeScores["B"]; !ok {
		t.Fatalf("expected score for node B")
	}
	if _, ok := pr.NodeScores["C"]; !ok {
		t.Fatalf("expected score for node C")
	}

	// Scores should be finite and non-negative.
	for k, v := range pr.NodeScores {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("expected finite score for %s, got %v", k, v)
		}
		if v < 0 {
			t.Fatalf("expected non-negative score for %s, got %v", k, v)
		}
	}

	// PageRank should form a probability distribution (approximately sums to 1).
	total := sumScores(pr.NodeScores)
	if !almostEqual(total, 1.0, 1e-6) {
		t.Fatalf("expected scores to sum to ~1.0, got %.10f (%#v)", total, pr.NodeScores)
	}

	// In a chain with a dangling end, C should typically have the highest rank.
	if !(pr.NodeScores["C"] > pr.NodeScores["B"] && pr.NodeScores["B"] > pr.NodeScores["A"]) {
		t.Fatalf("expected C > B > A, got A=%.6f B=%.6f C=%.6f",
			pr.NodeScores["A"], pr.NodeScores["B"], pr.NodeScores["C"])
	}
}

func TestPageRankAlgorithm_DanglingNode_StillSumsToOne(t *testing.T) {
	g := graph.NewGraph()

	// A -> B, and C is dangling and disconnected.
	g.G["A"] = []graph.GraphNode{graph.NewGraphNode("B")}
	g.G["B"] = []graph.GraphNode{}
	g.G["C"] = []graph.GraphNode{}

	pr := algorithm.NewPageRankAlgorithm()
	pr.CalculateScore(g)

	if len(pr.NodeScores) != 3 {
		t.Fatalf("expected 3 scores, got %d: %#v", len(pr.NodeScores), pr.NodeScores)
	}

	for k, v := range pr.NodeScores {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			t.Fatalf("expected finite score for %s, got %v", k, v)
		}
		if v < 0 {
			t.Fatalf("expected non-negative score for %s, got %v", k, v)
		}
	}

	total := sumScores(pr.NodeScores)
	if !almostEqual(total, 1.0, 1e-6) {
		t.Fatalf("expected scores to sum to ~1.0, got %.10f (%#v)", total, pr.NodeScores)
	}
}
