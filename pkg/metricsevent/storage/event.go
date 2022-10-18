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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package storage

import (
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

type Storage interface {
	SaveGetEvaluationLatencyMetricsEvent(tag, status string, duration time.Duration)
	SaveGetEvaluationSizeMetricsEvent(tag, status string, sizeByte int32)
	SaveTimeoutErrorCountMetricsEvent(tag string)
	SaveInternalErrorCountMetricsEvent(tag string)
	SaveLatencyMetricsEvent(tag, status, api string, duration time.Duration)
	SaveSizeMetricsEvent(tag, status, api string, sizeByte int32)
	SaveTimeoutErrorMetricsEvent(tag, api string)
	SaveInternalErrorMetricsEvent(tag, api string)
	SaveNetworkErrorMetricsEvent(tag, api string)
	SaveInternalSdkErrorMetricsEvent(tag, api string)
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

func (s *storage) SaveLatencyMetricsEvent(tag, status, api string, duration time.Duration) {
	sdkLatencyHistogram.WithLabelValues(tag, status, api).Observe(duration.Seconds())
}

func (s *storage) SaveSizeMetricsEvent(tag, status, api string, sizeByte int32) {
	sdkSizeHistogram.WithLabelValues(tag, status, api).Observe(float64(sizeByte))
}

func (s *storage) SaveTimeoutErrorMetricsEvent(tag, api string) {
	sdkTimeoutError.WithLabelValues(tag, api).Inc()
}

func (s *storage) SaveInternalErrorMetricsEvent(tag, api string) {
	sdkInternalError.WithLabelValues(tag, api).Inc()
}

func (s *storage) SaveNetworkErrorMetricsEvent(tag, api string) {
	sdkNetworkError.WithLabelValues(tag, api).Inc()
}

func (s *storage) SaveInternalSdkErrorMetricsEvent(tag, api string) {
	sdkInternalSdkError.WithLabelValues(tag, api).Inc()
}
