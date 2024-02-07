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
	"time"

	"golang.org/x/time/rate"
)

type RateLimitedPuller interface {
	Run(context.Context) error
	MessageCh() <-chan *Message
}

type rateLimitedPuller struct {
	puller  Puller
	msgCh   chan *Message
	limiter *rate.Limiter
}

func NewRateLimitedPuller(puller Puller, maxMPS int) RateLimitedPuller {
	return &rateLimitedPuller{
		puller:  puller,
		msgCh:   make(chan *Message),
		limiter: rate.NewLimiter(rate.Limit(maxMPS), maxMPS),
	}
}

func (p *rateLimitedPuller) Run(ctx context.Context) error {
	err := p.puller.Pull(ctx, func(ctx context.Context, msg *Message) {
		rv := p.limiter.Reserve()
		time.Sleep(rv.Delay())
		select {
		case p.msgCh <- msg:
		case <-ctx.Done():
		}
	})
	close(p.msgCh)
	return err
}

func (p *rateLimitedPuller) MessageCh() <-chan *Message {
	return p.msgCh
}
