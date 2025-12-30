package dag

import (
	"github.com/manosriram/wingman/internal/types"
)

type GraphNode struct {
	NodeValue string // interface?
}

func NewGraphNode(nodeValue string) GraphNode {
	return GraphNode{
		NodeValue: nodeValue,
	}
}

type DAG struct {
	Graph         map[string][]GraphNode
	GraphInDegree map[string]int
}

func NewDAG() *DAG {
	return &DAG{
		Graph:         make(map[string][]GraphNode),
		GraphInDegree: make(map[string]int),
	}
}

func (d *DAG) addEdge(src, dest GraphNode) {
	if _, ok := d.Graph[src.NodeValue]; !ok {
		d.Graph[src.NodeValue] = []GraphNode{}
	}
	d.Graph[src.NodeValue] = append(d.Graph[src.NodeValue], dest)
}

func (d *DAG) GetOutNodesOfNode(nodeKey string) []GraphNode {
	return d.Graph[nodeKey]
}

func (d *DAG) GetInNodesOfNode(nodeKey string) []string {
	var inNodes []string

	for k, v := range d.Graph {
		for _, imp := range v {
			if imp.NodeValue == nodeKey {
				inNodes = append(inNodes, k)
			}
		}
	}

	return inNodes
}

func (d *DAG) removeEdge(src, dest GraphNode) {}

func (d *DAG) BuildGraphFromImports(imports []types.NodeImport) {
	for _, i := range imports {
		d.addEdge(
			NewGraphNode(i.FilePath),
			NewGraphNode(i.ImportPath),
		)
	}
}
