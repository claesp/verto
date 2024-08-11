package parser

import (
	"github.com/claesp/verto/internal/importer"
	"github.com/claesp/verto/internal/types"
)

func Parse(importer importer.Importer, data []byte) (types.VertoDevice, error) {
	var d types.VertoDevice

	ie := importer.Import(string(data))
	if ie != nil {
		return d, ie
	}

	pe := importer.Parse()
	if pe != nil {
		return d, pe
	}

	d = importer.ExtractDevice()

	return d, nil
}
