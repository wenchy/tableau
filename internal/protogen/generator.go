package protogen

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/fs"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/golang/protobuf/proto"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
)

var mapRegexp *regexp.Regexp
var listRegexp *regexp.Regexp
var structRegexp *regexp.Regexp

const (
	tableauProtoPath   = "tableau/protobuf/tableau.proto"
	timestampProtoPath = "google/protobuf/timestamp.proto"
	durationProtoPath  = "google/protobuf/duration.proto"
)

func init() {
	mapRegexp = regexp.MustCompile(`^map<(.+),(.+)>`)  // e.g.: map<uint32,Element>
	listRegexp = regexp.MustCompile(`^\[(.*)\](.+)`)   // e.g.: [Element]uint32
	structRegexp = regexp.MustCompile(`^\{(.+)\}(.+)`) // e.g.: {Element}uint32
}

type Generator struct {
	ProtoPackage string   // protobuf package name.
	GoPackage    string   // golang package name.
	InputDir     string   // input dir of workbooks.
	OutputDir    string   // output dir of generated protoconf files.
	Imports      []string // imported common proto file paths

	Header *options.HeaderOption // header settings.
}

func (gen *Generator) Generate() {

	if existed, err := fs.Exists(gen.OutputDir); err != nil {
		atom.Log.Panic(err)
	} else {
		if existed {
			// remove all *.proto file but not Imports
			imports := make(map[string]int)
			for _, path := range gen.Imports {
				imports[path] = 1
			}
			files, err := os.ReadDir(gen.OutputDir)
			if err != nil {
				atom.Log.Panic(err)
			}
			for _, file := range files {
				if !strings.HasSuffix(file.Name(), ".proto") {
					continue
				}
				if _, ok := imports[file.Name()]; ok {
					continue
				}
				fpath := filepath.Join(gen.OutputDir, file.Name())
				err := os.Remove(fpath)
				if err != nil {
					atom.Log.Panic(err)
				}
			}

		} else {
			// create output dir
			err = os.MkdirAll(gen.OutputDir, 0700)
			if err != nil {
				atom.Log.Panic(err)
			}
		}
	}

	files, err := os.ReadDir(gen.InputDir)
	if err != nil {
		atom.Log.Panic(err)
	}
	for _, wbFile := range files {
		if strings.HasPrefix(wbFile.Name(), "~$") {
			// ignore xlsx temp file named with prefix "~$"
			continue
		}
		wbPath := filepath.Join(gen.InputDir, wbFile.Name())
		atom.Log.Debugf("workbook: %s", wbPath)
		f, err := excelize.OpenFile(wbPath)
		if err != nil {
			atom.Log.Panic(err)
			return
		}
		wbProtoName := strcase.ToSnake(strings.TrimSuffix(wbFile.Name(), filepath.Ext(wbFile.Name())))
		book := &book{
			wb: &tableaupb.Workbook{
				Options: &tableaupb.WorkbookOptions{
					Name: wbFile.Name(),
				},
				Worksheets: []*tableaupb.Worksheet{},
				Name:       wbProtoName,
				Imports: map[string]int32{
					tableauProtoPath: 1, // default import
				},
			},
			withNote: false,
		}

		for _, path := range gen.Imports {
			book.wb.Imports[path] = 1 // custom imports
		}

		for _, sheetName := range f.GetSheetList() {
			rows, err := f.GetRows(sheetName)
			if err != nil {
				atom.Log.Panic(err)
			}
			ws := &tableaupb.Worksheet{
				Options: &tableaupb.WorksheetOptions{
					Name:      sheetName,
					Namerow:   gen.Header.Namerow,
					Typerow:   gen.Header.Typerow,
					Noterow:   gen.Header.Noterow,
					Datarow:   gen.Header.Datarow,
					Transpose: false,
					Tags:      "",
				},
				Fields: []*tableaupb.Field{},
				Name:   sheetName,
			}
			header := &sheetHeader{
				namerow: rows[0],
				typerow: rows[1],
				noterow: rows[2],
			}

			var ok bool
			for cursor := 0; cursor < len(header.namerow); cursor++ {
				field := &tableaupb.Field{}
				cursor, ok = book.parseField(field, header, cursor, "")
				if ok {
					ws.Fields = append(ws.Fields, field)
				}
			}
			book.wb.Worksheets = append(book.wb.Worksheets, ws)
		}
		if err := gen.exportWorkbook(book.wb); err != nil {
			atom.Log.Panic(err)
		}
	}
}

func (gen *Generator) exportWorkbook(wb *tableaupb.Workbook) error {
	atom.Log.Debug(proto.MarshalTextString(wb))
	path := filepath.Join(gen.OutputDir, wb.Name+".proto")
	atom.Log.Debugf("output: %s", path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := NewWriter(f)
	defer w.Flush()

	w.P("syntax = \"proto3\";")
	w.P("package %s;", gen.ProtoPackage)
	w.P("option go_package = \"%s\";", gen.GoPackage)
	w.P("")

	// keep the elements ordered by sheet name
	set := treeset.NewWithStringComparator()
	for key := range wb.Imports {
		set.Add(key)
	}
	for _, key := range set.Values() {
		w.P("import \"%s\";", key)
	}
	w.P("")
	w.P("option (tableau.workbook) = {%s};", genPrototext(wb.Options))
	w.P("")

	for i, ws := range wb.Worksheets {
		exporter := &sheetExporter{
			ws:             ws,
			w:              w,
			isLastSheet:    i == len(wb.Worksheets)-1,
			nestedMessages: make(map[string]*tableaupb.Field),
		}
		if err := exporter.export(); err != nil {
			return err
		}
	}

	return nil
}
