package controller

import (
	"github.com/landrisek/cisco/src/repository"
)

// Paths returns a slice of paths in the graph starting from the given node. Each path is represented as a slice of nodes.
// If the node is nil, an empty slice is returned.
func Paths(node repository.GNode) [][]repository.GNode {
	var paths [][]repository.GNode
	if nil == node {
		return paths
	}
	var path []repository.GNode
	findBottom(node, &path, &paths)
	return paths
}

// findBottom is a recursive function that finds all paths in the graph starting from the given node and appends them to the paths slice.
// It keeps track of the current path in the 'path' parameter and appends the completed paths to the 'paths' parameter.
// If the node has no children, the current path is considered complete and added to the 'paths' slice.
// After processing a child node, the function removes it from the current path to prepare for the next iteration.
func findBottom(node repository.GNode, path *[]repository.GNode, paths *[][]repository.GNode) {
	*path = append(*path, node)
	childrens := node.GetChildren()
	if len(childrens) == 0 {
		*paths = append(*paths, append([]repository.GNode{}, *path...))
	} else {
		for _, child := range childrens {
			findBottom(child, path, paths)
		}
	}
	// HINT: as path [A,B] was added, but next call will be on C, we need path looks like [A] where we add C so it will be [A,C]
	*path = (*path)[:len(*path)-1]
}
