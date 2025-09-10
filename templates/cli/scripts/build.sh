#!/bin/bash
set -euo pipefail

# Cross-platform build script for CLI tools
# Supports Go, Rust, and Node.js projects

APP_NAME="${APP_NAME:-my-cli-tool}"
VERSION="${VERSION:-0.1.0}"
BUILD_DIR="${BUILD_DIR:-dist}"

# Platform targets
PLATFORMS=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
)

log_info() {
    echo "ℹ️  $1"
}

log_success() {
    echo "✅ $1"
}

log_error() {
    echo "❌ $1" >&2
}

build_go() {
    local platform="$1"
    local goos="${platform%/*}"
    local goarch="${platform#*/}"
    local output="${BUILD_DIR}/${APP_NAME}-${goos}-${goarch}"
    
    if [[ "$goos" == "windows" ]]; then
        output="${output}.exe"
    fi
    
    log_info "Building for $platform..."
    GOOS="$goos" GOARCH="$goarch" go build \
        -ldflags "-X main.version=$VERSION" \
        -o "$output" \
        ./src
    
    log_success "Built $output"
}

build_rust() {
    log_info "Building Rust project..."
    cargo build --release
    cp target/release/"$APP_NAME" "$BUILD_DIR/"
    log_success "Built $BUILD_DIR/$APP_NAME"
}

build_node() {
    log_info "Building Node.js project..."
    npm run build
    log_success "Built Node.js project"
}

main() {
    mkdir -p "$BUILD_DIR"
    
    if [[ -f "go.mod" ]]; then
        log_info "Detected Go project"
        for platform in "${PLATFORMS[@]}"; do
            build_go "$platform"
        done
    elif [[ -f "Cargo.toml" ]]; then
        log_info "Detected Rust project"
        build_rust
    elif [[ -f "package.json" ]]; then
        log_info "Detected Node.js project"
        build_node
    else
        log_error "No supported project type detected (go.mod, Cargo.toml, or package.json)"
        exit 1
    fi
    
    log_success "Build complete! Artifacts in $BUILD_DIR/"
}

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi