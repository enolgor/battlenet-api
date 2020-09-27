package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/blizzard"
)

type Item struct {
	ID   int64                    `json:"id"`
	Name *blizzard.LocalizedField `json:"name"`
}

type itemApi interface {
	GetItem(id int64) (*Item, error)
}

func (wc *wowClientImpl) GetItem(id int64) (*Item, error) {
	item := &Item{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/item/%d", id), blizzard.Static, item)
	return item, err
}
