# slippi-go

Native Go library for parsing and analyzing Slippi replay data.

## Module

- Module path: `github.com/ethangamma24/slippi-go`
- Public package: `github.com/ethangamma24/slippi-go/pkg/slippi`

Example:

```go
import "github.com/ethangamma24/slippi-go/pkg/slippi"
```

## Current structure

- `pkg/slippi`: public Go API facade (`Game`)
- `internal/stats`: native stats aggregation primitives
- `internal/realtime`: native frame selection/realtime helpers
- `internal/io`: native replay file writing utilities
- `docs/parity_contract.md`: comparator rules and parity expectations

## Run

Run tests:

```bash
go test ./...
```

Run benchmark gate directly:

```bash
go test ./pkg/slippi -run TestPerformanceGate -count=1
```

## Publish as a Go module

1. Ensure tests pass locally:

```bash
go test ./...
```

2. Commit and push your default branch to GitHub.

3. Create a semver tag and push it (first stable release example):

```bash
git tag v1.0.0
git push origin v1.0.0
```

4. Verify the module resolves from the Go proxy:

```bash
go list -m github.com/ethangamma24/slippi-go@v1.0.0
```

5. (Optional) Trigger faster indexing if the version is not visible yet:

```bash
GOPROXY=https://proxy.golang.org go list -m github.com/ethangamma24/slippi-go@v1.0.0
```