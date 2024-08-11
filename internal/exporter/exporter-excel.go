package exporter

import (
	"github.com/claesp/verto/internal/types"
)

type ExcelExporter struct {
}

func NewExcelExporter() ExcelExporter {
	return ExcelExporter{}
}

func (e ExcelExporter) Export(vd types.VertoDevice) {
}
