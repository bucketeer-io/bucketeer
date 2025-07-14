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

package api

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func TestAuthService_GetDeploymentStatus(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		setup       func()
		expectedErr error
		expected    *authproto.GetDeploymentStatusResponse
	}{
		{
			desc: "success: true",
			setup: func() {
				err := os.Setenv("DEMO_SITE_ENABLED", "true")
				if err != nil {
					t.Fatalf("failed to set environment variable: %v", err)
				}
			},
			expectedErr: nil,
			expected: &authproto.GetDeploymentStatusResponse{
				IsDemoSiteEnabled: true,
			},
		},
		{
			desc: "success: false",
			setup: func() {
				err := os.Setenv("DEMO_SITE_ENABLED", "false")
				if err != nil {
					t.Fatalf("failed to set environment variable: %v", err)
				}
			},
			expectedErr: nil,
			expected: &authproto.GetDeploymentStatusResponse{
				IsDemoSiteEnabled: false,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p.setup()
			options := defaultOptions
			s := &authService{
				logger: options.logger,
			}
			resp, err := s.GetDeploymentStatus(context.Background(), nil)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, resp)
		})
	}
}
