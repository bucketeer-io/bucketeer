// Copyright 2025 The Bucketeer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v4.23.4
// source: proto/feature/feature_last_used_info.proto

package feature

import (
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

type FeatureLastUsedInfo_Status int32

const (
	FeatureLastUsedInfo_UNKNOWN     FeatureLastUsedInfo_Status = 0
	FeatureLastUsedInfo_NEW         FeatureLastUsedInfo_Status = 1
	FeatureLastUsedInfo_ACTIVE      FeatureLastUsedInfo_Status = 2
	FeatureLastUsedInfo_NO_ACTIVITY FeatureLastUsedInfo_Status = 3
)

// Enum value maps for FeatureLastUsedInfo_Status.
var (
	FeatureLastUsedInfo_Status_name = map[int32]string{
		0: "UNKNOWN",
		1: "NEW",
		2: "ACTIVE",
		3: "NO_ACTIVITY",
	}
	FeatureLastUsedInfo_Status_value = map[string]int32{
		"UNKNOWN":     0,
		"NEW":         1,
		"ACTIVE":      2,
		"NO_ACTIVITY": 3,
	}
)

func (x FeatureLastUsedInfo_Status) Enum() *FeatureLastUsedInfo_Status {
	p := new(FeatureLastUsedInfo_Status)
	*p = x
	return p
}

func (x FeatureLastUsedInfo_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FeatureLastUsedInfo_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_feature_feature_last_used_info_proto_enumTypes[0].Descriptor()
}

func (FeatureLastUsedInfo_Status) Type() protoreflect.EnumType {
	return &file_proto_feature_feature_last_used_info_proto_enumTypes[0]
}

func (x FeatureLastUsedInfo_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FeatureLastUsedInfo_Status.Descriptor instead.
func (FeatureLastUsedInfo_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_feature_feature_last_used_info_proto_rawDescGZIP(), []int{0, 0}
}

type FeatureLastUsedInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FeatureId           string `protobuf:"bytes,1,opt,name=feature_id,json=featureId,proto3" json:"feature_id"`
	Version             int32  `protobuf:"varint,2,opt,name=version,proto3" json:"version"`
	LastUsedAt          int64  `protobuf:"varint,3,opt,name=last_used_at,json=lastUsedAt,proto3" json:"last_used_at"`
	CreatedAt           int64  `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	ClientOldestVersion string `protobuf:"bytes,5,opt,name=client_oldest_version,json=clientOldestVersion,proto3" json:"client_oldest_version"`
	ClientLatestVersion string `protobuf:"bytes,6,opt,name=client_latest_version,json=clientLatestVersion,proto3" json:"client_latest_version"`
}

func (x *FeatureLastUsedInfo) Reset() {
	*x = FeatureLastUsedInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_feature_last_used_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FeatureLastUsedInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FeatureLastUsedInfo) ProtoMessage() {}

func (x *FeatureLastUsedInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_feature_last_used_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FeatureLastUsedInfo.ProtoReflect.Descriptor instead.
func (*FeatureLastUsedInfo) Descriptor() ([]byte, []int) {
	return file_proto_feature_feature_last_used_info_proto_rawDescGZIP(), []int{0}
}

func (x *FeatureLastUsedInfo) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

func (x *FeatureLastUsedInfo) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *FeatureLastUsedInfo) GetLastUsedAt() int64 {
	if x != nil {
		return x.LastUsedAt
	}
	return 0
}

func (x *FeatureLastUsedInfo) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *FeatureLastUsedInfo) GetClientOldestVersion() string {
	if x != nil {
		return x.ClientOldestVersion
	}
	return ""
}

func (x *FeatureLastUsedInfo) GetClientLatestVersion() string {
	if x != nil {
		return x.ClientLatestVersion
	}
	return ""
}

var File_proto_feature_feature_last_used_info_proto protoreflect.FileDescriptor

var file_proto_feature_feature_last_used_info_proto_rawDesc = []byte{
	0x0a, 0x2a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65,
	0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22,
	0xb4, 0x02, 0x0a, 0x13, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x55,
	0x73, 0x65, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x12, 0x20, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x73, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x6c, 0x64, 0x65,
	0x73, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4f, 0x6c, 0x64, 0x65, 0x73, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x15, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4c, 0x61, 0x74, 0x65,
	0x73, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x3b, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00,
	0x12, 0x07, 0x0a, 0x03, 0x4e, 0x45, 0x57, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x54,
	0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x4e, 0x4f, 0x5f, 0x41, 0x43, 0x54, 0x49,
	0x56, 0x49, 0x54, 0x59, 0x10, 0x03, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69,
	0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_feature_feature_last_used_info_proto_rawDescOnce sync.Once
	file_proto_feature_feature_last_used_info_proto_rawDescData = file_proto_feature_feature_last_used_info_proto_rawDesc
)

func file_proto_feature_feature_last_used_info_proto_rawDescGZIP() []byte {
	file_proto_feature_feature_last_used_info_proto_rawDescOnce.Do(func() {
		file_proto_feature_feature_last_used_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_feature_feature_last_used_info_proto_rawDescData)
	})
	return file_proto_feature_feature_last_used_info_proto_rawDescData
}

var file_proto_feature_feature_last_used_info_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_feature_feature_last_used_info_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_feature_feature_last_used_info_proto_goTypes = []interface{}{
	(FeatureLastUsedInfo_Status)(0), // 0: bucketeer.feature.FeatureLastUsedInfo.Status
	(*FeatureLastUsedInfo)(nil),     // 1: bucketeer.feature.FeatureLastUsedInfo
}
var file_proto_feature_feature_last_used_info_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_feature_feature_last_used_info_proto_init() }
func file_proto_feature_feature_last_used_info_proto_init() {
	if File_proto_feature_feature_last_used_info_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_feature_feature_last_used_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FeatureLastUsedInfo); i {
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
			RawDescriptor: file_proto_feature_feature_last_used_info_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_feature_feature_last_used_info_proto_goTypes,
		DependencyIndexes: file_proto_feature_feature_last_used_info_proto_depIdxs,
		EnumInfos:         file_proto_feature_feature_last_used_info_proto_enumTypes,
		MessageInfos:      file_proto_feature_feature_last_used_info_proto_msgTypes,
	}.Build()
	File_proto_feature_feature_last_used_info_proto = out.File
	file_proto_feature_feature_last_used_info_proto_rawDesc = nil
	file_proto_feature_feature_last_used_info_proto_goTypes = nil
	file_proto_feature_feature_last_used_info_proto_depIdxs = nil
}
