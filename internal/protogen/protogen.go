package protogen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"encoding/xml"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen"
	"github.com/Wenchy/tableau/internal/fs"
	"github.com/Wenchy/tableau/internal/importer"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

const (
	app                 = "tableauc"
	version             = "0.1.1"
	TableauProtoPackage = "tableau"
)

type IGenerator interface {
	GetSuffix() string // get filename extension
	convert(dir, filename string) error // convert one input file to proto
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
	Header *options.HeaderOption // header settings.

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
		Header:         opts.Header,
	}
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

func (gen *Generator) Generate() error {
	if err := gen.PrepareOutpuDir(); err != nil {
		return errors.Wrapf(err, "failed to prepare output dir: %s", gen.OutputDir)
	}
	return gen.generate(gen.InputDir)
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
		if err := gen.convert(dir, entry.Name()); err != nil {
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

func (gen *Generator) convert(dir, filename string) error {
	return gen.i.convert(dir, filename)
}

type XlsxGenerator struct {
	Generator
}

func NewXlsxGenerator(protoPackage, goPackage, indir, outdir string, setters ...options.Option) *XlsxGenerator {
	g := &XlsxGenerator{
		Generator: *newGenerator(protoPackage, goPackage, indir, outdir, setters...),
	}
	g.i = g
	return g
}

func (gen XlsxGenerator) GetSuffix() string {
	return ".xlsx"
}

func (gen *XlsxGenerator) convert(dir, filename string) error {
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

func (gen *XlsxGenerator) GenOneWorkbook(relativeWorkbookPath string) error {
	absPath := filepath.Join(gen.InputDir, relativeWorkbookPath)
	return gen.convert(filepath.Dir(absPath), filepath.Base(absPath))
}

type sheetHeader struct {
	meta    *tableaupb.SheetMeta
	namerow []string
	typerow []string
	noterow []string
}

func getCell(row []string, cursor int, line int32) string {
	// empty cell may be not in list
	if cursor >= len(row) {
		return ""
	}
	return importer.ExtractFromCell(row[cursor], line)
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

func (gen *XmlGenerator) convert(dir, filename string) error {
	// open xml file and parse the document
	xmlPath := filepath.Join(dir, filename)
	atom.Log.Debugf("xml: %s", xmlPath)
	buf, err := os.ReadFile(xmlPath)
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", xmlPath)
	}
	// replacement for `<` and `>` not allowed in attribute values
	attrValRegexp := regexp.MustCompile(`"\S+"`)
	nudeRegexp := regexp.MustCompile(`<([A-Za-z0-9]+)>`)
	keywordRegexp := regexp.MustCompile(`([</]+)@([A-Z]+)`)
	replaced_str := attrValRegexp.ReplaceAllStringFunc(string(buf), func(s string) string {
		var buf bytes.Buffer
		xml.EscapeText(&buf, []byte(s[1:len(s)-1]))
		return fmt.Sprintf("\"%s\"", buf.String())
	})
	replaced_str = nudeRegexp.ReplaceAllString(replaced_str, `<$1 TableauPlaceholder="0">`)
	replaced_str = keywordRegexp.ReplaceAllString(replaced_str, `$1$2`)
	s := strings.NewReader(replaced_str)
	p, err := xmlquery.CreateStreamParser(s, "/")
	if err != nil {
		return errors.Wrapf(err, "failed to create parser for string %s", s)
	}
	n, err := p.Read()
	if err != nil {
		return errors.Wrapf(err, "failed to read from string %s", s)
	}
	root := xmlquery.CreateXPathNavigator(n)
	// relatice path with filename
	relativePath, err := getRelCleanSlashPath(gen.InputDir, dir, filename)
	if err != nil {
		return errors.WithMessagef(err, "get relative path failed")
	}
	wbBasePath := filepath.Join(gen.InputDir, ".xml2xlsx")
	xlsxGen := xlsxgen.Generator{
		OutputDir: filepath.Join(wbBasePath, filepath.Dir(relativePath)),
		Workbook: strings.ReplaceAll(filename, gen.GetSuffix(), ".xlsx"),
	}
	// generate sheet `@TABLEAU`
	metaSheet := xlsxgen.NewMetaSheet(excel.TableauSheetName, &options.HeaderOption{
		Namerow: 1,
		Datarow: 2,
	}, false)
	metaSheet.Rows = []xlsxgen.Row{{
			Index: int(metaSheet.Namerow)-1,
			Cells: []xlsxgen.Cell{
				{Data: "Sheet"},
				{Data: "Alias"},
				{Data: "Nameline"},
				{Data: "Typeline"},
				{Data: "Nested"},
			},
		},{
			Index: int(metaSheet.Datarow)-1,
			Cells: []xlsxgen.Cell{
				{Data: root.LocalName()},
				{Data: root.LocalName()},
				{Data: "1"},
				{Data: "1"},
				{Data: "true"},
			},
		},
	}
	xlsxGen.ExportSheet(metaSheet)
	// generate data sheet
	metaSheet = xlsxgen.NewMetaSheet(root.LocalName(), gen.Header, false)
	if err := gen.parseXml(root, metaSheet, int(metaSheet.Datarow)-1); err != nil {
		return errors.Wrapf(err, "parseXml for root node %s failed", root.LocalName())
	}
	xlsxGen.ExportSheet(metaSheet)
	// generate proto by XlsxGenerator.convert
	wbPath := filepath.Join(xlsxGen.OutputDir, xlsxGen.Workbook)
	xlsx2ProtoGen := NewXlsxGenerator(
						gen.ProtoPackage,
						gen.GoPackage,
						wbBasePath,
						gen.OutputDir,
						options.Header(gen.Header),
						options.Imports(gen.Imports),
						options.Output(
							&options.OutputOption{
								FilenameSuffix: gen.FilenameSuffix,
								FilenameWithSubdirPrefix: gen.FilenameWithSubdirPrefix,
							},
						),
						options.LocationName(gen.LocationName),
					)
	if err := xlsx2ProtoGen.convert(xlsxGen.OutputDir, xlsxGen.Workbook); err != nil {
		return errors.Wrapf(err, "failed to parse workbook: %s", wbPath)
	}

	return nil
}