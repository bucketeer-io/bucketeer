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
	"strings"

	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/bucketeer-io/bucketeer/pkg/account/command"
	"github.com/bucketeer-io/bucketeer/pkg/account/domain"
	v2as "github.com/bucketeer-io/bucketeer/pkg/account/storage/v2"
	domainevent "github.com/bucketeer-io/bucketeer/pkg/domainevent/domain"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/log"
	"github.com/bucketeer-io/bucketeer/pkg/storage"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	tagdomain "github.com/bucketeer-io/bucketeer/pkg/tag/domain"
	accountproto "github.com/bucketeer-io/bucketeer/proto/account"
	eventproto "github.com/bucketeer-io/bucketeer/proto/event/domain"
	tagproto "github.com/bucketeer-io/bucketeer/proto/tag"
)

func (s *AccountService) CreateAccountV2(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
) (*accountproto.CreateAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if req.Command == nil {
		return s.createAccountV2NoCommand(ctx, req, localizer, editor)
	}
	if err := validateCreateAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}
	account := domain.NewAccountV2(
		req.Command.Email,
		req.Command.Name,
		req.Command.FirstName,
		req.Command.LastName,
		req.Command.Language,
		req.Command.AvatarImageUrl,
		req.Command.Tags,
		req.OrganizationId,
		req.Command.OrganizationRole,
		req.Command.EnvironmentRoles,
	)
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		// TODO: temporary implementation: double write account v2 ---
		exist, err := s.accountStorage.GetAccountV2(ctx, account.Email, req.OrganizationId)
		if err != nil && !errors.Is(err, v2as.ErrAccountNotFound) {
			return err
		}
		if exist != nil {
			handler, err := command.NewAccountV2CommandHandler(editor, exist, s.publisher, req.OrganizationId)
			if err != nil {
				return err
			}
			cmd := &accountproto.ChangeAccountV2EnvironmentRolesCommand{
				Roles:     account.EnvironmentRoles,
				WriteType: accountproto.ChangeAccountV2EnvironmentRolesCommand_WriteType_PATCH,
			}
			if err := handler.Handle(ctx, cmd); err != nil {
				return err
			}
			return s.accountStorage.UpdateAccountV2(ctx, exist)
		}
		// TODO: temporary implementation end ---
		handler, err := command.NewAccountV2CommandHandler(editor, account, s.publisher, req.OrganizationId)
		if err != nil {
			return err
		}
		if err := handler.Handle(ctx, req.Command); err != nil {
			return err
		}
		return s.accountStorage.CreateAccountV2(ctx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountAlreadyExists) {
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
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.Any("environmentRoles", req.Command.EnvironmentRoles),
				zap.String("email", req.Command.Email),
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
	// Upsert tags
	for _, envRole := range req.Command.EnvironmentRoles {
		if err := s.upsertTags(ctx, req.Command.Tags, envRole.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to upsert account tags",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationId", req.OrganizationId),
					zap.String("environmentId", envRole.EnvironmentId),
					zap.String("email", req.Command.Email),
					zap.Strings("tags", req.Command.Tags),
				)...,
			)
			return nil, statusInternal.Err()
		}
	}
	return &accountproto.CreateAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) createAccountV2NoCommand(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*accountproto.CreateAccountV2Response, error) {
	err := validateCreateAccountV2NoCommandRequest(req, localizer)
	if err != nil {
		s.logger.Error(
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
			)...,
		)
		return nil, err
	}
	account := domain.NewAccountV2(
		req.Email,
		req.Name,
		req.FirstName,
		req.LastName,
		req.Language,
		req.AvatarImageUrl,
		req.Tags,
		req.OrganizationId,
		req.OrganizationRole,
		req.EnvironmentRoles,
	)
	var createAccountEvent *eventproto.Event
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		// TODO: temporary implementation: double write account v2 ---
		exist, err := s.accountStorage.GetAccountV2(ctx, account.Email, req.OrganizationId)
		if err != nil && !errors.Is(err, v2as.ErrAccountNotFound) {
			return err
		}
		if exist != nil {
			return s.changeExistedAccountV2EnvironmentRoles(ctx, req, exist, editor)
		}
		// TODO: temporary implementation end ---

		createAccountEvent, err = domainevent.NewEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_CREATED,
			&eventproto.AccountV2CreatedEvent{
				Email:            account.Email,
				FirstName:        account.FirstName,
				LastName:         account.LastName,
				Language:         account.Language,
				AvatarImageUrl:   account.AvatarImageUrl,
				OrganizationId:   account.OrganizationId,
				OrganizationRole: account.OrganizationRole,
				EnvironmentRoles: account.EnvironmentRoles,
				Disabled:         account.Disabled,
				CreatedAt:        account.CreatedAt,
				UpdatedAt:        account.UpdatedAt,
			},
			storage.AdminEnvironmentID,
			account,
			nil,
		)
		if err != nil {
			return err
		}
		return s.accountStorage.CreateAccountV2(ctx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountAlreadyExists) {
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
			"Failed to create account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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

	if err = s.publisher.Publish(ctx, createAccountEvent); err != nil {
		s.logger.Error(
			"Failed to publish create account event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	// Upsert tags
	for _, envRole := range req.EnvironmentRoles {
		if err := s.upsertTags(ctx, req.Tags, envRole.EnvironmentId); err != nil {
			s.logger.Error(
				"Failed to upsert account tags",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationId", req.OrganizationId),
					zap.String("environmentId", envRole.EnvironmentId),
					zap.String("email", req.Email),
					zap.Strings("tags", req.Tags),
				)...,
			)
			return nil, statusInternal.Err()
		}
	}
	return &accountproto.CreateAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) changeExistedAccountV2EnvironmentRoles(
	ctx context.Context,
	req *accountproto.CreateAccountV2Request,
	account *domain.AccountV2,
	editor *eventproto.Editor,
) error {
	var updateAccountEvent *eventproto.Event
	updated := &domain.AccountV2{}
	if err := copier.Copy(updated, account); err != nil {
		return err
	}
	err := updated.PatchEnvironmentRole(req.EnvironmentRoles)
	if err != nil {
		return err
	}

	updateAccountEvent, err = domainevent.NewEvent(
		editor,
		eventproto.Event_PUSH,
		updated.Email,
		eventproto.Event_ACCOUNT_V2_ENVIRONMENT_ROLES_CHANGED,
		&eventproto.AccountV2EnvironmentRolesChangedEvent{
			Email:            updated.Email,
			EnvironmentRoles: updated.EnvironmentRoles,
		},
		storage.AdminEnvironmentID,
		updated,
		account,
	)
	if err != nil {
		s.logger.Error(
			"Failed to create update account event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("email", req.Email),
			)...,
		)
		return err
	}
	if err = s.publisher.Publish(ctx, updateAccountEvent); err != nil {
		return err
	}
	err = s.accountStorage.UpdateAccountV2(ctx, updated)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountService) upsertTags(
	ctx context.Context,
	tags []string,
	environmentID string,
) error {
	for _, tag := range tags {
		trimed := strings.TrimSpace(tag)
		if trimed == "" {
			continue
		}
		t, err := tagdomain.NewTag(trimed, environmentID, tagproto.Tag_ACCOUNT)
		if err != nil {
			s.logger.Error(
				"Failed to create domain tag",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("environmentId", environmentID),
					zap.String("tagId", tag),
				)...,
			)
			return err
		}
		if err := s.tagStorage.UpsertTag(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

func (s *AccountService) UpdateAccountV2(
	ctx context.Context,
	req *accountproto.UpdateAccountV2Request,
) (*accountproto.UpdateAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	isAdmin := false
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		// If not admin, check if user is updating their own account
		if editor.Email != req.Email {
			dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.PermissionDenied),
			})
			if err != nil {
				return nil, statusInternal.Err()
			}
			return nil, dt.Err()
		}
		editor, err = s.checkOrganizationRole(
			ctx,
			accountproto.AccountV2_Role_Organization_MEMBER,
			req.OrganizationId,
			localizer,
		)
		if err != nil {
			return nil, err
		}
	} else {
		isAdmin = true
	}

	if !isAdmin {
		if err := s.checkRestrictedCommands(req, localizer); err != nil {
			s.logger.Error(
				"Member user is not allowed to update organization role or environment roles",
				log.FieldsFromImcomingContext(ctx).AddFields(
					zap.Error(err),
					zap.String("organizationID", req.OrganizationId),
					zap.String("email", req.Email),
				)...,
			)
			return nil, err
		}
	}

	if isNoUpdateAccountV2Command(req) {
		return s.updateAccountV2NoCommand(ctx, req, localizer, editor)
	}
	commands := s.getUpdateAccountV2Commands(req)
	if err := validateUpdateAccountV2Request(req, commands, localizer); err != nil {
		s.logger.Error(
			"Failed to update account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	updatedAccountPb, err := s.updateAccountV2MySQL(ctx, editor, commands, req.Email, req.OrganizationId)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to update account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	// Upsert tags
	if req.ChangeTagsCommand != nil {
		for _, envRole := range updatedAccountPb.EnvironmentRoles {
			if err := s.upsertTags(ctx, req.ChangeTagsCommand.Tags, envRole.EnvironmentId); err != nil {
				s.logger.Error(
					"Failed to upsert account tags",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("organizationId", req.OrganizationId),
						zap.String("environmentId", envRole.EnvironmentId),
						zap.String("email", updatedAccountPb.Email),
						zap.Strings("tags", req.ChangeTagsCommand.Tags),
					)...,
				)
				return nil, statusInternal.Err()
			}
		}
	}
	return &accountproto.UpdateAccountV2Response{
		Account: updatedAccountPb,
	}, nil
}

