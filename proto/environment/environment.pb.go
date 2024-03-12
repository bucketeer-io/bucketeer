// Copyright 2024 The Bucketeer Authors.
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
// source: proto/environment/environment.proto

package environment

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

type EnvironmentV2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name           string                  `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	UrlCode        string                  `protobuf:"bytes,3,opt,name=url_code,json=urlCode,proto3" json:"url_code"`
	Description    string                  `protobuf:"bytes,4,opt,name=description,proto3" json:"description"` // optional
	ProjectId      string                  `protobuf:"bytes,5,opt,name=project_id,json=projectId,proto3" json:"project_id"`
	Archived       bool                    `protobuf:"varint,6,opt,name=archived,proto3" json:"archived"`
	CreatedAt      int64                   `protobuf:"varint,7,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt      int64                   `protobuf:"varint,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	OrganizationId string                  `protobuf:"bytes,9,opt,name=organization_id,json=organizationId,proto3" json:"organization_id"`
	Settings       *EnvironmentV2_Settings `protobuf:"bytes,10,opt,name=settings,proto3" json:"settings"`
}

func (x *EnvironmentV2) Reset() {
	*x = EnvironmentV2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_environment_environment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentV2) ProtoMessage() {}

func (x *EnvironmentV2) ProtoReflect() protoreflect.Message {
	mi := &file_proto_environment_environment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentV2.ProtoReflect.Descriptor instead.
func (*EnvironmentV2) Descriptor() ([]byte, []int) {
	return file_proto_environment_environment_proto_rawDescGZIP(), []int{0}
}

func (x *EnvironmentV2) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *EnvironmentV2) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *EnvironmentV2) GetUrlCode() string {
	if x != nil {
		return x.UrlCode
	}
	return ""
}

func (x *EnvironmentV2) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *EnvironmentV2) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *EnvironmentV2) GetArchived() bool {
	if x != nil {
		return x.Archived
	}
	return false
}

func (x *EnvironmentV2) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *EnvironmentV2) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *EnvironmentV2) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *EnvironmentV2) GetSettings() *EnvironmentV2_Settings {
	if x != nil {
		return x.Settings
	}
	return nil
}

type EnvironmentV2_Settings struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequireComment bool `protobuf:"varint,1,opt,name=require_comment,json=requireComment,proto3" json:"require_comment"`
}

func (x *EnvironmentV2_Settings) Reset() {
	*x = EnvironmentV2_Settings{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_environment_environment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentV2_Settings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentV2_Settings) ProtoMessage() {}

func (x *EnvironmentV2_Settings) ProtoReflect() protoreflect.Message {
	mi := &file_proto_environment_environment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentV2_Settings.ProtoReflect.Descriptor instead.
func (*EnvironmentV2_Settings) Descriptor() ([]byte, []int) {
	return file_proto_environment_environment_proto_rawDescGZIP(), []int{0, 0}
}

func (x *EnvironmentV2_Settings) GetRequireComment() bool {
	if x != nil {
		return x.RequireComment
	}
	return false
}

var File_proto_environment_environment_proto protoreflect.FileDescriptor

var file_proto_environment_environment_proto_rawDesc = []byte{
	0x0a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x92, 0x03, 0x0a,
	0x0d, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x32, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x72, 0x6c, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x75, 0x72, 0x6c, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x12, 0x49, 0x0a, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e,
	0x67, 0x73, 0x52, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x1a, 0x33, 0x0a, 0x08,
	0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x71, 0x75,
	0x69, 0x72, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_environment_environment_proto_rawDescOnce sync.Once
	file_proto_environment_environment_proto_rawDescData = file_proto_environment_environment_proto_rawDesc
)

func file_proto_environment_environment_proto_rawDescGZIP() []byte {
	file_proto_environment_environment_proto_rawDescOnce.Do(func() {
		file_proto_environment_environment_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_environment_environment_proto_rawDescData)
	})
	return file_proto_environment_environment_proto_rawDescData
}

var file_proto_environment_environment_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_environment_environment_proto_goTypes = []interface{}{
	(*EnvironmentV2)(nil),          // 0: bucketeer.environment.EnvironmentV2
	(*EnvironmentV2_Settings)(nil), // 1: bucketeer.environment.EnvironmentV2.Settings
}
var file_proto_environment_environment_proto_depIdxs = []int32{
	1, // 0: bucketeer.environment.EnvironmentV2.settings:type_name -> bucketeer.environment.EnvironmentV2.Settings
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_environment_environment_proto_init() }
func file_proto_environment_environment_proto_init() {
	if File_proto_environment_environment_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_environment_environment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentV2); i {
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
		file_proto_environment_environment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentV2_Settings); i {
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
			RawDescriptor: file_proto_environment_environment_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_environment_environment_proto_goTypes,
		DependencyIndexes: file_proto_environment_environment_proto_depIdxs,
		MessageInfos:      file_proto_environment_environment_proto_msgTypes,
	}.Build()
	File_proto_environment_environment_proto = out.File
	file_proto_environment_environment_proto_rawDesc = nil
	file_proto_environment_environment_proto_goTypes = nil
	file_proto_environment_environment_proto_depIdxs = nil
}
