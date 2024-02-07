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
	"bytes"
	"context"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	serviceeventproto "github.com/bucketeer-io/bucketeer/proto/event/service"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func (s *FeatureService) AddSegmentUser(
	ctx context.Context,
	req *featureproto.AddSegmentUserRequest,
) (*featureproto.AddSegmentUserResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateAddSegmentUserRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := validateAddSegmentUserCommand(req.Command, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		req.EnvironmentNamespace,
		localizer,
	); err != nil {
		return nil, err
	}
	return &featureproto.AddSegmentUserResponse{}, nil
}

func (s *FeatureService) DeleteSegmentUser(
	ctx context.Context,
	req *featureproto.DeleteSegmentUserRequest,
) (*featureproto.DeleteSegmentUserResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteSegmentUserRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := validateDeleteSegmentUserCommand(req.Command, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
		req.EnvironmentNamespace,
		localizer,
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
	environmentNamespace string,
	localizer locale.Localizer,
) error {
	segmentUsers := make([]*featureproto.SegmentUser, 0, len(userIDs))
	for _, userID := range userIDs {
		userID = strings.TrimSpace(userID)
		user := domain.NewSegmentUser(segmentID, userID, state, deleted)
		segmentUsers = append(segmentUsers, user.SegmentUser)
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		segmentStorage := v2fs.NewSegmentStorage(tx)
		segment, _, err := segmentStorage.GetSegment(ctx, segmentID, environmentNamespace)
		if err != nil {
			s.logger.Error(
				"Failed to get segment",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		segmentUserStorage := v2fs.NewSegmentUserStorage(tx)
		if err := segmentUserStorage.UpsertSegmentUsers(ctx, segmentUsers, environmentNamespace); err != nil {
			s.logger.Error(
				"Failed to store segment user",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		handler := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			environmentNamespace,
		)
		if err := handler.Handle(ctx, cmd); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", environmentNamespace),
				)...,
			)
			return err
		}
		if err := segmentStorage.UpdateSegment(ctx, segment, environmentNamespace); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if err == v2fs.ErrSegmentNotFound || err == v2fs.ErrSegmentUnexpectedAffectedRows {
			dt, err := statusSegmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to upsert segment user",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", environmentNamespace),
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
	return nil
}

func (s *FeatureService) GetSegmentUser(
	ctx context.Context,
	req *featureproto.GetSegmentUserRequest,
) (*featureproto.GetSegmentUserResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetSegmentUserRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	segmentUserStorage := v2fs.NewSegmentUserStorage(s.mysqlClient)
	id := domain.SegmentUserID(req.SegmentId, req.UserId, req.State)
	user, err := segmentUserStorage.GetSegmentUser(ctx, id, req.EnvironmentNamespace)
	if err != nil {
		if err == v2fs.ErrSegmentUserNotFound {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get segment user",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &featureproto.GetSegmentUserResponse{
		User: user.SegmentUser,
	}, nil
}

func (s *FeatureService) ListSegmentUsers(
	ctx context.Context,
	req *featureproto.ListSegmentUsersRequest,
) (*featureproto.ListSegmentUsersResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateListSegmentUsersRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("segment_id", "=", req.SegmentId),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
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
		dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "cursor"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	segmentUserStorage := v2fs.NewSegmentUserStorage(s.mysqlClient)
	users, nextCursor, err := segmentUserStorage.ListSegmentUsers(
		ctx,
		whereParts,
		nil,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &featureproto.ListSegmentUsersResponse{
		Users:  users,
		Cursor: strconv.Itoa(nextCursor),
	}, nil
}

func (s *FeatureService) BulkUploadSegmentUsers(
	ctx context.Context,
	req *featureproto.BulkUploadSegmentUsersRequest,
) (*featureproto.BulkUploadSegmentUsersResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateBulkUploadSegmentUsersRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	if err := validateBulkUploadSegmentUsersCommand(req.Command, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
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
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		segmentStorage := v2fs.NewSegmentStorage(tx)
		segment, _, err := segmentStorage.GetSegment(ctx, req.SegmentId, req.EnvironmentNamespace)
		if err != nil {
			return err
		}
		if segment.IsInUseStatus {
			dt, err := statusSegmentInUse.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.SegmentInUse),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if segment.Status == featureproto.Segment_UPLOADING {
			dt, err := statusSegmentUsersAlreadyUploading.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.SegmentUsersAlreadyUploading),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		handler := command.NewSegmentCommandHandler(
			editor,
			segment,
			s.domainPublisher,
			req.EnvironmentNamespace,
		)
		if err := handler.Handle(ctx, req.Command); err != nil {
			s.logger.Error(
				"Failed to handle command",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", req.EnvironmentNamespace),
				)...,
			)
			return err
		}
		if err := segmentStorage.UpdateSegment(ctx, segment, req.EnvironmentNamespace); err != nil {
			return err
		}
		return s.publishBulkSegmentUsersReceivedEvent(
			ctx,
			editor,
			req.EnvironmentNamespace,
			req.SegmentId,
			req.Command.Data,
			req.Command.State,
		)
	})
	if err != nil {
		if err == v2fs.ErrSegmentNotFound || err == v2fs.ErrFeatureUnexpectedAffectedRows {
			dt, err := statusSegmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		if status.Code(err) == codes.FailedPrecondition {
			return nil, err
		}
		s.logger.Error(
			"Failed to bulk upload segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	return &featureproto.BulkUploadSegmentUsersResponse{}, nil
}

func (s *FeatureService) publishBulkSegmentUsersReceivedEvent(
	ctx context.Context,
	editor *eventproto.Editor,
	environmentNamespace string,
	segmentID string,
	data []byte,
	state featureproto.SegmentUser_State,
) error {
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	e := &serviceeventproto.BulkSegmentUsersReceivedEvent{
		Id:                   id.String(),
		EnvironmentNamespace: environmentNamespace,
		SegmentId:            segmentID,
		Data:                 data,
		State:                state,
		Editor:               editor,
	}
	return s.segmentUsersPublisher.Publish(ctx, e)
}

func (s *FeatureService) BulkDownloadSegmentUsers(
	ctx context.Context,
	req *featureproto.BulkDownloadSegmentUsersRequest,
) (*featureproto.BulkDownloadSegmentUsersResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateBulkDownloadSegmentUsersRequest(req, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	segmentStorage := v2fs.NewSegmentStorage(s.mysqlClient)
	segment, _, err := segmentStorage.GetSegment(ctx, req.SegmentId, req.EnvironmentNamespace)
	if err != nil {
		if err == v2fs.ErrSegmentNotFound {
			dt, err := statusSegmentNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to get segment",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	if segment.Status != featureproto.Segment_SUCEEDED {
		dt, err := statusSegmentStatusNotSuceeded.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.SegmentStatusNotSucceeded),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("segment_id", "=", req.SegmentId),
		mysql.NewFilter("state", "=", int32(req.State)),
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_namespace", "=", req.EnvironmentNamespace),
	}
	segmentUserStorage := v2fs.NewSegmentUserStorage(s.mysqlClient)
	users, _, err := segmentUserStorage.ListSegmentUsers(
		ctx,
		whereParts,
		nil,
		mysql.QueryNoLimit,
		mysql.QueryNoOffset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list segment users",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", req.EnvironmentNamespace),
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
	var buf bytes.Buffer
	for _, user := range users {
		buf.WriteString(user.UserId + "\n")
	}
	return &featureproto.BulkDownloadSegmentUsersResponse{
		Data: buf.Bytes(),
	}, nil
}
