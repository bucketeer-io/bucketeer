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

package api

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	gstatus "google.golang.org/grpc/status"

	ftstorage "github.com/bucketeer-io/bucketeer/pkg/feature/storage"
	ftmock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	bigtable "github.com/bucketeer-io/bucketeer/pkg/storage/v2/bigtable"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestGetUserEvaluations(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		desc        string
		setup       func(context.Context, ftstorage.UserEvaluationsStorage)
		role        accountproto.Account_Role
		req         *featureproto.GetUserEvaluationsRequest
		expected    *featureproto.GetUserEvaluationsResponse
		expectedErr error
	}{
		{
			desc:  "ErrMissingFeatureTag",
			setup: nil,
			role:  accountproto.Account_EDITOR,
			req: &featureproto.GetUserEvaluationsRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  "",
				UserId:               userID,
			},
			expected:    nil,
			expectedErr: createError(statusMissingFeatureTag, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc:  "ErrMissingUserID",
			setup: nil,
			role:  accountproto.Account_EDITOR,
			req: &featureproto.GetUserEvaluationsRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  tag,
				UserId:               "",
			},
			expected:    nil,
			expectedErr: createError(statusMissingUserID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_id")),
		},
		{
			desc: "ErrInternal",
			setup: func(ctx context.Context, s ftstorage.UserEvaluationsStorage) {
				s.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					userID,
					environmentNamespace,
					tag,
				).Return(nil, bigtable.ErrInternal).Times(1)
			},
			role: accountproto.Account_EDITOR,
			req: &featureproto.GetUserEvaluationsRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  tag,
				UserId:               userID,
			},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(ctx context.Context, s ftstorage.UserEvaluationsStorage) {
				s.(*ftmock.MockUserEvaluationsStorage).EXPECT().GetUserEvaluations(
					ctx,
					userID,
					environmentNamespace,
					tag,
				).Return([]*featureproto.Evaluation{}, nil).Times(1)
			},
			role: accountproto.Account_EDITOR,
			req: &featureproto.GetUserEvaluationsRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  tag,
				UserId:               userID,
			},
			expected: &featureproto.GetUserEvaluationsResponse{
				Evaluations: []*featureproto.Evaluation{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		service := createFeatureServiceNew(mockController)
		ctx = setToken(ctx, p.role)
		if p.setup != nil {
			p.setup(ctx, service.userEvaluationStorage)
		}
		resp, err := service.GetUserEvaluations(
			ctx,
			p.req,
		)
		assert.Equal(t, p.expected, resp, p.desc)
		assert.Equal(t, p.expectedErr, err, p.desc)
	}
}

func TestUpsertUserEvaluation(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	patterns := []struct {
		desc        string
		setup       func(context.Context, ftstorage.UserEvaluationsStorage)
		role        accountproto.Account_Role
		req         *featureproto.UpsertUserEvaluationRequest
		expected    *featureproto.UpsertUserEvaluationResponse
		expectedErr error
	}{
		{
			desc:  "ErrPermissionDenied",
			setup: nil,
			role:  accountproto.Account_UNASSIGNED,
			req: &featureproto.UpsertUserEvaluationRequest{
				EnvironmentNamespace: environmentNamespace,
				Evaluation:           evaluation,
				Tag:                  tag,
			},
			expected:    nil,
			expectedErr: createError(statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied)),
		},
		{
			desc:  "ErrMissingFeatureTag",
			setup: nil,
			role:  accountproto.Account_EDITOR,
			req: &featureproto.UpsertUserEvaluationRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  "",
				Evaluation:           evaluation,
			},
			expected:    nil,
			expectedErr: createError(statusMissingFeatureTag, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag")),
		},
		{
			desc:  "ErrMissingEvaluation",
			setup: nil,
			role:  accountproto.Account_EDITOR,
			req: &featureproto.UpsertUserEvaluationRequest{
				EnvironmentNamespace: environmentNamespace,
				Tag:                  tag,
				Evaluation:           nil,
			},
			expected:    nil,
			expectedErr: createError(statusMissingEvaluation, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "evaluation")),
		},
		{
			desc: "ErrInternal",
			setup: func(ctx context.Context, s ftstorage.UserEvaluationsStorage) {
				s.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					tag,
				).Return(bigtable.ErrInternal).Times(1)
			},
			role: accountproto.Account_EDITOR,
			req: &featureproto.UpsertUserEvaluationRequest{
				EnvironmentNamespace: environmentNamespace,
				Evaluation:           evaluation,
				Tag:                  tag,
			},
			expected:    nil,
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(ctx context.Context, s ftstorage.UserEvaluationsStorage) {
				s.(*ftmock.MockUserEvaluationsStorage).EXPECT().UpsertUserEvaluation(
					ctx,
					evaluation,
					environmentNamespace,
					tag,
				).Return(nil).Times(1)
			},
			role: accountproto.Account_EDITOR,
			req: &featureproto.UpsertUserEvaluationRequest{
				EnvironmentNamespace: environmentNamespace,
				Evaluation:           evaluation,
				Tag:                  tag,
			},
			expected:    &featureproto.UpsertUserEvaluationResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		service := createFeatureServiceNew(mockController)
		ctx = setToken(ctx, p.role)
		if p.setup != nil {
			p.setup(ctx, service.userEvaluationStorage)
		}
		resp, err := service.UpsertUserEvaluation(
			ctx,
			p.req,
		)
		assert.Equal(t, p.expected, resp, p.desc)
		assert.Equal(t, p.expectedErr, err, p.desc)
	}
}
