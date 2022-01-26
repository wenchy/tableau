// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.4
// source: tableau.proto

package tableaupb

import (
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

// field type.
type Type int32

const (
	// auto deduced protobuf types
	Type_TYPE_DEFAULT Type = 0
	//////////////////////////////
	/////Build-in Types///////////
	//////////////////////////////
	// interger
	Type_TYPE_INT32  Type = 1
	Type_TYPE_UINT32 Type = 2
	Type_TYPE_INT64  Type = 3
	Type_TYPE_UINT64 Type = 4
	// floating-point number
	Type_TYPE_DOUBLE Type = 5
	Type_TYPE_FLOAT  Type = 6
	// bool
	Type_TYPE_BOOL Type = 7
	// string
	Type_TYPE_STRING Type = 8
	////////////////////////
	/////Struct Type//////
	////////////////////////
	Type_TYPE_STRUCT Type = 10
	////////////////////////
	/////Extended Types/////
	////////////////////////
	// time
	Type_TYPE_DATE     Type = 21 // format: "yyyy-MM-dd"
	Type_TYPE_TIME     Type = 22 // format: "HH:mm:ss"
	Type_TYPE_DATETIME Type = 23 // format: "yyyy-MM-dd HH:mm:ss"
	// list in a cell:
	// - the list **item** must be **built-in** type
	// - format: ',' separated items
	Type_TYPE_INCELL_LIST Type = 24
	// map in a cell:
	// - both the **key** and **value** must be **built-in** type
	// - format: key-value pairs is separated by ',', and
	//           key and value is separated by ':'
	Type_TYPE_INCELL_MAP Type = 25
	// struct in a cell
	Type_TYPE_INCELL_STRUCT Type = 26
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0:  "TYPE_DEFAULT",
		1:  "TYPE_INT32",
		2:  "TYPE_UINT32",
		3:  "TYPE_INT64",
		4:  "TYPE_UINT64",
		5:  "TYPE_DOUBLE",
		6:  "TYPE_FLOAT",
		7:  "TYPE_BOOL",
		8:  "TYPE_STRING",
		10: "TYPE_STRUCT",
		21: "TYPE_DATE",
		22: "TYPE_TIME",
		23: "TYPE_DATETIME",
		24: "TYPE_INCELL_LIST",
		25: "TYPE_INCELL_MAP",
		26: "TYPE_INCELL_STRUCT",
	}
	Type_value = map[string]int32{
		"TYPE_DEFAULT":       0,
		"TYPE_INT32":         1,
		"TYPE_UINT32":        2,
		"TYPE_INT64":         3,
		"TYPE_UINT64":        4,
		"TYPE_DOUBLE":        5,
		"TYPE_FLOAT":         6,
		"TYPE_BOOL":          7,
		"TYPE_STRING":        8,
		"TYPE_STRUCT":        10,
		"TYPE_DATE":          21,
		"TYPE_TIME":          22,
		"TYPE_DATETIME":      23,
		"TYPE_INCELL_LIST":   24,
		"TYPE_INCELL_MAP":    25,
		"TYPE_INCELL_STRUCT": 26,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_tableau_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_tableau_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{0}
}

// layout of composite types, such as list and map.
type Layout int32

const (
	Layout_LAYOUT_DEFAULT    Layout = 0 // default direction: vertical for map, horizontal for list
	Layout_LAYOUT_VERTICAL   Layout = 1 // vertical direction
	Layout_LAYOUT_HORIZONTAL Layout = 2 // horizontal direction
)

// Enum value maps for Layout.
var (
	Layout_name = map[int32]string{
		0: "LAYOUT_DEFAULT",
		1: "LAYOUT_VERTICAL",
		2: "LAYOUT_HORIZONTAL",
	}
	Layout_value = map[string]int32{
		"LAYOUT_DEFAULT":    0,
		"LAYOUT_VERTICAL":   1,
		"LAYOUT_HORIZONTAL": 2,
	}
)

func (x Layout) Enum() *Layout {
	p := new(Layout)
	*p = x
	return p
}

func (x Layout) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Layout) Descriptor() protoreflect.EnumDescriptor {
	return file_tableau_proto_enumTypes[1].Descriptor()
}

func (Layout) Type() protoreflect.EnumType {
	return &file_tableau_proto_enumTypes[1]
}

func (x Layout) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Layout.Descriptor instead.
func (Layout) EnumDescriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{1}
}

type WorkbookOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // workbook name
}

func (x *WorkbookOptions) Reset() {
	*x = WorkbookOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkbookOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkbookOptions) ProtoMessage() {}

