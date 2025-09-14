#!/usr/bin/env bash
# tgs/adapters/claude-code.sh
# Tiny adapter to call Claude Code with context files + task prompt.
# - Builds an ordered, hashed "context snapshot" for reproducibility
# - Passes files via repeated --context flags
# - Supports return mode: patch-or-text | text
# - Writes to stdout or a target file; can auto-route patches to tgs/suggestions/
# - Exits non-zero on obvious misconfig (missing files, empty prompt)

set -euo pipefail

# --- Defaults ---
CLAUDE_CMD=${CLAUDE_CMD:-"claude"}
RETURN_MODE="patch_or_text"     # patch_or_text | text
TIMEOUT_SEC=""                  # empty = no timeout wrapper
OUT_PATH=""                     # if empty → stdout
PROMPT_TEXT=""                  # direct text (highest precedence)
PROMPT_FILE=""                  # file path (fallback if PROMPT_TEXT empty)
CONTEXT_LIST_FILE=""            # newline-separated list of files
CONTEXT_GLOB=""                 # optional glob (expanded deterministically)
SUGGESTIONS_DIR="tgs/suggestions"
VERBOSE=false

# --- Helpers ---
die() { echo "claude-code.sh: $*" >&2; exit 2; }
log() { $VERBOSE && echo "claude-code.sh: $*" >&2 || true; }

have_sha256() {
  if command -v sha256sum >/dev/null 2>&1; then echo "sha256sum"; return 0; fi
  if command -v shasum >/dev/null 2>&1; then echo "shasum -a 256"; return 0; fi
  return 1
}

calc_file_hash() {
  local f=$1; local tool; tool=$(have_sha256) || die "Need sha256sum or shasum"
  # hash of "path\n<sha-of-content>"
  local content_hash; content_hash=$($tool "$f" | awk '{print $1}')
  printf "%s\n%s" "$f" "$content_hash" | ($tool | awk '{print $1}')
}

calc_snapshot_hash() {
  local files=("$@"); local tool; tool=$(have_sha256) || die "Need sha256sum or shasum"
  # stable: path + content hash per file, then hash the concatenation
  local tmp
  tmp=$(mktemp)
  for f in "${files[@]}"; do
    [ -f "$f" ] || die "Context file missing: $f"
    local h; h=$(calc_file_hash "$f")
    printf "%s %s\n" "$h" "$f" >> "$tmp"
  done
  local snap; snap=$($tool "$tmp" | awk '{print $1}')
  rm -f "$tmp"
  echo "$snap"
}

is_patch() {
  # Heuristic: unified diff starts with 'diff --git ' or 'Index: ' etc.
  grep -qE '^(diff --git |Index: |---[[:space:]]|+++[[:space:]])' <<<"$1"
}

ensure_dir() {
  local d=$1; [ -d "$d" ] || mkdir -p "$d"
}

# --- Usage ---
usage() {
  cat <<'EOF'
Usage: claude-code.sh [options]

Options:
  --prompt-file PATH          Prompt template/text file to send to Claude
  --prompt-text TEXT          Prompt text (overrides --prompt-file if both provided)
  --context-list PATH         Newline-separated file list to include as --context
  --context-glob "PATTERN"    Shell glob for context files (will be expanded, sorted)
  --return-mode MODE          patch_or_text | text  (default: patch_or_text)
  --timeout SEC               Seconds timeout for Claude execution (uses 'timeout' if present)
  --out PATH                  Write output to file; if omitted, prints to stdout
  --suggestions-dir DIR       Where to place .patch/.txt if --out omitted and we detect type (default: tgs/suggestions)
  --claude-cmd CMD            Claude CLI command (default: "claude")
  --verbose                   Extra logs to stderr
  -h, --help                  Show this help

Environment overrides:
  CLAUDE_CMD                  Same as --claude-cmd
  RETURN_MODE                 Same as --return-mode
  TIMEOUT_SEC                 Same as --timeout
  OUT_PATH                    Same as --out
  CONTEXT_FILES               Newline-separated file list (alternative to --context-list/glob)
  PROMPT_TEXT                 Same as --prompt-text
EOF
}

# --- Parse args ---
while [[ $# -gt 0 ]]; do
  case "$1" in
    --prompt-file) PROMPT_FILE=${2:-}; shift 2 ;;
    --prompt-text) PROMPT_TEXT=${2:-}; shift 2 ;;
    --context-list) CONTEXT_LIST_FILE=${2:-}; shift 2 ;;
    --context-glob) CONTEXT_GLOB=${2:-}; shift 2 ;;
    --return-mode) RETURN_MODE=${2:-}; shift 2 ;;
    --timeout) TIMEOUT_SEC=${2:-}; shift 2 ;;
    --out) OUT_PATH=${2:-}; shift 2 ;;
    --suggestions-dir) SUGGESTIONS_DIR=${2:-}; shift 2 ;;
    --claude-cmd) CLAUDE_CMD=${2:-}; shift 2 ;;
    --verbose) VERBOSE=true; shift ;;
    -h|--help) usage; exit 0 ;;
    *) die "Unknown arg: $1 (use -h)";;
  esac
done

