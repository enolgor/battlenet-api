package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/blizzard"
)

type AuctionTimeLeft string

const (
	SHORT     AuctionTimeLeft = "SHORT"
	LONG      AuctionTimeLeft = "LONG"
	VERY_LONG AuctionTimeLeft = "VERY_LONG"
)

type Auction struct {
	ID          uint64          `json:"id"`
	AuctionItem *AuctionItem    `json:"item"`
	Quantity    uint64          `json:"quantity"`   // quantity != 1 then unit_price != 0, buyout = 0, bid = 0
	UnitPrice   uint64          `json:"unit_price"` // unit_price != 0 then quantity != 1
	Buyout      uint64          `json:"buyout"`     // buyout != 0 then quantity = 1
	Bid         uint64          `json:"bid"`        //bid != 0 then quantity = 1
	TimeLeft    AuctionTimeLeft `json:"time_left"`
}

type AuctionItem struct {
	ID           uint64                `json:"id"`
	Context      uint8                 `json:"context"`
	Modifiers    []AuctionItemModifier `json:"modifiers,omitempty"`
	PetBreedID   *uint16               `json:"pet_breed_id,omitempty"`
	PetLevel     *uint8                `json:"pet_level,omitempty"`
	PetQualityID *uint8                `json:"pet_quality_id,omitempty"`
	PetSpeciesID *uint64               `json:"pet_species_id,omitempty"`
	BonusLists   []uint64              `json:"bonus_lists"`
}

type AuctionItemModifier struct {
	Type  uint16 `json:"type"`
	Value uint64 `json:"value"`
}

type auctionApi interface {
	GetAuctions(connectedRealmID uint16) ([]Auction, error)
}

func (wc *wowClientImpl) GetAuctions(connectedRealmID uint16) ([]Auction, error) {
	auctionsResp := struct {
		Auctions []Auction `json:"auctions"`
	}{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/connected-realm/%d/auctions", connectedRealmID), blizzard.Dynamic, &auctionsResp)
	return auctionsResp.Auctions, err
}
