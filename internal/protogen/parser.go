package protogen

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"fmt"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/types"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/iancoleman/strcase"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
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
var numRegex *regexp.Regexp

func init() {
	mapRegexp = regexp.MustCompile(`^map<(.+),(.+)>`)  // e.g.: map<uint32,Type>
	listRegexp = regexp.MustCompile(`^\[(.*)\](.+)`)   // e.g.: [Type]uint32
	structRegexp = regexp.MustCompile(`^\{(.+)\}(.+)`) // e.g.: {Type}uint32
	enumRegexp = regexp.MustCompile(`^enum<(.+)>`)     // e.g.: enum<Type>
}

type bookParser struct {
	wb       *tableaupb.Workbook
	types    map[string]bool // type name -> existed
	withNote bool
}

func newBookParser(relSlashPath string, filenameWithSubdirPrefix bool, imports []string) *bookParser {
	// atom.Log.Debugf("filenameWithSubdirPrefix: %v", filenameWithSubdirPrefix)
	ext := filepath.Ext(relSlashPath)
	filename := ""
	if filenameWithSubdirPrefix {
		snakePath := strcase.ToSnake(strings.TrimSuffix(relSlashPath, ext))
		filename = strings.ReplaceAll(snakePath, "/", "__")
	} else {
		workbookName := filepath.Base(relSlashPath)
		filename = strcase.ToSnake(strings.TrimSuffix(workbookName, ext))
	}
	bp := &bookParser{
		wb: &tableaupb.Workbook{
			Options: &tableaupb.WorkbookOptions{
				// NOTE(wenchyzhu): all OS platforms use path slash separator `/`
				// see: https://stackoverflow.com/questions/9371031/how-do-i-create-crossplatform-file-paths-in-go
				Name: relSlashPath,
			},
			Worksheets: []*tableaupb.Worksheet{},
			Name:       filename,
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

func (p *bookParser) parseField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string, nested bool) (cur int, ok bool) {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)
	// atom.Log.Debugf("column: %d, name: %s, type: %s", cursor, nameCell, typeCell)
	if nameCell == "" || typeCell == "" {
		atom.Log.Warnf("no need to parse column %d, as name(%s) or type(%s) is empty", cursor, nameCell, typeCell)
		return cursor, false
	}
	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	if matches := types.MatchMap(typeCell); len(matches) > 0 {
		cursor = p.parseMapField(field, header, cursor, prefix, nested)
	} else if matches := types.MatchList(typeCell); len(matches) > 0 {
		cursor = p.parseListField(field, header, cursor, prefix, nested)
	} else if matches := types.MatchStruct(typeCell); len(matches) > 0 {
		cursor = p.parseStructField(field, header, cursor, prefix, nested)
	} else {
		// scalar
		*field = *p.parseScalarField(trimmedNameCell, typeCell, noteCell)
	}

	return cursor, true
}

func (p *bookParser) parseSubField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string, nested bool) int {
	subField := &tableaupb.Field{}
	cursor, ok := p.parseField(subField, header, cursor, prefix, nested)
	if ok {
		field.Fields = append(field.Fields, subField)
		if field.Options.Layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
			field.Options.ListMaxLen /= int32(len(field.Fields))
		}
	}
	return cursor
}

