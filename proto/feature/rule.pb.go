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
// source: proto/feature/rule.proto

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

type Rule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Strategy *Strategy `protobuf:"bytes,2,opt,name=strategy,proto3" json:"strategy"`
	Clauses  []*Clause `protobuf:"bytes,3,rep,name=clauses,proto3" json:"clauses"`
}

func (x *Rule) Reset() {
	*x = Rule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_rule_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Rule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Rule) ProtoMessage() {}

func (x *Rule) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_rule_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Rule.ProtoReflect.Descriptor instead.
func (*Rule) Descriptor() ([]byte, []int) {
	return file_proto_feature_rule_proto_rawDescGZIP(), []int{0}
}

func (x *Rule) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Rule) GetStrategy() *Strategy {
	if x != nil {
		return x.Strategy
	}
	return nil
}

func (x *Rule) GetClauses() []*Clause {
	if x != nil {
		return x.Clauses
	}
	return nil
}

type RuleListValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values []*Rule `protobuf:"bytes,1,rep,name=values,proto3" json:"values"`
}

func (x *RuleListValue) Reset() {
	*x = RuleListValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_feature_rule_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RuleListValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RuleListValue) ProtoMessage() {}

func (x *RuleListValue) ProtoReflect() protoreflect.Message {
	mi := &file_proto_feature_rule_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RuleListValue.ProtoReflect.Descriptor instead.
func (*RuleListValue) Descriptor() ([]byte, []int) {
	return file_proto_feature_rule_proto_rawDescGZIP(), []int{1}
}

func (x *RuleListValue) GetValues() []*Rule {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_proto_feature_rule_proto protoreflect.FileDescriptor

var file_proto_feature_rule_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f,
	0x72, 0x75, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x1a, 0x1a, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x63, 0x6c, 0x61,
	0x75, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x84, 0x01, 0x0a, 0x04, 0x52, 0x75, 0x6c, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x37, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66,
	0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x52,
	0x08, 0x73, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x33, 0x0a, 0x07, 0x63, 0x6c, 0x61,
	0x75, 0x73, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74, 0x75, 0x72, 0x65, 0x2e, 0x43,
	0x6c, 0x61, 0x75, 0x73, 0x65, 0x52, 0x07, 0x63, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x73, 0x22, 0x40,
	0x0a, 0x0d, 0x52, 0x75, 0x6c, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x2e, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x65, 0x61, 0x74,
	0x75, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_feature_rule_proto_rawDescOnce sync.Once
	file_proto_feature_rule_proto_rawDescData = file_proto_feature_rule_proto_rawDesc
)

func file_proto_feature_rule_proto_rawDescGZIP() []byte {
	file_proto_feature_rule_proto_rawDescOnce.Do(func() {
		file_proto_feature_rule_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_feature_rule_proto_rawDescData)
	})
	return file_proto_feature_rule_proto_rawDescData
}

var file_proto_feature_rule_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_feature_rule_proto_goTypes = []interface{}{
	(*Rule)(nil),          // 0: bucketeer.feature.Rule
	(*RuleListValue)(nil), // 1: bucketeer.feature.RuleListValue
	(*Strategy)(nil),      // 2: bucketeer.feature.Strategy
	(*Clause)(nil),        // 3: bucketeer.feature.Clause
}
var file_proto_feature_rule_proto_depIdxs = []int32{
	2, // 0: bucketeer.feature.Rule.strategy:type_name -> bucketeer.feature.Strategy
	3, // 1: bucketeer.feature.Rule.clauses:type_name -> bucketeer.feature.Clause
	0, // 2: bucketeer.feature.RuleListValue.values:type_name -> bucketeer.feature.Rule
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_feature_rule_proto_init() }
func file_proto_feature_rule_proto_init() {
	if File_proto_feature_rule_proto != nil {
		return
	}
	file_proto_feature_clause_proto_init()
	file_proto_feature_strategy_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_feature_rule_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Rule); i {
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
		file_proto_feature_rule_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RuleListValue); i {
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
			RawDescriptor: file_proto_feature_rule_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_feature_rule_proto_goTypes,
		DependencyIndexes: file_proto_feature_rule_proto_depIdxs,
		MessageInfos:      file_proto_feature_rule_proto_msgTypes,
	}.Build()
	File_proto_feature_rule_proto = out.File
	file_proto_feature_rule_proto_rawDesc = nil
	file_proto_feature_rule_proto_goTypes = nil
	file_proto_feature_rule_proto_depIdxs = nil
}
