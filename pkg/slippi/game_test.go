package slippi

import (
	"context"
	"testing"
)

func TestGameSettingsParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/test.slp", "settings", func(ctx context.Context, g *Game) (any, error) {
		return g.GetSettings(ctx)
	})
}

func TestGameMetadataParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/test.slp", "metadata", func(ctx context.Context, g *Game) (any, error) {
		return g.GetMetadata(ctx)
	})
}

func TestGameFramesParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/test.slp", "frames", func(ctx context.Context, g *Game) (any, error) {
		return g.GetFrames(ctx)
	})
}

func TestGameSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/unranked_game1.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
