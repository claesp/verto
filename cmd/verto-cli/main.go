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

type ConfigArguments struct {
	VendorSpecified bool
	VendorName      string
	ReadFromFile    bool
	ReadFilename    string
}

func parseArgs() ConfigArguments {
	a := ConfigArguments{}

	for i, arg := range os.Args {
		switch arg {
		case "-t":
			a.VendorSpecified = true
			if len(os.Args) >= i+2 {
				a.VendorName = os.Args[i+1]
			} else {
				e("missing vendor")
			}
		case "-f":
			a.ReadFromFile = true
			if len(os.Args) >= i+2 {
				a.ReadFilename = os.Args[i+1]
			} else {
				e("missing filename")
			}
		}
	}

	return a
}

func loadDataFromFile(filename string) ([]byte, error) {
	d, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return d, nil
}

func main() {
	fmt.Fprintf(os.Stdout, "%s %s\n", AppName, Version)
	a := parseArgs()

	var d []byte
	if a.ReadFromFile {
		var loadErr error
		d, loadErr = loadDataFromFile(a.ReadFilename)
		if loadErr != nil {
			e(loadErr.Error())
		}
	} else {
		d = []byte{}
	}

	var i importer.Importer
	if a.VendorSpecified {
		switch a.VendorName {
		case "fortigate":
			i = importer.FortiOSImporter{}
		default:
			e("unknown vendor")
		}
	} else {
		e("undefined vendor")
	}

	parseErr := parser.Parse(i, string(d))
	if parseErr != nil {
		e(parseErr.Error())
	}

	os.Exit(0)
}
