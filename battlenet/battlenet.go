package battlenet

import (
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
	GetAccessToken() (string, error)
	NewBlizzardAPIRequest(path string, locale blizzard.Locale, namespace blizzard.Namespace, searchQuery blizzard.SearchQuery, token string) (*http.Request, error)
	NewBattleNetRequest(path string, token string) (*http.Request, error)
	BlizzardAPIGet(path string, locale blizzard.Locale, namespace blizzard.Namespace, token string, v interface{}) error
	BlizzardAPISearch(path string, query blizzard.SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace, token string, v interface{}) (*blizzard.SearchResult, error)
	BattleNetGet(path string, token string, v interface{}) error
	GetAuthorizationURI(redirectURI string, scopes ...string) (*url.URL, error)
	ParseAuthorizationResponse(req *http.Request) (string, error)
	GetAuthorizationToken(code string, redirectURI string, scopes ...string) (string, *time.Time, error)
	UserInfo(token string) (*UserInfo, error)
}

type UserInfo struct {
	Sub       string `json:"sub"`
	ID        uint64 `json:"id"`
	BattleTag string `json:"battletag"`
}

type battleNetClientImpl struct {
	battlenetLogger
	region        blizzard.Region
	clientID      string
	clientSecret  string
	accessToken   string
	expiration    time.Time
	httpClient    *http.Client
	oAuthStateGen *oAuthStateGen
}

func NewBattleNetClient(region blizzard.Region, clientID string, clientSecret string) (BattleNetClient, error) {
	oAuthStateGen, err := newOAuthStateGen(clientID, clientSecret)
	if err != nil {
		return nil, err
	}
	bnc := &battleNetClientImpl{region: region, clientID: clientID, clientSecret: clientSecret, httpClient: &http.Client{}, oAuthStateGen: oAuthStateGen}
	bnc.battlenetLogger = newBattleNetLogger()
	return bnc, nil
}

func (bnci *battleNetClientImpl) GetAccessToken() (string, error) {
	if bnci.accessToken == "" || bnci.expiration.Before(time.Now()) {
		respStruct := struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int64  `json:"expires_in"`
		}{}
		ep := NewEndpoint(bnci.region, "/oauth/token")
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

func (bnci *battleNetClientImpl) GetAuthorizationURI(redirectURI string, scopes ...string) (*url.URL, error) {
	state := bnci.oAuthStateGen.generate()
	ep := NewEndpoint(bnci.region, "/oauth/authorize")
	scopeStr := strings.Join(scopes, " ")
	values := ep.Query()
	values.Add("client_id", bnci.clientID)
	values.Add("scope", scopeStr)
	values.Add("state", state)
	values.Add("redirect_uri", redirectURI)
	values.Add("response_type", "code")
	ep.RawQuery = values.Encode()
	return ep, nil
}

func (bnci *battleNetClientImpl) ParseAuthorizationResponse(req *http.Request) (string, error) {
	state := req.URL.Query().Get("state")
	if state == "" {
		return "", fmt.Errorf("Missing state value in request query")
	}
	code := req.URL.Query().Get("code")
	if code == "" {
		return "", fmt.Errorf("Missing code value in request query")
	}
	ok, err := bnci.oAuthStateGen.parse(state)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("State does not match with client")
	}
	return code, nil
}

func (bnci *battleNetClientImpl) GetAuthorizationToken(code string, redirectURI string, scopes ...string) (string, *time.Time, error) {
	scopeStr := strings.Join(scopes, " ")
	if bnci.clientSecret == "" {
		return "", nil, fmt.Errorf("Must specify client secret first")
	}
	respStruct := struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}{}
	ep := NewEndpoint(bnci.region, "/oauth/token")
	ep.User = url.UserPassword(bnci.clientID, bnci.clientSecret)
	response, err := http.PostForm(ep.String(), url.Values{
		"redirect_uri": {redirectURI},
		"scope":        {scopeStr},
		"grant_type":   {"authorization_code"},
		"code":         {code},
	})
	if err != nil {
		return "", nil, err
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	if err = dec.Decode(&respStruct); err != nil {
		return "", nil, err
	}
	expiration := time.Now().Add(time.Second * time.Duration(respStruct.ExpiresIn))
	return respStruct.AccessToken, &expiration, nil
}

func (bnci *battleNetClientImpl) NewBattleNetRequest(path string, token string) (*http.Request, error) {
	req := &http.Request{}
	req.URL = NewEndpoint(bnci.region, path)
	req.URL.Path = path
	req.Method = http.MethodGet
	req.Header = make(http.Header)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return req, nil
}

func (bnci *battleNetClientImpl) NewBlizzardAPIRequest(path string, locale blizzard.Locale, namespace blizzard.Namespace, searchQuery blizzard.SearchQuery, token string) (*http.Request, error) {
	req := &http.Request{}
	req.URL = blizzard.NewAPIEndpoint(bnci.region)
	req.URL.Path = path
	req.Method = http.MethodGet
	req.Header = make(http.Header)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if namespace != blizzard.NoNamespace {
		req.Header.Add("Battlenet-Namespace", namespace.ForRegion(bnci.region))
	}
	if locale != blizzard.NoLocale {
		query := req.URL.Query()
		query.Add("locale", string(locale))
		req.URL.RawQuery = query.Encode()
	}
	if searchQuery != nil {
		searchQuery.AddToRequest(req)
	}
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

func (bnci *battleNetClientImpl) BlizzardAPIGet(path string, locale blizzard.Locale, namespace blizzard.Namespace, token string, receiver interface{}) error {
	var req *http.Request
	var err error
	if req, err = bnci.NewBlizzardAPIRequest(path, locale, namespace, nil, token); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	return bnci.doGet(req, receiver)
}

func (bnci *battleNetClientImpl) BlizzardAPISearch(path string, query blizzard.SearchQuery, locale blizzard.Locale, namespace blizzard.Namespace, token string, receiver interface{}) (*blizzard.SearchResult, error) {
	var err error
	var req *http.Request
	if req, err = bnci.NewBlizzardAPIRequest(path, locale, namespace, query, token); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return nil, err
	}
	searchResult := struct {
		blizzard.SearchResult
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

func (bnci *battleNetClientImpl) BattleNetGet(path string, token string, receiver interface{}) error {
	var req *http.Request
	var err error
	if req, err = bnci.NewBattleNetRequest(path, token); err != nil {
		bnci.getLogger(ERROR).Print(err)
		return err
	}
	return bnci.doGet(req, receiver)
}

func (bnci *battleNetClientImpl) UserInfo(token string) (*UserInfo, error) {
	userInfo := &UserInfo{}
	if err := bnci.BattleNetGet("/oauth/userinfo", token, userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}
