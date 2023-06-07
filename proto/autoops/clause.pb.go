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
// source: proto/autoops/clause.proto

package autoops

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OpsEventRateClause_Operator int32

const (
	OpsEventRateClause_GREATER_OR_EQUAL OpsEventRateClause_Operator = 0
	OpsEventRateClause_LESS_OR_EQUAL    OpsEventRateClause_Operator = 1
)

// Enum value maps for OpsEventRateClause_Operator.
var (
	OpsEventRateClause_Operator_name = map[int32]string{
		0: "GREATER_OR_EQUAL",
		1: "LESS_OR_EQUAL",
	}
	OpsEventRateClause_Operator_value = map[string]int32{
		"GREATER_OR_EQUAL": 0,
		"LESS_OR_EQUAL":    1,
	}
)

func (x OpsEventRateClause_Operator) Enum() *OpsEventRateClause_Operator {
	p := new(OpsEventRateClause_Operator)
	*p = x
	return p
}

func (x OpsEventRateClause_Operator) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OpsEventRateClause_Operator) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_autoops_clause_proto_enumTypes[0].Descriptor()
}

func (OpsEventRateClause_Operator) Type() protoreflect.EnumType {
	return &file_proto_autoops_clause_proto_enumTypes[0]
}

func (x OpsEventRateClause_Operator) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OpsEventRateClause_Operator.Descriptor instead.
func (OpsEventRateClause_Operator) EnumDescriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{1, 0}
}

type WebhookClause_Condition_Operator int32

const (
	WebhookClause_Condition_EQUAL              WebhookClause_Condition_Operator = 0
	WebhookClause_Condition_NOT_EQUAL          WebhookClause_Condition_Operator = 1
	WebhookClause_Condition_MORE_THAN          WebhookClause_Condition_Operator = 2
	WebhookClause_Condition_MORE_THAN_OR_EQUAL WebhookClause_Condition_Operator = 3
	WebhookClause_Condition_LESS_THAN          WebhookClause_Condition_Operator = 4
	WebhookClause_Condition_LESS_THAN_OR_EQUAL WebhookClause_Condition_Operator = 5
)

// Enum value maps for WebhookClause_Condition_Operator.
var (
	WebhookClause_Condition_Operator_name = map[int32]string{
		0: "EQUAL",
		1: "NOT_EQUAL",
		2: "MORE_THAN",
		3: "MORE_THAN_OR_EQUAL",
		4: "LESS_THAN",
		5: "LESS_THAN_OR_EQUAL",
	}
	WebhookClause_Condition_Operator_value = map[string]int32{
		"EQUAL":              0,
		"NOT_EQUAL":          1,
		"MORE_THAN":          2,
		"MORE_THAN_OR_EQUAL": 3,
		"LESS_THAN":          4,
		"LESS_THAN_OR_EQUAL": 5,
	}
)

func (x WebhookClause_Condition_Operator) Enum() *WebhookClause_Condition_Operator {
	p := new(WebhookClause_Condition_Operator)
	*p = x
	return p
}

func (x WebhookClause_Condition_Operator) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WebhookClause_Condition_Operator) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_autoops_clause_proto_enumTypes[1].Descriptor()
}

func (WebhookClause_Condition_Operator) Type() protoreflect.EnumType {
	return &file_proto_autoops_clause_proto_enumTypes[1]
}

func (x WebhookClause_Condition_Operator) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use WebhookClause_Condition_Operator.Descriptor instead.
func (WebhookClause_Condition_Operator) EnumDescriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{3, 0, 0}
}

type ProgressiveRolloutTemplateScheduleClause_Interval int32

const (
	ProgressiveRolloutTemplateScheduleClause_UNKNOWN ProgressiveRolloutTemplateScheduleClause_Interval = 0
	ProgressiveRolloutTemplateScheduleClause_HOURLY  ProgressiveRolloutTemplateScheduleClause_Interval = 1
	ProgressiveRolloutTemplateScheduleClause_DAILY   ProgressiveRolloutTemplateScheduleClause_Interval = 2
	ProgressiveRolloutTemplateScheduleClause_WEEKLY  ProgressiveRolloutTemplateScheduleClause_Interval = 3
)

