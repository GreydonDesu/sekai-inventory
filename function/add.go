package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"sort"
)

// Add introduces new cards to the user's inventory with default initial values.
//
// It supports adding multiple cards in a single operation and keeps the
// inventory sorted by card ID.
//
// New cards are initialized with:
//   - Level:       1
//   - MasteryRank: 0
//   - SkillLevel:  1
//   - SideStory1:  false
//   - SideStory2:  false
//   - Painting:    false
//
// Add performs the following steps:
//
//  1. Load the current inventory.
//  2. Validate each card ID against cards.json.
//  3. Add new cards with default values.
//  4. Sort the inventory by card ID.
//  5. Save the updated inventory and update the modification timestamp.
//
// It prints a detailed summary of added cards, already-owned cards, and
// card IDs that are missing from the database; file operation errors
// are reported as error messages.
func Add(cardIDs ...int) {
	inventory, err := tools.LoadInventory()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading inventory: %v\n", err))
		return
	}

	cards, err := tools.LoadCards()
	if err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error loading card data: %v\n", err))
		return
	}

	// Load character data for prettier reporting (non-fatal if it fails).
	characters, charErr := tools.LoadCharacters()
	var characterMap map[int]model.Character
	if charErr == nil {
		characterMap = tools.CreateCharacterMap(characters)
	}

	cardMap := make(map[int]model.Card, len(cards))
	for _, card := range cards {
		cardMap[card.ID] = card
	}

	added, existing, missingIDs := classifyCardIDs(cardIDs, inventory, cardMap)

	sort.Slice(inventory.Cards, func(i, j int) bool {
		return inventory.Cards[i].ID < inventory.Cards[j].ID
	})

	// Save even when nothing changed to stay consistent.
	if err = tools.SaveInventory(inventory); err != nil {
		tools.PrintErrorMessage(fmt.Sprintf("Error saving inventory: %v\n", err))
		return
	}

	printAddReport(added, existing, missingIDs, characterMap)
}

// classifyCardIDs iterates over the requested IDs, appending new cards to
// inventory and returning three buckets: successfully added, already-owned,
// and IDs missing from the game database.
func classifyCardIDs(cardIDs []int, inventory *model.Inventory, cardMap map[int]model.Card) (added, existing []model.CardEntity, missingIDs []int) {
	for _, cardID := range cardIDs {
		var found *model.CardEntity
		for i := range inventory.Cards {
			if inventory.Cards[i].ID == cardID {
				found = &inventory.Cards[i]
				break
			}
		}

		if found != nil {
			existing = append(existing, *found)
			continue
		}

		cardData, ok := cardMap[cardID]
		if !ok {
			missingIDs = append(missingIDs, cardID)
			continue
		}

		newCard := model.CardEntity{
			Card: model.Card{
				ID:             cardData.ID,
				CharacterID:    cardData.CharacterID,
				CardRarityType: cardData.CardRarityType,
				Attr:           cardData.Attr,
				SupportUnit:    cardData.SupportUnit,
				Prefix:         cardData.Prefix,
			},
			Level:       1,
			MasteryRank: 0,
			SkillLevel:  1,
			SideStory1:  false,
			SideStory2:  false,
			Painting:    false,
		}
		inventory.Cards = append(inventory.Cards, newCard)
		added = append(added, newCard)
	}
	return
}

// printAddReport prints the outcome of an Add operation: which cards were
// added, which were already owned, and which IDs were not found in the
// game database.
func printAddReport(added, existing []model.CardEntity, missingIDs []int, characterMap map[int]model.Character) {
	if len(added) > 0 {
		tools.PrintSuccessMessage(fmt.Sprintf("Added %d card(s):", len(added)))
		for _, c := range added {
			fmt.Println(tools.FormatCardLabel(c, characterMap))
		}
		_ = tools.UpdateTimeSet()
	}
	if len(existing) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Already in inventory (%d card(s)):", len(existing)))
		for _, c := range existing {
			fmt.Println(tools.FormatCardLabel(c, characterMap))
		}
	}
	if len(missingIDs) > 0 {
		tools.PrintWarningMessage(fmt.Sprintf("Not found in database (%d ID(s)):", len(missingIDs)))
		for _, id := range missingIDs {
			fmt.Println(id)
		}
	}
}
