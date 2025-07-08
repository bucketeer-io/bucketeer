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
	"errors"
	"slices"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	gstatus "google.golang.org/grpc/status"

	accclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/team/domain"
	"github.com/bucketeer-io/bucketeer/pkg/team/storage"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	accproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	proto "github.com/bucketeer-io/bucketeer/proto/team"
)

var (
	statusInternal         = gstatus.New(codes.Internal, "team: internal")
	statusNameRequired     = gstatus.New(codes.InvalidArgument, "team: name must be specified")
	statusInvalidCursor    = gstatus.New(codes.InvalidArgument, "team: cursor is invalid")
	statusInvalidOrderBy   = gstatus.New(codes.InvalidArgument, "team: order_by is invalid")
	statusUnauthenticated  = gstatus.New(codes.Unauthenticated, "team: unauthenticated")
	statusPermissionDenied = gstatus.New(codes.PermissionDenied, "team: permission denied")
	statusTeamNotFound     = gstatus.New(codes.NotFound, "team: not found")
	statusTeamInUsed       = gstatus.New(codes.FailedPrecondition, "team: team is in use by an account")
)

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type TeamService struct {
	mysqlClient   mysql.Client
	teamStorage   storage.TeamStorage
	accountClient accclient.Client
	publisher     publisher.Publisher
	flightgroup   singleflight.Group
	opts          *options
	logger        *zap.Logger
}

func NewTeamService(
	mysqlClient mysql.Client,
	accountClient accclient.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *TeamService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &TeamService{
		mysqlClient:   mysqlClient,
		teamStorage:   storage.NewTeamStorage(mysqlClient),
		accountClient: accountClient,
		publisher:     publisher,
		opts:          dopts,
		logger:        dopts.logger.Named("api"),
	}
}

func (s *TeamService) Register(server *grpc.Server) {
	proto.RegisterTeamServiceServer(server, s)
}

func (s *TeamService) CreateTeam(
	ctx context.Context,
	req *proto.CreateTeamRequest,
) (*proto.CreateTeamResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx, req.OrganizationId,
		accproto.AccountV2_Role_Organization_ADMIN, localizer)
	if err != nil {
		return nil, err
	}
	err = s.validateCreateTeamRequest(req, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to validate create team request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", req.OrganizationId),
				zap.String("name", req.Name),
			)...,
		)
		return nil, err
	}
	team, err := domain.NewTeam(
		req.Name,
		req.Description,
		req.OrganizationId,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create new team",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", req.OrganizationId),
				zap.String("name", req.Name),
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
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		if err := s.teamStorage.UpsertTeam(ctxWithTx, team); err != nil {
			return err
		}
		// Publish event
		event, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_TEAM,
			team.Id,
			eventproto.Event_TEAM_CREATED,
			&eventproto.TeamCreatedEvent{
				Id:          team.Id,
				Name:        team.Name,
				Description: team.Description,
				CreatedAt:   team.CreatedAt,
				UpdatedAt:   team.UpdatedAt,
			},
			team,
			team,
		)
		if err != nil {
			return err
		}
		return s.publisher.Publish(ctx, event)
	})
	if err != nil {
		s.logger.Error(
			"Failed to upsert team",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", req.OrganizationId),
				zap.String("name", req.Name),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.OrganizationId, localizer)
	}
	return &proto.CreateTeamResponse{
		Team: team.Team,
	}, nil
}

func (s *TeamService) validateCreateTeamRequest(req *proto.CreateTeamRequest, localizer locale.Localizer) error {
	if len(strings.TrimSpace(req.Name)) == 0 {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "name"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *TeamService) DeleteTeam(
	ctx context.Context,
	req *proto.DeleteTeamRequest,
) (*proto.DeleteTeamResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx, req.OrganizationId,
		accproto.AccountV2_Role_Organization_ADMIN, localizer)
	if err != nil {
		return nil, err
	}
	err = s.validateDeleteTeamRequest(req, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to validate delete team request",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", req.OrganizationId),
				zap.String("teamId", req.Id),
			)...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(ctxWithTx context.Context, tx mysql.Transaction) error {
		team, err := s.teamStorage.GetTeam(ctxWithTx, req.Id, req.OrganizationId)
		if err != nil {
			return err
		}

		// Check if team is in use by any account
		acs, err, _ := s.flightgroup.Do(
			req.OrganizationId,
			func() (interface{}, error) {
				return s.listAccountsFromOrganization(ctxWithTx, req.OrganizationId)
			},
		)
		if err != nil {
			return err
		}
		accounts := acs.([]*accproto.AccountV2)
		var inUsed = false
		for _, a := range accounts {
			if slices.Contains(a.Teams, team.Name) {
				inUsed = true
				break
			}
		}
		if inUsed {
			s.logger.Error(
				"Failed to delete the team because it is in use by an account",
				log.FieldsFromImcomingContext(ctxWithTx).AddFields(
					zap.String("organizationID", req.OrganizationId),
					zap.Any("team", team),
				)...,
			)
			return statusTeamInUsed.Err()
		}

		if err := s.teamStorage.DeleteTeam(ctxWithTx, req.Id); err != nil {
			return err
		}
		event, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_TEAM,
			req.Id,
			eventproto.Event_TEAM_DELETED,
			&eventproto.TeamDeletedEvent{
				Id:             team.Id,
				OrganizationId: req.OrganizationId,
			},
			nil,
			nil,
		)
		if err != nil {
			return err
		}
		return s.publisher.Publish(ctx, event)
	})
	if err != nil {
		if errors.Is(err, statusTeamInUsed.Err()) {
			dt, err := statusTeamInUsed.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.Team),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if errors.Is(err, storage.ErrTeamNotFound) {
			dt, err := statusTeamNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, s.reportInternalServerError(ctx, err, req.OrganizationId, localizer)
	}
	return &proto.DeleteTeamResponse{}, nil
}

