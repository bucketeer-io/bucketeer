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
)

func TestFeatureServiceGenerateTriggerURL(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	featureService := createFeatureServiceNew(mockController)
	triggerURL, err := featureService.generateTriggerURL(context.Background(), nil, true)
	if err != nil {
		t.Errorf("generateTriggerURL() error = %v", err)
	}
	if triggerURL == "" {
		t.Errorf("generateTriggerURL() triggerURL is empty")
	}
	t.Logf("generateTriggerURL() triggerURL = %v", triggerURL)

}
