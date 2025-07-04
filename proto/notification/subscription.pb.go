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
// source: proto/notification/subscription.proto

package notification

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

type Subscription_SourceType int32

const (
	Subscription_DOMAIN_EVENT_FEATURE             Subscription_SourceType = 0
	Subscription_DOMAIN_EVENT_GOAL                Subscription_SourceType = 1
	Subscription_DOMAIN_EVENT_EXPERIMENT          Subscription_SourceType = 2
	Subscription_DOMAIN_EVENT_ACCOUNT             Subscription_SourceType = 3
	Subscription_DOMAIN_EVENT_APIKEY              Subscription_SourceType = 4
	Subscription_DOMAIN_EVENT_SEGMENT             Subscription_SourceType = 5
	Subscription_DOMAIN_EVENT_ENVIRONMENT         Subscription_SourceType = 6
	Subscription_DOMAIN_EVENT_ADMIN_ACCOUNT       Subscription_SourceType = 7
	Subscription_DOMAIN_EVENT_AUTOOPS_RULE        Subscription_SourceType = 8
	Subscription_DOMAIN_EVENT_PUSH                Subscription_SourceType = 9
	Subscription_DOMAIN_EVENT_SUBSCRIPTION        Subscription_SourceType = 10
	Subscription_DOMAIN_EVENT_ADMIN_SUBSCRIPTION  Subscription_SourceType = 11
	Subscription_DOMAIN_EVENT_PROJECT             Subscription_SourceType = 12
	Subscription_DOMAIN_EVENT_WEBHOOK             Subscription_SourceType = 13
	Subscription_DOMAIN_EVENT_PROGRESSIVE_ROLLOUT Subscription_SourceType = 14
	Subscription_DOMAIN_EVENT_ORGANIZATION        Subscription_SourceType = 15
	Subscription_DOMAIN_EVENT_FLAG_TRIGGER        Subscription_SourceType = 16
	Subscription_DOMAIN_EVENT_TAG                 Subscription_SourceType = 17
	Subscription_DOMAIN_EVENT_CODEREF             Subscription_SourceType = 18
	Subscription_DOMAIN_EVENT_TEAM                Subscription_SourceType = 19
	Subscription_FEATURE_STALE                    Subscription_SourceType = 100
	Subscription_EXPERIMENT_RUNNING               Subscription_SourceType = 200
	Subscription_MAU_COUNT                        Subscription_SourceType = 300
)

