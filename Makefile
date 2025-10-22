SHELL := /bin/bash

.PHONY: build test tidy

build:
	go build -o ./bin/tgs ./src

test:
	go test -coverprofile=coverage.out -covermode=atomic -coverpkg=./... ./...
	@go tool cover -func=coverage.out | tee coverage.txt
	@go tool cover -html=coverage.out -o coverage.html

tidy:
	go mod tidy

.PHONY: help new-thought bootstrap test-bootstrap clean-templates install-cli install-cli-dry-run test-install ears-gen

help:
	@echo "Targets:"
	@echo "  new-thought title=\"short title\" [spec=\"idea\"]   Scaffold a new TGS thought directory"
	@echo "  bootstrap                          Test bootstrap script locally"
	@echo "  test-bootstrap                     Test bootstrap with all templates"
	@echo "  clean-templates                    Clean template artifacts"
	@echo "  install-cli                        Install latest tgs via scripts/install.sh"
	@echo "  install-cli TAG=vX.Y.Z             Install specific tag (requires published release)"
	@echo "  install-cli-dry-run                Print resolved URL/paths without installing"
	@echo "  test-install                       Dry-run + URL HEAD checks for installer"


new-thought:
	@if ! command -v git >/dev/null; then echo "git not found in PATH"; exit 2; fi; \
	if [ -z "$(title)" ]; then echo "Usage: make new-thought title=\"short title\" [spec=\"idea\"]"; exit 1; fi; \
	if [ ! -d "tgs/agentops/tgs" ]; then echo "Templates missing at tgs/agentops/tgs"; exit 2; fi; \
	HASH=$$(git rev-parse --short HEAD); \
	SLUG=$$(printf "%s" "$(title)" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '-' | sed 's/^-*//' | sed 's/-*$$//'); \
	DIR="tgs/thoughts/$$HASH-$$SLUG"; \
	mkdir -p "$$DIR"; \
	for f in tgs/agentops/tgs/*; do bn=$$(basename "$$f"); if [ ! -e "$$DIR/$$bn" ]; then cp "$$f" "$$DIR/"; fi; done; \
	if [ ! -f "$$DIR/README.md" ]; then \
		{ \
			printf "# %s - %s\n\n" "$$HASH" "$(title)"; \
			printf -- "- Base Hash: \`%s\`\n\n" "$$HASH"; \
			printf "## Quick Links\n- [research.md](./research.md)\n- [plan.md](./plan.md)\n- [implementation.md](./implementation.md)\n\n"; \
			if [ -n "$(spec)" ]; then printf "## Idea Spec\n%s\n" "$(spec)"; fi; \
		} > "$$DIR/README.md"; \
	fi; \
	echo "Created $$DIR"

bootstrap:
	@echo "Testing bootstrap script locally..."
	@if [ ! -f "bootstrap.sh" ]; then echo "bootstrap.sh not found"; exit 1; fi
	@echo "Run: ./bootstrap.sh"
	@echo "Note: This will clone the current directory structure"

test-bootstrap:
	@echo "Testing bootstrap script with dry run..."
	@./bootstrap.sh --dry-run 2>/dev/null || echo "Add --dry-run support to bootstrap.sh for testing"

clean-templates:
	@echo "Cleaning template artifacts..."
	@find templates/ -name "node_modules" -type d -exec rm -rf {} + 2>/dev/null || true
	@find templates/ -name "target" -type d -exec rm -rf {} + 2>/dev/null || true
	@find templates/ -name "dist" -type d -exec rm -rf {} + 2>/dev/null || true
	@find templates/ -name "*.log" -type f -delete 2>/dev/null || true
	@echo "Template artifacts cleaned"

install-cli:
	@bash ./scripts/install.sh

install-cli-dry-run:
	@DRY_RUN=1 bash ./scripts/install.sh

test-install:
	@bash ./scripts/test-install.sh


# Generate ANTLR4 Go lexer/parser for EARS grammar
ears-gen:
	@if ! command -v antlr4 >/dev/null && ! command -v antlr >/dev/null; then \
		echo "antlr4 not found. Install via brew: brew install antlr && echo 'export CLASSPATH=\"$$CLASSPATH:$$(brew --prefix)/libexec/antlr-4.13.1-complete.jar\"' >> ~/.zshrc'"; exit 2; \
	fi; \
	GEN_DIR=src/core/ears/gen; mkdir -p $$GEN_DIR; \
	ANTLR_CMD=$$(command -v antlr4 || command -v antlr); \
	$$ANTLR_CMD -Dlanguage=Go -o $$GEN_DIR src/core/ears/ears.g4


