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

	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/feature"
)

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
