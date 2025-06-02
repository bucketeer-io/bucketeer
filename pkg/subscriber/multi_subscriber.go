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

package subscriber

import (
	"context"

	"go.uber.org/zap"
)

type MultiSubscriber struct {
	subscribers []Subscriber
	opts        options
	logger      *zap.Logger
}

func NewMultiSubscriber(opts ...Option) *MultiSubscriber {
	dopts := options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(&dopts)
	}
	return &MultiSubscriber{
		subscribers: make([]Subscriber, 0, 10),
		opts:        dopts,
		logger:      dopts.logger.Named("multi_subscriber"),
	}
}

func (m *MultiSubscriber) AddSubscriber(subscriber Subscriber) {
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *MultiSubscriber) Start(ctx context.Context) {
	for _, subscriber := range m.subscribers {
		go subscriber.Run(ctx)
	}
}

func (m *MultiSubscriber) Stop() {
	for _, subscriber := range m.subscribers {
		subscriber.Stop()
	}
}