func (p *bookParser) parseMapField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string, nested bool) int {
	// refer: https://developers.google.com/protocol-buffers/docs/proto3#maps
	//
	//	map<key_type, value_type> map_field = N;
	//
	// where the key_type can be any integral or string type (so, any scalar type
	// except for floating point types and bytes). Note that enum is not a valid
	// key_type. The value_type can be any type except another map.

	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	// map syntax pattern
	matches := types.MatchMap(typeCell)
	keyType := strings.TrimSpace(matches[1])
	valueType := strings.TrimSpace(matches[2])

	parsedKeyType := keyType
	if types.MatchEnum(keyType) != nil {
		// NOTE: support enum as map key, convert key type as `int32`.
		parsedKeyType = "int32"
	}
	parsedValueType, valueTypeDefined := p.parseType(valueType)
	mapType := fmt.Sprintf("map<%s, %s>", parsedKeyType, parsedValueType)

	isScalarType := types.IsScalarType(parsedValueType)
	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// preprocess
	layout := tableaupb.Layout_LAYOUT_VERTICAL // default layout is vertical.
	index := -1
	if index = strings.Index(trimmedNameCell, "1"); index > 0 {
		layout = tableaupb.Layout_LAYOUT_HORIZONTAL
		if cursor+1 < len(header.namerow) {
			// Header:
			//
			// TaskParamMap1		TaskParamMap2		TaskParamMap3
			// map<int32, int32>	map<int32, int32>	map<int32, int32>

			// check next cursor
			nextNameCell := header.getNameCell(cursor + 1)
			trimmedNextNameCell := strings.TrimPrefix(nextNameCell, prefix)
			if index2 := strings.Index(trimmedNextNameCell, "2"); index2 > 0 {
				nextTypeCell := header.getTypeCell(cursor + 1)
				if matches := types.MatchMap(nextTypeCell); len(matches) > 0 {
					// The next type cell is also a map type declaration.
					if isScalarType {
						layout = tableaupb.Layout_LAYOUT_DEFAULT // incell map
					}
				}
			} else {
				// only one map item, treat it as incell map
				if isScalarType {
					layout = tableaupb.Layout_LAYOUT_DEFAULT // incell map
				}
			}
		}
	} else {
		if isScalarType {
			layout = tableaupb.Layout_LAYOUT_DEFAULT // incell map
		}
	}

	switch layout {
	case tableaupb.Layout_LAYOUT_VERTICAL:
		if nested {
			prefix += parsedValueType // add prefix with value type
		}
		field.Name = strcase.ToSnake(parsedValueType) + "_map"
		field.Type = mapType
		// For map type, TypeDefined indicates the ValueType of map has been defined.
		field.TypeDefined = valueTypeDefined
		field.MapEntry = &tableaupb.MapEntry{
			KeyType:   parsedKeyType,
			ValueType: parsedValueType,
		}

		trimmedNameCell := strings.TrimPrefix(nameCell, prefix)
		field.Options = &tableaupb.FieldOptions{
			Key:    trimmedNameCell,
			Layout: layout,
		}
		if nested {
			field.Options.Name = parsedValueType
		}
		field.Fields = append(field.Fields, p.parseScalarField(trimmedNameCell, keyType, noteCell))
		for cursor++; cursor < len(header.namerow); cursor++ {
			if nested {
				nameCell := header.getNameCell(cursor)
				if !strings.HasPrefix(nameCell, prefix) {
					cursor--
					return cursor
				}
			}
			cursor = p.parseSubField(field, header, cursor, prefix, nested)
		}
	case tableaupb.Layout_LAYOUT_HORIZONTAL:
		if nested {
			prefix += parsedValueType // add prefix with value type
		}
		field.Name = strcase.ToSnake(parsedValueType) + "_map"
		field.Type = mapType
		// For map type, TypeDefined indicates the ValueType of map has been defined.
		field.TypeDefined = valueTypeDefined
		field.MapEntry = &tableaupb.MapEntry{
			KeyType:   parsedKeyType,
			ValueType: parsedValueType,
		}

		trimmedNameCell := strings.TrimPrefix(nameCell, prefix+"1")
		field.Options = &tableaupb.FieldOptions{
			Key:    trimmedNameCell,
			Layout: layout,
		}
		if nested {
			field.Options.Name = parsedValueType
		}

		name := strings.TrimPrefix(nameCell, prefix+"1")
		field.Fields = append(field.Fields, p.parseScalarField(name, keyType, noteCell))
		for cursor++; cursor < len(header.namerow); cursor++ {
			nameCell := header.getNameCell(cursor)
			if strings.HasPrefix(nameCell, prefix+"1") {
				cursor = p.parseSubField(field, header, cursor, prefix+"1", nested)
			} else if strings.HasPrefix(nameCell, prefix) {
				continue
			} else {
				cursor--
				return cursor
			}
		}

	case tableaupb.Layout_LAYOUT_DEFAULT:
		// incell map
		trimmedNameCell := strings.TrimPrefix(nameCell, prefix)
		field.Name = strcase.ToSnake(trimmedNameCell)
		field.Type = mapType
		// For map type, TypeDefined indicates the ValueType of map has been defined.
		field.TypeDefined = valueTypeDefined
		field.Options = &tableaupb.FieldOptions{
			Name: trimmedNameCell,
			Type: tableaupb.Type_TYPE_INCELL_MAP,
		}
	}

	return cursor
}

