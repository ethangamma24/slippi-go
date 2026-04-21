package slippi

import (
	"context"
	"testing"
)

func TestStagesSettingsParity(t *testing.T) {
	assertOperationParity(t, "testdata/slp/stadiumTransformations.slp", "settings", func(ctx context.Context, g *Game) (any, error) {
		return g.GetSettings(ctx)
	})
}
