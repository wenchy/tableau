package protogen

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/types"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/iancoleman/strcase"
)

const (
	tableauProtoPath   = "tableau/protobuf/tableau.proto"
	timestampProtoPath = "google/protobuf/timestamp.proto"
	durationProtoPath  = "google/protobuf/duration.proto"
)

var mapRegexp *regexp.Regexp
var listRegexp *regexp.Regexp
var structRegexp *regexp.Regexp
var enumRegexp *regexp.Regexp


func init() {
	mapRegexp = regexp.MustCompile(`^map<(.+),(.+)>`)              // e.g.: map<uint32,Type>
	listRegexp = regexp.MustCompile(`^\[(.*)\](.+)`)               // e.g.: [Type]uint32
	structRegexp = regexp.MustCompile(`^\{(.+)\}(.+)`)             // e.g.: {Type}uint32
	enumRegexp = regexp.MustCompile(`^enum<(.+)>`)                 // e.g.: enum<Type>
}

type bookParser struct {
	wb       *tableaupb.Workbook
	withNote bool
	types    map[string]bool // type name -> existed
}

func newBookParser(workbookName string, imports []string) *bookParser {
	wbProtoName := strcase.ToSnake(strings.TrimSuffix(workbookName, filepath.Ext(workbookName)))
	bp := &bookParser{
		wb: &tableaupb.Workbook{
			Options: &tableaupb.WorkbookOptions{
				Name: workbookName,
			},
			Worksheets: []*tableaupb.Worksheet{},
			Name:       wbProtoName,
			Imports:    make(map[string]int32),
		},
		withNote: false,
	}

	// custom imports
	for _, path := range imports {
		bp.wb.Imports[path] = 1
	}
	return bp
}

func (p *bookParser) parseField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string) (cur int, ok bool) {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)
	// atom.Log.Debugf("column: %d, name: %s, type: %s", cursor, nameCell, typeCell)
	if nameCell == "" || typeCell == "" {
		atom.Log.Warnf("no need to parse column %d, as name(%s) or type(%s) is empty", cursor, nameCell, typeCell)
		return cursor, false
	}
	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	if matches := mapRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		cursor = p.parseMapField(field, header, cursor, prefix)
	} else if matches := listRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		cursor = p.parseListField(field, header, cursor, prefix)
	} else if matches := structRegexp.FindStringSubmatch(typeCell); len(matches) > 0 {
		cursor = p.parseStructField(field, header, cursor, prefix)
	} else {
		// scalar
		*field = *p.parseScalarField(trimmedNameCell, typeCell, noteCell)
	}

	return cursor, true
}

func (p *bookParser) parseSubField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string) int {
	subField := &tableaupb.Field{}
	cursor, ok := p.parseField(subField, header, cursor, prefix)
	if ok {
		field.Fields = append(field.Fields, subField)
	}
	return cursor
}

func (p *bookParser) parseMapField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string) int {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// map syntax pattern
	matches := mapRegexp.FindStringSubmatch(typeCell)
	keyType := strings.TrimSpace(matches[1])
	valueType := strings.TrimSpace(matches[2])

	if types.IsScalarType(valueType) {
		// incell map
		field.Name = strcase.ToSnake(trimmedNameCell)
		field.Type, field.TypeDefined = p.parseType(typeCell)
		field.Options = &tableaupb.FieldOptions{
			Name: trimmedNameCell,
			Type: tableaupb.Type_TYPE_INCELL_MAP,
		}
	} else {
		field.Name = strcase.ToSnake(valueType) + "_map"
		field.Type, field.TypeDefined = p.parseType(typeCell)
		field.MapEntry = &tableaupb.MapEntry{
			KeyType:   keyType,
			ValueType: valueType,
		}
		field.Options = &tableaupb.FieldOptions{
			Key: trimmedNameCell,
		}
		field.Fields = append(field.Fields, p.parseScalarField(trimmedNameCell, keyType, noteCell))
		for cursor++; cursor < len(header.namerow); cursor++ {
			cursor = p.parseSubField(field, header, cursor, prefix)
		}
	}
	return cursor
}

