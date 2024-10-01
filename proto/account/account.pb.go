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
// source: proto/account/account.proto

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

type Account_Role int32

const (
	Account_VIEWER     Account_Role = 0
	Account_EDITOR     Account_Role = 1
	Account_OWNER      Account_Role = 2
	Account_UNASSIGNED Account_Role = 99
)

// Enum value maps for Account_Role.
var (
	Account_Role_name = map[int32]string{
		0:  "VIEWER",
		1:  "EDITOR",
		2:  "OWNER",
		99: "UNASSIGNED",
	}
	Account_Role_value = map[string]int32{
		"VIEWER":     0,
		"EDITOR":     1,
		"OWNER":      2,
		"UNASSIGNED": 99,
	}
)

func (x Account_Role) Enum() *Account_Role {
	p := new(Account_Role)
	*p = x
	return p
}

func (x Account_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Account_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_account_proto_enumTypes[0].Descriptor()
}

func (Account_Role) Type() protoreflect.EnumType {
	return &file_proto_account_account_proto_enumTypes[0]
}

func (x Account_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Account_Role.Descriptor instead.
func (Account_Role) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{0, 0}
}

type AccountV2_Role_Environment int32

const (
	AccountV2_Role_Environment_UNASSIGNED AccountV2_Role_Environment = 0
	AccountV2_Role_Environment_VIEWER     AccountV2_Role_Environment = 1
	AccountV2_Role_Environment_EDITOR     AccountV2_Role_Environment = 2
)

// Enum value maps for AccountV2_Role_Environment.
var (
	AccountV2_Role_Environment_name = map[int32]string{
		0: "Environment_UNASSIGNED",
		1: "Environment_VIEWER",
		2: "Environment_EDITOR",
	}
	AccountV2_Role_Environment_value = map[string]int32{
		"Environment_UNASSIGNED": 0,
		"Environment_VIEWER":     1,
		"Environment_EDITOR":     2,
	}
)

func (x AccountV2_Role_Environment) Enum() *AccountV2_Role_Environment {
	p := new(AccountV2_Role_Environment)
	*p = x
	return p
}

func (x AccountV2_Role_Environment) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountV2_Role_Environment) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_account_proto_enumTypes[1].Descriptor()
}

func (AccountV2_Role_Environment) Type() protoreflect.EnumType {
	return &file_proto_account_account_proto_enumTypes[1]
}

func (x AccountV2_Role_Environment) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountV2_Role_Environment.Descriptor instead.
func (AccountV2_Role_Environment) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1, 0, 0}
}

type AccountV2_Role_Organization int32

const (
	AccountV2_Role_Organization_UNASSIGNED AccountV2_Role_Organization = 0
	AccountV2_Role_Organization_MEMBER     AccountV2_Role_Organization = 1
	AccountV2_Role_Organization_ADMIN      AccountV2_Role_Organization = 2
	AccountV2_Role_Organization_OWNER      AccountV2_Role_Organization = 3
)

// Enum value maps for AccountV2_Role_Organization.
var (
	AccountV2_Role_Organization_name = map[int32]string{
		0: "Organization_UNASSIGNED",
		1: "Organization_MEMBER",
		2: "Organization_ADMIN",
		3: "Organization_OWNER",
	}
	AccountV2_Role_Organization_value = map[string]int32{
		"Organization_UNASSIGNED": 0,
		"Organization_MEMBER":     1,
		"Organization_ADMIN":      2,
		"Organization_OWNER":      3,
	}
)

func (x AccountV2_Role_Organization) Enum() *AccountV2_Role_Organization {
	p := new(AccountV2_Role_Organization)
	*p = x
	return p
}

func (x AccountV2_Role_Organization) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AccountV2_Role_Organization) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_account_account_proto_enumTypes[2].Descriptor()
}

func (AccountV2_Role_Organization) Type() protoreflect.EnumType {
	return &file_proto_account_account_proto_enumTypes[2]
}

func (x AccountV2_Role_Organization) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AccountV2_Role_Organization.Descriptor instead.
func (AccountV2_Role_Organization) EnumDescriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1, 0, 1}
}

