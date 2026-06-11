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

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	gatewayproto "github.com/bucketeer-io/bucketeer/v2/proto/gateway"
	userproto "github.com/bucketeer-io/bucketeer/v2/proto/user"
)

func TestValidateStreamEvaluationsRequest(t *testing.T) {
	t.Parallel()
	patterns := []struct {
		desc        string
		req         *gatewayproto.StreamEvaluationsRequest
		expectedErr error
	}{
		{
			desc:        "tag missing",
			req:         &gatewayproto.StreamEvaluationsRequest{User: &userproto.User{Id: "u"}},
			expectedErr: errTagRequired,
		},
		{
			desc:        "user missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t"},
			expectedErr: errUserRequired,
		},
		{
			desc:        "user id missing",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{}},
			expectedErr: errUserIDRequired,
		},
		{
			desc:        "valid",
			req:         &gatewayproto.StreamEvaluationsRequest{Tag: "t", User: &userproto.User{Id: "u"}},
			expectedErr: nil,
		},
	}
	gs := &gatewayService{}
	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			assert.Equal(t, p.expectedErr, gs.validateStreamEvaluationsRequest(p.req))
		})
	}
}
