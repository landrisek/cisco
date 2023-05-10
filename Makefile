.PHONY: help build

help:
	@echo "make help                        Show this help message"
	@echo "make build                       Build standalone binary for linux, darwin or windows for all tasks driven by flag, see README.md"
	@echo "make walk-graph                  Take a structure of graph in input_graph.json, collect all nodes and prints their name in terminal"
	@echo "make test-walk-graph             Run test suite for walk graph functionality"
	@echo "make paths                       Take a structure of graph in input_graph.json and return all possible path until bottom is reached"
	@echo "make test-path                   Run test suite for paths functionality"

r:
	./scripts/install-go.sh && ./scripts/test-rest-api.sh

build:
	./scripts/install-go.sh && ./scripts/build.sh

walk-graph:
	./scripts/install-go.sh && ./scripts/walk-graph.sh

test-walk-graph:
	./scripts/install-go.sh && ./scripts/test-walk-graph.sh

paths:
	./scripts/install-go.sh && ./scripts/paths.sh

test-paths:
	./scripts/install-go.sh && ./scripts/test-paths.sh

rest-api:
	./scripts/install-go.sh && ./scripts/rest-api.sh

test-rest-api:
	./scripts/install-go.sh && ./scripts/test-rest-api.sh