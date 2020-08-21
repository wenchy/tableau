package converter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Wenchy/tableau/pkg/tableaupb"
	"github.com/iancoleman/strcase"
	"github.com/tealeg/xlsx/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Format int

// file format
const (
	JSON      Format = 0
	Protobin         = 1
	Prototext        = 2
	// Xlsx             = 3
)

type metasheet struct {
	worksheet string // worksheet name
	captrow   int32  // exact row number of caption at worksheet
	descrow   int32  // exact row number of description at wooksheet
	datarow   int32  // start row number of data
	transpose bool   // interchange the rows and columns
}

type Tableaux struct {
	ProtoPackageName          string    // protobuf package name.
	InputPath                 string    // root dir of workbooks.
	OutputPath                string    // output path of generated files.
	OutputFilenameAsSnakeCase bool      // output filename as snake case, default is camel case same as the protobuf message name.
	OutputFormat              Format    // output format: json, protobin, or prototext, and default is json.
	OutputPretty              bool      // output pretty format, with mulitline and indent.
	metasheet                 metasheet // meta info of worksheet
}

var specialMessageMap = map[string]int{
	"google.protobuf.Timestamp": 1,
	"google.protobuf.Duration":  1,
}

func (tbx *Tableaux) Convert() {
	// parseActivity()
	// parseItem()
	// numFiles := protoregistry.GlobalFiles.NumFiles()
	// fmt.Println("numFiles", numFiles)
	// protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
	// 	fmt.Printf("filepath: %s\n", fd.Path())
	// 	return true
	// })
	// fmt.Println("====================")

	// create oupput dir
	err := os.MkdirAll(tbx.OutputPath, 0700)
	if err != nil {
		panic(err)

	}

	protoPackage := protoreflect.FullName(tbx.ProtoPackageName)
	protoregistry.GlobalFiles.RangeFilesByPackage(protoPackage, func(fd protoreflect.FileDescriptor) bool {
		fmt.Printf("filepath: %s\n", fd.Path())
		opts := fd.Options().(*descriptorpb.FileOptions)
		workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(string)
		if workbook == "" {
			return true
		}

		fmt.Printf("proto: %s => workbook: %s\n", fd.Path(), workbook)
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			// fmt.Printf("%s\n", md.FullName())
			opts := md.Options().(*descriptorpb.MessageOptions)
			worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
			if worksheet != "" {
				fmt.Printf("message: %s, worksheet: %s\n", md.FullName(), worksheet)
			}
			newMsg := dynamicpb.NewMessage(md)
			tbx.Export(newMsg)
		}
		return true
	})
	err = IllegalFieldType{FieldType: "Map", Line: 10}
	fmt.Println(err)
}

// Export the protomsg message.
func (tbx *Tableaux) Export(protomsg proto.Message) {
	md := protomsg.ProtoReflect().Descriptor()
	msg := protomsg.ProtoReflect()
	_, workbook := TestParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	msgName, worksheet, captrow, descrow, datarow, transpose := TestParseMessageOptions(md)
	tbx.metasheet.worksheet = worksheet
	tbx.metasheet.captrow = captrow
	tbx.metasheet.descrow = descrow
	tbx.metasheet.datarow = datarow
	tbx.metasheet.transpose = transpose

	fmt.Println("==================")
	sheet := ReadSheet(tbx.InputPath+workbook, worksheet)
	if transpose {
		// col caprow: caption row
		// col [datarow, MaxRow]: data
		for ncol := 0; ncol < sheet.MaxCol; ncol++ {
			if ncol >= int(datarow)-1 {
				// row, err := sheet.Row(nrow)
				// if err != nil {
				// 	panic(err)
				// }
				kv := make(map[string]string)
				for i := 0; i < sheet.MaxRow; i++ {
					captionCell, err := sheet.Cell(i, int(captrow)-1)
					if err != nil {
						panic(err)
					}
					key := captionCell.Value
					dataCell, err := sheet.Cell(i, ncol)
					if err != nil {
						panic(err)
					}
					value := dataCell.Value
					kv[key] = value
				}
				tbx.TestParseFieldOptions(msg, kv, 0, "")
			}
			fmt.Println()
		}

	} else {
		// row captrow: caption row
		// row [datarow, MaxRow]: data row
		for nrow := 0; nrow < sheet.MaxRow; nrow++ {
			if nrow >= int(datarow)-1 {
				// row, err := sheet.Row(nrow)
				// if err != nil {
				// 	panic(err)
				// }
				kv := make(map[string]string)
				for i := 0; i < sheet.MaxCol; i++ {
					captionCell, err := sheet.Cell(int(captrow)-1, i)
					if err != nil {
						panic(err)
					}
					key := captionCell.Value
					dataCell, err := sheet.Cell(nrow, i)
					if err != nil {
						panic(err)
					}
					value := dataCell.Value
					kv[key] = value
				}
				tbx.TestParseFieldOptions(msg, kv, 0, "")
			}
			fmt.Println()
		}
	}
	fmt.Println("==================")
	filename := msgName
	if tbx.OutputFilenameAsSnakeCase {
		filename = strcase.ToSnake(msgName)
	}
	filePath := tbx.OutputPath + filename
	switch tbx.OutputFormat {
	case JSON:
		exportJSON(protomsg, filePath, tbx.OutputPretty)
	case Protobin:
		exportProtobin(protomsg, filePath)
	case Prototext:
		exportPrototext(protomsg, filePath, tbx.OutputPretty)
	default:
		fmt.Println("unknown format, default to JSON")
		exportJSON(protomsg, filePath, tbx.OutputPretty)
	}
}

