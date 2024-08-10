package main

import (
	"fmt"
	"github.com/claesp/verto/internal/importer"
	"github.com/claesp/verto/internal/parser"
	"os"
)

var (
	AppName string = "verto-cli"
	Version string = "0.0.0-ffffff"
)

func e(text string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", text)
	os.Exit(1)
}

type AppArguments struct {
	VendorSpecified   bool
	VendorName        string
	ReadFileSpecified bool
	ReadFile          string
}

func parseArgs() AppArguments {
	a := AppArguments{}

	for i, arg := range os.Args {
		switch arg {
		case "-f":
			a.ReadFileSpecified = true
			if len(os.Args) >= i+2 {
				a.ReadFile = os.Args[i+1]
			} else {
				e("missing filename")
			}
		case "-t":
			a.VendorSpecified = true
			if len(os.Args) >= i+2 {
				a.VendorName = os.Args[i+1]
			} else {
				e("missing vendor")
			}
		case "-v":
			fmt.Fprintf(os.Stdout, "%s %s\n", AppName, Version)
			os.Exit(0)
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

func parseDataFromFile(importer importer.Importer, filename string) error {
	d, loadErr := loadDataFromFile(filename)
	if loadErr != nil {
		return loadErr
	}

	parseErr := parser.Parse(importer, string(d))
	if parseErr != nil {
		return parseErr
	}

	return nil
}

func selectVendorImporter(vendor string) (importer.Importer, error) {
	switch vendor {
	case "fortigate":
		return importer.FortiOSImporter{}, nil
	default:
		return nil, fmt.Errorf("unknown vendor")
	}
}

func main() {
	a := parseArgs()

	var i importer.Importer
	if a.VendorSpecified {
		var vendErr error
		i, vendErr = selectVendorImporter(a.VendorName)
		if vendErr != nil {
			e(vendErr.Error())
		}
	}

	if a.ReadFileSpecified {
		parseErr := parseDataFromFile(i, a.ReadFile)
		if parseErr != nil {
			e(parseErr.Error())
		}
	}

	os.Exit(0)
}