// Enum value maps for ProgressiveRolloutTemplateScheduleClause_Interval.
var (
	ProgressiveRolloutTemplateScheduleClause_Interval_name = map[int32]string{
		0: "UNKNOWN",
		1: "HOURLY",
		2: "DAILY",
		3: "WEEKLY",
	}
	ProgressiveRolloutTemplateScheduleClause_Interval_value = map[string]int32{
		"UNKNOWN": 0,
		"HOURLY":  1,
		"DAILY":   2,
		"WEEKLY":  3,
	}
)

func (x ProgressiveRolloutTemplateScheduleClause_Interval) Enum() *ProgressiveRolloutTemplateScheduleClause_Interval {
	p := new(ProgressiveRolloutTemplateScheduleClause_Interval)
	*p = x
	return p
}

func (x ProgressiveRolloutTemplateScheduleClause_Interval) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProgressiveRolloutTemplateScheduleClause_Interval) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_autoops_clause_proto_enumTypes[2].Descriptor()
}

func (ProgressiveRolloutTemplateScheduleClause_Interval) Type() protoreflect.EnumType {
	return &file_proto_autoops_clause_proto_enumTypes[2]
}

func (x ProgressiveRolloutTemplateScheduleClause_Interval) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProgressiveRolloutTemplateScheduleClause_Interval.Descriptor instead.
func (ProgressiveRolloutTemplateScheduleClause_Interval) EnumDescriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{6, 0}
}

type Clause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Clause *anypb.Any `protobuf:"bytes,2,opt,name=clause,proto3" json:"clause,omitempty"`
}

func (x *Clause) Reset() {
	*x = Clause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Clause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Clause) ProtoMessage() {}

func (x *Clause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Clause.ProtoReflect.Descriptor instead.
func (*Clause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{0}
}

func (x *Clause) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Clause) GetClause() *anypb.Any {
	if x != nil {
		return x.Clause
	}
	return nil
}

type OpsEventRateClause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VariationId     string                      `protobuf:"bytes,2,opt,name=variation_id,json=variationId,proto3" json:"variation_id,omitempty"`
	GoalId          string                      `protobuf:"bytes,3,opt,name=goal_id,json=goalId,proto3" json:"goal_id,omitempty"`
	MinCount        int64                       `protobuf:"varint,4,opt,name=min_count,json=minCount,proto3" json:"min_count,omitempty"`
	ThreadsholdRate float64                     `protobuf:"fixed64,5,opt,name=threadshold_rate,json=threadsholdRate,proto3" json:"threadshold_rate,omitempty"`
	Operator        OpsEventRateClause_Operator `protobuf:"varint,6,opt,name=operator,proto3,enum=bucketeer.autoops.OpsEventRateClause_Operator" json:"operator,omitempty"`
}

func (x *OpsEventRateClause) Reset() {
	*x = OpsEventRateClause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OpsEventRateClause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OpsEventRateClause) ProtoMessage() {}

