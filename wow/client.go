package wow

import (
	"fmt"

	"github.com/enolgor/battlenet-api/battlenet"
	"github.com/enolgor/battlenet-api/blizzard"
)

type WoWClient interface {
	DynamicAPI
	StaticAPI
	ProfileAPI
	SetLocale(locale blizzard.Locale) error
	Dynamic() DynamicAPI
	Static() StaticAPI
	Profile() ProfileAPI
}

type context struct {
	locale blizzard.Locale
}

type wowClientImpl struct {
	battleNetClient battlenet.BattleNetClient
	namespace       blizzard.Namespace
	context         *context
	static          *wowClientImpl
	dynamic         *wowClientImpl
	profile         *wowClientImpl
}

//var currentLocale blizzard.Locale = blizzard.NoLocale

func NewWoWClient(battlenetClient battlenet.BattleNetClient) WoWClient {
	context := &context{locale: blizzard.NoLocale}
	defaultWowClient := &wowClientImpl{battleNetClient: battlenetClient, context: context, namespace: blizzard.NoNamespace}
	staticClient := *defaultWowClient
	dynamicClient := *defaultWowClient
	profileClient := *defaultWowClient
	staticClient.namespace = blizzard.Static
	dynamicClient.namespace = blizzard.Dynamic
	profileClient.namespace = blizzard.Profile
	defaultWowClient.static = &staticClient
	defaultWowClient.dynamic = &dynamicClient
	defaultWowClient.profile = &profileClient
	return defaultWowClient
}

func (wc *wowClientImpl) SetLocale(locale blizzard.Locale) error {
	if locale == blizzard.PtPT {
		return fmt.Errorf("Unsupported locale")
	}
	wc.context.locale = locale
	return nil
}

func (wc *wowClientImpl) getGameData(path string, defaultNamespace blizzard.Namespace, receiver interface{}) error {
	var requestNamespace blizzard.Namespace
	if wc.namespace == blizzard.NoNamespace {
		requestNamespace = defaultNamespace
	} else {
		requestNamespace = wc.namespace
	}
	return wc.battleNetClient.GetGameData(path, wc.context.locale, requestNamespace, receiver)
}

func (wc *wowClientImpl) Dynamic() DynamicAPI {
	return wc.dynamic
}

func (wc *wowClientImpl) Static() StaticAPI {
	return wc.static
}

func (wc *wowClientImpl) Profile() ProfileAPI {
	return wc.profile
}
