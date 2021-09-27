package protogen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/pkg/tableaupb"
	"github.com/golang/protobuf/proto"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
)

var mapRegexp *regexp.Regexp
var listRegexp *regexp.Regexp

func init() {
	mapRegexp = regexp.MustCompile(`^map<(.+),(.+)>`) // e.g.: map<uint32,MessageType>
	listRegexp = regexp.MustCompile(`^\[(.+)\](.+)`)  // e.g.: [Section]uint32
}

type Generator struct {
	ProtoPackage string // protobuf package name.
	GoPackage    string // golang package name.
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated protoconf files.

}

func (gen *Generator) Generate() {
	err := os.RemoveAll(gen.OutputDir)
	if err != nil {
		panic(err)
	}
	// create output dir
	err = os.MkdirAll(gen.OutputDir, 0700)
	if err != nil {
		panic(err)
	}

	files, err := os.ReadDir(gen.InputDir)
	if err != nil {
		atom.Log.Fatal(err)
	}
	for _, wbFile := range files {
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
					"tableau_options.proto": 1, // default import
				},
			},
		}

		for _, sheetName := range f.GetSheetMap() {
			rows, err := f.GetRows(sheetName)
			if err != nil {
				atom.Log.Panic(err)
			}
			ws := &tableaupb.Worksheet{
				Options: &tableaupb.WorksheetOptions{
					Name:      sheetName,
					Namerow:   1,
					Typerow:   2,
					Descrow:   3,
					Datarow:   4,
					Transpose: false,
					Tags:      "",
				},
				Fields: []*tableaupb.Field{},
				Name:   sheetName,
			}
			namerow := rows[0]
			typerow := rows[1]

			for i := 0; i < len(namerow); i++ {
				nameCell := strings.TrimSpace(namerow[i])
				// typeCell := strings.TrimSpace(typerow[i])
				if nameCell == "" {
					continue
				}
				field := &tableaupb.Field{}
				cursor, err := book.parseField(i, namerow, typerow, field)
				if err != nil {
					atom.Log.Panic(err)
				}
				i = cursor
				ws.Fields = append(ws.Fields, field)
			}
			book.wb.Worksheets = append(book.wb.Worksheets, ws)
		}
		if err := gen.exportWorkbook(book.wb); err != nil {
			atom.Log.Panic(err)
		}
	}
}

type book struct {
	wb *tableaupb.Workbook
}

func (b *book) parseField(cursor int, namerow, typerow []string, field *tableaupb.Field) (int, error) {
	nameCell := strings.TrimSpace(namerow[cursor])
	typeCell := strings.TrimSpace(typerow[cursor])
	atom.Log.Debugf("column|name: %s, type: %s", nameCell, typeCell)
	var err error
	if matches := mapRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		// map
		keyType := matches[1]
		valueType := matches[2]

		field.Name = strcase.ToSnake(valueType) + "_map"
		field.Type = typeCell
		field.MapEntry = &tableaupb.MapEntry{
			KeyType:   keyType,
			ValueType: valueType,
		}
		field.Options = &tableaupb.FieldOptions{
			Key: nameCell,
		}
		field.Fields = append(field.Fields, b.parseScalarField(nameCell, keyType))
		for cursor++; cursor < len(namerow); cursor++ {
			nameCell := strings.TrimSpace(namerow[cursor])
			if nameCell == "" {
				continue
			}
			subField := &tableaupb.Field{}
			cursor, err = b.parseField(cursor, namerow, typerow, subField)
			if err != nil {
				atom.Log.Panic(err)
			}
			field.Fields = append(field.Fields, subField)
		}
		return cursor, nil
	} else if matches := listRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		// list
		elemType := matches[1]
		colType := matches[2]
		// preprocess
		layout := tableaupb.Layout_LAYOUT_VERTICAL // default layout is vertical.
		index := -1
		if index = strings.Index(nameCell, "1"); index > 0 {
			layout = tableaupb.Layout_LAYOUT_HORIZONTAL
		}

		if layout == tableaupb.Layout_LAYOUT_VERTICAL {
			// vertical list: all columns belong to this list after this cursor.
			field.Card = "repeated"
			field.Name = strcase.ToSnake(elemType) + "_list"
			field.Type = elemType
			field.Options = &tableaupb.FieldOptions{
				Name:   "", // name is empty for vertical list
				Layout: layout,
			}
			field.Fields = append(field.Fields, b.parseScalarField(nameCell, colType))

			for cursor++; cursor < len(namerow); cursor++ {
				nameCell := strings.TrimSpace(namerow[cursor])
				if nameCell == "" {
					continue
				}
				subField := &tableaupb.Field{}
				cursor, err = b.parseField(cursor, namerow, typerow, subField)
				if err != nil {
					atom.Log.Panic(err)
				}
				field.Fields = append(field.Fields, subField)
			}
		} else {
			// horizontal list: continuous N columns belong to this list after this cursor.
			prefix := nameCell[:index]
			name := prefix

			field.Card = "repeated"
			field.Name = strcase.ToSnake(name) + "_list"
			field.Type = elemType
			field.Options = &tableaupb.FieldOptions{
				Name:   prefix,
				Layout: layout,
			}
			camelCaseName := nameCell[index+1:]
			field.Fields = append(field.Fields, b.parseScalarField(camelCaseName, colType))

			for cursor++; cursor < len(namerow); cursor++ {
				nameCell := strings.TrimSpace(namerow[cursor])
				typeCell := strings.TrimSpace(typerow[cursor])
				if nameCell == "" {
					continue
				}
				if strings.HasPrefix(nameCell, prefix+"1") {
					camelCaseName = nameCell[index+1:]
					field.Fields = append(field.Fields, b.parseScalarField(camelCaseName, typeCell))
				} else if strings.HasPrefix(nameCell, prefix) {
					continue
				} else {
					cursor--
					break
				}
			}
		}
	} else {
		// scalar
		*field = *b.parseScalarField(nameCell, typeCell)
	}

	return cursor, nil
}

