package function

import (
	"fmt"
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"strconv"

	"github.com/fatih/color"
)

// Change modifies specific attributes of a card in the user's inventory.
//
// It supports updating multiple fields in a single operation. Supported fields
// and their valid values are:
//
//   - level:       integer between 1 and 60.
//   - skillLevel:  integer between 1 and 5.
//   - masteryRank: integer between 0 and 5.
//   - sideStory1:  boolean (true/false).
//   - sideStory2:  boolean (true/false).
//   - painting:    boolean (true/false).
//
// Change returns an error if the card does not exist, a field name is unknown,
// a value is invalid, or the inventory cannot be saved. On success it prints a
// detailed, colorized summary of the changes.
func Change(cardID int, updates map[string]string) error {
	// Load the inventory.
	inventory, err := tools.LoadInventory()
	if err != nil {
		return fmt.Errorf("error loading inventory: %v", err)
	}

	// Find the card in the inventory.
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

	// Keep a copy of the original state for comparison.
	original := *card

	// Track whether any actual change was made.
	changed := false

	// Apply updates.
	for field, value := range updates {
		switch field {
		case "level":
			level, err := strconv.Atoi(value)
			if err != nil || level < 1 || level > 60 {
				return fmt.Errorf("invalid value for 'level': %s. Must be an integer between 1 and 60", value)
			}
			card.Level = level
			if card.Level != original.Level {
				changed = true
			}
		case "skillLevel":
			skillLevel, err := strconv.Atoi(value)
			if err != nil || skillLevel < 1 || skillLevel > 5 {
				return fmt.Errorf("invalid value for 'skillLevel': %s. Must be an integer between 1 and 5", value)
			}
			card.SkillLevel = skillLevel
			if card.SkillLevel != original.SkillLevel {
				changed = true
			}
		case "masteryRank":
			masteryRank, err := strconv.Atoi(value)
			if err != nil || masteryRank < 0 || masteryRank > 5 {
				return fmt.Errorf("invalid value for 'masteryRank': %s. Must be an integer between 0 and 5", value)
			}
			card.MasteryRank = masteryRank
			if card.MasteryRank != original.MasteryRank {
				changed = true
			}
		case "sideStory1":
			sideStory1, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for 'sideStory1': %s. Must be 'true' or 'false'", value)
			}
			card.SideStory1 = sideStory1
			if card.SideStory1 != original.SideStory1 {
				changed = true
			}
		case "sideStory2":
			sideStory2, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for 'sideStory2': %s. Must be 'true' or 'false'", value)
			}
			card.SideStory2 = sideStory2
			if card.SideStory2 != original.SideStory2 {
				changed = true
			}
		case "painting":
			painting, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("invalid value for 'painting': %s. Must be 'true' or 'false'", value)
			}
			card.Painting = painting
			if card.Painting != original.Painting {
				changed = true
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}

	// If nothing actually changed, do not touch the file.
	if !changed {
		fmt.Printf("No changes made for card with ID %d.\n", cardID)
		return nil
	}

	// Save the updated inventory.
	err = tools.SaveInventory(inventory)
	if err != nil {
		return fmt.Errorf("error saving inventory: %v", err)
	}

	// Print a detailed summary of changes.
	printChangeSummary(card, &original)

	// Update timestamp.
	_ = tools.UpdateTimeSet()

	return nil
}

// printChangeSummary prints a user-friendly, colorized summary of the changes
// made to a card. It shows numeric fields as "old > new" and boolean fields
// as colored checkboxes.
func printChangeSummary(card, original *model.CardEntity) {
	// Precompute colors (same as in FormatCardDetails).
	rGreen, gGreen, bGreen, _ := tools.HexToRGB("#00ff00")
	rRed, gRed, bRed, _ := tools.HexToRGB("#ff0000")

	green := color.RGB(rGreen, gGreen, bGreen)
	red := color.RGB(rRed, gRed, bRed)

	// Try to load character data for a nicer header; fall back if it fails.
	characters, err := tools.LoadCharacters()
	var header string
	if err != nil {
		// Fallback header without character info.
		header = fmt.Sprintf("Changes for ID %d:", card.ID)
	} else {
		characterMap := tools.CreateCharacterMap(characters)
		character, exists := characterMap[card.CharacterID]

		characterName := "Unknown Character"
		if exists {
			if character.FirstName == "" {
				characterName = character.GivenName
			} else {
				characterName = fmt.Sprintf("%s %s", character.FirstName, character.GivenName)
			}
		}

		// Determine unit abbreviation (card.SupportUnit first, then character.Unit).
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

	// Fixed width for the label column so everything lines up.
	const labelWidth = 13

	// Helper for numeric "old > new" with color.
	printNumericChange := func(label string, oldVal, newVal int) {
		fmt.Printf(
			"  %-*s %s > %s\n",
			labelWidth,
			label,
			red.Sprintf("%d", oldVal),
			green.Sprintf("%d", newVal),
		)
	}

	// Boolean fields: show final state ☑/☐ in color.
	boolMark := func(b bool) string {
		if b {
			return green.Sprint("☑")
		}
		return red.Sprint("☐")
	}

	// Numeric fields: show old > new if they actually changed.
	if original.Level != card.Level {
		printNumericChange("Level", original.Level, card.Level)
	}
	if original.MasteryRank != card.MasteryRank {
		printNumericChange("Mastery Rank", original.MasteryRank, card.MasteryRank)
	}
	if original.SkillLevel != card.SkillLevel {
		printNumericChange("Skill Level", original.SkillLevel, card.SkillLevel)
	}

	// Boolean fields: only print if they changed.
	if original.SideStory1 != card.SideStory1 {
		fmt.Printf("  %-*s %s\n", labelWidth, "Side Story 1", boolMark(card.SideStory1))
	}
	if original.SideStory2 != card.SideStory2 {
		fmt.Printf("  %-*s %s\n", labelWidth, "Side Story 2", boolMark(card.SideStory2))
	}
	if original.Painting != card.Painting {
		fmt.Printf("  %-*s %s\n", labelWidth, "Painting", boolMark(card.Painting))
	}
}
