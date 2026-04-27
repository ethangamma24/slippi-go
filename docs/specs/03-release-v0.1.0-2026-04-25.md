# Spec 03: Release `v0.1.0` (`slippi-go`)

**Repository:** [`github.com/ethangamma24/slippi-go`](https://github.com/ethangamma24/slippi-go)
**Scope:** Ship the first public, semver-versioned release of `slippi-go` after specs 01 and 02 land. Update README, add a CHANGELOG, tag `v0.1.0`, and verify resolution from the Go module proxy so downstream consumers can `go get` it.

---

## Implementation Status

- **Status:** 📝 Not Started
- **Completed:** —
- **Implemented by:** —
- **Notes:** —

---

## Spec Dependencies

**Depends on:**
- [`01-public-api-promotion-2026-04-25.md`](./01-public-api-promotion-2026-04-25.md) — types/melee/stats moved to `pkg/`.
- [`02-typed-accessors-and-bytes-constructor-2026-04-25.md`](./02-typed-accessors-and-bytes-constructor-2026-04-25.md) — typed accessors and bytes constructors are present.

**Blocks:** Backend migration spec at `dontgethit/backend/docs/specs/slippi-go-parser-migration-2026-04-25.md`. The backend will pin against `github.com/ethangamma24/slippi-go@v0.1.0`.

---

## Context

`slippi-go` has no published release tags today. The README walks through the
tagging procedure (`git tag v1.0.0 && git push origin v1.0.0`) but does not
commit to a version. The DontGetHit backend cannot pin against a stable version
until something is tagged.

Before tagging, document the public API surface introduced by specs 01 and 02
(README) and the breaking renames within those specs (CHANGELOG). After this
spec lands, the public API surface becomes a stability commitment and any
follow-up changes should ship as `v0.x.y` (additive) or `v0.(x+1).0` /
`v1.0.0` (breaking).

---

## Goals

1. A `CHANGELOG.md` at repo root documenting `v0.1.0` — including the
   `LRASInitiatior` → `LRASInitiator` rename and the `(any, error)` JSON-key
   knock-on it implies.
2. A README that describes the public API: `pkg/slippi.Game` (untyped facade),
   `pkg/slippi.NewGameFromBytes` / `NewGameFromReader`, `ParseMetaFromBytes`,
   typed accessors on `Game`, `pkg/slippi/types`, `pkg/slippi/melee`,
   `pkg/slippi/stats`.
3. `go.mod` (and any `go.sum` updates) committed.
4. A `v0.1.0` tag pushed.
5. Verification that `go list -m github.com/ethangamma24/slippi-go@v0.1.0`
   resolves through `proxy.golang.org`.

## Non-goals

- No further code changes. This spec is **only** docs + tag + verification.

---

## Required changes

### 3.1 `CHANGELOG.md`

Create `CHANGELOG.md` at repo root using
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) shape. Suggested
contents for this release:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

## [v0.1.0] — 2026-MM-DD

First semver-tagged release. The public API surface is now considered stable
within the v0.x line, with breaking changes reserved for future v0.(x+1).0 /
v1.0.0 releases.

### Added

- Public types under `github.com/ethangamma24/slippi-go/pkg/slippi/types` covering
  the entire decoded Slippi game (`Game`, `Data`, `GameStart`, `Player`, `GameEnd`,
  `PlayerPlacement`, `Frame`, `PreFrameUpdate`, `PostFrameUpdate`, `SelfInducedSpeeds`,
  `ItemUpdate`, `FrameStart`, `FrameBookend`, `Metadata`, `PlayerMeta`, `PlayersMeta`,
  `Names`, `Character`, `Characters`, plus all enums and `FirstFrame` /
  `FirstPlayableFrame` constants).
- Public Melee enums under `github.com/ethangamma24/slippi-go/pkg/slippi/melee`
  (`ExternalCharacterID`, `InternalCharacterID`, `Stage`, `EnabledItem`, `Item`).
- `(Stage|ExternalCharacterID|InternalCharacterID|Item).DisplayName() string`
  helpers returning slippi-js-parity display strings.
- Public stats output types under `github.com/ethangamma24/slippi-go/pkg/slippi/stats`
  (`Stats`, `Overall`, `ActionCounts`, `AttackCount`, `GrabCount`, `ThrowCount`,
  `GroundTechCount`, `WallTechCount`, `EdgeCancelCount`, `LCancelCount`, `Stock`,
  `Conversion`, `Combo`, `MoveLanded`, `Ratio`, `InputCounts`) plus the public
  `Compute(types.Game) Stats` entry point.
- `pkg/slippi.NewGameFromBytes(name string, data []byte) *Game`.
- `pkg/slippi.NewGameFromReader(name string, r io.Reader) (*Game, error)`.
- `pkg/slippi.ParseMetaFromBytes(name string, data []byte) (types.Metadata, error)`.
- Typed accessors on `*Game`: `Parsed`, `SettingsTyped`, `MetadataTyped`,
  `StatsTyped`, `FramesTyped`, `GameEndTyped`, `LatestFrameTyped`. These return
  the public typed values directly without the existing JSON-`any` round-trip.

### Changed

- **Field rename:** `pkg/slippi/types.GameEnd.LRASInitiatior` →
  `pkg/slippi/types.GameEnd.LRASInitiator` (typo fix). This also changes the JSON
  payload key emitted by the existing `(any, error)` getters from
  `"LRASInitiatior"` to `"LRASInitiator"`. Consumers of the JSON facade should
  update accordingly.
