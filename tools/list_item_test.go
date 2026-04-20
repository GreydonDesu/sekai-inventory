package tools

import (
	"os"
	"sekai-inventory/model"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func TestMain(m *testing.M) {
	color.NoColor = true
	os.Exit(m.Run())
}

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex     string
		r, g, b int
		wantErr bool
	}{
		{"#ff0000", 255, 0, 0, false},
		{"#00ff00", 0, 255, 0, false},
		{"#0000ff", 0, 0, 255, false},
		{"ff0000", 255, 0, 0, false},
		{"#000000", 0, 0, 0, false},
		{"#ffffff", 255, 255, 255, false},
		{"invalid!", 0, 0, 0, true},
	}
	for _, tt := range tests {
		r, g, b, err := HexToRGB(tt.hex)
		if (err != nil) != tt.wantErr {
			t.Errorf("HexToRGB(%q) error = %v, wantErr %v", tt.hex, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && (r != tt.r || g != tt.g || b != tt.b) {
			t.Errorf("HexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)", tt.hex, r, g, b, tt.r, tt.g, tt.b)
		}
	}
}

func TestFormatBool(t *testing.T) {
	if got := FormatBool(true); got != "☑" {
		t.Errorf("FormatBool(true) = %q, want %q", got, "☑")
	}
	if got := FormatBool(false); got != "☐" {
		t.Errorf("FormatBool(false) = %q, want %q", got, "☐")
	}
}

func TestFormatCharacterName(t *testing.T) {
	tests := []struct {
		c    model.Character
		want string
	}{
		{model.Character{FirstName: "Kiritani", GivenName: "Haruka"}, "Kiritani Haruka"},
		{model.Character{FirstName: "", GivenName: "Miku"}, "Miku"},
	}
	for _, tt := range tests {
		if got := FormatCharacterName(tt.c); got != tt.want {
			t.Errorf("FormatCharacterName(%v) = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestFormatRarity(t *testing.T) {
	tests := []struct {
		rarity string
		want   string
	}{
		{model.RarityType1, "★"},
		{model.RarityType2, "★★"},
		{model.RarityType3, "★★★"},
		{model.RarityType4, "★★★★"},
		{model.RarityTypeBirthday, "୨୧"},
		{"unknown", "unknown"},
	}
	for _, tt := range tests {
		if got := FormatRarity(tt.rarity); got != tt.want {
			t.Errorf("FormatRarity(%q) = %q, want %q", tt.rarity, got, tt.want)
		}
	}
}

func TestFormatUnit(t *testing.T) {
	tests := []struct {
		unit string
		want string
	}{
		{"light_sound", "L/N"},
		{"idol", "MMJ"},
		{"street", "VBS"},
		{"theme_park", "WxS"},
		{"school_refusal", "N25"},
		{"piapro", "VS"},
		{"unknown", ""},
	}
	for _, tt := range tests {
		if got := FormatUnit(tt.unit); got != tt.want {
			t.Errorf("FormatUnit(%q) = %q, want %q", tt.unit, got, tt.want)
		}
	}
}

func TestFormatAttribute(t *testing.T) {
	tests := []struct {
		attr string
		want string
	}{
		{"cool", "Cool"},
		{"cute", "Cute"},
		{"happy", "Happy"},
		{"pure", "Pure"},
		{"mysterious", "Myst"},
		{"other", "other"},
	}
	for _, tt := range tests {
		if got := FormatAttribute(tt.attr); got != tt.want {
			t.Errorf("FormatAttribute(%q) = %q, want %q", tt.attr, got, tt.want)
		}
	}
}

func TestFormatLevel(t *testing.T) {
	tests := []struct {
		rarity string
		level  int
		want   string
	}{
		{model.RarityType1, 0, "Lvl 0"},
		{model.RarityType1, 20, "Lvl 20"},
		{model.RarityType2, 30, "Lvl 30"},
		{model.RarityType3, 40, "Lvl 40"},
		{model.RarityType3, 50, "Lvl 50"},
		{model.RarityType4, 50, "Lvl 50"},
		{model.RarityType4, 60, "Lvl 60"},
		{model.RarityTypeBirthday, 60, "Lvl 60"},
		{model.RarityType1, 10, "Lvl 10"},
	}
	for _, tt := range tests {
		if got := FormatLevel(tt.rarity, tt.level); got != tt.want {
			t.Errorf("FormatLevel(%q, %d) = %q, want %q", tt.rarity, tt.level, got, tt.want)
		}
	}
}

func TestParseCardID(t *testing.T) {
	tests := []struct {
		arg     string
		want    int
		wantErr bool
	}{
		{"123", 123, false},
		{"1", 1, false},
		{"0", 0, true},
		{"-1", 0, true},
		{"abc", 0, true},
	}
	for _, tt := range tests {
		got, err := ParseCardID(tt.arg)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseCardID(%q) error = %v, wantErr %v", tt.arg, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("ParseCardID(%q) = %d, want %d", tt.arg, got, tt.want)
		}
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		str, substr string
		want        bool
	}{
		{"Kiritani Haruka", "haruka", true},
		{"Kiritani Haruka", "Haruka", true},
		{"Kiritani Haruka", "HARUKA", true},
		{"Kiritani Haruka", "xyz", false},
		{"", "x", false},
		{"abc", "", true},
	}
	for _, tt := range tests {
		if got := ContainsIgnoreCase(tt.str, tt.substr); got != tt.want {
			t.Errorf("ContainsIgnoreCase(%q, %q) = %v, want %v", tt.str, tt.substr, got, tt.want)
		}
	}
}

func TestFormatCardLabel(t *testing.T) {
	card := model.CardEntity{
		Card: model.Card{
			ID:             1,
			CharacterID:    1,
			CardRarityType: model.RarityType4,
			SupportUnit:    "idol",
			Prefix:         "Test Card",
		},
		Level: 1,
	}
	characterMap := map[int]model.Character{
		1: {ID: 1, FirstName: "Test", GivenName: "User"},
	}
	label := FormatCardLabel(card, characterMap)
	for _, substr := range []string{"[1]", "★★★★", "Test User", "(MMJ)", "Test Card"} {
		if !strings.Contains(label, substr) {
			t.Errorf("FormatCardLabel() output %q missing %q", label, substr)
		}
	}
}

func TestFormatCardLabelUnknownCharacter(t *testing.T) {
	card := model.CardEntity{
		Card: model.Card{
			ID:             99,
			CardRarityType: model.RarityType1,
			Prefix:         "Solo",
		},
	}
	label := FormatCardLabel(card, nil)
	if !strings.Contains(label, "Unknown Character") {
		t.Errorf("FormatCardLabel() with nil map = %q, want 'Unknown Character'", label)
	}
}

func TestFormatCardDetails(t *testing.T) {
	card := model.CardEntity{
		Card: model.Card{
			ID:             5,
			CharacterID:    2,
			CardRarityType: model.RarityType3,
			Attr:           "cool",
			SupportUnit:    "light_sound",
			Prefix:         "Detail Card",
		},
		Level:      30,
		MasterRank: 3,
		SkillLevel: 2,
		SideStory1: true,
		SideStory2: false,
		Painting:   true,
	}
	characterMap := map[int]model.Character{
		2: {ID: 2, FirstName: "Hoshino", GivenName: "Ichika"},
	}
	details := FormatCardDetails(card, characterMap)
	for _, substr := range []string{"[5]", "★★★", "Lvl 30", "Cool", "MR3", "SL2", "☑", "☐", "Hoshino Ichika", "(L/N)", "Detail Card"} {
		if !strings.Contains(details, substr) {
			t.Errorf("FormatCardDetails() output %q missing %q", details, substr)
		}
	}
}
