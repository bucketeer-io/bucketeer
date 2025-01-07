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

package api

type EventType int

const (
	GoalEventType EventType = iota + 1 // eventType starts from 1 for validation.
	// Do NOT remove the goalBatchEventType because the go-server-sdk depends on the same order
	// https://github.com/ca-dp/bucketeer-go-server-sdk/blob/master/pkg/bucketeer/api/rest.go#L35
	GoalBatchEventType // nolint:deadcode,unused,varcheck
	EvaluationEventType
	MetricsEventType
)

type metricsDetailEventType int

const (
	latencyMetricsEventType metricsDetailEventType = iota + 1
	sizeMetricsEventType
	timeoutErrorMetricsEventType
	internalErrorMetricsEventType
	networkErrorMetricsEventType
	internalSdkErrorMetricsEventType
)
