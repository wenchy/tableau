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
	// Redact clears every sensitive field in pb.
	// var item testpb.Item
	item := testpb.Item{}
	item.Id = 1
	item.Num = 2
	Redact(item.ProtoReflect().Interface())
}

// Redact a pb message
func Redact(pb proto.Message) {
	m := pb.ProtoReflect()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		opts := fd.Options().(*descriptorpb.FieldOptions)
		col := proto.GetExtension(opts, testpb.E_Col).(string)
		if col != "" {
			fmt.Println(fd.FullName().Name(), col)
			// fmt.Println(fd.ContainingMessage().FullName())
		}
		return true
	})
}
