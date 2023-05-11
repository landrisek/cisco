package controller

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/landrisek/cisco/src/repository"
)

// Log logs an error message along with the provided error.
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg+": %s", err)
	}
}

type testCase struct {
	name            string
	input           repository.GNode
	expectedLength  int
	expectedContent []string
	expectedPaths   [][]string
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
	fmt.Fprintln(file, "<head><title>Documentation of "+name+" function</title></head>")
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