func (x *WorkbookOptions) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use WorkbookOptions.ProtoReflect.Descriptor instead.
func (*WorkbookOptions) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{0}
}

func (x *WorkbookOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type WorksheetOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`            // worksheet name
	Namerow   int32  `protobuf:"varint,2,opt,name=namerow,proto3" json:"namerow,omitempty"`     // [default = 1]; // exact row number of name at worksheet.
	Typerow   int32  `protobuf:"varint,3,opt,name=typerow,proto3" json:"typerow,omitempty"`     // [default = 2]; // exact row number of type at worksheet, for generating protos.
	Noterow   int32  `protobuf:"varint,4,opt,name=noterow,proto3" json:"noterow,omitempty"`     // [default = 3]; // exact row number of note at worksheet.
	Datarow   int32  `protobuf:"varint,5,opt,name=datarow,proto3" json:"datarow,omitempty"`     // [default = 4]; // start row number of data at worksheet.
	Transpose bool   `protobuf:"varint,6,opt,name=transpose,proto3" json:"transpose,omitempty"` // [default = false]; // interchange the rows and columns
	Tags      string `protobuf:"bytes,7,opt,name=tags,proto3" json:"tags,omitempty"`            // [default = ""]; // tags for usage, e.g.: "1,2" specifying loading servers. Speciallly, "*"
	// means all servers.
	Nameline int32 `protobuf:"varint,8,opt,name=nameline,proto3" json:"nameline,omitempty"` // [default = 0]; // specify which line in cell as name, '0' means the whole cell is name.
	Typeline int32 `protobuf:"varint,9,opt,name=typeline,proto3" json:"typeline,omitempty"` // [default = 0]; // specify which line in cell as type,'0' means the whole cell is type.
	Nested   bool  `protobuf:"varint,10,opt,name=nested,proto3" json:"nested,omitempty"`    // [default = false]; // whether the naming of name row is nested.
}

func (x *WorksheetOptions) Reset() {
	*x = WorksheetOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorksheetOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorksheetOptions) ProtoMessage() {}

func (x *WorksheetOptions) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use WorksheetOptions.ProtoReflect.Descriptor instead.
func (*WorksheetOptions) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{1}
}

func (x *WorksheetOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *WorksheetOptions) GetNamerow() int32 {
	if x != nil {
		return x.Namerow
	}
	return 0
}

func (x *WorksheetOptions) GetTyperow() int32 {
	if x != nil {
		return x.Typerow
	}
	return 0
}

func (x *WorksheetOptions) GetNoterow() int32 {
	if x != nil {
		return x.Noterow
	}
	return 0
}

func (x *WorksheetOptions) GetDatarow() int32 {
	if x != nil {
		return x.Datarow
	}
	return 0
}

func (x *WorksheetOptions) GetTranspose() bool {
	if x != nil {
		return x.Transpose
	}
	return false
}

func (x *WorksheetOptions) GetTags() string {
	if x != nil {
		return x.Tags
	}
	return ""
}

func (x *WorksheetOptions) GetNameline() int32 {
	if x != nil {
		return x.Nameline
	}
	return 0
}

func (x *WorksheetOptions) GetTypeline() int32 {
	if x != nil {
		return x.Typeline
	}
	return 0
}

func (x *WorksheetOptions) GetNested() bool {
	if x != nil {
		return x.Nested
	}
	return false
}

type FieldOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                          // scalar type's varible name or composite type's varible name (prefix)
	Note     string `protobuf:"bytes,2,opt,name=note,proto3" json:"note,omitempty"`                          // note of name, maybe in another language (Chinese) other than name (English)
	Type     Type   `protobuf:"varint,3,opt,name=type,proto3,enum=tableau.Type" json:"type,omitempty"`       // default: TYPE_DEFAULT
	Key      string `protobuf:"bytes,4,opt,name=key,proto3" json:"key,omitempty"`                            // only set when type is map
	Layout   Layout `protobuf:"varint,5,opt,name=layout,proto3,enum=tableau.Layout" json:"layout,omitempty"` // default: LAYOUT_DEFAULT
	Sep      string `protobuf:"bytes,6,opt,name=sep,proto3" json:"sep,omitempty"`                            // separator, default: ","
	Subsep   string `protobuf:"bytes,7,opt,name=subsep,proto3" json:"subsep,omitempty"`                      // sub separator, default: ":"
	Optional bool   `protobuf:"varint,8,opt,name=optional,proto3" json:"optional,omitempty"`                 // whether the field is optional.
	/////////////////////////////
	// Simple Validators Below //
	/////////////////////////////
	// Different meanings:
	// repeated: size range of array
	// integer: value range
	// string: count of utf-8 code point
	Min   int32  `protobuf:"varint,11,opt,name=min,proto3" json:"min,omitempty"`    // min value
	Max   int32  `protobuf:"varint,12,opt,name=max,proto3" json:"max,omitempty"`    // max value
	Range string `protobuf:"bytes,13,opt,name=range,proto3" json:"range,omitempty"` // format like set description: [1,10], (1,10], [1,10), [1,~]
}

func (x *FieldOptions) Reset() {
	*x = FieldOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOptions) ProtoMessage() {}

func (x *FieldOptions) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use FieldOptions.ProtoReflect.Descriptor instead.
func (*FieldOptions) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{2}
}

func (x *FieldOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *FieldOptions) GetNote() string {
	if x != nil {
		return x.Note
	}
	return ""
}

func (x *FieldOptions) GetType() Type {
	if x != nil {
		return x.Type
	}
	return Type_TYPE_DEFAULT
}

func (x *FieldOptions) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *FieldOptions) GetLayout() Layout {
	if x != nil {
		return x.Layout
	}
	return Layout_LAYOUT_DEFAULT
}

func (x *FieldOptions) GetSep() string {
	if x != nil {
		return x.Sep
	}
	return ""
}

func (x *FieldOptions) GetSubsep() string {
	if x != nil {
		return x.Subsep
	}
	return ""
}

func (x *FieldOptions) GetOptional() bool {
	if x != nil {
		return x.Optional
	}
	return false
}

func (x *FieldOptions) GetMin() int32 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *FieldOptions) GetMax() int32 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *FieldOptions) GetRange() string {
	if x != nil {
		return x.Range
	}
	return ""
}

type EnumOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // alias name
}

func (x *EnumOptions) Reset() {
	*x = EnumOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumOptions) ProtoMessage() {}

func (x *EnumOptions) ProtoReflect() protoreflect.Message {
	mi := &file_tableau_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumOptions.ProtoReflect.Descriptor instead.
func (*EnumOptions) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{3}
}

func (x *EnumOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type EnumValueOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"` // alias name
}

