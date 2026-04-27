package melee

var stageNames = map[Stage]string{
	StageFountainOfDreams:     "Fountain of Dreams",
	StagePokemonStadium:       "Pokémon Stadium",
	StagePrincessPeachsCastle: "Princess Peach's Castle",
	StageKongoJungle:          "Kongo Jungle",
	StageBrinstar:             "Brinstar",
	StageCorneria:             "Corneria",
	StageYoshisStory:          "Yoshi's Story",
	StageOnett:                "Onett",
	StageMuteCity:             "Mute City",
	StageRainbowCruise:        "Rainbow Cruise",
	StageJungleJapes:          "Jungle Japes",
	StageGreatBay:             "Great Bay",
	StageHyruleTemple:         "Hyrule Temple",
	StageBrinstarDepths:       "Brinstar Depths",
	StageYoshisIsland:         "Yoshi's Island",
	StageGreenGreens:          "Green Greens",
	StageFourside:             "Fourside",
	StageMushroomKingdomI:     "Mushroom Kingdom I",
	StageMushroomKingdomII:    "Mushroom Kingdom II",
	StageVenom:                "Venom",
	StagePokeFloats:           "Poké Floats",
	StageBigBlue:              "Big Blue",
	StageIcicleMountain:       "Icicle Mountain",
	StageICETOP:               "Icetop",
	StageFlatZone:             "Flat Zone",
	StageDreamLandN64:         "Dream Land N64",
	StageYoshisIslandN64:      "Yoshi's Island N64",
	StageKongoJungleN64:       "Kongo Jungle N64",
	StageBattlefield:          "Battlefield",
	StageFinalDestination:     "Final Destination",
	StageRaceToTheFinish:      "Race to the Finish",
	StageGrabTheTrophies:      "Grab the Trophies",
	StageHomeRunContest:       "Home-Run Contest",
	StageAllStarLobby:         "All-Star Lobby",
	StageMultiManMelee:        "Multi-Man Melee",
}

// DisplayName returns the canonical Slippi/Melee display name for the stage.
// Returns "Unknown Stage" for stage IDs not present in the table.
func (s Stage) DisplayName() string {
	if name, ok := stageNames[s]; ok {
		return name
	}
	return "Unknown Stage"
}

var externalCharacterNames = map[ExternalCharacterID]string{
	Ext_CaptainFalcon:   "Captain Falcon",
	Ext_DonkeyKong:      "Donkey Kong",
	Ext_Fox:             "Fox",
	Ext_GameAndWatch:    "Mr. Game & Watch",
	Ext_Kirby:           "Kirby",
	Ext_Bowser:          "Bowser",
	Ext_Link:            "Link",
	Ext_Luigi:           "Luigi",
	Ext_Mario:           "Mario",
	Ext_Marth:           "Marth",
	Ext_Mewtwo:          "Mewtwo",
	Ext_Ness:            "Ness",
	Ext_Peach:           "Peach",
	Ext_Pikachu:         "Pikachu",
	Ext_IceClimbers:     "Ice Climbers",
	Ext_Jigglypuff:      "Jigglypuff",
	Ext_Samus:           "Samus",
	Ext_Yoshi:           "Yoshi",
	Ext_Zelda:           "Zelda",
	Ext_Sheik:           "Sheik",
	Ext_Falco:           "Falco",
	Ext_YoungLink:       "Young Link",
	Ext_DrMario:         "Dr. Mario",
	Ext_Roy:             "Roy",
	Ext_Pichu:           "Pichu",
	Ext_Ganondorf:       "Ganondorf",
	Ext_MasterHand:      "Master Hand",
	Ext_WireframeMale:   "Wireframe (Male)",
	Ext_WireframeFemale: "Wireframe (Female)",
	Ext_GigaBowser:      "Gigabowser",
	Ext_CrazyHand:       "Crazy Hand",
	Ext_Sandbag:         "Sandbag",
	Ext_Popo:            "Popo",
}

