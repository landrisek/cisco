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

// NewNode creates and returns a new instance of MyNode.
// The returned MyNode pointer can be used to build a graph structure.
func NewNode() *MyNode {
	return &MyNode{
		children: make([]GNode, 0),
	}
}

func (n MyNode) GetName() string {
	return n.name
}

// SetName sets the name of the MyNode to the given name.
// It enforces immutability by checking if the name is already set.
// If the name is already set, it returns nil indicating that the operation is not allowed.
// The returned *MyNode pointer is used to maintain method chaining.
// Using a pointer receiver allows modifying the value of the node in the context of this method.

// HINT: I decided here to go with immutable setter though theres should be definitely same broader
// discussion on this. First but not the least, the only other setter SetChildren allows
// mutability, which somebody can consider for heteregenous design. In my POV I would go to enforcing
// immutability where it is possible and reasonable and go with mutability where it is practical
// in terms of readability and performance. If I would be sure, this is not diverting from
// acceptance criteria (which I am not), I would propose this setters to be part of GNode interface.
func (n *MyNode) SetName(name string) *MyNode {
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
	n.children = children
	return n
}

// MarshalJSON is exposing json for rest api server purposes.
func (n *MyNode) MarshalJSON() ([]byte, error) {
	type node struct {
		Name     string   `json:"name"`
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

// GetSubTags will return fist occurence of tag.
// It does not expect tags with duplicite names in data structures.
func GetSubTags(ctx context.Context, node GNode, tag string) MyNode {
	result := make(chan MyNode, 1)
	active := int32(1)
	done := make(chan struct{})
	innerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool := NewPool(10)
	defer pool.Wait()
	go lookupChildrens(ctx, innerCtx, node, tag, result, done, pool, &active)

	select {
	case <-ctx.Done():
		return MyNode{}
	case <-innerCtx.Done():
		return MyNode{}
	case res := <-result:
		return res
	case <-done:
		return MyNode{}
	}
}

// lookupChildrens performs a recursive search for a specific tag within a node and its children,
// utilizing parallel processing for improved performance in the hot spot section.
// It sends matching nodes to the result channel and tracks the number of active goroutines using the active counter.
// The ctx context is used for cancellation and termination.
// Once all goroutines have completed, a signal is sent to the done channel.
// HINT: only channels for writing as input parameters, method is not draining them.
func lookupChildrens(ctx context.Context, innerCtx context.Context, node GNode, tag string, result chan<- MyNode, done chan<- struct{}, pool *Pool, active *int32) {
	defer func() {
		if atomic.AddInt32(active, -1) == 0 {
			done <- struct{}{}
		}
	}()

	if node.GetName() == tag {
		myNode, ok := node.(MyNode)
		if !ok {
			return
		}
		result <- myNode
		return
	}

	if len(node.GetChildren()) >= 10 {
		for _, child := range node.GetChildren() {
			select {
			case <-ctx.Done():
				return
			case <-innerCtx.Done():
				return
			default:
				atomic.AddInt32(active, 1)
				pool.Schedule(func() {
					lookupChildrens(ctx, innerCtx, child, tag, result, done, pool, active)
				})
			}
		}
	} else {
		for _, child := range node.GetChildren() {
			atomic.AddInt32(active, 1)
			lookupChildrens(ctx, innerCtx, child, tag, result, done, pool, active)
		}
	}
}

var tokenCache = map[string]bool{
	"XXX": true,
	"YYY": true,
}

// IsAuthenticated checks if a token is present in the tokenCache map,
// indicating that the token is valid and represents an authenticated user.
// It returns true if the token is found in the cache, otherwise false.
func IsAuthenticated(token string) bool {
	if _, ok := tokenCache[token]; ok {
		return true
	}
	return false
}

// GetValidToken retrieves a valid token from the tokenCache map.
// It returns the first token found in the cache, or an empty string if the cache is empty.
func GetValidToken() string {
	for token, _ := range tokenCache {
		return token
	}
	return ""
}
