// Generated by tableauc 0.1.1. DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: feature_toggle.proto

package testpb

import (
	_ "github.com/Wenchy/tableau/proto/tableaupb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ToggleCfg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeatureToggleMap map[int32]*ToggleCfg_FeatureToggle `protobuf:"bytes,1,rep,name=feature_toggle_map,json=featureToggleMap,proto3" json:"feature_toggle_map,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ToggleCfg) Reset() {
	*x = ToggleCfg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feature_toggle_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToggleCfg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToggleCfg) ProtoMessage() {}

func (x *ToggleCfg) ProtoReflect() protoreflect.Message {
	mi := &file_feature_toggle_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToggleCfg.ProtoReflect.Descriptor instead.
func (*ToggleCfg) Descriptor() ([]byte, []int) {
	return file_feature_toggle_proto_rawDescGZIP(), []int{0}
}

func (x *ToggleCfg) GetFeatureToggleMap() map[int32]*ToggleCfg_FeatureToggle {
	if x != nil {
		return x.FeatureToggleMap
	}
	return nil
}

type ToggleCfg_FeatureToggle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvId      int32                             `protobuf:"varint,1,opt,name=env_id,json=envId,proto3" json:"env_id,omitempty"`
	Desc       string                            `protobuf:"bytes,2,opt,name=desc,proto3" json:"desc,omitempty"`
	ToggleList []*ToggleCfg_FeatureToggle_Toggle `protobuf:"bytes,3,rep,name=toggle_list,json=toggleList,proto3" json:"toggle_list,omitempty"`
}

func (x *ToggleCfg_FeatureToggle) Reset() {
	*x = ToggleCfg_FeatureToggle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feature_toggle_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToggleCfg_FeatureToggle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToggleCfg_FeatureToggle) ProtoMessage() {}

