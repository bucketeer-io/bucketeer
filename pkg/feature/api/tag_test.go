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

package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	mysqlmock "github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql/mock"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestListTagsMySQL(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	ctx := createContextWithToken()
	service := createFeatureServiceNew(mockController)
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	createError := func(msg string, status *status.Status) error {
		status, err := status.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: msg,
		})
		require.NoError(t, err)
		return status.Err()
	}

	patterns := []struct {
		desc        string
		setup       func(*FeatureService)
		input       *featureproto.ListTagsRequest
		expected    *featureproto.ListTagsResponse
		expectedErr error
	}{
		{
			desc:        "errInvalidCursor",
			setup:       nil,
			input:       &featureproto.ListTagsRequest{EnvironmentNamespace: environmentNamespace, Cursor: "foo"},
			expected:    nil,
			expectedErr: createError(localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"), statusInvalidCursor),
		},
		{
			desc: "errInternal",
			setup: func(fs *FeatureService) {
				fs.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, errors.New("test"))
			},
			input:       &featureproto.ListTagsRequest{EnvironmentNamespace: environmentNamespace},
			expected:    nil,
			expectedErr: createError(localizer.MustLocalize(locale.InternalServerError), statusInternal),
		},
		{
			desc: "success",
			setup: func(fs *FeatureService) {
				rows := mysqlmock.NewMockRows(mockController)
				rows.EXPECT().Close().Return(nil)
				rows.EXPECT().Next().Return(false)
				rows.EXPECT().Err().Return(nil)
				fs.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(rows, nil)
				row := mysqlmock.NewMockRow(mockController)
				row.EXPECT().Scan(gomock.Any()).Return(nil)
				fs.mysqlClient.(*mysqlmock.MockClient).EXPECT().QueryRowContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(row)
			},
			input: &featureproto.ListTagsRequest{
				PageSize:             2,
				Cursor:               "",
				EnvironmentNamespace: environmentNamespace,
			},
			expected:    &featureproto.ListTagsResponse{Tags: []*featureproto.Tag{}, Cursor: "0"},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(service)
			}
			actual, err := service.ListTags(ctx, p.input)
			assert.Equal(t, p.expected, actual, p.desc)
			assert.Equal(t, p.expectedErr, err, p.desc)
		})
	}
}

func TestUpsertTags(t *testing.T) {
	t.Parallel()
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	service := createFeatureServiceNew(mockController)
	ctx := createContextWithToken()
	transaction := mysqlmock.NewMockTransaction(mockController)
	internalErr := errors.New("test")
	patterns := []struct {
		desc        string
		tags        []string
		setup       func(*mysqlmock.MockTransaction)
		expectedErr error
	}{
		{
			desc: "error: internal error when creating tag",
			tags: []string{"tag"},
			setup: func(mt *mysqlmock.MockTransaction) {
				mt.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, internalErr)
			},
			expectedErr: internalErr,
		},
		{
			desc:        "success: tag with whitespaces",
			tags:        []string{" "},
			expectedErr: nil,
		},
		{
			desc: "success: create new tag",
			tags: []string{"tag"},
			setup: func(mt *mysqlmock.MockTransaction) {
				mt.EXPECT().ExecContext(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
			},
			expectedErr: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			if p.setup != nil {
				p.setup(transaction)
			}
			err := service.upsertTags(ctx, transaction, p.tags, environmentNamespace)
			assert.Equal(t, p.expectedErr, err)
		})
	}
}
