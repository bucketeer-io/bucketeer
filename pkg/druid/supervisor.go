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

package druid

import (
	"fmt"

	"github.com/ca-dp/godruid"

	storagedruid "github.com/bucketeer-io/bucketeer/pkg/storage/druid"
	storagekafka "github.com/bucketeer-io/bucketeer/pkg/storage/kafka"
)

type supervisor struct {
	dataType           string
	kafkaTopicDataType string
	metricsSpec        []*godruid.MetricsSpec
}

var (
	eventSupervisors = []supervisor{
		{
			dataType:           "evaluation_events",
			kafkaTopicDataType: "evaluation-events",
			metricsSpec: []*godruid.MetricsSpec{
				{
					Type: "count",
					Name: "count",
				},
				{
					Name:      "userIdHllSketch",
					Type:      "HLLSketchBuild",
					FieldName: "metric.userId",
				},
				{
					Name:      "userIdThetaSketch",
					Type:      "thetaSketch",
					FieldName: "metric.userId",
				},
			},
		},
		{
			dataType:           "goal_events",
			kafkaTopicDataType: "goal-events",
			metricsSpec: []*godruid.MetricsSpec{
				{
					Name: "count",
					Type: "count",
				},
				{
					Name:      "valueSum",
					Type:      "doubleSum",
					FieldName: "value",
				},
				{
					Name:      "valueVariance",
					Type:      "variance",
					FieldName: "value",
					InputType: "double",
					Estimator: "population",
				},
				{
					Name:      "userIdHllSketch",
					Type:      "HLLSketchBuild",
					FieldName: "metric.userId",
				},
				{
					Name:      "userIdThetaSketch",
					Type:      "thetaSketch",
					FieldName: "metric.userId",
				},
			},
		},
		{
			dataType:           "user_events",
			kafkaTopicDataType: "user-events",
		},
	}
)

func EventSupervisors(
	druidDatasourcePrefix,
	kafkaTopicPrefix,
	kafkaURL,
	kafkaUsername,
	kafkaPassword string,
	maxRowsPerSegment int,
) []*godruid.SupervisorKafka {
	supervisors := []*godruid.SupervisorKafka{}
	for _, es := range eventSupervisors {
		kafkaTopic := storagekafka.TopicName(kafkaTopicPrefix, es.kafkaTopicDataType)
		datasource := storagedruid.Datasource(druidDatasourcePrefix, es.dataType)
		supervisors = append(
			supervisors,
			eventSupervisor(
				kafkaTopic,
				kafkaURL,
				kafkaUsername,
				kafkaPassword,
				datasource,
				maxRowsPerSegment,
				es.metricsSpec,
			),
		)
	}
	return supervisors
}

func eventSupervisor(
	kafkaTopic, kafkaURL, kafkaUsername, kafkaPassword, datasource string,
	maxRowsPerSegment int,
	metricsSpec []*godruid.MetricsSpec,
) *godruid.SupervisorKafka {
	sasLJAASConfig := fmt.Sprintf(
		"org.apache.kafka.common.security.scram.ScramLoginModule required username=\"%s\" password=\"%s\";",
		kafkaUsername,
		kafkaPassword,
	)
	return &godruid.SupervisorKafka{
		Type: "kafka",
		IOConfig: godruid.IOConfigKafka(
			kafkaTopic,
			"PT24H",
			"PT25H",
			kafkaURL,
			"SCRAM-SHA-512",
			"SASL_PLAINTEXT",
			sasLJAASConfig,
		),
		TuningConfig: godruid.TuningConfigKafka(maxRowsPerSegment),
		DataSchema: &godruid.DataSchema{
			Datasource:      datasource,
			GranularitySpec: godruid.GranUniform("HOUR", "HOUR", true),
			Parser: &godruid.Parser{
				Type: "string",
				ParseSpec: &godruid.ParseSpec{
					Format:        "json",
					TimestampSpec: &godruid.TimestampSpec{Column: "timestamp", Format: "iso"},
					DimensionSpec: &godruid.EmptyDimension{},
				},
			},
			MetricsSpec: metricsSpec,
		},
	}
}

func retentionRule(period string) *godruid.RetentionRules {
	return &godruid.RetentionRules{
		Rules: []*godruid.RetentionRule{
			{
				Type:             "loadByPeriod",
				Period:           period,
				IncludeFuture:    false,
				TieredReplicants: &godruid.TieredReplicants{DefaultTier: 2},
			},
			{Type: "dropForever"},
		},
	}
}
