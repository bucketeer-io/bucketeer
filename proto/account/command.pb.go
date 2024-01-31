// Copyright 2023 The Bucketeer Authors.
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
// source: proto/account/command.proto

package account

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

type ChangeAccountV2EnvironmentRolesCommand_WriteType int32

const (
	ChangeAccountV2EnvironmentRolesCommand_WriteType_UNSPECIFIED ChangeAccountV2EnvironmentRolesCommand_WriteType = 0
	ChangeAccountV2EnvironmentRolesCommand_WriteType_OVERRIDE    ChangeAccountV2EnvironmentRolesCommand_WriteType = 1
	ChangeAccountV2EnvironmentRolesCommand_WriteType_PATCH       ChangeAccountV2EnvironmentRolesCommand_WriteType = 2
)

// Enum value maps for ChangeAccountV2EnvironmentRolesCommand_WriteType.
var (
	ChangeAccountV2EnvironmentRolesCommand_WriteType_name = map[int32]string{
		0: "WriteType_UNSPECIFIED",
		1: "WriteType_OVERRIDE",
		2: "WriteType_PATCH",
	}
	ChangeAccountV2EnvironmentRolesCommand_WriteType_value = map[string]int32{
		"WriteType_UNSPECIFIED": 0,
		"WriteType_OVERRIDE":    1,
		"WriteType_PATCH":       2,
	}
)

func (x ChangeAccountV2EnvironmentRolesCommand_WriteType) Enum() *ChangeAccountV2EnvironmentRolesCommand_WriteType {
	p := new(ChangeAccountV2EnvironmentRolesCommand_WriteType)
	*p = x
	return p
}

func (x ChangeAccountV2EnvironmentRolesCommand_WriteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ChangeAccountV2EnvironmentRolesCommand_WriteType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_command_proto_enumTypes[0].Descriptor()
}

func (ChangeAccountV2EnvironmentRolesCommand_WriteType) Type() protoreflect.EnumType {
	return &file_proto_account_command_proto_enumTypes[0]
}

func (x ChangeAccountV2EnvironmentRolesCommand_WriteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ChangeAccountV2EnvironmentRolesCommand_WriteType.Descriptor instead.
func (ChangeAccountV2EnvironmentRolesCommand_WriteType) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{4, 0}
}

type CreateAccountV2Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email            string                       `protobuf:"bytes,1,opt,name=email,proto3" json:"email"`
	Name             string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	AvatarImageUrl   string                       `protobuf:"bytes,3,opt,name=avatar_image_url,json=avatarImageUrl,proto3" json:"avatar_image_url"`
	OrganizationRole AccountV2_Role_Organization  `protobuf:"varint,4,opt,name=organization_role,json=organizationRole,proto3,enum=bucketeer.account.AccountV2_Role_Organization" json:"organization_role"`
	EnvironmentRoles []*AccountV2_EnvironmentRole `protobuf:"bytes,5,rep,name=environment_roles,json=environmentRoles,proto3" json:"environment_roles"`
}

func (x *CreateAccountV2Command) Reset() {
	*x = CreateAccountV2Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountV2Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountV2Command) ProtoMessage() {}

func (x *CreateAccountV2Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountV2Command.ProtoReflect.Descriptor instead.
func (*CreateAccountV2Command) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{0}
}

func (x *CreateAccountV2Command) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateAccountV2Command) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAccountV2Command) GetAvatarImageUrl() string {
	if x != nil {
		return x.AvatarImageUrl
	}
	return ""
}

func (x *CreateAccountV2Command) GetOrganizationRole() AccountV2_Role_Organization {
	if x != nil {
		return x.OrganizationRole
	}
	return AccountV2_Role_Organization_UNASSIGNED
}

func (x *CreateAccountV2Command) GetEnvironmentRoles() []*AccountV2_EnvironmentRole {
	if x != nil {
		return x.EnvironmentRoles
	}
	return nil
}

type ChangeAccountV2NameCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
}

func (x *ChangeAccountV2NameCommand) Reset() {
	*x = ChangeAccountV2NameCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccountV2NameCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccountV2NameCommand) ProtoMessage() {}

