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
//

package api

import (
	"context"
	"testing"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2/mock"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestGetFlagTrigger(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

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

func TestFlagTriggerWebhook(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := createContextWithToken()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})

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
	}{
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
	triggerURL, err := featureService.generateTriggerURL(context.Background(), trigger.Token, false)
	if err != nil {
		t.Errorf("generateTriggerURL() [full] error = %v", err)
	}
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [full] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [full] triggerURL = %v", triggerURL)
	triggerURL, err = featureService.generateTriggerURL(context.Background(), trigger.Token, true)
	if err != nil {
		t.Errorf("generateTriggerURL() [masked] error = %v", err)
	}
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() [masked] triggerURL is empty")
	}
	t.Logf("generateTriggerURL() [masked] triggerURL = %v", triggerURL)
}
