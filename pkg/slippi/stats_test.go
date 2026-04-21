package slippi

import (
	"context"
	"testing"
)

func TestStatsParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/lCancel.slp", "stats", func(ctx context.Context, g *Game) (any, error) {
		return g.GetStats(ctx)
	})
}
