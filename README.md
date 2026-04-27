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
| `internal/realtime/` | Frame selection helpers (not public) |
| `internal/io/` | Replay file utilities (not public) |
| `docs/parity_contract.md` | Comparator rules and parity expectations |

## Testing

```sh
go test ./...
go test ./pkg/slippi -run TestPerformanceGate -count=1
```

## Publish as a Go module

1. Ensure tests pass locally:

```bash
go test ./...
```

2. Commit and push your default branch to GitHub.

3. Create a semver tag and push it:

```bash
git tag v0.1.0
git push origin v0.1.0
```

4. Verify the module resolves from the Go proxy:

```bash
GOPROXY=https://proxy.golang.org go list -m github.com/ethangamma24/slippi-go@v0.1.0
```
