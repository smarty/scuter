#!/usr/bin/make -f

test: fmt
	GOEXPERIMENT=jsonv2 go test -timeout=1s -short -race -covermode=atomic ./...

fmt:
	go fmt ./... && go mod tidy

compile:
	GOEXPERIMENT=jsonv2 go build ./...

build: test compile

.PHONY: test fmt compile build
