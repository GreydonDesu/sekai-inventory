package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
)

// Add introduces new cards to the user's inventory with default initial values.
// This function supports adding multiple cards in a single operation and maintains
// the inventory in a sorted state.
//
// Default values for new cards:
//   - Level: 1
//   - MasterRank: 0
//   - SkillLevel: 1
//   - SideStory1: false
//   - SideStory2: false
//
// The function:
//  1. Loads the current inventory
//  2. Validates each card ID against the game's database
//  3. Adds new cards with default values
//  4. Maintains inventory sorted by card ID
//  5. Updates the modification timestamp
//
// Success/Error reporting:
//   - Successfully added cards are listed
//   - Already existing cards are reported as warnings
//   - Cards not found in the database are reported as warnings
//   - File operation errors are reported as errors
//
// Parameters:
//   - cardIDs: One or more card IDs to add to the inventory
func Add(cardIDs ...int) {
	// Load the inventory
	inventory, err := tools.LoadInventory()
	if err != nil {
		message := fmt.Sprintf("Error loading inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Load the card data from cards.json
	cards, err := tools.LoadCards()
	if err != nil {
		message := fmt.Sprintf("Error loading card data: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Create a map of CardID to Card for quick lookup
	cardMap := make(map[int]model.Card)
	for _, card := range cards {
		cardMap[card.ID] = card
	}

	// Track cards that were successfully added and those that already exist
	addedCards := []int{}
	existingCards := []int{}
	missingCards := []int{}

	// Iterate over the provided cardIDs
	for _, cardID := range cardIDs {
		// Check if the card already exists in the inventory
		cardExists := false
		for _, card := range inventory.Cards {
			if card.ID == cardID {
				cardExists = true
				break
			}
		}

		if cardExists {
			// Card already exists
			existingCards = append(existingCards, cardID)
			continue
		}

		// Fetch card data from cards.json
		cardData, exists := cardMap[cardID]
		if !exists {
			// Card not found in cards.json
			missingCards = append(missingCards, cardID)
			continue
		}

		// Create a new CardEntity with data from cards.json and default values
		newCard := model.CardEntity{
			Card: model.Card{
				ID:             cardData.ID,
				CharacterID:    cardData.CharacterID,
				CardRarityType: cardData.CardRarityType,
				Attr:           cardData.Attr,
				SupportUnit:    cardData.SupportUnit,
				Prefix:         cardData.Prefix,
			},
			Level:      1,
			MasterRank: 0,
			SkillLevel: 1,
			SideStory1: false,
			SideStory2: false,
		}

		// Add the new card to the inventory
		inventory.Cards = append(inventory.Cards, newCard)
		addedCards = append(addedCards, cardID)
	}

	// Sort the inventory by card ID
	sort.Slice(inventory.Cards, func(i, j int) bool {
		return inventory.Cards[i].ID < inventory.Cards[j].ID
	})

	// Save the updated inventory
	err = tools.SaveInventory(inventory)
	if err != nil {
		message := fmt.Sprintf("Error saving inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Print success message for added cards
	if len(addedCards) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Added cards with IDs: %v", addedCards))
		tools.UpdateTimeSet()
	}

	// Print warning message for cards that already exist
	if len(existingCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Cards with IDs already exist: %v", existingCards))
	}

	// Print warning message for cards not found in cards.json
	if len(missingCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Cards with IDs not found in cards.json: %v", missingCards))
	}

}
