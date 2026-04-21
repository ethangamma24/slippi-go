package slippi

import (
	"context"
	"testing"
)

func TestDoublesSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/doubles.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
