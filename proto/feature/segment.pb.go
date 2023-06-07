// Copyright 2022 The Bucketeer Authors.
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
// 	protoc        v3.20.3
// source: proto/feature/segment.proto

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

type Segment_Status int32

const (
	Segment_INITIAL   Segment_Status = 0
	Segment_UPLOADING Segment_Status = 1
	Segment_SUCEEDED  Segment_Status = 2
	Segment_FAILED    Segment_Status = 3
)

// Enum value maps for Segment_Status.
var (
	Segment_Status_name = map[int32]string{
		0: "INITIAL",
		1: "UPLOADING",
		2: "SUCEEDED",
		3: "FAILED",
	}
	Segment_Status_value = map[string]int32{
		"INITIAL":   0,
		"UPLOADING": 1,
		"SUCEEDED":  2,
		"FAILED":    3,
	}
)

func (x Segment_Status) Enum() *Segment_Status {
	p := new(Segment_Status)
	*p = x
	return p
}

func (x Segment_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Segment_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_feature_segment_proto_enumTypes[0].Descriptor()
}

func (Segment_Status) Type() protoreflect.EnumType {
	return &file_proto_feature_segment_proto_enumTypes[0]
}

func (x Segment_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Segment_Status.Descriptor instead.
func (Segment_Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_feature_segment_proto_rawDescGZIP(), []int{0, 0}
}

type SegmentUser_State int32

const (
	SegmentUser_INCLUDED SegmentUser_State = 0
	// Deprecated: Do not use.
	SegmentUser_EXCLUDED SegmentUser_State = 1
)

// Enum value maps for SegmentUser_State.
var (
	SegmentUser_State_name = map[int32]string{
		0: "INCLUDED",
		1: "EXCLUDED",
	}
	SegmentUser_State_value = map[string]int32{
		"INCLUDED": 0,
		"EXCLUDED": 1,
	}
)

func (x SegmentUser_State) Enum() *SegmentUser_State {
	p := new(SegmentUser_State)
	*p = x
	return p
}

func (x SegmentUser_State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SegmentUser_State) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_feature_segment_proto_enumTypes[1].Descriptor()
}

func (SegmentUser_State) Type() protoreflect.EnumType {
	return &file_proto_feature_segment_proto_enumTypes[1]
}

func (x SegmentUser_State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SegmentUser_State.Descriptor instead.
func (SegmentUser_State) EnumDescriptor() ([]byte, []int) {
	return file_proto_feature_segment_proto_rawDescGZIP(), []int{1, 0}
}

type Segment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Rules       []*Rule `protobuf:"bytes,4,rep,name=rules,proto3" json:"rules,omitempty"`
	CreatedAt   int64   `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   int64   `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	// Deprecated: Do not use.
	Version           int64 `protobuf:"varint,7,opt,name=version,proto3" json:"version,omitempty"`
	Deleted           bool  `protobuf:"varint,8,opt,name=deleted,proto3" json:"deleted,omitempty"`
	IncludedUserCount int64 `protobuf:"varint,9,opt,name=included_user_count,json=includedUserCount,proto3" json:"included_user_count,omitempty"`
	// Deprecated: Do not use.
	ExcludedUserCount int64          `protobuf:"varint,10,opt,name=excluded_user_count,json=excludedUserCount,proto3" json:"excluded_user_count,omitempty"`
	Status            Segment_Status `protobuf:"varint,11,opt,name=status,proto3,enum=bucketeer.feature.Segment_Status" json:"status,omitempty"`
	IsInUseStatus     bool           `protobuf:"varint,12,opt,name=is_in_use_status,json=isInUseStatus,proto3" json:"is_in_use_status,omitempty"` // This field is set only when APIs return.
}

func (x *Segment) Reset() {
	*x = Segment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_segment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Segment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Segment) ProtoMessage() {}

func (x *Segment) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_segment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Segment.ProtoReflect.Descriptor instead.
func (*Segment) Descriptor() ([]byte, []int) {
	return file_proto_feature_segment_proto_rawDescGZIP(), []int{0}
}

func (x *Segment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Segment) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Segment) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Segment) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

func (x *Segment) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Segment) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

// Deprecated: Do not use.
func (x *Segment) GetVersion() int64 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Segment) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Segment) GetIncludedUserCount() int64 {
	if x != nil {
		return x.IncludedUserCount
	}
	return 0
}

// Deprecated: Do not use.
func (x *Segment) GetExcludedUserCount() int64 {
	if x != nil {
		return x.ExcludedUserCount
	}
	return 0
}

func (x *Segment) GetStatus() Segment_Status {
	if x != nil {
		return x.Status
	}
	return Segment_INITIAL
}

func (x *Segment) GetIsInUseStatus() bool {
	if x != nil {
		return x.IsInUseStatus
	}
	return false
}