func (x *OpsEventRateClause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OpsEventRateClause.ProtoReflect.Descriptor instead.
func (*OpsEventRateClause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{1}
}

func (x *OpsEventRateClause) GetVariationId() string {
	if x != nil {
		return x.VariationId
	}
	return ""
}

func (x *OpsEventRateClause) GetGoalId() string {
	if x != nil {
		return x.GoalId
	}
	return ""
}

func (x *OpsEventRateClause) GetMinCount() int64 {
	if x != nil {
		return x.MinCount
	}
	return 0
}

func (x *OpsEventRateClause) GetThreadsholdRate() float64 {
	if x != nil {
		return x.ThreadsholdRate
	}
	return 0
}

func (x *OpsEventRateClause) GetOperator() OpsEventRateClause_Operator {
	if x != nil {
		return x.Operator
	}
	return OpsEventRateClause_GREATER_OR_EQUAL
}

type DatetimeClause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time int64 `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *DatetimeClause) Reset() {
	*x = DatetimeClause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DatetimeClause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DatetimeClause) ProtoMessage() {}

func (x *DatetimeClause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DatetimeClause.ProtoReflect.Descriptor instead.
func (*DatetimeClause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{2}
}

func (x *DatetimeClause) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

type WebhookClause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WebhookId  string                     `protobuf:"bytes,1,opt,name=webhook_id,json=webhookId,proto3" json:"webhook_id,omitempty"`
	Conditions []*WebhookClause_Condition `protobuf:"bytes,2,rep,name=conditions,proto3" json:"conditions,omitempty"`
}

func (x *WebhookClause) Reset() {
	*x = WebhookClause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookClause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookClause) ProtoMessage() {}

func (x *WebhookClause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookClause.ProtoReflect.Descriptor instead.
func (*WebhookClause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{3}
}

func (x *WebhookClause) GetWebhookId() string {
	if x != nil {
		return x.WebhookId
	}
	return ""
}

func (x *WebhookClause) GetConditions() []*WebhookClause_Condition {
	if x != nil {
		return x.Conditions
	}
	return nil
}

type ProgressiveRolloutSchedule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScheduleId  string `protobuf:"bytes,1,opt,name=schedule_id,json=scheduleId,proto3" json:"schedule_id,omitempty"`
	ExecuteAt   int64  `protobuf:"varint,2,opt,name=execute_at,json=executeAt,proto3" json:"execute_at,omitempty"`
	Weight      int32  `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty"`
	TriggeredAt int64  `protobuf:"varint,4,opt,name=triggered_at,json=triggeredAt,proto3" json:"triggered_at,omitempty"`
}

func (x *ProgressiveRolloutSchedule) Reset() {
	*x = ProgressiveRolloutSchedule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgressiveRolloutSchedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgressiveRolloutSchedule) ProtoMessage() {}

func (x *ProgressiveRolloutSchedule) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgressiveRolloutSchedule.ProtoReflect.Descriptor instead.
func (*ProgressiveRolloutSchedule) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{4}
}

func (x *ProgressiveRolloutSchedule) GetScheduleId() string {
	if x != nil {
		return x.ScheduleId
	}
	return ""
}

func (x *ProgressiveRolloutSchedule) GetExecuteAt() int64 {
	if x != nil {
		return x.ExecuteAt
	}
	return 0
}

func (x *ProgressiveRolloutSchedule) GetWeight() int32 {
	if x != nil {
		return x.Weight
	}
	return 0
}

func (x *ProgressiveRolloutSchedule) GetTriggeredAt() int64 {
	if x != nil {
		return x.TriggeredAt
	}
	return 0
}

type ProgressiveRolloutManualScheduleClause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Schedules   []*ProgressiveRolloutSchedule `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules,omitempty"`
	VariationId string                        `protobuf:"bytes,2,opt,name=variation_id,json=variationId,proto3" json:"variation_id,omitempty"`
}

func (x *ProgressiveRolloutManualScheduleClause) Reset() {
	*x = ProgressiveRolloutManualScheduleClause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgressiveRolloutManualScheduleClause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgressiveRolloutManualScheduleClause) ProtoMessage() {}

