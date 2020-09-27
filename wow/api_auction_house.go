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
	ID        int64           `json:"id"`
	Item      *Item           `json:"item"`
	Quantity  int64           `json:"quantity"`   // quantity != 1 then unit_price != 0, buyout = 0, bid = 0
	UnitPrice int64           `json:"unit_price"` // unit_price != 0 then quantity != 1
	Buyout    int64           `json:"buyout"`     // buyout != 0 then quantity = 1
	Bid       int64           `json:"bid"`        //bid != 0 then quantity = 1
	TimeLeft  AuctionTimeLeft `json:"time_left"`
}

type auctionApi interface {
	GetAuctions(connectedRealmID int64) ([]Auction, error)
}

func (wc *wowClientImpl) GetAuctions(connectedRealmID int64) ([]Auction, error) {
	auctionsResp := struct {
		Auctions []Auction `json:"auctions"`
	}{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/connected-realm/%d/auctions", connectedRealmID), blizzard.Dynamic, &auctionsResp)
	return auctionsResp.Auctions, err
}
