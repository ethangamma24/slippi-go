# Changelog

All notable changes to this project will be documented in this file.

## [v0.1.0] — 2026-04-26

First semver-tagged release. The public API surface is now considered stable
within the v0.x line, with breaking changes reserved for future v0.(x+1).0 /
v1.0.0 releases.

### Added

- Public types under `github.com/ethangamma24/slippi-go/pkg/slippi/types` covering
  the entire decoded Slippi game (`Game`, `Data`, `GameStart`, `Player`, `GameEnd`,
  `PlayerPlacement`, `Frame`, `PreFrameUpdate`, `PostFrameUpdate`, `SelfInducedSpeeds`,
  `ItemUpdate`, `FrameStart`, `FrameBookend`, `Metadata`, `PlayerMeta`, `PlayersMeta`,
  `Names`, `Character`, `Characters`, plus all enums and `FirstFrame` /
  `FirstPlayableFrame` constants).
- Public Melee enums under `github.com/ethangamma24/slippi-go/pkg/slippi/melee`
  (`ExternalCharacterID`, `InternalCharacterID`, `Stage`, `EnabledItem`, `Item`).
- `(Stage|ExternalCharacterID|InternalCharacterID|Item).DisplayName() string`
  helpers returning slippi-js-parity display strings.
- Public stats output types under `github.com/ethangamma24/slippi-go/pkg/slippi/stats`
  (`Stats`, `Overall`, `ActionCounts`, `AttackCount`, `GrabCount`, `ThrowCount`,
  `GroundTechCount`, `WallTechCount`, `EdgeCancelCount`, `LCancelCount`, `Stock`,
  `Conversion`, `Combo`, `MoveLanded`, `Ratio`, `InputCounts`) plus the public
  `Compute(types.Game) Stats` entry point.
- `pkg/slippi.NewGameFromBytes(name string, data []byte) *Game`.
- `pkg/slippi.NewGameFromReader(name string, r io.Reader) (*Game, error)`.
- `pkg/slippi.ParseMetaFromBytes(name string, data []byte) (types.Metadata, error)`.
- Typed accessors on `*Game`: `Parsed`, `SettingsTyped`, `MetadataTyped`,
  `StatsTyped`, `FramesTyped`, `GameEndTyped`, `LatestFrameTyped`. These return
  the public typed values directly without the existing JSON-`any` round-trip.

### Changed

- **Field rename:** `pkg/slippi/types.GameEnd.LRASInitiatior` →
  `pkg/slippi/types.GameEnd.LRASInitiator` (typo fix). This also changes the JSON
  payload key emitted by the existing `(any, error)` getters from
  `"LRASInitiatior"` to `"LRASInitiator"`. Consumers of the JSON facade should
  update accordingly.
- The typed model and Melee enums moved out of `internal/`. The packages at
  `internal/goslippi/slippi/...` and `internal/goslippi/slippi/melee/...` are
  gone; their contents now live at `pkg/slippi/types` and `pkg/slippi/melee`.
  External code already could not import the `internal/` paths, so this is not
  a breaking change for external consumers.
- The `internal/stats` package's exported types now live at `pkg/slippi/stats`;
  `Compute` is exported there. Behavior unchanged.

### Unchanged (compatibility)

- `pkg/slippi.NewGame(fixturePath string)` and all existing `(any, error)`
  accessors (`GetSettings`, `GetMetadata`, `GetStats`, `GetFrames`,
  `GetLatestFrame`, `GetGameEnd`, `Summary`) keep their signatures and behavior,
  modulo the `LRASInitiatior → LRASInitiator` JSON key rename above.
- `TestStatsParityAgainstReferenceJS` and the performance gate continue to pass.
