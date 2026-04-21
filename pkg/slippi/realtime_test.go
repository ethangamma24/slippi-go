package slippi

import (
	"context"
	"testing"
)

func TestRealtimeLatestFrameParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/realtimeTest.slp", "latestFrame", func(ctx context.Context, g *Game) (any, error) {
		return g.GetLatestFrame(ctx)
	})
}
