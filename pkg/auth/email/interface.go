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

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// EmailService defines the interface for sending emails
type EmailService interface {
	SendPasswordChangedNotification(ctx context.Context, to string) error
	SendPasswordSetupEmail(ctx context.Context, to, setupURL string, ttl time.Duration) error
}

// NewEmailService EmailServiceFactory creates an email service based on configuration
func NewEmailService(config auth.EmailServiceConfig, logger *zap.Logger) (EmailService, error) {
	switch config.Provider {
	case "smtp":
		return NewSMTPEmailService(config, logger), nil
	case "sendgrid":
		return NewSendGridEmailService(config, logger), nil
	case "ses":
		return NewSESEmailService(config, logger)
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", config.Provider)
	}
}

// NoOpEmailService is a no-operation email service for testing or when email is disabled
type NoOpEmailService struct {
	logger *zap.Logger
}

// NewNoOpEmailService creates a no-operation email service
func NewNoOpEmailService(logger *zap.Logger) EmailService {
	return &NoOpEmailService{logger: logger}
}

func (s *NoOpEmailService) SendPasswordChangedNotification(ctx context.Context, to string) error {
	s.logger.Info("No-op email service: password changed notification not sent",
		zap.String("to", to),
	)
	return nil
}

func (s *NoOpEmailService) SendPasswordSetupEmail(ctx context.Context, to, setupURL string, ttl time.Duration) error {
	s.logger.Info("No-op email service: password setup email not sent",
		zap.String("to", to),
		zap.String("setupURL", setupURL),
	)
	return nil
}
