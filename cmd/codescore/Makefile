PROJECT=$(shell basename $(CURDIR))

VERSION ?= $(shell git tag -l --sort=-v:refname | head -n1 2>/dev/null || echo dev)
COMMIT  = $(shell git rev-parse --short HEAD 2>/dev/null || echo dev-unknown-commit)
DATE    = $(shell git log -1 --format=%cd --date=format:"%Y-%m-%d" HEAD 2>/dev/null || echo dev-unknown-date)
LDFLAGS = -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

all: build

check:
	gofumpt -w .
	gofmt -w -s .
	go fix ./...
	go vet ./...
	staticcheck

run:
	go build -ldflags "$(LDFLAGS)" -o codescore .
	./codescore

build:
	go build -ldflags "$(LDFLAGS)" -o codescore .

clean:
	rm -f codescore
