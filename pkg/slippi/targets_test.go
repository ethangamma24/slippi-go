package slippi

import (
	"context"
	"testing"
)

func TestTargetsSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/BTTDK.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
