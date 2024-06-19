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

package api

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	accountproto "github.com/bucketeer-io/bucketeer/proto/account"

	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	"github.com/bucketeer-io/bucketeer/pkg/token"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateSegmentMySQL(t *testing.T) {
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
		setup                func(*FeatureService)
		cmd                  *featureproto.CreateSegmentCommand
		environmentNamespace string
		expected             error
	}{
		{
			setup:                nil,
			cmd:                  nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			setup: nil,
			cmd: &featureproto.CreateSegmentCommand{
				Name:        "",
				Description: "description",
			},
			environmentNamespace: "ns0",
			expected:             createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			cmd: &featureproto.CreateSegmentCommand{
				Name:        "name",
				Description: "description",
			},
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.CreateSegmentRequest{Command: tc.cmd, EnvironmentNamespace: tc.environmentNamespace}
		_, err := service.CreateSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestDeleteSegmentMySQL(t *testing.T) {
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
		setup                func(*FeatureService)
		id                   string
		cmd                  *featureproto.DeleteSegmentCommand
		environmentNamespace string
		expected             error
	}{
		{
			setup:                nil,
			id:                   "",
			cmd:                  nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup:                nil,
			id:                   "id",
			cmd:                  nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			id:                   "id",
			cmd:                  &featureproto.DeleteSegmentCommand{},
			environmentNamespace: "ns0",
			expected:             createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			id:                   "id",
			cmd:                  &featureproto.DeleteSegmentCommand{},
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.DeleteSegmentRequest{
			Id:                   tc.id,
			Command:              tc.cmd,
			EnvironmentNamespace: tc.environmentNamespace,
		}
		_, err := service.DeleteSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestUpdateSegmentMySQL(t *testing.T) {
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

	changeSegmentNameCmd, err := ptypes.MarshalAny(&featureproto.ChangeSegmentNameCommand{Name: "name"})
	require.NoError(t, err)
	testcases := []struct {
		setup                func(*FeatureService)
		id                   string
		cmds                 []*featureproto.Command
		environmentNamespace string
		expected             error
	}{
		{
			setup:                nil,
			id:                   "",
			cmds:                 nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup:                nil,
			id:                   "id",
			cmds:                 nil,
			environmentNamespace: "ns0",
			expected:             createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "command")),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			id: "id",
			cmds: []*featureproto.Command{
				{Command: changeSegmentNameCmd},
			},
			environmentNamespace: "ns0",
			expected:             nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.UpdateSegmentRequest{
			Id:                   tc.id,
			Commands:             tc.cmds,
			EnvironmentNamespace: tc.environmentNamespace,
		}
		_, err := service.UpdateSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestGetSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		desc                 string
		setup                func(*FeatureService)
		service              *FeatureService
		context              context.Context
		id                   string
		environmentNamespace string
		getExpectedErr       func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: missing id",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:                nil,
			id:                   "",
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"), localizer)
			},
		},
		{
			desc:    "error: segment not found",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id",
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusNotFound, localizer.MustLocalize(locale.NotFoundError), localizer)
			},
		},
		{
			desc:    "success",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			id:                   "id",
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success with Viewer account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
			},
			id:                   "id",
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:                func(s *FeatureService) {},
			id:                   "id",
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := tc.service
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx := tc.context
			localizer := locale.NewLocalizer(ctx)

			req := &featureproto.GetSegmentRequest{Id: tc.id, EnvironmentNamespace: tc.environmentNamespace}
			_, err := service.GetSegment(ctx, req)
			assert.Equal(t, tc.getExpectedErr(localizer), err)
		})
	}
}

func TestListSegmentsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		desc                 string
		service              *FeatureService
		context              context.Context
		setup                func(*FeatureService)
		pageSize             int64
		environmentNamespace string
		getExpectedErr       func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: exceeded max page size per request",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:                nil,
			pageSize:             int64(maxPageSizePerRequest + 1),
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusExceededMaxPageSizePerRequest, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "page_size"), localizer)
			},
		},
		{
			desc:    "success",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil).Times(2)
				rows.EXPECT().Next().Return(false).Times(2)
				rows.EXPECT().Err().Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil).Times(2)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
			},
			pageSize:             int64(maxPageSizePerRequest),
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "success with Viewer account",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_MEMBER, accountproto.AccountV2_Role_Environment_VIEWER),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil).Times(2)
				rows.EXPECT().Next().Return(false).Times(2)
				rows.EXPECT().Err().Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil).Times(2)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(2)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row).Times(2)
			},
			pageSize:             int64(maxPageSizePerRequest),
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return nil
			},
		},
		{
			desc:    "errPermissionDenied",
			service: createFeatureServiceWithGetAccountByEnvironmentMock(mockController, accountproto.AccountV2_Role_Organization_UNASSIGNED, accountproto.AccountV2_Role_Environment_UNASSIGNED),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:                func(s *FeatureService) {},
			pageSize:             int64(maxPageSizePerRequest),
			environmentNamespace: "ns0",
			getExpectedErr: func(localizer locale.Localizer) error {
				return createError(t, statusPermissionDenied, localizer.MustLocalize(locale.PermissionDenied), localizer)
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			service := tc.service
			if tc.setup != nil {
				tc.setup(service)
			}
			ctx := tc.context
			localizer := locale.NewLocalizer(ctx)

			req := &featureproto.ListSegmentsRequest{PageSize: tc.pageSize, EnvironmentNamespace: tc.environmentNamespace}
			_, err := service.ListSegments(ctx, req)
			assert.Equal(t, tc.getExpectedErr(localizer), err)
		})
	}
}

func setToken(ctx context.Context) context.Context {
	t := &token.AccessToken{
		Issuer:   "issuer",
		Subject:  "sub",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	return context.WithValue(ctx, rpc.Key, t)
}
