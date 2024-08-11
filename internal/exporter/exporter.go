package exporter

import (
	"github.com/claesp/verto/internal/types"
)

type Exporter interface {
	Export(types.VertoDevice)
}
