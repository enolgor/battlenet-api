package blizzard

import (
	"fmt"
)

type Namespace string

const (
	NoNamespace Namespace = ""
	Dynamic     Namespace = "dynamic"
	Static      Namespace = "static"
	Profile     Namespace = "profile"
)

func (nt Namespace) ForRegion(reg Region) string {
	return fmt.Sprintf("%s-%s", nt, reg.String())
}
