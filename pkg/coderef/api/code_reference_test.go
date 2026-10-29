// Copyright 2026 The Bucketeer Authors.
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
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage"
	storagemock "github.com/bucketeer-io/bucketeer/v2/pkg/coderef/storage/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/rpc"
	dbmock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/token"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	proto "github.com/bucketeer-io/bucketeer/v2/proto/coderef"
)

const (
	testEditorEmail     = "editor@example.com"
	testEnvironmentID   = "env-a"
	testCodeReferenceID = "cr-id"
)

func createCodeReferenceService(t *testing.T, mockController *gomock.Controller) *CodeReferenceService {
	t.Helper()
	return &CodeReferenceService{
		accountClient:  accountclientmock.NewMockClient(mockController),
		dbClient:       dbmock.NewMockClient(mockController),
		codeRefStorage: storagemock.NewMockCodeReferenceStorage(mockController),
		publisher:      publishermock.NewMockPublisher(mockController),
		logger:         zap.NewNop().Named("api"),
	}
}

func createContextWithToken(t *testing.T, email string) context.Context {
	t.Helper()
	tok := &token.AccessToken{
		Issuer:   "issuer",
		Audience: "audience",
		Expiry:   time.Now().AddDate(100, 0, 0),
		IssuedAt: time.Now(),
		Email:    email,
	}
	return context.WithValue(context.TODO(), rpc.AccessTokenKey, tok)
}

// expectEditorInEnvironment makes the caller an environment EDITOR in the given
// environment only (a plain organization MEMBER otherwise).
func expectEditorInEnvironment(s *CodeReferenceService, email, environmentID string) {
	s.accountClient.(*accountclientmock.MockClient).EXPECT().
		GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).
		Return(&accountproto.GetAccountV2ByEnvironmentIDResponse{
			Account: &accountproto.AccountV2{
				Email:            email,
				OrganizationId:   "org0",
				OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
				EnvironmentRoles: []*accountproto.AccountV2_EnvironmentRole{
					{EnvironmentId: environmentID, Role: accountproto.AccountV2_Role_Environment_EDITOR},
				},
			},
		}, nil)
}

func codeReferenceInEnvironment(environmentID string) *domain.CodeReference {
	return &domain.CodeReference{CodeReference: proto.CodeReference{
		Id:               testCodeReferenceID,
		FeatureId:        "feature-1",
		FilePath:         "src/main.go",
		FileExtension:    "go",
		LineNumber:       10,
		CodeSnippet:      "client.BoolVariation(\"feature-1\", false)",
		ContentHash:      "hash-1",
		RepositoryName:   "repo",
		RepositoryOwner:  "owner",
		RepositoryType:   proto.CodeReference_GITHUB,
		RepositoryBranch: "main",
		CommitHash:       "abc123",
		EnvironmentId:    environmentID,
		CreatedAt:        1,
		UpdatedAt:        1,
	}}
}

func TestGetCodeReferenceScopedToEnvironment(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		setup   func(*CodeReferenceService)
		wantErr error
	}{
		{
			name: "returns not found when the id does not belong to the requested environment",
			setup: func(s *CodeReferenceService) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(nil, storage.ErrCodeReferenceNotFound)
			},
			wantErr: statusNotFound.Err(),
		},
		{
			name: "returns the code reference that belongs to the requested environment",
			setup: func(s *CodeReferenceService) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(codeReferenceInEnvironment(testEnvironmentID), nil)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createCodeReferenceService(t, mockController)
			expectEditorInEnvironment(service, testEditorEmail, testEnvironmentID)
			tc.setup(service)

			resp, err := service.GetCodeReference(
				createContextWithToken(t, testEditorEmail),
				&proto.GetCodeReferenceRequest{
					Id:            testCodeReferenceID,
					EnvironmentId: testEnvironmentID,
				},
			)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp.CodeReference)
			assert.Equal(t, testCodeReferenceID, resp.CodeReference.Id)
			assert.Equal(t, testEnvironmentID, resp.CodeReference.EnvironmentId)
			assert.NotEmpty(t, resp.CodeReference.SourceUrl)
		})
	}
}

func TestUpdateCodeReferenceScopedToEnvironment(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		setup   func(*CodeReferenceService, **domain.CodeReference)
		wantErr error
	}{
		{
			name: "returns not found and writes nothing when the id does not belong to the requested environment",
			setup: func(s *CodeReferenceService, _ **domain.CodeReference) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(nil, storage.ErrCodeReferenceNotFound)
			},
			wantErr: statusNotFound.Err(),
		},
		{
			name: "updates the code reference that belongs to the requested environment",
			setup: func(s *CodeReferenceService, persisted **domain.CodeReference) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(codeReferenceInEnvironment(testEnvironmentID), nil)
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					UpdateCodeReference(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, cr *domain.CodeReference) error {
						*persisted = cr
						return nil
					})
				s.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createCodeReferenceService(t, mockController)
			expectEditorInEnvironment(service, testEditorEmail, testEnvironmentID)
			service.dbClient.(*dbmock.MockClient).EXPECT().
				RunInTransactionV2(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})
			var persisted *domain.CodeReference
			tc.setup(service, &persisted)

			resp, err := service.UpdateCodeReference(
				createContextWithToken(t, testEditorEmail),
				&proto.UpdateCodeReferenceRequest{
					Id:            testCodeReferenceID,
					EnvironmentId: testEnvironmentID,
					FilePath:      "src/updated.go",
					LineNumber:    20,
					CodeSnippet:   "client.BoolVariation(\"feature-1\", true)",
					ContentHash:   "hash-2",
				},
			)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				assert.Nil(t, persisted)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, persisted)
			assert.Equal(t, testCodeReferenceID, persisted.Id)
			assert.Equal(t, testEnvironmentID, persisted.EnvironmentId)
			assert.Equal(t, "src/updated.go", persisted.FilePath)
			assert.Equal(t, "src/updated.go", resp.CodeReference.FilePath)
		})
	}
}

func TestDeleteCodeReferenceScopedToEnvironment(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		setup   func(*CodeReferenceService)
		wantErr error
	}{
		{
			name: "returns not found and deletes nothing when the id does not belong to the requested environment",
			setup: func(s *CodeReferenceService) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(nil, storage.ErrCodeReferenceNotFound)
			},
			wantErr: statusNotFound.Err(),
		},
		{
			name: "deletes the code reference that belongs to the requested environment",
			setup: func(s *CodeReferenceService) {
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					GetCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(codeReferenceInEnvironment(testEnvironmentID), nil)
				s.codeRefStorage.(*storagemock.MockCodeReferenceStorage).EXPECT().
					DeleteCodeReference(gomock.Any(), testCodeReferenceID, testEnvironmentID).
					Return(nil)
				s.publisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			service := createCodeReferenceService(t, mockController)
			expectEditorInEnvironment(service, testEditorEmail, testEnvironmentID)
			tc.setup(service)

			resp, err := service.DeleteCodeReference(
				createContextWithToken(t, testEditorEmail),
				&proto.DeleteCodeReferenceRequest{
					Id:            testCodeReferenceID,
					EnvironmentId: testEnvironmentID,
				},
			)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Nil(t, resp)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
