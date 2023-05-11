#!/bin/bash

set -e

export CGO_ENABLED=0
export GO111MODULE="on"
export GOPATH="$HOME/go"

# Build for all supported platforms
GOOS=darwin GOARCH=amd64 go build -o darwin-app ./src
GOOS=linux GOARCH=amd64 go build -o linux-app ./src
GOOS=windows GOARCH=amd64 go build -o windows-app ./src