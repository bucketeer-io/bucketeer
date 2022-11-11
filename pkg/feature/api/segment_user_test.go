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
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

	localizer := locale.NewLocalizer(locale.NewLocale(locale.JaJP))
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
			expectedErr:          errMissingSegmentIDJaJP,
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
			expectedErr:          errMissingSegmentUsersDataJaJP,
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
			expectedErr: errExceededMaxSegmentUsersDataSizeJaJP,
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
			expectedErr: errUnknownSegmentUserStateJaJP,
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
			expectedErr: errSegmentNotFoundJaJP,
		},
		{
			desc: "ErrSegmentUsersAlreadyUploading",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(localizedError(statusSegmentUsersAlreadyUploading, locale.JaJP))
			},
			environmentNamespace: "ns0",
			role:                 accountproto.Account_OWNER,
			segmentID:            "id",
			cmd: &featureproto.BulkUploadSegmentUsersCommand{
				Data:  []byte("data"),
				State: featureproto.SegmentUser_INCLUDED,
			},
			expectedErr: errSegmentUsersAlreadyUploadingJaJP,
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
			ctx := setToken(context.Background(), tc.role)
			req := &featureproto.BulkUploadSegmentUsersRequest{
				EnvironmentNamespace: tc.environmentNamespace,
				SegmentId:            tc.segmentID,
				Command:              tc.cmd,
			}
			_, err := service.BulkUploadSegmentUsers(ctx, req)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBulkDownloadSegmentUsersMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

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
			expectedErr:          errMissingSegmentIDJaJP,
		},
		{
			desc:                 "ErrUnknownSegmentUserState",
			setup:                nil,
			environmentNamespace: "ns0",
			segmentID:            "id",
			state:                featureproto.SegmentUser_State(99),
			expectedErr:          errUnknownSegmentUserStateJaJP,
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
			expectedErr:          errSegmentNotFoundJaJP,
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
			expectedErr:          errSegmentStatusNotSuceededJaJP,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := createFeatureService(mockController)
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
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
