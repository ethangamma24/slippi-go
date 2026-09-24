package slippi

import (
	"context"
	"testing"
)

func TestFilterConversionsByOpeningMove(t *testing.T) {
	root := mustProjectRoot(t)
	fixtures := mustFixtures(t, root+"/testdata/slp")

	ctx := context.Background()

	noFilter := ComputeOptions{}
	withFilter := ComputeOptions{ExcludedMoves: []uint8{NeutralBMoveID}}

	totalUnfiltered := 0
	totalFiltered := 0
	difference := 0

	for _, fixture := range fixtures {
		game := NewGame(fixture)

		unfiltered, err := game.StatsTypedWithOptions(ctx, noFilter)
		if err != nil {
			t.Fatalf("unfiltered stats %s: %v", fixture, err)
		}

		filtered, err := game.StatsTypedWithOptions(ctx, withFilter)
		if err != nil {
			t.Fatalf("filtered stats %s: %v", fixture, err)
		}

		totalUnfiltered += len(unfiltered.Conversions)
		totalFiltered += len(filtered.Conversions)
		difference += len(unfiltered.Conversions) - len(filtered.Conversions)

		for _, c := range filtered.Conversions {
			if len(c.Moves) > 0 && c.Moves[0].MoveID == NeutralBMoveID {
				t.Errorf("%s: conversion still has excluded opening move %d", fixture, NeutralBMoveID)
			}
		}
	}

	if totalUnfiltered == 0 {
		t.Skip("no conversions found in fixtures, cannot verify filter")
	}

	t.Logf("unfiltered=%d filtered=%d removed=%d", totalUnfiltered, totalFiltered, difference)

	if difference == 0 {
		t.Log("warning: no conversions were filtered — fixtures may not contain neutral-B-opening conversions")
	}
}

func TestFilterConversionsEmptyOptions(t *testing.T) {
	root := mustProjectRoot(t)
	fixtures := mustFixtures(t, root+"/testdata/slp")

	ctx := context.Background()

	for _, fixture := range fixtures {
		game := NewGame(fixture)

		defaultOpts, err := game.StatsTyped(ctx)
		if err != nil {
			t.Fatalf("default stats %s: %v", fixture, err)
		}

		emptyOpts, err := game.StatsTypedWithOptions(ctx, ComputeOptions{})
		if err != nil {
			t.Fatalf("empty options stats %s: %v", fixture, err)
		}

		if len(defaultOpts.Conversions) != len(emptyOpts.Conversions) {
			t.Errorf("%s: default opts gave %d conversions but empty opts gave %d", fixture, len(defaultOpts.Conversions), len(emptyOpts.Conversions))
		}
	}
}

func TestFilterConversionsNilExcludedMoves(t *testing.T) {
	root := mustProjectRoot(t)
	fixtures := mustFixtures(t, root+"/testdata/slp")

	ctx := context.Background()

	for _, fixture := range fixtures {
		game := NewGame(fixture)

		defaultOpts, err := game.StatsTyped(ctx)
		if err != nil {
			t.Fatalf("default stats %s: %v", fixture, err)
		}

		nilFilter, err := game.StatsTypedWithOptions(ctx, ComputeOptions{ExcludedMoves: nil})
		if err != nil {
			t.Fatalf("nil filter stats %s: %v", fixture, err)
		}

		if len(defaultOpts.Conversions) != len(nilFilter.Conversions) {
			t.Errorf("%s: default opts gave %d conversions but nil filter gave %d", fixture, len(defaultOpts.Conversions), len(nilFilter.Conversions))
		}
	}
}
