all:
	go run main.go 

check: 
	gofumpt -w .
	gofmt -w -s .
	go fix ./...
	go vet ./...
	staticcheck
