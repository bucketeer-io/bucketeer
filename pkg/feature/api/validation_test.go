// Copyright 2024 The Bucketeer Authors.
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
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"

	"github.com/bucketeer-io/bucketeer/pkg/locale"
	featureproto "github.com/bucketeer-io/bucketeer/proto/feature"
)

func TestValidateUpdateFeatureTargetingRequest(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctx = metadata.NewIncomingContext(ctx, metadata.MD{
		"accept-language": []string{"ja"},
	})
	localizer := locale.NewLocalizer(ctx)
	patterns := []struct {
		desc     string
		request  *featureproto.UpdateFeatureTargetingRequest
		expected error
	}{
		{
			desc: "error: missing feature id",
			request: &featureproto.UpdateFeatureTargetingRequest{
				Id: "",
			},
			expected: createError(
				t,
				statusMissingID,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "id"),
				localizer,
			),
		},
		{
			desc: "error: missing from",
			request: &featureproto.UpdateFeatureTargetingRequest{
				Id: "feature-id",
			},
			expected: createError(t,
				statusMissingFrom,
				localizer.MustLocalizeWithTemplate(locale.RequiredFieldTemplate, "from"),
				localizer,
			),
		},
		{
			desc: "success: request from user",
			request: &featureproto.UpdateFeatureTargetingRequest{
				Id:   "feature-id",
				From: featureproto.UpdateFeatureTargetingRequest_USER,
			},
			expected: nil,
		},
		{
			desc: "success: request from ops",
			request: &featureproto.UpdateFeatureTargetingRequest{
				Id:   "feature-id",
				From: featureproto.UpdateFeatureTargetingRequest_OPS,
			},
			expected: nil,
		},
	}

	for _, p := range patterns {
		t.Run(p.desc, func(t *testing.T) {
			err := validateUpdateFeatureTargetingRequest(p.request, localizer)
			assert.Equal(t, p.expected, err)
		})
	}
}

func createError(
	t *testing.T,
	status *gstatus.Status,
	msg string,
	localizer locale.Localizer,
) error {
	t.Helper()
	st, err := status.WithDetails(&errdetails.LocalizedMessage{
		Locale:  localizer.GetLocale(),
		Message: msg,
	})
	require.NoError(t, err)
	return st.Err()
}