// Enum value maps for Subscription_SourceType.
var (
	Subscription_SourceType_name = map[int32]string{
		0:   "DOMAIN_EVENT_FEATURE",
		1:   "DOMAIN_EVENT_GOAL",
		2:   "DOMAIN_EVENT_EXPERIMENT",
		3:   "DOMAIN_EVENT_ACCOUNT",
		4:   "DOMAIN_EVENT_APIKEY",
		5:   "DOMAIN_EVENT_SEGMENT",
		6:   "DOMAIN_EVENT_ENVIRONMENT",
		7:   "DOMAIN_EVENT_ADMIN_ACCOUNT",
		8:   "DOMAIN_EVENT_AUTOOPS_RULE",
		9:   "DOMAIN_EVENT_PUSH",
		10:  "DOMAIN_EVENT_SUBSCRIPTION",
		11:  "DOMAIN_EVENT_ADMIN_SUBSCRIPTION",
		12:  "DOMAIN_EVENT_PROJECT",
		13:  "DOMAIN_EVENT_WEBHOOK",
		14:  "DOMAIN_EVENT_PROGRESSIVE_ROLLOUT",
		15:  "DOMAIN_EVENT_ORGANIZATION",
		16:  "DOMAIN_EVENT_FLAG_TRIGGER",
		17:  "DOMAIN_EVENT_TAG",
		18:  "DOMAIN_EVENT_CODEREF",
		19:  "DOMAIN_EVENT_TEAM",
		100: "FEATURE_STALE",
		200: "EXPERIMENT_RUNNING",
		300: "MAU_COUNT",
	}
	Subscription_SourceType_value = map[string]int32{
		"DOMAIN_EVENT_FEATURE":             0,
		"DOMAIN_EVENT_GOAL":                1,
		"DOMAIN_EVENT_EXPERIMENT":          2,
		"DOMAIN_EVENT_ACCOUNT":             3,
		"DOMAIN_EVENT_APIKEY":              4,
		"DOMAIN_EVENT_SEGMENT":             5,
		"DOMAIN_EVENT_ENVIRONMENT":         6,
		"DOMAIN_EVENT_ADMIN_ACCOUNT":       7,
		"DOMAIN_EVENT_AUTOOPS_RULE":        8,
		"DOMAIN_EVENT_PUSH":                9,
		"DOMAIN_EVENT_SUBSCRIPTION":        10,
		"DOMAIN_EVENT_ADMIN_SUBSCRIPTION":  11,
		"DOMAIN_EVENT_PROJECT":             12,
		"DOMAIN_EVENT_WEBHOOK":             13,
		"DOMAIN_EVENT_PROGRESSIVE_ROLLOUT": 14,
		"DOMAIN_EVENT_ORGANIZATION":        15,
		"DOMAIN_EVENT_FLAG_TRIGGER":        16,
		"DOMAIN_EVENT_TAG":                 17,
		"DOMAIN_EVENT_CODEREF":             18,
		"DOMAIN_EVENT_TEAM":                19,
		"FEATURE_STALE":                    100,
		"EXPERIMENT_RUNNING":               200,
		"MAU_COUNT":                        300,
	}
)

func (x Subscription_SourceType) Enum() *Subscription_SourceType {
	p := new(Subscription_SourceType)
	*p = x
	return p
}

func (x Subscription_SourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Subscription_SourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_notification_subscription_proto_enumTypes[0].Descriptor()
}

func (Subscription_SourceType) Type() protoreflect.EnumType {
	return &file_proto_notification_subscription_proto_enumTypes[0]
}

func (x Subscription_SourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Subscription_SourceType.Descriptor instead.
func (Subscription_SourceType) EnumDescriptor() ([]byte, []int) {
	return file_proto_notification_subscription_proto_rawDescGZIP(), []int{0, 0}
}

type Subscription struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                    `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	CreatedAt       int64                     `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt       int64                     `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	Disabled        bool                      `protobuf:"varint,4,opt,name=disabled,proto3" json:"disabled"`
	SourceTypes     []Subscription_SourceType `protobuf:"varint,5,rep,packed,name=source_types,json=sourceTypes,proto3,enum=bucketeer.notification.Subscription_SourceType" json:"source_types"`
	Recipient       *Recipient                `protobuf:"bytes,6,opt,name=recipient,proto3" json:"recipient"`
	Name            string                    `protobuf:"bytes,7,opt,name=name,proto3" json:"name"`
	EnvironmentId   string                    `protobuf:"bytes,8,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
	EnvironmentName string                    `protobuf:"bytes,9,opt,name=environment_name,json=environmentName,proto3" json:"environment_name"`
	FeatureFlagTags []string                  `protobuf:"bytes,10,rep,name=feature_flag_tags,json=featureFlagTags,proto3" json:"feature_flag_tags"`
}

func (x *Subscription) Reset() {
	*x = Subscription{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_notification_subscription_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Subscription) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Subscription) ProtoMessage() {}

func (x *Subscription) ProtoReflect() protoreflect.Message {
	mi := &file_proto_notification_subscription_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Subscription.ProtoReflect.Descriptor instead.
func (*Subscription) Descriptor() ([]byte, []int) {
	return file_proto_notification_subscription_proto_rawDescGZIP(), []int{0}
}

func (x *Subscription) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Subscription) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *Subscription) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *Subscription) GetDisabled() bool {
	if x != nil {
		return x.Disabled
	}
	return false
}