# Env fallbacks
RETURN_MODE=${RETURN_MODE:-${RETURN_MODE:-patch_or_text}}
TIMEOUT_SEC=${TIMEOUT_SEC:-${TIMEOUT_SEC:-}}
OUT_PATH=${OUT_PATH:-${OUT_PATH:-}}
CLAUDE_CMD=${CLAUDE_CMD:-"claude"}

# --- Validate basics ---
if [ -z "$PROMPT_TEXT" ]; then
  [ -n "$PROMPT_FILE" ] || die "Provide --prompt-text or --prompt-file"
  [ -f "$PROMPT_FILE" ] || die "Prompt file not found: $PROMPT_FILE"
  PROMPT_TEXT="$(cat "$PROMPT_FILE")"
fi
[ -n "$PROMPT_TEXT" ] || die "Prompt is empty"

# Gather context files (ordered deterministically)
declare -a CONTEXT_FILES_ARR=()
if [ -n "${CONTEXT_FILES:-}" ]; then
  # From env: newline-separated
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    CONTEXT_FILES_ARR+=("$line")
  done <<< "$CONTEXT_FILES"
fi
if [ -n "$CONTEXT_LIST_FILE" ]; then
  [ -f "$CONTEXT_LIST_FILE" ] || die "Context list file not found: $CONTEXT_LIST_FILE"
  while IFS= read -r line; do
    [ -z "$line" ] && continue
    CONTEXT_FILES_ARR+=("$line")
  done < <(grep -v '^[[:space:]]*$' "$CONTEXT_LIST_FILE")
fi
if [ -n "$CONTEXT_GLOB" ]; then
  # Expand and sort
  # shellcheck disable=SC2206
  GLOB_EXPANDED=($CONTEXT_GLOB)
  if [ ${#GLOB_EXPANDED[@]} -gt 0 ]; then
    tmp_glob="$(mktemp)"
    for g in "${GLOB_EXPANDED[@]}"; do
      printf "%s\n" "$g"
    done | LC_ALL=C sort > "$tmp_glob"
    while IFS= read -r line; do
      [ -z "$line" ] && continue
      CONTEXT_FILES_ARR+=("$line")
    done < "$tmp_glob"
    rm -f "$tmp_glob"
  fi
fi

# Deduplicate while preserving order (bash 3 compatible)
declare -a FINAL_CTX=()
for f in "${CONTEXT_FILES_ARR[@]:-}"; do
  found=0
  for g in "${FINAL_CTX[@]:-}"; do
    if [ "$g" = "$f" ]; then
      found=1
      break
    fi
  done
  if [ $found -eq 1 ]; then
    continue
  fi
  FINAL_CTX+=("$f")
done

[ ${#FINAL_CTX[@]:-0} -gt 0 ] || die "No context files provided (use --context-list, --context-glob, or CONTEXT_FILES)"

# Verify existence
for f in "${FINAL_CTX[@]:-}"; do
  [ -f "$f" ] || die "Context file missing: $f"
done

# Snapshot hash for reproducibility
SNAP_HASH=$(calc_snapshot_hash "${FINAL_CTX[@]:-}")
log "Context snapshot: CTX-${SNAP_HASH}"

# --- Build Claude args (non-interactive print) ---
# Compose a combined prompt that includes the list of context file paths
CTX_LIST=""
for f in "${FINAL_CTX[@]:-}"; do
  CTX_LIST+="$f\n"
done
FINAL_PROMPT="$PROMPT_TEXT\n\nContext files (paths):\n$CTX_LIST"

# Unique directories for --add-dir
declare -a ADD_DIRS=()
for f in "${FINAL_CTX[@]:-}"; do
  d=$(dirname "$f")
  found=0
  for g in "${ADD_DIRS[@]:-}"; do
    if [ "$g" = "$d" ]; then found=1; break; fi
  done
  [ $found -eq 1 ] || ADD_DIRS+=("$d")
done

declare -a ARGS=( -p --output-format text )
for d in "${ADD_DIRS[@]:-}"; do
  ARGS+=( --add-dir "$d" )
done

# Build command (reads prompt from stdin)
RUN=("$CLAUDE_CMD" "${ARGS[@]}")
if [ -n "$TIMEOUT_SEC" ] && command -v timeout >/dev/null 2>&1; then
  RUN=( timeout "$TIMEOUT_SEC" "$CLAUDE_CMD" "${ARGS[@]}" )
fi

log "Invoking: ${RUN[*]}"
set +e
OUTPUT="$(printf "%s" "$FINAL_PROMPT" | "${RUN[@]}")"
STATUS=$?
set -e

[ $STATUS -eq 0 ] || die "Claude command failed with exit code $STATUS"

# Decide destination
if [ -n "$OUT_PATH" ]; then
  # Honor explicit path
  printf "%s" "$OUTPUT" > "$OUT_PATH"
  log "Wrote output to $OUT_PATH"
  exit 0
fi

# Auto-route to suggestions dir if looks like a patch
if is_patch "$OUTPUT"; then
  ensure_dir "$SUGGESTIONS_DIR"
  TS=$(date +"%Y%m%d_%H%M%S")
  DEST="$SUGGESTIONS_DIR/CTX-${SNAP_HASH:0:8}_$TS.patch"
  printf "%s" "$OUTPUT" > "$DEST"
  echo "$DEST"
  exit 0
fi

# Otherwise print text to stdout
printf "%s" "$OUTPUT"
exit 0
