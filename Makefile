.PHONY: default build test clean tidy doc-test help

# Default target (help)
default: help

# Build target
build: clean tidy test
	@echo "Building..."
	@go build .

# Test target
test:
	@echo "Running tests..."
	@go test ./...

# Tidy target
tidy:
	@echo "Tidying go.mod and go.sum..."
	@go mod tidy

# Clean target
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf go.sum

# Run go doc server to test locally to test documentation
doc-test:
	@echo "Running locally in Docker on port 8080..."
	@echo "Run following command in you home directory 'go get golang.org/x/tools/cmd/godoc && go install golang.org/x/tools/cmd/godoc'"
	@godoc -http :8080

# Help target to display available targets and their descriptions
help:
	@echo "Available targets:"
	@echo "  help    : Show this help message"
	@echo "  build      : Clean, tidy, test, and build"
	@echo "  test       : Run tests"
	@echo "  tidy       : Tidy go.mod and go.sum"
	@echo "  clean      : Clean build artifacts"
	@echo "  doc-test : Run go doc server at 8080 to test locally to test documentation(NOTE: needs godoc )"
