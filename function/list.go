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
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
		return
	}

	characters, err := tools.LoadCharacters()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading character data: %v\n", err))
		return
	}

	characterMap := tools.CreateCharacterMap(characters)

	var filteredCards []model.CardEntity
	for _, card := range inventory.Cards {
		if matchesFilters(card, filters, characterMap) {
			filteredCards = append(filteredCards, card)
		}
	}

	sort.Slice(filteredCards, func(i, j int) bool {
		return filteredCards[i].ID < filteredCards[j].ID
	})

	if len(filteredCards) == 0 {
		tools.PrintWarningMessage("No matching cards found.")
		return
	}

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
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_4"), rarityCounts["rarity_4"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_birthday"), rarityCounts["rarity_birthday"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_3"), rarityCounts["rarity_3"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_2"), rarityCounts["rarity_2"])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity("rarity_1"), rarityCounts["rarity_1"])

	tools.PrintSuccessMessage("--- Inventory List ---")
	for _, card := range filteredCards {
		fmt.Println(tools.FormatCardDetails(card, characterMap))
	}
}

// matchesFilters reports whether card satisfies all provided filters.
func matchesFilters(card model.CardEntity, filters map[string]string, characterMap map[int]model.Character) bool {
	for field, value := range filters {
		switch field {
		case "character":
			character, exists := characterMap[card.CharacterID]
			if !exists || !tools.ContainsIgnoreCase(character.FirstName+" "+character.GivenName, value) {
				return false
			}
		case "rarity":
			if card.CardRarityType != tools.RarityToKey[value] {
				return false
			}
		case "group":
			expectedGroup := tools.GroupToKey[value]
			if card.SupportUnit != expectedGroup {
				character, exists := characterMap[card.CharacterID]
				if !exists || character.Unit != expectedGroup {
					return false
				}
			}
		case "painting":
			expectedPainting, err := strconv.ParseBool(value)
			if err != nil {
				tools.PrintWarningMessage(fmt.Sprintf("Invalid value for 'painting': %s. Must be 'true' or 'false'", value))
				return false
			}
			if card.Painting != expectedPainting {
				return false
			}
		default:
			tools.PrintWarningMessage(fmt.Sprintf("Unknown filter field: %s", field))
			return false
		}
	}
	return true
}
