package battlenet

import (
	"net/url"

	"github.com/enolgor/battlenet-api/blizzard"
)

var oauthTokenEndpointUS, oauthTokenEndpointEU, oauthTokenEndpointKR, oauthTokenEndpointTW, oauthTokenEndpointCN *url.URL
var oauthAuthorizeEndpointUS, oauthAuthorizeEndpointEU, oauthAuthorizeEndpointKR, oauthAuthorizeEndpointTW, oauthAuthorizeEndpointCN *url.URL

func init() {
	mustParse := func(s string) (url *url.URL) {
		var err error
		if url, err = url.Parse(s); err != nil {
			panic(err)
		}
		return
	}

	oauthTokenEndpointUS = mustParse("https://us.battle.net/oauth/token")
	oauthTokenEndpointEU = mustParse("https://eu.battle.net/oauth/token")
	oauthTokenEndpointKR = mustParse("https://kr.battle.net/oauth/token")
	oauthTokenEndpointTW = mustParse("https://tw.battle.net/oauth/token")
	oauthTokenEndpointCN = mustParse("https://www.battlenet.com.cn/oauth/token")

	oauthAuthorizeEndpointUS = mustParse("https://us.battle.net/oauth/authorize")
	oauthAuthorizeEndpointEU = mustParse("https://eu.battle.net/oauth/authorize")
	oauthAuthorizeEndpointKR = mustParse("https://kr.battle.net/oauth/authorize")
	oauthAuthorizeEndpointTW = mustParse("https://tw.battle.net/oauth/authorize")
	oauthAuthorizeEndpointCN = mustParse("https://www.battlenet.com.cn/oauth/authorize")
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

func newOAuthTokenEndpoint(region blizzard.Region) *url.URL {
	switch region {
	case blizzard.US:
		return cloneURL(oauthTokenEndpointUS)
	case blizzard.EU:
		return cloneURL(oauthTokenEndpointEU)
	case blizzard.KR:
		return cloneURL(oauthTokenEndpointKR)
	case blizzard.TW:
		return cloneURL(oauthTokenEndpointTW)
	case blizzard.CN:
		return cloneURL(oauthTokenEndpointCN)
	}
	panic("Missing region in switch")
}

func newOAuthAuthorizeEndpoint(region blizzard.Region) *url.URL {
	switch region {
	case blizzard.US:
		return cloneURL(oauthAuthorizeEndpointUS)
	case blizzard.EU:
		return cloneURL(oauthAuthorizeEndpointEU)
	case blizzard.KR:
		return cloneURL(oauthAuthorizeEndpointKR)
	case blizzard.TW:
		return cloneURL(oauthAuthorizeEndpointTW)
	case blizzard.CN:
		return cloneURL(oauthAuthorizeEndpointCN)
	}
	panic("Missing region in switch")
}

const OAuthScopeWoWProfile string = "wow.profile"
const OAuthScopeD3Profile string = "d3.profile"
const OAuthScopeSC2Profile string = "sc2.profile"
const OAuthScopeOpenID string = "openid"
