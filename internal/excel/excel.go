package excel

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
)

type Parser interface {
	Parse(protomsg proto.Message, sheet *Sheet, namerow, datarow int32, transpose bool) error
}

// tableauSheetName defines the meta data of each worksheet.
const tableauSheetName = "@TABLEAU"

type Sheet struct {
	Name   string
	MaxRow int
	MaxCol int

	Rows [][]string

	Meta *tableaupb.SheetMeta
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
	Filename     string
	file         *excelize.File
	WorkbookMeta *tableaupb.WorkbookMeta
	Sheets       map[string]*Sheet // sheet name -> Sheet
	parser       Parser
}

func NewBook(filename string, parser Parser) (*Book, error) {
	book := &Book{
		Filename: filename,
		WorkbookMeta: &tableaupb.WorkbookMeta{
			SheetMetaMap: make(map[string]*tableaupb.SheetMeta),
		},
		Sheets: make(map[string]*Sheet),
		parser: parser,
	}

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open workbook: %s", filename)
	}
	book.file = file

	if err := book.parseWorkbookMeta(); err != nil {
		return nil, errors.Wrapf(err, "failed to parse tableau sheet: %s", tableauSheetName)
	}

	if err := book.parse(); err != nil {
		return nil, errors.WithMessagef(err, "failed to parse workbook: %s", filename)
	}
	return book, nil
}

func (b *Book) parseWorkbookMeta() error {
	if b.file.GetSheetIndex(tableauSheetName) == -1 {
		atom.Log.Debugf("workbook %s has no sheet named %s", b.Filename, tableauSheetName)
		return nil
	}

	sheet, err := b.parseSheet(tableauSheetName)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse sheet: %s#%s", b.Filename, tableauSheetName)
	}

	if sheet.MaxRow <= 1 {
		for _, sheetName := range b.file.GetSheetList() {
			if sheetName != tableauSheetName {
				b.WorkbookMeta.SheetMetaMap[sheetName] = &tableaupb.SheetMeta{
					Sheet: sheetName,
				}
			}
		}
		return nil
	}
	if err := b.parser.Parse(b.WorkbookMeta, sheet, 1, 2, false); err != nil {
		return errors.WithMessagef(err, "failed to parse sheet: %s#%s", b.Filename, tableauSheetName)
	}
	atom.Log.Debugf("%s#%s: %+v", b.Filename, tableauSheetName, b.WorkbookMeta)
	return nil
}

func (b *Book) parse() error {
	for _, sheetMeta := range b.WorkbookMeta.SheetMetaMap {
		sheetName := sheetMeta.Sheet
		s, err := b.parseSheet(sheetName)
		if err != nil {
			return errors.WithMessagef(err, "failed to get rows of s: %s#%s", b.Filename, sheetName)
		}
		b.Sheets[sheetName] = s
	}
	return nil
}

func (b *Book) parseSheet(sheetName string) (*Sheet, error) {
	rows, err := b.file.GetRows(sheetName)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get rows of sheet: %s#%s", b.Filename, sheetName)
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
		Meta:   b.WorkbookMeta.SheetMetaMap[sheetName],
	}
	return s, nil
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
		return fmt.Sprintf("(%d,%d)%s:%s", r.Row+1, 1, name, "")
	}
	return fmt.Sprintf("(%d,%d)%s:%s", r.Row+1, rc.Col+1, name, rc.Data)
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
