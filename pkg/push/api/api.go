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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	accountclient "github.com/bucketeer-io/bucketeer/v2/pkg/account/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	experimentclient "github.com/bucketeer-io/bucketeer/v2/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/v2/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/v2/pkg/push/domain"
	v2ps "github.com/bucketeer-io/bucketeer/v2/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/role"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	pushproto "github.com/bucketeer-io/bucketeer/v2/proto/push"
)

var errTagDuplicated = errors.New("push: tag is duplicated")

type options struct {
	logger *zap.Logger
}

type Option func(*options)

func WithLogger(l *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = l
	}
}

type PushService struct {
	mysqlClient      mysql.Client
	pushStorage      v2ps.PushStorage
	featureClient    featureclient.Client
	experimentClient experimentclient.Client
	accountClient    accountclient.Client
	publisher        publisher.Publisher
	opts             *options
	logger           *zap.Logger
}

func NewPushService(
	mysqlClient mysql.Client,
	featureClient featureclient.Client,
	experimentClient experimentclient.Client,
	accountClient accountclient.Client,
	publisher publisher.Publisher,
	opts ...Option,
) *PushService {
	dopts := &options{
		logger: zap.NewNop(),
	}
	for _, opt := range opts {
		opt(dopts)
	}
	return &PushService{
		mysqlClient:      mysqlClient,
		pushStorage:      v2ps.NewPushStorage(mysqlClient),
		featureClient:    featureClient,
		experimentClient: experimentClient,
		accountClient:    accountClient,
		publisher:        publisher,
		opts:             dopts,
		logger:           dopts.logger.Named("api"),
	}
}

func (s *PushService) Register(server *grpc.Server) {
	pushproto.RegisterPushServiceServer(server, s)
}

func (s *PushService) CreatePush(
	ctx context.Context,
	req *pushproto.CreatePushRequest,
) (*pushproto.CreatePushResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := s.validateCreatePushRequest(req); err != nil {
		return nil, err
	}
	push, err := domain.NewPush(
		req.Name,
		string(req.FcmServiceAccount),
		req.Tags,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Strings("tags", req.Tags),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	pushes, err := s.listAllPushes(ctx, req.EnvironmentId)
	if err != nil {
		return nil, api.NewGRPCStatus(err).Err()
	}
	if err := s.checkFCMServiceAccount(ctx, pushes, req.FcmServiceAccount); err != nil {
		return nil, err
	}
	err = s.containsTags(pushes, req.Tags)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, statusTagAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to validate tag existence",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Strings("tags", req.Tags),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := s.pushStorage.CreatePush(contextWithTx, push, req.EnvironmentId); err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_PUSH,
			push.Id,
			eventproto.Event_PUSH_CREATED,
			&eventproto.PushCreatedEvent{
				FcmServiceAccount: push.FcmServiceAccount,
				Tags:              push.Tags,
				Name:              push.Name,
			},
			req.EnvironmentId,
			push,
			nil,
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, event); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, v2ps.ErrPushAlreadyExists) {
			return nil, statusAlreadyExists.Err()
		}
		s.logger.Error(
			"Failed to create push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	// For security reasons we remove the service account from the API response
	push.FcmServiceAccount = ""

	return &pushproto.CreatePushResponse{
		Push: push.Push,
	}, nil
}

func (s *PushService) validateCreatePushRequest(req *pushproto.CreatePushRequest) error {
	if string(req.FcmServiceAccount) == "" {
		return statusFCMServiceAccountRequired.Err()
	}
	if req.Name == "" {
		return statusNameRequired.Err()
	}
	return nil
}

