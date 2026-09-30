package function

import (
	"sekai-inventory/model"
	"testing"
)

func TestParseIntField(t *testing.T) {
	tests := []struct {
		value, name string
		min, max    int
		want        int
		wantErr     bool
	}{
		{"5", "level", 1, 60, 5, false},
		{"1", "level", 1, 60, 1, false},
		{"60", "level", 1, 60, 60, false},
		{"0", "level", 1, 60, 0, true},
		{"61", "level", 1, 60, 0, true},
		{"abc", "level", 1, 60, 0, true},
	}
	for _, tt := range tests {
		got, err := parseIntField(tt.value, tt.name, tt.min, tt.max)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseIntField(%q, %q, %d, %d) error = %v, wantErr %v", tt.value, tt.name, tt.min, tt.max, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("parseIntField(%q, ...) = %d, want %d", tt.value, got, tt.want)
		}
	}
}

func TestParseBoolField(t *testing.T) {
	tests := []struct {
		value   string
		want    bool
		wantErr bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"TRUE", true, false},
		{"1", true, false},
		{"0", false, false},
		{"yes", false, true},
		{"", false, true},
	}
	for _, tt := range tests {
		got, err := parseBoolField(tt.value, "field")
		if (err != nil) != tt.wantErr {
			t.Errorf("parseBoolField(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("parseBoolField(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestApplyCardField(t *testing.T) {
	tests := []struct {
		field   string
		value   string
		check   func(*model.CardEntity) bool
		wantErr bool
	}{
		{"level", "30", func(c *model.CardEntity) bool { return c.Level == 30 }, false},
		{"skillLevel", "4", func(c *model.CardEntity) bool { return c.SkillLevel == 4 }, false},
		{"masterRank", "3", func(c *model.CardEntity) bool { return c.MasterRank == 3 }, false},
		{"sideStory1", "true", func(c *model.CardEntity) bool { return c.SideStory1 }, false},
		{"sideStory2", "true", func(c *model.CardEntity) bool { return c.SideStory2 }, false},
		{"painting", "true", func(c *model.CardEntity) bool { return c.Painting }, false},
		{"skillLevel", "5", nil, true},
		{"level", "99", nil, true},
		{"level", "abc", nil, true},
		{"unknown", "x", nil, true},
	}
	for _, tt := range tests {
		card := &model.CardEntity{Card: model.Card{ID: 1}, Level: 1, SkillLevel: 1}
		err := applyCardField(card, tt.field, tt.value)
		if (err != nil) != tt.wantErr {
			t.Errorf("applyCardField(%q, %q) error = %v, wantErr %v", tt.field, tt.value, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && !tt.check(card) {
			t.Errorf("applyCardField(%q, %q) check failed", tt.field, tt.value)
		}
	}
}
