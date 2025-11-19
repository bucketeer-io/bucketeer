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

package email

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Service defines the interface for sending emails
type Service interface {
	SendPasswordChangedNotification(ctx context.Context, to string, language string) error
	SendPasswordSetupEmail(ctx context.Context, to, setupURL string, ttl time.Duration, language string) error
	SendPasswordResetEmail(ctx context.Context, to, resetURL string, ttl time.Duration, language string) error
	SendWelcomeEmail(ctx context.Context, to string, language string) error
}

// NewService creates an email service based on configuration
func NewService(config Config, logger *zap.Logger) (Service, error) {
	if !config.Enabled {
		return NewNoOpService(logger), nil
	}

	switch config.Provider {
	case "smtp":
		return NewSMTPService(config, logger), nil
	case "sendgrid":
		return NewSendGridService(config, logger), nil
	case "ses":
		return NewSESService(config, logger)
	case "mailersend":
		return NewMailerSendService(config, logger), nil
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
	}
}

// NoOpService is a no-operation email service for testing or when email is disabled
type NoOpService struct {
	logger *zap.Logger
}

// NewNoOpService creates a no-operation email service
func NewNoOpService(logger *zap.Logger) Service {
	return &NoOpService{logger: logger}
}

func (s *NoOpService) SendPasswordChangedNotification(ctx context.Context, to string, language string) error {
	s.logger.Info("No-op email service: password changed notification not sent",
		zap.String("to", to),
		zap.String("language", language),
	)
	return nil
}

func (s *NoOpService) SendPasswordSetupEmail(
	ctx context.Context, to, setupURL string, ttl time.Duration, language string,
) error {
	s.logger.Info("No-op email service: password setup email not sent",
		zap.String("to", to),
		zap.String("setupURL", setupURL),
		zap.String("language", language),
	)
	return nil
}

func (s *NoOpService) SendPasswordResetEmail(
	ctx context.Context, to, resetURL string, ttl time.Duration, language string,
) error {
	s.logger.Info("No-op email service: password reset email not sent",
		zap.String("to", to),
		zap.String("resetURL", resetURL),
		zap.String("language", language),
	)
	return nil
}

func (s *NoOpService) SendWelcomeEmail(ctx context.Context, to string, language string) error {
	s.logger.Info("No-op email service: welcome email not sent",
		zap.String("to", to),
		zap.String("language", language),
	)
	return nil
}
