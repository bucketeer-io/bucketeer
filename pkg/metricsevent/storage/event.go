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
	SaveGetEvaluationLatencyMetricsEvent(string, string, time.Duration)
	SaveGetEvaluationSizeMetricsEvent(string, string, int32)
	SaveTimeoutErrorCountMetricsEvent(string)
	SaveInternalErrorCountMetricsEvent(string)
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