func (x *ChangeAccountV2NameCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccountV2NameCommand.ProtoReflect.Descriptor instead.
func (*ChangeAccountV2NameCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{1}
}

func (x *ChangeAccountV2NameCommand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ChangeAccountV2AvatarImageUrlCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AvatarImageUrl string `protobuf:"bytes,1,opt,name=avatar_image_url,json=avatarImageUrl,proto3" json:"avatar_image_url"`
}

func (x *ChangeAccountV2AvatarImageUrlCommand) Reset() {
	*x = ChangeAccountV2AvatarImageUrlCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccountV2AvatarImageUrlCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccountV2AvatarImageUrlCommand) ProtoMessage() {}

func (x *ChangeAccountV2AvatarImageUrlCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccountV2AvatarImageUrlCommand.ProtoReflect.Descriptor instead.
func (*ChangeAccountV2AvatarImageUrlCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{2}
}

func (x *ChangeAccountV2AvatarImageUrlCommand) GetAvatarImageUrl() string {
	if x != nil {
		return x.AvatarImageUrl
	}
	return ""
}

type ChangeAccountV2OrganizationRoleCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role AccountV2_Role_Organization `protobuf:"varint,1,opt,name=role,proto3,enum=bucketeer.account.AccountV2_Role_Organization" json:"role"`
}

func (x *ChangeAccountV2OrganizationRoleCommand) Reset() {
	*x = ChangeAccountV2OrganizationRoleCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccountV2OrganizationRoleCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccountV2OrganizationRoleCommand) ProtoMessage() {}

func (x *ChangeAccountV2OrganizationRoleCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccountV2OrganizationRoleCommand.ProtoReflect.Descriptor instead.
func (*ChangeAccountV2OrganizationRoleCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{3}
}

func (x *ChangeAccountV2OrganizationRoleCommand) GetRole() AccountV2_Role_Organization {
	if x != nil {
		return x.Role
	}
	return AccountV2_Role_Organization_UNASSIGNED
}

type ChangeAccountV2EnvironmentRolesCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Roles     []*AccountV2_EnvironmentRole                     `protobuf:"bytes,1,rep,name=roles,proto3" json:"roles"`
	WriteType ChangeAccountV2EnvironmentRolesCommand_WriteType `protobuf:"varint,2,opt,name=write_type,json=writeType,proto3,enum=bucketeer.account.ChangeAccountV2EnvironmentRolesCommand_WriteType" json:"write_type"`
}

func (x *ChangeAccountV2EnvironmentRolesCommand) Reset() {
	*x = ChangeAccountV2EnvironmentRolesCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAccountV2EnvironmentRolesCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAccountV2EnvironmentRolesCommand) ProtoMessage() {}

func (x *ChangeAccountV2EnvironmentRolesCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAccountV2EnvironmentRolesCommand.ProtoReflect.Descriptor instead.
func (*ChangeAccountV2EnvironmentRolesCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{4}
}

func (x *ChangeAccountV2EnvironmentRolesCommand) GetRoles() []*AccountV2_EnvironmentRole {
	if x != nil {
		return x.Roles
	}
	return nil
}

func (x *ChangeAccountV2EnvironmentRolesCommand) GetWriteType() ChangeAccountV2EnvironmentRolesCommand_WriteType {
	if x != nil {
		return x.WriteType
	}
	return ChangeAccountV2EnvironmentRolesCommand_WriteType_UNSPECIFIED
}

type EnableAccountV2Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EnableAccountV2Command) Reset() {
	*x = EnableAccountV2Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnableAccountV2Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnableAccountV2Command) ProtoMessage() {}

func (x *EnableAccountV2Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnableAccountV2Command.ProtoReflect.Descriptor instead.
func (*EnableAccountV2Command) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{5}
}

type DisableAccountV2Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisableAccountV2Command) Reset() {
	*x = DisableAccountV2Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisableAccountV2Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisableAccountV2Command) ProtoMessage() {}

func (x *DisableAccountV2Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisableAccountV2Command.ProtoReflect.Descriptor instead.
func (*DisableAccountV2Command) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{6}
}

type DeleteAccountV2Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteAccountV2Command) Reset() {
	*x = DeleteAccountV2Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAccountV2Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAccountV2Command) ProtoMessage() {}

func (x *DeleteAccountV2Command) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAccountV2Command.ProtoReflect.Descriptor instead.
func (*DeleteAccountV2Command) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{7}
}

type CreateAPIKeyCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
	Role APIKey_Role `protobuf:"varint,2,opt,name=role,proto3,enum=bucketeer.account.APIKey_Role" json:"role"`
}

func (x *CreateAPIKeyCommand) Reset() {
	*x = CreateAPIKeyCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAPIKeyCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAPIKeyCommand) ProtoMessage() {}

