package battlenet

import (
	"encoding/json"
	"fmt"
	"io"
)

type BattlenetError struct {
	Code   int64  `json:"code"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

func (be *BattlenetError) Error() string {
	return fmt.Sprintf("code: %d, type: %s, detail: %s", be.Code, be.Type, be.Detail)
}

func newBattlenetError(reader io.Reader) error {
	berr := &BattlenetError{}
	if err := json.NewDecoder(reader).Decode(berr); err != nil {
		return err
	}
	return berr
}
