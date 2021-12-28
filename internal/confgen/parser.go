package confgen

import (
	"strconv"
	"strings"
	"time"

	"github.com/Wenchy/tableau/internal/atom"
	"github.com/Wenchy/tableau/internal/confgen/mexporter"
	"github.com/Wenchy/tableau/internal/excel"
	"github.com/Wenchy/tableau/internal/printer"
	"github.com/Wenchy/tableau/options"
	"github.com/Wenchy/tableau/proto/tableaupb"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type sheetExporter struct {
	OutputDir string
	Output    *options.OutputOption // output settings.

}

func NewSheetExporter(outputDir string, output *options.OutputOption) *sheetExporter {
	return &sheetExporter{
		OutputDir: outputDir,
		Output:    output,
	}
}

// export the protomsg message.
func (x *sheetExporter) Export(book *excel.Book, parser *sheetParser, protomsg proto.Message) error {
	md := protomsg.ProtoReflect().Descriptor()
	// _, workbook := parseFileOptions(md.ParentFile())
	msgName, wsOpts := parseMessageOptions(md)

	sheet := book.Sheets[wsOpts.Name]
	if sheet == nil {
		return errors.Errorf("not found worksheet: %s", wsOpts.Name)
	}

	if err := parser.Parse(protomsg, sheet, wsOpts); err != nil {
		return errors.WithMessage(err, "failed to parse sheet")
	}

	exporter := mexporter.New(msgName, protomsg, x.OutputDir, x.Output)
	return exporter.Export()
}

type sheetParser struct {
	ProtoPackage string
	LocationName string
}

func NewSheetParser(protoPackage, locationName string) *sheetParser {
	return &sheetParser{
		ProtoPackage: protoPackage,
		LocationName: locationName,
	}
}

func (sp *sheetParser) Parse(protomsg proto.Message, sheet *excel.Sheet, wsOpts *tableaupb.WorksheetOptions) error {
	atom.Log.Debugf("parse sheet: %s", sheet.Name)
	msg := protomsg.ProtoReflect()
	if wsOpts.Transpose {
		// interchange the rows and columns
		// namerow: name column
		// [datarow, MaxCol]: data column
		// kvRow := make(map[string]string)
		for col := int(wsOpts.Datarow) - 1; col < sheet.MaxCol; col++ {
			rc := excel.NewRowCells(col)
			for row := 0; row < sheet.MaxRow; row++ {
				nameCol := int(wsOpts.Namerow) - 1
				nameCell, err := sheet.Cell(row, nameCol)
				if err != nil {
					return errors.WithMessagef(err, "failed to get name cell: %d, %d", row, nameCol)
				}
				name := excel.ExtractFromCell(nameCell, wsOpts.Nameline)
				data, err := sheet.Cell(row, col)
				if err != nil {
					return errors.WithMessagef(err, "failed to get data cell: %d, %d", row, col)
				}
				rc.SetCell(name, row, data)
			}
			err := sp.parseFieldOptions(msg, rc, 0, "")
			if err != nil {
				return err
			}
		}
	} else {
		// namerow: name row
		// [datarow, MaxRow]: data row
		for row := int(wsOpts.Datarow) - 1; row < sheet.MaxRow; row++ {
			rc := excel.NewRowCells(row)
			for col := 0; col < sheet.MaxCol; col++ {
				nameRow := int(wsOpts.Namerow) - 1
				nameCell, err := sheet.Cell(nameRow, col)
				if err != nil {
					return errors.WithMessagef(err, "failed to get name cell: %d, %d", nameRow, col)
				}
				name := excel.ExtractFromCell(nameCell, wsOpts.Nameline)
				data, err := sheet.Cell(row, col)
				if err != nil {
					return errors.WithMessagef(err, "failed to get data cell: %d, %d", row, col)
				}
				rc.SetCell(name, col, data)
			}
			err := sp.parseFieldOptions(msg, rc, 0, "")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Field struct {
	fd   protoreflect.FieldDescriptor
	opts *tableaupb.FieldOptions
}

// parseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func (sp *sheetParser) parseFieldOptions(msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	md := msg.Descriptor()
	pkg := md.ParentFile().Package()
	// opts := md.Options().(*descriptorpb.MessageOptions)
	// worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
	// worksheetName := ""
	// if worksheet != nil {
	// 	worksheetName = worksheet.Name
	// }
	// atom.Log.Debugf("%s// %s, '%s', %v, %v, %v", printer.Indent(depth), md.FullName(), worksheetName, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if string(pkg) != sp.ProtoPackage && pkg != "google.protobuf" {
			atom.Log.Debugf("%s// no need to proces package: %v", printer.Indent(depth), pkg)
			return nil
		}

		// default value
		name := strcase.ToCamel(string(fd.FullName().Name()))
		note := ""
		etype := tableaupb.Type_TYPE_DEFAULT
		key := ""
		layout := tableaupb.Layout_LAYOUT_DEFAULT
		sep := ""
		subsep := ""
		optional := false

		opts := fd.Options().(*descriptorpb.FieldOptions)
		fieldOpts := proto.GetExtension(opts, tableaupb.E_Field).(*tableaupb.FieldOptions)
		if fieldOpts != nil {
			name = fieldOpts.Name
			note = fieldOpts.Note
			etype = fieldOpts.Type
			key = fieldOpts.Key
			layout = fieldOpts.Layout
			sep = fieldOpts.Sep
			subsep = fieldOpts.Subsep
			optional = fieldOpts.Optional
		} else {
			// default processing
			if fd.IsList() {
				// truncate suffix `List` (CamelCase) corresponding to `_list` (snake_case)
				name = strings.TrimSuffix(name, "List")
			} else if fd.IsMap() {
				// truncate suffix `Map` (CamelCase) corresponding to `_map` (snake_case)
				// name = strings.TrimSuffix(name, "Map")
				name = ""
				key = "Key"
			}
		}
		if sep == "" {
			sep = ","
		}
		if subsep == "" {
			subsep = ":"
		}

		// msgName := ""
		// if fd.Kind() == protoreflect.MessageKind {
		// 	msgName = string(fd.Message().FullName())
		// }
		// atom.Log.Debugf("%s%s(%v) %s(%s) %s = %d [(name) = \"%s\", (type) = %s, (key) = \"%s\", (layout) = \"%s\", (sep) = \"%s\"];",
		// 	printer.Indent(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+name, etype.String(), key, layout.String(), sep)
		// atom.Log.Debugw("field metadata",
		// 	"tabs", depth,
		// 	"cardinality", fd.Cardinality().String(),
		// 	"isMap", fd.IsMap(),
		// 	"kind", fd.Kind().String(),
		// 	"msgName", msgName,
		// 	"fullName", fd.FullName(),
		// 	"number", fd.Number(),
		// 	"name", prefix+name,
		// 	"note", note,
		// 	"type", etype.String(),
		// 	"key", key,
		// 	"layout", layout.String(),
		// 	"sep", sep,
		// )

		field := &Field{
			fd: fd,
			opts: &tableaupb.FieldOptions{
				Name:     name,
				Note:     note,
				Type:     etype,
				Key:      key,
				Layout:   layout,
				Sep:      sep,
				Subsep:   subsep,
				Optional: optional,
			},
		}
		err = sp.parseField(field, msg, rc, depth, prefix)
		if err != nil {
			return errors.WithMessagef(err, "failed to parse field: %s, opts: %v", fd.FullName().Name(), field.opts)
		}
	}
	return nil
}

func (sp *sheetParser) parseField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	// atom.Log.Debug(field.fd.ContainingMessage().FullName())
	if field.fd.IsMap() {
		return sp.parseMapField(field, msg, rc, depth, prefix)
	} else if field.fd.IsList() {
		return sp.parseListField(field, msg, rc, depth, prefix)
	} else if field.fd.Kind() == protoreflect.MessageKind {
		return sp.parseStructField(field, msg, rc, depth, prefix)
	} else if field.fd.Kind() == protoreflect.EnumKind {
		return sp.parseEnumField(field, msg, rc, depth, prefix)
	} else {
		return sp.parseScalarField(field, msg, rc, depth, prefix)
	}
}

