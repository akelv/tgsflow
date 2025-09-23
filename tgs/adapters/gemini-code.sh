#!/usr/bin/env bash
# tgs/adapters/gemini-code.sh
# Tiny adapter to call Gemini CLI with context files + task prompt.
# Parity goals with claude-code.sh: deterministic context expansion, snapshot hash,
# prompt via stdin, suggestions routing for patches, non-interactive behavior.

set -euo pipefail

# --- Defaults ---
GEMINI_CMD=${GEMINI_CMD:-"gemini"}
RETURN_MODE="patch_or_text"     # patch_or_text | text (kept for parity)
TIMEOUT_SEC=""
OUT_PATH=""
PROMPT_TEXT=""
PROMPT_FILE=""
CONTEXT_LIST_FILE=""
CONTEXT_GLOB=""
SUGGESTIONS_DIR="tgs/suggestions"
VERBOSE=false
WORKDIR="."
SANDBOX_IMAGE=""
APPLY_PATCH=false
CREATE_PR=false
NEW_BRANCH=""
COMMIT_MESSAGE=""

# --- Helpers ---
die() { echo "gemini-code.sh: $*" >&2; exit 2; }
log() { $VERBOSE && echo "gemini-code.sh: $*" >&2 || true; }

have_sha256() {
  if command -v sha256sum >/dev/null 2>&1; then echo "sha256sum"; return 0; fi
  if command -v shasum >/dev/null 2>&1; then echo "shasum -a 256"; return 0; fi
  return 1
}

calc_file_hash() {
  local f=$1; local tool; tool=$(have_sha256) || die "Need sha256sum or shasum"
  local content_hash; content_hash=$($tool "$f" | awk '{print $1}')
  printf "%s\n%s" "$f" "$content_hash" | ($tool | awk '{print $1}')
}

calc_snapshot_hash() {
  local files=("$@"); local tool; tool=$(have_sha256) || die "Need sha256sum or shasum"
  local tmp; tmp=$(mktemp)
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
  grep -qE '^(diff --git |Index: |---[[:space:]]|\+\+\+[[:space:]])' <<<"$1"
}

ensure_dir() {
  local d=$1; [ -d "$d" ] || mkdir -p "$d"
}

apply_patch_robust() {
  # reads patch from STDIN, tries multiple strategies
  if git apply -p1 -v; then return 0; fi
  if git apply -p0 -v; then return 0; fi
  if command -v patch >/dev/null 2>&1; then
    if patch -p1 -s -N; then return 0; fi
    if patch -p0 -s -N; then return 0; fi
  fi
  return 1
}

