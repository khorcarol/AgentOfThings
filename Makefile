.DEFAULT_GOAL=build
.PHONY: clean format build generate-bundled build-android

SRCS:=$(wildcard *.go)

GO_BINARY_PATH:=$(shell go env GOPATH)/bin

build: format
	go build -o build/agentofthings

build-android:
	GOFLAGS="-ldflags=-checklinkname=0" $(GO_BINARY_PATH)/fyne package -os android -appID com.groupalpha.agentofthings -icon assets/golang.png

generate-bundled:
	go generate ./frontend

clean:
	rm -rf build/*

format:
	go mod tidy
	go fmt ./...