func (sp *sheetParser) parseMapField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	// Mutable returns a mutable reference to a composite type.
	newValue := msg.Mutable(field.fd)
	reflectMap := newValue.Map()
	// reflectMap := msg.Mutable(field.fd).Map()
	keyFd := field.fd.MapKey()
	valueFd := field.fd.MapValue()
	// newKey := protoreflect.ValueOf(int32(1)).MapKey()
	// newKey := sp.parseFieldValue(keyFd, "1111001").MapKey()
	if field.opts.Type == tableaupb.Type_TYPE_INCELL_MAP {
		colName := prefix + field.opts.Name
		cell := rc.Cell(colName, field.opts.Optional)
		if cell == nil {
			return errors.Errorf("%s|column not found", rc.CellDebugString(colName))
		}
		if valueFd.Kind() == protoreflect.MessageKind {
			return errors.Errorf("%s|incell map: message type not supported", rc.CellDebugString(colName))
		}

		if cell.Data != "" {
			// If s does not contain sep and sep is not empty, Split returns a
			// slice of length 1 whose only element is s.
			splits := strings.Split(cell.Data, field.opts.Sep)
			for _, pair := range splits {
				kv := strings.Split(pair, field.opts.Subsep)
				if len(kv) != 2 {
					return errors.Errorf("%s|incell map: illegal key-value pair: %s", rc.CellDebugString(colName), pair)
				}
				key, value := kv[0], kv[1]

				fieldValue, err := sp.parseFieldValue(keyFd, key)
				if err != nil {
					return errors.WithMessagef(err, "%s|incell map: failed to parse field value: %s", rc.CellDebugString(colName), key)
				}
				newMapKey := fieldValue.MapKey()

				fieldValue, err = sp.parseFieldValue(valueFd, value)
				if err != nil {
					return errors.WithMessagef(err, "%s|incell map: failed to parse field value: %s", rc.CellDebugString(colName), value)
				}
				newMapValue := reflectMap.NewValue()
				newMapValue = fieldValue

				reflectMap.Set(newMapKey, newMapValue)
			}
		}
	} else {
		emptyMapValue := reflectMap.NewValue()
		if valueFd.Kind() == protoreflect.MessageKind {
			if field.opts.Layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
				size := rc.GetCellCountWithPrefix(prefix + field.opts.Name)
				// atom.Log.Debug("prefix size: ", size)
				for i := 1; i <= size; i++ {
					keyColName := prefix + field.opts.Name + strconv.Itoa(i) + field.opts.Key
					cell := rc.Cell(keyColName, field.opts.Optional)
					if cell == nil {
						return errors.Errorf("%s|horizontal map: key column not found", rc.CellDebugString(keyColName))
					}
					fieldValue, err := sp.parseFieldValue(keyFd, cell.Data)
					if err != nil {
						return errors.WithMessagef(err, "%s|horizontal map: failed to parse field value: %s", rc.CellDebugString(keyColName), cell.Data)
					}
					newMapKey := fieldValue.MapKey()

					var newMapValue protoreflect.Value
					if reflectMap.Has(newMapKey) {
						newMapValue = reflectMap.Mutable(newMapKey)
					} else {
						newMapValue = reflectMap.NewValue()
					}
					err = sp.parseFieldOptions(newMapValue.Message(), rc, depth+1, prefix+field.opts.Name+strconv.Itoa(i))
					if err != nil {
						return errors.WithMessagef(err, "%s|horizontal map: failed to parse field options with prefix: %s", rc.CellDebugString(keyColName), prefix+field.opts.Name+strconv.Itoa(i))
					}
					if !MessageValueEqual(emptyMapValue, newMapValue) {
						reflectMap.Set(newMapKey, newMapValue)
					}
				}
			} else {
				keyColName := prefix + field.opts.Name + field.opts.Key
				cell := rc.Cell(keyColName, field.opts.Optional)
				if cell == nil {
					return errors.Errorf("%s|vertical map: key column not found", rc.CellDebugString(keyColName))
				}
				fieldValue, err := sp.parseFieldValue(keyFd, cell.Data)
				if err != nil {
					return errors.WithMessagef(err, "%s|vertical map: failed to parse field value: %s", rc.CellDebugString(keyColName), cell.Data)
				}
				newMapKey := fieldValue.MapKey()

				var newMapValue protoreflect.Value
				if reflectMap.Has(newMapKey) {
					newMapValue = reflectMap.Mutable(newMapKey)
				} else {
					newMapValue = reflectMap.NewValue()
				}
				err = sp.parseFieldOptions(newMapValue.Message(), rc, depth+1, prefix+field.opts.Name)
				if err != nil {
					return errors.WithMessagef(err, "%s|vertical map: failed to parse field options with prefix: %s", rc.CellDebugString(keyColName), prefix+field.opts.Name)
				}
				if !MessageValueEqual(emptyMapValue, newMapValue) {
					reflectMap.Set(newMapKey, newMapValue)
				}
			}
		} else {
			// value is scalar type
			key := "Key"     // default key name
			value := "Value" // default value name
			// key cell
			keyColName := prefix + field.opts.Name + key
			cell := rc.Cell(keyColName, field.opts.Optional)
			if cell == nil {
				return errors.Errorf("%s|vertical map(scalar): key column not found", rc.CellDebugString(keyColName))
			}
			fieldValue, err := sp.parseFieldValue(keyFd, cell.Data)
			if err != nil {
				return errors.WithMessagef(err, "%s|failed to parse field value: %s", rc.CellDebugString(keyColName), cell.Data)
			}
			newMapKey := fieldValue.MapKey()
			var newMapValue protoreflect.Value
			if reflectMap.Has(newMapKey) {
				newMapValue = reflectMap.Mutable(newMapKey)
			} else {
				newMapValue = reflectMap.NewValue()
			}
			// value cell
			valueColName := prefix + field.opts.Name + value
			cell = rc.Cell(valueColName, field.opts.Optional)
			if cell == nil {
				return errors.Errorf("%s|vertical map(scalar): value colum not found", rc.CellDebugString(valueColName))
			}
			newMapValue, err = sp.parseFieldValue(field.fd, cell.Data)
			if err != nil {
				return errors.WithMessagef(err, "%s|vertical map(scalar): failed to parse field value: %s", rc.CellDebugString(valueColName), cell.Data)
			}
			if !reflectMap.Has(newMapKey) {
				reflectMap.Set(newMapKey, newMapValue)
			}
		}
	}
	if !msg.Has(field.fd) && reflectMap.Len() != 0 {
		msg.Set(field.fd, newValue)
	}
	return nil
}

