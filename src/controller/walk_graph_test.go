package controller

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/landrisek/cisco/src/repository"
)

/*
			   A
		/      |    \
	   B       C     D
	  / \     /|\    |
	  E  F   G H I   J
*/
func testWalkGraph(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := WalkGraph(tc.input)

			// check if returned nodes slice has the expected length
			if len(nodes) != tc.expectedLength {
				t.Errorf("Expected nodes length of %d, but got %d", tc.expectedLength, len(nodes))
			}
			// check if returned nodes slice has the expected elements
			for i, node := range nodes {
				if node.GetName() != tc.expectedContent[i] {
					t.Errorf("Expected node name %s, but got %s", tc.expectedContent[i], node.GetName())
				}
			}
		})
	}
	// If test passes, old documenation is overwriten by new one keeping what actually is microservices doing
	generateHTMLDoc(testCases, "WalkGraph")
}

func TestWalkGraph(t *testing.T) {
	testCases := []testCase{
		{
			name:            "Test acyclic graph with no node",
			input:           nil,
			expectedLength:  0,
			expectedContent: []string{},
		},
		{
			name:            "Test acyclic graph with one empty node",
			input:           repository.NewNode(),
			expectedLength:  1,
			expectedContent: []string{""},
		},
		{
			name:            "Test acyclic graph with one node with values",
			input:           repository.NewNode().SetName("A").SetChildren([]repository.GNode{}),
			expectedLength:  1,
			expectedContent: []string{"A"},
		},
		{
			name: "Test acyclic graph with more nodes",
			input: repository.NewNode().SetName("A").SetChildren([]repository.GNode{
				repository.NewNode().SetName("B").SetChildren([]repository.GNode{
					repository.NewNode().SetName("E"),
					repository.NewNode().SetName("F"),
				}),
				repository.NewNode().SetName("C").SetChildren([]repository.GNode{
					repository.NewNode().SetName("G"),
					repository.NewNode().SetName("H"),
					repository.NewNode().SetName("I"),
				}),
				repository.NewNode().SetName("D").SetChildren([]repository.GNode{
					repository.NewNode().SetName("J"),
				}),
			}),
			expectedLength:  10,
			expectedContent: []string{"A", "B", "E", "F", "C", "G", "H", "I", "D", "J"},
		},
		{
			name: "Test disconnected graph",
			input: repository.NewNode().
				SetName("H").
				SetChildren([]repository.GNode{
					repository.NewNode().
						SetName("E").
						SetChildren([]repository.GNode{
							repository.NewNode().SetName("B"),
							repository.NewNode().SetName("C"),
						}),
					repository.NewNode().SetName("F"),
				}),
			expectedLength:  5,
			expectedContent: []string{"H", "E", "B", "C", "F"},
		},
		generateRandomTestCase(10),
		// Add more test cases as needed
	}

	// Call the testWalkGraph function with the test cases
	testWalkGraph(t, testCases)
}

func generateRandomTestCase(maxNodes int) testCase {
	rand.Seed(42)

	nodes := make([]*repository.MyNode, maxNodes)

	root := repository.NewNode().SetName("A")
	nodes[0] = root

	numChildren := rand.Intn(maxNodes-1) + 1
	for i := 0; i < numChildren; i++ {
		node := repository.NewNode().SetName(string(byte('B' + i)))
		root.SetChildren(append(root.GetChildren(), node))
		nodes[i+1] = node
	}

	for i := numChildren + 1; i < maxNodes; i++ {
		node := repository.NewNode().SetName(string(byte('B' + i)))
		parent := nodes[rand.Intn(i-numChildren)]
		parent.SetChildren(append(root.GetChildren(), node))
		nodes[i] = node
	}

	expectedContent := make([]string, 0, maxNodes)
	for _, node := range nodes {
		expectedContent = append(expectedContent, node.GetName())
	}

	return testCase{
		name:            fmt.Sprintf("Random acyclic graph with %d nodes", maxNodes),
		input:           root,
		expectedLength:  maxNodes,
		expectedContent: expectedContent,
	}
}
