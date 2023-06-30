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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAccountCommandHandler(t *testing.T) {
	t.Parallel()
	a := NewAccountCommandHandler(nil, nil, nil, "")
	assert.IsType(t, &accountCommandHandler{}, a)
}

func TestHandle(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountCommandHandler)
		input       Command
		expectedErr error
	}{
		{
			desc: "CreateAccountCommand: success",
			setup: func(h *accountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.CreateAccountCommand{},
			expectedErr: nil,
		},
		{
			desc: "ChangeAccountRoleCommand: success",
			setup: func(h *accountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAccountRoleCommand{},
			expectedErr: nil,
		},
		{
			desc: "EnableAccountCommand: success",
			setup: func(h *accountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.EnableAccountCommand{},
			expectedErr: nil,
		},
		{
			desc: "DisableAccountCommand: success",
			setup: func(h *accountCommandHandler) {
				a, err := domain.NewAccount("email", accountproto.Account_VIEWER)
				require.NoError(t, err)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.DisableAccountCommand{},
			expectedErr: nil,
		},
		{
			desc:        "ErrBadCommand",
			input:       nil,
			expectedErr: ErrBadCommand,
		},
	}
	for _, p := range patterns {
		h := newAccountCommandHandlerWithMock(t, mockController)
		if p.setup != nil {
			p.setup(h)
		}
		err := h.Handle(context.Background(), p.input)
		assert.Equal(t, p.expectedErr, err, "%s", p.desc)
	}
}

func newAccountCommandHandlerWithMock(t *testing.T, mockController *gomock.Controller) *accountCommandHandler {
	return &accountCommandHandler{
		publisher:            publishermock.NewMockPublisher(mockController),
		environmentNamespace: "ns0",
	}
}