func exportJSON(protomsg proto.Message, filePath string, pretty bool) {
	var out []byte
	var err error
	if pretty {
		o, err := protojson.Marshal(protomsg)
		if err != nil {
			panic(err)
		}
		// fmt.Println("json: ", string(output))
		var buf bytes.Buffer
		json.Indent(&buf, o, "", "    ")
		out = buf.Bytes()
	} else {
		out, err = protojson.Marshal(protomsg)
		if err != nil {
			panic(err)
		}
	}
	err = ioutil.WriteFile(filePath+".json", out, 0644)
	if err != nil {
		panic(err)
	}
	// out.WriteTo(os.Stdout)
	fmt.Println()
}

func exportProtobin(protomsg proto.Message, filePath string) {
	out, err := proto.Marshal(protomsg)
	if err != nil {
		log.Fatalln("Failed to encode protomsg:", err)
	}
	if err := ioutil.WriteFile(filePath+".protobin", out, 0644); err != nil {
		log.Fatalln("Failed to write file:", err)
	}
	// out.WriteTo(os.Stdout)
	fmt.Println()
}

func exportPrototext(protomsg proto.Message, filePath string, pretty bool) {
	var out []byte
	var err error
	if pretty {
		opts := prototext.MarshalOptions{
			Multiline: true,
			Indent:    "    ",
		}
		out, err = opts.Marshal(protomsg)
		if err != nil {
			log.Fatalln("Failed to encode protomsg:", err)
		}
	} else {
		out, err = prototext.Marshal(protomsg)
		if err != nil {
			panic(err)
		}
	}
	if err := ioutil.WriteFile(filePath+".prototext", out, 0644); err != nil {
		log.Fatalln("Failed to write file:", err)
	}
	// out.WriteTo(os.Stdout)
	fmt.Println()
}

func getTabStr(depth int) string {
	tab := ""
	for i := 0; i < depth; i++ {
		tab += "\t"
	}
	return tab
}

// ReadSheet read a sheet from specified workbook.
func ReadSheet(workbook string, worksheet string) *xlsx.Sheet {
	// open an existing file
	wb, err := xlsx.OpenFile(workbook)
	if err != nil {
		panic(err)
	}
	sh, ok := wb.Sheet[worksheet]
	if !ok {
		fmt.Printf("Sheet %s does not exist in %s\n", worksheet, workbook)
		panic("sheet not foound")
	}
	exportSheet(sh)
	fmt.Println("----")
	return sh
}

func exportSheet(sheet *xlsx.Sheet) {
	fmt.Printf("MaxCol: %d, MaxRow: %d\n", sheet.MaxCol, sheet.MaxRow)
	// row 0: captrow
	// row 1 - MaxRow: datarow
	for nrow := 0; nrow < sheet.MaxRow; nrow++ {
		for ncol := 0; ncol < sheet.MaxCol; ncol++ {
			// get the Cell in D1, which is row 0, column 3
			cell, err := sheet.Cell(nrow, ncol)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s ", cell.Value)
		}
		fmt.Println()
	}
}

