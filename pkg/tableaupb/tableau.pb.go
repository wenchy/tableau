// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: tableau.proto

package tableaupb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Cardinality of field
type Card int32

const (
	Card_CARD_REQUIRED Card = 0 // appears exactly one time
	Card_CARD_OPTIONAL Card = 1 // appears zero or one times
	Card_CARD_REPEATED Card = 2 // appears zero or more times
	Card_CARD_MAP      Card = 3 // appears zero or more times
)

// Enum value maps for Card.
var (
	Card_name = map[int32]string{
		0: "CARD_REQUIRED",
		1: "CARD_OPTIONAL",
		2: "CARD_REPEATED",
		3: "CARD_MAP",
	}
	Card_value = map[string]int32{
		"CARD_REQUIRED": 0,
		"CARD_OPTIONAL": 1,
		"CARD_REPEATED": 2,
		"CARD_MAP":      3,
	}
)

func (x Card) Enum() *Card {
	p := new(Card)
	*p = x
	return p
}

func (x Card) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Card) Descriptor() protoreflect.EnumDescriptor {
	return file_tableau_proto_enumTypes[0].Descriptor()
}

func (Card) Type() protoreflect.EnumType {
	return &file_tableau_proto_enumTypes[0]
}

func (x Card) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Card.Descriptor instead.
func (Card) EnumDescriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{0}
}

type Workbook struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options    *WorkbookOptions `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
	Worksheets []*Worksheet     `protobuf:"bytes,2,rep,name=worksheets,proto3" json:"worksheets,omitempty"`
}

func (x *Workbook) Reset() {
	*x = Workbook{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Workbook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Workbook) ProtoMessage() {}

func (x *Workbook) ProtoReflect() protoreflect.Message {
	mi := &file_tableau_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Workbook.ProtoReflect.Descriptor instead.
func (*Workbook) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{0}
}

func (x *Workbook) GetOptions() *WorkbookOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Workbook) GetWorksheets() []*Worksheet {
	if x != nil {
		return x.Worksheets
	}
	return nil
}

type Worksheet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options *WorksheetOptions `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
	Fields  []*Field          `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *Worksheet) Reset() {
	*x = Worksheet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Worksheet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Worksheet) ProtoMessage() {}

func (x *Worksheet) ProtoReflect() protoreflect.Message {
	mi := &file_tableau_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Worksheet.ProtoReflect.Descriptor instead.
func (*Worksheet) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{1}
}

func (x *Worksheet) GetOptions() *WorksheetOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Worksheet) GetFields() []*Field {
	if x != nil {
		return x.Fields
	}
	return nil
}

type Field struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Options *FieldOptions `protobuf:"bytes,1,opt,name=options,proto3" json:"options,omitempty"`
	Card    Card          `protobuf:"varint,2,opt,name=card,proto3,enum=tableau.Card" json:"card,omitempty"`
}

func (x *Field) Reset() {
	*x = Field{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Field) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Field) ProtoMessage() {}

func (x *Field) ProtoReflect() protoreflect.Message {
	mi := &file_tableau_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Field.ProtoReflect.Descriptor instead.
func (*Field) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{2}
}

func (x *Field) GetOptions() *FieldOptions {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Field) GetCard() Card {
	if x != nil {
		return x.Card
	}
	return Card_CARD_REQUIRED
}

var File_tableau_proto protoreflect.FileDescriptor

var file_tableau_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x1a, 0x15, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61,
	0x75, 0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x72, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x12, 0x32, 0x0a, 0x07, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12,
	0x32, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x65, 0x65, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x57, 0x6f,
	0x72, 0x6b, 0x73, 0x68, 0x65, 0x65, 0x74, 0x52, 0x0a, 0x77, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x65,
	0x65, 0x74, 0x73, 0x22, 0x68, 0x0a, 0x09, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x65, 0x65, 0x74,
	0x12, 0x33, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x57, 0x6f, 0x72, 0x6b,
	0x73, 0x68, 0x65, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x5b, 0x0a,
	0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x2f, 0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61,
	0x75, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x07,
	0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x0a, 0x04, 0x63, 0x61, 0x72, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e,
	0x43, 0x61, 0x72, 0x64, 0x52, 0x04, 0x63, 0x61, 0x72, 0x64, 0x2a, 0x4d, 0x0a, 0x04, 0x43, 0x61,
	0x72, 0x64, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x49,
	0x52, 0x45, 0x44, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x4f, 0x50,
	0x54, 0x49, 0x4f, 0x4e, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x43, 0x41, 0x52, 0x44,
	0x5f, 0x52, 0x45, 0x50, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x43,
	0x41, 0x52, 0x44, 0x5f, 0x4d, 0x41, 0x50, 0x10, 0x03, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x65, 0x6e, 0x63, 0x68, 0x79, 0x2f, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x61, 0x75, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tableau_proto_rawDescOnce sync.Once
	file_tableau_proto_rawDescData = file_tableau_proto_rawDesc
)

func file_tableau_proto_rawDescGZIP() []byte {
	file_tableau_proto_rawDescOnce.Do(func() {
		file_tableau_proto_rawDescData = protoimpl.X.CompressGZIP(file_tableau_proto_rawDescData)
	})
	return file_tableau_proto_rawDescData
}

var file_tableau_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tableau_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tableau_proto_goTypes = []interface{}{
	(Card)(0),                // 0: tableau.Card
	(*Workbook)(nil),         // 1: tableau.Workbook
	(*Worksheet)(nil),        // 2: tableau.Worksheet
	(*Field)(nil),            // 3: tableau.Field
	(*WorkbookOptions)(nil),  // 4: tableau.WorkbookOptions
	(*WorksheetOptions)(nil), // 5: tableau.WorksheetOptions
	(*FieldOptions)(nil),     // 6: tableau.FieldOptions
}
var file_tableau_proto_depIdxs = []int32{
	4, // 0: tableau.Workbook.options:type_name -> tableau.WorkbookOptions
	2, // 1: tableau.Workbook.worksheets:type_name -> tableau.Worksheet
	5, // 2: tableau.Worksheet.options:type_name -> tableau.WorksheetOptions
	3, // 3: tableau.Worksheet.fields:type_name -> tableau.Field
	6, // 4: tableau.Field.options:type_name -> tableau.FieldOptions
	0, // 5: tableau.Field.card:type_name -> tableau.Card
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_tableau_proto_init() }
func file_tableau_proto_init() {
	if File_tableau_proto != nil {
		return
	}
	file_tableau_options_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tableau_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Workbook); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tableau_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Worksheet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_tableau_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Field); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tableau_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tableau_proto_goTypes,
		DependencyIndexes: file_tableau_proto_depIdxs,
		EnumInfos:         file_tableau_proto_enumTypes,
		MessageInfos:      file_tableau_proto_msgTypes,
	}.Build()
	File_tableau_proto = out.File
	file_tableau_proto_rawDesc = nil
	file_tableau_proto_goTypes = nil
	file_tableau_proto_depIdxs = nil
}
