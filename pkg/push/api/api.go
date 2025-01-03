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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"golang.org/x/oauth2/google"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	accountclient "github.com/bucketeer-io/bucketeer/pkg/account/client"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	featureclient "github.com/bucketeer-io/bucketeer/pkg/feature/client"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/pubsub/publisher"
	"github.com/bucketeer-io/bucketeer/pkg/push/command"
	"github.com/bucketeer-io/bucketeer/pkg/push/domain"
	v2ps "github.com/bucketeer-io/bucketeer/pkg/push/storage/v2"
	"github.com/bucketeer-io/bucketeer/pkg/role"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	pushproto "github.com/bucketeer-io/bucketeer/proto/push"
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
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createPushNoCommand(ctx, req, localizer, editor)
	}

	if err := s.validateCreatePushRequest(req, localizer); err != nil {
		return nil, err
	}
	push, err := domain.NewPush(
		req.Command.Name,
		string(req.Command.FcmServiceAccount),
		req.Command.Tags,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create a new push",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Strings("tags", req.Command.Tags),
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
	pushes, err := s.listAllPushes(ctx, req.EnvironmentId, localizer)
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
	if err := s.checkFCMServiceAccount(ctx, pushes, req.Command.FcmServiceAccount, localizer); err != nil {
		return nil, err
	}
	err = s.containsTags(pushes, req.Command.Tags, localizer)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			dt, err := statusTagAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to validate tag existence",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.Strings("tags", req.Command.Tags),
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
		pushStorage := v2ps.NewPushStorage(tx)
		if err := pushStorage.CreatePush(ctx, push, req.EnvironmentId); err != nil {
			return err
		}
		handler, err := command.NewPushCommandHandler(editor, push, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, v2ps.ErrPushAlreadyExists) {
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create push",
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

	// For security reasons we remove the service account from the API response
	push.Push.FcmServiceAccount = ""

	return &pushproto.CreatePushResponse{
		Push: push.Push,
	}, nil
}

// createPushNoCommand implement logic without command
func (s *PushService) createPushNoCommand(
	ctx context.Context,
	req *pushproto.CreatePushRequest,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*pushproto.CreatePushResponse, error) {
	if err := s.validateCreatePushNoCommand(req, localizer); err != nil {
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Strings("tags", req.Tags),
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
	pushes, err := s.listAllPushes(ctx, req.EnvironmentId, localizer)
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
	if err := s.checkFCMServiceAccount(ctx, pushes, req.FcmServiceAccount, localizer); err != nil {
		return nil, err
	}
	err = s.containsTags(pushes, req.Tags, localizer)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			dt, err := statusTagAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to validate tag existence",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("name", req.Name),
				zap.Strings("tags", req.Tags),
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

	var event *eventproto.Event
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
		pushStorage := v2ps.NewPushStorage(tx)
		if err := pushStorage.CreatePush(ctx, push, req.EnvironmentId); err != nil {
			return err
		}
		prev := &domain.Push{}
		if err = copier.Copy(prev, push); err != nil {
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
			prev,
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
			dt, err := statusAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		s.logger.Error(
			"Failed to create push",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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

	// For security reasons we remove the service account from the API response
	push.Push.FcmServiceAccount = ""

	return &pushproto.CreatePushResponse{
		Push: push.Push,
	}, nil
}

func (s *PushService) validateCreatePushRequest(req *pushproto.CreatePushRequest, localizer locale.Localizer) error {
	if string(req.Command.FcmServiceAccount) == "" {
		dt, err := statusFCMServiceAccountRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fcm_service_account"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Command.Tags) == 0 {
		dt, err := statusTagsRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Command.Name == "" {
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

func (s *PushService) validateCreatePushNoCommand(req *pushproto.CreatePushRequest, localizer locale.Localizer) error {
	if string(req.FcmServiceAccount) == "" {
		dt, err := statusFCMServiceAccountRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "fcm_service_account"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if len(req.Tags) == 0 {
		dt, err := statusTagsRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.Name == "" {
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

func (s *PushService) UpdatePush(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
) (*pushproto.UpdatePushResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}

	if s.isNoUpdatePushCommand(req) {
		return s.updatePushNoCommand(ctx, req, localizer, editor)
	}

	if err := s.validateUpdatePushRequest(ctx, req, localizer); err != nil {
		return nil, err
	}

	var updatedPushPb *pushproto.Push
	commands := s.createUpdatePushCommands(req)
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
		pushStorage := v2ps.NewPushStorage(tx)
		push, err := pushStorage.GetPush(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		handler, err := command.NewPushCommandHandler(editor, push, s.publisher, req.EnvironmentId)
		if err != nil {
			return err
		}
		for _, command := range commands {
			if err := handler.Handle(ctx, command); err != nil {
				return err
			}
		}
		updatedPushPb = push.Push
		return pushStorage.UpdatePush(ctx, push, req.EnvironmentId)
	})
	if err != nil {
		switch {
		case errors.Is(err, v2ps.ErrPushNotFound):
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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

	if updatedPushPb != nil {
		// For security reasons we remove the service account from the API response
		updatedPushPb.FcmServiceAccount = ""
	}
	return &pushproto.UpdatePushResponse{
		Push: updatedPushPb,
	}, nil
}

func (s *PushService) updatePushNoCommand(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*pushproto.UpdatePushResponse, error) {
	if err := s.validateUpdatePushRequestNoCommand(ctx, req, localizer); err != nil {
		return nil, err
	}
	var updatedPushPb *pushproto.Push
	var updatePushEvent *eventproto.Event
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
		pushStorage := v2ps.NewPushStorage(tx)
		push, err := pushStorage.GetPush(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}

		updated, err := push.Update(req.Name, req.Tags, req.Disabled)
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

		return pushStorage.UpdatePush(ctx, updated, req.EnvironmentId)
	})
	if err != nil {
		switch {
		case errors.Is(err, v2ps.ErrPushNotFound):
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
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
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
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

	if updatedPushPb != nil {
		// For security reasons we remove the service account from the API response
		updatedPushPb.FcmServiceAccount = ""
	}

	return &pushproto.UpdatePushResponse{
		Push: updatedPushPb,
	}, nil
}

func (s *PushService) validateUpdatePushRequest(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if req.DeletePushTagsCommand != nil && len(req.DeletePushTagsCommand.Tags) == 0 {
		dt, err := statusTagsRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	if err := s.validateAddPushTagsCommand(ctx, req, localizer); err != nil {
		return err
	}
	if req.RenamePushCommand != nil && req.RenamePushCommand.Name == "" {
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

func (s *PushService) validateUpdatePushRequestNoCommand(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}

	return nil
}

func (s *PushService) validateAddPushTagsCommand(
	ctx context.Context,
	req *pushproto.UpdatePushRequest,
	localizer locale.Localizer,
) error {
	if req.AddPushTagsCommand == nil {
		return nil
	}
	if len(req.AddPushTagsCommand.Tags) == 0 {
		dt, err := statusTagsRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "tag"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	pushes, err := s.listAllPushes(ctx, req.EnvironmentId, localizer)
	if err != nil {
		dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InternalServerError),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	err = s.containsTags(pushes, req.AddPushTagsCommand.Tags, localizer)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			dt, err := statusTagAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "tag"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		s.logger.Error(
			"Failed to validate tag existence",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentId", req.EnvironmentId),
				zap.String("id", req.Id),
				zap.Strings("tags", req.AddPushTagsCommand.Tags),
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

func (s *PushService) isNoUpdatePushCommand(req *pushproto.UpdatePushRequest) bool {
	return req.AddPushTagsCommand == nil &&
		req.DeletePushTagsCommand == nil &&
		req.RenamePushCommand == nil
}

func (s *PushService) DeletePush(
	ctx context.Context,
	req *pushproto.DeletePushRequest,
) (*pushproto.DeletePushResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_EDITOR,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	if err := validateDeletePushRequest(req, localizer); err != nil {
		return nil, err
	}

	var event *eventproto.Event
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
		pushStorage := v2ps.NewPushStorage(tx)
		push, err := pushStorage.GetPush(ctx, req.Id, req.EnvironmentId)
		if err != nil {
			return err
		}
		prev := &domain.Push{}
		if err = copier.Copy(prev, push); err != nil {
			return err
		}
		push.Deleted = true
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
			push,
			prev,
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, event); err != nil {
			return err
		}
		return pushStorage.UpdatePush(ctx, push, req.EnvironmentId)
	})
	if err != nil {
		switch {
		case errors.Is(err, v2ps.ErrPushNotFound):
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		case errors.Is(err, v2ps.ErrPushUnexpectedAffectedRows):
			return &pushproto.DeletePushResponse{}, nil
		}
		s.logger.Error(
			"Failed to delete push",
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
	return &pushproto.DeletePushResponse{}, nil
}

func (s *PushService) GetPush(
	ctx context.Context,
	req *pushproto.GetPushRequest,
) (*pushproto.GetPushResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	if err := s.validateGetPushRequest(req, localizer); err != nil {
		return nil, err
	}

	pushStorage := v2ps.NewPushStorage(s.mysqlClient)
	push, err := pushStorage.GetPush(ctx, req.Id, req.EnvironmentId)
	if err != nil {
		if errors.Is(err, v2ps.ErrPushNotFound) {
			dt, err := statusNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			s.logger.Error(
				"Failed to get push",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("id", req.Id),
					zap.String("environmentId", req.EnvironmentId),
				)...,
			)
			return nil, dt.Err()
		}
	}

	if push.Push != nil {
		// For security reasons we remove the service account from the API response
		push.Push.FcmServiceAccount = ""
	}

	return &pushproto.GetPushResponse{
		Push: push.Push,
	}, nil
}

func (s *PushService) validateGetPushRequest(
	req *pushproto.GetPushRequest,
	localizer locale.Localizer,
) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func validateDeletePushRequest(req *pushproto.DeletePushRequest, localizer locale.Localizer) error {
	if req.Id == "" {
		dt, err := statusIDRequired.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *PushService) createUpdatePushCommands(req *pushproto.UpdatePushRequest) []command.Command {
	commands := make([]command.Command, 0)
	if req.DeletePushTagsCommand != nil {
		commands = append(commands, req.DeletePushTagsCommand)
	}
	if req.AddPushTagsCommand != nil {
		commands = append(commands, req.AddPushTagsCommand)
	}
	if req.RenamePushCommand != nil {
		commands = append(commands, req.RenamePushCommand)
	}
	return commands
}

func (s *PushService) containsTags(
	pushes []*pushproto.Push,
	tags []string,
	localizer locale.Localizer,
) error {
	m, err := s.tagMap(pushes)
	if err != nil {
		return err
	}
	for _, t := range tags {
		if _, ok := m[t]; ok {
			dt, err := statusTagAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "tag"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
	}
	return nil
}

func (s *PushService) checkFCMServiceAccount(
	ctx context.Context,
	pushes []*pushproto.Push,
	fcmServiceAccount []byte,
	localizer locale.Localizer,
) error {
	// Check if the JSON is a service account file
	_, err := google.CredentialsFromJSON(
		ctx,
		fcmServiceAccount,
		"https://www.googleapis.com/auth/firebase.messaging",
	)
	if err != nil {
		s.logger.Error("failed to get credentials from JSON", zap.Error(err))
		dt, err := statusFCMServiceAccountInvalid.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalizeWithTemplate(locale.InvalidArgumentError, "FCM service account"),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	// Check if the service account already exists in the database
	for _, push := range pushes {
		equal, err := s.compareJSON(push.FcmServiceAccount, string(fcmServiceAccount))
		if err != nil {
			s.logger.Error("failed to compare the JSON", zap.Error(err))
			dt, err := statusInternal.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InternalServerError),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
		}
		if equal {
			s.logger.Error("fcm service account already exists in the database")
			dt, err := statusFCMServiceAccountAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalizeWithTemplate(locale.AlreadyExistsError, "FCM service account"),
			})
			if err != nil {
				return statusInternal.Err()
			}
			return dt.Err()
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
	localizer locale.Localizer,
) ([]*pushproto.Push, error) {
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", environmentId),
	}
	pushes, _, _, err := s.listPushes(
		ctx,
		mysql.QueryNoLimit,
		"",
		environmentId,
		whereParts,
		nil,
		localizer,
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
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkEnvironmentRole(
		ctx, accountproto.AccountV2_Role_Environment_VIEWER,
		req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("deleted", "=", false),
		mysql.NewFilter("environment_id", "=", req.EnvironmentId),
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"name"}, req.SearchKeyword))
	}
	orders, err := s.newListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
		)
		return nil, err
	}
	pushes, cursor, totalCount, err := s.listPushes(
		ctx,
		req.PageSize,
		req.Cursor,
		req.EnvironmentId,
		whereParts,
		orders,
		localizer,
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

func (s *PushService) newListOrders(
	orderBy pushproto.ListPushesRequest_OrderBy,
	orderDirection pushproto.ListPushesRequest_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case pushproto.ListPushesRequest_DEFAULT,
		pushproto.ListPushesRequest_NAME:
		column = "name"
	case pushproto.ListPushesRequest_CREATED_AT:
		column = "created_at"
	case pushproto.ListPushesRequest_UPDATED_AT:
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
	if orderDirection == pushproto.ListPushesRequest_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}

func (s *PushService) listPushes(
	ctx context.Context,
	pageSize int64,
	cursor string,
	environmentId string,
	whereParts []mysql.WherePart,
	orders []*mysql.Order,
	localizer locale.Localizer,
) ([]*pushproto.Push, string, int64, error) {
	limit := int(pageSize)
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
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	pushStorage := v2ps.NewPushStorage(s.mysqlClient)
	pushes, nextCursor, totalCount, err := pushStorage.ListPushes(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list pushes",
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
			return nil, "", 0, statusInternal.Err()
		}
		return nil, "", 0, dt.Err()
	}
	return pushes, strconv.Itoa(nextCursor), totalCount, nil
}

func (s *PushService) checkEnvironmentRole(
	ctx context.Context,
	requiredRole accountproto.AccountV2_Role_Environment,
	environmentId string,
	localizer locale.Localizer,
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
			s.logger.Info(
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
			s.logger.Info(
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
