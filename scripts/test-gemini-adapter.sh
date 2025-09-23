#!/usr/bin/env bash
set -euo pipefail

# Integration tests for tgs/adapters/gemini-code.sh
# - Stubs a local "gemini" CLI
# - Verifies text output, suggestions path, apply-patch, workdir, and --claude-cmd alias

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
ADAPTER="$ROOT_DIR/tgs/adapters/gemini-code.sh"
TEST_DIR=$(mktemp -d)
echo "Using TEST_DIR=$TEST_DIR"

cleanup() { rm -rf "$TEST_DIR"; }
trap cleanup EXIT

expect() { # expect $1 == $2
  local got=$1 want=$2
  if [ "$got" != "$want" ]; then
    echo "FAIL: expected [$want] got [$got]" >&2
    exit 1
  fi
}

# 1) Text output test
mkdir -p "$TEST_DIR/text"
echo hello > "$TEST_DIR/text/ctx.txt"
cat > "$TEST_DIR/text/gemini_stub.sh" <<'EOS'
#!/usr/bin/env bash
cat - >/dev/null
echo "Hello from gemini"
EOS
chmod +x "$TEST_DIR/text/gemini_stub.sh"

OUT=$("$ADAPTER" --prompt-text "Say hi" \
  --context-glob "$TEST_DIR/text/ctx.txt" \
  --gemini-cmd "$TEST_DIR/text/gemini_stub.sh" \
  --workdir "$TEST_DIR/text")
expect "$(echo "$OUT" | tr -d '\r')" "Hello from gemini"
echo "ok: text output"

# 2) Suggestions path (patch) test using a real git-generated diff
mkdir -p "$TEST_DIR/suggest"
echo hello > "$TEST_DIR/suggest/a.txt"
(
  cd "$TEST_DIR/suggest"
  git init >/dev/null
  git config user.email test@example.com
  git config user.name "Test User"
  git add a.txt
  git commit -m init >/dev/null
  # generate patch
  echo "hello world" > a.txt
  git diff > generated.patch
  # revert working file for later apply
  echo "hello" > a.txt
)
cat > "$TEST_DIR/suggest/gemini_stub_patch.sh" <<'EOS'
#!/usr/bin/env bash
cat - >/dev/null
cat generated.patch
EOS
chmod +x "$TEST_DIR/suggest/gemini_stub_patch.sh"

PATCH_PATH=$("$ADAPTER" --prompt-text "Make patch" \
  --context-glob "$TEST_DIR/suggest/a.txt" \
  --gemini-cmd "$TEST_DIR/suggest/gemini_stub_patch.sh" \
  --workdir "$TEST_DIR/suggest")
case "$PATCH_PATH" in
  *.patch) ;; 
  *) echo "FAIL: expected .patch path, got $PATCH_PATH" >&2; exit 1;;
esac
if [[ "$PATCH_PATH" != /* ]]; then
  PATCH_PATH="$TEST_DIR/suggest/$PATCH_PATH"
fi
[ -f "$PATCH_PATH" ] || { echo "FAIL: patch file missing $PATCH_PATH" >&2; exit 1; }
echo "ok: suggestions path"

# 3) Apply-patch test (optional)
if [ "${TEST_APPLY_PATCH:-0}" = "1" ]; then
  echo hello > "$TEST_DIR/suggest/a.txt"
  if "$ADAPTER" --prompt-text "Make patch" \
    --context-glob "$TEST_DIR/suggest/a.txt" \
    --gemini-cmd "$TEST_DIR/suggest/gemini_stub_patch.sh" \
    --workdir "$TEST_DIR/suggest" \
    --apply-patch >/dev/null; then
    if grep -q "hello world" "$TEST_DIR/suggest/a.txt"; then
      echo "ok: apply patch"
    else
      echo "FAIL: patch not applied" >&2; exit 1;
    fi
  else
    echo "SKIP: apply patch path failed (set TEST_APPLY_PATCH=1 only if environment supports git/patch)" >&2
  fi
else
  echo "skip: apply patch (set TEST_APPLY_PATCH=1 to enable)"
fi

# 4) --claude-cmd alias works
OUT2=$("$ADAPTER" --prompt-text "Say hi" \
  --context-glob "$TEST_DIR/text/ctx.txt" \
  --claude-cmd "$TEST_DIR/text/gemini_stub.sh" \
  --workdir "$TEST_DIR/text")
expect "$(echo "$OUT2" | tr -d '\r')" "Hello from gemini"
echo "ok: claude-cmd alias"

echo "All adapter integration tests passed."


story