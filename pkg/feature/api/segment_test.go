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
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountproto "github.com/bucketeer-io/bucketeer/proto/account"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	storagemock "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
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
		setup         func(*FeatureService)
		cmd           *featureproto.CreateSegmentCommand
		environmentId string
		expected      error
	}{
		{
			setup: nil,
			cmd: &featureproto.CreateSegmentCommand{
				Name:        "",
				Description: "description",
			},
			environmentId: "ns0",
			expected:      createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().CreateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			cmd: &featureproto.CreateSegmentCommand{
				Name:        "name",
				Description: "description",
			},
			environmentId: "ns0",
			expected:      nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.CreateSegmentRequest{Command: tc.cmd, EnvironmentId: tc.environmentId}
		_, err := service.CreateSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestCreateSegmentNoCommandMySQL(t *testing.T) {
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
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.CreateSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing name",
			setup: nil,
			req: &featureproto.CreateSegmentRequest{
				Name:          "",
				Description:   "description",
				EnvironmentId: "ns0",
			},
			expected: createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().CreateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.CreateSegmentRequest{
				Name:          "name",
				Description:   "description",
				EnvironmentId: "ns0",
			},
			expected: nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		_, err := service.CreateSegment(ctx, tc.req)
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
		setup         func(*FeatureService)
		id            string
		cmd           *featureproto.DeleteSegmentCommand
		environmentId string
		expected      error
	}{
		{
			setup:         nil,
			id:            "",
			cmd:           &featureproto.DeleteSegmentCommand{},
			environmentId: "ns0",
			expected:      createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			id:            "id",
			cmd:           &featureproto.DeleteSegmentCommand{},
			environmentId: "ns0",
			expected:      createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			id:            "id",
			cmd:           &featureproto.DeleteSegmentCommand{},
			environmentId: "ns0",
			expected:      nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.DeleteSegmentRequest{
			Id:            tc.id,
			Command:       tc.cmd,
			EnvironmentId: tc.environmentId,
		}
		_, err := service.DeleteSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestDeleteSegmentNoCommandMySQL(t *testing.T) {
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
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.DeleteSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing id",
			setup: nil,
			req: &featureproto.DeleteSegmentRequest{
				Id:            "",
				EnvironmentId: "ns0",
			},
			expected: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "error: segment not found",
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Return(v2fs.ErrSegmentNotFound)
			},
			req: &featureproto.DeleteSegmentRequest{
				Id:            "id",
				EnvironmentId: "ns0",
			},
			expected: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "success",
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
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().DeleteSegment(
					gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.DeleteSegmentRequest{
				Id:            "id",
				EnvironmentId: "ns0",
			},
			expected: nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		_, err := service.DeleteSegment(ctx, tc.req)
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
		setup         func(*FeatureService)
		id            string
		cmds          []*featureproto.Command
		environmentId string
		expected      error
	}{
		{
			setup:         nil,
			id:            "",
			cmds:          nil,
			environmentId: "ns0",
			expected:      createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			id: "id",
			cmds: []*featureproto.Command{
				{Command: changeSegmentNameCmd},
			},
			environmentId: "ns0",
			expected:      nil,
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		req := &featureproto.UpdateSegmentRequest{
			Id:            tc.id,
			Commands:      tc.cmds,
			EnvironmentId: tc.environmentId,
		}
		_, err := service.UpdateSegment(ctx, req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestUpdateSegmentNoCommandMySQL(t *testing.T) {
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
		desc     string
		setup    func(*FeatureService)
		req      *featureproto.UpdateSegmentRequest
		expected error
	}{
		{
			desc:  "error: missing id",
			setup: nil,
			req: &featureproto.UpdateSegmentRequest{
				EnvironmentId: "ns0",
			},
			expected: createError(statusMissingID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc:  "error: update empty name",
			setup: nil,
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String(""),
			},
			expected: createError(statusMissingName, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name")),
		},
		{
			desc: "success update name",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Name:          wrapperspb.String("new-name"),
			},
			expected: nil,
		},
		{
			desc: "success update description",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransactionV2(
					gomock.Any(), gomock.Any(),
				).Do(func(ctx context.Context, fn func(ctx context.Context, tx mysql.Transaction) error) {
					err := fn(ctx, nil)
					require.NoError(t, err)
				}).Return(nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.All(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{
						Id: "id0",
					},
				}, nil, nil)
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().UpdateSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &featureproto.UpdateSegmentRequest{
				Id:            "id0",
				EnvironmentId: "ns0",
				Description:   wrapperspb.String("new-description"),
			},
		},
	}
	for _, tc := range testcases {
		service := createFeatureService(mockController)
		if tc.setup != nil {
			tc.setup(service)
		}
		ctx = setToken(ctx)
		_, err := service.UpdateSegment(ctx, tc.req)
		assert.Equal(t, tc.expected, err)
	}
}

func TestGetSegmentMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	testcases := []struct {
		desc           string
		setup          func(*FeatureService)
		service        *FeatureService
		context        context.Context
		id             string
		environmentId  string
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: missing id",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithToken(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         nil,
			id:            "",
			environmentId: "ns0",
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
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil, v2fs.ErrSegmentNotFound)
			},
			id:            "id",
			environmentId: "ns0",
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
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{}}, nil, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			id:            "id",
			environmentId: "ns0",
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
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().GetSegment(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Segment{
					Segment: &featureproto.Segment{}}, nil, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			id:            "id",
			environmentId: "ns0",
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
			setup:         func(s *FeatureService) {},
			id:            "id",
			environmentId: "ns0",
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

			req := &featureproto.GetSegmentRequest{Id: tc.id, EnvironmentId: tc.environmentId}
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
		desc           string
		service        *FeatureService
		context        context.Context
		setup          func(*FeatureService)
		pageSize       int64
		environmentId  string
		getExpectedErr func(localizer locale.Localizer) error
	}{
		{
			desc:    "error: exceeded max page size per request",
			service: createFeatureService(mockController),
			context: metadata.NewIncomingContext(
				createContextWithTokenRoleUnassigned(),
				metadata.MD{"accept-language": []string{"ja"}},
			),
			setup:         nil,
			pageSize:      int64(maxPageSizePerRequest + 1),
			environmentId: "ns0",
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
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().ListSegments(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Segment{
					{
						Id: "id",
					},
				}, 0, int64(0), map[string][]string{
					"id": {"id"},
				}, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
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
				s.segmentStorage.(*storagemock.MockSegmentStorage).EXPECT().ListSegments(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Segment{
					{
						Id: "id",
					},
				}, 0, int64(0), map[string][]string{
					"id": {"id"},
				}, nil)
				s.featureStorage.(*storagemock.MockFeatureStorage).EXPECT().ListFeatures(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return([]*featureproto.Feature{
					{
						Id: "id",
					},
				}, 0, int64(0), nil)
			},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
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
			setup:         func(s *FeatureService) {},
			pageSize:      int64(maxPageSizePerRequest),
			environmentId: "ns0",
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

			req := &featureproto.ListSegmentsRequest{PageSize: tc.pageSize, EnvironmentId: tc.environmentId}
			_, err := service.ListSegments(ctx, req)
			assert.Equal(t, tc.getExpectedErr(localizer), err)
		})
	}
}

func setToken(ctx context.Context) context.Context {
	t := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    "email",
	}
	return context.WithValue(ctx, rpc.Key, t)
}
