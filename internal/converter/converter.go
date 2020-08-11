package converter

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Wenchy/tableau/tableaupb"
	"github.com/tealeg/xlsx/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const TableauPackageName = "test"

// WorkbookRootDir is root dir of workbooks.
const WorkbookRootDir = "./testdata/"

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
	// row 0: metarow
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
func TestParseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgFullName := string(md.FullName())
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	metarow := proto.GetExtension(opts, tableaupb.E_Metarow).(int32)
	descrow := proto.GetExtension(opts, tableaupb.E_Descrow).(int32)
	datarow := proto.GetExtension(opts, tableaupb.E_Datarow).(int32)
	fmt.Printf("message:%s, worksheet:%s, metarow:%d, descrow:%d, datarow:%d\n", msgFullName, worksheet, metarow, descrow, datarow)
	return msgFullName, worksheet, metarow, descrow, datarow
}

// TestParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func TestParseFieldOptions(msg protoreflect.Message, row map[string]string, depth int, prefix string) {
	md := msg.Descriptor()
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, tableaupb.E_Worksheet).(string)
	pkg := md.ParentFile().Package()
	fmt.Printf("%s// %s, '%s', %v, %v, %v\n", getTabStr(depth), md.FullName(), worksheet, md.IsMapEntry(), prefix, pkg)
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if pkg != TableauPackageName && pkg != "google.protobuf" {
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
		sep := proto.GetExtension(opts, tableaupb.E_Sep).(string)
		if sep == "" {
			sep = ","
		}
		subsep := proto.GetExtension(opts, tableaupb.E_Subsep).(string)
		if subsep == "" {
			subsep = ":"
		}
		fmt.Printf("%s%s(%v) %s(%s) %s = %d [(caption) = \"%s\", (type) = %s, (key) = \"%s\", (sep) = \"%s\"];\n",
			getTabStr(depth), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), prefix+caption, etype.String(), key, sep)
		// fmt.Println(fd.ContainingMessage().FullName())

		// if fd.Cardinality() == protoreflect.Repeated && fd.Kind() == protoreflect.MessageKind {
		// 	msg := fd.Message().New()
		// }
		if fd.IsMap() {
			keyFd := fd.MapKey()
			valueFd := fd.MapValue()
			reflectMap := msg.Mutable(fd).Map()
			// newKey := protoreflect.ValueOf(int32(1)).MapKey()
			// newKey := getScalarFieldValue(keyFd, "1111001").MapKey()
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
					newKey := getScalarFieldValue(keyFd, kv[0]).MapKey()
					newValue := reflectMap.NewValue()
					newValue = getScalarFieldValue(valueFd, kv[1])
					reflectMap.Set(newKey, newValue)
				}
			} else {
				newKey := keyFd.Default().MapKey()
				cellValue, ok := row[prefix+caption+key]
				if ok {
					newKey = getScalarFieldValue(keyFd, cellValue).MapKey()
				} else {
					panic(fmt.Sprintf("key not found: %s\n", prefix+caption+key))
				}
				var newValue protoreflect.Value
				if reflectMap.Has(newKey) {
					newValue = reflectMap.Mutable(newKey)
				} else {
					newValue = reflectMap.NewValue()
					reflectMap.Set(newKey, newValue)
				}
				// check if newValue is message type
				if valueFd.Kind() == protoreflect.MessageKind {
					newMsg := newValue.Message()
					TestParseFieldOptions(newMsg, row, depth+1, prefix+caption)
				} else {
					cellValue, ok := row[prefix+caption]
					if ok {
						newValue = getScalarFieldValue(fd, cellValue)
					}
				}
			}

		} else if fd.IsList() {
			reflectList := msg.Mutable(fd).List()
			if fd.Kind() == protoreflect.MessageKind {
				listSize := getListSize(row, prefix+caption)
				// fmt.Println("list size", listSize)
				for i := 1; i <= listSize; i++ {
					newElement := reflectList.NewElement()
					subMsg := newElement.Message()
					TestParseFieldOptions(subMsg, row, depth+1, prefix+caption+strconv.Itoa(i))
					reflectList.Append(newElement)
				}
			} else {
				if etype == tableaupb.FieldType_FIELD_TYPE_CELL_LIST {
					cellValue, ok := row[prefix+caption]
					if ok {
						splits := strings.Split(cellValue, sep)
						for _, v := range splits {
							value := getScalarFieldValue(fd, v)
							reflectList.Append(value)
						}
					} else {
						panic(fmt.Sprintf("key not found: %s\n", key))
					}

				} else {
					panic(fmt.Sprintf("unknown list type: %v\n", etype))
				}
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				subMsg := msg.Mutable(fd).Message()
				subMd := subMsg.Descriptor()
				// fmt.Println("subMsg FullName: ", subMd.FullName())
				subMsgName := subMd.FullName()
				switch subMsgName {
				case "google.protobuf.Timestamp":
					cellValue, ok := row[prefix+caption]
					if !ok {
						panic(fmt.Sprintf("not found column caption: %v\n", prefix+caption))
					}
					// layout := "2006-01-02T15:04:05.000Z"
					layout := "2006-01-02 15:04:05"
					t, err := time.Parse(layout, cellValue)
					if err != nil {
						panic(fmt.Sprintf("illegal timestamp string format: %v, err: %v\n", cellValue, err))
					}
					for i := 0; i < subMd.Fields().Len(); i++ {
						fd := subMd.Fields().Get(i)
						// fmt.Println("fd.FullName().Name(): ", fd.FullName().Name())
						if fd.FullName().Name() == "seconds" {
							value := getScalarFieldValue(fd, strconv.FormatInt(t.Unix(), 10))
							subMsg.Set(fd, value)
							break
						}
					}
				case "google.protobuf.Duration":
					cellValue, ok := row[prefix+caption]
					if !ok {
						panic(fmt.Sprintf("not found column: %v\n", prefix+caption))
					}
					for i := 0; i < subMd.Fields().Len(); i++ {
						fd := subMd.Fields().Get(i)
						// fmt.Println("fd.FullName().Name(): ", fd.FullName().Name())
						if fd.FullName().Name() == "seconds" {
							value := getScalarFieldValue(fd, cellValue)
							subMsg.Set(fd, value)
							break
						}
					}
				default:
					subPkg := subMd.ParentFile().Package()
					if subPkg != TableauPackageName {
						panic(fmt.Sprintf("unknown message %v in package %v", subMsgName, subPkg))
					}
					subMsg := msg.Mutable(fd).Message()
					TestParseFieldOptions(subMsg, row, depth+1, prefix+caption)
				}
			} else {
				cellValue, ok := row[prefix+caption]
				if ok {
					value := getScalarFieldValue(fd, cellValue)
					msg.Set(fd, value)
				} else {
					panic(fmt.Sprintf("not found column caption: %v\n", prefix+caption))
				}
			}
		}
	}
}

