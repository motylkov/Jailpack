.PHONY: build clean test lint install

# Build the application for FreeBSD
build:
	GOOS=freebsd GOARCH=amd64 go build -o bin/jailpack .

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter
lint:
	golangci-lint run

# Install the application
install:
	go install .

# Build for FreeBSD
build-freebsd: clean
	GOOS=freebsd GOARCH=amd64 go build -o bin/jailpack .

# Build from cmd/jailpack.go
build-cmd: clean
	GOOS=freebsd GOARCH=amd64 go build -o bin/jailpack ./cmd/jailpack.go

# Development
dev: build
	./bin/jailpack --help

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./... 