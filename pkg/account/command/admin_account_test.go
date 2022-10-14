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

func TestNewAdminAccountCommandHandler(t *testing.T) {
	t.Parallel()
	a := NewAdminAccountCommandHandler(nil, nil, nil)
	assert.IsType(t, &adminAccountCommandHandler{}, a)
}

func TestHandleAdmin(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*adminAccountCommandHandler)
		input       Command
		expectedErr error
	}{
		{
			desc: "CreateAdminAccountCommand: success",
			setup: func(h *adminAccountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.CreateAdminAccountCommand{},
			expectedErr: nil,
		},
		{
			desc: "EnableAdminAccountCommand: success",
			setup: func(h *adminAccountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.EnableAdminAccountCommand{},
			expectedErr: nil,
		},
		{
			desc:        "ErrBadCommand",
			input:       nil,
			expectedErr: ErrBadCommand,
		},
	}
	for _, p := range patterns {
		h := newAdminAccountCommandHandlerWithMock(t, mockController)
		if p.setup != nil {
			p.setup(h)
		}
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func newAdminAccountCommandHandlerWithMock(t *testing.T, mockController *gomock.Controller) *adminAccountCommandHandler {
	return &adminAccountCommandHandler{
		publisher: publishermock.NewMockPublisher(mockController),
	}
}
