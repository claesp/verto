package types

import (
	"fmt"
)

type VertoDevice struct {
	Hostname string
}

func (d VertoDevice) String() string {
	return fmt.Sprintf("%s", d.Hostname)
}
