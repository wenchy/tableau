package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Wenchy/tableau/testpb"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tealeg/xlsx/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const tableauPackageName = "test"
const workbookRootDir = "./tests/"

func parseActivity() {
	// fmt.Println("Hello, world.")
	// item := testpb.Item{}
	// Redact(item.ProtoReflect().Interface())

	conf := testpb.ActivityConf{
		Activities: map[int32]*testpb.ActivityConf_Activity{
			1: &testpb.ActivityConf_Activity{
				Chapters: map[int32]*testpb.ActivityConf_Chapter{
					2: &testpb.ActivityConf_Chapter{
						Sections: map[int32]*testpb.ActivityConf_Row{
							3: &testpb.ActivityConf_Row{
								ActivityId: 1,
								ChapterId:  2,
								SectionId:  3,
								Desc:       "aha",
								Items: []*testpb.Item{
									&testpb.Item{
										Id:  1,
										Num: 2,
									},
									&testpb.Item{
										Id:  2,
										Num: 2,
									},
								},
								BeginTime: &timestamp.Timestamp{
									Seconds: 1596985188,
									Nanos:   1,
								},
								EndTime: &timestamp.Timestamp{
									Seconds: 1699985188,
									Nanos:   1,
								},
								Duration: &duration.Duration{
									Seconds: 3600,
									Nanos:   2,
								},
							},
						},
					},
				},
			},
		},
	}

	output, err := protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	var out bytes.Buffer
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()

	desc := conf.Activities[1].Chapters[2].Sections[3].Desc
	fmt.Printf("desc: %s\n", desc)

	md := conf.ProtoReflect().Descriptor()

	ParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	ParseMessageOptions(md)
	fmt.Println("==================")
	ParseFieldOptions(md, 0)
	fmt.Println("==================")
}

// ParseFileOptions is aimed to parse the options of a protobuf definition file.
func ParseFileOptions(fd protoreflect.FileDescriptor) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	workbook := proto.GetExtension(opts, testpb.E_Workbook).(string)
	fmt.Printf("file:%s.proto, workbook:%s\n", fd.FullName(), workbook)
}

// ParseMessageOptions is aimed to parse the options of a protobuf message.
func ParseMessageOptions(md protoreflect.MessageDescriptor) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, testpb.E_Worksheet).(string)
	metarow := proto.GetExtension(opts, testpb.E_Metarow).(int32)
	descrow := proto.GetExtension(opts, testpb.E_Descrow).(int32)
	datarow := proto.GetExtension(opts, testpb.E_Datarow).(int32)
	fmt.Printf("message:%s, worksheet:%s, metarow:%d, descrow:%d, datarow:%d\n", md.FullName(), worksheet, metarow, descrow, datarow)
}

func getTabStr(level int) string {
	tab := ""
	for i := 0; i < level; i++ {
		tab += "\t"
	}
	return tab
}

// ParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func ParseFieldOptions(md protoreflect.MessageDescriptor, level int) {
	fmt.Printf("%s// %s\n", getTabStr(level), md.FullName())
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ParentFile().Package() != tableauPackageName {
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
			// fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			// ParseFieldOptions(fd.Message(), level+1)
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		col := proto.GetExtension(opts, testpb.E_Col).(string)
		etype := proto.GetExtension(opts, testpb.E_Type).(testpb.FieldType)
		key := proto.GetExtension(opts, testpb.E_Key).(string)
		fmt.Printf("%s%s %s(%s) %s = %d [(col) = \"%s\", (type) = %s, (key) = \"%s\"];\n", getTabStr(level), fd.Cardinality().String(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), col, etype.String(), key)
		// fmt.Println(fd.ContainingMessage().FullName())

		if fd.Kind() == protoreflect.MessageKind {
			ParseFieldOptions(fd.Message(), level+1)
		}
	}

	// m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
	// 	opts := fd.Options().(*descriptorpb.FieldOptions)
	// 	col := proto.GetExtension(opts, testpb.E_Col).(string)
	// 	if col != "" {
	// 		fmt.Println(fd.FullName().Name(), col)
	// 		// fmt.Println(fd.ContainingMessage().FullName())
	// 	}
	// 	return true
	// })
}

func readSheet(workbook string, worksheet string) *xlsx.Sheet {
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
			// get the Cell in D1, which is row 0, col 3
			cell, err := sheet.Cell(nrow, ncol)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s ", cell.Value)
		}
		fmt.Println()
	}
}

func parseItem() {
	conf := testpb.ItemConf{
		Items: map[int32]*testpb.ItemConf_Row{
			1: &testpb.ItemConf_Row{
				Id:     1,
				Name:   "金币",
				Desc:   "通用货币",
				IconId: 5001,
				Attributes: []*testpb.ItemConf_Attribute{
					&testpb.ItemConf_Attribute{
						Id:    1,
						Value: 2,
					},
					&testpb.ItemConf_Attribute{
						Id:    2,
						Value: 2,
					},
				},
				Effects: []int32{1, 2, 3},
				ExpiryTime: &timestamp.Timestamp{
					Seconds: 1596985188,
					Nanos:   1,
				},
			},
		},
	}

	output, err := protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	var out bytes.Buffer
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()

	// desc := conf.Activities[1].Chapters[2].Sections[3].Desc
	// fmt.Printf("desc: %s\n", desc)

	md := conf.ProtoReflect().Descriptor()
	msg := conf.ProtoReflect()
	_, workbook := testParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	_, worksheet, _, _, _ := testParseMessageOptions(md)
	fmt.Println("==================")
	sheet := readSheet(workbookRootDir+workbook, worksheet)
	// row 0: metarow
	// row 1 - MaxRow: datarow
	for nrow := 0; nrow < sheet.MaxRow; nrow++ {
		if nrow >= 1 {
			row, err := sheet.Row(nrow)
			if err != nil {
				panic(err)
			}
			testParseFieldOptions(msg, row, 0)
		}
		fmt.Println()
	}
	fmt.Println("==================")

	output, err = protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()
}

