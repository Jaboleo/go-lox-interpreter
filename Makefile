.PHONY: all build test lint clean

all: format test lint build

format:
	go fmt

build:
	go build

test:
	go test -v

lint:
	golangci-lint run