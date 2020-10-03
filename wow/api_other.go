package wow

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/enolgor/battlenet-api/blizzard"
)

type otherStatic interface {
	GetRetailVersion() (string, error)
}

func (wc *wowClientImpl) GetRetailVersion() (string, error) {
	respStruct := struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}{}
	if err := wc.getGameData("/data/wow/playable-class/index", blizzard.Static, &respStruct); err != nil {
		return "", err
	}
	href, _ := url.Parse(respStruct.Links.Self.Href)
	namespace := href.Query().Get("namespace")
	if namespace == "" {
		return "", fmt.Errorf("Namespace not found")
	}
	parts := strings.Split(namespace, "-")
	return parts[1], nil
}