func (x *EnumValueOptions) Reset() {
	*x = EnumValueOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tableau_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnumValueOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnumValueOptions) ProtoMessage() {}

func (x *EnumValueOptions) ProtoReflect() protoreflect.Message {
	mi := &file_tableau_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnumValueOptions.ProtoReflect.Descriptor instead.
func (*EnumValueOptions) Descriptor() ([]byte, []int) {
	return file_tableau_proto_rawDescGZIP(), []int{4}
}

func (x *EnumValueOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var file_tableau_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptor.FileOptions)(nil),
		ExtensionType: (*WorkbookOptions)(nil),
		Field:         50000,
		Name:          "tableau.workbook",
		Tag:           "bytes,50000,opt,name=workbook",
		Filename:      "tableau.proto",
	},
	{
		ExtendedType:  (*descriptor.MessageOptions)(nil),
		ExtensionType: (*WorksheetOptions)(nil),
		Field:         50000,
		Name:          "tableau.worksheet",
		Tag:           "bytes,50000,opt,name=worksheet",
		Filename:      "tableau.proto",
	},
	{
		ExtendedType:  (*descriptor.FieldOptions)(nil),
		ExtensionType: (*FieldOptions)(nil),
		Field:         50000,
		Name:          "tableau.field",
		Tag:           "bytes,50000,opt,name=field",
		Filename:      "tableau.proto",
	},
	{
		ExtendedType:  (*descriptor.EnumOptions)(nil),
		ExtensionType: (*EnumOptions)(nil),
		Field:         50000,
		Name:          "tableau.enum",
		Tag:           "bytes,50000,opt,name=enum",
		Filename:      "tableau.proto",
	},
	{
		ExtendedType:  (*descriptor.EnumValueOptions)(nil),
		ExtensionType: (*EnumValueOptions)(nil),
		Field:         50000,
		Name:          "tableau.evalue",
		Tag:           "bytes,50000,opt,name=evalue",
		Filename:      "tableau.proto",
	},
}

// Extension fields to descriptor.FileOptions.
var (
	// optional tableau.WorkbookOptions workbook = 50000;
	E_Workbook = &file_tableau_proto_extTypes[0]
)