// TestParseFileOptions is aimed to parse the options of a protobuf definition file.
func TestParseFileOptions(fd protoreflect.FileDescriptor) (string, string) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	protofile := string(fd.FullName())
	workbook := proto.GetExtension(opts, tableaupb.E_Workbook).(string)
	fmt.Printf("file:%s.proto, workbook:%s\n", protofile, workbook)
	return protofile, workbook
}

// TestParseMessageOptions is aimed to parse the options of a protobuf message.
func TestParseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32, bool) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgName := string(md.Name())
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	captrow := proto.GetExtension(opts, tableaupb.E_Captrow).(int32)
	if captrow == 0 {
		captrow = 1 // default
	}
	descrow := proto.GetExtension(opts, tableaupb.E_Descrow).(int32)
	if descrow == 0 {
		descrow = 1 // default
	}
	datarow := proto.GetExtension(opts, tableaupb.E_Datarow).(int32)
	if datarow == 0 {
		datarow = 2 // default
	}
	transpose := proto.GetExtension(opts, tableaupb.E_Transpose).(bool)
	fmt.Printf("message:%s, worksheet:%s, captrow:%d, descrow:%d, datarow:%d, transpose:%v\n", msgName, worksheet, captrow, descrow, datarow, transpose)
	return msgName, worksheet, captrow, descrow, datarow, transpose
}

// TestParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func (tbx *Tableaux) TestParseFieldOptions(msg protoreflect.Message, row map[string]string, depth int, prefix string) {
	md := msg.Descriptor()
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	pkg := md.ParentFile().Package()
	fmt.Printf("%s// %s, '%s', %v, %v, %v\n", getTabStr(depth), md.FullName(), worksheet, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if string(pkg) != tbx.ProtoPackageName && pkg != "google.protobuf" {
			fmt.Printf("%s// no need to proces package: %v\n", getTabStr(depth), pkg)
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
			// fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			// ParseFieldOptions(fd.Message(), depth+1)
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		caption := proto.GetExtension(opts, tableaupb.E_Caption).(string)
		etype := proto.GetExtension(opts, tableaupb.E_Type).(tableaupb.FieldType)
		key := proto.GetExtension(opts, tableaupb.E_Key).(string)
		layout := proto.GetExtension(opts, tableaupb.E_Layout).(tableaupb.CompositeLayout)
		sep := proto.GetExtension(opts, tableaupb.E_Sep).(string)
		if sep == "" {
			sep = ","
		}
		subsep := proto.GetExtension(opts, tableaupb.E_Subsep).(string)
		if subsep == "" {
			subsep = ":"
		}
		fmt.Printf("%s%s(%v) %s(%s) %s = %d [(caption) = \"%s\", (type) = %s, (key) = \"%s\", (layout) = \"%s\", (sep) = \"%s\"];\n",
			getTabStr(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+caption, etype.String(), layout.String(), key, sep)
		// fmt.Println(fd.ContainingMessage().FullName())

		// if fd.Cardinality() == protoreflect.Repeated && fd.Kind() == protoreflect.MessageKind {
		// 	msg := fd.Message().New()
		// }
		if fd.IsMap() {
			keyFd := fd.MapKey()
			valueFd := fd.MapValue()
			reflectMap := msg.Mutable(fd).Map()
			// newKey := protoreflect.ValueOf(int32(1)).MapKey()
			// newKey := getFieldValue(keyFd, "1111001").MapKey()
			if etype == tableaupb.FieldType_FIELD_TYPE_CELL_MAP {
				if valueFd.Kind() == protoreflect.MessageKind {
					panic("in-cell map do not support value as message type")
				}
				cellValue, ok := row[prefix+caption]
				if !ok {
					panic(fmt.Sprintf("column caption not found: %s\n", prefix+caption))
				}
				splits := strings.Split(cellValue, sep)
				for _, pair := range splits {
					kv := strings.Split(pair, subsep)
					if len(kv) != 2 {
						panic(fmt.Sprintf("illegal key-value pair: %v, %v\n", prefix+caption, pair))
					}
					newKey := getFieldValue(keyFd, kv[0]).MapKey()
					newValue := reflectMap.NewValue()
					newValue = getFieldValue(valueFd, kv[1])
					reflectMap.Set(newKey, newValue)
				}
			} else {
				// check if newValue is message type
				if valueFd.Kind() == protoreflect.MessageKind {
					if layout == tableaupb.CompositeLayout_COMPOSITE_LAYOUT_HORIZONTAL {
						size := getPrefixSize(row, prefix+caption)
						// fmt.Println("prefix size: ", size)
						for i := 1; i <= size; i++ {
							newKey := keyFd.Default().MapKey()
							cellValue, ok := row[prefix+caption+strconv.Itoa(i)+key]
							if !ok {
								panic(fmt.Sprintf("key not found: %s\n", prefix+caption+key))
							}
							newKey = getFieldValue(keyFd, cellValue).MapKey()
							var newValue protoreflect.Value
							if reflectMap.Has(newKey) {
								newValue = reflectMap.Mutable(newKey)
							} else {
								newValue = reflectMap.NewValue()
								reflectMap.Set(newKey, newValue)
							}
							newMsg := newValue.Message()
							tbx.TestParseFieldOptions(newMsg, row, depth+1, prefix+caption+strconv.Itoa(i))
						}
					} else {
						newKey := keyFd.Default().MapKey()
						cellValue, ok := row[prefix+caption+key]
						if !ok {
							panic(fmt.Sprintf("key not found: %s\n", prefix+caption+key))
						}
						newKey = getFieldValue(keyFd, cellValue).MapKey()
						var newValue protoreflect.Value
						if reflectMap.Has(newKey) {
							newValue = reflectMap.Mutable(newKey)
						} else {
							newValue = reflectMap.NewValue()
							reflectMap.Set(newKey, newValue)
						}
						newMsg := newValue.Message()
						tbx.TestParseFieldOptions(newMsg, row, depth+1, prefix+caption)
					}
				} else {
					newKey := keyFd.Default().MapKey()
					cellValue, ok := row[prefix+caption+key]
					if !ok {
						panic(fmt.Sprintf("key not found: %s\n", prefix+caption+key))
					}
					newKey = getFieldValue(keyFd, cellValue).MapKey()
					var newValue protoreflect.Value
					if reflectMap.Has(newKey) {
						newValue = reflectMap.Mutable(newKey)
					} else {
						newValue = reflectMap.NewValue()
						reflectMap.Set(newKey, newValue)
					}
					newValue = getFieldValue(fd, cellValue)
				}
			}

		} else if fd.IsList() {
			reflectList := msg.Mutable(fd).List()
			if fd.Kind() == protoreflect.MessageKind {
				if layout == tableaupb.CompositeLayout_COMPOSITE_LAYOUT_VERTICAL {
					newElement := reflectList.NewElement()
					subMsg := newElement.Message()
					tbx.TestParseFieldOptions(subMsg, row, depth+1, prefix+caption)
					reflectList.Append(newElement)
				} else {
					size := getPrefixSize(row, prefix+caption)
					// fmt.Println("prefix size: ", size)
					for i := 1; i <= size; i++ {
						newElement := reflectList.NewElement()
						subMsg := newElement.Message()
						tbx.TestParseFieldOptions(subMsg, row, depth+1, prefix+caption+strconv.Itoa(i))
						reflectList.Append(newElement)
					}
				}
			} else {
				if etype == tableaupb.FieldType_FIELD_TYPE_CELL_LIST {
					cellValue, ok := row[prefix+caption]
					if ok {
						splits := strings.Split(cellValue, sep)
						for _, v := range splits {
							value := getFieldValue(fd, v)
							reflectList.Append(value)
						}
					} else {
						panic(fmt.Sprintf("caption not found: %s\n", prefix+caption))
					}

				} else {
					panic(fmt.Sprintf("unknown list type: %v\n", etype))
				}
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				subMsg := msg.Mutable(fd).Message()
				subMd := subMsg.Descriptor()

				if etype == tableaupb.FieldType_FIELD_TYPE_CELL_MESSAGE {
					cellValue, ok := row[prefix+caption]
					if !ok {
						panic(fmt.Sprintf("not found column caption: %v\n", prefix+caption))
					}
					splits := strings.Split(cellValue, sep)
					if len(splits) != subMd.Fields().Len() {
						// TODO(wenchyzhu): more clear error message
						panic("in-cell message fields len not equal to cell splits len")
					}
					for i := 0; i < subMd.Fields().Len(); i++ {
						fd := subMd.Fields().Get(i)
						// fmt.Println("fd.FullName().Name(): ", fd.FullName().Name())
						value := getFieldValue(fd, splits[i])
						subMsg.Set(fd, value)
					}
				} else {
					subMsgName := string(fd.Message().FullName())
					_, found := specialMessageMap[subMsgName]
					if found {
						cellValue, ok := row[prefix+caption]
						if !ok {
							panic(fmt.Sprintf("not found column caption: %v\n", prefix+caption))
						}
						value := getFieldValue(fd, cellValue)
						msg.Set(fd, value)
					} else {
						subPkg := subMd.ParentFile().Package()
						if string(subPkg) != tbx.ProtoPackageName {
							panic(fmt.Sprintf("unknown message %v in package %v", subMsgName, subPkg))
						}
						tbx.TestParseFieldOptions(subMsg, row, depth+1, prefix+caption)
					}
				}
			} else {
				cellValue, ok := row[prefix+caption]
				if ok {
					value := getFieldValue(fd, cellValue)
					msg.Set(fd, value)
				} else {
					panic(fmt.Sprintf("not found column caption: %v\n", prefix+caption))
				}
			}
		}
	}
}

func getPrefixSize(row map[string]string, prefix string) int {
	// fmt.Println("caption prefix: ", prefix)
	size := 0
	for caption := range row {
		if strings.HasPrefix(caption, prefix) {
			num := 0
			// fmt.Println("caption: ", caption)
			colSuffix := caption[len(prefix):]
			// fmt.Println("caption: suffix ", colSuffix)
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

func getFieldValue(fd protoreflect.FieldDescriptor, cellValue string) protoreflect.Value {
	switch fd.Kind() {
	case protoreflect.Int32Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				fmt.Println("cellValue: ", cellValue)
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sint32Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sfixed32Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Uint32Kind:
		var val uint64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Fixed32Kind:
		var val uint64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Int64Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sint64Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sfixed64Kind:
		var val int64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseInt(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Uint64Kind:
		var val uint64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseUint(cellValue, 10, 64)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(uint64(val))
	case protoreflect.Fixed64Kind:
		var val uint64
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
		var val bool
		var err error
		if cellValue != "" {
			val, err = strconv.ParseBool(cellValue)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(val)
	case protoreflect.FloatKind:
		var val float64
		var err error
		if cellValue != "" {
			val, err = strconv.ParseFloat(cellValue, 32)
			if err != nil {
				panic(err)
			}
		}
		return protoreflect.ValueOf(float32(val))
	case protoreflect.DoubleKind:
		var val float64
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
			// location name: "Asia/Shanghai" or "Asia/Chongqing".
			// NOTE(wenchy): There is no "Asia/Beijing" location name. Whoa!!! Big surprize?
			t, err := parseTimeWithLocation("Asia/Shanghai", cellValue)
			if err != nil {
				panic(fmt.Sprintf("illegal timestamp string format: %v, err: %v\n", cellValue, err))
			}
			// fmt.Printf("timeStr: %v, unix timestamp: %v\n", cellValue, t.Unix())
			ts := timestamppb.New(t)
			// make use of t as a *timestamppb.Timestamp
			if err = ts.CheckValid(); err != nil {
				panic(fmt.Sprintf("invalid timestamp: %v\n", err))
			}
			return protoreflect.ValueOf(ts.ProtoReflect())
		case "google.protobuf.Duration":
			d, err := time.ParseDuration(cellValue)
			if err != nil {
				panic(fmt.Sprintf("ParseDuration failed, illegal format: %v\n", cellValue))
			}
			dur := durationpb.New(d)
			// make use of d as a *durationpb.Duration
			if err = dur.CheckValid(); err != nil {
				panic(fmt.Sprintf("duration CheckValid failed: %v\n", err))
			}
			return protoreflect.ValueOf(dur.ProtoReflect())
		default:
			panic(fmt.Sprintf("not supported message type: %s", msgName))
		}
	default:
		panic(fmt.Sprintf("not supported scalar type: %s", fd.Kind().String()))
		// case protoreflect.EnumKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// case protoreflect.GroupKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
	}
}

func parseTimeWithLocation(locationName string, timeStr string) (time.Time, error) {
	// see https://golang.org/pkg/time/#LoadLocation
	if location, err := time.LoadLocation(locationName); err != nil {
		panic(fmt.Sprintf("LoadLocation failed: %s", err))
	} else {
		timeLayout := "2006-01-02 15:04:05"
		t, err := time.ParseInLocation(timeLayout, timeStr, location)
		if err != nil {
			panic(fmt.Sprintf("ParseInLocation failed:%v ,timeStr: %v, locationName: %v\n", err, timeStr, locationName))
		}
		return t, nil
	}
}
