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
	"time"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
	"github.com/bucketeer-io/bucketeer/pkg/auth/storage"
	"github.com/bucketeer-io/bucketeer/pkg/locale"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func (s *authService) UpdatePassword(
	ctx context.Context,
	request *authproto.UpdatePasswordRequest,
) (*authproto.UpdatePasswordResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateUpdatePasswordRequest(request, localizer)
	if err != nil {
		s.logger.Error("UpdatePassword request validation failed", zap.Error(err))
		return nil, err
	}

	// Check if password authentication is enabled
	if !s.config.PasswordAuth.Enabled {
		s.logger.Error("Password authentication not enabled")
		dt, err := auth.StatusInvalidEmailConfig.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Get email from context (user must be authenticated)
	email := extractEmailFromContext(ctx)
	if email == "" {
		s.logger.Error("No email in context")
		return nil, auth.StatusUnauthenticated.Err()
	}

	// Get current credentials
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("No password found for user", zap.String("email", email))
			dt, err := auth.StatusPasswordNotFound.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.NotFoundError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error("Failed to get credentials", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Verify current password
	if !auth.ValidatePassword(request.CurrentPassword, credentials.PasswordHash) {
		s.logger.Error("Current password mismatch", zap.String("email", email))
		dt, err := auth.StatusPasswordMismatch.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Validate new password complexity
	err = auth.ValidatePasswordComplexity(request.NewPassword, s.config.PasswordAuth)
	if err != nil {
		s.logger.Error("New password complexity validation failed", zap.Error(err))
		dt, err := auth.StatusPasswordTooWeak.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: err.Error(),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Hash new password
	newPasswordHash, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash new password", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Update password
	err = s.credentialsStorage.UpdatePassword(ctx, email, newPasswordHash)
	if err != nil {
		s.logger.Error("Failed to update password", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Send notification email if email service is enabled
	if s.config.PasswordAuth.EmailServiceEnabled && s.emailService != nil {
		err = s.emailService.SendPasswordChangedNotification(ctx, email)
		if err != nil {
			s.logger.Warn("Failed to send password changed notification",
				zap.Error(err),
				zap.String("email", email),
			)
			// Don't fail the password update if email sending fails
		}
	}

	s.logger.Info("Password updated successfully", zap.String("email", email))
	return &authproto.UpdatePasswordResponse{}, nil
}

func (s *authService) InitiatePasswordSetup(
	ctx context.Context,
	request *authproto.InitiatePasswordSetupRequest,
) (*authproto.InitiatePasswordSetupResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateInitiatePasswordSetupRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Check if password authentication and email service are enabled
	if !s.config.PasswordAuth.Enabled || !s.config.PasswordAuth.EmailServiceEnabled {
		s.logger.Error("Password setup not available")
		dt, err := auth.StatusEmailServiceUnavailable.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	email := request.Email

	// Validate that the user has organizations (i.e., account exists)
	organizations, err := s.getOrganizationsByEmail(ctx, email, localizer)
	if err != nil || len(organizations) == 0 {
		s.logger.Warn("Password setup attempted for non-existent account", zap.String("email", email))
		return &authproto.InitiatePasswordSetupResponse{
			Message: localizer.MustLocalize(locale.PasswordSetupEmailSent),
		}, nil
	}

	// Check if credentials already exist (user already has a password)
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err == nil && credentials.PasswordHash != "" {
		// Password already exists, don't reveal this for security
		s.logger.Warn("Password setup attempted for account with existing password", zap.String("email", email))
		return &authproto.InitiatePasswordSetupResponse{
			Message: localizer.MustLocalize(locale.PasswordSetupEmailSent),
		}, nil
	}
	if !errors.Is(err, storage.ErrCredentialsNotFound) {
		s.logger.Error("Failed to check credentials for password setup", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Create empty credentials record for password setup (no password hash yet)
	err = s.credentialsStorage.CreateCredentials(ctx, email, "")
	if err != nil {
		s.logger.Error("Failed to create credentials record for password setup", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Generate secure setup token
	setupToken, err := auth.GenerateSecureToken()
	if err != nil {
		s.logger.Error("Failed to generate setup token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Store setup token with longer expiration (use PasswordSetupTokenTTL)
	expiresAt := time.Now().Add(s.config.PasswordAuth.PasswordSetupTokenTTL).Unix()
	err = s.credentialsStorage.SetPasswordResetToken(ctx, email, setupToken, expiresAt)
	if err != nil {
		s.logger.Error("Failed to store setup token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Send setup email
	if s.emailService != nil {
		setupPath := s.config.PasswordAuth.EmailServiceConfig.PasswordSetupPath
		if setupPath == "" {
			s.logger.Error("Password setup path not configured")
			return nil, auth.StatusInternal.Err()
		}
		setupParam := s.config.PasswordAuth.EmailServiceConfig.PasswordSetupParam
		if setupParam == "" {
			s.logger.Error("Password setup parameter not configured")
			return nil, auth.StatusInternal.Err()
		}
		setupURL := fmt.Sprintf("%s%s?%s=%s",
			s.config.PasswordAuth.EmailServiceConfig.BaseURL, setupPath, setupParam, setupToken)
		err = s.emailService.SendPasswordSetupEmail(ctx, email, setupURL, s.config.PasswordAuth.PasswordSetupTokenTTL)
		if err != nil {
			s.logger.Error("Failed to send password setup email",
				zap.Error(err),
				zap.String("email", email),
			)
			// Don't return error to user for security reasons
		}
	}

	s.logger.Info("Password setup initiated", zap.String("email", email))
	return &authproto.InitiatePasswordSetupResponse{
		Message: localizer.MustLocalize(locale.PasswordSetupEmailSent),
	}, nil
}

func (s *authService) SetupPassword(
	ctx context.Context,
	request *authproto.SetupPasswordRequest,
) (*authproto.SetupPasswordResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateSetupPasswordRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Check if password authentication is enabled
	if !s.config.PasswordAuth.Enabled {
		s.logger.Error("Password authentication not enabled")
		dt, err := auth.StatusInvalidEmailConfig.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Get and validate setup token (reusing password reset token infrastructure)
	setupToken, err := s.credentialsStorage.GetPasswordResetToken(ctx, request.SetupToken)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordResetTokenNotFound) {
			s.logger.Error("Invalid setup token", zap.String("token", request.SetupToken))
			dt, err := auth.StatusInvalidResetToken.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error("Failed to get setup token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if token is expired
	if setupToken.IsExpired() {
		s.logger.Error("Expired setup token", zap.String("email", setupToken.Email))
		dt, err := auth.StatusExpiredResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Validate new password complexity
	err = auth.ValidatePasswordComplexity(request.NewPassword, s.config.PasswordAuth)
	if err != nil {
		s.logger.Error("Password complexity validation failed", zap.Error(err))
		dt, err := auth.StatusPasswordTooWeak.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: err.Error(),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Check if credentials already exist (prevent double setup)
	credentials, err := s.credentialsStorage.GetCredentials(ctx, setupToken.Email)
	if err == nil && credentials.PasswordHash != "" {
		s.logger.Error("Setup attempted for account with existing password", zap.String("email", setupToken.Email))
		dt, err := auth.StatusPasswordAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.AlreadyExistsError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}
	if !errors.Is(err, storage.ErrCredentialsNotFound) {
		s.logger.Error("Failed to check credentials during setup", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Hash new password
	passwordHash, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash password during setup", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Create credentials
	err = s.credentialsStorage.CreateCredentials(ctx, setupToken.Email, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsAlreadyExists) {
			s.logger.Error("Password setup attempted but credentials already exist", zap.String("email", setupToken.Email))
			dt, err := auth.StatusPasswordAlreadyExists.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.AlreadyExistsError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error("Failed to create credentials during setup", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Delete setup token
	err = s.credentialsStorage.DeletePasswordResetToken(ctx, request.SetupToken)
	if err != nil {
		s.logger.Warn("Failed to delete setup token", zap.Error(err))
		// Don't fail the setup if token deletion fails
	}

	// Send welcome email if email service is enabled
	if s.config.PasswordAuth.EmailServiceEnabled && s.emailService != nil {
		err = s.emailService.SendPasswordChangedNotification(ctx, setupToken.Email)
		if err != nil {
			s.logger.Warn("Failed to send password setup completion notification",
				zap.Error(err),
				zap.String("email", setupToken.Email),
			)
			// Don't fail the setup if email sending fails
		}
	}

	s.logger.Info("Password setup completed successfully", zap.String("email", setupToken.Email))
	return &authproto.SetupPasswordResponse{}, nil
}

func (s *authService) ValidatePasswordSetupToken(
	ctx context.Context,
	request *authproto.ValidatePasswordSetupTokenRequest,
) (*authproto.ValidatePasswordSetupTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validatePasswordSetupTokenRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Get setup token (reusing password reset token infrastructure)
	setupToken, err := s.credentialsStorage.GetPasswordResetToken(ctx, request.SetupToken)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordResetTokenNotFound) {
			return &authproto.ValidatePasswordSetupTokenResponse{
				IsValid: false,
				Email:   "",
			}, nil
		}
		s.logger.Error("Failed to get setup token for validation", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if token is valid (not expired)
	isValid := setupToken.IsValid()
	email := ""
	if isValid {
		// Additional validation: check if account still needs password setup
		credentials, err := s.credentialsStorage.GetCredentials(ctx, setupToken.Email)
		if err == nil && credentials.PasswordHash != "" {
			// Credentials already exist, token is no longer valid for setup
			isValid = false
		} else if !errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("Failed to check credentials during token validation", zap.Error(err))
			return nil, auth.StatusInternal.Err()
		} else {
			// Credentials don't exist, token is valid for setup
			email = setupToken.Email
		}
	}

	return &authproto.ValidatePasswordSetupTokenResponse{
		IsValid: isValid,
		Email:   email,
	}, nil
}

// extractEmailFromContext extracts email from the authentication context
func (s *authService) InitiatePasswordReset(
	ctx context.Context,
	request *authproto.InitiatePasswordResetRequest,
) (*authproto.InitiatePasswordResetResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateInitiatePasswordResetRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Check if password authentication and email service are enabled
	if !s.config.PasswordAuth.Enabled || !s.config.PasswordAuth.EmailServiceEnabled {
		s.logger.Error("Password reset not available")
		dt, err := auth.StatusEmailServiceUnavailable.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	email := request.Email

	// Check if credentials exist with a password (only allow reset for existing password users)
	credentials, err := s.credentialsStorage.GetCredentials(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			// Don't reveal whether the account exists for security
			s.logger.Warn("Password reset attempted for non-existent account", zap.String("email", email))
			return &authproto.InitiatePasswordResetResponse{
				Message: localizer.MustLocalize(locale.PasswordResetEmailSent),
			}, nil
		}
		s.logger.Error("Failed to check credentials for password reset", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if user actually has a password set
	if credentials.PasswordHash == "" {
		// Don't reveal this information for security
		s.logger.Warn("Password reset attempted for account without password", zap.String("email", email))
		return &authproto.InitiatePasswordResetResponse{
			Message: localizer.MustLocalize(locale.PasswordResetEmailSent),
		}, nil
	}

	// Generate secure reset token
	resetToken, err := auth.GenerateSecureToken()
	if err != nil {
		s.logger.Error("Failed to generate reset token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Store reset token with expiration (use PasswordResetTokenTTL)
	expiresAt := time.Now().Add(s.config.PasswordAuth.PasswordResetTokenTTL).Unix()
	err = s.credentialsStorage.SetPasswordResetToken(ctx, email, resetToken, expiresAt)
	if err != nil {
		s.logger.Error("Failed to store reset token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Send reset email
	if s.emailService != nil {
		resetPath := s.config.PasswordAuth.EmailServiceConfig.PasswordResetPath
		if resetPath == "" {
			s.logger.Error("Password reset path not configured")
			return nil, auth.StatusInternal.Err()
		}
		resetParam := s.config.PasswordAuth.EmailServiceConfig.PasswordResetParam
		if resetParam == "" {
			s.logger.Error("Password reset parameter not configured")
			return nil, auth.StatusInternal.Err()
		}
		resetURL := fmt.Sprintf("%s%s?%s=%s",
			s.config.PasswordAuth.EmailServiceConfig.BaseURL, resetPath, resetParam, resetToken)
		err = s.emailService.SendPasswordResetEmail(ctx, email, resetURL, s.config.PasswordAuth.PasswordResetTokenTTL)
		if err != nil {
			s.logger.Error("Failed to send password reset email",
				zap.Error(err),
				zap.String("email", email),
			)
			// Don't return error to user for security reasons
		}
	}

	s.logger.Info("Password reset initiated", zap.String("email", email))
	return &authproto.InitiatePasswordResetResponse{
		Message: localizer.MustLocalize(locale.PasswordResetEmailSent),
	}, nil
}

func (s *authService) ResetPassword(
	ctx context.Context,
	request *authproto.ResetPasswordRequest,
) (*authproto.ResetPasswordResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validateResetPasswordRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Check if password authentication is enabled
	if !s.config.PasswordAuth.Enabled {
		s.logger.Error("Password authentication not enabled")
		dt, err := auth.StatusInvalidEmailConfig.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.PermissionDenied),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Get and validate reset token (reusing password reset token infrastructure)
	resetToken, err := s.credentialsStorage.GetPasswordResetToken(ctx, request.ResetToken)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordResetTokenNotFound) {
			s.logger.Error("Invalid reset token", zap.String("token", request.ResetToken))
			dt, err := auth.StatusInvalidResetToken.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error("Failed to get reset token", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if token is expired
	if resetToken.IsExpired() {
		s.logger.Error("Expired reset token", zap.String("email", resetToken.Email))
		dt, err := auth.StatusExpiredResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Validate new password complexity
	err = auth.ValidatePasswordComplexity(request.NewPassword, s.config.PasswordAuth)
	if err != nil {
		s.logger.Error("Password complexity validation failed", zap.Error(err))
		dt, err := auth.StatusPasswordTooWeak.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: err.Error(),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Check if credentials exist (should exist since reset token was valid)
	credentials, err := s.credentialsStorage.GetCredentials(ctx, resetToken.Email)
	if err != nil {
		if errors.Is(err, storage.ErrCredentialsNotFound) {
			s.logger.Error("Reset attempted for account without credentials", zap.String("email", resetToken.Email))
			dt, err := auth.StatusInvalidResetToken.WithDetails(&errdetails.LocalizedMessage{
				Locale:  localizer.GetLocale(),
				Message: localizer.MustLocalize(locale.InvalidArgumentError),
			})
			if err != nil {
				return nil, err
			}
			return nil, dt.Err()
		}
		s.logger.Error("Failed to check credentials during reset", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if user has a password (should have one since reset was initiated)
	if credentials.PasswordHash == "" {
		s.logger.Error("Reset attempted for account without password", zap.String("email", resetToken.Email))
		dt, err := auth.StatusInvalidResetToken.WithDetails(&errdetails.LocalizedMessage{
			Locale:  localizer.GetLocale(),
			Message: localizer.MustLocalize(locale.InvalidArgumentError),
		})
		if err != nil {
			return nil, err
		}
		return nil, dt.Err()
	}

	// Hash new password
	passwordHash, err := auth.HashPassword(request.NewPassword)
	if err != nil {
		s.logger.Error("Failed to hash password during reset", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Update password
	err = s.credentialsStorage.UpdatePassword(ctx, resetToken.Email, passwordHash)
	if err != nil {
		s.logger.Error("Failed to update password during reset", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Delete reset token
	err = s.credentialsStorage.DeletePasswordResetToken(ctx, request.ResetToken)
	if err != nil {
		s.logger.Warn("Failed to delete reset token", zap.Error(err))
		// Don't fail the reset if token deletion fails
	}

	// Send password changed notification email if email service is enabled
	if s.config.PasswordAuth.EmailServiceEnabled && s.emailService != nil {
		err = s.emailService.SendPasswordChangedNotification(ctx, resetToken.Email)
		if err != nil {
			s.logger.Warn("Failed to send password changed notification",
				zap.Error(err),
				zap.String("email", resetToken.Email),
			)
			// Don't fail the reset if email sending fails
		}
	}

	s.logger.Info("Password reset completed successfully", zap.String("email", resetToken.Email))
	return &authproto.ResetPasswordResponse{}, nil
}

func (s *authService) ValidatePasswordResetToken(
	ctx context.Context,
	request *authproto.ValidatePasswordResetTokenRequest,
) (*authproto.ValidatePasswordResetTokenResponse, error) {
	localizer := locale.NewLocalizer(ctx)
	err := validatePasswordResetTokenRequest(request, localizer)
	if err != nil {
		return nil, err
	}

	// Get reset token (reusing password reset token infrastructure)
	resetToken, err := s.credentialsStorage.GetPasswordResetToken(ctx, request.ResetToken)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordResetTokenNotFound) {
			return &authproto.ValidatePasswordResetTokenResponse{
				IsValid: false,
				Email:   "",
			}, nil
		}
		s.logger.Error("Failed to get reset token for validation", zap.Error(err))
		return nil, auth.StatusInternal.Err()
	}

	// Check if token is valid (not expired)
	isValid := resetToken.IsValid()
	email := ""
	if isValid {
		// Additional validation: check if account still has a password to reset
		credentials, err := s.credentialsStorage.GetCredentials(ctx, resetToken.Email)
		if err != nil {
			if errors.Is(err, storage.ErrCredentialsNotFound) {
				// No credentials exist, token is no longer valid for reset
				isValid = false
			} else {
				s.logger.Error("Failed to check credentials during token validation", zap.Error(err))
				return nil, auth.StatusInternal.Err()
			}
		} else if credentials.PasswordHash == "" {
			// No password exists, token is no longer valid for reset
			isValid = false
		} else {
			// Credentials exist with password, token is valid for reset
			email = resetToken.Email
		}
	}

	return &authproto.ValidatePasswordResetTokenResponse{
		IsValid: isValid,
		Email:   email,
	}, nil
}

func extractEmailFromContext(ctx context.Context) string {
	accessToken, ok := rpc.GetAccessToken(ctx)
	if !ok || accessToken == nil {
		return ""
	}
	return accessToken.Email
}
