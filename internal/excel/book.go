package excel

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/camelcase"
	"github.com/Wenchy/tableau/internal/types"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

var newlineRegex *regexp.Regexp

func init() {
	newlineRegex = regexp.MustCompile(`\r?\n?`)
}

func clearNewline(s string) string {
	return newlineRegex.ReplaceAllString(s, "")
}

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
	Sheets   map[string]*Sheet // sheet name -> sheet
}

func NewBook(filename string, sheets []string) (*Book, error) {
	b := &Book{
		Filename: filename,
		Sheets:   make(map[string]*Sheet),
	}

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open workbook: %s", filename)
	}

	for _, sheetName := range sheets {
		s, err := parseSheet(file, sheetName)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to get rows of s: %s#%s", b.Filename, sheetName)
		}
		b.Sheets[sheetName] = s
	}
	return b, nil
}

func parseSheet(file *excelize.File, sheetName string) (*Sheet, error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to get rows of sheet: %s", sheetName)
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
	return s, nil
}

type RowCells struct {
	// This previous row cells is for auto-populating the currernt row's missing data.
	// As the user doesn't fill the duplicate map key for easy use and clear reading.
	//
	// ServerName			ServerConfName
	// map<string, Server>	[Conf]string
	//
	// gamesvr				HeadFrameConf
	// activitysvr			ActivityConf
	// *MISSING-KEY*		ChapterConf
	// *MISSING-KEY*		CollectionConf

	prev *RowCells

	Row          int                 // row number
	cells        map[string]*RowCell // name -> RowCell
	indexedCells map[int]*RowCell    // column index -> RowCell
}

func NewRowCells(row int, prev *RowCells) *RowCells {
	return &RowCells{
		prev: prev,

		Row:          row,
		cells:        make(map[string]*RowCell),
		indexedCells: make(map[int]*RowCell),
	}
}

type RowCell struct {
	Col           int    // cell column (0-based)
	Data          string // cell data
	Type          string // cell type
	Name          string // cell name
	autoPopulated bool   // auto-populated
}

func (r *RowCells) Cell(name string, optional bool) *RowCell {
	c := r.cells[name]
	if c == nil && optional {
		// if optional, return an empty cell.
		c = &RowCell{
			Col:  -1,
			Data: "",
		}
	}
	return c
}

func (r *RowCells) CellDebugString(name string) string {
	rc := r.Cell(name, false)
	if rc == nil {
		return fmt.Sprintf("(%d,%d)%s:%s", r.Row+1, -1, name, "")
	}
	dataFlag := ""
	if rc.autoPopulated {
		dataFlag = "~"
	}
	return fmt.Sprintf("(%d,%d)%s:%s%s", r.Row+1, rc.Col+1, name, rc.Data, dataFlag)
}

func (r *RowCells) SetCell(name string, col int, data, typ string) {
	cell := &RowCell{
		Col:  col,
		Data: data,
		Type: typ,
		Name: name,
	}

	// TODO: Parser(first-pass), check if this sheet is nested.
	if data == "" {
		if (types.MatchMap(typ) != nil || types.MatchKeyedList(typ) != nil) && r.prev != nil {
			// NOTE: populate the missing map key from the prev row's corresponding cell.
			// TODO(wenchy): this is a flawed hack, need to be taken into more consideration.
			// Check: reverse backward to find the previous same nested-level keyed cell and
			// compare them to make sure they are the same.
			prefix := ""
			splits := camelcase.Split(name)
			if len(splits) >= 2 {
				prefix = strings.Join(splits[:len(splits)-2], "")
			}
			needPopulate := false
			if prefix == "" {
				needPopulate = true
			} else {
				for i := col - 1; i >= 0; i-- {
					// prevData := r.prev.indexedCells[col].Data
					backCell := r.indexedCells[i]
					if !strings.HasPrefix(backCell.Name, prefix) {
						break
					}
					if types.MatchMap(backCell.Type) != nil || types.MatchKeyedList(backCell.Type) != nil {
						if r.prev.indexedCells[i].Data == r.indexedCells[i].Data {
							needPopulate = true
							break
						}
					}
				}
			}

			if needPopulate {
				if prevCell := r.prev.Cell(name, false); prevCell != nil {
					cell.Data = prevCell.Data
					cell.autoPopulated = true
				} else {
					atom.Log.Errorf("failed to find prev cell for name: %s, row: %d", name, r.Row)
				}
			}
		}
	}

	// add new cell
	r.cells[name] = cell
	r.indexedCells[col] = cell
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

func ExtractFromCell(cell string, line int32) string {
	if line == 0 {
		// line 0 means the whole cell.
		return clearNewline(strings.TrimSpace(cell))
	}

	lines := strings.Split(cell, "\n")
	if int32(len(lines)) >= line {
		return strings.TrimSpace(lines[line-1])
	}
	// atom.Log.Debugf("No enough lines in cell: %s, want at least %d lines", cell, line)
	return ""
}
