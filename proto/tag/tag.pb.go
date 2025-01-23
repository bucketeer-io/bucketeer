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
// source: proto/tag/tag.proto

package tag

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

type Tag_EntityType int32

const (
	Tag_UNSPECIFIED  Tag_EntityType = 0
	Tag_FEATURE_FLAG Tag_EntityType = 1
	Tag_ACCOUNT      Tag_EntityType = 2
)

// Enum value maps for Tag_EntityType.
var (
	Tag_EntityType_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "FEATURE_FLAG",
		2: "ACCOUNT",
	}
	Tag_EntityType_value = map[string]int32{
		"UNSPECIFIED":  0,
		"FEATURE_FLAG": 1,
		"ACCOUNT":      2,
	}
)

func (x Tag_EntityType) Enum() *Tag_EntityType {
	p := new(Tag_EntityType)
	*p = x
	return p
}

func (x Tag_EntityType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Tag_EntityType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_tag_tag_proto_enumTypes[0].Descriptor()
}

func (Tag_EntityType) Type() protoreflect.EnumType {
	return &file_proto_tag_tag_proto_enumTypes[0]
}

func (x Tag_EntityType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Tag_EntityType.Descriptor instead.
func (Tag_EntityType) EnumDescriptor() ([]byte, []int) {
	return file_proto_tag_tag_proto_rawDescGZIP(), []int{0, 0}
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name          string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	CreatedAt     int64          `protobuf:"varint,3,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt     int64          `protobuf:"varint,4,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	EntityType    Tag_EntityType `protobuf:"varint,5,opt,name=entity_type,json=entityType,proto3,enum=bucketeer.tag.Tag_EntityType" json:"entity_type"`
	EnvironmentId string         `protobuf:"bytes,6,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tag_tag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tag_tag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Tag.ProtoReflect.Descriptor instead.
func (*Tag) Descriptor() ([]byte, []int) {
	return file_proto_tag_tag_proto_rawDescGZIP(), []int{0}
}

func (x *Tag) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Tag) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Tag) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Tag) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Tag) GetEntityType() Tag_EntityType {
	if x != nil {
		return x.EntityType
	}
	return Tag_UNSPECIFIED
}

func (x *Tag) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

type EnvironmentTag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvironmentId string `protobuf:"bytes,1,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
	Tags          []*Tag `protobuf:"bytes,2,rep,name=tags,proto3" json:"tags"`
}

func (x *EnvironmentTag) Reset() {
	*x = EnvironmentTag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_tag_tag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentTag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentTag) ProtoMessage() {}

func (x *EnvironmentTag) ProtoReflect() protoreflect.Message {
	mi := &file_proto_tag_tag_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentTag.ProtoReflect.Descriptor instead.
func (*EnvironmentTag) Descriptor() ([]byte, []int) {
	return file_proto_tag_tag_proto_rawDescGZIP(), []int{1}
}

func (x *EnvironmentTag) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

func (x *EnvironmentTag) GetTags() []*Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

var File_proto_tag_tag_proto protoreflect.FileDescriptor

var file_proto_tag_tag_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x61, 0x67, 0x2f, 0x74, 0x61, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x74, 0x61, 0x67, 0x22, 0x8c, 0x02, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3e,
	0x0a, 0x0b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x74, 0x61, 0x67, 0x2e, 0x54, 0x61, 0x67, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0a, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x25,
	0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x3c, 0x0a, 0x0a, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x46, 0x45, 0x41, 0x54, 0x55, 0x52, 0x45, 0x5f,
	0x46, 0x4c, 0x41, 0x47, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e,
	0x54, 0x10, 0x02, 0x22, 0x5f, 0x0a, 0x0e, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x54, 0x61, 0x67, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x74, 0x61, 0x67, 0x2e, 0x54, 0x61, 0x67, 0x52, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x74, 0x61, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_tag_tag_proto_rawDescOnce sync.Once
	file_proto_tag_tag_proto_rawDescData = file_proto_tag_tag_proto_rawDesc
)

func file_proto_tag_tag_proto_rawDescGZIP() []byte {
	file_proto_tag_tag_proto_rawDescOnce.Do(func() {
		file_proto_tag_tag_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_tag_tag_proto_rawDescData)
	})
	return file_proto_tag_tag_proto_rawDescData
}

var file_proto_tag_tag_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_tag_tag_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_tag_tag_proto_goTypes = []interface{}{
	(Tag_EntityType)(0),    // 0: bucketeer.tag.Tag.EntityType
	(*Tag)(nil),            // 1: bucketeer.tag.Tag
	(*EnvironmentTag)(nil), // 2: bucketeer.tag.EnvironmentTag
}
var file_proto_tag_tag_proto_depIdxs = []int32{
	0, // 0: bucketeer.tag.Tag.entity_type:type_name -> bucketeer.tag.Tag.EntityType
	1, // 1: bucketeer.tag.EnvironmentTag.tags:type_name -> bucketeer.tag.Tag
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_tag_tag_proto_init() }
func file_proto_tag_tag_proto_init() {
	if File_proto_tag_tag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_tag_tag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Tag); i {
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
		file_proto_tag_tag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentTag); i {
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
			RawDescriptor: file_proto_tag_tag_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_tag_tag_proto_goTypes,
		DependencyIndexes: file_proto_tag_tag_proto_depIdxs,
		EnumInfos:         file_proto_tag_tag_proto_enumTypes,
		MessageInfos:      file_proto_tag_tag_proto_msgTypes,
	}.Build()
	File_proto_tag_tag_proto = out.File
	file_proto_tag_tag_proto_rawDesc = nil
	file_proto_tag_tag_proto_goTypes = nil
	file_proto_tag_tag_proto_depIdxs = nil
}