usage() {
  cat <<'EOF'
Usage: gemini-code.sh [options]

Options:
  --prompt-file PATH          Prompt template/text file to send to Gemini
  --prompt-text TEXT          Prompt text (overrides --prompt-file if both provided)
  --context-list PATH         Newline-separated file list to include as context
  --context-glob "PATTERN"    Shell glob for context files (expanded, sorted)
  --return-mode MODE          patch_or_text | text  (default: patch_or_text)
  --timeout SEC               Seconds timeout for Gemini execution (uses 'timeout' if present)
  --out PATH                  Write output to file; if omitted, prints to stdout
  --suggestions-dir DIR       Where to place .patch/.txt if --out omitted and we detect type (default: tgs/suggestions)
  --gemini-cmd CMD            Gemini CLI command (default: "gemini")
  --workdir PATH              Change to this directory before running (default: current)
  --sandbox-image IMAGE       Run under Docker with IMAGE, mounting PWD to /work
  --apply-patch               If output looks like a patch, apply with git apply
  --create-pr                 Create PR workflow: new branch, apply patch, commit, push, open PR
  --new-branch NAME           Branch name to create when using --create-pr
  --commit-message MSG        Commit message when using --create-pr
  --verbose                   Extra logs to stderr
  -h, --help                  Show this help

Environment overrides:
  GEMINI_CMD                  Same as --gemini-cmd
  RETURN_MODE                 Same as --return-mode
  TIMEOUT_SEC                 Same as --timeout
  OUT_PATH                    Same as --out
  CONTEXT_FILES               Newline-separated file list (alternative to --context-list/glob)
  PROMPT_TEXT                 Same as --prompt-text
  WORKDIR                     Same as --workdir
  SANDBOX_IMAGE               Same as --sandbox-image
  APPLY_PATCH                 Set to true to apply patch output
  CREATE_PR                   Set to true to enable PR workflow
  NEW_BRANCH                  Same as --new-branch
  COMMIT_MESSAGE              Same as --commit-message
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
    --gemini-cmd) GEMINI_CMD=${2:-}; shift 2 ;;
    --claude-cmd) GEMINI_CMD=${2:-}; shift 2 ;;
    --workdir) WORKDIR=${2:-}; shift 2 ;;
    --sandbox-image) SANDBOX_IMAGE=${2:-}; shift 2 ;;
    --apply-patch) APPLY_PATCH=true; shift ;;
    --create-pr) CREATE_PR=true; shift ;;
    --new-branch) NEW_BRANCH=${2:-}; shift 2 ;;
    --commit-message) COMMIT_MESSAGE=${2:-}; shift 2 ;;
    --verbose) VERBOSE=true; shift ;;
    -h|--help) usage; exit 0 ;;
    *) die "Unknown arg: $1 (use -h)";;
  esac
done

# Env fallbacks
RETURN_MODE=${RETURN_MODE:-${RETURN_MODE:-patch_or_text}}
TIMEOUT_SEC=${TIMEOUT_SEC:-${TIMEOUT_SEC:-}}
OUT_PATH=${OUT_PATH:-${OUT_PATH:-}}
GEMINI_CMD=${GEMINI_CMD:-"gemini"}
WORKDIR=${WORKDIR:-${WORKDIR:-"."}}
SANDBOX_IMAGE=${SANDBOX_IMAGE:-${SANDBOX_IMAGE:-}}
APPLY_PATCH=${APPLY_PATCH:-${APPLY_PATCH:-false}}
CREATE_PR=${CREATE_PR:-${CREATE_PR:-false}}
NEW_BRANCH=${NEW_BRANCH:-${NEW_BRANCH:-}}
COMMIT_MESSAGE=${COMMIT_MESSAGE:-${COMMIT_MESSAGE:-}}

# Change working directory if requested
if [ "$WORKDIR" != "." ]; then
  cd "$WORKDIR" || die "Could not change to working directory: $WORKDIR"
fi

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
  shopt -s nullglob
  # shellcheck disable=SC2206,SC2086
  GLOB_EXPANDED=( $CONTEXT_GLOB )
  shopt -u nullglob
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

# Deduplicate while preserving order
declare -a FINAL_CTX=()
for f in "${CONTEXT_FILES_ARR[@]:-}"; do
  found=0
  for g in "${FINAL_CTX[@]:-}"; do
    if [ "$g" = "$f" ]; then found=1; break; fi
  done
  [ $found -eq 1 ] || FINAL_CTX+=("$f")
done

