package stats

import (
	"math/bits"
	"sort"

	types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

const (
	punishResetFrames      = 45
	comboStringResetFrames = 45
)


type playerPair struct {
	PlayerIndex   uint8
	OpponentIndex uint8
}

type playerActionState struct {
	Counts            *ActionCounts
	Animations        []uint16
	ActionFrameCounts []float32
	PositionsY        []float32
	LastLCancelStatus uint8
}

type playerConversionState struct {
	// ConversionIdx and MoveIdx are indices into the outer slices (-1 = none). We cannot
	// retain raw *Conversion / *MoveLanded pointers because the backing arrays reallocate
	// when subsequent entries are appended, which would silently discard later updates
	// (for example, EndFrame/EndPercent when a conversion terminates well after it was
	// appended). the reference implementation sidesteps this by mutating heap-allocated objects; we emulate
	// that by resolving via index on every access.
	ConversionIdx    int
	MoveIdx          int
	ResetCounter     int
	LastHitAnimation *uint16
}

type playerComboState struct {
	ComboIdx         int
	MoveIdx          int
	ResetCounter     int
	LastHitAnimation *uint16
}

type playerStockState struct {
	Stock *Stock
}

type playerInputState struct {
	PlayerIndex        int
	OpponentIndex      int
	InputCount         int
	JoystickInputCount int
	CStickInputCount   int
	ButtonInputCount   int
	TriggerInputCount  int
}

func Compute(game types.Game) Stats {
	lastFrame := latestFrameNumber(game.Data.Frames)
	playableFrameCount := lastFrame - types.FirstPlayableFrame
	if playableFrameCount < 0 {
		playableFrameCount = 0
	}

	players := activePlayers(game.Data.GameStart.Players)
	perms := singlesPlayerPermutations(players)
	actionStates := newActionStates(perms)
	conversionStates := newConversionStates(perms)
	comboStates := newComboStates(perms)
	stockStates := newStockStates(perms)
	inputStates := newInputStates(perms)
	stocks := make([]Stock, 0, 32)
	conversions := make([]Conversion, 0, 64)
	combos := make([]Combo, 0, 64)

	requiredPlayers := make([]uint8, 0, len(players))
	for _, p := range players {
		requiredPlayers = append(requiredPlayers, p.Index)
	}
	for frameNum := types.FirstFrame; frameNum <= lastFrame; frameNum++ {
		frame, ok := game.Data.Frames[frameNum]
		if !ok {
			break
		}
		if !isCompletedFrame(frameNum, requiredPlayers, frame) {
			break
		}
		for _, pair := range perms {
			player, okPlayer := frame.Players[pair.PlayerIndex]
			opp, okOpp := frame.Players[pair.OpponentIndex]
			if !okPlayer || !okOpp {
				continue
			}
			processActions(actionStates[pair], player.Post, opp.Post)
			processStocks(game.Data.Frames, stockStates[pair], pair, player.Post, &stocks)
			processConversions(game.Data.Frames, conversionStates[pair], pair, player.Post, opp.Post, &conversions)
			processCombos(game.Data.Frames, comboStates[pair], pair, player.Post, opp.Post, &combos)
			processInputs(game.Data.Frames, inputStates[pair], pair, player.Pre)
		}
	}

	gameComplete := game.Data.GameEnd.GameEndMethod != types.GameEndUnresolved

	populateOpeningTypes(conversions)
	for i := range conversions {
		if conversions[i].OpeningType == "unknown" {
			// Still-open conversions have no endFrame in output; openingType should never remain "unknown"
			// once getStats returns. Default to neutral-win to match observed edge-tail behavior.
			conversions[i].OpeningType = "neutral-win"
		}
	}
	overall := generateOverall(players, inputStates, conversions, playableFrameCount, game.Data.GameStart.IsTeams)

	return Stats{
		GameComplete:       gameComplete,
		LastFrame:          lastFrame,
		PlayableFrameCount: playableFrameCount,
		Stocks:             stocks,
		Conversions:        conversions,
		Combos:             combos,
		ActionCounts:       collectActionCounts(perms, actionStates),
		Overall:            overall,
	}
}

func latestFrameNumber(frames map[int]types.Frame) int {
	latest := 0
	for k := range frames {
		if k > latest {
			latest = k
		}
	}
	return latest
}

func activePlayers(players []types.Player) []types.Player {
	out := make([]types.Player, 0, len(players))
	for _, p := range players {
		if p.PlayerType != types.PlayerTypeEmpty {
			out = append(out, p)
		}
	}
	return out
}

func singlesPlayerPermutations(players []types.Player) []playerPair {
	if len(players) != 2 {
		return nil
	}
	return []playerPair{
		{PlayerIndex: players[0].Index, OpponentIndex: players[1].Index},
		{PlayerIndex: players[1].Index, OpponentIndex: players[0].Index},
	}
}

func newActionStates(perms []playerPair) map[playerPair]*playerActionState {
	out := make(map[playerPair]*playerActionState, len(perms))
	for _, pair := range perms {
		counts := &ActionCounts{PlayerIndex: int(pair.PlayerIndex)}
		out[pair] = &playerActionState{Counts: counts}
	}
	return out
}

func newConversionStates(perms []playerPair) map[playerPair]*playerConversionState {
	out := make(map[playerPair]*playerConversionState, len(perms))
	for _, pair := range perms {
		out[pair] = &playerConversionState{ConversionIdx: -1, MoveIdx: -1}
	}
	return out
}

func newComboStates(perms []playerPair) map[playerPair]*playerComboState {
	out := make(map[playerPair]*playerComboState, len(perms))
	for _, pair := range perms {
		out[pair] = &playerComboState{ComboIdx: -1, MoveIdx: -1}
	}
	return out
}

func newStockStates(perms []playerPair) map[playerPair]*playerStockState {
	out := make(map[playerPair]*playerStockState, len(perms))
	for _, pair := range perms {
		out[pair] = &playerStockState{}
	}
	return out
}

func newInputStates(perms []playerPair) map[playerPair]*playerInputState {
	out := make(map[playerPair]*playerInputState, len(perms))
	for _, pair := range perms {
		out[pair] = &playerInputState{PlayerIndex: int(pair.PlayerIndex), OpponentIndex: int(pair.OpponentIndex)}
	}
	return out
}

// -- helpers and processors omitted here for brevity in patch readability --
// The implementation below follows the expected stats behavior used in tests.

func processActions(state *playerActionState, player, opponent types.PostFrameUpdate) {
	currentAnimation := player.ActionStateID
	state.Animations = append(state.Animations, currentAnimation)
	state.ActionFrameCounts = append(state.ActionFrameCounts, player.ActionStateFrameCounter)
	state.PositionsY = append(state.PositionsY, player.YPos)
	if len(state.Animations) < 2 {
		return
	}
	prevAnimation := state.Animations[len(state.Animations)-2]
	prevFrameCounter := state.ActionFrameCounts[len(state.ActionFrameCounts)-2]
	isNewAction := currentAnimation != prevAnimation || prevFrameCounter > player.ActionStateFrameCounter
	if !isNewAction {
		return
	}
	last3 := lastAnimations(state.Animations, 3)
	if len(last3) == 3 && last3[0] == 0x14 && last3[1] == 0x12 && last3[2] == 0x14 {
		state.Counts.DashDanceCount++
	}
	if currentAnimation == 0xe9 || currentAnimation == 0xea {
		state.Counts.RollCount++
	}
	if currentAnimation == 0xeb {
		state.Counts.SpotDodgeCount++
	}
	if currentAnimation == 0xec {
		state.Counts.AirDodgeCount++
	}
	if currentAnimation == 0xfc {
		state.Counts.LedgegrabCount++
	}
	isGrabbing := currentAnimation == 0xd4 || currentAnimation == 0xd6
	isGrabAction := currentAnimation > 0xd4 && currentAnimation <= 0xde && currentAnimation != 0xd6
	if (prevAnimation == 0xd4 || prevAnimation == 0xd6) && isGrabAction {
		state.Counts.GrabCount.Success++
	}
	if (prevAnimation == 0xd4 || prevAnimation == 0xd6) && !isGrabAction {
		state.Counts.GrabCount.Fail++
	}
	if currentAnimation == 0xd6 && prevAnimation == 0x32 {
		state.Counts.AttackCount.Dash--
	}
	if isGrabbing {
		// no-op to keep structure similar
	}
	incrementAttacks(state.Counts, currentAnimation, int(player.CharacterID))
	incrementThrows(state.Counts, currentAnimation)
	incrementTechs(state.Counts, currentAnimation, player, opponent)
	if currentAnimation >= 0x46 && currentAnimation <= 0x4a {
		if player.LCancelStatus == 1 {
			state.Counts.LCancelCount.Success++
		}
		if player.LCancelStatus == 2 {
			state.Counts.LCancelCount.Fail++
		}
	}
	if isAerialLanding(prevAnimation) && (currentAnimation == 0x1d || currentAnimation == 0xf5) {
		landingFrames := prevAnimationRunLength(state.Animations)
		landingLag, lCancelLag, okLag := getAerialLandingLags(int(player.CharacterID), prevAnimation)
		if okLag {
			if landingFrames < lCancelLag {
				state.Counts.EdgeCancelCount.Success++
			}
			if landingFrames >= lCancelLag && landingFrames < landingLag {
				state.Counts.EdgeCancelCount.Slow++
			}
			if landingFrames <= lCancelLag && state.LastLCancelStatus == 2 && state.Counts.LCancelCount.Fail > 0 {
				state.Counts.LCancelCount.Fail--
			}
		}
	}
	handleWavedash(state.Counts, state.Animations, state.PositionsY)
	if player.LCancelStatus > 0 {
		state.LastLCancelStatus = player.LCancelStatus
	}
}

func processStocks(frames map[int]types.Frame, state *playerStockState, pair playerPair, player types.PostFrameUpdate, stocks *[]Stock) {
	currentFrame := player.FrameNumber
	prevFrame := currentFrame - 1
	var prevPlayer *types.PostFrameUpdate
	if pf, ok := frames[prevFrame]; ok {
		if p, ok := pf.Players[pair.PlayerIndex]; ok {
			prevPlayer = &p.Post
		}
	}
	if state.Stock == nil {
		if isDead(player.ActionStateID) {
			return
		}
		s := Stock{
			PlayerIndex:    int(pair.PlayerIndex),
			StartFrame:     currentFrame,
			StartPercent:   0,
			CurrentPercent: 0,
			Count:          int(player.StocksRemaining),
		}
		*stocks = append(*stocks, s)
		state.Stock = &(*stocks)[len(*stocks)-1]
		return
	}
	if prevPlayer != nil && didLoseStock(player, *prevPlayer) {
		end := player.FrameNumber
		endPercent := float64(prevPlayer.Percent)
		death := int(player.ActionStateID)
		state.Stock.EndFrame = &end
		state.Stock.EndPercent = &endPercent
		state.Stock.DeathAnimation = &death
		state.Stock = nil
		return
	}
	state.Stock.CurrentPercent = float64(player.Percent)
}

func processConversions(frames map[int]types.Frame, state *playerConversionState, pair playerPair, player, opponent types.PostFrameUpdate, conversions *[]Conversion) {
	prevPlayer, prevOpponent := prevFrames(frames, pair, player.FrameNumber-1)
	oppState := opponent.ActionStateID
	oppDamaged := isDamaged(oppState)
	oppGrabbed := isGrabbed(oppState)
	oppCmdGrabbed := isCommandGrabbed(oppState)
	damageTaken := calcDamageTaken(opponent, prevOpponent)

	if state.LastHitAnimation != nil && (player.ActionStateID != *state.LastHitAnimation ||
		(prevPlayer != nil && player.ActionStateFrameCounter < prevPlayer.ActionStateFrameCounter)) {
		state.LastHitAnimation = nil
	}

	if oppDamaged || oppGrabbed || oppCmdGrabbed {
		if state.ConversionIdx < 0 {
			lastHitBy := int(pair.PlayerIndex)
			startPercent := 0.0
			if prevOpponent != nil {
				startPercent = float64(prevOpponent.Percent)
			}
			c := Conversion{
				Combo: Combo{
					PlayerIndex:    int(pair.OpponentIndex),
					LastHitBy:      &lastHitBy,
					StartFrame:     player.FrameNumber,
					StartPercent:   startPercent,
					CurrentPercent: float64(opponent.Percent),
					Moves:          []MoveLanded{},
				},
				OpeningType: "unknown",
			}
			*conversions = append(*conversions, c)
			state.ConversionIdx = len(*conversions) - 1
			state.MoveIdx = -1
		}
		conv := &(*conversions)[state.ConversionIdx]
		if damageTaken > 0 {
			if state.LastHitAnimation == nil {
				mv := MoveLanded{
					PlayerIndex: int(pair.PlayerIndex),
					Frame:       player.FrameNumber,
					MoveID:      int(player.LastHittingAttackID),
				}
				conv.Moves = append(conv.Moves, mv)
				state.MoveIdx = len(conv.Moves) - 1
			}
			if state.MoveIdx >= 0 && state.MoveIdx < len(conv.Moves) {
				mv := &conv.Moves[state.MoveIdx]
				mv.HitCount++
				mv.Damage += damageTaken
			}
			if prevPlayer != nil {
				a := prevPlayer.ActionStateID
				state.LastHitAnimation = &a
			}
		}
	}
	if state.ConversionIdx < 0 {
		return
	}
	conv := &(*conversions)[state.ConversionIdx]
	oppDidLoseStock := prevOpponent != nil && didLoseStock(opponent, *prevOpponent)
	if !oppDidLoseStock {
		conv.CurrentPercent = float64(opponent.Percent)
	}
	if oppDamaged || oppGrabbed || oppCmdGrabbed {
		state.ResetCounter = 0
	}
	oppInControl := isInControl(oppState)
	if (state.ResetCounter == 0 && oppInControl) || state.ResetCounter > 0 {
		state.ResetCounter++
	}
	shouldTerminate := false
	if oppDidLoseStock {
		conv.DidKill = true
		shouldTerminate = true
	}
	if state.ResetCounter > punishResetFrames {
		shouldTerminate = true
	}
	if shouldTerminate {
		end := player.FrameNumber
		conv.EndFrame = &end
		endPercent := 0.0
		if prevOpponent != nil {
			endPercent = float64(prevOpponent.Percent)
		}
		conv.EndPercent = &endPercent
		state.ConversionIdx = -1
		state.MoveIdx = -1
	}
}

func processCombos(frames map[int]types.Frame, state *playerComboState, pair playerPair, player, opponent types.PostFrameUpdate, combos *[]Combo) {
	prevPlayer, prevOpponent := prevFrames(frames, pair, player.FrameNumber-1)
	oppState := opponent.ActionStateID
	oppDamaged := isDamaged(oppState)
	oppGrabbed := isGrabbed(oppState)
	oppCmdGrabbed := isCommandGrabbed(oppState)
	damageTaken := calcDamageTaken(opponent, prevOpponent)

	if state.LastHitAnimation != nil && (player.ActionStateID != *state.LastHitAnimation ||
		(prevPlayer != nil && player.ActionStateFrameCounter < prevPlayer.ActionStateFrameCounter)) {
		state.LastHitAnimation = nil
	}
	if oppDamaged || oppGrabbed || oppCmdGrabbed {
		if state.ComboIdx < 0 {
			lastHitBy := int(pair.PlayerIndex)
			startPercent := 0.0
			if prevOpponent != nil {
				startPercent = float64(prevOpponent.Percent)
			}
			c := Combo{
				PlayerIndex:    int(pair.OpponentIndex),
				StartFrame:     player.FrameNumber,
				StartPercent:   startPercent,
				CurrentPercent: float64(opponent.Percent),
				Moves:          []MoveLanded{},
				LastHitBy:      &lastHitBy,
			}
			*combos = append(*combos, c)
			state.ComboIdx = len(*combos) - 1
			state.MoveIdx = -1
		}
		cb := &(*combos)[state.ComboIdx]
		if damageTaken > 0 {
			if state.LastHitAnimation == nil {
				mv := MoveLanded{
					PlayerIndex: int(pair.PlayerIndex),
					Frame:       player.FrameNumber,
					MoveID:      int(player.LastHittingAttackID),
				}
				cb.Moves = append(cb.Moves, mv)
				state.MoveIdx = len(cb.Moves) - 1
			}
			if state.MoveIdx >= 0 && state.MoveIdx < len(cb.Moves) {
				mv := &cb.Moves[state.MoveIdx]
				mv.HitCount++
				mv.Damage += damageTaken
			}
			if prevPlayer != nil {
				a := prevPlayer.ActionStateID
				state.LastHitAnimation = &a
			}
		}
	}
	if state.ComboIdx < 0 {
		return
	}
	cb := &(*combos)[state.ComboIdx]
	oppTeching := isTeching(oppState)
	oppDown := isDown(oppState)
	oppDying := isDead(oppState)
	oppDidLoseStock := prevOpponent != nil && didLoseStock(opponent, *prevOpponent)
	if !oppDidLoseStock {
		cb.CurrentPercent = float64(opponent.Percent)
	}
	if oppDamaged || oppGrabbed || oppCmdGrabbed || oppTeching || oppDown || oppDying {
		state.ResetCounter = 0
	} else {
		state.ResetCounter++
	}
	shouldTerminate := false
	if oppDidLoseStock {
		cb.DidKill = true
		shouldTerminate = true
	}
	if state.ResetCounter > comboStringResetFrames {
		shouldTerminate = true
	}
	if shouldTerminate {
		end := player.FrameNumber
		cb.EndFrame = &end
		endPercent := 0.0
		if prevOpponent != nil {
			endPercent = float64(prevOpponent.Percent)
		}
		cb.EndPercent = &endPercent
		state.ComboIdx = -1
		state.MoveIdx = -1
	}
}

func processInputs(frames map[int]types.Frame, state *playerInputState, pair playerPair, player types.PreFrameUpdate) {
	// Must use pre.frame so prev-frame joins use the same frame source as input counting.
	currentFrame := player.FrameNumber
	prevFrame := currentFrame - 1
	if currentFrame < types.FirstPlayableFrame {
		return
	}
	pf, ok := frames[prevFrame]
	if !ok {
		return
	}
	prevPlayer, ok := pf.Players[pair.PlayerIndex]
	if !ok {
		return
	}
	prev := prevPlayer.Pre
	newPressed := risingEdgeButtonCount(prev.PhysicalButtons, player.PhysicalButtons)
	state.InputCount += newPressed
	state.ButtonInputCount += newPressed
	prevRegion := joystickRegion(prev.JoyStickX, prev.JoyStickY)
	currRegion := joystickRegion(player.JoyStickX, player.JoyStickY)
	if prevRegion != currRegion && currRegion != 0 {
		state.InputCount++
		state.JoystickInputCount++
	}
	prevC := joystickRegion(prev.CStickX, prev.CStickY)
	currC := joystickRegion(player.CStickX, player.CStickY)
	if prevC != currC && currC != 0 {
		state.InputCount++
		state.CStickInputCount++
	}
	if float64(prev.PhysicalTriggerL) < 0.3 && float64(player.PhysicalTriggerL) >= 0.3 {
		state.InputCount++
		state.TriggerInputCount++
	}
	if float64(prev.PhysicalTriggerR) < 0.3 && float64(player.PhysicalTriggerR) >= 0.3 {
		state.InputCount++
		state.TriggerInputCount++
	}
}

func prevFrames(frames map[int]types.Frame, pair playerPair, frame int) (*types.PostFrameUpdate, *types.PostFrameUpdate) {
	pf, ok := frames[frame]
	if !ok {
		return nil, nil
	}
	player, okP := pf.Players[pair.PlayerIndex]
	opp, okO := pf.Players[pair.OpponentIndex]
	if !okP || !okO {
		return nil, nil
	}
	p := player.Post
	o := opp.Post
	return &p, &o
}

func collectActionCounts(perms []playerPair, states map[playerPair]*playerActionState) []ActionCounts {
	out := make([]ActionCounts, 0, len(perms))
	for _, pair := range perms {
		if state, ok := states[pair]; ok && state.Counts != nil {
			out = append(out, *state.Counts)
		}
	}
	return out
}

func populateOpeningTypes(conversions []Conversion) {
	// Populate only unknown openings,
	// grouped by startFrame, ordered by startFrame, preserving per-group discovery order.
	lastEndFrameByOpp := map[int]int{}
	unknownIdx := make([]int, 0)
	for i := range conversions {
		if conversions[i].OpeningType == "unknown" {
			unknownIdx = append(unknownIdx, i)
		}
	}
	grouped := make(map[int][]int)
	for _, i := range unknownIdx {
		s := conversions[i].StartFrame
		grouped[s] = append(grouped[s], i)
	}
	starts := make([]int, 0, len(grouped))
	for k := range grouped {
		starts = append(starts, k)
	}
	sort.Ints(starts)
	for _, start := range starts {
		indices := grouped[start]
		isTrade := len(indices) >= 2
		for _, idx := range indices {
			c := &conversions[idx]
			if c.EndFrame != nil {
				lastEndFrameByOpp[c.PlayerIndex] = *c.EndFrame
			} else {
				delete(lastEndFrameByOpp, c.PlayerIndex)
			}
			if isTrade {
				c.OpeningType = "trade"
				continue
			}
			lastMovePlayer := c.PlayerIndex
			if len(c.Moves) > 0 {
				lastMovePlayer = c.Moves[0].PlayerIndex
			}
			oppEnd := lastEndFrameByOpp[lastMovePlayer]
			// isCounterAttack = oppEndFrame && oppEndFrame > conversion.startFrame
			// (0 is falsy, so a conversion ending on frame 0 is not a counter-attack)
			if oppEnd != 0 && oppEnd > c.StartFrame {
				c.OpeningType = "counter-attack"
			} else {
				c.OpeningType = "neutral-win"
			}
		}
	}
}

func generateOverall(players []types.Player, inputStates map[playerPair]*playerInputState, conversions []Conversion, playableFrameCount int, isTeams bool) []Overall {
	inputByPlayer := map[int]*playerInputState{}
	for _, state := range inputStates {
		inputByPlayer[state.PlayerIndex] = state
	}
	gameMinutes := float64(playableFrameCount) / 3600.0
	out := make([]Overall, 0, len(players))
	for _, p := range players {
		playerIndex := int(p.Index)
		input := inputByPlayer[playerIndex]
		inputCounts := InputCounts{}
		if input != nil {
			inputCounts = InputCounts{
				Buttons:  input.ButtonInputCount,
				Triggers: input.TriggerInputCount,
				CStick:   input.CStickInputCount,
				Joystick: input.JoystickInputCount,
				Total:    input.InputCount,
			}
		}
		conversionCount := 0
		successful := 0
		totalDamage := 0.0
		killCount := 0
		for _, c := range conversions {
			if c.PlayerIndex == playerIndex {
				continue
			}
			conversionCount++
			if c.DidKill && c.LastHitBy != nil && *c.LastHitBy == playerIndex {
				killCount++
			}
			if len(c.Moves) > 1 && c.Moves[0].PlayerIndex == playerIndex {
				successful++
			}
			for _, mv := range c.Moves {
				if mv.PlayerIndex == playerIndex {
					totalDamage += mv.Damage
				}
			}
		}
		opponents := make([]int, 0, len(players)-1)
		for _, opp := range players {
			if opp.Index == p.Index {
				continue
			}
			if isTeams && opp.TeamColour == p.TeamColour {
				continue
			}
			opponents = append(opponents, int(opp.Index))
		}
		out = append(out, Overall{
			PlayerIndex:            playerIndex,
			InputCounts:            inputCounts,
			ConversionCount:        conversionCount,
			TotalDamage:            totalDamage,
			KillCount:              killCount,
			SuccessfulConversions:  ratio(float64(successful), float64(conversionCount)),
			InputsPerMinute:        ratio(float64(inputCounts.Total), gameMinutes),
			DigitalInputsPerMinute: ratio(float64(inputCounts.Buttons), gameMinutes),
			OpeningsPerKill:        ratio(float64(conversionCount), float64(killCount)),
			DamagePerOpening:       ratio(totalDamage, float64(conversionCount)),
			NeutralWinRatio:        openingRatio(conversions, playerIndex, opponents, "neutral-win"),
			CounterHitRatio:        openingRatio(conversions, playerIndex, opponents, "counter-attack"),
			BeneficialTradeRatio:   beneficialTradeRatio(conversions, playerIndex, opponents),
		})
	}
	return out
}

func openingRatio(conversions []Conversion, player int, opponents []int, openingType string) Ratio {
	openings := 0
	opponentOpenings := 0
	for _, c := range conversions {
		if c.OpeningType != openingType {
			continue
		}
		// Group by conv.moves[0]?.playerIndex; conversions with no moves end
		// up keyed by `undefined` and are excluded from every player's bucket.
		if len(c.Moves) == 0 {
			continue
		}
		attacker := c.Moves[0].PlayerIndex
		if attacker == player {
			openings++
			continue
		}
		for _, opp := range opponents {
			if attacker == opp {
				opponentOpenings++
				break
			}
		}
	}
	return ratio(float64(openings), float64(openings+opponentOpenings))
}

func beneficialTradeRatio(conversions []Conversion, player int, opponents []int) Ratio {
	playerTrades := make([]Conversion, 0)
	opponentTrades := make([]Conversion, 0)
	for _, c := range conversions {
		if c.OpeningType != "trade" {
			continue
		}
		if len(c.Moves) == 0 {
			continue
		}
		attacker := c.Moves[0].PlayerIndex
		if attacker == player {
			playerTrades = append(playerTrades, c)
			continue
		}
		for _, opp := range opponents {
			if attacker == opp {
				opponentTrades = append(opponentTrades, c)
				break
			}
		}
	}
	benefits := 0
	limit := len(playerTrades)
	if len(opponentTrades) < limit {
		limit = len(opponentTrades)
	}
	for i := 0; i < limit; i++ {
		p := playerTrades[i]
		o := opponentTrades[i]
		pDamage := p.CurrentPercent - p.StartPercent
		oDamage := o.CurrentPercent - o.StartPercent
		if p.DidKill && !o.DidKill || pDamage > oDamage {
			benefits++
		}
	}
	return ratio(float64(benefits), float64(len(playerTrades)))
}

func ratio(count, total float64) Ratio {
	var r *float64
	if total != 0 {
		v := count / total
		r = &v
	}
	return Ratio{Count: count, Total: total, Ratio: r}
}

func calcDamageTaken(frame types.PostFrameUpdate, prev *types.PostFrameUpdate) float64 {
	if prev == nil {
		return 0
	}
	return float64(frame.Percent) - float64(prev.Percent)
}

func didLoseStock(frame, prev types.PostFrameUpdate) bool {
	return int(prev.StocksRemaining)-int(frame.StocksRemaining) > 0
}

func isInControl(state uint16) bool {
	ground := state >= 0x0e && state <= 0x18
	squat := state >= 0x27 && state <= 0x29
	groundAttack := state > 0x2c && state <= 0x40
	isGrab := state == 0xd4
	return ground || squat || groundAttack || isGrab
}

func isTeching(state uint16) bool { return state >= 0xc7 && state <= 0xcc }
func isDown(state uint16) bool    { return state >= 0xb7 && state <= 0xc6 }
func isDamaged(state uint16) bool {
	return (state >= 0x4b && state <= 0x5b) || state == 0x26 || state == 0xb9 || state == 0xc1
}
func isGrabbed(state uint16) bool { return state >= 0xdf && state <= 0xe8 }
func isCommandGrabbed(state uint16) bool {
	return (((state >= 0x10a && state <= 0x130) || (state >= 0x147 && state <= 0x152)) && state != 0x125)
}
func isDead(state uint16) bool { return state >= 0x00 && state <= 0x0a }

func lastAnimations(animations []uint16, n int) []uint16 {
	if len(animations) <= n {
		return animations
	}
	return animations[len(animations)-n:]
}

func isCompletedFrame(frameNum int, players []uint8, frame types.Frame) bool {
	// Require Post data for each tracked player frame; also require coherent frame indexing for
	// frame index so we never run stats on half-written PlayerFrameUpdate entries.
	for _, player := range players {
		pl, ok := frame.Players[player]
		if !ok {
			return false
		}
		if pl.Post.FrameNumber != frameNum {
			return false
		}
	}
	return true
}

func isAerialLanding(animation uint16) bool {
	return animation >= 0x46 && animation <= 0x4a
}

func incrementThrows(counts *ActionCounts, anim uint16) {
	switch anim {
	case 0xdd:
		counts.ThrowCount.Up++
	case 0xdb:
		counts.ThrowCount.Forward++
	case 0xde:
		counts.ThrowCount.Down++
	case 0xdc:
		counts.ThrowCount.Back++
	}
}

func incrementAttacks(counts *ActionCounts, anim uint16, charID int) {
	switch anim {
	case 0x2c:
		counts.AttackCount.Jab1++
	case 0x2d:
		counts.AttackCount.Jab2++
	case 0x2e:
		counts.AttackCount.Jab3++
	case 0x2f:
		counts.AttackCount.Jabm++
	case 0x32:
		counts.AttackCount.Dash++
	case 0x38:
		counts.AttackCount.Utilt++
	case 0x39:
		counts.AttackCount.Dtilt++
	case 0x3f:
		counts.AttackCount.Usmash++
	case 0x40:
		counts.AttackCount.Dsmash++
	case 0x41:
		counts.AttackCount.Nair++
	case 0x42:
		counts.AttackCount.Fair++
	case 0x43:
		counts.AttackCount.Bair++
	case 0x44:
		counts.AttackCount.Uair++
	case 0x45:
		counts.AttackCount.Dair++
	}
	if anim >= 0x33 && anim <= 0x37 {
		counts.AttackCount.Ftilt++
	}
	if anim >= 0x3a && anim <= 0x3e {
		counts.AttackCount.Fsmash++
	}
	if charID == 0x18 {
		switch anim {
		case 0x155:
			counts.AttackCount.Jab1++
		case 0x156:
			counts.AttackCount.Jabm++
		case 0x159:
			counts.AttackCount.Dtilt++
		case 0x15a:
			counts.AttackCount.Fsmash++
		case 0x15b:
			counts.AttackCount.Nair++
		case 0x15c:
			counts.AttackCount.Bair++
		case 0x15d:
			counts.AttackCount.Uair++
		}
	}
	if charID == 0x09 && (anim == 0x15d || anim == 0x15e || anim == 0x15f) {
		counts.AttackCount.Fsmash++
	}
}

func incrementTechs(counts *ActionCounts, anim uint16, player, opponent types.PostFrameUpdate) {
	opponentDir := float32(1)
	if player.XPos > opponent.XPos {
		opponentDir = -1
	}
	facingOpponent := player.FacingDirection == opponentDir
	if anim == 0xbf || anim == 0xb7 {
		counts.GroundTechCount.Fail++
	}
	if (anim == 0xc8 && facingOpponent) || (anim == 0xc9 && !facingOpponent) {
		counts.GroundTechCount.In++
	}
	if anim == 0xc7 {
		counts.GroundTechCount.Neutral++
	}
	if (anim == 0xc9 && facingOpponent) || (anim == 0xc8 && !facingOpponent) {
		counts.GroundTechCount.Away++
	}
	if anim == 0xca {
		counts.WallTechCount.Success++
	}
	if anim == 0xf7 {
		counts.WallTechCount.Fail++
	}
}

func handleWavedash(counts *ActionCounts, animations []uint16, positionsY []float32) {
	// Aligns with slippi-js: uses an 8-frame lookback and simple presence checks
	// (no Y-displacement heuristic) to classify special-landings after air-dodge.
	if len(animations) < 2 {
		return
	}
	current := animations[len(animations)-1]
	prev := animations[len(animations)-2]
	isSpecialLanding := current == 0x2b
	isAcceptablePrev := isWavedashInitiationAnimation(prev)
	if !isSpecialLanding || !isAcceptablePrev {
		return
	}
	const lookbackFrames = 8
	start := len(animations) - lookbackFrames
	if start < 0 {
		start = 0
	}
	recent := animations[start:]
	hasAirDodge := false
	hasKneeBend := false
	unique := map[uint16]struct{}{}
	for _, a := range recent {
		unique[a] = struct{}{}
		if a == 0xec {
			hasAirDodge = true
		}
		if a == 0x18 {
			hasKneeBend = true
		}
	}
	if len(unique) == 2 && hasAirDodge {
		return
	}
	if hasAirDodge {
		counts.AirDodgeCount--
	}
	if hasKneeBend {
		counts.WavedashCount++
	} else {
		counts.WavelandCount++
	}
}

func isWavedashInitiationAnimation(animation uint16) bool {
	if animation == 0xec {
		return true
	}
	return animation >= 0x18 && animation <= 0x22
}

func risingEdgeButtonCount(prev, curr uint16) int {
	// Apply ToInt32-style bitwise ~ on physicalButtons, then & curr & 0xfff.
	p := int32(uint16(prev))
	c := int32(uint16(curr))
	changes := (^p) & c & 0x0fff
	return bits.OnesCount32(uint32(changes))
}

func prevAnimationRunLength(animations []uint16) int {
	if len(animations) < 2 {
		return 0
	}
	startIndex := len(animations) - 2
	target := animations[startIndex]
	i := startIndex
	frameCount := 0
	for i >= 0 && animations[i] == target {
		frameCount++
		i--
	}
	return frameCount
}

func joystickRegion(x, y float32) int {
	// Stick values arrive as float32 but comparisons are made in float64.
	// Comparisons against 0.2875 happen in double precision. If we compared in float32 the
	// boundary value 0.2875f32 (== 0.28749999403953552 as a double) would satisfy
	// `x >= 0.2875` in Go but not in JS, causing off-by-one joystick-region transitions.
	xd, yd := float64(x), float64(y)
	switch {
	case xd >= 0.2875 && yd >= 0.2875:
		return 1
	case xd >= 0.2875 && yd <= -0.2875:
		return 2
	case xd <= -0.2875 && yd <= -0.2875:
		return 3
	case xd <= -0.2875 && yd >= 0.2875:
		return 4
	case yd >= 0.2875:
		return 5
	case xd >= 0.2875:
		return 6
	case yd <= -0.2875:
		return 7
	case xd <= -0.2875:
		return 8
	default:
		return 0
	}
}
