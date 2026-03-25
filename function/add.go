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
//   - Painting: false
//
// The function:
//  1. Loads the current inventory
//  2. Validates each card ID against the game's database
//  3. Adds new cards with default values
//  4. Maintains inventory sorted by card ID
//  5. Updates the modification timestamp
//
// Success/Error reporting:
//   - Successfully added cards are listed with detailed info
//   - Already existing cards are reported as warnings with detailed info
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

	// Load character data for prettier reporting (non-fatal if it fails)
	characters, charErr := tools.LoadCharacters()
	var characterMap map[int]model.Character
	if charErr == nil {
		characterMap = tools.CreateCharacterMap(characters)
	}

	// Create a map of CardID to Card for quick lookup
	cardMap := make(map[int]model.Card)
	for _, card := range cards {
		cardMap[card.ID] = card
	}

	// Helper to create a nice one-line label for a card
	cardLabel := func(card model.CardEntity) string {
		// Rarity (colored)
		rarity := tools.FormatRarity(card.CardRarityType)

		// Character name
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

		// Unit abbreviation: from card.SupportUnit, fallback to character.Unit
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

		return fmt.Sprintf("[%d] %s	%s%s \"%s\"",
			card.ID,
			rarity,
			characterName,
			unitPart,
			card.Prefix,
		)
	}

	// Track cards that were successfully added and those that already exist
	addedCards := []model.CardEntity{}
	existingCards := []model.CardEntity{}
	missingCards := []int{}

	// Iterate over the provided cardIDs
	for _, cardID := range cardIDs {
		// Check if the card already exists in the inventory
		var existing *model.CardEntity
		for i := range inventory.Cards {
			if inventory.Cards[i].ID == cardID {
				existing = &inventory.Cards[i]
				break
			}
		}

		if existing != nil {
			// Card already exists
			existingCards = append(existingCards, *existing)
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
			Painting:   false,
		}

		// Add the new card to the inventory
		inventory.Cards = append(inventory.Cards, newCard)
		addedCards = append(addedCards, newCard)
	}

	// Sort the inventory by card ID
	sort.Slice(inventory.Cards, func(i, j int) bool {
		return inventory.Cards[i].ID < inventory.Cards[j].ID
	})

	// Save the updated inventory (even if nothing changed, to be consistent)
	err = tools.SaveInventory(inventory)
	if err != nil {
		message := fmt.Sprintf("Error saving inventory: %v\n", err)
		tools.PrintErrorMessage(message)
		return
	}

	// Print success message for added cards
	if len(addedCards) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Added %d card(s):", len(addedCards)))
		for _, c := range addedCards {
			fmt.Printf("  %s\n", cardLabel(c))
		}
		_ = tools.UpdateTimeSet()
	}

	// Print warning message for cards that already exist
	if len(existingCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Already in inventory (%d card(s)):", len(existingCards)))
		for _, c := range existingCards {
			fmt.Printf("  %s\n", cardLabel(c))
		}
	}

	// Print warning message for cards not found in Database (cards.json)
	if len(missingCards) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Not found in Database (%d ID(s)):", len(missingCards)))
		for _, id := range missingCards {
			fmt.Printf("  %d\n", id)
		}
	}
}
