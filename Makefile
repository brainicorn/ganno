all: local
.PHONY: all clean local init build test coverage lint

init:
	@go mod tidy

# Cleans our project: deletes binaries
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; echo "removing ${BINARY}"; fi
	@if [ -f ${DEBUGBINARY} ] ; then rm ${DEBUGBINARY} ; echo "removing ${DEBUGBINARY}"; fi

local: init lint test

build:
	CGO_ENABLED=0 go build ./...

test:
	@go test ./...

coverage:
	@go test -coverprofile=coverage.txt -covermode=atomic

lint:
	@golangci-lint run ./...
