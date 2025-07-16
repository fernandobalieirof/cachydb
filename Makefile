test:
	@go test -v ./...
fmt:
	@go fmt ./...
run: build
	@bin/cachydb
build:
	@go build -o bin/cachydb .
