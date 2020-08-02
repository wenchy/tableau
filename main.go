package main

import (
	"fmt"

	"github.com/Wenchy/tableau/testpb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func main() {
	// fmt.Println("Hello, world.")
	// item := testpb.Item{}
	// Redact(item.ProtoReflect().Interface())

	testConf := testpb.TestConf{}
	Redact(testConf.ProtoReflect().Descriptor())
}

// Redact a pb message
func Redact(md protoreflect.MessageDescriptor) {
	for i := 0; i < md.Fields().Len(); i++ {
		fd := md.Fields().Get(i)
		if fd.ParentFile().Package() != "test" {
			return
		}
		if fd.Kind() == protoreflect.MessageKind {
			fmt.Println(fd.Cardinality().String(), fd.Kind().String(), fd.FullName(), fd.Number())
			Redact(fd.Message())
		}

		// if fd.IsList() {
		// 	fmt.Println("repeated", fd.Kind().String(), fd.FullName().Name())
		// 	// Redact(fd.Options().ProtoReflect().Interface())
		// }
		opts := fd.Options().(*descriptorpb.FieldOptions)
		col := proto.GetExtension(opts, testpb.E_Col).(string)
		etype := proto.GetExtension(opts, testpb.E_Type).(testpb.TableauFieldType)
		fmt.Printf("%s, %s %s %s = %d [(col) = \"%s\", (type) = %s];\n", fd.ParentFile().Package(), fd.Cardinality().String(), fd.Kind().String(), fd.FullName().Name(), fd.Number(), col, etype.String())
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