func (x *ProgressiveRolloutManualScheduleClause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgressiveRolloutManualScheduleClause.ProtoReflect.Descriptor instead.
func (*ProgressiveRolloutManualScheduleClause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{5}
}

func (x *ProgressiveRolloutManualScheduleClause) GetSchedules() []*ProgressiveRolloutSchedule {
	if x != nil {
		return x.Schedules
	}
	return nil
}

func (x *ProgressiveRolloutManualScheduleClause) GetVariationId() string {
	if x != nil {
		return x.VariationId
	}
	return ""
}

type ProgressiveRolloutTemplateScheduleClause struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The reason of setting `schedules` is to save `triggered_at` in each
	// schedule.
	Schedules   []*ProgressiveRolloutSchedule                     `protobuf:"bytes,1,rep,name=schedules,proto3" json:"schedules,omitempty"`
	Interval    ProgressiveRolloutTemplateScheduleClause_Interval `protobuf:"varint,2,opt,name=interval,proto3,enum=bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause_Interval" json:"interval,omitempty"`
	Increments  int64                                             `protobuf:"varint,3,opt,name=increments,proto3" json:"increments,omitempty"`
	VariationId string                                            `protobuf:"bytes,4,opt,name=variation_id,json=variationId,proto3" json:"variation_id,omitempty"`
}

func (x *ProgressiveRolloutTemplateScheduleClause) Reset() {
	*x = ProgressiveRolloutTemplateScheduleClause{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProgressiveRolloutTemplateScheduleClause) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProgressiveRolloutTemplateScheduleClause) ProtoMessage() {}

func (x *ProgressiveRolloutTemplateScheduleClause) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProgressiveRolloutTemplateScheduleClause.ProtoReflect.Descriptor instead.
func (*ProgressiveRolloutTemplateScheduleClause) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{6}
}

func (x *ProgressiveRolloutTemplateScheduleClause) GetSchedules() []*ProgressiveRolloutSchedule {
	if x != nil {
		return x.Schedules
	}
	return nil
}

func (x *ProgressiveRolloutTemplateScheduleClause) GetInterval() ProgressiveRolloutTemplateScheduleClause_Interval {
	if x != nil {
		return x.Interval
	}
	return ProgressiveRolloutTemplateScheduleClause_UNKNOWN
}

func (x *ProgressiveRolloutTemplateScheduleClause) GetIncrements() int64 {
	if x != nil {
		return x.Increments
	}
	return 0
}

func (x *ProgressiveRolloutTemplateScheduleClause) GetVariationId() string {
	if x != nil {
		return x.VariationId
	}
	return ""
}

