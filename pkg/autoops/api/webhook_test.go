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
	"encoding/base64"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	autoopspb "github.com/bucketeer-io/bucketeer/proto/autoops"
)

func TestCreateWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string) error {
		st, err := statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}
	baseSetup := func(s *AutoOpsService) {
		s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
		s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
			gomock.Any(), gomock.Any(), gomock.Any(),
		).Return(nil)
	}

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopspb.CreateWebhookRequest
		resp        *autoopspb.CreateWebhookResponse
		expectedErr error
	}{
		{
			desc:        "err: ErrNoCommand",
			req:         &autoopspb.CreateWebhookRequest{},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrWebhookNameRequired",
			req: &autoopspb.CreateWebhookRequest{
				Command: &autoopspb.CreateWebhookCommand{},
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook name")),
		},
		{
			desc:  "success",
			setup: baseSetup,
			req: &autoopspb.CreateWebhookRequest{
				Command: &autoopspb.CreateWebhookCommand{
					Name:        "name",
					Description: "description",
				},
			},
			resp: &autoopspb.CreateWebhookResponse{
				Webhook: &autoopspb.Webhook{
					Name:        "name",
					Description: "description",
				},
				Url: "https://bucketeer.io/hook?auth=secret",
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createAutoOpsService(mockController, nil)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.CreateWebhook(ctx, p.req)
			if p.resp != nil {
				assert.Equal(t, p.resp.Webhook.Name, resp.Webhook.Name)
				assert.Equal(t, p.resp.Webhook.Description, resp.Webhook.Description)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string) error {
		status, err := statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return status.Err()
	}
	service := createAutoOpsService(mockController, nil)

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopspb.GetWebhookRequest
		res         *autoopspb.GetWebhookResponse
		expectedErr error
	}{
		{
			desc: "err: ErrNoId",
			req: &autoopspb.GetWebhookRequest{
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			req: &autoopspb.GetWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			res: &autoopspb.GetWebhookResponse{
				Webhook: &autoopspb.Webhook{
					Name:        "",
					Description: "",
				},
				Url: "https://bucketeer.io/hook?auth=secret",
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.GetWebhook(ctx, p.req)
			if p.res != nil {
				assert.Equal(t, p.res.Webhook.Name, resp.Webhook.Name)
				assert.Equal(t, p.res.Webhook.Description, resp.Webhook.Name)
			}
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListWebhooks(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string) error {
		status, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return status.Err()
	}
	service := createAutoOpsService(mockController, nil)

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopspb.ListWebhooksRequest
		res         *autoopspb.ListWebhooksResponse
		expectedErr error
	}{
		{
			desc: "err: ErrInvalidArgument",
			req: &autoopspb.ListWebhooksRequest{
				EnvironmentNamespace: "ns0",
				PageSize:             int64(500),
				Cursor:               "abc",
				OrderBy:              autoopspb.ListWebhooksRequest_DEFAULT,
				OrderDirection:       autoopspb.ListWebhooksRequest_ASC,
				SearchKeyword:        "",
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
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
			},
			req: &autoopspb.ListWebhooksRequest{
				EnvironmentNamespace: "ns0",
				PageSize:             int64(500),
				Cursor:               "",
				OrderBy:              autoopspb.ListWebhooksRequest_DEFAULT,
				OrderDirection:       autoopspb.ListWebhooksRequest_ASC,
				SearchKeyword:        "",
			},
			res: &autoopspb.ListWebhooksResponse{
				Webhooks:   []*autoopspb.Webhook{},
				Cursor:     "0",
				TotalCount: 0,
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.ListWebhooks(ctx, p.req)
			assert.Equal(t, p.res, resp)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string) error {
		status, err := statusInvalidRequest.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return status.Err()
	}
	service := createAutoOpsService(mockController, nil)

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopspb.UpdateWebhookRequest
		res         *autoopspb.UpdateWebhookResponse
		expectedErr error
	}{
		{
			desc: "err: ErrNoId",
			req: &autoopspb.UpdateWebhookRequest{
				EnvironmentNamespace:            "ns0",
				ChangeWebhookDescriptionCommand: &autoopspb.ChangeWebhookDescriptionCommand{},
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id")),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopspb.UpdateWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command")),
		},
		{
			desc: "err: ErrNoName",
			req: &autoopspb.UpdateWebhookRequest{
				Id:                       "id-0",
				EnvironmentNamespace:     "ns0",
				ChangeWebhookNameCommand: &autoopspb.ChangeWebhookNameCommand{},
			},
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "webhook name")),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopspb.UpdateWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
				ChangeWebhookNameCommand: &autoopspb.ChangeWebhookNameCommand{
					Name: "name",
				},
				ChangeWebhookDescriptionCommand: &autoopspb.ChangeWebhookDescriptionCommand{},
			},
			res:         &autoopspb.UpdateWebhookResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.UpdateWebhook(ctx, p.req)
			assert.Equal(t, p.res, resp)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithTokenRoleOwner(t)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string, status *status.Status) error {
		status, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return status.Err()
	}
	service := createAutoOpsService(mockController, nil)

	patterns := []struct {
		desc        string
		setup       func(*AutoOpsService)
		req         *autoopspb.DeleteWebhookRequest
		res         *autoopspb.DeleteWebhookResponse
		expectedErr error
	}{
		{
			desc: "err: ErrNoId",
			req: &autoopspb.DeleteWebhookRequest{
				EnvironmentNamespace: "ns0",
				Command:              &autoopspb.DeleteWebhookCommand{},
			},
			expectedErr: createError(
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
				statusInvalidRequest,
			),
		},
		{
			desc: "err: ErrNoCommand",
			req: &autoopspb.DeleteWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
			},
			expectedErr: createError(
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "command"),
				statusInvalidRequest,
			),
		},
		{
			desc: "err: InternalErr",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("Internal error"))
			},
			req: &autoopspb.DeleteWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
				Command:              &autoopspb.DeleteWebhookCommand{},
			},
			expectedErr: createError(
				localizer.MustLocalize(locale.InternalServerError),
				statusInternal,
			),
		},
		{
			desc: "success",
			setup: func(s *AutoOpsService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			req: &autoopspb.DeleteWebhookRequest{
				Id:                   "id-0",
				EnvironmentNamespace: "ns0",
				Command:              &autoopspb.DeleteWebhookCommand{},
			},
			res:         &autoopspb.DeleteWebhookResponse{},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			resp, err := service.DeleteWebhook(ctx, p.req)
			assert.Equal(t, p.res, resp)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

type dummyWebhookSecret struct {
	WebhookID            string `json:"webhook_id"`
	EnvironmentNamespace string `json:"environment_namespace"`
}

func TestGenerateWebhookSecret(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	service := createAutoOpsService(mockController, nil)
	ctx := context.TODO()

	testcases := []struct {
		desc                 string
		id                   string
		environmentNamespace string
	}{
		{
			desc:                 "success",
			id:                   "id-1",
			environmentNamespace: "ns-1",
		},
	}
	for _, p := range testcases {
		t.Run(p.desc, func(t *testing.T) {
			secret, err := service.generateWebhookSecret(ctx, p.id, p.environmentNamespace)
			require.NoError(t, err)
			ws := dummyWebhookSecret{}
			decoded, err := base64.RawURLEncoding.DecodeString(secret)
			require.NoError(t, err)
			err = json.Unmarshal(decoded, &ws)
			require.NoError(t, err)
			assert.Equal(t, p.environmentNamespace, ws.EnvironmentNamespace)
			assert.Equal(t, p.id, ws.WebhookID)
		})
	}
}
