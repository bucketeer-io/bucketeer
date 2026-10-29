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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	accountclientmock "github.com/bucketeer-io/bucketeer/v2/pkg/account/client/mock"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2/mock"
	publishermock "github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher/mock"
	databasemock "github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/database/mock"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func TestCloneFeatureRequiresSourceEnvironmentRole(t *testing.T) {
	t.Parallel()

	const (
		// Email of the token created by createContextWithTokenRoleUnassigned.
		callerEmail = "email"
		featureID   = "feature-0"
		sourceEnvID = "env-source"
		targetEnvID = "env-target"
		secondEnvID = "env-target-2"
		// An environment of an organization the caller has no account in.
		foreignEnvID = "env-foreign"

		roleViewer = accountproto.AccountV2_Role_Environment_VIEWER
		roleEditor = accountproto.AccountV2_Role_Environment_EDITOR
	)

	sourceFeature := func() *domain.Feature {
		return &domain.Feature{Feature: &featureproto.Feature{
			Id:         featureID,
			Name:       "source",
			Variations: createFeatureVariations(),
			Targets: []*featureproto.Target{
				{Variation: "variation_id_1", Users: []string{"user-1"}},
			},
			DefaultStrategy: &featureproto.Strategy{
				Type:          featureproto.Strategy_FIXED,
				FixedStrategy: &featureproto.FixedStrategy{Variation: "variation_id_1"},
			},
		}}
	}

	cases := []struct {
		name string
		// Environment roles of the caller's account in the organization that
		// owns the environments. Environments missing from the map are
		// UNASSIGNED for the caller.
		roles map[string]accountproto.AccountV2_Role_Environment
		call  func(*testing.T, context.Context, *FeatureService) error
		// Environments where the cloned feature is expected to be stored.
		wantCreated []string
		wantErr     error
	}{
		{
			name:  "clone rejects a caller without a role on the source environment",
			roles: map[string]accountproto.AccountV2_Role_Environment{targetEnvID: roleEditor},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:                  featureID,
					EnvironmentId:       sourceEnvID,
					TargetEnvironmentId: targetEnvID,
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name:  "clone rejects a source environment of another organization",
			roles: map[string]accountproto.AccountV2_Role_Environment{targetEnvID: roleEditor},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:                  featureID,
					EnvironmentId:       foreignEnvID,
					TargetEnvironmentId: targetEnvID,
				})
				return err
			},
			wantErr: statusUnauthenticated.Err(),
		},
		{
			name: "clone still rejects a caller without editor on the target environment",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				sourceEnvID: roleEditor,
				targetEnvID: roleViewer,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:                  featureID,
					EnvironmentId:       sourceEnvID,
					TargetEnvironmentId: targetEnvID,
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "clone succeeds with viewer on the source and editor on the target",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				sourceEnvID: roleViewer,
				targetEnvID: roleEditor,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:                  featureID,
					EnvironmentId:       sourceEnvID,
					TargetEnvironmentId: targetEnvID,
				})
				return err
			},
			wantCreated: []string{targetEnvID},
		},
		{
			name: "clone succeeds with editor on both environments",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				sourceEnvID: roleEditor,
				targetEnvID: roleEditor,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.CloneFeature(ctx, &featureproto.CloneFeatureRequest{
					Id:                  featureID,
					EnvironmentId:       sourceEnvID,
					TargetEnvironmentId: targetEnvID,
				})
				return err
			},
			wantCreated: []string{targetEnvID},
		},
		{
			name: "bulk clone rejects a caller without a role on the source environment",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				targetEnvID: roleEditor,
				secondEnvID: roleEditor,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.BulkCloneFeature(ctx, &featureproto.BulkCloneFeatureRequest{
					Id:                   featureID,
					EnvironmentId:        sourceEnvID,
					TargetEnvironmentIds: []string{targetEnvID, secondEnvID},
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "bulk clone rejects a source environment of another organization",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				targetEnvID: roleEditor,
				secondEnvID: roleEditor,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.BulkCloneFeature(ctx, &featureproto.BulkCloneFeatureRequest{
					Id:                   featureID,
					EnvironmentId:        foreignEnvID,
					TargetEnvironmentIds: []string{targetEnvID, secondEnvID},
				})
				return err
			},
			wantErr: statusUnauthenticated.Err(),
		},
		{
			name: "bulk clone still rejects a caller without editor on one target environment",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				sourceEnvID: roleViewer,
				targetEnvID: roleEditor,
				secondEnvID: roleViewer,
			},
			call: func(_ *testing.T, ctx context.Context, s *FeatureService) error {
				_, err := s.BulkCloneFeature(ctx, &featureproto.BulkCloneFeatureRequest{
					Id:                   featureID,
					EnvironmentId:        sourceEnvID,
					TargetEnvironmentIds: []string{targetEnvID, secondEnvID},
				})
				return err
			},
			wantErr: statusPermissionDenied.Err(),
		},
		{
			name: "bulk clone succeeds with viewer on the source and editor on every target",
			roles: map[string]accountproto.AccountV2_Role_Environment{
				sourceEnvID: roleViewer,
				targetEnvID: roleEditor,
				secondEnvID: roleEditor,
			},
			call: func(t *testing.T, ctx context.Context, s *FeatureService) error {
				resp, err := s.BulkCloneFeature(ctx, &featureproto.BulkCloneFeatureRequest{
					Id:                   featureID,
					EnvironmentId:        sourceEnvID,
					TargetEnvironmentIds: []string{targetEnvID, secondEnvID},
				})
				if err != nil {
					return err
				}
				assert.Equal(t, []*featureproto.BulkCloneFeatureResult{
					{
						EnvironmentId: targetEnvID,
						Status:        featureproto.BulkCloneFeatureStatus_BULK_CLONE_FEATURE_STATUS_SUCCESS,
					},
					{
						EnvironmentId: secondEnvID,
						Status:        featureproto.BulkCloneFeatureStatus_BULK_CLONE_FEATURE_STATUS_SUCCESS,
					},
				}, resp.Results)
				return nil
			},
			wantCreated: []string{targetEnvID, secondEnvID},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			// The token is not a system admin; the caller's roles come from the
			// account returned by the account client below.
			ctx := createContextWithTokenRoleUnassigned()
			ctx = metadata.NewIncomingContext(ctx, metadata.MD{})

			service := createFeatureServiceNew(mockController)

			envRoles := make([]*accountproto.AccountV2_EnvironmentRole, 0, len(tc.roles))
			for envID, r := range tc.roles {
				envRoles = append(envRoles, &accountproto.AccountV2_EnvironmentRole{
					EnvironmentId: envID,
					Role:          r,
				})
			}
			accountClient := accountclientmock.NewMockClient(mockController)
			accountClient.EXPECT().
				GetAccountV2ByEnvironmentID(gomock.Any(), gomock.Any()).
				DoAndReturn(func(
					_ context.Context,
					req *accountproto.GetAccountV2ByEnvironmentIDRequest,
					_ ...grpc.CallOption,
				) (*accountproto.GetAccountV2ByEnvironmentIDResponse, error) {
					// The account is resolved in the organization that owns the
					// requested environment, so there is none for foreignEnvID.
					if req.Email != callerEmail || req.EnvironmentId == foreignEnvID {
						return nil, status.Error(codes.NotFound, "account not found")
					}
					return &accountproto.GetAccountV2ByEnvironmentIDResponse{
						Account: &accountproto.AccountV2{
							Email:            callerEmail,
							OrganizationId:   "org0",
							OrganizationRole: accountproto.AccountV2_Role_Organization_MEMBER,
							EnvironmentRoles: envRoles,
						},
					}, nil
				}).AnyTimes()
			service.accountClient = accountClient

			featureStorage := service.featureStorage.(*mock.MockFeatureStorage)
			created := map[string]*domain.Feature{}
			if tc.wantErr != nil {
				// A denied request must neither read the source feature nor
				// write to any target environment.
				featureStorage.EXPECT().GetFeature(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				featureStorage.EXPECT().CreateFeature(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			} else {
				featureStorage.EXPECT().
					GetFeature(gomock.Any(), featureID, sourceEnvID).
					Return(sourceFeature(), nil)
				service.dbClient.(*databasemock.MockClient).EXPECT().
					RunInTransactionV2(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
						return fn(ctx)
					}).Times(len(tc.wantCreated))
				featureStorage.EXPECT().
					CreateFeature(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, f *domain.Feature, envID string) error {
						created[envID] = f
						return nil
					}).Times(len(tc.wantCreated))
				service.domainPublisher.(*publishermock.MockPublisher).EXPECT().
					Publish(gomock.Any(), gomock.Any()).
					Return(nil).Times(len(tc.wantCreated))
			}

			err := tc.call(t, ctx, service)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				assert.Empty(t, created)
				return
			}
			require.NoError(t, err)
			require.Len(t, created, len(tc.wantCreated))
			for _, envID := range tc.wantCreated {
				cloned, ok := created[envID]
				require.True(t, ok, "feature was not stored in %s", envID)
				assert.Equal(t, featureID, cloned.Id)
				assert.Equal(t, callerEmail, cloned.Maintainer)
				assert.Len(t, cloned.Variations, len(createFeatureVariations()))
				require.Len(t, cloned.Targets, 1)
				assert.Equal(t, []string{"user-1"}, cloned.Targets[0].Users)
			}
		})
	}
}
