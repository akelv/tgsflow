#!/bin/bash
set -euo pipefail

# Repository Bootstrap Script
# Clones this repository and sets up a new project with selected template

REPO_URL="https://github.com/akelv/tgsflow.git"
SCRIPT_VERSION="1.2.0"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

# Error handling
cleanup() {
    if [[ -n "${TEMP_DIR:-}" && -d "$TEMP_DIR" ]]; then
        log_info "Cleaning up temporary directory: $TEMP_DIR"
        rm -rf "$TEMP_DIR"
    fi
}

error_exit() {
    log_error "$1"
    cleanup
    exit 1
}

trap cleanup EXIT
trap 'error_exit "Script interrupted"' INT TERM

# Global flags (can be set via CLI)
DECORATE=0
DRY_RUN=0
FORCE=0
WITH_TEMPLATES="none"

# Helpers to handle dry-run copies
safe_mkdir_p() {
    local dir="$1"
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo -e "${BLUE}[DRY-RUN]${NC} mkdir -p \"$dir\""
        return 0
    fi
    mkdir -p "$dir"
}

safe_copy() {
    local src="$1"
    local dst="$2"
    local dst_dir
    dst_dir=$(dirname "$dst")
    if [[ ! -e "$src" ]]; then
        log_warning "Source not found, skipping: $src"
        return 0
    fi
    if [[ -e "$dst" && "$FORCE" -ne 1 ]]; then
        log_info "Exists, skipping (use --force to overwrite): $dst"
        return 0
    fi
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo -e "${BLUE}[DRY-RUN]${NC} cp -a \"$src\" \"$dst\""
        return 0
    fi
    safe_mkdir_p "$dst_dir"
    cp -a "$src" "$dst"
}

append_if_missing() {
    local file="$1"
    local line="$2"
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo -e "${BLUE}[DRY-RUN]${NC} ensure line in $file: $line"
        return 0
    fi
    if [[ ! -f "$file" ]]; then
        printf "%s\n" "$line" > "$file"
        return 0
    fi
    if ! grep -Fq "$line" "$file" 2>/dev/null; then
        printf "%s\n" "$line" >> "$file"
    fi
}

