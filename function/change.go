package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"strconv"
)

// Change modifies specific attributes of a card in the user's inventory.
// This function supports updating multiple fields in a single operation.
//
// Supported fields and their valid values:
//   - level: Integer between 1 and 60
//   - skillLevel: Integer between 1 and 5
//   - masterRank: Integer between 0 and 5
//   - sideStory1: Boolean (true/false)
//   - sideStory2: Boolean (true/false)
//
// Parameters:
//   - cardID: The ID of the card to modify
//   - updates: Map of field names to their new values (as strings)
//
// Returns:
//   - nil if all updates were successful
//   - error if:
//     * Card not found in inventory
//     * Invalid field name
//     * Invalid value for a field
//     * File operations fail
func Change(cardID int, updates map[string]string) error {
	// Load the inventory
	inventory, err := tools.LoadInventory()
	if err != nil {
		return fmt.Errorf("error loading inventory: %v", err)
	}

	// Find the card in the inventory
	var card *model.CardEntity
	for i := range inventory.Cards {
		if inventory.Cards[i].ID == cardID {
			card = &inventory.Cards[i]
			break
		}
	}

	if card == nil {
		return fmt.Errorf("no card found with ID %d", cardID)
	}

	// Track changes
	changes := make(map[string]string)

	// Apply updates
	for field, value := range updates {
		switch field {
		case "level":
			level, err := strconv.Atoi(value)
			if err != nil || level < 1 || level > 60 {
				return fmt.Errorf("invalid value for 'level': %s. Must be an integer between 1 and 60", value)
			}
			card.Level = level
			changes["level"] = value
		case "skillLevel":
			skillLevel, err := strconv.Atoi(value)
			if err != nil || skillLevel < 1 || skillLevel > 5 {
				return fmt.Errorf("invalid value for 'skillLevel': %s. Must be an integer between 1 and 5", value)
			}
			card.SkillLevel = skillLevel
			changes["skillLevel"] = value
		case "masterRank":
			masterRank, err := strconv.Atoi(value)
			if err != nil || masterRank < 0 || masterRank > 5 {
				return fmt.Errorf("invalid value for 'masterRank': %s. Must be an integer between 0 and 5", value)
			}
			card.MasterRank = masterRank
			changes["masterRank"] = value
		case "sideStory1":
			sideStory1, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for 'sideStory1': %s. Must be 'true' or 'false'", value)
			}
			card.SideStory1 = sideStory1
			changes["sideStory1"] = value
		case "sideStory2":
			sideStory2, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for 'sideStory2': %s. Must be 'true' or 'false'", value)
			}
			card.SideStory2 = sideStory2
			changes["sideStory2"] = value
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}

	// Save the updated inventory
	err = tools.SaveInventory(inventory)
	if err != nil {
		return fmt.Errorf("error saving inventory: %v", err)
	}

	// Report changes
	if len(changes) > 0 {
		fmt.Println("Changes made:")
		for field, value := range changes {
			fmt.Printf("  - Field '%s' changed to '%s'.\n", field, value)
		}
		tools.UpdateTimeSet()
	}

	return nil
}