func (s *PushService) UpdatePush(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
) (*pushproto.UpdatePushResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if err := s.validateUpdatePushRequest(req); err != nil {
		return nil, err
	}
	var updatedPushPb *pushproto.Push
	var updatePushEvent *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		push, err := s.pushStorage.GetPush(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		updated, err := push.Update(req.Name, req.TagChanges, req.Disabled)
		if err != nil {
			return err
		}
		updatePushEvent, err = domainevent.NewEvent(
			editor,
			eventproto.Event_PUSH,
			push.Id,
			eventproto.Event_PUSH_UPDATED,
			&eventproto.PushUpdatedEvent{
				Name: req.Name,
				Tags: req.Tags,
			},
			req.EnvironmentId,
			updated,
			push,
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, updatePushEvent); err != nil {
			return err
		}
		updatedPushPb = updated.Push

		return s.pushStorage.UpdatePush(contextWithTx, updated, req.EnvironmentId)
	})
	if err != nil {
		switch {
		case errors.Is(err, v2ps.ErrPushNotFound):
			return nil, statusNotFound.Err()
		case errors.Is(err, v2ps.ErrPushUnexpectedAffectedRows):
			if updatedPushPb != nil {
				// For security reasons we remove the service account from the API response
				updatedPushPb.FcmServiceAccount = ""
			}
			return &pushproto.UpdatePushResponse{
				Push: updatedPushPb,
			}, nil
		}
		s.logger.Error(
			"Failed to update push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}

	if updatedPushPb != nil {
		// For security reasons we remove the service account from the API response
		updatedPushPb.FcmServiceAccount = ""
	}

	return &pushproto.UpdatePushResponse{
		Push: updatedPushPb,
	}, nil
}

func (s *PushService) validateUpdatePushRequest(
	req *pushproto.UpdatePushRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}

	return nil
}

func (s *PushService) DeletePush(
	ctx context.Context,
	req *pushproto.DeletePushRequest,
) (*pushproto.DeletePushResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateDeletePushRequest(req); err != nil {
		return nil, err
	}

	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		push, err := s.pushStorage.GetPush(contextWithTx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_PUSH,
			push.Id,
			eventproto.Event_PUSH_DELETED,
			&eventproto.PushCreatedEvent{
				FcmServiceAccount: push.FcmServiceAccount,
				Tags:              push.Tags,
				Name:              push.Name,
			},
			req.EnvironmentId,
			nil,  // Current state: entity no longer exists
			push, // Previous state: what was deleted
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return s.pushStorage.DeletePush(contextWithTx, push.Id, req.EnvironmentId)
	})
	if err != nil {
		switch {
		case errors.Is(err, v2ps.ErrPushNotFound):
			return nil, statusNotFound.Err()
		case errors.Is(err, v2ps.ErrPushUnexpectedAffectedRows):
			return &pushproto.DeletePushResponse{}, nil
		}
		s.logger.Error(
			"Failed to delete push",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &pushproto.DeletePushResponse{}, nil
}

func (s *PushService) GetPush(
	ctx context.Context,
	req *pushproto.GetPushRequest,
) (*pushproto.GetPushResponse, error) {
	if err := s.validateGetPushRequest(req); err != nil {
		return nil, err
	}

	push, err := s.pushStorage.GetPush(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2ps.ErrPushNotFound) {
			s.logger.Error(
				"Failed to get push",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, statusNotFound.Err()
		}
	}

	if push.Push != nil {
		// For security reasons we remove the service account from the API response
		push.FcmServiceAccount = ""
	}

	return &pushproto.GetPushResponse{
		Push: push.Push,
	}, nil
}

func (s *PushService) validateGetPushRequest(
	req *pushproto.GetPushRequest,
) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}

func validateDeletePushRequest(req *pushproto.DeletePushRequest) error {
	if req.Id == "" {
		return statusIDRequired.Err()
	}
	return nil
}

func (s *PushService) containsTags(
	pushes []*pushproto.Push,
	tags []string,
) error {
	m, err := s.tagMap(pushes)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if _, ok := m[t]; ok {
			return statusTagAlreadyExists.Err()
		}
	}
	return nil
}

func (s *PushService) checkFCMServiceAccount(
	ctx context.Context,
	pushes []*pushproto.Push,
	fcmServiceAccount []byte,
) error {
	// Check if the JSON is a service account file
	_, err := google.CredentialsFromJSON(
		ctx,
		fcmServiceAccount,
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		s.logger.Error("failed to get credentials from JSON", zap.Error(err))
		return statusFCMServiceAccountInvalid.Err()
	}
	// Check if the service account already exists in the database
	for _, push := range pushes {
		equal, err := s.compareJSON(push.FcmServiceAccount, string(fcmServiceAccount))
		if err != nil {
			s.logger.Error("failed to compare the JSON", zap.Error(err))
			return statusInternal.Err()
		}
		if equal {
			s.logger.Error("fcm service account already exists in the database")
			return statusFCMServiceAccountAlreadyExists.Err()
		}
	}
	return nil
}

