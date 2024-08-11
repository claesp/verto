package importer

import (
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
	s, err := loadTestFile("testdata/fortigate-cgustave.conf")
	if err != nil {
		t.Fatal(err)
	}

	i := NewFortiOSImporter()
	err = i.Import(string(s))
	if err != nil {
		t.Fatal(err)
	}
}
