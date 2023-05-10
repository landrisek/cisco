package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
)

type GNode interface {
	GetName() string
	GetChildren() []GNode
}

type MyNode struct {
	name     string  `json:"name"`
	children []GNode `json:"children"`
}

func (n MyNode) GetName() string {
	return n.name
}

func (n MyNode) GetChildren() []GNode {
	return n.children
}