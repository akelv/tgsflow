.PHONY: help new-thought bootstrap test-bootstrap clean-templates

help:
	@echo "Targets:"
	@echo "  new-thought title=\"short title\"   Scaffold a new TGS thought directory"
	@echo "  bootstrap                          Test bootstrap script locally"
	@echo "  test-bootstrap                     Test bootstrap with all templates"
	@echo "  clean-templates                    Clean template artifacts"

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


