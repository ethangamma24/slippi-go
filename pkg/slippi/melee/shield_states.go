package melee

// Shield action state IDs from the standard Melee action state table.
const (
	StateGuardOn     uint16 = 178
	StateGuard       uint16 = 179
	StateGuardDamage uint16 = 180
	StateGuardOff    uint16 = 181
	StateGuardSetOff uint16 = 182
)

// IsShieldState returns true if the action state is any shield-related state
// (GuardOn, Guard, GuardDamage, GuardOff, GuardSetOff).
func IsShieldState(state uint16) bool {
	return state >= StateGuardOn && state <= StateGuardSetOff
}

// IsShieldStun returns true only for GuardDamage (180), the direct signal
// that a character's shield was hit and is currently in shield stun.
func IsShieldStun(state uint16) bool {
	return state == StateGuardDamage
}