func (sp *sheetParser) parseListField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	// Mutable returns a mutable reference to a composite type.
	newValue := msg.Mutable(field.fd)
	reflectList := newValue.List()
	if field.opts.Type == tableaupb.Type_TYPE_INCELL_LIST {
		// incell list
		colName := prefix + field.opts.Name
		cell := rc.Cell(colName, field.opts.Optional)
		if cell == nil {
			return errors.Errorf("%s|incell list: column not found", rc.CellDebugString(colName))
		}
		if cell.Data != "" {
			// If s does not contain sep and sep is not empty, Split returns a
			// slice of length 1 whose only element is s.
			splits := strings.Split(cell.Data, field.opts.Sep)
			for _, incell := range splits {
				value, err := sp.parseFieldValue(field.fd, incell)
				if err != nil {
					return errors.WithMessagef(err, "%s|incell list: failed to parse field value: %s", rc.CellDebugString(colName), incell)
				}
				reflectList.Append(value)
			}
		}
	} else {
		emptyListValue := reflectList.NewElement()
		if field.opts.Layout == tableaupb.Layout_LAYOUT_VERTICAL {
			newListValue := reflectList.NewElement()
			if field.fd.Kind() == protoreflect.MessageKind {
				// struct list
				err = sp.parseFieldOptions(newListValue.Message(), rc, depth+1, prefix+field.opts.Name)
				if err != nil {
					return errors.WithMessagef(err, "...|vertical list: failed to parse field options with prefix: %s", prefix+field.opts.Name)
				}
				if !MessageValueEqual(emptyListValue, newListValue) {
					reflectList.Append(newListValue)
				}
			} else {
				// TODO: support list of scalar type when layout is vertical?
				// NOTE(wenchyzhu): we don't support list of scalar type when layout is vertical
			}
		} else {
			size := rc.GetCellCountWithPrefix(prefix + field.opts.Name)
			if size <= 0 {
				return errors.Errorf("%s|horizontal list: no cell found with digit suffix", rc.CellDebugString(prefix+field.opts.Name))
			}
			for i := 1; i <= size; i++ {
				newListValue := reflectList.NewElement()
				if field.fd.Kind() == protoreflect.MessageKind {
					// struct list
					err = sp.parseFieldOptions(newListValue.Message(), rc, depth+1, prefix+field.opts.Name+strconv.Itoa(i))
					if err != nil {
						return errors.WithMessagef(err, "...|horizontal list: failed to parse field options with prefix: %s", prefix+field.opts.Name+strconv.Itoa(i))
					}
					if !MessageValueEqual(emptyListValue, newListValue) {
						reflectList.Append(newListValue)
					}
				} else {
					// scalar list
					colName := prefix + field.opts.Name + strconv.Itoa(i)
					cell := rc.Cell(colName, field.opts.Optional)
					if cell == nil {
						return errors.Errorf("%s|horizontal list(scalar): column not found", rc.CellDebugString(colName))
					}
					newListValue, err = sp.parseFieldValue(field.fd, cell.Data)
					if err != nil {
						return errors.WithMessagef(err, "%s|horizontal list(scalar): failed to parse field value: %s", rc.CellDebugString(colName), cell.Data)
					}
					reflectList.Append(newListValue)
				}
			}
		}
	}
	if !msg.Has(field.fd) && reflectList.Len() != 0 {
		msg.Set(field.fd, newValue)
	}

	return nil
}

