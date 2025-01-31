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
// source: proto/experiment/goal.proto

package experiment

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

type Goal_ConnectionType int32

const (
	Goal_UNKNOWN    Goal_ConnectionType = 0
	Goal_EXPERIMENT Goal_ConnectionType = 1
	Goal_OPERATION  Goal_ConnectionType = 2
)

// Enum value maps for Goal_ConnectionType.
var (
	Goal_ConnectionType_name = map[int32]string{
		0: "UNKNOWN",
		1: "EXPERIMENT",
		2: "OPERATION",
	}
	Goal_ConnectionType_value = map[string]int32{
		"UNKNOWN":    0,
		"EXPERIMENT": 1,
		"OPERATION":  2,
	}
)

func (x Goal_ConnectionType) Enum() *Goal_ConnectionType {
	p := new(Goal_ConnectionType)
	*p = x
	return p
}

func (x Goal_ConnectionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Goal_ConnectionType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_experiment_goal_proto_enumTypes[0].Descriptor()
}

func (Goal_ConnectionType) Type() protoreflect.EnumType {
	return &file_proto_experiment_goal_proto_enumTypes[0]
}

func (x Goal_ConnectionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Goal_ConnectionType.Descriptor instead.
func (Goal_ConnectionType) EnumDescriptor() ([]byte, []int) {
	return file_proto_experiment_goal_proto_rawDescGZIP(), []int{0, 0}
}

type Goal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                       `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name           string                       `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description    string                       `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	Deleted        bool                         `protobuf:"varint,4,opt,name=deleted,proto3" json:"deleted"`
	CreatedAt      int64                        `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt      int64                        `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	IsInUseStatus  bool                         `protobuf:"varint,7,opt,name=is_in_use_status,json=isInUseStatus,proto3" json:"is_in_use_status"` // This field is set only when APIs return.
	Archived       bool                         `protobuf:"varint,8,opt,name=archived,proto3" json:"archived"`
	ConnectionType Goal_ConnectionType          `protobuf:"varint,9,opt,name=connection_type,json=connectionType,proto3,enum=bucketeer.experiment.Goal_ConnectionType" json:"connection_type"`
	Experiments    []*Goal_ExperimentReference  `protobuf:"bytes,10,rep,name=experiments,proto3" json:"experiments"`
	AutoOpsRules   []*Goal_AutoOpsRuleReference `protobuf:"bytes,11,rep,name=auto_ops_rules,json=autoOpsRules,proto3" json:"auto_ops_rules"`
}

func (x *Goal) Reset() {
	*x = Goal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_experiment_goal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Goal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Goal) ProtoMessage() {}

func (x *Goal) ProtoReflect() protoreflect.Message {
	mi := &file_proto_experiment_goal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Goal.ProtoReflect.Descriptor instead.
func (*Goal) Descriptor() ([]byte, []int) {
	return file_proto_experiment_goal_proto_rawDescGZIP(), []int{0}
}

func (x *Goal) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Goal) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Goal) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Goal) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

func (x *Goal) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Goal) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Goal) GetIsInUseStatus() bool {
	if x != nil {
		return x.IsInUseStatus
	}
	return false
}

func (x *Goal) GetArchived() bool {
	if x != nil {
		return x.Archived
	}
	return false
}

func (x *Goal) GetConnectionType() Goal_ConnectionType {
	if x != nil {
		return x.ConnectionType
	}
	return Goal_UNKNOWN
}

func (x *Goal) GetExperiments() []*Goal_ExperimentReference {
	if x != nil {
		return x.Experiments
	}
	return nil
}

func (x *Goal) GetAutoOpsRules() []*Goal_AutoOpsRuleReference {
	if x != nil {
		return x.AutoOpsRules
	}
	return nil
}

type Goal_ExperimentReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name      string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	FeatureId string            `protobuf:"bytes,3,opt,name=feature_id,json=featureId,proto3" json:"feature_id"`
	Status    Experiment_Status `protobuf:"varint,4,opt,name=status,proto3,enum=bucketeer.experiment.Experiment_Status" json:"status"`
}

func (x *Goal_ExperimentReference) Reset() {
	*x = Goal_ExperimentReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_experiment_goal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Goal_ExperimentReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Goal_ExperimentReference) ProtoMessage() {}

func (x *Goal_ExperimentReference) ProtoReflect() protoreflect.Message {
	mi := &file_proto_experiment_goal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Goal_ExperimentReference.ProtoReflect.Descriptor instead.
func (*Goal_ExperimentReference) Descriptor() ([]byte, []int) {
	return file_proto_experiment_goal_proto_rawDescGZIP(), []int{0, 0}
}

func (x *Goal_ExperimentReference) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Goal_ExperimentReference) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Goal_ExperimentReference) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

func (x *Goal_ExperimentReference) GetStatus() Experiment_Status {
	if x != nil {
		return x.Status
	}
	return Experiment_WAITING
}

type Goal_AutoOpsRuleReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FeatureId string `protobuf:"bytes,2,opt,name=feature_id,json=featureId,proto3" json:"feature_id"`
}

func (x *Goal_AutoOpsRuleReference) Reset() {
	*x = Goal_AutoOpsRuleReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_experiment_goal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Goal_AutoOpsRuleReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Goal_AutoOpsRuleReference) ProtoMessage() {}

func (x *Goal_AutoOpsRuleReference) ProtoReflect() protoreflect.Message {
	mi := &file_proto_experiment_goal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Goal_AutoOpsRuleReference.ProtoReflect.Descriptor instead.
func (*Goal_AutoOpsRuleReference) Descriptor() ([]byte, []int) {
	return file_proto_experiment_goal_proto_rawDescGZIP(), []int{0, 1}
}

func (x *Goal_AutoOpsRuleReference) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Goal_AutoOpsRuleReference) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

var File_proto_experiment_goal_proto protoreflect.FileDescriptor

var file_proto_experiment_goal_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x2f, 0x67, 0x6f, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d,
	0x65, 0x6e, 0x74, 0x1a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72,
	0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x87, 0x06, 0x0a, 0x04, 0x47, 0x6f, 0x61, 0x6c, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x27, 0x0a,
	0x10, 0x69, 0x73, 0x5f, 0x69, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x69, 0x73, 0x49, 0x6e, 0x55, 0x73, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76,
	0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76,
	0x65, 0x64, 0x12, 0x52, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x47, 0x6f, 0x61, 0x6c, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x50, 0x0a, 0x0b, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x2e, 0x47, 0x6f, 0x61, 0x6c, 0x2e, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0b, 0x65, 0x78, 0x70,
	0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x55, 0x0a, 0x0e, 0x61, 0x75, 0x74, 0x6f,
	0x5f, 0x6f, 0x70, 0x73, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x70,
	0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x47, 0x6f, 0x61, 0x6c, 0x2e, 0x41, 0x75, 0x74,
	0x6f, 0x4f, 0x70, 0x73, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x52, 0x0c, 0x61, 0x75, 0x74, 0x6f, 0x4f, 0x70, 0x73, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x1a,
	0x99, 0x01, 0x0a, 0x13, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64, 0x12, 0x3f, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x27, 0x2e, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x45, 0x0a, 0x14, 0x41,
	0x75, 0x74, 0x6f, 0x4f, 0x70, 0x73, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65,
	0x6e, 0x63, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x49, 0x64, 0x22, 0x3c, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x58, 0x50, 0x45, 0x52, 0x49, 0x4d, 0x45, 0x4e, 0x54, 0x10,
	0x01, 0x12, 0x0d, 0x0a, 0x09, 0x4f, 0x50, 0x45, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02,
	0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x70, 0x65,
	0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_experiment_goal_proto_rawDescOnce sync.Once
	file_proto_experiment_goal_proto_rawDescData = file_proto_experiment_goal_proto_rawDesc
)

func file_proto_experiment_goal_proto_rawDescGZIP() []byte {
	file_proto_experiment_goal_proto_rawDescOnce.Do(func() {
		file_proto_experiment_goal_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_experiment_goal_proto_rawDescData)
	})
	return file_proto_experiment_goal_proto_rawDescData
}

var file_proto_experiment_goal_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_experiment_goal_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_experiment_goal_proto_goTypes = []interface{}{
	(Goal_ConnectionType)(0),          // 0: bucketeer.experiment.Goal.ConnectionType
	(*Goal)(nil),                      // 1: bucketeer.experiment.Goal
	(*Goal_ExperimentReference)(nil),  // 2: bucketeer.experiment.Goal.ExperimentReference
	(*Goal_AutoOpsRuleReference)(nil), // 3: bucketeer.experiment.Goal.AutoOpsRuleReference
	(Experiment_Status)(0),            // 4: bucketeer.experiment.Experiment.Status
}
var file_proto_experiment_goal_proto_depIdxs = []int32{
	0, // 0: bucketeer.experiment.Goal.connection_type:type_name -> bucketeer.experiment.Goal.ConnectionType
	2, // 1: bucketeer.experiment.Goal.experiments:type_name -> bucketeer.experiment.Goal.ExperimentReference
	3, // 2: bucketeer.experiment.Goal.auto_ops_rules:type_name -> bucketeer.experiment.Goal.AutoOpsRuleReference
	4, // 3: bucketeer.experiment.Goal.ExperimentReference.status:type_name -> bucketeer.experiment.Experiment.Status
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_experiment_goal_proto_init() }
func file_proto_experiment_goal_proto_init() {
	if File_proto_experiment_goal_proto != nil {
		return
	}
	file_proto_experiment_experiment_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_experiment_goal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Goal); i {
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
		file_proto_experiment_goal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Goal_ExperimentReference); i {
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
		file_proto_experiment_goal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Goal_AutoOpsRuleReference); i {
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
			RawDescriptor: file_proto_experiment_goal_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_experiment_goal_proto_goTypes,
		DependencyIndexes: file_proto_experiment_goal_proto_depIdxs,
		EnumInfos:         file_proto_experiment_goal_proto_enumTypes,
		MessageInfos:      file_proto_experiment_goal_proto_msgTypes,
	}.Build()
	File_proto_experiment_goal_proto = out.File
	file_proto_experiment_goal_proto_rawDesc = nil
	file_proto_experiment_goal_proto_goTypes = nil
	file_proto_experiment_goal_proto_depIdxs = nil
}