// compareJSON compares two JSON strings and returns true if they are equivalent
func (s *PushService) compareJSON(jsonStr1, jsonStr2 string) (bool, error) {
	var obj1, obj2 json.RawMessage
	// Unmarshal the JSON strings into Go data structures
	if err := json.Unmarshal([]byte(jsonStr1), &obj1); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(jsonStr2), &obj2); err != nil {
		return false, err
	}
	// Marshal the Go data structures into canonical JSON format
	json1, err := json.Marshal(obj1)
	if err != nil {
		return false, err
	}
	json2, err := json.Marshal(obj2)
	if err != nil {
		return false, err
	}
	// Compare the canonical JSON representations
	return bytes.Equal(json1, json2), nil
}

func (s *PushService) tagMap(pushes []*pushproto.Push) (map[string]struct{}, error) {
	m := make(map[string]struct{})
	for _, p := range pushes {
		for _, t := range p.Tags {
			if _, ok := m[t]; ok {
				return nil, errTagDuplicated
			}
			m[t] = struct{}{}
		}
	}
	return m, nil
}

func (s *PushService) listAllPushes(
	ctx context.Context,
	environmentId string,
) ([]*pushproto.Push, error) {
	pushes, _, _, err := s.listPushes(
		ctx,
		mysql.QueryNoLimit,
		"",
		"",
		[]string{environmentId},
		"",
		wrapperspb.Bool(false),
		nil,
	)
	if err != nil {
		return nil, err
	}
	return pushes, nil
}

func (s *PushService) ListPushes(
	ctx context.Context,
	req *pushproto.ListPushesRequest,
) (*pushproto.ListPushesResponse, error) {
	var filterEnvironmentIDs []string
	if req.OrganizationId != "" {
		// console v3
		editor, err := s.checkOrganizationRole(
			ctx, accountproto.AccountV2_Role_Organization_MEMBER,
			req.OrganizationId,
		)
		if err != nil {
			return nil, err
		}
		filterEnvironmentIDs = s.getAllowedEnvironments(req.EnvironmentIds, editor)
	} else {
		// console v2
		_, err := s.checkEnvironmentRole(
			ctx, accountproto.AccountV2_Role_Environment_VIEWER,
			req.EnvironmentId)
		if err != nil {
			return nil, err
		}
		filterEnvironmentIDs = append(filterEnvironmentIDs, req.EnvironmentId)
	}

	orders, err := s.newListOrders(req.OrderBy, req.OrderDirection)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	pushes, cursor, totalCount, err := s.listPushes(
		ctx,
		req.PageSize,
		req.Cursor,
		req.OrganizationId,
		filterEnvironmentIDs,
		req.SearchKeyword,
		req.Disabled,
		orders,
	)
	if err != nil {
		return nil, err
	}
	// For security reasons we remove the service account from the API response
	for _, p := range pushes {
		p.FcmServiceAccount = ""
	}
	return &pushproto.ListPushesResponse{
		Pushes:     pushes,
		Cursor:     cursor,
		TotalCount: totalCount,
	}, nil
}

func (s *PushService) getAllowedEnvironments(
	reqEnvironmentIDs []string,
	editor *eventproto.Editor,
) []string {
	filterEnvironmentIDs := make([]string, 0)
	if editor.OrganizationRole == accountproto.AccountV2_Role_Organization_MEMBER {
		// only show API keys in allowed environments for member.
		if len(reqEnvironmentIDs) > 0 {
			for _, id := range reqEnvironmentIDs {
				for _, e := range editor.EnvironmentRoles {
					if e.EnvironmentId == id {
						filterEnvironmentIDs = append(filterEnvironmentIDs, id)
						break
					}
				}
			}
		} else {
			for _, e := range editor.EnvironmentRoles {
				filterEnvironmentIDs = append(filterEnvironmentIDs, e.EnvironmentId)
			}
		}
	} else {
		// if the user is an admin or owner, no need to filter environments.
		filterEnvironmentIDs = append(filterEnvironmentIDs, reqEnvironmentIDs...)
	}
	return filterEnvironmentIDs
}

