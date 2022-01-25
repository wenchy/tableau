package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

// metaSheetName defines the meta data of each worksheet.
const metaSheetName = "@TABLEAU"

type ExcelImporter struct {
	filename   string
	sheetMap   map[string]*Sheet // sheet name -> sheet
	sheetNames []string
	includeMetaSheet bool

	Meta       *tableaupb.WorkbookMeta
	metaParser SheetParser
}
// TODO: options
func NewExcelImporter(filename string, sheets []string, parser SheetParser, includeMetaSheet bool) *ExcelImporter {
	return &ExcelImporter{
		filename:   filename,
		sheetNames: sheets,
		includeMetaSheet: includeMetaSheet,
		metaParser: parser,
		Meta: &tableaupb.WorkbookMeta{
			SheetMetaMap: make(map[string]*tableaupb.SheetMeta),
		},
	}
}

func (x *ExcelImporter) GetSheets() ([]*Sheet, error) {
	if x.sheetMap == nil {
		if err := x.parse(); err != nil {
			return nil, errors.WithMessagef(err, "failed to parse %s", x.filename)
		}
	}

	sheets := []*Sheet{}
	for _, name := range x.sheetNames {
		sheet, err := x.GetSheet(name)
		if err != nil {
			return nil, errors.WithMessagef(err, "get sheet failed: %s", name)
		}
		sheets = append(sheets, sheet)
	}
	return sheets, nil
}

// GetSheet returns a Sheet of the specified sheet name.
func (x *ExcelImporter) GetSheet(name string) (*Sheet, error) {
	if x.sheetMap == nil {
		if err := x.parse(); err != nil {
			return nil, errors.WithMessagef(err, "failed to parse %s", x.filename)
		}
	}

	sheet, ok := x.sheetMap[name]
	if !ok {
		return nil, errors.Errorf("sheet %s not found", name)
	}
	return sheet, nil
}

func (x *ExcelImporter) parse() error {
	x.sheetMap = make(map[string]*Sheet)
	file, err := excelize.OpenFile(x.filename)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", x.filename)
	}

	if err := x.parseWorkbookMeta(file); err != nil {
		return errors.Wrapf(err, "failed to parse workbook meta: %s", metaSheetName)
	}

	if err := x.collectSheetsInOrder(file); err != nil {
		return errors.WithMessagef(err, "failed to collectSheetsInOrder: %s", x.filename)
	}

	for _, sheetName := range x.sheetNames {
		s, err := x.parseSheet(file, sheetName)
		if err != nil {
			return errors.WithMessagef(err, "failed to parse sheet: %s#%s", x.filename, sheetName)
		}
		x.sheetMap[sheetName] = s
	}
	return nil
}
func (x *ExcelImporter) NeedParseMeta() bool {
	return x.metaParser != nil
}

func (x *ExcelImporter) parseWorkbookMeta(file *excelize.File) error {
	if !x.NeedParseMeta() {
		atom.Log.Debugf("skip parsing workbook meta: %s", x.filename)
		return nil
	}

	if file.GetSheetIndex(metaSheetName) == -1 {
		atom.Log.Debugf("workbook %s has no sheet named %s", x.filename, metaSheetName)
		return nil
	}

	sheet, err := x.parseSheet(file, metaSheetName)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse sheet: %s#%s", x.filename, metaSheetName)
	}

	if sheet.MaxRow <= 1 {
		for _, sheetName := range file.GetSheetList() {
			x.Meta.SheetMetaMap[sheetName] = &tableaupb.SheetMeta{
				Sheet: sheetName,
			}
		}
		return nil
	}
	wsOpts := &tableaupb.WorksheetOptions{
		Name:    sheet.Name,
		Namerow: 1,
		Datarow: 2,
	}
	if err := x.metaParser.Parse(x.Meta, sheet, wsOpts); err != nil {
		return errors.WithMessagef(err, "failed to parse sheet: %s#%s", x.filename, metaSheetName)
	}

	atom.Log.Debugf("%s#%s: %+v", x.filename, metaSheetName, x.Meta)
	return nil
}

func (x *ExcelImporter) parseSheet(file *excelize.File, sheetName string) (*Sheet, error) {
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get rows of sheet: %s#%s", x.filename, sheetName)
	}
	sheet := NewSheet(sheetName, rows)

	if x.NeedParseMeta() {
		sheet.Meta = x.Meta.SheetMetaMap[sheetName]
	}
	return sheet, nil
}

func (x *ExcelImporter) collectSheetsInOrder(file *excelize.File) error {
	sortedMap := treemap.NewWithIntComparator()
	if x.NeedParseMeta() {
		for sheetName := range x.Meta.SheetMetaMap {
			index := file.GetSheetIndex(sheetName)
			if index == -1 {
				return errors.Errorf("sheet %s not found in workbook %s", sheetName, x.filename)
			}
			sortedMap.Put(index, sheetName)
		}
	} else {
		// Import all sheets except `@TABLEAU` if x.sheetNames is empty.
		if x.sheetNames == nil {
			for index, name := range file.GetSheetMap() {
				sortedMap.Put(index, name)
			}
		}

		for _, name := range x.sheetNames {
			index := file.GetSheetIndex(name)
			if index == -1 {
				return errors.Errorf("sheet %s not found in workbook %s", name, x.filename)
			}
			sortedMap.Put(index, name)
		}

	}

	// Clear before re-assign.
	x.sheetNames = nil
	for _, val := range sortedMap.Values() {
		sheetName := val.(string)
		if  sheetName != metaSheetName || (x.includeMetaSheet && sheetName == metaSheetName) {
			// exclude meta sheet
			x.sheetNames = append(x.sheetNames, sheetName)
		}
	}
	return nil
}

func (x *ExcelImporter) ExportCSV() error {
	ext := filepath.Ext(x.filename)
	basename := strings.TrimSuffix(x.filename, ext)
	sheets, err := x.GetSheets()
	if err != nil {
		return errors.WithMessagef(err, "failed to get sheets: %s", x.filename)
	}
	for _, sheet := range sheets {
		path := fmt.Sprintf("%s#%s.csv", basename, sheet.Name)
		f, err := os.Create(path)
		if err != nil {
			return errors.Wrapf(err, "failed to create csv file: %s", path)
		}
		defer f.Close()

		if err := sheet.ExportCSV(f); err != nil {
			return errors.WithMessagef(err, "export sheet %s to excel failed", sheet.Name)
		}
	}
	return nil
}

func OpenExcel(filename string, sheetName string) (*excelize.File, error) {
	var wb *excelize.File
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("create file: ", filename)
		wb = excelize.NewFile()
		t := time.Now()
		datetime := t.Format(time.RFC3339)
		err := wb.SetDocProps(&excelize.DocProperties{
			Category:       "category",
			ContentStatus:  "Draft",
			Created:        datetime,
			Creator:        "Tableau",
			Description:    "This file was created by Tableau",
			Identifier:     "xlsx",
			Keywords:       "Spreadsheet",
			LastModifiedBy: "Tableau",
			Modified:       datetime,
			Revision:       "0",
			Subject:        "Configuration",
			Title:          filepath.Base(filename),
			Language:       "en-US",
			Version:        "1.0.0",
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to set doc props: %s", filename)
		}
		// The newly created workbook will by default contain a worksheet named `Sheet1`.
		wb.SetSheetName("Sheet1", sheetName)
		wb.SetDefaultFont("Courier")
	} else {
		fmt.Println("existed file: ", filename)
		wb, err = excelize.OpenFile(filename)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file: %s", filename)
		}
		wb.NewSheet(sheetName)
	}
	return wb, nil
}
