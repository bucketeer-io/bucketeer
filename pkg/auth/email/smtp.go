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
	"net/smtp"
	"time"

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// SMTPEmailService implements EmailService using SMTP
type SMTPEmailService struct {
	config   auth.EmailServiceConfig
	logger   *zap.Logger
	renderer *TemplateRenderer
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(config auth.EmailServiceConfig, logger *zap.Logger) EmailService {
	return &SMTPEmailService{
		config:   config,
		logger:   logger,
		renderer: NewTemplateRenderer(config),
	}
}

func (s *SMTPEmailService) SendPasswordChangedNotification(ctx context.Context, to string, language string) error {
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

func (s *SMTPEmailService) SendPasswordSetupEmail(
	ctx context.Context, to, setupURL string, ttl time.Duration, language string,
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

func (s *SMTPEmailService) SendPasswordResetEmail(
	ctx context.Context, to, resetURL string, ttl time.Duration, language string,
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

func (s *SMTPEmailService) sendEmail(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.SMTPUsername, s.config.SMTP.SMTPPassword, s.config.SMTP.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, s.config.FromEmail, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.SMTPHost, s.config.SMTP.SMTPPort)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}