// Extension fields to descriptor.MessageOptions.
var (
	// optional tableau.WorksheetOptions worksheet = 50000;
	E_Worksheet = &file_tableau_proto_extTypes[1]
)

// Extension fields to descriptor.FieldOptions.
var (
	// optional tableau.FieldOptions field = 50000;
	E_Field = &file_tableau_proto_extTypes[2]
)

// Extension fields to descriptor.EnumOptions.
var (
	// optional tableau.EnumOptions enum = 50000;
	E_Enum = &file_tableau_proto_extTypes[3]
)

// Extension fields to descriptor.EnumValueOptions.
var (
	// optional tableau.EnumValueOptions evalue = 50000;
	E_Evalue = &file_tableau_proto_extTypes[4]
)

var File_tableau_proto protoreflect.FileDescriptor

var file_tableau_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0f, 0x57, 0x6f,
	0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x90, 0x02, 0x0a, 0x10, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x68, 0x65, 0x65, 0x74, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x61,
	0x6d, 0x65, 0x72, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6e, 0x61, 0x6d,
	0x65, 0x72, 0x6f, 0x77, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x79, 0x70, 0x65, 0x72, 0x6f, 0x77, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x79, 0x70, 0x65, 0x72, 0x6f, 0x77, 0x12, 0x18,
	0x0a, 0x07, 0x6e, 0x6f, 0x74, 0x65, 0x72, 0x6f, 0x77, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x6e, 0x6f, 0x74, 0x65, 0x72, 0x6f, 0x77, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x61, 0x74, 0x61,
	0x72, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x64, 0x61, 0x74, 0x61, 0x72,
	0x6f, 0x77, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x73, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x70, 0x6f, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x61, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x61, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x74, 0x79, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x74, 0x79, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x6e, 0x65,
	0x73, 0x74, 0x65, 0x64, 0x22, 0x94, 0x02, 0x0a, 0x0c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x6f, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x6f, 0x74, 0x65, 0x12, 0x21, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x27, 0x0a, 0x06, 0x6c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x4c, 0x61, 0x79,
	0x6f, 0x75, 0x74, 0x52, 0x06, 0x6c, 0x61, 0x79, 0x6f, 0x75, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x73,
	0x65, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x65, 0x70, 0x12, 0x16, 0x0a,
	0x06, 0x73, 0x75, 0x62, 0x73, 0x65, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73,
	0x75, 0x62, 0x73, 0x65, 0x70, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x61,
	0x6c, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x0d,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x21, 0x0a, 0x0b, 0x45,
	0x6e, 0x75, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x26,
	0x0a, 0x10, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x2a, 0xa0, 0x02, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x10, 0x0a, 0x0c, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x46, 0x41, 0x55, 0x4c, 0x54, 0x10,
	0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x33, 0x32, 0x10,
	0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x49, 0x4e, 0x54, 0x33, 0x32,
	0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x54, 0x36, 0x34,
	0x10, 0x03, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x49, 0x4e, 0x54, 0x36,
	0x34, 0x10, 0x04, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x4f, 0x55, 0x42,
	0x4c, 0x45, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x4c, 0x4f,
	0x41, 0x54, 0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x4f, 0x4f,
	0x4c, 0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x52, 0x49,
	0x4e, 0x47, 0x10, 0x08, 0x12, 0x0f, 0x0a, 0x0b, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x52,
	0x55, 0x43, 0x54, 0x10, 0x0a, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x41,
	0x54, 0x45, 0x10, 0x15, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x49, 0x4d,
	0x45, 0x10, 0x16, 0x12, 0x11, 0x0a, 0x0d, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x41, 0x54, 0x45,
	0x54, 0x49, 0x4d, 0x45, 0x10, 0x17, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49,
	0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x5f, 0x4c, 0x49, 0x53, 0x54, 0x10, 0x18, 0x12, 0x13, 0x0a, 0x0f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x43, 0x45, 0x4c, 0x4c, 0x5f, 0x4d, 0x41, 0x50, 0x10,
	0x19, 0x12, 0x16, 0x0a, 0x12, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x49, 0x4e, 0x43, 0x45, 0x4c, 0x4c,
	0x5f, 0x53, 0x54, 0x52, 0x55, 0x43, 0x54, 0x10, 0x1a, 0x2a, 0x48, 0x0a, 0x06, 0x4c, 0x61, 0x79,
	0x6f, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x41, 0x59, 0x4f, 0x55, 0x54, 0x5f, 0x44, 0x45,
	0x46, 0x41, 0x55, 0x4c, 0x54, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x41, 0x59, 0x4f, 0x55,
	0x54, 0x5f, 0x56, 0x45, 0x52, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11,
	0x4c, 0x41, 0x59, 0x4f, 0x55, 0x54, 0x5f, 0x48, 0x4f, 0x52, 0x49, 0x5a, 0x4f, 0x4e, 0x54, 0x41,
	0x4c, 0x10, 0x02, 0x3a, 0x54, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e,
	0x57, 0x6f, 0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x08, 0x77, 0x6f, 0x72, 0x6b, 0x62, 0x6f, 0x6f, 0x6b, 0x3a, 0x5a, 0x0a, 0x09, 0x77, 0x6f, 0x72,
	0x6b, 0x73, 0x68, 0x65, 0x65, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x19, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x73, 0x68,
	0x65, 0x65, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x6b,
	0x73, 0x68, 0x65, 0x65, 0x74, 0x3a, 0x4c, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x1d,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x3a, 0x48, 0x0a, 0x04, 0x65, 0x6e, 0x75, 0x6d, 0x12, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e,
	0x75, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x45, 0x6e, 0x75, 0x6d,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x04, 0x65, 0x6e, 0x75, 0x6d, 0x3a, 0x56, 0x0a,
	0x06, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd0, 0x86, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x45, 0x6e, 0x75,
	0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x06, 0x65,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x65, 0x6e, 0x63, 0x68, 0x79, 0x2f, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x61, 0x75, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_tableau_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_tableau_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_tableau_proto_goTypes = []interface{}{
	(Type)(0),                           // 0: tableau.Type
	(Layout)(0),                         // 1: tableau.Layout
	(*WorkbookOptions)(nil),             // 2: tableau.WorkbookOptions
	(*WorksheetOptions)(nil),            // 3: tableau.WorksheetOptions
	(*FieldOptions)(nil),                // 4: tableau.FieldOptions
	(*EnumOptions)(nil),                 // 5: tableau.EnumOptions
	(*EnumValueOptions)(nil),            // 6: tableau.EnumValueOptions
	(*descriptor.FileOptions)(nil),      // 7: google.protobuf.FileOptions
	(*descriptor.MessageOptions)(nil),   // 8: google.protobuf.MessageOptions
	(*descriptor.FieldOptions)(nil),     // 9: google.protobuf.FieldOptions
	(*descriptor.EnumOptions)(nil),      // 10: google.protobuf.EnumOptions
	(*descriptor.EnumValueOptions)(nil), // 11: google.protobuf.EnumValueOptions
}
var file_tableau_proto_depIdxs = []int32{
	0,  // 0: tableau.FieldOptions.type:type_name -> tableau.Type
	1,  // 1: tableau.FieldOptions.layout:type_name -> tableau.Layout
	7,  // 2: tableau.workbook:extendee -> google.protobuf.FileOptions
	8,  // 3: tableau.worksheet:extendee -> google.protobuf.MessageOptions
	9,  // 4: tableau.field:extendee -> google.protobuf.FieldOptions
	10, // 5: tableau.enum:extendee -> google.protobuf.EnumOptions
	11, // 6: tableau.evalue:extendee -> google.protobuf.EnumValueOptions
	2,  // 7: tableau.workbook:type_name -> tableau.WorkbookOptions
	3,  // 8: tableau.worksheet:type_name -> tableau.WorksheetOptions
	4,  // 9: tableau.field:type_name -> tableau.FieldOptions
	5,  // 10: tableau.enum:type_name -> tableau.EnumOptions
	6,  // 11: tableau.evalue:type_name -> tableau.EnumValueOptions
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	7,  // [7:12] is the sub-list for extension type_name
	2,  // [2:7] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_tableau_proto_init() }
func file_tableau_proto_init() {
	if File_tableau_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tableau_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkbookOptions); i {
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
			switch v := v.(*WorksheetOptions); i {
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
			switch v := v.(*FieldOptions); i {
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
		file_tableau_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumOptions); i {
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
		file_tableau_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnumValueOptions); i {
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
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 5,
			NumServices:   0,
		},
		GoTypes:           file_tableau_proto_goTypes,
		DependencyIndexes: file_tableau_proto_depIdxs,
		EnumInfos:         file_tableau_proto_enumTypes,
		MessageInfos:      file_tableau_proto_msgTypes,
		ExtensionInfos:    file_tableau_proto_extTypes,
	}.Build()
	File_tableau_proto = out.File
	file_tableau_proto_rawDesc = nil
	file_tableau_proto_goTypes = nil
	file_tableau_proto_depIdxs = nil
}
