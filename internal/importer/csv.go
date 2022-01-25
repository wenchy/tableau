package importer

import (
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
)

type CSVImporter struct {
	filename string
	sheets   map[string]*Sheet // sheet name -> sheet
}

func NewCSVImporter(filename string) *CSVImporter {
	return &CSVImporter{
		filename: filename,
	}
}

func (x *CSVImporter) GetSheets() ([]*Sheet, error){
	sheet, err := x.GetSheet(x.filename)
	if err != nil {
		return nil, errors.WithMessagef(err, "get sheet failed: %s", x.filename)
	}
	return []*Sheet{sheet}, nil
}

// GetSheet returns a Sheet of the specified sheet name.
func (x *CSVImporter) GetSheet(name string) (*Sheet, error) {
	if x.sheets == nil {
		if err := x.parse(); err != nil {
			return nil, errors.WithMessagef(err, "failed to parse %s", x.filename)
		}
	}

	sheet, ok := x.sheets[name]
	if !ok {
		return nil, errors.Errorf("sheet %s not found", name)
	}
	return sheet, nil
}

func (x *CSVImporter) parse() error {
	x.sheets = make(map[string]*Sheet)
	f, err := os.Open(x.filename)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", x.filename)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return errors.Wrapf(err, "failed to read file %s", x.filename)
	}

	// NOTE: For CSV, sheet name is the same as filename.
	x.sheets[x.filename] = NewSheet(x.filename, records)
	return nil
}
