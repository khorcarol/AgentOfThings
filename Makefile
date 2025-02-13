.DEFAULT_GOAL=build
.PHONY: clean format build

SRCS:=$(wildcard *.go)

build: format
	go build -o build/agentofthings

clean:
	rm -rf build/*

format: 
	go fmt ./...
