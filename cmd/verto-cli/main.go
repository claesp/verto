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

func main() {
	fmt.Fprintf(os.Stdout, "%s %s\n", AppName, Version)
	i := importer.FortiOSImporter{}
	err := parser.Parse(i, "test")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}
	os.Exit(0)
}
