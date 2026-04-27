# Spec 02: Typed accessors and bytes/reader constructors (`slippi-go`)

**Repository:** [`github.com/ethangamma24/slippi-go`](https://github.com/ethangamma24/slippi-go)
**Scope:** Add public, type-safe entry points to `slippi-go` so external consumers can:
1. Hand a `[]byte` or `io.Reader` directly to the parser instead of always going through a file path.
2. Get **typed** Slippi data (`types.Game`, `stats.Stats`, etc.) out of the parser without the existing JSON-`any` round-trip.

---

## Implementation Status

- **Status:** 📝 Not Started
- **Completed:** —
- **Implemented by:** —
- **Notes:** —

---

## Spec Dependencies

**Depends on:** [`01-public-api-promotion-2026-04-25.md`](./01-public-api-promotion-2026-04-25.md) (the typed packages must already exist at `pkg/slippi/types`, `pkg/slippi/melee`, `pkg/slippi/stats`).
**Blocks:** [`03-release-v0.1.0-2026-04-25.md`](./03-release-v0.1.0-2026-04-25.md), and the consuming backend migration spec.

---

## Context

After spec 01, the typed model is importable from external modules but `pkg/slippi.Game` still:

- Only accepts a file path (`NewGame(fixturePath string)`), so consumers that hold replay bytes in memory (HTTP uploads, S3 streams, in-memory test fixtures) must touch disk to use the library.
- Returns every accessor's result as `(any, error)` after `json.Marshal → json.Unmarshal`, which (a) discards static types — including the typed `melee.Stage`, `melee.InternalCharacterID`, `time.Duration`-style values — and (b) costs ~2× the marshaling work per call.

This spec adds the missing constructors and typed accessors. It does **not** change parsing or stats math, and it does **not** remove or break the existing facade methods.

---

## Goals

1. Provide `NewGameFromBytes(name string, data []byte) *Game` and `NewGameFromReader(name string, r io.Reader) (*Game, error)`.
2. Provide typed accessors on `*Game` returning the public types from spec 01:
   - `Parsed(ctx context.Context) (types.Game, error)`
   - `SettingsTyped(ctx context.Context) (types.GameStart, error)`
   - `MetadataTyped(ctx context.Context) (types.Metadata, error)`
   - `StatsTyped(ctx context.Context) (stats.Stats, error)`
   - `FramesTyped(ctx context.Context) (map[int]types.Frame, error)`
   - `GameEndTyped(ctx context.Context) (types.GameEnd, error)`
   - `LatestFrameTyped(ctx context.Context) (types.Frame, bool, error)`
3. Add a metadata-only fast path: `ParseMetaFromBytes(name string, data []byte) (types.Metadata, error)` (mirrors the existing internal `goslippi.ParseMeta(filePath)`).
4. Keep all existing `(any, error)` accessors and `NewGame(fixturePath string)` working unchanged.
5. `TestStatsParityAgainstReferenceJS` and `TestPerformanceGate` continue to pass.

## Non-goals

- No deprecation of the `(any, error)` accessors. They remain as-is.
- No changes to `Compute` or any event handlers.
- No new third-party dependencies.

---

## Required changes

### 2.1 Refactor `goslippi.ParseGame` to share a bytes core

`internal/goslippi/parse.go::ParseGame(filePath string)` currently does
`os.ReadFile` then `ubjson.Unmarshal`. Extract the bytes path into a private
helper and have both the file-path and bytes entry points share it.

After the refactor `internal/goslippi/parse.go` exposes:

```go
// ParseGame reads the .slp file at filePath and returns the decoded game.
func ParseGame(filePath string) (types.Game, error) { ... }

// ParseGameFromBytes parses an in-memory .slp byte slice and returns the decoded game.
// `name` is used in error wrapping so callers can attach a logical filename.
func ParseGameFromBytes(name string, data []byte) (types.Game, error) { ... }

// ParseGameFromReader reads all bytes from r and parses them. ReadAll is bounded by
// the standard io.ReadAll behavior; callers wanting a hard size limit should pre-buffer.
func ParseGameFromReader(name string, r io.Reader) (types.Game, error) { ... }
```

Apply the same treatment to `internal/goslippi/parse_meta_only.go`:

```go
func ParseMeta(filePath string) (types.Metadata, error)                    { ... }
func ParseMetaFromBytes(name string, data []byte) (types.Metadata, error)   { ... }
```

These are still under `internal/`; the public entry point is via `pkg/slippi.Game`
(below).

### 2.2 Extend `pkg/slippi.Game`

Update `pkg/slippi/game.go`. The struct gains a `data []byte` field and an internal
"how do I get my bytes" tag so `ensureParsed` can pick the right path on first use.

```go
type Game struct {
    name        string        // human-readable name for error wrapping
    fixturePath string        // optional; empty when constructed from bytes/reader
    data        []byte        // optional; populated by NewGameFromBytes / NewGameFromReader

    once     sync.Once
    parsed   types.Game
    parseErr error
}
```

Add the new constructors:

```go
// NewGame keeps its existing signature for backward compatibility.
func NewGame(fixturePath string) *Game {
    return &Game{name: filepath.Base(fixturePath), fixturePath: fixturePath}
}

// NewGameFromBytes returns a Game backed by an in-memory .slp byte slice.
// `name` is included in error messages (e.g. "parse game replay-123.slp: …").
func NewGameFromBytes(name string, data []byte) *Game {
    return &Game{name: name, data: data}
}

// NewGameFromReader fully reads r into memory and returns a Game backed by those bytes.
// Returns an error if the read fails. Pass a size-limited reader (e.g. io.LimitReader)
// if you need an upper bound.
func NewGameFromReader(name string, r io.Reader) (*Game, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, fmt.Errorf("read replay %s: %w", name, err)
    }
    return NewGameFromBytes(name, data), nil
}
```

Update `ensureParsed` to handle both cases:

```go
func (g *Game) ensureParsed(ctx context.Context) error {
    if err := ctx.Err(); err != nil {
        return err
    }
    g.once.Do(func() {
        switch {
        case len(g.data) > 0:
            g.parsed, g.parseErr = goslippi.ParseGameFromBytes(g.name, g.data)
        case g.fixturePath != "":
            absPath, err := resolveFixturePath(g.fixturePath)
            if err != nil {
                g.parseErr = err
                return
            }
            g.parsed, g.parseErr = goslippi.ParseGame(absPath)
        default:
            g.parseErr = fmt.Errorf("slippi: Game has no source (use NewGame, NewGameFromBytes, or NewGameFromReader)")
        }
    })
    if g.parseErr != nil {
        return fmt.Errorf("parse game %s: %w", g.name, g.parseErr)
    }
    return nil
}
```

> Once we accept bytes, holding a reference to `g.data` keeps the entire input
> alive for the lifetime of `*Game`. That is intentional: callers using
> `NewGameFromBytes` are signaling they want re-parse to be possible from the
> same buffer if `parseErr` was transient. If memory pressure becomes a concern,
> we can clear `g.data = nil` at the end of `ensureParsed` after a successful
> parse — make this an internal optimization, *not* part of the public contract.

### 2.3 Add typed accessors

Add to `pkg/slippi/game.go`:

```go
// Parsed returns the entire decoded game in typed form. The returned value
// shares no mutable state with the Game; callers may freely retain it.
func (g *Game) Parsed(ctx context.Context) (types.Game, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return types.Game{}, err
    }
    return g.parsed, nil
}

// SettingsTyped returns the GameStart event payload in typed form.
func (g *Game) SettingsTyped(ctx context.Context) (types.GameStart, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return types.GameStart{}, err
    }
    return g.parsed.Data.GameStart, nil
}

// MetadataTyped returns the parsed UBJSON metadata block in typed form.
func (g *Game) MetadataTyped(ctx context.Context) (types.Metadata, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return types.Metadata{}, err
    }
    return g.parsed.Meta, nil
}

// StatsTyped returns slippi-js-parity stats in typed form.
func (g *Game) StatsTyped(ctx context.Context) (stats.Stats, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return stats.Stats{}, err
    }
    return stats.Compute(g.parsed), nil
}

// FramesTyped returns the per-frame data in typed form. The map is the same one
// held by the Game and SHOULD NOT be mutated by the caller.
func (g *Game) FramesTyped(ctx context.Context) (map[int]types.Frame, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return nil, err
    }
    return g.parsed.Data.Frames, nil
}

// GameEndTyped returns the GameEnd event payload in typed form.
func (g *Game) GameEndTyped(ctx context.Context) (types.GameEnd, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return types.GameEnd{}, err
    }
    return g.parsed.Data.GameEnd, nil
}

// LatestFrameTyped returns the highest-numbered finalized frame, plus a bool
// indicating whether such a frame exists.
func (g *Game) LatestFrameTyped(ctx context.Context) (types.Frame, bool, error) {
    if err := g.ensureParsed(ctx); err != nil {
        return types.Frame{}, false, err
    }
    f, ok := realtime.LatestFrame(g.parsed.Data.Frames)
    return f, ok, nil
}
```

### 2.4 Public `ParseMetaFromBytes`

Add a top-level helper to `pkg/slippi/parse.go` (new file):

```go
package slippi

import (
    "github.com/ethangamma24/slippi-go/internal/goslippi"
    types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

// ParseMetaFromBytes parses just the UBJSON metadata block from a .slp byte slice.
// This is significantly cheaper than full parse + Game.MetadataTyped when only
// metadata is needed (e.g. dashboard listings).
func ParseMetaFromBytes(name string, data []byte) (types.Metadata, error) {
    return goslippi.ParseMetaFromBytes(name, data)
}
```

### 2.5 Tests

Add new test file `pkg/slippi/typed_test.go` covering:

| Test | Asserts |
|---|---|
| `TestNewGameFromBytes_Parsed_MatchesAnyFacade` | Picks one fixture, calls both `NewGame(path).Summary()` (`any`) and `NewGameFromBytes(name, data).Parsed()` (typed); marshals each to JSON; asserts byte-for-byte equality of the JSON projections of `Settings`, `Metadata`, `GameEnd`. |
| `TestStatsTyped_MatchesAnyFacade` | Same fixture; asserts `json.Marshal(StatsTyped()) == GetStats() then json.Marshal(any)`. |
| `TestNewGameFromReader_HandlesIOErrors` | Pass a reader that returns `io.ErrUnexpectedEOF` mid-stream; assert constructor returns wrapped error and parsed Game can be GC'd without panic. |
| `TestGame_NoSource_ErrorsCleanly` | `&Game{}` (no constructor) → `Parsed(ctx)` returns the explicit "no source" error. |
| `TestParseMetaFromBytes_OnlyTouchesMetadata` | Ensure `ParseMetaFromBytes` decodes `Metadata.StartAt` and `LastFrame` without consuming the `raw` array (compare against `MetadataTyped` from `NewGame(path)` for the same fixture). |
| `TestNewGameFromBytes_ContextCancellation` | Pre-cancel context, call `Parsed`, expect `context.Canceled`. |

### 2.6 Backward compatibility

The existing methods `GetSettings/GetMetadata/GetStats/GetFrames/GetLatestFrame/GetGameEnd/Summary` (`(any, error)`) must continue to work identically. Internally they should be rewritten as thin shims over the typed accessors:

```go
func (g *Game) GetSettings(ctx context.Context) (any, error) {
    s, err := g.SettingsTyped(ctx)
    if err != nil {
        return nil, err
    }
    return toAny(s)
}
```

Confirm via test that `Summary(ctx)` returns the same shape after this refactor.

---

## Acceptance criteria

- [ ] `go test ./...` passes from a clean checkout.
- [ ] `go test ./pkg/slippi -run TestPerformanceGate -count=1` still passes; the typed-accessor path is allowed to be **faster** but must not regress.
- [ ] `TestStatsParityAgainstReferenceJS` still passes.
- [ ] All tests in §2.5 are present and pass.
- [ ] Public API additions are exactly: `NewGameFromBytes`, `NewGameFromReader`, `ParseMetaFromBytes`, plus `Parsed`, `SettingsTyped`, `MetadataTyped`, `StatsTyped`, `FramesTyped`, `GameEndTyped`, `LatestFrameTyped` methods. No other public symbols change.
- [ ] `git diff --stat origin/main` shows no edits to `internal/goslippi/event/handler/handlers/*` (the parser code itself is untouched).
- [ ] No new third-party dependencies in `go.mod`.

## Out of scope

- Streaming/realtime mode (consuming events as they arrive over a socket). The current architecture decodes the full UBJSON `raw` array up front; lifting that is a future spec.
- A `Close()` / explicit `Free()` method on `Game`. Standard Go GC handles the byte buffer.
- Generics / type-parametric API. Keep the surface concrete and easy to consume.

## References

- Existing facade: [`pkg/slippi/game.go`](https://github.com/ethangamma24/slippi-go/blob/main/pkg/slippi/game.go)
- Existing parse path: [`internal/goslippi/parse.go`](https://github.com/ethangamma24/slippi-go/blob/main/internal/goslippi/parse.go)
- Metadata-only path: [`internal/goslippi/parse_meta_only.go`](https://github.com/ethangamma24/slippi-go/blob/main/internal/goslippi/parse_meta_only.go)
- Realtime helper: [`internal/realtime/realtime.go`](https://github.com/ethangamma24/slippi-go/blob/main/internal/realtime/realtime.go)
