# Spec: upstream-parity-bugs

Scope: feature

# Spec: Upstream Parity Bug Fixes

**Repository:** `github.com/ethangamma24/slippi-go`
**Scope:** Fix two upstream bugs in slippi-go's stats engine that cause corpus parity mismatches vs slippi-js. These fixes are external to the dontgethit backend; the backend will consume them via a version bump after release.

---

## Implementation Status

- **Status:** ✅ Done
- **Completed:** 2026-04-27
- **Implemented by:** OpenCode
- **Notes:** Both bugs fixed in `pkg/slippi/stats/compute.go` and released as v0.1.1.

---

## Depends on

- `backend/docs/specs/adapter-parity-fixes-2026-04-27.md` (or equivalent) — adapter-side fixes must be complete so remaining mismatches are purely upstream.
- `docs/specs/slippi-go-parser-migration-2026-04-25.md` — establishes baseline and mismatch categories.

---

## Bugs

### Bug 1 — handleWavedash heuristic divergence

**Location:** `pkg/slippi/stats/compute.go` (lines ~916–994)
**Symptom:** Cross-displacements between `wavedash`, `waveland`, and `airDodge` counts vs slippi-js. Example: Go waveland=5 / wavedash=9 vs JS waveland=6 / wavedash=8.
**Root cause:** When slippi-go detects a special-landing after an air-dodge, it decrements `AirDodgeCount` and increments either `WavedashCount` or `WavelandCount` based on a 15-frame lookback + Y-displacement heuristic. The thresholds and frame-counter logic differ from slippi-js.
**Fix:** Align the heuristic with slippi-js — compare the exact frame window, Y-velocity thresholds, and landing detection logic.
**Files to produce:**
- Minimal reproduction `.slp` replay that triggers the divergence
- Expected vs actual counts table
- Proposed code diff for `compute.go`

### Bug 2 — populateOpeningTypes uses last-move playerIndex

**Location:** `pkg/slippi/stats/compute.go` (lines ~567–614)
**Symptom:** Occasional 1-count flips between `neutral-win` and `counter-attack` opening types. Example: Go classifies 3 conversions as neutral-win; JS classifies 2 as neutral-win, 1 as counter-attack.
**Root cause:** `populateOpeningTypes` uses `c.Moves[len(c.Moves)-1].PlayerIndex` (last move) to determine the attacker, but slippi-js groups conversions by `c.Moves[0].PlayerIndex` (first move). Conversions where the victim hits back at the end are misclassified.
**Fix:** Change `lastMovePlayer` assignment to use `c.Moves[0].PlayerIndex` when available.
**Files to produce:**
- Minimal reproduction `.slp` replay with a late hit from the victim
- Expected vs actual openingType counts
- Proposed code diff for `compute.go`

---

## Acceptance Criteria

- [ ] Both bugs reproduced with minimal `.slp` fixtures.
- [ ] Issues filed on `github.com/ethangamma24/slippi-go` with clear reproduction steps.
- [ ] PRs opened with proposed fixes and tests.
- [ ] PRs merged and released as a new version tag.
- [ ] Backend `go.mod` bumped to the fixed version.
- [ ] Corpus parity re-run confirms zero exact mismatches.
- [ ] `backend/docs/specs/slippi-go-parser-migration-2026-04-25.md` updated with upstream fix details.

---

## Out of Scope

- Adapter-side fixes (handled separately).
- Changes to dontgethit domain models or API contracts.
- Any feature enrichment beyond parity alignment.