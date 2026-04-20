package function

import (
	"sekai-inventory/model"
	"sekai-inventory/tools"
	"testing"
)

func TestMatchesFilters(t *testing.T) {
	characterMap := map[int]model.Character{
		1: {ID: 1, FirstName: "Kiritani", GivenName: "Haruka", Unit: "idol"},
	}
	card := model.CardEntity{
		Card: model.Card{
			ID:             1,
			CharacterID:    1,
			CardRarityType: model.RarityType4,
			SupportUnit:    "idol",
		},
		Painting: true,
	}

	tests := []struct {
		filters map[string]string
		want    bool
	}{
		{map[string]string{}, true},
		{map[string]string{"character": "Haruka"}, true},
		{map[string]string{"character": "haruka"}, true},
		{map[string]string{"character": "Miku"}, false},
		{map[string]string{"rarity": "4"}, true},
		{map[string]string{"rarity": "3"}, false},
		{map[string]string{"group": "MMJ"}, true},
		{map[string]string{"group": "VBS"}, false},
		{map[string]string{"painting": "true"}, true},
		{map[string]string{"painting": "false"}, false},
		{map[string]string{"painting": "invalid"}, false},
		{map[string]string{"unknown": "x"}, false},
	}
	for _, tt := range tests {
		got := matchesFilters(card, tt.filters, characterMap)
		if got != tt.want {
			t.Errorf("matchesFilters(filters=%v) = %v, want %v", tt.filters, got, tt.want)
		}
	}
}

func TestMatchesFiltersGroupFallback(t *testing.T) {
	characterMap := map[int]model.Character{
		2: {ID: 2, GivenName: "Miku", Unit: "piapro"},
	}
	card := model.CardEntity{
		Card: model.Card{
			ID:          2,
			CharacterID: 2,
			SupportUnit: "",
		},
	}
	if !matchesFilters(card, map[string]string{"group": "VS"}, characterMap) {
		t.Error("matchesFilters() group should fall back to character.Unit")
	}
}

func TestMatchesFiltersRarityConsistency(t *testing.T) {
	for short, full := range tools.RarityToKey {
		card := model.CardEntity{
			Card: model.Card{CardRarityType: full},
		}
		if !matchesFilters(card, map[string]string{"rarity": short}, nil) {
			t.Errorf("rarity filter %q should match card with rarity %q", short, full)
		}
	}
}
