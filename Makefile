.PHONY: all build test lint clean

all: format build

format:
	go fmt

build:
	go build

test:
	go test -v

lint:
	golangci-lint run