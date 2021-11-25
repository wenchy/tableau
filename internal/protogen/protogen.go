package protogen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/fs"
	"github.com/Wenchy/tableau/internal/types"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/golang/protobuf/proto"
	"github.com/iancoleman/strcase"
	"github.com/xuri/excelize/v2"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/reflect/protoreflect"
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
			namerow := rows[0]
			typerow := rows[1]
			noterow := rows[2]

			var ok bool
			for cursor := 0; cursor < len(namerow); cursor++ {
				field := &tableaupb.Field{}
				cursor, ok = book.parseField(cursor, namerow, typerow, noterow, field, "")
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

type book struct {
	wb       *tableaupb.Workbook
	withNote bool
}

func (b *book) parseField(cursor int, namerow, typerow, noterow []string, field *tableaupb.Field, prefix string) (cur int, ok bool) {
	nameCell := strings.TrimSpace(namerow[cursor])
	typeCell := strings.TrimSpace(typerow[cursor])
	noteCell := strings.TrimSpace(noterow[cursor])
	atom.Log.Debugf("column|name: %s, type: %s", nameCell, typeCell)
	if nameCell == "" || typeCell == "" {
		atom.Log.Warnf("no need to parse column, as name(%s) or type(%s) is empty", nameCell, typeCell)
		return cursor, false
	}
	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	if matches := mapRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		// map
		keyType := strings.TrimSpace(matches[1])
		valueType := strings.TrimSpace(matches[2])

		if types.IsScalarType(valueType) {
			// incell map
			field.Name = strcase.ToSnake(trimmedNameCell)
			field.Type, field.TypeDefined = ParseType(typeCell)
			field.Options = &tableaupb.FieldOptions{
				Name: trimmedNameCell,
				Type: tableaupb.Type_TYPE_INCELL_MAP,
			}
		} else {
			field.Name = strcase.ToSnake(valueType) + "_map"
			field.Type, field.TypeDefined = ParseType(typeCell)
			field.MapEntry = &tableaupb.MapEntry{
				KeyType:   keyType,
				ValueType: valueType,
			}
			field.Options = &tableaupb.FieldOptions{
				Key: trimmedNameCell,
			}
			field.Fields = append(field.Fields, b.parseScalarField(trimmedNameCell, keyType, noteCell))
			for cursor++; cursor < len(namerow); cursor++ {
				cursor = b.parseSubField(cursor, namerow, typerow, noterow, field, prefix)
			}
		}
	} else if matches := listRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		// list
		colType := strings.TrimSpace(matches[2])
		var isScalarType bool
		elemType := strings.TrimSpace(matches[1])
		if elemType == "" {
			// scalar type, such as int32, string, etc.
			elemType = colType
			isScalarType = true
		}

		// preprocess
		layout := tableaupb.Layout_LAYOUT_VERTICAL // default layout is vertical.
		index := -1
		if index = strings.Index(trimmedNameCell, "1"); index > 0 {
			layout = tableaupb.Layout_LAYOUT_HORIZONTAL
		} else {
			if isScalarType {
				// incell list
				layout = tableaupb.Layout_LAYOUT_DEFAULT
			}
		}

		if layout == tableaupb.Layout_LAYOUT_VERTICAL {
			// vertical list: all columns belong to this list after this cursor.
			field.Card = "repeated"
			field.Name = strcase.ToSnake(elemType) + "_list"
			field.Type, field.TypeDefined = ParseType(elemType)
			field.Options = &tableaupb.FieldOptions{
				Name:   "", // name is empty for vertical list
				Layout: layout,
			}

			if isScalarType {
				// TODO: support list of scalar type when lyout is vertical?
				// NOTE(wenchyzhu): we don't support list of scalar type when layout is vertical
			} else {
				field.Fields = append(field.Fields, b.parseScalarField(trimmedNameCell, colType, noteCell))
				for cursor++; cursor < len(namerow); cursor++ {
					cursor = b.parseSubField(cursor, namerow, typerow, noterow, field, prefix)
				}
			}
		} else if layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
			// horizontal list: continuous N columns belong to this list after this cursor.
			listName := trimmedNameCell[:index]
			prefix += listName

			field.Card = "repeated"
			field.Name = strcase.ToSnake(listName) + "_list"
			field.Type, field.TypeDefined = ParseType(elemType)
			field.Options = &tableaupb.FieldOptions{
				Name:   listName,
				Layout: layout,
			}
			if isScalarType {
				for cursor++; cursor < len(namerow); cursor++ {
					nameCell := strings.TrimSpace(namerow[cursor])
					if nameCell == "" {
						continue
					}
					if strings.HasPrefix(nameCell, prefix) {
						continue
					} else {
						cursor--
						break
					}
				}
			} else {
				name := strings.TrimPrefix(nameCell, prefix+"1")
				field.Fields = append(field.Fields, b.parseScalarField(name, colType, noteCell))
				for cursor++; cursor < len(namerow); cursor++ {
					nameCell := strings.TrimSpace(namerow[cursor])
					if strings.HasPrefix(nameCell, prefix+"1") {
						cursor = b.parseSubField(cursor, namerow, typerow, noterow, field, prefix+"1")
					} else if strings.HasPrefix(nameCell, prefix) {
						continue
					} else {
						cursor--
						break
					}
				}
			}
		} else if layout == tableaupb.Layout_LAYOUT_DEFAULT {
			// incell list
			field.Card = "repeated"
			field.Name = strcase.ToSnake(trimmedNameCell)
			field.Type, field.TypeDefined = ParseType(elemType)
			field.Options = &tableaupb.FieldOptions{
				Name: trimmedNameCell,
				Type: tableaupb.Type_TYPE_INCELL_LIST,
			}
		}
	} else if matches := structRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		// struct
		elemType := strings.TrimSpace(matches[1])
		colType := strings.TrimSpace(matches[2])

		fieldPairs := ParseIncellStruct(elemType)
		if fieldPairs == nil {
			// cross cell struct
			field.Type, field.TypeDefined = ParseType(elemType)
			field.Name = strcase.ToSnake(field.Type)
			index := len(field.Type)
			structName := trimmedNameCell[:index]
			field.Options = &tableaupb.FieldOptions{
				Name: structName,
			}
			prefix += structName

			name := strings.TrimPrefix(nameCell, prefix)
			field.Fields = append(field.Fields, b.parseScalarField(name, colType, noteCell))
			for cursor++; cursor < len(namerow); cursor++ {
				nameCell := strings.TrimSpace(namerow[cursor])
				if strings.HasPrefix(nameCell, prefix) {
					cursor = b.parseSubField(cursor, namerow, typerow, noterow, field, prefix)
				} else {
					break
				}
			}
		} else {
			// incell struct
			field.Name = strcase.ToSnake(trimmedNameCell)
			field.Type, field.TypeDefined = ParseType(colType)
			field.Options = &tableaupb.FieldOptions{
				Name: trimmedNameCell,
				Type: tableaupb.Type_TYPE_INCELL_STRUCT,
			}

			for i := 0; i < len(fieldPairs); i += 2 {
				fieldType := fieldPairs[i]
				fieldName := fieldPairs[i+1]
				field.Fields = append(field.Fields, b.parseScalarField(fieldName, fieldType, ""))
			}
		}

	} else {
		// scalar
		*field = *b.parseScalarField(trimmedNameCell, typeCell, noteCell)
	}

	return cursor, true
}

