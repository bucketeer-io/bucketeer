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

// SendGridEmailService implements EmailService using SendGrid
type SendGridEmailService struct {
	config auth.EmailServiceConfig
	logger *zap.Logger
}

// NewSendGridEmailService creates a new SendGrid email service
func NewSendGridEmailService(config auth.EmailServiceConfig, logger *zap.Logger) EmailService {
	return &SendGridEmailService{
		config: config,
		logger: logger,
	}
}

func (s *SendGridEmailService) SendPasswordResetEmail(ctx context.Context, to, resetToken, resetURL string) error {
	s.logger.Warn("SendGrid email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SendGrid email service not implemented")
}

func (s *SendGridEmailService) SendPasswordChangedNotification(ctx context.Context, to string) error {
	s.logger.Warn("SendGrid email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SendGrid email service not implemented")
}

func (s *SendGridEmailService) SendWelcomeEmail(ctx context.Context, to, tempPassword string) error {
	s.logger.Warn("SendGrid email service not implemented",
		zap.String("to", to),
	)
	return fmt.Errorf("SendGrid email service not implemented")
}