// DisplayName returns the canonical display name for the external character ID.
// Returns "Unknown" for character IDs not present in the table.
func (c ExternalCharacterID) DisplayName() string {
	if name, ok := externalCharacterNames[c]; ok {
		return name
	}
	return "Unknown"
}

var internalCharacterNames = map[InternalCharacterID]string{
	Int_Mario:           "Mario",
	Int_Fox:             "Fox",
	Int_CaptainFalcon:   "Captain Falcon",
	Int_DonkeyKong:      "Donkey Kong",
	Int_Kirby:           "Kirby",
	Int_Bowser:          "Bowser",
	Int_Link:            "Link",
	Int_Sheik:           "Sheik",
	Int_Ness:            "Ness",
	Int_Peach:           "Peach",
	Int_Popo:            "Popo",
	Int_Nana:            "Nana",
	Int_Pikachu:         "Pikachu",
	Int_Samus:           "Samus",
	Int_Yoshi:           "Yoshi",
	Int_Jigglypuff:      "Jigglypuff",
	Int_Mewtwo:          "Mewtwo",
	Int_Luigi:           "Luigi",
	Int_Marth:           "Marth",
	Int_Zelda:           "Zelda",
	Int_YoungLink:       "Young Link",
	Int_DrMario:         "Dr. Mario",
	Int_Falco:           "Falco",
	Int_Pichu:           "Pichu",
	Int_GameAndWatch:    "Mr. Game & Watch",
	Int_Ganondorf:       "Ganondorf",
	Int_Roy:             "Roy",
	Int_MasterHand:      "Master Hand",
	Int_CrazyHand:       "Crazy Hand",
	Int_WireFrameMale:   "Wireframe (Male)",
	Int_WireFrameFemale: "Wireframe (Female)",
	Int_GigaBowser:      "Gigabowser",
	Int_Sandbag:         "Sandbag",
}

// DisplayName returns the canonical display name for the internal character ID.
// Returns "Unknown" for character IDs not present in the table.
func (c InternalCharacterID) DisplayName() string {
	if name, ok := internalCharacterNames[c]; ok {
		return name
	}
	return "Unknown"
}

var itemNames = map[Item]string{
	ItemCapsule:        "Capsule",
	ItemBox:            "Box",
	ItemBarrel:         "Barrel",
	ItemEgg:            "Egg",
	ItemPartyBall:      "Party Ball",
	ItemBarrelCannon:   "Barrel Cannon",
	ItemBobOmb:         "Bob-omb",
	ItemMrSaturn:       "Mr. Saturn",
	ItemHeartContainer: "Heart Container",
	ItemMaximTomato:    "Maxim Tomato",
	ItemStarman:        "Starman",
	ItemHomeRunBat:     "Home-Run Bat",
	ItemBeamSword:      "Beam Sword",
	ItemParasol:        "Parasol",
	ItemGreenShell:     "Green Shell",
	ItemRedShell:       "Red Shell",
	ItemRayGun:         "Ray Gun",
	ItemFreezie:        "Freezie",
	ItemFood:           "Food",
	ItemProximityMine:  "Proximity Mine",
	ItemFlipper:        "Flipper",
	ItemSuperScope:     "Super Scope",
	ItemStarRod:        "Star Rod",
	ItemLipsStick:      "Lip's Stick",
	ItemFan:            "Fan",
	ItemFireFlower:     "Fire Flower",
	ItemSuperMushroom:  "Super Mushroom",
	ItemMiniMushroom:   "Mini Mushroom",
	ItemWarpStar:       "Warp Star",
	ItemScrewAttack:    "Screw Attack",
	ItemBunnyHood:      "Bunny Hood",
	ItemMetalBox:       "Metal Box",
	ItemCloakingDevice: "Cloaking Device",
	ItemPokeBall:       "Poké Ball",
}

// DisplayName returns the canonical display name for the item.
// Returns "" (empty string) for item IDs not present in the table.
func (i Item) DisplayName() string {
	if name, ok := itemNames[i]; ok {
		return name
	}
	return ""
}
