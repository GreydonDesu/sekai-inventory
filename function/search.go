package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
)

// Search identifies cards that exist in the game but are not in the user's inventory.
// This helps users discover which cards they're missing and can acquire.
//
// Supported filters:
//   - character: Filter by character name (case-insensitive partial match)
//   - rarity: Filter by card rarity (1, 2, 3, 4, bd)
//   - group: Filter by unit/group (L/N, MMJ, VBS, WxS, N25, VS)
//
// The function:
//  1. Loads the user's inventory and game's card database
//  2. Identifies cards not in the inventory
//  3. Applies any specified filters
//  4. Displays matching cards in a sorted, formatted list
//
// Error cases handled:
//   - Failed inventory loading
//   - Failed card database loading
//   - Failed character data loading
//   - Invalid filter fields
//   - No matching cards
func Search(filters map[string]string) {
	// Load the inventory
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
		return
	}

	// Load the full list of cards from cards.json
	allCards, err := tools.LoadCards()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading cards.json: %v\n", err))
		return
	}

	// Load the character data
	characters, err := tools.LoadCharacters()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading character data: %v\n", err))
		return
	}

	// Create a map of CharacterID to Character
	characterMap := tools.CreateCharacterMap(characters)

	// Create a map of card IDs in the inventory for quick lookup
	inventoryCardIDs := make(map[int]bool)
	for _, card := range inventory.Cards {
		inventoryCardIDs[card.ID] = true
	}

	// Find cards that are in cards.json but not in the inventory
	var missingCards []model.Card
	for _, card := range allCards {
		if !inventoryCardIDs[card.ID] {
			missingCards = append(missingCards, card)
		}
	}

	// Filter the missing cards based on the provided filters
	var filteredCards []model.Card
	for _, card := range missingCards {
		matches := true
		for field, value := range filters {
			switch field {
			case "character":
				// Match by character's given name
				character, exists := characterMap[card.CharacterID]
				if !exists || !tools.ContainsIgnoreCase(character.FirstName+" "+character.GivenName, value) {
					matches = false
				}
			case "rarity":
				// Match by card rarity
				expectedRarity := map[string]string{
					"1":  "rarity_1",
					"2":  "rarity_2",
					"3":  "rarity_3",
					"4":  "rarity_4",
					"bd": "rarity_birthday",
				}[value]
				if card.CardRarityType != expectedRarity {
					matches = false
				}
			case "group":
				// Match by support unit (from card or character data)
				expectedGroup := map[string]string{
					"L/N": "light_sound",
					"MMJ": "idol",
					"VBS": "street",
					"WxS": "theme_park",
					"N25": "school_refusal",
					"VS":  "piapro",
				}[value]

				// Check supportUnit in the card data
				if card.SupportUnit != expectedGroup {
					// Fallback: Check Unit in the character data
					character, exists := characterMap[card.CharacterID]
					if !exists || character.Unit != expectedGroup {
						matches = false
					}
				}
			default:
				// Handle unknown filter fields
				tools.PrintWarningMessage(fmt.Sprintf("Unknown filter field: %s", field))
				matches = false
			}

			// If any filter doesn't match, skip this card
			if !matches {
				break
			}
		}

		// Add the card to the filtered list if it matches all filters
		if matches {
			filteredCards = append(filteredCards, card)
		}
	}

	// Sort the filtered cards by card ID
	sort.Slice(filteredCards, func(i, j int) bool {
		return filteredCards[i].ID < filteredCards[j].ID
	})

	// Display the filtered and sorted cards
	if len(filteredCards) == 0 {
		tools.PrintWarningMessage("No matching cards found.")
		return
	}

	tools.PrintSuccessMessage(fmt.Sprintf("Found %d matching cards:", len(filteredCards)))
	for _, card := range filteredCards {
		// Use the utility function to format the card details
		cardDetails := tools.FormatCardDetails(model.CardEntity{Card: card}, characterMap)
		fmt.Println(cardDetails)
	}
}
