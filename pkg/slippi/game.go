package slippi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	goslippi "github.com/ethangamma24/slippi-go/internal/goslippi"
	"github.com/ethangamma24/slippi-go/internal/realtime"
	"github.com/ethangamma24/slippi-go/pkg/slippi/stats"
	types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

// Game provides a compatibility facade for replay inspection.
type Game struct {
	name        string
	fixturePath string
	data        []byte

	once     sync.Once
	parsed   types.Game
	parseErr error
}

// NewGame keeps its existing signature for backward compatibility.
func NewGame(fixturePath string) *Game {
	return &Game{name: filepath.Base(fixturePath), fixturePath: fixturePath}
}

// NewGameFromBytes returns a Game backed by an in-memory .slp byte slice.
// `name` is included in error messages (e.g. "parse game replay-123.slp: …").
func NewGameFromBytes(name string, data []byte) *Game {
	return &Game{name: name, data: data}
}

// NewGameFromReader fully reads r into memory and returns a Game backed by those bytes.
// Returns an error if the read fails. Pass a size-limited reader (e.g. io.LimitReader)
// if you need an upper bound.
func NewGameFromReader(name string, r io.Reader) (*Game, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read replay %s: %w", name, err)
	}
	return NewGameFromBytes(name, data), nil
}

// Parsed returns the entire decoded game in typed form. The returned value
// shares no mutable state with the Game; callers may freely retain it.
func (g *Game) Parsed(ctx context.Context) (types.Game, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return types.Game{}, err
	}
	return g.parsed, nil
}

// SettingsTyped returns the GameStart event payload in typed form.
func (g *Game) SettingsTyped(ctx context.Context) (types.GameStart, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return types.GameStart{}, err
	}
	return g.parsed.Data.GameStart, nil
}

// MetadataTyped returns the parsed UBJSON metadata block in typed form.
func (g *Game) MetadataTyped(ctx context.Context) (types.Metadata, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return types.Metadata{}, err
	}
	return g.parsed.Meta, nil
}

// StatsTyped returns slippi-js-parity stats in typed form.
func (g *Game) StatsTyped(ctx context.Context) (stats.Stats, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return stats.Stats{}, err
	}
	return stats.Compute(g.parsed), nil
}

// FramesTyped returns the per-frame data in typed form. The map is the same one
// held by the Game and SHOULD NOT be mutated by the caller.
func (g *Game) FramesTyped(ctx context.Context) (map[int]types.Frame, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return g.parsed.Data.Frames, nil
}

// GameEndTyped returns the GameEnd event payload in typed form.
func (g *Game) GameEndTyped(ctx context.Context) (types.GameEnd, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return types.GameEnd{}, err
	}
	return g.parsed.Data.GameEnd, nil
}

// LatestFrameTyped returns the highest-numbered finalized frame, plus a bool
// indicating whether such a frame exists.
func (g *Game) LatestFrameTyped(ctx context.Context) (types.Frame, bool, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return types.Frame{}, false, err
	}
	f, ok := realtime.LatestFrame(g.parsed.Data.Frames)
	return f, ok, nil
}

// GetSettings returns the GameStart event payload via the JSON-any facade.
func (g *Game) GetSettings(ctx context.Context) (any, error) {
	s, err := g.SettingsTyped(ctx)
	if err != nil {
		return nil, err
	}
	return toAny(s)
}

// GetMetadata returns the parsed UBJSON metadata block via the JSON-any facade.
func (g *Game) GetMetadata(ctx context.Context) (any, error) {
	m, err := g.MetadataTyped(ctx)
	if err != nil {
		return nil, err
	}
	return toAny(m)
}

// GetStats returns slippi-js-parity stats via the JSON-any facade.
func (g *Game) GetStats(ctx context.Context) (any, error) {
	s, err := g.StatsTyped(ctx)
	if err != nil {
		return nil, err
	}
	return toAny(s)
}

// GetFrames returns the per-frame data via the JSON-any facade.
func (g *Game) GetFrames(ctx context.Context) (any, error) {
	f, err := g.FramesTyped(ctx)
	if err != nil {
		return nil, err
	}
	return toAny(f)
}

// GetLatestFrame returns the highest-numbered finalized frame via the JSON-any facade.
func (g *Game) GetLatestFrame(ctx context.Context) (any, error) {
	f, ok, err := g.LatestFrameTyped(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return toAny(f)
}

// GetGameEnd returns the GameEnd event payload via the JSON-any facade.
func (g *Game) GetGameEnd(ctx context.Context) (any, error) {
	e, err := g.GameEndTyped(ctx)
	if err != nil {
		return nil, err
	}
	return toAny(e)
}

// Summary returns a composite map of settings, metadata, stats, gameEnd and latestFrame.
func (g *Game) Summary(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	settings, _ := g.GetSettings(ctx)
	metadata, _ := g.GetMetadata(ctx)
	st, _ := g.GetStats(ctx)
	gameEnd, _ := g.GetGameEnd(ctx)
	latestFrame, _ := g.GetLatestFrame(ctx)
	return map[string]any{
		"settings":    settings,
		"metadata":    metadata,
		"stats":       st,
		"gameEnd":     gameEnd,
		"latestFrame": latestFrame,
	}, nil
}

func (g *Game) ensureParsed(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.once.Do(func() {
		switch {
		case len(g.data) > 0:
			g.parsed, g.parseErr = goslippi.ParseGameFromBytes(g.name, g.data)
		case g.fixturePath != "":
			absPath, err := resolveFixturePath(g.fixturePath)
			if err != nil {
				g.parseErr = err
				return
			}
			g.parsed, g.parseErr = goslippi.ParseGame(absPath)
		default:
			g.parseErr = fmt.Errorf("slippi: Game has no source (use NewGame, NewGameFromBytes, or NewGameFromReader)")
		}
	})
	if g.parseErr != nil {
		return fmt.Errorf("parse game %s: %w", g.name, g.parseErr)
	}
	return nil
}

func resolveFixturePath(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	root, err := projectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, filepath.FromSlash(p)), nil
}

func projectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod from %s", dir)
		}
		dir = parent
	}
}

func toAny(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return out, nil
}
