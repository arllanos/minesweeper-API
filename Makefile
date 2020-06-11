VERSION := 0.1.0
PROJECT := $(shell basename "$(PWD)")
BINNAME := minesweeper-api
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build: fmt vet
	@rm -f $(BINNAME)
	go build $(LDFLAGS) -o $(BINNAME) main.go

run:
	@go run main.go

fmt:
	go fmt ./...

vet:
	go vet ./...
