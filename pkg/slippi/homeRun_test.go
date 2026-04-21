package slippi

import (
	"context"
	"testing"
)

func TestHomeRunStatsParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/homeRun_positive.slp", "stats", func(ctx context.Context, g *Game) (any, error) {
		return g.GetStats(ctx)
	})
}
