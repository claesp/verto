package importer

import (
	"fmt"
	"github.com/claesp/verto/internal/types"
	"os"
	"strings"
)

type FortiOSRow struct {
	RowNumber int
	Data      string
	Cols      []string
}

func (r FortiOSRow) String() string {
	return fmt.Sprintf("%6d: '%s'", r.RowNumber, r.Data)
}

type FortiOSSection struct {
	Name     string
	Rows     []FortiOSRow
	Sections []FortiOSSection
	Parent   *FortiOSSection
	Commands []FortiOSCommand
}

func NewFortiOSSection(row FortiOSRow, parent *FortiOSSection) FortiOSSection {
	o := FortiOSSection{}
	o.Name = strings.Join(row.Cols, " ")
	o.Rows = make([]FortiOSRow, 0)
	o.Rows = append(o.Rows, row)
	o.Sections = make([]FortiOSSection, 0)
	o.Parent = parent
	o.Commands = make([]FortiOSCommand, 0)
	o.Commands = append(o.Commands, NewFortiOSCommand(row))

	return o
}

func (s FortiOSSection) String() string {
	var o string
	for _, cmd := range s.Commands {
		o = fmt.Sprintf("%s\n%s: %s", o, fmt.Sprintf("%s/%s", s.Parent.Name, s.Name), cmd)
	}

	for _, sec := range s.Sections {
		o = fmt.Sprintf("%s\n%s", o, sec)
	}

	return o
}

func (s *FortiOSSection) RowCount() int {
	r := 0

	r += len(s.Rows)

	for _, u := range s.Sections {
		r += u.RowCount()
	}

	return r
}

type FortiOSImporter struct {
	Section FortiOSSection
	Device  types.VertoDevice
}

func NewFortiOSImporter() FortiOSImporter {
	return NewFortiOSImporterWithSection(NewFortiOSSection(FortiOSRow{Data: "root"}, &FortiOSSection{Name: "root", Parent: &FortiOSSection{Name: "-"}}))
}

func NewFortiOSImporterWithSection(section FortiOSSection) FortiOSImporter {
	o := FortiOSImporter{}
	o.Section = section

	return o
}

func (f FortiOSImporter) parseRows(inRows []string) []FortiOSRow {
	outRows := make([]FortiOSRow, 0)

	for idx := 0; idx < len(inRows); idx++ {
		inRow := inRows[idx]
		cols := strings.Split(strings.TrimLeft(inRow, " "), " ")
		first := cols[0]
		rest := strings.Join(cols[1:], " ")
		outRow := FortiOSRow{
			RowNumber: idx + 1,
			Data:      inRow,
			Cols:      []string{first, rest},
		}
		outRows = append(outRows, outRow)
	}

	return outRows
}

type FortiOSCommand struct {
	Type    FortiOSCommandType
	Command string
	Row     *FortiOSRow
}

func (c FortiOSCommand) String() string {
	if c.Row != nil {
		return fmt.Sprintf("%s: %s (%d)", c.Type, c.Command, c.Row.RowNumber)
	}

	return fmt.Sprintf("%s: %s", c.Type, c.Command)
}

type FortiOSCommandType int

func (t FortiOSCommandType) String() string {
	switch t {
	case FortiOSCommandTypeUnknown:
		return "UNKNOWN"
	case FortiOSCommandTypeSet:
		return "SET"
	case FortiOSCommandTypeNext:
		return "NEXT"
	case FortiOSCommandTypeEnd:
		return "END"
	case FortiOSCommandTypeEdit:
		return "EDIT"
	case FortiOSCommandTypeConfig:
		return "CONFIG"
	case FortiOSCommandTypeComment:
		return "COMMENT"
	case FortiOSCommandTypeUnset:
		return "UNSET"
	default:
		return "UNDEFINED"
	}
}

const (
	FortiOSCommandTypeUnknown = iota
	FortiOSCommandTypeSet
	FortiOSCommandTypeNext
	FortiOSCommandTypeEnd
	FortiOSCommandTypeEdit
	FortiOSCommandTypeConfig
	FortiOSCommandTypeComment
	FortiOSCommandTypeUnset
)

func NewFortiOSCommand(row FortiOSRow) FortiOSCommand {
	o := FortiOSCommand{}
	if len(row.Cols) != 0 {
		switch row.Cols[0] {
		case "set":
			o.Type = FortiOSCommandTypeSet
		case "next":
			o.Type = FortiOSCommandTypeNext
		case "end":
			o.Type = FortiOSCommandTypeEnd
		case "edit":
			o.Type = FortiOSCommandTypeEdit
		case "config":
			o.Type = FortiOSCommandTypeConfig
		case "unset":
			o.Type = FortiOSCommandTypeUnset
		default:
			o.Type = FortiOSCommandTypeUnknown
		}
	} else {
		o.Type = FortiOSCommandTypeUnknown
	}

	if len(row.Cols) > 0 {
		o.Command = strings.Join(row.Cols, " ")
	}

	if o.Type == FortiOSCommandTypeUnknown {
		if len(o.Command) > 0 && o.Command[0] == '#' {
			o.Type = FortiOSCommandTypeComment
		}
	}

	o.Row = &row

	return o
}

func (f FortiOSImporter) parseCommand(row FortiOSRow) FortiOSCommand {
	return NewFortiOSCommand(row)
}

func (f FortiOSImporter) parseSections(rows []FortiOSRow, section *FortiOSSection) {
	for i := 0; i < len(rows); i++ {
		row := rows[i]

		if row.Cols[0] == "config" || row.Cols[0] == "edit" {
			subSection := NewFortiOSSection(row, section)
			f.parseSections(rows[i+1:], &subSection)
			section.Sections = append(section.Sections, subSection)
			i = section.RowCount() - 2
			continue
		}

		section.Rows = append(section.Rows, row)
		command := f.parseCommand(row)
		section.Commands = append(section.Commands, command)

		if row.Cols[0] == "end" || row.Cols[0] == "next" {
			break
		}
	}
}

func (f FortiOSImporter) Import(s string) error {
	rows := f.parseRows(strings.Split(s, "\n"))
	f.parseSections(rows[1:], &f.Section)
	fmt.Fprintf(os.Stdout, "%s", f.Section)

	return nil
}

func (f FortiOSImporter) Parse() error {
	return fmt.Errorf("not implemented")
}

func (f FortiOSImporter) ExtractDevice() types.VertoDevice {
	return f.Device
}
