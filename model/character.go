package model

// Character represents a Project Sekai character with basic information used
// for display and filtering in the Sekai Inventory Manager.
type Character struct {
	// ID is the unique identifier for the character.
	ID int `json:"id"`

	// FirstName is the character's family name (may be empty for some
	// characters or data sources).
	FirstName string `json:"firstName"`

	// GivenName is the character's given name.
	GivenName string `json:"givenName"`

	// Unit is the internal name of the group or unit the character belongs to
	// (for example "light_sound", "idol", "street", etc.).
	Unit string `json:"unit"`
}
