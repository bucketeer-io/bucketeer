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

	"go.uber.org/zap"

	"github.com/bucketeer-io/bucketeer/pkg/auth"
)

// SESEmailService implements EmailService using Amazon SES
type SESEmailService struct {
	config   auth.EmailServiceConfig
	logger   *zap.Logger
	renderer *TemplateRenderer
}

// NewSESEmailService creates a new SES email service
func NewSESEmailService(config auth.EmailServiceConfig, logger *zap.Logger) (EmailService, error) {
	return &SESEmailService{
		config:   config,
		logger:   logger,
		renderer: NewTemplateRenderer(config),
	}, nil
}

func (s *SESEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
	s.logger.Warn("SES email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SES email service not implemented")
}

func (s *SESEmailService) SendPasswordChangedNotification(ctx context.Context, to string) error {
	s.logger.Warn("SES email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SES email service not implemented")
}

func (s *SESEmailService) SendWelcomeEmail(ctx context.Context, to, tempPassword string) error {
	s.logger.Warn("SES email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SES email service not implemented")
}
