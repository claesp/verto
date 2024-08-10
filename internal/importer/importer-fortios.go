package importer

import (
	"github.com/claesp/verto/internal/types"
)

type FortiOSImporter struct {
}

func (f FortiOSImporter) ImportFromText(s string) types.VertoDevice {
	d := types.VertoDevice{}

	d.Hostname = s

	return d
}
