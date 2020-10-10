package blizzard

import (
	"net/url"
)

var apiEndpointUS, apiEndpointEU, apiEndpointKR, apiEndpointTW, apiEndpointCN *url.URL

func init() {
	mustParse := func(s string) (url *url.URL) {
		var err error
		if url, err = url.Parse(s); err != nil {
			panic(err)
		}
		return
	}
	apiEndpointUS = mustParse("https://us.api.blizzard.com/")
	apiEndpointEU = mustParse("https://eu.api.blizzard.com/")
	apiEndpointKR = mustParse("https://kr.api.blizzard.com/")
	apiEndpointTW = mustParse("https://tw.api.blizzard.com/")
	apiEndpointCN = mustParse("https://gateway.battlenet.com.cn/")
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

func NewAPIEndpoint(region Region) *url.URL {
	switch region {
	case US:
		return cloneURL(apiEndpointUS)
	case EU:
		return cloneURL(apiEndpointEU)
	case KR:
		return cloneURL(apiEndpointKR)
	case TW:
		return cloneURL(apiEndpointTW)
	case CN:
		return cloneURL(apiEndpointCN)
	}
	panic("Missing region in switch")
}
