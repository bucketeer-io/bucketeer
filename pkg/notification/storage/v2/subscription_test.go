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

package v2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/notification/domain"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	proto "github.com/bucketeer-io/bucketeer/proto/notification"
)

func TestNewSubscriptionStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	client := mock.NewMockClient(mockController)
	storage := NewSubscriptionStorage(client)
	assert.IsType(t, &subscriptionStorage{}, storage)
}

func TestCreateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	sourceTypes := []proto.Subscription_SourceType{5}
	recipient := &proto.Recipient{Type: 0, SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "slack"}}
	name := "name-0"
	featureFlagTags := []string{"tag"}
	envNamespace := "ns"
	patterns := []struct {
		desc          string
		setup         func(*subscriptionStorage)
		input         *domain.Subscription
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrSubscriptionAlreadyExists",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			environmentId: "ns",
			expectedErr:   ErrSubscriptionAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			environmentId: "ns",
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					id, int64(1), int64(2), false, mysql.JSONObject{Val: sourceTypes}, mysql.JSONObject{Val: recipient}, name, mysql.JSONObject{Val: featureFlagTags}, envNamespace,
				).Return(nil, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{
					Id:          id,
					CreatedAt:   1,
					UpdatedAt:   2,
					Disabled:    false,
					SourceTypes: sourceTypes,
					Recipient:   recipient,
					Name:        name, FeatureFlagTags: featureFlagTags,
				},
			},
			environmentId: envNamespace,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateSubscription(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	sourceTypes := []proto.Subscription_SourceType{5}
	recipient := &proto.Recipient{Type: 0, SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "slack"}}
	name := "name-0"
	featureFlagTags := []string{"tag"}
	envNamespace := "ns"
	patterns := []struct {
		desc          string
		setup         func(*subscriptionStorage)
		input         *domain.Subscription
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrSubscriptionUnexpectedAffectedRows",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			environmentId: envNamespace,
			expectedErr:   ErrSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			environmentId: envNamespace,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					int64(2), false, mysql.JSONObject{Val: sourceTypes}, mysql.JSONObject{Val: recipient}, name, mysql.JSONObject{Val: featureFlagTags}, id, envNamespace,
				).Return(result, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{
					Id:              id,
					CreatedAt:       1,
					UpdatedAt:       2,
					Disabled:        false,
					SourceTypes:     sourceTypes,
					Recipient:       recipient,
					Name:            name,
					FeatureFlagTags: featureFlagTags,
				},
			},
			environmentId: envNamespace,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateSubscription(context.Background(), p.input, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	envNamespace := "ns"
	patterns := []struct {
		desc          string
		setup         func(*subscriptionStorage)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrSubscriptionUnexpectedAffectedRows",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   ErrSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					id, envNamespace,
				).Return(result, nil)
			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteSubscription(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	envNamespace := "ns"
	patterns := []struct {
		desc          string
		setup         func(*subscriptionStorage)
		id            string
		environmentId string
		expectedErr   error
	}{
		{
			desc: "ErrSubscriptionNotFound",
			setup: func(s *subscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   ErrSubscriptionNotFound,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					id, envNamespace,
				).Return(row)
			},
			id:            id,
			environmentId: envNamespace,
			expectedErr:   nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetSubscription(context.Background(), p.id, p.environmentId)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListSubscriptions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	getSize := 2
	offset := 5
	limit := 10
	updatedAt := 8
	disable := false
	patterns := []struct {
		desc           string
		setup          func(*subscriptionStorage)
		listOpts       *mysql.ListOptions
		expectedCount  int
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			listOpts:       nil,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= getSize
				}).Times(getSize + 1)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(getSize)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					updatedAt, disable,
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					updatedAt, disable,
				).Return(row)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			listOpts: &mysql.ListOptions{
				Limit:  limit,
				Offset: offset,
				Filters: []*mysql.FilterV2{
					{
						Column:   "updated_at",
						Operator: mysql.OperatorGreaterThanOrEqual,
						Value:    updatedAt,
					},
					{
						Column:   "disabled",
						Operator: mysql.OperatorEqual,
						Value:    disable,
					},
				},
				InFilters:   nil,
				NullFilters: nil,
				JSONFilters: nil,
				SearchQuery: nil,
				Orders: []*mysql.Order{
					{
						Column:    "id",
						Direction: mysql.OrderDirectionAsc,
					},
					{
						Column:    "create_at",
						Direction: mysql.OrderDirectionDesc,
					},
				},
			},
			expectedCount:  getSize,
			expectedCursor: offset + getSize,
			expectedErr:    nil,
		},
		{
			desc: "Success:No wereParts and no orderParts and no limit and no offset",
			setup: func(s *subscriptionStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					[]interface{}{},
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					[]interface{}{},
				).Return(row)
			},
			listOpts:       nil,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			subscriptions, cursor, _, err := storage.ListSubscriptions(
				context.Background(),
				p.listOpts,
			)
			assert.Equal(t, p.expectedCount, len(subscriptions))
			if subscriptions != nil {
				assert.IsType(t, []*proto.Subscription{}, subscriptions)
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newsubscriptionStorageWithMock(t *testing.T, mockController *gomock.Controller) *subscriptionStorage {
	t.Helper()
	return &subscriptionStorage{mock.NewMockQueryExecer(mockController)}
}
