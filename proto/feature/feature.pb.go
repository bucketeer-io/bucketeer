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
// source: proto/feature/feature.proto

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

type Feature_VariationType int32

const (
	Feature_STRING  Feature_VariationType = 0
	Feature_BOOLEAN Feature_VariationType = 1
	Feature_NUMBER  Feature_VariationType = 2
	Feature_JSON    Feature_VariationType = 3
)

// Enum value maps for Feature_VariationType.
var (
	Feature_VariationType_name = map[int32]string{
		0: "STRING",
		1: "BOOLEAN",
		2: "NUMBER",
		3: "JSON",
	}
	Feature_VariationType_value = map[string]int32{
		"STRING":  0,
		"BOOLEAN": 1,
		"NUMBER":  2,
		"JSON":    3,
	}
)

func (x Feature_VariationType) Enum() *Feature_VariationType {
	p := new(Feature_VariationType)
	*p = x
	return p
}

func (x Feature_VariationType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Feature_VariationType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_feature_feature_proto_enumTypes[0].Descriptor()
}

func (Feature_VariationType) Type() protoreflect.EnumType {
	return &file_proto_feature_feature_proto_enumTypes[0]
}

func (x Feature_VariationType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Feature_VariationType.Descriptor instead.
func (Feature_VariationType) EnumDescriptor() ([]byte, []int) {
	return file_proto_feature_feature_proto_rawDescGZIP(), []int{0, 0}
}

type Feature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Enabled     bool   `protobuf:"varint,4,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Deleted     bool   `protobuf:"varint,5,opt,name=deleted,proto3" json:"deleted,omitempty"`
	// Deprecated: Do not use.
	EvaluationUndelayable bool                  `protobuf:"varint,6,opt,name=evaluation_undelayable,json=evaluationUndelayable,proto3" json:"evaluation_undelayable,omitempty"`
	Ttl                   int32                 `protobuf:"varint,7,opt,name=ttl,proto3" json:"ttl,omitempty"`
	Version               int32                 `protobuf:"varint,8,opt,name=version,proto3" json:"version,omitempty"`
	CreatedAt             int64                 `protobuf:"varint,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt             int64                 `protobuf:"varint,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Variations            []*Variation          `protobuf:"bytes,11,rep,name=variations,proto3" json:"variations,omitempty"`
	Targets               []*Target             `protobuf:"bytes,12,rep,name=targets,proto3" json:"targets,omitempty"`
	Rules                 []*Rule               `protobuf:"bytes,13,rep,name=rules,proto3" json:"rules,omitempty"`
	DefaultStrategy       *Strategy             `protobuf:"bytes,14,opt,name=default_strategy,json=defaultStrategy,proto3" json:"default_strategy,omitempty"`
	OffVariation          string                `protobuf:"bytes,15,opt,name=off_variation,json=offVariation,proto3" json:"off_variation,omitempty"`
	Tags                  []string              `protobuf:"bytes,16,rep,name=tags,proto3" json:"tags,omitempty"`
	LastUsedInfo          *FeatureLastUsedInfo  `protobuf:"bytes,17,opt,name=last_used_info,json=lastUsedInfo,proto3" json:"last_used_info,omitempty"`
	Maintainer            string                `protobuf:"bytes,18,opt,name=maintainer,proto3" json:"maintainer,omitempty"`
	VariationType         Feature_VariationType `protobuf:"varint,19,opt,name=variation_type,json=variationType,proto3,enum=bucketeer.feature.Feature_VariationType" json:"variation_type,omitempty"`
	Archived              bool                  `protobuf:"varint,20,opt,name=archived,proto3" json:"archived,omitempty"`
	Prerequisites         []*Prerequisite       `protobuf:"bytes,21,rep,name=prerequisites,proto3" json:"prerequisites,omitempty"`
	SamplingSeed          string                `protobuf:"bytes,22,opt,name=sampling_seed,json=samplingSeed,proto3" json:"sampling_seed,omitempty"`
}

func (x *Feature) Reset() {
	*x = Feature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_feature_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Feature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Feature) ProtoMessage() {}

func (x *Feature) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_feature_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Feature.ProtoReflect.Descriptor instead.
func (*Feature) Descriptor() ([]byte, []int) {
	return file_proto_feature_feature_proto_rawDescGZIP(), []int{0}
}

func (x *Feature) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Feature) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Feature) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Feature) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

