package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
)

// Remove deletes one or more cards from the user's inventory.
//
// It supports batch removal of multiple cards in a single operation and prints
// a detailed summary of removed cards and IDs that were not found. Remove
// performs the following steps:
//
//  1. Load the current inventory.
//  2. Attempt to remove each specified card ID.
//  3. Save the updated inventory.
//  4. Update the modification timestamp.
//  5. Print a list of removed cards and a warning for card IDs not found.
//
// File operation errors are reported as error messages.
func Remove(cardIDs ...int) {
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v", err))
		return
	}

	// Load character data for prettier reporting (non-fatal if it fails).
	characters, charErr := tools.LoadCharacters()
	var characterMap map[int]model.Character
	if charErr == nil {
		characterMap = tools.CreateCharacterMap(characters)
	}

	var removedCards []model.CardEntity
	var notFoundCards []int

	for _, cardID := range cardIDs {
		cardIndex := -1
		for i, card := range inventory.Cards {
			if card.ID == cardID {
				cardIndex = i
				break
			}
		}

		if cardIndex == -1 {
			notFoundCards = append(notFoundCards, cardID)
		} else {
			removedCards = append(removedCards, inventory.Cards[cardIndex])
			inventory.Cards = append(inventory.Cards[:cardIndex], inventory.Cards[cardIndex+1:]...)
		}
	}

	if err = tools.SaveInventory(inventory); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error saving inventory: %v", err))
		return
	}

	if len(removedCards) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Removed %d card(s):", len(removedCards)))
		for _, c := range removedCards {
			fmt.Println(tools.FormatCardLabel(c, characterMap))
		}
		_ = tools.UpdateTimeSet()
	}

	if len(notFoundCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Not found in inventory (%d ID(s)):", len(notFoundCards)))
		for _, id := range notFoundCards {
			fmt.Println(id)
		}
	}
}