func (b *book) parseSubField(cursor int, namerow, typerow, noterow []string, field *tableaupb.Field, prefix string) int {
	subField := &tableaupb.Field{}
	cursor, ok := b.parseField(cursor, namerow, typerow, noterow, subField, prefix)
	if ok {
		field.Fields = append(field.Fields, subField)
	}
	return cursor
}

func (b *book) genNote(note string) string {
	if b.withNote {
		return note
	}
	return ""
}

func (b *book) parseScalarField(name, typ, note string) *tableaupb.Field {
	if typ == "timestamp" {
		typ = "google.protobuf.Timestamp"
		b.wb.Imports[timestampProtoPath] = 1
	} else if typ == "duration" {
		typ = "google.protobuf.Duration"
		b.wb.Imports[durationProtoPath] = 1
	}
	return &tableaupb.Field{
		Name: strcase.ToSnake(name),
		Type: typ,
		Options: &tableaupb.FieldOptions{
			Name: name,
			Note: b.genNote(note),
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

	// keep the elements ordered by sheet name
	set := treeset.NewWithStringComparator()
	for key := range wb.Imports {
		set.Add(key)
	}
	for _, key := range set.Values() {
		w.WriteString(fmt.Sprintf("import \"%s\";\n", key))
	}
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf("option (tableau.workbook) = {%s};\n", genPrototext(wb.Options)))
	w.WriteString("\n")

	for i, ws := range wb.Worksheets {
		isLastSheet := false
		if i == len(wb.Worksheets)-1 {
			isLastSheet = true
		}
		sheet := &sheet{
			ws:             ws,
			writer:         w,
			isLastSheet:    isLastSheet,
			nestedMessages: make(map[string]*tableaupb.Field),
		}
		if err := sheet.export(); err != nil {
			return err
		}
	}

	return nil
}

type sheet struct {
	ws             *tableaupb.Worksheet
	writer         *bufio.Writer
	isLastSheet    bool
	nestedMessages map[string]*tableaupb.Field // type name -> field
}

func (s *sheet) export() error {
	s.writer.WriteString(fmt.Sprintf("message %s {\n", s.ws.Name))
	s.writer.WriteString(fmt.Sprintf("  option (tableau.worksheet) = {%s};\n", genPrototext(s.ws.Options)))
	s.writer.WriteString("\n")
	// generate the fields
	depth := 1
	for i, field := range s.ws.Fields {
		tagid := i + 1
		if err := s.exportField(depth, tagid, field); err != nil {
			return err
		}
	}
	s.writer.WriteString("}\n")
	if !s.isLastSheet {
		s.writer.WriteString("\n")
	}
	return nil
}

func genPrototext(m protoreflect.ProtoMessage) []byte {
	// text := proto.CompactTextString(field.Options)
	text, err := prototext.Marshal(m)
	if err != nil {
		panic(err)
	}
	return text
}

func (s *sheet) exportField(depth int, tagid int, field *tableaupb.Field) error {
	head := "%s%s"
	if field.Card != "" {
		head += " " // cardinality exists
	}
	s.writer.WriteString(fmt.Sprintf(head+"%s %s = %d [(tableau.field) = {%s}];\n", indent(depth), field.Card, field.Type, field.Name, tagid, genPrototext(field.Options)))

	if !field.TypeDefined && field.Fields != nil {
		// iff field is a map or list and message type is not imported.
		nestedMsgName := field.Type
		if field.MapEntry != nil {
			nestedMsgName = field.MapEntry.ValueType
		}

		if isSameFieldMessageType(field, s.nestedMessages[nestedMsgName]) {
			// if the nested message is the same as the previous one,
			// just use the previous one, and don't generate a new one.
			return nil
		}

		// bookkeeping this nested msessage, so we can check if we can reuse it later.
		s.nestedMessages[nestedMsgName] = field

		s.writer.WriteString("\n")
		s.writer.WriteString(fmt.Sprintf("%smessage %s {\n", indent(depth), nestedMsgName))
		for i, f := range field.Fields {
			tagid := i + 1
			if err := s.exportField(depth+1, tagid, f); err != nil {
				return err
			}
		}
		s.writer.WriteString(fmt.Sprintf("%s}\n", indent(depth)))
	}
	return nil
}

func indent(depth int) string {
	return strings.Repeat("  ", depth)
}

func isSameFieldMessageType(left, right *tableaupb.Field) bool {
	if left == nil || right == nil {
		return false
	}
	if left.Fields == nil || right.Fields == nil {
		return false
	}
	if len(left.Fields) != len(right.Fields) ||
		left.Type != right.Type ||
		left.Card != right.Card {
		return false
	}

	for i, l := range left.Fields {
		r := right.Fields[i]
		if !proto.Equal(l, r) {
			return false
		}
	}
	return true
}

func ParseType(msgName string) (string, bool) {
	if strings.Contains(msgName, ".") {
		// This messge type is defined in imported proto
		msgName = strings.TrimPrefix(msgName, ".")
		return msgName, true
	}
	// if matches := mapRegexp.FindStringSubmatch(msgName); len(matches) > 0 {
	// 	// map
	// 	keyType := strings.TrimSpace(matches[1])
	// 	valueType := strings.TrimSpace(matches[2])
	// 	return msgName, types.IsScalarType(keyType) && types.IsScalarType(valueType)
	// }
	return msgName, false
}

func ParseIncellStruct(elemType string) []string {
	fields := strings.Split(elemType, ",")
	if len(fields) == 1 && len(strings.Split(fields[0], " ")) == 1 {
		// cross cell struct
		return nil
	}

	fieldPairs := make([]string, 0)
	for _, pair := range strings.Split(elemType, ",") {
		kv := strings.Split(pair, " ")
		if len(kv) != 2 {
			atom.Log.Panicf("illegal type-variable pair: %v in incell struct: %s", pair, elemType)
		}
		fieldPairs = append(fieldPairs, kv...)
	}
	return fieldPairs
}
