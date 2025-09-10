.PHONY: build build-all test clean install help new-thought

# Variables
APP_NAME=my-cli-tool
VERSION=0.1.0
BUILD_DIR=dist

# Platform-specific variables
GOOS_TARGETS=darwin linux windows
GOARCH_TARGETS=amd64 arm64

# Default target
help:
	@echo "Available targets:"
	@echo "  build       Build for current platform"
	@echo "  build-all   Build for all target platforms"
	@echo "  test        Run tests"
	@echo "  clean       Clean build artifacts"
	@echo "  install     Install dependencies"
	@echo "  new-thought Create a new TGS thought directory"

# Build for current platform
build:
	@echo "Building $(APP_NAME) for current platform..."
	@mkdir -p $(BUILD_DIR)
	@if command -v go >/dev/null 2>&1; then \
		go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) ./src; \
	elif command -v cargo >/dev/null 2>&1; then \
		cargo build --release; \
		cp target/release/$(APP_NAME) $(BUILD_DIR)/; \
	elif command -v node >/dev/null 2>&1; then \
		npm run build; \
	else \
		echo "No supported build tool found (go, cargo, or node)"; \
		exit 1; \
	fi

# Build for all platforms (Go example)
build-all:
	@echo "Building $(APP_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@if command -v go >/dev/null 2>&1; then \
		for goos in $(GOOS_TARGETS); do \
			for goarch in $(GOARCH_TARGETS); do \
				echo "Building for $$goos/$$goarch..."; \
				GOOS=$$goos GOARCH=$$goarch go build \
					-ldflags "-X main.version=$(VERSION)" \
					-o $(BUILD_DIR)/$(APP_NAME)-$$goos-$$goarch \
					./src; \
				if [ "$$goos" = "windows" ]; then \
					mv $(BUILD_DIR)/$(APP_NAME)-$$goos-$$goarch $(BUILD_DIR)/$(APP_NAME)-$$goos-$$goarch.exe; \
				fi; \
			done; \
		done; \
	else \
		echo "Cross-compilation requires Go. Use 'make build' for current platform."; \
	fi

# Run tests
test:
	@echo "Running tests..."
	@if command -v go >/dev/null 2>&1; then \
		go test ./...; \
	elif command -v cargo >/dev/null 2>&1; then \
		cargo test; \
	elif command -v npm >/dev/null 2>&1; then \
		npm test; \
	else \
		echo "No test runner found"; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@if command -v go >/dev/null 2>&1; then go clean; fi
	@if command -v cargo >/dev/null 2>&1; then cargo clean; fi
	@if [ -d "node_modules" ]; then rm -rf node_modules; fi

# Install dependencies
install:
	@echo "Installing dependencies..."
	@if command -v go >/dev/null 2>&1; then \
		go mod tidy; \
		go mod download; \
	elif command -v cargo >/dev/null 2>&1; then \
		cargo fetch; \
	elif command -v npm >/dev/null 2>&1; then \
		npm install; \
	else \
		echo "No package manager found"; \
	fi

# Create a new TGS thought directory
new-thought:
	@if ! command -v git >/dev/null; then echo "git not found in PATH"; exit 2; fi; \
	if [ -z "$(title)" ]; then echo "Usage: make new-thought title=\"short title\""; exit 1; fi; \
	if [ ! -d "agentops/tgs" ]; then echo "Templates missing at agentops/tgs"; exit 2; fi; \
	HASH=$$(git rev-parse --short HEAD); \
	SLUG=$$(printf "%s" "$(title)" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-z0-9]+/-/g' | sed -E 's/^-+|-+$$//g'); \
	DIR="tgs/$$HASH-$$SLUG"; \
	mkdir -p "$$DIR"; \
	for f in agentops/tgs/*; do bn=$$(basename "$$f"); if [ ! -e "$$DIR/$$bn" ]; then cp "$$f" "$$DIR/"; fi; done; \
	if [ ! -f "$$DIR/README.md" ]; then echo "# $$HASH - $(title)" > "$$DIR/README.md"; fi; \
	echo "Created $$DIR"