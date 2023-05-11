package controller

import (
	"context"
	"encoding/json"
	"fmt"
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
	//http.Handle("/heap", pprof.Handler("heap").ServeHTTP)

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