func (x *Subscription) GetSourceTypes() []Subscription_SourceType {
	if x != nil {
		return x.SourceTypes
	}
	return nil
}

func (x *Subscription) GetRecipient() *Recipient {
	if x != nil {
		return x.Recipient
	}
	return nil
}

func (x *Subscription) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Subscription) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

func (x *Subscription) GetEnvironmentName() string {
	if x != nil {
		return x.EnvironmentName
	}
	return ""
}

func (x *Subscription) GetFeatureFlagTags() []string {
	if x != nil {
		return x.FeatureFlagTags
	}
	return nil
}

var File_proto_notification_subscription_proto protoreflect.FileDescriptor

var file_proto_notification_subscription_proto_rawDesc = []byte{
	0x0a, 0x25, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65,
	0x65, 0x72, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x1a,
	0x22, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9c, 0x08, 0x0a, 0x0c, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f,
	0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64,
	0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x64, 0x69, 0x73, 0x61, 0x62, 0x6c, 0x65, 0x64, 0x12, 0x52,
	0x0a, 0x0c, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x0e, 0x32, 0x2f, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x75,
	0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x73, 0x12, 0x3f, 0x0a, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72,
	0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x29,
	0x0a, 0x10, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f,
	0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2a, 0x0a, 0x11, 0x66, 0x65, 0x61,
	0x74, 0x75, 0x72, 0x65, 0x5f, 0x66, 0x6c, 0x61, 0x67, 0x5f, 0x74, 0x61, 0x67, 0x73, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x46, 0x6c, 0x61,
	0x67, 0x54, 0x61, 0x67, 0x73, 0x22, 0xfa, 0x04, 0x0a, 0x0a, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45,
	0x56, 0x45, 0x4e, 0x54, 0x5f, 0x46, 0x45, 0x41, 0x54, 0x55, 0x52, 0x45, 0x10, 0x00, 0x12, 0x15,
	0x0a, 0x11, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x47,
	0x4f, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f,
	0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x45, 0x58, 0x50, 0x45, 0x52, 0x49, 0x4d, 0x45, 0x4e, 0x54,
	0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45,
	0x4e, 0x54, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13,
	0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x41, 0x50, 0x49,
	0x4b, 0x45, 0x59, 0x10, 0x04, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f,
	0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x45, 0x47, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x05, 0x12,
	0x1c, 0x0a, 0x18, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f,
	0x45, 0x4e, 0x56, 0x49, 0x52, 0x4f, 0x4e, 0x4d, 0x45, 0x4e, 0x54, 0x10, 0x06, 0x12, 0x1e, 0x0a,
	0x1a, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x41, 0x44,
	0x4d, 0x49, 0x4e, 0x5f, 0x41, 0x43, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10, 0x07, 0x12, 0x1d, 0x0a,
	0x19, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x41, 0x55,
	0x54, 0x4f, 0x4f, 0x50, 0x53, 0x5f, 0x52, 0x55, 0x4c, 0x45, 0x10, 0x08, 0x12, 0x15, 0x0a, 0x11,
	0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x55, 0x53,
	0x48, 0x10, 0x09, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56,
	0x45, 0x4e, 0x54, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49, 0x50, 0x54, 0x49, 0x4f, 0x4e,
	0x10, 0x0a, 0x12, 0x23, 0x0a, 0x1f, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45,
	0x4e, 0x54, 0x5f, 0x41, 0x44, 0x4d, 0x49, 0x4e, 0x5f, 0x53, 0x55, 0x42, 0x53, 0x43, 0x52, 0x49,
	0x50, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0b, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49,
	0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x4a, 0x45, 0x43, 0x54, 0x10,
	0x0c, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e,
	0x54, 0x5f, 0x57, 0x45, 0x42, 0x48, 0x4f, 0x4f, 0x4b, 0x10, 0x0d, 0x12, 0x24, 0x0a, 0x20, 0x44,
	0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x47,
	0x52, 0x45, 0x53, 0x53, 0x49, 0x56, 0x45, 0x5f, 0x52, 0x4f, 0x4c, 0x4c, 0x4f, 0x55, 0x54, 0x10,
	0x0e, 0x12, 0x1d, 0x0a, 0x19, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e,
	0x54, 0x5f, 0x4f, 0x52, 0x47, 0x41, 0x4e, 0x49, 0x5a, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x0f,
	0x12, 0x1d, 0x0a, 0x19, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54,
	0x5f, 0x46, 0x4c, 0x41, 0x47, 0x5f, 0x54, 0x52, 0x49, 0x47, 0x47, 0x45, 0x52, 0x10, 0x10, 0x12,
	0x14, 0x0a, 0x10, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x41, 0x47, 0x10, 0x11, 0x12, 0x18, 0x0a, 0x14, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f,
	0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x52, 0x45, 0x46, 0x10, 0x12, 0x12,
	0x15, 0x0a, 0x11, 0x44, 0x4f, 0x4d, 0x41, 0x49, 0x4e, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x45, 0x41, 0x4d, 0x10, 0x13, 0x12, 0x11, 0x0a, 0x0d, 0x46, 0x45, 0x41, 0x54, 0x55, 0x52,
	0x45, 0x5f, 0x53, 0x54, 0x41, 0x4c, 0x45, 0x10, 0x64, 0x12, 0x17, 0x0a, 0x12, 0x45, 0x58, 0x50,
	0x45, 0x52, 0x49, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10,
	0xc8, 0x01, 0x12, 0x0e, 0x0a, 0x09, 0x4d, 0x41, 0x55, 0x5f, 0x43, 0x4f, 0x55, 0x4e, 0x54, 0x10,
	0xac, 0x02, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x6f,
	0x74, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_notification_subscription_proto_rawDescOnce sync.Once
	file_proto_notification_subscription_proto_rawDescData = file_proto_notification_subscription_proto_rawDesc
)

