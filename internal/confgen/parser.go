package confgen

import (
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

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

type sheetParser struct {
	ProtoPackage string
	LocationName string
	InputDir     string
	OutputDir    string
	Output       *options.OutputOption // output settings.

	protomsg proto.Message
}
// export the protomsg message.
func (sh *sheetParser) Export() error {
	md := sh.protomsg.ProtoReflect().Descriptor()
	msg := sh.protomsg.ProtoReflect()
	_, workbook := parseFileOptions(md.ParentFile())
	msgName, worksheetName, namerow, _, datarow, transpose := parseMessageOptions(md)
	// msgName, worksheetName, namerow, noterow, datarow, transpose := parseMessageOptions(md)

	wbPath := sh.InputDir + workbook.Name
	book, err := excel.NewBook(wbPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create new workbook: %s", wbPath)
	}

	sheet, ok := book.Sheets[worksheetName]
	if !ok {
		return errors.Wrapf(err, "not found worksheet: %s", worksheetName)
	}

	if transpose {
		// interchange the rows and columns
		// namerow: name column
		// [datarow, MaxCol]: data column
		kvRow := make(map[string]string)
		for col := int(datarow) - 1; col < sheet.MaxCol; col++ {
			for row := 0; row < sheet.MaxRow; row++ {
				nameCol := int(namerow) - 1
				name, err := sheet.Cell(row, nameCol)
				if err != nil {
					return errors.Wrapf(err, "failed to get name cell: %d, %d", row, nameCol)
				}
				name = clearNewline(name)
				data, err := sheet.Cell(row, col)
				if err != nil {
					return errors.Wrapf(err, "failed to get data cell: %d, %d", row, col)
				}
				kvRow[name] = data
			}
			sh.parseFieldOptions(msg, kvRow, 0, "")
		}
	} else {
		// namerow: name row
		// [datarow, MaxRow]: data row
		kvRow := make(map[string]string)
		for row := int(datarow) - 1; row < sheet.MaxRow; row++ {
			for col := 0; col < sheet.MaxCol; col++ {
				nameRow := int(namerow) - 1
				name, err := sheet.Cell(nameRow, col)
				if err != nil {
					return errors.Wrapf(err, "failed to get name cell: %d, %d", nameRow, col)
				}
				name = clearNewline(name)
				data, err := sheet.Cell(row, col)
				if err != nil {
					return errors.Wrapf(err, "failed to get data cell: %d, %d", row, col)
				}
				kvRow[name] = data
			}
			sh.parseFieldOptions(msg, kvRow, 0, "")
		}
	}
	x := mexporter.New(msgName, sh.protomsg, sh.OutputDir, sh.Output)
	return x.Export()
}

// parseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func (sh *sheetParser) parseFieldOptions(msg protoreflect.Message, row map[string]string, depth int, prefix string) (err error) {
	md := msg.Descriptor()
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
	worksheetName := ""
	if worksheet != nil {
		worksheetName = worksheet.Name
	}

	pkg := md.ParentFile().Package()
	atom.Log.Debugf("%s// %s, '%s', %v, %v, %v", printer.Indent(depth), md.FullName(), worksheetName, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if string(pkg) != sh.ProtoPackage && pkg != "google.protobuf" {
			atom.Log.Debugf("%s// no need to proces package: %v", printer.Indent(depth), pkg)
			return nil
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
			// atom.Log.Debug(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			// ParseFieldOptions(fd.Message(), depth+1)
		}

		// if fd.IsList() {
		// 	atom.Log.Debug("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }

		// default value
		name := strcase.ToCamel(string(fd.FullName().Name()))
		etype := tableaupb.Type_TYPE_DEFAULT
		key := ""
		layout := tableaupb.Layout_LAYOUT_DEFAULT
		sep := ""
		subsep := ""

		opts := fd.Options().(*descriptorpb.FieldOptions)
		field := proto.GetExtension(opts, tableaupb.E_Field).(*tableaupb.FieldOptions)
		if field != nil {
			name = field.Name
			etype = field.Type
			key = field.Key
			layout = field.Layout
			sep = field.Sep
			subsep = field.Subsep
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
		atom.Log.Debugf("%s%s(%v) %s(%s) %s = %d [(name) = \"%s\", (type) = %s, (key) = \"%s\", (layout) = \"%s\", (sep) = \"%s\"];",
			printer.Indent(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+name, etype.String(), key, layout.String(), sep)
		atom.Log.Debugw("field metadata",
			"tabs", depth,
			"cardinality", fd.Cardinality().String(),
			"isMap", fd.IsMap(),
			"kind", fd.Kind().String(),
			"msgName", msgName,
			"fullName", fd.FullName(),
			"number", fd.Number(),
			"name", prefix+name,
			"type", etype.String(),
			"key", key,
			"layout", layout.String(),
			"sep", sep,
		)

		// atom.Log.Debug(fd.ContainingMessage().FullName())

		// if fd.Cardinality() == protoreflect.Repeated && fd.Kind() == protoreflect.MessageKind {
		// 	msg := fd.Message().New()
		// }

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
		emptyValue := msg.NewField(fd)
		newValue := msg.NewField(fd)
		if fd.IsMap() {
			// Mutable returns a mutable reference to a composite type.
			if msg.Has(fd) {
				newValue = msg.Mutable(fd)
			}
			reflectMap := newValue.Map()
			// reflectMap := msg.Mutable(fd).Map()
			keyFd := fd.MapKey()
			valueFd := fd.MapValue()
			// newKey := protoreflect.ValueOf(int32(1)).MapKey()
			// newKey := sh.parseFieldValue(keyFd, "1111001").MapKey()
			if etype == tableaupb.Type_TYPE_INCELL_MAP {
				if valueFd.Kind() == protoreflect.MessageKind {
					atom.Log.Panicf("in-cell map do not support value as message type")
				}
				cellValue, ok := row[prefix+name]
				if !ok {
					atom.Log.Panicf("column name not found: %s", prefix+name)
				}
				if cellValue != "" {
					// If s does not contain sep and sep is not empty, Split returns a
					// slice of length 1 whose only element is s.
					splits := strings.Split(cellValue, sep)
					for _, pair := range splits {
						kv := strings.Split(pair, subsep)
						if len(kv) != 2 {
							atom.Log.Panicf("illegal key-value pair: %v, %v", prefix+name, pair)
						}
						key, value := kv[0], kv[1]

						fieldValue, err := sh.parseFieldValue(keyFd, key)
						if err != nil {
							return errors.Wrapf(err, "failed to parse field value: %v", key)
						}
						newMapKey := fieldValue.MapKey()

						fieldValue, err = sh.parseFieldValue(valueFd, value)
						if err != nil {
							return errors.Wrapf(err, "failed to parse field value: %v", value)
						}
						newMapValue := reflectMap.NewValue()
						newMapValue = fieldValue

						reflectMap.Set(newMapKey, newMapValue)
					}
				}
			} else {
				emptyMapValue := reflectMap.NewValue()
				if valueFd.Kind() == protoreflect.MessageKind {
					if layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
						size := getPrefixSize(row, prefix+name)
						// atom.Log.Debug("prefix size: ", size)
						for i := 1; i <= size; i++ {
							cellValue, ok := row[prefix+name+strconv.Itoa(i)+key]
							if !ok {
								atom.Log.Panicf("key not found: %s", prefix+name+key)
							}
							fieldValue, err := sh.parseFieldValue(keyFd, cellValue)
							if err != nil {
								return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
							}
							newMapKey := fieldValue.MapKey()

							var newMapValue protoreflect.Value
							if reflectMap.Has(newMapKey) {
								newMapValue = reflectMap.Mutable(newMapKey)
							} else {
								newMapValue = reflectMap.NewValue()
							}
							sh.parseFieldOptions(newMapValue.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
							if !MessageValueEqual(emptyMapValue, newMapValue) {
								reflectMap.Set(newMapKey, newMapValue)
							}
						}
					} else {
						newMapKey := keyFd.Default().MapKey()
						cellValue, ok := row[prefix+name+key]
						if !ok {
							atom.Log.Panicf("key not found: %s", prefix+name+key)
						}
						fieldValue, err := sh.parseFieldValue(keyFd, cellValue)
						if err != nil {
							return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
						}
						newMapKey = fieldValue.MapKey()

						var newMapValue protoreflect.Value
						if reflectMap.Has(newMapKey) {
							newMapValue = reflectMap.Mutable(newMapKey)
						} else {
							newMapValue = reflectMap.NewValue()
						}
						sh.parseFieldOptions(newMapValue.Message(), row, depth+1, prefix+name)
						if !MessageValueEqual(emptyMapValue, newMapValue) {
							reflectMap.Set(newMapKey, newMapValue)
						}
					}
				} else {
					// value is scalar type
					key := "Key"     // deafult key name
					value := "Value" // deafult value name
					newMapKey := keyFd.Default().MapKey()
					// key cell
					cellValue, ok := row[prefix+name+key]
					if !ok {
						atom.Log.Panicf("key not found: %s", prefix+name+key)
					}
					fieldValue, err := sh.parseFieldValue(keyFd, cellValue)
					if err != nil {
						return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
					}
					newMapKey = fieldValue.MapKey()
					var newMapValue protoreflect.Value
					if reflectMap.Has(newMapKey) {
						newMapValue = reflectMap.Mutable(newMapKey)
					} else {
						newMapValue = reflectMap.NewValue()
					}
					// value cell
					cellValue, ok = row[prefix+name+value]
					if !ok {
						atom.Log.Panicf("value not found: %s", prefix+name+value)
					}
					newMapValue, err = sh.parseFieldValue(fd, cellValue)
					if err != nil {
						return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
					}
					if !reflectMap.Has(newMapKey) {
						reflectMap.Set(newMapKey, newMapValue)
					}
				}
			}
			if !msg.Has(fd) && reflectMap.Len() != 0 {
				msg.Set(fd, newValue)
			}
		} else if fd.IsList() {
			// Mutable returns a mutable reference to a composite type.
			if msg.Has(fd) {
				newValue = msg.Mutable(fd)
			}
			reflectList := newValue.List()
			if etype == tableaupb.Type_TYPE_INCELL_LIST {
				cellValue, ok := row[prefix+name]
				if !ok {
					atom.Log.Panicf("name not found: %s", prefix+name)
				}
				if cellValue != "" {
					// If s does not contain sep and sep is not empty, Split returns a
					// slice of length 1 whose only element is s.
					splits := strings.Split(cellValue, sep)
					for _, incellValue := range splits {
						value, err := sh.parseFieldValue(fd, incellValue)
						if err != nil {
							return errors.Wrapf(err, "failed to parse field value: %v", incellValue)
						}
						reflectList.Append(value)
					}
				}
			} else {
				emptyListValue := reflectList.NewElement()
				if layout == tableaupb.Layout_LAYOUT_VERTICAL {
					newListValue := reflectList.NewElement()
					if fd.Kind() == protoreflect.MessageKind {
						sh.parseFieldOptions(newListValue.Message(), row, depth+1, prefix+name)
						if !MessageValueEqual(emptyListValue, newListValue) {
							reflectList.Append(newListValue)
						}
					} else {
						// TODO: support list of scalar type when layout is vertical?
						// NOTE(wenchyzhu): we don't support list of scalar type when layout is vertical
					}
				} else {
					size := getPrefixSize(row, prefix+name)
					// atom.Log.Debug("prefix size: ", size)
					for i := 1; i <= size; i++ {
						newListValue := reflectList.NewElement()
						if fd.Kind() == protoreflect.MessageKind {
							sh.parseFieldOptions(newListValue.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
							if !MessageValueEqual(emptyListValue, newListValue) {
								reflectList.Append(newListValue)
							}
						} else {
							cellValue, ok := row[prefix+name+strconv.Itoa(i)]
							if !ok {
								atom.Log.Panicf("not found column name: %v", prefix+name)
							}
							newListValue, err = sh.parseFieldValue(fd, cellValue)
							if err != nil {
								return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
							}
							reflectList.Append(newListValue)
						}
					}
				}
			}
			if !msg.Has(fd) && reflectList.Len() != 0 {
				msg.Set(fd, newValue)
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				if etype == tableaupb.Type_TYPE_INCELL_STRUCT {
					cellValue, ok := row[prefix+name]
					if !ok {
						atom.Log.Panicf("not found column name: %v", prefix+name)
					}
					if cellValue != "" {
						// If s does not contain sep and sep is not empty, Split returns a
						// slice of length 1 whose only element is s.
						splits := strings.Split(cellValue, sep)
						subMd := newValue.Message().Descriptor()
						for i := 0; i < subMd.Fields().Len() && i < len(splits); i++ {
							fd := subMd.Fields().Get(i)
							// atom.Log.Debugf("fd.FullName().Name(): ", fd.FullName().Name())
							incellValue := splits[i]
							value, err := sh.parseFieldValue(fd, incellValue)
							if err != nil {
								return errors.Wrapf(err, "failed to parse field value: %v", incellValue)
							}
							newValue.Message().Set(fd, value)
						}
					}
				} else {
					subMsgName := string(fd.Message().FullName())
					_, found := specialMessageMap[subMsgName]
					if found {
						cellValue, ok := row[prefix+name]
						if !ok {
							atom.Log.Panicf("not found column name: %v", prefix+name)
						}
						newValue, err = sh.parseFieldValue(fd, cellValue)
						if err != nil {
							return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
						}
					} else {
						pkgName := newValue.Message().Descriptor().ParentFile().Package()
						if string(pkgName) != sh.ProtoPackage {
							atom.Log.Panicf("unknown message %v in package %v", subMsgName, pkgName)
						}
						sh.parseFieldOptions(newValue.Message(), row, depth+1, prefix+name)
					}
				}
				if !MessageValueEqual(emptyValue, newValue) {
					msg.Set(fd, newValue)
				}
			} else {
				cellValue, ok := row[prefix+name]
				if !ok {
					atom.Log.Panicf("not found column name: %v", prefix+name)
				}
				newValue, err = sh.parseFieldValue(fd, cellValue)
				if err != nil {
					return errors.Wrapf(err, "failed to parse field value: %v", cellValue)
				}
				msg.Set(fd, newValue)
			}
		}
	}
	return nil
}

func (sh *sheetParser) parseFieldValue(fd protoreflect.FieldDescriptor, value string) (protoreflect.Value, error) {
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
				t, err := parseTimeWithLocation(sh.LocationName, value)
				if err != nil {
					return protoreflect.ValueOf(ts.ProtoReflect()), errors.Wrapf(err, "illegal timestamp string format: %v", value)
				}
				// atom.Log.Debugf("timeStr: %v, unix timestamp: %v", value, t.Unix())
				ts = timestamppb.New(t)
				if err := ts.CheckValid(); err != nil {
					return protoreflect.ValueOf(ts.ProtoReflect()), errors.Wrapf(err, "invalid timestamp: %v", value)
				}
			}
			return protoreflect.ValueOf(ts.ProtoReflect()), nil
		case "google.protobuf.Duration":
			// make use of d as a *durationpb.Duration
			du := &durationpb.Duration{} // default
			if value != "" {
				d, err := time.ParseDuration(value)
				if err != nil {
					return protoreflect.ValueOf(du.ProtoReflect()), errors.Wrapf(err, "illegal duration string format: %v", value)
				}
				du = durationpb.New(d)
				if err := du.CheckValid(); err != nil {
					return protoreflect.ValueOf(du.ProtoReflect()), errors.Wrapf(err, "invalid duration: %v", value)
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
	atom.Log.Debugf("file:%s.proto, workbook:%s", protofile, workbook)
	return protofile, workbook
}

// parseMessageOptions is aimed to parse the options of a protobuf message.
func parseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32, bool) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgName := string(md.Name())
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
	worksheetName := worksheet.Name
	namerow := worksheet.Namerow
	if worksheet.Namerow == 0 {
		namerow = 1 // default
	}
	noterow := worksheet.Noterow
	if noterow == 0 {
		noterow = 1 // default
	}
	datarow := worksheet.Datarow
	if datarow == 0 {
		datarow = 2 // default
	}
	transpose := worksheet.Transpose
	atom.Log.Debugf("msgName:%s, worksheetName:%s, namerow:%d, noterow:%d, datarow:%d, transpose:%v\n", msgName, worksheetName, namerow, noterow, datarow, transpose)
	return msgName, worksheetName, namerow, noterow, datarow, transpose
}

func parseTimeWithLocation(locationName string, timeStr string) (time.Time, error) {
	// see https://golang.org/pkg/time/#LoadLocation
	if location, err := time.LoadLocation(locationName); err != nil {
		atom.Log.Panicf("LoadLocation failed: %s", err)
		return time.Time{}, err
	} else {
		timeLayout := "2006-01-02 15:04:05"
		t, err := time.ParseInLocation(timeLayout, timeStr, location)
		if err != nil {
			atom.Log.Panicf("ParseInLocation failed:%v ,timeStr: %v, locationName: %v", err, timeStr, locationName)
		}
		return t, nil
	}
}

func MessageValueEqual(v1, v2 protoreflect.Value) bool {
	if proto.Equal(v1.Message().Interface(), v2.Message().Interface()) {
		atom.Log.Debug("empty message exists")
		return true
	}
	return false
}

func getPrefixSize(row map[string]string, prefix string) int {
	// atom.Log.Debug("name prefix: ", prefix)
	size := 0
	for name := range row {
		if strings.HasPrefix(name, prefix) {
			num := 0
			// atom.Log.Debug("name: ", name)
			colSuffix := name[len(prefix):]
			// atom.Log.Debug("name: suffix ", colSuffix)
			for _, r := range colSuffix {
				if unicode.IsDigit(r) {
					num = num*10 + int(r-'0')
				} else {
					break
				}
			}
			size = int(math.Max(float64(size), float64(num)))
		}
	}
	return size
}
