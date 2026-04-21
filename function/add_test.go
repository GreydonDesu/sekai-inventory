package function

import (
	"sekai-inventory/model"
	"testing"
)

func TestClassifyCardIDs(t *testing.T) {
	cardMap := map[int]model.Card{
		1: {ID: 1, CardRarityType: model.RarityType4},
		2: {ID: 2, CardRarityType: model.RarityType3},
	}
	inventory := &model.Inventory{
		Cards: []model.CardEntity{
			{Card: model.Card{ID: 1}},
		},
	}

	added, existing, missing := classifyCardIDs([]int{1, 2, 99}, inventory, cardMap)

	if len(added) != 1 || added[0].ID != 2 {
		t.Errorf("added = %v, want 1 card with ID 2", added)
	}
	if len(existing) != 1 || existing[0].ID != 1 {
		t.Errorf("existing = %v, want 1 card with ID 1", existing)
	}
	if len(missing) != 1 || missing[0] != 99 {
		t.Errorf("missing = %v, want [99]", missing)
	}
}

func TestClassifyCardIDsDefaults(t *testing.T) {
	cardMap := map[int]model.Card{
		5: {ID: 5, CardRarityType: model.RarityType2},
	}
	inventory := &model.Inventory{}

	added, existing, missing := classifyCardIDs([]int{5}, inventory, cardMap)

	if len(added) != 1 {
		t.Fatalf("expected 1 added card, got %d", len(added))
	}
	if added[0].Level != 1 {
		t.Errorf("default Level = %d, want 1", added[0].Level)
	}
	if added[0].MasteryRank != 0 {
		t.Errorf("default MasteryRank = %d, want 0", added[0].MasteryRank)
	}
	if added[0].SkillLevel != 1 {
		t.Errorf("default SkillLevel = %d, want 1", added[0].SkillLevel)
	}
	if len(existing) != 0 || len(missing) != 0 {
		t.Errorf("unexpected existing=%v missing=%v", existing, missing)
	}
}

func TestClassifyCardIDsAllMissing(t *testing.T) {
	cardMap := map[int]model.Card{}
	inventory := &model.Inventory{}

	added, existing, missing := classifyCardIDs([]int{10, 20}, inventory, cardMap)

	if len(added) != 0 || len(existing) != 0 {
		t.Errorf("expected no added/existing, got added=%v existing=%v", added, existing)
	}
	if len(missing) != 2 {
		t.Errorf("missing len = %d, want 2", len(missing))
	}
}
