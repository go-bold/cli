# Makefile
.PHONY: build install test clean

# Build the CLI binary
build:
	go build -o bin/bold main.go

# Install the CLI globally
install:
	go install .

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Development: build and install locally
dev: build
	sudo cp bin/bold /usr/local/bin/

# Cross-compile for different platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/bold-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/bold-darwin-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/bold-windows-amd64.exe main.go