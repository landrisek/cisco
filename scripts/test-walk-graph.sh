export GO111MODULE="on"
export GOPATH="$HOME/go"

go test -v ./src/controller -run TestWalkGraph