func (s *TeamService) listAccountsFromOrganization(
	ctx context.Context,
	organizationID string,
) ([]*accproto.AccountV2, error) {
	resp, err := s.accountClient.ListAccountsV2(ctx, &accproto.ListAccountsV2Request{
		OrganizationId: organizationID,
	})
	if err != nil {
		s.logger.Error(
			"Failed to list accounts from organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationId", organizationID),
			)...,
		)
		return nil, err
	}
	if resp == nil || resp.Accounts == nil {
		s.logger.Warn(
			"No accounts found in organization",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.String("organizationId", organizationID),
			)...,
		)
		return nil, nil
	}
	return resp.Accounts, nil
}

func (s *TeamService) validateDeleteTeamRequest(
	req *proto.DeleteTeamRequest,
	localizer locale.Localizer,
) error {
	if len(strings.TrimSpace(req.Id)) == 0 {
		dt, err := statusNameRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "team_id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *TeamService) ListTeams(
	ctx context.Context,
	req *proto.ListTeamsRequest,
) (*proto.ListTeamsResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(
		ctx, req.OrganizationId,
		accproto.AccountV2_Role_Organization_MEMBER, localizer)
	if err != nil {
		return nil, err
	}

	filters := make([]*mysql.FilterV2, 0)
	filters = append(filters, &mysql.FilterV2{
		Column:   "team.organization_id",
		Operator: mysql.OperatorEqual,
		Value:    req.OrganizationId,
	})

	var searchQuery *mysql.SearchQuery
	if req.SearchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"team.name, team.description"},
			Keyword: req.SearchKeyword,
		}
	}

	orders, err := s.newListTeamsOrdersMySQL(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to valid list teams API. Invalid argument.",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}

	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	options := &mysql.ListOptions{
		Filters:     filters,
		SearchQuery: searchQuery,
		Orders:      orders,
		Limit:       limit,
		Offset:      offset,
		JSONFilters: nil,
		InFilters:   nil,
		NullFilters: nil,
		OrFilters:   nil,
	}
	teams, nextOffset, totalCount, err := s.teamStorage.ListTeams(ctx, options)
	if err != nil {
		s.logger.Error(
			"Failed to list teams",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, s.reportInternalServerError(ctx, err, req.OrganizationId, localizer)
	}
	return &proto.ListTeamsResponse{
		Teams:      teams,
		TotalCount: totalCount,
		NextCursor: strconv.Itoa(nextOffset),
	}, nil
}

func (s *TeamService) newListTeamsOrdersMySQL(
	orderBy proto.ListTeamsRequest_OrderBy,
	orderDirection proto.ListTeamsRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case proto.ListTeamsRequest_DEFAULT,
		proto.ListTeamsRequest_NAME:
		column = "team.name"
	case proto.ListTeamsRequest_CREATED_AT:
		column = "team.created_at"
	case proto.ListTeamsRequest_UPDATED_AT:
		column = "team.updated_at"
	case proto.ListTeamsRequest_ORGANIZATION:
		column = "team.organization_id"
	default:
		dt, err := statusInvalidOrderBy.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "order_by"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == proto.ListTeamsRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *TeamService) checkOrganizationRole(
	ctx context.Context,
	organizationID string,
	requiredRole accountproto.AccountV2_Role_Organization,
	localizer locale.Localizer,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(
		ctx,
		requiredRole,
		func(email string) (*accountproto.GetAccountV2Response, error) {
			return s.accountClient.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
				Email:          email,
				OrganizationId: organizationID,
			})
		},
	)
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

func (s *TeamService) reportInternalServerError(
	ctx context.Context,
	err error,
	organizationID string,
	localizer locale.Localizer,
) error {
	s.logger.Error(
		"Internal server error",
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
		return statusInternal.Err()
	}
	return dt.Err()
}