[ ${#FINAL_CTX[@]} -gt 0 ] || die "No context files provided (use --context-list, --context-glob, or CONTEXT_FILES)"

# Verify existence
for f in "${FINAL_CTX[@]:-}"; do
  [ -f "$f" ] || die "Context file missing: $f"
done

# Snapshot hash for reproducibility
SNAP_HASH=$(calc_snapshot_hash "${FINAL_CTX[@]:-}")
log "Context snapshot: CTX-${SNAP_HASH}"

# Compose combined prompt with context file paths
CTX_LIST=""
for f in "${FINAL_CTX[@]:-}"; do
  CTX_LIST+="$f\\n"
done
FINAL_PROMPT="$PROMPT_TEXT\\n\\nContext files (paths):\\n$CTX_LIST"

# Unique directories for potential CLI hints
declare -a ADD_DIRS=()
for f in "${FINAL_CTX[@]:-}"; do
  d=$(dirname "$f")
  found=0
  for g in "${ADD_DIRS[@]:-}"; do if [ "$g" = "$d" ]; then found=1; break; fi; done
  [ $found -eq 1 ] || ADD_DIRS+=("$d")
done

# Build Gemini args; pipe prompt via stdin. Add directory hints if supported.
declare -a ARGS=()
for d in "${ADD_DIRS[@]:-}"; do
  ARGS+=( --include-directories "$d" )
done

# Build command (reads prompt from stdin)
RUN=("$GEMINI_CMD" "${ARGS[@]}")

# Sandbox via Docker if requested
if [ -n "$SANDBOX_IMAGE" ]; then
  RUN=( docker run --rm -i -v "$PWD:/work" -w /work "$SANDBOX_IMAGE" "${RUN[@]}" )
fi

# Add timeout wrapper last
if [ -n "$TIMEOUT_SEC" ] && command -v timeout >/dev/null 2>&1; then
  RUN=( timeout "$TIMEOUT_SEC" "${RUN[@]}" )
fi

log "Invoking: ${RUN[*]}"

if [ "$CREATE_PR" = true ]; then
  [ -n "$NEW_BRANCH" ] || die "--new-branch is required with --create-pr"
  [ -n "$COMMIT_MESSAGE" ] || die "--commit-message is required with --create-pr"
  # Ensure we're in a git repo
  git rev-parse --git-dir >/dev/null 2>&1 || die "Not a git repository"
  # Create and switch to branch
  git checkout -b "$NEW_BRANCH" || die "Could not create branch $NEW_BRANCH"
  set +e
  OUTPUT="$(printf "%s" "$FINAL_PROMPT" | "${RUN[@]}")"
  STATUS=$?
  set -e
  [ $STATUS -eq 0 ] || die "Gemini command failed with exit code $STATUS"
  if is_patch "$OUTPUT"; then
    if ! printf "%s" "$OUTPUT" | apply_patch_robust; then
      die "Failed to apply patch."
    fi
  else
    echo "Output is not a patch. Cannot create PR." >&2
    git checkout - >/dev/null 2>&1 || true
    exit 1
  fi
  git add .
  git commit -m "$COMMIT_MESSAGE" || die "Commit failed"
  git push -u origin "$NEW_BRANCH" || die "Push failed"
  if command -v gh >/dev/null 2>&1; then
    gh pr create --fill --body "PR created by tgsflow gemini-code.sh adapter." || die "Failed to create PR via gh"
  else
    log "gh not found; skipping PR creation"
  fi
  log "Successfully created PR for branch $NEW_BRANCH"
  exit 0
else
  set +e
  OUTPUT="$(printf "%s" "$FINAL_PROMPT" | "${RUN[@]}")"
  STATUS=$?
  set -e
  [ $STATUS -eq 0 ] || die "Gemini command failed with exit code $STATUS"

  if [ -n "$OUT_PATH" ]; then
    printf "%s" "$OUTPUT" > "$OUT_PATH"
    log "Wrote output to $OUT_PATH"
    exit 0
  fi

  if is_patch "$OUTPUT"; then
    if [ "$APPLY_PATCH" = true ]; then
      if ! printf "%s" "$OUTPUT" | apply_patch_robust; then
        die "Failed to apply patch."
      fi
      log "Patch applied successfully."
      exit 0
    fi
    ensure_dir "$SUGGESTIONS_DIR"
    TS=$(date +"%Y%m%d_%H%M%S")
    DEST="$SUGGESTIONS_DIR/CTX-${SNAP_HASH:0:8}_$TS.patch"
    printf "%s" "$OUTPUT" > "$DEST"
    echo "$DEST"
    exit 0
  fi

  printf "%s" "$OUTPUT"
  exit 0
fi


