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

package processor

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	subscriberAuditLog           = "AuditLog"
	subscriberDomainEvent        = "DomainEvent"
	subscriberEvaluationCount    = "EvaluationCount"
	subscriberEvaluationEventDWH = "EvaluationEventDWH"
	subscriberEvaluationEventOPS = "EvaluationEventOPS"
	subscriberGoalEventDWH       = "GoalEventDWH"
	subscriberGoalEventOPS       = "GoalEventOPS"
	subscriberMetricsEvent       = "MetricsEvent"
	subscriberPushSender         = "PushSender"
	subscriberSegmentUser        = "SegmentUser"
	subscriberUserEvent          = "UserEvent"
)

const (
	codeAutoOpsRuleNotFound                 = "ErrAutoOpsRuleNotFound"
	codeFailedToExtractOpsEventRateClauses  = "FailedToExtractOpsEventRateClauses"
	codeFailedToGetFeatures                 = "FailedToGetFeatures"
	codeFailedToFindFeatureVersion          = "FailedToFindFeatureVersion"
	codeFailedToListAutoOpsRules            = "FailedToListAutoOpsRules"
	codeFailedToUpdateUserCount             = "FailedToUpdateUserCount"
	codeFailedToGetUserEvaluation           = "FailedToGetUserEvaluation"
	codeFailedToStoreRetryMessage           = "FailedToStoreRetryMessage"
	codeRetryMessageAppendSuccess           = "RetryMessageAppendSuccess"
	codeRetryMessageNoEvaluations           = "RetryMessageNoEvaluations"
	codeRetryMessageNoGoalEvents            = "RetryMessageNoGoalEvents"
	codeGetFeaturesReturnedEmpty            = "GetFeaturesReturnedEmpty"
	codeEvaluationsAreEmpty                 = "EvaluationsAreEmpty"
	codeEventIssuedAfterExperimentEnded     = "EventIssuedAfterExperimentEnded"
	codeGoalEventIssuedBeforeEvaluation     = "GoalEventIssuedBeforeEvaluation"
	codeEventOlderThanExperiment            = "EventOlderThanExperiment"
	codeExperimentNotFound                  = "ExperimentNotFound"
	codeUserEvaluationNotFound              = "UserEvaluationNotFound"
	codeGoalEventIssuedAfterExperimentEnded = "GoalEventIssuedAfterExperimentEnded"
	codeFailedToEvaluateUser                = "FailedToEvaluateUser"
	codeFailedToListExperiments             = "FailedToListExperiments"
	codeFailedToAppendEvaluationEvents      = "FailedToAppendEvaluationEvents"
	codeFailedToAppendGoalEvents            = "FailedToAppendGoalEvents"
	codeLinked                              = "Linked"
	appVersion                              = "app_version"
)

var (
	subscriberReceivedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "subscriber",
			Name:      "received_total",
			Help:      "Total number of received messages",
		}, []string{"subscriber"})

	subscriberHandledCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "subscriber",
			Name:      "handled_total",
			Help:      "Total number of handled messages",
		}, []string{"subscriber", "code"})

	subscriberHandledHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "bucketeer",
			Subsystem: "subscriber",
			Name:      "handled_seconds",
			Help:      "Histogram of message handling duration (seconds)",
			Buckets:   prometheus.DefBuckets,
		}, []string{"subscriber", "code"})

	evaluationEventCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "bucketeer",
			Subsystem: "subscriber",
			Name:      "evaluation_event_total",
			Help:      "Total number of evaluation events",
		}, []string{"environment_id", "sdk_version", "feature_id", "app_version", "variation_id"})
)

func registerMetrics(r metrics.Registerer) {
	r.MustRegister(
		subscriberReceivedCounter,
		subscriberHandledCounter,
		subscriberHandledHistogram,
		evaluationEventCounter,
	)
}
