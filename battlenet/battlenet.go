package battlenet

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
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
	NewGameDataSearchRequest(path string, query SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace) (*http.Request, error)
	GetGameData(path string, locale blizzard.Locale, namespace blizzard.Namespace, v interface{}) error
	SearchGameData(path string, query SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace, v interface{}) (*SearchResult, error)
	GetAuthorizationURI(scopes ...string) (string, *url.URL, error)
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
		ep := newOAuthTokenEndpoint(bnci.region)
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

func (bnci *battleNetClientImpl) GetAuthorizationURI(scopes ...string) (string, *url.URL, error) {
	randBytes := make([]byte, 20)
	if _, err := rand.Read(randBytes); err != nil {
		return "", nil, err
	}
	randStr := base32.StdEncoding.EncodeToString(randBytes)
	ep := newOAuthAuthorizeEndpoint(bnci.region)
	scopeStr := strings.Join(scopes, " ")
	values := ep.Query()
	values.Add("client_id", bnci.clientID)
	values.Add("scope", scopeStr)
	values.Add("state", randStr)
	values.Add("redirect_uri", "https://localhost")
	values.Add("response_type", "code")
	ep.RawQuery = values.Encode()
	return randStr, ep, nil
}

func (bnci *battleNetClientImpl) GetRegion() blizzard.Region {
	return bnci.region
}

func (bnci *battleNetClientImpl) NewGameDataRequest(path string, locale blizzard.Locale, namespace blizzard.Namespace) (*http.Request, error) {
	req := http.Request{}
	req.URL = blizzard.NewAPIEndpoint(bnci.region)
	req.URL.Path = path
	req.Method = http.MethodGet
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

func (bnci *battleNetClientImpl) NewGameDataSearchRequest(path string, query SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace) (*http.Request, error) {
	req, err := bnci.NewGameDataRequest(path, locale, namespace)
	if err != nil {
		return nil, err
	}
	query.appendTo(req)
	return req, nil
}

func (bnci *battleNetClientImpl) doGet(req *http.Request, receiver interface{}) error {
	var resp *http.Response
	var err error
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

func (bnci *battleNetClientImpl) GetGameData(path string, locale blizzard.Locale, namespace blizzard.Namespace, receiver interface{}) error {
	var req *http.Request
	var err error
	if req, err = bnci.NewGameDataRequest(path, locale, namespace); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	return bnci.doGet(req, receiver)
}

func (bnci *battleNetClientImpl) SearchGameData(path string, query SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace, receiver interface{}) (*SearchResult, error) {
	var err error
	var req *http.Request
	if req, err = bnci.NewGameDataSearchRequest(path, query, locale, namespace); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return nil, err
	}
	searchResult := struct {
		SearchResult
		Results json.RawMessage `json:"results"`
	}{}
	if err := bnci.doGet(req, &searchResult); err != nil {
		return nil, err
	}
	searchResultData := []struct {
		Data json.RawMessage `json:"data"`
	}{}
	if err = json.Unmarshal(searchResult.Results, &searchResultData); err != nil {
		return nil, err
	}
	data := reflect.ValueOf(receiver).Elem()
	objType := reflect.TypeOf(receiver).Elem().Elem()
	for _, srd := range searchResultData {
		elem := reflect.New(objType).Interface()
		if err := json.Unmarshal(srd.Data, &elem); err != nil {
			return nil, err
		}
		data = reflect.Append(data, reflect.ValueOf(elem).Elem())
	}
	reflect.ValueOf(receiver).Elem().Set(data)
	return &searchResult.SearchResult, err
}
