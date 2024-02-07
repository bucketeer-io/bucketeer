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
package puller

import (
	"context"

	"cloud.google.com/go/pubsub"
	"go.uber.org/zap"
)

type Message struct {
	ID         string
	Data       []byte
	Attributes map[string]string

	Ack  func()
	Nack func()
}

type Puller interface {
	Pull(context.Context, func(context.Context, *Message)) error
}

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type puller struct {
	subscription *pubsub.Subscription
	logger       *zap.Logger
}

func NewPuller(sub *pubsub.Subscription, opts ...Option) Puller {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &puller{
		subscription: sub,
		logger:       dopts.logger.Named("puller"),
	}
}

func (p *puller) Pull(ctx context.Context, f func(context.Context, *Message)) error {
	err := p.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		f(ctx, &Message{
			ID:         msg.ID,
			Data:       msg.Data,
			Attributes: msg.Attributes,
			Ack:        msg.Ack,
			Nack:       msg.Nack})
	})
	if err != nil {
		p.logger.Error("Failed to receive message", zap.Error(err))
		return err
	}
	return nil
}
