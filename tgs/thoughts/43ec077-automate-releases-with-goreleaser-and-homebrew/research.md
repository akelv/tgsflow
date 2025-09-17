# Research: Automate releases with GoReleaser and Homebrew

## Problem
We need an automated, repeatable release pipeline that builds, tests, versions, and publishes cross-platform `tgs` binaries. New developers should be able to install `tgs` via a single-line curl script or Homebrew tap. Releases must be traceable and align with the TGS approval-gated workflow.

## Current State
- CI exists for build/test on push and PR (`.github/workflows/ci.yml`).
- No GoReleaser configuration.
- No GitHub Actions workflow for releases.
- No versioning strategy wired to tags.
- No Homebrew tap or formula automation.
- `make build/test/tidy` targets exist; binary builds to `./bin/tgs`.

## Goals / Requirements
- Build multi-OS/arch artifacts: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64.
- Embed version, commit, and date via `-ldflags` from git tags.
- Create checksums and release notes.
- Publish GitHub Release on semver tag push (`vX.Y.Z`).
- Optionally create pre-releases for `-rc`, `-beta` tags.
- Provide installation:
  - curl installer pointing to latest release binary.
  - Homebrew install via tap (auto-updated formula on release).

## Constraints
- Use GitHub Actions; avoid self-hosted runners.
- Use Go 1.23.x (project toolchain uses 1.23.6).
- Keep secrets minimal; Homebrew tap publish requires a token with repo access if using a separate tap repo.
- Keep implementation outside `tgs/` per workflow; only docs live here.

## Risks / Security
- Leaking tokens in CI: scope a dedicated token for Homebrew tap publishing if needed.
- Incorrect version injection: ensure GoReleaser uses `main.version` variable consistently.
- Binary notarization/signing: out of scope initially; can add cosign in a follow-up.
- Formula drift if manual: mitigate by using GoReleaser’s `brew` pipe to update tap automatically.

## Alternatives Considered
- Manual `goreleaser` run from local: error-prone, not reproducible.
- GitHub Actions without GoReleaser: more YAML and bespoke scripting; harder to maintain.
- Using GitHub Packages instead of Releases: reduces visibility for curl installs.

## Recommendation
Adopt GoReleaser with two workflows:
- CI (keep existing) for pull requests and branch pushes.
- Release workflow triggered on tag push `v*` using Go 1.23.x that:
  - Checks out code with full tags history.
  - Sets up Go and caches modules.
  - Runs `goreleaser release --clean`.

Create `goreleaser.yaml` to:
- Configure builds for common OS/arch.
- Set `main` package to `./src` with ldflags for `version`, `commit`, and `date`.
- Name artifact `tgs` with OS/arch suffixes.
- Generate checksums and release notes.
- Publish to GitHub Releases.
- Configure `brews` to update a Homebrew tap repo `akelv/homebrew-tgs` (configurable), creating/maintaining `tgs.rb` with `install` linking the `tgs` binary.

Provide a curl installer script (simple shim) that fetches the latest release and installs to `/usr/local/bin` or `$HOME/.local/bin` depending on OS.

## References
- GoReleaser Quick Start: https://goreleaser.com/quick-start/
- GoReleaser Homebrew Tap: https://goreleaser.com/customization/homebrew/
- Actions setup-go: https://github.com/actions/setup-go
- Example release workflow: https://goreleaser.com/ci/actions/
