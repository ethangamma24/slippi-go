package slippi_test

import (
	"testing"

	"github.com/ethangamma24/slippi-go/pkg/slippi/melee"
	types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

func TestPostFrameUpdateIsShieldStun(t *testing.T) {
	cases := []struct {
		name          string
		actionStateID uint16
		want          bool
	}{
		{"GuardOn", melee.StateGuardOn, false},
		{"Guard", melee.StateGuard, false},
		{"GuardDamage", melee.StateGuardDamage, true},
		{"GuardOff", melee.StateGuardOff, false},
		{"GuardSetOff", melee.StateGuardSetOff, false},
		{"NonShield 177", 177, false},
		{"NonShield 183", 183, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := types.PostFrameUpdate{ActionStateID: tc.actionStateID}
			if got := p.IsShieldStun(); got != tc.want {
				t.Errorf("PostFrameUpdate.IsShieldStun() = %v, want %v", got, tc.want)
			}
		})
	}
}
