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

package api

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	domain "github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/locale"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql/mock"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestBulkUploadSegmentUsersMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	testcases := []struct {
		desc          string
		setup         func(*FeatureService)
		environmentId string
		segmentID     string
		cmd           *featureproto.BulkUploadSegmentUsersCommand
		expectedErr   error
	}{
		{
			desc:          "ErrMissingSegmentID",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "",
			cmd:           &featureproto.BulkUploadSegmentUsersCommand{},
			expectedErr:   createError(statusMissingSegmentID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "segment_id")),
		},
		{
			desc:          "ErrMissingSegmentUsersData",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "id",
			cmd:           &featureproto.BulkUploadSegmentUsersCommand{},
			expectedErr:   createError(statusMissingSegmentUsersData, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_data")),
		},
		{
			desc:          "ErrExceededMaxSegmentUsersDataSize",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data: []byte(strings.Repeat("a", maxSegmentUsersDataSize+1)),
			},
			expectedErr: createError(statusExceededMaxSegmentUsersDataSize, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_data_state")),
		},
		{
			desc:          "ErrUnknownSegmentUserState",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_State(99),
			},
			expectedErr: createError(statusUnknownSegmentUserState, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_state")),
		},
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			environmentId: "ns0",
			segmentID:     "not_found_id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: createError(statusSegmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "ErrSegmentUsersAlreadyUploading",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)))
			},
			environmentId: "ns0",
			segmentID:     "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			environmentId: "ns0",
			segmentID:     "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if tc.setup != nil {
				tc.setup(service)
			}
			req := &featureproto.BulkUploadSegmentUsersRequest{
				EnvironmentId: tc.environmentId,
				SegmentId:     tc.segmentID,
				Command:       tc.cmd,
			}
			ctx = setToken(ctx)
			_, err := service.BulkUploadSegmentUsers(ctx, req)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBulkUploadSegmentUsersNoCommandMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	testcases := []struct {
		desc          string
		setup         func(*FeatureService)
		req           *featureproto.BulkUploadSegmentUsersRequest
		environmentId string
		expectedErr   error
	}{
		{
			desc:  "ErrMissingSegmentID",
			setup: nil,
			req: &featureproto.BulkUploadSegmentUsersRequest{
				EnvironmentId: "ns0",
				SegmentId:     "",
			},
			expectedErr: createError(statusMissingSegmentID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "segment_id")),
		},
		{
			desc:  "ErrMissingSegmentUsersData",
			setup: nil,
			req: &featureproto.BulkUploadSegmentUsersRequest{
				EnvironmentId: "ns0",
				SegmentId:     "id",
			},
			expectedErr: createError(statusMissingSegmentUsersData, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_data")),
		},
		{
			desc:  "ErrExceededMaxSegmentUsersDataSize",
			setup: nil,
			req: &featureproto.BulkUploadSegmentUsersRequest{
				Data:          []byte(strings.Repeat("a", maxSegmentUsersDataSize+1)),
				EnvironmentId: "ns0",
				SegmentId:     "id",
			},
			expectedErr: createError(statusExceededMaxSegmentUsersDataSize, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_data_state")),
		},
		{
			desc:  "ErrUnknownSegmentUserState",
			setup: nil,
			req: &featureproto.BulkUploadSegmentUsersRequest{
				Data:          []byte("data"),
				State:         featureproto.SegmentUser_State(99),
				EnvironmentId: "ns0",
				SegmentId:     "id",
			},
			expectedErr: createError(statusUnknownSegmentUserState, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_state")),
		},
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			req: &featureproto.BulkUploadSegmentUsersRequest{
				Data:          []byte("data"),
				State:         featureproto.SegmentUser_INCLUDED,
				EnvironmentId: "ns0",
				SegmentId:     "not_found_id",
			},
			expectedErr: createError(statusSegmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "ErrSegmentUsersAlreadyUploading",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)))
			},
			req: &featureproto.BulkUploadSegmentUsersRequest{
				Data:          []byte("data"),
				State:         featureproto.SegmentUser_INCLUDED,
				EnvironmentId: "ns0",
				SegmentId:     "id",
			},
			expectedErr: createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.BulkUploadSegmentUsersRequest{
				Data:          []byte("data"),
				State:         featureproto.SegmentUser_INCLUDED,
				EnvironmentId: "ns0",
				SegmentId:     "id",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx = setToken(ctx)
			_, err := service.BulkUploadSegmentUsers(ctx, tc.req)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBulkDownloadSegmentUsersMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	testcases := []struct {
		desc          string
		setup         func(*FeatureService)
		environmentId string
		segmentID     string
		state         featureproto.SegmentUser_State
		expectedErr   error
	}{
		{
			desc:          "ErrMissingSegmentID",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "",
			state:         featureproto.SegmentUser_INCLUDED,
			expectedErr:   createError(statusMissingSegmentID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "segment_id")),
		},
		{
			desc:          "ErrUnknownSegmentUserState",
			setup:         nil,
			environmentId: "ns0",
			segmentID:     "id",
			state:         featureproto.SegmentUser_State(99),
			expectedErr:   createError(statusUnknownSegmentUserState, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_state")),
		},
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil, v2fs.ErrSegmentNotFound)
			},
			environmentId: "ns0",
			segmentID:     "id",
			state:         featureproto.SegmentUser_INCLUDED,
			expectedErr:   createError(statusSegmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "ErrSegmentStatusNotSuceeded",
			setup: func(s *FeatureService) {
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{},
				}, nil, nil)
			},
			environmentId: "ns0",
			segmentID:     "id",
			state:         featureproto.SegmentUser_INCLUDED,
			expectedErr:   createError(statusSegmentStatusNotSuceeded, localizer.MustLocalize(locale.SegmentStatusNotSucceeded)),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx = setToken(ctx)
			req := &featureproto.BulkDownloadSegmentUsersRequest{
				EnvironmentId: tc.environmentId,
				SegmentId:     tc.segmentID,
				State:         tc.state,
			}
			_, err := service.BulkDownloadSegmentUsers(ctx, req)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
