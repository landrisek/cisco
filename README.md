# Introduction
Hello Max,

This code is a simple microservice that implements four tasks specified in the acceptance criteria. The microservice is driven by a Makefile, which means you can run the "make help" command to get descriptions of each task.

The Makefile runs shell scripts that build the Go application and execute the tasks accordingly. The scripts handle the build and execution processes, aiming for a user-friendly experience. Please make sure you have GNU Make installed on your system to use the Makefile.

Additionally, I have pre-built binary tools for Linux, Windows, and Darwin systems, which can be run independently with the provided flags (explained in each section below), eliminating the need for additional tools.

I have included pseudo code (or proposals) that served as the basis for the implemented solution. 
In the code, you will find two types of comments:

Proper comments that mirror the documentation status and are typically left in the code as part of the documentation. These comments provide explanations and details about the code implementation.

Comments starting with the word "HINT." These comments are not typically left in the final codebase, but they serve the purpose of illustrating the evolution of the code and sharing my thoughts during the development process.

# Task one: Walk graph

### Pseudo code
1. Create an output slice of GNodes.
2. Call the recursive function with the first parameter as the root node and the second parameter as a pointer to the slice of nodes.
In the recursive function:
3. Add the root node to the slice of GNodes.
4. Get the children of the root node and iterate through each element.
5. Call the recursive function on each child element.
6. After the recursive function has finished (no return value), return the slice of nodes.
	
### Implementation
Taks was function accepting interface node and returning dynamic arrays of strings
We can assume graph as follow (with one root node)
          A
     /    |    \
    B     C     D
   / \   /|\    | 
  E  F  G H I   J

Implement a function that accepts an interface node and returns a dynamic array of strings.
The function performs a preorder traversal of the graph represented by the node, assuming a single root node.
The expected result is an array of strings representing the nodes in preorder traversal order.
For example, given the graph:
         A
    /    |    \
   B     C     D
  / \   /|\    | 
 E  F  G H I   J
The function will print the result as ["A", "B", "E", "F", "C", "G", "H", "I", "D", "J"].
The preorder traversal method is used, but the implementation can be easily modified for postorder or inorder traversal.
To test with a different graph, you can alter the input_graph.json file (while maintaining acyclic graph rules), and the function will generate a new result based on the modified graph.

### How to run
You can run the functionality by:
1. executing the command "make walk-graph" in the terminal
2. if you did make changes in code which you like to test, building the code with "make build" and running ./<your_operation_system>-app -walk-graph
3. directly executing the pre-built solution running ./<your_operation_system>-app -walk-graph

### Testing

Resutls of tests are also generating documentation in WalkGraph-docs.html in root folder.
Tests were written in advance to ensure the functionality of the microservice.
The first use case is designed to match the given acceptance criteria precisely, as any failure in this case would be critical.
The tests cover various scenarios including graphs with zero nodes, one node, and multiple nodes.
Additionally, we test the disconnected graph and (generated) random test case to ensure correctness.
The tests serve as documentation for the microservice functionality, and if they pass, the current results will overwrite the documentation in HTML format. This approach ensures that the code remains the single source of truth for documentation.
You can run the test by executing the command "make test-walk-graph" in the terminal or by ./<your_operation_system>-app -test-walk-graph if you do not have GNU Make install and you do not want to allow it to install go on your machine.

# Task two: Paths 

### Pseudo code

1. Handle nil node
2. Data structures - create empty slice of slices with GNodes as "paths",
3. Data structures - create empty slice of GNodes as "path"/
4. Cll recursive function with input parameters pointer node, pointer to path and paths.
5. We checked if node is not nil in parent function, so add it to the path by dereferencing appending of slice.
6. If childrens of node is empty, we reached the bottom.
7. Appened with dereferencing path to slice path.
8. Replace with dereferencing path with empty slice.  
9. Return value slice of slices.
10. Check for edge cases / reasonability of return values.