func file_proto_notification_subscription_proto_rawDescGZIP() []byte {
	file_proto_notification_subscription_proto_rawDescOnce.Do(func() {
		file_proto_notification_subscription_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_notification_subscription_proto_rawDescData)
	})
	return file_proto_notification_subscription_proto_rawDescData
}

var file_proto_notification_subscription_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_notification_subscription_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_notification_subscription_proto_goTypes = []interface{}{
	(Subscription_SourceType)(0), // 0: bucketeer.notification.Subscription.SourceType
	(*Subscription)(nil),         // 1: bucketeer.notification.Subscription
	(*Recipient)(nil),            // 2: bucketeer.notification.Recipient
}
var file_proto_notification_subscription_proto_depIdxs = []int32{
	0, // 0: bucketeer.notification.Subscription.source_types:type_name -> bucketeer.notification.Subscription.SourceType
	2, // 1: bucketeer.notification.Subscription.recipient:type_name -> bucketeer.notification.Recipient
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_notification_subscription_proto_init() }
func file_proto_notification_subscription_proto_init() {
	if File_proto_notification_subscription_proto != nil {
		return
	}
	file_proto_notification_recipient_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_notification_subscription_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Subscription); i {
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
			RawDescriptor: file_proto_notification_subscription_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_notification_subscription_proto_goTypes,
		DependencyIndexes: file_proto_notification_subscription_proto_depIdxs,
		EnumInfos:         file_proto_notification_subscription_proto_enumTypes,
		MessageInfos:      file_proto_notification_subscription_proto_msgTypes,
	}.Build()
	File_proto_notification_subscription_proto = out.File
	file_proto_notification_subscription_proto_rawDesc = nil
	file_proto_notification_subscription_proto_goTypes = nil
	file_proto_notification_subscription_proto_depIdxs = nil
}
