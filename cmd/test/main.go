package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Wenchy/tableau/converter"
	"github.com/Wenchy/tableau/testpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func main() {
	// parseActivity()
	// parseItem()
	numFiles := protoregistry.GlobalFiles.NumFiles()
	fmt.Println("numFiles", numFiles)
	// protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
	// 	fmt.Printf("filepath: %s\n", fd.Path())
	// 	return true
	// })
	fmt.Println("====================")
	protoregistry.GlobalFiles.RangeFilesByPackage(protoreflect.FullName("test"), func(fd protoreflect.FileDescriptor) bool {
		// fmt.Printf("filepath: %s\n", fd.Path())
		opts := fd.Options().(*descriptorpb.FileOptions)
		workbook := proto.GetExtension(opts, testpb.E_Workbook).(string)
		if workbook == "" {
			return true
		}

		fmt.Printf("proto: %s => workbook: %s\n", fd.Path(), workbook)
		msgs := fd.Messages()
		for i := 0; i < msgs.Len(); i++ {
			md := msgs.Get(i)
			// fmt.Printf("%s\n", md.FullName())
			opts := md.Options().(*descriptorpb.MessageOptions)
			worksheet := proto.GetExtension(opts, testpb.E_Worksheet).(string)
			if worksheet != "" {
				fmt.Printf("message: %s, worksheet: %s\n", md.FullName(), worksheet)
			}
			newMsg := dynamicpb.NewMessage(md)
			export(newMsg)
		}
		return true
	})
}

func export(conf proto.Message) {
	md := conf.ProtoReflect().Descriptor()
	msg := conf.ProtoReflect()
	_, workbook := converter.TestParseFileOptions(md.ParentFile())
	fmt.Println("==================")
	_, worksheet, _, _, _ := converter.TestParseMessageOptions(md)
	fmt.Println("==================")
	sheet := converter.ReadSheet(converter.WorkbookRootDir+workbook, worksheet)
	// row 0: metarow
	// row 1 - MaxRow: datarow
	for nrow := 0; nrow < sheet.MaxRow; nrow++ {
		if nrow >= 1 {
			// row, err := sheet.Row(nrow)
			// if err != nil {
			// 	panic(err)
			// }
			kv := make(map[string]string)
			for i := 0; i < sheet.MaxCol; i++ {
				metaCell, err := sheet.Cell(0, i)
				if err != nil {
					panic(err)
				}
				key := metaCell.Value
				dataCell, err := sheet.Cell(nrow, i)
				if err != nil {
					panic(err)
				}
				value := dataCell.Value
				kv[key] = value
			}
			converter.TestParseFieldOptions(msg, kv, 0, "")
		}
		fmt.Println()
	}
	fmt.Println("==================")

	output, err := protojson.Marshal(conf.ProtoReflect().Interface())
	if err != nil {
		panic(err)
	}
	fmt.Println("json: ", string(output))
	var out bytes.Buffer
	json.Indent(&out, output, "", "    ")
	out.WriteTo(os.Stdout)
	fmt.Println()
}
