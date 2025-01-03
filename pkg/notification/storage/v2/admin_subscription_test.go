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

func TestNewAdminSubscriptionStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	storage := NewAdminSubscriptionStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &adminSubscriptionStorage{}, storage)
}

func TestCreateAdminSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	sourceTypes := []proto.Subscription_SourceType{5}
	recipient := &proto.Recipient{Type: 0, SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "slack"}}
	name := "name-0"
	patterns := []struct {
		desc        string
		setup       func(*adminSubscriptionStorage)
		input       *domain.Subscription
		expectedErr error
	}{
		{
			desc: "ErrAdminSubscriptionAlreadyExists",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			expectedErr: ErrAdminSubscriptionAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					id, int64(1), int64(2), false, mysql.JSONObject{Val: sourceTypes}, mysql.JSONObject{Val: recipient}, name,
				).Return(nil, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id, CreatedAt: 1, UpdatedAt: 2, Disabled: false, SourceTypes: sourceTypes, Recipient: recipient, Name: name},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminSubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateAdminSubscription(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateAdminSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	id := "id-0"
	sourceTypes := []proto.Subscription_SourceType{5}
	recipient := &proto.Recipient{Type: 0, SlackChannelRecipient: &proto.SlackChannelRecipient{WebhookUrl: "slack"}}
	name := "name-0"
	patterns := []struct {
		desc        string
		setup       func(*adminSubscriptionStorage)
		input       *domain.Subscription
		expectedErr error
	}{
		{
			desc: "ErrAdminSubscriptionUnexpectedAffectedRows",
			setup: func(s *adminSubscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			expectedErr: ErrAdminSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id},
			},
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *adminSubscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					int64(2), false, mysql.JSONObject{Val: sourceTypes}, mysql.JSONObject{Val: recipient}, name, id,
				).Return(result, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: id, CreatedAt: 1, UpdatedAt: 2, Disabled: false, SourceTypes: sourceTypes, Recipient: recipient, Name: name},
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminSubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateAdminSubscription(context.Background(), p.input)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteAdminSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*adminSubscriptionStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "ErrAdminSubscriptionUnexpectedAffectedRows",
			setup: func(s *adminSubscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(0), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:          "id-0",
			expectedErr: ErrAdminSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			id:          "id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *adminSubscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(),
					gomock.Any(),
					"id-0",
				).Return(result, nil)
			},
			id:          "id-0",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminSubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteAdminSubscription(context.Background(), p.id)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetAdminSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc        string
		setup       func(*adminSubscriptionStorage)
		id          string
		expectedErr error
	}{
		{
			desc: "ErrAdminSubscriptionNotFound",
			setup: func(s *adminSubscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(mysql.ErrNoRows)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:          "id-0",
			expectedErr: ErrAdminSubscriptionNotFound,
		},
		{
			desc: "Error",
			setup: func(s *adminSubscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(errors.New("error"))
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)

			},
			id:          "id-0",
			expectedErr: errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *adminSubscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					"id-0",
				).Return(row)
			},
			id:          "id-0",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminSubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			result, err := storage.GetAdminSubscription(context.Background(), p.id)

			if result != nil {
				assert.IsType(t, &domain.Subscription{}, result)
			}
			assert.Equal(t, p.expectedErr, err)

		})
	}
}

func TestListAdminSubscriptions(t *testing.T) {
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
		setup          func(*adminSubscriptionStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expectedCount  int
		expectedCursor int
		expectedErr    error
	}{
		{
			desc: "Error",
			setup: func(s *adminSubscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))
			},
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *adminSubscriptionStorage) {
				var nextCallCount = 0
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().DoAndReturn(func() bool {
					nextCallCount++
					return nextCallCount <= getSize
				}).Times(getSize + 1)
				rows.EXPECT().Err().Return(nil)
				rows.EXPECT().Scan(gomock.Any()).Return(nil).Times(getSize)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(),
					gomock.Any(),
					updatedAt, disable,
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(),
					gomock.Any(),
					updatedAt, disable,
				).Return(row)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("updated_at", ">=", updatedAt),
				mysql.NewFilter("disabled", "=", disable),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
				mysql.NewOrder("create_at", mysql.OrderDirectionDesc),
			},
			limit:          limit,
			offset:         offset,
			expectedCount:  getSize,
			expectedCursor: offset + getSize,
			expectedErr:    nil,
		},
		{
			desc: "Success:No wereParts and no orderParts and no limit and no offset",
			setup: func(s *adminSubscriptionStorage) {
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
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expectedCount:  0,
			expectedCursor: 0,
			expectedErr:    nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newAdminSubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			subscriptions, cursor, _, err := storage.ListAdminSubscriptions(
				context.Background(),
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expectedCount, len(subscriptions))
			if len(subscriptions) > 0 {
				assert.IsType(t, []*proto.Subscription{}, subscriptions)
			}
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newAdminSubscriptionStorageWithMock(t *testing.T, mockController *gomock.Controller) *adminSubscriptionStorage {
	t.Helper()
	return &adminSubscriptionStorage{mock.NewMockQueryExecer(mockController)}
}
