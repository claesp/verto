package parser

import (
	"fmt"
	"github.com/claesp/verto/internal/importer"
	"os"
)

func Parse(imp importer.Importer, filename string) error {
	d := imp.ImportFromText("test")

	fmt.Fprintf(os.Stdout, "%s\n", d.Hostname)

	return nil
}