// checkRestrictedCommands checks if the request contains any restricted values changed
// and returns a permission denied error if it does
func (s *AccountService) checkRestrictedCommands(
	req *accountproto.UpdateAccountV2Request,
	localizer locale.Localizer,
) error {
	if req.ChangeOrganizationRoleCommand != nil ||
		req.ChangeEnvironmentRolesCommand != nil ||
		req.OrganizationRole != nil ||
		req.EnvironmentRoles != nil {
		dt, err := statusPermissionDenied.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return statusInternal.Err()
		}
		return dt.Err()
	}
	return nil
}

func (s *AccountService) updateAccountV2NoCommand(
	ctx context.Context,
	req *accountproto.UpdateAccountV2Request,
	localizer locale.Localizer,
	editor *eventproto.Editor,
) (*accountproto.UpdateAccountV2Response, error) {
	err := validateUpdateAccountV2NoCommandRequest(req, localizer)
	if err != nil {
		return nil, err
	}
	updatedAccountPb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		req.Name,
		req.FirstName,
		req.LastName,
		req.Language,
		req.AvatarImageUrl,
		req.Avatar,
		req.Tags,
		req.OrganizationRole,
		req.EnvironmentRoles,
		req.Disabled,
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to update account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	// Upsert tags
	if req.Tags != nil {
		for _, envRole := range updatedAccountPb.EnvironmentRoles {
			if err := s.upsertTags(ctx, req.Tags, envRole.EnvironmentId); err != nil {
				s.logger.Error(
					"Failed to upsert account tags",
					log.FieldsFromImcomingContext(ctx).AddFields(
						zap.Error(err),
						zap.String("organizationId", req.OrganizationId),
						zap.String("environmentId", envRole.EnvironmentId),
						zap.String("email", req.Email),
						zap.Strings("tags", req.Tags),
					)...,
				)
				return nil, statusInternal.Err()
			}
		}
	}
	return &accountproto.UpdateAccountV2Response{
		Account: updatedAccountPb,
	}, nil
}

