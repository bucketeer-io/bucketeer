// Copyright 2023 The Bucketeer Authors.
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

	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
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
		desc                 string
		setup                func(*FeatureService)
		environmentNamespace string
		role                 accountproto.Account_Role
		segmentID            string
		cmd                  *featureproto.BulkUploadSegmentUsersCommand
		expectedErr          error
	}{
		{
			desc:                 "ErrMissingSegmentID",
			setup:                nil,
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "",
			cmd:                  nil,
			expectedErr:          createError(statusMissingSegmentID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "segment_id")),
		},
		{
			desc:                 "ErrMissingCommand",
			setup:                nil,
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd:                  nil,
			expectedErr:          createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			desc:                 "ErrMissingSegmentUsersData",
			setup:                nil,
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd:                  &featureproto.BulkUploadSegmentUsersCommand{},
			expectedErr:          createError(statusMissingSegmentUsersData, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "user_data")),
		},
		{
			desc:                 "ErrExceededMaxSegmentUsersDataSize",
			setup:                nil,
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data: []byte(strings.Repeat("a", maxSegmentUsersDataSize+1)),
			},
			expectedErr: createError(statusExceededMaxSegmentUsersDataSize, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_data_state")),
		},
		{
			desc:                 "ErrUnknownSegmentUserState",
			setup:                nil,
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_State(99),
			},
			expectedErr: createError(statusUnknownSegmentUserState, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_state")),
		},
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "not_found_id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: createError(statusSegmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "ErrSegmentUsersAlreadyUploading",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)))
			},
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: createError(statusSegmentUsersAlreadyUploading, localizer.MustLocalize(locale.SegmentUsersAlreadyUploading)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
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
				EnvironmentNamespace: tc.environmentNamespace,
				SegmentId:            tc.segmentID,
				Command:              tc.cmd,
			}
			ctx = setToken(ctx, tc.role)
			_, err := service.BulkUploadSegmentUsers(ctx, req)
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
		desc                 string
		setup                func(*FeatureService)
		environmentNamespace string
		segmentID            string
		state                featureproto.SegmentUser_State
		expectedErr          error
	}{
		{
			desc:                 "ErrMissingSegmentID",
			setup:                nil,
			environmentNamespace: "ns0",
			segmentID:            "",
			state:                featureproto.SegmentUser_INCLUDED,
			expectedErr:          createError(statusMissingSegmentID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "segment_id")),
		},
		{
			desc:                 "ErrUnknownSegmentUserState",
			setup:                nil,
			environmentNamespace: "ns0",
			segmentID:            "id",
			state:                featureproto.SegmentUser_State(99),
			expectedErr:          createError(statusUnknownSegmentUserState, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "user_state")),
		},
		{
			desc: "ErrSegmentNotFound",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			environmentNamespace: "ns0",
			segmentID:            "id",
			state:                featureproto.SegmentUser_INCLUDED,
			expectedErr:          createError(statusSegmentNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "ErrSegmentStatusNotSuceeded",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			environmentNamespace: "ns0",
			segmentID:            "id",
			state:                featureproto.SegmentUser_INCLUDED,
			expectedErr:          createError(statusSegmentStatusNotSuceeded, localizer.MustLocalize(locale.SegmentStatusNotSucceeded)),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx = setToken(ctx, accountproto.Account_UNASSIGNED)
			req := &featureproto.BulkDownloadSegmentUsersRequest{
				EnvironmentNamespace: tc.environmentNamespace,
				SegmentId:            tc.segmentID,
				State:                tc.state,
			}
			_, err := service.BulkDownloadSegmentUsers(ctx, req)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
