.DEFAULT_GOAL=build
.PHONY: clean format build generate-bundled

SRCS:=$(wildcard *.go)

build: format
	go build -o build/agentofthings

generate-bundled:
	go generate ./frontend

clean:
	rm -rf build/*

format: 
	go mod tidy
	go fmt ./...