func (x *CreateAPIKeyCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAPIKeyCommand.ProtoReflect.Descriptor instead.
func (*CreateAPIKeyCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{8}
}

func (x *CreateAPIKeyCommand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateAPIKeyCommand) GetRole() APIKey_Role {
	if x != nil {
		return x.Role
	}
	return APIKey_SDK
}

type ChangeAPIKeyNameCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
}

func (x *ChangeAPIKeyNameCommand) Reset() {
	*x = ChangeAPIKeyNameCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeAPIKeyNameCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeAPIKeyNameCommand) ProtoMessage() {}

func (x *ChangeAPIKeyNameCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeAPIKeyNameCommand.ProtoReflect.Descriptor instead.
func (*ChangeAPIKeyNameCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{9}
}

func (x *ChangeAPIKeyNameCommand) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type EnableAPIKeyCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EnableAPIKeyCommand) Reset() {
	*x = EnableAPIKeyCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EnableAPIKeyCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnableAPIKeyCommand) ProtoMessage() {}

func (x *EnableAPIKeyCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EnableAPIKeyCommand.ProtoReflect.Descriptor instead.
func (*EnableAPIKeyCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{10}
}

type DisableAPIKeyCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisableAPIKeyCommand) Reset() {
	*x = DisableAPIKeyCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_command_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisableAPIKeyCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisableAPIKeyCommand) ProtoMessage() {}

func (x *DisableAPIKeyCommand) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_command_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisableAPIKeyCommand.ProtoReflect.Descriptor instead.
func (*DisableAPIKeyCommand) Descriptor() ([]byte, []int) {
	return file_proto_account_command_proto_rawDescGZIP(), []int{11}
}

var File_proto_account_command_proto protoreflect.FileDescriptor

var file_proto_account_command_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x1a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x61, 0x70, 0x69,
	0x5f, 0x6b, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa4, 0x02, 0x0a, 0x16, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x28, 0x0a, 0x10, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x5b, 0x0a, 0x11, 0x6f, 0x72, 0x67,
	0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x56, 0x32, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x2e, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x59, 0x0a, 0x11, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2c, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e,
	0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x10, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65,
	0x73, 0x22, 0x30, 0x0a, 0x1a, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x56, 0x32, 0x4e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x22, 0x50, 0x0a, 0x24, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x55, 0x72, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x55, 0x72, 0x6c, 0x22, 0x6c, 0x0a, 0x26, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12,
	0x42, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x52, 0x6f, 0x6c, 0x65,
	0x2e, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x04, 0x72,
	0x6f, 0x6c, 0x65, 0x22, 0xa5, 0x02, 0x0a, 0x26, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x42,
	0x0a, 0x05, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x45, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x6f, 0x6c,
	0x65, 0x73, 0x12, 0x62, 0x0a, 0x0a, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x43, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65,
	0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x43, 0x68, 0x61, 0x6e, 0x67,
	0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x2e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x77, 0x72, 0x69,
	0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x53, 0x0a, 0x09, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16,
	0x0a, 0x12, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x5f, 0x4f, 0x56, 0x45, 0x52,
	0x52, 0x49, 0x44, 0x45, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x57, 0x72, 0x69, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x5f, 0x50, 0x41, 0x54, 0x43, 0x48, 0x10, 0x02, 0x22, 0x18, 0x0a, 0x16, 0x45,
	0x6e, 0x61, 0x62, 0x6c, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x19, 0x0a, 0x17, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x22, 0x18, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x56, 0x32, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x5d, 0x0a, 0x13, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x2e, 0x52,
	0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x2d, 0x0a, 0x17, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x15, 0x0a, 0x13, 0x45, 0x6e, 0x61, 0x62,
	0x6c, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22,
	0x16, 0x0a, 0x14, 0x44, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x41, 0x50, 0x49, 0x4b, 0x65, 0x79,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d,
	0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_account_command_proto_rawDescOnce sync.Once
	file_proto_account_command_proto_rawDescData = file_proto_account_command_proto_rawDesc
)

func file_proto_account_command_proto_rawDescGZIP() []byte {
	file_proto_account_command_proto_rawDescOnce.Do(func() {
		file_proto_account_command_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_account_command_proto_rawDescData)
	})
	return file_proto_account_command_proto_rawDescData
}

