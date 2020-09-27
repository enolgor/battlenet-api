package blizzard

import (
	"net/url"
)

var apiEndpointUS, apiEndpointEU, apiEndpointKR, apiEndpointTW, apiEndpointCN *url.URL
var oauthEndpointUS, oauthEndpointEU, oauthEndpointKR, oauthEndpointTW, oauthEndpointCN *url.URL

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

	oauthEndpointUS = mustParse("https://us.battle.net/oauth/token")
	oauthEndpointEU = mustParse("https://eu.battle.net/oauth/token")
	oauthEndpointKR = mustParse("https://kr.battle.net/oauth/token")
	oauthEndpointTW = mustParse("https://tw.battle.net/oauth/token")
	oauthEndpointCN = mustParse("https://www.battlenet.com.cn/oauth/token") //????
}

type Region string

const (
	NoRegion Region = ""
	US       Region = "us"
	EU       Region = "eu"
	KR       Region = "kr"
	TW       Region = "tw"
	CN       Region = "cn"
)

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

func (reg Region) oauthEndpoint() *url.URL {
	switch reg {
	case US:
		return cloneURL(oauthEndpointUS)
	case EU:
		return cloneURL(oauthEndpointEU)
	case KR:
		return cloneURL(oauthEndpointKR)
	case TW:
		return cloneURL(oauthEndpointTW)
	case CN:
		return cloneURL(oauthEndpointCN)
	}
	panic("Missing region in switch")
}

func (reg Region) Locales() []Locale {
	switch reg {
	case US:
		return []Locale{EnUS, EsMX, PtBR}
	case EU:
		return []Locale{EnGB, EsES, FrFR, RuRU, DeDE, PtPT, ItIT}
	case KR:
		return []Locale{KoKR}
	case TW:
		return []Locale{ZhTW}
	case CN:
		return []Locale{ZhCN}
	}
	panic("Missing region in switch")
}

type Locale string

const (
	NoLocale Locale = ""
	EnUS     Locale = "en_US"
	EsMX     Locale = "es_MX"
	PtBR     Locale = "pt_BR"
	EnGB     Locale = "en_GB"
	EsES     Locale = "es_ES"
	FrFR     Locale = "fr_FR"
	RuRU     Locale = "ru_RU"
	DeDE     Locale = "de_DE"
	PtPT     Locale = "pt_PT"
	ItIT     Locale = "it_IT"
	KoKR     Locale = "ko_KR"
	ZhTW     Locale = "zh_TW"
	ZhCN     Locale = "zh_CN"
)

func GetLocale(localeStr string) Locale {
	switch localeStr {
	case "en_US":
		return EnUS
	case "es_MX":
		return EsMX
	case "pt_BR":
		return PtBR
	case "en_GB":
		return EnGB
	case "es_ES":
		return EsES
	case "fr_FR":
		return FrFR
	case "ru_RU":
		return RuRU
	case "de_DE":
		return DeDE
	case "pt_PT":
		return PtPT
	case "it_IT":
		return ItIT
	case "ko_KR":
		return KoKR
	case "zh_TW":
		return ZhTW
	case "zh_CN":
		return ZhCN
	default:
		return NoLocale
	}
}
