package stats

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed framedata.json
var framedataJSON []byte

// framedataRoot is keyed by character display name (e.g. "Fox") then move key (e.g. "fair").
// We only read landingLag / lcancelledLandingLag for aerial landing parity.
type framedataMove struct {
	LandingLag           int `json:"landingLag"`
	LcancelledLandingLag int `json:"lcancelledLandingLag"`
}

var framedataByChar map[string]map[string]framedataMove

func init() {
	var root map[string]map[string]framedataMove
	if err := json.Unmarshal(framedataJSON, &root); err != nil {
		panic(fmt.Sprintf("pkg/slippi/stats: decode framedata.json: %v", err))
	}
	framedataByChar = root
}

var internalCharToFramedataName = map[int]string{
	0x00: "Mario",
	0x01: "Fox",
	0x02: "Captain Falcon",
	0x03: "Donkey Kong",
	0x04: "Kirby",
	0x05: "Bowser",
	0x06: "Link",
	0x07: "Sheik",
	0x08: "Ness",
	0x09: "Peach",
	0x0a: "Popo",
	0x0b: "Nana",
	0x0c: "Pikachu",
	0x0d: "Samus",
	0x0e: "Yoshi",
	0x0f: "Jigglypuff",
	0x10: "Mewtwo",
	0x11: "Luigi",
	0x12: "Marth",
	0x13: "Zelda",
	0x14: "Young Link",
	0x15: "Dr. Mario",
	0x16: "Falco",
	0x17: "Pichu",
	0x18: "Mr. Game & Watch",
	0x19: "Ganondorf",
	0x1a: "Roy",
}

func getAerialLandingLags(internalCharID int, landingAnim uint16) (landingLag, lcancelLag int, ok bool) {
	aerialName, okA := aerialLandingMoveName(landingAnim)
	if !okA {
		return 0, 0, false
	}
	charName, okC := internalCharToFramedataName[internalCharID]
	if !okC {
		return 0, 0, false
	}
	charMoves, okM := framedataByChar[charName]
	if !okM {
		return 0, 0, false
	}
	move, ok := charMoves[aerialName]
	if !ok || move.LandingLag == 0 {
		return 0, 0, false
	}
	return move.LandingLag, move.LcancelledLandingLag, true
}

// aerialLandingMoveName maps landing animation IDs to framedata.json aerial keys.
func aerialLandingMoveName(anim uint16) (string, bool) {
	switch anim {
	case 0x47:
		return "fair", true
	case 0x48:
		return "bair", true
	case 0x49:
		return "upair", true
	case 0x4a:
		return "dair", true
	case 0x46:
		return "nair", true
	default:
		return "", false
	}
}
