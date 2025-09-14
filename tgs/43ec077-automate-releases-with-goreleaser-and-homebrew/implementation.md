# Implementation: Automate releases with GoReleaser and Homebrew

## What & Why
Adds automated releases for the `tgs` CLI using GoReleaser and GitHub Actions, plus a Homebrew tap and a portable curl installer. Enables new developers to install via Homebrew or curl and ensures reproducible, versioned binaries across platforms.

## Changes
- Added GoReleaser configuration for multi-OS/arch builds, checksums, changelog, and Homebrew tap updates:
  - `.goreleaser.yml`
- Added GitHub Actions workflows:
  - `.github/workflows/release.yml` — runs on `v*` tags to publish a GitHub Release
  - `.github/workflows/ci.yml` — bumped Go to 1.23 and added a snapshot GoReleaser validation step on PRs
- Added installer script:
  - `scripts/install.sh` — downloads latest GitHub Release asset and installs to PATH
- Updated docs with install instructions:
  - `README.md` — brew tap/install and curl installer

## How Versioning Works
- `src/main.go` exposes `version`, `commit`, `date` variables.
- GoReleaser provides them via `-ldflags` using tag, commit, and build date.
- `tgs --version` prints values from release builds.

## Homebrew Tap
- GoReleaser brews config targets the tap repo `akelv/homebrew-tgs` and manages `Formula/tgs.rb`.
- Install via:
  - `brew tap akelv/tgs`
  - `brew install tgs`

## Required Secrets & Setup
- In GitHub repo settings → Secrets and variables → Actions, add:
  - `GORELEASER_GITHUB_TOKEN`: A token with `repo` access to `akelv/homebrew-tgs` (tap) so GoReleaser can push formula updates.
- Ensure the tap repository exists: `github.com/akelv/homebrew-tgs` with default branch, empty ok.

## How to Release
1. Merge changes to `main`.
2. Create a semver tag and push:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```
3. GitHub Actions will:
   - Build and publish artifacts for linux/darwin/windows (amd64, arm64)
   - Generate checksums and changelog
   - Create a GitHub Release
   - Update the Homebrew formula in `akelv/homebrew-tgs`

## How to Test (Pre-release)
- Push a pre-release tag (e.g., `v0.1.0-rc.1`) and verify:
  - Release artifacts exist and `tgs --version` shows tag/commit/date
  - Homebrew tap updated; `brew install akelv/tgs/tgs` works
  - Curl installer works on macOS (arm64) and Linux (amd64)

## Rollback
- Delete the Release and tag in GitHub if needed.
- Revert any undesired commit in `akelv/homebrew-tgs` (GoReleaser bot commit).

## Follow-ups / Next Steps
- Optional: Add artifact signing (cosign) and provenance/SBOM.
- Optional: Add Windows MSI packaging.
- Optional: Add `brew uninstall` note and troubleshooting to docs.

## Links
- GoReleaser: https://goreleaser.com/
- Homebrew Tap docs: https://goreleaser.com/customization/homebrew/
- Release workflow: `.github/workflows/release.yml`
