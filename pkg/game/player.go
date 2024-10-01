package game

func CreatePlayer(id string, name string) *Player {
	return &Player{
		ID:      id,
		Name:    name,
		Health:  100,
		Defense: 100,
		Mana:    100,
		Cards:   []Card{},
	}
}
