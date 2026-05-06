package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"strconv"

	"github.com/fatih/color"
)

// Card field names used as keys in filter and update maps across commands.
const (
	fieldLevel       = "level"
	fieldSkillLevel  = "skillLevel"
	fieldMasteryRank = "masteryRank"
	fieldSideStory1  = "sideStory1"
	fieldSideStory2  = "sideStory2"
	fieldPainting    = "painting"
	fieldCharacter   = "character"
	fieldRarity      = "rarity"
	fieldGroup       = "group"
)

// Change modifies specific attributes of a card in the user's inventory.
//
// It supports updating multiple fields in a single operation. Supported fields
// and their valid values are:
//
//   - level:       integer between 1 and 60.
//   - skillLevel:  integer between 1 and 4.
//   - masteryRank: integer between 0 and 5.
//   - sideStory1:  boolean (true/false).
//   - sideStory2:  boolean (true/false).
//   - painting:    boolean (true/false).
//
// Change returns an error if the card does not exist, a field name is unknown,
// a value is invalid, or the inventory cannot be saved. On success it prints a
// detailed, colorized summary of the changes.
func Change(cardID int, updates map[string]string) error {
	inventory, err := tools.LoadInventory()
	if err != nil {
		return fmt.Errorf("error loading inventory: %v", err)
	}

	var card *model.CardEntity
	for i := range inventory.Cards {
		if inventory.Cards[i].ID == cardID {
			card = &inventory.Cards[i]
			break
		}
	}
	if card == nil {
		return fmt.Errorf("no card found with ID %d", cardID)
	}

	original := *card
	for field, value := range updates {
		if err := applyCardField(card, field, value); err != nil {
			return err
		}
	}

	if *card == original {
		fmt.Printf("No changes made for card with ID %d.\n", cardID)
		return nil
	}

	if err = tools.SaveInventory(inventory); err != nil {
		return fmt.Errorf("error saving inventory: %v", err)
	}

	printChangeSummary(card, &original)
	_ = tools.UpdateTimeSet()

	return nil
}

// applyCardField applies a single field update to card. It returns an error if
// the field name is unrecognized or the value fails validation.
func applyCardField(card *model.CardEntity, field, value string) error {
	switch field {
	case fieldLevel:
		return applyLevel(card, value)
	case fieldSkillLevel:
		return applySkillLevel(card, value)
	case fieldMasteryRank:
		return applyMasteryRank(card, value)
	case fieldSideStory1:
		return applySideStory1(card, value)
	case fieldSideStory2:
		return applySideStory2(card, value)
	case fieldPainting:
		return applyPainting(card, value)
	default:
		return fmt.Errorf("unknown field: %s", field)
	}
}

// applyLevel validates and sets the card's level (1–60).
func applyLevel(card *model.CardEntity, value string) error {
	v, err := parseIntField(value, fieldLevel, 1, 60)
	if err != nil {
		return err
	}
	card.Level = v
	return nil
}

// applySkillLevel validates and sets the card's skill level (1–4).
func applySkillLevel(card *model.CardEntity, value string) error {
	v, err := parseIntField(value, fieldSkillLevel, 1, 4)
	if err != nil {
		return err
	}
	card.SkillLevel = v
	return nil
}

// applyMasteryRank validates and sets the card's mastery rank (0–5).
func applyMasteryRank(card *model.CardEntity, value string) error {
	v, err := parseIntField(value, fieldMasteryRank, 0, 5)
	if err != nil {
		return err
	}
	card.MasteryRank = v
	return nil
}

// applySideStory1 validates and sets the card's side story 1 status.
func applySideStory1(card *model.CardEntity, value string) error {
	v, err := parseBoolField(value, fieldSideStory1)
	if err != nil {
		return err
	}
	card.SideStory1 = v
	return nil
}

// applySideStory2 validates and sets the card's side story 2 status.
func applySideStory2(card *model.CardEntity, value string) error {
	v, err := parseBoolField(value, fieldSideStory2)
	if err != nil {
		return err
	}
	card.SideStory2 = v
	return nil
}

// applyPainting validates and sets the card's painting status.
func applyPainting(card *model.CardEntity, value string) error {
	v, err := parseBoolField(value, fieldPainting)
	if err != nil {
		return err
	}
	card.Painting = v
	return nil
}

// parseIntField converts value to an integer and validates it falls within
// [min, max]. Returns a descriptive error if conversion or range check fails.
func parseIntField(value, fieldName string, minVal, maxVal int) (int, error) {
	v, err := strconv.Atoi(value)
	if err != nil || v < minVal || v > maxVal {
		return 0, fmt.Errorf("invalid value for '%s': %s. Must be an integer between %d and %d", fieldName, value, minVal, maxVal)
	}
	return v, nil
}

// parseBoolField converts value to a boolean. Returns a descriptive error if
// the value is not a valid boolean string.
func parseBoolField(value, fieldName string) (bool, error) {
	v, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("invalid value for '%s': %s. Must be 'true' or 'false'", fieldName, value)
	}
	return v, nil
}

// printChangeSummary prints a user-friendly, colorized summary of the changes
// made to a card. Numeric fields are shown as "old > new" and boolean fields
// as colorized checkboxes.
func printChangeSummary(card, original *model.CardEntity) {
	rGreen, gGreen, bGreen, _ := tools.HexToRGB("#00ff00")
	rRed, gRed, bRed, _ := tools.HexToRGB("#ff0000")
	green := color.RGB(rGreen, gGreen, bGreen)
	red := color.RGB(rRed, gRed, bRed)

	// Try to load character data for a nicer header; fall back if it fails.
	var header string
	if characters, err := tools.LoadCharacters(); err != nil {
		header = fmt.Sprintf("Changes for ID %d:", card.ID)
	} else {
		characterMap := tools.CreateCharacterMap(characters)
		character, exists := characterMap[card.CharacterID]

		characterName := "Unknown Character"
		if exists {
			characterName = tools.FormatCharacterName(character)
		}

		unitAbbrev := tools.FormatUnit(card.SupportUnit)
		if unitAbbrev == "" && exists {
			unitAbbrev = tools.FormatUnit(character.Unit)
		}

		unitPart := ""
		if unitAbbrev != "" {
			unitPart = fmt.Sprintf(" (%s)", unitAbbrev)
		}

		header = fmt.Sprintf("Changes for ID %d - %s%s \"%s\"", card.ID, characterName, unitPart, card.Prefix)
	}

	fmt.Println(header)

	const labelWidth = 13

	printNumericChange := func(label string, oldVal, newVal int) {
		fmt.Printf("  %-*s %s > %s\n", labelWidth, label,
			red.Sprintf("%d", oldVal),
			green.Sprintf("%d", newVal),
		)
	}

	if original.Level != card.Level {
		printNumericChange("Level", original.Level, card.Level)
	}
	if original.MasteryRank != card.MasteryRank {
		printNumericChange("Mastery Rank", original.MasteryRank, card.MasteryRank)
	}
	if original.SkillLevel != card.SkillLevel {
		printNumericChange("Skill Level", original.SkillLevel, card.SkillLevel)
	}
	if original.SideStory1 != card.SideStory1 {
		fmt.Printf("  %-*s %s\n", labelWidth, "Side Story 1", tools.FormatBool(card.SideStory1))
	}
	if original.SideStory2 != card.SideStory2 {
		fmt.Printf("  %-*s %s\n", labelWidth, "Side Story 2", tools.FormatBool(card.SideStory2))
	}
	if original.Painting != card.Painting {
		fmt.Printf("  %-*s %s\n", labelWidth, "Painting", tools.FormatBool(card.Painting))
	}
}
