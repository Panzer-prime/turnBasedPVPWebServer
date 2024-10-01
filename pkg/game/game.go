package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

func AdvanceTurn(gs *GameState) {
	gs.TurnIndex = (gs.TurnIndex + 1) % len(gs.Players)
}

func CalculateDmg(attackDmg int16, defense int16) int16 {
	damageReduction := float64(defense) / (float64(defense) + 100.0)

	finalDamage := float64(attackDmg) * (1.0 - damageReduction)

	return int16(finalDamage)
}

func IsGameOver(gameState *GameState) bool {

	for _, player := range gameState.Players {
		if player.Health <= 0 {
			gameState.IsGameOver = true
			return true
		}
	}

	return false
}

// here everything is calculated from the parspective of the player who has the turn
func HandleAction(action []byte, currentState *GameState) error {
	var newTurn Turn

	if err := json.Unmarshal(action, &newTurn); err != nil {
		return err
	}

	if newTurn.Skip {
		return nil
	}

	currentIndex := currentState.TurnIndex
	player := currentState.Players[currentIndex]

	nextPlayerIndex := (currentIndex + 1) % len(currentState.Players)
	nextPlayer := currentState.Players[nextPlayerIndex]
	//decrease the duration of the buffs e.g D = 5 after funcs it should be 4 meaning 4 more turns with the respective effect
	UpdateEffectDuration(&player.Buffs)
	UpdateEffectDuration(&player.Debuffs)

	//the values of the effect buff or debuff (buff is pozitive e.g +5, debuff negative -5)
	buffs := CalculateEffect(player.Buffs)
	debuffs := CalculateEffect(player.Debuffs)

	var totalDmg, manaCost int16 = 0, 0
	cards, err := getCardsByIDArr(newTurn.CardsId)

	if err != nil {
		return err
	}

	for _, card := range cards {
		switch card.CardType {
		case "attack":

			if debuffs["attack"] != 0 {
				card.Value += debuffs["attack"]
			}
			damageDealt := CalculateDmg(card.Value, nextPlayer.Defense)
			totalDmg += damageDealt

		case "defense":

			if debuffs["defense"] != 0 {
				card.Value += debuffs["defense"]
			}

			if player.Defense+card.Value > 100 {
				player.Defense = 100
			}
			player.Defense += card.Value
		case "heal":
			player.Health += card.Value
			if player.Health > 100 {
				player.Health = 100

			}

		case "buff":
			AddEffect(&player.Buffs, card)
		case "debuff":
			AddEffect(&nextPlayer.Debuffs, card)

		}

		manaCost += card.ManaCost
	}
	if totalDmg > 0 {
		nextPlayer.Health -= totalDmg + buffs["attack"]
		nextPlayer.Defense -= totalDmg
	}

	player.Mana -= manaCost
	return nil
}

func getCardsByIDArr(cardsID []string) ([]Card, error) {
	var foundCards []Card

	loadCards, err := LoadCards()
	if err != nil {
		return nil, err
	}

	mapCards := make(map[string]Card)

	for _, card := range loadCards {
		mapCards[card.ID] = card
	}

	for _, id := range cardsID {

		if card, ok := mapCards[id]; ok {
			foundCards = append(foundCards, card)
		}

	}

	return foundCards, nil
}

func AddEffect(effectList *[]Effect, card Card) {
	newEffect := Effect{
		EffectTupe: card.EffectType,
		Value:      card.EffectValue,
		Duration:   card.Duration,
		Stackable:  card.Stackable,
	}

	for i, effect := range *effectList {
		if effect.EffectTupe == newEffect.EffectTupe {
			if newEffect.Stackable {
				(*effectList)[i].Value += card.EffectValue

			} else {
				(*effectList)[i] = newEffect
			}
			return
		}
	}
	*effectList = append(*effectList, newEffect)

}

func UpdateEffectDuration(effectList *[]Effect) {
	var newEffectList []Effect

	for _, effect := range *effectList {
		effect.Duration--
		if effect.Duration > 0 {
			newEffectList = append(newEffectList, effect)
		}
	}

	*effectList = newEffectList
}

func CalculateEffect(effectList []Effect) map[string]int16 {
	effectStats := make(map[string]int16)

	for _, effect := range effectList {

		effectStats[effect.EffectTupe] += effect.Value
	}

	return effectStats
}

func GameInit(gamestate *GameState) (*GameState, error) {
	for _, player := range gamestate.Players {
		cards, err := PickCards(20) //it wont be hard coded here i will change it when i will add more gameplay modes
		if err != nil {
			return nil, err
		}

		player.Cards = cards
		if err != nil {
			return nil, err
		}
	}
	return gamestate, nil
}

func ShuffleCards(cards []Card) {
	for i := len(cards) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

// Function to pick cards
func PickCards(amountOfCards int16) ([]Card, error) {

	loadCards, err := LoadCards()
	if err != nil {
		return nil, err
	}

	// Shuffle the deck initially
	ShuffleCards(loadCards)

	cardCounter := make(map[string]int16)

	maxCards := map[string]int16{
		"attack":  int16(float64(amountOfCards) * 0.5),
		"defense": int16(float64(amountOfCards) * 0.2),
		"heal":    int16(float64(amountOfCards) * 0.1),
		"buff":    int16(float64(amountOfCards) * 0.1),
		"debuff":  int16(float64(amountOfCards) * 0.1),
		"special": int16(float64(amountOfCards) * 0.1),
	}

	var pickedCards []Card
	const maxAttempts = 100
	attempts := 0

	for int16(len(pickedCards)) < amountOfCards && attempts < maxAttempts {
		index := rand.Intn(len(loadCards))
		card := loadCards[index]

		// Check duplicate card and type limits
		if cardCounter[card.ID] >= 2 || cardCounter[card.CardType] >= maxCards[card.CardType] {
			fmt.Println("Card skipped due to limits:", card.ID, card.CardType)
			attempts++
			continue
		}

		cardCounter[card.ID]++
		cardCounter[card.CardType]++
		pickedCards = append(pickedCards, card)

		attempts = 0
	}

	if attempts >= maxAttempts {
		return nil, errors.New("failed to pick cards within max attempts, check card pool or constraints")
	}

	// Print picked cards for debugging
	for _, card := range pickedCards {
		fmt.Println("Picked:", card.CardType, card.ID)
	}

	return pickedCards, nil
}
