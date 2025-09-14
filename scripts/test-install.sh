#!/usr/bin/env bash
set -euo pipefail

: "${REPO:=akelv/tgsflow}"
: "${BINARY:=tgs}"

# Optional: limit to a specific tag via TAG env

check_combo() {
  local os="$1"
  local arch="$2"
  echo "-- Testing ${os}/${arch} --"
  OS="$os" ARCH="$arch" DRY_RUN=1 REPO="$REPO" BINARY="$BINARY" TAG="${TAG:-}" bash "$(dirname "$0")/install.sh" || true

  # Build expected URL and HEAD it if TAG is set (skips GitHub API call spam)
  local tag="${TAG:-}"
  if [[ -z "$tag" ]]; then
    tag=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep -m1 '"tag_name"' | sed -E 's/.*"tag_name"\s*:\s*"([^"]+)".*/\1/') || true
  fi
  if [[ -z "$tag" ]]; then
    echo "Note: no release tag found; skipping HEAD checks for ${os}/${arch}." >&2
    return 0
  fi
  local base="${BINARY}_${tag}_${os}_${arch}"
  local file="$base.tar.gz"
  if [[ "$os" == "windows" ]]; then file="$base.zip"; fi
  local url="https://github.com/${REPO}/releases/download/${tag}/${file}"
  echo "HEAD $url"
  if ! curl -sI "$url" | head -n 1 | grep -qE "HTTP/[0-9.]+ 200"; then
    echo "WARN: asset not reachable (might be expected if release not published for $os/$arch)" >&2
  else
    echo "OK: asset reachable"
  fi
}

# Test common combos
check_combo darwin arm64
check_combo darwin amd64
check_combo linux arm64
check_combo linux amd64
# windows is not supported by installer; skip

echo "Done"