write_tgs_mk() {
    local target="tgs.mk"
    if [[ -f "$target" && "$FORCE" -ne 1 ]]; then
        log_info "Exists, skipping tgs.mk (use --force to overwrite)"
        return 0
    fi
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo -e "${BLUE}[DRY-RUN]${NC} write $target with new-thought target"
        return 0
    fi
    cat > "$target" <<'EOF'
.PHONY: new-thought

new-thought:
	@if ! command -v git >/dev/null; then echo "git not found in PATH"; exit 2; fi; \
	if [ -z "$(title)" ]; then echo "Usage: make new-thought title=\"short title\" [spec=\"idea\"]"; exit 1; fi; \
	if [ ! -d "agentops/tgs" ]; then echo "Templates missing at agentops/tgs"; exit 2; fi; \
	HASH=$$(git rev-parse --short HEAD); \
	SLUG=$$(printf "%s" "$(title)" | tr '[:upper:]' '[:lower:]' | sed -E 's/[^a-zA-Z0-9]+/-/g' | sed -E 's/^-+|-+$$//g'); \
	DIR="tgs/$$HASH-$$SLUG"; \
	mkdir -p "$$DIR"; \
	for f in agentops/tgs/*; do bn=$$(basename "$$f"); if [ ! -e "$$DIR/$$bn" ]; then cp "$$f" "$$DIR/"; fi; done; \
	if [ ! -f "$$DIR/README.md" ]; then \
		{ \
			printf "# %s - %s\n\n" "$$HASH" "$(title)"; \
			printf "- Base Hash: \`%s\`\n\n" "$$HASH"; \
			printf "## Quick Links\n- [research.md](./research.md)\n- [plan.md](./plan.md)\n- [implementation.md](./implementation.md)\n\n"; \
			if [ -n "$(spec)" ]; then printf "## Idea Spec\n%s\n" "$(spec)"; fi; \
		} > "$$DIR/README.md"; \
	fi; \
	echo "Created $$DIR"
EOF
}

ensure_makefile_includes_tgs() {
    if [[ "$DRY_RUN" -eq 1 ]]; then
        echo -e "${BLUE}[DRY-RUN]${NC} ensure Makefile includes tgs.mk or create minimal Makefile"
        return 0
    fi
    if [[ ! -f "Makefile" ]]; then
        echo "include tgs.mk" > Makefile
        return 0
    fi
    if ! grep -Eq '^new-thought:' Makefile 2>/dev/null && ! grep -Fq 'include tgs.mk' Makefile 2>/dev/null; then
        echo "" >> Makefile
        echo "include tgs.mk" >> Makefile
    fi
}

# Usage
print_usage() {
    cat >&2 <<USAGE
Usage: ./bootstrap.sh [--decorate] [--dry-run] [--force] [--with-templates=<type>] [--help]

Modes:
  --decorate                Decorate the current directory with minimal TGS workflow files.

Flags:
  --dry-run                 Show actions without making changes.
  --force                   Overwrite existing files when copying/writing.
  --with-templates=TYPE     Optional template overlay (react|python|go|cli|none). Default: none.
  -h, --help                Show this help.
USAGE
}

# Validation functions
validate_dependencies() {
    local deps=("git" "curl")
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            error_exit "Required dependency '$dep' not found in PATH"
        fi
    done
}

# Name sanitization
sanitize_name() {
    local input_name="$1"
    # Remove control characters and trim whitespace
    input_name=$(printf '%s' "$input_name" | tr -d '\r\n\t' | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//')
    # Keep only allowed ASCII characters
    input_name=$(printf '%s' "$input_name" | LC_ALL=C tr -cd 'A-Za-z0-9_-')
    # Ensure starts with an alphanumeric
    input_name=$(printf '%s' "$input_name" | sed -E 's/^[^A-Za-z0-9]+//')
    echo "$input_name"
}

validate_project_name() {
    local name="$1"
    if [[ ! "$name" =~ ^[a-zA-Z0-9][a-zA-Z0-9_-]*$ ]]; then
        error_exit "Invalid project name. Use only letters, numbers, hyphens, and underscores. Must start with letter or number."
    fi
    
    if [[ -e "$name" ]]; then
        error_exit "Directory or file '$name' already exists in current location"
    fi
}

# Project selection menu
DETAILED_VIEW=1

print_help() {
    {
        echo
        echo -e "${PURPLE}How to choose a template:${NC}"
        echo "- Enter a number (1-5) or a name (react, python, go, cli, none)"
        echo "- Type 'd' to toggle detailed descriptions on/off"
        echo "- Type 'h' or '?' to view this help"
        echo "- Type 'q' to quit"
        echo
    } >&2
}

print_template_details() {
    case "$1" in
        react)
            {
                echo "   • Stack: React + TypeScript + Vite"
                echo "   • Includes: tsconfig, vite.config, starter App.tsx, basic CSS"
                echo "   • Good for: Frontend SPAs with fast dev server"
            } >&2
            ;;
        python)
            {
                echo "   • Stack: Python with src/ layout and pyproject.toml"
                echo "   • Includes: package metadata, example module and tests"
                echo "   • Good for: Libraries, CLIs, or services in Python"
            } >&2
            ;;
        go)
            {
                echo "   • Stack: Go modules with conventional cmd/ and pkg/"
                echo "   • Includes: example main under cmd/, Makefile"
                echo "   • Good for: Go services or binaries"
            } >&2
            ;;
        cli)
            {
                echo "   • Stack: Go-based cross-platform CLI"
                echo "   • Includes: Makefile, scripts/build.sh, single entrypoint"
                echo "   • Good for: Distributable command-line tools"
            } >&2
            ;;
        none)
            {
                echo "   • No application template"
                echo "   • Includes only the TGS workflow and AgentOps files"
                echo "   • Good for: Adopting the workflow in an existing codebase"
            } >&2
            ;;
    esac
}

show_project_types() {
    {
        echo
        echo -e "${PURPLE}Available Project Templates:${NC}"
        echo "1) React  — Modern React application with TypeScript (Vite)"
        echo "2) Python — Python project with pyproject.toml and tests"
        echo "3) Go     — Go application using modules and standard layout"
        echo "4) CLI    — Cross-platform CLI tool template (Go)"
        echo "5) None   — TGS workflow + AgentOps only (no app template)"
        echo
        if [[ "$DETAILED_VIEW" -eq 1 ]]; then
            echo -e "${BLUE}Details:${NC}"
            echo "- React:";  print_template_details react
            echo "- Python:"; print_template_details python
            echo "- Go:";     print_template_details go
            echo "- CLI:";    print_template_details cli
            echo "- None:";   print_template_details none
            echo
        fi
        echo "Tip: Enter number or name (e.g., '1' or 'react'). Type 'd' to toggle details, '?' for help."
        echo
    } >&2
}

get_project_selection() {
    while true; do
        show_project_types
        if ! read -r -p "Select template [1-5|react|python|go|cli|none]: " choice </dev/tty; then
            >&2 echo -e "${YELLOW}[WARNING]${NC} No TTY detected for input; defaulting to 'None'"
            echo "none"
            return
        fi
        normalized=$(echo "$choice" | tr '[:upper:]' '[:lower:]' | xargs)
        case $normalized in
            1|react)
                selected="react"
                ;;
            2|python)
                selected="python"
                ;;
            3|go)
                selected="go"
                ;;
            4|cli)
                selected="cli"
                ;;
            5|none|skip|no|n)
                selected="none"
                ;;
            d|detail|details)
                if [[ "$DETAILED_VIEW" -eq 0 ]]; then DETAILED_VIEW=1; else DETAILED_VIEW=0; fi
                {
                    echo
                    if [[ "$DETAILED_VIEW" -eq 1 ]]; then
                        echo -e "${BLUE}[INFO]${NC} Detailed view enabled"
                    else
                        echo -e "${BLUE}[INFO]${NC} Detailed view disabled"
                    fi
                } >&2
                continue
                ;;
            h|help|\?)
                print_help
                continue
                ;;
            q|quit|exit)
                >&2 echo -e "${BLUE}[INFO]${NC} Operation cancelled by user"
                exit 0
                ;;
            *)
                >&2 echo -e "${YELLOW}[WARNING]${NC} Invalid selection. Please choose 1-5 or a valid name."
                continue
                ;;
        esac
        echo "$selected"
        return
    done
}

# Directory naming
get_project_name() {
    local default_name="my-project"
    
    echo >&2
    if ! read -r -p "Enter project name (default: $default_name): " project_name </dev/tty; then
        project_name=""
    fi
    
    if [[ -z "$project_name" ]]; then
        project_name="$default_name"
    fi
    # Sanitize aggressively
    project_name=$(sanitize_name "$project_name")
    if [[ -z "$project_name" ]]; then project_name="$default_name"; fi
    validate_project_name "$project_name"
    echo "$project_name"
}

# Git operations
clone_repository() {
    local target_dir_raw="$1"
    local target_dir
    target_dir=$(sanitize_name "$target_dir_raw")
    
    log_info "Cloning repository..."
    TEMP_DIR=$(mktemp -d)
    
    if ! git clone --quiet "$REPO_URL" "$TEMP_DIR"; then
        error_exit "Failed to clone repository from $REPO_URL"
    fi
    
    log_info "Copying repository contents to $target_dir..."
    mkdir -p "$target_dir"
    # Copy contents (not the directory itself)
    cp -a "$TEMP_DIR"/. "$target_dir"/
    
    # Ensure no git history is carried over
    rm -rf "$target_dir/.git"
}

# Fetch repository into temporary directory only (no copy)
fetch_repository_to_temp() {
    log_info "Fetching repository contents..."
    TEMP_DIR=$(mktemp -d)
    if ! git clone --quiet "$REPO_URL" "$TEMP_DIR"; then
        error_exit "Failed to clone repository from $REPO_URL"
    fi
}

# Decorate mode implementation
decorate_current_directory() {
    log_info "Running decorate mode (minimal TGS workflow)"
    validate_dependencies
    fetch_repository_to_temp

    local src_root="$TEMP_DIR"

    # Essential files
    safe_copy "$src_root/agentops/AGENTOPS.md" "agentops/AGENTOPS.md"
    safe_copy "$src_root/tgs/README.md" "tgs/README.md"
    safe_mkdir_p "agentops/tgs"
    if [[ -d "$src_root/agentops/tgs" ]]; then
        for f in "$src_root"/agentops/tgs/*; do
            local_name=$(basename "$f")
            safe_copy "$f" "agentops/tgs/$local_name"
        done
    fi

    # Optional template overlay
    if [[ "$WITH_TEMPLATES" != "none" ]]; then
        local tmpl_dir="$src_root/templates/$WITH_TEMPLATES"
        if [[ -d "$tmpl_dir" ]]; then
            log_info "Applying template overlay: $WITH_TEMPLATES"
            # Create directories first
            if [[ "$DRY_RUN" -eq 1 ]]; then
                echo -e "${BLUE}[DRY-RUN]${NC} overlay from $tmpl_dir"
            fi
            while IFS= read -r -d '' d; do
                rel_path="${d#"$tmpl_dir/"}"
                safe_mkdir_p "$rel_path"
            done < <(find "$tmpl_dir" -type d -print0)
            while IFS= read -r -d '' f; do
                rel_path="${f#"$tmpl_dir/"}"
                safe_copy "$f" "$rel_path"
            done < <(find "$tmpl_dir" -type f -print0)
        else
            log_warning "Template directory not found: $tmpl_dir (skipping)"
        fi
    fi

    # Generate tgs.mk and ensure Makefile includes it
    write_tgs_mk
    ensure_makefile_includes_tgs

    log_success "Decoration complete"
}

# Prompt when existing repo is detected and not explicitly decorating
prompt_decorate_or_new() {
    if [[ -d ".git" && "$DECORATE" -eq 0 ]]; then
        {
            echo
            echo -e "${YELLOW}[WARNING]${NC} Detected existing git repository in current directory."
            echo "You can either decorate this repository with the minimal TGS workflow,"
            echo "or initialize a new project in a subdirectory."
            echo
        } >&2
        if ! read -r -p "Choose action: [d]ecorate current repo, [n]ew project, [q]uit [d/n/q]: " choice </dev/tty; then
            >&2 echo -e "${YELLOW}[WARNING]${NC} No TTY detected for input; defaulting to 'decorate'"
            DECORATE=1
            return
        fi
        local normalized
        normalized=$(echo "$choice" | tr '[:upper:]' '[:lower:]' | xargs)
        case "$normalized" in
            d|decorate)
                DECORATE=1
                ;;
            n|new)
                : # proceed with new project flow
                ;;
            q|quit|exit)
                >&2 echo -e "${BLUE}[INFO]${NC} Operation cancelled by user"
                exit 0
                ;;
            *)
                >&2 echo -e "${YELLOW}[WARNING]${NC} Invalid selection. Defaulting to 'decorate'"
                DECORATE=1
                ;;
        esac
    fi
}

# Template application
apply_template() {
    local template_type_raw="$1"
    local project_dir_raw="$2"
    local project_dir
    local template_type
    project_dir=$(sanitize_name "$project_dir_raw")
    template_type=$(sanitize_name "$template_type_raw")
    
    if [[ "$template_type" == "none" ]]; then
        log_info "No template selected - keeping base TGS workflow only"
        return
    fi
    
    # Resolve absolute project directory, with safe fallback
    if [[ "$project_dir" != /* ]]; then
        resolved_dir=$(cd "$project_dir" 2>/dev/null && pwd || echo "")
        if [[ -n "$resolved_dir" ]]; then
            project_dir="$resolved_dir"
        fi
    fi
    local template_dir="$project_dir/templates/$template_type"
    
    if [[ ! -d "$template_dir" ]]; then
        log_warning "Template directory not found: $template_dir"
        log_info "Proceeding with base TGS workflow only"
        return
    fi
    
    log_info "Applying $template_type template..."
    
    # Copy template files to project root
    if [[ -d "$template_dir" ]]; then
        cp -r "$template_dir"/* "$project_dir/" 2>/dev/null || true
        log_success "Applied $template_type template"
    fi
}

# Main execution
main() {
    # Parse arguments
    for arg in "$@"; do
        case "$arg" in
            --decorate)
                DECORATE=1
                ;;
            --dry-run)
                DRY_RUN=1
                ;;
            --force)
                FORCE=1
                ;;
            --with-templates=*)
                WITH_TEMPLATES="${arg#*=}"
                ;;
            -h|--help)
                print_usage
                exit 0
                ;;
            *)
                # ignore unknown to preserve backward compatibility
                ;;
        esac
    done

    echo -e "${GREEN}Repository Bootstrap Script v${SCRIPT_VERSION}${NC}"
    echo "This script will help you create a new project with established tooling and workflows."
    echo

    # Decorate mode: operate in current directory and exit
    if [[ "$DECORATE" -eq 1 ]]; then
        decorate_current_directory
        return
    fi

    # Validate environment
    validate_dependencies
    
    # If an existing git repo is detected, prompt user to choose decorate vs new
    prompt_decorate_or_new
    if [[ "$DECORATE" -eq 1 ]]; then
        decorate_current_directory
        return
    fi

    # If dry-run, simulate and exit success
    if [[ "$DRY_RUN" -eq 1 ]]; then
        if [[ -d ".git" ]]; then
            echo -e "${BLUE}[DRY-RUN]${NC} Would detect existing git repo and prompt to choose decorate or new"
        fi
        echo -e "${BLUE}[DRY-RUN]${NC} Would prompt for template selection"
        echo -e "${BLUE}[DRY-RUN]${NC} Would prompt for project name"
        echo -e "${BLUE}[DRY-RUN]${NC} Would clone repository and apply selected template"
        echo -e "${BLUE}[DRY-RUN]${NC} Would initialize a fresh git repository"
        return
    fi

    # Get user input
    project_type=$(get_project_selection)
    project_name=$(get_project_name)
    
    log_info "Project: $project_name"
    log_info "Template: $project_type"
    echo
    
    # Single confirmation
    if ! read -r -p "Proceed with project creation? (Y/n): " confirm </dev/tty; then
        confirm="y"
    fi
    if [[ "$confirm" =~ ^[Nn]$ ]]; then
        log_info "Operation cancelled by user"
        exit 0
    fi
    
    # Execute bootstrap
    clone_repository "$project_name"
    apply_template "$project_type" "$project_name"
    
    # Final setup
    log_info "Finalizing project setup..."
    cd "$project_name"
    
    # Update README with project-specific information
    if [[ -f "README.md" ]]; then
        tmp_readme=$(mktemp)
        awk -v name="$project_name" '{ if ($0=="## Toolings") { print "## " name } else { print $0 } }' README.md > "$tmp_readme" && mv "$tmp_readme" README.md || rm -f "$tmp_readme"
    fi
    
    # Initialize fresh git repository including any template changes
    log_info "Initializing fresh git repository..."
    git init --quiet
    git add .
    git commit --quiet -m "chore(bootstrap): initial project scaffold via bootstrap script"

    echo
    log_success "Project '$project_name' created successfully!"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. cd $project_name"
    echo "2. Review the README.md for project-specific instructions"
    echo "3. Start your first TGS thought: make new-thought title=\"My First Feature\""
    echo "4. Follow the AgentOps workflow in agentops/AGENTOPS.md"
    echo
    echo -e "${PURPLE}Happy coding! 🚀${NC}"
}

# Script entry point
# Handle invocation via 'bash <(curl ...)' where BASH_SOURCE may be unset when set -u is active
if [[ "$0" == "bash" || "$0" == "-bash" ]]; then
    main "$@"
elif [[ "${BASH_SOURCE[0]:-}" == "$0" ]]; then
    main "$@"
fi