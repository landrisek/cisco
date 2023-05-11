package controller

import (
	"github.com/landrisek/cisco/src/repository"
)

// WalkGraph traverses the graph starting from the given node and returns a slice of all visited nodes.
// If the provided node is nil, it returns an empty slice.
func WalkGraph(node repository.GNode) []repository.GNode {
	if node == nil {
		return []repository.GNode{}
	}
	var nodes []repository.GNode
	findNode(node, &nodes)
	return nodes
}

// findNode recursively traverses the graph starting from the given node and appends each visited node to the nodes slice.
// It updates the nodes slice by reference using a pointer.
func findNode(node repository.GNode, nodes *[]repository.GNode) {
	*nodes = append(*nodes, node)
	for _, children := range node.GetChildren() {
		findNode(children, nodes)
	}
}
