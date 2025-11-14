// Copyright 2026 The Bucketeer Authors.
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"go.uber.org/zap"
)

// SESService implements Service using Amazon SES
type SESService struct {
	config   Config
	logger   *zap.Logger
	renderer *TemplateRenderer
	client   *sesv2.Client
}

// NewSESService creates a new SES email service
func NewSESService(emailConfig Config, logger *zap.Logger) (Service, error) {
	// Create AWS config with explicit credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(emailConfig.SES.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			emailConfig.SES.AccessKey,
			emailConfig.SES.SecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sesv2.NewFromConfig(cfg)

	return &SESService{
		config:   emailConfig,
		logger:   logger,
		renderer: NewTemplateRenderer(emailConfig),
		client:   client,
	}, nil
}

func (s *SESService) SendPasswordChangedNotification(ctx context.Context, to string, language string) error {
	subject, body := s.renderer.RenderPasswordChangedEmail(language, to)

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

func (s *SESService) SendPasswordSetupEmail(
	ctx context.Context, to, setupURL string, ttl time.Duration, language string,
) error {
	subject, body := s.renderer.RenderPasswordSetupEmail(language, setupURL, ttl, to)

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

func (s *SESService) SendPasswordResetEmail(
	ctx context.Context, to, resetURL string, ttl time.Duration, language string,
) error {
	subject, body := s.renderer.RenderPasswordResetEmail(language, resetURL, ttl, to)

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

func (s *SESService) SendWelcomeEmail(ctx context.Context, to string, language string) error {
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

func (s *SESService) sendEmail(ctx context.Context, to, subject, body string) error {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.config.Sender.Email),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(body),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email via SES: %w", err)
	}

	return nil
}
