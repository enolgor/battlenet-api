package battlenet

import (
	"net/url"

	"github.com/enolgor/battlenet-api/blizzard"
)

var oauthEndpointUS, oauthEndpointEU, oauthEndpointKR, oauthEndpointTW, oauthEndpointCN *url.URL

func init() {
	mustParse := func(s string) (url *url.URL) {
		var err error
		if url, err = url.Parse(s); err != nil {
			panic(err)
		}
		return
	}

	oauthEndpointUS = mustParse("https://us.battle.net/oauth/token")
	oauthEndpointEU = mustParse("https://eu.battle.net/oauth/token")
	oauthEndpointKR = mustParse("https://kr.battle.net/oauth/token")
	oauthEndpointTW = mustParse("https://tw.battle.net/oauth/token")
	oauthEndpointCN = mustParse("https://www.battlenet.com.cn/oauth/token")
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

func newOAUTHEndpoint(region blizzard.Region) *url.URL {
	switch region {
	case blizzard.US:
		return cloneURL(oauthEndpointUS)
	case blizzard.EU:
		return cloneURL(oauthEndpointEU)
	case blizzard.KR:
		return cloneURL(oauthEndpointKR)
	case blizzard.TW:
		return cloneURL(oauthEndpointTW)
	case blizzard.CN:
		return cloneURL(oauthEndpointCN)
	}
	panic("Missing region in switch")
}