func isNoUpdateAccountV2Command(req *accountproto.UpdateAccountV2Request) bool {
	return req.ChangeNameCommand == nil &&
		req.ChangeFirstNameCommand == nil &&
		req.ChangeLastNameCommand == nil &&
		req.ChangeLanguageCommand == nil &&
		req.ChangeAvatarUrlCommand == nil &&
		req.ChangeAvatarCommand == nil &&
		req.ChangeTagsCommand == nil &&
		req.ChangeOrganizationRoleCommand == nil &&
		req.ChangeEnvironmentRolesCommand == nil &&
		req.ChangeLastSeenCommand == nil
}

func (s *AccountService) getUpdateAccountV2Commands(req *accountproto.UpdateAccountV2Request) []command.Command {
	commands := make([]command.Command, 0)
	if req.ChangeNameCommand != nil {
		commands = append(commands, req.ChangeNameCommand)
	}
	if req.ChangeFirstNameCommand != nil {
		commands = append(commands, req.ChangeFirstNameCommand)
	}
	if req.ChangeLastNameCommand != nil {
		commands = append(commands, req.ChangeLastNameCommand)
	}
	if req.ChangeLanguageCommand != nil {
		commands = append(commands, req.ChangeLanguageCommand)
	}
	if req.ChangeAvatarUrlCommand != nil {
		commands = append(commands, req.ChangeAvatarUrlCommand)
	}
	if req.ChangeAvatarCommand != nil {
		commands = append(commands, req.ChangeAvatarCommand)
	}
	if req.ChangeTagsCommand != nil {
		commands = append(commands, req.ChangeTagsCommand)
	}
	if req.ChangeOrganizationRoleCommand != nil {
		commands = append(commands, req.ChangeOrganizationRoleCommand)
	}
	if req.ChangeEnvironmentRolesCommand != nil {
		commands = append(commands, req.ChangeEnvironmentRolesCommand)
	}
	if req.ChangeLastSeenCommand != nil {
		commands = append(commands, req.ChangeLastSeenCommand)
	}
	return commands
}