func (x *Feature) GetDeleted() bool {
	if x != nil {
		return x.Deleted
	}
	return false
}

// Deprecated: Do not use.
func (x *Feature) GetEvaluationUndelayable() bool {
	if x != nil {
		return x.EvaluationUndelayable
	}
	return false
}

func (x *Feature) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

func (x *Feature) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Feature) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Feature) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Feature) GetVariations() []*Variation {
	if x != nil {
		return x.Variations
	}
	return nil
}

func (x *Feature) GetTargets() []*Target {
	if x != nil {
		return x.Targets
	}
	return nil
}

func (x *Feature) GetRules() []*Rule {
	if x != nil {
		return x.Rules
	}
	return nil
}

func (x *Feature) GetDefaultStrategy() *Strategy {
	if x != nil {
		return x.DefaultStrategy
	}
	return nil
}

func (x *Feature) GetOffVariation() string {
	if x != nil {
		return x.OffVariation
	}
	return ""
}

func (x *Feature) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Feature) GetLastUsedInfo() *FeatureLastUsedInfo {
	if x != nil {
		return x.LastUsedInfo
	}
	return nil
}

func (x *Feature) GetMaintainer() string {
	if x != nil {
		return x.Maintainer
	}
	return ""
}

func (x *Feature) GetVariationType() Feature_VariationType {
	if x != nil {
		return x.VariationType
	}
	return Feature_STRING
}

func (x *Feature) GetArchived() bool {
	if x != nil {
		return x.Archived
	}
	return false
}

func (x *Feature) GetPrerequisites() []*Prerequisite {
	if x != nil {
		return x.Prerequisites
	}
	return nil
}

func (x *Feature) GetSamplingSeed() string {
	if x != nil {
		return x.SamplingSeed
	}
	return ""
}

type Features struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Features []*Feature `protobuf:"bytes,1,rep,name=features,proto3" json:"features,omitempty"`
}

func (x *Features) Reset() {
	*x = Features{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_feature_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Features) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Features) ProtoMessage() {}

func (x *Features) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_feature_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Features.ProtoReflect.Descriptor instead.
func (*Features) Descriptor() ([]byte, []int) {
	return file_proto_feature_feature_proto_rawDescGZIP(), []int{1}
}

func (x *Features) GetFeatures() []*Feature {
	if x != nil {
		return x.Features
	}
	return nil
}

type Tag struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt int64  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt int64  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Tag) Reset() {
	*x = Tag{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_feature_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Tag) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Tag) ProtoMessage() {}

func (x *Tag) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_feature_proto_msgTypes[2]
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
	return file_proto_feature_feature_proto_rawDescGZIP(), []int{2}
}

func (x *Tag) GetId() string {
	if x != nil {
		return x.Id
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

var File_proto_feature_feature_proto protoreflect.FileDescriptor

var file_proto_feature_feature_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x1a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x2f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x2a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x75, 0x73, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x70,
	0x72, 0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xd2, 0x07, 0x0a, 0x07, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x65, 0x6e, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x39, 0x0a, 0x16, 0x65, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x6e, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x61, 0x62, 0x6c,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x42, 0x02, 0x18, 0x01, 0x52, 0x15, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x55, 0x6e, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x61, 0x62,
	0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x03, 0x74, 0x74, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x3c, 0x0a, 0x0a,
	0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a,
	0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x73, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e,
	0x54, 0x61, 0x72, 0x67, 0x65, 0x74, 0x52, 0x07, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x73, 0x12,
	0x2d, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x46,
	0x0a, 0x10, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65,
	0x67, 0x79, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x52, 0x0f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x53, 0x74,
	0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x66, 0x66, 0x5f, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6f,
	0x66, 0x66, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x10, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12,
	0x4c, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x75, 0x73, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x66,
	0x6f, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x4c, 0x61, 0x73, 0x74, 0x55, 0x73, 0x65, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x0c, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x73, 0x65, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1e, 0x0a,
	0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x18, 0x12, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x4f, 0x0a,
	0x0e, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x13, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x0d, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x18, 0x14, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x08, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x64, 0x12, 0x45, 0x0a, 0x0d, 0x70, 0x72,
	0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x65, 0x73, 0x18, 0x15, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x50, 0x72, 0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69,
	0x74, 0x65, 0x52, 0x0d, 0x70, 0x72, 0x65, 0x72, 0x65, 0x71, 0x75, 0x69, 0x73, 0x69, 0x74, 0x65,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x69, 0x6e, 0x67, 0x5f, 0x73, 0x65,
	0x65, 0x64, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x69,
	0x6e, 0x67, 0x53, 0x65, 0x65, 0x64, 0x22, 0x3e, 0x0a, 0x0d, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x54, 0x52, 0x49, 0x4e,
	0x47, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x42, 0x4f, 0x4f, 0x4c, 0x45, 0x41, 0x4e, 0x10, 0x01,
	0x12, 0x0a, 0x0a, 0x06, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04,
	0x4a, 0x53, 0x4f, 0x4e, 0x10, 0x03, 0x22, 0x42, 0x0a, 0x08, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x73, 0x12, 0x36, 0x0a, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x46, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x52, 0x08, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x73, 0x22, 0x53, 0x0a, 0x03, 0x54, 0x61,
	0x67, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x42,
	0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65,
	0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_feature_feature_proto_rawDescOnce sync.Once
	file_proto_feature_feature_proto_rawDescData = file_proto_feature_feature_proto_rawDesc
)

func file_proto_feature_feature_proto_rawDescGZIP() []byte {
	file_proto_feature_feature_proto_rawDescOnce.Do(func() {
		file_proto_feature_feature_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_feature_feature_proto_rawDescData)
	})
	return file_proto_feature_feature_proto_rawDescData
}

