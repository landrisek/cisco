#!/bin/bash

set -e

export CGO_ENABLED=0
export GO111MODULE="on"
export GOPATH="$HOME/go"

# Build for all supported platforms
GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/cisco-app ./src
GOOS=linux GOARCH=amd64 go build -o bin/darwin/amd64/cisco-app ./src
GOOS=windows GOARCH=amd64 go build -o bin/darwin/amd64/cisco-app ./src