package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/landrisek/cisco/src/repository"
)

func UploadJson(filename string) (repository.GNode, error) { 
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var node interface{}
	err = json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}
	return convertNode(node), nil
}

// Acceptance criteria imply by using getter in interface GNode
// that fields ("class variables") should stay private
// on this assumption there is this workaround
// otherwise with exported names it will work with unmarshall out of the box
func convertNode(n interface{}) repository.MyNode {
    node := repository.MyNode{}
	// we reach the bottom on subsequent branch of graph
    switch n := n.(type) {
	// for node processing
    case map[string]interface{}:
        for k, v := range n {
            if k == "name" {
                if ok := node.SetName(v.(string)); ok == nil {
					Log(fmt.Errorf("Immutability on tag`s name was broken"), fmt.Sprintf("Trying to replace %s with %s", node.GetName(), v.(string)))
				}
            } else if k == "children" {
                children := v.([]interface{})
                for _, child := range children {
                    node.SetChildren(append(node.GetChildren(), convertNode(child)))
                }
            }
        }
    }

    return node
}

func WalkGraph(node repository.GNode) []repository.GNode {
	if node == nil {
		return []repository.GNode{}
	}
	// 1. create ouput slice of GNodes
	var nodes []repository.GNode
	// 2. call recursive function with 1st parameter given root, 2nd pointer to slice of nodes
	findNode(node, &nodes)
	// 6. after recursive function with no return value finished, return value of slice nodes
	return nodes
}

func findNode(node repository.GNode, nodes *[]repository.GNode) {
	// in recursive function:
	// 3. add root to slice of GNodes
	*nodes = append(*nodes, node)
	// 4. get children of root and iterate all elements
	for _, children := range node.GetChildren() {
		// 5. call recursive on each element of children slice
		findNode(children, nodes)
	}
}

func Paths(node repository.GNode) [][]repository.GNode {
	// 1. data structrues - create empty slice of slices with GNodes as "paths"
	var paths [][]repository.GNode
	// 2. handle nil node
	if nil == node {
		return paths
	}
	// 3. data structrues - create empty slice of GNodes as "path"
	var path []repository.GNode
	// 4. call recursive function with input parameters pointer node, pointer to path and paths
	findBottom(node, &path, &paths)
	// 10. return value slice of slices
	return paths
	// 11. check for edge cases / reasonability of return values
}

func findBottom(node repository.GNode, path *[]repository.GNode, paths *[][]repository.GNode) {
	// 5. we checked if node is not nil, in parent function, so add it to the path by dereferencing appending of slice
	*path = append(*path, node)
	childrens := node.GetChildren()
	// 6. if childrens of node is empty, we reached the bottom
	if len(childrens) == 0 {
		// 7. append defefrenced path to derefefernce paths
		*paths = append(*paths, append([]repository.GNode{}, *path...))
	} else {
		// 8. otherwise call recursive function on each child
		for _, child := range childrens {
			findBottom(child, path, paths)
		}
	}
	// 9. as path [A,B] was added, but next call will be on C, we need path looks like [A] 
	// where we add C so it will be [A,C]
	*path = (*path)[:len(*path)-1]
}

func RestAPI(node repository.GNode) {

}

func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg + ": %s", err)
	}
}