// Deprecated: Do not use.
type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Email     string       `protobuf:"bytes,2,opt,name=email,proto3" json:"email"`
	Name      string       `protobuf:"bytes,3,opt,name=name,proto3" json:"name"`
	Role      Account_Role `protobuf:"varint,4,opt,name=role,proto3,enum=bucketeer.account.Account_Role" json:"role"`
	Disabled  bool         `protobuf:"varint,5,opt,name=disabled,proto3" json:"disabled"`
	CreatedAt int64        `protobuf:"varint,6,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt int64        `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	Deleted   bool         `protobuf:"varint,8,opt,name=deleted,proto3" json:"deleted"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{0}
}

func (x *Account) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Account) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *Account) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Account) GetRole() Account_Role {
	if x != nil {
		return x.Role
	}
	return Account_VIEWER
}

func (x *Account) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Account) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Account) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Account) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type AccountV2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email            string                       `protobuf:"bytes,1,opt,name=email,proto3" json:"email"`
	AvatarImageUrl   string                       `protobuf:"bytes,3,opt,name=avatar_image_url,json=avatarImageUrl,proto3" json:"avatar_image_url"`
	OrganizationId   string                       `protobuf:"bytes,4,opt,name=organization_id,json=organizationId,proto3" json:"organization_id"`
	OrganizationRole AccountV2_Role_Organization  `protobuf:"varint,5,opt,name=organization_role,json=organizationRole,proto3,enum=bucketeer.account.AccountV2_Role_Organization" json:"organization_role"`
	EnvironmentRoles []*AccountV2_EnvironmentRole `protobuf:"bytes,6,rep,name=environment_roles,json=environmentRoles,proto3" json:"environment_roles"`
	Disabled         bool                         `protobuf:"varint,7,opt,name=disabled,proto3" json:"disabled"`
	CreatedAt        int64                        `protobuf:"varint,8,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt        int64                        `protobuf:"varint,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	SearchFilters    []*SearchFilter              `protobuf:"bytes,10,rep,name=search_filters,json=searchFilters,proto3" json:"search_filters"`
	FirstName        string                       `protobuf:"bytes,11,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName         string                       `protobuf:"bytes,12,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Language         string                       `protobuf:"bytes,13,opt,name=language,proto3" json:"language"`
}

func (x *AccountV2) Reset() {
	*x = AccountV2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountV2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountV2) ProtoMessage() {}

func (x *AccountV2) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountV2.ProtoReflect.Descriptor instead.
func (*AccountV2) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1}
}

func (x *AccountV2) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *AccountV2) GetAvatarImageUrl() string {
	if x != nil {
		return x.AvatarImageUrl
	}
	return ""
}

func (x *AccountV2) GetOrganizationId() string {
	if x != nil {
		return x.OrganizationId
	}
	return ""
}

func (x *AccountV2) GetOrganizationRole() AccountV2_Role_Organization {
	if x != nil {
		return x.OrganizationRole
	}
	return AccountV2_Role_Organization_UNASSIGNED
}

func (x *AccountV2) GetEnvironmentRoles() []*AccountV2_EnvironmentRole {
	if x != nil {
		return x.EnvironmentRoles
	}
	return nil
}

func (x *AccountV2) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *AccountV2) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *AccountV2) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *AccountV2) GetSearchFilters() []*SearchFilter {
	if x != nil {
		return x.SearchFilters
	}
	return nil
}

func (x *AccountV2) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *AccountV2) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *AccountV2) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

type ConsoleAccount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email            string                            `protobuf:"bytes,1,opt,name=email,proto3" json:"email"`
	AvatarUrl        string                            `protobuf:"bytes,3,opt,name=avatar_url,json=avatarUrl,proto3" json:"avatar_url"`
	IsSystemAdmin    bool                              `protobuf:"varint,4,opt,name=is_system_admin,json=isSystemAdmin,proto3" json:"is_system_admin"`
	Organization     *environment.Organization         `protobuf:"bytes,5,opt,name=organization,proto3" json:"organization"`
	OrganizationRole AccountV2_Role_Organization       `protobuf:"varint,6,opt,name=organization_role,json=organizationRole,proto3,enum=bucketeer.account.AccountV2_Role_Organization" json:"organization_role"`
	EnvironmentRoles []*ConsoleAccount_EnvironmentRole `protobuf:"bytes,7,rep,name=environment_roles,json=environmentRoles,proto3" json:"environment_roles"`
	SearchFilters    []*SearchFilter                   `protobuf:"bytes,8,rep,name=search_filters,json=searchFilters,proto3" json:"search_filters"`
	FirstName        string                            `protobuf:"bytes,9,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName         string                            `protobuf:"bytes,10,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Language         string                            `protobuf:"bytes,11,opt,name=language,proto3" json:"language"`
}

