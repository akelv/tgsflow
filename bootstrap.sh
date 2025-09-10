#!/bin/bash
set -euo pipefail

# Repository Bootstrap Script
# Clones this repository and sets up a new project with selected template

REPO_URL="https://github.com/akelv/tgsflow.git"
SCRIPT_VERSION="1.0.0"

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
show_project_types() {
    echo
    echo -e "${PURPLE}Available Project Templates:${NC}"
    echo "1) React - Modern React application with TypeScript"
    echo "2) Python - Python project with pyproject.toml and modern tooling"
    echo "3) Go - Go application with modules and basic structure"
    echo "4) CLI - Cross-platform CLI tool template"
    echo "5) None - Just the TGS workflow and AgentOps (no project template)"
    echo
}

get_project_selection() {
    while true; do
        show_project_types
        read -p "Select project type (1-5): " choice
        case $choice in
            1) echo "react"; return ;;
            2) echo "python"; return ;;
            3) echo "go"; return ;;
            4) echo "cli"; return ;;
            5) echo "none"; return ;;
            *) log_warning "Invalid selection. Please choose 1-5." ;;
        esac
    done
}

# Directory naming
get_project_name() {
    local default_name="my-project"
    
    echo
    read -p "Enter project name (default: $default_name): " project_name
    
    if [[ -z "$project_name" ]]; then
        project_name="$default_name"
    fi
    
    validate_project_name "$project_name"
    echo "$project_name"
}

# Git operations
clone_repository() {
    local target_dir="$1"
    
    log_info "Cloning repository..."
    TEMP_DIR=$(mktemp -d)
    
    if ! git clone --quiet "$REPO_URL" "$TEMP_DIR"; then
        error_exit "Failed to clone repository from $REPO_URL"
    fi
    
    log_info "Copying repository contents to $target_dir..."
    cp -r "$TEMP_DIR" "$target_dir"
    
    # Clean up git history and create fresh repository
    log_info "Initializing fresh git repository..."
    rm -rf "$target_dir/.git"
    cd "$target_dir"
    git init --quiet
    git add .
    git commit --quiet -m "Initial commit from bootstrap

🤖 Generated with repository bootstrap script

Co-Authored-By: Bootstrap Script <noreply@bootstrap.local>"
}

# Template application
apply_template() {
    local template_type="$1"
    local project_dir="$2"
    
    if [[ "$template_type" == "none" ]]; then
        log_info "No template selected - keeping base TGS workflow only"
        return
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
    
    # Confirm before proceeding
    read -p "Proceed with project creation? (y/N): " confirm
    if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
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
        sed -i.bak "s/## Toolings/## $project_name/" README.md
        rm -f README.md.bak
    fi
    
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
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi