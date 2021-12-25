package protogen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen"
	"github.com/Wenchy/tableau/internal/fs"
	"github.com/Wenchy/tableau/internal/importer"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/antchfx/xmlquery"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

const (
	app                 = "tableauc"
	version             = "0.1.1"
	TableauProtoPackage = "tableau"
)

type IGenerator interface {
	GetSuffix() string // get filename extension
}

type Generator struct {
	ProtoPackage   string   // protobuf package name.
	GoPackage      string   // golang package name.
	// Location represents the collection of time offsets in use in a geographical area.
	// Default is "Asia/Shanghai".
	LocationName string
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated protoconf files.

	FilenameWithSubdirPrefix bool   // filename dir separator `/` or `\` is replaced by "__"
	FilenameSuffix           string // filename suffix of generated protoconf files.

	Imports        []string // imported common proto file paths

	i IGenerator
}

func (gen Generator) GetSuffix() string {
	return gen.i.GetSuffix()
}

func newGenerator(protoPackage, goPackage, indir, outdir string, setters ...options.Option) *Generator {
	opts := options.ParseOptions(setters...)
	return &Generator{
		ProtoPackage:   protoPackage,
		GoPackage:      goPackage,
		LocationName: opts.LocationName,
		InputDir:     indir,
		OutputDir:    outdir,

		FilenameWithSubdirPrefix: opts.Output.FilenameWithSubdirPrefix,
		FilenameSuffix:           opts.Output.FilenameSuffix,

		Imports: opts.Imports,
	}
}

type XlsxGenerator struct {
	Generator

	Header *options.HeaderOption // header settings.
}

func NewXlsxGenerator(protoPackage, goPackage, indir, outdir string, setters ...options.Option) *XlsxGenerator {
	opts := options.ParseOptions(setters...)
	g := &XlsxGenerator{
		Generator: *newGenerator(protoPackage, goPackage, indir, outdir, setters...),
		Header:         opts.Header,
	}
	g.i = g
	return g
}

func (gen XlsxGenerator) GetSuffix() string {
	return ".xlsx"
}

func (gen *XlsxGenerator) Generate() error {
	if err := gen.PrepareOutpuDir(); err != nil {
		return errors.Wrapf(err, "failed to prepare output dir: %s", gen.OutputDir)
	}
	return gen.generate(gen.InputDir)
}

func (gen *Generator) GenOneWorkbook(relativeWorkbookPath string) error {
	absPath := filepath.Join(gen.InputDir, relativeWorkbookPath)
	return gen.convertWorkbook(filepath.Dir(absPath), filepath.Base(absPath))
}

func (gen *Generator) generate(dir string) error {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return errors.Wrapf(err, "failed to read input dir: %s", gen.InputDir)
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			// scan and generate subdir recursively
			subdir := filepath.Join(dir, entry.Name())
			err := gen.generate(subdir)
			if err != nil {
				return errors.WithMessagef(err, "failed to generate subdir: %s", subdir)
			}
			continue
		}

		if strings.HasPrefix(entry.Name(), "~$") {
			// ignore xlsx temp file named with prefix "~$"
			continue
		}
		// atom.Log.Debugf("generating %s, %s", entry.Name(), filepath.Ext(entry.Name()))
		if filepath.Ext(entry.Name()) != gen.GetSuffix() {
			// ignore not xlsx files
			continue
		}
		if err := gen.convertWorkbook(dir, entry.Name()); err != nil {
			return errors.WithMessage(err, "failed to convert workbook")
		}
	}
	return nil
}

func getRelCleanSlashPath(rootdir, dir, filename string) (string, error) {
	relativeDir, err := filepath.Rel(rootdir, dir)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get relative path from %s to %s", rootdir, dir)
	}
	// relative slash separated path
	relativePath := filepath.Join(relativeDir, filename)
	relSlashPath := filepath.ToSlash(filepath.Clean(relativePath))
	return relSlashPath, nil
}

func (gen *Generator) convertWorkbook(dir, filename string) error {
	relativePath, err := getRelCleanSlashPath(gen.InputDir, dir, filename)
	if err != nil {
		return errors.WithMessagef(err, "get relative path failed")
	}
	absPath := filepath.Join(dir, filename)
	parser := confgen.NewSheetParser(TableauProtoPackage, gen.LocationName)
	imp := importer.New(absPath, importer.Parser(parser))
	sheets, err := imp.GetSheets()
	if err != nil {
		return errors.Wrapf(err, "failed to get sheet of workbook: %s", absPath)
	}
	if len(sheets) == 0 {
		return nil
	}
	atom.Log.Infof("workbook: %s, %s", gen.InputDir, relativePath)
	// creat a book parser
	bp := newBookParser(relativePath, gen.FilenameWithSubdirPrefix, gen.Imports)
	for _, sheet := range sheets {
		// parse sheet header
		sheetMsgName := sheet.Name
		if sheet.Meta.Alias != "" {
			sheetMsgName = sheet.Meta.Alias
		}
		// merge nameline and typeline from sheet meta and header option
		if sheet.Meta.Nameline == 0 {
			sheet.Meta.Nameline = gen.Header.Nameline
		}
		if sheet.Meta.Typeline == 0 {
			sheet.Meta.Typeline = gen.Header.Typeline
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
				Nested:    sheet.Meta.Nested,
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
			cursor, ok = bp.parseField(field, shHeader, cursor, "", sheet.Meta.Nested)
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
		return errors.WithMessagef(err, "failed to export workbook: %s", relativePath)
	}

	return nil
}

// func (gen *Generator) PrepareOutpuDir() error {
// 	existed, err := fs.Exists(gen.OutputDir)
// 	if err != nil {
// 		return errors.Wrapf(err, "failed to check existence of output dir: %s", gen.OutputDir)
// 	}
// 	if existed {
// 		// remove all *.proto file but not Imports
// 		imports := make(map[string]int)
// 		for _, path := range gen.Imports {
// 			imports[path] = 1
// 		}
// 		files, err := os.ReadDir(gen.OutputDir)
// 		if err != nil {
// 			return errors.Wrapf(err, "failed to read dir: %s", gen.OutputDir)
// 		}
// 		for _, file := range files {
// 			if !strings.HasSuffix(file.Name(), ".proto") {
// 				continue
// 			}
// 			if _, ok := imports[file.Name()]; ok {
// 				continue
// 			}
// 			fpath := filepath.Join(gen.OutputDir, file.Name())
// 			err := os.Remove(fpath)
// 			if err != nil {
// 				return errors.Wrapf(err, "failed to remove file: %s", fpath)
// 			}
// 		}

// 	} else {
// 		// create output dir
// 		err = os.MkdirAll(gen.OutputDir, 0700)
// 		if err != nil {
// 			return errors.Wrapf(err, "failed to create output dir: %s", gen.OutputDir)
// 		}
// 	}

// 	return nil
// }

func (gen *Generator) PrepareOutpuDir() error {
	existed, err := fs.Exists(gen.OutputDir)
	if err != nil {
		return errors.Wrapf(err, "failed to check existence of output dir: %s", gen.OutputDir)
	}
	if existed {
		files, err := os.ReadDir(gen.InputDir)
		if err != nil {
			return errors.Wrapf(err, "failed to read input dir: %s", gen.InputDir)
		}
		fileMap := make(map[string]bool)
		for _, wbFile := range files {
			if strings.HasPrefix(wbFile.Name(), "~$") || !strings.HasSuffix(wbFile.Name(), gen.GetSuffix()) {
				// ignore xlsx temp file named with prefix "~$"
				continue
			}
			fileMap[strcase.ToSnake(strings.ReplaceAll(wbFile.Name(), gen.GetSuffix(), ""))] = true
		}
		files, err = os.ReadDir(gen.OutputDir)
		if err != nil {
			return errors.Wrapf(err, "failed to read dir: %s", gen.OutputDir)
		}
		for _, file := range files {
			if _, existed := fileMap[strings.ReplaceAll(file.Name(), fmt.Sprintf("%s.proto", gen.FilenameSuffix), "")]; !existed {
				continue
			}
			fpath := filepath.Join(gen.OutputDir, file.Name())
			atom.Log.Debug(fpath)
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
	return importer.ExtractFromCell(cell, line)
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


type XmlGenerator struct {
	Generator

	fieldMap map[string]*tableaupb.Field
	nav *xmlquery.NodeNavigator
}

func NewXmlGenerator(protoPackage, goPackage, indir, outdir string, setters ...options.Option) *XmlGenerator {
	g := &XmlGenerator{
		Generator: *newGenerator(protoPackage, goPackage, indir, outdir, setters...),
	}
	g.i = g
	return g
}

func (gen XmlGenerator) GetSuffix() string {
	return ".xml"
}

func (gen *XmlGenerator) Generate() error {
	if err := gen.PrepareOutpuDir(); err != nil {
		return errors.Wrapf(err, "failed to prepare output dir: %s", gen.OutputDir)
	}
	files, err := os.ReadDir(gen.InputDir)
	if err != nil {
		atom.Log.Fatal(err)
	}
	for _, xmlFile := range files {
		// ignore temp file named with prefix "~$"
		if strings.HasPrefix(xmlFile.Name(), "~$") || !strings.HasSuffix(xmlFile.Name(), gen.GetSuffix()) {
			continue
		}
		// open xml file and parse the document
		xmlPath := filepath.Join(gen.InputDir, xmlFile.Name())
		atom.Log.Debugf("xml: %s", xmlPath)
		f, err := os.Open(xmlPath)
		if err != nil {
			atom.Log.Panic(err)
			continue
		}
		p, err := xmlquery.CreateStreamParser(f, "/")
		if err != nil {
			atom.Log.Panic(err)
			continue
		}
		// create xml proto meta struct
		xmlProtoName := strcase.ToSnake(strings.TrimSuffix(xmlFile.Name(), filepath.Ext(xmlFile.Name())))
		xml := &tableaupb.Workbook{
			Options: &tableaupb.WorkbookOptions{
				Name: xmlFile.Name(),
			},
			Name:       xmlProtoName,
			Imports: map[string]int32{
				tableauProtoPath: 1, // default import
			},
		}
		for _, path := range gen.Imports {
			xml.Imports[path] = 1 // custom imports
		}
		n, err := p.Read()
		if err != nil {
			atom.Log.Panic(err)
		}
		gen.fieldMap = make(map[string]*tableaupb.Field)
		root := xmlquery.CreateXPathNavigator(n)
		worksheet := &tableaupb.Worksheet{
			Options: &tableaupb.WorksheetOptions{
				Name: root.LocalName(),
			},
			Name: root.LocalName(),
		}
		field := &tableaupb.Field{}
		gen.parseNode(root, field, "")
		atom.Log.Debug(field)
		worksheet.Fields = append(worksheet.Fields, field.Fields...) // root节点变成了sheet
		xml.Worksheets = append(xml.Worksheets, worksheet)
		// export book
		be := newBookExporter(gen.ProtoPackage, gen.GoPackage, gen.OutputDir, gen.FilenameSuffix, gen.Imports, xml)
		if err := be.export(); err != nil {
			return errors.Wrapf(err, "failed to export workbook: %s", xmlPath)
		}
	}

	return nil
}