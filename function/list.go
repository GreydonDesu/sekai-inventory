package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
)

// List displays the user's card inventory with optional filtering.
// The output is sorted by card ID and includes detailed card information
// such as rarity, level, attributes, and character details.
//
// Supported filters:
//   - character: Filter by character name (case-insensitive partial match)
//   - rarity: Filter by card rarity (1, 2, 3, 4, bd)
//   - group: Filter by unit/group (L/N, MMJ, VBS, WxS, N25, VS)
//
// The function handles error cases such as:
//   - Failed inventory loading
//   - Failed character data loading
//   - Invalid filter fields
//   - No matching cards
func List(filters map[string]string) {
	// Load the inventory
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
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

	// Filter the inventory based on the provided filters
	var filteredCards []model.CardEntity
	for _, card := range inventory.Cards {
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

	// Sort the filtered inventory by card ID
	sort.Slice(filteredCards, func(i, j int) bool {
		return filteredCards[i].ID < filteredCards[j].ID
	})

	// Display the filtered and sorted inventory
	if len(filteredCards) == 0 {
		tools.PrintWarningMessage("No matching cards found.")
		return
	}

	tools.PrintSuccessMessage("Inventory:")
	for _, card := range filteredCards {
		// Use the utility function to format the card details
		cardDetails := tools.FormatCardDetails(card, characterMap)
		fmt.Println(cardDetails)
	}
}
