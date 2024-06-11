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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

type Storage interface {
	SaveGetEvaluationLatencyMetricsEvent(tag, status string, duration time.Duration)
	SaveGetEvaluationSizeMetricsEvent(tag, status string, sizeByte int32)
	SaveTimeoutErrorCountMetricsEvent(tag string)
	SaveInternalErrorCountMetricsEvent(tag string)
	SaveLatencyMetricsEvent(tag, status, sdkVersion string, api eventproto.ApiId, duration time.Duration)
	SaveSizeMetricsEvent(tag, status, sdkVersion string, api eventproto.ApiId, sizeByte int32)
	SaveTimeoutErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId)
	SaveInternalErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId)
	SaveNetworkErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId)
	SaveInternalSdkErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId)
}

type storage struct {
	logger *zap.Logger
}

func NewStorage(logger *zap.Logger, register metrics.Registerer) Storage {
	registerMetrics(register)
	return &storage{logger: logger.Named("storage")}
}

func (s *storage) SaveGetEvaluationLatencyMetricsEvent(tag, status string, duration time.Duration) {
	sdkGetEvaluationsLatencyHistogram.WithLabelValues(tag, status).Observe(duration.Seconds())
}

func (s *storage) SaveGetEvaluationSizeMetricsEvent(tag, status string, sizeByte int32) {
	sdkGetEvaluationsSizeHistogram.WithLabelValues(tag, status).Observe(float64(sizeByte))
}

func (s *storage) SaveTimeoutErrorCountMetricsEvent(tag string) {
	sdkTimeoutErrorCounter.WithLabelValues(tag).Inc()
}

func (s *storage) SaveInternalErrorCountMetricsEvent(tag string) {
	sdkInternalErrorCounter.WithLabelValues(tag).Inc()
}

func (s *storage) SaveLatencyMetricsEvent(
	tag, status, sdkVersion string,
	api eventproto.ApiId,
	duration time.Duration,
) {
	sdkLatencyHistogram.WithLabelValues(tag, status, api.String(), sdkVersion).Observe(duration.Seconds())
}

func (s *storage) SaveSizeMetricsEvent(tag, status, sdkVersion string, api eventproto.ApiId, sizeByte int32) {
	sdkSizeHistogram.WithLabelValues(tag, status, api.String(), sdkVersion).Observe(float64(sizeByte))
}

func (s *storage) SaveTimeoutErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId) {
	sdkTimeoutError.WithLabelValues(tag, api.String(), sdkVersion).Inc()
}

func (s *storage) SaveInternalErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId) {
	sdkInternalError.WithLabelValues(tag, api.String(), sdkVersion).Inc()
}

func (s *storage) SaveNetworkErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId) {
	sdkNetworkError.WithLabelValues(tag, api.String(), sdkVersion).Inc()
}

func (s *storage) SaveInternalSdkErrorMetricsEvent(tag, sdkVersion string, api eventproto.ApiId) {
	sdkInternalSdkError.WithLabelValues(tag, api.String(), sdkVersion).Inc()
}
