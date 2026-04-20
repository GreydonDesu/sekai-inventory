package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
)

// Search identifies cards that exist in the game but are not in the user's
// inventory. This helps users discover which cards they can still acquire.
//
// Supported filters:
//
//   - character: filter by character name (case-insensitive partial match).
//   - rarity:    filter by card rarity (1, 2, 3, 4, bd).
//   - group:     filter by unit/group (L/N, MMJ, VBS, WxS, N25, VS).
//
// Search reports errors if any required data file cannot be loaded, and prints
// a warning if no cards match the filters.
func Search(filters map[string]string) {
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
		return
	}

	allCards, err := tools.LoadCards()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading cards.json: %v\n", err))
		return
	}

	characters, err := tools.LoadCharacters()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading character data: %v\n", err))
		return
	}

	characterMap := tools.CreateCharacterMap(characters)

	inventoryIDs := make(map[int]bool, len(inventory.Cards))
	for _, card := range inventory.Cards {
		inventoryIDs[card.ID] = true
	}

	var filteredCards []model.Card
	for _, card := range allCards {
		if inventoryIDs[card.ID] {
			continue
		}
		if matchesCardFilters(card, filters, characterMap) {
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

	tools.PrintSuccessMessage(fmt.Sprintf("Found %d matching cards:", len(filteredCards)))
	for _, card := range filteredCards {
		fmt.Println(tools.FormatCardDetails(model.CardEntity{Card: card}, characterMap))
	}
}

// matchesCardFilters reports whether a Card (not in inventory) satisfies all
// provided filters. It mirrors matchesFilters in list.go but operates on the
// base Card type and omits the painting filter (not applicable to unowned cards).
func matchesCardFilters(card model.Card, filters map[string]string, characterMap map[int]model.Character) bool {
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
		default:
			tools.PrintWarningMessage(fmt.Sprintf("Unknown filter field: %s", field))
			return false
		}
	}
	return true
}
