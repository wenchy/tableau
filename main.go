package main

import (
	"fmt"

	"github.com/Wenchy/tableau/testpb"
	"github.com/tealeg/xlsx/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const tableauPackageName = "test"

func main() {
	// fmt.Println("Hello, world.")
	// item := testpb.Item{}
	// Redact(item.ProtoReflect().Interface())

	itemConf := testpb.ItemConf{}
	md := itemConf.ProtoReflect().Descriptor()

	ParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	ParseMessageOptions(md)
	fmt.Println("==================")
	ParseFieldOptions(md)
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

// ParseFieldOptions is aimed to parse the options of all the fields of a protobuf message.
func ParseFieldOptions(md protoreflect.MessageDescriptor) {
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ParentFile().Package() != tableauPackageName {
			return
		}
		if fd.Kind() == protoreflect.MessageKind {
			fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			ParseFieldOptions(fd.Message())
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		col := proto.GetExtension(opts, testpb.E_Col).(string)
		etype := proto.GetExtension(opts, testpb.E_Type).(testpb.FieldType)
		fmt.Printf("%s %s %s = %d [(col) = \"%s\", (type) = %s];\n", fd.Cardinality().String(), fd.Kind().String(), fd.FullName().Name(), fd.Number(), col, etype.String())
		// fmt.Println(fd.ContainingMessage().FullName())

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
