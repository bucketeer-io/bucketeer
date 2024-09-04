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
// source: proto/eventpersisterdwh/goal_event.proto

package eventcounter

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

type GoalEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	GoalId         string  `protobuf:"bytes,2,opt,name=goal_id,json=goalId,proto3" json:"goal_id"`
	Value          float32 `protobuf:"fixed32,3,opt,name=value,proto3" json:"value"`
	UserData       string  `protobuf:"bytes,4,opt,name=user_data,json=userData,proto3" json:"user_data"`
	UserId         string  `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Tag            string  `protobuf:"bytes,6,opt,name=tag,proto3" json:"tag"`
	SourceId       string  `protobuf:"bytes,7,opt,name=source_id,json=sourceId,proto3" json:"source_id"`
	Timestamp      int64   `protobuf:"varint,9,opt,name=timestamp,proto3" json:"timestamp"`
	FeatureId      string  `protobuf:"bytes,10,opt,name=feature_id,json=featureId,proto3" json:"feature_id"`
	FeatureVersion int32   `protobuf:"varint,11,opt,name=feature_version,json=featureVersion,proto3" json:"feature_version"`
	VariationId    string  `protobuf:"bytes,12,opt,name=variation_id,json=variationId,proto3" json:"variation_id"`
	Reason         string  `protobuf:"bytes,13,opt,name=reason,proto3" json:"reason"`
	EnvironmentId  string  `protobuf:"bytes,14,opt,name=environment_id,json=environmentId,proto3" json:"environment_id"`
}

func (x *GoalEvent) Reset() {
	*x = GoalEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eventpersisterdwh_goal_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoalEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoalEvent) ProtoMessage() {}

func (x *GoalEvent) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eventpersisterdwh_goal_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoalEvent.ProtoReflect.Descriptor instead.
func (*GoalEvent) Descriptor() ([]byte, []int) {
	return file_proto_eventpersisterdwh_goal_event_proto_rawDescGZIP(), []int{0}
}

func (x *GoalEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GoalEvent) GetGoalId() string {
	if x != nil {
		return x.GoalId
	}
	return ""
}

func (x *GoalEvent) GetValue() float32 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *GoalEvent) GetUserData() string {
	if x != nil {
		return x.UserData
	}
	return ""
}

func (x *GoalEvent) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GoalEvent) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

func (x *GoalEvent) GetSourceId() string {
	if x != nil {
		return x.SourceId
	}
	return ""
}

func (x *GoalEvent) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *GoalEvent) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

func (x *GoalEvent) GetFeatureVersion() int32 {
	if x != nil {
		return x.FeatureVersion
	}
	return 0
}

func (x *GoalEvent) GetVariationId() string {
	if x != nil {
		return x.VariationId
	}
	return ""
}

func (x *GoalEvent) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *GoalEvent) GetEnvironmentId() string {
	if x != nil {
		return x.EnvironmentId
	}
	return ""
}

var File_proto_eventpersisterdwh_goal_event_proto protoreflect.FileDescriptor

var file_proto_eventpersisterdwh_goal_event_proto_rawDesc = []byte{
	0x0a, 0x28, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x65, 0x72,
	0x73, 0x69, 0x73, 0x74, 0x65, 0x72, 0x64, 0x77, 0x68, 0x2f, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x22, 0xfd, 0x02, 0x0a, 0x09, 0x47, 0x6f, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x67, 0x6f, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x74, 0x61, 0x67, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x69, 0x64,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x49,
	0x64, 0x12, 0x27, 0x0a, 0x0f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e,
	0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x4a, 0x04, 0x08, 0x08,
	0x10, 0x09, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_eventpersisterdwh_goal_event_proto_rawDescOnce sync.Once
	file_proto_eventpersisterdwh_goal_event_proto_rawDescData = file_proto_eventpersisterdwh_goal_event_proto_rawDesc
)

func file_proto_eventpersisterdwh_goal_event_proto_rawDescGZIP() []byte {
	file_proto_eventpersisterdwh_goal_event_proto_rawDescOnce.Do(func() {
		file_proto_eventpersisterdwh_goal_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_eventpersisterdwh_goal_event_proto_rawDescData)
	})
	return file_proto_eventpersisterdwh_goal_event_proto_rawDescData
}

var file_proto_eventpersisterdwh_goal_event_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_eventpersisterdwh_goal_event_proto_goTypes = []interface{}{
	(*GoalEvent)(nil), // 0: bucketeer.eventcounter.GoalEvent
}
var file_proto_eventpersisterdwh_goal_event_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_eventpersisterdwh_goal_event_proto_init() }
func file_proto_eventpersisterdwh_goal_event_proto_init() {
	if File_proto_eventpersisterdwh_goal_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_eventpersisterdwh_goal_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoalEvent); i {
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
			RawDescriptor: file_proto_eventpersisterdwh_goal_event_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_eventpersisterdwh_goal_event_proto_goTypes,
		DependencyIndexes: file_proto_eventpersisterdwh_goal_event_proto_depIdxs,
		MessageInfos:      file_proto_eventpersisterdwh_goal_event_proto_msgTypes,
	}.Build()
	File_proto_eventpersisterdwh_goal_event_proto = out.File
	file_proto_eventpersisterdwh_goal_event_proto_rawDesc = nil
	file_proto_eventpersisterdwh_goal_event_proto_goTypes = nil
	file_proto_eventpersisterdwh_goal_event_proto_depIdxs = nil
}
