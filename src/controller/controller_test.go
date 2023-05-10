package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/landrisek/cisco/src/repository"
)

type testCase struct {
	name            string
	input           repository.GNode
	expectedLength  int
	expectedContent []string
	expectedPaths   [][]string
}

/*		   A
	/      |    \
   B       C     D
  / \     /|\    | 
  E  F   G H I   J */
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
			expectedLength: 5,
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

func TestPaths(t *testing.T) {
	testCases := []testCase{
		{
			name: "Test with no node",
			input: nil,
			expectedLength: 0,
			expectedPaths: [][]string{},
		},
		{
			name: "Test with one node",
			input: repository.NewNode().SetName("A"),
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

func generateHTMLDoc(testCases []testCase, name string) {	
	// Create a new HTML file
	file, err := os.Create("../../" + url.QueryEscape(name) + "-doc.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// Write the HTML header and body tags
	fmt.Fprintln(file, "<html>")
	fmt.Fprintln(file, "<head><title>Documentation of " + name + " function</title></head>")
	fmt.Fprintln(file, "<body>")

	// Iterate through all test cases and write the input and expected output to the HTML file
	for _, tc := range testCases {
		// Write the test case name to the HTML file
		fmt.Fprintf(file, "<h2>%s</h2>", tc.name)

		// Write the input node to the HTML file
		fmt.Fprintln(file, "<p><b>Input:</b></p>")
		writeNodeToHTML(file, tc.input)
		// Write the expected output to the HTML file
		fmt.Fprintf(file, "<p><b>Expected output:</b> %d elements</p>", tc.expectedLength)
		fmt.Fprintln(file, "<ul>")
		for _, name := range tc.expectedContent {
			fmt.Fprintf(file, "<li>%s</li>", name)
		}
		for _, name := range tc.expectedPaths {
			fmt.Fprintf(file, "<li>%s</li>", name)
		}
		fmt.Fprintln(file, "</ul>")

		// Write a horizontal rule to separate the test cases
		fmt.Fprintln(file, "<hr>")
	}

	// Write the HTML closing tags
	fmt.Fprintln(file, "</body>")
	fmt.Fprintln(file, "</html>")
}

func writeNodeToHTML(w *os.File, node repository.GNode) {
	if node == nil {
		return
	}
	fmt.Fprintf(w, "<li>%s", node.GetName())
	children := node.GetChildren()
	if len(children) > 0 {
		fmt.Fprint(w, "<ul>")
		for _, child := range children {
			writeNodeToHTML(w, child)
		}
		fmt.Fprint(w, "</ul>")
	}
	fmt.Fprint(w, "</li>")
}

func TestRestAPI(t *testing.T) {
	// dereferencing for server structure
    node := *repository.NewNode().SetName("root").SetChildren([]repository.GNode{
        *repository.NewNode().SetName("child1").SetChildren([]repository.GNode{
            *repository.NewNode().SetName("grandchild1"),
            *repository.NewNode().SetName("grandchild2"),
        }),
        *repository.NewNode().SetName("child2").SetChildren([]repository.GNode{
            *repository.NewNode().SetName("grandchild3"),
        }),
    })

    go RestAPI(node)

    time.Sleep(100 * time.Millisecond) // wait for server to start
	
	token := repository.GetValidToken()

    tests := []struct {
        name           string
        url            string
        token          string
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "missing tag parameter",
            url:            "http://localhost:8080/taggedContent?token=" + token,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Missing 'tag' parameter\n",
        },
        {
            name:           "unauthorized request",
            url:            "http://localhost:8080/taggedContent?tag=child1&token=invalid",
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   "Unauthorized\n",
        },
        {
            name:           "not found",
            url:            "http://localhost:8080/taggedContent?tag=unknown&token=" + token,
            expectedStatus: http.StatusBadRequest,
            expectedBody:   "Tag unknown was not found\n",
        },
        {
            name:           "success",
            url:            "http://localhost:8080/taggedContent?tag=child1&token=" + token,
            expectedStatus: http.StatusOK,
            expectedBody:   `{"name":"child1","children":[{"name":"grandchild1","children":null},{"name":"grandchild2","children":null}]}`,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            req, err := http.NewRequest("GET", tc.url + "", nil)
            if err != nil {
                t.Fatal(err)
            }
            req.Header.Set("token", tc.token)

            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                t.Fatal(err)
            }
            defer resp.Body.Close()

            if resp.StatusCode != tc.expectedStatus {
                t.Errorf("Expected status %d, but got %d", tc.expectedStatus, resp.StatusCode)
            }

            body, err := ioutil.ReadAll(resp.Body)
            if err != nil {
                t.Fatal(err)
            }

            if string(body) != tc.expectedBody {
                t.Errorf("Expected response body %q, but got %q", tc.expectedBody, string(body))
            }
        })
    }

    // Wait for server to shut down
    time.Sleep(100 * time.Millisecond)
}

