package blizzard

import (
	"fmt"
	"net/http"
)

/* From: https://develop.battle.net/documentation/world-of-warcraft/guides/search

AND	str=5&dex=10	An implicit AND operation is performed by multiple query parameters.
OR	str=5||10
type=man||bear||pig	An OR operation is performed by placing a pair of bars between two values.
NOT	race!=orc	A NOT operation is performed by using a combination of an exclamation mark and equal sign between two values. You can combine the NOT and OR operations in a single statement (for example, race!=orc||human).
RANGE	str=[2,99]
str=(2,99)	A RANGE operation is performed only on numeric field values by using either brackets (inclusive) or parentheses (exclusive) around comma-separated minimum and maximum values.
MIN	str=[41,]
str=(41,)	A MIN operation performs a minimum value check using the Range syntax.
MAX	str=[,77]
str=(,77)	A MAX operation performs a maximum value check using the Range syntax.
*/

type SearchResult struct {
	Page              int  `json:"page"`
	PageSize          int  `json:"pageSize"`
	MaxPageSize       int  `json:"maxPageSize"`
	PageCount         int  `json:"pageCount"`
	ResultCountCapped bool `json:"resultCountCapped"`
}

type SearchQuery []SearchParam

type SearchParam []string

func Query(params ...SearchParam) SearchQuery {
	return SearchQuery(params)
}

func Param(key, value string) SearchParam {
	return SearchParam([]string{key, value})
}

var Params = struct {
	Page        func(value int) SearchParam
	PageSize    func(value int) SearchParam
	OrderBy     func(field string) SearchParam
	OrderByAsc  func(field string) SearchParam
	OrderByDesc func(field string) SearchParam
}{
	Page: func(value int) SearchParam {
		return Param("_page", fmt.Sprintf("%d", value))
	},
	PageSize: func(value int) SearchParam {
		return Param("_pageSize", fmt.Sprintf("%d", value))
	},
	OrderBy: func(field string) SearchParam {
		return Param("orderby", field)
	},
	OrderByAsc: func(field string) SearchParam {
		return Param("orderby", fmt.Sprintf("%s:asc", field))
	},
	OrderByDesc: func(field string) SearchParam {
		return Param("orderby", fmt.Sprintf("%s:desc", field))
	},
}

func (sq SearchQuery) AddToRequest(req *http.Request) {
	q := req.URL.Query()
	for _, query := range sq {
		q.Add(query[0], query[1])
	}
	req.URL.RawQuery = q.Encode()
}
