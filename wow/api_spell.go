package wow

import "github.com/enolgor/battlenet-api/blizzard"

type Spell struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   uint64                   `json:"id"`
}
