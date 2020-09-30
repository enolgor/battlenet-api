package wow

import "github.com/enolgor/battlenet-api/blizzard"

type PlayableClassID uint8

const (
	HUNTER PlayableClassID = 3
)

type PlayableClass struct {
	Name *blizzard.LocalizedField `json:"name"`
	ID   *PlayableClassID         `json:"id"`
}
