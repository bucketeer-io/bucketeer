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
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/coderef/domain"
	"github.com/bucketeer-io/bucketeer/pkg/coderef/storage"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	proto "github.com/bucketeer-io/bucketeer/proto/coderef"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
)

// generateSourceURL generates a URL to view the code in the repository web interface
func generateSourceURL(codeRef *proto.CodeReference) string {
	switch codeRef.RepositoryType {
	case proto.CodeReference_GITHUB:
		return fmt.Sprintf("https://github.com/%s/%s/blob/%s/%s#L%d",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.CommitHash,
			codeRef.FilePath,
			codeRef.LineNumber)
	case proto.CodeReference_GITLAB:
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/blob/%s/%s#L%d",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.CommitHash,
			codeRef.FilePath,
			codeRef.LineNumber)
	case proto.CodeReference_BITBUCKET:
		return fmt.Sprintf("https://bitbucket.org/%s/%s/src/%s/%s#lines-%d",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.CommitHash,
			codeRef.FilePath,
			codeRef.LineNumber)
	default:
		return ""
	}
}

// generateBranchURL generates a URL to view the branch in the repository web interface
func generateBranchURL(codeRef *proto.CodeReference) string {
	switch codeRef.RepositoryType {
	case proto.CodeReference_GITHUB:
		return fmt.Sprintf("https://github.com/%s/%s/tree/%s",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.RepositoryBranch)
	case proto.CodeReference_GITLAB:
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/tree/%s",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.RepositoryBranch)
	case proto.CodeReference_BITBUCKET:
		return fmt.Sprintf("https://bitbucket.org/%s/%s/src/%s",
			codeRef.RepositoryOwner,
			codeRef.RepositoryName,
			codeRef.RepositoryBranch)
	default:
		return ""
	}
}

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
	codeRefStorage := storage.NewCodeReferenceStorage(s.mysqlClient)
	codeRef, err := codeRefStorage.GetCodeReference(ctx, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrCodeReferenceNotFound) {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	codeRef.CodeReference.SourceUrl = generateSourceURL(&codeRef.CodeReference)
	codeRef.CodeReference.BranchUrl = generateBranchURL(&codeRef.CodeReference)
	return &proto.GetCodeReferenceResponse{CodeReference: &codeRef.CodeReference}, nil
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
		mysql.NewFilter("feature_id", "=", req.FeatureId),
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
	if req.FileExtension != "" {
		whereParts = append(whereParts, mysql.NewFilter("file_extension", "=", req.FileExtension))
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
	codeRefStorage := storage.NewCodeReferenceStorage(s.mysqlClient)
	codeRefs, nextCursor, totalCount, err := codeRefStorage.ListCodeReferences(
		ctx,
		whereParts,
		orders,
		limit,
		cursor,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list code references",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
		ref.CodeReference.SourceUrl = generateSourceURL(&ref.CodeReference)
		ref.CodeReference.BranchUrl = generateBranchURL(&ref.CodeReference)
		protoRefs = append(protoRefs, &ref.CodeReference)
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
	codeRef, err := domain.NewCodeReference(
		req.FeatureId,
		req.FilePath,
		req.FileExtension,
		req.LineNumber,
		req.CodeSnippet,
		req.ContentHash,
		req.Aliases,
		req.RepositoryName,
		req.RepositoryOwner,
		req.RepositoryType,
		req.RepositoryBranch,
		req.CommitHash,
		req.EnvironmentId,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create code reference",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	codeRefStorage := storage.NewCodeReferenceStorage(s.mysqlClient)
	if err := codeRefStorage.CreateCodeReference(ctx, codeRef); err != nil {
		s.logger.Error(
			"Failed to create code reference",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", codeRef.Id),
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
	createEvent, err := domainevent.NewEvent(
		editor,
		eventproto.Event_CODEREF,
		codeRef.Id,
		eventproto.Event_CODE_REFERENCE_CREATED,
		&eventproto.CodeReferenceCreatedEvent{
			Id:               codeRef.Id,
			FeatureId:        codeRef.FeatureId,
			FilePath:         codeRef.FilePath,
			LineNumber:       codeRef.LineNumber,
			CodeSnippet:      codeRef.CodeSnippet,
			ContentHash:      codeRef.ContentHash,
			Aliases:          codeRef.Aliases,
			RepositoryName:   codeRef.RepositoryName,
			RepositoryOwner:  codeRef.RepositoryOwner,
			RepositoryType:   codeRef.RepositoryType,
			RepositoryBranch: codeRef.RepositoryBranch,
			CommitHash:       codeRef.CommitHash,
			EnvironmentId:    codeRef.EnvironmentId,
			CreatedAt:        codeRef.CreatedAt,
			UpdatedAt:        codeRef.UpdatedAt,
		},
		req.EnvironmentId,
		codeRef,
		&proto.CodeReference{},
	)
	if err != nil {
		s.logger.Error(
			"Failed to create event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	if err := s.publisher.Publish(ctx, createEvent); err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	return &proto.CreateCodeReferenceResponse{CodeReference: &codeRef.CodeReference}, nil
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
	codeRefStorage := storage.NewCodeReferenceStorage(s.mysqlClient)
	var codeRef *domain.CodeReference
	var updatedCodeRef *domain.CodeReference
	err = codeRefStorage.RunInTransaction(ctx, func() error {
		var err error
		codeRef, err = codeRefStorage.GetCodeReference(ctx, req.Id)
		if err != nil {
			if errors.Is(err, storage.ErrCodeReferenceNotFound) {
				dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalize(locale.NotFoundError),
				})
				if err != nil {
					return statusInternal.Err()
				}
				return dt.Err()
			}
			s.logger.Error(
				"Failed to get code reference",
				log.FieldsFromIncomingContext(ctx).AddFields(
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
				return statusInternal.Err()
			}
			return dt.Err()
		}
		updatedCodeRef, err = codeRef.Update(
			req.FilePath,
			req.FileExtension,
			req.LineNumber,
			req.CodeSnippet,
			req.ContentHash,
			req.Aliases,
			req.RepositoryBranch,
			req.CommitHash,
		)
		if err != nil {
			s.logger.Error(
				"Failed to update code reference domain object",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
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
		if err := codeRefStorage.UpdateCodeReference(ctx, updatedCodeRef); err != nil {
			s.logger.Error(
				"Failed to update code reference",
				log.FieldsFromIncomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
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
	})
	if err != nil {
		return nil, err
	}
	updateEvent, err := domainevent.NewEvent(
		editor,
		eventproto.Event_CODEREF,
		updatedCodeRef.Id,
		eventproto.Event_CODE_REFERENCE_UPDATED,
		&eventproto.CodeReferenceUpdatedEvent{
			Id:               updatedCodeRef.Id,
			FilePath:         updatedCodeRef.FilePath,
			LineNumber:       updatedCodeRef.LineNumber,
			CodeSnippet:      updatedCodeRef.CodeSnippet,
			ContentHash:      updatedCodeRef.ContentHash,
			Aliases:          updatedCodeRef.Aliases,
			RepositoryBranch: updatedCodeRef.RepositoryBranch,
			CommitHash:       updatedCodeRef.CommitHash,
			EnvironmentId:    updatedCodeRef.EnvironmentId,
			UpdatedAt:        updatedCodeRef.UpdatedAt,
		},
		req.EnvironmentId,
		updatedCodeRef,
		codeRef,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	if err := s.publisher.Publish(ctx, updateEvent); err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	return &proto.UpdateCodeReferenceResponse{CodeReference: &updatedCodeRef.CodeReference}, nil
}

func (s *CodeReferenceService) DeleteCodeReference(
	ctx context.Context,
	req *proto.DeleteCodeReferenceRequest,
) (*proto.DeleteCodeReferenceResponse, error) {
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
	if err := validateDeleteCodeReferenceRequest(req, localizer); err != nil {
		return nil, err
	}
	codeRefStorage := storage.NewCodeReferenceStorage(s.mysqlClient)
	codeRef, err := codeRefStorage.GetCodeReference(ctx, req.Id)
	if err != nil {
		if errors.Is(err, storage.ErrCodeReferenceNotFound) {
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
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	if err := codeRefStorage.DeleteCodeReference(ctx, codeRef.Id); err != nil {
		s.logger.Error(
			"Failed to delete code reference",
			log.FieldsFromIncomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", req.Id),
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
	deleteEvent, err := domainevent.NewEvent(
		editor,
		eventproto.Event_CODEREF,
		codeRef.Id,
		eventproto.Event_CODE_REFERENCE_DELETED,
		&eventproto.CodeReferenceDeletedEvent{
			Id:            codeRef.Id,
			EnvironmentId: codeRef.EnvironmentId,
		},
		req.EnvironmentId,
		nil,
		codeRef,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	if err := s.publisher.Publish(ctx, deleteEvent); err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromIncomingContext(ctx).AddFields(
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
	return &proto.DeleteCodeReferenceResponse{}, nil
}

func (s *CodeReferenceService) toMySQLOrderDirection(
	d proto.ListCodeReferencesRequest_OrderDirection,
) mysql.OrderDirection {
	if d == proto.ListCodeReferencesRequest_DESC {
		return mysql.OrderDirectionDesc
	}
	return mysql.OrderDirectionAsc
}
