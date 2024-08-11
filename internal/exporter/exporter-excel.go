package exporter

import (
	"github.com/claesp/verto/internal/types"
)

type Exporter interface {
	Export(types.VertoDevice)
}

type ExcelExporter struct {
}

func NewExcelExporter() ExcelExporter {
	return ExcelExporter{}
}

func (e ExcelExporter) Export(vd types.VertoDevice) {
}
