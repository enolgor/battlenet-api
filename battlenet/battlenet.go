package battlenet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/enolgor/battlenet-api/blizzard"
	"github.com/moul/http2curl"
)

type BattleNetClient interface {
	battlenetLogger
	GetRegion() blizzard.Region
	SetRegion(region blizzard.Region)
	GetAccessToken() (string, error)
	NewGameDataRequest(path string, locale blizzard.Locale, namespace blizzard.Namespace) (*http.Request, error)
	GetGameData(path string, locale blizzard.Locale, namespace blizzard.Namespace, v interface{}) error
}

type battleNetClientImpl struct {
	battlenetLogger
	region       blizzard.Region
	clientID     string
	clientSecret string
	accessToken  string
	expiration   time.Time
	httpClient   *http.Client
}

func NewBattleNetClient(region blizzard.Region, clientID string, clientSecret string) BattleNetClient {
	bnc := &battleNetClientImpl{region: region, clientID: clientID, clientSecret: clientSecret, httpClient: &http.Client{}}
	bnc.battlenetLogger = newBattleNetLogger()
	return bnc
}

func (bnci *battleNetClientImpl) SetRegion(region blizzard.Region) {
	bnci.region = region
	bnci.accessToken = ""
}

func (bnci *battleNetClientImpl) GetAccessToken() (string, error) {
	if bnci.clientSecret == "" {
		return "", fmt.Errorf("Must specify client secret first")
	}
	if bnci.accessToken == "" || bnci.expiration.Before(time.Now()) {
		respStruct := struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int64  `json:"expires_in"`
		}{}
		ep := newOAUTHEndpoint(bnci.region)
		ep.User = url.UserPassword(bnci.clientID, bnci.clientSecret)
		response, err := http.PostForm(ep.String(), url.Values{
			"grant_type": {"client_credentials"},
		})
		if err != nil {
			return "", err
		}
		defer response.Body.Close()
		dec := json.NewDecoder(response.Body)
		if err = dec.Decode(&respStruct); err != nil {
			return "", err
		}
		bnci.accessToken = respStruct.AccessToken
		bnci.expiration = time.Now().Add(time.Second * time.Duration(respStruct.ExpiresIn))
		return bnci.accessToken, nil
	}
	return bnci.accessToken, nil
}

func (bnci *battleNetClientImpl) GetRegion() blizzard.Region {
	return bnci.region
}

func (bnci *battleNetClientImpl) NewGameDataRequest(path string, locale blizzard.Locale, namespace blizzard.Namespace) (*http.Request, error) {
	req := http.Request{}
	req.URL = blizzard.NewAPIEndpoint(bnci.region)
	req.URL.Path = path
	at, err := bnci.GetAccessToken()
	if err != nil {
		bnci.getLogger(ERROR).Print(err)
		return nil, err
	}
	req.Header = make(http.Header)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", at))
	if namespace != blizzard.NoNamespace {
		req.Header.Add("Battlenet-Namespace", namespace.ForRegion(bnci.region))
	}
	if locale != blizzard.NoLocale {
		query := req.URL.Query()
		query.Add("locale", string(locale))
		req.URL.RawQuery = query.Encode()
	}
	return &req, nil
}

func (bnci *battleNetClientImpl) GetGameData(path string, locale blizzard.Locale, namespace blizzard.Namespace, receiver interface{}) error {
	var err error
	var req *http.Request
	var resp *http.Response
	if req, err = bnci.NewGameDataRequest(path, locale, namespace); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	req.Method = http.MethodGet
	curlReq, _ := http2curl.GetCurlCommand(req)
	bnci.getLogger(INFO).Printf("REQ: %s\n", curlReq)
	if resp, err = bnci.httpClient.Do(req); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	bnci.getLogger(INFO).Printf("Response StatusCode: %d\n", resp.StatusCode)
	var reader io.Reader = resp.Body
	bnci.getLogger(DEBUG).Println("Body>>>")
	reader = io.TeeReader(resp.Body, bnci.getLogger(DEBUG).Writer())
	defer bnci.getLogger(DEBUG).Println("<<<EndBody")
	defer bnci.getLogger(DEBUG).Writer().Write([]byte("\n"))
	defer resp.Body.Close()
	if resp.StatusCode/100 != http.StatusOK/100 {
		err = newBattlenetError(reader)
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	if err = json.NewDecoder(reader).Decode(&receiver); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	return nil
}
