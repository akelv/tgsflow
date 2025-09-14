# Plan: Automate releases with GoReleaser and Homebrew

## Objectives
- Build, test, version, and publish `tgs` binaries to GitHub Releases on semver tag pushes.
- Provide easy install via curl and Homebrew.

## Scope
- In: GoReleaser config, GitHub Actions release workflow, Homebrew tap automation, curl installer, docs.
- Out: Code signing/notarization, snap/apt/rpm packages, Windows MSI.

## Acceptance Criteria
- Pushing `vX.Y.Z` tag triggers a GitHub Release with multi-OS/arch assets and checksums.
- `tgs --version` prints injected version/commit/date from release builds.
- Homebrew tap updated with a working `brew install` command.
- Curl installer downloads latest release and installs binary to PATH.

## Phased Tasks
1) GoReleaser config
- Add `.goreleaser.yml` with multi-platform builds from `./src` main, `ldflags` for versioning, archives, checksums, changelog.

2) Release workflow
- Add `.github/workflows/release.yml` triggered on tag `v*` that runs `goreleaser release --clean` using Go 1.23.x.

3) Homebrew tap
- Configure `brews` in GoReleaser targeting `akelv/homebrew-tgs`, with `tgs` formula.
- Use `GORELEASER_GITHUB_TOKEN` for tap updates.

4) Curl installer
- Add `scripts/install.sh` to fetch latest release asset for the OS/arch and install to `/usr/local/bin` or `$HOME/.local/bin`.
- Document usage in `README.md`.

5) Docs and CI tweaks
- Update `README.md` with install instructions (curl and brew).
- Ensure existing CI remains intact.

## File-by-file Changes
- `.goreleaser.yml`: New file with build, release, brews config.
- `.github/workflows/release.yml`: New workflow for tag-based releases.
- `scripts/install.sh`: New portable curl/bash installer.
- `README.md`: Add install section and versioning note.

## Version Injection
- `src/main.go` already exposes `version`, `commit`, `date` variables. Pass via `-ldflags` using GoReleaser template values.

## Test Plan
- CI: On PRs, run `goreleaser build --snapshot --skip=publish --clean` to validate config.
- Dry-run release: Create a pre-release tag `v0.0.0-rc.1` and verify artifacts, version string, Homebrew tap PR/commit, and curl installer works on macOS (arm64) and Linux (amd64) runners.
- `tgs --version` output includes correct values.

## Rollout / Rollback
- Rollout: Merge to main, push a `vX.Y.Z` tag.
- Rollback: Delete the GitHub Release and tag if necessary; revert formula update via tap repo.