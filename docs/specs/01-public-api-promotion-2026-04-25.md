# Spec 01: Public API promotion (`slippi-go`)

**Repository:** [`github.com/ethangamma24/slippi-go`](https://github.com/ethangamma24/slippi-go)
**Scope:** Move the typed Slippi data model, the `melee` enum package, and the stats output types out of `internal/` into the public `pkg/slippi/…` surface. Rename one mis-spelled exported field (`LRASInitiatior` → `LRASInitiator`). Add display-name helpers on stage and character enums.

---

## Implementation Status

- **Status:** 📝 Not Started
- **Completed:** —
- **Implemented by:** —
- **Notes:** —

---

## Spec Dependencies

**Depends on:** None
**Blocks:** [`02-typed-accessors-and-bytes-constructor-2026-04-25.md`](./02-typed-accessors-and-bytes-constructor-2026-04-25.md), [`03-release-v0.1.0-2026-04-25.md`](./03-release-v0.1.0-2026-04-25.md), and the backend migration spec at `dontgethit/backend/docs/specs/slippi-go-parser-migration-2026-04-25.md`.

---

## Context

`slippi-go` exposes a single public package, `pkg/slippi`, with a `Game` facade
whose methods all return `(any, error)` after a `json.Marshal → json.Unmarshal`
round-trip (see `pkg/slippi/game.go::toAny`). The strongly-typed model and all
slippi-js-parity enums live under `internal/goslippi/slippi/...` and
`internal/goslippi/slippi/melee/...`, which Go forbids importing from any other
module.

Downstream consumers (specifically the [DontGetHit backend](https://github.com/dontgethit/dontgethit-go))
need to consume typed values — `slippi.Game`, `slippi.Frame`, `slippi.Player`,
`melee.Stage`, `stats.Stats`, etc. — both for type safety and to skip the
double-marshal cost on every replay.

This spec promotes those types to the public surface. It does **not** change any
parsing or stats behavior; it is purely an API-shape change.

---

## Goals

1. Make these packages importable from external Go modules:
   - `slippi.*` (Data model: `Game`, `Data`, `GameStart`, `Player`, `GameEnd`, `PlayerPlacement`, `Frame`, `PreFrameUpdate`, `PostFrameUpdate`, `SelfInducedSpeeds`, `ItemUpdate`, `FrameStart`, `FrameBookend`, `Metadata`, `PlayerMeta`, `PlayersMeta`, `Names`, `Character`, `Characters`, plus all enums: `TimerType`, `InGameMode`, `ItemSpawnBehaviour`, `PlayerType`, `TeamShade`, `TeamColour`, `GameMode`, `Language`, `GameEndMethod`, `HurtboxCollisionState`, `MissileType`, `TurnipFace`, `GeckoCodeList`, plus the `FirstFrame`/`FirstPlayableFrame` constants.)
   - `melee.*` (`ExternalCharacterID`, `InternalCharacterID`, `Stage`, `EnabledItem`, `Item`, plus all `Ext_*` / `Int_*` / `Stage*` / item constants).
   - The stats output types (`stats.Stats`, `stats.Overall`, `stats.ActionCounts`, `stats.Stock`, `stats.Conversion`, `stats.Combo`, `stats.MoveLanded`, `stats.Ratio`, `stats.InputCounts`).
2. Rename `slippi.GameEnd.LRASInitiatior` → `slippi.GameEnd.LRASInitiator`.
3. Add `DisplayName() string` helpers on `melee.Stage`, `melee.ExternalCharacterID`, `melee.InternalCharacterID`, `melee.Item`. (Move-name lookup is **not** included — the consuming backend will keep its own move-name table; see "Out of scope" below.)
4. Keep the existing `pkg/slippi.Game` `(any, error)` API working — it is additive only, no signatures change.
5. `TestStatsParityAgainstReferenceJS` and the `TestPerformanceGate` continue to pass.

## Non-goals

- No changes to UBJSON parsing, event handlers, or stats math.
- No new constructors yet (`NewGameFromBytes`, `NewGameFromReader` belong to spec 02).
- No new typed accessor methods (`*Typed()` belong to spec 02).
- No CHANGELOG / README / version-tag work (belongs to spec 03).

---

## Required changes

The work is laid out as four sections that can ship as separate PRs.

### 1.1 Promote the slippi data model

Move (`git mv`) the package directory:

```
internal/goslippi/slippi/                  →  pkg/slippi/slippi/
internal/goslippi/slippi/event/            →  internal/goslippi/event/         (stays internal)
internal/goslippi/slippi/event/handler/…   →  internal/goslippi/event/handler/ (stays internal)
internal/goslippi/slippi/melee/            →  pkg/slippi/melee/                (see §1.2)
```

The event decoder and handler packages stay internal because they are
implementation details. Update import paths everywhere — `internal/goslippi/parse.go`,
`internal/goslippi/parse_meta_only.go`, `internal/stats/stats.go`,
`internal/realtime/realtime.go`, `pkg/slippi/game.go`, all event handlers.

The new package import path is:

```go
import "github.com/ethangamma24/slippi-go/pkg/slippi/slippi"
```

That double-`slippi` is awkward. **Preferred final layout:**

- Public root: `pkg/slippi`
  - Existing facade: `pkg/slippi/game.go` (untouched in this spec)
  - Promoted types: keep them in a child sub-package named **`pkg/slippi/types`**.

So the actual moves are:

```
internal/goslippi/slippi/{frame,game,game_start,game_end,gecko_code,meta}.go
        →  pkg/slippi/types/{frame,game,game_start,game_end,gecko_code,meta}.go
internal/goslippi/slippi/melee/*.go
        →  pkg/slippi/melee/*.go
```

The `pkg/slippi.Game` facade in `pkg/slippi/game.go` should switch its
`nativeslippi` import to `slippitypes "github.com/ethangamma24/slippi-go/pkg/slippi/types"`
(or whatever short alias you prefer).

#### 1.1.1 `LRASInitiatior` rename

While the type is moving and the public API is being defined for the first time,
fix the existing typo:

- `pkg/slippi/types/game_end.go`: rename field `LRASInitiatior int8` → `LRASInitiator int8`.
- `internal/goslippi/event/handler/handlers/game_end_handler.go`: update the struct literal.

This is the only allowed *behavior-adjacent* change in this spec — the field's
semantics are identical, only the spelling changes. There are no other consumers
inside the module, and the `(any, error)` JSON facade exposes it via Go's default
field-name → JSON-key mapping, so the JSON payload key changes from
`LRASInitiatior` to `LRASInitiator`. **Document this rename loudly** in the
follow-up spec 03's CHANGELOG.

### 1.2 Promote the `melee` enum package

Move:

```
internal/goslippi/slippi/melee/  →  pkg/slippi/melee/
```

Update all importers (event handlers, `pkg/slippi/types/*`, `internal/stats/`).

Add a `melee/strings.go` (new file) implementing display-name lookups:

```go
package melee

// DisplayName returns the canonical Slippi/Melee display name for the stage.
// Returns "Unknown Stage" for stage IDs not present in the table.
func (s Stage) DisplayName() string { ... }

// DisplayName returns the canonical display name for the external character ID.
// Returns "Unknown" for character IDs not present in the table.
func (c ExternalCharacterID) DisplayName() string { ... }

// DisplayName returns the canonical display name for the internal character ID.
// Returns "Unknown" for character IDs not present in the table.
func (c InternalCharacterID) DisplayName() string { ... }

// DisplayName returns the canonical display name for the item.
// Returns "" (empty string) for item IDs not present in the table.
func (i Item) DisplayName() string { ... }
```

Reference data (already aligned with `slippi-js`):

- **Stages:** copy verbatim from the DontGetHit backend table at
  [`backend/internal/slippi/melee/stages.go`](https://github.com/dontgethit/dontgethit-go/blob/main/internal/slippi/melee/stages.go)
  (already verified slippi-js-parity per dontgethit Phase 1; covers IDs 2–32 plus 82, 83, 84, 85, 285, and the Target Test stages 33–58). Use the `StageNames` map there as the source — same display strings, including `Pokémon Stadium`, `Princess Peach's Castle`, `Poké Floats`, `Mr. Game & Watch`, `Home-Run Contest`, etc.
- **Characters:** Use slippi-js
  [`common/melee/characters.json`](https://github.com/project-slippi/slippi-js/blob/master/src/melee/characters.ts)
  for `ExternalCharacterID.DisplayName()`. The dontgethit table at
  [`backend/internal/slippi/melee/characters.go`](https://github.com/dontgethit/dontgethit-go/blob/main/internal/slippi/melee/characters.go)
  is also slippi-js-parity (note: ID 32 is `Popo`, ID 29 is `Gigabowser`).
- **Internal characters:** map IDs 0–31 (per `internal_character_ids.go` constants) to the same display name as the corresponding external character (use the table you already keep in `internal/stats/framedata_embed.go::internalCharToFramedataName` as a starting point; add `Master Hand`, `Crazy Hand`, `Wireframe (Male)`, `Wireframe (Female)`, `Gigabowser`, `Sandbag` for IDs 27–31).
- **Items:** display names are nice-to-have. If maintaining 200+ item entries is too much, ship the helper with only the most common items (`Capsule`, `Bob-omb`, `Banana Peel`, `Beam Sword`, `Hammer`, `Heart`, `Maxim Tomato`, `Mr. Saturn`, `Pokéball`, `Ray Gun`, `Star Rod`, `Super Mushroom`, `Poison Mushroom`, `Warpstar`, …) and let `DisplayName()` return `""` for the rest.

Add corresponding tests in `pkg/slippi/melee/strings_test.go` covering at least:

- All 28 unique stages between IDs 2 and 32 plus 82–85 and 285.
- All 26 external character IDs 0–25, plus `Popo` (32).
- A representative sample of internal characters (0, 1, 5, 10/11 = Popo/Nana, 31 = Sandbag).

### 1.3 Promote the stats output types

Move:

```
internal/stats/stats.go            →  split into pkg/slippi/stats/types.go (the public types) + internal/stats/compute.go (Compute() and unexported helpers stay internal)
internal/stats/framedata.json      →  stays in internal/stats/ (embedded as before)
internal/stats/framedata_embed.go  →  stays in internal/stats/
```

The split is along this line:

| New file | Contents |
|---|---|
| `pkg/slippi/stats/types.go` | `Ratio`, `Stock`, `MoveLanded`, `Combo`, `Conversion`, `ActionCounts` (with all nested anonymous structs flattened into named types: `AttackCount`, `GrabCount`, `ThrowCount`, `GroundTechCount`, `WallTechCount`, `EdgeCancelCount`, `LCancelCount`), `InputCounts`, `Overall`, `Stats`. |
| `internal/stats/compute.go` (renamed from current `stats.go` minus the type definitions) | `Compute(game types.Game) Stats` and every unexported helper (`processActions`, `processStocks`, …). |

Promote `Compute` to public by re-exporting it from a new top-level file `pkg/slippi/stats/compute.go`:

```go
package stats

import (
    types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
    impl "github.com/ethangamma24/slippi-go/internal/stats"
)

// Compute aggregates a Slippi game into the slippi-js-parity stats output.
func Compute(game types.Game) Stats {
    return impl.Compute(game)
}
```

> **About the nested anonymous structs:** the current `ActionCounts` definition
> uses anonymous nested structs (`AttackCount struct{Jab1 int ...}`). Moving these
> to a public package without naming them means external consumers can never
> declare a variable of those types, can never write a generic helper that takes
> `AttackCount`, and cannot mock/test them ergonomically. **Promote each nested
> struct to a named type** in `pkg/slippi/stats/types.go`:
>
> ```go
> type AttackCount struct { Jab1, Jab2, Jab3, Jabm, Dash, Ftilt, Utilt, Dtilt, Fsmash, Usmash, Dsmash, Nair, Fair, Bair, Uair, Dair int }
> type GrabCount struct { Success, Fail int }
> type ThrowCount struct { Up, Forward, Back, Down int }
> type GroundTechCount struct { Away, In, Neutral, Fail int }
> type WallTechCount struct { Success, Fail int }
> type EdgeCancelCount struct { Success, Slow int }
> type LCancelCount struct { Success, Fail int }
> ```
>
> Then `ActionCounts` references them by name. JSON tags must match the existing
> field names so the `(any, error)` getters keep emitting identical JSON
> (`"attackCount": {"jab1": ...}`).

### 1.4 Update the `Game` facade

`pkg/slippi/game.go` currently embeds `nativeslippi.Game` and routes everything
through `toAny`. After the moves above, that file's imports change but its
behavior must not. After-state import block:

```go
import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sync"

    goslippi "github.com/ethangamma24/slippi-go/internal/goslippi"
    "github.com/ethangamma24/slippi-go/internal/realtime"
    "github.com/ethangamma24/slippi-go/internal/stats"
    types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)
```

`Game.parsed` becomes `types.Game`. All `(any, error)` methods continue to return
`(any, error)` for backward compatibility. **No new methods are added in this
spec.**

---

## Acceptance criteria

- [ ] `go test ./...` passes from a clean checkout (no skipped tests).
- [ ] `go test ./pkg/slippi -run TestPerformanceGate -count=1` still passes.
- [ ] `go test ./pkg/slippi -run TestStatsParityAgainstReferenceJS` still passes (treat any new mismatch as a regression).
- [ ] A new tiny smoke test under `pkg/slippi/types/` imports the package from outside its own folder (e.g. via a `_test` package) and exercises one of every promoted type's zero value.
- [ ] A new smoke test under `pkg/slippi/melee/strings_test.go` covers stage / character display names per the table in §1.2.
- [ ] A new smoke test under `pkg/slippi/stats/types_test.go` zero-values `stats.Stats`, `stats.ActionCounts`, and confirms JSON round-trip key shape (`{"attackCount":{"jab1":0,…}}`) is unchanged.
- [ ] `git grep -n LRASInitiatior` returns zero hits.
- [ ] `git grep -n -E 'package (slippi|melee|stats)' -- internal/` returns zero hits for the *type-bearing* slippi/melee/stats packages (the event handler packages stay internal under their own names).
- [ ] No new third-party dependencies in `go.mod`.

## Out of scope

- Move-name table (slippi-js `attackId → name`). The DontGetHit backend will keep its own table at `internal/slippi/melee/moves.go`. If you want to add it later it should be a follow-up spec.
- Deprecating the `(any, error)` getters on `pkg/slippi.Game`. They stay; the typed accessors added in spec 02 are additive.
- Adding `String()` (vs `DisplayName()`) — `DisplayName()` is the chosen name to keep `String()` available for future debug formatting if desired.

## References

- Current model: [`internal/goslippi/slippi/`](https://github.com/ethangamma24/slippi-go/tree/main/internal/goslippi/slippi)
- Current melee enums: [`internal/goslippi/slippi/melee/`](https://github.com/ethangamma24/slippi-go/tree/main/internal/goslippi/slippi/melee)
- Current stats: [`internal/stats/`](https://github.com/ethangamma24/slippi-go/tree/main/internal/stats)
- DontGetHit melee tables to copy display strings from: [`backend/internal/slippi/melee/`](https://github.com/dontgethit/dontgethit-go/tree/main/internal/slippi/melee)
- slippi-js source for parity reference: [`project-slippi/slippi-js`](https://github.com/project-slippi/slippi-js)