func (p *bookParser) parseListField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string, nested bool) int {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// list syntax pattern
	matches := types.MatchList(typeCell)
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
	tmpCursor := cursor
	if index = strings.Index(trimmedNameCell, "1"); index > 0 {
		layout = tableaupb.Layout_LAYOUT_HORIZONTAL
		if cursor+1 < len(header.namerow) {
			// Header:
			//
			// TaskParamList1	TaskParamList2	TaskParamList3
			// []int32			[]int32			[]int32

			// check next cursor
			nextNameCell := header.getNameCell(cursor + 1)
			trimmedNextNameCell := strings.TrimPrefix(nextNameCell, prefix)
			if index2 := strings.Index(trimmedNextNameCell, "2"); index2 > 0 {
				nextTypeCell := header.getTypeCell(cursor + 1)
				if matches := types.MatchList(nextTypeCell); len(matches) > 0 {
					// The next type cell is also a list type declaration.
					if isScalarType {
						layout = tableaupb.Layout_LAYOUT_DEFAULT // incell list
					}
				}
			} else {
				// only one list item, treat it as incell list
				if isScalarType {
					layout = tableaupb.Layout_LAYOUT_DEFAULT // incell list
				}
			}
		}
		for tmpNameCell, listName := trimmedNameCell, trimmedNameCell[:index]; strings.Contains(tmpNameCell, listName); {
			if tmpCursor++; tmpCursor >= len(header.namerow) {
				break
			}
			tmpNameCell = strings.TrimPrefix(header.getNameCell(tmpCursor), prefix)
		}
	} else {
		if isScalarType {
			layout = tableaupb.Layout_LAYOUT_DEFAULT // incell list
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
			// TODO: support list of scalar type when layout is vertical?
			atom.Log.Errorf("vertical list of scalar type is not supported")
		} else {
			if nested {
				prefix += field.Type // add prefix with value type
				field.Options.Name = field.Type
			}
			trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

			if matches := types.MatchKeyedList(typeCell); matches != nil {
				// set column type and key if this is a keyed list.
				colType = strings.TrimSpace(matches[2])
				field.Options.Key = trimmedNameCell
			}

			field.Fields = append(field.Fields, p.parseScalarField(trimmedNameCell, colType, noteCell))
			for cursor++; cursor < len(header.namerow); cursor++ {
				if nested {
					nameCell := header.getNameCell(cursor)
					if !strings.HasPrefix(nameCell, prefix) {
						cursor--
						return cursor
					}
				}
				cursor = p.parseSubField(field, header, cursor, prefix, nested)
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
			ListMaxLen: int32(tmpCursor - cursor),
		}
		if isScalarType {
			for cursor++; cursor < len(header.namerow); cursor++ {
				nameCell := header.getNameCell(cursor)
				if nameCell == "" || strings.HasPrefix(nameCell, prefix) {
					continue
				} else {
					cursor--
					return cursor
				}
			}
		} else {
			name := strings.TrimPrefix(nameCell, prefix+"1")
			field.Fields = append(field.Fields, p.parseScalarField(name, colType, noteCell))
			for cursor++; cursor < len(header.namerow); cursor++ {
				nameCell := header.getNameCell(cursor)
				if strings.HasPrefix(nameCell, prefix+"1") {
					cursor = p.parseSubField(field, header, cursor, prefix+"1", nested)
				} else if strings.HasPrefix(nameCell, prefix) {
					continue
				} else {
					cursor--
					return cursor
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

func (p *bookParser) parseStructField(field *tableaupb.Field, header *sheetHeader, cursor int, prefix string, nested bool) int {
	nameCell := header.getNameCell(cursor)
	typeCell := header.getTypeCell(cursor)
	noteCell := header.getNoteCell(cursor)

	trimmedNameCell := strings.TrimPrefix(nameCell, prefix)

	// struct syntax pattern
	matches := types.MatchStruct(typeCell)
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
		// NOTE(wenchy): treated as nested named struct
		field.Type, field.TypeDefined = p.parseType(elemType)
		field.Name = strcase.ToSnake(field.Type)
		// index := len(field.Type)
		// structName := trimmedNameCell[:index]
		structName := field.Type
		field.Options = &tableaupb.FieldOptions{
			Name: structName,
		}
		prefix += structName

		name := strings.TrimPrefix(nameCell, prefix)
		field.Fields = append(field.Fields, p.parseScalarField(name, colType, noteCell))
		for cursor++; cursor < len(header.namerow); cursor++ {
			nameCell := header.getNameCell(cursor)
			if !strings.HasPrefix(nameCell, prefix) {
				cursor--
				return cursor
			}
			cursor = p.parseSubField(field, header, cursor, prefix, nested)
		}
	}

	return cursor
}

func (p *bookParser) parseScalarField(name, typ, note string) *tableaupb.Field {
	// enum syntax pattern
	if matches := types.MatchEnum(typ); len(matches) > 0 {
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
	if strings.HasPrefix(typ, ".") {
		// This messge type is defined in imported proto
		typ = strings.TrimPrefix(typ, ".")
		return typ, true
	}
	switch typ {
	case "datetime", "date":
		typ = "google.protobuf.Timestamp"
	case "duration", "time":
		typ = "google.protobuf.Duration"
	default:
		return typ, false
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


func (gen *XmlGenerator) parseNode(nav *xmlquery.NodeNavigator, element *tableaupb.Field, prefix string) error {
	gen.nav = nav
	element.Options = &tableaupb.FieldOptions{
		Name: nav.LocalName(),
	}
	// iterate over attributes
	for _, attr := range nav.Current().Attr {
		switch strings.ToLower(attr.Name.Local) {
		case "keycol":
			tagName := strings.Split(attr.Value, ".")[0]
			attrName := strings.Split(attr.Value, ".")[1]
			keyNode := xmlquery.FindOne(nav.Current(), fmt.Sprintf("%s/@%s", tagName, attrName))
			if keyNode == nil {
				atom.Log.Panicf("KeyCol:%s not found in the immediately following nodes of %s", attr.Value, nav.LocalName())
				continue
			}
			keyType, _ := gen.guessType(keyNode.InnerText())
			element.Options.Key = attrName
			element.MapEntry = &tableaupb.MapEntry{
				KeyType: keyType,
				ValueType: tagName,
			}
			element.Type = fmt.Sprintf("map<%s, %s>", keyType, tagName)
		case "desc":
		default:
			attrName := attr.Name.Local
			attrValue := attr.Value
			t, d := gen.guessType(attrValue)
			newAttr := &tableaupb.Field{
				Options: &tableaupb.FieldOptions{
					Name: attrName,
					Default: d,
				},
				Type: t,
				Name: strcase.ToSnake(attrName),
			}
			if matches := numRegex.FindStringSubmatch(attrName); len(matches) > 0 {
				if matches[0] != "1" {
					break
				}
				newAttr.Card = "repeated"
				newAttr.Options.Name = strings.ReplaceAll(newAttr.Options.Name, matches[0], "")
				newAttr.Name = strcase.ToSnake(newAttr.Options.Name) + "_list"
			}
			element.Fields = append(element.Fields, newAttr)
		}
	}
	// fmt.Println(element)

	// createMetaStruct create a meta struct by some pre-defined tag
	tagMap := make(map[string]bool)
	createMetaStruct := func(metaNode *xmlquery.Node) {
		metaNav := xmlquery.CreateXPathNavigator(metaNode.FirstChild)
		for metaNav.NodeType() != xpath.ElementNode {
			metaNav.MoveToNext()
		}
		newChild := &tableaupb.Field{
			Options: &tableaupb.FieldOptions{
				Name: metaNav.LocalName(),
			},
			Card: "repeated",
			Type: metaNav.LocalName(),
			Name: strcase.ToSnake(metaNav.LocalName()) + "_list",
		}
		element.Fields = append(element.Fields, newChild)
		gen.parseNode(metaNav, newChild, fmt.Sprintf("%s/%s", prefix, metaNode.Parent.Data))
		gen.fieldMap[fmt.Sprintf("%s/%s/%s", prefix, metaNode.Parent.Data, metaNav.LocalName())] = newChild
		tagMap[metaNode.Data] = true
	}
	// `StructSupplement` defines the default values of one tag
	if metaNode := xmlquery.FindOne(nav.Current(), "StructSupplement"); metaNode != nil {
		createMetaStruct(metaNode)
	}
	// `StructFormatSupplement` defines a meta struct
	if metaNode := xmlquery.FindOne(nav.Current(), "StructFormatSupplement"); metaNode != nil {
		createMetaStruct(metaNode)
	}
	// iterate over child nodes
	navCopy := *nav
	for flag := navCopy.MoveToChild(); flag; flag = navCopy.MoveToNext() {
		// commentNode, documentNode and other meaningless nodes should be filtered
		if navCopy.NodeType() != xpath.ElementNode {
			continue
		}
		tagName := navCopy.LocalName()
		if _, exist := tagMap[tagName]; exist {
			continue
		}
		newChild := &tableaupb.Field{
			Options: &tableaupb.FieldOptions{
				Name: tagName,
			},
			Type: tagName,
			Name: strcase.ToSnake(tagName),
		}
		if childList := xmlquery.Find(nav.Current(), tagName); len(childList) > 1 {
			newChild.Card = "repeated"
			newChild.Name = newChild.Name + "_list"
		}
		fatherPath := fmt.Sprintf("%s/%s", prefix, nav.Current().Data)
		gen.parseNode(&navCopy, newChild, fatherPath)

		// overwrite previous meta struct if necessary
		curPath := fatherPath + "/" + tagName
		if child, exist := gen.fieldMap[curPath]; exist {
			if child.Card == "" && newChild.Card == "repeated" {
				child.Card = "repeated"
				child.Name = child.Name + "_list"
			}
			fieldMap := make(map[string]bool)
			for _, c := range child.Fields {
				fieldMap[c.Options.Name] = true
			}
			for _, c := range newChild.Fields {
				if _, exist := fieldMap[c.Options.Name]; !exist {
					child.Fields = append(child.Fields, c)
				}
			}
		} else {
			gen.fieldMap[curPath] = newChild
			element.Fields = append(element.Fields, newChild)
		}
	}

	if element.MapEntry != nil {
		element.Name = strcase.ToSnake(element.MapEntry.ValueType) + "_map"
	} else if element.Card == "repeated" {
		element.Name = strcase.ToSnake(nav.LocalName()) + "_list"
	} else {
		element.Name = strcase.ToSnake(nav.LocalName())
	}
	return nil
}

func (gen *XmlGenerator) guessType(value string) (string, string) {
	var t, d string
	if _, err := strconv.Atoi(value); err == nil {
		t, d = "int32", "0"
	} else if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		t, d = "int64", "0"
	} else {
		t, d = "string", ""
	}
	if gen.nav.Current().Parent.Data == "StructSupplement" {
		d = value
	}
	return t, d
}