package model

// Rarity type constants mirror the game's internal rarity identifier strings.
const (
	RarityType1        = "rarity_1"
	RarityType2        = "rarity_2"
	RarityType3        = "rarity_3"
	RarityType4        = "rarity_4"
	RarityTypeBirthday = "rarity_birthday"
)

// Card represents the base card information from the game's master data.
//
// The fields of Card describe immutable properties that are the same for all
// copies of a given card (e.g. rarity, attribute, and supported unit).
type Card struct {
	// ID uniquely identifies the card in the game.
	ID int `json:"id"`

	// CharacterID links to the character associated with this card.
	CharacterID int `json:"characterId"`

	// CardRarityType indicates the card's rarity (for example "rarity_1",
	// "rarity_2", "rarity_3", "rarity_4", or "rarity_birthday").
	CardRarityType string `json:"cardRarityType"`

	// Attr represents the card's attribute (cool, cute, happy, pure, mysterious).
	Attr string `json:"attr"`

	// SupportUnit indicates which unit the card provides support for, using
	// the internal unit keys from the master data (e.g. "light_sound", "idol").
	SupportUnit string `json:"supportUnit"`

	// Prefix is the card's title or name prefix.
	Prefix string `json:"prefix"`
}

// Card episode part types as they appear in cardEpisodes.json.
const (
	EpisodePartFirst  = "first_part"
	EpisodePartSecond = "second_part"
)

// CardEpisode describes one side story part of a card in the game's master data.
//
// A card has side stories only if cardEpisodes.json contains entries for its ID.
type CardEpisode struct {
	// ID uniquely identifies the episode.
	ID int `json:"id"`

	// Seq is the episode's order within its card.
	Seq int `json:"seq"`

	// CardID links to the card this episode belongs to.
	CardID int `json:"cardId"`

	// CardEpisodePartType is "first_part" or "second_part".
	CardEpisodePartType string `json:"cardEpisodePartType"`
}

// CardEntity extends Card with user-specific properties that represent a card
// instance in the user's inventory.
//
// These fields can be modified over time as the user levels and invests in
// the card (level, master rank, skill level, side stories, painting).
type CardEntity struct {
	// Embed the base Card type.
	Card

	// Level is the card's current level (1–60).
	Level int `json:"level"`

	// MasterRank indicates the number of duplicate cards merged (0–5).
	MasterRank int `json:"masterRank"`

	// SkillLevel represents the card's skill enhancement level (1–4).
	SkillLevel int `json:"skillLevel"`

	// SideStory1 reports whether the first side story has been unlocked.
	SideStory1 bool `json:"sideStory1"`

	// SideStory2 reports whether the second side story has been unlocked.
	SideStory2 bool `json:"sideStory2"`

	// Painting reports whether the card's painting has been unlocked.
	Painting bool `json:"painting"`
}
