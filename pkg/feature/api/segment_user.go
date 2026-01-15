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
	"errors"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/v2/pkg/api/api"
	domainevent "github.com/bucketeer-io/bucketeer/v2/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/v2/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/v2/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/v2/pkg/log"
	"github.com/bucketeer-io/bucketeer/v2/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/v2/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/v2/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/domain"
	serviceeventproto "github.com/bucketeer-io/bucketeer/v2/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/v2/proto/feature"
)

func (s *FeatureService) AddSegmentUser(
	ctx context.Context,
	req *featureproto.AddSegmentUserRequest,
) (*featureproto.AddSegmentUserResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateAddSegmentUserRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := validateAddSegmentUserCommand(req.Command); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.updateSegmentUser(
		ctx,
		editor,
		req.Id,
		req.Command.UserIds,
		req.Command.State,
		false,
		req.Command,
		req.EnvironmentId,
	); err != nil {
		return nil, err
	}
	return &featureproto.AddSegmentUserResponse{}, nil
}

func (s *FeatureService) DeleteSegmentUser(
	ctx context.Context,
	req *featureproto.DeleteSegmentUserRequest,
) (*featureproto.DeleteSegmentUserResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteSegmentUserRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := validateDeleteSegmentUserCommand(req.Command); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := s.updateSegmentUser(
		ctx,
		editor,
		req.Id,
		req.Command.UserIds,
		req.Command.State,
		true,
		req.Command,
		req.EnvironmentId,
	); err != nil {
		return nil, err
	}
	return &featureproto.DeleteSegmentUserResponse{}, nil
}

func (s *FeatureService) updateSegmentUser(
	ctx context.Context,
	editor *eventproto.Editor,
	segmentID string,
	userIDs []string,
	state featureproto.SegmentUser_State,
	deleted bool,
	cmd command.Command,
	environmentId string,
) error {
	segmentUsers := make([]*featureproto.SegmentUser, 0, len(userIDs))
	for _, userID := range userIDs {
		userID = strings.TrimSpace(userID)
		user := domain.NewSegmentUser(segmentID, userID, state, deleted)
		segmentUsers = append(segmentUsers, user.SegmentUser)
	}
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, segmentID, environmentId)
		if err != nil {
			s.logger.Error(
				"Failed to get segment",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		if err := s.segmentUserStorage.UpsertSegmentUsers(contextWithTx, segmentUsers, environmentId); err != nil {
			s.logger.Error(
				"Failed to store segment user",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			environmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, cmd); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
			return err
		}
		if err := s.segmentStorage.UpdateSegment(contextWithTx, segment, environmentId); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2fs.ErrSegmentNotFound || err == v2fs.ErrSegmentUnexpectedAffectedRows {
			return statusSegmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to upsert segment user",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", environmentId),
			)...,
		)
		return api.NewGRPCStatus(err).Err()
	}
	return nil
}

func (s *FeatureService) GetSegmentUser(
	ctx context.Context,
	req *featureproto.GetSegmentUserRequest,
) (*featureproto.GetSegmentUserResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateGetSegmentUserRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	id := domain.SegmentUserID(req.SegmentId, req.UserId, req.State)
	user, err := s.segmentUserStorage.GetSegmentUser(ctx, id, req.EnvironmentId)
	if err != nil {
		if err == v2fs.ErrSegmentUserNotFound {
			return nil, statusSegmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to get segment user",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.GetSegmentUserResponse{
		User: user.SegmentUser,
	}, nil
}

func (s *FeatureService) ListSegmentUsers(
	ctx context.Context,
	req *featureproto.ListSegmentUsersRequest,
) (*featureproto.ListSegmentUsersResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateListSegmentUsersRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("segment_id", "=", req.SegmentId),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.State != nil {
		whereParts = append(whereParts, mysql.NewFilter("state", "=", req.State.GetValue()))
	}
	if req.UserId != "" {
		whereParts = append(whereParts, mysql.NewFilter("user_id", "=", req.UserId))
	}
	limit := int(req.PageSize)
	cursor := req.Cursor
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.Atoi(cursor)
	if err != nil {
		return nil, statusInvalidCursor.Err()
	}
	users, nextCursor, err := s.segmentUserStorage.ListSegmentUsers(
		ctx,
		whereParts,
		nil,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list segment users",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.ListSegmentUsersResponse{
		Users:  users,
		Cursor: strconv.Itoa(nextCursor),
	}, nil
}

func (s *FeatureService) BulkUploadSegmentUsers(
	ctx context.Context,
	req *featureproto.BulkUploadSegmentUsersRequest,
) (*featureproto.BulkUploadSegmentUsersResponse, error) {
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.bulkUploadSegmentUsersNoCommand(ctx, req, editor)
	}
	if err := validateBulkUploadSegmentUsersRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	if err := validateBulkUploadSegmentUsersCommand(req.Command); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, req.SegmentId, req.EnvironmentId)
		if err != nil {
			return err
		}
		if segment.Status == featureproto.Segment_UPLOADING {
			return statusSegmentUsersAlreadyUploading.Err()
		}
		handler, err := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			req.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return err
		}
		if err := s.segmentStorage.UpdateSegment(contextWithTx, segment, req.EnvironmentId); err != nil {
			return err
		}
		return s.publishBulkSegmentUsersReceivedEvent(
			ctx,
			editor,
			req.EnvironmentId,
			req.SegmentId,
			req.Command.Data,
			req.Command.State,
		)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) || errors.Is(err, v2fs.ErrFeatureUnexpectedAffectedRows) {
			return nil, statusSegmentNotFound.Err()
		}
		if status.Code(err) == codes.FailedPrecondition {
			return nil, err
		}
		s.logger.Error(
			"Failed to bulk upload segment users",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.BulkUploadSegmentUsersResponse{}, nil
}

