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

package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	v2 "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	auditlogstorage "github.com/bucketeer-io/bucketeer/pkg/auditlog/storage/v2"
	environmentclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagstorage "github.com/bucketeer-io/bucketeer/pkg/tag/storage"
	teamstorage "github.com/bucketeer-io/bucketeer/pkg/team/storage"
	proto "github.com/bucketeer-io/bucketeer/proto/account"
	environmentproto "github.com/bucketeer-io/bucketeer/proto/environment"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

const (
	listRequestPageSize = 500
)

type options struct {
	logger *zap.Logger
}

var defaultOptions = options{
	logger: zap.NewNop(),
}

type Option func(*options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

type AccountService struct {
	environmentClient    environmentclient.Client
	mysqlClient          mysql.Client
	accountStorage       v2.AccountStorage
	tagStorage           tagstorage.TagStorage
	teamStorage          teamstorage.TeamStorage
	adminAuditLogStorage auditlogstorage.AdminAuditLogStorage
	publisher            publisher.Publisher
	opts                 *options
	logger               *zap.Logger
}

func NewAccountService(
	e environmentclient.Client,
	mysqlClient mysql.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *AccountService {
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &AccountService{
		environmentClient:    e,
		mysqlClient:          mysqlClient,
		accountStorage:       v2.NewAccountStorage(mysqlClient),
		tagStorage:           tagstorage.NewTagStorage(mysqlClient),
		teamStorage:          teamstorage.NewTeamStorage(mysqlClient),
		adminAuditLogStorage: auditlogstorage.NewAdminAuditLogStorage(mysqlClient),
		publisher:            publisher,
		opts:                 &options,
		logger:               options.logger.Named("api"),
	}
}

func (s *AccountService) Register(server *grpc.Server) {
	proto.RegisterAccountServiceServer(server, s)
}

func (s *AccountService) makeEnvironmentSet(
	environments []*environmentproto.EnvironmentV2,
) map[string]*environmentproto.EnvironmentV2 {
	environmentSet := make(map[string]*environmentproto.EnvironmentV2)
	for _, e := range environments {
		environmentSet[e.Id] = e
	}
	return environmentSet
}

func (s *AccountService) makeProjectSet(projects []*environmentproto.Project) map[string]*environmentproto.Project {
	projectSet := make(map[string]*environmentproto.Project)
	for _, p := range projects {
		projectSet[p.Id] = p
	}
	return projectSet
}

func (s *AccountService) listProjectsByOrganizationID(
	ctx context.Context,
	organizationID string,
) ([]*environmentproto.Project, error) {
	var projects []*environmentproto.Project
	cursor := ""
	for {
		resp, err := s.environmentClient.ListProjects(ctx, &environmentproto.ListProjectsRequest{
			PageSize:        listRequestPageSize,
			Cursor:          cursor,
			OrganizationIds: []string{organizationID},
		})
		if err != nil {
			return nil, err
		}
		projects = append(projects, resp.Projects...)
		projectSize := len(resp.Projects)
		if projectSize == 0 || projectSize < listRequestPageSize {
			return projects, nil
		}
		cursor = resp.Cursor
	}
}

func (s *AccountService) listEnvironmentsByOrganizationID(
	ctx context.Context,
	organizationID string,
) ([]*environmentproto.EnvironmentV2, error) {
	var environments []*environmentproto.EnvironmentV2
	cursor := ""
	for {
		resp, err := s.environmentClient.ListEnvironmentsV2(ctx, &environmentproto.ListEnvironmentsV2Request{
			PageSize:       listRequestPageSize,
			Cursor:         cursor,
			OrganizationId: organizationID,
			Archived:       wrapperspb.Bool(false),
		})
		if err != nil {
			return nil, err
		}
		environments = append(environments, resp.Environments...)
		environmentSize := len(resp.Environments)
		if environmentSize == 0 || environmentSize < listRequestPageSize {
			return environments, nil
		}
		cursor = resp.Cursor
	}
}

func (s *AccountService) getOrganization(
	ctx context.Context,
	id string,
) (*environmentproto.Organization, error) {
	resp, err := s.environmentClient.GetOrganization(ctx, &environmentproto.GetOrganizationRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Organization, nil
}

func (s *AccountService) checkSystemAdminRole(
	ctx context.Context,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckSystemAdminRole(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}

func (s *AccountService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole proto.AccountV2_Role_Environment,
	environmentId string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*proto.AccountV2, error) {
			account, err := s.getAccountV2ByEnvironmentID(ctx, email, environmentId, localizer)
			if err != nil {
				return nil, err
			}
			return account.AccountV2, nil
		},
	)
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}

func (s *AccountService) checkOrganizationRole(
	ctx context.Context,
	requiredRole proto.AccountV2_Role_Organization,
	organizationID string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(ctx, requiredRole, func(email string) (*proto.GetAccountV2Response, error) {
		account, err := s.getAccountV2(ctx, email, organizationID, localizer)
		if err != nil {
			return nil, err
		}
		return &proto.GetAccountV2Response{Account: account.AccountV2}, nil
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}

func (s *AccountService) checkOrganizationRoleByEnvironmentID(
	ctx context.Context,
	requiredRole proto.AccountV2_Role_Organization,
	environmentID string,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(ctx, requiredRole, func(email string) (*proto.GetAccountV2Response, error) {
		account, err := s.getAccountV2ByEnvironmentID(ctx, email, environmentID, localizer)
		if err != nil {
			return nil, err
		}
		return &proto.GetAccountV2Response{Account: account.AccountV2}, nil
	})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			dt, err := statusUnauthenticated.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.UnauthenticatedError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentID", environmentID),
				)...,
			)
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
	}
	return editor, nil
}
