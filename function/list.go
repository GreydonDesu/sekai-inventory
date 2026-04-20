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
		model.RarityType1:        0,
		model.RarityType2:        0,
		model.RarityType3:        0,
		model.RarityType4:        0,
		model.RarityTypeBirthday: 0,
	}
	for _, card := range filteredCards {
		rarityCounts[card.CardRarityType]++
	}

	tools.PrintSuccessMessage(fmt.Sprintf("Inventory Stats (Total: %d):", len(filteredCards)))
	fmt.Printf("  %s\t%d\n", tools.FormatRarity(model.RarityType4), rarityCounts[model.RarityType4])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity(model.RarityTypeBirthday), rarityCounts[model.RarityTypeBirthday])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity(model.RarityType3), rarityCounts[model.RarityType3])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity(model.RarityType2), rarityCounts[model.RarityType2])
	fmt.Printf("  %s\t%d\n", tools.FormatRarity(model.RarityType1), rarityCounts[model.RarityType1])

	tools.PrintSuccessMessage("--- Inventory List ---")
	for _, card := range filteredCards {
		fmt.Println(tools.FormatCardDetails(card, characterMap))
	}
}

// matchesFilters reports whether card satisfies all provided filters.
func matchesFilters(card model.CardEntity, filters map[string]string, characterMap map[int]model.Character) bool {
	for field, value := range filters {
		switch field {
		case fieldCharacter:
			character, exists := characterMap[card.CharacterID]
			if !exists || !tools.ContainsIgnoreCase(tools.FormatCharacterName(character), value) {
				return false
			}
		case fieldRarity:
			if card.CardRarityType != tools.RarityToKey[value] {
				return false
			}
		case fieldGroup:
			expectedGroup := tools.GroupToKey[value]
			if card.SupportUnit != expectedGroup {
				character, exists := characterMap[card.CharacterID]
				if !exists || character.Unit != expectedGroup {
					return false
				}
			}
		case fieldPainting:
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
