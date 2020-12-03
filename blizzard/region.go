package blizzard

type Region uint8

const (
	NoRegion Region = 0
	US       Region = 1
	EU       Region = 3
	KR       Region = 2
	TW       Region = 4
	CN       Region = 5
)

func (reg Region) String() string {
	switch reg {
	case US:
		return "us"
	case EU:
		return "eu"
	case KR:
		return "kr"
	case TW:
		return "tw"
	case CN:
		return "cn"
	}
	panic("Missing region in switch")
}

func GetRegion(regionStr string) Region {
	switch regionStr {
	case "us":
		return US
	case "eu":
		return EU
	case "kr":
		return KR
	case "tw":
		return TW
	case "cn":
		return CN
	default:
		return NoRegion
	}
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
