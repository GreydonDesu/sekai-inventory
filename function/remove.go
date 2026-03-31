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
	// Load the inventory.
	inventory, err := tools.LoadInventory()
	if err != nil {
		message := fmt.Sprintf("Error loading inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Load character data for prettier reporting (non-fatal if it fails).
	characters, charErr := tools.LoadCharacters()
	var characterMap map[int]model.Character
	if charErr == nil {
		characterMap = tools.CreateCharacterMap(characters)
	}

	// Helper to create a nice one-line label for a card.
	cardLabel := func(card model.CardEntity) string {
		// Rarity (colored).
		rarity := tools.FormatRarity(card.CardRarityType)

		// Character name.
		characterName := "Unknown Character"
		if characterMap != nil {
			if c, ok := characterMap[card.CharacterID]; ok {
				if c.FirstName == "" {
					characterName = c.GivenName
				} else {
					characterName = fmt.Sprintf("%s %s", c.FirstName, c.GivenName)
				}
			}
		}

		// Unit abbreviation: from card.SupportUnit, fallback to character.Unit.
		unitAbbrev := tools.FormatUnit(card.SupportUnit)
		if unitAbbrev == "" && characterMap != nil {
			if c, ok := characterMap[card.CharacterID]; ok {
				unitAbbrev = tools.FormatUnit(c.Unit)
			}
		}
		unitPart := ""
		if unitAbbrev != "" {
			unitPart = fmt.Sprintf(" (%s)", unitAbbrev)
		}

		return fmt.Sprintf("[%d]\t%s\t%s%s \"%s\"",
			card.ID,
			rarity,
			characterName,
			unitPart,
			card.Prefix,
		)
	}

	// Track cards that were successfully removed and those that were not found.
	removedCards := []model.CardEntity{}
	notFoundCards := []int{}

	// Iterate over the provided cardIDs.
	for _, cardID := range cardIDs {
		// Check if the card exists in the inventory.
		cardIndex := -1
		for i, card := range inventory.Cards {
			if card.ID == cardID {
				cardIndex = i
				break
			}
		}

		if cardIndex == -1 {
			// Card not found.
			notFoundCards = append(notFoundCards, cardID)
		} else {
			// Remember the card for reporting.
			removedCard := inventory.Cards[cardIndex]

			// Remove the card from the inventory.
			inventory.Cards = append(inventory.Cards[:cardIndex], inventory.Cards[cardIndex+1:]...)
			removedCards = append(removedCards, removedCard)
		}
	}

	// Save the updated inventory.
	err = tools.SaveInventory(inventory)
	if err != nil {
		message := fmt.Sprintf("Error saving inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Print success message for removed cards.
	if len(removedCards) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Removed %d card(s):", len(removedCards)))
		for _, c := range removedCards {
			fmt.Printf("%s\n", cardLabel(c))
		}
		_ = tools.UpdateTimeSet()
	}

	// Print warning message for cards that were not found.
	if len(notFoundCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Not found in inventory (%d ID(s)):", len(notFoundCards)))
		for _, id := range notFoundCards {
			fmt.Printf("  %d\n", id)
		}
	}
}
