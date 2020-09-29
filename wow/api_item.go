package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/battlenet"
	"github.com/enolgor/battlenet-api/blizzard"
)

type Item struct {
	ID   int64                    `json:"id"`
	Name *blizzard.LocalizedField `json:"name"`
}

type itemApi interface {
	GetItem(id int64) (*Item, error)
	SearchItem(query battlenet.SearchQuery) ([]Item, *battlenet.SearchResult, error)
}

func (wc *wowClientImpl) GetItem(id int64) (*Item, error) {
	item := &Item{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/item/%d", id), blizzard.Static, item)
	return item, err
}

func (wc *wowClientImpl) SearchItem(query battlenet.SearchQuery) ([]Item, *battlenet.SearchResult, error) {
	items := []Item{}
	searchResult, err := wc.searchGameData("/data/wow/search/item", query, blizzard.Static, &items)
	if err != nil {
		return nil, nil, err
	}
	return items, searchResult, nil
}
