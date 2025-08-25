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

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	config   auth.EmailServiceConfig
	logger   *zap.Logger
	renderer *TemplateRenderer
}

// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService(config auth.EmailServiceConfig, logger *zap.Logger) EmailService {
	return &SendGridEmailService{
		config:   config,
		logger:   logger,
		renderer: NewTemplateRenderer(config),
	}
}

func (s *SendGridEmailService) SendPasswordChangedNotification(ctx context.Context, to string) error {
	subject, body := s.renderer.RenderPasswordChangedEmail()

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send password changed notification",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send password changed notification: %w", err)
	}

	s.logger.Info("Password changed notification sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SendGridEmailService) SendPasswordSetupEmail(
	ctx context.Context,
	to, setupURL string,
	ttl time.Duration,
) error {
	subject, body := s.renderer.RenderPasswordSetupEmail(setupURL, ttl)

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send password setup email",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send password setup email: %w", err)
	}

	s.logger.Info("Password setup email sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SendGridEmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	from := mail.NewEmail(s.config.FromName, s.config.FromEmail)
	toEmail := mail.NewEmail("", to)
	message := mail.NewSingleEmail(from, subject, toEmail, "", body)

	client := sendgrid.NewSendClient(s.config.SendGridAPIKey)
	response, err := client.SendWithContext(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid API error: status code %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
