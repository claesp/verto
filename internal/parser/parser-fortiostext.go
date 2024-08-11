package parser

import (
	"fmt"
	"github.com/claesp/verto/internal/importer"
	"github.com/claesp/verto/internal/types"
	"os"
)

type Parser interface {
	Parse(importer.Importer) (types.VertoDevice, error)
}

type FortiOSTextParser struct {
	Importer importer.FortiOSImporter
}

func NewFortiOSTextParser() FortiOSTextParser {
	return FortiOSTextParser{}
}

func (p *FortiOSTextParser) parseSections(section importer.FortiOSSection) {
	fmt.Println("parseSections", section)
	for _, cmd := range section.Commands {
		if cmd.Type == importer.FortiOSCommandTypeConfig {
			fmt.Fprintf(os.Stdout, "%s\n", cmd.Command)
		}

		for _, subSection := range section.Sections {
			fmt.Fprintf(os.Stdout, "%s\n", subSection)
			p.parseSections(subSection)
		}
	}
}

func (p FortiOSTextParser) Parse(imp importer.Importer) (types.VertoDevice, error) {
	fmt.Println("Parse", imp)
	d := types.VertoDevice{}
	p.Importer = imp.(importer.FortiOSImporter)
	fmt.Println("Parse", p.Importer)

	p.parseSections(p.Importer.Section)

	return d, nil
}
