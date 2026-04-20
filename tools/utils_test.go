package tools

import (
	"sekai-inventory/model"
	"testing"
)

func TestParseFilters(t *testing.T) {
	tests := []struct {
		args    []string
		want    map[string]string
		wantNil bool
	}{
		{
			[]string{"--character", "Miku"},
			map[string]string{"character": "Miku"},
			false,
		},
		{
			[]string{"--rarity", "4", "--group", "MMJ"},
			map[string]string{"rarity": "4", "group": "MMJ"},
			false,
		},
		{[]string{}, map[string]string{}, false},
		{[]string{"--character"}, nil, true},
		{[]string{"--rarity", "4", "--group"}, nil, true},
	}
	for _, tt := range tests {
		got := ParseFilters(tt.args)
		if tt.wantNil {
			if got != nil {
				t.Errorf("ParseFilters(%v) = %v, want nil", tt.args, got)
			}
			continue
		}
		if got == nil {
			t.Errorf("ParseFilters(%v) = nil, want %v", tt.args, tt.want)
			continue
		}
		for k, v := range tt.want {
			if got[k] != v {
				t.Errorf("ParseFilters(%v)[%q] = %q, want %q", tt.args, k, got[k], v)
			}
		}
	}
}

func TestCreateCharacterMap(t *testing.T) {
	chars := []model.Character{
		{ID: 1, FirstName: "A", GivenName: "B"},
		{ID: 2, FirstName: "C", GivenName: "D"},
	}
	m := CreateCharacterMap(chars)
	if len(m) != 2 {
		t.Fatalf("CreateCharacterMap() len = %d, want 2", len(m))
	}
	if m[1].GivenName != "B" {
		t.Errorf("CreateCharacterMap()[1].GivenName = %q, want %q", m[1].GivenName, "B")
	}
	if m[2].FirstName != "C" {
		t.Errorf("CreateCharacterMap()[2].FirstName = %q, want %q", m[2].FirstName, "C")
	}
}

func TestCreateCharacterMapEmpty(t *testing.T) {
	m := CreateCharacterMap(nil)
	if len(m) != 0 {
		t.Errorf("CreateCharacterMap(nil) len = %d, want 0", len(m))
	}
}

func TestRarityToKey(t *testing.T) {
	tests := []struct {
		k    string
		want string
	}{
		{"1", model.RarityType1},
		{"2", model.RarityType2},
		{"3", model.RarityType3},
		{"4", model.RarityType4},
		{"bd", model.RarityTypeBirthday},
	}
	for _, tt := range tests {
		if got := RarityToKey[tt.k]; got != tt.want {
			t.Errorf("RarityToKey[%q] = %q, want %q", tt.k, got, tt.want)
		}
	}
}

func TestGroupToKey(t *testing.T) {
	tests := []struct {
		k    string
		want string
	}{
		{"L/N", "light_sound"},
		{"MMJ", "idol"},
		{"VBS", "street"},
		{"WxS", "theme_park"},
		{"N25", "school_refusal"},
		{"VS", "piapro"},
	}
	for _, tt := range tests {
		if got := GroupToKey[tt.k]; got != tt.want {
			t.Errorf("GroupToKey[%q] = %q, want %q", tt.k, got, tt.want)
		}
	}
}
