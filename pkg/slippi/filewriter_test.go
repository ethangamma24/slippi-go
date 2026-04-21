package slippi

import (
	"context"
	"testing"
)

func TestFilewriterPathwayParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/unranked_game2.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
