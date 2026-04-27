package slippi_test

import (
	"testing"

	types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

func TestZeroValues(t *testing.T) {
	_ = types.Game{}
	_ = types.Data{}
	_ = types.GameStart{}
	_ = types.Player{}
	_ = types.GameEnd{}
	_ = types.PlayerPlacement{}
	_ = types.Frame{}
	_ = types.PreFrameUpdate{}
	_ = types.PostFrameUpdate{}
	_ = types.SelfInducedSpeeds{}
	_ = types.ItemUpdate{}
	_ = types.FrameStart{}
	_ = types.FrameBookend{}
	_ = types.Metadata{}
	_ = types.PlayerMeta{}
	_ = types.PlayersMeta{}
	_ = types.Names{}
	_ = types.Character{}
	_ = types.Characters{}
}
