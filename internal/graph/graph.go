package graph

import (
	"github.com/manosriram/wingman/internal/types"
)

type GraphNode struct {
	NodeValue string
}

func NewGraphNode(nodeValue string) GraphNode {
	return GraphNode{
		NodeValue: nodeValue,
	}
}

type Graph struct {
	G map[string][]GraphNode
}

func NewGraph() *Graph {
	return &Graph{
		G: make(map[string][]GraphNode),
	}
}

func (d *Graph) addEdge(src, dest GraphNode) {
	if _, ok := d.G[src.NodeValue]; !ok {
		d.G[src.NodeValue] = []GraphNode{}
	}
	if _, ok := d.G[dest.NodeValue]; !ok {
		d.G[dest.NodeValue] = []GraphNode{}
	}
	d.G[src.NodeValue] = append(d.G[src.NodeValue], dest)
}

func (d *Graph) GetOutNodesOfNode(nodeKey string) []GraphNode {
	return d.G[nodeKey]
}

func (d *Graph) GetInNodesOfNode(nodeKey string) []string {
	var inNodes []string

	for k, v := range d.G {
		for _, imp := range v {
			if imp.NodeValue == nodeKey {
				inNodes = append(inNodes, k)
			}
		}
	}

	return inNodes
}

func (d *Graph) BuildGraphFromImports(imports []types.NodeImport) {
	for _, i := range imports {
		d.addEdge(
			NewGraphNode(i.ImportPackage),
			NewGraphNode(i.FilePath),
		)
	}
}
