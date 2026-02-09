package tools

import (
	"fmt"
	"sekai-inventory/model"
	"strings"
	"time"

	"github.com/fatih/color"
)

// CreateCharacterMap converts a slice of Characters into a map indexed by CharacterID
// for efficient character lookup operations. This is particularly useful when processing
// card data that references characters by their IDs.
func CreateCharacterMap(characters []model.Character) map[int]model.Character {
	characterMap := make(map[int]model.Character)
	for _, char := range characters {
		characterMap[char.ID] = char
	}
	return characterMap
}

// PrintSuccessMessage displays a success message in green color to provide
// clear visual feedback for successful operations.
func PrintSuccessMessage(message string) {
	color.Green(message)
}

// PrintErrorMessage displays an error message in red color to highlight
// critical issues or failures that require user attention.
func PrintErrorMessage(message string) {
	color.Red(message)
}

// PrintWarningMessage displays a warning message in yellow color to indicate
// potential issues or important information that doesn't prevent operation.
func PrintWarningMessage(message string) {
	color.Yellow(message)
}

// UpdateTimeSet updates the UpdatedAt timestamp of the inventory file.
// This function should be called whenever the inventory contents are modified
// to maintain accurate modification tracking. It handles both loading and saving
// the inventory file, updating only the timestamp.
//
// Returns an error if there are issues accessing or modifying the inventory file.
func UpdateTimeSet() error {
	inventory, err := LoadInventory()
	if err != nil {
		return fmt.Errorf("error loading inventory: %v", err)
	}

	inventory.UpdatedAt = time.Now()

	if err = SaveInventory(inventory); err != nil {
		return fmt.Errorf("error saving inventory: %v", err)
	}

	return nil
}

// ParseFilters converts command-line arguments into a structured filter map.
// It expects arguments in pairs of --field value format (e.g., --character "Miku").
//
// Parameters:
//   - args: Command-line arguments to parse, excluding the command itself
//
// Returns:
//   - A map of field names to their values, or nil if the arguments are invalid
//     (e.g., odd number of arguments or missing values)
func ParseFilters(args []string) map[string]string {
	if len(args)%2 != 0 {
		// Invalid number of arguments (filters must be in key-value pairs)
		return nil
	}

	filters := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		field := strings.TrimPrefix(args[i], "--") // Remove "--" prefix
		if i+1 >= len(args) {
			// Missing value for the field
			return nil
		}
		value := args[i+1]
		filters[field] = value
	}

	return filters
}
