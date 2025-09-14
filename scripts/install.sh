#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-akelv/tgsflow}"
BINARY="${BINARY:-tgs}"
DRY_RUN="${DRY_RUN:-}" # if set to non-empty, skip download/extract/install and just print info

# Detect OS (allow override)
uname_s=$(uname -s | tr '[:upper:]' '[:lower:]')
case "${OS:-$uname_s}" in
  linux)   OS="linux" ;;
  darwin)  OS="darwin" ;;
  msys*|mingw*|cygwin*|windows) OS="windows" ;;
  *) echo "Unsupported OS: ${OS:-$uname_s}" >&2; exit 1 ;;
 esac

# Detect ARCH (allow override)
uname_m=$(uname -m)
case "${ARCH:-$uname_m}" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported ARCH: ${ARCH:-$uname_m}" >&2; exit 1 ;;
 esac

# Determine tag (allow override)
if [[ -n "${TAG:-}" ]]; then
  TAG="${TAG}"
else
  TAG=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name"\s*:\s*"([^"]+)".*/\1/')
fi
if [[ -z "${TAG}" ]]; then
  if [[ -n "$DRY_RUN" ]]; then
    echo "Note: Could not determine latest release tag (DRY_RUN mode)."
  else
    echo "Failed to determine latest release tag" >&2
    exit 1
  fi
fi

echo "Installing ${BINARY} ${TAG:-unknown} for ${OS}/${ARCH}..."

ASSET_BASE="${BINARY}_${TAG:-unknown}_${OS}_${ARCH}"
if [[ "${OS}" == "windows" ]]; then
  ASSET_FILE="${ASSET_BASE}.zip"
else
  ASSET_FILE="${ASSET_BASE}.tar.gz"
fi

URL="https://github.com/${REPO}/releases/download/${TAG}/${ASSET_FILE}"
INSTALL_DIR_DEFAULT="/usr/local/bin"
INSTALL_DIR="${INSTALL_DIR:-$INSTALL_DIR_DEFAULT}"

echo "Download URL: $URL"
echo "Install dir: $INSTALL_DIR"

if [[ -n "$DRY_RUN" ]]; then
  echo "DRY_RUN set; exiting before download/install"
  exit 0
fi

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT
cd "$TMPDIR"

curl -fL --retry 3 --connect-timeout 15 -o asset "$URL"

if [[ "${OS}" == "windows" ]]; then
  echo "Windows is not supported by this installer yet." >&2
  exit 1
else
  tar -xzf asset
fi

# Determine install dir
if [[ ! -w "$INSTALL_DIR" ]]; then
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
  echo "Installing to $INSTALL_DIR (add to your PATH if needed)"
fi

install -m 0755 "$BINARY" "$INSTALL_DIR/$BINARY"

echo "Installed $BINARY to $INSTALL_DIR"
"$INSTALL_DIR/$BINARY" --version || true