func (sp *sheetParser) parseStructField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	// NOTE(wenchy): `proto.Equal` treats a nil message as not equal to an empty one.
	// doc: [Equal](https://pkg.go.dev/google.golang.org/protobuf/proto?tab=doc#Equal)
	// issue: [APIv2: protoreflect: consider Message nilness test](https://github.com/golang/protobuf/issues/966)
	// ```
	// nilMessage = (*MyMessage)(nil)
	// emptyMessage = new(MyMessage)
	//
	// Equal(nil, nil)                   // true
	// Equal(nil, nilMessage)            // false
	// Equal(nil, emptyMessage)          // false
	// Equal(nilMessage, nilMessage)     // true
	// Equal(nilMessage, emptyMessage)   // ??? false
	// Equal(emptyMessage, emptyMessage) // true
	// ```
	//
	// Case: `subMsg := msg.Mutable(fd).Message()`
	// `Message.Mutable` will allocate new "empty message", and is not equal to "nil"
	//
	// Solution:
	// 1. spawn two values: `emptyValue` and `newValue`
	// 2. set `newValue` back to field if `newValue` is not equal to `emptyValue`
	emptyValue := msg.NewField(field.fd)
	newValue := msg.NewField(field.fd)

	colName := prefix + field.opts.Name
	if field.opts.Type == tableaupb.Type_TYPE_INCELL_STRUCT {
		// incell struct
		cell := rc.Cell(colName, field.opts.Optional)
		if cell == nil {
			return errors.Errorf("%s|incell struct: column not found", rc.CellDebugString(colName))
		}
		if cell.Data != "" {
			// If s does not contain sep and sep is not empty, Split returns a
			// slice of length 1 whose only element is s.
			splits := strings.Split(cell.Data, field.opts.Sep)
			subMd := newValue.Message().Descriptor()
			for i := 0; i < subMd.Fields().Len() && i < len(splits); i++ {
				fd := subMd.Fields().Get(i)
				// atom.Log.Debugf("fd.FullName().Name(): ", fd.FullName().Name())
				incell := splits[i]
				value, err := sp.parseFieldValue(fd, incell)
				if err != nil {
					return errors.WithMessagef(err, "%s|incell struct: failed to parse field value: %s", rc.CellDebugString(colName), incell)
				}
				newValue.Message().Set(fd, value)
			}
		}
	} else {
		// built-in struct
		subMsgName := string(field.fd.Message().FullName())
		_, found := specialMessageMap[subMsgName]
		if found {
			cell := rc.Cell(colName, field.opts.Optional)
			if cell == nil {
				return errors.Errorf("%s|builtin type: column not found", rc.CellDebugString(colName))
			}
			newValue, err = sp.parseFieldValue(field.fd, cell.Data)
			if err != nil {
				return errors.WithMessagef(err, "%s|builtin type: failed to parse field value: %s", rc.CellDebugString(colName), cell.Data)
			}
		} else {
			pkgName := newValue.Message().Descriptor().ParentFile().Package()
			if string(pkgName) != sp.ProtoPackage {
				return errors.Errorf("%s|builtin type: unknown message %v in package %s", rc.CellDebugString(colName), subMsgName, pkgName)
			}
			err = sp.parseFieldOptions(newValue.Message(), rc, depth+1, prefix+field.opts.Name)
			if err != nil {
				return errors.WithMessagef(err, "%s|builtin type: failed to parse field options with prefix: %s", rc.CellDebugString(colName), prefix+field.opts.Name)
			}
		}
	}
	if !MessageValueEqual(emptyValue, newValue) {
		msg.Set(field.fd, newValue)
	}
	return nil
}

