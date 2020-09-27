package wow

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/enolgor/battlenet-api/blizzard"
)

type ConnectedRealmStatus struct {
	Type string                   `json:"type"`
	Name *blizzard.LocalizedField `json:"name"`
}

type Population struct {
	Type string                   `json:"type"`
	Name *blizzard.LocalizedField `json:"name"`
}

type ConnectedRealm struct {
	ID         int64                 `json:"id"`
	HasQueue   bool                  `json:"has_queue"`
	Status     *ConnectedRealmStatus `json:"status"`
	Population *Population           `json:"population"`
	Realms     []RealmExtended       `json:"realms"`
}

type connectedRealmApi interface {
	GetConnectedRealmsIndex() ([]int64, error)
	GetConnectedRealm(id int64) (*ConnectedRealm, error)
}

func (wc *wowClientImpl) GetConnectedRealmsIndex() ([]int64, error) {
	respStruct := struct {
		ConnectedRealms []struct {
			Href string `json:"href"`
		} `json:"connected_realms"`
	}{}
	var err error
	if err = wc.getGameData("/data/wow/connected-realm/index", blizzard.Dynamic, &respStruct); err != nil {
		return nil, err
	}
	indices := make([]int64, 0, len(respStruct.ConnectedRealms))
	for _, cr := range respStruct.ConnectedRealms {
		indices = append(indices, parseConnectedRealmID(cr.Href))
	}
	return indices, nil
}

func (wc *wowClientImpl) GetConnectedRealm(id int64) (*ConnectedRealm, error) {
	connectedRealm := &ConnectedRealm{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/connected-realm/%d", id), blizzard.Dynamic, connectedRealm)
	for i, realm := range connectedRealm.Realms {
		realm.ConnectedRealmID = id
		connectedRealm.Realms[i] = realm
	}
	return connectedRealm, err
}

func parseConnectedRealmID(href string) int64 {
	hrefURL, _ := url.Parse(href)
	path := hrefURL.Path // /data/wow/connected-realm/531
	id, _ := strconv.ParseInt(strings.Split(path, "/")[4], 10, 64)
	return id
}
