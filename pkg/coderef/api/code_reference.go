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
	"context"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/coderef/command"
	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	v2 "github.com/bucketeer-io/bucketeer/pkg/coderef/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/coderef"
)

func (s *CodeReferenceService) GetCodeReference(
	ctx context.Context,
	req *proto.GetCodeReferenceRequest,
) (*proto.GetCodeReferenceResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetCodeReferenceRequest(req, localizer); err != nil {
		return nil, err
	}
	codeRef, err := s.codeRefStorage.GetCodeReference(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2.ErrCodeReferenceNotFound) {
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
			"Failed to get code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
				zap.String("environmentId", req.EnvironmentId),
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
	return &proto.GetCodeReferenceResponse{CodeReference: codeRef.ToProto()}, nil
}

func (s *CodeReferenceService) ListCodeReferences(
	ctx context.Context,
	req *proto.ListCodeReferencesRequest,
) (*proto.ListCodeReferencesResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateListCodeReferencesRequest(req, localizer); err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.FeatureId != "" {
		whereParts = append(whereParts, mysql.NewFilter("feature_id", "=", req.FeatureId))
	}
	if req.RepositoryName != "" {
		whereParts = append(whereParts, mysql.NewFilter("repository_name", "=", req.RepositoryName))
	}
	if req.RepositoryOwner != "" {
		whereParts = append(whereParts, mysql.NewFilter("repository_owner", "=", req.RepositoryOwner))
	}
	if req.RepositoryType != proto.CodeReference_REPOSITORY_TYPE_UNSPECIFIED {
		whereParts = append(whereParts, mysql.NewFilter("repository_type", "=", req.RepositoryType))
	}
	if req.RepositoryBranch != "" {
		whereParts = append(whereParts, mysql.NewFilter("repository_branch", "=", req.RepositoryBranch))
	}
	orders := []*mysql.Order{mysql.NewOrder("id", mysql.OrderDirectionAsc)}
	switch req.OrderBy {
	case proto.ListCodeReferencesRequest_CREATED_AT:
		orders = []*mysql.Order{mysql.NewOrder("created_at", s.toMySQLOrderDirection(req.OrderDirection))}
	case proto.ListCodeReferencesRequest_UPDATED_AT:
		orders = []*mysql.Order{mysql.NewOrder("updated_at", s.toMySQLOrderDirection(req.OrderDirection))}
	}
	limit := int(req.PageSize)
	cursor := 0
	if req.Cursor != "" {
		c, err := strconv.Atoi(req.Cursor)
		if err != nil {
			dt, err := statusInvalidCursor.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		cursor = c
	}
	codeRefs, nextCursor, totalCount, err := s.codeRefStorage.ListCodeReferences(
		ctx,
		whereParts,
		orders,
		limit,
		cursor,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list code references",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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
	protoRefs := make([]*proto.CodeReference, 0, len(codeRefs))
	for _, ref := range codeRefs {
		protoRefs = append(protoRefs, ref.ToProto())
	}
	return &proto.ListCodeReferencesResponse{
		CodeReferences: protoRefs,
		Cursor:         strconv.Itoa(nextCursor),
		TotalCount:     totalCount,
	}, nil
}

func (s *CodeReferenceService) CreateCodeReference(
	ctx context.Context,
	req *proto.CreateCodeReferenceRequest,
) (*proto.CreateCodeReferenceResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateCreateCodeReferenceRequest(req, localizer); err != nil {
		return nil, err
	}
	id, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(
			"Failed to generate uuid",
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
	now := time.Now()
	codeRef := domain.NewCodeReferenceWithTime(
		id.String(),
		req.Command.FeatureId,
		req.Command.FilePath,
		req.Command.LineNumber,
		req.Command.CodeSnippet,
		req.Command.ContentHash,
		req.Command.Aliases,
		req.Command.RepositoryName,
		req.Command.RepositoryOwner,
		req.Command.RepositoryType,
		req.Command.RepositoryBranch,
		req.Command.CommitHash,
		req.EnvironmentId,
		now,
	)
	handler, err := command.NewCodeReferenceCommandHandler(editor, codeRef, s.publisher, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to create command handler",
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
	if err := handler.Handle(ctx, req.Command); err != nil {
		s.logger.Error(
			"Failed to handle command",
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
	if err := s.codeRefStorage.CreateCodeReference(ctx, codeRef); err != nil {
		s.logger.Error(
			"Failed to create code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", id.String()),
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
	return &proto.CreateCodeReferenceResponse{CodeReference: codeRef.ToProto()}, nil
}

func (s *CodeReferenceService) UpdateCodeReference(
	ctx context.Context,
	req *proto.UpdateCodeReferenceRequest,
) (*proto.UpdateCodeReferenceResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateUpdateCodeReferenceRequest(req, localizer); err != nil {
		return nil, err
	}
	codeRef, err := s.codeRefStorage.GetCodeReference(ctx, req.Command.Id, req.EnvironmentId)
	if err != nil {
		if err == v2.ErrCodeReferenceNotFound {
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
			"Failed to get code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Command.Id),
				zap.String("environmentId", req.EnvironmentId),
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
	handler, err := command.NewCodeReferenceCommandHandler(editor, codeRef, s.publisher, req.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to create command handler",
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
	if err := handler.Handle(ctx, req.Command); err != nil {
		s.logger.Error(
			"Failed to handle command",
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
	if err := s.codeRefStorage.UpdateCodeReference(ctx, codeRef); err != nil {
		s.logger.Error(
			"Failed to update code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Command.Id),
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
	return &proto.UpdateCodeReferenceResponse{CodeReference: codeRef.ToProto()}, nil
}

func (s *CodeReferenceService) DeleteCodeReference(
	ctx context.Context,
	req *proto.DeleteCodeReferenceRequest,
) (*proto.DeleteCodeReferenceResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		req.Command.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteCodeReferenceRequest(req, localizer); err != nil {
		return nil, err
	}
	codeRef, err := s.codeRefStorage.GetCodeReference(ctx, req.Command.Id, req.Command.EnvironmentId)
	if err != nil {
		if err == v2.ErrCodeReferenceNotFound {
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
			"Failed to get code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Command.Id),
				zap.String("environmentId", req.Command.EnvironmentId),
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
	handler, err := command.NewCodeReferenceCommandHandler(editor, codeRef, s.publisher, req.Command.EnvironmentId)
	if err != nil {
		s.logger.Error(
			"Failed to create command handler",
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
	if err := handler.Handle(ctx, req.Command); err != nil {
		s.logger.Error(
			"Failed to handle command",
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
	if err := s.codeRefStorage.DeleteCodeReference(ctx, codeRef.Id, codeRef.EnvironmentId); err != nil {
		s.logger.Error(
			"Failed to delete code reference",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Command.Id),
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
	return &proto.DeleteCodeReferenceResponse{}, nil
}

func (s *CodeReferenceService) toMySQLOrderDirection(d proto.ListCodeReferencesRequest_OrderDirection) mysql.OrderDirection {
	if d == proto.ListCodeReferencesRequest_DESC {
		return mysql.OrderDirectionDesc
	}
	return mysql.OrderDirectionAsc
}
