package repository

import (
	"encoding/json"
	"fmt"
)

type GNode interface {
	GetName() string
	GetChildren() []GNode
}

type MyNode struct {
	name     string  `json:"name"`
	children []GNode `json:"children"`
}

func NewNode() *MyNode {
	return &MyNode{
		children: make([]GNode, 0),
	}
}

func (n MyNode) GetName() string {
	return n.name
}

// this setter would suits GNode interface
// effectively dividing "data-structure" functionality
// from application logic
// but as it was not in acceptance criteria it is not
// pointer is used to avoid assigning value to value
// though this have impact on heap stack
func (n *MyNode) SetName(name string) *MyNode {
	// set name should be immutable
	// because we want to enforce immutability
	// where it is possible and reasonable
	if n.name != "" {
		return nil
	}
	n.name = name
	return n
}

func (n MyNode) GetChildren() []GNode {
	return n.children
}

func (n *MyNode) SetChildren(children []GNode) *MyNode {
	// we allow mutability on childrens
	// as we expect appending values
	n.children = children
	return n
}

// expose for json
func (n *MyNode) MarshalJSON() ([]byte, error) {
	type node struct {
		Name     string  `json:"name"`
		Children []MyNode `json:"children"`
	}

	var children []MyNode
	for _, child := range n.children {
		myNode, ok := child.(MyNode)
		if !ok {
			return nil, fmt.Errorf("child is not of type MyNode")
		}
		children = append(children, myNode)
	}
	
	return json.Marshal(node{
		Name:     n.name,
		Children: children,
	})
}

// dummy
func GetValidToken() string {
	return "YYY"
}