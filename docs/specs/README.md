# `slippi-go` upstream specs

This directory holds specs to be executed **outside** the `dontgethit` monorepo, in the
[`ethangamma24/slippi-go`](https://github.com/ethangamma24/slippi-go) repository. They
are written for an agent who will open that repo, implement the changes, and ship a
tagged release.

## Why these live here

`slippi-go` is the canonical Go port of `slippi-js`. The DontGetHit backend wants to
replace its vendored Slippi parser/stats with `slippi-go`, but cannot consume the
library's typed model today because it lives under Go `internal/`. These specs
describe a one-shot batch of upstream changes that promotes the public API surface
the backend needs.

> **Design intent (read this first):** the maintainer plans for this batch to be one
> of the **last** sets of upstream changes to `slippi-go`, because `slippi-js` itself
> updates infrequently. Do not optimize for "small additive PRs over time" — these
> specs intentionally combine several public-API decisions that should be locked in
> together while the API can still break.

## Specs (implement in order)

| # | Spec | Summary |
|---|------|---------|
| 01 | [`01-public-api-promotion-2026-04-25.md`](./01-public-api-promotion-2026-04-25.md) | Promote types, melee enums, and stats output to `pkg/slippi/…`; rename one typo; add display-name helpers. |
| 02 | [`02-typed-accessors-and-bytes-constructor-2026-04-25.md`](./02-typed-accessors-and-bytes-constructor-2026-04-25.md) | Add `NewGameFromBytes` / `NewGameFromReader` and typed `*Typed()` accessors that bypass the JSON-`any` round-trip. |
| 03 | [`03-release-v0.1.0-2026-04-25.md`](./03-release-v0.1.0-2026-04-25.md) | CHANGELOG, README, and the `v0.1.0` tag that the backend will pin. |

## Consuming spec

The DontGetHit backend migration that consumes the result of this work lives at
[`backend/docs/specs/slippi-go-parser-migration-2026-04-25.md`](../../../backend/docs/specs/slippi-go-parser-migration-2026-04-25.md).
That spec assumes 01–03 have shipped and `github.com/ethangamma24/slippi-go@v0.1.0`
is resolvable from the Go module proxy.

## House rules for the upstream agent

- **No new product features.** These specs only re-export, rename, and add
  constructors/accessors; they do not change parsing or stats behavior. The
  `TestStatsParityAgainstReferenceJS` gate must keep passing.
- **Backward compatibility.** Existing `pkg/slippi.Game` getters that return
  `(any, error)` must keep working unchanged. Typed accessors are additive.
- **Single new external dep allowance.** The upstream module already depends on
  `github.com/toitware/ubjson`, `github.com/rs/zerolog`, `golang.org/x/text`, and
  `github.com/stretchr/testify`. Do not add new third-party deps to satisfy these
  specs; everything below is achievable with the current dep set.
- **Internal -> public surface choice.** Where a spec says "promote", the preferred
  technique is **moving** the package out of `internal/` rather than re-exporting via
  type aliases — type aliases would leave the canonical type still under `internal/`
  and consumers cannot use type assertions / package-level constants reliably across
  that boundary. If a move would create a large diff, an alias-only PR is acceptable
  *if* the spec's acceptance criteria still pass.

## Release status

- `v0.1.0` shipped on 2026-04-26.