type SegmentUser struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SegmentId string            `protobuf:"bytes,2,opt,name=segment_id,json=segmentId,proto3" json:"segment_id,omitempty"`
	UserId    string            `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	State     SegmentUser_State `protobuf:"varint,4,opt,name=state,proto3,enum=bucketeer.feature.SegmentUser_State" json:"state,omitempty"`
	Deleted   bool              `protobuf:"varint,5,opt,name=deleted,proto3" json:"deleted,omitempty"`
}

func (x *SegmentUser) Reset() {
	*x = SegmentUser{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_segment_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentUser) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentUser) ProtoMessage() {}

func (x *SegmentUser) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_segment_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentUser.ProtoReflect.Descriptor instead.
func (*SegmentUser) Descriptor() ([]byte, []int) {
	return file_proto_feature_segment_proto_rawDescGZIP(), []int{1}
}

func (x *SegmentUser) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SegmentUser) GetSegmentId() string {
	if x != nil {
		return x.SegmentId
	}
	return ""
}

func (x *SegmentUser) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SegmentUser) GetState() SegmentUser_State {
	if x != nil {
		return x.State
	}
	return SegmentUser_INCLUDED
}

func (x *SegmentUser) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

type SegmentUsers struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SegmentId string         `protobuf:"bytes,1,opt,name=segment_id,json=segmentId,proto3" json:"segment_id,omitempty"`
	Users     []*SegmentUser `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *SegmentUsers) Reset() {
	*x = SegmentUsers{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_segment_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SegmentUsers) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SegmentUsers) ProtoMessage() {}

func (x *SegmentUsers) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_segment_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SegmentUsers.ProtoReflect.Descriptor instead.
func (*SegmentUsers) Descriptor() ([]byte, []int) {
	return file_proto_feature_segment_proto_rawDescGZIP(), []int{2}
}

func (x *SegmentUsers) GetSegmentId() string {
	if x != nil {
		return x.SegmentId
	}
	return ""
}

func (x *SegmentUsers) GetUsers() []*SegmentUser {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_proto_feature_segment_proto protoreflect.FileDescriptor

var file_proto_feature_segment_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x1a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfc, 0x03, 0x0a, 0x07, 0x53,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x2d, 0x0a, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02, 0x18, 0x01, 0x52, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x12, 0x2e, 0x0a, 0x13, 0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11,
	0x69, 0x6e, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x32, 0x0a, 0x13, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x03, 0x42, 0x02,
	0x18, 0x01, 0x52, 0x11, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x39, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x27, 0x0a, 0x10, 0x69, 0x73, 0x5f, 0x69, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x49, 0x6e,
	0x55, 0x73, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3e, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x0b, 0x0a, 0x07, 0x49, 0x4e, 0x49, 0x54, 0x49, 0x41, 0x4c, 0x10, 0x00,
	0x12, 0x0d, 0x0a, 0x09, 0x55, 0x50, 0x4c, 0x4f, 0x41, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x0c, 0x0a, 0x08, 0x53, 0x55, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x22, 0xd4, 0x01, 0x0a, 0x0b, 0x53, 0x65,
	0x67, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x67,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x3a, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x24, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x22, 0x27, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x0c, 0x0a, 0x08, 0x49, 0x4e, 0x43, 0x4c, 0x55, 0x44, 0x45, 0x44, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x08, 0x45, 0x58, 0x43, 0x4c, 0x55, 0x44, 0x45, 0x44, 0x10, 0x01, 0x1a, 0x02, 0x08, 0x01,
	0x22, 0x63, 0x0a, 0x0c, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x34, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x53, 0x65, 0x67, 0x6d, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x05,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f,
	0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_feature_segment_proto_rawDescOnce sync.Once
	file_proto_feature_segment_proto_rawDescData = file_proto_feature_segment_proto_rawDesc
)

func file_proto_feature_segment_proto_rawDescGZIP() []byte {
	file_proto_feature_segment_proto_rawDescOnce.Do(func() {
		file_proto_feature_segment_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_feature_segment_proto_rawDescData)
	})
	return file_proto_feature_segment_proto_rawDescData
}

var file_proto_feature_segment_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_feature_segment_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_feature_segment_proto_goTypes = []interface{}{
	(Segment_Status)(0),    // 0: bucketeer.feature.Segment.Status
	(SegmentUser_State)(0), // 1: bucketeer.feature.SegmentUser.State
	(*Segment)(nil),        // 2: bucketeer.feature.Segment
	(*SegmentUser)(nil),    // 3: bucketeer.feature.SegmentUser
	(*SegmentUsers)(nil),   // 4: bucketeer.feature.SegmentUsers
	(*Rule)(nil),           // 5: bucketeer.feature.Rule
}
var file_proto_feature_segment_proto_depIdxs = []int32{
	5, // 0: bucketeer.feature.Segment.rules:type_name -> bucketeer.feature.Rule
	0, // 1: bucketeer.feature.Segment.status:type_name -> bucketeer.feature.Segment.Status
	1, // 2: bucketeer.feature.SegmentUser.state:type_name -> bucketeer.feature.SegmentUser.State
	3, // 3: bucketeer.feature.SegmentUsers.users:type_name -> bucketeer.feature.SegmentUser
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_feature_segment_proto_init() }
func file_proto_feature_segment_proto_init() {
	if File_proto_feature_segment_proto != nil {
		return
	}
	file_proto_feature_rule_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_feature_segment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Segment); i {
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
		file_proto_feature_segment_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentUser); i {
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
		file_proto_feature_segment_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SegmentUsers); i {
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
			RawDescriptor: file_proto_feature_segment_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_feature_segment_proto_goTypes,
		DependencyIndexes: file_proto_feature_segment_proto_depIdxs,
		EnumInfos:         file_proto_feature_segment_proto_enumTypes,
		MessageInfos:      file_proto_feature_segment_proto_msgTypes,
	}.Build()
	File_proto_feature_segment_proto = out.File
	file_proto_feature_segment_proto_rawDesc = nil
	file_proto_feature_segment_proto_goTypes = nil
	file_proto_feature_segment_proto_depIdxs = nil
}
