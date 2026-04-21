package slippi

import (
	"context"
	"testing"
)

func TestMeleeSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/ntsc.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
