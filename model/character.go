package model

// Character represents a Project Sekai character with their basic information
type Character struct {
	// ID is the unique identifier for the character
	ID int `json:"id"`
	// FirstName is the character's family name
	FirstName string `json:"firstName"`
	// GivenName is the character's given name
	GivenName string `json:"givenName"`
	// Unit is the group/unit the character belongs to (e.g., "L/N", "MMJ", "VBS", etc.)
	Unit string `json:"unit"`
}
