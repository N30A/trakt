.PHONY: build run test

build: test
	mkdir -p bin
	go build -o bin/ ./cmd/*

run: test build
	./bin/trakt

test:
	go test ./...
