package game_test

import (
	"arena-game/pkg/game"
	"testing"
)

func TestPickCards(t *testing.T) {
	// Test 1: Basic functionality
	pickedCards, err := game.PickCards(10)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(pickedCards) != 10 {
		t.Fatalf("Expected 10 cards, got %d", len(pickedCards))
	}

	// Test 2: Card type ratios
	cardCounter := make(map[string]int16)
	for _, card := range pickedCards {
		cardCounter[card.CardType]++
	}
	expectedRatios := map[string]float64{
		"attack":  0.5,
		"defense": 0.2,
		"heal":    0.1,
		"buff":    0.1,
		"debuff":  0.1,
	}

	for cardType, expectedRatio := range expectedRatios {
		expectedCount := int16(float64(10) * expectedRatio)
		if cardCounter[cardType] > expectedCount {
			t.Errorf("Too many %s cards: got %d, want at most %d", cardType, cardCounter[cardType], expectedCount)
		}
	}

	// Test 3: Duplicate handling
	duplicateCounter := make(map[string]int16)
	for _, card := range pickedCards {
		duplicateCounter[card.ID]++
		if duplicateCounter[card.ID] > 2 {
			t.Errorf("Duplicate card detected: %s", card.ID)
		}
	}

}
