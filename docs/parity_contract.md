# Behavior Contract

Comparison rules:

- Float tolerance is absolute-only: `abs <= 1e-9`.
- `undefined` and omitted fields are equivalent.
- `null` and omitted fields are different.
- Zero values (`0`, `""`, `false`) are different from missing fields.
- Array ordering is strict.
- Object/map key ordering is semantic (order-insensitive).
- Realtime semantics (event ordering and partial-read behavior) must match.
- Error behavior must preserve condition/category/message intent.

## Stats behavior (`GetStats`)

- Stats are computed in [`internal/stats`](internal/stats) and exposed through [`pkg/slippi.Game.GetStats`](pkg/slippi/game.go).
- Frame iteration processes consecutive frames from `Frames.FIRST` (`-123`) and stops at the first missing or incomplete frame.
- Aerial landing-lag data for edge-cancel behavior is embedded as [`internal/stats/framedata.json`](internal/stats/framedata.json).
- Joystick-region comparisons use `float64` precision so boundary values like `0.2875f32` are handled consistently.
- Opening-type ratios (`counterHitRatio`, `neutralWinRatio`, `beneficialTradeRatio`) group conversions by `conversion.moves[0].playerIndex`; conversions with no moves are excluded.
- Conversion/combo/move state is tracked by slice index rather than raw pointer, so append-driven reallocations cannot invalidate in-progress updates.
- `TestStatsParityAgainstReferenceJS` is the comparator gate; mismatches are treated as implementation bugs, not test normalization.
