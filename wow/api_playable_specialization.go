package wow

import "github.com/enolgor/battlenet-api/blizzard"

type Specialization struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   uint16                   `json:"id"`
}
