#!/bin/bash
set -euo pipefail

# Repository Bootstrap Script
# Clones this repository and sets up a new project with selected template

REPO_URL="https://github.com/akelv/tgsflow.git"
SCRIPT_VERSION="1.1.1"

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
    echo -e "${GREEN}Repository Bootstrap Script v${SCRIPT_VERSION}${NC}"
    echo "This script will help you create a new project with established tooling and workflows."
    echo

    # Validate environment
    validate_dependencies
    
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