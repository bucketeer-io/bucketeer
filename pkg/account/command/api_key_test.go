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

package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAPIKeyCommandHandler(t *testing.T) {
	t.Parallel()
	a := NewAPIKeyCommandHandler(nil, nil, nil, "")
	assert.IsType(t, &apiKeyCommandHandler{}, a)
}

func newAPIKeyCommandHandlerWithMock(t *testing.T, mockController *gomock.Controller) *apiKeyCommandHandler {
	return &apiKeyCommandHandler{
		publisher: publishermock.NewMockPublisher(mockController),
	}
}

func TestAPIKeyHandle(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := map[string]struct {
		setup       func(*apiKeyCommandHandler)
		input       Command
		expectedErr error
	}{
		"CreateAPIKeyCommand: success": {
			setup: func(h *apiKeyCommandHandler) {
				a, err := domain.NewAPIKey("email", accountproto.APIKey_SDK)
				require.NoError(t, err)
				h.apiKey = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.CreateAPIKeyCommand{},
			expectedErr: nil,
		},
		"ChangeAPIKeyNameCommand: success": {
			setup: func(h *apiKeyCommandHandler) {
				a, err := domain.NewAPIKey("email", accountproto.APIKey_SDK)
				require.NoError(t, err)
				h.apiKey = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAPIKeyNameCommand{},
			expectedErr: nil,
		},
		"EnableAPIKeyCommand: success": {
			setup: func(h *apiKeyCommandHandler) {
				a, err := domain.NewAPIKey("email", accountproto.APIKey_SDK)
				require.NoError(t, err)
				h.apiKey = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.EnableAPIKeyCommand{},
			expectedErr: nil,
		},
		"DisableAPIKeyCommand: success": {
			setup: func(h *apiKeyCommandHandler) {
				a, err := domain.NewAPIKey("email", accountproto.APIKey_SDK)
				require.NoError(t, err)
				h.apiKey = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.DisableAPIKeyCommand{},
			expectedErr: nil,
		},
		"ErrBadCommand": {
			input:       nil,
			expectedErr: ErrBadCommand,
		},
	}
	for msg, p := range patterns {
		h := newAPIKeyCommandHandlerWithMock(t, mockController)
		if p.setup != nil {
			p.setup(h)
		}
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expectedErr, err, "%s", msg)
	}
}