var file_proto_account_command_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_account_command_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_proto_account_command_proto_goTypes = []interface{}{
	(ChangeAccountV2EnvironmentRolesCommand_WriteType)(0), // 0: bucketeer.account.ChangeAccountV2EnvironmentRolesCommand.WriteType
	(*CreateAccountV2Command)(nil),                        // 1: bucketeer.account.CreateAccountV2Command
	(*ChangeAccountV2NameCommand)(nil),                    // 2: bucketeer.account.ChangeAccountV2NameCommand
	(*ChangeAccountV2AvatarImageUrlCommand)(nil),          // 3: bucketeer.account.ChangeAccountV2AvatarImageUrlCommand
	(*ChangeAccountV2OrganizationRoleCommand)(nil),        // 4: bucketeer.account.ChangeAccountV2OrganizationRoleCommand
	(*ChangeAccountV2EnvironmentRolesCommand)(nil),        // 5: bucketeer.account.ChangeAccountV2EnvironmentRolesCommand
	(*EnableAccountV2Command)(nil),                        // 6: bucketeer.account.EnableAccountV2Command
	(*DisableAccountV2Command)(nil),                       // 7: bucketeer.account.DisableAccountV2Command
	(*DeleteAccountV2Command)(nil),                        // 8: bucketeer.account.DeleteAccountV2Command
	(*CreateAPIKeyCommand)(nil),                           // 9: bucketeer.account.CreateAPIKeyCommand
	(*ChangeAPIKeyNameCommand)(nil),                       // 10: bucketeer.account.ChangeAPIKeyNameCommand
	(*EnableAPIKeyCommand)(nil),                           // 11: bucketeer.account.EnableAPIKeyCommand
	(*DisableAPIKeyCommand)(nil),                          // 12: bucketeer.account.DisableAPIKeyCommand
	(AccountV2_Role_Organization)(0),                      // 13: bucketeer.account.AccountV2.Role.Organization
	(*AccountV2_EnvironmentRole)(nil),                     // 14: bucketeer.account.AccountV2.EnvironmentRole
	(APIKey_Role)(0),                                      // 15: bucketeer.account.APIKey.Role
}
var file_proto_account_command_proto_depIdxs = []int32{
	13, // 0: bucketeer.account.CreateAccountV2Command.organization_role:type_name -> bucketeer.account.AccountV2.Role.Organization
	14, // 1: bucketeer.account.CreateAccountV2Command.environment_roles:type_name -> bucketeer.account.AccountV2.EnvironmentRole
	13, // 2: bucketeer.account.ChangeAccountV2OrganizationRoleCommand.role:type_name -> bucketeer.account.AccountV2.Role.Organization
	14, // 3: bucketeer.account.ChangeAccountV2EnvironmentRolesCommand.roles:type_name -> bucketeer.account.AccountV2.EnvironmentRole
	0,  // 4: bucketeer.account.ChangeAccountV2EnvironmentRolesCommand.write_type:type_name -> bucketeer.account.ChangeAccountV2EnvironmentRolesCommand.WriteType
	15, // 5: bucketeer.account.CreateAPIKeyCommand.role:type_name -> bucketeer.account.APIKey.Role
	6,  // [6:6] is the sub-list for method output_type
	6,  // [6:6] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_proto_account_command_proto_init() }
func file_proto_account_command_proto_init() {
	if File_proto_account_command_proto != nil {
		return
	}
	file_proto_account_account_proto_init()
	file_proto_account_api_key_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_account_command_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAccountV2Command); i {
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
		file_proto_account_command_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccountV2NameCommand); i {
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
		file_proto_account_command_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccountV2AvatarImageUrlCommand); i {
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
		file_proto_account_command_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccountV2OrganizationRoleCommand); i {
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
		file_proto_account_command_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAccountV2EnvironmentRolesCommand); i {
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
		file_proto_account_command_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnableAccountV2Command); i {
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
		file_proto_account_command_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisableAccountV2Command); i {
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
		file_proto_account_command_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAccountV2Command); i {
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
		file_proto_account_command_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAPIKeyCommand); i {
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
		file_proto_account_command_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeAPIKeyNameCommand); i {
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
		file_proto_account_command_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EnableAPIKeyCommand); i {
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
		file_proto_account_command_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisableAPIKeyCommand); i {
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
			RawDescriptor: file_proto_account_command_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_account_command_proto_goTypes,
		DependencyIndexes: file_proto_account_command_proto_depIdxs,
		EnumInfos:         file_proto_account_command_proto_enumTypes,
		MessageInfos:      file_proto_account_command_proto_msgTypes,
	}.Build()
	File_proto_account_command_proto = out.File
	file_proto_account_command_proto_rawDesc = nil
	file_proto_account_command_proto_goTypes = nil
	file_proto_account_command_proto_depIdxs = nil
}
