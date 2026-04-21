package slippi

import (
	"context"
	"testing"
)

func TestConversionSummaryParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/controllerFixes.slp", "summary", func(ctx context.Context, g *Game) (any, error) {
		return g.Summary(ctx)
	})
}
