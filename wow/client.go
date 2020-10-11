package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/battlenet"
	"github.com/enolgor/battlenet-api/blizzard"
)

var AvailableLocales []blizzard.Locale = []blizzard.Locale{
	blizzard.EnUS,
	blizzard.EsMX,
	blizzard.PtBR,
	blizzard.EnGB,
	blizzard.EsES,
	blizzard.FrFR,
	blizzard.RuRU,
	blizzard.DeDE,
	blizzard.ItIT,
	blizzard.KoKR,
	blizzard.ZhTW,
	blizzard.ZhCN,
}

type WoWClient interface {
	DynamicAPI
	StaticAPI
	ProfileAPI
}

type wowClientImpl struct {
	battleNetClient battlenet.BattleNetClient
	locale          blizzard.Locale
}

func NewWoWClient(battlenetClient battlenet.BattleNetClient, locale blizzard.Locale) (WoWClient, error) {
	if locale == blizzard.PtPT {
		return nil, fmt.Errorf("Unsupported locale")
	}
	defaultWowClient := &wowClientImpl{battleNetClient: battlenetClient}
	return defaultWowClient, nil
}

func (wc *wowClientImpl) getGameData(path string, namespace blizzard.Namespace, receiver interface{}) error {
	token, err := wc.battleNetClient.GetAccessToken()
	if err != nil {
		return err
	}
	return wc.battleNetClient.BlizzardAPIGet(path, wc.locale, namespace, token, receiver)
}

func (wc *wowClientImpl) searchGameData(path string, query blizzard.SearchQuery, namespace blizzard.Namespace, receiver interface{}) (*blizzard.SearchResult, error) {
	token, err := wc.battleNetClient.GetAccessToken()
	if err != nil {
		return nil, err
	}
	return wc.battleNetClient.BlizzardAPISearch(path, query, wc.locale, namespace, token, receiver)
}

func (wc *wowClientImpl) getProfileData(path string, token string, receiver interface{}) error {
	return wc.battleNetClient.BlizzardAPIGet(path, wc.locale, blizzard.Profile, token, receiver)
}