func (sp *sheetParser) parseEnumField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	newValue := msg.NewField(field.fd)
	colName := prefix + field.opts.Name
	cell := rc.Cell(colName, field.opts.Optional)
	if cell == nil {
		return errors.Errorf("%s|enum: column not found", rc.CellDebugString(colName))
	}
	newValue, err = func() (protoreflect.Value, error) {
		// default enum value
		defaultValue := protoreflect.ValueOfEnum(protoreflect.EnumNumber(0))
		if cell.Data == "" {
			return defaultValue, nil
		}
		ed := field.fd.Enum() // get enum descriptor
		// try enum value number
		val, err := strconv.ParseInt(cell.Data, 10, 32)
		if err == nil {
			evd := ed.Values().ByNumber(protoreflect.EnumNumber(val))
			if evd != nil {
				return protoreflect.ValueOfEnum(evd.Number()), nil
			}
			return defaultValue, errors.Errorf("%s|enum: enum type(%s) has no enum value name: %s", rc.CellDebugString(colName), field.opts.Type, cell.Data)
		}

		// try enum value name
		evd := ed.Values().ByName(protoreflect.Name(cell.Data))
		if evd != nil {
			return protoreflect.ValueOfEnum(evd.Number()), nil
		}
		// try enum value alias name
		for i := 0; i < ed.Values().Len(); i++ {
			// get enum value descriptor
			evd := ed.Values().Get(i)
			opts := evd.Options().(*descriptorpb.EnumValueOptions)
			evalueOpts := proto.GetExtension(opts, tableaupb.E_Evalue).(*tableaupb.EnumValueOptions)
			if evalueOpts == nil {
				return defaultValue, errors.Errorf("%s|enum: enum type(%s) has no enum value options", rc.CellDebugString(colName), field.opts.Type)
			}
			if evalueOpts.Name == cell.Data {
				return protoreflect.ValueOfEnum(evd.Number()), nil
			}
		}
		return defaultValue, errors.Errorf("%s|enum: enum value alias name(%s) not found", rc.CellDebugString(colName), cell.Data)
	}()

	if err != nil {
		return errors.WithMessagef(err, "%s|enum: failed to parse field value: %s", rc.CellDebugString(colName), cell.Data)
	}
	msg.Set(field.fd, newValue)
	return nil
}

