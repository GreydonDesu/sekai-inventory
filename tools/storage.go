package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"sekai-inventory/model"
)

// File paths for persistent storage
const (
	// InventoryFile is the path to the user's card inventory
	InventoryFile = "res/inventory.json"
	// CardsFile is the path to the game's card database
	CardsFile = "res/cards.json"
	// CharactersFile is the path to the game's character database
	CharactersFile = "res/gameCharacters.json"
)

// EnsureResDirectory creates the resources directory if it doesn't exist.
// This directory is required for storing game data and user inventory.
//
// Returns an error if directory creation fails or if there are permission issues.
func EnsureResDirectory() error {
	if _, err := os.Stat("res"); os.IsNotExist(err) {
		err := os.Mkdir("res", 0755)
		if err != nil {
			return fmt.Errorf("error creating res directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("error checking res directory: %w", err)
	}
	return nil
}

// LoadInventory reads and parses the user's card inventory from the JSON file.
// If the inventory file doesn't exist, it returns an empty inventory instead of an error.
// This behavior supports first-time usage where no inventory exists yet.
//
// Returns:
//   - A pointer to the Inventory structure containing the user's cards
//   - An error if file reading or JSON parsing fails
func LoadInventory() (*model.Inventory, error) {
	// Check if the inventory file exists
	if _, err := os.Stat(InventoryFile); os.IsNotExist(err) {
		// If the file doesn't exist, return an empty inventory
		return &model.Inventory{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to check inventory file: %w", err)
	}

	// Read the inventory file
	data, err := os.ReadFile(InventoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %w", err)
	}

	// Parse the JSON data into an Inventory struct
	var inv model.Inventory
	err = json.Unmarshal(data, &inv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inventory file: %w", err)
	}
	return &inv, nil
}

// SaveInventory writes the current inventory state to the JSON file.
// The JSON output is pretty-printed for better human readability.
//
// Parameters:
//   - inv: A pointer to the Inventory structure to be saved
//
// Returns an error if file creation or JSON encoding fails
func SaveInventory(inv *model.Inventory) error {
	file, err := os.Create(InventoryFile)
	if err != nil {
		return fmt.Errorf("failed to create inventory file: %w", err)
	}
	defer file.Close() // Ensure the file is closed

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	err = encoder.Encode(inv)
	if err != nil {
		return fmt.Errorf("failed to encode inventory: %w", err)
	}

	return nil
}

// LoadCards reads and parses the game's card database from cards.json.
// This file contains the master data for all available cards in the game.
// The file must exist for the application to function properly.
//
// Returns:
//   - A slice of Card structures containing all available cards
//   - An error if the file is missing or contains invalid data
func LoadCards() ([]model.Card, error) {
	// Check if the cards file exists
	if _, err := os.Stat(CardsFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("cards.json not found. Please run 'sekai-inventory update' to fetch the latest data")
	} else if err != nil {
		return nil, fmt.Errorf("failed to check cards.json file: %w", err)
	}

	// Read the cards file
	data, err := os.ReadFile(CardsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cards.json: %w", err)
	}

	// Parse the JSON data into a slice of Card structs
	var cards []model.Card
	err = json.Unmarshal(data, &cards)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cards.json: %w", err)
	}

	return cards, nil
}

// LoadCharacters reads and parses the game's character database.
// This data is required for displaying character names and unit affiliations.
// The file must be updated using the 'update' command if it doesn't exist.
//
// Returns:
//   - A slice of Character structures containing all game characters
//   - An error if the file is missing or contains invalid data
func LoadCharacters() ([]model.Character, error) {
	// Check if the characters file exists
	if _, err := os.Stat(CharactersFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("gameCharacters.json not found. Please run 'sekai-inventory update' to fetch the latest data")
	} else if err != nil {
		return nil, fmt.Errorf("failed to check gameCharacters.json file: %w", err)
	}

	// Read the characters file
	data, err := os.ReadFile(CharactersFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read gameCharacters.json: %w", err)
	}

	// Parse the JSON data into a slice of Character structs
	var characters []model.Character
	err = json.Unmarshal(data, &characters)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gameCharacters.json: %w", err)
	}
	return characters, nil
}
