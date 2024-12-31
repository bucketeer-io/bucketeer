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

package processor

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/codes"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber/storage"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
)

var (
	getEvaluationLatencyMetricsEvent = &eventproto.GetEvaluationLatencyMetricsEvent{}
	getEvaluationSizeMetricsEvent    = &eventproto.GetEvaluationSizeMetricsEvent{}
	timeoutErrorCountMetricsEvent    = &eventproto.TimeoutErrorCountMetricsEvent{}
	internalErrorCountMetricsEvent   = &eventproto.InternalErrorCountMetricsEvent{}
	latencyMetricsEvent              = &eventproto.LatencyMetricsEvent{}
	sizeMetricsEvent                 = &eventproto.SizeMetricsEvent{}
	timeoutErrorMetricsEvent         = &eventproto.TimeoutErrorMetricsEvent{}
	internalErrorMetricsEvent        = &eventproto.InternalErrorMetricsEvent{}
	networkErrorMetricsEvent         = &eventproto.NetworkErrorMetricsEvent{}
	internalSdkErrorMetricsEvent     = &eventproto.InternalSdkErrorMetricsEvent{}
)

type metricsEventPersister struct {
	storage storage.Storage
	logger  *zap.Logger
}

func NewMetricsEventPersister(
	registerer metrics.Registerer,
	logger *zap.Logger,
) subscriber.Processor {
	return &metricsEventPersister{
		storage: storage.NewStorage(logger, registerer),
		logger:  logger,
	}
}

func (m metricsEventPersister) Process(ctx context.Context, msgChan <-chan *puller.Message) error {
	record := func(code codes.Code, startTime time.Time) {
		subscriberHandledCounter.WithLabelValues(subscriberMetricsEvent, code.String()).Inc()
		subscriberHandledHistogram.WithLabelValues(
			subscriberMetricsEvent,
			code.String(),
		).Observe(time.Since(startTime).Seconds())
	}
	for {
		select {
		case msg, ok := <-msgChan:
			if !ok {
				return nil
			}
			subscriberReceivedCounter.WithLabelValues(subscriberMetricsEvent).Inc()
			startTime := time.Now()
			if id := msg.Attributes["id"]; id == "" {
				m.logger.Error("message has no id")
				msg.Ack()
				record(codes.MissingID, startTime)
				continue
			}
			err := m.handle(msg)
			if err != nil {
				msg.Ack()
				record(codes.NonRepeatableError, startTime)
				continue
			}
			msg.Ack()
			record(codes.OK, startTime)
		case <-ctx.Done():
			return nil
		}
	}
}

func (m metricsEventPersister) handle(message *puller.Message) error {
	metricsEvents, err := m.unmarshalMessage(message)
	if err != nil {
		m.logger.Error("message is bad")
		return err
	}
	err = m.saveMetrics(metricsEvents)
	if err != nil {
		m.logger.Error("could not store data to prometheus client", zap.Error(err))
		return err
	}
	return nil
}

func (m metricsEventPersister) unmarshalMessage(message *puller.Message) (*eventproto.MetricsEvent, error) {
	event := &eventproto.Event{}
	if err := proto.Unmarshal(message.Data, event); err != nil {
		m.logger.Error("ummarshal event failed",
			zap.Error(err),
			zap.Any("msg", message),
		)
		subscriberHandledCounter.WithLabelValues(subscriberMetricsEvent, codes.BadMessage.String()).Inc()
		return nil, err
	}
	me := &eventproto.MetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, me); err != nil {
		m.logger.Error("ummarshal metrics event failed",
			zap.Error(err),
			zap.Any("msg", message),
		)
		subscriberHandledCounter.WithLabelValues(subscriberMetricsEvent, codes.BadMessage.String()).Inc()
		return nil, err
	}
	return me, nil
}

