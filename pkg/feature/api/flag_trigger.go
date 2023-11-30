// Copyright 2023 The Bucketeer Authors.
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
//

package api

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const mask = "********"

var maskLength = len(mask)

func (s *FeatureService) CreateFlagTrigger(
	ctx context.Context,
	request *featureproto.CreateFlagTriggerRequest,
) (*featureproto.CreateFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err = validateCreateFlagTriggerCommand(request.CreateFlagTriggerCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	triggerID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(
			"Failed to generate trigger id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	id := triggerID.String()
	triggerUUID, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(
			"Failed to generate trigger uuid",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	flagTriggerUuid := triggerUUID.String()
	flagTrigger := domain.NewFlagTrigger(
		id,
		request.EnvironmentNamespace,
		flagTriggerUuid,
		request.CreateFlagTriggerCommand,
	)
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		if err := storage.CreateFlagTrigger(ctx, flagTrigger); err != nil {
			s.logger.Error(
				"Failed to create flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", request.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerAlreadyExists) {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	flagTriggerSecret := domain.NewFlagTriggerSecret(
		id,
		request.CreateFlagTriggerCommand.FeatureId,
		request.EnvironmentNamespace,
		flagTriggerUuid,
		int(request.CreateFlagTriggerCommand.Action),
	)
	triggerURL, err := s.generateTriggerURL(ctx, flagTriggerSecret, false)
	if err != nil {
		s.logger.Error(
			"Failed to generate trigger url",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return &featureproto.CreateFlagTriggerResponse{
		FlagTrigger: flagTrigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) UpdateFlagTrigger(
	ctx context.Context,
	request *featureproto.UpdateFlagTriggerRequest,
) (*featureproto.UpdateFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateUpdateFlagTriggerCommand(request.ChangeFlagTriggerDescriptionCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
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
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		if err := storage.UpdateFlagTrigger(
			ctx,
			request.ChangeFlagTriggerDescriptionCommand.Id,
			request.EnvironmentNamespace,
			request.ChangeFlagTriggerDescriptionCommand.Description,
		); err != nil {
			s.logger.Error(
				"Failed to update flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", request.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.UpdateFlagTriggerResponse{}, nil
}

func (s *FeatureService) EnableFlagTrigger(
	ctx context.Context,
	request *featureproto.EnableFlagTriggerRequest,
) (*featureproto.EnableFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateEnableFlagTriggerCommand(request.EnableFlagTriggerCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
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
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		if err := storage.EnableFlagTrigger(
			ctx,
			request.EnableFlagTriggerCommand.Id,
			request.EnvironmentNamespace,
		); err != nil {
			s.logger.Error(
				"Failed to enable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", request.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.EnableFlagTriggerResponse{}, nil
}

func (s *FeatureService) DisableFlagTrigger(
	ctx context.Context,
	request *featureproto.DisableFlagTriggerRequest,
) (*featureproto.DisableFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDisableFlagTriggerCommand(request.DisableFlagTriggerCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
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
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		if err := storage.DisableFlagTrigger(
			ctx,
			request.DisableFlagTriggerCommand.Id,
			request.EnvironmentNamespace,
		); err != nil {
			s.logger.Error(
				"Failed to disable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentNamespace", request.EnvironmentNamespace),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.DisableFlagTriggerResponse{}, nil
}

func (s *FeatureService) ResetFlagTrigger(
	ctx context.Context,
	request *featureproto.ResetFlagTriggerRequest,
) (*featureproto.ResetFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateResetFlagTriggerCommand(request.ResetFlagTriggerCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	trigger, err := v2fs.NewFlagTriggerStorage(s.mysqlClient).
		GetFlagTrigger(ctx, request.ResetFlagTriggerCommand.Id, request.EnvironmentNamespace)
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, err
	}
	newTriggerUuid, err := uuid.NewUUID()
	if err != nil {
		s.logger.Error(
			"Failed to generate new trigger uuid",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	newFlagTriggerId := newTriggerUuid.String()
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		err := v2fs.NewFlagTriggerStorage(tx).
			ResetFlagTrigger(ctx, trigger.Id, request.EnvironmentNamespace, newFlagTriggerId)
		if err != nil {
			s.logger.Error(
				"Failed to reset flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, err
	}
	flagTriggerSecret := domain.NewFlagTriggerSecret(
		trigger.Id,
		trigger.FeatureId,
		trigger.EnvironmentNamespace,
		newFlagTriggerId,
		int(trigger.Action),
	)
	triggerURL, err := s.generateTriggerURL(ctx, flagTriggerSecret, false)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return &featureproto.ResetFlagTriggerResponse{
		FlagTrigger: trigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) DeleteFlagTrigger(
	ctx context.Context,
	request *featureproto.DeleteFlagTriggerRequest,
) (*featureproto.DeleteFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_EDITOR, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteFlagTriggerCommand(request.DeleteFlagTriggerCommand, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	tx, err := s.mysqlClient.BeginTx(ctx)
	if err != nil {
		s.logger.Error(
			"Failed to begin transaction",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransaction(ctx, tx, func() error {
		storage := v2fs.NewFlagTriggerStorage(tx)
		if err := storage.DeleteFlagTrigger(
			ctx,
			request.DeleteFlagTriggerCommand.Id,
			request.EnvironmentNamespace,
		); err != nil {
			s.logger.Error(
				"Failed to delete flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, err
	}
	return &featureproto.DeleteFlagTriggerResponse{}, nil
}

func (s *FeatureService) GetFlagTrigger(
	ctx context.Context,
	request *featureproto.GetFlagTriggerRequest,
) (*featureproto.GetFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetFlagTriggerRequest(request, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentNamespace", request.EnvironmentNamespace),
			)...,
		)
		return nil, err
	}
	trigger, err := v2fs.NewFlagTriggerStorage(s.mysqlClient).GetFlagTrigger(
		ctx,
		request.Id,
		request.EnvironmentNamespace,
	)
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		return nil, err
	}
	flagTriggerSecret := domain.NewFlagTriggerSecret(
		trigger.Id,
		trigger.FeatureId,
		trigger.EnvironmentNamespace,
		trigger.Uuid,
		int(trigger.Action),
	)
	triggerURL, err := s.generateTriggerURL(ctx, flagTriggerSecret, true)
	if err != nil {
		s.logger.Error(
			"Failed to generate trigger url",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	return &featureproto.GetFlagTriggerResponse{
		FlagTrigger: trigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) ListFlagTriggers(
	ctx context.Context,
	request *featureproto.ListFlagTriggersRequest,
) (*featureproto.ListFlagTriggersResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkRole(ctx, accountproto.Account_VIEWER, request.EnvironmentNamespace, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateListFlagTriggersRequest(request, localizer); err != nil {
		s.logger.Info(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("feature_id", "=", request.FeatureId),
		mysql.NewFilter("environment_namespace", "=", request.EnvironmentNamespace),
		mysql.NewFilter("deleted", "=", false),
	}
	orders, err := s.newListFlagTriggerOrders(request.OrderBy, request.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to create order",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	limit := int(request.PageSize)
	offset := int(request.Cursor)
	flagTriggers, nextOffset, totalCount, _ := v2fs.NewFlagTriggerStorage(s.mysqlClient).ListFlagTriggers(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	triggerWithUrls := make([]*featureproto.ListFlagTriggersResponse_FlagTriggerWithUrl, 0, len(flagTriggers))
	for _, trigger := range flagTriggers {
		flagTriggerSecret := domain.NewFlagTriggerSecret(
			trigger.Id,
			trigger.FeatureId,
			trigger.EnvironmentNamespace,
			trigger.Uuid,
			int(trigger.Action),
		)
		triggerURL, err := s.generateTriggerURL(ctx, flagTriggerSecret, true)
		if err != nil {
			s.logger.Error(
				"Failed to generate trigger url",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return nil, err
		}
		triggerWithUrls = append(triggerWithUrls, &featureproto.ListFlagTriggersResponse_FlagTriggerWithUrl{
			FlagTrigger: trigger,
			Url:         triggerURL,
		})
	}
	return &featureproto.ListFlagTriggersResponse{
		FlagTriggers: triggerWithUrls,
		Cursor:       strconv.Itoa(nextOffset),
		TotalCount:   totalCount,
	}, nil
}

func (s *FeatureService) newListFlagTriggerOrders(
	orderBy featureproto.ListFlagTriggersRequest_OrderBy,
	orderDirection featureproto.ListFlagTriggersRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case featureproto.ListFlagTriggersRequest_DEFAULT, featureproto.ListFlagTriggersRequest_CREATED_AT:
		column = "created_at"
	case featureproto.ListFlagTriggersRequest_UPDATED_AT:
		column = "updated_at"
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
	if orderDirection == featureproto.ListFlagTriggersRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{
		mysql.NewOrder(column, direction),
	}, nil
}

func (s *FeatureService) generateTriggerURL(
	ctx context.Context,
	flagTriggerSecret *domain.FlagTriggerSecret,
	masked bool,
) (string, error) {
	encoded, err := flagTriggerSecret.Marshal()
	if err != nil {
		return "", err
	}
	encrypted, err := s.triggerCryptoUtil.Encrypt(ctx, encoded)
	if err != nil {
		return "", err
	}
	secret := base64.RawURLEncoding.EncodeToString(encrypted)
	if masked {
		return fmt.Sprintf("%s/%s", s.triggerURL, mask+secret[maskLength:]), nil
	} else {
		return fmt.Sprintf("%s/%s", s.triggerURL, secret), nil
	}
}
