# Agent Instructions

## Project

Native Go library for parsing and analyzing Slippi replay data.

- **Module:** `github.com/ethangamma24/slippi-go`
- **Go version:** 1.24.2
- **Public API:** `pkg/slippi` (entrypoint is the `Game` facade)

## Structure

- `pkg/slippi` — public API (`Game` struct wrapping parse + stats + realtime)
- `internal/goslippi` — low-level replay parser
- `internal/stats` — stats aggregation primitives; embeds `framedata.json` via `//go:embed`
- `internal/realtime` — frame-selection helpers
- `internal/io` — replay file writer
- `testdata/slp/` — `.slp` fixture files used by tests
- `docs/parity_contract.md` — comparator rules for parity tests

## Commands

```bash
# Run all tests
go test ./...

# Run only the performance gate
go test ./pkg/slippi -run TestPerformanceGate -count=1

# Benchmark
go test ./pkg/slippi -bench BenchmarkGoSummary -benchtime 5s
```

## CI

`.github/workflows/ci.yml` runs `go test ./...` with:

- `TZ=UTC`
- `LC_ALL=C`

## Testing Notes

- Fixture paths in tests are relative to the repo root (resolved by walking up to `go.mod`).
- `testdata/slp/incomplete.slp` is explicitly skipped when enumerating fixtures.
- The performance gate runs 5 measured iterations plus a warmup and compares the median against a static limit (currently 120s for the full fixture set).

## Known State

- `pkg/slippi/parity_test.go` is currently an empty file, which causes `go test ./...` to fail with a parse error (`expected 'package', found 'EOF'`).
- Most `pkg/slippi/*_test.go` files reference `assertOperationParity`, which is not currently defined in the repo, so parity tests will not compile even after fixing the empty file.
- The only self-contained tests in `pkg/slippi` right now are `TestPerformanceGate` and `BenchmarkGoSummary` (in `performance_test.go`).

## Parity Contract

When modifying stats or parser behavior, consult `docs/parity_contract.md`:

- Float tolerance is absolute-only: `abs <= 1e-9`
- Frame iteration starts at `Frames.FIRST` (`-123`) and stops at the first missing/incomplete frame
- Opening-type ratios group by `conversion.moves[0].playerIndex`; conversions with no moves are excluded
- Slice-index tracking is used for in-progress state updates to survive append reallocations