var file_proto_feature_feature_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_feature_feature_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_feature_feature_proto_goTypes = []interface{}{
	(Feature_VariationType)(0),  // 0: bucketeer.feature.Feature.VariationType
	(*Feature)(nil),             // 1: bucketeer.feature.Feature
	(*Features)(nil),            // 2: bucketeer.feature.Features
	(*Tag)(nil),                 // 3: bucketeer.feature.Tag
	(*Variation)(nil),           // 4: bucketeer.feature.Variation
	(*Target)(nil),              // 5: bucketeer.feature.Target
	(*Rule)(nil),                // 6: bucketeer.feature.Rule
	(*Strategy)(nil),            // 7: bucketeer.feature.Strategy
	(*FeatureLastUsedInfo)(nil), // 8: bucketeer.feature.FeatureLastUsedInfo
	(*Prerequisite)(nil),        // 9: bucketeer.feature.Prerequisite
}
var file_proto_feature_feature_proto_depIdxs = []int32{
	4, // 0: bucketeer.feature.Feature.variations:type_name -> bucketeer.feature.Variation
	5, // 1: bucketeer.feature.Feature.targets:type_name -> bucketeer.feature.Target
	6, // 2: bucketeer.feature.Feature.rules:type_name -> bucketeer.feature.Rule
	7, // 3: bucketeer.feature.Feature.default_strategy:type_name -> bucketeer.feature.Strategy
	8, // 4: bucketeer.feature.Feature.last_used_info:type_name -> bucketeer.feature.FeatureLastUsedInfo
	0, // 5: bucketeer.feature.Feature.variation_type:type_name -> bucketeer.feature.Feature.VariationType
	9, // 6: bucketeer.feature.Feature.prerequisites:type_name -> bucketeer.feature.Prerequisite
	1, // 7: bucketeer.feature.Features.features:type_name -> bucketeer.feature.Feature
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_proto_feature_feature_proto_init() }
func file_proto_feature_feature_proto_init() {
	if File_proto_feature_feature_proto != nil {
		return
	}
	file_proto_feature_rule_proto_init()
	file_proto_feature_target_proto_init()
	file_proto_feature_variation_proto_init()
	file_proto_feature_strategy_proto_init()
	file_proto_feature_feature_last_used_info_proto_init()
	file_proto_feature_prerequisite_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_feature_feature_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Feature); i {
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
		file_proto_feature_feature_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Features); i {
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
		file_proto_feature_feature_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_feature_feature_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_feature_feature_proto_goTypes,
		DependencyIndexes: file_proto_feature_feature_proto_depIdxs,
		EnumInfos:         file_proto_feature_feature_proto_enumTypes,
		MessageInfos:      file_proto_feature_feature_proto_msgTypes,
	}.Build()
	File_proto_feature_feature_proto = out.File
	file_proto_feature_feature_proto_rawDesc = nil
	file_proto_feature_feature_proto_goTypes = nil
	file_proto_feature_feature_proto_depIdxs = nil
}