func (sp *sheetParser) parseScalarField(field *Field, msg protoreflect.Message, rc *excel.RowCells, depth int, prefix string) (err error) {
	newValue := msg.NewField(field.fd)
	colName := prefix + field.opts.Name
	cell := rc.Cell(colName, field.opts.Optional)
	if cell == nil {
		return errors.Errorf("%s|scalar: column not found", rc.CellDebugString(colName))
	}
	newValue, err = sp.parseFieldValue(field.fd, cell.Data)
	if err != nil {
		return errors.WithMessagef(err, "%s|scalar: failed to parse field value: %s", rc.CellDebugString(colName), cell.Data)
	}
	msg.Set(field.fd, newValue)
	return nil
}

func (sp *sheetParser) parseFieldValue(fd protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if value != "" {
			val, err := strconv.ParseInt(value, 10, 32)
			return protoreflect.ValueOf(int32(val)), err
		}
		return protoreflect.ValueOf(int32(0)), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if value != "" {
			val, err := strconv.ParseUint(value, 10, 32)
			return protoreflect.ValueOf(uint32(val)), err
		}
		return protoreflect.ValueOf(uint32(0)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if value != "" {
			val, err := strconv.ParseInt(value, 10, 64)
			return protoreflect.ValueOf(int64(val)), err
		}
		return protoreflect.ValueOf(int64(0)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if value != "" {
			val, err := strconv.ParseUint(value, 10, 64)
			return protoreflect.ValueOf(uint64(val)), err
		}
		return protoreflect.ValueOf(uint64(0)), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOf(value), nil
	case protoreflect.BytesKind:
		return protoreflect.ValueOf([]byte(value)), nil
	case protoreflect.BoolKind:
		if value != "" {
			val, err := strconv.ParseBool(value)
			return protoreflect.ValueOf(val), err
		}
		return protoreflect.ValueOf(false), nil
	case protoreflect.FloatKind:
		if value != "" {
			val, err := strconv.ParseFloat(value, 32)
			return protoreflect.ValueOf(float32(val)), err
		}
		return protoreflect.ValueOf(float32(0)), nil
	case protoreflect.DoubleKind:
		if value != "" {
			val, err := strconv.ParseFloat(value, 64)
			return protoreflect.ValueOf(float64(val)), err
		}
		return protoreflect.ValueOf(float64(0)), nil
	// case protoreflect.EnumKind:
	// 	atom.Log.Panicf("not supported key type: %s", fd.Kind().String())
	// case protoreflect.GroupKind:
	// 	atom.Log.Panicf("not supported key type: %s", fd.Kind().String())
	// 	return protoreflect.Value{}
	case protoreflect.MessageKind:
		msgName := fd.Message().FullName()
		switch msgName {
		case "google.protobuf.Timestamp":
			// make use of t as a *timestamppb.Timestamp
			ts := &timestamppb.Timestamp{}
			if value != "" {
				// location name: "Asia/Shanghai" or "Asia/Chongqing".
				// NOTE(wenchy): There is no "Asia/Beijing" location name. Whoa!!! Big surprize?
				t, err := parseTimeWithLocation(sp.LocationName, value)
				if err != nil {
					return protoreflect.ValueOf(ts.ProtoReflect()), errors.WithMessagef(err, "illegal timestamp string format: %v", value)
				}
				// atom.Log.Debugf("timeStr: %v, unix timestamp: %v", value, t.Unix())
				ts = timestamppb.New(t)
				if err := ts.CheckValid(); err != nil {
					return protoreflect.ValueOf(ts.ProtoReflect()), errors.WithMessagef(err, "invalid timestamp: %v", value)
				}
			}
			return protoreflect.ValueOf(ts.ProtoReflect()), nil
		case "google.protobuf.Duration":
			// make use of d as a *durationpb.Duration
			du := &durationpb.Duration{} // default
			if value != "" {
				d, err := parseDuration(value)
				if err != nil {
					return protoreflect.ValueOf(du.ProtoReflect()), errors.WithMessagef(err, "illegal duration string format: %v", value)
				}
				du = durationpb.New(d)
				if err := du.CheckValid(); err != nil {
					return protoreflect.ValueOf(du.ProtoReflect()), errors.WithMessagef(err, "invalid duration: %v", value)
				}
			}
			return protoreflect.ValueOf(du.ProtoReflect()), nil

		default:
			return protoreflect.Value{}, errors.Errorf("not supported message type: %s", msgName)
		}
	default:
		return protoreflect.Value{}, errors.Errorf("not supported scalar type: %s", fd.Kind().String())
	}
}

// parseFileOptions is aimed to parse the options of a protobuf definition file.
func parseFileOptions(fd protoreflect.FileDescriptor) (string, *tableaupb.WorkbookOptions) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	protofile := string(fd.FullName())
	workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(*tableaupb.WorkbookOptions)
	return protofile, workbook
}

