# Spec: conversion-move-filter

Scope: feature

# Conversion Move Filter

## Goal

Allow consumers to exclude conversions that start with specific move IDs from stats output. Common use case: filtering out Falco's neutral B (MoveID 17) pokes that don't represent meaningful combo openings.

## Design Decisions

### Filter scope

- **Opening move only.** A conversion is excluded if and only if its first move (`Moves[0].MoveID`) matches an excluded ID. Conversions that contain an excluded move after the opener are unaffected.

### Output behavior

- **Remove from output.** Filtered conversions are not included in the returned `Stats.Conversions` slice. No boolean flag or separate list is emitted.

### Configuration

- **Configurable move ID list.** A new `ComputeOptions` struct in the `slippi` package exposes an `ExcludedMoves []uint8` field. The `uint8` type matches the raw `LastHittingAttackID` from replay data. Internally it is cast to `int` for comparison with `MoveID`.
- **Opt-in only.** An empty `ExcludedMoves` slice produces identical output to the current `Compute()` / `StatsTyped()` path. No default exclusions.

### Filter stage

- **Two-stage pipeline.** All conversions are built via the existing frame iteration loop. Filtering is applied once, after `populateOpeningTypes` and before `generateOverall` and ratio calculation. Downstream code sees only the filtered list and requires no changes.

### API surface

- **New method, no breaking changes.** `Game.StatsTypedWithOptions(ctx, ComputeOptions)` is added alongside the existing `StatsTyped()`. Internally it calls a new `stats.ComputeWithOptions(game, opts)` function.
- **`ComputeOptions` lives in `slippi` package.** Consumers import only the public API package. The internal stats function accepts it as a parameter.

## Key Moves

| MoveID | Attack | Notes |
|--------|--------|-------|
| 17 | Neutral B | Falco/Fox blaster, Pikachu thunder jolt, etc. — common filter target |

## Edge Cases

- A conversion with no moves (`len(Moves) == 0`) can never be filtered — there is no opening move to match against.
- If a game contains only excluded-move conversions, the filtered list is empty and all ratios are zero.
- If all conversions for a player are filtered, `ConversionCount` and all derived ratios reflect zero openings.