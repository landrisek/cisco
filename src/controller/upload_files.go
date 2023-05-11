package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/landrisek/cisco/src/repository"
)

// UploadJson reads the JSON data from the specified file and unmarshals it into a repository.GNode.
// It returns the root node of the JSON structure and an error if any occurred during the process.
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

// HINT: Acceptance criteria imply by using getter in interface GNode that fields ("class variables") should stay private.
// On this assumption there is this workaround, otherwise with exported names it will work with unmarshall out of the box.
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