func (x *ConsoleAccount) Reset() {
	*x = ConsoleAccount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsoleAccount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsoleAccount) ProtoMessage() {}

func (x *ConsoleAccount) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsoleAccount.ProtoReflect.Descriptor instead.
func (*ConsoleAccount) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{2}
}

func (x *ConsoleAccount) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ConsoleAccount) GetAvatarUrl() string {
	if x != nil {
		return x.AvatarUrl
	}
	return ""
}

func (x *ConsoleAccount) GetIsSystemAdmin() bool {
	if x != nil {
		return x.IsSystemAdmin
	}
	return false
}

func (x *ConsoleAccount) GetOrganization() *environment.Organization {
	if x != nil {
		return x.Organization
	}
	return nil
}

func (x *ConsoleAccount) GetOrganizationRole() AccountV2_Role_Organization {
	if x != nil {
		return x.OrganizationRole
	}
	return AccountV2_Role_Organization_UNASSIGNED
}

func (x *ConsoleAccount) GetEnvironmentRoles() []*ConsoleAccount_EnvironmentRole {
	if x != nil {
		return x.EnvironmentRoles
	}
	return nil
}

func (x *ConsoleAccount) GetSearchFilters() []*SearchFilter {
	if x != nil {
		return x.SearchFilters
	}
	return nil
}

func (x *ConsoleAccount) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *ConsoleAccount) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *ConsoleAccount) GetLanguage() string {
	if x != nil {
		return x.Language
	}
	return ""
}

type AccountV2_Role struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AccountV2_Role) Reset() {
	*x = AccountV2_Role{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountV2_Role) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountV2_Role) ProtoMessage() {}

func (x *AccountV2_Role) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountV2_Role.ProtoReflect.Descriptor instead.
func (*AccountV2_Role) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1, 0}
}

