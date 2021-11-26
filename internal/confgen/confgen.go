package confgen

import (
	"math"
	"os"
	"regexp"
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
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var newlineRegex *regexp.Regexp

func init() {
	newlineRegex = regexp.MustCompile(`\r?\n?`)
}

func clearNewline(s string) string {
	return newlineRegex.ReplaceAllString(s, "")
}

type metasheet struct {
	worksheet string // worksheet name
	namerow   int32  // exact row number of name at worksheet
	noterow   int32  // exact row number of description at wooksheet
	datarow   int32  // start row number of data
	transpose bool   // interchange the rows and columns
}

type Tableaux struct {
	ProtoPackage string // protobuf package name.
	InputDir     string // input dir of workbooks.
	OutputDir    string // output dir of generated files.

	LocationName string // Location represents the collection of time offsets in use in a geographical area. Default is "Asia/Shanghai".

	Output *options.OutputOption // output settings.

	metasheet metasheet // meta info of worksheet
}

var specialMessageMap = map[string]int{
	"google.protobuf.Timestamp": 1,
	"google.protobuf.Duration":  1,
}

func (tbx *Tableaux) Convert() error {
	// create output dir
	err := os.MkdirAll(tbx.OutputDir, 0700)
	if err != nil {
		return errors.Wrapf(err, "failed to create output dir: %s", tbx.OutputDir)
	}

	protoregistry.GlobalFiles.RangeFilesByPackage(
		protoreflect.FullName(tbx.ProtoPackage),
		func(fd protoreflect.FileDescriptor) bool {
			atom.Log.Debugf("filepath: %s", fd.Path())
			opts := fd.Options().(*descriptorpb.FileOptions)
			workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(*tableaupb.WorkbookOptions)
			if workbook == nil {
				return true
			}

			atom.Log.Debugf("proto: %s, workbook %s", fd.Path(), workbook)
			msgs := fd.Messages()
			for i := 0; i < msgs.Len(); i++ {
				md := msgs.Get(i)
				// atom.Log.Debugf("%s", md.FullName())
				opts := md.Options().(*descriptorpb.MessageOptions)
				worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(*tableaupb.WorksheetOptions)
				if worksheet == nil {
					continue
				}
				atom.Log.Infof("generate: %s, message: %s@%s, worksheet: %s@%s", md.Name(), fd.Path(), md.Name(), workbook.Name, worksheet.Name)
				newMsg := dynamicpb.NewMessage(md)
				tbx.export(newMsg)
			}
			return true
		})
	return nil
}

// export the protomsg message.
func (tbx *Tableaux) export(protomsg proto.Message) error {
	md := protomsg.ProtoReflect().Descriptor()
	msg := protomsg.ProtoReflect()
	_, workbook := parseFileOptions(md.ParentFile())
	msgName, worksheetName, namerow, noterow, datarow, transpose := parseMessageOptions(md)
	tbx.metasheet.worksheet = worksheetName
	tbx.metasheet.namerow = namerow
	tbx.metasheet.noterow = noterow
	tbx.metasheet.datarow = datarow
	tbx.metasheet.transpose = transpose

	wbPath := tbx.InputDir + workbook.Name
	book, err := excel.NewBook(wbPath)
	if err != nil {
		return errors.Wrapf(err, "failed to create new workbook: %s", wbPath)
	}

	sheet, ok := book.Sheets[worksheetName]
	if !ok {
		return errors.Wrapf(err, "not found worksheet: %s", worksheetName)
	}

	if transpose {
		// col caprow: name row
		// col [datarow, MaxRow]: data
		for col := 0; col < sheet.MaxCol; col++ {
			if col >= int(datarow)-1 {
				kv := make(map[string]string)
				for row := 0; row < sheet.MaxRow; row++ {
					name, err := sheet.Cell(row, int(namerow)-1)
					if err != nil {
						panic(err)
					}
					name = clearNewline(name)
					data, err := sheet.Cell(row, col)
					if err != nil {
						panic(err)
					}
					kv[name] = data
				}
				tbx.parseFieldOptions(msg, kv, 0, "")
			}
		}
	} else {
		// row namerow: name row
		// row [datarow, MaxRow]: data row
		for row := 0; row < sheet.MaxRow; row++ {
			if row >= int(datarow)-1 {
				kv := make(map[string]string)
				for col := 0; col < sheet.MaxCol; col++ {
					name, err := sheet.Cell(int(namerow)-1, col)
					if err != nil {
						panic(err)
					}
					name = clearNewline(name)
					data, err := sheet.Cell(row, col)
					if err != nil {
						panic(err)
					}
					kv[name] = data
				}
				tbx.parseFieldOptions(msg, kv, 0, "")
			}
		}
	}
	x := mexporter.New(msgName, protomsg, tbx.OutputDir, tbx.Output)
	return x.Export()
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

// parseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func (tbx *Tableaux) parseFieldOptions(msg protoreflect.Message, row map[string]string, depth int, prefix string) {
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
		if string(pkg) != tbx.ProtoPackage && pkg != "google.protobuf" {
			atom.Log.Debugf("%s// no need to proces package: %v", printer.Indent(depth), pkg)
			return
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
			// newKey := tbx.parseFieldValue(keyFd, "1111001").MapKey()
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
						key := tbx.parseFieldValue(keyFd, kv[0]).MapKey()
						val := reflectMap.NewValue()
						val = tbx.parseFieldValue(valueFd, kv[1])
						reflectMap.Set(key, val)
					}
				}
			} else {
				emptyMapValue := reflectMap.NewValue()
				if valueFd.Kind() == protoreflect.MessageKind {
					if layout == tableaupb.Layout_LAYOUT_HORIZONTAL {
						size := getPrefixSize(row, prefix+name)
						// atom.Log.Debug("prefix size: ", size)
						for i := 1; i <= size; i++ {
							newMapKey := keyFd.Default().MapKey()
							cellValue, ok := row[prefix+name+strconv.Itoa(i)+key]
							if !ok {
								atom.Log.Panicf("key not found: %s", prefix+name+key)
							}
							newMapKey = tbx.parseFieldValue(keyFd, cellValue).MapKey()
							var newMapValue protoreflect.Value
							if reflectMap.Has(newMapKey) {
								newMapValue = reflectMap.Mutable(newMapKey)
							} else {
								newMapValue = reflectMap.NewValue()
							}
							tbx.parseFieldOptions(newMapValue.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
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
						newMapKey = tbx.parseFieldValue(keyFd, cellValue).MapKey()
						var newMapValue protoreflect.Value
						if reflectMap.Has(newMapKey) {
							newMapValue = reflectMap.Mutable(newMapKey)
						} else {
							newMapValue = reflectMap.NewValue()
						}
						tbx.parseFieldOptions(newMapValue.Message(), row, depth+1, prefix+name)
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
					newMapKey = tbx.parseFieldValue(keyFd, cellValue).MapKey()
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
					newMapValue = tbx.parseFieldValue(fd, cellValue)
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
					for _, v := range splits {
						value := tbx.parseFieldValue(fd, v)
						reflectList.Append(value)
					}
				}
			} else {
				emptyListValue := reflectList.NewElement()
				if layout == tableaupb.Layout_LAYOUT_VERTICAL {
					newListValue := reflectList.NewElement()
					if fd.Kind() == protoreflect.MessageKind {
						tbx.parseFieldOptions(newListValue.Message(), row, depth+1, prefix+name)
						if !MessageValueEqual(emptyListValue, newListValue) {
							reflectList.Append(newListValue)
						}
					} else {
						// TODO: support list of scalar type when lyout is vertical?
						// NOTE(wenchyzhu): we don't support list of scalar type when layout is vertical
					}
				} else {
					size := getPrefixSize(row, prefix+name)
					// atom.Log.Debug("prefix size: ", size)
					for i := 1; i <= size; i++ {
						newListValue := reflectList.NewElement()
						if fd.Kind() == protoreflect.MessageKind {
							tbx.parseFieldOptions(newListValue.Message(), row, depth+1, prefix+name+strconv.Itoa(i))
							if !MessageValueEqual(emptyListValue, newListValue) {
								reflectList.Append(newListValue)
							}
						} else {
							cellValue, ok := row[prefix+name+strconv.Itoa(i)]
							if !ok {
								atom.Log.Panicf("not found column name: %v", prefix+name)
							}
							newListValue = tbx.parseFieldValue(fd, cellValue)
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
							value := tbx.parseFieldValue(fd, splits[i])
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
						newValue = tbx.parseFieldValue(fd, cellValue)
					} else {
						pkgName := newValue.Message().Descriptor().ParentFile().Package()
						if string(pkgName) != tbx.ProtoPackage {
							atom.Log.Panicf("unknown message %v in package %v", subMsgName, pkgName)
						}
						tbx.parseFieldOptions(newValue.Message(), row, depth+1, prefix+name)
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
				newValue = tbx.parseFieldValue(fd, cellValue)
				msg.Set(fd, newValue)
			}
		}
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

func (tbx *Tableaux) parseFieldValue(fd protoreflect.FieldDescriptor, cellValue string) protoreflect.Value {
	switch fd.Kind() {
	case protoreflect.Int32Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				atom.Log.Debug("cellValue:", cellValue)
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sint32Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sfixed32Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Uint32Kind:
		var val uint64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Fixed32Kind:
		var val uint64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Int64Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sint64Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sfixed64Kind:
		var val int64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Uint64Kind:
		var val uint64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint64(val))
	case protoreflect.Fixed64Kind:
		var val uint64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint64(val))
	case protoreflect.StringKind:
		return protoreflect.ValueOf(cellValue)
	case protoreflect.BytesKind:
		return protoreflect.ValueOf([]byte(cellValue))
	case protoreflect.BoolKind:
		var val bool // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseBool(cellValue)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(val)
	case protoreflect.FloatKind:
		var val float64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseFloat(cellValue, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(float32(val))
	case protoreflect.DoubleKind:
		var val float64 // default
		var err error
		if cellValue != "" {
			val, err = strconv.ParseFloat(cellValue, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(float64(val))
	case protoreflect.MessageKind:
		msgName := fd.Message().FullName()
		switch msgName {
		case "google.protobuf.Timestamp":
			ts := &timestamppb.Timestamp{}
			if cellValue != "" {
				// location name: "Asia/Shanghai" or "Asia/Chongqing".
				// NOTE(wenchy): There is no "Asia/Beijing" location name. Whoa!!! Big surprize?
				t, err := parseTimeWithLocation(tbx.LocationName, cellValue)
				if err != nil {
					atom.Log.Panicf("illegal timestamp string format: %v, err: %v", cellValue, err)
				}
				// atom.Log.Debugf("timeStr: %v, unix timestamp: %v", cellValue, t.Unix())
				ts = timestamppb.New(t)
				// make use of t as a *timestamppb.Timestamp
			}
			if err := ts.CheckValid(); err != nil {
				atom.Log.Panicf("invalid timestamp: %v", err)
			}
			return protoreflect.ValueOf(ts.ProtoReflect())
		case "google.protobuf.Duration":
			dur := &durationpb.Duration{} // default
			if cellValue != "" {
				d, err := time.ParseDuration(cellValue)
				if err != nil {
					atom.Log.Panicf("ParseDuration failed, illegal format: %v", cellValue)
				}
				dur = durationpb.New(d)
				// make use of d as a *durationpb.Duration
			}
			if err := dur.CheckValid(); err != nil {
				atom.Log.Panicf("duration CheckValid failed: %v", err)
			}
			return protoreflect.ValueOf(dur.ProtoReflect())
		default:
			atom.Log.Panicf("not supported message type: %s", msgName)
		}
	default:
		atom.Log.Panicf("not supported scalar type: %s", fd.Kind().String())
		// case protoreflect.EnumKind:
		// 	atom.Log.Panicf("not supported key type: %s", fd.Kind().String())
		// case protoreflect.GroupKind:
		// 	atom.Log.Panicf("not supported key type: %s", fd.Kind().String())
		// 	return protoreflect.Value{}
	}
	atom.Log.Panicf("should not go here")
	return protoreflect.Value{}
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
