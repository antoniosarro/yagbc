.PHONY: build test clean run demo coverage help

# Build the emulator
build:
	@echo "Building YAGBC..."
	@go build -o yagbc .

# Run all tests
test:
	@echo "Running tests..."
	@go test -v -count=1 ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

# Run the main program
run: build
	@./yagbc

# Run the CPU demo
demo:
	@echo "Running CPU demo..."
	@go run examples/cpu_demo.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f yagbc coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build    - Build the emulator"
	@echo "  test     - Run all tests"
	@echo "  coverage - Run tests with coverage report"
	@echo "  run      - Build and run the emulator"
	@echo "  demo     - Run the CPU demo"
	@echo "  clean    - Remove build artifacts"
	@echo "  fmt      - Format code"
	@echo "  lint     - Run linter"
	@echo "  help     - Show this help"