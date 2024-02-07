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
	storage := NewSubscriptionStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &subscriptionStorage{}, storage)
}

func TestCreateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*subscriptionStorage)
		input                *domain.Subscription
		environmentNamespace string
		expectedErr          error
	}{
		{
			desc: "ErrSubscriptionAlreadyExists",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, mysql.ErrDuplicateEntry)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrSubscriptionAlreadyExists,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreateSubscription(context.Background(), p.input, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestUpdateSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*subscriptionStorage)
		input                *domain.Subscription
		environmentNamespace string
		expectedErr          error
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
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          ErrSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			input: &domain.Subscription{
				Subscription: &proto.Subscription{Id: "id-0"},
			},
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.UpdateSubscription(context.Background(), p.input, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDeleteSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*subscriptionStorage)
		id                   string
		environmentNamespace string
		expectedErr          error
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
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          ErrSubscriptionUnexpectedAffectedRows,
		},
		{
			desc: "Error",
			setup: func(s *subscriptionStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("error"))

			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				result := mock.NewMockResult(mockController)
				result.EXPECT().RowsAffected().Return(int64(1), nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteSubscription(context.Background(), p.id, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestGetSubscription(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc                 string
		setup                func(*subscriptionStorage)
		id                   string
		environmentNamespace string
		expectedErr          error
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
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          ErrSubscriptionNotFound,
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
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			id:                   "id-0",
			environmentNamespace: "ns",
			expectedErr:          nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newsubscriptionStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			_, err := storage.GetSubscription(context.Background(), p.id, p.environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestListSubscriptions(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	patterns := []struct {
		desc           string
		setup          func(*subscriptionStorage)
		whereParts     []mysql.WherePart
		orders         []*mysql.Order
		limit          int
		offset         int
		expected       []*proto.Subscription
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
			whereParts:     nil,
			orders:         nil,
			limit:          0,
			offset:         0,
			expected:       nil,
			expectedCursor: 0,
			expectedErr:    errors.New("error"),
		},
		{
			desc: "Success",
			setup: func(s *subscriptionStorage) {
				rows := mock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				s.qe.(*mock.MockQueryExecer).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			whereParts: []mysql.WherePart{
				mysql.NewFilter("num", ">=", 5),
			},
			orders: []*mysql.Order{
				mysql.NewOrder("id", mysql.OrderDirectionAsc),
			},
			limit:          10,
			offset:         5,
			expected:       []*proto.Subscription{},
			expectedCursor: 5,
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
				p.whereParts,
				p.orders,
				p.limit,
				p.offset,
			)
			assert.Equal(t, p.expected, subscriptions)
			assert.Equal(t, p.expectedCursor, cursor)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newsubscriptionStorageWithMock(t *testing.T, mockController *gomock.Controller) *subscriptionStorage {
	t.Helper()
	return &subscriptionStorage{mock.NewMockQueryExecer(mockController)}
}
