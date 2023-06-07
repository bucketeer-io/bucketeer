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
// 	protoc        v3.18.1
// source: proto/eventcounter/experiment_count.proto

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

type ExperimentCount struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FeatureId      string `protobuf:"bytes,2,opt,name=feature_id,json=featureId,proto3" json:"feature_id,omitempty"`
	FeatureVersion int32  `protobuf:"varint,3,opt,name=feature_version,json=featureVersion,proto3" json:"feature_version,omitempty"`
	// Deprecated: Do not use.
	GoalId string `protobuf:"bytes,4,opt,name=goal_id,json=goalId,proto3" json:"goal_id,omitempty"`
	// Deprecated: Do not use.
	RealtimeCounts []*VariationCount `protobuf:"bytes,5,rep,name=realtime_counts,json=realtimeCounts,proto3" json:"realtime_counts,omitempty"`
	// Deprecated: Do not use.
	BatchCounts []*VariationCount `protobuf:"bytes,6,rep,name=batch_counts,json=batchCounts,proto3" json:"batch_counts,omitempty"`
	UpdatedAt   int64             `protobuf:"varint,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	GoalCounts  []*GoalCounts     `protobuf:"bytes,8,rep,name=goal_counts,json=goalCounts,proto3" json:"goal_counts,omitempty"`
}

func (x *ExperimentCount) Reset() {
	*x = ExperimentCount{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eventcounter_experiment_count_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExperimentCount) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExperimentCount) ProtoMessage() {}

func (x *ExperimentCount) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eventcounter_experiment_count_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExperimentCount.ProtoReflect.Descriptor instead.
func (*ExperimentCount) Descriptor() ([]byte, []int) {
	return file_proto_eventcounter_experiment_count_proto_rawDescGZIP(), []int{0}
}

func (x *ExperimentCount) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ExperimentCount) GetFeatureId() string {
	if x != nil {
		return x.FeatureId
	}
	return ""
}

func (x *ExperimentCount) GetFeatureVersion() int32 {
	if x != nil {
		return x.FeatureVersion
	}
	return 0
}

// Deprecated: Do not use.
func (x *ExperimentCount) GetGoalId() string {
	if x != nil {
		return x.GoalId
	}
	return ""
}

// Deprecated: Do not use.
func (x *ExperimentCount) GetRealtimeCounts() []*VariationCount {
	if x != nil {
		return x.RealtimeCounts
	}
	return nil
}

// Deprecated: Do not use.
func (x *ExperimentCount) GetBatchCounts() []*VariationCount {
	if x != nil {
		return x.BatchCounts
	}
	return nil
}

func (x *ExperimentCount) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *ExperimentCount) GetGoalCounts() []*GoalCounts {
	if x != nil {
		return x.GoalCounts
	}
	return nil
}

type GoalCounts struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoalId         string            `protobuf:"bytes,1,opt,name=goal_id,json=goalId,proto3" json:"goal_id,omitempty"`
	RealtimeCounts []*VariationCount `protobuf:"bytes,2,rep,name=realtime_counts,json=realtimeCounts,proto3" json:"realtime_counts,omitempty"`
	// Deprecated: Do not use.
	BatchCounts []*VariationCount `protobuf:"bytes,3,rep,name=batch_counts,json=batchCounts,proto3" json:"batch_counts,omitempty"`
}

func (x *GoalCounts) Reset() {
	*x = GoalCounts{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eventcounter_experiment_count_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoalCounts) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoalCounts) ProtoMessage() {}

func (x *GoalCounts) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eventcounter_experiment_count_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoalCounts.ProtoReflect.Descriptor instead.
func (*GoalCounts) Descriptor() ([]byte, []int) {
	return file_proto_eventcounter_experiment_count_proto_rawDescGZIP(), []int{1}
}

func (x *GoalCounts) GetGoalId() string {
	if x != nil {
		return x.GoalId
	}
	return ""
}

func (x *GoalCounts) GetRealtimeCounts() []*VariationCount {
	if x != nil {
		return x.RealtimeCounts
	}
	return nil
}

// Deprecated: Do not use.
func (x *GoalCounts) GetBatchCounts() []*VariationCount {
	if x != nil {
		return x.BatchCounts
	}
	return nil
}

var File_proto_eventcounter_experiment_count_proto protoreflect.FileDescriptor

var file_proto_eventcounter_experiment_count_proto_rawDesc = []byte{
	0x0a, 0x29, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x2f, 0x65, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x1a, 0x28, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x03,
	0x0a, 0x0f, 0x45, 0x78, 0x70, 0x65, 0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x49, 0x64,
	0x12, 0x27, 0x0a, 0x0f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x66, 0x65, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x07, 0x67, 0x6f, 0x61,
	0x6c, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x06,
	0x67, 0x6f, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x53, 0x0a, 0x0f, 0x72, 0x65, 0x61, 0x6c, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0e, 0x72, 0x65, 0x61,
	0x6c, 0x74, 0x69, 0x6d, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x4d, 0x0a, 0x0c, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0b, 0x62,
	0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x43, 0x0a, 0x0b, 0x67, 0x6f, 0x61,
	0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x47, 0x6f, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x73, 0x52, 0x0a, 0x67, 0x6f, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0xc5,
	0x01, 0x0a, 0x0a, 0x47, 0x6f, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x17, 0x0a,
	0x07, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x67, 0x6f, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x4f, 0x0a, 0x0f, 0x72, 0x65, 0x61, 0x6c, 0x74, 0x69,
	0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x0e, 0x72, 0x65, 0x61, 0x6c, 0x74, 0x69, 0x6d,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x12, 0x4d, 0x0a, 0x0c, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x02, 0x18, 0x01, 0x52, 0x0b, 0x62, 0x61, 0x74, 0x63, 0x68,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69,
	0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_eventcounter_experiment_count_proto_rawDescOnce sync.Once
	file_proto_eventcounter_experiment_count_proto_rawDescData = file_proto_eventcounter_experiment_count_proto_rawDesc
)

func file_proto_eventcounter_experiment_count_proto_rawDescGZIP() []byte {
	file_proto_eventcounter_experiment_count_proto_rawDescOnce.Do(func() {
		file_proto_eventcounter_experiment_count_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_eventcounter_experiment_count_proto_rawDescData)
	})
	return file_proto_eventcounter_experiment_count_proto_rawDescData
}

var file_proto_eventcounter_experiment_count_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_eventcounter_experiment_count_proto_goTypes = []interface{}{
	(*ExperimentCount)(nil), // 0: bucketeer.eventcounter.ExperimentCount
	(*GoalCounts)(nil),      // 1: bucketeer.eventcounter.GoalCounts
	(*VariationCount)(nil),  // 2: bucketeer.eventcounter.VariationCount
}
var file_proto_eventcounter_experiment_count_proto_depIdxs = []int32{
	2, // 0: bucketeer.eventcounter.ExperimentCount.realtime_counts:type_name -> bucketeer.eventcounter.VariationCount
	2, // 1: bucketeer.eventcounter.ExperimentCount.batch_counts:type_name -> bucketeer.eventcounter.VariationCount
	1, // 2: bucketeer.eventcounter.ExperimentCount.goal_counts:type_name -> bucketeer.eventcounter.GoalCounts
	2, // 3: bucketeer.eventcounter.GoalCounts.realtime_counts:type_name -> bucketeer.eventcounter.VariationCount
	2, // 4: bucketeer.eventcounter.GoalCounts.batch_counts:type_name -> bucketeer.eventcounter.VariationCount
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_eventcounter_experiment_count_proto_init() }
func file_proto_eventcounter_experiment_count_proto_init() {
	if File_proto_eventcounter_experiment_count_proto != nil {
		return
	}
	file_proto_eventcounter_variation_count_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_eventcounter_experiment_count_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExperimentCount); i {
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
		file_proto_eventcounter_experiment_count_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoalCounts); i {
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
			RawDescriptor: file_proto_eventcounter_experiment_count_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_eventcounter_experiment_count_proto_goTypes,
		DependencyIndexes: file_proto_eventcounter_experiment_count_proto_depIdxs,
		MessageInfos:      file_proto_eventcounter_experiment_count_proto_msgTypes,
	}.Build()
	File_proto_eventcounter_experiment_count_proto = out.File
	file_proto_eventcounter_experiment_count_proto_rawDesc = nil
	file_proto_eventcounter_experiment_count_proto_goTypes = nil
	file_proto_eventcounter_experiment_count_proto_depIdxs = nil
}
