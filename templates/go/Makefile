.PHONY: build test clean run install help new-thought

# Variables
BINARY_NAME=my-go-project
BUILD_DIR=bin
MAIN_PATH=./cmd/main

# Default target
help:
	@echo "Available targets:"
	@echo "  build       Build the binary"
	@echo "  test        Run tests"
	@echo "  clean       Clean build artifacts"
	@echo "  run         Run the application"
	@echo "  install     Install dependencies"
	@echo "  new-thought Create a new TGS thought directory"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run the application
run:
	@go run $(MAIN_PATH)

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Create a new TGS thought directory
new-thought:
	@if ! command -v git >/dev/null; then echo "git not found in PATH"; exit 2; fi; \
	if [ -z "$(title)" ]; then echo "Usage: make new-thought title=\"short title\""; exit 1; fi; \
	if [ ! -d "tgs/agentops/tgs" ]; then echo "Templates missing at tgs/agentops/tgs"; exit 2; fi; \
	HASH=$$(git rev-parse --short HEAD); \
	SLUG=$$(printf "%s" "$(title)" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' | sed 's/^-*//' | sed 's/-*$$//'); \
	DIR="tgs/$$HASH-$$SLUG"; \
	mkdir -p "$$DIR"; \
	for f in tgs/agentops/tgs/*; do bn=$$(basename "$$f"); if [ ! -e "$$DIR/$$bn" ]; then cp "$$f" "$$DIR/"; fi; done; \
	if [ ! -f "$$DIR/README.md" ]; then echo "# $$HASH - $(title)" > "$$DIR/README.md"; fi; \
	echo "Created $$DIR"