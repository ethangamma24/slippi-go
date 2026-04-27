# Spec: parity-gap-v011

Scope: repo

# Spec: v0.1.1 Parity Gap Analysis

**Repository:** backend
**Scope:** Document and track the remaining corpus parity mismatches after bumping `slippi-go` from `v0.1.0` to `v0.1.1`. Provide actionable findings for the upstream agent.

---

## Implementation Status

- **Status:** 📝 Not Started
- **Completed:** —
- **Implemented by:** —
- **Notes:** —

---

## Depends on

- `backend/docs/specs/slippi-go-parser-migration-2026-04-25.md` (completed) — establishes migration baseline.
- `backend/docs/specs/adapter-parity-fixes-*.md` (completed) — adapter-side fixes reduced mismatches from 102 to 61.
- `docs/specs/slippi-go/upstream-parity-bugs-*.md` (partially completed) — upstream v0.1.1 released but did not fix all issues.

---

## v0.1.1 — What Was Fixed

Verified diff of `pkg/slippi/stats/compute.go` between v0.1.0 and v0.1.1:

1. **`populateOpeningTypes` attacker resolution** (line 602):
   - v0.1.0: `lastMovePlayer = c.Moves[len(c.Moves)-1].PlayerIndex` (last move)
   - v0.1.1: `lastMovePlayer = c.Moves[0].PlayerIndex` (first move)
   - **Status:** ✅ Fixed. This was one of the two identified upstream bugs.

2. **`handleWavedash` refactor** (lines 926–1000):
   - Extracted `isWavedashInitiationAnimation(animation uint16) bool` helper.
   - Renamed `lookback` → `lookbackFrames`.
   - **Status:** ❌ Not fixed. The heuristic logic is byte-for-byte identical to v0.1.0.

---

## v0.1.1 — What Was NOT Fixed

### Bug A: `handleWavedash` heuristic divergence (largest remaining issue)

**Impact:** 50 files (waveland), 45 files (airDodge), 29 files (wavedash)

**Root cause:** When slippi-go detects a special-landing after an air-dodge, it decrements `AirDodgeCount` and increments either `WavedashCount` or `WavelandCount` based on a 15-frame lookback + Y-displacement heuristic. The thresholds and frame-counter logic differ from slippi-js, causing misclassification at the margin.

**Concrete examples:**
- `Game_20260416T111441.slp`: Go wavedash=9/waveland=5 vs JS wavedash=8/waveland=6 (1 landing misclassified)
- `Game_20260416T113816.slp`: Go wavedash=23/waveland=10 vs JS wavedash=25/waveland=8 (2 landings misclassified)
- `Game_20260420T224550.slp`: Go wavedash=39/waveland=6 vs JS wavedash=30/waveland=15 (massive divergence)

**Fix needed:** Align `handleWavedash` with slippi-js — compare exact frame window, Y-velocity thresholds, and landing detection logic. The v0.1.1 refactor did not change behavior; the agent must actually modify the heuristic.

---

### Bug B: Open-conversion handling in `populateOpeningTypes`

**Impact:** 5 files (neutral-win), 1 file (trade), 4 files (counter-attack)

**Root cause:** slippi-go's `populateOpeningTypes` skips open conversions entirely:
```go
if c.EndFrame == nil {
    continue
}
```
slippi-js handles open conversions differently — they are classified based on available state, not skipped.

**Concrete example:**
- `Game_20260416T114847.slp`: Go trade=5 vs JS trade=6. The 6th JS trade is conversion #66 (victim=1, attacker=0, start=18019, end=undefined) which slippi-go skips because `EndFrame == nil`.
- `Game_20260416T115359.slp`: Go counter-attack=0/neutral-win=3 vs JS counter-attack=1/neutral-win=2. The JS counter-attack is an open conversion that slippi-go skips.

**Fix needed:** Do not skip open conversions in `populateOpeningTypes`. slippi-js classifies them using the same logic as closed conversions (checking `oppEndFrame` against `startFrame`). An open conversion with `oppEndFrame > startFrame` should still be classified as `counter-attack`.

---

### Bug C: `edgeCancelSlow` enrichment (acceptable deviation)

**Impact:** 8 files

**Root cause:** slippi-go exposes `EdgeCancelCount.Slow` which slippi-js does not. Our adapter keeps this enrichment. The parity script tracks it separately.

**Decision:** This is an acceptable deviation — slippi-go provides richer data. No upstream fix needed.

---

## Corpus Parity Snapshot (post-v0.1.1)

```
Total files compared: 102
Exact match count: 41
Partial mismatch count: 61
Parse failures: 0

## Mismatch categories (strict)
counts.actionCounts.waveland: 50
counts.actionCounts.airDodge: 45
counts.actionCounts.wavedash: 29
openingTypes.neutral-win: 5
openingTypes.trade: 1
openingTypes.counter-attack: 4
counts.actionCounts.edgeCancelSlow: 8
players.0.actionCount: 1
players.1.actionCount: 1
```

The `players.*.actionCount` mismatches (2 files) are downstream effects of `airDodge` diffs — when the airDodge count differs by 1, the total actionCount also differs by 1.

---

## Acceptance Criteria

- [ ] Upstream `handleWavedash` heuristic aligned with slippi-js (Bug A).
- [ ] Upstream `populateOpeningTypes` handles open conversions (Bug B).
- [ ] New slippi-go version released with both fixes.
- [ ] Backend `go.mod` bumped to the new version.
- [ ] Corpus parity re-run confirms exact match count rises from 41 toward 102.
- [ ] `counts.actionCounts.wavedash`, `waveland`, `airDodge` no longer appear in mismatch categories.
- [ ] `openingTypes.*` no longer appears in mismatch categories.
- [ ] This spec marked `✅ Implemented` with final exact-match count.

---

## Files to Provide Upstream Agent

1. **Reproduction replays** (minimal `.slp` files that trigger each bug):
   - `Game_20260416T111441.slp` — wavedash/waveland misclassification
   - `Game_20260416T114847.slp` — open-conversion skipped in openingTypes
2. **Expected vs actual tables** for each reproduction replay.
3. **Proposed code diffs** for `pkg/slippi/stats/compute.go`.

---

## Out of Scope

- `edgeCancelSlow` enrichment (Bug C) — acceptable deviation, no fix needed.
- Adapter-side changes — all adapter fixes are complete.
- Changes to dontgethit domain models or API contracts.