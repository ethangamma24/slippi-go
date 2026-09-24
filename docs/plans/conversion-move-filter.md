---
plan name: conversion-move-filter
plan description: Add optional filtering of conversions by opening move ID
plan status: active
---

## Idea
Allow consumers to exclude conversions that start with specific move IDs (e.g., Falco's neutral B, MoveID 17) from stats computation. The filter is applied after all conversions are built but before downstream ratios/overall stats are calculated. Conversions whose `Moves[0].MoveID` matches any excluded ID are removed entirely from output. The feature is fully opt-in and backwards-compatible.

## Implementation
- Add ComputeOptions struct with ExcludedMoves []uint8 to pkg/slippi
- Add stats.ComputeWithOptions(game, opts) that builds conversions then filters before downstream processing
- Add Game.StatsTypedWithOptions(ctx, opts) forwarding to stats.ComputeWithOptions
- Add helper constant or docs identifying neutral B = MoveID 17
- Add tests: no opts returns all conversions, excluding 17 removes neutral-B-opening conversions
- Run go test ./... to verify backwards compatibility

## Required Specs
<!-- SPECS_START -->
- conversion-move-filter
<!-- SPECS_END -->