package controller

import (
	"testing"

	"github.com/landrisek/cisco/src/repository"
)

func TestPaths(t *testing.T) {
	testCases := []testCase{
		{
			name:           "Test with no node",
			input:          nil,
			expectedLength: 0,
			expectedPaths:  [][]string{},
		},
		{
			name:           "Test with one node",
			input:          repository.NewNode().SetName("A"),
			expectedLength: 1,
			expectedPaths: [][]string{
				{"A"},
			},
		},
		{
			name: "Test with multiple nodes",
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
			expectedLength: 6,
			expectedPaths: [][]string{
				{"A", "B", "E"},
				{"A", "B", "F"},
				{"A", "C", "G"},
				{"A", "C", "H"},
				{"A", "C", "I"},
				{"A", "D", "J"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			paths := Paths(tc.input)
			if len(paths) != tc.expectedLength {
				t.Errorf("Expected %d paths, but got %d", tc.expectedLength, len(paths))
			}
			for i, path := range paths {
				expectedPaths := tc.expectedPaths[i]
				for j, node := range path {
					if node.GetName() != expectedPaths[j] {
						t.Errorf("Path %d is not equal to the expected value.\nExpected: %v\nActual: %v\n", i, expectedPaths[j], node.GetName())
					}
				}
			}
		})
	}
	// If test passes, old documenation is overwriten by new one keeping what actually is microservices doing
	generateHTMLDoc(testCases, "Paths")
}
