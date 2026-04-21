package slippi

import (
	"context"
	"testing"
)

func TestSlpReaderFramesParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/test.slp", "frames", func(ctx context.Context, g *Game) (any, error) {
		return g.GetFrames(ctx)
	})
}
