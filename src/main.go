package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/landrisek/cisco/src/controller"
)

func main() {
	// Define command line flags
	walkGraph := flag.Bool("walk-graph", false, "Walk the graph")
	pathsGraph := flag.Bool("paths", false, "Find all paths in graph")
	restAPI := flag.Bool("rest-api", false, "Run server with rest API for tags")

	// Parse command line flags
	flag.Parse()

	// Handle "walk-graph" flag
	if *walkGraph {
		// Call UploadJson function to read the input JSON and create the graph
		graph, err := controller.UploadJson("input_graph.json")
		controller.Log(err, "Error uploading JSON")

		// Call WalkGraph function to traverse the graph and get all nodes
		nodes := controller.WalkGraph(graph)

		// Print all nodes
		for _, node := range nodes {
			fmt.Println(node.GetName())
		}
	}

	// Handle "paths" flag
	if *pathsGraph {
		// Call UploadJson function to read the input JSON and create the graph
		graph, err := controller.UploadJson("input_graph.json")
		controller.Log(err, "Error uploading JSON")

		// Call WalkGraph function to traverse the graph and get all paths
		paths := controller.Paths(graph)
		// Print the paths in the desired format
		fmt.Print("paths(A) = (")
		for _, path := range paths {
			fmt.Print(" (")
			for i, node := range path {
				fmt.Print(node.GetName())
				if i != len(path)-1 {
					fmt.Print(" ")
				}
			}
			fmt.Print(")")
		}
		fmt.Print(" )")
	}

	// Handle "paths" flag
	if *restAPI {
		defer func() {
			fmt.Printf("container died on %v"+"\n", time.Now())
		}()
		fmt.Printf("container started on %v"+"\n", time.Now())
		// Call UploadJson function to read the input JSON and create the graph
		tags, err := controller.UploadJson("input_tags.json")
		controller.Log(err, "Error uploading JSON")
		controller.RestAPI(tags)
	}
}
