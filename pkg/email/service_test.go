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
	"testing"

	"go.uber.org/zap"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()

	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "disabled email service returns NoOpService",
			config: Config{
				Enabled: false,
			},
			expectError: false,
		},
		{
			name: "smtp provider",
			config: Config{
				Enabled:  true,
				Provider: "smtp",
				SMTP: SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					Username: "user",
					Password: "pass",
				},
			},
			expectError: false,
		},
		{
			name: "sendgrid provider",
			config: Config{
				Enabled:  true,
				Provider: "sendgrid",
				SendGrid: SendGridConfig{
					APIKey: "test-key",
				},
			},
			expectError: false,
		},
		{
			name: "unsupported provider",
			config: Config{
				Enabled:  true,
				Provider: "unsupported",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, err := NewService(tt.config, logger)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if service == nil {
				t.Error("expected service but got nil")
			}
		})
	}
}

func TestNoOpService(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	service := NewNoOpService(logger)
	ctx := context.Background()

	t.Run("SendWelcomeEmail", func(t *testing.T) {
		err := service.SendWelcomeEmail(ctx, "test@example.com", "en")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
