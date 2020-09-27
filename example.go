package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/enolgor/battlenet-api/battlenet"
	"github.com/enolgor/battlenet-api/blizzard"
	"github.com/enolgor/battlenet-api/wow"
)

func main() {
	bnc := battlenet.NewBattleNetClient(blizzard.EU, "client_id", "client_secret")
	bnc.SetLogOutput(os.Stdout, battlenet.ERROR, battlenet.INFO)
	wowclient := wow.NewWoWClient(bnc)
	print(wowclient.GetAuctions(1403))
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