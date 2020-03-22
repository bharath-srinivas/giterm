BINARY = giterm
BUILD_DIR = $(shell pwd)
GO_VERSION = $(shell go version | cut -d ' ' -f 3 | sed -e 's/go//')
LDFLAGS = -ldflags="-s -w"

# Disable CGO
export CGO_ENABLED=0

all: clean test build

# Test all packages.
test:
	@go test -cover ./...
.PHONY: test

# Build linux binaries.
linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}_linux_amd64 ${LDFLAGS} ./cmd/giterm;
.PHONY: linux

linux-386:
	GOOS=linux GOARCH=386 go build -o ${BINARY}_linux_386 ${LDFLAGS} ./cmd/giterm;
.PHONY: linux-386

# Build mac binaries.
darwin:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY}_darwin_amd64 ${LDFLAGS} ./cmd/giterm;
.PHONY: darwin

darwin-386:
	GOOS=darwin GOARCH=386 go build -o ${BINARY}_darwin_386 ${LDFLAGS} ./cmd/giterm;
.PHONY: darwin-386

# Build release binaries.
build: linux linux-386 darwin darwin-386
.PHONY: build

# Clean build artifacts.
clean:
	@rm -f ${BINARY}_linux*
	@rm -f ${BINARY}_darwin*
.PHONY: clean