// ParseFileOptions is aimed to parse the options of a protobuf definition file.
func testParseFileOptions(fd protoreflect.FileDescriptor) (string, string) {
	opts := fd.Options().(*descriptorpb.FileOptions)
	protofile := string(fd.FullName())
	workbook := proto.GetExtension(opts, testpb.E_Workbook).(string)
	fmt.Printf("file:%s.proto, workbook:%s\n", protofile, workbook)
	return protofile, workbook
}

// ParseMessageOptions is aimed to parse the options of a protobuf message.
func testParseMessageOptions(md protoreflect.MessageDescriptor) (string, string, int32, int32, int32) {
	opts := md.Options().(*descriptorpb.MessageOptions)
	msgFullName := string(md.FullName())
	worksheet := proto.GetExtension(opts, testpb.E_Worksheet).(string)
	metarow := proto.GetExtension(opts, testpb.E_Metarow).(int32)
	descrow := proto.GetExtension(opts, testpb.E_Descrow).(int32)
	datarow := proto.GetExtension(opts, testpb.E_Datarow).(int32)
	fmt.Printf("message:%s, worksheet:%s, metarow:%d, descrow:%d, datarow:%d\n", msgFullName, worksheet, metarow, descrow, datarow)
	return msgFullName, worksheet, metarow, descrow, datarow
}

// ParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func testParseFieldOptions(msg protoreflect.Message, row *xlsx.Row, level int) {
	md := msg.Descriptor()
	opts := md.Options().(*descriptorpb.MessageOptions)
	worksheet := proto.GetExtension(opts, testpb.E_Worksheet).(string)
	fmt.Printf("%s// %s, '%s', %v\n", getTabStr(level), md.FullName(), worksheet, md.IsMapEntry())
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ParentFile().Package() != tableauPackageName {
			return
		}
		msgName := ""
		if fd.Kind() == protoreflect.MessageKind {
			msgName = string(fd.Message().FullName())
			// fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			// ParseFieldOptions(fd.Message(), level+1)
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		col := proto.GetExtension(opts, testpb.E_Col).(string)
		etype := proto.GetExtension(opts, testpb.E_Type).(testpb.FieldType)
		key := proto.GetExtension(opts, testpb.E_Key).(string)
		fmt.Printf("%s%s(%v) %s(%s) %s = %d [(col) = \"%s\", (type) = %s, (key) = \"%s\"];\n", getTabStr(level), fd.Cardinality().String(), fd.IsMap(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), col, etype.String(), key)
		// fmt.Println(fd.ContainingMessage().FullName())

		// if fd.Cardinality() == protoreflect.Repeated && fd.Kind() == protoreflect.MessageKind {
		// 	msg := fd.Message().New()
		// }
		if fd.IsMap() {
			// TODO(wenchyzhu): add new empty item
			keyFd := fd.MapKey()
			valueFd := fd.MapValue()
			reflectMap := msg.Mutable(fd).Map()
			newValue := reflectMap.NewValue()
			// newKey := protoreflect.ValueOf(int32(1)).MapKey()
			// newKey := keyFd.Default().MapKey()
			newKey := getScalarFieldValue(keyFd, "1111001").MapKey()
			// check if newValue is message type
			if valueFd.Kind() == protoreflect.MessageKind {
				newMsg := newValue.Message()
				testParseFieldOptions(newMsg, row, level+1)
			} else {
				newValue = getScalarFieldValue(fd, "1111001")
			}
			reflectMap.Set(newKey, newValue)
		} else if fd.IsList() {
			// TODO(wenchyzhu): add new empty item
			reflectList := msg.Mutable(fd).List()
			if fd.Kind() == protoreflect.MessageKind {
				newElement := reflectList.NewElement()
				subMsg := newElement.Message()
				testParseFieldOptions(subMsg, row, level+1)
			} else {
				value := getScalarFieldValue(fd, "1111001")
				reflectList.Append(value)
			}
		} else {
			if fd.Kind() == protoreflect.MessageKind {
				subMsg := msg.Mutable(fd).Message()
				testParseFieldOptions(subMsg, row, level+1)
			} else {
				value := getScalarFieldValue(fd, "1111001")
				msg.Set(fd, value)
			}
		}
	}

	// m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
	// 	opts := fd.Options().(*descriptorpb.FieldOptions)
	// 	col := proto.GetExtension(opts, testpb.E_Col).(string)
	// 	if col != "" {
	// 		fmt.Println(fd.FullName().Name(), col)
	// 		// fmt.Println(fd.ContainingMessage().FullName())
	// 	}
	// 	return true
	// })
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

	case protoreflect.BoolKind:
		panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		return protoreflect.Value{}
	case protoreflect.EnumKind:
		panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		return protoreflect.Value{}
	case protoreflect.DoubleKind:
		panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		return protoreflect.Value{}
	case protoreflect.MessageKind:
		panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		return protoreflect.Value{}
	case protoreflect.GroupKind:
		panic(fmt.Sprintf("not supported key type: %s", fd.Kind().String()))
		return protoreflect.Value{}
	}
	return protoreflect.Value{}
}