func (s *PushService) newListOrders(
	orderBy pushproto.ListPushesRequest_OrderBy,
	orderDirection pushproto.ListPushesRequest_OrderDirection,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case pushproto.ListPushesRequest_DEFAULT,
		pushproto.ListPushesRequest_NAME:
		column = "push.name"
	case pushproto.ListPushesRequest_CREATED_AT:
		column = "push.created_at"
	case pushproto.ListPushesRequest_UPDATED_AT:
		column = "push.updated_at"
	case pushproto.ListPushesRequest_ENVIRONMENT:
		column = "env.name"
	case pushproto.ListPushesRequest_STATE:
		column = "push.disabled"
	default:
		return nil, statusInvalidOrderBy.Err()
	}
	direction := mysql.OrderDirectionAsc
	if orderDirection == pushproto.ListPushesRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *PushService) listPushes(
	ctx context.Context,
	pageSize int64,
	cursor string,
	organizationId string,
	environmentIDs []string,
	searchKeyword string,
	disabled *wrapperspb.BoolValue,
	orders []*mysql.Order,
) ([]*pushproto.Push, string, int64, error) {
	var filters []*mysql.FilterV2
	var inFilters []*mysql.InFilter
	if organizationId != "" {
		// console v3
		filters = append(filters, &mysql.FilterV2{
			Column:   "env.organization_id",
			Operator: mysql.OperatorEqual,
			Value:    organizationId,
		})
		if len(environmentIDs) > 0 {
			envIDs := make([]interface{}, 0, len(environmentIDs))
			for _, id := range environmentIDs {
				envIDs = append(envIDs, id)
			}
			inFilters = append(inFilters, &mysql.InFilter{
				Column: "push.environment_id",
				Values: envIDs,
			})
		}
	} else {
		// console v2
		if len(environmentIDs) > 0 {
			filters = append(filters, &mysql.FilterV2{
				Column:   "push.environment_id",
				Operator: mysql.OperatorEqual,
				Value:    environmentIDs[0],
			})
		}
	}
	if disabled != nil {
		filters = append(filters, &mysql.FilterV2{
			Column:   "push.disabled",
			Operator: mysql.OperatorEqual,
			Value:    disabled.Value,
		})
	}
	var searchQuery *mysql.SearchQuery
	if searchKeyword != "" {
		searchQuery = &mysql.SearchQuery{
			Columns: []string{"push.name"},
			Keyword: searchKeyword,
		}
	}
	limit := int(pageSize)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, "", 0, statusInvalidCursor.Err()
	}

	options := &mysql.ListOptions{
		Limit:       limit,
		Offset:      offset,
		Filters:     filters,
		SearchQuery: searchQuery,
		InFilters:   inFilters,
		Orders:      orders,
		JSONFilters: nil,
		NullFilters: nil,
	}
	pushes, nextCursor, totalCount, err := s.pushStorage.ListPushes(
		ctx,
		options,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list pushes",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.Strings("environmentId", environmentIDs),
			)...,
		)
		return nil, "", 0, statusInternal.Err()
	}
	return pushes, strconv.Itoa(nextCursor), totalCount, nil
}

func (s *PushService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckEnvironmentRole(
		ctx,
		requiredRole,
		environmentId,
		func(email string) (*accountproto.AccountV2, error) {
			resp, err := s.accountClient.GetAccountV2ByEnvironmentID(ctx, &accountproto.GetAccountV2ByEnvironmentIDRequest{
				Email:         email,
				EnvironmentId: environmentId,
			})
			if err != nil {
				return nil, err
			}
			return resp.Account, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}

func (s *PushService) checkOrganizationRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Organization,
	organizationID string,
) (*eventproto.Editor, error) {
	editor, err := role.CheckOrganizationRole(
		ctx,
		requiredRole,
		func(email string) (*accountproto.GetAccountV2Response, error) {
			resp, err := s.accountClient.GetAccountV2(ctx, &accountproto.GetAccountV2Request{
				Email:          email,
				OrganizationId: organizationID,
			})
			if err != nil {
				return nil, err
			}
			return resp, nil
		})
	if err != nil {
		switch status.Code(err) {
		case codes.Unauthenticated:
			s.logger.Error(
				"Unauthenticated",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			return nil, statusUnauthenticated.Err()
		case codes.PermissionDenied:
			s.logger.Error(
				"Permission denied",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			return nil, statusPermissionDenied.Err()
		default:
			s.logger.Error(
				"Failed to check role",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", organizationID),
				)...,
			)
			return nil, api.NewGRPCStatus(err).Err()
		}
	}
	return editor, nil
}
