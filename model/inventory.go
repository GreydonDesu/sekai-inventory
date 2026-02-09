package model

import "time"

// Inventory represents the user's collection of cards and tracks modification times.
// It serves as the root structure for persistence and maintains timestamps for
// tracking when the inventory was created and last modified.
type Inventory struct {
	// Cards holds all the cards currently in the user's inventory
	Cards []CardEntity `json:"cards"`
	// CreatedAt stores the timestamp when the inventory was first created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt stores the timestamp of the last modification to the inventory
	UpdatedAt time.Time `json:"updated_at"`
}
