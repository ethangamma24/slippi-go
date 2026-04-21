package slippi

import (
	"context"
	"testing"
)

func TestItemsStatsParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/itemExport.slp", "stats", func(ctx context.Context, g *Game) (any, error) {
		return g.GetStats(ctx)
	})
}
