package game

import (
	"encoding/json"
	"io"
	"os"
	
)



func LoadCards() ([]Card, error) {
	jsonFile, err := os.Open("../assets/cards.json")
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var cards []Card

	if err := json.Unmarshal(byteValue, &cards); err != nil {
		return nil, err
	}
	jsonFile.Close()
	return cards, nil
}
