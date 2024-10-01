package game_test

import (
	"arena-game/pkg/game"

	"testing"
)

func TestAddEffect(t *testing.T) {
	// Test 1: Add a new effect
	playerEffects := []game.Effect{}
	card := game.Card{
		CardType:    "attack",
		EffectValue: 10,
		Duration:    3,
		Stackable:   true,
	}

	game.AddEffect(&playerEffects, card)
	if len(playerEffects) != 1 {
		t.Errorf("Expected 1 effect, got %d", len(playerEffects))
	}

	if playerEffects[0].Value != 10 {
		t.Errorf("Expected effect value to be 10, got %d", playerEffects[0].Value)
	}

	// Test 2: Stackable effect
	card2 := game.Card{
		CardType:    "attack",
		EffectValue: 5,
		Duration:    2,
		Stackable:   true,
	}

	game.AddEffect(&playerEffects, card2)

	if playerEffects[0].Value != 15 {
		t.Errorf("Expected stacked value to be 15, got %d", playerEffects[0].Value)
	}

	// Test 3: Non-stackable effect replaces the existing one
	card3 := game.Card{
		CardType:    "attack",
		EffectValue: 20,
		Duration:    2,
		Stackable:   false,
	}

	game.AddEffect(&playerEffects, card3)

	if playerEffects[0].Value != 20 {
		t.Errorf("Expected non-stackable effect value to replace, got %d", playerEffects[0].Value)
	}
}

func TestCalculateEffect(t *testing.T) {
	// Test case 1: Single effect
	effectList := []game.Effect{
		{EffectTupe: "attack", Value: 10},
	}

	effectStats := game.CalculateEffect(effectList)

	if len(effectStats) != 1 {
		t.Errorf("Expected 1 effect type, got %d", len(effectStats))
	}

	if effectStats["attack"] != 10 {
		t.Errorf("Expected attack value to be 10, got %d", effectStats["attack"])
	}

	// Test case 2: Multiple effects of different types
	effectList = []game.Effect{
		{EffectTupe: "attack", Value: 10},
		{EffectTupe: "defense", Value: 5},
		{EffectTupe: "heal", Value: 15},
	}

	effectStats = game.CalculateEffect(effectList)

	if len(effectStats) != 3 {
		t.Errorf("Expected 3 effect types, got %d", len(effectStats))
	}

	if effectStats["attack"] != 10 {
		t.Errorf("Expected attack value to be 10, got %d", effectStats["attack"])
	}

	if effectStats["defense"] != 5 {
		t.Errorf("Expected defense value to be 5, got %d", effectStats["defense"])
	}

	if effectStats["heal"] != 15 {
		t.Errorf("Expected heal value to be 15, got %d", effectStats["heal"])
	}

	// Test case 3: Multiple effects of the same type
	effectList = []game.Effect{
		{EffectTupe: "attack", Value: 10},
		{EffectTupe: "attack", Value: 5},
		{EffectTupe: "attack", Value: 15},
	}

	effectStats = game.CalculateEffect(effectList)

	if len(effectStats) != 1 {
		t.Errorf("Expected 1 effect type, got %d", len(effectStats))
	}

	if effectStats["attack"] != 30 {
		t.Errorf("Expected attack value to be 30, got %d", effectStats["attack"])
	}

	// Test case 4: No effects
	effectList = []game.Effect{}

	effectStats = game.CalculateEffect(effectList)

	if len(effectStats) != 0 {
		t.Errorf("Expected no effect types, got %d", len(effectStats))
	}
}
