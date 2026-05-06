package tools

import (
	"fmt"
	"sekai-inventory/model"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// rarity label map used by FormatRarity; avoids repeating string literals.
var rarityLabels = map[string]string{
	model.RarityType1:        "★",
	model.RarityType2:        "★★",
	model.RarityType3:        "★★★",
	model.RarityType4:        "★★★★",
	model.RarityTypeBirthday: "୨୧",
}

// FormatBool returns a colorized checkbox symbol representing a boolean value:
// ☑ (green) for true and ☐ (red) for false.
func FormatBool(b bool) string {
	if b {
		r, g, bl, _ := HexToRGB("#00ff00") //nolint:errcheck // hardcoded valid hex
		return color.RGB(r, g, bl).Sprint("☑")
	}
	r, g, bl, _ := HexToRGB("#ff0000") //nolint:errcheck // hardcoded valid hex
	return color.RGB(r, g, bl).Sprint("☐")
}

// FormatCharacterName returns the display name for a character. For characters
// with both a family name and a given name the result is "FirstName GivenName".
// For characters with no family name only the given name is returned.
func FormatCharacterName(c model.Character) string {
	if c.FirstName == "" {
		return c.GivenName
	}
	return c.FirstName + " " + c.GivenName
}

// FormatCardLabel returns a compact, colorized one-liner for a card, used in
// add and remove operation summaries. It shows the card ID, rarity stars,
// character name, unit abbreviation, and card title prefix.
//
// A nil characterMap is safe; affected fields fall back to "Unknown Character".
func FormatCardLabel(card model.CardEntity, characterMap map[int]model.Character) string {
	rarity := FormatRarity(card.CardRarityType)
	characterName := "Unknown Character"
	unitAbbrev := FormatUnit(card.SupportUnit)

	if c, ok := characterMap[card.CharacterID]; ok {
		characterName = FormatCharacterName(c)
		if unitAbbrev == "" {
			unitAbbrev = FormatUnit(c.Unit)
		}
	}

	unitPart := ""
	if unitAbbrev != "" {
		unitPart = fmt.Sprintf(" (%s)", unitAbbrev)
	}

	return fmt.Sprintf("[%d]\t%s\t%s%s \"%s\"",
		card.ID,
		rarity,
		characterName,
		unitPart,
		card.Prefix,
	)
}

// FormatCardDetails returns a human-readable, colorized string representation
// of a card. It combines card properties with character information and applies
// formatting based on rarity, level, side stories, and painting status.
//
// The characterMap is used to resolve CharacterID to character names and units.
func FormatCardDetails(card model.CardEntity, characterMap map[int]model.Character) string {
	character, exists := characterMap[card.CharacterID]
	characterName := "Unknown Character"
	if exists {
		characterName = FormatCharacterName(character)
	}

	unitAbbrev := FormatUnit(card.SupportUnit)
	if unitAbbrev == "" && exists {
		unitAbbrev = FormatUnit(character.Unit)
	}

	var unitPart string
	if unitAbbrev != "" {
		unitPart = fmt.Sprintf(" (%s)", unitAbbrev)
	}

	// Painting.
	var painting string
	if card.Painting {
		r, g, b, _ := HexToRGB("#00ff00")
		painting = color.RGB(r, g, b).Sprint("☑")
	} else {
		r, g, b, _ := HexToRGB("#ff0000")
		painting = color.RGB(r, g, b).Sprint("☐")
	}

	// Highlight Mastery Rank and Skill Level with RGB green at max values.
	rGreen, gGreen, bGreen, _ := HexToRGB("#00ff00")
	green := color.RGB(rGreen, gGreen, bGreen)

	masteryRank := fmt.Sprintf("MR%d", card.MasteryRank)
	if card.MasteryRank == 5 {
		masteryRank = color.RGB(rGreen, gGreen, bGreen).Sprint(masteryRank)
	}

	skillLevel := fmt.Sprintf("SL%d", card.SkillLevel)
	if card.SkillLevel == 4 {
		skillLevel = green.Sprint(skillLevel)
	}

	return fmt.Sprintf("[%d]\t%s\t%s\t%s\t| %s | %s | Side Story 1: %s | Side Story 2: %s | Painting: %s | %s%s \"%s\"",
		card.ID,
		FormatRarity(card.CardRarityType),
		FormatLevel(card.CardRarityType, card.Level),
		FormatAttribute(card.Attr),
		masteryRank,
		skillLevel,
		FormatBool(card.SideStory1),
		FormatBool(card.SideStory2),
		painting,
		characterName,
		unitPart,
		card.Prefix,
	)
}

// FormatRarity converts a raw rarity value into a visually appealing string.
// Each rarity level is represented by colored stars; birthday cards use a
// special symbol.
//
// Rarity levels:
//
//   - rarity_1:        ★    (yellow)
//   - rarity_2:       ★★   (yellow)
//   - rarity_3:      ★★★  (yellow)
//   - rarity_4:     ★★★★ (yellow)
//   - rarity_birthday: ୨୧   (magenta)
func FormatRarity(rarity string) string {
	label, ok := rarityLabels[rarity]
	if !ok {
		return rarity
	}
	if rarity == model.RarityTypeBirthday {
		r, g, b, _ := HexToRGB("#ff00ff")
		return color.RGB(r, g, b).Sprint(label)
	}
	rYellow, gYellow, bYellow, _ := HexToRGB("#ffd700")
	return color.RGB(rYellow, gYellow, bYellow).Sprint(label)
}

// FormatLevel formats a card's level with color coding based on rarity-specific
// maximum level thresholds. The color indicates the card's progression:
//
//   - Green:  at maximum level for the card's rarity.
//   - Yellow: approaching maximum level (rarity_3 at 40, rarity_4 at 50).
//   - Plain:  below the thresholds above.
//
// Level thresholds by rarity:
//
//   - rarity_1:        20 (green/max)
//   - rarity_2:        30 (green/max)
//   - rarity_3:        40 (yellow), 50 (green/max)
//   - rarity_4:        50 (yellow), 60 (green/max)
//   - rarity_birthday: 60 (green/max)
func FormatLevel(rarity string, level int) string {
	if level <= 0 {
		return "Lvl 0"
	}

	rGreen, gGreen, bGreen, _ := HexToRGB("#00ff00")
	rYellow, gYellow, bYellow, _ := HexToRGB("#ffff00")
	green := color.RGB(rGreen, gGreen, bGreen)
	yellow := color.RGB(rYellow, gYellow, bYellow)

	switch level {
	case 20:
		if rarity == model.RarityType1 {
			return green.Sprintf("Lvl %d", level)
		}
	case 30:
		if rarity == model.RarityType2 {
			return green.Sprintf("Lvl %d", level)
		}
	case 40:
		if rarity == model.RarityType3 {
			return yellow.Sprintf("Lvl %d", level)
		}
	case 50:
		if rarity == model.RarityType3 {
			return green.Sprintf("Lvl %d", level)
		} else if rarity == model.RarityType4 {
			return yellow.Sprintf("Lvl %d", level)
		}
	case 60:
		if rarity == model.RarityType4 || rarity == model.RarityTypeBirthday {
			return green.Sprintf("Lvl %d", level)
		}
	}

	return fmt.Sprintf("Lvl %d", level)
}

// FormatAttribute applies color coding to card attributes using the game's
// official colors. Each attribute has its own distinct color for easy visual
// identification:
//
//   - cool:       blue   (#2545ec)
//   - cute:       pink   (#FF65AA)
//   - happy:      orange (#fe8100)
//   - pure:       green  (#009632)
//   - mysterious: purple (#713fc1)
func FormatAttribute(attr string) string {
	switch attr {
	case "cool":
		r, g, b, _ := HexToRGB("#2545ec")
		return color.RGB(r, g, b).Sprint("Cool")
	case "cute":
		r, g, b, _ := HexToRGB("#FF65AA")
		return color.RGB(r, g, b).Sprint("Cute")
	case "happy":
		r, g, b, _ := HexToRGB("#fe8100")
		return color.RGB(r, g, b).Sprint("Happy")
	case "pure":
		r, g, b, _ := HexToRGB("#009632")
		return color.RGB(r, g, b).Sprint("Pure")
	case "mysterious":
		r, g, b, _ := HexToRGB("#713fc1")
		return color.RGB(r, g, b).Sprint("Myst")
	default:
		return attr
	}
}

// FormatUnit converts internal unit names to their official abbreviations.
//
// Project Sekai units and their abbreviations:
//
//   - light_sound:    L/N  (Leo/need)
//   - idol:           MMJ  (MORE MORE JUMP!)
//   - street:         VBS  (Vivid BAD SQUAD)
//   - theme_park:     WxS  (Wonderlands×Showtime)
//   - school_refusal: N25  (Nightcord at 25:00)
//   - piapro:         VS   (Virtual Singer)
func FormatUnit(supportUnit string) string {
	switch supportUnit {
	case "light_sound":
		return "L/N"
	case "idol":
		return "MMJ"
	case "street":
		return "VBS"
	case "theme_park":
		return "WxS"
	case "school_refusal":
		return "N25"
	case "piapro":
		return "VS"
	default:
		return ""
	}
}

// ParseCardID parses a card ID from a string and validates that it is a
// positive integer. It returns an error with a human-readable message if
// parsing fails or the value is not positive.
func ParseCardID(arg string) (int, error) {
	cardID, err := strconv.Atoi(arg)
	if err != nil {
		return 0, fmt.Errorf("invalid cardID '%s'. cardID must be a number", arg)
	}
	if cardID <= 0 {
		return 0, fmt.Errorf("cardID must be a positive number")
	}
	return cardID, nil
}

// ContainsIgnoreCase reports whether substr is contained within str,
// ignoring case. Useful for case-insensitive name matching.
func ContainsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// HexToRGB converts a hexadecimal color code into its red, green, and blue
// components in the range [0, 255]. The input may include a leading "#" prefix.
func HexToRGB(hex string) (int, int, int, error) {
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}
	r, err := strconv.ParseInt(hex[0:2], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.ParseInt(hex[2:4], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.ParseInt(hex[4:6], 16, 0)
	if err != nil {
		return 0, 0, 0, err
	}
	return int(r), int(g), int(b), nil
}
