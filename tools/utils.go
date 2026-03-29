package tools

import (
	"fmt"
	"sekai-inventory/model"
	"strings"
	"time"

	"github.com/fatih/color"
)

// CreateCharacterMap converts a slice of Character values into a map indexed
// by CharacterID. It is useful for efficient character lookup when processing
// card data.
func CreateCharacterMap(characters []model.Character) map[int]model.Character {
	characterMap := make(map[int]model.Character)
	for _, char := range characters {
		characterMap[char.ID] = char
	}
	return characterMap
}

// PrintSuccessMessage displays a success message in green to provide clear
// visual feedback for successful operations.
func PrintSuccessMessage(message string) {
	color.Green(message)
}

// PrintErrorMessage displays an error message in red to highlight critical
// issues or failures that require user attention.
func PrintErrorMessage(message string) {
	color.Red(message)
}

// PrintWarningMessage displays a warning message in yellow to indicate
// potential issues or important information that does not prevent operation.
func PrintWarningMessage(message string) {
	color.Yellow(message)
}

// UpdateTimeSet updates the UpdatedAt timestamp of the inventory file to
// the current time.
//
// It should be called whenever the inventory contents are modified to maintain
// accurate modification tracking. UpdateTimeSet returns an error if there are
// problems accessing or saving the inventory file.
func UpdateTimeSet() error {
	inventory, err := LoadInventory()
	if err != nil {
		return fmt.Errorf("error loading inventory: %v", err)
	}

	inventory.UpdatedAt = time.Now()

	if err := SaveInventory(inventory); err != nil {
		return fmt.Errorf("error saving inventory: %v", err)
	}

	return nil
}

// ParseFilters converts command-line arguments into a map of filter names to
// values.
//
// It expects arguments in pairs of "--field value", for example:
//
//	--character Miku --rarity 4
//
// ParseFilters returns nil if the arguments are invalid (e.g. odd number of
// tokens or missing value for a field).
func ParseFilters(args []string) map[string]string {
	if len(args)%2 != 0 {
		// Invalid number of arguments (filters must be in key-value pairs).
		return nil
	}

	filters := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		field := strings.TrimPrefix(args[i], "--") // Remove "--" prefix.
		if i+1 >= len(args) {
			// Missing value for the field.
			return nil
		}
		value := args[i+1]
		filters[field] = value
	}

	return filters
}
