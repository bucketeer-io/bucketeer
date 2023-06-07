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
// source: proto/eventcounter/variation_result.proto

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

type VariationResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VariationId                                string               `protobuf:"bytes,1,opt,name=variation_id,json=variationId,proto3" json:"variation_id,omitempty"`
	ExperimentCount                            *VariationCount      `protobuf:"bytes,2,opt,name=experiment_count,json=experimentCount,proto3" json:"experiment_count,omitempty"`
	EvaluationCount                            *VariationCount      `protobuf:"bytes,3,opt,name=evaluation_count,json=evaluationCount,proto3" json:"evaluation_count,omitempty"`
	CvrProbBest                                *DistributionSummary `protobuf:"bytes,4,opt,name=cvr_prob_best,json=cvrProbBest,proto3" json:"cvr_prob_best,omitempty"`
	CvrProbBeatBaseline                        *DistributionSummary `protobuf:"bytes,5,opt,name=cvr_prob_beat_baseline,json=cvrProbBeatBaseline,proto3" json:"cvr_prob_beat_baseline,omitempty"`
	CvrProb                                    *DistributionSummary `protobuf:"bytes,6,opt,name=cvr_prob,json=cvrProb,proto3" json:"cvr_prob,omitempty"`
	EvaluationUserCountTimeseries              *Timeseries          `protobuf:"bytes,7,opt,name=evaluation_user_count_timeseries,json=evaluationUserCountTimeseries,proto3" json:"evaluation_user_count_timeseries,omitempty"`
	EvaluationEventCountTimeseries             *Timeseries          `protobuf:"bytes,8,opt,name=evaluation_event_count_timeseries,json=evaluationEventCountTimeseries,proto3" json:"evaluation_event_count_timeseries,omitempty"`
	GoalUserCountTimeseries                    *Timeseries          `protobuf:"bytes,9,opt,name=goal_user_count_timeseries,json=goalUserCountTimeseries,proto3" json:"goal_user_count_timeseries,omitempty"`
	GoalEventCountTimeseries                   *Timeseries          `protobuf:"bytes,10,opt,name=goal_event_count_timeseries,json=goalEventCountTimeseries,proto3" json:"goal_event_count_timeseries,omitempty"`
	GoalValueSumTimeseries                     *Timeseries          `protobuf:"bytes,11,opt,name=goal_value_sum_timeseries,json=goalValueSumTimeseries,proto3" json:"goal_value_sum_timeseries,omitempty"`
	CvrMedianTimeseries                        *Timeseries          `protobuf:"bytes,12,opt,name=cvr_median_timeseries,json=cvrMedianTimeseries,proto3" json:"cvr_median_timeseries,omitempty"`
	CvrPercentile025Timeseries                 *Timeseries          `protobuf:"bytes,13,opt,name=cvr_percentile025_timeseries,json=cvrPercentile025Timeseries,proto3" json:"cvr_percentile025_timeseries,omitempty"`
	CvrPercentile975Timeseries                 *Timeseries          `protobuf:"bytes,14,opt,name=cvr_percentile975_timeseries,json=cvrPercentile975Timeseries,proto3" json:"cvr_percentile975_timeseries,omitempty"`
	CvrTimeseries                              *Timeseries          `protobuf:"bytes,15,opt,name=cvr_timeseries,json=cvrTimeseries,proto3" json:"cvr_timeseries,omitempty"`
	GoalValueSumPerUserTimeseries              *Timeseries          `protobuf:"bytes,16,opt,name=goal_value_sum_per_user_timeseries,json=goalValueSumPerUserTimeseries,proto3" json:"goal_value_sum_per_user_timeseries,omitempty"`
	GoalValueSumPerUserProb                    *DistributionSummary `protobuf:"bytes,17,opt,name=goal_value_sum_per_user_prob,json=goalValueSumPerUserProb,proto3" json:"goal_value_sum_per_user_prob,omitempty"`
	GoalValueSumPerUserProbBest                *DistributionSummary `protobuf:"bytes,18,opt,name=goal_value_sum_per_user_prob_best,json=goalValueSumPerUserProbBest,proto3" json:"goal_value_sum_per_user_prob_best,omitempty"`
	GoalValueSumPerUserProbBeatBaseline        *DistributionSummary `protobuf:"bytes,19,opt,name=goal_value_sum_per_user_prob_beat_baseline,json=goalValueSumPerUserProbBeatBaseline,proto3" json:"goal_value_sum_per_user_prob_beat_baseline,omitempty"`
	GoalValueSumPerUserMedianTimeseries        *Timeseries          `protobuf:"bytes,20,opt,name=goal_value_sum_per_user_median_timeseries,json=goalValueSumPerUserMedianTimeseries,proto3" json:"goal_value_sum_per_user_median_timeseries,omitempty"`
	GoalValueSumPerUserPercentile025Timeseries *Timeseries          `protobuf:"bytes,21,opt,name=goal_value_sum_per_user_percentile025_timeseries,json=goalValueSumPerUserPercentile025Timeseries,proto3" json:"goal_value_sum_per_user_percentile025_timeseries,omitempty"`
	GoalValueSumPerUserPercentile975Timeseries *Timeseries          `protobuf:"bytes,22,opt,name=goal_value_sum_per_user_percentile975_timeseries,json=goalValueSumPerUserPercentile975Timeseries,proto3" json:"goal_value_sum_per_user_percentile975_timeseries,omitempty"`
}

