package blizzard

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LocalizedField map[Locale]string

func (lf LocalizedField) String() string {
	if len(lf) == 1 {
		return lf[NoLocale]
	}
	elems := make([]string, 0, len(lf))
	for loc, v := range lf {
		elems = append(elems, fmt.Sprintf("%s: %s", loc, v))
	}
	return fmt.Sprintf("[%s]", strings.Join(elems, ", "))
}

func (lf *LocalizedField) UnmarshalJSON(data []byte) error {
	v := make(map[Locale]string)
	if data[0] == '{' {
		d := make(map[string]string)
		json.Unmarshal(data, &d)
		for key, value := range d {
			v[Locale(key)] = value
		}
	} else {
		recv := ""
		json.Unmarshal(data, &recv)
		v[NoLocale] = recv
	}
	*lf = v
	return nil
}

func (lf *LocalizedField) MarshalJSON() ([]byte, error) {
	if len(*lf) == 1 {
		return json.Marshal(string((*lf)[NoLocale]))
	}
	data := make(map[string]string)
	for k, v := range *lf {
		data[string(k)] = v
	}
	return json.Marshal(data)
}