- The typed model and Melee enums moved out of `internal/`. The packages at
  `internal/goslippi/slippi/...` and `internal/goslippi/slippi/melee/...` are
  gone; their contents now live at `pkg/slippi/types` and `pkg/slippi/melee`.
  External code already could not import the `internal/` paths, so this is not
  a breaking change for external consumers.
- The `internal/stats` package's exported types now live at `pkg/slippi/stats`;
  `Compute` is re-exported there. Behavior unchanged.

### Unchanged (compatibility)

- `pkg/slippi.NewGame(fixturePath string)` and all existing `(any, error)`
  accessors (`GetSettings`, `GetMetadata`, `GetStats`, `GetFrames`,
  `GetLatestFrame`, `GetGameEnd`, `Summary`) keep their signatures and behavior,
  modulo the `LRASInitiatior → LRASInitiator` JSON key rename above.
- `TestStatsParityAgainstReferenceJS` and the performance gate continue to pass.
```

> **Date stamp:** the agent should fill in the actual ISO date (`YYYY-MM-DD`) of
> the merge for `v0.1.0`, not `2026-MM-DD`.

### 3.2 README updates

Edit `README.md` at the repo root. Replace the existing "Current structure" and
"Module" sections so that they reflect the post-promotion layout. Add a "Public
API" section. Suggested replacement (merge with whatever else lives in README):

````markdown
# slippi-go

Native Go library for parsing and analyzing Slippi (Super Smash Bros. Melee)
replay files.

## Module

- Module path: `github.com/ethangamma24/slippi-go`
- Public packages:
  - `github.com/ethangamma24/slippi-go/pkg/slippi` — facade and entry points
  - `github.com/ethangamma24/slippi-go/pkg/slippi/types` — typed Slippi data model
  - `github.com/ethangamma24/slippi-go/pkg/slippi/melee` — Melee enums and display names
  - `github.com/ethangamma24/slippi-go/pkg/slippi/stats` — slippi-js-parity stats output

## Quick start

```go
import "github.com/ethangamma24/slippi-go/pkg/slippi"

// From a file
game := slippi.NewGame("path/to/replay.slp")
parsed, err := game.Parsed(ctx)            // typed
stats, err := game.StatsTyped(ctx)         // typed slippi-js-parity stats
settings, err := game.GetSettings(ctx)     // back-compat untyped facade

// From bytes (e.g. HTTP upload)
g := slippi.NewGameFromBytes("upload.slp", data)
metadata, err := g.MetadataTyped(ctx)
```

## Layout

| Path | Purpose |
|---|---|
| `pkg/slippi/` | Facade (`Game`), entry points, typed accessors |
| `pkg/slippi/types/` | Typed Slippi data model |
| `pkg/slippi/melee/` | Melee enums + `DisplayName()` helpers |
| `pkg/slippi/stats/` | slippi-js-parity stats output + `Compute` |
| `internal/goslippi/` | UBJSON parsing orchestrator (not public) |
| `internal/stats/` | Stats computation implementation (not public) |
| `internal/realtime/` | Frame selection helpers (not public) |
| `internal/io/` | Replay file utilities (not public) |
| `docs/parity_contract.md` | Comparator rules and parity expectations |

## Testing

```sh
go test ./...
go test ./pkg/slippi -run TestPerformanceGate -count=1
```
````

### 3.3 Tag and publish

After both the CHANGELOG and README updates merge to the default branch:

```sh
# 1. Verify clean state and a green test run on the tip of main
git checkout main && git pull && go test ./...

# 2. Tag and push
git tag v0.1.0
git push origin v0.1.0

# 3. Trigger / verify Go module proxy resolution
GOPROXY=https://proxy.golang.org go list -m github.com/ethangamma24/slippi-go@v0.1.0
```

The `go list -m …@v0.1.0` invocation should print the module path and version
without error. If it fails with "unknown version", wait a minute and retry —
the proxy is eventually consistent.

### 3.4 Notify the consumer

Drop a one-line note in this directory's `README.md` (`docs/specs/slippi-go/README.md`)
indicating that `v0.1.0` is live. The downstream backend migration spec
(`backend/docs/specs/slippi-go-parser-migration-2026-04-25.md`) gates on this.

---

## Acceptance criteria

- [ ] `CHANGELOG.md` exists at repo root, follows Keep-a-Changelog 1.1.0, and contains a `v0.1.0` entry covering the items in §3.1.
- [ ] `README.md` reflects the post-promotion package layout and shows both the typed and back-compat usage.
- [ ] `git tag --list 'v0.*'` includes `v0.1.0` on `origin`.
- [ ] `GOPROXY=https://proxy.golang.org go list -m github.com/ethangamma24/slippi-go@v0.1.0` resolves successfully.
- [ ] The `dontgethit` repo's slippi-go specs README is updated with the release date.
- [ ] No code changes in this PR — diff is limited to `CHANGELOG.md`, `README.md`, and (optionally) `docs/`.

## Out of scope

- Promoting to `v1.0.0`. The maintainer wants the option to break the API again before locking in 1.x.
- Publishing to anywhere outside the Go module proxy (no GitHub Releases entry, no announcement post). Optional follow-up.
- Backporting tags to historical commits.

## References

- Keep-a-Changelog: https://keepachangelog.com/en/1.1.0/
- Semantic versioning: https://semver.org/
- Go modules tagging: https://go.dev/ref/mod#vcs-version
- Go module proxy: https://proxy.golang.org
