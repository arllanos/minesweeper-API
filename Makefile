VERSION := 0.1.0
PROJECT := $(shell basename "$(PWD)")
BINNAME := bin/minesweeper-api
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

build: fmt vet
	@rm -f $(BINNAME)
	go build $(LDFLAGS) -o $(BINNAME) ./cmd/minesweeper-api

run:
	@go run ./cmd/minesweeper-api

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v ./...

up:
	docker compose up -d --build

down:
	docker compose down --remove-orphans