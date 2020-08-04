package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Wenchy/tableau/testpb"
	"github.com/tealeg/xlsx/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const tableauPackageName = "test"

func main() {
	// fmt.Println("Hello, world.")
	// item := testpb.Item{}
	// Redact(item.ProtoReflect().Interface())

	conf := testpb.RewardConf{
		Activities: map[int32]*testpb.RewardConf_Activity{
			1: &testpb.RewardConf_Activity{
				Chapters: map[int32]*testpb.RewardConf_Chapter{
					1: &testpb.RewardConf_Chapter{
						Sections: map[int32]*testpb.RewardConf_Row{
							1: &testpb.RewardConf_Row{
								ActivityId: 1,
								ChapterId:  2,
								SectionId:  3,
								Desc:       "aha",
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

	desc := conf.Activities[1].Chapters[1].Sections[1].Desc
	fmt.Printf("desc: %s\n", desc)

	md := conf.ProtoReflect().Descriptor()

	ParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	ParseMessageOptions(md)
	fmt.Println("==================")
	ParseFieldOptions(md, 0)
	fmt.Println("==================")

	readSheet("tests/Test.xlsx")
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

func GetTabStr(level int) string {
	tab := ""
	for i := 0; i < level; i++ {
		tab += "\t"
	}
	return tab
}

// ParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func ParseFieldOptions(md protoreflect.MessageDescriptor, level int) {
	fmt.Printf("%s// %s\n", GetTabStr(level), md.FullName())
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
		fmt.Printf("%s%s %s(%s) %s = %d [(col) = \"%s\", (type) = %s, (key) = \"%s\"];\n", GetTabStr(level), fd.Cardinality().String(), fd.Kind().String(), msgName, fd.FullName().Name(), fd.Number(), col, etype.String(), key)
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

func readSheet(workbook string) {
	// open an existing file
	wb, err := xlsx.OpenFile(workbook)
	if err != nil {
		panic(err)
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook
	fmt.Println("Sheets in this file:")
	for i, sh := range wb.Sheets {
		fmt.Println(i, sh.Name)
	}
	fmt.Println("----")
}
