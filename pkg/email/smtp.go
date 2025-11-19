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

	"go.uber.org/zap"
)

// SMTPService implements Service using SMTP
type SMTPService struct {
	config   Config
	logger   *zap.Logger
	renderer *TemplateRenderer
}

// NewSMTPService creates a new SMTP email service
func NewSMTPService(config Config, logger *zap.Logger) Service {
	return &SMTPService{
		config:   config,
		logger:   logger,
		renderer: NewTemplateRenderer(config),
	}
}

func (s *SMTPService) SendWelcomeEmail(ctx context.Context, to string, language string) error {
	subject, body := s.renderer.RenderWelcomeEmail(language, to)

	err := s.sendEmail(ctx, to, subject, body)
	if err != nil {
		s.logger.Error("Failed to send welcome email",
			zap.Error(err),
			zap.String("to", to),
		)
		return fmt.Errorf("failed to send welcome email: %w", err)
	}

	s.logger.Info("Welcome email sent successfully",
		zap.String("to", to),
	)
	return nil
}

func (s *SMTPService) sendEmail(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", to, s.config.Sender.Email, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	return smtp.SendMail(addr, auth, s.config.Sender.Email, []string{to}, msg)
}
