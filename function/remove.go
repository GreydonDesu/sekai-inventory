package function

import (
	"fmt"
	"sekai-inventory/tools"
)

// Remove deletes one or more cards from the user's inventory.
// This function supports batch removal of multiple cards in a single operation.
//
// The function:
//  1. Loads the current inventory
//  2. Attempts to remove each specified card
//  3. Updates the inventory file
//  4. Updates the modification timestamp
//  5. Provides feedback on successful removals and not-found cards
//
// Parameters:
//   - cardIDs: One or more card IDs to remove from the inventory
//
// Success/Error reporting:
//   - Successfully removed cards are listed in a success message
//   - Cards not found in the inventory are listed in a warning message
//   - File operation errors are reported as error messages
func Remove(cardIDs ...int) {
	// Load the inventory
	inventory, err := tools.LoadInventory()
	if err != nil {
		message := fmt.Sprintf("Error loading inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Track cards that were successfully removed and those that were not found
	removedCards := []int{}
	notFoundCards := []int{}

	// Iterate over the provided cardIDs
	for _, cardID := range cardIDs {
		// Check if the card exists in the inventory
		cardIndex := -1
		for i, card := range inventory.Cards {
			if card.ID == cardID {
				cardIndex = i
				break
			}
		}

		if cardIndex == -1 {
			// Card not found
			notFoundCards = append(notFoundCards, cardID)
		} else {
			// Remove the card from the inventory
			inventory.Cards = append(inventory.Cards[:cardIndex], inventory.Cards[cardIndex+1:]...)
			removedCards = append(removedCards, cardID)
		}
	}

	// Save the updated inventory
	err = tools.SaveInventory(inventory)
	if err != nil {
		message := fmt.Sprintf("Error saving inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Print success message for removed cards
	if len(removedCards) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Removed cards with IDs: %v", removedCards))
		tools.UpdateTimeSet()
	}

	// Print warning message for cards that were not found
	if len(notFoundCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Cards with IDs not found: %v", notFoundCards))
	}
}
