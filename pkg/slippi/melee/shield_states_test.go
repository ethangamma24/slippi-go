package melee

import "testing"

func TestIsShieldState(t *testing.T) {
	cases := []struct {
		state uint16
		want  bool
	}{
		{177, false},
		{StateGuardOn, true},
		{StateGuard, true},
		{StateGuardDamage, true},
		{StateGuardOff, true},
		{StateGuardSetOff, true},
		{183, false},
	}

	for _, tc := range cases {
		got := IsShieldState(tc.state)
		if got != tc.want {
			t.Errorf("IsShieldState(%d) = %v, want %v", tc.state, got, tc.want)
		}
	}
}

func TestIsShieldStun(t *testing.T) {
	cases := []struct {
		state uint16
		want  bool
	}{
		{177, false},
		{StateGuardOn, false},
		{StateGuard, false},
		{StateGuardDamage, true},
		{StateGuardOff, false},
		{StateGuardSetOff, false},
		{183, false},
	}

	for _, tc := range cases {
		got := IsShieldStun(tc.state)
		if got != tc.want {
			t.Errorf("IsShieldStun(%d) = %v, want %v", tc.state, got, tc.want)
		}
	}
}