type AccountV2_EnvironmentRole struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnvironmentId string                     `protobuf:"bytes,1,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
	Role          AccountV2_Role_Environment `protobuf:"varint,2,opt,name=role,proto3,enum=bucketeer.account.AccountV2_Role_Environment" json:"role"`
}

func (x *AccountV2_EnvironmentRole) Reset() {
	*x = AccountV2_EnvironmentRole{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountV2_EnvironmentRole) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountV2_EnvironmentRole) ProtoMessage() {}

func (x *AccountV2_EnvironmentRole) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountV2_EnvironmentRole.ProtoReflect.Descriptor instead.
func (*AccountV2_EnvironmentRole) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{1, 1}
}

func (x *AccountV2_EnvironmentRole) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

func (x *AccountV2_EnvironmentRole) GetRole() AccountV2_Role_Environment {
	if x != nil {
		return x.Role
	}
	return AccountV2_Role_Environment_UNASSIGNED
}

type ConsoleAccount_EnvironmentRole struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Environment *environment.EnvironmentV2 `protobuf:"bytes,1,opt,name=environment,proto3" json:"environment"`
	Project     *environment.Project       `protobuf:"bytes,2,opt,name=project,proto3" json:"project"`
	Role        AccountV2_Role_Environment `protobuf:"varint,3,opt,name=role,proto3,enum=bucketeer.account.AccountV2_Role_Environment" json:"role"`
}

func (x *ConsoleAccount_EnvironmentRole) Reset() {
	*x = ConsoleAccount_EnvironmentRole{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_account_account_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsoleAccount_EnvironmentRole) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsoleAccount_EnvironmentRole) ProtoMessage() {}

func (x *ConsoleAccount_EnvironmentRole) ProtoReflect() protoreflect.Message {
	mi := &file_proto_account_account_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsoleAccount_EnvironmentRole.ProtoReflect.Descriptor instead.
func (*ConsoleAccount_EnvironmentRole) Descriptor() ([]byte, []int) {
	return file_proto_account_account_proto_rawDescGZIP(), []int{2, 0}
}

func (x *ConsoleAccount_EnvironmentRole) GetEnvironment() *environment.EnvironmentV2 {
	if x != nil {
		return x.Environment
	}
	return nil
}

func (x *ConsoleAccount_EnvironmentRole) GetProject() *environment.Project {
	if x != nil {
		return x.Project
	}
	return nil
}

func (x *ConsoleAccount_EnvironmentRole) GetRole() AccountV2_Role_Environment {
	if x != nil {
		return x.Role
	}
	return AccountV2_Role_Environment_UNASSIGNED
}

var File_proto_account_account_proto protoreflect.FileDescriptor

var file_proto_account_account_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x1a, 0x23, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6e,
	0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69,
	0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x21, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xab, 0x02, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e,
	0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69,
	0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69,
	0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x39,
	0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x56, 0x49, 0x45, 0x57, 0x45, 0x52,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x45, 0x44, 0x49, 0x54, 0x4f, 0x52, 0x10, 0x01, 0x12, 0x09,
	0x0a, 0x05, 0x4f, 0x57, 0x4e, 0x45, 0x52, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x4e, 0x41,
	0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x63, 0x3a, 0x02, 0x18, 0x01, 0x22, 0x83, 0x07,
	0x0a, 0x09, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69,
	0x6c, 0x12, 0x28, 0x0a, 0x10, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x27, 0x0a, 0x0f, 0x6f,
	0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x64, 0x12, 0x5b, 0x0a, 0x11, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x2e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x52, 0x6f,
	0x6c, 0x65, 0x2e, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x10, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6c,
	0x65, 0x12, 0x59, 0x0a, 0x11, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2c, 0x2e, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x45, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x10, 0x65, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x46, 0x0a, 0x0e, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f,
	0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52,
	0x0d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61,
	0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x1a, 0xd7, 0x01, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65, 0x22,
	0x59, 0x0a, 0x0b, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1a,
	0x0a, 0x16, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x55, 0x4e,
	0x41, 0x53, 0x53, 0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x6e,
	0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x56, 0x49, 0x45, 0x57, 0x45, 0x52,
	0x10, 0x01, 0x12, 0x16, 0x0a, 0x12, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x5f, 0x45, 0x44, 0x49, 0x54, 0x4f, 0x52, 0x10, 0x02, 0x22, 0x74, 0x0a, 0x0c, 0x4f, 0x72,
	0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x17, 0x4f, 0x72,
	0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x55, 0x4e, 0x41, 0x53, 0x53,
	0x49, 0x47, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x4f, 0x72, 0x67, 0x61, 0x6e,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x01,
	0x12, 0x16, 0x0a, 0x12, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x4f, 0x72, 0x67, 0x61,
	0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x4f, 0x57, 0x4e, 0x45, 0x52, 0x10, 0x03,
	0x1a, 0x7b, 0x0a, 0x0f, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x6f, 0x6c, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x2e, 0x45, 0x6e, 0x76, 0x69,
	0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x4a, 0x04, 0x08,
	0x02, 0x10, 0x03, 0x22, 0xf2, 0x05, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1d, 0x0a, 0x0a,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x55, 0x72, 0x6c, 0x12, 0x26, 0x0a, 0x0f, 0x69,
	0x73, 0x5f, 0x73, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x5f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x12, 0x47, 0x0a, 0x0c, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x4f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c,
	0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5b, 0x0a, 0x11,
	0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x6f, 0x6c,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x2e, 0x4f, 0x72, 0x67, 0x61, 0x6e,
	0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x10, 0x6f, 0x72, 0x67, 0x61, 0x6e, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x5e, 0x0a, 0x11, 0x65, 0x6e, 0x76,
	0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x73, 0x18, 0x07,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65,
	0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x10, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x73, 0x12, 0x46, 0x0a, 0x0e, 0x73, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x52, 0x0d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x1a, 0xd6, 0x01, 0x0a, 0x0f, 0x45, 0x6e,
	0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x12, 0x46, 0x0a,
	0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x24, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x56, 0x32, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x38, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65,
	0x65, 0x72, 0x2e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12,
	0x41, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2d, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x32, 0x2e, 0x52, 0x6f, 0x6c, 0x65,
	0x2e, 0x45, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_account_account_proto_rawDescOnce sync.Once
	file_proto_account_account_proto_rawDescData = file_proto_account_account_proto_rawDesc
)

func file_proto_account_account_proto_rawDescGZIP() []byte {
	file_proto_account_account_proto_rawDescOnce.Do(func() {
		file_proto_account_account_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_account_account_proto_rawDescData)
	})
	return file_proto_account_account_proto_rawDescData
}

var file_proto_account_account_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_account_account_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_account_account_proto_goTypes = []interface{}{
	(Account_Role)(0),                      // 0: bucketeer.account.Account.Role
	(AccountV2_Role_Environment)(0),        // 1: bucketeer.account.AccountV2.Role.Environment
	(AccountV2_Role_Organization)(0),       // 2: bucketeer.account.AccountV2.Role.Organization
	(*Account)(nil),                        // 3: bucketeer.account.Account
	(*AccountV2)(nil),                      // 4: bucketeer.account.AccountV2
	(*ConsoleAccount)(nil),                 // 5: bucketeer.account.ConsoleAccount
	(*AccountV2_Role)(nil),                 // 6: bucketeer.account.AccountV2.Role
	(*AccountV2_EnvironmentRole)(nil),      // 7: bucketeer.account.AccountV2.EnvironmentRole
	(*ConsoleAccount_EnvironmentRole)(nil), // 8: bucketeer.account.ConsoleAccount.EnvironmentRole
	(*SearchFilter)(nil),                   // 9: bucketeer.account.SearchFilter
	(*environment.Organization)(nil),       // 10: bucketeer.environment.Organization
	(*environment.EnvironmentV2)(nil),      // 11: bucketeer.environment.EnvironmentV2
	(*environment.Project)(nil),            // 12: bucketeer.environment.Project
}
var file_proto_account_account_proto_depIdxs = []int32{
	0,  // 0: bucketeer.account.Account.role:type_name -> bucketeer.account.Account.Role
	2,  // 1: bucketeer.account.AccountV2.organization_role:type_name -> bucketeer.account.AccountV2.Role.Organization
	7,  // 2: bucketeer.account.AccountV2.environment_roles:type_name -> bucketeer.account.AccountV2.EnvironmentRole
	9,  // 3: bucketeer.account.AccountV2.search_filters:type_name -> bucketeer.account.SearchFilter
	10, // 4: bucketeer.account.ConsoleAccount.organization:type_name -> bucketeer.environment.Organization
	2,  // 5: bucketeer.account.ConsoleAccount.organization_role:type_name -> bucketeer.account.AccountV2.Role.Organization
	8,  // 6: bucketeer.account.ConsoleAccount.environment_roles:type_name -> bucketeer.account.ConsoleAccount.EnvironmentRole
	9,  // 7: bucketeer.account.ConsoleAccount.search_filters:type_name -> bucketeer.account.SearchFilter
	1,  // 8: bucketeer.account.AccountV2.EnvironmentRole.role:type_name -> bucketeer.account.AccountV2.Role.Environment
	11, // 9: bucketeer.account.ConsoleAccount.EnvironmentRole.environment:type_name -> bucketeer.environment.EnvironmentV2
	12, // 10: bucketeer.account.ConsoleAccount.EnvironmentRole.project:type_name -> bucketeer.environment.Project
	1,  // 11: bucketeer.account.ConsoleAccount.EnvironmentRole.role:type_name -> bucketeer.account.AccountV2.Role.Environment
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_proto_account_account_proto_init() }
func file_proto_account_account_proto_init() {
	if File_proto_account_account_proto != nil {
		return
	}
	file_proto_account_search_filter_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_account_account_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
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
		file_proto_account_account_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountV2); i {
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
		file_proto_account_account_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsoleAccount); i {
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
		file_proto_account_account_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountV2_Role); i {
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
		file_proto_account_account_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountV2_EnvironmentRole); i {
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
		file_proto_account_account_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConsoleAccount_EnvironmentRole); i {
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
			RawDescriptor: file_proto_account_account_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_account_account_proto_goTypes,
		DependencyIndexes: file_proto_account_account_proto_depIdxs,
		EnumInfos:         file_proto_account_account_proto_enumTypes,
		MessageInfos:      file_proto_account_account_proto_msgTypes,
	}.Build()
	File_proto_account_account_proto = out.File
	file_proto_account_account_proto_rawDesc = nil
	file_proto_account_account_proto_goTypes = nil
	file_proto_account_account_proto_depIdxs = nil
}