func (p *bookParser) parseListField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string) int {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// list syntax pattern
	matches := listRegexp.FindStringSubmatch(typeCell)
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

	switch layout {
	case tableaupb.Layout_LAYOUT_VERTICAL:
		// vertical list: all columns belong to this list after this cursor.
		field.Card = "repeated"
		field.Name = strcase.ToSnake(elemType) + "_list"
		field.Type, field.TypeDefined = p.parseType(elemType)
		field.Options = &tableaupb.FieldOptions{
			Name:   "", // name is empty for vertical list
			Layout: layout,
		}

		if isScalarType {
			// TODO: support list of scalar type when lyout is vertical?
			// NOTE(wenchyzhu): we don't support list of scalar type when layout is vertical
		} else {
			field.Fields = append(field.Fields, p.parseScalarField(trimmedNameCell, colType, noteCell))
			for cursor++; cursor < len(header.namerow); cursor++ {
				cursor = p.parseSubField(field, header, cursor, prefix)
			}
		}
	case tableaupb.Layout_LAYOUT_HORIZONTAL:
		// horizontal list: continuous N columns belong to this list after this cursor.
		listName := trimmedNameCell[:index]
		prefix += listName

		field.Card = "repeated"
		field.Name = strcase.ToSnake(listName) + "_list"
		field.Type, field.TypeDefined = p.parseType(elemType)
		field.Options = &tableaupb.FieldOptions{
			Name:   listName,
			Layout: layout,
		}
		if isScalarType {
			for cursor++; cursor < len(header.namerow); cursor++ {
				nameCell := header.getNameCell(cursor)
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
			field.Fields = append(field.Fields, p.parseScalarField(name, colType, noteCell))
			for cursor++; cursor < len(header.namerow); cursor++ {
				nameCell := header.getNameCell(cursor)
				if strings.HasPrefix(nameCell, prefix+"1") {
					cursor = p.parseSubField(field, header, cursor, prefix+"1")
				} else if strings.HasPrefix(nameCell, prefix) {
					continue
				} else {
					cursor--
					break
				}
			}
		}
	case tableaupb.Layout_LAYOUT_DEFAULT:
		// incell list
		field.Card = "repeated"
		field.Name = strcase.ToSnake(trimmedNameCell)
		field.Type, field.TypeDefined = p.parseType(elemType)
		field.Options = &tableaupb.FieldOptions{
			Name: trimmedNameCell,
			Type: tableaupb.Type_TYPE_INCELL_LIST,
		}
	}
	return cursor
}

func (p *bookParser) parseStructField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string) int {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// struct syntax pattern
	matches := structRegexp.FindStringSubmatch(typeCell)
	elemType := strings.TrimSpace(matches[1])
	colType := strings.TrimSpace(matches[2])

	if fieldPairs := ParseIncellStruct(elemType); fieldPairs != nil {
		// incell struct
		field.Name = strcase.ToSnake(trimmedNameCell)
		field.Type, field.TypeDefined = p.parseType(colType)
		field.Options = &tableaupb.FieldOptions{
			Name: trimmedNameCell,
			Type: tableaupb.Type_TYPE_INCELL_STRUCT,
		}

		for i := 0; i < len(fieldPairs); i += 2 {
			fieldType := fieldPairs[i]
			fieldName := fieldPairs[i+1]
			field.Fields = append(field.Fields, p.parseScalarField(fieldName, fieldType, ""))
		}
	} else {
		// cross cell struct
		field.Type, field.TypeDefined = p.parseType(elemType)
		field.Name = strcase.ToSnake(field.Type)
		index := len(field.Type)
		structName := trimmedNameCell[:index]
		field.Options = &tableaupb.FieldOptions{
			Name: structName,
		}
		prefix += structName

		name := strings.TrimPrefix(nameCell, prefix)
		field.Fields = append(field.Fields, p.parseScalarField(name, colType, noteCell))
		for cursor++; cursor < len(header.namerow); cursor++ {
			nameCell := header.getNameCell(cursor)
			if strings.HasPrefix(nameCell, prefix) {
				cursor = p.parseSubField(field, header, cursor, prefix)
			} else {
				cursor--
				break
			}
		}
	}

	return cursor
}

func (p *bookParser) parseScalarField(name, typ, note string) *tableaupb.Field {
	// enum syntax pattern
	if matches := enumRegexp.FindStringSubmatch(typ); len(matches) > 0 {
		enumType := strings.TrimSpace(matches[1])
		typ = enumType
	}
	typ, typeDefined := p.parseType(typ)

	return &tableaupb.Field{
		Name:        strcase.ToSnake(name),
		Type:        typ,
		TypeDefined: typeDefined,
		Options: &tableaupb.FieldOptions{
			Name: name,
			Note: p.genNote(note),
		},
	}
}

func (p *bookParser) genNote(note string) string {
	if p.withNote {
		return note
	}
	return ""
}

func (p *bookParser) parseType(typ string) (string, bool) {
	if strings.Contains(typ, ".") {
		// This messge type is defined in imported proto
		typ = strings.TrimPrefix(typ, ".")
		return typ, true
	}
	if typ == "timestamp" {
		typ = "google.protobuf.Timestamp"
		// p.wb.Imports[timestampProtoPath] = 1
	} else if typ == "duration" {
		typ = "google.protobuf.Duration"
		// p.wb.Imports[durationProtoPath] = 1
	}
	return typ, false
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
