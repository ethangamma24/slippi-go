# Spec: upstream-shield-state

Scope: feature

# Spec: Upstream Shield-State Constants and Helpers

**Repository:** `github.com/ethangamma24/slippi-go`
**Scope:** Add shield action-state constants and a typed `IsShieldStun()` helper to the public `pkg/slippi/melee/` and `pkg/slippi/types/` packages. Enables downstream consumers (DontGetHit backend) to detect shield hits precisely instead of using coarse heuristics.

---

## Implementation Status

- **Status:** ✅ Implemented
- **Completed:** 2026-04-27
- **Implemented by:** OpenCode
- **Notes:** All constants, helpers, and tests created. `go test ./...` passes; performance gate passes. No new dependencies.

---

## Depends on

- `docs/specs/slippi-go/01-public-api-promotion-2026-04-25.md` (completed) — public `pkg/slippi/melee/` and `pkg/slippi/types/` packages must exist.

---

## Goals

1. Add canonical Melee shield action-state constants to `pkg/slippi/melee/`:
   - `StateGuardOn` = 178
   - `StateGuard` = 179
   - `StateGuardDamage` = 180
   - `StateGuardOff` = 181
   - `StateGuardSetOff` = 182
2. Add `IsShieldState(state uint16) bool` helper that returns true for states 178–182.
3. Add `IsShieldStun(state uint16) bool` helper that returns true **only** for `StateGuardDamage` (180).
4. Add `IsShieldStun() bool` method on `types.PostFrameUpdate` for ergonomic frame-level checks.
5. Keep all existing tests passing; add new tests for the helpers.

---

## Required Changes

### 1. New file: `pkg/slippi/melee/shield_states.go`

```go
package melee

// Shield action state IDs from the standard Melee action state table.
const (
	StateGuardOn     uint16 = 178
	StateGuard       uint16 = 179
	StateGuardDamage uint16 = 180
	StateGuardOff    uint16 = 181
	StateGuardSetOff uint16 = 182
)

// IsShieldState returns true if the action state is any shield-related state
// (GuardOn, Guard, GuardDamage, GuardOff, GuardSetOff).
func IsShieldState(state uint16) bool {
	return state >= StateGuardOn && state <= StateGuardSetOff
}

// IsShieldStun returns true only for GuardDamage (180), the direct signal
// that a character's shield was hit and is currently in shield stun.
func IsShieldStun(state uint16) bool {
	return state == StateGuardDamage
}
```

### 2. New method on `types.PostFrameUpdate`

**File:** `pkg/slippi/types/frame.go`

Add to the `PostFrameUpdate` struct (or as a method on the type):

```go
// IsShieldStun returns true if the player is in GuardDamage (shield stun),
// indicating their shield was hit on this frame.
func (p *PostFrameUpdate) IsShieldStun() bool {
	return melee.IsShieldStun(p.ActionStateID)
}
```

### 3. Tests

**File:** `pkg/slippi/melee/shield_states_test.go`

Cover:
- `IsShieldState` returns true for 178–182, false for 177 and 183.
- `IsShieldStun` returns true only for 180, false for all others.

**File:** `pkg/slippi/types/frame_test.go` (new or append to existing)

Cover:
- `PostFrameUpdate.IsShieldStun()` returns true when `ActionStateID == 180`.
- Returns false when `ActionStateID` is 179, 181, etc.

---

## Acceptance Criteria

- [ ] `pkg/slippi/melee/shield_states.go` exists with constants and helpers.
- [ ] `pkg/slippi/melee/shield_states_test.go` passes.
- [ ] `types.PostFrameUpdate.IsShieldStun()` exists and is tested.
- [ ] `go test ./...` passes from a clean checkout.
- [ ] No new third-party dependencies.
- [ ] Existing public API surface remains backward-compatible.

---

## Out of Scope

- Changes to stats computation (`pkg/slippi/stats/compute.go`).
- Changes to the `(any, error)` getters on `pkg/slippi.Game`.
- DontGetHit backend consumption (covered by `backend-shield-hit` repo spec).
- Hitbox-level shield analysis using `InstanceHitBy`.

---

## References

- Melee action state table: GuardOn (178), Guard (179), GuardDamage (180), GuardOff (181), GuardSetOff (182).
- Current downstream heuristic: `actionState >= 178 && actionState <= 182` in `dontgethit/backend/internal/domain/services/replay_report_service.go`.