func (s *AccountService) EnableAccountV2(
	ctx context.Context,
	req *accountproto.EnableAccountV2Request,
) (*accountproto.EnableAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateEnableAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	accountV2Pb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		[]string{},
		nil,
		nil,
		wrapperspb.Bool(false),
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to enable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.EnableAccountV2Response{
		Account: accountV2Pb,
	}, nil
}

func (s *AccountService) DisableAccountV2(
	ctx context.Context,
	req *accountproto.DisableAccountV2Request,
) (*accountproto.DisableAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDisableAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	accountV2Pb, err := s.updateAccountV2NoCommandMysql(
		ctx,
		editor,
		req.Email,
		req.OrganizationId,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		[]string{},
		nil,
		nil,
		wrapperspb.Bool(true),
	)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to disable account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.DisableAccountV2Response{
		Account: accountV2Pb,
	}, nil
}

func (s *AccountService) updateAccountV2MySQL(
	ctx context.Context,
	editor *eventproto.Editor,
	commands []command.Command,
	email, organizationID string,
) (*accountproto.AccountV2, error) {
	var updatedAccountPb *accountproto.AccountV2
	err := s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
		if err != nil {
			return err
		}
		handler, err := command.NewAccountV2CommandHandler(editor, account, s.publisher, organizationID)
		if err != nil {
			return err
		}
		for _, c := range commands {
			if err := handler.Handle(ctx, c); err != nil {
				return err
			}
		}
		updatedAccountPb = account.AccountV2
		return s.accountStorage.UpdateAccountV2(ctx, account)
	})
	return updatedAccountPb, err
}

// updateAccountV2NoCommandMysql updates account properties, if the value is nil, it will not be updated.
func (s *AccountService) updateAccountV2NoCommandMysql(
	ctx context.Context,
	editor *eventproto.Editor,
	email, organizationID string,
	name, firstName, lastName, language, avatarImageURL *wrapperspb.StringValue,
	avatar *accountproto.UpdateAccountV2Request_AccountV2Avatar,
	tags []string,
	organizationRole *accountproto.UpdateAccountV2Request_OrganizationRoleValue,
	environmentRoles []*accountproto.AccountV2_EnvironmentRole,
	isDisabled *wrapperspb.BoolValue,
) (*accountproto.AccountV2, error) {
	var updatedAccountPb *accountproto.AccountV2
	var updateAccountV2Event *eventproto.Event
	err := s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
		if err != nil {
			return err
		}
		updated, err := account.Update(
			name,
			firstName,
			lastName,
			language,
			avatarImageURL,
			avatar,
			tags,
			organizationRole,
			environmentRoles,
			isDisabled,
		)
		if err != nil {
			return err
		}
		updateAccountV2Event, err = domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_UPDATED,
			&eventproto.AccountV2UpdatedEvent{
				Email:          updated.Email,
				OrganizationId: updated.OrganizationId,
			},
			updated,
			account,
		)
		if err != nil {
			return err
		}
		updatedAccountPb = updated.AccountV2
		return s.accountStorage.UpdateAccountV2(ctx, updated)
	})
	if err != nil {
		return nil, err
	}
	if err = s.publisher.Publish(ctx, updateAccountV2Event); err != nil {
		s.logger.Error(
			"Failed to publish update account event",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", organizationID),
				zap.String("email", email),
			)...,
		)
		return nil, err
	}
	return updatedAccountPb, err
}

