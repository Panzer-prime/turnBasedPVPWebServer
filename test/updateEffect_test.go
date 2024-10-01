package game_test

import (
	"arena-game/pkg/game"
	"testing"
)

func TestUpdateEffectDuration(t *testing.T) {
	// Test case 1: Effect with non-zero duration should be retained
	effectList := []game.Effect{
		{EffectTupe: "attack", Duration: 2, Value: 10},
		{EffectTupe: "defense", Duration: 1, Value: 5},
	}

	game.UpdateEffectDuration(&effectList)

	if len(effectList) != 1 {
		t.Errorf("Expected 1 effect to remain, got %d", len(effectList))
	}

	if effectList[0].EffectTupe != "attack" || effectList[0].Duration != 1 {
		t.Errorf("Expected 'attack' effect with duration 1, got %v with duration %d", effectList[0].EffectTupe, effectList[0].Duration)
	}

	// Test case 2: Effect with zero duration should be removed
	effectList = []game.Effect{
		{EffectTupe: "heal", Duration: 1, Value: 15},
		{EffectTupe: "buff", Duration: 1, Value: 8},
	}

	game.UpdateEffectDuration(&effectList)

	if len(effectList) != 0 {
		t.Errorf("Expected no effects to remain, but got %d effects", len(effectList))
	}

	// Test case 3: Multiple effects with different durations
	effectList = []game.Effect{
		{EffectTupe: "buff", Duration: 3, Value: 10},
		{EffectTupe: "debuff", Duration: 1, Value: -5},
		{EffectTupe: "attack", Duration: 2, Value: 20},
	}

	game.UpdateEffectDuration(&effectList)

	if len(effectList) != 2 {
		t.Errorf("Expected 2 effects to remain, but got %d effects", len(effectList))
	}

	if effectList[0].EffectTupe != "buff" || effectList[0].Duration != 2 {
		t.Errorf("Expected 'buff' effect with duration 2, got %v with duration %d", effectList[0].EffectTupe, effectList[0].Duration)
	}

	if effectList[1].EffectTupe != "attack" || effectList[1].Duration != 1 {
		t.Errorf("Expected 'attack' effect with duration 1, got %v with duration %d", effectList[1].EffectTupe, effectList[1].Duration)
	}
}
