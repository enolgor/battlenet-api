package wow

import "github.com/enolgor/battlenet-api/blizzard"

type Pet struct {
	ID   uint64 `json:"id"`
	Name *blizzard.LocalizedField
}

type PetExtended struct {
	IsCapturable bool `json:"is_capturable"`
	IsTradable   bool `json:"is_tradable"`
	IsBattlePet  bool `json:"is_battlepet"`
	//...
}

type petApi interface {
	GetPetIndex() ([]Pet, error)
}

func (wc *wowClientImpl) GetPetIndex() ([]Pet, error) {
	respStruct := struct {
		Pets []Pet `json:"pets"`
	}{}
	if err := wc.getGameData("/data/wow/pet/index", blizzard.Static, &respStruct); err != nil {
		return nil, err
	}
	return respStruct.Pets, nil
}
