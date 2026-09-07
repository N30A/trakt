.PHONY: build run test

build: test
	go build -o bin/trakt .

run: test
	go run .

test:
	go test ./...