func (x *ToggleCfg_FeatureToggle) ProtoReflect() protoreflect.Message {
	mi := &file_feature_toggle_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToggleCfg_FeatureToggle.ProtoReflect.Descriptor instead.
func (*ToggleCfg_FeatureToggle) Descriptor() ([]byte, []int) {
	return file_feature_toggle_proto_rawDescGZIP(), []int{0, 1}
}

func (x *ToggleCfg_FeatureToggle) GetEnvId() int32 {
	if x != nil {
		return x.EnvId
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

func (x *ToggleCfg_FeatureToggle) GetToggleList() []*ToggleCfg_FeatureToggle_Toggle {
	if x != nil {
		return x.ToggleList
	}
	return nil
}

type ToggleCfg_FeatureToggle_Toggle struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             FeatureToggleMacroType `protobuf:"varint,1,opt,name=id,proto3,enum=testxml.FeatureToggleMacroType" json:"id,omitempty"`
	OpenRate       int32                  `protobuf:"varint,2,opt,name=open_rate,json=openRate,proto3" json:"open_rate,omitempty"`
	WorldId        uint32                 `protobuf:"varint,3,opt,name=world_id,json=worldId,proto3" json:"world_id,omitempty"`
	ZoneId         uint32                 `protobuf:"varint,4,opt,name=zone_id,json=zoneId,proto3" json:"zone_id,omitempty"`
	OpenTime       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=open_time,json=openTime,proto3" json:"open_time,omitempty"`
	CloseTime      *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=close_time,json=closeTime,proto3" json:"close_time,omitempty"`
	UserLimitType  int32                  `protobuf:"varint,7,opt,name=user_limit_type,json=userLimitType,proto3" json:"user_limit_type,omitempty"`
	UserSourceType int32                  `protobuf:"varint,8,opt,name=user_source_type,json=userSourceType,proto3" json:"user_source_type,omitempty"`
	SysKey_1       int64                  `protobuf:"varint,9,opt,name=sys_key_1,json=sysKey1,proto3" json:"sys_key_1,omitempty"`
	SysKey_2       int64                  `protobuf:"varint,10,opt,name=sys_key_2,json=sysKey2,proto3" json:"sys_key_2,omitempty"`
	SysKey_3       int64                  `protobuf:"varint,11,opt,name=sys_key_3,json=sysKey3,proto3" json:"sys_key_3,omitempty"`
	NotifyClosed   bool                   `protobuf:"varint,12,opt,name=notify_closed,json=notifyClosed,proto3" json:"notify_closed,omitempty"`
	ErrCode        Code                   `protobuf:"varint,13,opt,name=err_code,json=errCode,proto3,enum=testxml.Code" json:"err_code,omitempty"`
	Name           string                 `protobuf:"bytes,14,opt,name=name,proto3" json:"name,omitempty"`
	Desc           string                 `protobuf:"bytes,15,opt,name=desc,proto3" json:"desc,omitempty"`
}

func (x *ToggleCfg_FeatureToggle_Toggle) Reset() {
	*x = ToggleCfg_FeatureToggle_Toggle{}
	if protoimpl.UnsafeEnabled {
		mi := &file_feature_toggle_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToggleCfg_FeatureToggle_Toggle) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToggleCfg_FeatureToggle_Toggle) ProtoMessage() {}

func (x *ToggleCfg_FeatureToggle_Toggle) ProtoReflect() protoreflect.Message {
	mi := &file_feature_toggle_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToggleCfg_FeatureToggle_Toggle.ProtoReflect.Descriptor instead.
func (*ToggleCfg_FeatureToggle_Toggle) Descriptor() ([]byte, []int) {
	return file_feature_toggle_proto_rawDescGZIP(), []int{0, 1, 0}
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetId() FeatureToggleMacroType {
	if x != nil {
		return x.Id
	}
	return FeatureToggleMacroType_TOGGLE_SAMPLE
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetOpenRate() int32 {
	if x != nil {
		return x.OpenRate
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetWorldId() uint32 {
	if x != nil {
		return x.WorldId
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetZoneId() uint32 {
	if x != nil {
		return x.ZoneId
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetOpenTime() *timestamppb.Timestamp {
	if x != nil {
		return x.OpenTime
	}
	return nil
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetCloseTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CloseTime
	}
	return nil
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetUserLimitType() int32 {
	if x != nil {
		return x.UserLimitType
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetUserSourceType() int32 {
	if x != nil {
		return x.UserSourceType
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetSysKey_1() int64 {
	if x != nil {
		return x.SysKey_1
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetSysKey_2() int64 {
	if x != nil {
		return x.SysKey_2
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetSysKey_3() int64 {
	if x != nil {
		return x.SysKey_3
	}
	return 0
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetNotifyClosed() bool {
	if x != nil {
		return x.NotifyClosed
	}
	return false
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetErrCode() Code {
	if x != nil {
		return x.ErrCode
	}
	return Code_SUCCESS
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ToggleCfg_FeatureToggle_Toggle) GetDesc() string {
	if x != nil {
		return x.Desc
	}
	return ""
}

var File_feature_toggle_proto protoreflect.FileDescriptor

var file_feature_toggle_proto_rawDesc = []byte{
	0x0a, 0x14, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x74, 0x65, 0x73, 0x74, 0x78, 0x6d, 0x6c, 0x1a,
	0x10, 0x63, 0x73, 0x5f, 0x63, 0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1e, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xc4, 0x09, 0x0a, 0x09, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x43, 0x66, 0x67,
	0x12, 0x74, 0x0a, 0x12, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x74, 0x6f, 0x67, 0x67,
	0x6c, 0x65, 0x5f, 0x6d, 0x61, 0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x78, 0x6d, 0x6c, 0x2e, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x43, 0x66, 0x67,
	0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4d, 0x61,
	0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x42, 0x1c, 0x82, 0xb5, 0x18, 0x18, 0x0a, 0x0d, 0x46, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x22, 0x05, 0x45, 0x6e, 0x76,
	0x49, 0x44, 0x28, 0x01, 0x52, 0x10, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67,
	0x67, 0x6c, 0x65, 0x4d, 0x61, 0x70, 0x1a, 0x65, 0x0a, 0x15, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x36, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x78, 0x6d, 0x6c, 0x2e, 0x54, 0x6f, 0x67, 0x67, 0x6c,
	0x65, 0x43, 0x66, 0x67, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67,
	0x6c, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0xba, 0x07,
	0x0a, 0x0d, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x12,
	0x22, 0x0a, 0x06, 0x65, 0x6e, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x42,
	0x0b, 0x82, 0xb5, 0x18, 0x07, 0x0a, 0x05, 0x45, 0x6e, 0x76, 0x49, 0x44, 0x52, 0x05, 0x65, 0x6e,
	0x76, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x04, 0x64, 0x65, 0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x42, 0x0a, 0x82, 0xb5, 0x18, 0x06, 0x0a, 0x04, 0x44, 0x65, 0x73, 0x63, 0x52, 0x04, 0x64,
	0x65, 0x73, 0x63, 0x12, 0x5c, 0x0a, 0x0b, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x78,
	0x6d, 0x6c, 0x2e, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x43, 0x66, 0x67, 0x2e, 0x46, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x54, 0x6f, 0x67, 0x67, 0x6c,
	0x65, 0x42, 0x12, 0x82, 0xb5, 0x18, 0x0e, 0x0a, 0x06, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x22,
	0x02, 0x49, 0x64, 0x28, 0x01, 0x52, 0x0a, 0x74, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x73,
	0x74, 0x1a, 0x86, 0x06, 0x0a, 0x06, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x12, 0x39, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x78,
	0x6d, 0x6c, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65,
	0x4d, 0x61, 0x63, 0x72, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x42, 0x08, 0x82, 0xb5, 0x18, 0x04, 0x0a,
	0x02, 0x49, 0x64, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2b, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x6e, 0x5f,
	0x72, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x0e, 0x82, 0xb5, 0x18, 0x0a,
	0x0a, 0x08, 0x4f, 0x70, 0x65, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x6e,
	0x52, 0x61, 0x74, 0x65, 0x12, 0x28, 0x0a, 0x08, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x0d, 0x82, 0xb5, 0x18, 0x09, 0x0a, 0x07, 0x57, 0x6f,
	0x72, 0x6c, 0x64, 0x49, 0x44, 0x52, 0x07, 0x77, 0x6f, 0x72, 0x6c, 0x64, 0x49, 0x64, 0x12, 0x25,
	0x0a, 0x07, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x42,
	0x0c, 0x82, 0xb5, 0x18, 0x08, 0x0a, 0x06, 0x5a, 0x6f, 0x6e, 0x65, 0x49, 0x44, 0x52, 0x06, 0x7a,
	0x6f, 0x6e, 0x65, 0x49, 0x64, 0x12, 0x47, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x42, 0x0e, 0x82, 0xb5, 0x18, 0x0a, 0x0a, 0x08, 0x4f, 0x70, 0x65, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x4a,
	0x0a, 0x0a, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x42, 0x0f,
	0x82, 0xb5, 0x18, 0x0b, 0x0a, 0x09, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x52,
	0x09, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x05, 0x42, 0x13, 0x82, 0xb5, 0x18, 0x0f, 0x0a, 0x0d, 0x55, 0x73, 0x65, 0x72, 0x4c,
	0x69, 0x6d, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x3e, 0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x42, 0x14, 0x82, 0xb5, 0x18, 0x10, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x53, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x29, 0x0a, 0x09, 0x73, 0x79, 0x73, 0x5f, 0x6b,
	0x65, 0x79, 0x5f, 0x31, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0d, 0x82, 0xb5, 0x18, 0x09,
	0x0a, 0x07, 0x53, 0x79, 0x73, 0x4b, 0x65, 0x79, 0x31, 0x52, 0x07, 0x73, 0x79, 0x73, 0x4b, 0x65,
	0x79, 0x31, 0x12, 0x29, 0x0a, 0x09, 0x73, 0x79, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x32, 0x18,
	0x0a, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0d, 0x82, 0xb5, 0x18, 0x09, 0x0a, 0x07, 0x53, 0x79, 0x73,
	0x4b, 0x65, 0x79, 0x32, 0x52, 0x07, 0x73, 0x79, 0x73, 0x4b, 0x65, 0x79, 0x32, 0x12, 0x29, 0x0a,
	0x09, 0x73, 0x79, 0x73, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x33, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x0d, 0x82, 0xb5, 0x18, 0x09, 0x0a, 0x07, 0x53, 0x79, 0x73, 0x4b, 0x65, 0x79, 0x33, 0x52,
	0x07, 0x73, 0x79, 0x73, 0x4b, 0x65, 0x79, 0x33, 0x12, 0x37, 0x0a, 0x0d, 0x6e, 0x6f, 0x74, 0x69,
	0x66, 0x79, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x42,
	0x12, 0x82, 0xb5, 0x18, 0x0e, 0x0a, 0x0c, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x43, 0x6c, 0x6f,
	0x73, 0x65, 0x64, 0x52, 0x0c, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x64, 0x12, 0x37, 0x0a, 0x08, 0x65, 0x72, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x78, 0x6d, 0x6c, 0x2e, 0x43, 0x6f,
	0x64, 0x65, 0x42, 0x0d, 0x82, 0xb5, 0x18, 0x09, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64,
	0x65, 0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0x82, 0xb5, 0x18, 0x06, 0x0a, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x64, 0x65,
	0x73, 0x63, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0x82, 0xb5, 0x18, 0x06, 0x0a, 0x04,
	0x44, 0x65, 0x73, 0x63, 0x52, 0x04, 0x64, 0x65, 0x73, 0x63, 0x3a, 0x1d, 0x82, 0xb5, 0x18, 0x19,
	0x0a, 0x09, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x43, 0x66, 0x67, 0x10, 0x01, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x04, 0x40, 0x01, 0x48, 0x01, 0x50, 0x01, 0x42, 0x43, 0x5a, 0x29, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x65, 0x6e, 0x63, 0x68, 0x79, 0x2f, 0x74,
	0x61, 0x62, 0x6c, 0x65, 0x61, 0x75, 0x2f, 0x63, 0x6d, 0x64, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f,
	0x74, 0x65, 0x73, 0x74, 0x70, 0x62, 0x82, 0xb5, 0x18, 0x14, 0x0a, 0x12, 0x46, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x2e, 0x78, 0x6c, 0x73, 0x78, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_feature_toggle_proto_rawDescOnce sync.Once
	file_feature_toggle_proto_rawDescData = file_feature_toggle_proto_rawDesc
)

func file_feature_toggle_proto_rawDescGZIP() []byte {
	file_feature_toggle_proto_rawDescOnce.Do(func() {
		file_feature_toggle_proto_rawDescData = protoimpl.X.CompressGZIP(file_feature_toggle_proto_rawDescData)
	})
	return file_feature_toggle_proto_rawDescData
}

var file_feature_toggle_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_feature_toggle_proto_goTypes = []interface{}{
	(*ToggleCfg)(nil),                      // 0: testxml.ToggleCfg
	nil,                                    // 1: testxml.ToggleCfg.FeatureToggleMapEntry
	(*ToggleCfg_FeatureToggle)(nil),        // 2: testxml.ToggleCfg.FeatureToggle
	(*ToggleCfg_FeatureToggle_Toggle)(nil), // 3: testxml.ToggleCfg.FeatureToggle.Toggle
	(FeatureToggleMacroType)(0),            // 4: testxml.FeatureToggleMacroType
	(*timestamppb.Timestamp)(nil),          // 5: google.protobuf.Timestamp
	(Code)(0),                              // 6: testxml.Code
}
var file_feature_toggle_proto_depIdxs = []int32{
	1, // 0: testxml.ToggleCfg.feature_toggle_map:type_name -> testxml.ToggleCfg.FeatureToggleMapEntry
	2, // 1: testxml.ToggleCfg.FeatureToggleMapEntry.value:type_name -> testxml.ToggleCfg.FeatureToggle
	3, // 2: testxml.ToggleCfg.FeatureToggle.toggle_list:type_name -> testxml.ToggleCfg.FeatureToggle.Toggle
	4, // 3: testxml.ToggleCfg.FeatureToggle.Toggle.id:type_name -> testxml.FeatureToggleMacroType
	5, // 4: testxml.ToggleCfg.FeatureToggle.Toggle.open_time:type_name -> google.protobuf.Timestamp
	5, // 5: testxml.ToggleCfg.FeatureToggle.Toggle.close_time:type_name -> google.protobuf.Timestamp
	6, // 6: testxml.ToggleCfg.FeatureToggle.Toggle.err_code:type_name -> testxml.Code
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_feature_toggle_proto_init() }
func file_feature_toggle_proto_init() {
	if File_feature_toggle_proto != nil {
		return
	}
	file_cs_com_def_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_feature_toggle_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToggleCfg); i {
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
		file_feature_toggle_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToggleCfg_FeatureToggle); i {
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
		file_feature_toggle_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToggleCfg_FeatureToggle_Toggle); i {
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
			RawDescriptor: file_feature_toggle_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_feature_toggle_proto_goTypes,
		DependencyIndexes: file_feature_toggle_proto_depIdxs,
		MessageInfos:      file_feature_toggle_proto_msgTypes,
	}.Build()
	File_feature_toggle_proto = out.File
	file_feature_toggle_proto_rawDesc = nil
	file_feature_toggle_proto_goTypes = nil
	file_feature_toggle_proto_depIdxs = nil
}