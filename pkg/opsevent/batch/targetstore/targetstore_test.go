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

package targetstore

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	autoopsclientmock "github.com/bucketeer-io/bucketeer/pkg/autoops/client/mock"
	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentclientmock "github.com/bucketeer-io/bucketeer/pkg/environment/client/mock"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestWithRefreshInterval(t *testing.T) {
	t.Parallel()
	dur := time.Second
	f := WithRefreshInterval(dur)
	opt := &options{}
	f(opt)
	assert.Equal(t, dur, opt.refreshInterval)
}

func TestWithMetrics(t *testing.T) {
	t.Parallel()
	metrics := metrics.NewMetrics(
		9999,
		"/metrics",
	)
	reg := metrics.DefaultRegisterer()
	f := WithMetrics(reg)
	opt := &options{}
	f(opt)
	assert.Equal(t, reg, opt.metrics)
}

func TestWithLogger(t *testing.T) {
	t.Parallel()
	logger, err := log.NewLogger()
	require.NoError(t, err)
	f := WithLogger(logger)
	opt := &options{}
	f(opt)
	assert.Equal(t, logger, opt.logger)
}

func TestNewTargetStore(t *testing.T) {
	g := NewTargetStore(nil, nil)
	assert.IsType(t, &targetStore{}, g)
}

func TestListEnvironments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*targetStore)
		expected    []*environmentproto.Environment
		expectedErr error
	}{
		{
			desc: "enable",
			setup: func(ts *targetStore) {
				ts.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsResponse{Environments: []*environmentproto.Environment{
						{Id: "ns0", Namespace: "ns0"},
					}}, nil)
			},
			expected: []*environmentproto.Environment{{Id: "ns0", Namespace: "ns0"}},
		},
		{
			desc: "list environments fails",
			setup: func(ts *targetStore) {
				ts.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Unknown, "test"))
			},
			expectedErr: status.Errorf(codes.Unknown, "test"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newTargetStoreWithMock(t, mockController)
			p.setup(s)
			actual, err := s.listEnvironments(context.Background())
			assert.Equal(t, p.expected, actual)
			if err != nil {
				assert.Equal(t, p.expectedErr, err)
			}
		})
	}
}

func TestListAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc                 string
		setup                func(*targetStore)
		environmentNamespace string
		expected             []*autoopsproto.AutoOpsRule
		expectedErr          error
	}{
		{
			desc: "success",
			setup: func(ts *targetStore) {
				ts.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ListAutoOpsRules(gomock.Any(), gomock.Any()).Return(
					&autoopsproto.ListAutoOpsRulesResponse{AutoOpsRules: []*autoopsproto.AutoOpsRule{
						{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{},
						},
					}}, nil)
			},
			expected: []*autoopsproto.AutoOpsRule{{Id: "id-0", FeatureId: "fid-0", Clauses: []*autoopsproto.Clause{}}},
		},
		{
			desc: "failure",
			setup: func(ts *targetStore) {
				ts.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ListAutoOpsRules(gomock.Any(), gomock.Any()).Return(
					nil, status.Errorf(codes.Unknown, "test"))
			},
			expectedErr: status.Errorf(codes.Unknown, "test"),
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newTargetStoreWithMock(t, mockController)
			p.setup(s)
			actual, err := s.listAutoOpsRules(context.Background(), p.environmentNamespace)
			assert.Equal(t, p.expected, actual)
			if err != nil {
				assert.Equal(t, p.expectedErr, err)
			}
		})
	}
}

func TestRefreshEnvironments(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc     string
		setup    func(*targetStore)
		expected []*environmentdomain.Environment
	}{
		{
			desc: "enable",
			setup: func(ts *targetStore) {
				ts.environmentClient.(*environmentclientmock.MockClient).EXPECT().ListEnvironments(gomock.Any(), gomock.Any()).Return(
					&environmentproto.ListEnvironmentsResponse{Environments: []*environmentproto.Environment{
						{Id: "ns0", Namespace: "ns0"},
						{Id: "ns1", Namespace: "ns1"},
					}}, nil)
			},
			expected: []*environmentdomain.Environment{
				{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
				{Environment: &environmentproto.Environment{Id: "ns1", Namespace: "ns1"}},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newTargetStoreWithMock(t, mockController)
			p.setup(s)
			_ = s.refreshEnvironments(context.Background())
			actual := s.GetEnvironments(context.Background())
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestRefreshAutoOpsRules(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	c1, c2, c3 := newOpsEventRateClauses(t)
	patterns := []struct {
		desc     string
		setup    func(*targetStore)
		expected []*autoopsdomain.AutoOpsRule
	}{
		{
			desc: "enable",
			setup: func(ts *targetStore) {
				ts.environments.Store([]*environmentdomain.Environment{{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}}})
				ts.autoOpsClient.(*autoopsclientmock.MockClient).EXPECT().ListAutoOpsRules(gomock.Any(), gomock.Any()).Return(
					&autoopsproto.ListAutoOpsRulesResponse{AutoOpsRules: []*autoopsproto.AutoOpsRule{
						{Id: "id-0", FeatureId: "fid-0",
							Clauses: []*autoopsproto.Clause{{Clause: c1}, {Clause: c2}, {Clause: c3}},
						},
						{Id: "id-1", FeatureId: "fid-1", TriggeredAt: time.Now().Unix(),
							Clauses: []*autoopsproto.Clause{{Clause: c1}, {Clause: c2}, {Clause: c3}},
						},
					}}, nil)
			},
			expected: []*autoopsdomain.AutoOpsRule{
				{AutoOpsRule: &autoopsproto.AutoOpsRule{Id: "id-0", FeatureId: "fid-0",
					Clauses: []*autoopsproto.Clause{{Clause: c1}, {Clause: c2}, {Clause: c3}},
				}},
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := newTargetStoreWithMock(t, mockController)
			p.setup(s)
			_ = s.refreshAutoOpsRules(context.Background())
			actual := s.GetAutoOpsRules(context.Background(), "ns0")
			assert.Equal(t, p.expected, actual)
		})
	}
}

func newTargetStoreWithMock(t *testing.T, mockController *gomock.Controller) *targetStore {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	store := &targetStore{
		environmentClient: environmentclientmock.NewMockClient(mockController),
		autoOpsClient:     autoopsclientmock.NewMockClient(mockController),
		autoOpsRules:      make(map[string][]*autoopsdomain.AutoOpsRule),
		logger:            logger,
		ctx:               ctx,
		cancel:            cancel,
		doneCh:            make(chan struct{}),
	}
	store.environments.Store(make([]*environmentdomain.Environment, 0))
	return store
}

func newOpsEventRateClauses(t *testing.T) (*any.Any, *any.Any, *any.Any) {
	c1, err := ptypes.MarshalAny(&autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	})
	require.NoError(t, err)
	c2, err := ptypes.MarshalAny(&autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid2",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	})
	require.NoError(t, err)
	c3, err := ptypes.MarshalAny(&autoopsproto.OpsEventRateClause{
		VariationId:     "vid1",
		GoalId:          "gid1",
		MinCount:        int64(10),
		ThreadsholdRate: float64(0.5),
		Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
	})
	return c1, c2, c3
}
