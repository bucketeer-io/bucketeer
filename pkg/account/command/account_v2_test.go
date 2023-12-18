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

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	publishermock "github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher/mock"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
)

func TestNewAccountV2CommandHandler(t *testing.T) {
	t.Parallel()
	a := NewAccountV2CommandHandler(nil, nil, nil, "")
	assert.IsType(t, &accountV2CommandHandler{}, a)
}

func TestHandleV2(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*accountV2CommandHandler)
		input       Command
		expectedErr error
	}{
		{
			desc: "CreateAccountV2Command: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.CreateAccountV2Command{},
			expectedErr: nil,
		},
		{
			desc: "ChangeAccountV2NameCommand: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAccountV2NameCommand{},
			expectedErr: nil,
		},
		{
			desc: "ChangeAccountV2AvatarImageURLCommand: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAccountV2AvatarImageUrlCommand{},
			expectedErr: nil,
		},
		{
			desc: "ChangeAccountV2OrganizationRoleCommand: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAccountV2OrganizationRoleCommand{},
			expectedErr: nil,
		},
		{
			desc: "ChangeAccountV2NameCommand: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.ChangeAccountV2NameCommand{},
			expectedErr: nil,
		},
		{
			desc: "EnableAccountV2Command: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.EnableAccountV2Command{},
			expectedErr: nil,
		},
		{
			desc: "DisableAccountV2Command: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.DisableAccountV2Command{},
			expectedErr: nil,
		},
		{
			desc: "DeleteAccountV2Command: success",
			setup: func(h *accountV2CommandHandler) {
				a := domain.NewAccountV2(
					"email",
					"name",
					"avatarImageURL",
					"organizationID",
					accountproto.AccountV2_Role_Organization_MEMBER,
					[]*accountproto.AccountV2_EnvironmentRole{},
				)
				h.account = a
				h.publisher.(*publishermock.MockPublisher).EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			input:       &accountproto.DeleteAccountV2Command{},
			expectedErr: nil,
		},
		{
			desc:        "ErrBadCommand",
			input:       nil,
			expectedErr: ErrBadCommand,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			h := newAccountV2CommandHandlerWithMock(t, mockController)
			if p.setup != nil {
				p.setup(h)
			}
			err := h.Handle(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err, "%s", p.desc)
		})
	}
}

func newAccountV2CommandHandlerWithMock(t *testing.T, mockController *gomock.Controller) *accountV2CommandHandler {
	return &accountV2CommandHandler{
		publisher:      publishermock.NewMockPublisher(mockController),
		organizationID: "org0",
	}
}
