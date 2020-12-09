package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/enolgor/battlenet-api/battlenet"
	"github.com/enolgor/battlenet-api/blizzard"
	"github.com/enolgor/battlenet-api/wow"
)

func main() {
	bnc, err := battlenet.NewBattleNetClient(blizzard.EU, "client_id", "client_secret")
	if err != nil {
		panic(err)
	}
	//bnc.SetLogger(log.New(os.Stdout, "", log.LstdFlags), battlenet.ERROR, battlenet.INFO, battlenet.DEBUG)
	/*wowclient, err := wow.NewWoWClient(bnc, blizzard.NoLocale)
	if err != nil {
		panic(err)
	}
	url, err := bnc.GetAuthorizationURI("https://localhost:8123", "wow.profile")
	fmt.Println(url)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		code, err := bnc.ParseAuthorizationResponse(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, expiration, err := bnc.GetAuthorizationToken(code, "https://localhost:8123", "wow.profile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		userinfo, err := bnc.UserInfo(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(err)
			return
		}
		profile, err := wowclient.GetProfile(token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			encoder.Encode(err)
			return
		}
		w.Write([]byte(code + "\n"))
		w.Write([]byte(token + "\n"))
		w.Write([]byte(expiration.Format(time.RFC822) + "\n"))
		encoder.Encode(userinfo)
		w.Write([]byte("\n"))
		encoder.Encode(profile)
	})
	log.Fatal(http.ListenAndServeTLS(":8123", "server.crt", "server.key", nil))*/

	//wowclient := wow.NewWoWClient(bnc, blizzard.EnUS)
	//print(wowclient.GetAuctions(1403))

	wowclient, err := wow.NewWoWClient(bnc, blizzard.EsES)
	if err != nil {
		log.Fatal(err)
	}
	skillTier, err := wowclient.GetSkillTier(164, 2751)
	if err != nil {
		log.Fatal(err)
	}
	for _, cat := range skillTier.Categories {
		for _, rec := range cat.Recipes {
			recipe, err := wowclient.GetRecipe(rec.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(recipe)
		}
	}

}

func print(data interface{}, err error) {
	if err != nil {
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(data); err != nil {
		log.Fatal(err)
	}
}
