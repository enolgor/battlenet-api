package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/blizzard"
)

type InventoryType string

const (
	AMMO           InventoryType = "AMMO"
	BAG            InventoryType = "BAG"
	BODY           InventoryType = "BODY"
	CHEST          InventoryType = "CHEST"
	CLOAK          InventoryType = "CLOAK"
	FEET           InventoryType = "FEET"
	FINGER         InventoryType = "FINGER"
	HAND           InventoryType = "HAND"
	HEAD           InventoryType = "HEAD"
	HOLDABLE       InventoryType = "HOLDABLE"
	LEGS           InventoryType = "LEGS"
	NECK           InventoryType = "NECK"
	NON_EQUIP      InventoryType = "NON_EQUIP"
	RANGED         InventoryType = "RANGED"
	RANGEDRIGHT    InventoryType = "RANGEDRIGHT"
	ROBE           InventoryType = "ROBE"
	SHIELD         InventoryType = "SHIELD"
	SHOULDER       InventoryType = "SHOULDER"
	TABARD         InventoryType = "TABARD"
	THROWN         InventoryType = "THROWN"
	TRINKET        InventoryType = "TRINKET"
	TWOHWEAPON     InventoryType = "TWOHWEAPON"
	WAIST          InventoryType = "WAIST"
	WEAPON         InventoryType = "WEAPON"
	WEAPONMAINHAND InventoryType = "WEAPONMAINHAND"
	WEAPONOFFHAND  InventoryType = "WEAPONOFFHAND"
	WRIST          InventoryType = "WRIST"
)

type QualityType string

const (
	ARTIFACT  QualityType = "ARTIFACT"
	COMMON    QualityType = "COMMON"
	EPIC      QualityType = "EPIC"
	HEIRLOOM  QualityType = "HEIRLOOM"
	LEGENDARY QualityType = "LEGENDARY"
	POOR      QualityType = "POOR"
	RARE      QualityType = "RARE"
	UNCOMMON  QualityType = "UNCOMMON"
	WOWTOKEN  QualityType = "WOWTOKEN"
)

type ItemMedia struct {
	ID uint64 `json:"id"`
}

type Item struct {
	ID            uint64                   `json:"id"`
	Name          *blizzard.LocalizedField `json:"name"`
	IsEquippable  bool                     `json:"is_equippable"`
	IsStackable   bool                     `json:"is_stackable"`
	IsAzeriteItem *bool                    `json:"is_azerite_item,omitempty"` //just appears in the heart of azeroth item
	Level         uint32                   `json:"level"`
	RequiredLevel uint32                   `json:"required_level"`
	InventoryType *struct {
		Name *blizzard.LocalizedField `json:"name"`
		Type *InventoryType           `json:"type"`
	} `json:"inventory_type,omitempty"`
	Quality *struct {
		Name *blizzard.LocalizedField `json:"name"`
		Type *QualityType             `json:"type"`
	} `json:"quality,omitempty"`
	MaxCount         uint64        `json:"max_count"`
	PurchasePrice    uint64        `json:"purchase_price"`
	PurchaseQuantity uint64        `json:"purchase_quantity"`
	SellPrice        uint64        `json:"sell_price"`
	Media            *ItemMedia    `json:"media"`
	ItemClass        *ItemClass    `json:"item_class"`
	ItemSubClass     *ItemSubClass `json:"item_subclass"`
}

type ItemClass struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   uint64                   `json:"id"`
}

type ItemSubClass struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   uint64                   `json:"id"`
}

type AzeriteClassPower struct {
	PlayableClass *PlayableClass `json:"playable_class"`
}

type AzeritePower struct {
	Spell                  *Spell           `json:"spell"`
	Tier                   uint8            `json:"tier"`
	ID                     uint64           `json:"id"`
	AllowedSpecializations []Specialization `json:"allowed_specializations,omitempty"`
}

type itemApi interface {
	GetItem(id uint64) (*Item, error)
	SearchItem(query blizzard.SearchQuery) ([]Item, *blizzard.SearchResult, error)
}

func (wc *wowClientImpl) GetItem(id uint64) (*Item, error) {
	item := &Item{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/item/%d", id), blizzard.Static, item)
	return item, err
}

func (wc *wowClientImpl) SearchItem(query blizzard.SearchQuery) ([]Item, *blizzard.SearchResult, error) {
	items := []Item{}
	searchResult, err := wc.searchGameData("/data/wow/search/item", query, blizzard.Static, &items)
	if err != nil {
		return nil, nil, err
	}
	return items, searchResult, nil
}
