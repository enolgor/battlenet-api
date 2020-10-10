package battlenet

import (
	"net/url"

	"github.com/enolgor/battlenet-api/blizzard"
)

var usEndpoint, euEndpoint, krEndpoint, twEndpoint, cnEndpoint *url.URL

func init() {
	mustParse := func(s string) (url *url.URL) {
		var err error
		if url, err = url.Parse(s); err != nil {
			panic(err)
		}
		return
	}

	usEndpoint = mustParse("https://us.battle.net")
	euEndpoint = mustParse("https://eu.battle.net")
	krEndpoint = mustParse("https://kr.battle.net")
	twEndpoint = mustParse("https://tw.battle.net")
	cnEndpoint = mustParse("https://www.battlenet.com.cn")
}

// from https://golang.org/src/net/http/clone.go?m=text
func cloneURL(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	u2 := new(url.URL)
	*u2 = *u
	if u.User != nil {
		u2.User = new(url.Userinfo)
		*u2.User = *u.User
	}
	return u2
}

const OAuthScopeWoWProfile string = "wow.profile"
const OAuthScopeD3Profile string = "d3.profile"
const OAuthScopeSC2Profile string = "sc2.profile"
const OAuthScopeOpenID string = "openid"

func NewEndpoint(region blizzard.Region, path string) *url.URL {
	var endpoint *url.URL
	switch region {
	case blizzard.US:
		endpoint = cloneURL(usEndpoint)
	case blizzard.EU:
		endpoint = cloneURL(euEndpoint)
	case blizzard.KR:
		endpoint = cloneURL(krEndpoint)
	case blizzard.TW:
		endpoint = cloneURL(twEndpoint)
	case blizzard.CN:
		endpoint = cloneURL(cnEndpoint)
	default:
		panic("Missing region in switch")
	}
	if path != "" {
		endpoint.Path = path
	}
	return endpoint
}