func (s *AccountService) DeleteAccountV2(
	ctx context.Context,
	req *accountproto.DeleteAccountV2Request,
) (*accountproto.DeleteAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	editor, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_ADMIN,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateDeleteAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to delete account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	err = s.accountStorage.RunInTransaction(ctx, func() error {
		account, err := s.accountStorage.GetAccountV2(ctx, req.Email, req.OrganizationId)
		if err != nil {
			return err
		}
		deleteAccount := &domain.AccountV2{}
		if err := copier.Copy(deleteAccount, account); err != nil {
			return err
		}
		deleteAccountV2Event, err := domainevent.NewAdminEvent(
			editor,
			eventproto.Event_ACCOUNT,
			account.Email,
			eventproto.Event_ACCOUNT_V2_DELETED,
			&eventproto.AccountV2UpdatedEvent{
				Email:          account.Email,
				OrganizationId: account.OrganizationId,
			},
			deleteAccount,
			account,
		)
		if err != nil {
			return err
		}
		if err = s.publisher.Publish(ctx, deleteAccountV2Event); err != nil {
			return err
		}
		return s.accountStorage.DeleteAccountV2(ctx, account)
	})
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) || errors.Is(err, v2as.ErrAccountUnexpectedAffectedRows) {
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
			"Failed to delete account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
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
	return &accountproto.DeleteAccountV2Response{}, nil
}

