package main

import (
	"fmt"
	"os"

	"github.com/claesp/verto/internal/exporter"
	"github.com/claesp/verto/internal/importer"
	"github.com/claesp/verto/internal/parser"
	"github.com/claesp/verto/internal/types"
)

var (
	AppName string = "verto-cli"
	Version string = "0.0.0-ffffff"
)

func e(text string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", text)
	os.Exit(1)
}

func o(text string) {
	fmt.Fprintf(os.Stdout, "%s\n", text)
	os.Exit(0)
}

type AppArguments struct {
	ImportVendorSpecified bool
	ImportVendorName      string
	ExportVendorSpecified bool
	ExportVendorName      string
	ReadFileSpecified     bool
	ReadFile              string
}

func help() {
	s := fmt.Sprintf("usage: %s [-h | -v ] -i import_type -e export_type -f source_file", AppName)
	o(s)
}

func parseArgs() AppArguments {
	a := AppArguments{}

	for i, arg := range os.Args {
		switch arg {
		case "-e":
			a.ExportVendorSpecified = true
			if len(os.Args) >= i+2 {
				a.ExportVendorName = os.Args[i+1]
			} else {
				e("export: missing export vendor")
			}
		case "-f":
			a.ReadFileSpecified = true
			if len(os.Args) >= i+2 {
				a.ReadFile = os.Args[i+1]
			} else {
				e("file: missing filename")
			}
		case "-i":
			a.ImportVendorSpecified = true
			if len(os.Args) >= i+2 {
				a.ImportVendorName = os.Args[i+1]
			} else {
				e("import: missing import vendor")
			}
		case "-v":
			o(fmt.Sprintf("%s %s", AppName, Version))
		case "-h":
			help()
		}
	}

	return a
}

func loadDataFromFile(filename string) ([]byte, error) {
	d, readErr := os.ReadFile(filename)
	if readErr != nil {
		return []byte{}, readErr
	}

	return d, nil
}

func parseDataFromFile(imp importer.Importer, par parser.Parser, filename string) (types.VertoDevice, error) {
	var device types.VertoDevice

	data, loadErr := loadDataFromFile(filename)
	if loadErr != nil {
		return device, loadErr
	}

	importErr := imp.Import(data)
	if importErr != nil {
		return device, importErr
	}

	device, parseErr := par.Parse(imp)
	if parseErr != nil {
		return device, parseErr
	}

	return device, nil
}

func selectVendorImporter(vendor string) (importer.Importer, parser.Parser, error) {
	switch vendor {
	case "fortigate":
		return importer.NewFortiOSImporter(), parser.NewFortiOSTextParser(), nil
	default:
		return nil, nil, fmt.Errorf("unknown import vendor")
	}
}

func selectVendorExporter(vendor string) (exporter.Exporter, error) {
	switch vendor {
	case "excel":
		return exporter.NewExcelExporter(), nil
	default:
		return nil, fmt.Errorf("unknown export vendor")
	}
}

func main() {
	if len(os.Args) == 1 {
		help()
	}

	a := parseArgs()

	var impImpl importer.Importer
	var importParser parser.Parser
	if a.ImportVendorSpecified {
		var importErr error
		impImpl, importParser, importErr = selectVendorImporter(a.ImportVendorName)
		if importErr != nil {
			e(importErr.Error())
		}
	} else {
		e("missing importer")
	}

	var expImpl exporter.Exporter
	if a.ExportVendorSpecified {
		var exportErr error
		expImpl, exportErr = selectVendorExporter(a.ExportVendorName)
		if exportErr != nil {
			e(exportErr.Error())
		}
	} else {
		e("missing exporter")
	}

	if a.ReadFileSpecified {
		device, parseErr := parseDataFromFile(impImpl, importParser, a.ReadFile)
		if parseErr != nil {
			e(parseErr.Error())
		}

		expImpl.Export(device)
	} else {
		e("no filename specified")
	}

	os.Exit(0)
}
