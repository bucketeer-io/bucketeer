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
//

package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/feature/command"
	"github.com/bucketeer-io/bucketeer/pkg/feature/domain"
	v2fs "github.com/bucketeer-io/bucketeer/pkg/feature/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

const (
	numOfSecretCharsToShow = 5
	maskURI                = "********"
)

var webhookEditor = &eventproto.Editor{
	Email:   "webhook",
	IsAdmin: false,
}

func (s *FeatureService) CreateFlagTrigger(
	ctx context.Context,
	request *featureproto.CreateFlagTriggerRequest,
) (*featureproto.CreateFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}

	if request.CreateFlagTriggerCommand == nil {
		return s.createFlagTriggerNoCommand(ctx, request, editor, localizer)
	}

	if err = validateCreateFlagTriggerCommand(request.CreateFlagTriggerCommand, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	flagTrigger, err := domain.NewFlagTrigger(
		request.EnvironmentId,
		request.CreateFlagTriggerCommand.FeatureId,
		request.CreateFlagTriggerCommand.Type,
		request.CreateFlagTriggerCommand.Action,
		request.CreateFlagTriggerCommand.Description,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create flag trigger",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewFlagTriggerCommandHandler(
			editor,
			flagTrigger,
			s.domainPublisher,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, request.CreateFlagTriggerCommand); err != nil {
			s.logger.Error(
				"Failed to create flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		if err := s.flagTriggerStorage.CreateFlagTrigger(contextWithTx, flagTrigger); err != nil {
			s.logger.Error(
				"Failed to create flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
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
	triggerURL := s.generateTriggerURL(ctx, flagTrigger.Token, false)
	flagTrigger.FlagTrigger.Token = ""
	return &featureproto.CreateFlagTriggerResponse{
		FlagTrigger: flagTrigger.FlagTrigger,
		Url:         triggerURL,
	}, nil
}

func (s *FeatureService) createFlagTriggerNoCommand(
	ctx context.Context,
	req *featureproto.CreateFlagTriggerRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*featureproto.CreateFlagTriggerResponse, error) {
	if err := validateCreateFlagTriggerNoCommand(req, localizer); err != nil {
		s.logger.Error(
			"Error validating create flag trigger request",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	flagTrigger, err := domain.NewFlagTrigger(
		req.EnvironmentId,
		req.FeatureId,
		req.Type,
		req.Action,
		req.Description,
	)
	if err != nil {
		s.logger.Error(
			"Error creating flag trigger",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		if err := flagTrigger.GenerateToken(); err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FEATURE_CREATED,
			&eventproto.FlagTriggerCreatedEvent{
				Id:            flagTrigger.Id,
				FeatureId:     flagTrigger.FeatureId,
				Type:          flagTrigger.Type,
				Action:        flagTrigger.Action,
				Description:   flagTrigger.Description,
				Token:         flagTrigger.Token,
				CreatedAt:     flagTrigger.CreatedAt,
				UpdatedAt:     flagTrigger.UpdatedAt,
				EnvironmentId: flagTrigger.EnvironmentId,
			},
			req.EnvironmentId,
			flagTrigger,
			&domain.FlagTrigger{},
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.CreateFlagTrigger(contextWithTx, flagTrigger); err != nil {
			s.logger.Error(
				"Failed to create flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", req.EnvironmentId),
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
	triggerURL := s.generateTriggerURL(ctx, flagTrigger.Token, false)
	flagTrigger.FlagTrigger.Token = ""

	if err = s.domainPublisher.Publish(ctx, event); err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
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
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if request.ChangeFlagTriggerDescriptionCommand == nil {
		return s.updateFlagTriggerNoCommand(ctx, request, editor, localizer)
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to get flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewFlagTriggerCommandHandler(
			editor,
			flagTrigger,
			s.domainPublisher,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, request.ChangeFlagTriggerDescriptionCommand); err != nil {
			s.logger.Error(
				"Failed to update flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		if err := s.flagTriggerStorage.UpdateFlagTrigger(
			contextWithTx,
			flagTrigger,
		); err != nil {
			s.logger.Error(
				"Failed to update flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) ||
			errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
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

func (s *FeatureService) updateFlagTriggerNoCommand(
	ctx context.Context,
	request *featureproto.UpdateFlagTriggerRequest,
	editor *eventproto.Editor,
	localizer locale.Localizer,
) (*featureproto.UpdateFlagTriggerResponse, error) {
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		updated, err := flagTrigger.UpdateFlagTrigger(
			request.Description,
			request.Reset_,
			request.Disabled,
		)
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_UPDATED,
			&eventproto.FlagTriggerUpdateEvent{
				Id:          updated.Id,
				FeatureId:   updated.FeatureId,
				Description: request.Description,
				Reset_:      request.Reset_,
				Disabled:    request.Disabled,
			},
			request.EnvironmentId,
			updated,
			flagTrigger,
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.UpdateFlagTrigger(
			contextWithTx,
			updated,
		); err != nil {
			s.logger.Error(
				"Failed to update flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2fs.ErrFlagTriggerUnexpectedAffectedRows) ||
			errors.Is(err, v2fs.ErrFlagTriggerNotFound) {
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
	if err = s.domainPublisher.Publish(ctx, event); err != nil {
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
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableFlagTriggerCommand(request.EnableFlagTriggerCommand, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to get flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewFlagTriggerCommandHandler(
			editor,
			flagTrigger,
			s.domainPublisher,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, request.EnableFlagTriggerCommand); err != nil {
			s.logger.Error(
				"Failed to enable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		if err := s.flagTriggerStorage.UpdateFlagTrigger(
			contextWithTx,
			flagTrigger,
		); err != nil {
			s.logger.Error(
				"Failed to enable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
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
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableFlagTriggerCommand(request.DisableFlagTriggerCommand, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to get flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		handler, err := command.NewFlagTriggerCommandHandler(
			editor,
			flagTrigger,
			s.domainPublisher,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, request.DisableFlagTriggerCommand); err != nil {
			s.logger.Error(
				"Failed to enable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		if err := s.flagTriggerStorage.UpdateFlagTrigger(contextWithTx, flagTrigger); err != nil {
			s.logger.Error(
				"Failed to disable flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
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
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateResetFlagTriggerCommand(request.ResetFlagTriggerCommand, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	trigger, err := s.flagTriggerStorage.GetFlagTrigger(ctx, request.Id, request.EnvironmentId)
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
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		handler, err := command.NewFlagTriggerCommandHandler(
			editor,
			trigger,
			s.domainPublisher,
			request.EnvironmentId,
		)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, request.ResetFlagTriggerCommand); err != nil {
			s.logger.Error(
				"Failed to reset flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
			)
			return err
		}
		err = s.flagTriggerStorage.UpdateFlagTrigger(contextWithTx, trigger)
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
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	triggerURL := s.generateTriggerURL(ctx, trigger.Token, false)
	trigger.FlagTrigger.Token = ""
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
	editor, err := s.checkEnvironmentRole(
		ctx,
		accountproto.AccountV2_Role_Environment_EDITOR,
		request.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	var event *eventproto.Event
	err = s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		flagTrigger, err := s.flagTriggerStorage.GetFlagTrigger(
			contextWithTx,
			request.Id,
			request.EnvironmentId,
		)
		if err != nil {
			s.logger.Error(
				"Failed to get flag trigger",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", request.EnvironmentId),
				)...,
			)
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_DELETED,
			&eventproto.FlagTriggerDeletedEvent{
				Id:            flagTrigger.Id,
				FeatureId:     flagTrigger.FeatureId,
				EnvironmentId: flagTrigger.EnvironmentId,
			},
			request.EnvironmentId,
			&domain.FlagTrigger{},
			flagTrigger,
		)
		if err != nil {
			return err
		}
		if err := s.flagTriggerStorage.DeleteFlagTrigger(
			contextWithTx,
			flagTrigger.Id,
			flagTrigger.EnvironmentId,
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
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return &featureproto.DeleteFlagTriggerResponse{}, nil
}

func (s *FeatureService) GetFlagTrigger(
	ctx context.Context,
	request *featureproto.GetFlagTriggerRequest,
) (*featureproto.GetFlagTriggerResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		request.EnvironmentId,
		localizer)
	if err != nil {
		return nil, err
	}
	if err := validateGetFlagTriggerRequest(request, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", request.EnvironmentId),
			)...,
		)
		return nil, err
	}
	trigger, err := s.flagTriggerStorage.GetFlagTrigger(
		ctx,
		request.Id,
		request.EnvironmentId,
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
	triggerURL := s.generateTriggerURL(ctx, trigger.Token, true)
	trigger.FlagTrigger.Token = ""
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
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		request.EnvironmentId,
		localizer)
	if err != nil {
		return nil, err
	}
	if err := validateListFlagTriggersRequest(request, localizer); err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("feature_id", "=", request.FeatureId),
		mysql.NewFilter("environment_id", "=", request.EnvironmentId),
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
	cursor := request.Cursor
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
	flagTriggers, nextOffset, totalCount, _ := s.flagTriggerStorage.ListFlagTriggers(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	triggerWithUrls := make([]*featureproto.ListFlagTriggersResponse_FlagTriggerWithUrl, 0, len(flagTriggers))
	for _, trigger := range flagTriggers {
		triggerURL := s.generateTriggerURL(ctx, trigger.Token, true)
		trigger.Token = ""
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

func (s *FeatureService) FlagTriggerWebhook(
	ctx context.Context,
	request *featureproto.FlagTriggerWebhookRequest,
) (*featureproto.FlagTriggerWebhookResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	token := request.GetToken()
	resp := &featureproto.FlagTriggerWebhookResponse{}
	if token == "" {
		s.logger.Error(
			"Failed to get secret from query",
			log.FieldsFromImcomingContext(ctx)...,
		)
		dt, err := statusSecretRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "secret"),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	trigger, err := s.flagTriggerStorage.GetFlagTriggerByToken(ctx, token)
	if err != nil {
		s.logger.Error(
			"Failed to get flag trigger",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusTriggerNotFound.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.NotFoundError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	if trigger.GetDisabled() {
		s.logger.Error(
			"Flag trigger is disabled",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusTriggerAlreadyDisabled.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	feature, err := s.featureStorage.GetFeature(ctx, trigger.FeatureId, trigger.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2fs.ErrFeatureNotFound) {
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
			"Failed to get feature",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("id", trigger.FeatureId),
				zap.String("environmentId", trigger.EnvironmentId),
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
	if trigger.GetAction() == featureproto.FlagTrigger_Action_ON {
		// check if feature is already enabled
		if !feature.GetEnabled() {
			err := s.enableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentId(), localizer)
			if err != nil {
				dt, err := statusTriggerEnableFailed.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalize(locale.InternalServerError),
				})
				if err != nil {
					return nil, statusInternal.Err()
				}
				return nil, dt.Err()
			}
		}
	} else if trigger.GetAction() == featureproto.FlagTrigger_Action_OFF {
		// check if feature is already disabled
		if feature.GetEnabled() {
			err := s.disableFeature(ctx, trigger.GetFeatureId(), trigger.GetEnvironmentId(), localizer)
			if err != nil {
				dt, err := statusTriggerDisableFailed.WithDetails(&errdetails.LocalizedMessage{
					Locale:  localizer.GetLocale(),
					Message: localizer.MustLocalize(locale.InternalServerError),
				})
				if err != nil {
					return nil, statusInternal.Err()
				}
				return nil, dt.Err()
			}
		}
	} else {
		s.logger.Error(
			"Invalid action",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		dt, err := statusTriggerActionInvalid.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	err = s.updateTriggerUsageInfo(ctx, webhookEditor, trigger)
	if err != nil {
		dt, err := statusTriggerUsageUpdateFailed.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return nil, statusInternal.Err()
		}
		return nil, dt.Err()
	}
	return resp, nil
}

func (s *FeatureService) updateTriggerUsageInfo(
	ctx context.Context,
	editor *eventproto.Editor,
	flagTrigger *domain.FlagTrigger,
) error {
	var event *eventproto.Event
	err := s.mysqlClient.RunInTransactionV2(ctx, func(contextWithTx context.Context, _ mysql.Transaction) error {
		prev := &domain.FlagTrigger{}
		if err := copier.Copy(prev, flagTrigger); err != nil {
			return err
		}
		err := flagTrigger.UpdateTriggerUsage()
		if err != nil {
			return err
		}
		event, err = domainevent.NewEvent(
			editor,
			eventproto.Event_FLAG_TRIGGER,
			flagTrigger.Id,
			eventproto.Event_FLAG_TRIGGER_USAGE_UPDATED,
			&eventproto.FlagTriggerUsageUpdatedEvent{
				Id:              flagTrigger.Id,
				FeatureId:       flagTrigger.FeatureId,
				LastTriggeredAt: flagTrigger.LastTriggeredAt,
				TriggerTimes:    flagTrigger.TriggerCount,
				EnvironmentId:   flagTrigger.EnvironmentId,
			},
			flagTrigger.EnvironmentId,
			flagTrigger,
			prev,
		)
		if err != nil {
			return err
		}
		err = s.flagTriggerStorage.UpdateFlagTrigger(contextWithTx, flagTrigger)
		if err != nil {
			s.logger.Error(
				"Failed to update flag trigger usage",
				log.FieldsFromImcomingContext(ctx).
					AddFields(zap.Error(err))...,
			)
			return err
		}
		return nil
	})
	if err != nil {
		s.logger.Error(
			"Failed to update flag trigger usage",
			log.FieldsFromImcomingContext(ctx).
				AddFields(zap.Error(err))...,
		)
		return err
	}
	err = s.domainPublisher.Publish(ctx, event)
	if err != nil {
		s.logger.Error(
			"Failed to publish event",
			log.FieldsFromImcomingContext(ctx).
				AddFields(zap.Error(err))...,
		)
		return err
	}
	return nil
}

func (s *FeatureService) enableFeature(
	ctx context.Context,
	featureId, environmentId string,
	localizer locale.Localizer,
) error {
	if err := s.updateFeature(
		ctx,
		&featureproto.EnableFeatureCommand{},
		featureId,
		environmentId,
		"",
		localizer,
		webhookEditor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to enable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
		}
		return err
	}
	return nil
}

func (s *FeatureService) disableFeature(
	ctx context.Context,
	featureId, environmentId string,
	localizer locale.Localizer,
) error {
	if err := s.updateFeature(
		ctx,
		&featureproto.DisableFeatureCommand{},
		featureId,
		environmentId,
		"",
		localizer,
		webhookEditor,
	); err != nil {
		if status.Code(err) == codes.Internal {
			s.logger.Error(
				"Failed to disable feature",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentId),
				)...,
			)
		}
		return err
	}
	return nil
}

func (s *FeatureService) generateTriggerURL(
	ctx context.Context,
	token string,
	masked bool,
) string {
	if masked {
		return fmt.Sprintf("%s/%s", s.triggerURL, token[:numOfSecretCharsToShow]+maskURI)
	} else {
		return fmt.Sprintf("%s/%s", s.triggerURL, token)
	}
}
