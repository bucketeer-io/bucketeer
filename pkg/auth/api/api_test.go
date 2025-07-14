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
	"testing"

	"github.com/stretchr/testify/assert"

	authproto "github.com/bucketeer-io/bucketeer/proto/auth"
)

func TestAuthService_GetDeploymentStatus(t *testing.T) {
	t.Parallel()
	options := defaultOptions
	service := &authService{
		logger: options.logger,
	}
	patterns := []struct {
		desc        string
		setup       func(s *authService)
		expectedErr error
		expected    *authproto.GetDemoSiteStatusResponse
	}{
		{
			desc: "success: true",
			setup: func(s *authService) {
				s.opts.isDemoSiteEnabled = true
			},
			expectedErr: nil,
			expected: &authproto.GetDemoSiteStatusResponse{
				IsDemoSiteEnabled: true,
			},
		},
		{
			desc: "success: false",
			setup: func(s *authService) {
				s.opts.isDemoSiteEnabled = false
			},
			expectedErr: nil,
			expected: &authproto.GetDemoSiteStatusResponse{
				IsDemoSiteEnabled: false,
			},
		},
	}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			p.setup(service)
			resp, err := service.GetDemoSiteStatus(context.Background(), nil)
			assert.Equal(t, p.expectedErr, err)
			assert.Equal(t, p.expected, resp)
		})
	}
}
