package blizzard

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
