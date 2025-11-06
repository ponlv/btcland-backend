.PHONY: dev build run clean install-air

# Development commands
dev:
	@echo "Starting development server with hot reload..."
	air

# Build the application
build:
	@echo "Building application..."
	go build -o ./tmp/main .

# Run the application (without hot reload)
run:
	@echo "Running application..."
	go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ./tmp/
	rm -f build-errors.log

# Install Air hot reload tool
install-air:
	@echo "Installing Air hot reload tool..."
	go install github.com/air-verse/air@latest

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Help
help:
	@echo "Available commands:"
	@echo "  dev          - Start development server with hot reload"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application (without hot reload)"
	@echo "  clean        - Clean build artifacts"
	@echo "  install-air  - Install Air hot reload tool"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  help         - Show this help message"