func (m metricsEventPersister) saveMetrics(event *eventproto.MetricsEvent) error {
	if ptypes.Is(event.Event, getEvaluationLatencyMetricsEvent) {
		return m.saveGetEvaluationLatencyMetricsEvent(event)
	}
	if ptypes.Is(event.Event, getEvaluationSizeMetricsEvent) {
		return m.saveGetEvaluationSizeMetricsEvent(event)
	}
	if ptypes.Is(event.Event, timeoutErrorCountMetricsEvent) {
		return m.saveTimeoutErrorCountMetricsEvent(event)
	}
	if ptypes.Is(event.Event, internalErrorCountMetricsEvent) {
		return m.saveInternalErrorCountMetricsEvent(event)
	}
	if ptypes.Is(event.Event, latencyMetricsEvent) {
		return m.saveLatencyMetricsEvent(event)
	}
	if ptypes.Is(event.Event, sizeMetricsEvent) {
		return m.saveSizeMetricsEvent(event)
	}
	if ptypes.Is(event.Event, timeoutErrorMetricsEvent) {
		return m.saveTimeoutError(event)
	}
	if ptypes.Is(event.Event, internalErrorMetricsEvent) {
		return m.saveInternalError(event)
	}
	if ptypes.Is(event.Event, networkErrorMetricsEvent) {
		return m.saveNetworkError(event)
	}
	if ptypes.Is(event.Event, internalSdkErrorMetricsEvent) {
		return m.saveInternalSdkError(event)
	}
	return ErrUnknownEvent
}

func (m metricsEventPersister) saveGetEvaluationLatencyMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.GetEvaluationLatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.Duration == nil {
		return ErrInvalidDuration
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return ErrInvalidDuration
	}
	m.storage.SaveGetEvaluationLatencyMetricsEvent(tag, status, dur)
	return nil
}

func (m metricsEventPersister) saveGetEvaluationSizeMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.GetEvaluationSizeMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	m.storage.SaveGetEvaluationSizeMetricsEvent(tag, status, ev.SizeByte)
	return nil
}

func (m metricsEventPersister) saveTimeoutErrorCountMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.TimeoutErrorCountMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	m.storage.SaveTimeoutErrorCountMetricsEvent(ev.Tag)
	return nil
}

func (m metricsEventPersister) saveInternalErrorCountMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.InternalErrorCountMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	m.storage.SaveInternalErrorCountMetricsEvent(ev.Tag)
	return nil
}

func (m metricsEventPersister) saveLatencyMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.LatencyMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	dur, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return ErrInvalidDuration
	}
	m.storage.SaveLatencyMetricsEvent(tag, status, event.SdkVersion, ev.ApiId, dur)
	return nil
}

func (m metricsEventPersister) saveSizeMetricsEvent(event *eventproto.MetricsEvent) error {
	ev := &eventproto.SizeMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag, status string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
		status = ev.Labels["state"]
	}
	m.storage.SaveSizeMetricsEvent(tag, status, event.SdkVersion, ev.ApiId, ev.SizeByte)
	return nil
}

func (m metricsEventPersister) saveTimeoutError(event *eventproto.MetricsEvent) error {
	ev := &eventproto.TimeoutErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	m.storage.SaveTimeoutErrorMetricsEvent(tag, event.SdkVersion, ev.ApiId)
	return nil
}

func (m metricsEventPersister) saveInternalError(event *eventproto.MetricsEvent) error {
	ev := &eventproto.InternalErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	m.storage.SaveInternalErrorMetricsEvent(tag, event.SdkVersion, ev.ApiId)
	return nil
}

func (m metricsEventPersister) saveNetworkError(event *eventproto.MetricsEvent) error {
	ev := &eventproto.NetworkErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	m.storage.SaveNetworkErrorMetricsEvent(tag, event.SdkVersion, ev.ApiId)
	return nil
}

func (m metricsEventPersister) saveInternalSdkError(event *eventproto.MetricsEvent) error {
	ev := &eventproto.InternalSdkErrorMetricsEvent{}
	if err := ptypes.UnmarshalAny(event.Event, ev); err != nil {
		return err
	}
	if ev.ApiId == eventproto.ApiId_UNKNOWN_API {
		return ErrUnknownApiId
	}
	var tag string
	if ev.Labels != nil {
		tag = ev.Labels["tag"]
	}
	m.storage.SaveInternalSdkErrorMetricsEvent(tag, event.SdkVersion, ev.ApiId)
	return nil
}
