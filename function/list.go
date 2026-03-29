package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
	"strconv"
)

// List displays the user's card inventory with optional filtering.
//
// The output is sorted by card ID and includes detailed card information such
// as rarity, level, attributes, and character details. It also prints a rarity
// summary for the filtered result before listing individual cards.
//
// Supported filters:
//
//   - character: filter by character name (case-insensitive partial match).
//   - rarity:    filter by card rarity (1, 2, 3, 4, bd).
//   - group:     filter by unit/group (L/N, MMJ, VBS, WxS, N25, VS).
//   - painting:  filter by painting status (true/false).
//
// List reports errors if the inventory or character data cannot be loaded, and
// prints a warning if no cards match the filters.
func List(filters map[string]string) {
	// Load the inventory.
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
		return
	}

	// Load the character data.
	characters, err := tools.LoadCharacters()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading character data: %v\n", err))
		return
	}

	// Create a map of CharacterID to Character.
	characterMap := tools.CreateCharacterMap(characters)

	// Filter the inventory based on the provided filters.
	var filteredCards []model.CardEntity
	for _, card := range inventory.Cards {
		matches := true
		for field, value := range filters {
			switch field {
			case "character":
				// Match by character's full name.
				character, exists := characterMap[card.CharacterID]
				if !exists || !tools.ContainsIgnoreCase(character.FirstName+" "+character.GivenName, value) {
					matches = false
				}
			case "rarity":
				// Match by card rarity.
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
				// Match by support unit (from card or character data).
				expectedGroup := map[string]string{
					"L/N": "light_sound",
					"MMJ": "idol",
					"VBS": "street",
					"WxS": "theme_park",
					"N25": "school_refusal",
					"VS":  "piapro",
				}[value]

				// Check supportUnit in the card data.
				if card.SupportUnit != expectedGroup {
					// Fallback: check Unit in the character data.
					character, exists := characterMap[card.CharacterID]
					if !exists || character.Unit != expectedGroup {
						matches = false
					}
				}
			case "painting":
				// Match by painting status.
				expectedPainting, err := strconv.ParseBool(value)
				if err != nil {
					tools.PrintWarningMessage(fmt.Sprintf("Invalid value for 'painting': %s. Must be 'true' or 'false'", value))
					matches = false
				} else if card.Painting != expectedPainting {
					matches = false
				}
			default:
				// Handle unknown filter fields.
				tools.PrintWarningMessage(fmt.Sprintf("Unknown filter field: %s", field))
				matches = false
			}

			// If any filter does not match, skip this card.
			if !matches {
				break
			}
		}

		// Add the card to the filtered list if it matches all filters.
		if matches {
			filteredCards = append(filteredCards, card)
		}
	}

	// Sort the filtered inventory by card ID.
	sort.Slice(filteredCards, func(i, j int) bool {
		return filteredCards[i].ID < filteredCards[j].ID
	})

	// Display the filtered and sorted inventory.
	if len(filteredCards) == 0 {
		tools.PrintWarningMessage("No matching cards found.")
		return
	}

	// ---- Inventory stats header ----

	// Count rarities in the filtered result.
	rarityCounts := map[string]int{
		"rarity_1":        0,
		"rarity_2":        0,
		"rarity_3":        0,
		"rarity_4":        0,
		"rarity_birthday": 0,
	}

	for _, card := range filteredCards {
		rarityCounts[card.CardRarityType]++
	}

	tools.PrintSuccessMessage(fmt.Sprintf("Inventory Stats (Total: %d):", len(filteredCards)))
	// Order chosen to highlight higher rarities first; adjust if desired.
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_4"), rarityCounts["rarity_4"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_birthday"), rarityCounts["rarity_birthday"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_3"), rarityCounts["rarity_3"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_2"), rarityCounts["rarity_2"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_1"), rarityCounts["rarity_1"])

	// ---- Inventory list ----
	tools.PrintSuccessMessage("--- Inventory List ---")
	for _, card := range filteredCards {
		cardDetails := tools.FormatCardDetails(card, characterMap)
		fmt.Println(cardDetails)
	}
}