func (b *book) parseScalarField(name, typ string) *tableaupb.Field {
	if typ == "timestamp" {
		typ = "google.protobuf.Timestamp"
		b.wb.Imports["google/protobuf/timestamp.proto"] = 1
	} else if typ == "duration" {
		typ = "google.protobuf.Duration"
		b.wb.Imports["google/protobuf/duration.proto"] = 1
	}

	return &tableaupb.Field{
		Name: strcase.ToSnake(name),
		Type: typ,
		Options: &tableaupb.FieldOptions{
			Name: name,
		},
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

	w := bufio.NewWriter(f)
	defer w.Flush()

	w.WriteString("syntax = \"proto3\";\n")
	w.WriteString(fmt.Sprintf("package %s;\n", gen.ProtoPackage))
	w.WriteString(fmt.Sprintf("option go_package = \"%s\";\n", gen.GoPackage))
	w.WriteString("\n")
	for key, _ := range wb.Imports {
		w.WriteString(fmt.Sprintf("import \"%s\";\n", key))
	}
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf("option (tableau.workbook) = {%s};\n", proto.CompactTextString(wb.Options)))
	w.WriteString("\n")

	for _, ws := range wb.Worksheets {
		if err := gen.exportWorksheet(w, ws); err != nil {
			return err
		}
	}

	return nil
}

func (gen *Generator) exportWorksheet(w *bufio.Writer, ws *tableaupb.Worksheet) error {
	w.WriteString(fmt.Sprintf("message %s {\n", ws.Name))
	w.WriteString(fmt.Sprintf("  option (tableau.worksheet) = {%s};\n", proto.CompactTextString(ws.Options)))
	w.WriteString("\n")

	depth := 1
	for i, f := range ws.Fields {
		tagid := i + 1
		if err := gen.exportField(depth, w, tagid, f); err != nil {
			return err
		}
	}

	w.WriteString("}\n")
	return nil
}

func (gen *Generator) exportField(depth int, w *bufio.Writer, tagid int, field *tableaupb.Field) error {
	head := "%s%s"
	if field.Card != "" {
		head += " " // cardinality exists
	}
	w.WriteString(fmt.Sprintf(head+"%s %s = %d [(tableau.field) = {%s}];\n", indent(depth), field.Card, field.Type, field.Name, tagid, proto.CompactTextString(field.Options)))

	if field.Fields != nil { // iff field is a map or list.
		embbedMsgName := field.Type
		if field.MapEntry != nil {
			embbedMsgName = field.MapEntry.ValueType
		}
		w.WriteString("\n")
		w.WriteString(fmt.Sprintf("%smessage %s {\n", indent(depth), embbedMsgName))
		for i, f := range field.Fields {
			tagid := i + 1
			if err := gen.exportField(depth+1, w, tagid, f); err != nil {
				return err
			}
		}
		w.WriteString(fmt.Sprintf("%s}\n", indent(depth)))
	}
	return nil
}

func indent(depth int) string {
	return strings.Repeat("  ", depth)
}