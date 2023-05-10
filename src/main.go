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
		nodes := controller.WalkGraph()

		// Print all nodes
		for _, node := range nodes {
			fmt.Println(node.GetName())
		}
	}

	// Handle "paths" flag
	if *pathsGraph {

		// Call WalkGraph function to traverse the graph and get all paths
		paths := controller.Paths()

		for _, path := range paths {
			for i, node := range path {
				fmt.Prinln(node.GetName())
			}
		}
	}

	// Handle "paths" flag
	if *restAPI {
		controller.RestAPI()
	}
}