func (x *VariationResult) Reset() {
	*x = VariationResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_eventcounter_variation_result_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VariationResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VariationResult) ProtoMessage() {}

func (x *VariationResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_eventcounter_variation_result_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VariationResult.ProtoReflect.Descriptor instead.
func (*VariationResult) Descriptor() ([]byte, []int) {
	return file_proto_eventcounter_variation_result_proto_rawDescGZIP(), []int{0}
}

func (x *VariationResult) GetVariationId() string {
	if x != nil {
		return x.VariationId
	}
	return ""
}

func (x *VariationResult) GetExperimentCount() *VariationCount {
	if x != nil {
		return x.ExperimentCount
	}
	return nil
}

func (x *VariationResult) GetEvaluationCount() *VariationCount {
	if x != nil {
		return x.EvaluationCount
	}
	return nil
}

func (x *VariationResult) GetCvrProbBest() *DistributionSummary {
	if x != nil {
		return x.CvrProbBest
	}
	return nil
}

func (x *VariationResult) GetCvrProbBeatBaseline() *DistributionSummary {
	if x != nil {
		return x.CvrProbBeatBaseline
	}
	return nil
}

func (x *VariationResult) GetCvrProb() *DistributionSummary {
	if x != nil {
		return x.CvrProb
	}
	return nil
}

func (x *VariationResult) GetEvaluationUserCountTimeseries() *Timeseries {
	if x != nil {
		return x.EvaluationUserCountTimeseries
	}
	return nil
}

func (x *VariationResult) GetEvaluationEventCountTimeseries() *Timeseries {
	if x != nil {
		return x.EvaluationEventCountTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalUserCountTimeseries() *Timeseries {
	if x != nil {
		return x.GoalUserCountTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalEventCountTimeseries() *Timeseries {
	if x != nil {
		return x.GoalEventCountTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumTimeseries() *Timeseries {
	if x != nil {
		return x.GoalValueSumTimeseries
	}
	return nil
}

func (x *VariationResult) GetCvrMedianTimeseries() *Timeseries {
	if x != nil {
		return x.CvrMedianTimeseries
	}
	return nil
}

func (x *VariationResult) GetCvrPercentile025Timeseries() *Timeseries {
	if x != nil {
		return x.CvrPercentile025Timeseries
	}
	return nil
}

func (x *VariationResult) GetCvrPercentile975Timeseries() *Timeseries {
	if x != nil {
		return x.CvrPercentile975Timeseries
	}
	return nil
}

func (x *VariationResult) GetCvrTimeseries() *Timeseries {
	if x != nil {
		return x.CvrTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserTimeseries() *Timeseries {
	if x != nil {
		return x.GoalValueSumPerUserTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserProb() *DistributionSummary {
	if x != nil {
		return x.GoalValueSumPerUserProb
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserProbBest() *DistributionSummary {
	if x != nil {
		return x.GoalValueSumPerUserProbBest
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserProbBeatBaseline() *DistributionSummary {
	if x != nil {
		return x.GoalValueSumPerUserProbBeatBaseline
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserMedianTimeseries() *Timeseries {
	if x != nil {
		return x.GoalValueSumPerUserMedianTimeseries
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserPercentile025Timeseries() *Timeseries {
	if x != nil {
		return x.GoalValueSumPerUserPercentile025Timeseries
	}
	return nil
}

func (x *VariationResult) GetGoalValueSumPerUserPercentile975Timeseries() *Timeseries {
	if x != nil {
		return x.GoalValueSumPerUserPercentile975Timeseries
	}
	return nil
}

var File_proto_eventcounter_variation_result_proto protoreflect.FileDescriptor

var file_proto_eventcounter_variation_result_proto_rawDesc = []byte{
	0x0a, 0x29, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x75, 0x63,
	0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x65, 0x72, 0x1a, 0x28, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2d, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65,
	0x72, 0x2f, 0x64, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73,
	0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xac, 0x11, 0x0a, 0x0f, 0x56, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61, 0x72, 0x69, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x76, 0x61, 0x72,
	0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x51, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x65,
	0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x0f, 0x65, 0x78, 0x70, 0x65,
	0x72, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x51, 0x0a, 0x10, 0x65,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x0f, 0x65,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x4f,
	0x0a, 0x0d, 0x63, 0x76, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x62, 0x5f, 0x62, 0x65, 0x73, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44,
	0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61,
	0x72, 0x79, 0x52, 0x0b, 0x63, 0x76, 0x72, 0x50, 0x72, 0x6f, 0x62, 0x42, 0x65, 0x73, 0x74, 0x12,
	0x60, 0x0a, 0x16, 0x63, 0x76, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x62, 0x5f, 0x62, 0x65, 0x61, 0x74,
	0x5f, 0x62, 0x61, 0x73, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2b, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x13, 0x63, 0x76,
	0x72, 0x50, 0x72, 0x6f, 0x62, 0x42, 0x65, 0x61, 0x74, 0x42, 0x61, 0x73, 0x65, 0x6c, 0x69, 0x6e,
	0x65, 0x12, 0x46, 0x0a, 0x08, 0x63, 0x76, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x62, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x73,
	0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x52, 0x07, 0x63, 0x76, 0x72, 0x50, 0x72, 0x6f, 0x62, 0x12, 0x6b, 0x0a, 0x20, 0x65, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x1d, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x6d, 0x0a, 0x21, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x1e, 0x65, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x5f, 0x0a, 0x1a, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72,
	0x69, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x17, 0x67,
	0x6f, 0x61, 0x6c, 0x55, 0x73, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x61, 0x0a, 0x1b, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75,
	0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x18, 0x67, 0x6f, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x5d, 0x0a, 0x19, 0x67, 0x6f, 0x61,
	0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x16, 0x67, 0x6f, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x75, 0x6d, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x56, 0x0a, 0x15, 0x63, 0x76, 0x72, 0x5f,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65,
	0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74,
	0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x13, 0x63, 0x76, 0x72,
	0x4d, 0x65, 0x64, 0x69, 0x61, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x12, 0x64, 0x0a, 0x1c, 0x63, 0x76, 0x72, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x69,
	0x6c, 0x65, 0x30, 0x32, 0x35, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65,
	0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x1a, 0x63, 0x76, 0x72, 0x50,
	0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c, 0x65, 0x30, 0x32, 0x35, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x64, 0x0a, 0x1c, 0x63, 0x76, 0x72, 0x5f, 0x70, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c, 0x65, 0x39, 0x37, 0x35, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x1a, 0x63, 0x76, 0x72, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c, 0x65, 0x39,
	0x37, 0x35, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x49, 0x0a, 0x0e,
	0x63, 0x76, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72,
	0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x0d, 0x63, 0x76, 0x72, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x6d, 0x0a, 0x22, 0x67, 0x6f, 0x61, 0x6c, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75, 0x73,
	0x65, 0x72, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x10, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x1d, 0x67, 0x6f, 0x61, 0x6c, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x6a, 0x0a, 0x1c, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x70, 0x72, 0x6f, 0x62, 0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x17, 0x67, 0x6f, 0x61, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x62, 0x12, 0x73, 0x0a, 0x21, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x72,
	0x6f, 0x62, 0x5f, 0x62, 0x65, 0x73, 0x74, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e,
	0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x1b, 0x67, 0x6f, 0x61, 0x6c,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x50,
	0x72, 0x6f, 0x62, 0x42, 0x65, 0x73, 0x74, 0x12, 0x84, 0x01, 0x0a, 0x2a, 0x67, 0x6f, 0x61, 0x6c,
	0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x62, 0x5f, 0x62, 0x65, 0x61, 0x74, 0x5f, 0x62, 0x61,
	0x73, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x44, 0x69, 0x73, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x52, 0x23, 0x67, 0x6f, 0x61, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72,
	0x6f, 0x62, 0x42, 0x65, 0x61, 0x74, 0x42, 0x61, 0x73, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x7a,
	0x0a, 0x29, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d,
	0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x6e,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x23, 0x67, 0x6f, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x88, 0x01, 0x0a, 0x30, 0x67,
	0x6f, 0x61, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65,
	0x72, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c,
	0x65, 0x30, 0x32, 0x35, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18,
	0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65,
	0x72, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52, 0x2a, 0x67, 0x6f, 0x61, 0x6c, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x53, 0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x50, 0x65,
	0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c, 0x65, 0x30, 0x32, 0x35, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x88, 0x01, 0x0a, 0x30, 0x67, 0x6f, 0x61, 0x6c, 0x5f, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x5f, 0x73, 0x75, 0x6d, 0x5f, 0x70, 0x65, 0x72, 0x5f, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x69, 0x6c, 0x65, 0x39, 0x37, 0x35, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x2a, 0x67, 0x6f, 0x61, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x53,
	0x75, 0x6d, 0x50, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74,
	0x69, 0x6c, 0x65, 0x39, 0x37, 0x35, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73,
	0x42, 0x36, 0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x65, 0x65, 0x72, 0x2d, 0x69, 0x6f, 0x2f, 0x62, 0x75, 0x63, 0x6b,
	0x65, 0x74, 0x65, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_eventcounter_variation_result_proto_rawDescOnce sync.Once
	file_proto_eventcounter_variation_result_proto_rawDescData = file_proto_eventcounter_variation_result_proto_rawDesc
)

func file_proto_eventcounter_variation_result_proto_rawDescGZIP() []byte {
	file_proto_eventcounter_variation_result_proto_rawDescOnce.Do(func() {
		file_proto_eventcounter_variation_result_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_eventcounter_variation_result_proto_rawDescData)
	})
	return file_proto_eventcounter_variation_result_proto_rawDescData
}

var file_proto_eventcounter_variation_result_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_proto_eventcounter_variation_result_proto_goTypes = []interface{}{
	(*VariationResult)(nil),     // 0: bucketeer.eventcounter.VariationResult
	(*VariationCount)(nil),      // 1: bucketeer.eventcounter.VariationCount
	(*DistributionSummary)(nil), // 2: bucketeer.eventcounter.DistributionSummary
	(*Timeseries)(nil),          // 3: bucketeer.eventcounter.Timeseries
}
var file_proto_eventcounter_variation_result_proto_depIdxs = []int32{
	1,  // 0: bucketeer.eventcounter.VariationResult.experiment_count:type_name -> bucketeer.eventcounter.VariationCount
	1,  // 1: bucketeer.eventcounter.VariationResult.evaluation_count:type_name -> bucketeer.eventcounter.VariationCount
	2,  // 2: bucketeer.eventcounter.VariationResult.cvr_prob_best:type_name -> bucketeer.eventcounter.DistributionSummary
	2,  // 3: bucketeer.eventcounter.VariationResult.cvr_prob_beat_baseline:type_name -> bucketeer.eventcounter.DistributionSummary
	2,  // 4: bucketeer.eventcounter.VariationResult.cvr_prob:type_name -> bucketeer.eventcounter.DistributionSummary
	3,  // 5: bucketeer.eventcounter.VariationResult.evaluation_user_count_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 6: bucketeer.eventcounter.VariationResult.evaluation_event_count_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 7: bucketeer.eventcounter.VariationResult.goal_user_count_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 8: bucketeer.eventcounter.VariationResult.goal_event_count_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 9: bucketeer.eventcounter.VariationResult.goal_value_sum_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 10: bucketeer.eventcounter.VariationResult.cvr_median_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 11: bucketeer.eventcounter.VariationResult.cvr_percentile025_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 12: bucketeer.eventcounter.VariationResult.cvr_percentile975_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 13: bucketeer.eventcounter.VariationResult.cvr_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 14: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	2,  // 15: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_prob:type_name -> bucketeer.eventcounter.DistributionSummary
	2,  // 16: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_prob_best:type_name -> bucketeer.eventcounter.DistributionSummary
	2,  // 17: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_prob_beat_baseline:type_name -> bucketeer.eventcounter.DistributionSummary
	3,  // 18: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_median_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 19: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_percentile025_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	3,  // 20: bucketeer.eventcounter.VariationResult.goal_value_sum_per_user_percentile975_timeseries:type_name -> bucketeer.eventcounter.Timeseries
	21, // [21:21] is the sub-list for method output_type
	21, // [21:21] is the sub-list for method input_type
	21, // [21:21] is the sub-list for extension type_name
	21, // [21:21] is the sub-list for extension extendee
	0,  // [0:21] is the sub-list for field type_name
}

func init() { file_proto_eventcounter_variation_result_proto_init() }
func file_proto_eventcounter_variation_result_proto_init() {
	if File_proto_eventcounter_variation_result_proto != nil {
		return
	}
	file_proto_eventcounter_variation_count_proto_init()
	file_proto_eventcounter_distribution_summary_proto_init()
	file_proto_eventcounter_timeseries_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_eventcounter_variation_result_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VariationResult); i {
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
			RawDescriptor: file_proto_eventcounter_variation_result_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_eventcounter_variation_result_proto_goTypes,
		DependencyIndexes: file_proto_eventcounter_variation_result_proto_depIdxs,
		MessageInfos:      file_proto_eventcounter_variation_result_proto_msgTypes,
	}.Build()
	File_proto_eventcounter_variation_result_proto = out.File
	file_proto_eventcounter_variation_result_proto_rawDesc = nil
	file_proto_eventcounter_variation_result_proto_goTypes = nil
	file_proto_eventcounter_variation_result_proto_depIdxs = nil
}
