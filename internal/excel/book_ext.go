package excel

import (
	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/proto"
)

type Parser interface {
	Parse(protomsg proto.Message, sheet *Sheet, wsOpts *tableaupb.WorksheetOptions) error
}

// tableauSheetName defines the meta data of each worksheet.
const tableauSheetName = "@TABLEAU"

type SheetExt struct {
	*Sheet

	Meta *tableaupb.SheetMeta
}

type BookExt struct {
	Filename string
	file     *excelize.File
	Sheets   []*SheetExt
	Meta     *tableaupb.WorkbookMeta
	parser   Parser
}

func NewBookExt(filename string, parser Parser) (*BookExt, error) {
	book := &BookExt{
		Filename: filename,
		Meta: &tableaupb.WorkbookMeta{
			SheetMetaMap: make(map[string]*tableaupb.SheetMeta),
		},
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

func (b *BookExt) GetSheet(sheetName string) *SheetExt {
	for _, s := range b.Sheets {
		if s.Name == sheetName {
			return s
		}
	}
	return nil
}

func (b *BookExt) parseSheet(sheetName string) (*SheetExt, error) {
	s, err := parseSheet(b.file, sheetName)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to parse sheet: %s#%s", b.Filename, sheetName)
	}
	sheet := &SheetExt{
		Sheet: s,
	}

	// Special `Meta` field process for "@TABLEAU" sheet.
	if sheetName == tableauSheetName {
		sheet.Meta = &tableaupb.SheetMeta{
			Sheet: sheetName,
		}
	} else {
		sheet.Meta = b.Meta.SheetMetaMap[sheetName]
	}
	return sheet, nil
}

func (b *BookExt) parseWorkbookMeta() error {
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
				b.Meta.SheetMetaMap[sheetName] = &tableaupb.SheetMeta{
					Sheet: sheetName,
				}
			}
		}
		return nil
	}
	wsOpts := &tableaupb.WorksheetOptions{
		Name: sheet.Name,
		Namerow: 1,
		Datarow: 2,
	}
	if err := b.parser.Parse(b.Meta, sheet.Sheet, wsOpts); err != nil {
		return errors.WithMessagef(err, "failed to parse sheet: %s#%s", b.Filename, tableauSheetName)
	}

	atom.Log.Debugf("%s#%s: %+v", b.Filename, tableauSheetName, b.Meta)
	return nil
}

func (b *BookExt) parse() error {
	// Target: keep the order of sheets.
	set := treeset.NewWithStringComparator()
	for sheetName := range b.Meta.SheetMetaMap {
		set.Add(sheetName) // default must import path
	}

	for _, val := range set.Values() {
		sheetName := val.(string)
		s, err := b.parseSheet(sheetName)
		if err != nil {
			return errors.WithMessagef(err, "failed to get rows of s: %s#%s", b.Filename, sheetName)
		}
		b.Sheets = append(b.Sheets, s)
	}
	return nil
}
