package slippi

import (
	"context"
	"testing"
)

func TestPlacingsSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/placementsTest/ffa_1p2p4p_winner_4p.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
