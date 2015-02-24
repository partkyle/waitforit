.PHONY: all build-all

VERSION = $(shell git describe --tag)
LDFLAGS = -ldflags "-X main.VERSION $(VERSION)"

all:
	fig build
	fig pull
	fig run --rm build make test build-all

test:
	go test ./...

build-all: build/waitforit-darwin-amd64 build/waitforit-linux-amd64

build/waitforit-darwin-amd64: *.go .git
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o build/waitforit-darwin-amd64

build/waitforit-linux-amd64: *.go .git
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o build/waitforit-linux-amd64

clean:
	rm -rf build
