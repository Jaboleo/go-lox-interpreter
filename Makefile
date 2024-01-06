.PHONY: all build test lint clean

all: test build lint

build:
	go build

test:
	go test -v

lint:
	golangci-lint run