package excel

import (
	"fmt"
	"math"
	"strings"
	"unicode"

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
		return nil, errors.WithMessagef(err, "failed to parse workbook: %s", filename)
	}
	return book, nil
}

func (b *Book) parse() error {
	file, err := excelize.OpenFile(b.Filename)
	if err != nil {
		return errors.WithMessagef(err, "failed to open workbook: %s", b.Filename)
	}
	for _, sheetName := range file.GetSheetList() {
		rows, err := file.GetRows(sheetName)
		if err != nil {
			return errors.WithMessagef(err, "failed to get rows of s: %s@%s", b.Filename, sheetName)
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

type RowCells struct {
	Row   int                 // row number
	cells map[string]*RowCell // name row -> data row cell
}

func NewRowCells(row int) *RowCells {
	return &RowCells{
		Row:   row,
		cells: make(map[string]*RowCell),
	}
}

type RowCell struct {
	Col  int    // colum number
	Data string // cell data
}

func (r *RowCells) Cell(name string) *RowCell {
	return r.cells[name]
}

func (r *RowCells) CellDebugString(name string) string {
	rc := r.Cell(name)
	if rc == nil {
		return fmt.Sprintf("(%d,%d)%s:%s", r.Row, 0, name, "")
	}
	return fmt.Sprintf("(%d,%d)%s:%s", r.Row, rc.Col, name, rc.Data)
}

func (r *RowCells) SetCell(name string, col int, data string) {
	r.cells[name] = &RowCell{
		Col:  col,
		Data: data,
	}
}

func (r *RowCells) GetCellCountWithPrefix(prefix string) int {
	// atom.Log.Debug("name prefix: ", prefix)
	size := 0
	for name := range r.cells {
		if strings.HasPrefix(name, prefix) {
			num := 0
			// atom.Log.Debug("name: ", name)
			colSuffix := name[len(prefix):]
			// atom.Log.Debug("name: suffix ", colSuffix)
			for _, r := range colSuffix {
				if unicode.IsDigit(r) {
					num = num*10 + int(r-'0')
				} else {
					break
				}
			}
			size = int(math.Max(float64(size), float64(num)))
		}
	}
	return size
}
