package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"

	"github.com/landrisek/cisco/src/repository"
)

type tagServer struct {
	tags repository.GNode
	ctx  context.Context
}

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

// WalkGraph traverses the graph starting from the given node and returns a slice of all visited nodes.
// If the provided node is nil, it returns an empty slice.
func WalkGraph(node repository.GNode) []repository.GNode {
	if node == nil {
		return []repository.GNode{}
	}
	var nodes []repository.GNode
	findNode(node, &nodes)
	return nodes
}

// findNode recursively traverses the graph starting from the given node and appends each visited node to the nodes slice.
// It updates the nodes slice by reference using a pointer.
func findNode(node repository.GNode, nodes *[]repository.GNode) {
	*nodes = append(*nodes, node)
	for _, children := range node.GetChildren() {
		findNode(children, nodes)
	}
}

// Paths returns a slice of paths in the graph starting from the given node. Each path is represented as a slice of nodes.
// If the node is nil, an empty slice is returned.
func Paths(node repository.GNode) [][]repository.GNode {
	var paths [][]repository.GNode
	if nil == node {
		return paths
	}
	var path []repository.GNode
	findBottom(node, &path, &paths)
	return paths
}

// findBottom is a recursive function that finds all paths in the graph starting from the given node and appends them to the paths slice.
// It keeps track of the current path in the 'path' parameter and appends the completed paths to the 'paths' parameter.
// If the node has no children, the current path is considered complete and added to the 'paths' slice.
// After processing a child node, the function removes it from the current path to prepare for the next iteration.
func findBottom(node repository.GNode, path *[]repository.GNode, paths *[][]repository.GNode) {
	*path = append(*path, node)
	childrens := node.GetChildren()
	if len(childrens) == 0 {
		*paths = append(*paths, append([]repository.GNode{}, *path...))
	} else {
		for _, child := range childrens {
			findBottom(child, path, paths)
		}
	}
	// HINT: as path [A,B] was added, but next call will be on C, we need path looks like [A] where we add C so it will be [A,C]
	*path = (*path)[:len(*path)-1]
}

// RestAPI starts an HTTP server that exposes a REST API for interacting with the given node in the graph.
// It handles requests to the "/taggedContent" endpoint by serving the tagServer handler with the provided node and context.
// The server is started in a separate goroutine and listens for incoming requests.
// It gracefully shuts down the server and canceling the context.
// This function creates a context and a cancel function to control the server and goroutines.
func RestAPI(node repository.GNode) {
	ctx, cancel := context.WithCancel(context.Background())
	http.Handle("/taggedContent", &tagServer{
		tags: node,
		ctx:  ctx,
	})

	// HINT: this is on discussion
	//http.Handle("/heap", pprof.Handler("heap").ServeHTTP) // DO NOT PUSH TO PROD

	server := &http.Server{
		Addr: "localhost:8080",
		// HINT: These are also preventing DDoS atacks
		ReadTimeout:    60 * time.Second, // DDoS
		WriteTimeout:   60 * time.Second, // DDoS
		MaxHeaderBytes: 1 << 20,          // DDos
	}
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			Log(err, "Error running tag server")
		}
	}()

	// HINT: Wait for SIGINT (Ctrl+C) signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		Log(err, "Error shutting down tag server")
	}

	cancel()
}

// ServeHTTP handles HTTP requests for the tagServer handler.
// It checks the request method, headers, and parameters for valid CORS, authentication, and tag information.
// It optimizes goroutine utilization by setting GOMAXPROCS based on the number of available CPUs.
// It retrieves the subtags from the repository using GetSubTags and returns an error if not found.
// It encodes the subtags as JSON and writes the response to the client with appropriate headers.
func (server tagServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	headers := writer.Header()
	headers.Set("Access-Control-Allow-Origin", "http://localhost")
	headers.Set("Access-Control-Allow-Methods", "GET")
	headers.Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")

	if request.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parameters := request.URL.Query()
	token := parameters.Get("token")
	if !repository.IsAuthenticated(token) {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tag := parameters.Get("tag")
	if tag == "" {
		http.Error(writer, "Missing 'tag' parameter", http.StatusBadRequest)
		return
	}

	if runtime.NumCPU() > 1 {
		// HINT: if there are multiple CPUs, let`s be not too greedy and use less than the total available
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	} else {
		runtime.GOMAXPROCS(1)
	}

	subtags := repository.GetSubTags(server.ctx, server.tags, tag)

	if subtags.GetName() == "" {
		http.Error(writer, fmt.Sprintf("Tag %s was not found", tag), http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(&subtags)
	if err != nil {
		http.Error(writer, "Error encoding response as JSON", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	// HINT: let`s help a client
	writer.Header().Set("Content-Length", strconv.Itoa(len(jsonBytes)))
	writer.Write(jsonBytes)
}

// Log logs an error message along with the provided error.
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg+": %s", err)
	}
}
