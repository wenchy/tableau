package excel

import (
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type Sheet struct {
	Name   string
	MaxRow int
	MaxCol int

	Rows [][]string
}

func (s *Sheet) Cell(row, col int) (string, error) {
	if row < 0 || row >= s.MaxRow {
		return "", errors.Errorf("row %d out of range", row)
	}
	if col < 0 || col >= s.MaxCol {
		return "", errors.Errorf("col %d out of range", col)
	}
	// MOTE: different row may have different length.
	if col >= len(s.Rows[row]) {
		return "", nil
	}
	return s.Rows[row][col], nil
}

func (s *Sheet) String() string {
	str := ""
	for row := 0; row < s.MaxRow; row++ {
		for col := 0; col < s.MaxCol; col++ {
			cell, _ := s.Cell(row, col)
			str += cell + "|"
		}
		str += "\n"
	}
	return str
}

type Book struct {
	Filename string
	file     *excelize.File
	Sheets   map[string]*Sheet // s name -> Sheet
}

func NewBook(filename string) (*Book, error) {
	book := &Book{
		Filename: filename,
		Sheets:   make(map[string]*Sheet),
	}
	err := book.parse()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse workbook: %s", filename)
	}
	return book, nil
}

func (b *Book) parse() error {
	file, err := excelize.OpenFile(b.Filename)
	if err != nil {
		return errors.Wrapf(err, "failed to open workbook: %s", b.Filename)
	}
	for _, sheetName := range file.GetSheetList() {
		rows, err := file.GetRows(sheetName)
		if err != nil {
			return errors.Wrapf(err, "failed to get rows of s: %s@%s", b.Filename, sheetName)
		}
		maxRow := len(rows)
		maxCol := 0
		// MOTE: different row may have different length.
		// We need to find the max col.
		for _, row := range rows {
			n := len(row)
			if n > maxCol {
				maxCol = n
			}
		}
		s := &Sheet{
			Name:   sheetName,
			MaxRow: maxRow,
			MaxCol: maxCol,
			Rows:   rows,
		}
		b.Sheets[sheetName] = s
	}
	return nil
}
