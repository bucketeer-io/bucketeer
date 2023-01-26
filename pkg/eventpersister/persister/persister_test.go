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

package persister

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	pullermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/puller/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/client"
	esproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

var defaultOptions = options{
	logger: zap.NewNop(),
}

func TestUpsertMAU(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(context.Context, *gomock.Controller) *Persister
		input       proto.Message
		expectedErr error
	}{
		{
			desc: "not executed: mysqlClient is nil",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				return newPersister(ctrl)
			},
			input:       &esproto.UserEvent{},
			expectedErr: nil,
		},
		{
			desc: "not executed: message is not UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				return newPersisterWithMysqlClient(ctrl)
			},
			input:       &eventproto.EvaluationEvent{},
			expectedErr: nil,
		},
		{
			desc: "success upsert UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				p := newPersisterWithMysqlClient(ctrl)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				return p
			},
			input:       &esproto.UserEvent{},
			expectedErr: nil,
		},
		{
			desc: "error upsert UserEvent",
			setup: func(ctx context.Context, ctrl *gomock.Controller) *Persister {
				p := newPersisterWithMysqlClient(ctrl)
				p.mysqlClient.(*mysqlmock.MockClient).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("internal"))
				return p
			},
			input:       &esproto.UserEvent{},
			expectedErr: errors.New("internal"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := p.setup(context.Background(), mockController)
			actualErr := persister.upsertMAU(context.Background(), p.input, "ns")
			assert.Equal(t, p.expectedErr, actualErr)
		})
	}
}

func TestEvaluationCountkey(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureID := "feature_id"
	variationID := "variation_id"
	unix := time.Now().Unix()
	environmentNamespace := "en-1"
	now := time.Unix(unix, 0)
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jpLocation)
	patterns := []struct {
		desc                 string
		kind                 string
		featureID            string
		variationID          string
		environmentNamespace string
		timestamp            int64
		expected             string
	}{
		{
			desc:                 "userCount",
			kind:                 userCountKey,
			featureID:            featureID,
			variationID:          variationID,
			environmentNamespace: environmentNamespace,
			timestamp:            unix,
			expected:             fmt.Sprintf("%s:%s:%d:%s:%s", environmentNamespace, userCountKey, date.Unix(), featureID, variationID),
		},
		{
			desc:                 "eventCount",
			kind:                 eventCountKey,
			featureID:            featureID,
			variationID:          variationID,
			environmentNamespace: environmentNamespace,
			timestamp:            unix,
			expected:             fmt.Sprintf("%s:%s:%d:%s:%s", environmentNamespace, eventCountKey, date.Unix(), featureID, variationID),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			persister := newPersister(mockController)
			actual := persister.newEvaluationCountkey(p.kind, p.featureID, p.variationID, p.environmentNamespace, p.timestamp)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestGetVariationID(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		variationID string
		reason      featureproto.Reason_Type
		expected    string
	}{
		{
			desc:        "get given variation id if off variation",
			variationID: "vID1",
			reason:      featureproto.Reason_OFF_VARIATION,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if target",
			variationID: "vID1",
			reason:      featureproto.Reason_TARGET,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if rule",
			variationID: "vID1",
			reason:      featureproto.Reason_RULE,
			expected:    "vID1",
		},
		{
			desc:        "get given variation id if prerequisite",
			variationID: "vID1",
			reason:      featureproto.Reason_PREREQUISITE,
			expected:    "vID1",
		},
		{
			desc:        "get default variation id if client",
			variationID: "vID1",
			reason:      featureproto.Reason_CLIENT,
			expected:    defaultVariationID,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			actual := getVariationID(p.reason, p.variationID)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newPersister(c *gomock.Controller) *Persister {
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		puller: pullermock.NewMockRateLimitedPuller(c),
		opts:   &defaultOptions,
		logger: defaultOptions.logger,
		ctx:    ctx,
		cancel: cancel,
		doneCh: make(chan struct{}),
	}
}

func newPersisterWithMysqlClient(c *gomock.Controller) *Persister {
	ctx, cancel := context.WithCancel(context.Background())
	return &Persister{
		puller:      pullermock.NewMockRateLimitedPuller(c),
		opts:        &defaultOptions,
		logger:      defaultOptions.logger,
		ctx:         ctx,
		cancel:      cancel,
		doneCh:      make(chan struct{}),
		mysqlClient: mysqlmock.NewMockClient(c),
	}
}