// parseMessageOptions is aimed to parse the options of a protobuf message.
func parseMessageOptions(md protoreflect.MessageDescriptor) (string, *tableaupb.WorksheetOptions) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgName := string(md.Name())
	wsOpts := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
	if wsOpts.Namerow == 0 {
		wsOpts.Namerow = 1 // default
	}
	if wsOpts.Typerow == 0 {
		wsOpts.Typerow = 2 // default
	}

	if wsOpts.Noterow == 0 {
		wsOpts.Noterow = 3 // default
	}

	if wsOpts.Datarow == 0 {
		wsOpts.Datarow = 4 // default
	}
	// atom.Log.Debugf("msg: %v, wsOpts: %+v", msgName, wsOpts)
	return msgName, wsOpts
}

func parseTimeWithLocation(locationName string, timeStr string) (time.Time, error) {
	// see https://golang.org/pkg/time/#LoadLocation
	if location, err := time.LoadLocation(locationName); err != nil {
		return time.Time{}, errors.Wrapf(err, "LoadLocation failed: %s", locationName)
	} else {
		timeStr = strings.TrimSpace(timeStr)
		layout := "2006-01-02 15:04:05"
		if strings.Contains(timeStr, " ") {
			layout = "2006-01-02 15:04:05"
		} else {
			layout = "2006-01-02"
			if !strings.Contains(timeStr, "-") && len(timeStr) == 8 {
				// convert "yyyymmdd" to "yyyy-mm-dd"
				timeStr = timeStr[0:4] + "-" + timeStr[4:6] + "-" + timeStr[6:8]
			}
		}
		t, err := time.ParseInLocation(layout, timeStr, location)
		if err != nil {
			return time.Time{}, errors.Wrapf(err, "ParseInLocation failed, timeStr: %v, locationName: %v", timeStr, locationName)
		}
		return t, nil
	}
}

func parseDuration(duration string) (time.Duration, error) {
	duration = strings.TrimSpace(duration)
	if !strings.ContainsAny(duration, ":hmsµu") && len(duration) == 6 {
		duration = duration[0:2] + "h" + duration[2:4] + "m" + duration[4:6] + "s"
	} else if strings.Contains(duration, ":") && len(duration) == 8 {
		// convert "hh:mm:ss" to "<hh>h<mm>m:<ss>s"
		duration = duration[0:2] + "h" + duration[3:5] + "m" + duration[6:8] + "s"
	}

	return time.ParseDuration(duration)
}

func MessageValueEqual(v1, v2 protoreflect.Value) bool {
	if proto.Equal(v1.Message().Interface(), v2.Message().Interface()) {
		// atom.Log.Debug("empty message exists")
		return true
	}
	return false
}