type WebhookClause_Condition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filter   string                           `protobuf:"bytes,1,opt,name=filter,proto3" json:"filter,omitempty"`
	Value    string                           `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Operator WebhookClause_Condition_Operator `protobuf:"varint,3,opt,name=operator,proto3,enum=bucketeer.autoops.WebhookClause_Condition_Operator" json:"operator,omitempty"`
}

func (x *WebhookClause_Condition) Reset() {
	*x = WebhookClause_Condition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_autoops_clause_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookClause_Condition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookClause_Condition) ProtoMessage() {}

func (x *WebhookClause_Condition) ProtoReflect() protoreflect.Message {
	mi := &file_proto_autoops_clause_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookClause_Condition.ProtoReflect.Descriptor instead.
func (*WebhookClause_Condition) Descriptor() ([]byte, []int) {
	return file_proto_autoops_clause_proto_rawDescGZIP(), []int{3, 0}
}

func (x *WebhookClause_Condition) GetFilter() string {
	if x != nil {
		return x.Filter
	}
	return ""
}

func (x *WebhookClause_Condition) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

func (x *WebhookClause_Condition) GetOperator() WebhookClause_Condition_Operator {
	if x != nil {
		return x.Operator
	}
	return WebhookClause_Condition_EQUAL
}

var File_proto_autoops_clause_proto protoreflect.FileDescriptor

var file_proto_autoops_clause_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70, 0x73, 0x2f,
	0x63, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70, 0x73, 0x1a,
	0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x06, 0x43, 0x6c,
	0x61, 0x75, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x2c, 0x0a, 0x06, 0x63, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x06, 0x63, 0x6c, 0x61, 0x75,
	0x73, 0x65, 0x22, 0x9f, 0x02, 0x0a, 0x12, 0x4f, 0x70, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52,
	0x61, 0x74, 0x65, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61, 0x72,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07,
	0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x67,
	0x6f, 0x61, 0x6c, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6d, 0x69, 0x6e, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x69, 0x6e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x73, 0x68, 0x6f, 0x6c,
	0x64, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x74, 0x68,
	0x72, 0x65, 0x61, 0x64, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x52, 0x61, 0x74, 0x65, 0x12, 0x4a, 0x0a,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x2e, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x6f, 0x70, 0x73, 0x2e, 0x4f, 0x70, 0x73, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x61, 0x74, 0x65,
	0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x22, 0x33, 0x0a, 0x08, 0x4f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x52, 0x45, 0x41, 0x54, 0x45, 0x52,
	0x5f, 0x4f, 0x52, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4c,
	0x45, 0x53, 0x53, 0x5f, 0x4f, 0x52, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x01, 0x4a, 0x04,
	0x08, 0x01, 0x10, 0x02, 0x22, 0x24, 0x0a, 0x0e, 0x44, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65,
	0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xfb, 0x02, 0x0a, 0x0d, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x49, 0x64, 0x12, 0x4a, 0x0a, 0x0a, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f,
	0x6f, 0x70, 0x73, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x43, 0x6c, 0x61, 0x75, 0x73,
	0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6e,
	0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0xfe, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x4f, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x33, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70, 0x73, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72,
	0x61, 0x74, 0x6f, 0x72, 0x22, 0x72, 0x0a, 0x08, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x12, 0x09, 0x0a, 0x05, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4e,
	0x4f, 0x54, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x4d, 0x4f,
	0x52, 0x45, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x4d, 0x4f, 0x52,
	0x45, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x5f, 0x4f, 0x52, 0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10,
	0x03, 0x12, 0x0d, 0x0a, 0x09, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x10, 0x04,
	0x12, 0x16, 0x0a, 0x12, 0x4c, 0x45, 0x53, 0x53, 0x5f, 0x54, 0x48, 0x41, 0x4e, 0x5f, 0x4f, 0x52,
	0x5f, 0x45, 0x51, 0x55, 0x41, 0x4c, 0x10, 0x05, 0x22, 0x97, 0x01, 0x0a, 0x1a, 0x50, 0x72, 0x6f,
	0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x65, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x65, 0x41, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x12,
	0x21, 0x0a, 0x0c, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x65, 0x64,
	0x41, 0x74, 0x22, 0x98, 0x01, 0x0a, 0x26, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69,
	0x76, 0x65, 0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x4d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x12, 0x4b, 0x0a,
	0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2d, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74,
	0x6f, 0x6f, 0x70, 0x73, 0x2e, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65,
	0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52,
	0x09, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xd8, 0x02,
	0x0a, 0x28, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x52, 0x6f, 0x6c,
	0x6c, 0x6f, 0x75, 0x74, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65,
	0x64, 0x75, 0x6c, 0x65, 0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x12, 0x4b, 0x0a, 0x09, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2d, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70,
	0x73, 0x2e, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x52, 0x6f, 0x6c,
	0x6c, 0x6f, 0x75, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x09, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x60, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x44, 0x2e, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70, 0x73, 0x2e, 0x50, 0x72,
	0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x69, 0x76, 0x65, 0x52, 0x6f, 0x6c, 0x6c, 0x6f, 0x75, 0x74,
	0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65,
	0x43, 0x6c, 0x61, 0x75, 0x73, 0x65, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x52,
	0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x63,
	0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x69,
	0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61, 0x72,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x08,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x4f, 0x55, 0x52, 0x4c, 0x59, 0x10,
	0x01, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x41, 0x49, 0x4c, 0x59, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06,
	0x57, 0x45, 0x45, 0x4b, 0x4c, 0x59, 0x10, 0x03, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x75, 0x74, 0x6f, 0x6f, 0x70, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_autoops_clause_proto_rawDescOnce sync.Once
	file_proto_autoops_clause_proto_rawDescData = file_proto_autoops_clause_proto_rawDesc
)

func file_proto_autoops_clause_proto_rawDescGZIP() []byte {
	file_proto_autoops_clause_proto_rawDescOnce.Do(func() {
		file_proto_autoops_clause_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_autoops_clause_proto_rawDescData)
	})
	return file_proto_autoops_clause_proto_rawDescData
}

var file_proto_autoops_clause_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_proto_autoops_clause_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_autoops_clause_proto_goTypes = []interface{}{
	(OpsEventRateClause_Operator)(0),                       // 0: bucketeer.autoops.OpsEventRateClause.Operator
	(WebhookClause_Condition_Operator)(0),                  // 1: bucketeer.autoops.WebhookClause.Condition.Operator
	(ProgressiveRolloutTemplateScheduleClause_Interval)(0), // 2: bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause.Interval
	(*Clause)(nil),                                   // 3: bucketeer.autoops.Clause
	(*OpsEventRateClause)(nil),                       // 4: bucketeer.autoops.OpsEventRateClause
	(*DatetimeClause)(nil),                           // 5: bucketeer.autoops.DatetimeClause
	(*WebhookClause)(nil),                            // 6: bucketeer.autoops.WebhookClause
	(*ProgressiveRolloutSchedule)(nil),               // 7: bucketeer.autoops.ProgressiveRolloutSchedule
	(*ProgressiveRolloutManualScheduleClause)(nil),   // 8: bucketeer.autoops.ProgressiveRolloutManualScheduleClause
	(*ProgressiveRolloutTemplateScheduleClause)(nil), // 9: bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause
	(*WebhookClause_Condition)(nil),                  // 10: bucketeer.autoops.WebhookClause.Condition
	(*anypb.Any)(nil),                                // 11: google.protobuf.Any
}
var file_proto_autoops_clause_proto_depIdxs = []int32{
	11, // 0: bucketeer.autoops.Clause.clause:type_name -> google.protobuf.Any
	0,  // 1: bucketeer.autoops.OpsEventRateClause.operator:type_name -> bucketeer.autoops.OpsEventRateClause.Operator
	10, // 2: bucketeer.autoops.WebhookClause.conditions:type_name -> bucketeer.autoops.WebhookClause.Condition
	7,  // 3: bucketeer.autoops.ProgressiveRolloutManualScheduleClause.schedules:type_name -> bucketeer.autoops.ProgressiveRolloutSchedule
	7,  // 4: bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause.schedules:type_name -> bucketeer.autoops.ProgressiveRolloutSchedule
	2,  // 5: bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause.interval:type_name -> bucketeer.autoops.ProgressiveRolloutTemplateScheduleClause.Interval
	1,  // 6: bucketeer.autoops.WebhookClause.Condition.operator:type_name -> bucketeer.autoops.WebhookClause.Condition.Operator
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_proto_autoops_clause_proto_init() }
func file_proto_autoops_clause_proto_init() {
	if File_proto_autoops_clause_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_autoops_clause_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Clause); i {
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
		file_proto_autoops_clause_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OpsEventRateClause); i {
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
		file_proto_autoops_clause_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DatetimeClause); i {
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
		file_proto_autoops_clause_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookClause); i {
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
		file_proto_autoops_clause_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgressiveRolloutSchedule); i {
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
		file_proto_autoops_clause_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgressiveRolloutManualScheduleClause); i {
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
		file_proto_autoops_clause_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProgressiveRolloutTemplateScheduleClause); i {
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
		file_proto_autoops_clause_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookClause_Condition); i {
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
			RawDescriptor: file_proto_autoops_clause_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_autoops_clause_proto_goTypes,
		DependencyIndexes: file_proto_autoops_clause_proto_depIdxs,
		EnumInfos:         file_proto_autoops_clause_proto_enumTypes,
		MessageInfos:      file_proto_autoops_clause_proto_msgTypes,
	}.Build()
	File_proto_autoops_clause_proto = out.File
	file_proto_autoops_clause_proto_rawDesc = nil
	file_proto_autoops_clause_proto_goTypes = nil
	file_proto_autoops_clause_proto_depIdxs = nil
}
