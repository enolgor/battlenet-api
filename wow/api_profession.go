package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/blizzard"
)

type Profession struct {
	ID   uint64                   `json:"id"`
	Name *blizzard.LocalizedField `json:"name"`
}

type ProfessionType string

const (
	PRIMARY   ProfessionType = "PRIMARY"
	SECONDARY ProfessionType = "SECONDARY"
)

type ProfessionExtended struct {
	Profession
	Description *blizzard.LocalizedField `json:"description"`
	Type        struct {
		Type ProfessionType           `json:"type"`
		Name *blizzard.LocalizedField `json:"name"`
	} `json:"type"`
	SkillTiers []SkillTier `json:"skill_tiers"`
}

type SkillTier struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   uint64                   `json:"id"`
}

type SkillTierExtended struct {
	SkillTier
	MinimumSkillLevel uint16              `json:"minimum_skill_level"`
	MaximumSkillLevel uint16              `json:"maximum_skill_level"`
	Categories        []SkillTierCategory `json:"categories"`
}

type SkillTierCategory struct {
	Name    *blizzard.LocalizedField `json:"name"`
	Recipes []Recipe                 `json:"recipes"`
}

type Recipe struct {
	Name *blizzard.LocalizedField
	ID   uint64 `json:"id"`
}

type Reagent struct {
	ReagentItem *Item  `json:"reagent"`
	Quantity    uint16 `json:"quantity"`
}

type RecipeExtended struct {
	Recipe
	CraftedItem     *Item `json:"crafted_item"`
	CraftedQuantity struct {
		Value float32 `json:"value"`
	} `json:"crafted_quantity"`
	Reagents []Reagent `json:"reagents"`
}

type professionApi interface {
	GetProfessionIndex() ([]Profession, error)
	GetProfession(id uint64) (*ProfessionExtended, error)
	GetSkillTier(professionId uint64, skillTierId uint64) (*SkillTierExtended, error)
	GetRecipe(recipeId uint64) (*RecipeExtended, error)
}

func (wc *wowClientImpl) GetProfessionIndex() ([]Profession, error) {
	respStruct := struct {
		Professions []Profession `json:"professions"`
	}{}
	if err := wc.getGameData("/data/wow/profession/index", blizzard.Static, &respStruct); err != nil {
		return nil, err
	}
	return respStruct.Professions, nil
}

func (wc *wowClientImpl) GetProfession(id uint64) (*ProfessionExtended, error) {
	profession := &ProfessionExtended{}
	if err := wc.getGameData(fmt.Sprintf("/data/wow/profession/%d", id), blizzard.Static, profession); err != nil {
		return nil, err
	}
	return profession, nil
}

func (wc *wowClientImpl) GetSkillTier(professionId uint64, skillTierId uint64) (*SkillTierExtended, error) {
	skillTier := &SkillTierExtended{}
	if err := wc.getGameData(fmt.Sprintf("/data/wow/profession/%d/skill-tier/%d", professionId, skillTierId), blizzard.Static, skillTier); err != nil {
		return nil, err
	}
	return skillTier, nil
}

func (wc *wowClientImpl) GetRecipe(id uint64) (*RecipeExtended, error) {
	recipe := &RecipeExtended{}
	if err := wc.getGameData(fmt.Sprintf("/data/wow/recipe/%d", id), blizzard.Static, recipe); err != nil {
		return nil, err
	}
	return recipe, nil
}
