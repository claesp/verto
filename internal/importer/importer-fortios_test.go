package importer

import (
	"fmt"
	"os"
	"testing"
)

func loadTestFile(filename string) ([]byte, error) {
	f, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return f, nil
}

func TestNewFortiOSImporter(t *testing.T) {
	s, e := loadTestFile("testdata/fortigate.conf")
	if e != nil {
		t.Fatalf("%v\n", e)
	}
	i := FortiOSImporter{}
	d := i.ImportFromText(string(s))
	fmt.Fprintf(os.Stdout, "%s\n", d)
}