func (s *FeatureService) bulkUploadSegmentUsersNoCommand(
	ctx context.Context,
	req *featureproto.BulkUploadSegmentUsersRequest,
	editor *eventproto.Editor,
) (*featureproto.BulkUploadSegmentUsersResponse, error) {
	if err := validateBulkUploadSegmentUsersNoCommandRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		segment, _, err := s.segmentStorage.GetSegment(contextWithTx, req.SegmentId, req.EnvironmentId)
		if err != nil {
			return err
		}
		if segment.Status == featureproto.Segment_UPLOADING {
			return statusSegmentUsersAlreadyUploading.Err()
		}
		prev := &domain.Segment{}
		if err := copier.Copy(prev, segment); err != nil {
			return err
		}
		segment.SetStatus(featureproto.Segment_UPLOADING)
		if err := s.segmentStorage.UpdateSegment(contextWithTx, segment, req.EnvironmentId); err != nil {
			return err
		}
		e, err := domainevent.NewEvent(
			editor,
			eventproto.Event_SEGMENT,
			segment.Id,
			eventproto.Event_SEGMENT_BULK_UPLOAD_USERS,
			&eventproto.SegmentBulkUploadUsersEvent{
				SegmentId: segment.Id,
				Status:    featureproto.Segment_UPLOADING,
				State:     req.State,
			},
			req.EnvironmentId,
			segment.Segment,
			prev,
		)
		if err != nil {
			return err
		}
		err = s.domainPublisher.Publish(ctx, e)
		if err != nil {
			return err
		}
		return s.publishBulkSegmentUsersReceivedEvent(
			ctx,
			editor,
			req.EnvironmentId,
			req.SegmentId,
			req.Data,
			req.State,
		)
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) || errors.Is(err, v2fs.ErrFeatureUnexpectedAffectedRows) {
			return nil, statusSegmentNotFound.Err()
		}
		if status.Code(err) == codes.FailedPrecondition {
			return nil, err
		}
		s.logger.Error(
			"Failed to bulk upload segment users",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	return &featureproto.BulkUploadSegmentUsersResponse{}, nil
}

func (s *FeatureService) publishBulkSegmentUsersReceivedEvent(
	ctx context.Context,
	editor *eventproto.Editor,
	environmentId string,
	segmentID string,
	data []byte,
	state featureproto.SegmentUser_State,
) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	e := &serviceeventproto.BulkSegmentUsersReceivedEvent{
		Id:            id.String(),
		EnvironmentId: environmentId,
		SegmentId:     segmentID,
		Data:          data,
		State:         state,
		Editor:        editor,
	}
	return s.segmentUsersPublisher.Publish(ctx, e)
}

func (s *FeatureService) BulkDownloadSegmentUsers(
	ctx context.Context,
	req *featureproto.BulkDownloadSegmentUsersRequest,
) (*featureproto.BulkDownloadSegmentUsersResponse, error) {
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId)
	if err != nil {
		return nil, err
	}
	if err := validateBulkDownloadSegmentUsersRequest(req); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, err
	}
	segment, _, err := s.segmentStorage.GetSegment(ctx, req.SegmentId, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrSegmentNotFound) {
			return nil, statusSegmentNotFound.Err()
		}
		s.logger.Error(
			"Failed to get segment",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	if segment.Status != featureproto.Segment_SUCEEDED {
		return nil, statusSegmentStatusNotSuceeded.Err()
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("segment_id", "=", req.SegmentId),
		mysql.NewFilter("state", "=", int32(req.State)),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	users, _, err := s.segmentUserStorage.ListSegmentUsers(
		ctx,
		whereParts,
		nil,
		mysql.QueryNoLimit,
		mysql.QueryNoOffset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list segment users",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
			)...,
		)
		return nil, api.NewGRPCStatus(err).Err()
	}
	var buf bytes.Buffer
	for _, user := range users {
		buf.WriteString(user.UserId + "\n")
	}
	return &featureproto.BulkDownloadSegmentUsersResponse{
		Data: buf.Bytes(),
	}, nil
}
