# Project Templates

This directory contains project templates used by the bootstrap script. Each template provides a starting structure for different types of projects.

## Available Templates

- **react/** - Modern React application with TypeScript, Vite, and common tooling
- **python/** - Python project with pyproject.toml, modern packaging, and development tools
- **go/** - Go application with modules, basic structure, and build configuration
- **cli/** - Cross-platform CLI tool template with build scripts

## Template Structure

Each template directory contains:
- Project configuration files (package.json, pyproject.toml, go.mod, etc.)
- Basic project structure and example code
- Development tooling configuration
- README.md with template-specific instructions

## Usage

Templates are automatically applied by the bootstrap script based on user selection. The bootstrap script copies template files to the project root and preserves the TGS workflow structure.

## Adding New Templates

1. Create a new directory under `templates/`
2. Add necessary configuration and source files
3. Include a README.md explaining the template
4. Update the bootstrap script to include the new template option