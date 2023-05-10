# Introduction
Hello Max, 
this code is simple microservice implementing four following tasks which was send in acceptance criteria. Microsevice is "Makefile" driven, which means that you run "make help" command 
which will provide you description of separate command for each task. Makefile is running simple shell scripts which are building go app and run 
them apropriately to each task separately. Scripts should take care of build and run, so I hope for overall user-friendly result.
However, it will require to have GNU Make installed on your system. For case it will be an issue, I pre-build binary tool for linux, windows and 
darwin, which can be run completly stand-alone on base of provided flag (explained in each section below), disregarding any additional tool. I want you to see in my way of thinking and help you to review, so I am also providing pseudo code (or in other words, written proposal) on which I based implemented solution. Most comments I left in code as is to enlight process how code evolved.
Some comments just shared my thoughts, they are like messages for you or code reviewer, I would normally remove them.

# Task one 

### Pseudo code
	// 1. create ouput slice of GNodes
	// 2. call recursive function with 1st parameter given root, 2nd pointer to slice of nodes
	// in recursive function:
	// 3. add root to slice of GNodes
	// 4. get children of root and iterate all elements
	// 5. call recursive on each element of children slice
	// 6. after recursive function with no return value finished, return value of slice nodes
	
### Implement function accepting interface node and returning dynamic arrays of strings
We can assume graph as follow (with one root node)
          A
     /    |    \
    B     C     D
   / \   /|\    | 
  E  F  G H I   J
will print result as this: ["A", "B", "E", "F", "C", "G", "H", "I", "D", "J"] in preorder traversal method. Preorder method
was used, but in same way we could used postorder or inorder.
If you want to get processed different graph, you can alter input_graph.json (keeping acyclic graph rules) and it will reprint new result.
Can be run by command "make walk-graph" in terminal or by "make build" (after code was changed) and then "./bin/<OS>/amd64/cisco-app -walk-graph" or simple by "./bin/<OS>/amd64/cisco-app -walk-graph" using pre-build solution.

### Testing. This is first-written proposal for testing and edge cases
Tests were written in advance. First use case should match exactly given acceptance criteria, because most embarassing would be
fail in that. We are testing base scenarios - graph with zero node, one node, multiple nodes. We are also testing 
Can be run by command "make test-walk-graph" in terminal. We consider tests for documenation of microservice functionality. Thus, 
if tests passed, they current result will ovewrite documenation in html format so we keep one source of true for documenation which 
is code.

# Task two 

### Pseudo code
	// 1. handle nil node
	// 2. data structrues - create empty slice of slices with GNodes as "paths"
	// 2. data structrues - create empty slice of GNodes as "path"
	// 3. call recursive function with input parameters pointer node, pointer to path and paths
	// 4. we checked if node is not nil, in parent function, so add it to the path by dereferencing appending of slice
	// 5. if childrens of node is empty, we reached the bottom
	// 6. appened with dereferencing path to slice path 
	// 7. replace with dereferencing path with empty slice  
	// 10. return value slice of slices
	// 11. check for edge cases / reasonability of return values

# Task three

### Pseudo code
1. Input validation

2. Authentication

3. Error handling

4. Secure communication

5. Data encryption: The server should encrypt sensitive data, such as passwords and authentication tokens, to prevent unauthorized access.

4. Logging and monitoring

