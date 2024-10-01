package game

type GameState struct {
	Players    []*Player `json:"players"`
	TurnIndex  int       `json:"turnIndex"`
	IsGameOver bool      `json:"isGameOver"`
	Turn       Turn
}

type Turn struct {
	Skip    bool     `json:"skip"` // Type of action  "skip"
	CardsId []string `json:"cardsID"`
	Target  string   `json:"target"` // Target of the effect ("self", "opponent")
}

type Player struct {
	ID      string   `json:"id"`
	Health  int16    `json:"health"`
	Defense int16    `json:"defense"`
	Mana    int16    `json:"mana"`
	Name    string   `json:"name"`
	Cards   []Card   `json:"items"`
	Buffs   []Effect `json:"buffs"`
	Debuffs []Effect `json:"deBuffs"`
}

type Effect struct {
	EffectTupe string `json:"effectType"`   // the type of buff debuff are wa talkin in cosideration attack defense
	Value      int16  `json:"effectValue"` //the value in int
	Duration   int16  `json:"duration"`     //and the turation per turns
	Stackable  bool   `json:"stackable"`
}

type Card struct {
	ID          string `json:"id"`          // Unique ID for the card
	Name        string `json:"name"`        // Card name
	CardType    string `json:"cardType"`    // "attack", "defense", "heal", "special"
	Value       int16  `json:"value"`       // Damage, defense, or healing amount
	ManaCost    int16  `json:"manaCost"`    // Mana cost
	EffectType  string `json:"effect"`      // Status effect, buff, debuff
	EffectValue int16  `json:"effectValue"` //for buffs they will be posite and negative for debuffs e.g +5 for buff -23 debuff
	Duration    int16  `json:"duration"`    // here the duration is sized in turns more used for buffs /debuffs or long term cards
	Stackable   bool   `json:"stackable"`   //if the card is stackable their values will be added up
	Rarity      string `json:"rarity"`   //for the future
	Description string `json:"description"`
}


