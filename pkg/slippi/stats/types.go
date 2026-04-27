package stats

type Ratio struct {
	Count float64  `json:"count"`
	Total float64  `json:"total"`
	Ratio *float64 `json:"ratio,omitempty"`
}

type Stock struct {
	PlayerIndex    int      `json:"playerIndex"`
	StartFrame     int      `json:"startFrame"`
	EndFrame       *int     `json:"endFrame,omitempty"`
	StartPercent   float64  `json:"startPercent"`
	CurrentPercent float64  `json:"currentPercent"`
	EndPercent     *float64 `json:"endPercent,omitempty"`
	Count          int      `json:"count"`
	DeathAnimation *int     `json:"deathAnimation,omitempty"`
}

type MoveLanded struct {
	PlayerIndex int     `json:"playerIndex"`
	Frame       int     `json:"frame"`
	MoveID      int     `json:"moveId"`
	HitCount    int     `json:"hitCount"`
	Damage      float64 `json:"damage"`
}

type Combo struct {
	PlayerIndex    int          `json:"playerIndex"`
	StartFrame     int          `json:"startFrame"`
	EndFrame       *int         `json:"endFrame,omitempty"`
	StartPercent   float64      `json:"startPercent"`
	CurrentPercent float64      `json:"currentPercent"`
	EndPercent     *float64     `json:"endPercent,omitempty"`
	Moves          []MoveLanded `json:"moves"`
	DidKill        bool         `json:"didKill"`
	LastHitBy      *int         `json:"lastHitBy,omitempty"`
}

type Conversion struct {
	Combo
	OpeningType string `json:"openingType"`
}

type AttackCount struct {
	Jab1   int `json:"jab1"`
	Jab2   int `json:"jab2"`
	Jab3   int `json:"jab3"`
	Jabm   int `json:"jabm"`
	Dash   int `json:"dash"`
	Ftilt  int `json:"ftilt"`
	Utilt  int `json:"utilt"`
	Dtilt  int `json:"dtilt"`
	Fsmash int `json:"fsmash"`
	Usmash int `json:"usmash"`
	Dsmash int `json:"dsmash"`
	Nair   int `json:"nair"`
	Fair   int `json:"fair"`
	Bair   int `json:"bair"`
	Uair   int `json:"uair"`
	Dair   int `json:"dair"`
}

type GrabCount struct {
	Success int `json:"success"`
	Fail    int `json:"fail"`
}

type ThrowCount struct {
	Up      int `json:"up"`
	Forward int `json:"forward"`
	Back    int `json:"back"`
	Down    int `json:"down"`
}

type GroundTechCount struct {
	Away    int `json:"away"`
	In      int `json:"in"`
	Neutral int `json:"neutral"`
	Fail    int `json:"fail"`
}

type WallTechCount struct {
	Success int `json:"success"`
	Fail    int `json:"fail"`
}

type EdgeCancelCount struct {
	Success int `json:"success"`
	Slow    int `json:"slow"`
}

type LCancelCount struct {
	Success int `json:"success"`
	Fail    int `json:"fail"`
}

type ActionCounts struct {
	PlayerIndex int `json:"playerIndex"`

	WavedashCount  int `json:"wavedashCount"`
	WavelandCount  int `json:"wavelandCount"`
	AirDodgeCount  int `json:"airDodgeCount"`
	DashDanceCount int `json:"dashDanceCount"`
	SpotDodgeCount int `json:"spotDodgeCount"`
	LedgegrabCount int `json:"ledgegrabCount"`
	RollCount      int `json:"rollCount"`

	EdgeCancelCount EdgeCancelCount `json:"edgeCancelCount"`
	LCancelCount    LCancelCount    `json:"lCancelCount"`
	AttackCount     AttackCount     `json:"attackCount"`
	GrabCount       GrabCount       `json:"grabCount"`
	ThrowCount      ThrowCount      `json:"throwCount"`
	GroundTechCount GroundTechCount `json:"groundTechCount"`
	WallTechCount   WallTechCount   `json:"wallTechCount"`
}

type InputCounts struct {
	Buttons  int `json:"buttons"`
	Triggers int `json:"triggers"`
	Joystick int `json:"joystick"`
	CStick   int `json:"cstick"`
	Total    int `json:"total"`
}

type Overall struct {
	PlayerIndex            int         `json:"playerIndex"`
	InputCounts            InputCounts `json:"inputCounts"`
	ConversionCount        int         `json:"conversionCount"`
	TotalDamage            float64     `json:"totalDamage"`
	KillCount              int         `json:"killCount"`
	SuccessfulConversions  Ratio       `json:"successfulConversions"`
	InputsPerMinute        Ratio       `json:"inputsPerMinute"`
	DigitalInputsPerMinute Ratio       `json:"digitalInputsPerMinute"`
	OpeningsPerKill        Ratio       `json:"openingsPerKill"`
	DamagePerOpening       Ratio       `json:"damagePerOpening"`
	NeutralWinRatio        Ratio       `json:"neutralWinRatio"`
	CounterHitRatio        Ratio       `json:"counterHitRatio"`
	BeneficialTradeRatio   Ratio       `json:"beneficialTradeRatio"`
}

type Stats struct {
	GameComplete       bool           `json:"gameComplete"`
	LastFrame          int            `json:"lastFrame"`
	PlayableFrameCount int            `json:"playableFrameCount"`
	Stocks             []Stock        `json:"stocks"`
	Conversions        []Conversion   `json:"conversions"`
	Combos             []Combo        `json:"combos"`
	ActionCounts       []ActionCounts `json:"actionCounts"`
	Overall            []Overall      `json:"overall"`
}