### Implementation
Input data are taken from input_graph.json. 

### How to run
You can run the functionality by:
1. executing the command "make walk-graph" in the terminal
2. if you did make changes in code which you like to test, building the code with "make build" and running ./<your_operation_system>-app -walk-graph
3. directly executing the pre-built solution running ./<your_operation_system>-app -walk-graph

### Testing
Resutls of tests are also generating documentation in Paths-docs.html in root folder. Similary to previous tasks, is covered main scenarios with empty node, one node and multiple nodes.
You can run the test by executing the command "make test-paths" in the terminal or by ./<your_operation_system>-app -test-paths if you do not have GNU Make install and you do not want to allow it to install go on your machine.

# Task three: Rest API Server

### Pseudo code
1. RestAPI starts an HTTP server that exposes a REST API for interacting with the given node in the graph.
2. Handle requests to the "/taggedContent" endpoint by serving the tagServer handler with the provided node and context.
3. Start server in a separate goroutine and listens for incoming requests.
4. Also listens for the SIGINT signal (Ctrl+C) to gracefully shutdown the server when the signal is received.
5. Server to be shutdown by calling server.Shutdown() and the context is canceled to signal to other goroutines to exit.
6. Create a context and a cancel function to control the server and goroutines
7. ServeHTTP handles HTTP requests for the tagServer handler.
8. Checks the request method and headers for CORS (Cross-Origin Resource Sharing) and authentication.
9. If the request method is not GET, return a "Method not allowed" error.
10. If the authentication token is invalid, return an "Unauthorized" error.
11. If the 'tag' parameter is missing in the request, return a "Missing 'tag' parameter" error.
12. Sets the GOMAXPROCS value based on the number of available CPUs to optimize goroutine utilization.
13. Retrieves the subtags from the repository using the GetSubTags function.
14. If the subtags are not found, returns an error indicating that the tag was not found.
15. Encodes the subtags as JSON and writes the response to the client.
16. Sets the Content-Type header to "application/json" and the Content-Length header to the length of the JSON data.

### Implementation

The GetSubTags function is part of the repository package and provides the functionality to find the first occurrence of a specific tag within a given node and its children. The reason for placing it in the repository package is that it aligns with the purpose of the package, which is to handle data-related operations. It utilizes goroutines and parallel processing to improve performance in the hot spot section of the code, where the search for the tag is performed.
The function creates a channel to receive the result, initializes an active counter to track the number of active goroutines, and sets up a done channel for signaling when all goroutines have completed their execution.
Inside the lookupChildrens function, which is called recursively, the node and its children are examined to find a match with the tag. If a match is found, the corresponding node is sent to the result channel. If the context or inner context is canceled, the function exits gracefully. In the case of parallel processing, goroutines are spawned to process each child node.
The GetSubTags function waits for the result by listening to various channels, including the context from the server, the inner context for processing nodes, the result channel, and the done channel. Depending on the scenario, it returns the appropriate result or an empty MyNode struct to indicate no match was found.
By placing this functionality in the repository package, it follows a logical grouping of operations related to finding and retrieving data from the underlying data structures. It promotes code organization and separation of concerns, making the codebase more maintainable and understandable.

### How to run
Exposed on http://localhost:8080/taggedContent?tag=animals&token=YYY
You can start server by:
1. executing the command "make rest-api" in the terminal
2. if you did make changes in code which you like to test, building the code with "make build" and running ./<your_operation_system>-app -rest-api
3. directly executing the pre-built solution running ./<your_operation_system>-app -rest-api

### Testing
We are covering main test cases like not providing tag or token, or provided one or both of them invalid. Return of main http statuses
like bad request, unauthorize request, not found and success.
You can run the test by executing the command "make test-rest-api" in the terminal or by ./<your_operation_system>-app -test-rest-api if you do not have GNU Make install and you do not want to allow it to install go on your machine.