package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/landrisek/cisco/src/repository"
)

func UploadJson(filename string) (repository.GNode, error) { 
	// Read the JSON data from file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a slice of Nodes
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
	return nil
}

func RestAPI(node repository.GNode) {

}

func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg + ": %s", err)
	}
}