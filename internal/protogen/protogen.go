package protogen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen"
	"github.com/Wenchy/tableau/internal/excel"
	"github.com/Wenchy/tableau/internal/fs"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/pkg/errors"
)

const (
	version = "v0.1.0"
)

type Generator struct {
	ProtoPackage   string   // protobuf package name.
	GoPackage      string   // golang package name.
	InputDir       string   // input dir of workbooks.
	OutputDir      string   // output dir of generated protoconf files.
	FilenameSuffix string   // filename suffix of generated protoconf files.
	Imports        []string // imported common proto file paths

	Header *options.HeaderOption // header settings.
}

func (gen *Generator) Generate() error {
	if err := gen.PrepareOutpuDir(); err != nil {
		return errors.Wrapf(err, "failed to prepare output dir: %s", gen.OutputDir)
	}
	files, err := os.ReadDir(gen.InputDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read input dir: %s", gen.InputDir)
	}
	for _, wbFile := range files {
		if strings.HasPrefix(wbFile.Name(), "~$") {
			// ignore xlsx temp file named with prefix "~$"
			continue
		}
		wbPath := filepath.Join(gen.InputDir, wbFile.Name())
		atom.Log.Debugf("workbook: %s", wbPath)
		book, err := excel.NewBookExt(wbPath, confgen.NewSheetParser("tableau", ""))
		if err != nil {
			return errors.Wrapf(err, "failed to create new workbook: %s", wbPath)
		}

		if len(book.Sheets) == 0 {
			continue
		}
		// creat a book parser
		bp := newBookParser(wbFile.Name(), gen.Imports)
		for _, sheet := range book.Sheets {
			// parse sheet header
			sheetMsgName := sheet.Name
			if sheet.Meta.Alias != "" {
				sheetMsgName = sheet.Meta.Alias
			}
			ws := &tableaupb.Worksheet{
				Options: &tableaupb.WorksheetOptions{
					Name:      sheet.Name,
					Namerow:   gen.Header.Namerow,
					Typerow:   gen.Header.Typerow,
					Noterow:   gen.Header.Noterow,
					Datarow:   gen.Header.Datarow,
					Transpose: false,
					Tags:      "",
					Nameline:  sheet.Meta.Nameline,
					Typeline:  sheet.Meta.Typeline,
				},
				Fields: []*tableaupb.Field{},
				Name:   sheetMsgName,
			}
			shHeader := &sheetHeader{
				meta:    sheet.Meta,
				namerow: sheet.Rows[gen.Header.Namerow-1],
				typerow: sheet.Rows[gen.Header.Typerow-1],
				noterow: sheet.Rows[gen.Header.Noterow-1],
			}

			var ok bool
			for cursor := 0; cursor < len(shHeader.namerow); cursor++ {
				field := &tableaupb.Field{}
				cursor, ok = bp.parseField(field, shHeader, cursor, "")
				if ok {
					ws.Fields = append(ws.Fields, field)
				}
			}
			// append parsed sheet to workbook
			bp.wb.Worksheets = append(bp.wb.Worksheets, ws)
		}
		// export book
		be := newBookExporter(gen.ProtoPackage, gen.GoPackage, gen.OutputDir, gen.FilenameSuffix, gen.Imports, bp.wb)
		if err := be.export(); err != nil {
			return errors.Wrapf(err, "failed to export workbook: %s", wbPath)
		}
	}
	return nil
}

func (gen *Generator) PrepareOutpuDir() error {
	existed, err := fs.Exists(gen.OutputDir)
	if err != nil {
		return errors.Wrapf(err, "failed to check existence of output dir: %s", gen.OutputDir)
	}
	if existed {
		// remove all *.proto file but not Imports
		imports := make(map[string]int)
		for _, path := range gen.Imports {
			imports[path] = 1
		}
		files, err := os.ReadDir(gen.OutputDir)
		if err != nil {
			return errors.Wrapf(err, "failed to read dir: %s", gen.OutputDir)
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
				return errors.Wrapf(err, "failed to remove file: %s", fpath)
			}
		}

	} else {
		// create output dir
		err = os.MkdirAll(gen.OutputDir, 0700)
		if err != nil {
			return errors.Wrapf(err, "failed to create output dir: %s", gen.OutputDir)
		}
	}

	return nil
}

type sheetHeader struct {
	meta    *tableaupb.SheetMeta
	namerow []string
	typerow []string
	noterow []string
}

func getCell(row []string, cursor int, line int32) string {
	cell := row[cursor]
	return excel.ExtractFromCell(cell, line)
}

func (sh *sheetHeader) getNameCell(cursor int) string {
	return getCell(sh.namerow, cursor, sh.meta.Nameline)
}

func (sh *sheetHeader) getTypeCell(cursor int) string {
	return getCell(sh.typerow, cursor, sh.meta.Typeline)
}
func (sh *sheetHeader) getNoteCell(cursor int) string {
	return getCell(sh.noterow, cursor, 1) // default note line is 1
}

type GeneratedBuf struct {
	buf bytes.Buffer
}

// NewGeneratedFile creates a new generated file with the given filename.
func NewGeneratedBuf() *GeneratedBuf {
	return &GeneratedBuf{}
}

// P prints a line to the generated output. It converts each parameter to a
// string following the same rules as fmt.Print. It never inserts spaces
// between parameters.
func (g *GeneratedBuf) P(v ...interface{}) {
	for _, x := range v {
		fmt.Fprint(&g.buf, x)
	}
	fmt.Fprintln(&g.buf)
}

// Content returns the contents of the generated file.
func (g *GeneratedBuf) Content() []byte {
	return g.buf.Bytes()
}
