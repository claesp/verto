package importer

import (
	"github.com/claesp/verto/internal/types"
)

type Importer interface {
	ImportFromText(s string) types.VertoDevice
}
