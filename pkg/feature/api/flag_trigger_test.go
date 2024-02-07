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
//

package api

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestCreateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.CreateFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:  "Error Invalid Argument",
			setup: nil,
			input: &proto.CreateFlagTriggerRequest{
				EnvironmentNamespace: "namespace",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "create_flag_trigger_command")),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			input: &proto.CreateFlagTriggerRequest{
				EnvironmentNamespace: "namespace",
				CreateFlagTriggerCommand: &proto.CreateFlagTriggerCommand{
					FeatureId: "id-1",
					Type:      proto.FlagTrigger_Type_WEBHOOK,
					Action:    proto.FlagTrigger_Action_ON,
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalize(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.CreateFlagTriggerRequest{
				EnvironmentNamespace: "namespace",
				CreateFlagTriggerCommand: &proto.CreateFlagTriggerCommand{
					FeatureId: "id-1",
					Type:      proto.FlagTrigger_Type_WEBHOOK,
					Action:    proto.FlagTrigger_Action_ON,
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.CreateFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if resp != nil {
				assert.True(t, len(resp.FlagTrigger.Id) > 0)
				assert.Equal(t, p.input.CreateFlagTriggerCommand.FeatureId, resp.FlagTrigger.FeatureId)
				assert.Equal(t, p.input.CreateFlagTriggerCommand.Type, resp.FlagTrigger.Type)
				assert.Equal(t, p.input.CreateFlagTriggerCommand.Action, resp.FlagTrigger.Action)
				assert.True(t, len(resp.Url) > 0)
			}
		})
	}
}

func TestGetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	baseFlagTrigger := &proto.FlagTrigger{
		Id:                   "1",
		FeatureId:            "featureId",
		EnvironmentNamespace: "namespace",
		Type:                 proto.FlagTrigger_Type_WEBHOOK,
		Action:               proto.FlagTrigger_Action_ON,
		Description:          "base",
		TriggerCount:         100,
		LastTriggeredAt:      500,
		Token:                "test-token",
		Disabled:             false,
		CreatedAt:            200,
		UpdatedAt:            300,
	}

	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.GetFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:        "Error Validate",
			setup:       nil,
			input:       &proto.GetFlagTriggerRequest{},
			expectedErr: createError(statusMissingTriggerID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate)),
		},
		{
			desc: "Error Not Found",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, v2fs.ErrFlagTriggerNotFound)
			},
			input:       &proto.GetFlagTriggerRequest{Id: "1", EnvironmentNamespace: "namespace"},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},

		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTrigger(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: baseFlagTrigger,
				}, nil)
			},
			input:       &proto.GetFlagTriggerRequest{Id: baseFlagTrigger.Id, EnvironmentNamespace: baseFlagTrigger.EnvironmentNamespace},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.GetFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestUpdateFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.UpdateFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:  "Error Invalid Argument",
			setup: func(s *FeatureService) {},
			input: &proto.UpdateFlagTriggerRequest{
				Id:                   "id",
				EnvironmentNamespace: "namespace",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "change_flag_trigger_description_command")),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			input: &proto.UpdateFlagTriggerRequest{
				Id:                   "id",
				EnvironmentNamespace: "namespace",
				ChangeFlagTriggerDescriptionCommand: &proto.ChangeFlagTriggerDescriptionCommand{
					Description: "description",
				},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.UpdateFlagTriggerRequest{
				Id:                   "id",
				EnvironmentNamespace: "namespace",
				ChangeFlagTriggerDescriptionCommand: &proto.ChangeFlagTriggerDescriptionCommand{
					Description: "description",
				},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.UpdateFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestEnableFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.EnableFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:  "Error Invalid Argument",
			setup: func(s *FeatureService) {},
			input: &proto.EnableFlagTriggerRequest{
				Id:                   "id",
				EnvironmentNamespace: "namespace",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError)),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			input: &proto.EnableFlagTriggerRequest{
				Id:                       "id",
				EnvironmentNamespace:     "namespace",
				EnableFlagTriggerCommand: &proto.EnableFlagTriggerCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.EnableFlagTriggerRequest{
				Id:                       "id",
				EnvironmentNamespace:     "namespace",
				EnableFlagTriggerCommand: &proto.EnableFlagTriggerCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.EnableFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestDisableFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.DisableFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:  "Error Invalid Argument",
			setup: func(s *FeatureService) {},
			input: &proto.DisableFlagTriggerRequest{
				Id:                   "id",
				EnvironmentNamespace: "namespace",
			},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError)),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			input: &proto.DisableFlagTriggerRequest{
				Id:                        "id",
				EnvironmentNamespace:      "namespace",
				DisableFlagTriggerCommand: &proto.DisableFlagTriggerCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.DisableFlagTriggerRequest{
				Id:                        "id",
				EnvironmentNamespace:      "namespace",
				DisableFlagTriggerCommand: &proto.DisableFlagTriggerCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.DisableFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestResetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.ResetFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:        "Error Invalid Argument",
			setup:       func(s *FeatureService) {},
			input:       &proto.ResetFlagTriggerRequest{},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError)),
		},
		{
			desc: "Error GetFlagTrigger",
			setup: func(s *FeatureService) {
				row := mysqlmock.NewMockRow(mockController)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
			},
			input: &proto.ResetFlagTriggerRequest{
				ResetFlagTriggerCommand: &proto.ResetFlagTriggerCommand{},
			},
			expectedErr: createError(statusNotFound, localizer.MustLocalize(locale.NotFoundError)),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &proto.ResetFlagTriggerRequest{
				ResetFlagTriggerCommand: &proto.ResetFlagTriggerCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &proto.ResetFlagTriggerRequest{
				ResetFlagTriggerCommand: &proto.ResetFlagTriggerCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.ResetFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if resp != nil {
				assert.Equal(t, resp.FlagTrigger.Token, "")
			}
		})
	}
}

func TestDeleteFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.DeleteFlagTriggerRequest
		expectedErr error
	}{
		{
			desc:        "Error Invalid Argument",
			setup:       func(s *FeatureService) {},
			input:       &proto.DeleteFlagTriggerRequest{},
			expectedErr: createError(statusMissingCommand, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError)),
		},
		{
			desc: "Error Internal",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
			},
			input: &proto.DeleteFlagTriggerRequest{
				DeleteFlagTriggerCommand: &proto.DeleteFlagTriggerCommand{},
			},
			expectedErr: createError(statusInternal, localizer.MustLocalizeWithTemplate(locale.InternalServerError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input: &proto.DeleteFlagTriggerRequest{
				DeleteFlagTriggerCommand: &proto.DeleteFlagTriggerCommand{},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.DeleteFlagTrigger(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestFlagTriggerWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
	createError := func(status *gstatus.Status, msg string) error {
		st, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return st.Err()
	}

	baseFlagTrigger := &proto.FlagTrigger{
		Id:                   "1",
		FeatureId:            "featureId",
		EnvironmentNamespace: "namespace",
		Type:                 proto.FlagTrigger_Type_WEBHOOK,
		Action:               proto.FlagTrigger_Action_ON,
		Description:          "base",
		TriggerCount:         100,
		LastTriggeredAt:      500,
		Token:                "test-token",
		Disabled:             false,
		CreatedAt:            200,
		UpdatedAt:            300,
	}

	patterns := []struct {
		desc        string
		setup       func(service *FeatureService)
		input       *proto.FlagTriggerWebhookRequest
		expectedErr error
	}{{
		desc:        "Error Invalid Argument",
		setup:       func(s *FeatureService) {},
		input:       &proto.FlagTriggerWebhookRequest{},
		expectedErr: createError(statusSecretRequired, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "secret")),
	},
		{
			desc: "Error Not Found",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTriggerByToken(
					gomock.Any(), gomock.Any(),
				).Return(nil, v2fs.ErrFlagTriggerNotFound)
			},
			input:       &proto.FlagTriggerWebhookRequest{Token: "token"},
			expectedErr: createError(statusTriggerNotFound, localizer.MustLocalizeWithTemplate(locale.NotFoundError)),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				s.flagTriggerStorage.(*mock.MockFlagTriggerStorage).EXPECT().GetFlagTriggerByToken(
					gomock.Any(), gomock.Any(),
				).Return(&domain.FlagTrigger{
					FlagTrigger: baseFlagTrigger,
				}, nil)

				s.featureStorage.(*mock.MockFeatureStorage).EXPECT().GetFeature(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(&domain.Feature{
					Feature: &proto.Feature{
						Id:      "id",
						Name:    "test feature",
						Version: 1,
						Enabled: true,
					},
				}, nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().BeginTx(gomock.Any()).Return(nil, nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().RunInTransaction(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil)
			},
			input:       &proto.FlagTriggerWebhookRequest{Token: "token"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			resp, err := s.FlagTriggerWebhook(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			if err == nil {
				assert.NotNil(t, resp)
			}
		})
	}
}

func TestListFlagTriggers(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := metadata.NewIncomingContext(
		createContextWithToken(),
		metadata.MD{"accept-language": []string{"ja"}},
	)
	localizer := locale.NewLocalizer(ctx)
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
		setup       func(service *FeatureService)
		input       *proto.ListFlagTriggersRequest
		expected    *proto.ListFlagTriggersResponse
		expectedErr error
	}{
		{
			desc:        "Error Validate",
			setup:       nil,
			input:       &proto.ListFlagTriggersRequest{},
			expected:    nil,
			expectedErr: createError(statusMissingTriggerFeatureID, localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate)),
		},
		{
			desc:        "Error Invalid Argument",
			setup:       nil,
			input:       &proto.ListFlagTriggersRequest{FeatureId: "1", Cursor: "XXX"},
			expected:    nil,
			expectedErr: createError(statusInvalidCursor, localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor")),
		},
		{
			desc: "Success",
			setup: func(s *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Close().Return(nil)
				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)

				s.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input:       &proto.ListFlagTriggersRequest{FeatureId: "1", PageSize: 2, Cursor: ""},
			expected:    &proto.ListFlagTriggersResponse{FlagTriggers: []*proto.ListFlagTriggersResponse_FlagTriggerWithUrl{}, Cursor: "0"},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			s := createFeatureServiceNew(mockController)
			if p.setup != nil {
				p.setup(s)
			}
			actual, err := s.ListFlagTriggers(ctx, p.input)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, actual)
		})
	}
}

func TestFeatureServiceGenerateTriggerURL(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureService := createFeatureServiceNew(mockController)
	trigger, err := domain.NewFlagTrigger(
		"test",
		&proto.CreateFlagTriggerCommand{
			FeatureId:   "test",
			Type:        proto.FlagTrigger_Type_WEBHOOK,
			Action:      proto.FlagTrigger_Action_ON,
			Description: "test",
		},
	)
	if err != nil {
		t.Errorf("NewFlagTrigger() error = %v", err)
	}
	err = trigger.GenerateToken()
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}
	t.Logf("GenerateToken() token = %v", trigger.Token)
	triggerURL := featureService.generateTriggerURL(context.Background(), trigger.Token, false)
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [full] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [full] triggerURL = %v", triggerURL)
	triggerURL = featureService.generateTriggerURL(context.Background(), trigger.Token, true)
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [masked] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [masked] triggerURL = %v", triggerURL)
}
