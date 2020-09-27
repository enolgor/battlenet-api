package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/blizzard"
)

type Realm struct {
	ID   int64                    `json:"id"`
	Slug string                   `json:"slug"`
	Name *blizzard.LocalizedField `json:"name"`
}

type RealmExtended struct {
	Realm
	Region           *RealmRegion             `json:"region"`
	Category         *blizzard.LocalizedField `json:"category"`
	Locale           string                   `json:"locale"`
	Timezone         string                   `json:"timezone"`
	Type             *RealmType               `json:"type"`
	IsTournament     bool                     `json:"is_tournament"`
	ConnectedRealmID int64                    `json:"connected_realm_id"`
}

type RealmType struct {
	Type string                   `json:"type"`
	Name *blizzard.LocalizedField `json:"name"`
}

type RealmRegion struct {
	ID   int64                    `json:"id"`
	Name *blizzard.LocalizedField `json:"name"`
}

type realmApi interface {
	GetRealmsIndex() ([]Realm, error)
	GetRealm(realmSlug string) (*RealmExtended, error)
}

func (wc *wowClientImpl) GetRealmsIndex() ([]Realm, error) {
	respStruct := struct {
		Realms []Realm `json:"realms"`
	}{}
	if err := wc.getGameData("/data/wow/realm/index", blizzard.Dynamic, &respStruct); err != nil {
		return nil, err
	}
	return respStruct.Realms, nil
}

func (wc *wowClientImpl) GetRealm(realmSlug string) (*RealmExtended, error) {
	realmWithConnected := struct {
		RealmExtended
		ConnectedRealm struct {
			Href string `json:"href"`
		} `json:"connected_realm"`
	}{}
	err := wc.getGameData(fmt.Sprintf("/data/wow/realm/%s", realmSlug), blizzard.Dynamic, &realmWithConnected)
	realmWithConnected.RealmExtended.ConnectedRealmID = parseConnectedRealmID(realmWithConnected.ConnectedRealm.Href)
	return &realmWithConnected.RealmExtended, err
}