func getListSize(row map[string]string, prefix string) int {
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

func getScalarFieldValue(fd protoreflect.FieldDescriptor, cellVal string) protoreflect.Value {
	switch fd.Kind() {
	case protoreflect.Int32Kind:
		val, err := strconv.ParseInt(cellVal, 10, 32)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sint32Kind:
		val, err := strconv.ParseInt(cellVal, 10, 32)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Sfixed32Kind:
		val, err := strconv.ParseInt(cellVal, 10, 32)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int32(val))
	case protoreflect.Uint32Kind:
		val, err := strconv.ParseUint(cellVal, 10, 32)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Fixed32Kind:
		val, err := strconv.ParseUint(cellVal, 10, 32)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(uint32(val))
	case protoreflect.Int64Kind:
		val, err := strconv.ParseInt(cellVal, 10, 64)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sint64Kind:
		val, err := strconv.ParseInt(cellVal, 10, 64)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Sfixed64Kind:
		val, err := strconv.ParseInt(cellVal, 10, 64)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(int64(val))
	case protoreflect.Uint64Kind:
		val, err := strconv.ParseUint(cellVal, 10, 64)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(uint64(val))
	case protoreflect.Fixed64Kind:
		val, err := strconv.ParseUint(cellVal, 10, 64)
		if err != nil {
			panic(err)
		}
		return protoreflect.ValueOf(uint64(val))
	case protoreflect.StringKind:
		return protoreflect.ValueOf(string(cellVal))
	case protoreflect.BytesKind:
		return protoreflect.ValueOf([]byte(cellVal))
	default:
		panic(fmt.Sprintf("not supported scalar type: %s", fd.Kind().String()))
		// case protoreflect.BoolKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
		// case protoreflect.EnumKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
		// case protoreflect.DoubleKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
		// case protoreflect.MessageKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
		// case protoreflect.GroupKind:
		// 	panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		// 	return protoreflect.Value{}
	}
	// return protoreflect.Value{}
}