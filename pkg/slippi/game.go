package slippi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	goslippi "github.com/ethangamma24/slippi-go/internal/goslippi"
	nativeslippi "github.com/ethangamma24/slippi-go/internal/goslippi/slippi"
	"github.com/ethangamma24/slippi-go/internal/realtime"
	nativestats "github.com/ethangamma24/slippi-go/internal/stats"
)

// Game provides a compatibility facade for replay inspection.
type Game struct {
	fixturePath string
	once        sync.Once
	parsed      nativeslippi.Game
	parseErr    error
}

func NewGame(fixturePath string) *Game {
	return &Game{fixturePath: fixturePath}
}

func (g *Game) GetSettings(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return toAny(g.parsed.Data.GameStart)
}

func (g *Game) GetMetadata(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return toAny(g.parsed.Meta)
}

func (g *Game) GetStats(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return toAny(nativestats.Compute(g.parsed))
}

func (g *Game) GetFrames(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return toAny(g.parsed.Data.Frames)
}

func (g *Game) GetLatestFrame(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	var latest int
	for k := range g.parsed.Data.Frames {
		if k > latest {
			latest = k
		}
	}
	if latest == 0 {
		return nil, nil
	}
	frame, ok := realtime.LatestFrame(g.parsed.Data.Frames)
	if !ok {
		return nil, nil
	}
	return toAny(frame)
}

func (g *Game) GetGameEnd(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	return toAny(g.parsed.Data.GameEnd)
}

func (g *Game) Summary(ctx context.Context) (any, error) {
	if err := g.ensureParsed(ctx); err != nil {
		return nil, err
	}
	settings, _ := g.GetSettings(ctx)
	metadata, _ := g.GetMetadata(ctx)
	stats, _ := g.GetStats(ctx)
	gameEnd, _ := g.GetGameEnd(ctx)
	latestFrame, _ := g.GetLatestFrame(ctx)
	return map[string]any{
		"settings":    settings,
		"metadata":    metadata,
		"stats":       stats,
		"gameEnd":     gameEnd,
		"latestFrame": latestFrame,
	}, nil
}

func (g *Game) ensureParsed(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	g.once.Do(func() {
		absPath, err := resolveFixturePath(g.fixturePath)
		if err != nil {
			g.parseErr = err
			return
		}
		g.parsed, g.parseErr = goslippi.ParseGame(absPath)
	})
	if g.parseErr != nil {
		return fmt.Errorf("parse game %s: %w", g.fixturePath, g.parseErr)
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
