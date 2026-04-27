package melee

import "testing"

func TestStageDisplayName(t *testing.T) {
	tests := []struct {
		stage Stage
		want  string
	}{
		{StageFountainOfDreams, "Fountain of Dreams"},
		{StagePokemonStadium, "Pokémon Stadium"},
		{StageFinalDestination, "Final Destination"},
		{StageBattlefield, "Battlefield"},
		{StageDreamLandN64, "Dream Land N64"},
		{StageHomeRunContest, "Home-Run Contest"},
		{StageMultiManMelee, "Multi-Man Melee"},
		{Stage(999), "Unknown Stage"},
	}
	for _, tt := range tests {
		got := tt.stage.DisplayName()
		if got != tt.want {
			t.Errorf("Stage(%d).DisplayName() = %q, want %q", tt.stage, got, tt.want)
		}
	}
}

func TestExternalCharacterDisplayName(t *testing.T) {
	tests := []struct {
		char ExternalCharacterID
		want string
	}{
		{Ext_CaptainFalcon, "Captain Falcon"},
		{Ext_Fox, "Fox"},
		{Ext_Marth, "Marth"},
		{Ext_Popo, "Popo"},
		{Ext_Ganondorf, "Ganondorf"},
		{ExternalCharacterID(99), "Unknown"},
	}
	for _, tt := range tests {
		got := tt.char.DisplayName()
		if got != tt.want {
			t.Errorf("ExternalCharacterID(%d).DisplayName() = %q, want %q", tt.char, got, tt.want)
		}
	}
}

func TestInternalCharacterDisplayName(t *testing.T) {
	tests := []struct {
		char InternalCharacterID
		want string
	}{
		{Int_Mario, "Mario"},
		{Int_Fox, "Fox"},
		{Int_Bowser, "Bowser"},
		{Int_Popo, "Popo"},
		{Int_Nana, "Nana"},
		{Int_GigaBowser, "Gigabowser"},
		{Int_Sandbag, "Sandbag"},
		{InternalCharacterID(99), "Unknown"},
	}
	for _, tt := range tests {
		got := tt.char.DisplayName()
		if got != tt.want {
			t.Errorf("InternalCharacterID(%d).DisplayName() = %q, want %q", tt.char, got, tt.want)
		}
	}
}

func TestItemDisplayName(t *testing.T) {
	tests := []struct {
		item Item
		want string
	}{
		{ItemCapsule, "Capsule"},
		{ItemBobOmb, "Bob-omb"},
		{ItemBeamSword, "Beam Sword"},
		{ItemWarpStar, "Warp Star"},
		{Item(999), ""},
	}
	for _, tt := range tests {
		got := tt.item.DisplayName()
		if got != tt.want {
			t.Errorf("Item(%d).DisplayName() = %q, want %q", tt.item, got, tt.want)
		}
	}
}
