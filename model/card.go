package model

// Card represents the base card information from the game's master data.
// This structure contains the immutable properties of a card that are consistent
// across all instances of the same card.
type Card struct {
	// ID uniquely identifies the card in the game
	ID int `json:"id"`
	// CharacterID links to the character associated with this card
	CharacterID int `json:"characterId"`
	// CardRarityType indicates the card's rarity (1*, 2*, 3*, 4*, or Birthday)
	CardRarityType string `json:"cardRarityType"`
	// Attr represents the card's attribute (Cool, Cute, Happy, Mysterious, Pure)
	Attr string `json:"attr"`
	// SupportUnit indicates which unit the card provides support for
	SupportUnit string `json:"supportUnit"`
	// Prefix is the card's title or name prefix
	Prefix string `json:"prefix"`
}

// CardEntity extends the base Card type with user-specific properties.
// It represents a card instance in the user's inventory and includes
// properties that can be modified by the user, such as level and skill levels.
type CardEntity struct {
	Card // Embed the base Card type for inheritance
	// Level represents the card's current level (1-60)
	Level int `json:"level"`
	// MasterRank indicates the number of duplicate cards merged (0-5)
	MasterRank int `json:"masterRank"`
	// SkillLevel represents the card's skill enhancement level (1-4)
	SkillLevel int `json:"skillLevel"`
	// SideStory1 indicates whether the first side story has been unlocked
	SideStory1 bool `json:"sideStory1"`
	// SideStory2 indicates whether the second side story has been unlocked
	SideStory2 bool `json:"sideStory2"`
}
