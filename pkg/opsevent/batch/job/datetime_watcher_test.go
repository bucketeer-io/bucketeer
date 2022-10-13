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

package job

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	autoopsdomain "github.com/bucketeer-io/bucketeer/pkg/autoops/domain"
	environmentdomain "github.com/bucketeer-io/bucketeer/pkg/environment/domain"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	executormock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/executor/mock"
	targetstoremock "github.com/bucketeer-io/bucketeer/pkg/opsevent/batch/targetstore/mock"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
)

func TestNewDatetimeWatcher(t *testing.T) {
	w := NewDatetimeWatcher(nil, nil)
	assert.IsType(t, &datetimeWatcher{}, w)
}

func newNewDatetimeWatcherWithMock(t *testing.T, mockController *gomock.Controller) *datetimeWatcher {
	logger, err := log.NewLogger()
	require.NoError(t, err)
	return &datetimeWatcher{
		environmentLister: targetstoremock.NewMockEnvironmentLister(mockController),
		autoOpsRuleLister: targetstoremock.NewMockAutoOpsRuleLister(mockController),
		autoOpsExecutor:   executormock.NewMockAutoOpsExecutor(mockController),
		logger:            logger,
		opts: &options{
			timeout: time.Minute,
		},
	}
}

func TestRunDatetimeWatcher(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*datetimeWatcher)
		expectedErr error
	}{
		{
			desc: "success: assess: false",
			setup: func(w *datetimeWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.DatetimeClause{Time: time.Now().AddDate(0, 0, 1).Unix()}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.autoOpsRuleLister.(*targetstoremock.MockAutoOpsRuleLister).EXPECT().GetAutoOpsRules(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.AutoOpsRule{
						{AutoOpsRule: &autoopsproto.AutoOpsRule{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{{Clause: c}},
						}},
					},
				)
			},
			expectedErr: nil,
		},
		{
			desc: "success: assess: true",
			setup: func(w *datetimeWatcher) {
				w.environmentLister.(*targetstoremock.MockEnvironmentLister).EXPECT().GetEnvironments(gomock.Any()).Return(
					[]*environmentdomain.Environment{
						{Environment: &environmentproto.Environment{Id: "ns0", Namespace: "ns0"}},
					},
				)
				dc := &autoopsproto.DatetimeClause{Time: time.Now().Unix()}
				c, err := ptypes.MarshalAny(dc)
				require.NoError(t, err)
				w.autoOpsRuleLister.(*targetstoremock.MockAutoOpsRuleLister).EXPECT().GetAutoOpsRules(gomock.Any(), "ns0").Return(
					[]*autoopsdomain.AutoOpsRule{
						{AutoOpsRule: &autoopsproto.AutoOpsRule{
							Id:        "id-0",
							FeatureId: "fid-0",
							Clauses:   []*autoopsproto.Clause{{Clause: c}},
						}},
					},
				)
				w.autoOpsExecutor.(*executormock.MockAutoOpsExecutor).EXPECT().Execute(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newNewDatetimeWatcherWithMock(t, mockController)
			if p.setup != nil {
				p.setup(w)
			}
			err := w.Run(context.Background())
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
func TestDatetimeWatcherAssessRule(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc           string
		datetimeClause *autoopsproto.DatetimeClause
		nowTimestamp   int64
		expected       bool
	}{
		{
			desc: "false",
			datetimeClause: &autoopsproto.DatetimeClause{
				Time: 1000000001,
			},
			nowTimestamp: 1000000000,
			expected:     false,
		},
		{
			desc: "true",
			datetimeClause: &autoopsproto.DatetimeClause{
				Time: 1000000000,
			},
			nowTimestamp: 1000000000,
			expected:     true,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			w := newNewDatetimeWatcherWithMock(t, mockController)
			actual := w.assessRule(p.datetimeClause, p.nowTimestamp)
			assert.Equal(t, p.expected, actual)
		})
	}
}
