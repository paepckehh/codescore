PROJECT=$(shell basename $(CURDIR))

all:
	make -C cmd/$(PROJECT) all

deps: 
	rm go.mod go.sum
	go mod init paepcke.de/$(PROJECT)
	go mod tidy -v	

check: 
	gofumpt -w .
	gofmt -w -s .
	go fix ./...
	go vet ./...
	staticcheck
	make -C cmd/$(PROJECT) check
