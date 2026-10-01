package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sekai-inventory/model"
)

// File paths for persistent storage.
const (
	// InventoryFile is the path to the user's card inventory.
	InventoryFile = "res/inventory.json"

	// CardsFile is the path to the game's card database.
	CardsFile = "res/cards.json"

	// CharactersFile is the path to the game's character database.
	CharactersFile = "res/gameCharacters.json"

	// CardEpisodesFile is the path to the game's card episode (side story) database.
	CardEpisodesFile = "res/cardEpisodes.json"
)

// EnsureResDirectory creates the "res" directory if it does not exist.
//
// It returns an error if directory creation fails or if there are permission
// issues when checking or creating the directory.
func EnsureResDirectory() error {
	if _, err := os.Stat("res"); os.IsNotExist(err) {
		if err := os.Mkdir("res", 0o755); err != nil {
			return fmt.Errorf("error creating res directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking res directory: %w", err)
	}
	return nil
}

// LoadInventory reads and parses the user's card inventory from InventoryFile.
//
// If the inventory file does not exist, LoadInventory returns an empty
// Inventory value instead of an error. This supports first-time usage where
// no inventory exists yet.
//
// Additionally, LoadInventory checks whether the inventory uses the latest
// schema (i.e. includes the "painting" field on cards). If the field is
// missing and there are cards present, it prints a message instructing the
// user to run the "convert" command and returns an error.
func LoadInventory() (*model.Inventory, error) {
	// Check if the inventory file exists.
	if _, err := os.Stat(InventoryFile); os.IsNotExist(err) {
		// If the file does not exist, return an empty inventory.
		return &model.Inventory{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check inventory file: %w", err)
	}

	// Read the inventory file.
	data, err := os.ReadFile(InventoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %w", err)
	}

	// Schema check: does the file contain the "painting" field?
	// If not, and if there are cards, instruct the user to run "convert".
	if !bytes.Contains(data, []byte(`"painting"`)) {
		// Probe only the cards array to see if there is anything to migrate.
		var probe struct {
			Cards []json.RawMessage `json:"cards"`
		}
		if err := json.Unmarshal(data, &probe); err != nil {
			return nil, fmt.Errorf("failed to parse inventory file: %w", err)
		}

		if len(probe.Cards) > 0 {
			PrintErrorMessage("Your inventory file uses an outdated schema (missing 'painting' field).")
			PrintWarningMessage("Please run 'sekai-inventory convert' once to migrate your inventory.")
			return nil, fmt.Errorf("inventory schema outdated: missing 'painting'")
		}
		// If there are no cards, treat it as effectively compatible and continue.
	}

	// Parse the JSON data into an Inventory struct.
	var inv model.Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, fmt.Errorf("failed to parse inventory file: %w", err)
	}
	return &inv, nil
}

// SaveInventory serializes inv as JSON and writes it to InventoryFile.
//
// The JSON output is pretty-printed for better human readability. SaveInventory
// returns an error if file creation or JSON encoding fails.
func SaveInventory(inv *model.Inventory) error {
	file, err := os.OpenFile(InventoryFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("failed to create inventory file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON.
	if err := encoder.Encode(inv); err != nil {
		return fmt.Errorf("failed to encode inventory: %w", err)
	}

	return nil
}

// LoadCards reads and parses the game's card database from CardsFile.
//
// The cards.json file contains the master data for all available cards in the
// game. LoadCards returns an error if the file is missing or contains invalid
// data.
func LoadCards() ([]model.Card, error) {
	// Check if the cards file exists.
	if _, err := os.Stat(CardsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("cards.json not found. Please run 'sekai-inventory update' to fetch the latest data")
	} else if err != nil {
		return nil, fmt.Errorf("failed to check cards.json file: %w", err)
	}

	// Read the cards file.
	data, err := os.ReadFile(CardsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cards.json: %w", err)
	}

	// Parse the JSON data into a slice of Card structs.
	var cards []model.Card
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("failed to parse cards.json: %w", err)
	}

	return cards, nil
}

// LoadCharacters reads and parses the game's character database from
// CharactersFile.
//
// The gameCharacters.json file is required for displaying character names and
// unit affiliations. LoadCharacters returns an error if the file is missing or
// contains invalid data.
func LoadCharacters() ([]model.Character, error) {
	// Check if the characters file exists.
	if _, err := os.Stat(CharactersFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("gameCharacters.json not found. Please run 'sekai-inventory update' to fetch the latest data")
	} else if err != nil {
		return nil, fmt.Errorf("failed to check gameCharacters.json file: %w", err)
	}

	// Read the characters file.
	data, err := os.ReadFile(CharactersFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read gameCharacters.json: %w", err)
	}

	// Parse the JSON data into a slice of Character structs.
	var characters []model.Character
	if err := json.Unmarshal(data, &characters); err != nil {
		return nil, fmt.Errorf("failed to parse gameCharacters.json: %w", err)
	}
	return characters, nil
}

// LoadCardEpisodes reads and parses the game's card episode database from
// CardEpisodesFile.
//
// The cardEpisodes.json file lists the side story parts of each card. It is
// used to tell which cards have side stories at all. LoadCardEpisodes returns
// an error if the file is missing or contains invalid data.
func LoadCardEpisodes() ([]model.CardEpisode, error) {
	// Check if the card episodes file exists.
	if _, err := os.Stat(CardEpisodesFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("cardEpisodes.json not found. Please run 'sekai-inventory update' to fetch the latest data")
	} else if err != nil {
		return nil, fmt.Errorf("failed to check cardEpisodes.json file: %w", err)
	}

	// Read the card episodes file.
	data, err := os.ReadFile(CardEpisodesFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cardEpisodes.json: %w", err)
	}

	// Parse the JSON data into a slice of CardEpisode structs.
	var episodes []model.CardEpisode
	if err := json.Unmarshal(data, &episodes); err != nil {
		return nil, fmt.Errorf("failed to parse cardEpisodes.json: %w", err)
	}
	return episodes, nil
}

// CreateSideStorySet returns the set of card IDs that have side stories, that
// is, every card ID that appears in the given episodes. In the current master
// data a card has either both side story parts or none, so one flag per card
// is enough.
func CreateSideStorySet(episodes []model.CardEpisode) map[int]bool {
	set := make(map[int]bool, len(episodes)/2+1)
	for _, e := range episodes {
		set[e.CardID] = true
	}
	return set
}
