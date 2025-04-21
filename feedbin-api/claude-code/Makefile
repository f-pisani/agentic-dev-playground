.PHONY: build test lint fmt vet clean

# Default target
all: build

# Build the code
build:
	go build ./...

# Run tests
test:
	go test -v ./...

# Format the code
fmt:
	go fmt ./...

# Check code style
lint:
	go vet ./...

# Run go vet
vet:
	go vet ./...

# Clean build artifacts
clean:
	go clean