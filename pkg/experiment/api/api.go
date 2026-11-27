// Copyright 2026 The Bucketeer Authors.
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

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	autoopsclient "github.com/bucketeer-io/bucketeer/v2/pkg/autoops/client"
	storage "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/storage/v2"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/experiment"
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type experimentService struct {
	featureClient     featureclient.Client
	autoOpsClient     autoopsclient.Client
	dbClient          database.Client
	experimentStorage storage.ExperimentStorage
	goalStorage       storage.GoalStorage
	publisher         publisher.Publisher
	opts              *options
	logger            *zap.Logger
}

func NewExperimentService(
	featureClient featureclient.Client,
	autoOpsClient autoopsclient.Client,
	dbClient database.Client,
	experimentStorage storage.ExperimentStorage,
	goalStorage storage.GoalStorage,
	publisher publisher.Publisher,
	opts ...Option,
) rpc.Service {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &experimentService{
		featureClient:     featureClient,
		autoOpsClient:     autoOpsClient,
		dbClient:          dbClient,
		experimentStorage: experimentStorage,
		goalStorage:       goalStorage,
		publisher:         publisher,
		opts:              dopts,
		logger:            dopts.logger.Named("api"),
	}
}

func (s *experimentService) Register(server *grpc.Server) {
	proto.RegisterExperimentServiceServer(server, s)
}

// publishDomainEvent publishes a domain event after the transaction that
// produced it has committed. Publishing inside the transaction would emit the
// event even when the transaction is rolled back afterwards.
// On failure it logs the publisher error and returns statusInternal, so the
// caller does not receive a success response for a committed change whose
// event was never published, and the publisher's error text does not leak
// into the response.
func (s *experimentService) publishDomainEvent(
	ctx context.Context,
	event *eventproto.Event,
) error {
	if err := s.publisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"Failed to publish domain event",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", event.EnvironmentId),
			)...,
		)
		return statusInternal.Err()
	}
	return nil
}
