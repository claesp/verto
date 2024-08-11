package importer

import (
	"github.com/claesp/verto/internal/types"
)

type Importer interface {
	Import(s string) error
	Parse() error
	ExtractDevice() types.VertoDevice
}