func (s *AccountService) GetAccountV2(
	ctx context.Context,
	req *accountproto.GetAccountV2Request,
) (*accountproto.GetAccountV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2Request(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2(ctx, req.Email, req.OrganizationId, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2Response{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2(
	ctx context.Context,
	email, organizationID string,
	localizer locale.Localizer,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2(ctx, email, organizationID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
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
			"Failed to get account",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", organizationID),
				zap.String("email", email),
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
	return account, nil
}

func (s *AccountService) GetAccountV2ByEnvironmentID(
	ctx context.Context,
	req *accountproto.GetAccountV2ByEnvironmentIDRequest,
) (*accountproto.GetAccountV2ByEnvironmentIDResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRoleByEnvironmentID(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.EnvironmentId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	if err := validateGetAccountV2ByEnvironmentIDRequest(req, localizer); err != nil {
		s.logger.Error(
			"Failed to get account by environment id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("EnvironmentId", req.EnvironmentId),
				zap.String("email", req.Email),
			)...,
		)
		return nil, err
	}
	account, err := s.getAccountV2ByEnvironmentID(ctx, req.Email, req.EnvironmentId, localizer)
	if err != nil {
		return nil, err
	}
	return &accountproto.GetAccountV2ByEnvironmentIDResponse{Account: account.AccountV2}, nil
}

func (s *AccountService) getAccountV2ByEnvironmentID(
	ctx context.Context,
	email, environmentID string,
	localizer locale.Localizer,
) (*domain.AccountV2, error) {
	account, err := s.accountStorage.GetAccountV2ByEnvironmentID(ctx, email, environmentID)
	if err != nil {
		if errors.Is(err, v2as.ErrAccountNotFound) {
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
			"Failed to get account by environment id",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("environmentID", environmentID),
				zap.String("email", email),
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
	return account, nil
}

func (s *AccountService) ListAccountsV2(
	ctx context.Context,
	req *accountproto.ListAccountsV2Request,
) (*accountproto.ListAccountsV2Response, error) {
	localizer := locale.NewLocalizer(ctx)
	_, err := s.checkOrganizationRole(
		ctx,
		accountproto.AccountV2_Role_Organization_MEMBER,
		req.OrganizationId,
		localizer,
	)
	if err != nil {
		return nil, err
	}
	whereParts := []mysql.WherePart{
		mysql.NewFilter("organization_id", "=", req.OrganizationId),
	}
	if req.Disabled != nil {
		whereParts = append(whereParts, mysql.NewFilter("disabled", "=", req.Disabled.Value))
	}
	tagValues := make([]interface{}, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagValues = append(tagValues, tag)
	}
	if len(tagValues) > 0 {
		whereParts = append(
			whereParts,
			mysql.NewJSONFilter("tags", mysql.JSONContainsString, tagValues),
		)
	}
	if req.OrganizationRole != nil {
		whereParts = append(whereParts, mysql.NewFilter("organization_role", "=", req.OrganizationRole.Value))
	}
	if req.EnvironmentId != nil && req.EnvironmentRole != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"environment_id\": \"%s\", \"role\": %d}", req.EnvironmentId.Value, req.EnvironmentRole.Value) // nolint:lll
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
	} else if req.EnvironmentId != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"environment_id\": \"%s\"}", req.EnvironmentId.Value)
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
	} else if req.EnvironmentRole != nil {
		values := make([]interface{}, 1)
		values[0] = fmt.Sprintf("{\"role\": %d}", req.EnvironmentRole.Value)
		whereParts = append(whereParts, mysql.NewJSONFilter("environment_roles", mysql.JSONContainsJSON, values))
	}
	if req.SearchKeyword != "" {
		whereParts = append(whereParts, mysql.NewSearchQuery([]string{"email", "first_name", "last_name"}, req.SearchKeyword))
	}
	orders, err := s.newAccountV2ListOrders(req.OrderBy, req.OrderDirection, localizer)
	if err != nil {
		s.logger.Error(
			"Invalid argument",
			log.FieldsFromImcomingContext(ctx).AddFields(zap.Error(err))...,
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
	accounts, nextCursor, totalCount, err := s.accountStorage.ListAccountsV2(
		ctx,
		whereParts,
		orders,
		limit,
		offset,
	)
	if err != nil {
		s.logger.Error(
			"Failed to list accounts",
			log.FieldsFromImcomingContext(ctx).AddFields(
				zap.Error(err),
				zap.String("organizationID", req.OrganizationId),
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
	return &accountproto.ListAccountsV2Response{
		Accounts:   accounts,
		Cursor:     strconv.Itoa(nextCursor),
		TotalCount: totalCount,
	}, nil
}

func (s *AccountService) newAccountV2ListOrders(
	orderBy accountproto.ListAccountsV2Request_OrderBy,
	orderDirection accountproto.ListAccountsV2Request_OrderDirection,
	localizer locale.Localizer,
) ([]*mysql.Order, error) {
	var column string
	switch orderBy {
	case accountproto.ListAccountsV2Request_DEFAULT,
		accountproto.ListAccountsV2Request_EMAIL:
		column = "email"
	case accountproto.ListAccountsV2Request_CREATED_AT:
		column = "created_at"
	case accountproto.ListAccountsV2Request_UPDATED_AT:
		column = "updated_at"
	case accountproto.ListAccountsV2Request_ORGANIZATION_ROLE:
		column = "organization_role"
	case accountproto.ListAccountsV2Request_ENVIRONMENT_COUNT:
		column = "environment_count"
	case accountproto.ListAccountsV2Request_LAST_SEEN:
		column = "last_seen"
	case accountproto.ListAccountsV2Request_STATE:
		column = "disabled"
	case accountproto.ListAccountsV2Request_TAGS:
		column = "tags"
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
	if orderDirection == accountproto.ListAccountsV2Request_DESC {
		direction = mysql.OrderDirectionDesc
	}
	return []*mysql.Order{mysql.NewOrder(column, direction)}, nil
}
