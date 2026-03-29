package tools

import (
	"fmt"
	"sekai-inventory/model"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// FormatCardDetails returns a human-readable, colorized string representation
// of a card. It combines card properties with character information and applies
// formatting based on rarity, level, side stories, and painting status.
//
// The characterMap is used to resolve CharacterID to character names and units.
func FormatCardDetails(card model.CardEntity, characterMap map[int]model.Character) string {
	// Get the character name.
	character, exists := characterMap[card.CharacterID]
	characterName := "Unknown Character"
	if exists {
		if character.FirstName == "" {
			characterName = character.GivenName
		} else {
			characterName = fmt.Sprintf("%s %s", character.FirstName, character.GivenName)
		}
	}

	// Map the support unit to its abbreviation.
	supportUnitAbbreviation := FormatUnit(card.SupportUnit)
	if supportUnitAbbreviation == "" && exists {
		// Fallback: use the Unit from the character data.
		supportUnitAbbreviation = FormatUnit(character.Unit)
	}

	var supportUnitCard string
	if supportUnitAbbreviation != "" {
		supportUnitCard = fmt.Sprintf(" (%s)", supportUnitAbbreviation)
	}

	// Format the rarity.
	formattedRarity := FormatRarity(card.CardRarityType)

	// Format the level.
	formattedLevel := FormatLevel(card.CardRarityType, card.Level)

	// Colors.
	greenHex := "#00ff00"
	redHex := "#ff0000"

	// Side Story 1.
	var sideStory1 string
	if card.SideStory1 {
		r, g, b, _ := HexToRGB(greenHex)
		sideStory1 = color.RGB(r, g, b).Sprint("☑")
	} else {
		r, g, b, _ := HexToRGB(redHex)
		sideStory1 = color.RGB(r, g, b).Sprint("☐")
	}

	// Side Story 2.
	var sideStory2 string
	if card.SideStory2 {
		r, g, b, _ := HexToRGB(greenHex)
		sideStory2 = color.RGB(r, g, b).Sprint("☑")
	} else {
		r, g, b, _ := HexToRGB(redHex)
		sideStory2 = color.RGB(r, g, b).Sprint("☐")
	}

	// Painting.
	var painting string
	if card.Painting {
		r, g, b, _ := HexToRGB(greenHex)
		painting = color.RGB(r, g, b).Sprint("☑")
	} else {
		r, g, b, _ := HexToRGB(redHex)
		painting = color.RGB(r, g, b).Sprint("☐")
	}

	// Highlight Master Rank and Skill Level with RGB green at max values.
	rGreen, gGreen, bGreen, _ := HexToRGB("#00ff00")

	masterRank := fmt.Sprintf("MR%d", card.MasterRank)
	if card.MasterRank == 5 {
		masterRank = color.RGB(rGreen, gGreen, bGreen).Sprint(masterRank)
	}

	skillLevel := fmt.Sprintf("SL%d", card.SkillLevel)
	if card.SkillLevel == 4 {
		skillLevel = color.RGB(rGreen, gGreen, bGreen).Sprint(skillLevel)
	}

	// Format the card details as a one-liner.
	return fmt.Sprintf("[%d]\t%s\t%s\t%s\t| %s | %s | Side Story 1: %s | Side Story 2: %s | Painting: %s | %s%s \"%s\"",
		card.ID,
		formattedRarity,
		formattedLevel,
		FormatAttribute(card.Attr),
		masterRank,
		skillLevel,
		sideStory1,
		sideStory2,
		painting,
		characterName,
		supportUnitCard,
		card.Prefix,
	)
}

// FormatRarity converts a raw rarity value into a visually appealing string.
// Each rarity level is represented by colored stars, with birthday cards using
// a special symbol.
//
// Rarity levels:
//
//   - rarity_1:        ★ (yellow)
//   - rarity_2:       ★★ (yellow)
//   - rarity_3:      ★★★ (yellow)
//   - rarity_4:     ★★★★ (yellow)
//   - rarity_birthday: ୨୧ (magenta)
func FormatRarity(rarity string) string {
	yellowHex := "#ffd700"
	magentaHex := "#ff00ff"

	switch rarity {
	case "rarity_1":
		r, g, b, _ := HexToRGB(yellowHex)
		return color.RGB(r, g, b).Sprint("★")
	case "rarity_2":
		r, g, b, _ := HexToRGB(yellowHex)
		return color.RGB(r, g, b).Sprint("★★")
	case "rarity_3":
		r, g, b, _ := HexToRGB(yellowHex)
		return color.RGB(r, g, b).Sprint("★★★")
	case "rarity_4":
		r, g, b, _ := HexToRGB(yellowHex)
		return color.RGB(r, g, b).Sprint("★★★★")
	case "rarity_birthday":
		r, g, b, _ := HexToRGB(magentaHex)
		return color.RGB(r, g, b).Sprint("୨୧")
	default:
		return rarity
	}
}

// FormatLevel formats a card's level with color coding based on rarity-specific
// maximum level thresholds. The color indicates the card's progression:
//
//   - Green:  at maximum level for the card's rarity.
//   - Yellow: approaching maximum level.
//   - Plain:  below maximum level.
//
// Maximum levels by rarity:
//
//   - rarity_1:        20
//   - rarity_2:        30
//   - rarity_3:        50
//   - rarity_4:        60
//   - rarity_birthday: 60
func FormatLevel(rarity string, level int) string {
	if level <= 0 {
		return "Lvl 0"
	}

	greenHex := "#00ff00"
	yellowHex := "#ffff00"

	switch level {
	case 20:
		if rarity == "rarity_1" {
			r, g, b, _ := HexToRGB(greenHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		}
	case 30:
		if rarity == "rarity_2" {
			r, g, b, _ := HexToRGB(greenHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		}
	case 40:
		if rarity == "rarity_3" {
			r, g, b, _ := HexToRGB(yellowHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		}
	case 50:
		if rarity == "rarity_3" {
			r, g, b, _ := HexToRGB(greenHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		} else if rarity == "rarity_4" {
			r, g, b, _ := HexToRGB(yellowHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		}
	case 60:
		if rarity == "rarity_4" || rarity == "rarity_birthday" {
			r, g, b, _ := HexToRGB(greenHex)
			return color.RGB(r, g, b).Sprintf("Lvl %d", level)
		}
	}

	return fmt.Sprintf("Lvl %d", level)
}

// FormatAttribute applies color coding to card attributes using the game's
// official colors. Each attribute has its own distinct color for easy visual
// identification:
//
//   - cool:       blue  (#2545ec)
//   - cute:       pink  (#FF65AA)
//   - happy:      orange (#fe8100)
//   - pure:       green (#009632)
//   - mysterious: purple (#713fc1)
func FormatAttribute(attr string) string {
	switch attr {
	case "cool":
		r, g, b, _ := HexToRGB("#2545ec") // Blue
		return color.RGB(r, g, b).Sprint("Cool")
	case "cute":
		r, g, b, _ := HexToRGB("#FF65AA") // Pink
		return color.RGB(r, g, b).Sprint("Cute")
	case "happy":
		r, g, b, _ := HexToRGB("#fe8100") // Orange
		return color.RGB(r, g, b).Sprint("Happy")
	case "pure":
		r, g, b, _ := HexToRGB("#009632") // Green
		return color.RGB(r, g, b).Sprint("Pure")
	case "mysterious":
		r, g, b, _ := HexToRGB("#713fc1") // Purple
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
		// Return an empty string if the unit is not recognized.
		return ""
	}
}

// ParseCardID parses a card ID from a string and validates that it is a positive
// integer. It returns an error with a human-readable message if parsing fails.
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

// ContainsIgnoreCase reports whether substr is contained within str, ignoring
// case. It is useful for case-insensitive matching of names and other text.
func ContainsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// HexToRGB converts a hexadecimal color code into its red, green, and blue
// components. The input may include a leading "#" prefix.
//
// It returns the RGB components in the range [0, 255] and an error if the
// hex string is invalid.
func HexToRGB(hex string) (int, int, int, error) {
	// Remove the "#" if it exists.
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	// Parse the red, green, and blue components.
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
