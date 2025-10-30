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

	"github.com/mailersend/mailersend-go"
	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/v2/pkg/auth"
)

// MailerSendEmailService implements EmailService using MailerSend
type MailerSendEmailService struct {
	config   auth.EmailServiceConfig
	logger   *zap.Logger
	renderer *TemplateRenderer
	client   *mailersend.Mailersend
}

// NewMailerSendEmailService creates a new MailerSend email service
func NewMailerSendEmailService(config auth.EmailServiceConfig, logger *zap.Logger) EmailService {
	client := mailersend.NewMailersend(config.MailerSend.MailerSendAPIKey)

	return &MailerSendEmailService{
		config:   config,
		logger:   logger,
		renderer: NewTemplateRenderer(config),
		client:   client,
	}
}

func (s *MailerSendEmailService) SendPasswordChangedNotification(ctx context.Context, to string, language string) error {
	subject, body := s.renderer.RenderPasswordChangedEmail(language)

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

func (s *MailerSendEmailService) SendPasswordSetupEmail(
	ctx context.Context,
	to, setupURL string,
	ttl time.Duration,
	language string,
) error {
	subject, body := s.renderer.RenderPasswordSetupEmail(language, setupURL, ttl)

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

func (s *MailerSendEmailService) SendPasswordResetEmail(
	ctx context.Context,
	to, resetURL string,
	ttl time.Duration,
	language string,
) error {
	subject, body := s.renderer.RenderPasswordResetEmail(language, resetURL, ttl)

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send password reset email",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send password reset email: %w", err)
	}

	s.logger.Info("Password reset email sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *MailerSendEmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	// Create the message using MailerSend's message builder
	message := s.client.Email.NewMessage()

	// Set sender
	from := mailersend.From{
		Name:  s.config.FromName,
		Email: s.config.FromEmail,
	}
	message.SetFrom(from)

	// Set recipient
	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}
	message.SetRecipients(recipients)

	// Set subject and HTML body
	message.SetSubject(subject)
	message.SetHTML(body)

	// Send the email
	_, err := s.client.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email via MailerSend: %w", err)
	}

	return nil
}
