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
// source: proto/account/api_key.proto

package account

import (
	environment "github.com/bucketeer-io/bucketeer/proto/environment"
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

type APIKey_Role int32

const (
	APIKey_UNKNOWN              APIKey_Role = 0
	APIKey_SDK_CLIENT           APIKey_Role = 1
	APIKey_SDK_SERVER           APIKey_Role = 2
	APIKey_PUBLIC_API_READ_ONLY APIKey_Role = 3
	APIKey_PUBLIC_API_WRITE     APIKey_Role = 4
	// For sensitive data
	APIKey_PUBLIC_API_ADMIN APIKey_Role = 5
)

// Enum value maps for APIKey_Role.
var (
	APIKey_Role_name = map[int32]string{
		0: "UNKNOWN",
		1: "SDK_CLIENT",
		2: "SDK_SERVER",
		3: "PUBLIC_API_READ_ONLY",
		4: "PUBLIC_API_WRITE",
		5: "PUBLIC_API_ADMIN",
	}
	APIKey_Role_value = map[string]int32{
		"UNKNOWN":              0,
		"SDK_CLIENT":           1,
		"SDK_SERVER":           2,
		"PUBLIC_API_READ_ONLY": 3,
		"PUBLIC_API_WRITE":     4,
		"PUBLIC_API_ADMIN":     5,
	}
)

func (x APIKey_Role) Enum() *APIKey_Role {
	p := new(APIKey_Role)
	*p = x
	return p
}

func (x APIKey_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (APIKey_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_api_key_proto_enumTypes[0].Descriptor()
}

func (APIKey_Role) Type() protoreflect.EnumType {
	return &file_proto_account_api_key_proto_enumTypes[0]
}

func (x APIKey_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use APIKey_Role.Descriptor instead.
func (APIKey_Role) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_api_key_proto_rawDescGZIP(), []int{0, 0}
}

type APIKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name      string      `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Role      APIKey_Role `protobuf:"varint,3,opt,name=role,proto3,enum=bucketeer.account.APIKey_Role" json:"role"`
	Disabled  bool        `protobuf:"varint,4,opt,name=disabled,proto3" json:"disabled"`
	CreatedAt int64       `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt int64       `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
}

func (x *APIKey) Reset() {
	*x = APIKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_api_key_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *APIKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*APIKey) ProtoMessage() {}

func (x *APIKey) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_api_key_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use APIKey.ProtoReflect.Descriptor instead.
func (*APIKey) Descriptor() ([]byte, []int) {
	return file_proto_account_api_key_proto_rawDescGZIP(), []int{0}
}

func (x *APIKey) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *APIKey) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *APIKey) GetRole() APIKey_Role {
	if x != nil {
		return x.Role
	}
	return APIKey_UNKNOWN
}

func (x *APIKey) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *APIKey) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *APIKey) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

type EnvironmentAPIKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Deprecated: Do not use.
	EnvironmentNamespace string  `protobuf:"bytes,1,opt,name=environment_namespace,json=environmentNamespace,proto3" json:"environment_namespace"`
	ApiKey               *APIKey `protobuf:"bytes,2,opt,name=api_key,json=apiKey,proto3" json:"api_key"`
	EnvironmentDisabled  bool    `protobuf:"varint,3,opt,name=environment_disabled,json=environmentDisabled,proto3" json:"environment_disabled"`
	// Deprecated: Do not use.
	ProjectId      string                     `protobuf:"bytes,4,opt,name=project_id,json=projectId,proto3" json:"project_id"`
	Environment    *environment.EnvironmentV2 `protobuf:"bytes,5,opt,name=environment,proto3" json:"environment"`
	ProjectUrlCode string                     `protobuf:"bytes,6,opt,name=project_url_code,json=projectUrlCode,proto3" json:"project_url_code"`
}

func (x *EnvironmentAPIKey) Reset() {
	*x = EnvironmentAPIKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_api_key_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnvironmentAPIKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvironmentAPIKey) ProtoMessage() {}

func (x *EnvironmentAPIKey) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_api_key_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnvironmentAPIKey.ProtoReflect.Descriptor instead.
func (*EnvironmentAPIKey) Descriptor() ([]byte, []int) {
	return file_proto_account_api_key_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Do not use.
func (x *EnvironmentAPIKey) GetEnvironmentNamespace() string {
	if x != nil {
		return x.EnvironmentNamespace
	}
	return ""
}

func (x *EnvironmentAPIKey) GetApiKey() *APIKey {
	if x != nil {
		return x.ApiKey
	}
	return nil
}

func (x *EnvironmentAPIKey) GetEnvironmentDisabled() bool {
	if x != nil {
		return x.EnvironmentDisabled
	}
	return false
}

// Deprecated: Do not use.
func (x *EnvironmentAPIKey) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *EnvironmentAPIKey) GetEnvironment() *environment.EnvironmentV2 {
	if x != nil {
		return x.Environment
	}
	return nil
}

func (x *EnvironmentAPIKey) GetProjectUrlCode() string {
	if x != nil {
		return x.ProjectUrlCode
	}
	return ""
}

var File_proto_account_api_key_proto protoreflect.FileDescriptor

var file_proto_account_api_key_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f,
	0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x1a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb5, 0x02, 0x0a, 0x06, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x2e, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61,
	0x62, 0x6c, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61,
	0x62, 0x6c, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x22, 0x79, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x44, 0x4b, 0x5f, 0x43,
	0x4c, 0x49, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x44, 0x4b, 0x5f, 0x53,
	0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x50, 0x55, 0x42, 0x4c, 0x49,
	0x43, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x52, 0x45, 0x41, 0x44, 0x5f, 0x4f, 0x4e, 0x4c, 0x59, 0x10,
	0x03, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x41, 0x50, 0x49, 0x5f,
	0x57, 0x52, 0x49, 0x54, 0x45, 0x10, 0x04, 0x12, 0x14, 0x0a, 0x10, 0x50, 0x55, 0x42, 0x4c, 0x49,
	0x43, 0x5f, 0x41, 0x50, 0x49, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10, 0x05, 0x22, 0xc8, 0x02,
	0x0a, 0x11, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x50, 0x49,
	0x4b, 0x65, 0x79, 0x12, 0x37, 0x0a, 0x15, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x14, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x07,
	0x61, 0x70, 0x69, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79,
	0x12, 0x31, 0x0a, 0x14, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x13,
	0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x69, 0x73, 0x61, 0x62,
	0x6c, 0x65, 0x64, 0x12, 0x21, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x09, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x46, 0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56,
	0x32, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x28,
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x55, 0x72, 0x6c, 0x43, 0x6f, 0x64, 0x65, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_account_api_key_proto_rawDescOnce sync.Once
	file_proto_account_api_key_proto_rawDescData = file_proto_account_api_key_proto_rawDesc
)

func file_proto_account_api_key_proto_rawDescGZIP() []byte {
	file_proto_account_api_key_proto_rawDescOnce.Do(func() {
		file_proto_account_api_key_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_account_api_key_proto_rawDescData)
	})
	return file_proto_account_api_key_proto_rawDescData
}

var file_proto_account_api_key_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_account_api_key_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_account_api_key_proto_goTypes = []interface{}{
	(APIKey_Role)(0),                  // 0: bucketeer.account.APIKey.Role
	(*APIKey)(nil),                    // 1: bucketeer.account.APIKey
	(*EnvironmentAPIKey)(nil),         // 2: bucketeer.account.EnvironmentAPIKey
	(*environment.EnvironmentV2)(nil), // 3: bucketeer.environment.EnvironmentV2
}
var file_proto_account_api_key_proto_depIdxs = []int32{
	0, // 0: bucketeer.account.APIKey.role:type_name -> bucketeer.account.APIKey.Role
	1, // 1: bucketeer.account.EnvironmentAPIKey.api_key:type_name -> bucketeer.account.APIKey
	3, // 2: bucketeer.account.EnvironmentAPIKey.environment:type_name -> bucketeer.environment.EnvironmentV2
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_account_api_key_proto_init() }
func file_proto_account_api_key_proto_init() {
	if File_proto_account_api_key_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_account_api_key_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*APIKey); i {
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
		file_proto_account_api_key_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnvironmentAPIKey); i {
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
			RawDescriptor: file_proto_account_api_key_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_account_api_key_proto_goTypes,
		DependencyIndexes: file_proto_account_api_key_proto_depIdxs,
		EnumInfos:         file_proto_account_api_key_proto_enumTypes,
		MessageInfos:      file_proto_account_api_key_proto_msgTypes,
	}.Build()
	File_proto_account_api_key_proto = out.File
	file_proto_account_api_key_proto_rawDesc = nil
	file_proto_account_api_key_proto_goTypes = nil
	file_proto_account_api_key_proto_depIdxs = nil
}
