package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
)

func TestNewMAUStorage(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	db := NewMAUStorage(mock.NewMockQueryExecer(mockController))
	assert.IsType(t, &mauStorage{}, db)
}

func TestDeleteRecords(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*mauStorage)
		inputTarget string
		expectedErr error
	}{
		{
			desc: "fail",
			setup: func(s *mauStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			inputTarget: "invalid year month",
			expectedErr: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(s *mauStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			inputTarget: "202301",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newMAUStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DeleteRecords(context.Background(), p.inputTarget)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestRebuildPartition(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*mauStorage)
		inputTarget string
		expectedErr error
	}{
		{
			desc: "fail",
			setup: func(s *mauStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			inputTarget: "invalid year month",
			expectedErr: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(s *mauStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			inputTarget: "202301",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newMAUStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.RebuildPartition(context.Background(), p.inputTarget)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestDropPartition(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*mauStorage)
		inputTarget string
		expectedErr error
	}{
		{
			desc: "fail",
			setup: func(s *mauStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			inputTarget: "invalid year month",
			expectedErr: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(s *mauStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			inputTarget: "202301",
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newMAUStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.DropPartition(context.Background(), p.inputTarget)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func TestCreatePartition(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	patterns := []struct {
		desc        string
		setup       func(*mauStorage)
		expectedErr error
	}{
		{
			desc: "fail",
			setup: func(s *mauStorage) {
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			expectedErr: errors.New("test"),
		},
		{
			desc: "success",
			setup: func(s *mauStorage) {
				result := mock.NewMockResult(mockController)
				s.qe.(*mock.MockQueryExecer).EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(result, nil)
			},
			expectedErr: nil,
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			storage := newMAUStorageWithMock(t, mockController)
			if p.setup != nil {
				p.setup(storage)
			}
			err := storage.CreatePartition(context.Background(), "testp", "testlt")
			assert.Equal(t, p.expectedErr, err)
		})
	}
}

func newMAUStorageWithMock(t *testing.T, m *gomock.Controller) *mauStorage {
	t.Helper()
	return &mauStorage{mock.NewMockQueryExecer(m)}
}
