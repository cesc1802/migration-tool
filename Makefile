.PHONY: build run test lint clean

BINARY=migrate-tool
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY) ./cmd/migrate-tool

run:
	go run ./cmd/migrate-tool $(ARGS)

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
