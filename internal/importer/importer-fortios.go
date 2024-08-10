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
	return fmt.Sprintf("%6d: %s", r.RowNumber, r.Cols)
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

	for _, row := range s.Rows {
		o = fmt.Sprintf("%s\n%s: %s", o, fmt.Sprintf("%s/%s", s.Parent.Name, s.Name), row)
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
}

func (f FortiOSImporter) parseRows(inRows []string) []FortiOSRow {
	outRows := make([]FortiOSRow, 0)

	for idx := 0; idx < len(inRows); idx++ {
		inRow := inRows[idx]
		outRow := FortiOSRow{
			RowNumber: idx + 1,
			Data:      inRow,
			Cols:      strings.Split(strings.TrimLeft(inRow, " "), " "),
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
)

func NewFortiOSCommand(row FortiOSRow) FortiOSCommand {
	o := FortiOSCommand{}
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
	default:
		o.Type = FortiOSCommandTypeUnknown
	}

	if len(row.Cols) > 1 {
		o.Command = row.Cols[1]
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

func (f FortiOSImporter) ImportFromText(s string) types.VertoDevice {
	rows := strings.Split(s, "\n")
	osRows := f.parseRows(rows)
	firstSection := NewFortiOSSection(osRows[0], &FortiOSSection{Name: "root"})
	f.parseSections(osRows[1:], &firstSection)
	fmt.Fprintf(os.Stdout, "%s\n", firstSection)

	d := types.VertoDevice{}

	return d
}
